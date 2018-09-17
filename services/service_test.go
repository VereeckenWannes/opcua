// Copyright 2018 gopcua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package services

import (
	"testing"
	"time"

	"github.com/wmnsk/gopcua/datatypes"
)

var testServiceBytes = [][]byte{
	{ // OpenSecureChannelRequest
		// TypeID
		0x01, 0x00, 0xbe, 0x01,
		// RequestHeader
		0x00, 0x00, 0x00, 0x98, 0x67, 0xdd, 0xfd, 0x30,
		0xd4, 0x01, 0x01, 0x00, 0x00, 0x00, 0xff, 0x03,
		0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00,
		// ClientProtocolVersion
		0x00, 0x00, 0x00, 0x00,
		// SecurityTokenRequestType
		0x00, 0x00, 0x00, 0x00,
		// MessageSecurityMode
		0x01, 0x00, 0x00, 0x00,
		// ClientNonce
		0xff, 0xff, 0xff, 0xff,
		// RequestedLifetime
		0x80, 0x8d, 0x5b, 0x00,
	},
	{ // OpenSecureChannelResponse
		// TypeID
		0x01, 0x00, 0xc1, 0x01,
		// ResponseHeader
		0x00, 0x98, 0x67, 0xdd, 0xfd, 0x30, 0xd4, 0x01,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x02, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00,
		0x00, 0x66, 0x6f, 0x6f, 0x03, 0x00, 0x00, 0x00,
		0x62, 0x61, 0x72, 0x00, 0x00, 0x00,
		// ServerProtocolVersion
		0x00, 0x00, 0x00, 0x00,
		// SecurityToken
		0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00,
		0x00, 0x98, 0x67, 0xdd, 0xfd, 0x30, 0xd4, 0x01,
		0x80, 0x8d, 0x5b, 0x00,
		// ServerNonce
		0x01, 0x00, 0x00, 0x00, 0xff,
	},
	{ // GetEndpointsRequest
		// TypeID
		0x01, 0x00, 0xac, 0x01,
		// RequestHeader
		0x00, 0x00, 0x00, 0x98, 0x67, 0xdd, 0xfd, 0x30,
		0xd4, 0x01, 0x01, 0x00, 0x00, 0x00, 0xff, 0x03,
		0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00,
		// ClientProtocolVersion
		0x26, 0x00, 0x00, 0x00, 0x6f, 0x70, 0x63, 0x2e,
		0x74, 0x63, 0x70, 0x3a, 0x2f, 0x2f, 0x77, 0x6f,
		0x77, 0x2e, 0x69, 0x74, 0x73, 0x2e, 0x65, 0x61,
		0x73, 0x79, 0x3a, 0x31, 0x31, 0x31, 0x31, 0x31,
		0x2f, 0x55, 0x41, 0x2f, 0x53, 0x65, 0x72, 0x76,
		0x65, 0x72,
		// LocaleIDs
		0x00, 0x00, 0x00, 0x00,
		// ProfileURIs
		0x00, 0x00, 0x00, 0x00,
	},
	{ // GetEndpointsResponse
		// TypeID
		0x01, 0x00, 0xaf, 0x01,
		// ResponseHeader
		0x00, 0x98, 0x67, 0xdd, 0xfd, 0x30, 0xd4, 0x01,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// Endpoints
		// ArraySize: 2
		0x02, 0x00, 0x00, 0x00,
		// EndpointURI
		0x06, 0x00, 0x00, 0x00, 0x65, 0x70, 0x2d, 0x75, 0x72, 0x6c,
		// Server (ApplicationDescription)
		// ApplicationURI
		0x07, 0x00, 0x00, 0x00, 0x61, 0x70, 0x70, 0x2d, 0x75, 0x72, 0x69,
		// ProductURI
		0x08, 0x00, 0x00, 0x00, 0x70, 0x72, 0x6f, 0x64, 0x2d, 0x75, 0x72, 0x69,
		// ApplicationName
		0x02, 0x08, 0x00, 0x00, 0x00, 0x61, 0x70, 0x70, 0x2d,
		0x6e, 0x61, 0x6d, 0x65,
		// ApplicationType
		0x00, 0x00, 0x00, 0x00,
		// GatewayServerURI
		0x06, 0x00, 0x00, 0x00, 0x67, 0x77, 0x2d, 0x75, 0x72, 0x69,
		// DiscoveryProfileURI
		0x08, 0x00, 0x00, 0x00, 0x70, 0x72, 0x6f, 0x66, 0x2d, 0x75, 0x72, 0x69,
		// DiscoveryURIs
		0x02, 0x00, 0x00, 0x00,
		0x0c, 0x00, 0x00, 0x00, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x2d, 0x75, 0x72, 0x69, 0x2d, 0x31,
		0x0c, 0x00, 0x00, 0x00, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x2d, 0x75, 0x72, 0x69, 0x2d, 0x32,
		// ServerCertificate
		0xff, 0xff, 0xff, 0xff,
		// MessageSecurityMode
		0x01, 0x00, 0x00, 0x00,
		// SecurityPolicyURI
		0x07, 0x00, 0x00, 0x00, 0x73, 0x65, 0x63, 0x2d, 0x75, 0x72, 0x69,
		// UserIdentityTokens
		// ArraySize
		0x02, 0x00, 0x00, 0x00,
		// PolicyID
		0x01, 0x00, 0x00, 0x00, 0x31,
		// TokenType
		0x00, 0x00, 0x00, 0x00,
		// IssuedTokenType
		0x0c, 0x00, 0x00, 0x00, 0x69, 0x73, 0x73, 0x75, 0x65, 0x64, 0x2d, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
		// IssuerEndpointURI
		0x0a, 0x00, 0x00, 0x00, 0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x2d, 0x75, 0x72, 0x69,
		// SecurityPolicyURI
		0x07, 0x00, 0x00, 0x00, 0x73, 0x65, 0x63, 0x2d, 0x75, 0x72, 0x69,
		// PolicyID
		0x01, 0x00, 0x00, 0x00, 0x31,
		// TokenType
		0x00, 0x00, 0x00, 0x00,
		// IssuedTokenType
		0x0c, 0x00, 0x00, 0x00, 0x69, 0x73, 0x73, 0x75, 0x65, 0x64, 0x2d, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
		// IssuerEndpointURI
		0x0a, 0x00, 0x00, 0x00, 0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x2d, 0x75, 0x72, 0x69,
		// SecurityPolicyURI
		0x07, 0x00, 0x00, 0x00, 0x73, 0x65, 0x63, 0x2d, 0x75, 0x72, 0x69,
		// TransportProfileURI
		0x09, 0x00, 0x00, 0x00, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x2d, 0x75, 0x72, 0x69,
		// SecurityLevel
		0x00,
		// EndpointURI
		0x06, 0x00, 0x00, 0x00, 0x65, 0x70, 0x2d, 0x75, 0x72, 0x6c,
		// Server (ApplicationDescription)
		// ApplicationURI
		0x07, 0x00, 0x00, 0x00, 0x61, 0x70, 0x70, 0x2d, 0x75, 0x72, 0x69,
		// ProductURI
		0x08, 0x00, 0x00, 0x00, 0x70, 0x72, 0x6f, 0x64, 0x2d, 0x75, 0x72, 0x69,
		// ApplicationName
		0x02, 0x08, 0x00, 0x00, 0x00, 0x61, 0x70, 0x70, 0x2d,
		0x6e, 0x61, 0x6d, 0x65,
		// ApplicationType
		0x00, 0x00, 0x00, 0x00,
		// GatewayServerURI
		0x06, 0x00, 0x00, 0x00, 0x67, 0x77, 0x2d, 0x75, 0x72, 0x69,
		// DiscoveryProfileURI
		0x08, 0x00, 0x00, 0x00, 0x70, 0x72, 0x6f, 0x66, 0x2d, 0x75, 0x72, 0x69,
		// DiscoveryURIs
		0x02, 0x00, 0x00, 0x00,
		0x0c, 0x00, 0x00, 0x00, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x2d, 0x75, 0x72, 0x69, 0x2d, 0x31,
		0x0c, 0x00, 0x00, 0x00, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x2d, 0x75, 0x72, 0x69, 0x2d, 0x32,
		// ServerCertificate
		0xff, 0xff, 0xff, 0xff,
		// MessageSecurityMode
		0x01, 0x00, 0x00, 0x00,
		// SecurityPolicyURI
		0x07, 0x00, 0x00, 0x00, 0x73, 0x65, 0x63, 0x2d, 0x75, 0x72, 0x69,
		// UserIdentityTokens
		// ArraySize
		0x02, 0x00, 0x00, 0x00,
		// PolicyID
		0x01, 0x00, 0x00, 0x00, 0x31,
		// TokenType
		0x00, 0x00, 0x00, 0x00,
		// IssuedTokenType
		0x0c, 0x00, 0x00, 0x00, 0x69, 0x73, 0x73, 0x75, 0x65, 0x64, 0x2d, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
		// IssuerEndpointURI
		0x0a, 0x00, 0x00, 0x00, 0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x2d, 0x75, 0x72, 0x69,
		// SecurityPolicyURI
		0x07, 0x00, 0x00, 0x00, 0x73, 0x65, 0x63, 0x2d, 0x75, 0x72, 0x69,
		// PolicyID
		0x01, 0x00, 0x00, 0x00, 0x31,
		// TokenType
		0x00, 0x00, 0x00, 0x00,
		// IssuedTokenType
		0x0c, 0x00, 0x00, 0x00, 0x69, 0x73, 0x73, 0x75, 0x65, 0x64, 0x2d, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
		// IssuerEndpointURI
		0x0a, 0x00, 0x00, 0x00, 0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x2d, 0x75, 0x72, 0x69,
		// SecurityPolicyURI
		0x07, 0x00, 0x00, 0x00, 0x73, 0x65, 0x63, 0x2d, 0x75, 0x72, 0x69,
		// TransportProfileURI
		0x09, 0x00, 0x00, 0x00, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x2d, 0x75, 0x72, 0x69,
		// SecurityLevel
		0x00,
	},
	{ // CreateSessionRequest
		// TypeID
		0x01, 0x00, 0xcd, 0x01,
		// RequestHeader
		// AuthenticationToken
		0x00, 0x00,
		// Timestamp
		0x00, 0x98, 0x67, 0xdd, 0xfd, 0x30, 0xd4, 0x01,
		// RequestHandle
		0x01, 0x00, 0x00, 0x00,
		// ReturnDiagnostics
		0xff, 0x03, 0x00, 0x00,
		// AuditEntryID
		0xff, 0xff, 0xff, 0xff,
		// TimeoutHint
		0x00, 0x00, 0x00, 0x00,
		// AdditionalHeader
		0x00, 0x00, 0x00,
		// ClientDescription: ApplicationDescription
		// ApplicationURI
		0x07, 0x00, 0x00, 0x00, 0x61, 0x70, 0x70, 0x2d, 0x75, 0x72, 0x69,
		// ProductURI
		0x08, 0x00, 0x00, 0x00, 0x70, 0x72, 0x6f, 0x64, 0x2d, 0x75, 0x72, 0x69,
		// ApplicationName
		0x02, 0x08, 0x00, 0x00, 0x00, 0x61, 0x70, 0x70, 0x2d, 0x6e, 0x61, 0x6d, 0x65,
		// ApplicationType
		0x01, 0x00, 0x00, 0x00,
		// GatewayServerURI
		0xff, 0xff, 0xff, 0xff,
		// DiscoveryProfileURI
		0xff, 0xff, 0xff, 0xff,
		// DiscoveryURLs
		0x00, 0x00, 0x00, 0x00,
		// ServerURI
		0x0a, 0x00, 0x00, 0x00, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x75, 0x72, 0x69,
		// EndpointURL
		0x0c, 0x00, 0x00, 0x00, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2d, 0x75, 0x72, 0x6c,
		// SessionName
		0x0c, 0x00, 0x00, 0x00, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2d, 0x6e, 0x61, 0x6d, 0x65,
		// ClientNonce
		0xff, 0xff, 0xff, 0xff,
		// ClientCertificate
		0xff, 0xff, 0xff, 0xff,
		// RequestedTimeout
		0x80, 0x8d, 0x5b, 0x00, 0x00, 0x00, 0x00, 0x00,
		// MaxResponseMessageSize
		0xfe, 0xff, 0x00, 0x00,
	},
	{ // CreateSessionResponse
		// TypeID
		0x01, 0x00, 0xd0, 0x01,
		// ResponseHeader
		// Timestamp
		0x00, 0x98, 0x67, 0xdd, 0xfd, 0x30, 0xd4, 0x01,
		// RequestHandle
		0x01, 0x00, 0x00, 0x00,
		// ServiceResult
		0x00, 0x00, 0x00, 0x00,
		// ServiceDiagnostics
		0x00,
		// StringTable
		0x00, 0x00, 0x00, 0x00,
		// AdditionalHeader
		0x00, 0x00, 0x00,
		// SessionID
		0x02, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
		// AuthenticationToken
		0x05, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x08,
		0x22, 0x87, 0x62, 0xba, 0x81, 0xe1, 0x11, 0xa6,
		0x43, 0xf8, 0x77, 0x7b, 0xc6, 0x2f, 0xc8,
		// RevisedSessionTimeout
		0x80, 0x8d, 0x5b, 0x00, 0x00, 0x00, 0x00, 0x00,
		// ServerNonce
		0xff, 0xff, 0xff, 0xff,
		// ServerCertificate
		0xff, 0xff, 0xff, 0xff,
		// ServerEndpoints
		// ArraySize
		0x02, 0x00, 0x00, 0x00,
		// EndpointURL
		0x06, 0x00, 0x00, 0x00, 0x65, 0x70, 0x2d, 0x75, 0x72, 0x6c,
		// Server
		// ApplicationURI
		0x07, 0x00, 0x00, 0x00, 0x61, 0x70, 0x70, 0x2d, 0x75, 0x72, 0x69,
		// ProductURI
		0x08, 0x00, 0x00, 0x00, 0x70, 0x72, 0x6f, 0x64, 0x2d, 0x75, 0x72, 0x69,
		// ApplicationName
		0x02, 0x08, 0x00, 0x00, 0x00, 0x61, 0x70, 0x70, 0x2d, 0x6e, 0x61, 0x6d, 0x65,
		// ApplicationType
		0x00, 0x00, 0x00, 0x00,
		// GatewayServerURI
		0x06, 0x00, 0x00, 0x00, 0x67, 0x77, 0x2d, 0x75, 0x72, 0x69,
		// DiscoveryProfileURI
		0x08, 0x00, 0x00, 0x00, 0x70, 0x72, 0x6f, 0x66, 0x2d, 0x75, 0x72, 0x69,
		// DiscoveryURLs
		0x02, 0x00, 0x00, 0x00,
		0x0c, 0x00, 0x00, 0x00, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x2d, 0x75, 0x72, 0x69, 0x2d, 0x31,
		0x0c, 0x00, 0x00, 0x00, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x2d, 0x75, 0x72, 0x69, 0x2d, 0x32,
		// ServerCertificate
		0xff, 0xff, 0xff, 0xff,
		// MessageSecurityMode
		0x01, 0x00, 0x00, 0x00,
		// SecurityPolicyURI
		0x07, 0x00, 0x00, 0x00, 0x73, 0x65, 0x63, 0x2d, 0x75, 0x72, 0x69,
		// UserIdentityTokens
		0x02, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00, 0x31, 0x00, 0x00, 0x00, 0x00,
		0x0c, 0x00, 0x00, 0x00, 0x69, 0x73, 0x73, 0x75, 0x65, 0x64, 0x2d, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
		0x0a, 0x00, 0x00, 0x00, 0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x2d, 0x75, 0x72, 0x69,
		0x07, 0x00, 0x00, 0x00, 0x73, 0x65, 0x63, 0x2d, 0x75, 0x72, 0x69,
		0x01, 0x00, 0x00, 0x00, 0x31, 0x00, 0x00, 0x00, 0x00,
		0x0c, 0x00, 0x00, 0x00, 0x69, 0x73, 0x73, 0x75, 0x65, 0x64, 0x2d, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
		0x0a, 0x00, 0x00, 0x00, 0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x2d, 0x75, 0x72, 0x69,
		0x07, 0x00, 0x00, 0x00, 0x73, 0x65, 0x63, 0x2d, 0x75, 0x72, 0x69,
		// TransportProfileURI
		0x09, 0x00, 0x00, 0x00, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x2d, 0x75, 0x72, 0x69,
		// SecurityLevel
		0x00,
		// EndpointURL
		0x06, 0x00, 0x00, 0x00, 0x65, 0x70, 0x2d, 0x75, 0x72, 0x6c,
		// Server
		// ApplicationURI
		0x07, 0x00, 0x00, 0x00, 0x61, 0x70, 0x70, 0x2d, 0x75, 0x72, 0x69,
		// ProductURI
		0x08, 0x00, 0x00, 0x00, 0x70, 0x72, 0x6f, 0x64, 0x2d, 0x75, 0x72, 0x69,
		// ApplicationName
		0x02, 0x08, 0x00, 0x00, 0x00, 0x61, 0x70, 0x70, 0x2d, 0x6e, 0x61, 0x6d, 0x65,
		// ApplicationType
		0x00, 0x00, 0x00, 0x00,
		// GatewayServerURI
		0x06, 0x00, 0x00, 0x00, 0x67, 0x77, 0x2d, 0x75, 0x72, 0x69,
		// DiscoveryProfileURI
		0x08, 0x00, 0x00, 0x00, 0x70, 0x72, 0x6f, 0x66, 0x2d, 0x75, 0x72, 0x69,
		// DiscoveryURLs
		0x02, 0x00, 0x00, 0x00,
		0x0c, 0x00, 0x00, 0x00, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x2d, 0x75, 0x72, 0x69, 0x2d, 0x31,
		0x0c, 0x00, 0x00, 0x00, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x2d, 0x75, 0x72, 0x69, 0x2d, 0x32,
		// ServerCertificate
		0xff, 0xff, 0xff, 0xff,
		// MessageSecurityMode
		0x01, 0x00, 0x00, 0x00,
		// SecurityPolicyURI
		0x07, 0x00, 0x00, 0x00, 0x73, 0x65, 0x63, 0x2d, 0x75, 0x72, 0x69,
		// UserIdentityTokens
		0x02, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00, 0x31, 0x00, 0x00, 0x00, 0x00,
		0x0c, 0x00, 0x00, 0x00, 0x69, 0x73, 0x73, 0x75, 0x65, 0x64, 0x2d, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
		0x0a, 0x00, 0x00, 0x00, 0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x2d, 0x75, 0x72, 0x69,
		0x07, 0x00, 0x00, 0x00, 0x73, 0x65, 0x63, 0x2d, 0x75, 0x72, 0x69,
		0x01, 0x00, 0x00, 0x00, 0x31, 0x00, 0x00, 0x00, 0x00,
		0x0c, 0x00, 0x00, 0x00, 0x69, 0x73, 0x73, 0x75, 0x65, 0x64, 0x2d, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
		0x0a, 0x00, 0x00, 0x00, 0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x2d, 0x75, 0x72, 0x69,
		0x07, 0x00, 0x00, 0x00, 0x73, 0x65, 0x63, 0x2d, 0x75, 0x72, 0x69,
		// TransportProfileURI
		0x09, 0x00, 0x00, 0x00, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x2d, 0x75, 0x72, 0x69,
		// SecurityLevel
		0x00,
		// ServerSoftwareCertificates
		0x00, 0x00, 0x00, 0x00,
		// ServerSignature
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		// MaxRequestMessageSize
		0xfe, 0xff, 0x00, 0x00,
	},
	{ // CloseSecureChannelRequest
		// TypeID
		0x01, 0x00, 0xc4, 0x01,
		// RequestHeader
		0x00, 0x00, 0x00, 0x98, 0x67, 0xdd, 0xfd, 0x30,
		0xd4, 0x01, 0x01, 0x00, 0x00, 0x00, 0xff, 0x03,
		0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00,
		// SecureChannelID
		0x01, 0x00, 0x00, 0x00,
	},
	{ // CloseSecureChannelResponse
		// TypeID
		0x01, 0x00, 0xc7, 0x01,
		// ResponseHeader
		0x00, 0x98, 0x67, 0xdd, 0xfd, 0x30, 0xd4, 0x01,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x02, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00,
		0x00, 0x66, 0x6f, 0x6f, 0x03, 0x00, 0x00, 0x00,
		0x62, 0x61, 0x72, 0x00, 0x00, 0x00,
	},
	{ // CloseSessionRequest
		// TypeID
		0x01, 0x00, 0xd9, 0x01,
		// RequestHeader
		0x05, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x08,
		0x22, 0x87, 0x62, 0xba, 0x81, 0xe1, 0x11, 0xa6,
		0x43, 0xf8, 0x77, 0x7b, 0xc6, 0x2f, 0xc8, 0x00,
		0x98, 0x67, 0xdd, 0xfd, 0x30, 0xd4, 0x01, 0x01,
		0x00, 0x00, 0x00, 0xff, 0x03, 0x00, 0x00, 0xff,
		0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00,
		// DeleteSubscription
		0x01,
	},
	{ // CloseSessionResponse
		// TypeID
		0x01, 0x00, 0xdc, 0x01,
		// ResponseHeader
		0x00, 0x98, 0x67, 0xdd, 0xfd, 0x30, 0xd4, 0x01,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x02, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00,
		0x00, 0x66, 0x6f, 0x6f, 0x03, 0x00, 0x00, 0x00,
		0x62, 0x61, 0x72, 0x00, 0x00, 0x00,
	},
}

