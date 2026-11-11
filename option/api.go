package option

import (
	"reflect"

	"github.com/sagernet/sing-box/schema"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badoption"
)

type APIServiceOptions struct {
	ListenOptions
	// Secret for the API.
	//
	// Clients authenticate with the standard `authorization: Bearer <secret>` gRPC metadata header.
	//
	// If empty, authentication is disabled.
	Secret string `json:"secret,omitempty"`
	// CORS allowed origins, `*` will be used if empty.
	AccessControlAllowOrigin badoption.Listable[string] `json:"access_control_allow_origin,omitempty" examples:"http://sing-box-dashboard.sagernet.org/,https://sing-box-dashboard.sagernet.org/"`
	// Allow access from private network.
	AccessControlAllowPrivateNetwork bool `json:"access_control_allow_private_network,omitempty"`
	// Web dashboard downloaded and served over the API listener at `/dashboard/`; other browser
	// requests are redirected to it.
	//
	// @info:
	//
	//     The object can be replaced with a boolean value (equivalent to `{ "enabled": <bool> }`),
	//     or with a string path (equivalent to `{ "enabled": true, "path": "<string>" }`).
	//
	// ##### enabled
	//
	// Enable the dashboard.
	//
	// ##### path
	//
	// Directory the dashboard files are stored in.
	//
	// `dashboard` in the working directory will be used by default.
	//
	// If the directory contains files not downloaded by sing-box, they are served as-is.
	//
	// ##### download_url
	//
	// Download URL of the dashboard archive (zip).
	//
	// `https://github.com/SagerNet/sing-box-dashboard/archive/refs/heads/gh-pages.zip` will be used by default.
	//
	// ##### http_client
	//
	// HTTP client used to download the dashboard.
	//
	// See [HTTP Client Fields](/configuration/shared/http-client/) for details.
	//
	// ##### update_interval
	//
	// Update interval of the dashboard.
	//
	// `1d` will be used by default.
	Dashboard *APIDashboardOptions `json:"dashboard,omitempty"`
	InboundTLSOptionsContainer
}

type _APIDashboardOptions struct {
	Enabled        bool               `json:"enabled,omitempty"`
	Path           string             `json:"path,omitempty"`
	DownloadURL    string             `json:"download_url,omitempty"`
	HTTPClient     *HTTPClientOptions `json:"http_client,omitempty"`
	UpdateInterval badoption.Duration `json:"update_interval,omitempty"`
}

type APIDashboardOptions _APIDashboardOptions

func (o APIDashboardOptions) MarshalJSON() ([]byte, error) {
	if o.DownloadURL == "" && o.HTTPClient == nil && o.UpdateInterval == 0 {
		if o.Path == "" {
			return json.Marshal(o.Enabled)
		}
		if o.Enabled {
			return json.Marshal(o.Path)
		}
	}
	return json.Marshal(_APIDashboardOptions(o))
}

func (o *APIDashboardOptions) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, &o.Enabled)
	if err == nil {
		return nil
	}
	err = json.Unmarshal(bytes, &o.Path)
	if err == nil {
		o.Enabled = true
		return nil
	}
	return json.UnmarshalDisallowUnknownFields(bytes, (*_APIDashboardOptions)(o))
}

func (o APIDashboardOptions) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	objectForm := schema.StrictObject()
	err := builder.FlattenStruct(objectForm, reflect.TypeFor[APIDashboardOptions]())
	if err != nil {
		return nil, err
	}
	return schema.AnyOf(schema.BooleanNode(), schema.StringNode(), objectForm), nil
}
