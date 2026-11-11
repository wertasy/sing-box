package option

import (
	"net/url"
	"reflect"
	"slices"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/schema"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badjson"
	"github.com/sagernet/sing/common/json/badoption"
	M "github.com/sagernet/sing/common/metadata"
)

type Hysteria2InboundOptions struct {
	ListenOptions
	UpMbps   int            `json:"up_mbps,omitempty"`
	DownMbps int            `json:"down_mbps,omitempty"`
	Obfs     *Hysteria2Obfs `json:"obfs,omitempty"`
	// Hysteria2 users
	Users []Hysteria2User `json:"users,omitempty"`
	// *When `up_mbps` and `down_mbps` are not set*:
	//
	// Commands clients to use the BBR CC instead of Hysteria CC.
	//
	// *When `up_mbps` and `down_mbps` are set*:
	//
	// Deny clients to use the BBR CC.
	IgnoreClientBandwidth bool `json:"ignore_client_bandwidth,omitempty"`
	InboundTLSOptionsContainer
	QUICOptions
	// HTTP3 server behavior (URL string configuration) when authentication fails.
	//
	// | Scheme       | Example                 | Description        |
	// |--------------|-------------------------|--------------------|
	// | `file`       | `file:///var/www`       | As a file server   |
	// | `http/https` | `http://127.0.0.1:8080` | As a reverse proxy |
	//
	// Conflict with `masquerade.type`.
	//
	// A 404 page will be returned if masquerade is not configured.
	Masquerade *Hysteria2Masquerade `json:"masquerade,omitempty"`
	// @since 1.14.0
	//
	// BBR congestion control algorithm profile, one of `conservative` `standard` `aggressive`.
	//
	// `standard` is used by default.
	BBRProfile string `json:"bbr_profile,omitempty" enum:"standard,conservative,aggressive"`
	// Enable debug information logging for Hysteria Brutal CC.
	BrutalDebug bool `json:"brutal_debug,omitempty"`
	// @since 1.14.0
	//
	// Register this inbound to a Hysteria Realm rendezvous service to enable NAT traversal.
	//
	// See [Hysteria Realm](/configuration/service/hysteria-realm/) for the rendezvous service.
	Realm *Hysteria2InboundRealm `json:"realm,omitempty"`
}

type Hysteria2Realm struct {
	// @required
	//
	// Realm rendezvous service URL.
	ServerURL string `json:"server_url"`
	// Bearer token for the realm. Must match one of `users[].token` configured on the realm.
	Token string `json:"token,omitempty"`
	// @required
	//
	// The same slot identifier the target Hysteria2 server registered.
	RealmID string `json:"realm_id"`
	// @required
	//
	// List of STUN servers (`host` or `host:port`) used to discover this client's public addresses.
	//
	// Domain names are resolved using [`domain_resolver`](/configuration/shared/dial/#domain_resolver) from Dial Fields.
	STUNServers LegacyListable[string] `json:"stun_servers"`
	// Restrict realm connections (STUN, hole punching, and the resulting QUIC path) to a single IP version.
	//
	// `4` or `6`. Both are used if empty.
	IPVersion int `json:"ip_version,omitempty" enum:"0,4,6"`
	// Maintain a UDP port mapping on the local gateway via UPnP or NAT-PMP.
	//
	// Improves hole-punching reliability behind gateways that support it.
	//
	// Requires IPv4: conflicts with `"ip_version": 6`.
	PortMapping *Hysteria2RealmPortMapping `json:"port_mapping,omitempty"`
	// HTTP client used to talk to the realm.
	//
	// See [HTTP Client](/configuration/shared/http-client/) for details.
	HTTPClient *HTTPClientOptions `json:"http_client,omitempty"`
}

func (r Hysteria2Realm) STUNServersIsDomain() bool {
	return slices.ContainsFunc(r.STUNServers, func(server string) bool {
		return M.ParseSocksaddr(server).IsDomain()
	})
}

type Hysteria2RealmPortMapping struct {
	// Enable port mapping.
	Enabled bool `json:"enabled,omitempty"`
	// Timeout for gateway discovery and mapping operations.
	//
	// `10s` is used by default.
	Timeout badoption.Duration `json:"timeout,omitempty"`
	// Lease lifetime of the mapping.
	//
	// `10m` is used by default.
	Lifetime badoption.Duration `json:"lifetime,omitempty"`
}

type Hysteria2InboundRealm struct {
	Hysteria2Realm
	// Set domain resolver to use for resolving STUN server domain names.
	//
	// This option uses the same format as the [route DNS rule action](/configuration/dns/rule_action/#route) without the `action` field.
	//
	// Setting this option directly to a string is equivalent to setting `server` of this options.
	//
	// If empty, the default domain resolver is used.
	STUNDomainResolver *DomainResolveOptions `json:"stun_domain_resolver,omitempty"`
}

