package option

import (
	"net/url"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/sagernet/sing-box/common/ipset"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/schema"
	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/domain"
	E "github.com/sagernet/sing/common/exceptions"
	F "github.com/sagernet/sing/common/format"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badjson"
	"github.com/sagernet/sing/common/json/badoption"
)

type _RuleSet struct {
	Type          string                     `json:"type,omitempty" enum:"inline,local,remote"`
	Tag           badoption.Listable[string] `json:"tag"`
	Format        string                     `json:"format,omitempty" enum:"source,binary"`
	InlineOptions PlainRuleSet               `json:"-"`
	LocalOptions  LocalRuleSet               `json:"-"`
	RemoteOptions RemoteRuleSet              `json:"-"`
}

type RuleSet _RuleSet

func (r RuleSet) MarshalJSON() ([]byte, error) {
	if r.Type != C.RuleSetTypeInline {
		var defaultFormat string
		switch r.Type {
		case C.RuleSetTypeLocal:
			defaultFormat = ruleSetDefaultFormat(r.LocalOptions.Path)
		case C.RuleSetTypeRemote:
			defaultFormat = ruleSetDefaultFormat(r.RemoteOptions.URL)
		}
		if r.Format == defaultFormat {
			r.Format = ""
		}
	}
	var v any
	switch r.Type {
	case "", C.RuleSetTypeInline:
		r.Type = ""
		v = r.InlineOptions
	case C.RuleSetTypeLocal:
		v = r.LocalOptions
	case C.RuleSetTypeRemote:
		v = r.RemoteOptions
	default:
		return nil, E.New("unknown rule-set type: " + r.Type)
	}
	return badjson.MarshallObjects(_RuleSet(r), v)
}

func (r *RuleSet) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, (*_RuleSet)(r))
	if err != nil {
		return err
	}
	if len(r.Tag) == 0 || common.Any(r.Tag, func(tag string) bool { return tag == "" }) {
		return E.New("missing tag")
	}
	var v any
	switch r.Type {
	case "", C.RuleSetTypeInline:
		r.Type = C.RuleSetTypeInline
		v = &r.InlineOptions
	case C.RuleSetTypeLocal:
		v = &r.LocalOptions
	case C.RuleSetTypeRemote:
		v = &r.RemoteOptions
	default:
		return E.New("unknown rule-set type: " + r.Type)
	}
	err = badjson.UnmarshallExcluded(bytes, (*_RuleSet)(r), v)
	if err != nil {
		return err
	}
	if r.Type != C.RuleSetTypeInline {
		if r.Format == "" {
			switch r.Type {
			case C.RuleSetTypeLocal:
				r.Format = ruleSetDefaultFormat(r.LocalOptions.Path)
			case C.RuleSetTypeRemote:
				r.Format = ruleSetDefaultFormat(r.RemoteOptions.URL)
			}
		}
		switch r.Format {
		case "":
			return E.New("missing format")
		case C.RuleSetFormatSource, C.RuleSetFormatBinary:
		default:
			return E.New("unknown rule-set format: " + r.Format)
		}
	} else {
		r.Format = ""
	}
	if len(r.Tag) > 1 {
		switch r.Type {
		case C.RuleSetTypeInline:
			return E.New("inline rule-set does not support multiple tags")
		case C.RuleSetTypeLocal:
			if !strings.Contains(r.LocalOptions.Path, C.RuleSetTagPlaceholder) {
				return E.New("missing ", C.RuleSetTagPlaceholder, " placeholder in path")
			}
		case C.RuleSetTypeRemote:
			if !strings.Contains(r.RemoteOptions.URL, C.RuleSetTagPlaceholder) {
				return E.New("missing ", C.RuleSetTagPlaceholder, " placeholder in url")
			}
			if r.RemoteOptions.InitialPath != "" && !strings.Contains(r.RemoteOptions.InitialPath, C.RuleSetTagPlaceholder) {
				return E.New("missing ", C.RuleSetTagPlaceholder, " placeholder in initial_path")
			}
		}
	}
	return nil
}

func ruleSetDefaultFormat(path string) string {
	if pathURL, err := url.Parse(path); err == nil {
		path = pathURL.Path
	}
	switch filepath.Ext(path) {
	case ".json":
		return C.RuleSetFormatSource
	case ".srs":
		return C.RuleSetFormatBinary
	default:
		return ""
	}
}

