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
)

type _RuleAction struct {
	Action              string                    `json:"action,omitempty" enum:"route,route-options,bypass,reject,hijack-dns,sniff,resolve"`
	RouteOptions        RouteActionOptions        `json:"-"`
	RouteOptionsOptions RouteOptionsActionOptions `json:"-"`
	BypassOptions       RouteActionOptions        `json:"-"`
	RejectOptions       RejectActionOptions       `json:"-"`
	SniffOptions        RouteActionSniff          `json:"-"`
	ResolveOptions      RouteActionResolve        `json:"-"`
}

type RuleAction _RuleAction

func (r RuleAction) MarshalJSON() ([]byte, error) {
	if r.Action == "" {
		return json.Marshal(struct{}{})
	}
	var v any
	switch r.Action {
	case C.RuleActionTypeRoute:
		r.Action = ""
		v = r.RouteOptions
	case C.RuleActionTypeRouteOptions:
		v = r.RouteOptionsOptions
	case C.RuleActionTypeBypass:
		v = r.BypassOptions
	case C.RuleActionTypeReject:
		v = r.RejectOptions
	case C.RuleActionTypeHijackDNS:
		v = nil
	case C.RuleActionTypeSniff:
		v = r.SniffOptions
	case C.RuleActionTypeResolve:
		v = r.ResolveOptions
	default:
		return nil, E.New("unknown rule action: " + r.Action)
	}
	return badjson.MarshallObjects(_RuleAction(r), v)
}

func (r *RuleAction) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, (*_RuleAction)(r))
	if err != nil {
		return err
	}
	var v any
	switch r.Action {
	case "", C.RuleActionTypeRoute:
		r.Action = C.RuleActionTypeRoute
		v = &r.RouteOptions
	case C.RuleActionTypeRouteOptions:
		v = &r.RouteOptionsOptions
	case C.RuleActionTypeBypass:
		v = &r.BypassOptions
	case C.RuleActionTypeReject:
		v = &r.RejectOptions
	case C.RuleActionTypeHijackDNS:
		v = nil
	case C.RuleActionTypeSniff:
		v = &r.SniffOptions
	case C.RuleActionTypeResolve:
		v = &r.ResolveOptions
	default:
		return E.New("unknown rule action: " + r.Action)
	}
	if v == nil {
		// check unknown fields
		return json.UnmarshalDisallowUnknownFields(data, &_RuleAction{})
	}
	err = badjson.UnmarshallExcluded(data, (*_RuleAction)(r), v)
	if err != nil {
		return err
	}
	return nil
}

type _DNSRuleAction struct {
	// The action to perform. `route` will be used by default.
	Action string `json:"action,omitempty" enum:"route,evaluate,respond,route-options,reject,predefined"`

	// Only available with `route`, `respond`, `reject` and `predefined` actions.
	//
	// Requires [`match_response`](/configuration/dns/rule/#match_response) (for logical rules, in sub-rules).
	// Conflict with `speculative`.
	//
	// By default, rules are matched one after another in listed order: a rule with `match_response`
	// waits for its referenced responses, and no later rule is matched until it has been judged.
	//
	// A rule with `race` enabled does not hold this order: rule matching continues past it while its
	// referenced responses are still pending, so the matching of race rules runs in parallel — with
	// each other and with the rules after them. Each race rule is judged once its referenced
	// responses are available, and the first race rule that matches terminates rule evaluation
	// immediately; the remaining queries are canceled.
	//
	// Rules without `race` still take effect strictly in listed order: while a preceding race rule is
	// not yet judged, the action of any other matched rule is held until none of the race rules
	// matched.
	//
	// @Since 1.14.0
	Race                bool                         `json:"race,omitempty"`
	RouteOptions        DNSRouteActionOptions        `json:"-"`
	EvaluateOptions     DNSEvaluateActionOptions     `json:"-"`
	RouteOptionsOptions DNSRouteOptionsActionOptions `json:"-"`
	RejectOptions       RejectActionOptions          `json:"-"`
	PredefinedOptions   DNSRouteActionPredefined     `json:"-"`
}

