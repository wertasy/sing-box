package option

type VLESSInboundOptions struct {
	ListenOptions

	// VLESS users.
	//
	// @Example:
	//    [
	//      {
	//        "name": "sekai",
	//        "uuid": "bf000d23-0752-40b4-affe-68f7707a9661",
	//        "flow": ""
	//      }
	//    ]
	//
	// @Required
	Users []VLESSUser `json:"users,omitempty"`
	InboundTLSOptionsContainer

	// See [Multiplex](/configuration/shared/multiplex#inbound) for details.
	Multiplex *InboundMultiplexOptions `json:"multiplex,omitempty"`

	// V2Ray Transport configuration, see [V2Ray Transport](/configuration/shared/v2ray-transport/).
	Transport *V2RayTransportOptions `json:"transport,omitempty"`
}

type VLESSUser struct {
	Name string `json:"name,omitempty"`

	// VLESS user id.
	//
	// @Required
	UUID string `json:"uuid"`

	// VLESS Sub-protocol.
	//
	// Available values:
	//
	// * `xtls-rprx-vision`
	Flow string `json:"flow,omitempty"`
}

type VLESSOutboundOptions struct {
	DialerOptions

	// @Override server:
	//    @Example: "127.0.0.1"
	//
	// @Override server_port:
	//    @Example: 1080
	ServerOptions

	// VLESS user id.
	//
	// @Example: "bf000d23-0752-40b4-affe-68f7707a9661"
	//
	// @Required
	UUID string `json:"uuid"`

	// VLESS Sub-protocol.
	//
	// Available values:
	//
	// * `xtls-rprx-vision`
	//
	// @Example: "xtls-rprx-vision"
	Flow string `json:"flow,omitempty"`

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

	// UDP packet encoding, xudp will be used by default.
	//
	// | Encoding   | Description           |
	// |------------|-----------------------|
	// | (none)     | Disabled              |
	// | packetaddr | Supported by v2ray 5+ |
	// | xudp       | Supported by xray     |
	//
	// @Example: ""
	PacketEncoding *string `json:"packet_encoding,omitempty"`
}
