package option

// @Tag bridge-out
type BridgeOutboundOptions struct {
	// Interface name for forwarded traffic to egress.
	//
	// The default interface will be used by default.
	//
	// @Example: ""
	Interface string `json:"interface,omitempty"`

	// Custom bridge TUN interface name prefix, `bridge` will be used by default.
	//
	// Not effective on Apple platforms.
	//
	// @Example: ""
	BridgeName string `json:"bridge_name,omitempty"`

	// Linux iproute2 table index.
	//
	// `2200` + instance index will be used by default.
	//
	// @Quote:
	//
	//     Only supported on Linux.
	//
	// @Example: 0
	IPRoute2TableIndex int `json:"iproute2_table_index,omitempty"`

	// Linux iproute2 rule start index.
	//
	// @Default: 100
	//
	// @Quote:
	//
	//     Only supported on Linux.
	//
	// @Example: 0
	IPRoute2RuleIndex int `json:"iproute2_rule_index,omitempty"`
}
