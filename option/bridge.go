package option

type BridgeOutboundOptions struct {
	// Interface name for forwarded traffic to egress.
	//
	// The default interface will be used by default.
	Interface string `json:"interface,omitempty"`
	// Custom bridge TUN interface name prefix, `bridge` is used by default.
	//
	// Not effective on Apple platforms.
	BridgeName string `json:"bridge_name,omitempty"`
	// @quote:
	//
	//     Only supported on Linux.
	//
	// Linux iproute2 table index.
	//
	// `2200` + instance index is used by default.
	IPRoute2TableIndex int `json:"iproute2_table_index,omitempty"`
	// @quote:
	//
	//     Only supported on Linux.
	//
	// Linux iproute2 rule start index.
	//
	// `100` is used by default.
	IPRoute2RuleIndex int `json:"iproute2_rule_index,omitempty"`
}
