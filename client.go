// Copyright 2018-2019 opcua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package opcua

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gopcua/opcua/debug"
	"github.com/gopcua/opcua/id"
	"github.com/gopcua/opcua/ua"
	"github.com/gopcua/opcua/uacp"
	"github.com/gopcua/opcua/uasc"
)

// GetEndpoints returns the available endpoint descriptions for the server.
func GetEndpoints(endpoint string) ([]*ua.EndpointDescription, error) {
	c := NewClient(endpoint)
	if err := c.Dial(); err != nil {
		return nil, err
	}
	defer c.Close()
	res, err := c.GetEndpoints()
	if err != nil {
		return nil, err
	}
	return res.Endpoints, nil
}

// Client is a high-level client for an OPC/UA server.
// It establishes a secure channel and a session.
type Client struct {
	// endpointURL is the endpoint URL the client connects to.
	endpointURL string

	// cfg is the configuration for the secure channel.
	cfg *uasc.Config

	// sessionCfg is the configuration for the session.
	sessionCfg *uasc.SessionConfig

	// sechan is the open secure channel.
	sechan *uasc.SecureChannel

	// session is the active session.
	session atomic.Value // *Session

	// once initializes session
	once sync.Once

	// map of active subscriptions managed by this client
	subscriptions map[uint32]Subscription
}

// NewClient creates a new Client.
//
// When no options are provided the new client is created from
// DefaultClientConfig() and DefaultSessionConfig(). If no authentication method
// is configured, a UserIdentityToken for anonymous authentication will be set.
// See #Client.CreateSession for details.
//
// To modify configuration you can provide any number of Options as opts. See
// #Option for details.
//
// https://godoc.org/github.com/gopcua/opcua#Option
func NewClient(endpoint string, opts ...Option) *Client {
	c := &Client{
		endpointURL:   endpoint,
		cfg:           DefaultClientConfig(),
		sessionCfg:    DefaultSessionConfig(),
		subscriptions: make(map[uint32]Subscription),
	}
	for _, opt := range opts {
		opt(c.cfg, c.sessionCfg)
	}

	return c
}

// Connect establishes a secure channel and creates a new session.
func (c *Client) Connect() (err error) {
	if c.sechan != nil {
		return fmt.Errorf("already connected")
	}
	if err := c.Dial(); err != nil {
		return err
	}
	s, err := c.CreateSession(c.sessionCfg)
	if err != nil {
		_ = c.Close()
		return err
	}
	if err := c.ActivateSession(s); err != nil {
		_ = c.Close()
		return err
	}
	return nil
}

// Dial establishes a secure channel.
func (c *Client) Dial() error {
	c.once.Do(func() { c.session.Store((*Session)(nil)) })
	if c.sechan != nil {
		return fmt.Errorf("secure channel already connected")
	}
	conn, err := uacp.Dial(context.Background(), c.endpointURL)
	if err != nil {
		return err
	}
	sechan, err := uasc.NewSecureChannel(c.endpointURL, conn, c.cfg)
	if err != nil {
		_ = conn.Close()
		return err
	}
	if err := sechan.Open(); err != nil {
		_ = conn.Close()
		return err
	}
	c.sechan = sechan
	return nil
}

// Close closes the session and the secure channel.
func (c *Client) Close() error {
	// try to close the session but ignore any error
	// so that we close the underlying channel and connection.
	_ = c.CloseSession()
	return c.sechan.Close()
}

// Session returns the active session.
func (c *Client) Session() *Session {
	return c.session.Load().(*Session)
}

// Session is a OPC/UA session as described in Part 4, 5.6.
type Session struct {
	cfg *uasc.SessionConfig

	// resp is the response to the CreateSession request which contains all
	// necessary parameters to activate the session.
	resp *ua.CreateSessionResponse

	// serverCertificate is the certificate used to generate the signatures for
	// the ActivateSessionRequest methods
	serverCertificate []byte

	// serverNonce is the secret nonce received from the server during Create and Activate
	// Session respones. Used to generate the signatures for the ActivateSessionRequest
	// and User Authorization
	serverNonce []byte
}

