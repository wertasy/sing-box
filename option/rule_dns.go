package option

import (
	"context"
	"reflect"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/schema"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badjson"
	"github.com/sagernet/sing/common/json/badoption"
)

type _DNSRule struct {
	Type           string         `json:"type,omitempty" enum:"default,logical"`
	DefaultOptions DefaultDNSRule `json:"-"`
	LogicalOptions LogicalDNSRule `json:"-"`
}

type DNSRule _DNSRule

func (r DNSRule) MarshalJSON() ([]byte, error) {
	var v any
	switch r.Type {
	case C.RuleTypeDefault:
		r.Type = ""
		v = r.DefaultOptions
	case C.RuleTypeLogical:
		v = r.LogicalOptions
	default:
		return nil, E.New("unknown rule type: " + r.Type)
	}
	return badjson.MarshallObjects(_DNSRule(r), v)
}

func (r *DNSRule) UnmarshalJSONContext(ctx context.Context, bytes []byte) error {
	err := json.UnmarshalContext(ctx, bytes, (*_DNSRule)(r))
	if err != nil {
		return err
	}
	var v any
	switch r.Type {
	case "", C.RuleTypeDefault:
		r.Type = C.RuleTypeDefault
		v = &r.DefaultOptions
	case C.RuleTypeLogical:
		v = &r.LogicalOptions
	default:
		return E.New("unknown rule type: " + r.Type)
	}
	err = badjson.UnmarshallExcludedContext(ctx, bytes, (*_DNSRule)(r), v)
	if err != nil {
		return err
	}
	return nil
}

func (r DNSRule) IsValid() bool {
	switch r.Type {
	case C.RuleTypeDefault:
		return r.DefaultOptions.IsValid()
	case C.RuleTypeLogical:
		return r.LogicalOptions.IsValid()
	default:
		panic("unknown DNS rule type: " + r.Type)
	}
}

func (r DNSRule) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return builder.Define("DNSRule", func() (*schema.Node, error) {
		actionRef, err := builder.Define("DNSRuleAction", func() (*schema.Node, error) {
			return dnsActionUnion(builder)
		})
		if err != nil {
			return nil, err
		}
		nestedRef, err := builder.Define("NestedDNSRule", func() (*schema.Node, error) {
			return nestedRuleUnion(builder, reflect.TypeFor[RawDefaultDNSRule](), "NestedDNSRule")
		})
		if err != nil {
			return nil, err
		}
		return ruleUnion(builder, reflect.TypeFor[RawDefaultDNSRule](), nestedRef, actionRef)
	})
}

type DNSRuleMatchResponse struct {
	Enabled bool
	Tag     string
}

func (m *DNSRuleMatchResponse) UnmarshalJSON(content []byte) error {
	var boolValue bool
	err := json.Unmarshal(content, &boolValue)
	if err == nil {
		m.Enabled = boolValue
		m.Tag = ""
		return nil
	}
	var stringValue string
	err = json.Unmarshal(content, &stringValue)
	if err != nil {
		return E.New("invalid match_response value")
	}
	if stringValue == "" {
		return E.New("empty match_response tag")
	}
	m.Enabled = true
	m.Tag = stringValue
	return nil
}

func (m DNSRuleMatchResponse) MarshalJSON() ([]byte, error) {
	if m.Tag != "" {
		return json.Marshal(m.Tag)
	}
	return json.Marshal(m.Enabled)
}

func (m *DNSRuleMatchResponse) IsEnabled() bool {
	return m != nil && m.Enabled
}

func (m *DNSRuleMatchResponse) ResponseTag() string {
	if m == nil {
		return ""
	}
	return m.Tag
}

func (m DNSRuleMatchResponse) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return schema.AnyOf(schema.BooleanNode(), schema.StringNode()), nil
}

