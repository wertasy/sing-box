package option

import (
	"github.com/sagernet/sing/common/byteformats"
	"github.com/sagernet/sing/common/json/badoption"
)

type ExperimentalOptions struct {
	CacheFile *CacheFileOptions `json:"cache_file,omitempty"`
	ClashAPI  *ClashAPIOptions  `json:"clash_api,omitempty"`
	V2RayAPI  *V2RayAPIOptions  `json:"v2ray_api,omitempty"`
	Debug     *DebugOptions     `json:"debug,omitempty"`
}

type CacheFileOptions struct {
	// Enable cache file.
	//
	// @Example: true
	Enabled bool `json:"enabled,omitempty"`
	// Path to the cache file.
	//
	// @Default: cache.db
	Path string `json:"path,omitempty"`
	// Identifier in the cache file
	//
	// If not empty, configuration specified data will use a separate store keyed by it.
	CacheID string `json:"cache_id,omitempty"`
	// Store fakeip in the cache file
	StoreFakeIP bool `json:"store_fakeip,omitempty"`
	// Store DNS cache in the cache file.
	//
	// @Since 1.14.0
	StoreDNS bool `json:"store_dns,omitempty"`
	// Size of the write buffer.
	//
	// @Default: 1MB
	//
	// @Example: ""
	//
	// @Since 1.15.0
	BufferSize *byteformats.MemoryBytes `json:"buffer_size,omitempty"`
	// Interval for flushing the write buffer automatically.
	//
	// Disabled by default.
	//
	// @Since 1.15.0
	FlushInterval badoption.Duration `json:"flush_interval,omitempty"`

	// Store rejected DNS response cache in the cache file
	//
	// The check results of [Legacy Address Filter Fields](/configuration/dns/rule/#legacy-address-filter-fields)
	// will be cached until expiration.
	//
	// @Failure: Deprecated in sing-box 1.14.0
	//
	//     `store_rdrc` is deprecated and will be removed in sing-box 1.16.0, check [Migration](/migration/#migrate-store-rdrc).
	//
	// @Example: false
	//
	// @Deprecated: replaced by store_dns
	StoreRDRC bool `json:"store_rdrc,omitempty" schema:"omit"`
	// Timeout of rejected DNS response cache.
	//
	// @Failure: Deprecated in sing-box 1.14.0
	//
	//     `rdrc_timeout` is deprecated and will be removed in sing-box 1.16.0, check [Migration](/migration/#migrate-store-rdrc).
	//
	// @Default: 7d
	//
	// @Deprecated: replaced by store_dns
	RDRCTimeout badoption.Duration `json:"rdrc_timeout,omitempty" schema:"omit"`
}

type ClashAPIOptions struct {
	// RESTful web API listening address. Clash API will be disabled if empty.
	ExternalController string `json:"external_controller,omitempty"`
	// A relative path to the configuration directory or an absolute path to a
	// directory in which you put some static web resource. sing-box will then
	// serve it at `http://{{external-controller}}/ui`.
	ExternalUI string `json:"external_ui,omitempty"`
	// ZIP download URL for the external UI, will be used if the specified `external_ui` directory is empty.
	//
	// @Default: https://github.com/MetaCubeX/Yacd-meta/archive/gh-pages.zip
	ExternalUIDownloadURL string `json:"external_ui_download_url,omitempty"`
	// The tag of the outbound to download the external UI.
	//
	// Default outbound will be used if empty.
	ExternalUIDownloadDetour string `json:"external_ui_download_detour,omitempty" reference:"outbound"`
	// Secret for the RESTful API (optional)
	// Authenticate by spedifying HTTP header `Authorization: Bearer ${secret}`
	// ALWAYS set a secret if RESTful API is listening on 0.0.0.0
	Secret string `json:"secret,omitempty"`
	// Default mode in clash, `Rule` will be used if empty.
	//
	// This setting has no direct effect, but can be used in routing and DNS rules via the `clash_mode` rule item.
	DefaultMode string `json:"default_mode,omitempty"`
	// CORS allowed origins, `*` will be used if empty.
	//
	// To access the Clash API on a private network from a public website, you must explicitly specify it in `access_control_allow_origin` instead of using `*`.
	//
	// @Since 1.10.0
	AccessControlAllowOrigin badoption.Listable[string] `json:"access_control_allow_origin,omitempty"`
	// Allow access from private network.
	//
	// To access the Clash API on a private network from a public website, `access_control_allow_private_network` must be enabled.
	//
	// @Since 1.10.0
	AccessControlAllowPrivateNetwork bool `json:"access_control_allow_private_network,omitempty"`

	// Cache file path, `cache.db` will be used if empty.
	//
	// @Failure: Deprecated in sing-box 1.8.0
	//
	//     `cache_file` is deprecated in Clash API and migrated to `cache_file.enabled` and `cache_file.path`.
	//
	// @Deprecated: migrated to global cache file
	CacheFile string `json:"cache_file,omitempty" schema:"omit"`
	// Identifier in cache file.
	//
	// If not empty, configuration specified data will use a separate store keyed by it.
	//
	// @Failure: Deprecated in sing-box 1.8.0
	//
	//     `cache_id` is deprecated in Clash API and migrated to `cache_file.cache_id`.
	//
	// @Deprecated: migrated to global cache file
	CacheID string `json:"cache_id,omitempty" schema:"omit"`
	// Store Clash mode in cache file.
	//
	// @Failure: Deprecated in sing-box 1.8.0
	//
	//     `store_mode` is deprecated in Clash API and enabled by default if `cache_file.enabled`.
	//
	// @Deprecated: migrated to global cache file
	StoreMode bool `json:"store_mode,omitempty" schema:"omit"`
	// Store selected outbound for the `Selector` outbound in cache file.
	//
	// @Failure: Deprecated in sing-box 1.8.0
	//
	//     `store_selected` is deprecated in Clash API and enabled by default if `cache_file.enabled`.
	//
	// @Deprecated: migrated to global cache file
	//
	// @Note:
	//
	//     The tag must be set for target outbounds.
	StoreSelected bool `json:"store_selected,omitempty" schema:"omit"`
	// Store fakeip in cache file.
	//
	// @Failure: Deprecated in sing-box 1.8.0
	//
	//     `store_selected` is deprecated in Clash API and migrated to `cache_file.store_fakeip`.
	//
	// @Deprecated: migrated to global cache file
	StoreFakeIP bool `json:"store_fakeip,omitempty" schema:"omit"`
}

type V2RayAPIOptions struct {
	// gRPC API listening address. V2Ray API will be disabled if empty.
	//
	// @Example: "127.0.0.1:8080"
	Listen string `json:"listen,omitempty"`
	// Traffic statistics service settings.
	//
	// @Example:
	//    {
	//      "enabled": true,
	//      "inbounds": [
	//        "socks-in"
	//      ],
	//      "outbounds": [
	//        "proxy",
	//        "direct"
	//      ],
	//      "users": [
	//        "sekai"
	//      ]
	//    }
	Stats *V2RayStatsServiceOptions `json:"stats,omitempty"`
}

type V2RayStatsServiceOptions struct {
	// Enable statistics service.
	Enabled bool `json:"enabled,omitempty"`
	// Inbound list to count traffic.
	Inbounds []string `json:"inbounds,omitempty"`
	// Outbound list to count traffic.
	Outbounds []string `json:"outbounds,omitempty"`
	// User list to count traffic.
	Users []string `json:"users,omitempty"`
}