// CreateSession creates a new session which is not yet activated and not
// associated with the client. Call ActivateSession to both activate and
// associate the session with the client.
//
// If no UserIdentityToken is given explicitly before calling CreateSesion,
// it automatically sets anonymous identity token with the same PolicyID
// that the server sent in Create Session Response. The default PolicyID
// "Anonymous" wii be set if it's missing in response.
//
// See Part 4, 5.6.2
func (c *Client) CreateSession(cfg *uasc.SessionConfig) (*Session, error) {
	if c.sechan == nil {
		return nil, fmt.Errorf("secure channel not connected")
	}

	nonce := make([]byte, 32)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	req := &ua.CreateSessionRequest{
		ClientDescription:       cfg.ClientDescription,
		EndpointURL:             c.endpointURL,
		SessionName:             fmt.Sprintf("gopcua-%d", time.Now().UnixNano()),
		ClientNonce:             nonce,
		ClientCertificate:       c.cfg.Certificate,
		RequestedSessionTimeout: float64(cfg.SessionTimeout / time.Millisecond),
	}

	var s *Session
	// for the CreateSessionRequest the authToken is always nil.
	// use c.sechan.Send() to enforce this.
	err := c.sechan.Send(req, nil, func(v interface{}) error {
		var res *ua.CreateSessionResponse
		if err := safeAssign(v, &res); err != nil {
			return err
		}

		err := c.sechan.VerifySessionSignature(res.ServerCertificate, nonce, res.ServerSignature.Signature)
		if err != nil {
			log.Printf("error verifying session signature: %s", err)
			return nil
		}

		// Ensure we have a valid identity token that the server will accept before trying to activate a session
		if c.sessionCfg.UserIdentityToken == nil {
			opt := AuthAnonymous()
			opt(c.cfg, c.sessionCfg)

			p := anonymousPolicyID(res.ServerEndpoints)
			opt = AuthPolicyID(p)
			opt(c.cfg, c.sessionCfg)
		}

		s = &Session{
			cfg:               cfg,
			resp:              res,
			serverNonce:       res.ServerNonce,
			serverCertificate: res.ServerCertificate,
		}

		return nil
	})
	return s, err
}

const defaultAnonymousPolicyID = "Anonymous"

func anonymousPolicyID(endpoints []*ua.EndpointDescription) string {
	for _, e := range endpoints {
		if e.SecurityMode != ua.MessageSecurityModeNone || e.SecurityPolicyURI != ua.SecurityPolicyURINone {
			continue
		}

		for _, t := range e.UserIdentityTokens {
			if t.TokenType == ua.UserTokenTypeAnonymous {
				return t.PolicyID
			}
		}
	}

	return defaultAnonymousPolicyID
}