type RawDefaultDNSRule struct {
	// Tags of [Inbound](/configuration/inbound/).
	Inbound badoption.Listable[string] `json:"inbound,omitempty" reference:"inbound"`
	// @quote: Changes in sing-box 1.14.0
	//
	//     This field now also applies when a DNS rule is matched from an internal
	//     domain resolution that does not target a specific DNS server, such as a
	//     [`resolve`](../../route/rule_action/#resolve) route rule action without a
	//     `server` set. See
	//     [Migration](/migration/#ip_version-and-query_type-behavior-changes-in-dns-rules).
	//
	//     Setting this field makes the DNS rule incompatible in the same DNS
	//     configuration with Legacy Address Filter Fields in DNS rules, the Legacy
	//     `strategy` DNS rule action option, and the Legacy
	//     `rule_set_ip_cidr_accept_empty` DNS rule item. To combine with
	//     address-based filtering, use the [`evaluate`](../rule_action/#evaluate)
	//     action and [`match_response`](#match_response).
	//
	// 4 (A DNS query) or 6 (AAAA DNS query).
	//
	// Not limited if empty.
	IPVersion int `json:"ip_version,omitempty" enum:"4,6"`
	// @quote: Changes in sing-box 1.14.0
	//
	//     This field now also applies when a DNS rule is matched from an internal
	//     domain resolution that does not target a specific DNS server, such as a
	//     [`resolve`](../../route/rule_action/#resolve) route rule action without a
	//     `server` set. See
	//     [Migration](/migration/#ip_version-and-query_type-behavior-changes-in-dns-rules).
	//
	//     Setting this field makes the DNS rule incompatible in the same DNS
	//     configuration with Legacy Address Filter Fields in DNS rules, the Legacy
	//     `strategy` DNS rule action option, and the Legacy
	//     `rule_set_ip_cidr_accept_empty` DNS rule item. To combine with
	//     address-based filtering, use the [`evaluate`](../rule_action/#evaluate)
	//     action and [`match_response`](#match_response).
	//
	// DNS query type. Values can be integers or type name strings.
	QueryType badoption.Listable[DNSQueryType] `json:"query_type,omitempty"`
	// @since 1.14.0
	//
	// Match the `edns0-subnet` OPT extra record (EDNS Client Subnet) in the query.
	//
	// A listed prefix matches when it is no more specific than the received client subnet and contains its address.
	//
	// If value is an IP address instead of prefix, `/32` or `/128` will be appended automatically.
	QueryClientSubnet badoption.Listable[*badoption.Prefixable] `json:"query_client_subnet,omitempty"`
	// @since 1.14.0
	//
	// Match queries with the DNSSEC OK (`DO`) bit set.
	QueryDNSSEC bool `json:"query_dnssec,omitempty"`
	// `tcp` or `udp`.
	Network badoption.Listable[string] `json:"network,omitempty" enum:"tcp,udp"`
	// Username, see each inbound for details.
	AuthUser badoption.Listable[string] `json:"auth_user,omitempty"`
	// Sniffed protocol, see [Sniff](/configuration/route/sniff/) for details.
	Protocol badoption.Listable[string] `json:"protocol,omitempty" enum:"tls,http,quic,dns,stun,bittorrent,dtls,ssh,rdp,ntp"`
	// Match full domain.
	Domain badoption.Listable[string] `json:"domain,omitempty"`
	// Match domain suffix.
	DomainSuffix badoption.Listable[string] `json:"domain_suffix,omitempty"`
	// Match domain using keyword.
	DomainKeyword badoption.Listable[string] `json:"domain_keyword,omitempty"`
	// Match domain using regular expression.
	DomainRegex badoption.Listable[string] `json:"domain_regex,omitempty"`
	// Match source IP CIDR.
	SourceIPCIDR badoption.Listable[*badoption.Prefixable] `json:"source_ip_cidr,omitempty"`
	// @since 1.8.0
	//
	// Match non-public source IP.
	SourceIPIsPrivate bool `json:"source_ip_is_private,omitempty"`
	// Match source port.
	SourcePort badoption.Listable[uint16] `json:"source_port,omitempty"`
	// Match source port range.
	SourcePortRange badoption.Listable[string] `json:"source_port_range,omitempty"`
	// Match port.
	Port badoption.Listable[uint16] `json:"port,omitempty"`
	// Match port range.
	PortRange badoption.Listable[string] `json:"port_range,omitempty"`
	// @quote:
	//
	//     Only supported on Linux, Windows, and macOS.
	//
	// Match process name.
	ProcessName badoption.Listable[string] `json:"process_name,omitempty"`
	// @quote:
	//
	//     Only supported on Linux, Windows, and macOS.
	//
	// Match process path.
	ProcessPath badoption.Listable[string] `json:"process_path,omitempty"`
	// @since 1.10.0
	//
	// @quote:
	//
	//     Only supported on Linux, Windows, and macOS.
	//
	// Match process path using regular expression.
	ProcessPathRegex badoption.Listable[string] `json:"process_path_regex,omitempty"`
	// Match android package name.
	PackageName badoption.Listable[string] `json:"package_name,omitempty"`
	// @since 1.14.0
	//
	// Match android package name using regular expression.
	PackageNameRegex badoption.Listable[string] `json:"package_name_regex,omitempty"`
	// @quote:
	//
	//     Only supported on Linux.
	//
	// Match user name.
	User badoption.Listable[string] `json:"user,omitempty"`
	// @quote:
	//
	//     Only supported on Linux.
	//
	// Match user id.
	UserID badoption.Listable[int32] `json:"user_id,omitempty"`
	// @failure: Deprecated in sing-box 1.12.0
	//
	//     `outbound` rule items are deprecated and will be removed in sing-box 1.14.0, check [Migration](/migration/#migrate-outbound-dns-rule-items-to-domain-resolver).
	//
	// Match outbound.
	//
	// `any` can be used as a value to match any outbound.
	Outbound badoption.Listable[string] `json:"outbound,omitempty" reference:"outbound" schema:"omit"`
	// Match Clash mode.
	ClashMode string `json:"clash_mode,omitempty"`
	// @since 1.11.0
	//
	// @quote:
	//
	//     Only supported in graphical clients on Android and Apple platforms.
	//
	// Match network type.
	//
	// Available values: `wifi`, `cellular`, `ethernet` and `other`.
	NetworkType badoption.Listable[InterfaceType] `json:"network_type,omitempty"`
	// @since 1.11.0
	//
	// @quote:
	//
	//     Only supported in graphical clients on Android and Apple platforms.
	//
	// Match if network is considered Metered (on Android) or considered expensive,
	// such as Cellular or a Personal Hotspot (on Apple platforms).
	NetworkIsExpensive bool `json:"network_is_expensive,omitempty"`
	// @since 1.11.0
	//
	// @quote:
	//
	//     Only supported in graphical clients on Apple platforms.
	//
	// Match if network is in Low Data Mode.
	NetworkIsConstrained bool `json:"network_is_constrained,omitempty"`
	// @quote:
	//
	//     Only supported in graphical clients on Android and Apple platforms, or on Linux.
	//
	// Match WiFi SSID.
	WIFISSID badoption.Listable[string] `json:"wifi_ssid,omitempty"`
	// @quote:
	//
	//     Only supported in graphical clients on Android and Apple platforms, or on Linux.
	//
	// Match WiFi BSSID.
	WIFIBSSID badoption.Listable[string] `json:"wifi_bssid,omitempty"`
	// @since 1.13.0
	//
	// @quote:
	//
	//     Only supported on Linux, Windows, and macOS.
	//
	// Match interface address.
	InterfaceAddress *badjson.TypedMap[string, badoption.Listable[*badoption.Prefixable]] `json:"interface_address,omitempty"`
	// @since 1.13.0
	//
	// @quote:
	//
	//     Only supported in graphical clients on Android and Apple platforms.
	//
	// Matches network interface (same values as `network_type`) address.
	NetworkInterfaceAddress *badjson.TypedMap[InterfaceType, badoption.Listable[*badoption.Prefixable]] `json:"network_interface_address,omitempty"`
	// @since 1.13.0
	//
	// @quote:
	//
	//     Only supported on Linux, Windows, and macOS.
	//
	// Match default interface address.
	DefaultInterfaceAddress badoption.Listable[*badoption.Prefixable] `json:"default_interface_address,omitempty"`
	// @since 1.14.0
	//
	// @quote:
	//
	//     Only supported on Linux, macOS, or in graphical clients on Android and macOS. See [Neighbor Resolution](/configuration/shared/neighbor/) for setup.
	//
	// Match source device MAC address.
	SourceMACAddress badoption.Listable[string] `json:"source_mac_address,omitempty"`
	// @since 1.14.0
	//
	// @quote:
	//
	//     Only supported on Linux, macOS, or in graphical clients on Android and macOS. See [Neighbor Resolution](/configuration/shared/neighbor/) for setup.
	//
	// Match source device hostname from DHCP leases.
	SourceHostname badoption.Listable[string] `json:"source_hostname,omitempty"`
	// @since 1.14.0
	//
	// Match specified DNS servers' preferred domains.
	//
	// | Type          | Match                                                                        |
	// |---------------|------------------------------------------------------------------------------|
	// | `hosts`       | Match predefined entries and entries in hosts files                          |
	// | `local`       | Match hosts entries, neighbor-resolved hosts, and mDNS local domains         |
	// | `mdns`        | Match mDNS local domains (`*.local.` and IPv4/IPv6 link-local reverse zones) |
	// | `tailscale`   | Match MagicDNS hosts and DNS route suffixes                                  |
	// | `openconnect` | Match split DNS and search domains pushed by the VPN server                  |
	// | `resolved`    | Match split DNS and search domains from systemd-resolved links               |
	PreferredBy badoption.Listable[string] `json:"preferred_by,omitempty"`
	// @since 1.15.0
	//
	// Match specified DNS servers' server addresses.
	//
	// | Type          | Match                                         |
	// |---------------|-----------------------------------------------|
	// | `local`       | Match system DNS servers                      |
	// | `dhcp`        | Match DNS servers from DHCP                   |
	// | `resolved`    | Match DNS servers from systemd-resolved links |
	// | `tailscale`   | Match DNS resolvers of the tailnet            |
	// | `openvpn`     | Match DNS servers pushed by the VPN server    |
	// | `openconnect` | Match DNS servers pushed by the VPN server    |
	DNSServerAddress *badjson.TypedMap[string, badoption.Listable[*badoption.Prefixable]] `json:"dns_server_address,omitempty"`
	// @since 1.15.0
	//
	// Match specified DNS servers' search domains.
	//
	// | Type          | Match                                            |
	// |---------------|--------------------------------------------------|
	// | `local`       | Match system search domains                      |
	// | `dhcp`        | Match search domains from DHCP                   |
	// | `resolved`    | Match search domains from systemd-resolved links |
	// | `tailscale`   | Match search domains of the tailnet              |
	// | `openvpn`     | Match search domains pushed by the VPN server    |
	// | `openconnect` | Match search domains pushed by the VPN server    |
	DNSSearchDomain *badjson.TypedMap[string, badoption.Listable[string]] `json:"dns_search_domain,omitempty"`
	// @since 1.8.0
	//
	// Match [rule-set](/configuration/route/#rule_set).
	RuleSet badoption.Listable[string] `json:"rule_set,omitempty" reference:"rule_set"`
	// @since 1.10.0
	//
	// Make `ip_cidr` rule items in rule-sets match the source IP.
	RuleSetIPCIDRMatchSource bool `json:"rule_set_ip_cidr_match_source,omitempty"`
	// @since 1.14.0
	//
	// Enable response-based matching. When enabled, this rule matches against the evaluated response
	// (set by a preceding [`evaluate`](/configuration/dns/rule_action/#evaluate) action)
	// instead of only matching the original query.
	//
	// `true` or the `tag` of an `evaluate` action: `true` matches against the response of the latest
	// `evaluate` action without `tag`; a tag matches against the response of the `evaluate` action with the tag.
	//
	// The evaluated response can also be returned directly by a later [`respond`](/configuration/dns/rule_action/#respond) action;
	// in a rule with a `match_response` tag, `respond` returns the tagged response.
	//
	// Required for Response Match Fields (`response_rcode`, `response_answer`, `response_ns`, `response_extra`).
	// Also required for `ip_cidr`, `ip_is_private`, and `ip_accept_any` when used with `evaluate` or Response Match Fields.
	MatchResponse *DNSRuleMatchResponse `json:"match_response,omitempty"`
	// @since 1.9.0
	//
	// Match IP CIDR with query response.
	//
	// As a Legacy Address Filter Field, deprecated. Use with `match_response` instead,
	// check [Migration](/migration/#migrate-address-filter-fields-to-response-matching).
	IPCIDR badoption.Listable[*badoption.Prefixable] `json:"ip_cidr,omitempty"`
	// @since 1.9.0
	//
	// Match private IP with query response.
	//
	// As a Legacy Address Filter Field, deprecated. Use with `match_response` instead,
	// check [Migration](/migration/#migrate-address-filter-fields-to-response-matching).
	IPIsPrivate bool `json:"ip_is_private,omitempty"`
	// @since 1.12.0
	//
	// Match when the DNS query response contains at least one address.
	IPAcceptAny bool `json:"ip_accept_any,omitempty"`
	// Match DNS response code.
	//
	// Accepted values are the same as in the [predefined action rcode](/configuration/dns/rule_action/#rcode).
	ResponseRcode *DNSRCode `json:"response_rcode,omitempty"`
	// Match DNS answer records.
	//
	// Record format is the same as in [predefined action answer](/configuration/dns/rule_action/#answer).
	ResponseAnswer badoption.Listable[DNSRecordOptions] `json:"response_answer,omitempty"`
	// Match DNS name server records.
	//
	// Record format is the same as in [predefined action ns](/configuration/dns/rule_action/#ns).
	ResponseNs badoption.Listable[DNSRecordOptions] `json:"response_ns,omitempty"`
	// Match DNS extra records.
	//
	// Record format is the same as in [predefined action extra](/configuration/dns/rule_action/#extra).
	ResponseExtra badoption.Listable[DNSRecordOptions] `json:"response_extra,omitempty"`
	// Invert match result.
	Invert bool `json:"invert,omitempty"`

	// Deprecated: removed in sing-box 1.12.0
	// @failure: Deprecated in sing-box 1.8.0
	//
	//     Geosite is deprecated and will be removed in sing-box 1.12.0, check [Migration](/migration/#migrate-geosite-to-rule-sets).
	//
	// Match geosite.
	Geosite badoption.Listable[string] `json:"geosite,omitempty" schema:"omit"`
	// @failure: Deprecated in sing-box 1.8.0
	//
	//     GeoIP is deprecated and will be removed in sing-box 1.12.0, check [Migration](/migration/#migrate-geoip-to-rule-sets).
	//
	// Match source geoip.
	SourceGeoIP badoption.Listable[string] `json:"source_geoip,omitempty" schema:"omit"`
	// @failure: Removed in sing-box 1.12.0
	//
	//     GeoIP is deprecated in sing-box 1.8.0 and removed in sing-box 1.12.0, check [Migration](/migration/#migrate-geoip-to-rule-sets).
	//
	// Match GeoIP with query response.
	GeoIP badoption.Listable[string] `json:"geoip,omitempty" schema:"omit"`
	// Deprecated: removed in sing-box 1.11.0
	// @since 1.10.0
	//
	// @failure: Deprecated in sing-box 1.14.0
	//
	//     `rule_set_ip_cidr_accept_empty` is deprecated and will be removed in sing-box 1.16.0,
	//     check [Migration](/migration/#migrate-address-filter-fields-to-response-matching).
	//
	// Make `ip_cidr` rules in rule-sets accept empty query response.
	RuleSetIPCIDRAcceptEmpty bool `json:"rule_set_ip_cidr_accept_empty,omitempty" schema:"omit"`
	// Deprecated: renamed to rule_set_ip_cidr_match_source
	// @since 1.9.0
	//
	// @failure: Deprecated in sing-box 1.10.0
	//
	//     `rule_set_ipcidr_match_source` is renamed to `rule_set_ip_cidr_match_source` and will be remove in sing-box 1.11.0.
	//
	// Make `ip_cidr` rule items in rule-sets match the source IP.
	Deprecated_RulesetIPCIDRMatchSource bool `json:"rule_set_ipcidr_match_source,omitempty" schema:"omit"`
}