type Hysteria2ObfsGecko struct {
	// @since 1.14.0
	//
	// Minimum on-wire packet size in bytes. Gecko only.
	//
	// `512` is used by default.
	MinPacketSize int `json:"min_packet_size,omitempty"`
	// @since 1.14.0
	//
	// Maximum on-wire packet size in bytes. Gecko only.
	//
	// `1200` is used by default.
	MaxPacketSize int `json:"max_packet_size,omitempty"`
}

type _Hysteria2Obfs struct {
	// QUIC traffic obfuscator type, one of `salamander` `gecko`.
	//
	// Disabled if empty.
	Type string `json:"type,omitempty" enum:"salamander,gecko"`
	// QUIC traffic obfuscator password.
	Password     string             `json:"password,omitempty"`
	GeckoOptions Hysteria2ObfsGecko `json:"-"`
}

type Hysteria2Obfs _Hysteria2Obfs

func (o Hysteria2Obfs) MarshalJSON() ([]byte, error) {
	var v any
	switch o.Type {
	case C.Hysteria2ObfsTypeSalamander:
	case C.Hysteria2ObfsTypeGecko:
		v = o.GeckoOptions
	default:
		return nil, E.New("unknown obfs type: ", o.Type)
	}
	return badjson.MarshallObjects(_Hysteria2Obfs(o), v)
}

func (o *Hysteria2Obfs) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, (*_Hysteria2Obfs)(o))
	if err != nil {
		return err
	}
	var v any
	switch o.Type {
	case C.Hysteria2ObfsTypeSalamander:
	case C.Hysteria2ObfsTypeGecko:
		v = &o.GeckoOptions
	default:
		return E.New("unknown obfs type: ", o.Type)
	}
	if v == nil {
		return nil
	}
	return badjson.UnmarshallExcluded(bytes, (*_Hysteria2Obfs)(o), v)
}

func (o Hysteria2Obfs) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return schema.DiscriminatedUnion(builder, "type", true, []schema.UnionVariant{
		{Value: C.Hysteria2ObfsTypeSalamander},
		{Value: C.Hysteria2ObfsTypeGecko, StructType: reflect.TypeFor[Hysteria2ObfsGecko]()},
	}, func(variant *schema.Node) error {
		variant.Properties.Put("password", schema.StringNode())
		return nil
	})
}

type Hysteria2User struct {
	Name string `json:"name,omitempty"`
	// Authentication password
	Password string `json:"password,omitempty"`
}

type _Hysteria2Masquerade struct {
	// HTTP3 server behavior (Object configuration) when authentication fails.
	//
	// | Type     | Description                 | Fields                              |
	// |----------|-----------------------------|-------------------------------------|
	// | `file`   | As a file server            | `directory`                         |
	// | `proxy`  | As a reverse proxy          | `url`, `rewrite_host`               |
	// | `string` | Reply with a fixed response | `status_code`, `headers`, `content` |
	//
	// Conflict with `masquerade`.
	//
	// A 404 page will be returned if masquerade is not configured.
	Type          string                    `json:"type,omitempty" enum:"file,proxy,string"`
	FileOptions   Hysteria2MasqueradeFile   `json:"-"`
	ProxyOptions  Hysteria2MasqueradeProxy  `json:"-"`
	StringOptions Hysteria2MasqueradeString `json:"-"`
}

type Hysteria2Masquerade _Hysteria2Masquerade

func (m Hysteria2Masquerade) MarshalJSON() ([]byte, error) {
	var v any
	switch m.Type {
	case C.Hysterai2MasqueradeTypeFile:
		v = m.FileOptions
	case C.Hysterai2MasqueradeTypeProxy:
		v = m.ProxyOptions
	case C.Hysterai2MasqueradeTypeString:
		v = m.StringOptions
	default:
		return nil, E.New("unknown masquerade type: ", m.Type)
	}
	return badjson.MarshallObjects(_Hysteria2Masquerade(m), v)
}

func (m *Hysteria2Masquerade) UnmarshalJSON(bytes []byte) error {
	var urlString string
	err := json.Unmarshal(bytes, &urlString)
	if err == nil {
		masqueradeURL, err := url.Parse(urlString)
		if err != nil {
			return E.Cause(err, "invalid masquerade URL")
		}
		switch masqueradeURL.Scheme {
		case "file":
			m.Type = C.Hysterai2MasqueradeTypeFile
			m.FileOptions.Directory = masqueradeURL.Path
		case "http", "https":
			m.Type = C.Hysterai2MasqueradeTypeProxy
			m.ProxyOptions.URL = urlString
		default:
			return E.New("unknown masquerade URL scheme: ", masqueradeURL.Scheme)
		}
		return nil
	}
	err = json.Unmarshal(bytes, (*_Hysteria2Masquerade)(m))
	if err != nil {
		return err
	}
	var v any
	switch m.Type {
	case C.Hysterai2MasqueradeTypeFile:
		v = &m.FileOptions
	case C.Hysterai2MasqueradeTypeProxy:
		v = &m.ProxyOptions
	case C.Hysterai2MasqueradeTypeString:
		v = &m.StringOptions
	default:
		return E.New("unknown masquerade type: ", m.Type)
	}
	return badjson.UnmarshallExcluded(bytes, (*_Hysteria2Masquerade)(m), v)
}

