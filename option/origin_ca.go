package option

import (
	"strings"

	"github.com/sagernet/sing-box/schema"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badoption"
)

type CloudflareOriginCACertificateProviderOptions struct {
	// @required
	//
	// List of domain names or wildcard domain names to include in the certificate.
	Domain badoption.Listable[string] `json:"domain,omitempty"`
	// Root directory used to store the issued certificate, private key, and metadata.
	//
	// If empty, sing-box uses the same default data directory as the ACME certificate provider:
	// `$XDG_DATA_HOME/certmagic` or `$HOME/.local/share/certmagic`.
	DataDirectory string `json:"data_directory,omitempty"`
	// Cloudflare API token used to create the certificate.
	//
	// Get or create one in [Cloudflare Dashboard > My Profile > API Tokens](https://dash.cloudflare.com/profile/api-tokens).
	//
	// Requires the `Zone / SSL and Certificates / Edit` permission.
	//
	// Conflict with `origin_ca_key`.
	APIToken string `json:"api_token,omitempty"`
	// Cloudflare Origin CA Key.
	//
	// Get it in [Cloudflare Dashboard > My Profile > API Tokens > API Keys > Origin CA Key](https://dash.cloudflare.com/profile/api-tokens).
	//
	// Conflict with `api_token`.
	OriginCAKey string `json:"origin_ca_key,omitempty"`
	// The signature type to request from Cloudflare.
	//
	// | Value                | Type        |
	// |----------------------|-------------|
	// | `origin-rsa`         | RSA         |
	// | `origin-ecc`         | ECDSA P-256 |
	//
	// `origin-rsa` is used if empty.
	RequestType CloudflareOriginCARequestType `json:"request_type,omitempty" enum:"origin-rsa,origin-ecc"`
	// The requested certificate validity in days.
	//
	// Available values: `7`, `30`, `90`, `365`, `730`, `1095`, `5475`.
	//
	// `5475` days (15 years) is used if empty.
	RequestedValidity CloudflareOriginCARequestValidity `json:"requested_validity,omitempty" enum:"0,7,30,90,365,730,1095,5475"`
	// HTTP Client for all provider HTTP requests.
	//
	// See [HTTP Client Fields](/configuration/shared/http-client/) for details.
	HTTPClient *HTTPClientOptions `json:"http_client,omitempty"`
}

type CloudflareOriginCARequestType string

const (
	CloudflareOriginCARequestTypeOriginRSA = CloudflareOriginCARequestType("origin-rsa")
	CloudflareOriginCARequestTypeOriginECC = CloudflareOriginCARequestType("origin-ecc")
)

func (t *CloudflareOriginCARequestType) UnmarshalJSON(data []byte) error {
	var value string
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}
	value = strings.ToLower(value)
	switch CloudflareOriginCARequestType(value) {
	case "", CloudflareOriginCARequestTypeOriginRSA, CloudflareOriginCARequestTypeOriginECC:
		*t = CloudflareOriginCARequestType(value)
	default:
		return E.New("unsupported Cloudflare Origin CA request type: ", value)
	}
	return nil
}

func (t CloudflareOriginCARequestType) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return schema.StringEnum("", "origin-rsa", "origin-ecc"), nil
}

type CloudflareOriginCARequestValidity uint16

const (
	CloudflareOriginCARequestValidity7    = CloudflareOriginCARequestValidity(7)
	CloudflareOriginCARequestValidity30   = CloudflareOriginCARequestValidity(30)
	CloudflareOriginCARequestValidity90   = CloudflareOriginCARequestValidity(90)
	CloudflareOriginCARequestValidity365  = CloudflareOriginCARequestValidity(365)
	CloudflareOriginCARequestValidity730  = CloudflareOriginCARequestValidity(730)
	CloudflareOriginCARequestValidity1095 = CloudflareOriginCARequestValidity(1095)
	CloudflareOriginCARequestValidity5475 = CloudflareOriginCARequestValidity(5475)
)

func (v *CloudflareOriginCARequestValidity) UnmarshalJSON(data []byte) error {
	var value uint16
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}
	switch CloudflareOriginCARequestValidity(value) {
	case 0,
		CloudflareOriginCARequestValidity7,
		CloudflareOriginCARequestValidity30,
		CloudflareOriginCARequestValidity90,
		CloudflareOriginCARequestValidity365,
		CloudflareOriginCARequestValidity730,
		CloudflareOriginCARequestValidity1095,
		CloudflareOriginCARequestValidity5475:
		*v = CloudflareOriginCARequestValidity(value)
	default:
		return E.New("unsupported Cloudflare Origin CA requested validity: ", value)
	}
	return nil
}

func (v CloudflareOriginCARequestValidity) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return &schema.Node{Type: "integer", Enum: []any{0, 7, 30, 90, 365, 730, 1095, 5475}}, nil
}