type DefaultDNSRule struct {
	RawDefaultDNSRule
	DNSRuleAction
}

func (r DefaultDNSRule) MarshalJSON() ([]byte, error) {
	return badjson.MarshallObjects(r.RawDefaultDNSRule, r.DNSRuleAction)
}

func (r *DefaultDNSRule) UnmarshalJSONContext(ctx context.Context, data []byte) error {
	rawAction, routeOptions, err := inspectDNSRuleAction(ctx, data)
	if err != nil {
		return err
	}
	err = rejectNestedDNSRuleAction(ctx, data)
	if err != nil {
		return err
	}
	depth := nestedRuleDepth(ctx)
	err = json.UnmarshalContext(ctx, data, &r.RawDefaultDNSRule)
	if err != nil {
		return err
	}
	err = badjson.UnmarshallExcludedContext(ctx, data, &r.RawDefaultDNSRule, &r.DNSRuleAction)
	if err != nil {
		return err
	}
	if depth > 0 && rawAction == "" && routeOptions == (DNSRouteActionOptions{}) {
		r.DNSRuleAction = DNSRuleAction{}
	}
	return nil
}

func (r DefaultDNSRule) IsValid() bool {
	var defaultValue DefaultDNSRule
	defaultValue.Invert = r.Invert
	return !reflect.DeepEqual(r, defaultValue)
}

