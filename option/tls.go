package option

import (
	"crypto/tls"
	"encoding/json"
	"strings"

	"github.com/sagernet/sing-box/schema"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json/badoption"
)

type InboundTLSOptions struct {
	//
	// Enable TLS.
	Enabled bool `json:"enabled,omitempty"`
	//
	// Used to verify the hostname on the returned certificates unless insecure is given.
	//
	// It is also included in the client's handshake to support virtual hosting unless it is an IP address.
	ServerName string `json:"server_name,omitempty"`
	// Accepts any server certificate.
	//
	// @Indicator: Client only
	Insecure bool `json:"insecure,omitempty"`
	//
	// List of supported application level protocols, in order of preference.
	//
	// If both peers support ALPN, the selected protocol will be one from this list, and the connection will fail if there is
	// no mutually supported protocol.
	//
	// See [Application-Layer Protocol Negotiation](https://en.wikipedia.org/wiki/Application-Layer_Protocol_Negotiation).
	ALPN badoption.Listable[string] `json:"alpn,omitempty" examples:"http/1.1,h2,h3"`
	//
	// The minimum TLS version that is acceptable.
	//
	// TLS 1.2 is used by default.
	MinVersion string `json:"min_version,omitempty" enum:"1.0,1.1,1.2,1.3"`
	//
	// The maximum TLS version that is acceptable.
	//
	// By default, the maximum version is currently TLS 1.3.
	MaxVersion string `json:"max_version,omitempty" enum:"1.0,1.1,1.2,1.3"`
	//
	// List of enabled TLS 1.0–1.2 cipher suites. The order of the list is ignored.
	// Note that TLS 1.3 cipher suites are not configurable.
	//
	// If empty, a safe default list is used. The default cipher suites might change over time.
	CipherSuites LegacyListable[string] `json:"cipher_suites,omitempty"`
	// Set of supported key exchange mechanisms. The order of the list is ignored, and key exchange mechanisms are chosen
	// from this list using an internal preference order by Golang.
	//
	// Available values, also the default list:
	//
	// * `P256`
	// * `P384`
	// * `P521`
	// * `X25519`
	// * `X25519MLKEM768`
	//
	// @Since 1.13.0
	CurvePreferences LegacyListable[CurvePreference] `json:"curve_preferences,omitempty"`
	//
	// Server certificates chain line array, in PEM format.
	Certificate badoption.Listable[string] `json:"certificate,omitempty"`
	// The path to server certificate chain, in PEM format.
	//
	// @Note:
	//
	//     Will be automatically reloaded if file modified.
	CertificatePath string `json:"certificate_path,omitempty"`
	// The type of client authentication to use.
	//
	// Available values:
	//
	// * `no` (default)
	// * `request`
	// * `require-any`
	// * `verify-if-given`
	// * `require-and-verify`
	//
	// One of `client_certificate`, `client_certificate_path`, `client_certificate_sha256`, or `client_certificate_public_key_sha256` is required
	// if this option is set to `verify-if-given`, or `require-and-verify`.
	//
	// @Since 1.13.0
	//
	// @Indicator: Server only
	ClientAuthentication ClientAuthType `json:"client_authentication,omitempty"`
	// Client certificate chain line array, in PEM format.
	//
	// Client certificate chain line array, in PEM format.
	//
	// @Since 1.13.0
	//
	// @Indicator: Client only
	//
	// @Since 1.13.0
	//
	// @Indicator: Server only
	ClientCertificate badoption.Listable[string] `json:"client_certificate,omitempty"`
	// The path to client certificate chain, in PEM format.
	//
	// List of path to client certificate chain, in PEM format.
	//
	// @Since 1.13.0
	//
	// @Indicator: Client only
	//
	// @Since 1.13.0
	//
	// @Indicator: Server only
	//
	// @Note:
	//
	//     Will be automatically reloaded if file modified.
	ClientCertificatePath badoption.Listable[string] `json:"client_certificate_path,omitempty"`
	// List of SHA-256 hashes of client certificates, in base64 format.
	//
	// The hash is computed over the whole DER-encoded certificate, see [certificate_sha256](#certificate_sha256).
	//
	// @Since 1.15.0
	//
	// @Indicator: Server only
	ClientCertificateSHA256 badoption.Listable[[]byte] `json:"client_certificate_sha256,omitempty"`
	// List of SHA-256 hashes of client certificate public keys, in base64 format.
	//
	// To generate the SHA-256 hash for a certificate's public key, use the following commands:
	//
	// ```bash
	// # For a certificate file
	// openssl x509 -in certificate.pem -pubkey -noout | openssl pkey -pubin -outform der | openssl dgst -sha256 -binary | openssl enc -base64
	//
	// # For a certificate from a remote server
	// echo | openssl s_client -servername example.com -connect example.com:443 2>/dev/null | openssl x509 -pubkey -noout | openssl pkey -pubin -outform der | openssl dgst -sha256 -binary | openssl enc -base64
	// ```
	//
	// @Since 1.13.0
	//
	// @Indicator: Server only
	ClientCertificatePublicKeySHA256 badoption.Listable[[]byte] `json:"client_certificate_public_key_sha256,omitempty"`
	// The server private key line array, in PEM format.
	//
	// @Indicator: Server only
	Key badoption.Listable[string] `json:"key,omitempty"`
	// The path to the server private key, in PEM format.
	//
	// @Indicator: Server only
	//
	// @Note:
	//
	//     Will be automatically reloaded if file modified.
	KeyPath string `json:"key_path,omitempty"`
	// Enable kernel TLS transmit support.
	//
	// @Since 1.13.0
	//
	// @Quote:
	//
	//     Only supported on Linux 5.1+, use a newer kernel if possible.
	//
	// @Quote:
	//
	//     Only TLS 1.3 is supported.
	//
	// @Warning:
	//
	//     kTLS TX may only improve performance when `splice(2)` is available (both ends must be TCP or TLS without additional protocols after handshake); otherwise, it will definitely degrade performance.
	KernelTx bool `json:"kernel_tx,omitempty"`
	// Enable kernel TLS receive support.
	//
	// @Failure:
	//
	//     kTLS RX will definitely degrade performance even if `splice(2)` is in use, so enabling it is not recommended.
	//
	// @Since 1.13.0
	//
	// @Quote:
	//
	//     Only supported on Linux 5.1+, use a newer kernel if possible.
	//
	// @Quote:
	//
	//     Only TLS 1.3 is supported.
	KernelRx bool `json:"kernel_rx,omitempty"`
	// TLS handshake timeout, in golang's Duration format.
	//
	// @Default: 15s
	//
	// @Since 1.14.0
	HandshakeTimeout badoption.Duration `json:"handshake_timeout,omitempty"`
	// A string or an object.
	//
	// When string, the tag of a shared [Certificate Provider](/configuration/shared/certificate-provider/).
	//
	// When object, an inline certificate provider. See [Certificate Provider](/configuration/shared/certificate-provider/) for available types and fields.
	//
	// ## Custom TLS support
	//
	// @Info: QUIC support
	//
	//     Only ECH is supported in QUIC.
	//
	// @Since 1.14.0
	//
	// @Indicator: Server only
	CertificateProvider *CertificateProviderOptions `json:"certificate_provider,omitempty"`

	// @Omit
	// @Deprecated: use certificate_provider
	ACME *InboundACMEOptions `json:"acme,omitempty" schema:"omit"`

	// @Omit
	ECH *InboundECHOptions `json:"ech,omitempty"`
	// @Omit
	Reality *InboundRealityOptions `json:"reality,omitempty"`
}

