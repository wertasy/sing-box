package option

type TrojanInboundOptions struct {
	ListenOptions
	// @required
	//
	// Trojan users.
	Users []TrojanUser `json:"users,omitempty"`
	InboundTLSOptionsContainer
	// @failure:
	//
	//     There is no evidence that GFW detects and blocks Trojan servers based on HTTP responses, and opening the standard http/s port on the server is a much bigger signature.
	//
	// Fallback server configuration. Disabled if `fallback` and `fallback_for_alpn` are empty.
	Fallback *ServerOptions `json:"fallback,omitempty"`
	// Fallback server configuration for specified ALPN.
	//
	// If not empty, TLS fallback requests with ALPN not in this table will be rejected.
	FallbackForALPN map[string]*ServerOptions `json:"fallback_for_alpn,omitempty"`
	// See [Multiplex](/configuration/shared/multiplex#inbound) for details.
	Multiplex *InboundMultiplexOptions `json:"multiplex,omitempty"`
	// V2Ray Transport configuration, see [V2Ray Transport](/configuration/shared/v2ray-transport/).
	Transport *V2RayTransportOptions `json:"transport,omitempty"`
}

type TrojanUser struct {
	// User name.
	Name string `json:"name,omitempty"`
	// @required
	//
	// The Trojan password.
	Password string `json:"password"`
}

type TrojanOutboundOptions struct {
	DialerOptions
	ServerOptions
	// @required
	//
	// The Trojan password.
	Password string `json:"password"`
	// Enabled network
	//
	// One of `tcp` `udp`.
	//
	// Both is enabled by default.
	Network NetworkList `json:"network,omitempty"`
	OutboundTLSOptionsContainer
	// See [Multiplex](/configuration/shared/multiplex#outbound) for details.
	Multiplex *OutboundMultiplexOptions `json:"multiplex,omitempty"`
	// V2Ray Transport configuration, see [V2Ray Transport](/configuration/shared/v2ray-transport/).
	Transport *V2RayTransportOptions `json:"transport,omitempty"`
}
