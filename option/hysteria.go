package option

import (
	"github.com/sagernet/sing/common/byteformats"
	"github.com/sagernet/sing/common/json/badoption"
)

type HysteriaInboundOptions struct {
	ListenOptions
	Up       *byteformats.NetworkBytesCompat `json:"up,omitempty"`
	UpMbps   int                             `json:"up_mbps,omitempty"`
	Down     *byteformats.NetworkBytesCompat `json:"down,omitempty"`
	DownMbps int                             `json:"down_mbps,omitempty"`
	// Obfuscated password.
	Obfs string `json:"obfs,omitempty"`
	// Hysteria users
	Users []HysteriaUser `json:"users,omitempty"`
	// Deprecated: use QUIC fields instead
	ReceiveWindowConn uint64 `json:"recv_window_conn,omitempty" schema:"omit"`
	// Deprecated: use QUIC fields instead
	ReceiveWindowClient uint64 `json:"recv_window_client,omitempty" schema:"omit"`
	// Deprecated: use QUIC fields instead
	MaxConnClient int `json:"max_conn_client,omitempty" schema:"omit"`
	// Deprecated: use QUIC fields instead
	DisableMTUDiscovery bool `json:"disable_mtu_discovery,omitempty" schema:"omit"`
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
	ServerOptions
	// @since 1.12.0
	//
	// Server port range list.
	//
	// Conflicts with `server_port`.
	ServerPorts LegacyListable[string] `json:"server_ports,omitempty"`
	// @since 1.12.0
	//
	// Port hopping interval.
	//
	// `30s` is used by default.
	HopInterval badoption.Duration              `json:"hop_interval,omitempty"`
	Up          *byteformats.NetworkBytesCompat `json:"up,omitempty"`
	UpMbps      int                             `json:"up_mbps,omitempty"`
	Down        *byteformats.NetworkBytesCompat `json:"down,omitempty"`
	DownMbps    int                             `json:"down_mbps,omitempty"`
	// Obfuscated password.
	Obfs string `json:"obfs,omitempty"`
	// Authentication password, in base64.
	Auth []byte `json:"auth,omitempty"`
	// Authentication password.
	AuthString string `json:"auth_str,omitempty"`
	// Deprecated: use QUIC fields instead
	ReceiveWindowConn uint64 `json:"recv_window_conn,omitempty" schema:"omit"`
	// Deprecated: use QUIC fields instead
	ReceiveWindow uint64 `json:"recv_window,omitempty" schema:"omit"`
	// Deprecated: use QUIC fields instead
	DisableMTUDiscovery bool `json:"disable_mtu_discovery,omitempty" schema:"omit"`
	// Enabled network
	//
	// One of `tcp` `udp`.
	//
	// Both is enabled by default.
	Network NetworkList `json:"network,omitempty"`
	OutboundTLSOptionsContainer
	QUICOptions
}
