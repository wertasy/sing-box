package option

type InboundMultiplexOptions struct {
	// Enable multiplex support.
	Enabled bool `json:"enabled,omitempty"`
	// If enabled, non-padded connections will be rejected.
	Padding bool `json:"padding,omitempty"`
	// See [TCP Brutal](/configuration/shared/tcp-brutal/) for details.
	Brutal *BrutalOptions `json:"brutal,omitempty"`
}

type OutboundMultiplexOptions struct {
	// Enable multiplex.
	Enabled bool `json:"enabled,omitempty"`
	// Multiplex protocol.
	//
	// | Protocol | Description                        |
	// |----------|------------------------------------|
	// | smux     | https://github.com/xtaci/smux      |
	// | yamux    | https://github.com/hashicorp/yamux |
	// | h2mux    | https://golang.org/x/net/http2     |
	//
	// h2mux is used by default.
	Protocol string `json:"protocol,omitempty" enum:"h2mux,smux,yamux"`
	// Maximum connections.
	//
	// Conflict with `max_streams`.
	MaxConnections int `json:"max_connections,omitempty"`
	// Minimum multiplexed streams in a connection before opening a new connection.
	//
	// Conflict with `max_streams`.
	MinStreams int `json:"min_streams,omitempty"`
	// Maximum multiplexed streams in a connection before opening a new connection.
	//
	// Conflict with `max_connections` and `min_streams`.
	MaxStreams int `json:"max_streams,omitempty"`
	// @info
	//
	//     Requires sing-box server version 1.3-beta9 or later.
	//
	// Enable padding.
	Padding bool `json:"padding,omitempty"`
	// See [TCP Brutal](/configuration/shared/tcp-brutal/) for details.
	Brutal *BrutalOptions `json:"brutal,omitempty"`
}

type BrutalOptions struct {
	// Enable TCP Brutal congestion control algorithm。
	Enabled  bool `json:"enabled,omitempty"`
	UpMbps   int  `json:"up_mbps,omitempty"`
	DownMbps int  `json:"down_mbps,omitempty"`
}
