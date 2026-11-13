package option

type TrojanInboundOptions struct {
	ListenOptions

	// Trojan users.
	//
	// @Example:
	//    [
	//      {
	//        "name": "sekai",
	//        "password": "8JCsPssfgS8tiRwiMlhARg=="
	//      }
	//    ]
	//
	// @Required
	Users []TrojanUser `json:"users,omitempty"`
	InboundTLSOptionsContainer

	// Fallback server configuration. Disabled if `fallback` and `fallback_for_alpn` are empty.
	//
	// @Failure:
	//
	//     There is no evidence that GFW detects and blocks Trojan servers based on HTTP responses, and opening the standard http/s port on the server is a much bigger signature.
	//
	// @Example:
	//    {
	//      "server": "127.0.0.1",
	//      "server_port": 8080
	//    }
	Fallback *ServerOptions `json:"fallback,omitempty"`

	// Fallback server configuration for specified ALPN.
	//
	// If not empty, TLS fallback requests with ALPN not in this table will be rejected.
	//
	// @Example:
	//    {
	//      "http/1.1": {
	//        "server": "127.0.0.1",
	//        "server_port": 8081
	//      }
	//    }
	FallbackForALPN map[string]*ServerOptions `json:"fallback_for_alpn,omitempty"`

	// See [Multiplex](/configuration/shared/multiplex#inbound) for details.
	Multiplex *InboundMultiplexOptions `json:"multiplex,omitempty"`

	// V2Ray Transport configuration, see [V2Ray Transport](/configuration/shared/v2ray-transport/).
	Transport *V2RayTransportOptions `json:"transport,omitempty"`
}

type TrojanUser struct {
	// User name.
	Name string `json:"name,omitempty"`

	// The Trojan password.
	//
	// @Required
	Password string `json:"password"`
}

type TrojanOutboundOptions struct {
	DialerOptions

	// @Override server:
	//    @Example: "127.0.0.1"
	//
	// @Override server_port:
	//    @Example: 1080
	ServerOptions

	// The Trojan password.
	//
	// @Example: "8JCsPssfgS8tiRwiMlhARg=="
	//
	// @Required
	Password string `json:"password"`

	// Enabled network
	//
	// One of `tcp` `udp`.
	//
	// Both is enabled by default.
	//
	// @Example: "tcp"
	Network NetworkList `json:"network,omitempty"`
	OutboundTLSOptionsContainer

	// See [Multiplex](/configuration/shared/multiplex#outbound) for details.
	//
	// @Example: {}
	Multiplex *OutboundMultiplexOptions `json:"multiplex,omitempty"`

	// V2Ray Transport configuration, see [V2Ray Transport](/configuration/shared/v2ray-transport/).
	//
	// @Example: {}
	Transport *V2RayTransportOptions `json:"transport,omitempty"`
}