type ClientAuthType tls.ClientAuthType

func (t ClientAuthType) MarshalJSON() ([]byte, error) {
	var stringValue string
	switch t {
	case ClientAuthType(tls.NoClientCert):
		stringValue = "no"
	case ClientAuthType(tls.RequestClientCert):
		stringValue = "request"
	case ClientAuthType(tls.RequireAnyClientCert):
		stringValue = "require-any"
	case ClientAuthType(tls.VerifyClientCertIfGiven):
		stringValue = "verify-if-given"
	case ClientAuthType(tls.RequireAndVerifyClientCert):
		stringValue = "require-and-verify"
	default:
		return nil, E.New("unknown client authentication type: ", int(t))
	}
	return json.Marshal(stringValue)
}

func (t *ClientAuthType) UnmarshalJSON(data []byte) error {
	var stringValue string
	err := json.Unmarshal(data, &stringValue)
	if err != nil {
		return err
	}
	switch stringValue {
	case "no":
		*t = ClientAuthType(tls.NoClientCert)
	case "request":
		*t = ClientAuthType(tls.RequestClientCert)
	case "require-any":
		*t = ClientAuthType(tls.RequireAnyClientCert)
	case "verify-if-given":
		*t = ClientAuthType(tls.VerifyClientCertIfGiven)
	case "require-and-verify":
		*t = ClientAuthType(tls.RequireAndVerifyClientCert)
	default:
		return E.New("unknown client authentication type: ", stringValue)
	}
	return nil
}

