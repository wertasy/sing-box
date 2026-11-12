package option

import (
	"reflect"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/schema"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badjson"
	"github.com/sagernet/sing/common/json/badoption"
)

type InboundACMEOptions struct {
	// List of domain.
	//
	// ACME will be disabled if empty.
	Domain badoption.Listable[string] `json:"domain,omitempty"`
	// The directory to store ACME data.
	//
	// `$XDG_DATA_HOME/certmagic|$HOME/.local/share/certmagic` will be used if empty.
	DataDirectory string `json:"data_directory,omitempty"`
	// Server name to use when choosing a certificate if the ClientHello's ServerName field is empty.
	DefaultServerName string `json:"default_server_name,omitempty"`
	// The email address to use when creating or selecting an existing ACME server account
	Email string `json:"email,omitempty"`
	// The ACME CA provider to use.
	//
	// | Value                   | Provider      |
	// |-------------------------|---------------|
	// | `letsencrypt (default)` | Let's Encrypt |
	// | `zerossl`               | ZeroSSL       |
	// | `https://...`           | Custom        |
	Provider string `json:"provider,omitempty"`
	// Disable all HTTP challenges.
	DisableHTTPChallenge bool `json:"disable_http_challenge,omitempty"`
	// Disable all TLS-ALPN challenges
	DisableTLSALPNChallenge bool `json:"disable_tls_alpn_challenge,omitempty"`
	// The alternate port to use for the ACME HTTP challenge; if non-empty, this port will be used instead of 80 to spin up a listener for the HTTP challenge.
	AlternativeHTTPPort uint16 `json:"alternative_http_port,omitempty"`
	// The alternate port to use for the ACME TLS-ALPN challenge; the system must forward 443 to this port for challenge to succeed.
	AlternativeTLSPort uint16 `json:"alternative_tls_port,omitempty"`
	// EAB (External Account Binding) contains information necessary to bind or map an ACME account to some other account known by the CA.
	//
	// External account bindings are "used to associate an ACME account with an existing account in a non-ACME system, such as a CA customer database.
	//
	// To enable ACME account binding, the CA operating the ACME server needs to provide the ACME client with a MAC key and a key identifier, using some mechanism outside of ACME. §7.3.4
	ExternalAccount *ACMEExternalAccountOptions `json:"external_account,omitempty"`
	// ACME DNS01 challenge field. If configured, other challenge methods will be disabled.
	//
	// See [DNS01 Challenge Fields](/configuration/shared/dns01_challenge/) for details.
	DNS01Challenge *ACMEDNS01ChallengeOptions `json:"dns01_challenge,omitempty"`
	// The ACME profile to use for certificate issuance.
	//
	// When empty and `provider` is Let's Encrypt, `shortlived` will be used automatically if any domain is an IP address.
	//
	// @Since 1.14.0
	Profile string `json:"profile,omitempty"`
}

type ACMEExternalAccountOptions struct {
	// The key identifier.
	KeyID string `json:"key_id,omitempty"`
	// The MAC key.
	MACKey string `json:"mac_key,omitempty"`
}

type _ACMEDNS01ChallengeOptions struct {
	Provider          string                     `json:"provider,omitempty" enum:"alidns,cloudflare,acmedns"`
	AliDNSOptions     ACMEDNS01AliDNSOptions     `json:"-"`
	CloudflareOptions ACMEDNS01CloudflareOptions `json:"-"`
	ACMEDNSOptions    ACMEDNS01ACMEDNSOptions    `json:"-"`
}

type ACMEDNS01ChallengeOptions _ACMEDNS01ChallengeOptions

func (o ACMEDNS01ChallengeOptions) MarshalJSON() ([]byte, error) {
	var v any
	switch o.Provider {
	case C.DNSProviderAliDNS:
		v = o.AliDNSOptions
	case C.DNSProviderCloudflare:
		v = o.CloudflareOptions
	case C.DNSProviderACMEDNS:
		v = o.ACMEDNSOptions
	case "":
		return nil, E.New("missing provider type")
	default:
		return nil, E.New("unknown provider type: " + o.Provider)
	}
	return badjson.MarshallObjects(_ACMEDNS01ChallengeOptions(o), v)
}

func (o *ACMEDNS01ChallengeOptions) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, (*_ACMEDNS01ChallengeOptions)(o))
	if err != nil {
		return err
	}
	var v any
	switch o.Provider {
	case C.DNSProviderAliDNS:
		v = &o.AliDNSOptions
	case C.DNSProviderCloudflare:
		v = &o.CloudflareOptions
	case C.DNSProviderACMEDNS:
		v = &o.ACMEDNSOptions
	default:
		return E.New("unknown provider type: " + o.Provider)
	}
	err = badjson.UnmarshallExcluded(bytes, (*_ACMEDNS01ChallengeOptions)(o), v)
	if err != nil {
		return err
	}
	return nil
}

func acmeDNS01Variants() []schema.UnionVariant {
	return []schema.UnionVariant{
		{Value: C.DNSProviderAliDNS, StructType: reflect.TypeFor[ACMEDNS01AliDNSOptions]()},
		{Value: C.DNSProviderCloudflare, StructType: reflect.TypeFor[ACMEDNS01CloudflareOptions]()},
		{Value: C.DNSProviderACMEDNS, StructType: reflect.TypeFor[ACMEDNS01ACMEDNSOptions]()},
	}
}

func (o ACMEDNS01ChallengeOptions) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return builder.Define("ACMEDNS01Challenge", func() (*schema.Node, error) {
		return schema.DiscriminatedUnion(builder, "provider", true, acmeDNS01Variants(), nil)
	})
}

type ACMEDNS01AliDNSOptions struct {
	AccessKeyID     string `json:"access_key_id,omitempty"`
	AccessKeySecret string `json:"access_key_secret,omitempty"`
	RegionID        string `json:"region_id,omitempty"`
	SecurityToken   string `json:"security_token,omitempty"`
}

type ACMEDNS01CloudflareOptions struct {
	APIToken  string `json:"api_token,omitempty"`
	ZoneToken string `json:"zone_token,omitempty"`
}

type ACMEDNS01ACMEDNSOptions struct {
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	Subdomain string `json:"subdomain,omitempty"`
	ServerURL string `json:"server_url,omitempty"`
}