type DNSRuleAction _DNSRuleAction

func (r DNSRuleAction) MarshalJSON() ([]byte, error) {
	if r.Action == "" {
		return json.Marshal(struct{}{})
	}
	var v any
	switch r.Action {
	case C.RuleActionTypeRoute:
		r.Action = ""
		v = r.RouteOptions
	case C.RuleActionTypeEvaluate:
		v = r.EvaluateOptions
	case C.RuleActionTypeRespond:
		v = nil
	case C.RuleActionTypeRouteOptions:
		v = r.RouteOptionsOptions
	case C.RuleActionTypeReject:
		v = r.RejectOptions
	case C.RuleActionTypePredefined:
		v = r.PredefinedOptions
	default:
		return nil, E.New("unknown DNS rule action: " + r.Action)
	}
	return badjson.MarshallObjects(_DNSRuleAction(r), v)
}

func (r *DNSRuleAction) UnmarshalJSONContext(ctx context.Context, data []byte) error {
	err := json.Unmarshal(data, (*_DNSRuleAction)(r))
	if err != nil {
		return err
	}
	var v any
	switch r.Action {
	case "", C.RuleActionTypeRoute:
		r.Action = C.RuleActionTypeRoute
		v = &r.RouteOptions
	case C.RuleActionTypeEvaluate:
		v = &r.EvaluateOptions
	case C.RuleActionTypeRespond:
		v = nil
	case C.RuleActionTypeRouteOptions:
		v = &r.RouteOptionsOptions
	case C.RuleActionTypeReject:
		v = &r.RejectOptions
	case C.RuleActionTypePredefined:
		v = &r.PredefinedOptions
	default:
		return E.New("unknown DNS rule action: " + r.Action)
	}
	if v == nil {
		return json.UnmarshalDisallowUnknownFields(data, &_DNSRuleAction{})
	}
	err = badjson.UnmarshallExcludedContext(ctx, data, (*_DNSRuleAction)(r), v)
	if err != nil {
		return err
	}
	return nil
}

// @Doc route-options Fields:
//
//	See `route-options` fields below.
type RouteActionOptions struct {
	// Tag of target outbound.
	//
	// @Required
	Outbound string `json:"outbound,omitempty" reference:"outbound"`
	RawRouteOptionsActionOptions
}

// BypassActionOptions documents the options accepted by the `bypass` action in
// the configuration reference. The runtime type is RouteActionOptions; this
// declaration only carries the bypass-specific wording.
type BypassActionOptions struct {
	// @Override outbound:
	//    Tag of target outbound.
	//
	//    If not specified, the rule only matches in [pre-match](/configuration/shared/pre-match/)
	//    from auto redirect, and will be skipped in other contexts.
	RouteActionOptions
}