func (t ClientAuthType) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return schema.StringEnum("no", "request", "require-any", "verify-if-given", "require-and-verify"), nil
}

type InboundTLSOptionsContainer struct {
	// TLS configuration, see [TLS](/configuration/shared/tls/#inbound).
	//
	TLS *InboundTLSOptions `json:"tls,omitempty"`
}

type InboundTLSOptionsWrapper interface {
	TakeInboundTLSOptions() *InboundTLSOptions
	ReplaceInboundTLSOptions(options *InboundTLSOptions)
}

func (o *InboundTLSOptionsContainer) TakeInboundTLSOptions() *InboundTLSOptions {
	return o.TLS
}

func (o *InboundTLSOptionsContainer) ReplaceInboundTLSOptions(options *InboundTLSOptions) {
	o.TLS = options
}

type OutboundTLSOptions struct {
	Enabled bool `json:"enabled,omitempty"`
	// TLS engine to use.
	//
	// Values:
	//
	// * `go` (default)
	// * `apple`
	// * `windows`
	//
	// Supported fields:
	//
	// * `server_name`
	// * `insecure`
	// * `alpn`
	// * `min_version`
	// * `max_version`
	// * `certificate` / `certificate_path`
	// * `certificate_sha256`
	// * `certificate_public_key_sha256`
	// * `handshake_timeout`
	//
	// Unsupported fields:
	//
	// * `disable_sni`
	// * `cipher_suites`
	// * `curve_preferences`
	// * `client_certificate` / `client_certificate_path` / `client_key` / `client_key_path`
	// * `fragment` / `record_fragment`
	// * `kernel_tx` / `kernel_rx`
	// * `ech`
	// * `utls`
	// * `reality`
	//
	// The default version range is TLS 1.2 to TLS 1.3, matching the `go` engine.
	//
	// @Since 1.14.0
	//
	// @Indicator: Client only
	//
	// @Note:
	//
	//     `windows` uses Schannel via SSPI. Only available on Windows build 17763 or later (Windows 10 version 1809, Windows Server 2019, or newer).
	//
	// @Note:
	//
	//     TLS 1.3 is only negotiated on Windows 11 or Windows Server 2022 and newer.
	Engine string `json:"engine,omitempty" enum:"go,apple,windows"`
	// Do not send server name in ClientHello.
	//
	// @Indicator: Client only
	DisableSNI       bool                            `json:"disable_sni,omitempty"`
	ServerName       string                          `json:"server_name,omitempty"`
	Insecure         bool                            `json:"insecure,omitempty"`
	ALPN             badoption.Listable[string]      `json:"alpn,omitempty" examples:"http/1.1,h2,h3"`
	MinVersion       string                          `json:"min_version,omitempty" enum:"1.0,1.1,1.2,1.3"`
	MaxVersion       string                          `json:"max_version,omitempty" enum:"1.0,1.1,1.2,1.3"`
	CipherSuites     LegacyListable[string]          `json:"cipher_suites,omitempty"`
	CurvePreferences LegacyListable[CurvePreference] `json:"curve_preferences,omitempty"`
	Certificate      badoption.Listable[string]      `json:"certificate,omitempty"`
	CertificatePath  string                          `json:"certificate_path,omitempty"`
	// List of SHA-256 hashes of server certificates, in base64 format.
	//
	// The hash is computed over the whole DER-encoded certificate. Use `certificate_public_key_sha256` to pin only the public key.
	//
	// To generate the SHA-256 hash for a certificate, use the following commands:
	//
	// ```bash
	// # For a certificate file
	// openssl x509 -in certificate.pem -outform der | openssl dgst -sha256 -binary | openssl enc -base64
	//
	// # For a certificate from a remote server
	// echo | openssl s_client -servername example.com -connect example.com:443 2>/dev/null | openssl x509 -outform der | openssl dgst -sha256 -binary | openssl enc -base64
	// ```
	//
	// @Since 1.15.0
	//
	// @Indicator: Client only
	CertificateSHA256 badoption.Listable[[]byte] `json:"certificate_sha256,omitempty"`
	// List of SHA-256 hashes of server certificate public keys, in base64 format.
	//
	// To generate the SHA-256 hash for a certificate's public key, use the following commands:
	//
	// ```bash
	// # For a certificate file
	// openssl x509 -in certificate.pem -pubkey -noout | openssl pkey -pubin -outform der | openssl dgst -sha256 -binary | openssl enc -base64
	//
	// # For a certificate from a remote server
	// echo | openssl s_client -servername example.com -connect example.com:443 2>/dev/null | openssl x509 -pubkey -noout | openssl pkey -pubin -outform der | openssl dgst -sha256 -binary | openssl enc -base64
	// ```
	//
	// @Since 1.13.0
	//
	// @Indicator: Client only
	CertificatePublicKeySHA256 badoption.Listable[[]byte] `json:"certificate_public_key_sha256,omitempty"`
	ClientCertificate          badoption.Listable[string] `json:"client_certificate,omitempty"`
	ClientCertificatePath      string                     `json:"client_certificate_path,omitempty"`
	// Client private key line array, in PEM format.
	//
	// @Since 1.13.0
	//
	// @Indicator: Client only
	ClientKey badoption.Listable[string] `json:"client_key,omitempty"`
	// The path to client private key, in PEM format.
	//
	// @Since 1.13.0
	//
	// @Indicator: Client only
	ClientKeyPath string `json:"client_key_path,omitempty"`
	// Fragment TLS handshakes to bypass firewalls.
	//
	// This feature is intended to circumvent simple firewalls based on **plaintext packet matching**,
	// and should not be used to circumvent real censorship.
	//
	// Due to poor performance, try `record_fragment` first, and only apply to server names known to be blocked.
	//
	// On Linux, Apple platforms, (administrator privileges required) Windows,
	// the wait time can be automatically detected. Otherwise, it will fall back to
	// waiting for a fixed time specified by `fragment_fallback_delay`.
	//
	// In addition, if the actual wait time is less than 20ms, it will also fall back to waiting for a fixed time,
	// because the target is considered to be local or behind a transparent proxy.
	//
	// @Since 1.12.0
	//
	// @Indicator: Client only
	Fragment bool `json:"fragment,omitempty"`
	// The fallback value used when TLS segmentation cannot automatically determine the wait time.
	//
	// @Default: 500ms
	//
	// @Since 1.12.0
	//
	// @Indicator: Client only
	FragmentFallbackDelay badoption.Duration `json:"fragment_fallback_delay,omitempty"`
	// Fragment TLS handshake into multiple TLS records to bypass firewalls.
	//
	// @Since 1.12.0
	//
	// @Indicator: Client only
	RecordFragment bool `json:"record_fragment,omitempty"`
	// Inject a forged TLS ClientHello carrying a whitelisted SNI before the real one,
	// to fool SNI-filtering middleboxes that permit specific hostnames.
	//
	// Requires `CAP_NET_RAW` and `CAP_NET_ADMIN` on Linux, root on macOS, and
	// Administrator on Windows. Windows on ARM64 is not supported.
	//
	// @Since 1.14.0
	//
	// @Indicator: Client only
	//
	// @Quote:
	//
	//     Only supported on Linux, macOS, and Windows, and requires elevated privileges.
	Spoof string `json:"spoof,omitempty"`
	// How the forged segment is rejected by the real server.
	//
	// | Value                      | Behavior                                                                                                       |
	// |----------------------------|----------------------------------------------------------------------------------------------------------------|
	// | `wrong-sequence` (default) | The forged segment's TCP sequence number is placed before the server's receive window.                         |
	// | `wrong-checksum`           | The forged segment's TCP checksum is deliberately invalid.                                                     |
	// | `wrong-ack`                | The forged segment's TCP acknowledgment number is placed before the server's send window.                      |
	// | `wrong-md5`                | The forged segment carries a TCP-MD5 signature option.                                                         |
	// | `wrong-timestamp`          | The forged segment carries a backdated TCP timestamp. Linux/Windows only; not supported on macOS.              |
	//
	// @Since 1.14.0
	//
	// @Indicator: Client only
	SpoofMethod      string             `json:"spoof_method,omitempty" enum:"wrong-sequence,wrong-checksum,wrong-ack,wrong-md5,wrong-timestamp"`
	KernelTx         bool               `json:"kernel_tx,omitempty"`
	KernelRx         bool               `json:"kernel_rx,omitempty"`
	HandshakeTimeout badoption.Duration `json:"handshake_timeout,omitempty"`
	// @Omit
	ECH *OutboundECHOptions `json:"ech,omitempty"`
	// uTLS is a fork of "crypto/tls", which provides ClientHello fingerprinting resistance.
	//
	// Available fingerprint values:
	//
	// * chrome
	// * firefox
	// * edge
	// * safari
	// * 360
	// * qq
	// * ios
	// * android
	// * random
	// * randomized
	//
	// Chrome fingerprint will be used if empty.
	//
	// @Failure: Not Recommended
	//
	//     uTLS has had repeated fingerprinting vulnerabilities discovered by researchers.
	//
	//     uTLS is a Go library that attempts to imitate browser TLS fingerprints by copying
	//     ClientHello structure. However, browsers use completely different TLS stacks
	//     (Chrome uses BoringSSL, Firefox uses NSS) with distinct implementation behaviors
	//     that cannot be replicated by simply copying the handshake format, making detection possible.
	//     Additionally, the library lacks active maintenance and has poor code quality,
	//     making it unsuitable for censorship circumvention.
	//
	//     For TLS fingerprint resistance, use [NaiveProxy](/configuration/inbound/naive/) instead.
	//
	// @Indicator: Client only
	//
	// @Warning: Removed since sing-box 1.10.0
	//
	//     Some legacy chrome fingerprints have been removed and will fallback to chrome:
	//
	//     :material-close: chrome_psk
	//     :material-close: chrome_psk_shuffle
	//     :material-close: chrome_padding_psk_shuffle
	//     :material-close: chrome_pq
	//     :material-close: chrome_pq_psk
	UTLS *OutboundUTLSOptions `json:"utls,omitempty"`
	// @Omit
	Reality *OutboundRealityOptions `json:"reality,omitempty"`
}

