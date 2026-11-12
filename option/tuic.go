package option

import "github.com/sagernet/sing/common/json/badoption"

type TUICInboundOptions struct {
	ListenOptions
	// TUIC users
	//
	// @Example:
	//    [
	//      {
	//        "name": "sekai",
	//        "uuid": "059032A9-7D40-4A96-9BB1-36823D848068",
	//        "password": "hello"
	//      }
	//    ]
	Users []TUICUser `json:"users,omitempty"`
	// QUIC congestion control algorithm
	//
	// One of: `cubic`, `new_reno`, `bbr`
	//
	// @Default: cubic
	//
	// @Example: "cubic"
	CongestionControl string `json:"congestion_control,omitempty" enum:"cubic,new_reno,bbr"`
	// How long the server should wait for the client to send the authentication command
	//
	// @Default: 3s
	//
	// @Example: "3s"
	AuthTimeout badoption.Duration `json:"auth_timeout,omitempty"`
	// Enable 0-RTT QUIC connection handshake on the client side
	// This is not impacting much on the performance, as the protocol is fully multiplexed
	//
	// @Warning:
	//     Disabling this is highly recommended, as it is vulnerable to replay attacks.
	//     See [Attack of the clones](https://blog.cloudflare.com/even-faster-connection-establishment-with-quic-0-rtt-resumption/#attack-of-the-clones)
	ZeroRTTHandshake bool `json:"zero_rtt_handshake,omitempty"`
	// Interval for sending heartbeat packets for keeping the connection alive
	//
	// @Default: 10s
	//
	// @Example: "10s"
	Heartbeat badoption.Duration `json:"heartbeat,omitempty"`
	// @Override:
	//    @Required
	//
	//    TLS configuration, see [TLS](/configuration/shared/tls/#inbound).
	InboundTLSOptionsContainer
	QUICOptions
}

type TUICUser struct {
	Name string `json:"name,omitempty"`
	// TUIC user uuid
	//
	// @Required
	UUID string `json:"uuid,omitempty"`
	// TUIC user password
	Password string `json:"password,omitempty"`
}

type TUICOutboundOptions struct {
	DialerOptions
	// @Override server:
	//    @Example: "127.0.0.1"
	//
	// @Override server_port:
	//    @Example: 1080
	ServerOptions
	// TUIC user uuid
	//
	// @Example: "2DD61D93-75D8-4DA4-AC0E-6AECE7EAC365"
	//
	// @Required
	UUID string `json:"uuid,omitempty"`
	// TUIC user password
	//
	// @Example: "hello"
	Password string `json:"password,omitempty"`
	// QUIC congestion control algorithm
	//
	// One of: `cubic`, `new_reno`, `bbr`
	//
	// @Default: cubic
	//
	// @Example: "cubic"
	CongestionControl string `json:"congestion_control,omitempty" enum:"cubic,new_reno,bbr"`
	// UDP packet relay mode
	//
	// | Mode   | Description                                                              |
	// |:-------|:-------------------------------------------------------------------------|
	// | native | native UDP characteristics                                               |
	// | quic   | lossless UDP relay using QUIC streams, additional overhead is introduced |
	//
	// Conflict with `udp_over_stream`.
	//
	// @Default: native
	//
	// @Example: "native"
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
	// @Warning:
	//
	//     Disabling this feature is highly recommended, as it is vulnerable to replay attacks.
	//     See [Attack of the clones](https://blog.cloudflare.com/even-faster-connection-establishment-with-quic-0-rtt-resumption/#attack-of-the-clones)
	ZeroRTTHandshake bool `json:"zero_rtt_handshake,omitempty"`
	// Heartbeat interval for keeping the connection alive
	//
	// @Example: "10s"
	Heartbeat badoption.Duration `json:"heartbeat,omitempty"`
	// Enabled network
	//
	// One of `tcp` `udp`.
	//
	// Both is enabled by default.
	//
	// @Example: "tcp"
	Network NetworkList `json:"network,omitempty"`
	// @Override:
	//    @Required
	//
	//    TLS configuration, see [TLS](/configuration/shared/tls/#outbound).
	OutboundTLSOptionsContainer
	QUICOptions
}