// @Doc network_type:
//
//	See [Dial Fields](/configuration/shared/dial/#network_type) for details.
//
// @Doc fallback_network_type:
//
//	See [Dial Fields](/configuration/shared/dial/#fallback_network_type) for details.
type RawRouteOptionsActionOptions struct {
	// Override the connection destination address.
	OverrideAddress string `json:"override_address,omitempty"`

	// Override the connection destination port.
	OverridePort uint16 `json:"override_port,omitempty"`

	// See [Dial Fields](/configuration/shared/dial/#network_strategy) for details.
	//
	// Only take effect if outbound is direct without `outbound.bind_interface`,
	// `outbound.inet4_bind_address` and `outbound.inet6_bind_address` set.
	NetworkStrategy *NetworkStrategy `json:"network_strategy,omitempty"`

	// See [Dial Fields](/configuration/shared/dial/#fallback_delay) for details.
	FallbackDelay uint32 `json:"fallback_delay,omitempty"`

	// If enabled, for UDP proxy requests addressed to a domain,
	// the original packet address will be sent in the response instead of the mapped domain.
	//
	// This option is used for compatibility with clients that
	// do not support receiving UDP packets with domain addresses, such as Surge.
	UDPDisableDomainUnmapping bool `json:"udp_disable_domain_unmapping,omitempty"`

	// If enabled, attempts to connect UDP connection to the destination instead of listen.
	UDPConnect bool `json:"udp_connect,omitempty"`

	// Timeout for UDP connections.
	//
	// Setting a larger value than the UDP timeout in inbounds will have no effect.
	//
	// Default value for protocol sniffed connections:
	//
	// | Timeout | Protocol             |
	// |---------|----------------------|
	// | `10s`   | `dns`, `ntp`, `stun` |
	// | `30s`   | `quic`, `dtls`       |
	//
	// If no protocol is sniffed, the following ports will be recognized as protocols by default:
	//
	// | Port | Protocol |
	// |------|----------|
	// | 53   | `dns`    |
	// | 123  | `ntp`    |
	// | 443  | `quic`   |
	// | 3478 | `stun`   |
	UDPTimeout badoption.Duration `json:"udp_timeout,omitempty"`

	// Fragment TLS handshakes to bypass firewalls.
	//
	// This feature is intended to circumvent simple firewalls based on **plaintext packet matching**,
	// and should not be used to circumvent real censorship.
	//
	// Due to poor performance, try `tls_record_fragment` first, and only apply to server names known to be blocked.
	//
	// On Linux, Apple platforms, (administrator privileges required) Windows,
	// the wait time can be automatically detected. Otherwise, it will fall back to
	// waiting for a fixed time specified by `tls_fragment_fallback_delay`.
	//
	// In addition, if the actual wait time is less than 20ms, it will also fall back to waiting for a fixed time,
	// because the target is considered to be local or behind a transparent proxy.
	//
	// @Since 1.12.0
	TLSFragment bool `json:"tls_fragment,omitempty"`

	// The fallback value used when TLS segmentation cannot automatically determine the wait time.
	//
	// @Default: 500ms
	//
	// @Since 1.12.0
	TLSFragmentFallbackDelay badoption.Duration `json:"tls_fragment_fallback_delay,omitempty"`

	// Fragment TLS handshake into multiple TLS records to bypass firewalls.
	//
	// @Since 1.12.0
	TLSRecordFragment bool `json:"tls_record_fragment,omitempty"`

	// Inject a forged TLS ClientHello carrying this SNI before the real one,
	// to fool SNI-filtering middleboxes that permit specific hostnames.
	//
	// See outbound TLS [`spoof`](/configuration/shared/tls/#spoof) for details
	// and required privileges.
	//
	// @Since 1.14.0
	//
	// @Quote:
	//
	//     Only supported on Linux, macOS, and Windows, and requires elevated privileges.
	TLSSpoof string `json:"tls_spoof,omitempty"`

	// How the forged segment is rejected by the real server. See outbound TLS
	// [`spoof_method`](/configuration/shared/tls/#spoof_method) for the full table
	// of accepted values and platform notes.
	//
	// @Since 1.14.0
	TLSSpoofMethod string `json:"tls_spoof_method,omitempty" enum:"wrong-sequence,wrong-checksum,wrong-ack,wrong-md5,wrong-timestamp"`
}

type RouteOptionsActionOptions RawRouteOptionsActionOptions

func (r *RouteOptionsActionOptions) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, (*RawRouteOptionsActionOptions)(r))
	if err != nil {
		return err
	}
	if *r == (RouteOptionsActionOptions{}) {
		return E.New("empty route option action")
	}
	if r.TLSFragment && r.TLSRecordFragment {
		return E.New("`tls_fragment` and `tls_record_fragment` are mutually exclusive")
	}
	return nil
}

