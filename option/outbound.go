package option

import (
	"context"
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

type OutboundOptionsRegistry interface {
	OptionTypes() []string
	CreateOptions(outboundType string) (any, bool)
}

type _Outbound struct {
	Type    string `json:"type"`
	Tag     string `json:"tag,omitempty"`
	Options any    `json:"-"`
}

type Outbound _Outbound

func (h *Outbound) MarshalJSONContext(ctx context.Context) ([]byte, error) {
	return badjson.MarshallObjectsContext(ctx, (*_Outbound)(h), h.Options)
}

func (h *Outbound) UnmarshalJSONContext(ctx context.Context, content []byte) error {
	err := json.UnmarshalContext(ctx, content, (*_Outbound)(h))
	if err != nil {
		return err
	}
	registry := service.FromContext[OutboundOptionsRegistry](ctx)
	if registry == nil {
		return E.New("missing outbound options registry in context")
	}
	switch h.Type {
	case C.TypeDNS:
		return E.New("dns outbound is deprecated in sing-box 1.11.0 and removed in sing-box 1.13.0, use rule actions instead")
	}
	options, loaded := registry.CreateOptions(h.Type)
	if !loaded {
		return E.New("unknown outbound type: ", h.Type)
	}
	err = badjson.UnmarshallExcludedContext(ctx, content, (*_Outbound)(h), options)
	if err != nil {
		return err
	}
	if listenWrapper, isListen := options.(ListenOptionsWrapper); isListen {
		//nolint:staticcheck
		if listenWrapper.TakeListenOptions().InboundOptions != (InboundOptions{}) {
			return E.New("legacy inbound fields are deprecated in sing-box 1.11.0 and removed in sing-box 1.13.0, use rule actions instead")
		}
	}
	h.Options = options
	return nil
}

func (h Outbound) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return builder.Define("Outbound", func() (*schema.Node, error) {
		registry := service.FromContext[OutboundOptionsRegistry](builder.Context())
		if registry == nil {
			return nil, E.New("missing outbound options registry in context")
		}
		return registryUnion(builder, registry, []string{C.TypeShadowsocksR, C.TypeWireGuard}, true)
	})
}

type DialerOptionsWrapper interface {
	TakeDialerOptions() DialerOptions
	ReplaceDialerOptions(options DialerOptions)
}

type DialerOptions struct {
	// The tag of the upstream outbound.
	//
	// If enabled, all other fields will be ignored.
	Detour string `json:"detour,omitempty" reference:"outbound"`
	AbstractDialerOptions
}

