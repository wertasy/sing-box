package option

import "github.com/sagernet/sing/common/json/badoption"

type RouteOptions struct {
	GeoIP   *GeoIPOptions   `json:"geoip,omitempty" schema:"omit"`
	Geosite *GeositeOptions `json:"geosite,omitempty" schema:"omit"`
	// List of [Route Rule](./rule/)
	Rules []Rule `json:"rules,omitempty"`
	// @since 1.8.0
	//
	// List of [rule-set](/configuration/rule-set/)
	RuleSet []RuleSet `json:"rule_set,omitempty"`
	// Default outbound tag. the first outbound will be used if empty.
	Final string `json:"final,omitempty" reference:"outbound"`
	// @quote:
	//
	//     Only supported on Linux, Windows, and macOS.
	//
	// Enable process search for logging when no `process_name`, `process_path`, `package_name`, `user` or `user_id` rules exist.
	FindProcess bool `json:"find_process,omitempty"`
	// @since 1.14.0
	//
	// @quote:
	//
	//     Only supported on Linux and macOS.
	//
	// Enable neighbor resolution for logging when no `source_mac_address` or `source_hostname` rules exist.
	//
	// See [Neighbor Resolution](/configuration/shared/neighbor/) for setup.
	FindNeighbor bool `json:"find_neighbor,omitempty"`
	// @since 1.14.0
	//
	// @quote:
	//
	//     Only supported on Linux and macOS.
	//
	// Custom DHCP lease file paths for hostname and MAC address resolution.
	//
	// Automatically detected from common DHCP servers (dnsmasq, odhcpd, ISC dhcpd, Kea) if empty.
	DHCPLeaseFiles LegacyListable[string] `json:"dhcp_lease_files,omitempty"`
	// @quote:
	//
	//     Only supported on Linux, Windows and macOS.
	//
	// Bind outbound connections to the default NIC by default to prevent routing loops under tun.
	//
	// Takes no effect if `outbound.bind_interface` is set.
	AutoDetectInterface bool `json:"auto_detect_interface,omitempty"`
	// @quote:
	//
	//     Only supported on Android.
	//
	// Accept Android VPN as upstream NIC when `auto_detect_interface` enabled.
	OverrideAndroidVPN bool `json:"override_android_vpn,omitempty"`
	// @quote:
	//
	//     Only supported on Linux, Windows and macOS.
	//
	// Bind outbound connections to the specified NIC by default to prevent routing loops under tun.
	//
	// Takes no effect if `auto_detect_interface` is set.
	DefaultInterface string `json:"default_interface,omitempty"`
	// @quote:
	//
	//     Only supported on Linux.
	//
	// Set routing mark by default.
	//
	// Takes no effect if `outbound.routing_mark` is set.
	DefaultMark FwMark `json:"default_mark,omitempty"`
	// @since 1.12.0
	//
	// See [Dial Fields](/configuration/shared/dial/#domain_resolver) for details.
	//
	// Can be overridden by `outbound.domain_resolver`.
	DefaultDomainResolver *DomainResolveOptions `json:"default_domain_resolver,omitempty"`
	// @since 1.11.0
	//
	// See [Dial Fields](/configuration/shared/dial/#network_strategy) for details.
	//
	// Takes no effect if `outbound.bind_interface`, `outbound.inet4_bind_address` or `outbound.inet6_bind_address` is set.
	//
	// Can be overridden by `outbound.network_strategy`.
	//
	// Conflicts with `default_interface`.
	DefaultNetworkStrategy *NetworkStrategy `json:"default_network_strategy,omitempty"`
	// @since 1.11.0
	//
	// See [Dial Fields](/configuration/shared/dial/#network_type) for details.
	DefaultNetworkType badoption.Listable[InterfaceType] `json:"default_network_type,omitempty"`
	// @since 1.11.0
	//
	// See [Dial Fields](/configuration/shared/dial/#fallback_network_type) for details.
	DefaultFallbackNetworkType badoption.Listable[InterfaceType] `json:"default_fallback_network_type,omitempty"`
	// @since 1.11.0
	//
	// See [Dial Fields](/configuration/shared/dial/#fallback_delay) for details.
	DefaultFallbackDelay badoption.Duration `json:"default_fallback_delay,omitempty"`
	// @since 1.14.0
	//
	// Tag of the default [HTTP Client](/configuration/shared/http-client/) used by remote rule-sets.
	//
	// If empty and `http_clients` is defined, the first HTTP client is used.
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