type RawLogicalDNSRule struct {
	// `and` or `or`
	Mode string `json:"mode" enum:"and,or"`
	// Included rules.
	Rules []DNSRule `json:"rules,omitempty"`
	// Invert match result.
	Invert bool `json:"invert,omitempty"`
}

type LogicalDNSRule struct {
	RawLogicalDNSRule
	DNSRuleAction
}

func (r LogicalDNSRule) MarshalJSON() ([]byte, error) {
	return badjson.MarshallObjects(r.RawLogicalDNSRule, r.DNSRuleAction)
}

func (r *LogicalDNSRule) UnmarshalJSONContext(ctx context.Context, data []byte) error {
	rawAction, routeOptions, err := inspectDNSRuleAction(ctx, data)
	if err != nil {
		return err
	}
	err = rejectNestedDNSRuleAction(ctx, data)
	if err != nil {
		return err
	}
	depth := nestedRuleDepth(ctx)
	err = json.UnmarshalContext(nestedRuleChildContext(ctx), data, &r.RawLogicalDNSRule)
	if err != nil {
		return err
	}
	err = badjson.UnmarshallExcludedContext(ctx, data, &r.RawLogicalDNSRule, &r.DNSRuleAction)
	if err != nil {
		return err
	}
	if depth > 0 && rawAction == "" && routeOptions == (DNSRouteActionOptions{}) {
		r.DNSRuleAction = DNSRuleAction{}
	}
	return nil
}

func (r *LogicalDNSRule) IsValid() bool {
	return len(r.Rules) > 0 && common.All(r.Rules, DNSRule.IsValid)
}
