package option

import (
	"context"
	"net/netip"
	"reflect"

	"github.com/sagernet/sing-box/schema"
	"github.com/sagernet/sing/common/auth"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badjson"
	"github.com/sagernet/sing/common/json/badoption"
)

type MASQUEEndpointOptions struct {
	// Use system interface.
	//
	// Requires privilege and cannot conflict with existing system interfaces.
	//
	// If disabled, sing-box uses the internal network stack.
	System bool `json:"system,omitempty"`

	// Custom interface name for system interface.
	//
	// An automatically generated `masque` interface name will be used by default.
	Name string `json:"name,omitempty"`

	// Tunnel MTU.
	//
	// @Default: 1280
	//
	// @Example: 1280
	MTU uint32 `json:"mtu,omitempty"`

	// @Omit
	UDPMapping UDPNATBehavior `json:"udp_mapping,omitempty"`

	// @Omit
	UDPFiltering UDPNATBehavior `json:"udp_filtering,omitempty"`

	// @Omit
	UDPNATMax uint32 `json:"udp_nat_max,omitempty"`
}

type _MASQUEClientEndpointOptions struct {
	DialerOptions

	// @Override server:
	//    @Example: "127.0.0.1"
	//
	// @Override server_port:
	//    @Example: 443
	ServerOptions
	MASQUEEndpointOptions

	// Basic authorization username.
	Username string `json:"username,omitempty"`

	// Basic authorization password.
	Password string `json:"password,omitempty"`

	// @Override:
	//    TLS configuration, see [TLS](/configuration/shared/tls/#outbound).
	//
	//    Required for HTTP/3.
	OutboundTLSOptionsContainer

	// URI template path of the IP proxying resource, may contain the `target` and `ipproto` variables.
	//
	// `/.well-known/masque/ip/{target}/{ipproto}/` will be used by default.
	Path string `json:"path,omitempty"`

	// Extra headers of HTTP request.
	Headers badoption.HTTPHeader `json:"headers,omitempty"`

	// HTTP version.
	//
	// Available values: `1`, `2`, `3`.
	//
	// When `2`, [QUIC Fields](#quic-fields) are replaced by [HTTP2 Fields](#http2-fields).
	//
	// @Default: 3
	Version int `json:"version,omitempty" enum:"0,1,2,3"`

	// Disable automatic fallback to lower HTTP version.
	DisableVersionFallback bool `json:"disable_version_fallback,omitempty"`

	// List of IP prefixes to advertise to the server.
	//
	// The server will route traffic for these prefixes into this endpoint, where it is handled as inbound traffic.
	AdvertiseRoutes []netip.Prefix `json:"advertise_routes,omitempty"`

	// @Omit
	UDPTimeout badoption.Duration `json:"udp_timeout,omitempty"`

	// Allow the endpoint to be disconnected when necessary.
	OnDemand     bool         `json:"on_demand,omitempty"`
	HTTP2Options HTTP2Options `json:"-"`
	HTTP3Options QUICOptions  `json:"-"`
}

// @Tag masque-client
type MASQUEClientEndpointOptions _MASQUEClientEndpointOptions

func (o MASQUEClientEndpointOptions) ResolvedVersion() int {
	if o.Version == 0 {
		return 3
	}
	return o.Version
}

func (o MASQUEClientEndpointOptions) MarshalJSON() ([]byte, error) {
	return badjson.MarshallObjects(_MASQUEClientEndpointOptions(o), httpVersionVariant(o.ResolvedVersion(), o.HTTP2Options, o.HTTP3Options))
}

func (o *MASQUEClientEndpointOptions) UnmarshalJSONContext(ctx context.Context, content []byte) error {
	err := json.UnmarshalContext(ctx, content, (*_MASQUEClientEndpointOptions)(o))
	if err != nil {
		return err
	}
	return unmarshalHTTPVersionOptions(ctx, content, (*_MASQUEClientEndpointOptions)(o), o.ResolvedVersion(), &o.HTTP2Options, &o.HTTP3Options)
}

func (o MASQUEClientEndpointOptions) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	node := schema.StrictObject()
	err := builder.FlattenStruct(node, reflect.TypeFor[MASQUEClientEndpointOptions]())
	if err != nil {
		return nil, err
	}
	err = builder.FlattenStruct(node, reflect.TypeFor[QUICOptions]())
	if err != nil {
		return nil, err
	}
	return node, nil
}

type _MASQUEServerEndpointOptions struct {
	ListenOptions
	MASQUEEndpointOptions

	// HTTP users, verified by the `Authorization` header.
	//
	// No authentication required if empty.
	//
	// @Example:
	//    [
	//      {
	//        "username": "",
	//        "password": ""
	//      }
	//    ]
	Users []auth.User `json:"users,omitempty"`

	// List of HTTP versions to serve.
	//
	// Available values: `1`, `2`, `3`.
	//
	// All versions are used by default.
	//
	// TLS is required for `3`.
	Version badoption.Listable[int] `json:"version,omitempty" enum:"1,2,3"`

	// @Override:
	//    TLS configuration, see [TLS](/configuration/shared/tls/#inbound).
	//
	//    IP proxying must be operated over TLS or QUIC. Leave it disabled only when the server is placed behind an HTTP intermediary that terminates TLS.
	InboundTLSOptionsContainer

	// URI template path of the IP proxying resource, may contain the `target` and `ipproto` variables.
	//
	// `/.well-known/masque/ip/{target}/{ipproto}/` will be used by default.
	Path string `json:"path,omitempty"`

	// List of IP prefixes of the tunnel network, at most one for each IP version.
	//
	// The address of the prefix is used by the server itself, other addresses in the prefix are assigned to clients.
	//
	// @Required
	Address badoption.Listable[netip.Prefix] `json:"address"`

	// List of IP prefixes to advertise to clients, in addition to the tunnel network.
	//
	// Traffic from clients to other destinations is rejected.
	//
	// All addresses are advertised by default.
	AdvertiseRoutes []netip.Prefix `json:"advertise_routes,omitempty"`
	HTTP2Options    HTTP2Options   `json:"-"`
	HTTP3Options    QUICOptions    `json:"-"`
}

// @Tag masque-server
type MASQUEServerEndpointOptions _MASQUEServerEndpointOptions

func (o MASQUEServerEndpointOptions) Versions() []int {
	if len(o.Version) > 0 {
		return o.Version
	}
	return []int{1, 2, 3}
}

func (o MASQUEServerEndpointOptions) MarshalJSON() ([]byte, error) {
	return badjson.MarshallObjects(_MASQUEServerEndpointOptions(o), httpVersionsVariant(o.Versions(), o.HTTP2Options, o.HTTP3Options))
}

func (o *MASQUEServerEndpointOptions) UnmarshalJSONContext(ctx context.Context, content []byte) error {
	err := json.UnmarshalContext(ctx, content, (*_MASQUEServerEndpointOptions)(o))
	if err != nil {
		return err
	}
	return unmarshalHTTPVersionsOptions(ctx, content, (*_MASQUEServerEndpointOptions)(o), o.Versions(), &o.HTTP2Options, &o.HTTP3Options)
}

func (o MASQUEServerEndpointOptions) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	node := schema.StrictObject()
	err := builder.FlattenStruct(node, reflect.TypeFor[MASQUEServerEndpointOptions]())
	if err != nil {
		return nil, err
	}
	err = builder.FlattenStruct(node, reflect.TypeFor[QUICOptions]())
	if err != nil {
		return nil, err
	}
	return node, nil
}
