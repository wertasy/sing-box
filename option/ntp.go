package option

import "github.com/sagernet/sing/common/json/badoption"

type NTPOptions struct {
	// Enable NTP service.
	Enabled bool `json:"enabled,omitempty"`
	// Time synchronization interval.
	//
	// 30 minutes is used by default.
	Interval badoption.Duration `json:"interval,omitempty"`
	// Update system time after synchronization.
	WriteToSystem bool `json:"write_to_system,omitempty"`
	ServerOptions
	DialerOptions
}
