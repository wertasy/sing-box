package option

import (
	"reflect"
	"strings"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/schema"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badjson"
	"github.com/sagernet/sing/common/json/badoption"
)

type ACMECertificateProviderOptions struct {
	// @required
	//
	// List of domains.
	Domain badoption.Listable[string] `json:"domain,omitempty"`
	// The directory to store ACME data.
	//
	// `$XDG_DATA_HOME/certmagic|$HOME/.local/share/certmagic` will be used if empty.
	DataDirectory string `json:"data_directory,omitempty"`
	// Server name to use when choosing a certificate if the ClientHello's ServerName field is empty.
	DefaultServerName string `json:"default_server_name,omitempty"`
	// The email address to use when creating or selecting an existing ACME server account.
	Email string `json:"email,omitempty"`
	// The ACME CA provider to use.
	//
	// | Value                   | Provider      |
	// |-------------------------|---------------|
	// | `letsencrypt (default)` | Let's Encrypt |
	// | `zerossl`               | ZeroSSL       |
	// | `https://...`           | Custom        |
	//
	// When `provider` is `zerossl`, sing-box will automatically request ZeroSSL EAB credentials if `email` is set and `external_account` is empty.
	//
	// When `provider` is `zerossl`, at least one of `external_account`, `email`, or `account_key` is required.
	Provider string `json:"provider,omitempty"`
	// @since 1.14.0
	//
	// The PEM-encoded private key of an existing ACME account.
	AccountKey string `json:"account_key,omitempty"`
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
	// External account bindings are used to associate an ACME account with an existing account in a non-ACME system, such as a CA customer database.
	//
	// To enable ACME account binding, the CA operating the ACME server needs to provide the ACME client with a MAC key and a key identifier, using some mechanism outside of ACME. §7.3.4
	ExternalAccount *ACMEExternalAccountOptions `json:"external_account,omitempty"`
	// ACME DNS01 challenge field. If configured, other challenge methods will be disabled.
	//
	// See [DNS01 Challenge Fields](/configuration/shared/dns01_challenge/) for details.
	DNS01Challenge *ACMEProviderDNS01ChallengeOptions `json:"dns01_challenge,omitempty"`
	// @since 1.14.0
	//
	// The private key type to generate for new certificates.
	//
	// | Value      | Type    |
	// |------------|---------|
	// | `ed25519`  | Ed25519 |
	// | `p256`     | P-256   |
	// | `p384`     | P-384   |
	// | `rsa2048`  | RSA     |
	// | `rsa4096`  | RSA     |
	KeyType ACMEKeyType `json:"key_type,omitempty" enum:"ed25519,p256,p384,rsa2048,rsa4096"`
	// @since 1.14.0
	//
	// The ACME profile to use for certificate issuance.
	//
	// When empty and `provider` is Let's Encrypt, `shortlived` will be used automatically if any domain is an IP address.
	Profile string `json:"profile,omitempty"`
	// @since 1.14.0
	//
	// HTTP Client for all provider HTTP requests.
	//
	// See [HTTP Client Fields](/configuration/shared/http-client/) for details.
	HTTPClient *HTTPClientOptions `json:"http_client,omitempty"`
}

type _ACMEProviderDNS01ChallengeOptions struct {
	AbstractACMEProviderDNS01ChallengeOptions
	Provider          string                     `json:"provider,omitempty" enum:"alidns,cloudflare,acmedns"`
	AliDNSOptions     ACMEDNS01AliDNSOptions     `json:"-"`
	CloudflareOptions ACMEDNS01CloudflareOptions `json:"-"`
	ACMEDNSOptions    ACMEDNS01ACMEDNSOptions    `json:"-"`
}

type AbstractACMEProviderDNS01ChallengeOptions struct {
	TTL                badoption.Duration     `json:"ttl,omitempty"`
	PropagationDelay   badoption.Duration     `json:"propagation_delay,omitempty"`
	PropagationTimeout badoption.Duration     `json:"propagation_timeout,omitempty"`
	Resolvers          LegacyListable[string] `json:"resolvers,omitempty"`
	OverrideDomain     string                 `json:"override_domain,omitempty"`
}

type ACMEProviderDNS01ChallengeOptions _ACMEProviderDNS01ChallengeOptions

func (o ACMEProviderDNS01ChallengeOptions) MarshalJSON() ([]byte, error) {
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
		return nil, E.New("unknown provider type: ", o.Provider)
	}
	return badjson.MarshallObjects(_ACMEProviderDNS01ChallengeOptions(o), v)
}

func (o *ACMEProviderDNS01ChallengeOptions) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, (*_ACMEProviderDNS01ChallengeOptions)(o))
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
	case "":
		return E.New("missing provider type")
	default:
		return E.New("unknown provider type: ", o.Provider)
	}
	return badjson.UnmarshallExcluded(bytes, (*_ACMEProviderDNS01ChallengeOptions)(o), v)
}

func (o ACMEProviderDNS01ChallengeOptions) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return builder.Define("ACMEProviderDNS01Challenge", func() (*schema.Node, error) {
		return schema.DiscriminatedUnion(builder, "provider", true, acmeDNS01Variants(), func(variant *schema.Node) error {
			return builder.FlattenStruct(variant, reflect.TypeFor[AbstractACMEProviderDNS01ChallengeOptions]())
		})
	})
}

type ACMEKeyType string

const (
	ACMEKeyTypeED25519 = ACMEKeyType("ed25519")
	ACMEKeyTypeP256    = ACMEKeyType("p256")
	ACMEKeyTypeP384    = ACMEKeyType("p384")
	ACMEKeyTypeRSA2048 = ACMEKeyType("rsa2048")
	ACMEKeyTypeRSA4096 = ACMEKeyType("rsa4096")
)

func (t *ACMEKeyType) UnmarshalJSON(data []byte) error {
	var value string
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}
	value = strings.ToLower(value)
	switch ACMEKeyType(value) {
	case "", ACMEKeyTypeED25519, ACMEKeyTypeP256, ACMEKeyTypeP384, ACMEKeyTypeRSA2048, ACMEKeyTypeRSA4096:
		*t = ACMEKeyType(value)
	default:
		return E.New("unknown ACME key type: ", value)
	}
	return nil
}

func (t ACMEKeyType) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return schema.StringEnum("", "ed25519", "p256", "p384", "rsa2048", "rsa4096"), nil
}
