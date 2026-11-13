package option

type VMessInboundOptions struct {
	ListenOptions

	// VMess users.
	//
	// | Alter ID | Description             |
	// |----------|-------------------------|
	// | 0        | Disable legacy protocol |
	// | > 0      | Enable legacy protocol  |
	//
	// @Example:
	//    [
	//      {
	//        "name": "sekai",
	//        "uuid": "bf000d23-0752-40b4-affe-68f7707a9661",
	//        "alterId": 0
	//      }
	//    ]
	//
	// @Required
	//
	// @Warning:
	//
	//     Legacy protocol support (VMess MD5 Authentication) is provided for compatibility purposes only, use of alterId > 1 is not recommended.
	Users []VMessUser `json:"users,omitempty"`
	InboundTLSOptionsContainer

	// See [Multiplex](/configuration/shared/multiplex#inbound) for details.
	Multiplex *InboundMultiplexOptions `json:"multiplex,omitempty"`

	// V2Ray Transport configuration, see [V2Ray Transport](/configuration/shared/v2ray-transport/).
	Transport *V2RayTransportOptions `json:"transport,omitempty"`
}

type VMessUser struct {
	Name    string `json:"name,omitempty"`
	UUID    string `json:"uuid"`
	AlterId int    `json:"alterId,omitempty"`
}

type VMessOutboundOptions struct {
	DialerOptions

	// @Override server:
	//    @Example: "127.0.0.1"
	//
	// @Override server_port:
	//    @Example: 1080
	ServerOptions

	// The VMess user id.
	//
	// @Example: "bf000d23-0752-40b4-affe-68f7707a9661"
	//
	// @Required
	UUID string `json:"uuid"`

	// Encryption methods:
	//
	// * `auto`
	// * `none`
	// * `zero`
	// * `aes-128-gcm`
	// * `chacha20-poly1305`
	//
	// Legacy encryption methods:
	//
	// * `aes-128-ctr`
	//
	// @Example: "auto"
	Security string `json:"security,omitempty" enum:"auto,none,zero,aes-128-cfb,aes-128-gcm,chacha20-poly1305"`

	// | Alter ID | Description         |
	// |----------|---------------------|
	// | 0        | Use AEAD protocol   |
	// | 1        | Use legacy protocol |
	// | > 1      | Unused, same as 1   |
	//
	// @Example: 0
	AlterId int `json:"alter_id,omitempty"`

	// Protocol parameter. Will waste traffic randomly if enabled (enabled by default in v2ray and cannot be disabled).
	//
	// @Example: false
	GlobalPadding bool `json:"global_padding,omitempty"`

	// Protocol parameter. Enable length block encryption.
	//
	// @Example: true
	AuthenticatedLength bool `json:"authenticated_length,omitempty"`

	// Enabled network
	//
	// One of `tcp` `udp`.
	//
	// Both is enabled by default.
	//
	// @Example: "tcp"
	Network NetworkList `json:"network,omitempty"`
	OutboundTLSOptionsContainer

	// UDP packet encoding.
	//
	// | Encoding   | Description           |
	// |------------|-----------------------|
	// | (none)     | Disabled              |
	// | packetaddr | Supported by v2ray 5+ |
	// | xudp       | Supported by xray     |
	//
	// @Example: ""
	PacketEncoding string `json:"packet_encoding,omitempty" enum:"packetaddr,xudp"`

	// See [Multiplex](/configuration/shared/multiplex#outbound) for details.
	//
	// @Example: {}
	Multiplex *OutboundMultiplexOptions `json:"multiplex,omitempty"`

	// V2Ray Transport configuration, see [V2Ray Transport](/configuration/shared/v2ray-transport/).
	//
	// @Example: {}
	Transport *V2RayTransportOptions `json:"transport,omitempty"`
}
