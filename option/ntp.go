package option

import "github.com/sagernet/sing/common/json/badoption"

type NTPOptions struct {
	// Enable NTP service.
	Enabled bool `json:"enabled,omitempty"`
	// Time synchronization interval.
	//
	// @Default: 30 minutes
	Interval badoption.Duration `json:"interval,omitempty"`
	// Update system time after synchronization.
	WriteToSystem bool `json:"write_to_system,omitempty"`
	// @Override server:
	//    @Required
	//
	//    NTP server address.
	//
	// @Override server_port:
	//    NTP server port.
	//
	//    123 is used by default.
	ServerOptions
	DialerOptions
}