// ActivateSession activates the session and associates it with the client. If
// the client already has a session it will be closed. To retain the current
// session call DetachSession.
//
// See Part 4, 5.6.3
func (c *Client) ActivateSession(s *Session) error {
	sig, sigAlg, err := c.sechan.NewSessionSignature(s.serverCertificate, s.serverNonce)
	if err != nil {
		log.Printf("error creating session signature: %s", err)
		return nil
	}

	switch tok := s.cfg.UserIdentityToken.(type) {
	case *ua.AnonymousIdentityToken:
		// nothing to do

	case *ua.UserNameIdentityToken:
		pass, passAlg, err := c.sechan.EncryptUserPassword(s.cfg.AuthPolicyURI, s.cfg.AuthPassword, s.serverCertificate, s.serverNonce)
		if err != nil {
			log.Printf("error encrypting user password: %s", err)
			return err
		}
		tok.Password = pass
		tok.EncryptionAlgorithm = passAlg

	case *ua.X509IdentityToken:
		tokSig, tokSigAlg, err := c.sechan.NewUserTokenSignature(s.cfg.AuthPolicyURI, s.serverCertificate, s.serverNonce)
		if err != nil {
			log.Printf("error creating session signature: %s", err)
			return err
		}
		s.cfg.UserTokenSignature = &ua.SignatureData{
			Algorithm: tokSigAlg,
			Signature: tokSig,
		}

	case *ua.IssuedIdentityToken:
		tok.EncryptionAlgorithm = ""
	}

	req := &ua.ActivateSessionRequest{
		ClientSignature: &ua.SignatureData{
			Algorithm: sigAlg,
			Signature: sig,
		},
		ClientSoftwareCertificates: nil,
		LocaleIDs:                  s.cfg.LocaleIDs,
		UserIdentityToken:          ua.NewExtensionObject(s.cfg.UserIdentityToken),
		UserTokenSignature:         s.cfg.UserTokenSignature,
	}
	return c.sechan.Send(req, s.resp.AuthenticationToken, func(v interface{}) error {
		var res *ua.ActivateSessionResponse
		if err := safeAssign(v, &res); err != nil {
			return err
		}

		// save the nonce for the next request
		s.serverNonce = res.ServerNonce

		if err := c.CloseSession(); err != nil {
			// try to close the newly created session but report
			// only the initial error.
			_ = c.closeSession(s)
			return err
		}
		c.session.Store(s)
		return nil
	})
}

// CloseSession closes the current session.
//
// See Part 4, 5.6.4
func (c *Client) CloseSession() error {
	if err := c.closeSession(c.Session()); err != nil {
		return err
	}
	c.session.Store((*Session)(nil))
	return nil
}

// closeSession closes the given session.
func (c *Client) closeSession(s *Session) error {
	if s == nil {
		return nil
	}
	req := &ua.CloseSessionRequest{DeleteSubscriptions: true}
	var res *ua.CloseSessionResponse
	return c.Send(req, func(v interface{}) error {
		return safeAssign(v, &res)
	})
}

// DetachSession removes the session from the client without closing it. The
// caller is responsible to close or re-activate the session. If the client
// does not have an active session the function returns no error.
func (c *Client) DetachSession() (*Session, error) {
	s := c.Session()
	c.session.Store((*Session)(nil))
	return s, nil
}

// Send sends the request via the secure channel and registers a handler for
// the response. If the client has an active session it injects the
// authenticaton token.
func (c *Client) Send(req interface{}, h func(interface{}) error) error {
	var authToken *ua.NodeID
	if s := c.Session(); s != nil {
		authToken = s.resp.AuthenticationToken
	}
	return c.sechan.Send(req, authToken, h)
}

// Node returns a node object which accesses its attributes
// through this client connection.
func (c *Client) Node(id *ua.NodeID) *Node {
	return &Node{ID: id, c: c}
}

func (c *Client) GetEndpoints() (*ua.GetEndpointsResponse, error) {
	req := &ua.GetEndpointsRequest{
		EndpointURL: c.endpointURL,
	}
	var res *ua.GetEndpointsResponse
	err := c.Send(req, func(v interface{}) error {
		return safeAssign(v, &res)
	})
	return res, err
}

// Read executes a synchronous read request.
//
// By default, the function requests the value of the nodes
// in the default encoding of the server.
func (c *Client) Read(req *ua.ReadRequest) (*ua.ReadResponse, error) {
	// clone the request and the ReadValueIDs to set defaults without
	// manipulating them in-place.
	rvs := make([]*ua.ReadValueID, len(req.NodesToRead))
	for i, rv := range req.NodesToRead {
		rc := &ua.ReadValueID{}
		*rc = *rv
		if rc.AttributeID == 0 {
			rc.AttributeID = ua.AttributeIDValue
		}
		if rc.DataEncoding == nil {
			rc.DataEncoding = &ua.QualifiedName{}
		}
		rvs[i] = rc
	}
	req = &ua.ReadRequest{
		MaxAge:             req.MaxAge,
		TimestampsToReturn: req.TimestampsToReturn,
		NodesToRead:        rvs,
	}

	var res *ua.ReadResponse
	err := c.Send(req, func(v interface{}) error {
		return safeAssign(v, &res)
	})
	return res, err
}

