package option

import (
	"context"
	"net/netip"
	"reflect"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/schema"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badjson"
	"github.com/sagernet/sing/common/json/badoption"
	M "github.com/sagernet/sing/common/metadata"
	"github.com/sagernet/sing/service"
)

type RawDNSOptions struct {
	// @Omit
	Servers []DNSServerOptions `json:"servers,omitempty"`
	// @Omit
	Rules []DNSRule `json:"rules,omitempty"`
	// Default dns server tag.
	//
	// The first server will be used if empty.
	Final string `json:"final,omitempty" reference:"dns_server"`
	// Stores a reverse mapping of IP addresses after responding to a DNS query in order to provide domain names when routing.
	//
	// Since this process relies on the act of resolving domain names by an application before making a request, it can be
	// problematic in environments such as macOS, where DNS is proxied and cached by the system.
	ReverseMapping bool `json:"reverse_mapping,omitempty"`
	DNSClientOptions
}

type DNSOptions struct {
	RawDNSOptions
}

const (
	legacyDNSFakeIPRemovedMessage = "legacy DNS fakeip options are deprecated in sing-box 1.12.0 and removed in sing-box 1.14.0, checkout migration: https://sing-box.sagernet.org/migration/#migrate-to-new-dns-server-formats"
	legacyDNSServerRemovedMessage = "legacy DNS server formats are deprecated in sing-box 1.12.0 and removed in sing-box 1.14.0, checkout migration: https://sing-box.sagernet.org/migration/#migrate-to-new-dns-server-formats"
)

type removedLegacyDNSOptions struct {
	FakeIP json.RawMessage `json:"fakeip,omitempty"`
}

func (o *DNSOptions) UnmarshalJSONContext(ctx context.Context, content []byte) error {
	var legacyOptions removedLegacyDNSOptions
	err := json.UnmarshalContext(ctx, content, &legacyOptions)
	if err != nil {
		return err
	}
	if len(legacyOptions.FakeIP) != 0 {
		return E.New(legacyDNSFakeIPRemovedMessage)
	}
	return badjson.UnmarshallExcludedContext(ctx, content, legacyOptions, &o.RawDNSOptions)
}

func (o DNSOptions) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return builder.Define("DNS", func() (*schema.Node, error) {
		node := schema.StrictObject()
		err := builder.FlattenStruct(node, reflect.TypeFor[RawDNSOptions]())
		if err != nil {
			return nil, err
		}
		return node, nil
	})
}

type DNSClientOptions struct {
	// Default domain strategy for resolving the domain names.
	//
	// One of `prefer_ipv4` `prefer_ipv6` `ipv4_only` `ipv6_only`.
	Strategy DomainStrategy `json:"strategy,omitempty"`
	// Default timeout for each DNS query.
	//
	// Can be overridden by `rules.[].timeout` (DNS rule action) or `domain_resolver.timeout`.
	//
	// @Default: 10s
	//
	// @Since 1.14.0
	Timeout badoption.Duration `json:"timeout,omitempty"`
	// Disable dns cache.
	//
	// Conflict with `optimistic`.
	DisableCache bool `json:"disable_cache,omitempty"`
	// Disable dns cache expire.
	//
	// Conflict with `optimistic`.
	DisableExpire bool `json:"disable_expire,omitempty"`
	// Make each DNS server's cache independent for special purposes. If enabled, will slightly degrade performance.
	//
	// @Failure: Deprecated in sing-box 1.14.0
	//
	//     `independent_cache` is deprecated and will be removed in sing-box 1.16.0, check [Migration](/migration/#migrate-independent-dns-cache).
	IndependentCache bool `json:"independent_cache,omitempty" schema:"omit"`
	// LRU cache capacity.
	//
	// Value less than 1024 will be ignored.
	//
	// @Since 1.11.0
	CacheCapacity uint32 `json:"cache_capacity,omitempty"`
	// Enable optimistic DNS caching. When a cached DNS entry has expired but is still within the timeout window,
	// the stale response is returned immediately while a background refresh is triggered.
	//
	// Conflict with `disable_cache` and `disable_expire`.
	//
	// Accepts a boolean or an object. When set to `true`, the default timeout of `3d` is used.
	//
	// ```json
	// {
	//   "enabled": true,
	//   "timeout": "3d"
	// }
	// ```
	//
	// ##### enabled
	//
	// Enable optimistic DNS caching.
	//
	// ##### timeout
	//
	// The maximum time an expired cache entry can be served optimistically.
	//
	// @Default: 3d
	//
	// @Since 1.14.0
	Optimistic *OptimisticDNSOptions `json:"optimistic,omitempty"`
	// Append a `edns0-subnet` OPT extra record with the specified IP prefix to every query by default.
	//
	// If value is an IP address instead of prefix, `/32` or `/128` will be appended automatically.
	//
	// Can be overridden by `servers.[].client_subnet` or `rules.[].client_subnet`.
	//
	// @Since 1.9.0
	ClientSubnet *badoption.Prefixable `json:"client_subnet,omitempty"`
}

