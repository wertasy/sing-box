package option

import "github.com/sagernet/sing/common/json/badoption"

type SelectorOutboundOptions struct {
	// @required
	//
	// List of outbound tags to select.
	Outbounds []string `json:"outbounds" reference:"outbound"`
	// The default outbound tag. The first outbound will be used if empty.
	Default string `json:"default,omitempty" reference:"outbound"`
	// Interrupt existing connections when the selected outbound has changed.
	//
	// Only inbound connections are affected by this setting, internal connections will always be interrupted.
	InterruptExistConnections bool `json:"interrupt_exist_connections,omitempty"`
}

type URLTestOutboundOptions struct {
	// @required
	//
	// List of outbound tags to test.
	Outbounds []string `json:"outbounds" reference:"outbound"`
	// The URL to test. `https://www.gstatic.com/generate_204` will be used if empty.
	URL string `json:"url,omitempty"`
	// The test interval. `3m` will be used if empty.
	Interval badoption.Duration `json:"interval,omitempty"`
	// The test tolerance in milliseconds. `50` will be used if empty.
	Tolerance uint16 `json:"tolerance,omitempty"`
	// The idle timeout. `30m` will be used if empty.
	IdleTimeout badoption.Duration `json:"idle_timeout,omitempty"`
	// Interrupt existing connections when the selected outbound has changed.
	//
	// Only inbound connections are affected by this setting, internal connections will always be interrupted.
	InterruptExistConnections bool `json:"interrupt_exist_connections,omitempty"`
}