type OutboundTLSOptionsContainer struct {
	// TLS configuration, see [TLS](/configuration/shared/tls/#outbound).
	//
	TLS *OutboundTLSOptions `json:"tls,omitempty"`
}

type OutboundTLSOptionsWrapper interface {
	TakeOutboundTLSOptions() *OutboundTLSOptions
	ReplaceOutboundTLSOptions(options *OutboundTLSOptions)
}

func (o *OutboundTLSOptionsContainer) TakeOutboundTLSOptions() *OutboundTLSOptions {
	return o.TLS
}

func (o *OutboundTLSOptionsContainer) ReplaceOutboundTLSOptions(options *OutboundTLSOptions) {
	o.TLS = options
}

type CurvePreference tls.CurveID

const (
	CurveP256      = 23
	CurveP384      = 24
	CurveP521      = 25
	X25519         = 29
	X25519MLKEM768 = 4588
)

func (c CurvePreference) MarshalJSON() ([]byte, error) {
	var stringValue string
	switch c {
	case CurvePreference(CurveP256):
		stringValue = "P256"
	case CurvePreference(CurveP384):
		stringValue = "P384"
	case CurvePreference(CurveP521):
		stringValue = "P521"
	case CurvePreference(X25519):
		stringValue = "X25519"
	case CurvePreference(X25519MLKEM768):
		stringValue = "X25519MLKEM768"
	default:
		return nil, E.New("unknown curve id: ", int(c))
	}
	return json.Marshal(stringValue)
}

