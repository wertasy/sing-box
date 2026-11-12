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

// @Deprecated: Use rule action instead
type InboundOptions struct {
	// Enable sniffing.
	//
	// See [Protocol Sniff](/configuration/route/sniff/) for details.
	//
	// @Failure: Deprecated in sing-box 1.11.0
	//
	//     Inbound fields are deprecated and will be removed in sing-box 1.13.0, check [Migration](/migration/#migrate-legacy-inbound-fields-to-rule-actions).
	SniffEnabled bool `json:"sniff,omitempty" schema:"omit"`
	// Override the connection destination address with the sniffed domain.
	//
	// If the domain name is invalid (like tor), this will not work.
	//
	// @Failure: Deprecated in sing-box 1.11.0
	//
	//     Inbound fields are deprecated and will be removed in sing-box 1.13.0.
	SniffOverrideDestination bool `json:"sniff_override_destination,omitempty" schema:"omit"`
	// Timeout for sniffing.
	//
	// @Failure: Deprecated in sing-box 1.11.0
	//
	//     Inbound fields are deprecated and will be removed in sing-box 1.13.0, check [Migration](/migration/#migrate-legacy-inbound-fields-to-rule-actions).
	//
	// @Default: 300ms
	SniffTimeout badoption.Duration `json:"sniff_timeout,omitempty" schema:"omit"`
	// One of `prefer_ipv4` `prefer_ipv6` `ipv4_only` `ipv6_only`.
	//
	// If set, the requested domain name will be resolved to IP before routing.
	//
	// If `sniff_override_destination` is in effect, its value will be taken as a fallback.
	//
	// @Failure: Deprecated in sing-box 1.11.0
	//
	//     Inbound fields are deprecated and will be removed in sing-box 1.13.0, check [Migration](/migration/#migrate-legacy-inbound-fields-to-rule-actions).
	DomainStrategy DomainStrategy `json:"domain_strategy,omitempty" schema:"omit"`
	// If enabled, for UDP proxy requests addressed to a domain,
	// the original packet address will be sent in the response instead of the mapped domain.
	//
	// This option is used for compatibility with clients that
	// do not support receiving UDP packets with domain addresses, such as Surge.
	//
	// @Failure: Deprecated in sing-box 1.11.0
	//
	//     Inbound fields are deprecated and will be removed in sing-box 1.13.0, check [Migration](/migration/#migrate-legacy-inbound-fields-to-rule-actions).
	UDPDisableDomainUnmapping bool `json:"udp_disable_domain_unmapping,omitempty" schema:"omit"`
}

// @Shared listen
type ListenOptions struct {
	// Listen address.
	//
	// @Required
	Listen *badoption.Addr `json:"listen,omitempty"`
	// Listen port.
	//
	ListenPort uint16 `json:"listen_port,omitempty"`
	// The network interface to bind to.
	//
	// @Since 1.12.0
	BindInterface string `json:"bind_interface,omitempty"`
	// Set netfilter routing mark.
	//
	// Integers (e.g. `1234`) and string hexadecimals (e.g. `"0x1234"`) are supported.
	//
	// @Since 1.12.0
	//
	// @Quote:
	//
	//     Only supported on Linux.
	RoutingMark FwMark `json:"routing_mark,omitempty"`
	// Reuse listener address.
	//
	// @Since 1.12.0
	ReuseAddr bool `json:"reuse_addr,omitempty"`
	// Set network namespace, name or path.
	//
	// Since sing-box 1.14.0, the tag of a [network namespace](/configuration/network-namespace/) can also be used.
	//
	// @Since 1.12.0
	//
	// @Quote:
	//
	//     Only supported on Linux.
	NetNs string `json:"netns,omitempty" reference:"network_namespace"`
	// Disable TCP keep alive.
	//
	// @Since 1.13.0
	DisableTCPKeepAlive bool `json:"disable_tcp_keep_alive,omitempty"`
	// TCP keep alive initial period.
	//
	// @Default: 5m
	//
	// @Since 1.13.0
	//
	//     Default value changed from `10m` to `5m`.
	TCPKeepAlive badoption.Duration `json:"tcp_keep_alive,omitempty"`
	// TCP keep alive interval.
	//
	// @Default: 75s
	TCPKeepAliveInterval badoption.Duration `json:"tcp_keep_alive_interval,omitempty"`
	// Enable TCP Fast Open.
	TCPFastOpen bool `json:"tcp_fast_open,omitempty"`
	// Enable TCP Multi Path.
	//
	// @Warning:
	//
	//     Go 1.21 required.
	TCPMultiPath bool `json:"tcp_multi_path,omitempty"`
	// Enable UDP fragmentation.
	UDPFragment        *bool `json:"udp_fragment,omitempty"`
	UDPFragmentDefault bool  `json:"-"`
	// UDP NAT expiration time.
	//
	// @Default: 5m
	UDPTimeout UDPTimeoutCompat `json:"udp_timeout,omitempty"`
	// If set, connections will be forwarded to the specified inbound.
	//
	// Requires target inbound support, see [Injectable](/configuration/inbound/#fields).
	Detour string `json:"detour,omitempty" reference:"inbound"`

	// @Deprecated: removed
	ProxyProtocol bool `json:"proxy_protocol,omitempty" schema:"omit"`
	// @Deprecated: removed
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
