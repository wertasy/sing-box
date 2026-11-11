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
	Enabled bool `json:"enabled,omitempty"`
	// Path to the cache file.
	//
	// `cache.db` will be used if empty.
	Path string `json:"path,omitempty"`
	// Identifier in the cache file
	//
	// If not empty, configuration specified data will use a separate store keyed by it.
	CacheID string `json:"cache_id,omitempty"`
	// Store fakeip in the cache file
	StoreFakeIP bool `json:"store_fakeip,omitempty"`
	// @since 1.14.0
	//
	// Store DNS cache in the cache file.
	StoreDNS bool `json:"store_dns,omitempty"`
	// @since 1.15.0
	//
	// Size of the write buffer.
	//
	// `1MB` is used by default.
	BufferSize *byteformats.MemoryBytes `json:"buffer_size,omitempty"`
	// @since 1.15.0
	//
	// Interval for flushing the write buffer automatically.
	//
	// Disabled by default.
	FlushInterval badoption.Duration `json:"flush_interval,omitempty"`

	// Deprecated: replaced by store_dns
	// @failure: Deprecated in sing-box 1.14.0
	//
	//     `store_rdrc` is deprecated and will be removed in sing-box 1.16.0, check [Migration](/migration/#migrate-store-rdrc).
	//
	// Store rejected DNS response cache in the cache file
	//
	// The check results of [Legacy Address Filter Fields](/configuration/dns/rule/#legacy-address-filter-fields)
	// will be cached until expiration.
	StoreRDRC bool `json:"store_rdrc,omitempty" schema:"omit"`
	// Deprecated: replaced by store_dns
	// @failure: Deprecated in sing-box 1.14.0
	//
	//     `rdrc_timeout` is deprecated and will be removed in sing-box 1.16.0, check [Migration](/migration/#migrate-store-rdrc).
	//
	// Timeout of rejected DNS response cache.
	//
	// `7d` is used by default.
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
	// `https://github.com/MetaCubeX/Yacd-meta/archive/gh-pages.zip` will be used if empty.
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
	// @since 1.10.0
	//
	// CORS allowed origins, `*` will be used if empty.
	//
	// To access the Clash API on a private network from a public website, you must explicitly specify it in `access_control_allow_origin` instead of using `*`.
	AccessControlAllowOrigin badoption.Listable[string] `json:"access_control_allow_origin,omitempty"`
	// @since 1.10.0
	//
	// Allow access from private network.
	//
	// To access the Clash API on a private network from a public website, `access_control_allow_private_network` must be enabled.
	AccessControlAllowPrivateNetwork bool `json:"access_control_allow_private_network,omitempty"`

	// Deprecated: migrated to global cache file
	// @failure: Deprecated in sing-box 1.8.0
	//
	//     `cache_file` is deprecated in Clash API and migrated to `cache_file.enabled` and `cache_file.path`.
	//
	// Cache file path, `cache.db` will be used if empty.
	CacheFile string `json:"cache_file,omitempty" schema:"omit"`
	// Deprecated: migrated to global cache file
	// @failure: Deprecated in sing-box 1.8.0
	//
	//     `cache_id` is deprecated in Clash API and migrated to `cache_file.cache_id`.
	//
	// Identifier in cache file.
	//
	// If not empty, configuration specified data will use a separate store keyed by it.
	CacheID string `json:"cache_id,omitempty" schema:"omit"`
	// Deprecated: migrated to global cache file
	// @failure: Deprecated in sing-box 1.8.0
	//
	//     `store_mode` is deprecated in Clash API and enabled by default if `cache_file.enabled`.
	//
	// Store Clash mode in cache file.
	StoreMode bool `json:"store_mode,omitempty" schema:"omit"`
	// Deprecated: migrated to global cache file
	// @failure: Deprecated in sing-box 1.8.0
	//
	//     `store_selected` is deprecated in Clash API and enabled by default if `cache_file.enabled`.
	//
	// @note:
	//
	//     The tag must be set for target outbounds.
	//
	// Store selected outbound for the `Selector` outbound in cache file.
	StoreSelected bool `json:"store_selected,omitempty" schema:"omit"`
	// Deprecated: migrated to global cache file
	// @failure: Deprecated in sing-box 1.8.0
	//
	//     `store_selected` is deprecated in Clash API and migrated to `cache_file.store_fakeip`.
	//
	// Store fakeip in cache file.
	StoreFakeIP bool `json:"store_fakeip,omitempty" schema:"omit"`
}

type V2RayAPIOptions struct {
	// gRPC API listening address. V2Ray API will be disabled if empty.
	Listen string `json:"listen,omitempty"`
	// Traffic statistics service settings.
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