func TestDecode(t *testing.T) {
	t.Run("open-sec-chan-req", func(t *testing.T) {
		t.Parallel()
		o, err := Decode(testServiceBytes[0])
		if err != nil {
			t.Fatalf("Failed to decode Service: %s", err)
		}

		osc, ok := o.(*OpenSecureChannelRequest)
		if !ok {
			t.Fatalf("Failed to assert type.")
		}

		switch {
		case o.ServiceType() != ServiceTypeOpenSecureChannelRequest:
			t.Errorf("ServiceType doesn't Match. Want: %d, Got: %d", ServiceTypeOpenSecureChannelRequest, o.ServiceType())
		case osc.ClientProtocolVersion != 0:
			t.Errorf("ClientProtocolVersion doesn't Match. Want: %d, Got: %d", 0, osc.ClientProtocolVersion)
		case osc.SecurityTokenRequestType != 0:
			t.Errorf("SecurityTokenRequestType doesn't Match. Want: %d, Got: %d", 0, osc.SecurityTokenRequestType)
		case osc.MessageSecurityMode != 1:
			t.Errorf("MessageSecurityMode doesn't Match. Want: %d, Got: %d", 1, osc.MessageSecurityMode)
		case osc.ClientNonce.Get() != nil:
			t.Errorf("ClientNonce doesn't Match. Want: %v, Got: %v", nil, osc.ClientNonce.Get())
		case osc.RequestedLifetime != 6000000:
			t.Errorf("RequestedLifetime doesn't Match. Want: %d, Got: %d", 6000000, osc.RequestedLifetime)
		}
		t.Log(o.String())
	})
	t.Run("open-sec-chan-res", func(t *testing.T) {
		t.Parallel()
		o, err := Decode(testServiceBytes[1])
		if err != nil {
			t.Fatalf("Failed to decode Service: %s", err)
		}

		osc, ok := o.(*OpenSecureChannelResponse)
		if !ok {
			t.Fatalf("Failed to assert type.")
		}

		switch {
		case o.ServiceType() != ServiceTypeOpenSecureChannelResponse:
			t.Errorf("ServiceType doesn't Match. Want: %d, Got: %d", ServiceTypeOpenSecureChannelResponse, o.ServiceType())
		case osc.ServerProtocolVersion != 0:
			t.Errorf("ServerProtocolVersion doesn't Match. Want: %d, Got: %d", 0, osc.ServerProtocolVersion)
		case osc.SecurityToken.ChannelID != 1:
			t.Errorf("SecurityToken.ChannelID doesn't Match. Want: %d, Got: %d", 1, osc.SecurityToken.ChannelID)
		case osc.SecurityToken.TokenID != 2:
			t.Errorf("SecurityToken.TokenID doesn't Match. Want: %d, Got: %d", 2, osc.SecurityToken.TokenID)
		case osc.SecurityToken.CreatedAt != time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC):
			t.Errorf("SecurityToken.CreatedAt doesn't Match. Want: %v, Got: %v", time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC), osc.SecurityToken.CreatedAt)
		case osc.SecurityToken.RevisedLifetime != 6000000:
			t.Errorf("SecurityToken.RevisedLifetime doesn't Match. Want: %d, Got: %d", 6000000, osc.SecurityToken.RevisedLifetime)
		case osc.ServerNonce.Get()[0] != 255:
			t.Errorf("ServerNonce doesn't Match. Want: %v, Got: %v", 255, osc.ServerNonce.Get()[0])
		}
		t.Log(o.String())
	})
	t.Run("get-endpoint-req", func(t *testing.T) {
		t.Parallel()
		g, err := Decode(testServiceBytes[2])
		if err != nil {
			t.Fatalf("Failed to decode Service: %s", err)
		}

		gep, ok := g.(*GetEndpointsRequest)
		if !ok {
			t.Fatalf("Failed to assert type.")
		}

		switch {
		case g.ServiceType() != ServiceTypeGetEndpointsRequest:
			t.Errorf("ServiceType doesn't Match. Want: %d, Got: %d", ServiceTypeGetEndpointsRequest, g.ServiceType())
		case gep.EndpointURL.Get() != "opc.tcp://wow.its.easy:11111/UA/Server":
			t.Errorf("EndpointURL doesn't Match. Want: %s, Got: %s", "opc.tcp://wow.its.easy:11111/UA/Server", gep.EndpointURL.Get())
		case gep.LocaleIDs.ArraySize != 0:
			t.Errorf("LocaleIDs.ArraySize doesn't Match. Want: %d, Got: %d", 0, gep.LocaleIDs.ArraySize)
		case gep.ProfileURIs.ArraySize != 0:
			t.Errorf("ProfileURIs.ArraySize doesn't Match. Want: %d, Got: %d", 0, gep.ProfileURIs.ArraySize)
		}
		t.Log(g.String())
	})
	t.Run("get-endpoint-res", func(t *testing.T) {
		t.Parallel()
		g, err := Decode(testServiceBytes[3])
		if err != nil {
			t.Fatalf("Failed to decode Service: %s", err)
		}

		gep, ok := g.(*GetEndpointsResponse)
		if !ok {
			t.Fatalf("Failed to assert type.")
		}

		if g.ServiceType() != ServiceTypeGetEndpointsResponse {
			t.Errorf("ServiceType doesn't Match. Want: %d, Got: %d", ServiceTypeGetEndpointsResponse, g.ServiceType())
		}

		for _, ep := range gep.Endpoints.EndpointDescriptions {
			switch {
			case ep.EndpointURL.Get() != "ep-url":
				t.Errorf("EndpointURL doesn't match. Want: %s, Got: %s", "ep-url", ep.EndpointURL.Get())
			case ep.ServerCertificate.Get() != nil:
				t.Errorf("ServerCertificate doesn't match. Want: %v, Got: %v", nil, ep.ServerCertificate.Get())
			case ep.MessageSecurityMode != SecModeNone:
				t.Errorf("MessageSecurityMode doesn't match. Want: %d, Got: %d", SecModeNone, ep.MessageSecurityMode)
			case ep.SecurityPolicyURI.Get() != "sec-uri":
				t.Errorf("SecurityPolicyURI doesn't match. Want: %s, Got: %s", "sec-uri", ep.SecurityPolicyURI.Get())
			case ep.TransportProfileURI.Get() != "trans-uri":
				t.Errorf("TransportProfileURI doesn't match. Want: %s, Got: %s", "trans-uri", ep.TransportProfileURI.Get())
			case ep.SecurityLevel != 0:
				t.Errorf("SecurityLevel doesn't match. Want: %d, Got: %d", 0, ep.SecurityLevel)
			}
			t.Log(ep.String())
		}

		t.Log(gep.String())
	})
	t.Run("create-session-req", func(t *testing.T) {
		t.Parallel()
		c, err := Decode(testServiceBytes[4])
		if err != nil {
			t.Fatalf("Failed to decode Service: %s", err)
		}

		cs, ok := c.(*CreateSessionRequest)
		if !ok {
			t.Fatalf("Failed to assert type.")
		}

		if c.ServiceType() != ServiceTypeCreateSessionRequest {
			t.Errorf("ServiceType doesn't Match. Want: %d, Got: %d", ServiceTypeCreateSessionRequest, c.ServiceType())
		}

		switch {
		case cs.ServerURI.Get() != "server-uri":
			t.Errorf("ServerURI doesn't match. Want: %s, Got: %s", "server-uri", cs.ServerURI.Get())
		case cs.EndpointURL.Get() != "endpoint-url":
			t.Errorf("EndpointURL doesn't match. Want: %s, Got: %s", "endpoint-url", cs.EndpointURL.Get())
		case cs.SessionName.Get() != "session-name":
			t.Errorf("SessionName doesn't match. Want: %s, Got: %s", "session-name", cs.SessionName.Get())
		case cs.ClientNonce.Get() != nil:
			t.Errorf("ClientNonce doesn't match. Want: %v, Got: %v", nil, cs.ClientNonce.Get())
		case cs.ClientCertificate.Get() != nil:
			t.Errorf("ClientCertificate doesn't match. Want: %v, Got: %v", nil, cs.ClientCertificate.Get())
		case cs.RequestedSessionTimeout != 6000000:
			t.Errorf("RequestedSessionTimeout doesn't match. Want: %d, Got: %d", 6000000, cs.RequestedSessionTimeout)
		case cs.MaxResponseMessageSize != 65534:
			t.Errorf("MaxResponseMessageSize doesn't match. Want: %d, Got: %d", 65534, cs.MaxResponseMessageSize)
		}
		t.Log(cs.String())
	})
	t.Run("create-session-res", func(t *testing.T) {
		t.Parallel()
		c, err := Decode(testServiceBytes[5])
		if err != nil {
			t.Fatalf("Failed to decode Service: %s", err)
		}

		cs, ok := c.(*CreateSessionResponse)
		if !ok {
			t.Fatalf("Failed to assert type.")
		}

		if c.ServiceType() != ServiceTypeCreateSessionResponse {
			t.Errorf("ServiceType doesn't Match. Want: %d, Got: %d", ServiceTypeCreateSessionResponse, c.ServiceType())
		}

		sessionID, ok := cs.SessionID.(*datatypes.NumericNodeID)
		if !ok {
			t.Fatalf("Failed to assert session id type.")
		}

		if _, ok = cs.AuthenticationToken.(*datatypes.OpaqueNodeID); !ok {
			t.Fatalf("Failed to assert session id type.")
		}

		switch {
		case sessionID.Identifier != 1:
			t.Errorf("SessionID doesn't match. Want: %d, Got: %d", 1, sessionID.Identifier)
		// case authenticationToken.Identifier != 1:
		// 	t.Errorf("AuthenticationToken doesn't match. Want: %d, Got: %d", 1, authenticationToken.Identifier)
		case cs.RevisedSessionTimeout != 6000000:
			t.Errorf("RevisedSessionTimeout doesn't match. Want: %d, Got: %d", 6000000, cs.RevisedSessionTimeout)
		case cs.ServerNonce.Get() != nil:
			t.Errorf("ServerNonce doesn't match. Want: %v, Got: %v", nil, cs.ServerNonce.Get())
		case cs.ServerCertificate.Get() != nil:
			t.Errorf("ServerCertificate doesn't match. Want: %v, Got: %v", nil, cs.ServerCertificate.Get())
		case cs.MaxRequestMessageSize != 65534:
			t.Errorf("MaxRequestMessageSize doesn't match. Want: %d, Got: %d", 65534, cs.MaxRequestMessageSize)
		}
		t.Log(cs.String())
	})
	t.Run("close-sec-chan-req", func(t *testing.T) {
		t.Parallel()
		c, err := Decode(testServiceBytes[6])
		if err != nil {
			t.Fatalf("Failed to decode Service: %s", err)
		}

		csc, ok := c.(*CloseSecureChannelRequest)
		if !ok {
			t.Fatalf("Failed to assert type.")
		}

		switch {
		case c.ServiceType() != ServiceTypeCloseSecureChannelRequest:
			t.Errorf("ServiceType doesn't Match. Want: %d, Got: %d", ServiceTypeCloseSecureChannelRequest, c.ServiceType())
		case csc.SecureChannelID != 1:
			t.Errorf("SecureChannelID doesn't Match. Want: %d, Got: %d", 1, csc.SecureChannelID)
		}
		t.Log(c.String())
	})
	t.Run("close-sec-chan-res", func(t *testing.T) {
		t.Parallel()
		c, err := Decode(testServiceBytes[7])
		if err != nil {
			t.Fatalf("Failed to decode Service: %s", err)
		}

		_, ok := c.(*CloseSecureChannelResponse)
		if !ok {
			t.Fatalf("Failed to assert type.")
		}

		switch {
		case c.ServiceType() != ServiceTypeCloseSecureChannelResponse:
			t.Errorf("ServiceType doesn't Match. Want: %d, Got: %d", ServiceTypeCloseSecureChannelResponse, c.ServiceType())
		}
		t.Log(c.String())
	})
	t.Run("close-session-req", func(t *testing.T) {
		t.Parallel()
		c, err := Decode(testServiceBytes[8])
		if err != nil {
			t.Fatalf("Failed to decode Service: %s", err)
		}

		csr, ok := c.(*CloseSessionRequest)
		if !ok {
			t.Fatalf("Failed to assert type.")
		}

		switch {
		case c.ServiceType() != ServiceTypeCloseSessionRequest:
			t.Errorf("ServiceType doesn't Match. Want: %d, Got: %d", ServiceTypeCloseSessionRequest, c.ServiceType())
		case csr.DeleteSubscriptions.String() != "TRUE":
			t.Errorf("DeleteSubscriptions doesn't Match. Want: %s, Got: %s", "TRUE", csr.DeleteSubscriptions.String())
		}
		t.Log(c.String())
	})
	t.Run("close-session-res", func(t *testing.T) {
		t.Parallel()
		c, err := Decode(testServiceBytes[9])
		if err != nil {
			t.Fatalf("Failed to decode Service: %s", err)
		}

		_, ok := c.(*CloseSessionResponse)
		if !ok {
			t.Fatalf("Failed to assert type.")
		}

		switch {
		case c.ServiceType() != ServiceTypeCloseSessionResponse:
			t.Errorf("ServiceType doesn't Match. Want: %d, Got: %d", ServiceTypeCloseSessionResponse, c.ServiceType())
		}
		t.Log(c.String())
	})
}

