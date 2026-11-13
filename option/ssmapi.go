package option

import (
	"github.com/sagernet/sing/common/json/badjson"
)

type SSMAPIServiceOptions struct {
	ListenOptions

	// A mapping Object from HTTP endpoints to [Shadowsocks Inbound](/configuration/inbound/shadowsocks) tags.
	//
	// Selected Shadowsocks inbounds must be configured with [managed](/configuration/inbound/shadowsocks#managed) enabled.
	//
	// Example:
	//
	// ```json
	// {
	//   "servers": {
	//     "/": "ss-in"
	//   }
	// }
	// ```
	//
	// @Required
	Servers *badjson.TypedMap[string, string] `json:"servers"`

	// If set, when the server is about to stop, traffic and user state will be saved to the specified JSON file
	// to be restored on the next startup.
	CachePath string `json:"cache_path,omitempty"`
	InboundTLSOptionsContainer
}
