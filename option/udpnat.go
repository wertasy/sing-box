package option

import (
	"github.com/sagernet/sing/common/json/badoption"
)

// UDPNATOptions documents the UDP NAT fields shared by tunnel, endpoint and
// outbound options in the configuration reference. The runtime declarations
// live on the individual option structs.
type UDPNATOptions struct {
	// UDP NAT expiration time.
	//
	// `5m` will be used by default.
	UDPTimeout badoption.Duration `json:"udp_timeout,omitempty"`

	// @Since 1.14.0
	//
	// UDP NAT mapping behavior.
	//
	// | Value                        | Behavior                                                                      |
	// |------------------------------|-------------------------------------------------------------------------------|
	// | `endpoint_independent`       | Reuse the same mapping for the same source address and port for all destinations. |
	// | `address_dependent`          | Use a separate mapping for each destination address.                          |
	// | `address_and_port_dependent` | Use a separate mapping for each destination address and port.                 |
	//
	// `endpoint_independent` will be used by default.
	UDPMapping UDPNATBehavior `json:"udp_mapping,omitempty"`

	// @Since 1.14.0
	//
	// UDP NAT filtering behavior.
	//
	// | Value                        | Behavior                                                                    |
	// |------------------------------|-----------------------------------------------------------------------------|
	// | `endpoint_independent`       | Accept packets from any remote endpoint.                                    |
	// | `address_dependent`          | Accept packets only from remote addresses to which packets have been sent.  |
	// | `address_and_port_dependent` | Accept packets only from remote addresses and ports to which packets have been sent. |
	//
	// `endpoint_independent` will be used by default.
	UDPFiltering UDPNATBehavior `json:"udp_filtering,omitempty"`

	// @Since 1.14.0
	//
	// Maximum number of UDP NAT sessions.
	//
	// When the limit is reached, the least recently used session is closed.
	//
	// When unset or set to `0`, `4096` is used on iOS. On other platforms, a value from `4096` to `16384` is selected based on total memory;
	// `16384` is used if total memory cannot be detected.
	UDPNATMax uint32 `json:"udp_nat_max,omitempty"`
}