func TestSerializeServices(t *testing.T) {
	t.Run("open-sec-chan-req", func(t *testing.T) {
		t.Parallel()
		o := NewOpenSecureChannelRequest(
			time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC),
			0, 1, 0, 0, "",
			0, ReqTypeIssue, SecModeNone, 6000000, nil,
		)
		o.SetDiagAll()

		serialized, err := o.Serialize()
		if err != nil {
			t.Fatalf("Failed to serialize Service: %s", err)
		}

		for i, s := range serialized {
			x := testServiceBytes[0][i]
			if s != x {
				t.Errorf("Bytes doesn't match. Want: %#x, Got: %#x at %dth", x, s, i)
			}
		}
		t.Logf("%x", serialized)
	})
	t.Run("open-sec-chan-res", func(t *testing.T) {
		t.Parallel()
		o := NewOpenSecureChannelResponse(
			time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC),
			1,
			0x00000000,
			NewNullDiagnosticInfo(),
			[]string{"foo", "bar"},
			0,
			NewChannelSecurityToken(
				1, 2, time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC), 6000000,
			),
			[]byte{0xff},
		)

		serialized, err := o.Serialize()
		if err != nil {
			t.Fatalf("Failed to serialize Service: %s", err)
		}

		for i, s := range serialized {
			x := testServiceBytes[1][i]
			if s != x {
				t.Errorf("Bytes doesn't match. Want: %#x, Got: %#x at %dth", x, s, i)
			}
		}
		t.Logf("%x", serialized)
	})
	t.Run("get-endpoint-req", func(t *testing.T) {
		t.Parallel()
		g := NewGetEndpointsRequest(
			time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC),
			1, 0, 0, "",
			"opc.tcp://wow.its.easy:11111/UA/Server",
			nil, nil,
		)
		g.SetDiagAll()

		serialized, err := g.Serialize()
		if err != nil {
			t.Fatalf("Failed to serialize Service: %s", err)
		}

		for i, s := range serialized {
			x := testServiceBytes[2][i]
			if s != x {
				t.Errorf("Bytes doesn't match. Want: %#x, Got: %#x at %dth", x, s, i)
			}
		}
		t.Logf("%x", serialized)
	})
	t.Run("get-endpoint-res", func(t *testing.T) {
		t.Parallel()
		g := NewGetEndpointsResponse(
			time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC),
			1, 0x00000000,
			NewNullDiagnosticInfo(),
			[]string{},
			NewEndpointDesctiption(
				"ep-url",
				NewApplicationDescription(
					"app-uri", "prod-uri", "app-name", AppTypeServer,
					"gw-uri", "prof-uri", []string{"discov-uri-1", "discov-uri-2"},
				),
				[]byte{},
				SecModeNone,
				"sec-uri",
				NewUserTokenPolicyArray(
					[]*UserTokenPolicy{
						NewUserTokenPolicy(
							"1", UserTokenAnonymous,
							"issued-token", "issuer-uri", "sec-uri",
						),
						NewUserTokenPolicy(
							"1", UserTokenAnonymous,
							"issued-token", "issuer-uri", "sec-uri",
						),
					},
				),
				"trans-uri",
				0,
			),
			NewEndpointDesctiption(
				"ep-url",
				NewApplicationDescription(
					"app-uri", "prod-uri", "app-name", AppTypeServer,
					"gw-uri", "prof-uri", []string{"discov-uri-1", "discov-uri-2"},
				),
				[]byte{},
				SecModeNone,
				"sec-uri",
				NewUserTokenPolicyArray(
					[]*UserTokenPolicy{
						NewUserTokenPolicy(
							"1", UserTokenAnonymous,
							"issued-token", "issuer-uri", "sec-uri",
						),
						NewUserTokenPolicy(
							"1", UserTokenAnonymous,
							"issued-token", "issuer-uri", "sec-uri",
						),
					},
				),
				"trans-uri",
				0,
			),
		)

		serialized, err := g.Serialize()
		if err != nil {
			t.Fatalf("Failed to serialize Service: %s", err)
		}

		for i, s := range serialized {
			x := testServiceBytes[3][i]
			if s != x {
				t.Errorf("Bytes doesn't match. Want: %#x, Got: %#x at %dth", x, s, i)
			}
		}
		t.Logf("%x", serialized)
	})
	t.Run("create-session-req", func(t *testing.T) {
		t.Parallel()
		c := NewCreateSessionRequest(
			time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC),
			"app-uri", "prod-uri", "app-name", AppTypeClient,
			"server-uri", "endpoint-url", "session-name",
			nil, nil, 6000000, 65534,
		)
		c.SetDiagAll()

		serialized, err := c.Serialize()
		if err != nil {
			t.Fatalf("Failed to serialize Service: %s, %T", err, err)
		}

		for i, s := range serialized {
			x := testServiceBytes[4][i]
			if s != x {
				t.Errorf("Bytes doesn't match. Want: %#x, Got: %#x at %dth", x, s, i)
			}
		}
		t.Logf("%x", serialized)
	})
	t.Run("create-session-res", func(t *testing.T) {
		t.Parallel()
		c := NewCreateSessionResponse(
			time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC),
			0, NewNullDiagnosticInfo(), 1, []byte{
				0x08, 0x22, 0x87, 0x62, 0xba, 0x81, 0xe1, 0x11,
				0xa6, 0x43, 0xf8, 0x77, 0x7b, 0xc6, 0x2f, 0xc8,
			},
			6000000, nil, nil, "", nil, 65534,
			NewEndpointDesctiption(
				"ep-url",
				NewApplicationDescription(
					"app-uri", "prod-uri", "app-name", AppTypeServer,
					"gw-uri", "prof-uri", []string{"discov-uri-1", "discov-uri-2"},
				),
				[]byte{},
				SecModeNone,
				"sec-uri",
				NewUserTokenPolicyArray(
					[]*UserTokenPolicy{
						NewUserTokenPolicy(
							"1", UserTokenAnonymous,
							"issued-token", "issuer-uri", "sec-uri",
						),
						NewUserTokenPolicy(
							"1", UserTokenAnonymous,
							"issued-token", "issuer-uri", "sec-uri",
						),
					},
				),
				"trans-uri",
				0,
			),
			NewEndpointDesctiption(
				"ep-url",
				NewApplicationDescription(
					"app-uri", "prod-uri", "app-name", AppTypeServer,
					"gw-uri", "prof-uri", []string{"discov-uri-1", "discov-uri-2"},
				),
				[]byte{},
				SecModeNone,
				"sec-uri",
				NewUserTokenPolicyArray(
					[]*UserTokenPolicy{
						NewUserTokenPolicy(
							"1", UserTokenAnonymous,
							"issued-token", "issuer-uri", "sec-uri",
						),
						NewUserTokenPolicy(
							"1", UserTokenAnonymous,
							"issued-token", "issuer-uri", "sec-uri",
						),
					},
				),
				"trans-uri",
				0,
			),
		)

		serialized, err := c.Serialize()
		if err != nil {
			t.Fatalf("Failed to serialize Service: %s, %T", err, err)
		}

		for i, s := range serialized {
			x := testServiceBytes[5][i]
			if s != x {
				t.Errorf("Bytes doesn't match. Want: %#x, Got: %#x at %dth", x, s, i)
			}
		}
		t.Logf("%x", serialized)
	})
	t.Run("close-sec-chan-req", func(t *testing.T) {
		t.Parallel()
		o := NewCloseSecureChannelRequest(
			time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC),
			0, 1, 0, 0, "", 1,
		)
		o.SetDiagAll()

		serialized, err := o.Serialize()
		if err != nil {
			t.Fatalf("Failed to serialize Service: %s", err)
		}

		for i, s := range serialized {
			x := testServiceBytes[6][i]
			if s != x {
				t.Errorf("Bytes doesn't match. Want: %#x, Got: %#x at %dth", x, s, i)
			}
		}
		t.Logf("%x", serialized)
	})
	t.Run("close-sec-chan-res", func(t *testing.T) {
		t.Parallel()
		o := NewCloseSecureChannelResponse(
			time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC),
			1,
			0x00000000,
			NewNullDiagnosticInfo(),
			[]string{"foo", "bar"},
		)

		serialized, err := o.Serialize()
		if err != nil {
			t.Fatalf("Failed to serialize Service: %s", err)
		}

		for i, s := range serialized {
			x := testServiceBytes[7][i]
			if s != x {
				t.Errorf("Bytes doesn't match. Want: %#x, Got: %#x at %dth", x, s, i)
			}
		}
		t.Logf("%x", serialized)
	})
	/*
		t.Run("close-session-req", func(t *testing.T) {
			t.Parallel()
			o := NewCloseSessionRequest(
				time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC),
				[]byte{
					0x08, 0x22, 0x87, 0x62, 0xba, 0x81, 0xe1, 0x11,
					0xa6, 0x43, 0xf8, 0x77, 0x7b, 0xc6, 0x2f, 0xc8,
				}, 1, 0, 0, "", true,
			)
			o.SetDiagAll()

			serialized, err := o.Serialize()
			if err != nil {
				t.Fatalf("Failed to serialize Service: %s", err)
			}

			for i, s := range serialized {
				x := testServiceBytes[8][i]
				if s != x {
					t.Errorf("Bytes doesn't match. Want: %#x, Got: %#x at %dth", x, s, i)
				}
			}
			t.Logf("%x", serialized)
		})
	*/
	t.Run("close-session-res", func(t *testing.T) {
		t.Parallel()
		o := NewCloseSessionResponse(
			time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC),
			1,
			0x00000000,
			NewNullDiagnosticInfo(),
			[]string{"foo", "bar"},
		)

		serialized, err := o.Serialize()
		if err != nil {
			t.Fatalf("Failed to serialize Service: %s", err)
		}

		for i, s := range serialized {
			x := testServiceBytes[9][i]
			if s != x {
				t.Errorf("Bytes doesn't match. Want: %#x, Got: %#x at %dth", x, s, i)
			}
		}
		t.Logf("%x", serialized)
	})
}