type DNSRouteActionOptions struct {
	// Tag of target server.
	//
	// @Required
	Server string `json:"server,omitempty" reference:"dns_server"`

	// Conflict with `race`. Has no effect without a preceding `race` rule.
	//
	// By default, no query is sent in parallel with pending race rules: a matched `route` action
	// holds its query until none of the race rules matched.
	//
	// When `speculative` is enabled, the query is sent as soon as the rule matches, in parallel with
	// the pending race rules; its response is used only after none of the race rules matched.
	//
	// @Since 1.14.0
	Speculative bool `json:"speculative,omitempty"`
	AbstractDNSRouteActionOptions
}

type DNSEvaluateActionOptions struct {
	// Tag of target server.
	//
	// @Required
	Server string `json:"server,omitempty" reference:"dns_server"`

	// Tag of the evaluated response.
	//
	// A tagged response is only referenced via [`match_response`](/configuration/dns/rule/#match_response) with the tag;
	// `match_response: true` references the response of the latest `evaluate` action without `tag`.
	Tag string `json:"tag,omitempty"`

	// Has no effect without a preceding `race` rule.
	//
	// By default, no query is sent in parallel with pending race rules: a matched `evaluate` action
	// holds its query, and rule matching stops there, until none of the race rules matched.
	//
	// When `speculative` is enabled, the query is sent as soon as the rule matches, in parallel with
	// the pending race rules, and rule matching continues without waiting for them.
	//
	// @Since 1.14.0
	Speculative bool `json:"speculative,omitempty"`
	AbstractDNSRouteActionOptions
}

type AbstractDNSRouteActionOptions struct {
	// Override the DNS query timeout for matched queries.
	//
	// Will override `dns.timeout`.
	//
	// @Since 1.14.0
	Timeout badoption.Duration `json:"timeout,omitempty"`

	// Set domain strategy for this query.
	//
	// One of `prefer_ipv4` `prefer_ipv6` `ipv4_only` `ipv6_only`.
	//
	// @Failure: Deprecated in sing-box 1.14.0
	//
	//     `strategy` is deprecated in sing-box 1.14.0 and will be removed in sing-box 1.16.0.
	//
	// @Since 1.12.0
	Strategy DomainStrategy `json:"strategy,omitempty" schema:"omit"`

	// Disable cache and save cache in this query.
	DisableCache bool `json:"disable_cache,omitempty"`

	// Disable optimistic DNS caching in this query.
	//
	// @Since 1.14.0
	DisableOptimisticCache bool `json:"disable_optimistic_cache,omitempty"`

	// Rewrite TTL in DNS responses.
	RewriteTTL *uint32 `json:"rewrite_ttl,omitempty"`

	// Append a `edns0-subnet` OPT extra record with the specified IP prefix to every query by default.
	//
	// If value is an IP address instead of prefix, `/32` or `/128` will be appended automatically.
	//
	// Will override `dns.client_subnet`.
	ClientSubnet *badoption.Prefixable `json:"client_subnet,omitempty"`

	// Remove the `edns0-subnet` OPT extra record from the query, and suppress `dns.client_subnet`.
	//
	// Conflict with `client_subnet`.
	//
	// @Since 1.14.0
	RemoveClientSubnet bool `json:"remove_client_subnet,omitempty"`
}

type DNSRouteOptionsActionOptions AbstractDNSRouteActionOptions

func (r *DNSRouteOptionsActionOptions) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, (*AbstractDNSRouteActionOptions)(r))
	if err != nil {
		return err
	}
	if *r == (DNSRouteOptionsActionOptions{}) {
		return E.New("empty DNS route option action")
	}
	return nil
}

type _RejectActionOptions struct {
	// For TCP and UDP connections:
	//
	// - `default`: Reply with TCP RST for TCP connections, and ICMP port unreachable for UDP packets.
	// - `drop`: Drop packets.
	//
	// For ICMP echo requests:
	//
	// - `default`: Reply with ICMP host unreachable.
	// - `drop`: Drop packets.
	// - `reply`: Reply with ICMP echo reply.
	Method string `json:"method,omitempty" enum:"default,drop,reply"`

	// If not enabled, `method` will be temporarily overwritten to `drop` after 50 triggers in 30s.
	//
	// Not available when `method` is set to drop.
	NoDrop bool `json:"no_drop,omitempty"`
}

