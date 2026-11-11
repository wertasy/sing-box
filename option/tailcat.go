package option

import (
	"reflect"

	"github.com/sagernet/sing-box/schema"
	"github.com/sagernet/sing/common/json"
)

type TailcatDERPOptions struct {
	// URL of the [DERP map](https://pkg.go.dev/tailscale.com/tailcfg#DERPMap).
	//
	// `https://tailcat.dev/derpmap.json` is used by default.
	DERPMapURL string `json:"derp_map_url,omitempty"`
	// DERP region ID in the DERP map.
	//
	// Conflicts with `derp_servers`.
	DERPRegion int `json:"derp_region,omitempty"`
	// Custom DERP servers, see [derp_servers](/configuration/inbound/tailcat/#derp_servers) in Tailcat inbound.
	//
	// Conflicts with `derp_map_url` and `derp_region`.
	DERPServers []TailcatDERPServer `json:"derp_servers,omitempty"`
	// HTTP client used to fetch the DERP map.
	//
	// See [HTTP Client](/configuration/shared/http-client/) for details.
	HTTPClient *HTTPClientOptions `json:"http_client,omitempty"`
}

type _TailcatDERPServer struct {
	Host     string `json:"host"`
	IPv4     string `json:"ipv4,omitempty"`
	IPv6     string `json:"ipv6,omitempty"`
	DERPPort uint16 `json:"derp_port,omitempty"`
	STUNPort uint16 `json:"stun_port,omitempty"`
	CertName string `json:"cert_name,omitempty"`
}

type TailcatDERPServer _TailcatDERPServer

func (s TailcatDERPServer) MarshalJSON() ([]byte, error) {
	if s == (TailcatDERPServer{Host: s.Host}) {
		return json.Marshal(s.Host)
	}
	return json.Marshal(_TailcatDERPServer(s))
}

func (s *TailcatDERPServer) UnmarshalJSON(bytes []byte) error {
	var host string
	err := json.Unmarshal(bytes, &host)
	if err == nil {
		*s = TailcatDERPServer{Host: host}
		return nil
	}
	return json.UnmarshalDisallowUnknownFields(bytes, (*_TailcatDERPServer)(s))
}

func (s TailcatDERPServer) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	objectForm := schema.StrictObject()
	err := builder.FlattenStruct(objectForm, reflect.TypeFor[TailcatDERPServer]())
	if err != nil {
		return nil, err
	}
	return schema.AnyOf(schema.StringNode(), objectForm), nil
}

type TailcatUser struct {
	Name string `json:"name,omitempty"`
	// @required
	//
	// Client public key.
	PublicKey string `json:"public_key"`
}

type TailcatInboundOptions struct {
	// @required
	//
	// Private key.
	//
	// Generate with `sing-box generate tailcat-keypair`.
	PrivateKey string `json:"private_key"`
	// Pre-shared key.
	PreSharedKey string `json:"pre_shared_key,omitempty"`
	// Tailcat users.
	//
	// Clients are not verified if empty.
	//
	// To use a DERP server with client verification, set `verify_client_inbound` or `verify_client_key` in
	// [DERP service](/configuration/service/derp/#verify_client_inbound), and clients must use fixed
	// private keys.
	Users []TailcatUser `json:"users,omitempty"`
	TailcatDERPOptions
	DialerOptions
}

type TailcatOutboundOptions struct {
	// Private key.
	//
	// A random key is used by default.
	//
	// Required when the server verifies clients with `users`, or the DERP server verifies clients with
	// `verify_client_inbound` or `verify_client_key`.
	PrivateKey string `json:"private_key,omitempty"`
	// @required
	//
	// Server public key.
	ServerPublicKey string `json:"server_public_key"`
	// @required
	//
	// Server disco public key.
	ServerDiscoKey string `json:"server_disco_key"`
	// Pre-shared key.
	PreSharedKey string `json:"pre_shared_key,omitempty"`
	TailcatDERPOptions
	DialerOptions
}
