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

type _Rule struct {
	Type           string      `json:"type,omitempty" enum:"default,logical"`
	DefaultOptions DefaultRule `json:"-"`
	LogicalOptions LogicalRule `json:"-"`
}

type Rule _Rule

func (r Rule) MarshalJSON() ([]byte, error) {
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
	return badjson.MarshallObjects(_Rule(r), v)
}

func (r *Rule) UnmarshalJSONContext(ctx context.Context, bytes []byte) error {
	err := json.UnmarshalContext(ctx, bytes, (*_Rule)(r))
	if err != nil {
		return err
	}
	payload, err := rulePayloadWithoutType(ctx, bytes)
	if err != nil {
		return err
	}
	switch r.Type {
	case "", C.RuleTypeDefault:
		r.Type = C.RuleTypeDefault
		return unmarshalDefaultRuleContext(ctx, payload, &r.DefaultOptions)
	case C.RuleTypeLogical:
		return unmarshalLogicalRuleContext(ctx, payload, &r.LogicalOptions)
	default:
		return E.New("unknown rule type: " + r.Type)
	}
}

func (r Rule) IsValid() bool {
	switch r.Type {
	case C.RuleTypeDefault:
		return r.DefaultOptions.IsValid()
	case C.RuleTypeLogical:
		return r.LogicalOptions.IsValid()
	default:
		panic("unknown rule type: " + r.Type)
	}
}

func (r Rule) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return builder.Define("Rule", func() (*schema.Node, error) {
		actionRef, err := builder.Define("RuleAction", func() (*schema.Node, error) {
			return routeActionUnion(builder)
		})
		if err != nil {
			return nil, err
		}
		nestedRef, err := builder.Define("NestedRule", func() (*schema.Node, error) {
			return nestedRuleUnion(builder, reflect.TypeFor[RawDefaultRule](), "NestedRule")
		})
		if err != nil {
			return nil, err
		}
		return ruleUnion(builder, reflect.TypeFor[RawDefaultRule](), nestedRef, actionRef)
	})
}

// ruleUnion builds the top-level rule schema: match fields composed with rule
// actions via unevaluatedProperties, mirroring the badjson.UnmarshallExcluded
// composition in DefaultRule / LogicalRule.
func ruleUnion(builder schema.Builder, matchType reflect.Type, nestedRef *schema.Node, actionRef *schema.Node) (*schema.Node, error) {
	defaultMatch := schema.LooseObject()
	defaultMatch.Properties.Put("type", schema.StringEnum(C.RuleTypeDefault, ""))
	err := builder.FlattenStruct(defaultMatch, matchType)
	if err != nil {
		return nil, err
	}
	defaultVariant := &schema.Node{
		Type:                  "object",
		AllOf:                 []*schema.Node{defaultMatch, actionRef},
		UnevaluatedProperties: false,
	}

	logicalMatch := schema.LooseObject()
	logicalMatch.Properties.Put("type", schema.StringConst(C.RuleTypeLogical))
	logicalProperties(logicalMatch, nestedRef)
	logicalMatch.Required = []string{"type", "mode", "rules"}
	logicalVariant := &schema.Node{
		Type:                  "object",
		AllOf:                 []*schema.Node{logicalMatch, actionRef},
		UnevaluatedProperties: false,
	}

	return schema.OneOf(defaultVariant, logicalVariant), nil
}

// nestedRuleUnion builds a match-only rule schema: nested rules reject rule
// actions, and headless rules never carry them.
func nestedRuleUnion(builder schema.Builder, matchType reflect.Type, selfName string) (*schema.Node, error) {
	defaultVariant := schema.StrictObject()
	defaultVariant.Properties.Put("type", schema.StringEnum(C.RuleTypeDefault, ""))
	err := builder.FlattenStruct(defaultVariant, matchType)
	if err != nil {
		return nil, err
	}

	logicalVariant := schema.StrictObject()
	logicalVariant.Properties.Put("type", schema.StringConst(C.RuleTypeLogical))
	logicalProperties(logicalVariant, schema.RefNode(selfName))
	logicalVariant.Required = []string{"type", "mode", "rules"}

	return schema.OneOf(defaultVariant, logicalVariant), nil
}

func logicalProperties(node *schema.Node, nestedRef *schema.Node) {
	node.Properties.Put("mode", schema.StringEnum(C.LogicalTypeAnd, C.LogicalTypeOr))
	node.Properties.Put("rules", &schema.Node{Type: "array", Items: nestedRef})
	node.Properties.Put("invert", schema.BooleanNode())
}

