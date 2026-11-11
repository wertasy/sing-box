package option

import (
	"context"
	"reflect"

	"github.com/sagernet/sing-box/schema"
	"github.com/sagernet/sing/common/auth"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badjson"
	"github.com/sagernet/sing/common/json/badoption"
)

type SocksInboundOptions struct {
	ListenOptions
	// SOCKS users.
	//
	// No authentication required if empty.
	Users          []auth.User           `json:"users,omitempty"`
	DomainResolver *DomainResolveOptions `json:"domain_resolver,omitempty"`
}

type HTTPMixedInboundOptions struct {
	ListenOptions
	// SOCKS and HTTP users.
	//
	// No authentication required if empty.
	Users          []auth.User           `json:"users,omitempty"`
	DomainResolver *DomainResolveOptions `json:"domain_resolver,omitempty"`
	// @quote:
	//
	//     Only supported on Linux, Android, Windows, and macOS.
	//
	// @warning:
	//
	//     To work on Android and Apple platforms without privileges, use tun.platform.http_proxy instead.
	//
	// Automatically set system proxy configuration when start and clean up when stop.
	SetSystemProxy bool `json:"set_system_proxy,omitempty"`
	InboundTLSOptionsContainer
}

type _HTTPInboundOptions struct {
	ListenOptions
	// HTTP users.
	//
	// No authentication required if empty.
	Users          []auth.User           `json:"users,omitempty"`
	DomainResolver *DomainResolveOptions `json:"domain_resolver,omitempty"`
	// @quote:
	//
	//     Only supported on Linux, Android, Windows, and macOS.
	//
	// @warning:
	//
	//     To work on Android and Apple platforms without privileges, use tun.platform.http_proxy instead.
	//
	// Automatically set system proxy configuration when start and clean up when stop.
	SetSystemProxy bool `json:"set_system_proxy,omitempty"`
	// @since 1.15.0
	//
	// List of HTTP versions to serve.
	//
	// Available values: `1`, `2`, `3`.
	//
	// `1` and `2` are used by default.
	//
	// TLS is required for `3`.
	Version badoption.Listable[int] `json:"version,omitempty" enum:"1,2,3"`
	InboundTLSOptionsContainer
	HTTP2Options HTTP2Options `json:"-"`
	HTTP3Options QUICOptions  `json:"-"`
}

type HTTPInboundOptions _HTTPInboundOptions

func (o HTTPInboundOptions) Versions() []int {
	if len(o.Version) > 0 {
		return o.Version
	}
	return []int{1, 2}
}

func (o HTTPInboundOptions) MarshalJSON() ([]byte, error) {
	return badjson.MarshallObjects(_HTTPInboundOptions(o), httpVersionsVariant(o.Versions(), o.HTTP2Options, o.HTTP3Options))
}

func (o *HTTPInboundOptions) UnmarshalJSONContext(ctx context.Context, content []byte) error {
	err := json.UnmarshalContext(ctx, content, (*_HTTPInboundOptions)(o))
	if err != nil {
		return err
	}
	return unmarshalHTTPVersionsOptions(ctx, content, (*_HTTPInboundOptions)(o), o.Versions(), &o.HTTP2Options, &o.HTTP3Options)
}

func (o HTTPInboundOptions) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	node := schema.StrictObject()
	err := builder.FlattenStruct(node, reflect.TypeFor[HTTPInboundOptions]())
	if err != nil {
		return nil, err
	}
	err = builder.FlattenStruct(node, reflect.TypeFor[QUICOptions]())
	if err != nil {
		return nil, err
	}
	return node, nil
}

type SOCKSOutboundOptions struct {
	DialerOptions
	ServerOptions
	// The SOCKS version, one of `4` `4a` `5`.
	//
	// SOCKS5 used by default.
	Version string `json:"version,omitempty" enum:"4,4a,5"`
	// SOCKS username.
	Username string `json:"username,omitempty"`
	// SOCKS5 password.
	Password string `json:"password,omitempty"`
	// Enabled network
	//
	// One of `tcp` `udp`.
	//
	// Both is enabled by default.
	Network NetworkList `json:"network,omitempty"`
	// UDP over TCP protocol settings.
	//
	// See [UDP Over TCP](/configuration/shared/udp-over-tcp/) for details.
	UDPOverTCP *UDPOverTCPOptions `json:"udp_over_tcp,omitempty"`
}

type _HTTPOutboundOptions struct {
	DialerOptions
	ServerOptions
	// Basic authorization username.
	Username string `json:"username,omitempty"`
	// Basic authorization password.
	Password string `json:"password,omitempty"`
	OutboundTLSOptionsContainer
	// Path of HTTP request.
	Path string `json:"path,omitempty"`
	// Extra headers of HTTP request.
	Headers badoption.HTTPHeader `json:"headers,omitempty"`
	// @since 1.15.0
	//
	// HTTP version.
	//
	// Available values: `1`, `2`, `3`.
	//
	// `2` is used by default, or `1` if `path` or the `Host` header is set.
	//
	// `path` and the `Host` header are only available for `1`.
	//
	// When `3`, [HTTP2 Fields](#http2-fields) are replaced by [QUIC Fields](#quic-fields).
	Version int `json:"version,omitempty" enum:"0,1,2,3"`
	// @since 1.15.0
	//
	// Disable automatic fallback to lower HTTP version.
	DisableVersionFallback bool         `json:"disable_version_fallback,omitempty"`
	HTTP2Options           HTTP2Options `json:"-"`
	HTTP3Options           QUICOptions  `json:"-"`
}

type HTTPOutboundOptions _HTTPOutboundOptions

func (o HTTPOutboundOptions) MarshalJSON() ([]byte, error) {
	return badjson.MarshallObjects(_HTTPOutboundOptions(o), httpVersionVariant(o.Version, o.HTTP2Options, o.HTTP3Options))
}

func (o *HTTPOutboundOptions) UnmarshalJSONContext(ctx context.Context, content []byte) error {
	err := json.UnmarshalContext(ctx, content, (*_HTTPOutboundOptions)(o))
	if err != nil {
		return err
	}
	return unmarshalHTTPVersionOptions(ctx, content, (*_HTTPOutboundOptions)(o), o.Version, &o.HTTP2Options, &o.HTTP3Options)
}

func (o HTTPOutboundOptions) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	node := schema.StrictObject()
	err := builder.FlattenStruct(node, reflect.TypeFor[HTTPOutboundOptions]())
	if err != nil {
		return nil, err
	}
	err = builder.FlattenStruct(node, reflect.TypeFor[QUICOptions]())
	if err != nil {
		return nil, err
	}
	return node, nil
}