type AbstractDialerOptions struct {
	// The network interface to bind to.
	BindInterface string `json:"bind_interface,omitempty"`
	// The IPv4 address to bind to.
	Inet4BindAddress *badoption.Addr `json:"inet4_bind_address,omitempty"`
	// The IPv6 address to bind to.
	Inet6BindAddress *badoption.Addr `json:"inet6_bind_address,omitempty"`
	// @since 1.13.0
	//
	// @quote:
	//
	//     Only supported on Linux.
	//
	// Do not reserve a port when binding to a source address.
	//
	// This allows reusing the same source port for multiple connections if the full 4-tuple (source IP, source port, destination IP, destination port) remains unique.
	BindAddressNoPort bool   `json:"bind_address_no_port,omitempty"`
	ProtectPath       string `json:"protect_path,omitempty"`
	// @quote:
	//
	//     Only supported on Linux.
	//
	// Set netfilter routing mark.
	//
	// Integers (e.g. `1234`) and string hexadecimals (e.g. `"0x1234"`) are supported.
	RoutingMark FwMark `json:"routing_mark,omitempty"`
	// Reuse listener address.
	ReuseAddr bool `json:"reuse_addr,omitempty"`
	// @since 1.12.0
	//
	// @quote:
	//
	//     Only supported on Linux.
	//
	// Set network namespace, name or path.
	//
	// Since sing-box 1.14.0, the tag of a [network namespace](/configuration/network-namespace/) can also be used.
	// Referencing an `unshare` network namespace should be avoided.
	NetNs string `json:"netns,omitempty" reference:"network_namespace"`
	// Connect timeout, in golang's Duration format.
	//
	// A duration string is a possibly signed sequence of
	// decimal numbers, each with optional fraction and a unit suffix,
	// such as "300ms", "-1.5h" or "2h45m".
	// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	ConnectTimeout badoption.Duration `json:"connect_timeout,omitempty"`
	// Enable TCP Fast Open.
	TCPFastOpen bool `json:"tcp_fast_open,omitempty"`
	// @warning:
	//
	//     Go 1.21 required.
	//
	// Enable TCP Multi Path.
	TCPMultiPath bool `json:"tcp_multi_path,omitempty"`
	// @since 1.13.0
	//
	// Disable TCP keep alive.
	DisableTCPKeepAlive bool `json:"disable_tcp_keep_alive,omitempty"`
	// @since 1.13.0
	//
	//     Default value changed from `10m` to `5m`.
	//
	// TCP keep alive initial period.
	//
	// `5m` will be used by default.
	TCPKeepAlive badoption.Duration `json:"tcp_keep_alive,omitempty"`
	// @since 1.13.0
	//
	// TCP keep alive interval.
	//
	// `75s` will be used by default.
	TCPKeepAliveInterval       badoption.Duration `json:"tcp_keep_alive_interval,omitempty"`
	TCPKeepAliveSystemDefaults bool               `json:"-"`
	UDPBindPort                uint16             `json:"-"`
	// Enable UDP fragmentation.
	UDPFragment        *bool `json:"udp_fragment,omitempty"`
	UDPFragmentDefault bool  `json:"-"`
	// @warning:
	//
	//     `outbound` DNS rule items are deprecated and will be removed in sing-box 1.14.0, so this item will be required for outbound/endpoints using domain name in server address since sing-box 1.14.0.
	//
	// @info:
	//
	//     `domain_resolver` or `route.default_domain_resolver` is optional when only one DNS server is configured.
	//
	// Set domain resolver to use for resolving domain names.
	//
	// This option uses the same format as the [route DNS rule action](/configuration/dns/rule_action/#route) without the `action` field.
	//
	// Setting this option directly to a string is equivalent to setting `server` of this options.
	//
	// | Outbound/Endpoints | Effected domains         |
	// |--------------------|--------------------------|
	// | `direct`           | Domain in request        |
	// | others             | Domain in server address |
	DomainResolver *DomainResolveOptions `json:"domain_resolver,omitempty"`
	// @since 1.11.0
	//
	// @quote:
	//
	//     Only supported in graphical clients on Android and Apple platforms with `auto_detect_interface` enabled.
	//
	// Strategy for selecting network interfaces.
	//
	// Available values:
	//
	// - `default` (default): Connect to default network or networks specified in `network_type` sequentially.
	// - `hybrid`: Connect to all networks or networks specified in `network_type` concurrently.
	// - `fallback`: Connect to default network or preferred networks specified in `network_type` concurrently, and try fallback networks when unavailable or timeout.
	//
	// For fallback, when preferred interfaces fails or times out,
	// it will enter a 15s fast fallback state (Connect to all preferred and fallback networks concurrently),
	// and exit immediately if preferred networks recover.
	//
	// Conflicts with `bind_interface`, `inet4_bind_address` and `inet6_bind_address`.
	NetworkStrategy *NetworkStrategy `json:"network_strategy,omitempty"`
	// @since 1.11.0
	//
	// @quote:
	//
	//     Only supported in graphical clients on Android and Apple platforms with `auto_detect_interface` enabled.
	//
	// Network types to use when using `default` or `hybrid` network strategy or
	// preferred network types to use when using `fallback` network strategy.
	//
	// Available values: `wifi`, `cellular`, `ethernet`, `other`.
	//
	// Device's default network is used by default.
	NetworkType badoption.Listable[InterfaceType] `json:"network_type,omitempty"`
	// @since 1.11.0
	//
	// @quote:
	//
	//     Only supported in graphical clients on Android and Apple platforms with `auto_detect_interface` enabled.
	//
	// Fallback network types when preferred networks are unavailable or timeout when using `fallback` network strategy.
	//
	// All other networks expect preferred are used by default.
	FallbackNetworkType badoption.Listable[InterfaceType] `json:"fallback_network_type,omitempty"`
	// @since 1.11.0
	//
	// @quote:
	//
	//     Only supported in graphical clients on Android and Apple platforms with `auto_detect_interface` enabled.
	//
	// The length of time to wait before spawning a RFC 6555 Fast Fallback connection.
	//
	// For `domain_strategy`, is the amount of time to wait for connection to succeed before assuming
	// that IPv4/IPv6 is misconfigured and falling back to other type of addresses.
	//
	// For `network_strategy`, is the amount of time to wait for connection to succeed before falling
	// back to other interfaces.
	//
	// Only take effect when `domain_strategy` or `network_strategy` is set.
	//
	// `300ms` is used by default.
	FallbackDelay badoption.Duration `json:"fallback_delay,omitempty"`

	// Deprecated: migrated to domain resolver
	// @failure: Deprecated in sing-box 1.12.0
	//
	//     `domain_strategy` is deprecated and will be removed in sing-box 1.14.0, check [Migration](/migration/#migrate-outbound-domain-strategy-option-to-domain-resolver).
	//
	// Available values: `prefer_ipv4`, `prefer_ipv6`, `ipv4_only`, `ipv6_only`.
	//
	// If set, the requested domain name will be resolved to IP before connect.
	//
	// | Outbound | Effected domains         | Fallback Value                            |
	// |----------|--------------------------|-------------------------------------------|
	// | `direct` | Domain in request        | Take `inbound.domain_strategy` if not set |
	// | others   | Domain in server address | /                                         |
	DomainStrategy DomainStrategy `json:"domain_strategy,omitempty" schema:"omit"`
}