// @Doc action:
//
//	@Required
//
//	See [Rule Actions](../rule_action/) for details.
//
// @Doc outbound:
//
//	@Failure: Deprecated in sing-box 1.11.0
//
//	    Moved to [Rule Action](../rule_action#route).
type RawDefaultRule struct {
	// Tags of [Inbound](/configuration/inbound/).
	Inbound badoption.Listable[string] `json:"inbound,omitempty" reference:"inbound"`

	// 4 or 6.
	//
	// Not limited if empty.
	IPVersion int `json:"ip_version,omitempty" enum:"4,6"`

	// Match network type.
	//
	// `tcp`, `udp` or `icmp`.
	//
	// @Quote: Changes in sing-box 1.13.0
	//
	//     Since sing-box 1.13.0, you can match ICMP echo (ping) requests via the new `icmp` network.
	//
	//     Such traffic originates from `TUN`, `WireGuard`, and `Tailscale` inbounds and can be routed to `Direct`, `WireGuard`, and `Tailscale` outbounds.
	Network badoption.Listable[string] `json:"network,omitempty" enum:"tcp,udp,icmp"`

	// Username, see each inbound for details.
	AuthUser badoption.Listable[string] `json:"auth_user,omitempty"`

	// Sniffed protocol, see [Protocol Sniff](/configuration/route/sniff/) for details.
	Protocol badoption.Listable[string] `json:"protocol,omitempty" enum:"tls,http,quic,dns,stun,bittorrent,dtls,ssh,rdp,ntp"`

	// Sniffed client type, see [Protocol Sniff](/configuration/route/sniff/) for details.
	//
	// @Since 1.10.0
	Client badoption.Listable[string] `json:"client,omitempty"`

	// Match full domain.
	Domain badoption.Listable[string] `json:"domain,omitempty"`

	// Match domain suffix.
	DomainSuffix badoption.Listable[string] `json:"domain_suffix,omitempty"`

	// Match domain using keyword.
	DomainKeyword badoption.Listable[string] `json:"domain_keyword,omitempty"`

	// Match domain using regular expression.
	DomainRegex badoption.Listable[string] `json:"domain_regex,omitempty"`

	// Match geosite.
	//
	// @Failure: Deprecated in sing-box 1.8.0
	//
	//     Geosite is deprecated and will be removed in sing-box 1.12.0, check [Migration](/migration/#migrate-geosite-to-rule-sets).
	Geosite badoption.Listable[string] `json:"geosite,omitempty" schema:"omit"`

	// Match source geoip.
	//
	// @Failure: Deprecated in sing-box 1.8.0
	//
	//     GeoIP is deprecated and will be removed in sing-box 1.12.0, check [Migration](/migration/#migrate-geoip-to-rule-sets).
	SourceGeoIP badoption.Listable[string] `json:"source_geoip,omitempty" schema:"omit"`

	// Match geoip.
	//
	// @Failure: Deprecated in sing-box 1.8.0
	//
	//     GeoIP is deprecated and will be removed in sing-box 1.12.0, check [Migration](/migration/#migrate-geoip-to-rule-sets).
	GeoIP badoption.Listable[string] `json:"geoip,omitempty" schema:"omit"`

	// Match source IP CIDR.
	SourceIPCIDR badoption.Listable[*badoption.Prefixable] `json:"source_ip_cidr,omitempty"`

	// Match non-public source IP.
	//
	// @Since 1.8.0
	SourceIPIsPrivate bool `json:"source_ip_is_private,omitempty"`

	// Match IP CIDR.
	IPCIDR badoption.Listable[*badoption.Prefixable] `json:"ip_cidr,omitempty"`

	// Match non-public IP.
	//
	// @Since 1.8.0
	IPIsPrivate bool `json:"ip_is_private,omitempty"`

	// Match source port.
	SourcePort badoption.Listable[uint16] `json:"source_port,omitempty"`

	// Match source port range.
	SourcePortRange badoption.Listable[string] `json:"source_port_range,omitempty"`

	// Match port.
	Port badoption.Listable[uint16] `json:"port,omitempty"`

	// Match port range.
	PortRange badoption.Listable[string] `json:"port_range,omitempty"`

	// Match process name.
	//
	// @Quote:
	//
	//     Only supported on Linux, Windows, and macOS.
	ProcessName badoption.Listable[string] `json:"process_name,omitempty"`

	// Match process path.
	//
	// @Quote:
	//
	//     Only supported on Linux, Windows, and macOS.
	ProcessPath badoption.Listable[string] `json:"process_path,omitempty"`

	// Match process path using regular expression.
	//
	// @Since 1.10.0
	//
	// @Quote:
	//
	//     Only supported on Linux, Windows, and macOS.
	ProcessPathRegex badoption.Listable[string] `json:"process_path_regex,omitempty"`

	// Match android package name.
	PackageName badoption.Listable[string] `json:"package_name,omitempty"`

	// Match android package name using regular expression.
	//
	// @Since 1.14.0
	PackageNameRegex badoption.Listable[string] `json:"package_name_regex,omitempty"`

	// Match user name.
	//
	// @Quote:
	//
	//     Only supported on Linux.
	User badoption.Listable[string] `json:"user,omitempty"`

	// Match user id.
	//
	// @Quote:
	//
	//     Only supported on Linux.
	UserID badoption.Listable[int32] `json:"user_id,omitempty"`

	// Match Clash mode.
	ClashMode string `json:"clash_mode,omitempty"`

	// Match network type.
	//
	// Available values: `wifi`, `cellular`, `ethernet` and `other`.
	//
	// @Since 1.11.0
	//
	// @Quote:
	//
	//     Only supported in graphical clients on Android and Apple platforms.
	NetworkType badoption.Listable[InterfaceType] `json:"network_type,omitempty"`

	// Match if network is considered Metered (on Android) or considered expensive,
	// such as Cellular or a Personal Hotspot (on Apple platforms).
	//
	// @Since 1.11.0
	//
	// @Quote:
	//
	//     Only supported in graphical clients on Android and Apple platforms.
	NetworkIsExpensive bool `json:"network_is_expensive,omitempty"`

	// Match if network is in Low Data Mode.
	//
	// @Since 1.11.0
	//
	// @Quote:
	//
	//     Only supported in graphical clients on Apple platforms.
	NetworkIsConstrained bool `json:"network_is_constrained,omitempty"`

	// Match WiFi SSID.
	//
	// See [Wi-Fi State](/configuration/shared/wifi-state/) for details.
	WIFISSID badoption.Listable[string] `json:"wifi_ssid,omitempty"`

	// Match WiFi BSSID.
	//
	// See [Wi-Fi State](/configuration/shared/wifi-state/) for details.
	WIFIBSSID badoption.Listable[string] `json:"wifi_bssid,omitempty"`

	// Match interface address.
	//
	// @Since 1.13.0
	//
	// @Quote:
	//
	//     Only supported on Linux, Windows, and macOS.
	InterfaceAddress *badjson.TypedMap[string, badoption.Listable[*badoption.Prefixable]] `json:"interface_address,omitempty"`

	// Matches network interface (same values as `network_type`) address.
	//
	// @Since 1.13.0
	//
	// @Quote:
	//
	//     Only supported in graphical clients on Android and Apple platforms.
	NetworkInterfaceAddress *badjson.TypedMap[InterfaceType, badoption.Listable[*badoption.Prefixable]] `json:"network_interface_address,omitempty"`

	// Match default interface address.
	//
	// @Since 1.13.0
	//
	// @Quote:
	//
	//     Only supported on Linux, Windows, and macOS.
	DefaultInterfaceAddress badoption.Listable[*badoption.Prefixable] `json:"default_interface_address,omitempty"`

	// Match source device MAC address.
	//
	// @Since 1.14.0
	//
	// @Quote:
	//
	//     Only supported on Linux, macOS, or in graphical clients on Android and macOS. See [Neighbor Resolution](/configuration/shared/neighbor/) for setup.
	SourceMACAddress badoption.Listable[string] `json:"source_mac_address,omitempty"`

	// Match source device hostname from DHCP leases.
	//
	// @Since 1.14.0
	//
	// @Quote:
	//
	//     Only supported on Linux, macOS, or in graphical clients on Android and macOS. See [Neighbor Resolution](/configuration/shared/neighbor/) for setup.
	SourceHostname badoption.Listable[string] `json:"source_hostname,omitempty"`

	// Match specified outbounds' preferred routes.
	//
	// | Type        | Match                                              |
	// |-------------|----------------------------------------------------|
	// | `tailscale` | Match MagicDNS domains and peers' allowed IPs      |
	// | `wireguard` | Match peers's allowed IPs                          |
	// | `bridge`    | Match all addresses except local addresses of the machine, only in [pre-match](/configuration/shared/pre-match/) |
	//
	// @Since 1.13.0
	PreferredBy badoption.Listable[string] `json:"preferred_by,omitempty"`

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
	//
	// @Since 1.15.0
	DNSServerAddress *badjson.TypedMap[string, badoption.Listable[*badoption.Prefixable]] `json:"dns_server_address,omitempty"`

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
	//
	// @Since 1.15.0
	DNSSearchDomain *badjson.TypedMap[string, badoption.Listable[string]] `json:"dns_search_domain,omitempty"`

	// Match [rule-set](/configuration/route/#rule_set).
	//
	// @Since 1.8.0
	RuleSet badoption.Listable[string] `json:"rule_set,omitempty" reference:"rule_set"`

	// Make `ip_cidr` in rule-sets match the source IP.
	//
	// @Since 1.10.0
	RuleSetIPCIDRMatchSource bool `json:"rule_set_ip_cidr_match_source,omitempty"`

	// Invert match result.
	Invert bool `json:"invert,omitempty"`

	// Make `ip_cidr` in rule-sets match the source IP.
	//
	// @Failure: Deprecated in sing-box 1.10.0
	//
	//     `rule_set_ipcidr_match_source` is renamed to `rule_set_ip_cidr_match_source` and will be remove in sing-box 1.11.0.
	//
	// Deprecated: renamed to rule_set_ip_cidr_match_source
	//
	// @Since 1.8.0
	Deprecated_RulesetIPCIDRMatchSource bool `json:"rule_set_ipcidr_match_source,omitempty" schema:"omit"`
}

