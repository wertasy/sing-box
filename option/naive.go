package option

import (
	"github.com/sagernet/sing/common/auth"
	"github.com/sagernet/sing/common/byteformats"
	"github.com/sagernet/sing/common/json/badoption"
)

type NaiveInboundOptions struct {
	ListenOptions
	// @required
	//
	// Naive users.
	Users []auth.User `json:"users,omitempty"`
	// Listen network, one of `tcp` `udp`.
	//
	// Both if empty.
	Network NetworkList `json:"network,omitempty"`
	// @since 1.13.0
	//
	// QUIC congestion control algorithm.
	//
	// | Algorithm      | Description                     |
	// |----------------|---------------------------------|
	// | `bbr`          | BBR                             |
	// | `cubic`        | CUBIC                           |
	// | `reno`         | New Reno                        |
	//
	// `bbr` is used by default.
	QUICCongestionControl string `json:"quic_congestion_control,omitempty" enum:"bbr,cubic,reno"`
	InboundTLSOptionsContainer
}

type NaiveOutboundOptions struct {
	DialerOptions
	ServerOptions
	// Authentication username.
	Username string `json:"username,omitempty"`
	// Authentication password.
	Password string `json:"password,omitempty"`
	// Number of concurrent tunnel connections. Multiple connections make the tunneling easier to detect through traffic analysis, which defeats the purpose of NaiveProxy's design to resist traffic analysis.
	InsecureConcurrency int `json:"insecure_concurrency,omitempty"`
	// Extra headers to send in HTTP requests.
	ExtraHeaders badoption.HTTPHeader `json:"extra_headers,omitempty"`
	// The flow control window.
	//
	// When `quic` is enabled, it sets the initial QUIC stream receive window, and `6 MB` is used by default.
	//
	// Otherwise, it sets the HTTP/2 session receive window, the stream receive window is set to half of it, and `4 MB` is used by default on iOS, `128 MB` on other platforms.
	ReceiveWindow *byteformats.MemoryBytes `json:"stream_receive_window,omitempty"`
	// UDP over TCP protocol settings.
	//
	// See [UDP Over TCP](/configuration/shared/udp-over-tcp/) for details.
	UDPOverTCP *UDPOverTCPOptions `json:"udp_over_tcp,omitempty"`
	// Use QUIC instead of HTTP/2.
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
	// `cubic` is used by default (the default of Chromium, which NaiveProxy is based on).
	QUICCongestionControl string `json:"quic_congestion_control,omitempty" enum:"bbr,bbr2,cubic,reno"`
	// @note:
	//
	//     Only used when `quic` is enabled.
	//
	// The initial QUIC session receive window.
	//
	// `15 MB` is used by default.
	QUICSessionReceiveWindow *byteformats.MemoryBytes `json:"quic_session_receive_window,omitempty"`
	OutboundTLSOptionsContainer
}
