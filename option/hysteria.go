package option

import (
	"github.com/sagernet/sing/common/byteformats"
	"github.com/sagernet/sing/common/json/badoption"
)

type HysteriaInboundOptions struct {
	ListenOptions
	// @Example: "100 Mbps"
	// @Omit
	Up *byteformats.NetworkBytesCompat `json:"up,omitempty"`
	// @Example: 100
	// @Omit
	UpMbps int `json:"up_mbps,omitempty"`
	// @Example: "100 Mbps"
	// @Omit
	Down *byteformats.NetworkBytesCompat `json:"down,omitempty"`
	// @Example: 100
	// @Omit
	DownMbps int `json:"down_mbps,omitempty"`
	// Obfuscated password.
	//
	// @Example: "fuck me till the daylight"
	Obfs string `json:"obfs,omitempty"`
	// Hysteria users
	//
	// @Example:
	//    [
	//      {
	//        "name": "sekai",
	//        "auth": "",
	//        "auth_str": "password"
	//      }
	//    ]
	Users []HysteriaUser `json:"users,omitempty"`
	// @Example: 0
	// @Deprecated: use QUIC fields instead
	ReceiveWindowConn uint64 `json:"recv_window_conn,omitempty" schema:"omit"`
	// @Example: 0
	// @Deprecated: use QUIC fields instead
	ReceiveWindowClient uint64 `json:"recv_window_client,omitempty" schema:"omit"`
	// @Example: 0
	// @Deprecated: use QUIC fields instead
	MaxConnClient int `json:"max_conn_client,omitempty" schema:"omit"`
	// @Example: false
	// @Deprecated: use QUIC fields instead
	DisableMTUDiscovery bool `json:"disable_mtu_discovery,omitempty" schema:"omit"`
	// @Override:
	//    @Required
	//
	//    TLS configuration, see [TLS](/configuration/shared/tls/#inbound).
	InboundTLSOptionsContainer
	QUICOptions
}

type HysteriaUser struct {
	Name string `json:"name,omitempty"`
	// Authentication password, in base64.
	Auth []byte `json:"auth,omitempty"`
	// Authentication password.
	AuthString string `json:"auth_str,omitempty"`
}

type HysteriaOutboundOptions struct {
	DialerOptions
	// @Override server:
	//    @Example: "127.0.0.1"
	//
	// @Override server_port:
	//    @Example: 1080
	ServerOptions
	// Server port range list.
	//
	// Conflicts with `server_port`.
	//
	// @Example:
	//    [
	//      "2080:3000"
	//    ]
	//
	// @Since 1.12.0
	ServerPorts LegacyListable[string] `json:"server_ports,omitempty"`
	// Port hopping interval.
	//
	// @Default: 30s
	//
	// @Since 1.12.0
	HopInterval badoption.Duration `json:"hop_interval,omitempty"`
	// @Example: "100 Mbps"
	// @Omit
	Up *byteformats.NetworkBytesCompat `json:"up,omitempty"`
	// @Example: 100
	// @Omit
	UpMbps int `json:"up_mbps,omitempty"`
	// @Example: "100 Mbps"
	// @Omit
	Down *byteformats.NetworkBytesCompat `json:"down,omitempty"`
	// @Example: 100
	// @Omit
	DownMbps int `json:"down_mbps,omitempty"`
	// Obfuscated password.
	//
	// @Example: "fuck me till the daylight"
	Obfs string `json:"obfs,omitempty"`
	// Authentication password, in base64.
	//
	// @Example: ""
	Auth []byte `json:"auth,omitempty"`
	// Authentication password.
	//
	// @Example: "password"
	AuthString string `json:"auth_str,omitempty"`
	// @Example: 0
	// @Deprecated: use QUIC fields instead
	ReceiveWindowConn uint64 `json:"recv_window_conn,omitempty" schema:"omit"`
	// @Example: 0
	// @Deprecated: use QUIC fields instead
	ReceiveWindow uint64 `json:"recv_window,omitempty" schema:"omit"`
	// @Example: false
	// @Deprecated: use QUIC fields instead
	DisableMTUDiscovery bool `json:"disable_mtu_discovery,omitempty" schema:"omit"`
	// Enabled network
	//
	// One of `tcp` `udp`.
	//
	// Both is enabled by default.
	//
	// @Example: ""
	Network NetworkList `json:"network,omitempty"`
	// @Override:
	//    @Required
	//
	//    TLS configuration, see [TLS](/configuration/shared/tls/#outbound).
	OutboundTLSOptionsContainer
	QUICOptions
}