type _DomainResolveOptions struct {
	Server                 string                `json:"server" reference:"dns_server"`
	Timeout                badoption.Duration    `json:"timeout,omitempty"`
	Strategy               DomainStrategy        `json:"strategy,omitempty"`
	DisableCache           bool                  `json:"disable_cache,omitempty"`
	DisableOptimisticCache bool                  `json:"disable_optimistic_cache,omitempty"`
	RewriteTTL             *uint32               `json:"rewrite_ttl,omitempty"`
	ClientSubnet           *badoption.Prefixable `json:"client_subnet,omitempty"`
}

type DomainResolveOptions _DomainResolveOptions

func (o DomainResolveOptions) MarshalJSON() ([]byte, error) {
	if o.Server == "" {
		return []byte("{}"), nil
	} else if o.Strategy == DomainStrategy(C.DomainStrategyAsIS) &&
		o.Timeout == 0 &&
		!o.DisableCache &&
		!o.DisableOptimisticCache &&
		o.RewriteTTL == nil &&
		o.ClientSubnet == nil {
		return json.Marshal(o.Server)
	} else {
		return json.Marshal(_DomainResolveOptions(o))
	}
}

func (o *DomainResolveOptions) UnmarshalJSON(bytes []byte) error {
	var stringValue string
	err := json.Unmarshal(bytes, &stringValue)
	if err == nil {
		o.Server = stringValue
		return nil
	}
	err = json.Unmarshal(bytes, (*_DomainResolveOptions)(o))
	if err != nil {
		return err
	}
	if o.Server == "" {
		return E.New("empty domain_resolver.server")
	}
	return nil
}

func (o DomainResolveOptions) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return builder.Define("DomainResolver", func() (*schema.Node, error) {
		objectForm := schema.StrictObject()
		err := builder.FlattenStruct(objectForm, reflect.TypeFor[DomainResolveOptions]())
		if err != nil {
			return nil, err
		}
		objectForm.Required = []string{"server"}
		return schema.AnyOf(schema.TagReferenceNode("dns_server"), objectForm), nil
	})
}

func (o *DialerOptions) TakeDialerOptions() DialerOptions {
	return *o
}

func (o *DialerOptions) ReplaceDialerOptions(options DialerOptions) {
	*o = options
}

type ServerOptionsWrapper interface {
	TakeServerOptions() ServerOptions
	ReplaceServerOptions(options ServerOptions)
}

type ServerOptions struct {
	// @required
	//
	// The server address.
	Server string `json:"server"`
	// @required
	//
	// The server port.
	ServerPort uint16 `json:"server_port"`
}

func (o ServerOptions) Build() M.Socksaddr {
	return M.ParseSocksaddrHostPort(o.Server, o.ServerPort)
}

func (o ServerOptions) ServerIsDomain() bool {
	return o.Build().IsDomain()
}

func (o *ServerOptions) TakeServerOptions() ServerOptions {
	return *o
}

func (o *ServerOptions) ReplaceServerOptions(options ServerOptions) {
	*o = options
}