type _OptimisticDNSOptions struct {
	Enabled bool               `json:"enabled,omitempty"`
	Timeout badoption.Duration `json:"timeout,omitempty"`
}

type OptimisticDNSOptions _OptimisticDNSOptions

func (o OptimisticDNSOptions) MarshalJSON() ([]byte, error) {
	if o.Timeout == 0 {
		return json.Marshal(o.Enabled)
	}
	return json.Marshal(_OptimisticDNSOptions(o))
}

func (o *OptimisticDNSOptions) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, &o.Enabled)
	if err == nil {
		return nil
	}
	return json.UnmarshalDisallowUnknownFields(bytes, (*_OptimisticDNSOptions)(o))
}

func (o OptimisticDNSOptions) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	objectForm := schema.StrictObject()
	err := builder.FlattenStruct(objectForm, reflect.TypeFor[OptimisticDNSOptions]())
	if err != nil {
		return nil, err
	}
	return schema.AnyOf(schema.BooleanNode(), objectForm), nil
}

type DNSTransportOptionsRegistry interface {
	OptionTypes() []string
	CreateOptions(transportType string) (any, bool)
}
type _DNSServerOptions struct {
	Type    string `json:"type,omitempty"`
	Tag     string `json:"tag,omitempty"`
	Options any    `json:"-"`
}

type DNSServerOptions _DNSServerOptions

func (o *DNSServerOptions) MarshalJSONContext(ctx context.Context) ([]byte, error) {
	return badjson.MarshallObjectsContext(ctx, (*_DNSServerOptions)(o), o.Options)
}

func (o *DNSServerOptions) UnmarshalJSONContext(ctx context.Context, content []byte) error {
	err := json.UnmarshalContext(ctx, content, (*_DNSServerOptions)(o))
	if err != nil {
		return err
	}
	registry := service.FromContext[DNSTransportOptionsRegistry](ctx)
	if registry == nil {
		return E.New("missing DNS transport options registry in context")
	}
	var options any
	switch o.Type {
	case "", C.DNSTypeLegacy:
		return E.New(legacyDNSServerRemovedMessage)
	default:
		var loaded bool
		options, loaded = registry.CreateOptions(o.Type)
		if !loaded {
			return E.New("unknown transport type: ", o.Type)
		}
	}
	err = badjson.UnmarshallExcludedContext(ctx, content, (*_DNSServerOptions)(o), options)
	if err != nil {
		return err
	}
	o.Options = options
	return nil
}

func (o DNSServerOptions) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return builder.Define("DNSServer", func() (*schema.Node, error) {
		registry := service.FromContext[DNSTransportOptionsRegistry](builder.Context())
		if registry == nil {
			return nil, E.New("missing DNS transport options registry in context")
		}
		return registryUnion(builder, registry, nil, true)
	})
}

type DNSServerAddressOptions struct {
	// The address of the DNS server.
	//
	// If domain name is used, `domain_resolver` must also be set to resolve IP address.
	//
	// @Required
	Server string `json:"server"`
	// The port of the DNS server.
	//
	// @Default: 53
	ServerPort uint16 `json:"server_port,omitempty"`
}

func (o DNSServerAddressOptions) Build() M.Socksaddr {
	return M.ParseSocksaddrHostPort(o.Server, o.ServerPort)
}

func (o DNSServerAddressOptions) ServerIsDomain() bool {
	return o.Build().IsDomain()
}

func (o *DNSServerAddressOptions) TakeServerOptions() ServerOptions {
	return ServerOptions(*o)
}