// Write executes a synchronous write request.
func (c *Client) Write(req *ua.WriteRequest) (*ua.WriteResponse, error) {
	var res *ua.WriteResponse
	err := c.Send(req, func(v interface{}) error {
		return safeAssign(v, &res)
	})
	return res, err
}

// Browse executes a synchronous browse request.
func (c *Client) Browse(req *ua.BrowseRequest) (*ua.BrowseResponse, error) {
	var res *ua.BrowseResponse
	err := c.Send(req, func(v interface{}) error {
		return safeAssign(v, &res)
	})
	return res, err
}

type Subscription struct {
	SubscriptionID            uint32
	RevisedPublishingInterval float64
	RevisedLifetimeCount      uint32
	RevisedMaxKeepAliveCount  uint32
	Channel                   chan PublishNotificationData
	stopPublishLoop           chan<- struct{}
}

type SubscriptionParameters struct {
	Interval                   time.Duration
	LifetimeCount              uint32
	MaxKeepAliveCount          uint32
	MaxNotificationsPerPublish uint32
	Priority                   uint8
	ChannelBufferSize          int
}

func NewDefaultSubscriptionParameters() *SubscriptionParameters {
	return &SubscriptionParameters{
		MaxNotificationsPerPublish: 10000,
		LifetimeCount:              10000,
		MaxKeepAliveCount:          3000,
		Interval:                   100 * time.Millisecond,
		Priority:                   0,
		ChannelBufferSize:          0,
	}
}

// Subscribe creates a Subscription with given parameters and starts one Publish loop to ensure
// there is at least one PublishLoop loop per Subscription. Additional Publish loops may be started
// and managed by clients by calling PublishLoop()
// see also NewDefaultSubscriptionParameters()
func (c *Client) Subscribe(params SubscriptionParameters) (*Subscription, error) {
	req := &ua.CreateSubscriptionRequest{
		RequestedPublishingInterval: float64(params.Interval / time.Millisecond),
		RequestedLifetimeCount:      params.LifetimeCount,
		RequestedMaxKeepAliveCount:  params.MaxKeepAliveCount,
		PublishingEnabled:           true,
		MaxNotificationsPerPublish:  params.MaxNotificationsPerPublish,
		Priority:                    params.Priority,
	}

	res, err := c.CreateSubscription(req)
	if err != nil {
		return nil, err
	}
	if res.ResponseHeader.ServiceResult != ua.StatusOK {
		return nil, res.ResponseHeader.ServiceResult
	}

	sub := Subscription{
		res.SubscriptionID,
		res.RevisedPublishingInterval,
		res.RevisedLifetimeCount,
		res.RevisedMaxKeepAliveCount,
		make(chan PublishNotificationData, params.ChannelBufferSize),
		c.PublishLoop(),
	}
	c.subscriptions[sub.SubscriptionID] = sub

	return &sub, nil
}

func (c *Client) CreateSubscription(req *ua.CreateSubscriptionRequest) (*ua.CreateSubscriptionResponse, error) {
	var res *ua.CreateSubscriptionResponse
	err := c.Send(req, func(v interface{}) error {
		return safeAssign(v, &res)
	})
	return res, err
}

// Unsubscribe() deletes the given Subscription from server and stops the Publish loop that was
// started by Subscribe()
func (c *Client) Unsubscribe(sub *Subscription) error {
	if registeredSub, ok := c.subscriptions[sub.SubscriptionID]; ok {
		close(registeredSub.stopPublishLoop)
		delete(c.subscriptions, sub.SubscriptionID)
	}

	res, err := c.DeleteSubscriptions([]uint32{sub.SubscriptionID})
	if err != nil {
		return err
	}
	if res.ResponseHeader.ServiceResult != ua.StatusOK {
		return res.ResponseHeader.ServiceResult
	}

	return nil
}

