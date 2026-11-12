package option

// @Tag tor-out
type TorOutboundOptions struct {
	DialerOptions
	// The path to the Tor executable.
	//
	// Embedded Tor will be ignored if set.
	//
	// @Example: "/usr/bin/tor"
	ExecutablePath string `json:"executable_path,omitempty"`
	// List of extra arguments passed to the Tor instance when started.
	ExtraArgs []string `json:"extra_args,omitempty"`
	// The data directory of Tor.
	//
	// Each start will be very slow if not specified.
	//
	// @Example: "$HOME/.cache/tor"
	//
	// @Indicator: Recommended
	DataDirectory string `json:"data_directory,omitempty"`
	// Map of torrc options.
	//
	// See [tor(1)](https://linux.die.net/man/1/tor) for details.
	//
	// @Example:
	//    {
	//      "ClientOnly": 1
	//    }
	Options map[string]string `json:"torrc,omitempty"`
}