type DefaultRule struct {
	RawDefaultRule
	RuleAction
}

func (r DefaultRule) MarshalJSON() ([]byte, error) {
	return badjson.MarshallObjects(r.RawDefaultRule, r.RuleAction)
}

func (r *DefaultRule) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &r.RawDefaultRule)
	if err != nil {
		return err
	}
	return badjson.UnmarshallExcluded(data, &r.RawDefaultRule, &r.RuleAction)
}

func (r DefaultRule) IsValid() bool {
	var defaultValue DefaultRule
	defaultValue.Invert = r.Invert
	return !reflect.DeepEqual(r, defaultValue)
}

// @Doc type:
//
//	`logical`
type RawLogicalRule struct {
	// `and` or `or`
	//
	// @Required
	Mode string `json:"mode" enum:"and,or"`

	// Included rules.
	//
	// @Required
	Rules []Rule `json:"rules,omitempty"`

	// Invert match result.
	Invert bool `json:"invert,omitempty"`
}

type LogicalRule struct {
	RawLogicalRule
	RuleAction
}

func (r LogicalRule) MarshalJSON() ([]byte, error) {
	return badjson.MarshallObjects(r.RawLogicalRule, r.RuleAction)
}

func (r *LogicalRule) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &r.RawLogicalRule)
	if err != nil {
		return err
	}
	return badjson.UnmarshallExcluded(data, &r.RawLogicalRule, &r.RuleAction)
}