func (o *DNSServerAddressOptions) ReplaceServerOptions(options ServerOptions) {
	*o = DNSServerAddressOptions(options)
}

type HostsDNSServerOptions struct {
	// List of paths to hosts files.
	//
	// `/etc/hosts` is used by default.
	//
	// `C:\Windows\System32\Drivers\etc\hosts` is used by default on Windows.
	//
	// Example:
	//
	// ```json
	// {
	//   // "path": "/etc/hosts"
	//
	//   "path": [
	//     "/etc/hosts",
	//     "$HOME/.hosts"
	//   ]
	// }
	// ```
	Path badoption.Listable[string] `json:"path,omitempty"`
	// Predefined hosts.
	//
	// Example:
	//
	// ```json
	// {
	//   "predefined": {
	//     "www.google.com": "127.0.0.1",
	//     "localhost": [
	//       "127.0.0.1",
	//       "::1"
	//     ]
	//   }
	// }
	// ```
	Predefined *badjson.TypedMap[string, badoption.Listable[netip.Addr]] `json:"predefined,omitempty"`
}

type RawLocalDNSServerOptions struct {
	DialerOptions
}

type LocalDNSServerOptions struct {
	RawLocalDNSServerOptions
	// When enabled, `local` DNS server will resolve DNS by dialing itself whenever possible.
	//
	// Specifically, it disables following behaviors which was added as features in sing-box 1.13.0:
	//
	// 1. On Apple platforms: Attempt to resolve A/AAAA requests using `getaddrinfo` in NetworkExtension.
	// 2. On Linux: Resolve through `systemd-resolvd`'s DBus interface when available.
	//
	// As a sole exception, it cannot disable the following behavior:
	//
	// 1. In the Android graphical client,
	// `local` will always resolve DNS through the platform interface,
	// as there is no other way to obtain upstream DNS servers;
	// On devices running Android versions lower than 10, this interface can only resolve A/AAAA requests.
	//
	// 2. On macOS, `local` will try DHCP first in Network Extension, since DHCP respects DIal Fields,
	// it will not be disabled by `prefer_go`.
	//
	// @Since 1.13.0
	PreferGo bool `json:"prefer_go,omitempty"`
	// A list of domain suffixes for which A/AAAA queries are answered from the
	// [neighbor resolver](/configuration/shared/neighbor/) instead of the upstream.
	//
	// Each entry must start with `.`. Only queries whose host part (the portion
	// before the suffix) contains no dots are matched; `.` matches any
	// single-label name such as `nas`.
	//
	// Example: `[".", ".lan"]`.
	//
	// @Since 1.14.0
	NeighborDomain badoption.Listable[string] `json:"neighbor_domain,omitempty"`
}

type RemoteDNSServerOptions struct {
	RawLocalDNSServerOptions
	DNSServerAddressOptions
}

type RemoteTLSDNSServerOptions struct {
	RemoteDNSServerOptions
	OutboundTLSOptionsContainer
}

type RemoteHTTPSDNSServerOptions struct {
	RemoteTLSDNSServerOptions
	// The path of the DNS server.
	//
	// `/dns-query` will be used by default.
	Path string `json:"path,omitempty"`
	// @Omit
	Method string `json:"method,omitempty"`
	// Additional headers to be sent to the DNS server.
	Headers badoption.HTTPHeader `json:"headers,omitempty"`
}

type FakeIPDNSServerOptions struct {
	// IPv4 address range for FakeIP.
	Inet4Range *badoption.Prefix `json:"inet4_range,omitempty" examples:"198.18.0.0/15"`
	// IPv6 address range for FakeIP.
	Inet6Range *badoption.Prefix `json:"inet6_range,omitempty" examples:"2001:2::/48"`
}

type DHCPDNSServerOptions struct {
	LocalDNSServerOptions
	// Interface name to listen on.
	//
	// Tge default interface will be used by default.
	Interface string `json:"interface,omitempty"`
}

type MDNSDNSServerOptions struct {
	LocalDNSServerOptions
	// List of network interface names to send mDNS queries on.
	//
	// When empty, all interfaces that are up, multicast-capable, and non-loopback are used.
	Interface badoption.Listable[string] `json:"interface,omitempty"`
}