type RejectActionOptions _RejectActionOptions

func (r RejectActionOptions) MarshalJSON() ([]byte, error) {
	switch r.Method {
	case C.RuleActionRejectMethodDefault:
		r.Method = ""
	}
	return json.Marshal(_RejectActionOptions(r))
}

func (r *RejectActionOptions) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, (*_RejectActionOptions)(r))
	if err != nil {
		return err
	}
	switch r.Method {
	case "", C.RuleActionRejectMethodDefault:
		r.Method = C.RuleActionRejectMethodDefault
	case C.RuleActionRejectMethodDrop:
	case C.RuleActionRejectMethodReply:
	default:
		return E.New("unknown reject method: " + r.Method)
	}
	if r.Method == C.RuleActionRejectMethodDrop && r.NoDrop {
		return E.New("no_drop is not available in current context")
	}
	return nil
}

// DNSRejectActionOptions documents the options accepted by the DNS `reject`
// action in the configuration reference. The runtime type is
// RejectActionOptions; this declaration only carries the DNS-specific wording.
type DNSRejectActionOptions struct {
	// @Override method:
	//    - `default`: Reply with REFUSED.
	//    - `drop`: Drop the request.
	//
	//    `default` will be used by default.
	RejectActionOptions
}

type RouteActionSniff struct {
	// Enabled sniffers.
	//
	// All sniffers enabled by default.
	//
	// Available protocol values an be found on in [Protocol Sniff](../sniff/)
	Sniffer badoption.Listable[string] `json:"sniffer,omitempty" enum:"tls,http,quic,dns,stun,bittorrent,dtls,ssh,rdp,ntp"`

	// Timeout for sniffing.
	//
	// @Default: 300ms
	Timeout badoption.Duration `json:"timeout,omitempty"`
}

type RouteActionResolve struct {
	// Specifies DNS server tag to use instead of selecting through DNS routing.
	Server string `json:"server,omitempty" reference:"dns_server"`

	// Override the DNS query timeout for this lookup.
	//
	// Will override `dns.timeout`.
	//
	// @Since 1.14.0
	Timeout badoption.Duration `json:"timeout,omitempty"`

	// DNS resolution strategy, available values are: `prefer_ipv4`, `prefer_ipv6`, `ipv4_only`, `ipv6_only`.
	//
	// `dns.strategy` will be used by default.
	Strategy DomainStrategy `json:"strategy,omitempty"`

	// Disable cache and save cache in this query.
	//
	// @Since 1.12.0
	DisableCache bool `json:"disable_cache,omitempty"`

	// Disable optimistic DNS caching in this query.
	//
	// @Since 1.14.0
	DisableOptimisticCache bool `json:"disable_optimistic_cache,omitempty"`

	// Rewrite TTL in DNS responses.
	//
	// @Since 1.12.0
	RewriteTTL *uint32 `json:"rewrite_ttl,omitempty"`

	// Append a `edns0-subnet` OPT extra record with the specified IP prefix to every query by default.
	//
	// If value is an IP address instead of prefix, `/32` or `/128` will be appended automatically.
	//
	// Will override `dns.client_subnet`.
	//
	// @Since 1.12.0
	ClientSubnet *badoption.Prefixable `json:"client_subnet,omitempty"`
}