func (c *CurvePreference) UnmarshalJSON(data []byte) error {
	var stringValue string
	err := json.Unmarshal(data, &stringValue)
	if err != nil {
		return err
	}
	switch strings.ToUpper(stringValue) {
	case "P256":
		*c = CurvePreference(CurveP256)
	case "P384":
		*c = CurvePreference(CurveP384)
	case "P521":
		*c = CurvePreference(CurveP521)
	case "X25519":
		*c = CurvePreference(X25519)
	case "X25519MLKEM768":
		*c = CurvePreference(X25519MLKEM768)
	default:
		return E.New("unknown curve name: ", stringValue)
	}
	return nil
}

func (c CurvePreference) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return schema.StringEnum("P256", "P384", "P521", "X25519", "X25519MLKEM768"), nil
}

type InboundRealityOptions struct {
	// @Omit
	Enabled bool `json:"enabled,omitempty"`
	// Handshake server address and [Dial Fields](/configuration/shared/dial/).
	//
	// @Indicator: Server only
	//
	// @Required
	Handshake InboundRealityHandshakeOptions `json:"handshake,omitempty"`
	// Private key, generated by `sing-box generate reality-keypair`.
	//
	// @Indicator: Server only
	//
	// @Required
	PrivateKey string `json:"private_key,omitempty"`
	// A hexadecimal string with zero to eight digits.
	//
	// @Required
	ShortID badoption.Listable[string] `json:"short_id,omitempty"`
	// The maximum time difference between the server and the client.
	//
	// Check disabled if empty.
	//
	// @Indicator: Server only
	MaxTimeDifference badoption.Duration `json:"max_time_difference,omitempty"`
}

