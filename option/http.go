package option

import (
	"context"
	"reflect"
	"slices"

	"github.com/sagernet/sing-box/schema"
	"github.com/sagernet/sing/common/byteformats"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badjson"
	"github.com/sagernet/sing/common/json/badoption"
)

type HTTP2Options struct {
	// Idle connection timeout, in golang's Duration format.
	IdleTimeout badoption.Duration `json:"idle_timeout,omitempty"`
	// Keep alive period, in golang's Duration format.
	KeepAlivePeriod badoption.Duration `json:"keep_alive_period,omitempty"`
	// HTTP2 stream-level flow-control receive window size.
	//
	// Accepts memory size format, e.g. `"64 MB"`.
	StreamReceiveWindow *byteformats.MemoryBytes `json:"stream_receive_window,omitempty"`
	// HTTP2 connection-level flow-control receive window size.
	//
	// Accepts memory size format, e.g. `"64 MB"`.
	ConnectionReceiveWindow *byteformats.MemoryBytes `json:"connection_receive_window,omitempty"`
	// Maximum concurrent streams per connection.
	MaxConcurrentStreams int `json:"max_concurrent_streams,omitempty"`
}

type QUICOptions struct {
	HTTP2Options
	// Initial QUIC packet size.
	InitialPacketSize int `json:"initial_packet_size,omitempty"`
	// Disable QUIC path MTU discovery.
	DisablePathMTUDiscovery bool `json:"disable_path_mtu_discovery,omitempty"`
}

type HTTPClientOptions struct {
	Tag string `json:"tag,omitempty"`
	// HTTP engine to use.
	//
	// Values:
	//
	// * `go` (default)
	// * `apple`
	//
	// `apple` uses NSURLSession, only available on Apple platforms.
	//
	// @warning:
	//
	//     Experimental only: do not use in hot paths on iOS and tvOS.
	//
	// Supported fields:
	//
	// * `headers`
	// * `tls.server_name` (must match request host)
	// * `tls.insecure`
	// * `tls.min_version` / `tls.max_version`
	// * `tls.certificate` / `tls.certificate_path`
	// * `tls.certificate_sha256`
	// * `tls.certificate_public_key_sha256`
	// * Dial Fields
	//
	// Unsupported fields:
	//
	// * `version`
	// * `disable_version_fallback`
	// * HTTP2 Fields
	// * QUIC Fields
	// * `tls.engine`
	// * `tls.alpn`
	// * `tls.disable_sni`
	// * `tls.cipher_suites`
	// * `tls.curve_preferences`
	// * `tls.client_certificate` / `tls.client_certificate_path` / `tls.client_key` / `tls.client_key_path`
	// * `tls.fragment` / `tls.record_fragment`
	// * `tls.kernel_tx` / `tls.kernel_rx`
	// * `tls.ech`
	// * `tls.utls`
	// * `tls.reality`
	Engine string `json:"engine,omitempty" enum:"go,apple"`
	// HTTP version.
	//
	// Available values: `1`, `2`, `3`.
	//
	// `2` is used by default.
	//
	// When `3`, [HTTP2 Fields](#http2-fields) are replaced by [QUIC Fields](#quic-fields).
	Version int `json:"version,omitempty" enum:"0,1,2,3"`
	// Disable automatic fallback to lower HTTP version.
	DisableVersionFallback bool `json:"disable_version_fallback,omitempty"`
	// Custom HTTP headers.
	//
	// `Host` header is used as request host.
	Headers                 badoption.HTTPHeader `json:"headers,omitempty"`
	HTTP2Options            HTTP2Options         `json:"-"`
	HTTP3Options            QUICOptions          `json:"-"`
	DefaultOutbound         bool                 `json:"-"`
	DisableEmptyDirectCheck bool                 `json:"-"`
	ResolveOnDetour         bool                 `json:"-"`
	DirectResolver          bool                 `json:"-"`
	OutboundTLSOptionsContainer
	DialerOptions
}

type HTTPClient HTTPClientOptions

func (h HTTPClient) Options() HTTPClientOptions {
	options := HTTPClientOptions(h)
	options.Tag = ""
	return options
}

func (o HTTPClientOptions) IsEmpty() bool {
	if o.Tag != "" {
		return false
	}
	o.DefaultOutbound = false
	o.ResolveOnDetour = false
	o.DirectResolver = false
	return reflect.ValueOf(HTTPClientOptions(o)).IsZero()
}

