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
	//
	// @Example:
	//    [
	//      {
	//        "username": "admin",
	//        "password": "admin"
	//      }
	//    ]
	Users []auth.User `json:"users,omitempty"`

	// @Omit
	DomainResolver *DomainResolveOptions `json:"domain_resolver,omitempty"`
}

// @Tag mixed-in
type HTTPMixedInboundOptions struct {
	ListenOptions

	// SOCKS and HTTP users.
	//
	// No authentication required if empty.
	//
	// @Example:
	//    [
	//      {
	//        "username": "admin",
	//        "password": "admin"
	//      }
	//    ]
	Users []auth.User `json:"users,omitempty"`

	// @Omit
	DomainResolver *DomainResolveOptions `json:"domain_resolver,omitempty"`

	// Automatically set system proxy configuration when start and clean up when stop.
	//
	// @Quote:
	//
	//     Only supported on Linux, Android, Windows, and macOS.
	//
	// @Warning:
	//
	//     To work on Android and Apple platforms without privileges, use tun.platform.http_proxy instead.
	SetSystemProxy bool `json:"set_system_proxy,omitempty"`
	InboundTLSOptionsContainer
}

type _HTTPInboundOptions struct {
	ListenOptions

	// HTTP users.
	//
	// No authentication required if empty.
	//
	// @Example:
	//    [
	//      {
	//        "username": "admin",
	//        "password": "admin"
	//      }
	//    ]
	Users []auth.User `json:"users,omitempty"`

	// @Omit
	DomainResolver *DomainResolveOptions `json:"domain_resolver,omitempty"`

	// Automatically set system proxy configuration when start and clean up when stop.
	//
	// @Quote:
	//
	//     Only supported on Linux, Android, Windows, and macOS.
	//
	// @Warning:
	//
	//     To work on Android and Apple platforms without privileges, use tun.platform.http_proxy instead.
	SetSystemProxy bool `json:"set_system_proxy,omitempty"`

	// List of HTTP versions to serve.
	//
	// Available values: `1`, `2`, `3`.
	//
	// `1` and `2` are used by default.
	//
	// TLS is required for `3`.
	//
	// @Since 1.15.0
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
	//
	// @Example: "5"
	Version string `json:"version,omitempty" enum:"4,4a,5"`

	// SOCKS username.
	//
	// @Example: "sekai"
	Username string `json:"username,omitempty"`

	// SOCKS5 password.
	//
	// @Example: "admin"
	Password string `json:"password,omitempty"`

	// Enabled network
	//
	// One of `tcp` `udp`.
	//
	// Both is enabled by default.
	//
	// @Example: "udp"
	Network NetworkList `json:"network,omitempty"`

	// UDP over TCP protocol settings.
	//
	// See [UDP Over TCP](/configuration/shared/udp-over-tcp/) for details.
	//
	// @Example: false | {}
	UDPOverTCP *UDPOverTCPOptions `json:"udp_over_tcp,omitempty"`
}

type _HTTPOutboundOptions struct {
	DialerOptions

	// @Override server:
	//    @Example: "127.0.0.1"
	//
	// @Override server_port:
	//    @Example: 1080
	// @Override server:
	//    @Example: "127.0.0.1"
	//
	// @Override server_port:
	//    @Example: 1080
	ServerOptions

	// Basic authorization username.
	//
	// @Example: "sekai"
	Username string `json:"username,omitempty"`

	// Basic authorization password.
	//
	// @Example: "admin"
	Password string `json:"password,omitempty"`
	OutboundTLSOptionsContainer

	// Path of HTTP request.
	Path string `json:"path,omitempty"`

	// Extra headers of HTTP request.
	Headers badoption.HTTPHeader `json:"headers,omitempty"`

	// HTTP version.
	//
	// Available values: `1`, `2`, `3`.
	//
	// `2` will be used by default, or `1` if `path` or the `Host` header is set.
	//
	// `path` and the `Host` header are only available for `1`.
	//
	// When `3`, [HTTP2 Fields](#http2-fields) are replaced by [QUIC Fields](#quic-fields).
	//
	// @Since 1.15.0
	Version int `json:"version,omitempty" enum:"0,1,2,3"`

	// Disable automatic fallback to lower HTTP version.
	//
	// @Since 1.15.0
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