type InboundRealityHandshakeOptions struct {
	ServerOptions
	DialerOptions
}

type InboundECHOptions struct {
	// @Omit
	Enabled bool `json:"enabled,omitempty"`
	// ECH key line array, in PEM format.
	//
	// @Indicator: Server only
	Key badoption.Listable[string] `json:"key,omitempty"`
	// The path to ECH key, in PEM format.
	//
	// @Indicator: Server only
	//
	// @Note:
	//
	//     Will be automatically reloaded if file modified.
	KeyPath string `json:"key_path,omitempty"`

	// @Deprecated: not supported by stdlib
	PQSignatureSchemesEnabled bool `json:"pq_signature_schemes_enabled,omitempty" schema:"omit"`
	// @Deprecated: added by fault
	DynamicRecordSizingDisabled bool `json:"dynamic_record_sizing_disabled,omitempty" schema:"omit"`
}

type OutboundECHOptions struct {
	// @Omit
	Enabled bool `json:"enabled,omitempty"`
	// ECH configuration line array, in PEM format.
	//
	// If empty, load from DNS will be attempted.
	//
	// @Indicator: Client only
	Config badoption.Listable[string] `json:"config,omitempty"`
	// The path to ECH configuration, in PEM format.
	//
	// If empty, load from DNS will be attempted.
	//
	// @Indicator: Client only
	ConfigPath string `json:"config_path,omitempty"`
	// Overrides the domain name used for ECH HTTPS record queries.
	//
	// If empty, `server_name` is used for queries.
	//
	// @Since 1.13.0
	//
	// @Indicator: Client only
	QueryServerName string `json:"query_server_name,omitempty"`

	// Enable support for post-quantum peer certificate signature schemes.
	//
	// @Failure: Deprecated in sing-box 1.12.0
	//
	//     `pq_signature_schemes_enabled` is deprecated in sing-box 1.12.0 and removed in sing-box 1.13.0.
	//
	// @Deprecated: not supported by stdlib
	PQSignatureSchemesEnabled bool `json:"pq_signature_schemes_enabled,omitempty" schema:"omit"`
	// Disables adaptive sizing of TLS records.
	//
	// When true, the largest possible TLS record size is always used.
	// When false, the size of TLS records may be adjusted in an attempt to improve latency.
	//
	// @Failure: Deprecated in sing-box 1.12.0
	//
	//     `dynamic_record_sizing_disabled` is deprecated in sing-box 1.12.0 and removed in sing-box 1.13.0.
	//
	// @Deprecated: added by fault
	DynamicRecordSizingDisabled bool `json:"dynamic_record_sizing_disabled,omitempty" schema:"omit"`
}

type OutboundUTLSOptions struct {
	Enabled     bool   `json:"enabled,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty" enum:"chrome_psk,chrome_psk_shuffle,chrome_padding_psk_shuffle,chrome_pq,chrome_pq_psk,chrome,firefox,edge,safari,360,qq,ios,android,random,randomized"`
}

type OutboundRealityOptions struct {
	// @Omit
	Enabled bool `json:"enabled,omitempty"`
	// Public key, generated by `sing-box generate reality-keypair`.
	//
	// @Indicator: Client only
	//
	// @Required
	PublicKey string `json:"public_key,omitempty"`
	ShortID   string `json:"short_id,omitempty"`
}