func (r RuleSet) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return builder.Define("RuleSet", func() (*schema.Node, error) {
		headlessRef, err := builder.Describe(reflect.TypeFor[HeadlessRule]())
		if err != nil {
			return nil, err
		}
		tagNode := schema.ListableOf(schema.StringNode())
		formatNode := schema.StringEnum(C.RuleSetFormatSource, C.RuleSetFormatBinary)

		inlineVariant := schema.StrictObject()
		inlineVariant.Properties.Put("type", schema.StringEnum(C.RuleSetTypeInline, ""))
		inlineVariant.Properties.Put("tag", tagNode)
		inlineVariant.Properties.Put("rules", &schema.Node{Type: "array", Items: headlessRef})
		inlineVariant.Required = []string{"tag"}

		localVariant := schema.StrictObject()
		localVariant.Properties.Put("type", schema.StringConst(C.RuleSetTypeLocal))
		localVariant.Properties.Put("tag", tagNode)
		localVariant.Properties.Put("format", formatNode)
		err = builder.FlattenStruct(localVariant, reflect.TypeFor[LocalRuleSet]())
		if err != nil {
			return nil, err
		}
		localVariant.Required = []string{"type", "tag"}

		remoteVariant := schema.StrictObject()
		remoteVariant.Properties.Put("type", schema.StringConst(C.RuleSetTypeRemote))
		remoteVariant.Properties.Put("tag", tagNode)
		remoteVariant.Properties.Put("format", formatNode)
		err = builder.FlattenStruct(remoteVariant, reflect.TypeFor[RemoteRuleSet]())
		if err != nil {
			return nil, err
		}
		remoteVariant.Required = []string{"type", "tag"}

		return schema.OneOf(inlineVariant, localVariant, remoteVariant), nil
	})
}

type LocalRuleSet struct {
	Path string `json:"path,omitempty"`
}

type RemoteRuleSet struct {
	URL            string             `json:"url"`
	InitialPath    string             `json:"initial_path,omitempty"`
	HTTPClient     *HTTPClientOptions `json:"http_client,omitempty"`
	UpdateInterval badoption.Duration `json:"update_interval,omitempty"`

	// Deprecated: use http_client instead
	DownloadDetour string `json:"download_detour,omitempty" reference:"outbound" schema:"omit"`
}

type _HeadlessRule struct {
	Type           string              `json:"type,omitempty" enum:"default,logical"`
	DefaultOptions DefaultHeadlessRule `json:"-"`
	LogicalOptions LogicalHeadlessRule `json:"-"`
}

type HeadlessRule _HeadlessRule

func (r HeadlessRule) MarshalJSON() ([]byte, error) {
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
	return badjson.MarshallObjects(_HeadlessRule(r), v)
}

func (r *HeadlessRule) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, (*_HeadlessRule)(r))
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
	err = badjson.UnmarshallExcluded(bytes, (*_HeadlessRule)(r), v)
	if err != nil {
		return err
	}
	return nil
}

func (r HeadlessRule) IsValid() bool {
	switch r.Type {
	case C.RuleTypeDefault, "":
		return r.DefaultOptions.IsValid()
	case C.RuleTypeLogical:
		return r.LogicalOptions.IsValid()
	default:
		panic("unknown rule type: " + r.Type)
	}
}

func (r HeadlessRule) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return builder.Define("HeadlessRule", func() (*schema.Node, error) {
		return nestedRuleUnion(builder, reflect.TypeFor[DefaultHeadlessRule](), "HeadlessRule")
	})
}

