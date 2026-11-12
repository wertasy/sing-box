package option

type RedirectInboundOptions struct {
	ListenOptions
}

// @Tag tproxy-in
type TProxyInboundOptions struct {
	ListenOptions
	// Listen network, one of `tcp` `udp`.
	//
	// Both if empty.
	//
	// @Example: "udp"
	Network NetworkList `json:"network,omitempty"`
	// @Omit
	UDPMapping UDPNATBehavior `json:"udp_mapping,omitempty"`
	// @Omit
	UDPFiltering UDPNATBehavior `json:"udp_filtering,omitempty"`
	// @Omit
	UDPNATMax uint32 `json:"udp_nat_max,omitempty"`
}