func (o HTTPClientOptions) MarshalJSON() ([]byte, error) {
	if o.Tag != "" {
		return json.Marshal(o.Tag)
	}
	return badjson.MarshallObjects(HTTPClientOptions(o), httpClientVariant(HTTPClientOptions(o)))
}

func (o *HTTPClientOptions) UnmarshalJSONContext(ctx context.Context, content []byte) error {
	if len(content) > 0 && content[0] == '"' {
		*o = HTTPClientOptions{}
		return json.Unmarshal(content, &o.Tag)
	}
	var options HTTPClientOptions
	err := json.UnmarshalContext(ctx, content, &options)
	if err != nil {
		return err
	}
	err = unmarshalHTTPVersionOptions(ctx, content, &options, options.Version, &options.HTTP2Options, &options.HTTP3Options)
	if err != nil {
		return err
	}
	options.Tag = ""
	*o = HTTPClientOptions(options)
	return nil
}

func (h HTTPClient) MarshalJSON() ([]byte, error) {
	return badjson.MarshallObjects(HTTPClientOptions(h), httpClientVariant(HTTPClientOptions(h)))
}

func (h *HTTPClient) UnmarshalJSONContext(ctx context.Context, content []byte) error {
	err := json.UnmarshalContext(ctx, content, (*HTTPClientOptions)(h))
	if err != nil {
		return err
	}
	return unmarshalHTTPVersionOptions(ctx, content, (*HTTPClientOptions)(h), h.Version, &h.HTTP2Options, &h.HTTP3Options)
}

func unmarshalHTTPVersionOptions(ctx context.Context, content []byte, baseStruct any, version int, http2Options *HTTP2Options, http3Options *QUICOptions) error {
	switch version {
	case 1:
		return json.UnmarshalContextDisallowUnknownFields(ctx, content, baseStruct)
	case 0, 2:
		return badjson.UnmarshallExcludedContext(ctx, content, baseStruct, http2Options)
	case 3:
		return badjson.UnmarshallExcludedContext(ctx, content, baseStruct, http3Options)
	default:
		return E.New("unknown HTTP version: ", version)
	}
}

func unmarshalHTTPVersionsOptions(ctx context.Context, content []byte, baseStruct any, versions []int, http2Options *HTTP2Options, http3Options *QUICOptions) error {
	for _, version := range versions {
		if version < 1 || version > 3 {
			return E.New("unknown HTTP version: ", version)
		}
	}
	switch {
	case slices.Contains(versions, 3):
		err := badjson.UnmarshallExcludedContext(ctx, content, baseStruct, http3Options)
		if err != nil {
			return err
		}
		*http2Options = http3Options.HTTP2Options
		return nil
	case slices.Contains(versions, 2):
		return badjson.UnmarshallExcludedContext(ctx, content, baseStruct, http2Options)
	default:
		return json.UnmarshalContextDisallowUnknownFields(ctx, content, baseStruct)
	}
}

func httpVersionsVariant(versions []int, http2Options HTTP2Options, http3Options QUICOptions) any {
	switch {
	case slices.Contains(versions, 3):
		return http3Options
	case slices.Contains(versions, 2):
		return http2Options
	default:
		return nil
	}
}

func httpVersionVariant(version int, http2Options HTTP2Options, http3Options QUICOptions) any {
	switch version {
	case 1:
		return nil
	case 0, 2:
		return http2Options
	case 3:
		return http3Options
	default:
		return nil
	}
}

func httpClientVariant(options HTTPClientOptions) any {
	return httpVersionVariant(options.Version, options.HTTP2Options, options.HTTP3Options)
}

func describeHTTPClientObject(builder schema.Builder) (*schema.Node, error) {
	node := schema.StrictObject()
	err := builder.FlattenStruct(node, reflect.TypeFor[HTTPClient]())
	if err != nil {
		return nil, err
	}
	err = builder.FlattenStruct(node, reflect.TypeFor[QUICOptions]())
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (h HTTPClient) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return builder.Define("HTTPClient", func() (*schema.Node, error) {
		return describeHTTPClientObject(builder)
	})
}

func (o HTTPClientOptions) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return builder.Define("HTTPClientReference", func() (*schema.Node, error) {
		clientObject, err := describeHTTPClientObject(builder)
		if err != nil {
			return nil, err
		}
		clientObject.Properties.Remove("tag")
		return schema.AnyOf(schema.TagReferenceNode("http_client"), clientObject), nil
	})
}
