package option

import "github.com/sagernet/sing/common/json/badoption"

type CloudflaredInboundOptions struct {
	// Base64-encoded tunnel token from the Cloudflare Zero Trust dashboard
	// (`Networks → Tunnels → Install connector`).
	//
	// @Required
	Token string `json:"token,omitempty"`
	// Number of high-availability connections to the Cloudflare edge.
	//
	// Capped by the number of discovered edge addresses.
	HighAvailabilityConnections int `json:"ha_connections,omitempty"`
	// Transport protocol for edge connections.
	//
	// One of `quic` `http2`.
	Protocol string `json:"protocol,omitempty" enum:"auto,quic,http2,h2mux"`
	// Enable post-quantum key exchange on the control connection.
	PostQuantum bool `json:"post_quantum,omitempty"`
	// IP version used when connecting to the Cloudflare edge.
	//
	// One of `0` (automatic) `4` `6`.
	EdgeIPVersion int `json:"edge_ip_version,omitempty" enum:"0,4,6"`
	// Datagram protocol version used for UDP proxying over QUIC.
	//
	// One of `v2` `v3`. Only meaningful when `protocol` is `quic`.
	DatagramVersion string `json:"datagram_version,omitempty" enum:"v2,v3"`
	// Graceful shutdown window for in-flight edge connections.
	GracePeriod badoption.Duration `json:"grace_period,omitempty"`
	// Cloudflare edge region selector.
	//
	// Conflict with endpoints embedded in `token`.
	Region string `json:"region,omitempty"`
	// [Dial Fields](/configuration/shared/dial/) used when the tunnel client dials the
	// Cloudflare control plane.
	//
	// @Example:
	//    {
	//      ... // Dial Fields
	//    }
	ControlDialer DialerOptions `json:"control_dialer,omitempty"`
	// [Dial Fields](/configuration/shared/dial/) used when the tunnel client dials the
	// Cloudflare edge data plane.
	//
	// @Example:
	//    {
	//      ... // Dial Fields
	//    }
	TunnelDialer DialerOptions `json:"tunnel_dialer,omitempty"`
}
