package option

import "github.com/sagernet/sing/common/json/badoption"

type RouteOptions struct {
	// @Omit
	GeoIP *GeoIPOptions `json:"geoip,omitempty" schema:"omit"`

	// @Omit
	Geosite *GeositeOptions `json:"geosite,omitempty" schema:"omit"`

	// List of [Route Rule](./rule/)
	Rules []Rule `json:"rules,omitempty"`

	// List of [rule-set](/configuration/rule-set/)
	//
	// @Since 1.8.0
	RuleSet []RuleSet `json:"rule_set,omitempty"`

	// Default outbound tag. the first outbound will be used if empty.
	Final string `json:"final,omitempty" reference:"outbound"`

	// Enable process search for logging when no `process_name`, `process_path`, `package_name`, `user` or `user_id` rules exist.
	//
	// @Quote:
	//
	//     Only supported on Linux, Windows, and macOS.
	FindProcess bool `json:"find_process,omitempty"`

	// Enable neighbor resolution for logging when no `source_mac_address` or `source_hostname` rules exist.
	//
	// See [Neighbor Resolution](/configuration/shared/neighbor/) for setup.
	//
	// @Since 1.14.0
	//
	// @Quote:
	//
	//     Only supported on Linux and macOS.
	FindNeighbor bool `json:"find_neighbor,omitempty"`

	// Custom DHCP lease file paths for hostname and MAC address resolution.
	//
	// Automatically detected from common DHCP servers (dnsmasq, odhcpd, ISC dhcpd, Kea) if empty.
	//
	// @Since 1.14.0
	//
	// @Quote:
	//
	//     Only supported on Linux and macOS.
	DHCPLeaseFiles LegacyListable[string] `json:"dhcp_lease_files,omitempty"`

	// Bind outbound connections to the default NIC by default to prevent routing loops under tun.
	//
	// Takes no effect if `outbound.bind_interface` is set.
	//
	// @Quote:
	//
	//     Only supported on Linux, Windows and macOS.
	AutoDetectInterface bool `json:"auto_detect_interface,omitempty"`

	// Accept Android VPN as upstream NIC when `auto_detect_interface` enabled.
	//
	// @Quote:
	//
	//     Only supported on Android.
	OverrideAndroidVPN bool `json:"override_android_vpn,omitempty"`

	// Bind outbound connections to the specified NIC by default to prevent routing loops under tun.
	//
	// Takes no effect if `auto_detect_interface` is set.
	//
	// @Quote:
	//
	//     Only supported on Linux, Windows and macOS.
	DefaultInterface string `json:"default_interface,omitempty"`

	// Set routing mark by default.
	//
	// Takes no effect if `outbound.routing_mark` is set.
	//
	// @Quote:
	//
	//     Only supported on Linux.
	DefaultMark FwMark `json:"default_mark,omitempty"`

	// See [Dial Fields](/configuration/shared/dial/#domain_resolver) for details.
	//
	// Can be overridden by `outbound.domain_resolver`.
	//
	// @Since 1.12.0
	DefaultDomainResolver *DomainResolveOptions `json:"default_domain_resolver,omitempty"`

	// See [Dial Fields](/configuration/shared/dial/#network_strategy) for details.
	//
	// Takes no effect if `outbound.bind_interface`, `outbound.inet4_bind_address` or `outbound.inet6_bind_address` is set.
	//
	// Can be overridden by `outbound.network_strategy`.
	//
	// Conflicts with `default_interface`.
	//
	// @Since 1.11.0
	DefaultNetworkStrategy *NetworkStrategy `json:"default_network_strategy,omitempty"`

	// See [Dial Fields](/configuration/shared/dial/#network_type) for details.
	//
	// @Since 1.11.0
	DefaultNetworkType badoption.Listable[InterfaceType] `json:"default_network_type,omitempty"`

	// See [Dial Fields](/configuration/shared/dial/#fallback_network_type) for details.
	//
	// @Since 1.11.0
	DefaultFallbackNetworkType badoption.Listable[InterfaceType] `json:"default_fallback_network_type,omitempty"`

	// See [Dial Fields](/configuration/shared/dial/#fallback_delay) for details.
	//
	// @Since 1.11.0
	DefaultFallbackDelay badoption.Duration `json:"default_fallback_delay,omitempty"`

	// Tag of the default [HTTP Client](/configuration/shared/http-client/) used by remote rule-sets.
	//
	// If empty and `http_clients` is defined, the first HTTP client is used.
	//
	// @Since 1.14.0
	DefaultHTTPClient string `json:"default_http_client,omitempty"`
}

type GeoIPOptions struct {
	Path           string `json:"path,omitempty"`
	DownloadURL    string `json:"download_url,omitempty"`
	DownloadDetour string `json:"download_detour,omitempty" reference:"outbound"`
}

type GeositeOptions struct {
	Path           string `json:"path,omitempty"`
	DownloadURL    string `json:"download_url,omitempty"`
	DownloadDetour string `json:"download_detour,omitempty" reference:"outbound"`
}