func (m Hysteria2Masquerade) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	union, err := schema.DiscriminatedUnion(builder, "type", true, []schema.UnionVariant{
		{Value: C.Hysterai2MasqueradeTypeFile, StructType: reflect.TypeFor[Hysteria2MasqueradeFile]()},
		{Value: C.Hysterai2MasqueradeTypeProxy, StructType: reflect.TypeFor[Hysteria2MasqueradeProxy]()},
		{Value: C.Hysterai2MasqueradeTypeString, StructType: reflect.TypeFor[Hysteria2MasqueradeString]()},
	}, nil)
	if err != nil {
		return nil, err
	}
	return schema.AnyOf(schema.StringNode(), union), nil
}

type Hysteria2MasqueradeFile struct {
	// File server root directory.
	Directory string `json:"directory"`
}

type Hysteria2MasqueradeProxy struct {
	// Reverse proxy target URL.
	URL string `json:"url"`
	// Rewrite the `Host` header to the target URL.
	RewriteHost bool `json:"rewrite_host,omitempty"`
}

type Hysteria2MasqueradeString struct {
	// Fixed response status code.
	StatusCode int `json:"status_code,omitempty"`
	// Fixed response headers.
	Headers badoption.HTTPHeader `json:"headers,omitempty"`
	// Fixed response content.
	Content string `json:"content"`
}

type Hysteria2OutboundOptions struct {
	DialerOptions
	ServerOptions
	// @since 1.11.0
	//
	// Server port range list.
	//
	// Conflicts with `server_port` and `realm`.
	ServerPorts LegacyListable[string] `json:"server_ports,omitempty"`
	// @since 1.11.0
	//
	// Port hopping interval.
	//
	// `30s` is used by default.
	HopInterval badoption.Duration `json:"hop_interval,omitempty"`
	// @since 1.14.0
	//
	// Maximum port hopping interval, used for randomization.
	//
	// If set, the actual hop interval will be randomly chosen between `hop_interval` and `hop_interval_max`.
	HopIntervalMax badoption.Duration `json:"hop_interval_max,omitempty"`
	UpMbps         int                `json:"up_mbps,omitempty"`
	DownMbps       int                `json:"down_mbps,omitempty"`
	Obfs           *Hysteria2Obfs     `json:"obfs,omitempty"`
	// Authentication password.
	Password string `json:"password,omitempty"`
	// Enabled network
	//
	// One of `tcp` `udp`.
	//
	// Both is enabled by default.
	Network NetworkList `json:"network,omitempty"`
	OutboundTLSOptionsContainer
	QUICOptions
	// @since 1.14.0
	//
	// BBR congestion control algorithm profile, one of `conservative` `standard` `aggressive`.
	//
	// `standard` is used by default.
	BBRProfile string `json:"bbr_profile,omitempty" enum:"standard,conservative,aggressive"`
	// Enable debug information logging for Hysteria Brutal CC.
	BrutalDebug bool `json:"brutal_debug,omitempty"`
	// @since 1.14.0
	//
	// Disable Chrome QUIC fingerprint parroting.
	//
	// If it is not disabled, the client's QUIC handshake is made to parrot Chrome's, so that Hysteria traffic
	// is harder to identify by handshake fingerprinting.
	//
	// To match Chrome, the client uses Chrome's own QUIC parameters, which override some settings:
	// `idle_timeout` is fixed at 30 seconds, `max_concurrent_streams` and `initial_packet_size` are replaced by
	// Chrome's values, and the receive windows start at Chrome's initial values before growing to the configured
	// maximums.
	//
	// @warning:
	//
	//     Chrome does not declare support for Ed25519, so a server using an Ed25519 certificate will fail the
	//     handshake. Use an ECDSA or RSA certificate instead.
	DisableChromeParrot bool `json:"disable_chrome_parrot,omitempty"`
	// @since 1.14.0
	//
	// Connect to a Hysteria2 server through a Hysteria Realm rendezvous service.
	//
	// Conflicts with `server`, `server_port` and `server_ports`.
	//
	// The TLS SNI defaults to the host portion of `server_url`. Set `tls.server_name` to match the certificate the Hysteria2 server presents.
	//
	// See [Hysteria Realm](/configuration/service/hysteria-realm/) for the rendezvous service.
	Realm *Hysteria2Realm `json:"realm,omitempty"`
}

type HysteriaRealmUser struct {
	// @required
	//
	// Username, used in logs and as the quota key.
	Name string `json:"name"`
	// @required
	//
	// Bearer token presented by Hysteria2 inbounds and outbounds via `Authorization: Bearer <token>`.
	Token string `json:"token"`
	// Maximum number of realm slots this user may hold concurrently.
	MaxRealms int `json:"max_realms,omitempty"`
}

type HysteriaRealmServiceOptions struct {
	ListenOptions
	InboundTLSOptionsContainer
	HTTP2Options
	// @required
	//
	// Authorized users.
	Users []HysteriaRealmUser `json:"users"`
}