type DNSRouteActionPredefined struct {
	// The response code.
	//
	// | Value      | Value in the legacy rcode server | Description     |
	// |------------|----------------------------------|-----------------|
	// | `NOERROR`  | `success`                        | Ok              |
	// | `FORMERR`  | `format_error`                   | Bad request     |
	// | `SERVFAIL` | `server_failure`                 | Server failure  |
	// | `NXDOMAIN` | `name_error`                     | Not found       |
	// | `NOTIMP`   | `not_implemented`                | Not implemented |
	// | `REFUSED`  | `refused`                        | Refused         |
	//
	// @Default: NOERROR
	Rcode *DNSRCode `json:"rcode,omitempty"`

	// List of text DNS record to respond as answers.
	//
	// Examples:
	//
	// | Record Type | Example                       |
	// |-------------|-------------------------------|
	// | `A`         | `localhost. IN A 127.0.0.1`   |
	// | `AAAA`      | `localhost. IN AAAA ::1`      |
	// | `TXT`       | `localhost. IN TXT \"Hello\"` |
	Answer badoption.Listable[DNSRecordOptions] `json:"answer,omitempty"`

	// List of text DNS record to respond as name servers.
	Ns badoption.Listable[DNSRecordOptions] `json:"ns,omitempty"`

	// List of text DNS record to respond as extra records.
	Extra badoption.Listable[DNSRecordOptions] `json:"extra,omitempty"`
}

type actionVariant struct {
	action         string
	actionOptional bool
	structType     reflect.Type
	build          func(variant *schema.Node) error
}

func actionUnion(builder schema.Builder, variants []actionVariant) (*schema.Node, error) {
	variantNodes := make([]*schema.Node, 0, len(variants))
	for _, variant := range variants {
		variantNode := schema.LooseObject()
		variantNode.Properties.Put("action", schema.StringConst(variant.action))
		if !variant.actionOptional {
			variantNode.Required = []string{"action"}
		}
		if variant.build != nil {
			err := variant.build(variantNode)
			if err != nil {
				return nil, err
			}
		}
		if variant.structType != nil {
			err := builder.FlattenStruct(variantNode, variant.structType)
			if err != nil {
				return nil, err
			}
		}
		variantNodes = append(variantNodes, variantNode)
	}
	return schema.OneOf(variantNodes...), nil
}

func rejectProperties(variant *schema.Node) error {
	variant.Properties.Put("method", schema.StringEnum(
		"",
		C.RuleActionRejectMethodDefault,
		C.RuleActionRejectMethodDrop,
		C.RuleActionRejectMethodReply,
	))
	variant.Properties.Put("no_drop", schema.BooleanNode())
	return nil
}

func routeActionUnion(builder schema.Builder) (*schema.Node, error) {
	return actionUnion(builder, []actionVariant{
		{action: C.RuleActionTypeRoute, actionOptional: true, structType: reflect.TypeFor[RouteActionOptions]()},
		{action: C.RuleActionTypeRouteOptions, structType: reflect.TypeFor[RawRouteOptionsActionOptions]()},
		{action: C.RuleActionTypeBypass, structType: reflect.TypeFor[RouteActionOptions]()},
		{action: C.RuleActionTypeReject, build: rejectProperties},
		{action: C.RuleActionTypeHijackDNS},
		{action: C.RuleActionTypeSniff, structType: reflect.TypeFor[RouteActionSniff]()},
		{action: C.RuleActionTypeResolve, structType: reflect.TypeFor[RouteActionResolve]()},
	})
}

func dnsActionUnion(builder schema.Builder) (*schema.Node, error) {
	raceProperty := func(variant *schema.Node) error {
		variant.Properties.Put("race", schema.BooleanNode())
		return nil
	}
	rejectWithRace := func(variant *schema.Node) error {
		err := raceProperty(variant)
		if err != nil {
			return err
		}
		return rejectProperties(variant)
	}
	return actionUnion(builder, []actionVariant{
		{action: C.RuleActionTypeRoute, actionOptional: true, structType: reflect.TypeFor[DNSRouteActionOptions](), build: raceProperty},
		{action: C.RuleActionTypeEvaluate, structType: reflect.TypeFor[DNSEvaluateActionOptions](), build: raceProperty},
		{action: C.RuleActionTypeRespond, build: raceProperty},
		{action: C.RuleActionTypeRouteOptions, structType: reflect.TypeFor[DNSRouteOptionsActionOptions](), build: raceProperty},
		{action: C.RuleActionTypeReject, build: rejectWithRace},
		{action: C.RuleActionTypePredefined, structType: reflect.TypeFor[DNSRouteActionPredefined](), build: raceProperty},
	})
}
