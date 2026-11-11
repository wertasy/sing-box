package option

import (
	"context"
	"time"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/schema"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badjson"
	"github.com/sagernet/sing/common/json/badoption"
	"github.com/sagernet/sing/service"
)

type InboundOptionsRegistry interface {
	OptionTypes() []string
	CreateOptions(outboundType string) (any, bool)
}

type _Inbound struct {
	Type    string `json:"type"`
	Tag     string `json:"tag,omitempty"`
	Options any    `json:"-"`
}

type Inbound _Inbound

func (h *Inbound) MarshalJSONContext(ctx context.Context) ([]byte, error) {
	return badjson.MarshallObjectsContext(ctx, (*_Inbound)(h), h.Options)
}

func (h *Inbound) UnmarshalJSONContext(ctx context.Context, content []byte) error {
	err := json.UnmarshalContext(ctx, content, (*_Inbound)(h))
	if err != nil {
		return err
	}
	registry := service.FromContext[InboundOptionsRegistry](ctx)
	if registry == nil {
		return E.New("missing inbound fields registry in context")
	}
	options, loaded := registry.CreateOptions(h.Type)
	if !loaded {
		return E.New("unknown inbound type: ", h.Type)
	}
	err = badjson.UnmarshallExcludedContext(ctx, content, (*_Inbound)(h), options)
	if err != nil {
		return err
	}
	if listenWrapper, isListen := options.(ListenOptionsWrapper); isListen {
		//nolint:staticcheck
		if listenWrapper.TakeListenOptions().InboundOptions != (InboundOptions{}) {
			return E.New("legacy inbound fields are deprecated in sing-box 1.11.0 and removed in sing-box 1.13.0, checkout migration: https://sing-box.sagernet.org/migration/#migrate-legacy-inbound-fields-to-rule-actions")
		}
	}
	h.Options = options
	return nil
}

func (h Inbound) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return builder.Define("Inbound", func() (*schema.Node, error) {
		registry := service.FromContext[InboundOptionsRegistry](builder.Context())
		if registry == nil {
			return nil, E.New("missing inbound options registry in context")
		}
		return registryUnion(builder, registry, []string{C.TypeShadowsocksR}, true)
	})
}

// Deprecated: Use rule action instead
type InboundOptions struct {
	// @failure: Deprecated in sing-box 1.11.0
	//
	//     Inbound fields are deprecated and will be removed in sing-box 1.13.0, check [Migration](/migration/#migrate-legacy-inbound-fields-to-rule-actions).
	//
	// Enable sniffing.
	//
	// See [Protocol Sniff](/configuration/route/sniff/) for details.
	SniffEnabled bool `json:"sniff,omitempty" schema:"omit"`
	// @failure: Deprecated in sing-box 1.11.0
	//
	//     Inbound fields are deprecated and will be removed in sing-box 1.13.0.
	//
	// Override the connection destination address with the sniffed domain.
	//
	// If the domain name is invalid (like tor), this will not work.
	SniffOverrideDestination bool `json:"sniff_override_destination,omitempty" schema:"omit"`
	// @failure: Deprecated in sing-box 1.11.0
	//
	//     Inbound fields are deprecated and will be removed in sing-box 1.13.0, check [Migration](/migration/#migrate-legacy-inbound-fields-to-rule-actions).
	//
	// Timeout for sniffing.
	//
	// `300ms` is used by default.
	SniffTimeout badoption.Duration `json:"sniff_timeout,omitempty" schema:"omit"`
	// @failure: Deprecated in sing-box 1.11.0
	//
	//     Inbound fields are deprecated and will be removed in sing-box 1.13.0, check [Migration](/migration/#migrate-legacy-inbound-fields-to-rule-actions).
	//
	// One of `prefer_ipv4` `prefer_ipv6` `ipv4_only` `ipv6_only`.
	//
	// If set, the requested domain name will be resolved to IP before routing.
	//
	// If `sniff_override_destination` is in effect, its value will be taken as a fallback.
	DomainStrategy DomainStrategy `json:"domain_strategy,omitempty" schema:"omit"`
	// @failure: Deprecated in sing-box 1.11.0
	//
	//     Inbound fields are deprecated and will be removed in sing-box 1.13.0, check [Migration](/migration/#migrate-legacy-inbound-fields-to-rule-actions).
	//
	// If enabled, for UDP proxy requests addressed to a domain,
	// the original packet address will be sent in the response instead of the mapped domain.
	//
	// This option is used for compatibility with clients that
	// do not support receiving UDP packets with domain addresses, such as Surge.
	UDPDisableDomainUnmapping bool `json:"udp_disable_domain_unmapping,omitempty" schema:"omit"`
}

