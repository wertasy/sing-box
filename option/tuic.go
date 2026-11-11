package option

import "github.com/sagernet/sing/common/json/badoption"

type TUICInboundOptions struct {
	ListenOptions
	// TUIC users
	Users []TUICUser `json:"users,omitempty"`
	// QUIC congestion control algorithm
	//
	// One of: `cubic`, `new_reno`, `bbr`
	//
	// `cubic` is used by default.
	CongestionControl string `json:"congestion_control,omitempty" enum:"cubic,new_reno,bbr"`
	// How long the server should wait for the client to send the authentication command
	//
	// `3s` is used by default.
	AuthTimeout badoption.Duration `json:"auth_timeout,omitempty"`
	// Enable 0-RTT QUIC connection handshake on the client side
	// This is not impacting much on the performance, as the protocol is fully multiplexed
	//
	// @warning:
	//     Disabling this is highly recommended, as it is vulnerable to replay attacks.
	//     See [Attack of the clones](https://blog.cloudflare.com/even-faster-connection-establishment-with-quic-0-rtt-resumption/#attack-of-the-clones)
	ZeroRTTHandshake bool `json:"zero_rtt_handshake,omitempty"`
	// Interval for sending heartbeat packets for keeping the connection alive
	//
	// `10s` is used by default.
	Heartbeat badoption.Duration `json:"heartbeat,omitempty"`
	InboundTLSOptionsContainer
	QUICOptions
}

type TUICUser struct {
	Name string `json:"name,omitempty"`
	// @required
	//
	// TUIC user uuid
	UUID string `json:"uuid,omitempty"`
	// TUIC user password
	Password string `json:"password,omitempty"`
}

type TUICOutboundOptions struct {
	DialerOptions
	ServerOptions
	// @required
	//
	// TUIC user uuid
	UUID string `json:"uuid,omitempty"`
	// TUIC user password
	Password string `json:"password,omitempty"`
	// QUIC congestion control algorithm
	//
	// One of: `cubic`, `new_reno`, `bbr`
	//
	// `cubic` is used by default.
	CongestionControl string `json:"congestion_control,omitempty" enum:"cubic,new_reno,bbr"`
	// UDP packet relay mode
	//
	// | Mode   | Description                                                              |
	// |:-------|:-------------------------------------------------------------------------|
	// | native | native UDP characteristics                                               |
	// | quic   | lossless UDP relay using QUIC streams, additional overhead is introduced |
	//
	// `native` is used by default.
	//
	// Conflict with `udp_over_stream`.
	UDPRelayMode string `json:"udp_relay_mode,omitempty" enum:"native,quic"`
	// This is the TUIC port of the [UDP over TCP protocol](/configuration/shared/udp-over-tcp/), designed to provide a QUIC
	// stream based UDP relay mode that TUIC does not provide. Since it is an add-on protocol, you will need to use sing-box or
	// another program compatible with the protocol as a server.
	//
	// This mode has no positive effect in a proper UDP proxy scenario and should only be applied to relay streaming UDP
	// traffic (basically QUIC streams).
	//
	// Conflict with `udp_relay_mode`.
	UDPOverStream bool `json:"udp_over_stream,omitempty"`
	// Enable 0-RTT QUIC connection handshake on the client.
	// This has little performance impact, because the protocol is fully multiplexed.
	//
	// @warning:
	//
	//     Disabling this feature is highly recommended, as it is vulnerable to replay attacks.
	//     See [Attack of the clones](https://blog.cloudflare.com/even-faster-connection-establishment-with-quic-0-rtt-resumption/#attack-of-the-clones)
	ZeroRTTHandshake bool `json:"zero_rtt_handshake,omitempty"`
	// Heartbeat interval for keeping the connection alive
	Heartbeat badoption.Duration `json:"heartbeat,omitempty"`
	// Enabled network
	//
	// One of `tcp` `udp`.
	//
	// Both is enabled by default.
	Network NetworkList `json:"network,omitempty"`
	OutboundTLSOptionsContainer
	QUICOptions
}