func rulePayloadWithoutType(ctx context.Context, data []byte) ([]byte, error) {
	var content badjson.JSONObject
	err := content.UnmarshalJSONContext(ctx, data)
	if err != nil {
		return nil, err
	}
	content.Remove("type")
	return content.MarshalJSONContext(ctx)
}

func unmarshalDefaultRuleContext(ctx context.Context, data []byte, rule *DefaultRule) error {
	rawAction, routeOptions, err := inspectRouteRuleAction(ctx, data)
	if err != nil {
		return err
	}
	err = rejectNestedRouteRuleAction(ctx, data)
	if err != nil {
		return err
	}
	depth := nestedRuleDepth(ctx)
	err = json.UnmarshalContext(ctx, data, &rule.RawDefaultRule)
	if err != nil {
		return err
	}
	err = badjson.UnmarshallExcludedContext(ctx, data, &rule.RawDefaultRule, &rule.RuleAction)
	if err != nil {
		return err
	}
	if depth > 0 && rawAction == "" && routeOptions == (RouteActionOptions{}) {
		rule.RuleAction = RuleAction{}
	}
	return nil
}

func unmarshalLogicalRuleContext(ctx context.Context, data []byte, rule *LogicalRule) error {
	rawAction, routeOptions, err := inspectRouteRuleAction(ctx, data)
	if err != nil {
		return err
	}
	err = rejectNestedRouteRuleAction(ctx, data)
	if err != nil {
		return err
	}
	depth := nestedRuleDepth(ctx)
	err = json.UnmarshalContext(nestedRuleChildContext(ctx), data, &rule.RawLogicalRule)
	if err != nil {
		return err
	}
	err = badjson.UnmarshallExcludedContext(ctx, data, &rule.RawLogicalRule, &rule.RuleAction)
	if err != nil {
		return err
	}
	if depth > 0 && rawAction == "" && routeOptions == (RouteActionOptions{}) {
		rule.RuleAction = RuleAction{}
	}
	return nil
}

func (r *LogicalRule) IsValid() bool {
	return len(r.Rules) > 0 && common.All(r.Rules, Rule.IsValid)
}
