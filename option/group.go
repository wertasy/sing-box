package option

import "github.com/sagernet/sing/common/json/badoption"

// @Tag select
type SelectorOutboundOptions struct {
	// List of outbound tags to select.
	//
	// @Example:
	//    [
	//      "proxy-a",
	//      "proxy-b",
	//      "proxy-c"
	//    ]
	//
	// @Required
	Outbounds []string `json:"outbounds" reference:"outbound"`

	// The default outbound tag. The first outbound will be used if empty.
	//
	// @Example: "proxy-c"
	Default string `json:"default,omitempty" reference:"outbound"`

	// Interrupt existing connections when the selected outbound has changed.
	//
	// Only inbound connections are affected by this setting, internal connections will always be interrupted.
	//
	// @Example: false
	InterruptExistConnections bool `json:"interrupt_exist_connections,omitempty"`
}

// @Tag auto
type URLTestOutboundOptions struct {
	// List of outbound tags to test.
	//
	// @Example:
	//    [
	//      "proxy-a",
	//      "proxy-b",
	//      "proxy-c"
	//    ]
	//
	// @Required
	Outbounds []string `json:"outbounds" reference:"outbound"`

	// The URL to test. `https://www.gstatic.com/generate_204` will be used if empty.
	//
	// @Example: ""
	URL string `json:"url,omitempty"`

	// The test interval. `3m` will be used if empty.
	//
	// @Example: ""
	Interval badoption.Duration `json:"interval,omitempty"`

	// The test tolerance in milliseconds. `50` will be used if empty.
	//
	// @Example: 0
	Tolerance uint16 `json:"tolerance,omitempty"`

	// The idle timeout. `30m` will be used if empty.
	//
	// @Example: ""
	IdleTimeout badoption.Duration `json:"idle_timeout,omitempty"`

	// Interrupt existing connections when the selected outbound has changed.
	//
	// Only inbound connections are affected by this setting, internal connections will always be interrupted.
	//
	// @Example: false
	InterruptExistConnections bool `json:"interrupt_exist_connections,omitempty"`
}