func (c *Client) DeleteSubscriptions(subIds []uint32) (*ua.DeleteSubscriptionsResponse, error) {
	req := &ua.DeleteSubscriptionsRequest{
		SubscriptionIDs: subIds,
	}
	var res *ua.DeleteSubscriptionsResponse
	err := c.Send(req, func(v interface{}) error {
		return safeAssign(v, &res)
	})
	return res, err
}

func NewMonitoredItemCreateRequestWithDefaults(nodeID *ua.NodeID, attributeID ua.AttributeID, clientHandle uint32) *ua.MonitoredItemCreateRequest {
	if attributeID == 0 {
		attributeID = ua.AttributeIDValue
	}
	readValueID := &ua.ReadValueID{
		NodeID:       nodeID,
		AttributeID:  attributeID,
		DataEncoding: &ua.QualifiedName{},
	}
	params := ua.MonitoringParameters{
		ClientHandle:     clientHandle,
		DiscardOldest:    true,
		Filter:           nil,
		QueueSize:        10,
		SamplingInterval: 0.0,
	}
	createReq := ua.MonitoredItemCreateRequest{
		ItemToMonitor:       readValueID,
		MonitoringMode:      ua.MonitoringModeReporting,
		RequestedParameters: &params,
	}
	return &createReq
}

type PublishNotificationData struct {
	SubscriptionID uint32
	Error          error
	Value          interface{}
}

// Publish() sends a single Publish request with given acknowledgements
func (c *Client) Publish(acks []*ua.SubscriptionAcknowledgement) (*ua.PublishResponse, error) {
	req := &ua.PublishRequest{
		SubscriptionAcknowledgements: acks,
	}

	var res *ua.PublishResponse
	err := c.Send(req, func(v interface{}) error {
		return safeAssign(v, &res)
	})
	return res, err

}

// PublishLoop() starts an infinite loop that sends PublishRequests and delivers received
// notifications to registered Subscriptions.
// Returns a channel which can be used to stop the loop
func (c *Client) PublishLoop() chan<- struct{} {
	quit := make(chan struct{})
	go func() {

		// Empty SubscriptionAcknowledgements for first PublishRequest
		var acks = make([]*ua.SubscriptionAcknowledgement, 0)

		for {
			select {
			case <-quit:
				return
			default:
				res, err := c.Publish(acks)
				if err != nil {
					if err == ua.StatusBadTimeout {
						continue
					} else if err == ua.StatusBadNoSubscription {
						// ignore it as probably the cause is that all subscriptions are already deleted,
						// but the publishing loop is still running and will be stopped shortly
						continue
					}
					errorData := PublishNotificationData{Error: err}
					// notify all subscriptions of error
					for _, sub := range c.subscriptions {
						go func(s Subscription) { s.Channel <- errorData }(sub)
					}
					continue
				}
				// Prepare SubscriptionAcknowledgement for next PublishRequest
				acks = make([]*ua.SubscriptionAcknowledgement, 0)
				for _, i := range res.AvailableSequenceNumbers {
					ack := &ua.SubscriptionAcknowledgement{
						SubscriptionID: res.SubscriptionID,
						SequenceNumber: i,
					}
					acks = append(acks, ack)
				}

				c.notifySubscription(res)
			}
		}
	}()
	return quit
}