type DefaultHeadlessRule struct {
	// DNS query type. Values can be integers or type name strings.
	//
	// @Quote: Changes in sing-box 1.14.0
	//
	//     When a DNS rule references this rule-set, this field now also applies
	//     when the DNS rule is matched from an internal domain resolution that
	//     does not target a specific DNS server. See
	//     [Migration](/migration/#ip_version-and-query_type-behavior-changes-in-dns-rules).
	//
	//     When a DNS rule references a rule-set containing this field, the DNS
	//     rule is incompatible in the same DNS configuration with Legacy Address
	//     Filter Fields in DNS rules, the Legacy `strategy` DNS rule action
	//     option, and the Legacy `rule_set_ip_cidr_accept_empty` DNS rule item.
	QueryType badoption.Listable[DNSQueryType] `json:"query_type,omitempty"`

	// `tcp` or `udp`.
	Network badoption.Listable[string] `json:"network,omitempty" enum:"tcp,udp,icmp"`

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

	// Match IP CIDR.
	//
	// @Info:
	//
	//     `ip_cidr` is an alias for `source_ip_cidr` when `rule_set_ipcidr_match_source` enabled in route/DNS rules.
	IPCIDR badoption.Listable[*badoption.Prefixable] `json:"ip_cidr,omitempty"`

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
	// @Quote:
	//
	//     Only supported in graphical clients on Android and Apple platforms.
	WIFISSID badoption.Listable[string] `json:"wifi_ssid,omitempty"`

	// Match WiFi BSSID.
	//
	// @Quote:
	//
	//     Only supported in graphical clients on Android and Apple platforms.
	WIFIBSSID badoption.Listable[string] `json:"wifi_bssid,omitempty"`

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

	// Invert match result.
	Invert bool `json:"invert,omitempty"`

	DomainMatcher *domain.Matcher `json:"-"`
	SourceIPSet   *ipset.Set      `json:"-"`
	IPSet         *ipset.Set      `json:"-"`

	AdGuardDomain        badoption.Listable[string] `json:"-"`
	AdGuardDomainMatcher *domain.AdGuardMatcher     `json:"-"`
}

func (r DefaultHeadlessRule) IsValid() bool {
	var defaultValue DefaultHeadlessRule
	defaultValue.Invert = r.Invert
	return !reflect.DeepEqual(r, defaultValue)
}

// @Doc type:
//
//	`logical`
type LogicalHeadlessRule struct {
	// `and` or `or`
	//
	// @Required
	Mode string `json:"mode" enum:"and,or"`

	// Included rules.
	//
	// @Required
	Rules []HeadlessRule `json:"rules,omitempty"`

	// Invert match result.
	Invert bool `json:"invert,omitempty"`
}

func (r LogicalHeadlessRule) IsValid() bool {
	return len(r.Rules) > 0 && common.All(r.Rules, HeadlessRule.IsValid)
}

type _PlainRuleSetCompat struct {
	Version    uint8           `json:"version" enum:"1,2,3,4,5"`
	Options    PlainRuleSet    `json:"-"`
	RawMessage json.RawMessage `json:"-"`
}

type PlainRuleSetCompat _PlainRuleSetCompat

func (r PlainRuleSetCompat) MarshalJSON() ([]byte, error) {
	var v any
	switch r.Version {
	case C.RuleSetVersion1, C.RuleSetVersion2, C.RuleSetVersion3, C.RuleSetVersion4, C.RuleSetVersion5:
		v = r.Options
	default:
		return nil, E.New("unknown rule-set version: ", r.Version)
	}
	return badjson.MarshallObjects(_PlainRuleSetCompat(r), v)
}

func (r *PlainRuleSetCompat) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, (*_PlainRuleSetCompat)(r))
	if err != nil {
		return err
	}
	var v any
	switch r.Version {
	case C.RuleSetVersion1, C.RuleSetVersion2, C.RuleSetVersion3, C.RuleSetVersion4, C.RuleSetVersion5:
		v = &r.Options
	case 0:
		return E.New("missing rule-set version")
	default:
		return E.New("unknown rule-set version: ", r.Version)
	}
	err = badjson.UnmarshallExcluded(bytes, (*_PlainRuleSetCompat)(r), v)
	if err != nil {
		return err
	}
	r.RawMessage = bytes
	return nil
}

func (r PlainRuleSetCompat) Upgrade() (PlainRuleSet, error) {
	switch r.Version {
	case C.RuleSetVersion1, C.RuleSetVersion2, C.RuleSetVersion3, C.RuleSetVersion4, C.RuleSetVersion5:
	default:
		return PlainRuleSet{}, E.New("unknown rule-set version: " + F.ToString(r.Version))
	}
	return r.Options, nil
}

type PlainRuleSet struct {
	Rules []HeadlessRule `json:"rules,omitempty"`
}