type ListenOptions struct {
	// @required
	//
	// Listen address.
	Listen *badoption.Addr `json:"listen,omitempty"`
	// Listen port.
	ListenPort uint16 `json:"listen_port,omitempty"`
	// @since 1.12.0
	//
	// The network interface to bind to.
	BindInterface string `json:"bind_interface,omitempty"`
	// @since 1.12.0
	//
	// @quote:
	//
	//     Only supported on Linux.
	//
	// Set netfilter routing mark.
	//
	// Integers (e.g. `1234`) and string hexadecimals (e.g. `"0x1234"`) are supported.
	RoutingMark FwMark `json:"routing_mark,omitempty"`
	// @since 1.12.0
	//
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
	NetNs string `json:"netns,omitempty" reference:"network_namespace"`
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
	// TCP keep alive interval.
	//
	// `75s` will be used by default.
	TCPKeepAliveInterval badoption.Duration `json:"tcp_keep_alive_interval,omitempty"`
	// Enable TCP Fast Open.
	TCPFastOpen bool `json:"tcp_fast_open,omitempty"`
	// @warning:
	//
	//     Go 1.21 required.
	//
	// Enable TCP Multi Path.
	TCPMultiPath bool `json:"tcp_multi_path,omitempty"`
	// Enable UDP fragmentation.
	UDPFragment        *bool `json:"udp_fragment,omitempty"`
	UDPFragmentDefault bool  `json:"-"`
	// UDP NAT expiration time.
	//
	// `5m` will be used by default.
	UDPTimeout UDPTimeoutCompat `json:"udp_timeout,omitempty"`
	// If set, connections will be forwarded to the specified inbound.
	//
	// Requires target inbound support, see [Injectable](/configuration/inbound/#fields).
	Detour string `json:"detour,omitempty" reference:"inbound"`

	// Deprecated: removed
	ProxyProtocol bool `json:"proxy_protocol,omitempty" schema:"omit"`
	// Deprecated: removed
	ProxyProtocolAcceptNoHeader bool `json:"proxy_protocol_accept_no_header,omitempty" schema:"omit"`
	// Legacy inbound fields are rejected since sing-box 1.13.0.
	//nolint:staticcheck
	InboundOptions `schema:"omit"`
}

type UDPNATBehavior uint8

const (
	UDPNATBehaviorEndpointIndependent UDPNATBehavior = iota
	UDPNATBehaviorAddressDependent
	UDPNATBehaviorAddressAndPortDependent
)

func (b UDPNATBehavior) MarshalJSON() ([]byte, error) {
	var value string
	switch b {
	case UDPNATBehaviorEndpointIndependent:
		value = "endpoint_independent"
	case UDPNATBehaviorAddressDependent:
		value = "address_dependent"
	case UDPNATBehaviorAddressAndPortDependent:
		value = "address_and_port_dependent"
	default:
		return nil, E.New("unknown UDP NAT behavior: ", uint8(b))
	}
	return json.Marshal(value)
}

func (b *UDPNATBehavior) UnmarshalJSON(data []byte) error {
	var value string
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}
	switch value {
	case "", "endpoint_independent":
		*b = UDPNATBehaviorEndpointIndependent
	case "address_dependent":
		*b = UDPNATBehaviorAddressDependent
	case "address_and_port_dependent":
		*b = UDPNATBehaviorAddressAndPortDependent
	default:
		return E.New("unknown UDP NAT behavior: ", value)
	}
	return nil
}

func (b UDPNATBehavior) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return schema.StringEnum("", "endpoint_independent", "address_dependent", "address_and_port_dependent"), nil
}

type UDPTimeoutCompat badoption.Duration

func (c UDPTimeoutCompat) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(c).String())
}

func (c *UDPTimeoutCompat) UnmarshalJSON(data []byte) error {
	var valueNumber int64
	err := json.Unmarshal(data, &valueNumber)
	if err == nil {
		*c = UDPTimeoutCompat(time.Second * time.Duration(valueNumber))
		return nil
	}
	return json.Unmarshal(data, (*badoption.Duration)(c))
}

func (c UDPTimeoutCompat) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return builder.Define("UDPTimeout", func() (*schema.Node, error) {
		return schema.AnyOf(schema.UnsignedNode(32), schema.DurationNode()), nil
	})
}

type ListenOptionsWrapper interface {
	TakeListenOptions() ListenOptions
	ReplaceListenOptions(options ListenOptions)
}

func (o *ListenOptions) TakeListenOptions() ListenOptions {
	return *o
}

func (o *ListenOptions) ReplaceListenOptions(options ListenOptions) {
	*o = options
}
