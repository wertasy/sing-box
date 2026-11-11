package option

type TorOutboundOptions struct {
	DialerOptions
	// The path to the Tor executable.
	//
	// Embedded Tor will be ignored if set.
	ExecutablePath string `json:"executable_path,omitempty"`
	// List of extra arguments passed to the Tor instance when started.
	ExtraArgs []string `json:"extra_args,omitempty"`
	// @indicator: Recommended
	//
	// The data directory of Tor.
	//
	// Each start will be very slow if not specified.
	DataDirectory string `json:"data_directory,omitempty"`
	// Map of torrc options.
	//
	// See [tor(1)](https://linux.die.net/man/1/tor) for details.
	Options map[string]string `json:"torrc,omitempty"`
}