func (c *Client) notifySubscription(response *ua.PublishResponse) {
	sub, ok := c.subscriptions[response.SubscriptionID]
	if !ok {
		debug.Printf("Unknown subscription: %v", response.SubscriptionID)
		return
	}

	// Check for errors
	status := ua.StatusOK
	for _, res := range response.Results {
		if res != ua.StatusOK {
			status = res
			break
		}
	}

	if status != ua.StatusOK {
		sub.Channel <- PublishNotificationData{
			SubscriptionID: response.SubscriptionID,
			Error:          status,
		}
		return
	}

	if response.NotificationMessage == nil {
		sub.Channel <- PublishNotificationData{
			SubscriptionID: response.SubscriptionID,
			Error:          fmt.Errorf("empty NotificationMessage"),
		}
		return
	}

	// Part 4, 7.21 NotificationMessage
	for _, data := range response.NotificationMessage.NotificationData {
		// Part 4, 7.20 NotificationData parameters
		if data == nil || data.Value == nil {
			sub.Channel <- PublishNotificationData{
				SubscriptionID: response.SubscriptionID,
				Error:          fmt.Errorf("missing NotificationData parameter"),
			}
			continue
		}

		switch data.Value.(type) {
		// Part 4, 7.20.2 DataChangeNotification parameter
		// Part 4, 7.20.3 EventNotificationList parameter
		// Part 4, 7.20.4 StatusChangeNotification parameter
		case *ua.DataChangeNotification,
			*ua.EventNotificationList,
			*ua.StatusChangeNotification:
			sub.Channel <- PublishNotificationData{
				SubscriptionID: response.SubscriptionID,
				Value:          data.Value,
			}

		// Error
		default:
			sub.Channel <- PublishNotificationData{
				SubscriptionID: response.SubscriptionID,
				Error:          fmt.Errorf("unknown NotificationData parameter: %T", data.Value),
			}
		}
	}
}

func (c *Client) CreateMonitoredItems(subID uint32, ts ua.TimestampsToReturn, items ...*ua.MonitoredItemCreateRequest) (*ua.CreateMonitoredItemsResponse, error) {
	if subID == 0 {
		return nil, ua.StatusBadSubscriptionIDInvalid
	}

	// Part 4, 5.12.2.2 CreateMonitoredItems Service Parameters
	req := &ua.CreateMonitoredItemsRequest{
		SubscriptionID:     subID,
		TimestampsToReturn: ts,
		ItemsToCreate:      items,
	}

	var res *ua.CreateMonitoredItemsResponse
	err := c.Send(req, func(v interface{}) error {
		return safeAssign(v, &res)
	})
	return res, err
}

func (c *Client) DeleteMonitoredItems(subID uint32, monitoredItemIDs ...uint32) (*ua.DeleteMonitoredItemsResponse, error) {
	req := &ua.DeleteMonitoredItemsRequest{
		MonitoredItemIDs: monitoredItemIDs,
		SubscriptionID:   subID,
	}
	var res *ua.DeleteMonitoredItemsResponse
	err := c.Send(req, func(v interface{}) error {
		return safeAssign(v, &res)
	})
	return res, err
}

func (c *Client) HistoryReadRawModified(nodes []*ua.HistoryReadValueID, details *ua.ReadRawModifiedDetails) (*ua.HistoryReadResponse, error) {
	// Part 4, 5.10.3 HistoryRead
	req := &ua.HistoryReadRequest{
		TimestampsToReturn: ua.TimestampsToReturnBoth,
		NodesToRead:        nodes,
		// Part 11, 6.4 HistoryReadDetails parameters
		HistoryReadDetails: &ua.ExtensionObject{
			TypeID:       ua.NewFourByteExpandedNodeID(0, id.ReadRawModifiedDetails_Encoding_DefaultBinary),
			EncodingMask: ua.ExtensionObjectBinary,
			Value:        details,
		},
	}

	var res *ua.HistoryReadResponse
	err := c.Send(req, func(v interface{}) error {
		return safeAssign(v, &res)
	})
	return res, err
}

// safeAssign implements a type-safe assign from T to *T.
func safeAssign(t, ptrT interface{}) error {
	if reflect.TypeOf(t) != reflect.TypeOf(ptrT).Elem() {
		return InvalidResponseTypeError{t, ptrT}
	}

	// this is *ptrT = t
	reflect.ValueOf(ptrT).Elem().Set(reflect.ValueOf(t))
	return nil
}

type InvalidResponseTypeError struct {
	got, want interface{}
}

func (e InvalidResponseTypeError) Error() string {
	return fmt.Sprintf("invalid response: got %T want %T", e.got, e.want)
}
