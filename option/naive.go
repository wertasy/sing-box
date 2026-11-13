package option

import (
	"github.com/sagernet/sing/common/auth"
	"github.com/sagernet/sing/common/byteformats"
	"github.com/sagernet/sing/common/json/badoption"
)

type NaiveInboundOptions struct {
	ListenOptions

	// Naive users.
	//
	// @Example:
	//    [
	//      {
	//        "username": "sekai",
	//        "password": "password"
	//      }
	//    ]
	//
	// @Required
	Users []auth.User `json:"users,omitempty"`

	// Listen network, one of `tcp` `udp`.
	//
	// Both if empty.
	//
	// @Example: "udp"
	Network NetworkList `json:"network,omitempty"`

	// QUIC congestion control algorithm.
	//
	// | Algorithm      | Description                     |
	// |----------------|---------------------------------|
	// | `bbr`          | BBR                             |
	// | `cubic`        | CUBIC                           |
	// | `reno`         | New Reno                        |
	//
	// @Default: bbr
	//
	// @Since 1.13.0
	QUICCongestionControl string `json:"quic_congestion_control,omitempty" enum:"bbr,cubic,reno"`
	InboundTLSOptionsContainer
}

type NaiveOutboundOptions struct {
	DialerOptions

	// @Override server:
	//    @Example: "127.0.0.1"
	//
	// @Override server_port:
	//    @Example: 443
	ServerOptions

	// Authentication username.
	//
	// @Example: "sekai"
	Username string `json:"username,omitempty"`

	// Authentication password.
	//
	// @Example: "password"
	Password string `json:"password,omitempty"`

	// Number of concurrent tunnel connections. Multiple connections make the tunneling easier to detect through traffic analysis, which defeats the purpose of NaiveProxy's design to resist traffic analysis.
	//
	// @Example: 0
	InsecureConcurrency int `json:"insecure_concurrency,omitempty"`

	// Extra headers to send in HTTP requests.
	//
	// @Example: {}
	ExtraHeaders badoption.HTTPHeader `json:"extra_headers,omitempty"`

	// The flow control window.
	//
	// When `quic` is enabled, it sets the initial QUIC stream receive window, and `6 MB` will be used by default.
	//
	// Otherwise, it sets the HTTP/2 session receive window, the stream receive window is set to half of it, and `4 MB` will be used by default on iOS, `128 MB` on other platforms.
	//
	// @Example: ""
	ReceiveWindow *byteformats.MemoryBytes `json:"stream_receive_window,omitempty"`

	// UDP over TCP protocol settings.
	//
	// See [UDP Over TCP](/configuration/shared/udp-over-tcp/) for details.
	//
	// @Example: false | {}
	UDPOverTCP *UDPOverTCPOptions `json:"udp_over_tcp,omitempty"`

	// Use QUIC instead of HTTP/2.
	//
	// @Example: false
	QUIC bool `json:"quic,omitempty"`

	// QUIC congestion control algorithm.
	//
	// | Algorithm | Description |
	// |-----------|-------------|
	// | `bbr` | BBR |
	// | `bbr2` | BBRv2 |
	// | `cubic` | CUBIC |
	// | `reno` | New Reno |
	//
	// `cubic` will be used by default (the default of Chromium, which NaiveProxy is based on).
	//
	// @Example: ""
	QUICCongestionControl string `json:"quic_congestion_control,omitempty" enum:"bbr,bbr2,cubic,reno"`

	// The initial QUIC session receive window.
	//
	// @Default: 15 MB
	//
	// @Example: ""
	//
	// @Note:
	//
	//     Only used when `quic` is enabled.
	QUICSessionReceiveWindow *byteformats.MemoryBytes `json:"quic_session_receive_window,omitempty"`

	// @Override:
	//    @Required
	//
	//    TLS configuration, see [TLS](/configuration/shared/tls/#outbound).
	//
	//    Only `server_name`, `certificate`, `certificate_path` and `ech` are supported.
	//
	//    Self-signed certificates change traffic behavior significantly, which defeats the purpose of NaiveProxy's design to resist traffic analysis, and should not be used in production.
	OutboundTLSOptionsContainer
}
