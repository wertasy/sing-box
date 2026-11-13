package option

// @Tag ss-in
type ShadowsocksInboundOptions struct {
	ListenOptions

	// Listen network, one of `tcp` `udp`.
	//
	// Both if empty.
	Network NetworkList `json:"network,omitempty"`

	// | Method                        | Key Length |
	// |-------------------------------|------------|
	// | 2022-blake3-aes-128-gcm       | 16         |
	// | 2022-blake3-aes-256-gcm       | 32         |
	// | 2022-blake3-chacha20-poly1305 | 32         |
	// | none                          | /          |
	// | aes-128-gcm                   | /          |
	// | aes-192-gcm                   | /          |
	// | aes-256-gcm                   | /          |
	// | chacha20-ietf-poly1305        | /          |
	// | xchacha20-ietf-poly1305       | /          |
	//
	// @Example: "2022-blake3-aes-128-gcm"
	//
	// @Required
	Method string `json:"method" enum:"none,aes-128-gcm,aes-192-gcm,aes-256-gcm,chacha20-ietf-poly1305,xchacha20-ietf-poly1305,2022-blake3-aes-128-gcm,2022-blake3-aes-256-gcm,2022-blake3-chacha20-poly1305"`

	// | Method        | Password Format                                |
	// |---------------|------------------------------------------------|
	// | none          | /                                              |
	// | 2022 methods  | `sing-box generate rand --base64 <Key Length>` |
	// | other methods | any string                                     |
	//
	// @Example: "8JCsPssfgS8tiRwiMlhARg=="
	//
	// @Required
	Password string `json:"password,omitempty"`

	// @Omit
	Users []ShadowsocksUser `json:"users,omitempty"`

	// @Omit
	Destinations []ShadowsocksDestination `json:"destinations,omitempty"`

	// See [Multiplex](/configuration/shared/multiplex#inbound) for details.
	Multiplex *InboundMultiplexOptions `json:"multiplex,omitempty"`

	// Defaults to `false`. Enable this when the inbound is managed by the [SSM API](/configuration/service/ssm-api) for dynamic user.
	Managed bool `json:"managed,omitempty"`
}

type ShadowsocksUser struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password"`
}

type ShadowsocksDestination struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password"`
	ServerOptions
}

// @Tag ss-out
type ShadowsocksOutboundOptions struct {
	DialerOptions

	// @Override server:
	//    @Example: "127.0.0.1"
	//
	// @Override server_port:
	//    @Example: 1080
	ServerOptions

	// Encryption methods:
	//
	// * `2022-blake3-aes-128-gcm`
	// * `2022-blake3-aes-256-gcm`
	// * `2022-blake3-chacha20-poly1305`
	// * `none`
	// * `aes-128-gcm`
	// * `aes-192-gcm`
	// * `aes-256-gcm`
	// * `chacha20-ietf-poly1305`
	// * `xchacha20-ietf-poly1305`
	//
	// Legacy encryption methods:
	//
	// * `aes-128-ctr`
	// * `aes-192-ctr`
	// * `aes-256-ctr`
	// * `aes-128-cfb`
	// * `aes-192-cfb`
	// * `aes-256-cfb`
	// * `rc4-md5`
	// * `chacha20-ietf`
	// * `xchacha20`
	//
	// @Example: "2022-blake3-aes-128-gcm"
	//
	// @Required
	Method string `json:"method" enum:"none,aes-128-gcm,aes-192-gcm,aes-256-gcm,chacha20-ietf-poly1305,xchacha20-ietf-poly1305,2022-blake3-aes-128-gcm,2022-blake3-aes-256-gcm,2022-blake3-chacha20-poly1305,aes-128-ctr,aes-192-ctr,aes-256-ctr,aes-128-cfb,aes-192-cfb,aes-256-cfb,rc4-md5,chacha20-ietf,xchacha20"`

	// The shadowsocks password.
	//
	// @Example: "8JCsPssfgS8tiRwiMlhARg=="
	//
	// @Required
	Password string `json:"password,omitempty"`

	// Shadowsocks SIP003 plugin, implemented in internal.
	//
	// Only `obfs-local` and `v2ray-plugin` are supported.
	//
	// @Example: ""
	Plugin string `json:"plugin,omitempty"`

	// Shadowsocks SIP003 plugin options.
	//
	// @Example: ""
	PluginOptions string `json:"plugin_opts,omitempty"`

	// Enabled network
	//
	// One of `tcp` `udp`.
	//
	// Both is enabled by default.
	//
	// @Example: "udp"
	Network NetworkList `json:"network,omitempty"`

	// UDP over TCP configuration.
	//
	// See [UDP Over TCP](/configuration/shared/udp-over-tcp/) for details.
	//
	// Conflict with `multiplex`.
	//
	// @Example: false | {}
	UDPOverTCP *UDPOverTCPOptions `json:"udp_over_tcp,omitempty"`

	// See [Multiplex](/configuration/shared/multiplex#outbound) for details.
	//
	// @Example: {}
	Multiplex *OutboundMultiplexOptions `json:"multiplex,omitempty"`
}
