package option

import "github.com/sagernet/sing/common/json/badoption"

type TunPlatformOptions struct {
	// System HTTP proxy settings.
	HTTPProxy *HTTPProxyOptions `json:"http_proxy,omitempty"`
}

type HTTPProxyOptions struct {
	// Enable system HTTP proxy.
	Enabled bool `json:"enabled,omitempty"`

	// @Override server:
	//    @Required
	//
	//    HTTP proxy server address.
	//
	// @Override server_port:
	//    @Required
	//
	//    HTTP proxy server port.
	ServerOptions

	// Hostnames that bypass the HTTP proxy.
	//
	// @Note:
	//
	//     On Apple platforms, `bypass_domain` items matches hostname **suffixes**.
	BypassDomain badoption.Listable[string] `json:"bypass_domain,omitempty"`

	// Hostnames that use the HTTP proxy.
	//
	// @Quote:
	//
	//     Only supported in graphical clients on Apple platforms.
	MatchDomain badoption.Listable[string] `json:"match_domain,omitempty"`
}
