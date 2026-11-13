package option

import (
	"reflect"

	"github.com/sagernet/sing-box/schema"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badjson"
)

type _SnellInboundOptions struct {
	// The Snell protocol version, one of `5` `6`.
	//
	// Version `5` supports HTTP obfuscation (`obfs_mode`); version `6` replaces it
	// with traffic shaping (`mode`) and requires a `psk` of 12 to 255 bytes.
	//
	// @Required
	//
	// @Note
	//
	//     Since we intentionally do not support the QUIC proxy mode of Snell v5, the v5 wire protocol
	//     is effectively identical to v4, so no separate v4 server or v5 client is provided.
	Version int `json:"version" enum:"5,6"`
	AbstractSnellInboundOptions
	ObfsOptions SnellObfsServerOptions `json:"-"`
	V6Options   SnellV6Options         `json:"-"`
}

type AbstractSnellInboundOptions struct {
	ListenOptions

	// The pre-shared key.
	//
	// @Required
	PSK string `json:"psk"`

	// Snell users.
	//
	// When set, the server runs in multi-user mode: each entry has a `name` (optional, used in
	// logs) and a `userkey` (the user's key).
	Users []SnellUser `json:"users,omitempty"`
}

type SnellInboundOptions _SnellInboundOptions

func (o *SnellInboundOptions) UnmarshalJSON(content []byte) error {
	err := json.Unmarshal(content, (*_SnellInboundOptions)(o))
	if err != nil {
		return err
	}
	var versionOptions any
	switch o.Version {
	case 5:
		versionOptions = &o.ObfsOptions
	case 6:
		versionOptions = &o.V6Options
	case 0:
		return E.New("snell: missing version")
	default:
		return E.New("snell: unsupported version: ", o.Version)
	}
	return badjson.UnmarshallExcluded(content, (*_SnellInboundOptions)(o), versionOptions)
}

func (o SnellInboundOptions) MarshalJSON() ([]byte, error) {
	var versionOptions any
	switch o.Version {
	case 5:
		versionOptions = o.ObfsOptions
	case 6:
		versionOptions = o.V6Options
	case 0:
		return nil, E.New("snell: missing version")
	default:
		return nil, E.New("snell: unsupported version: ", o.Version)
	}
	return badjson.MarshallObjects(_SnellInboundOptions(o), versionOptions)
}

func (o SnellInboundOptions) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return schema.DiscriminatedUnion(builder, "version", true, []schema.UnionVariant{
		{Value: 5, StructType: reflect.TypeFor[SnellObfsServerOptions]()},
		{Value: 6, StructType: reflect.TypeFor[SnellV6Options]()},
	}, func(variant *schema.Node) error {
		return builder.FlattenStruct(variant, reflect.TypeFor[AbstractSnellInboundOptions]())
	})
}

type _SnellOutboundOptions struct {
	// The Snell protocol version, one of `4` `6`.
	//
	// Version `4` supports HTTP obfuscation (`obfs_mode` / `obfs_host`); version `6`
	// replaces it with traffic shaping (`mode`) and requires a `psk` of 12 to 255
	// bytes.
	//
	// @Required
	//
	// @Note
	//
	//     Since we intentionally do not support the QUIC proxy mode of Snell v5, the v5 wire protocol
	//     is effectively identical to v4, so no separate v4 server or v5 client is provided.
	Version int `json:"version" enum:"4,6"`
	AbstractSnellOutboundOptions
	ObfsOptions SnellObfsClientOptions `json:"-"`
	V6Options   SnellV6Options         `json:"-"`
}

type AbstractSnellOutboundOptions struct {
	DialerOptions
	ServerOptions

	// The pre-shared key.
	//
	// @Required
	PSK string `json:"psk"`

	// The user key, used to authenticate against a multi-user server.
	UserKey string `json:"userkey,omitempty"`

	// Enable connection reuse (the Snell v2 `CONNECT` command).
	Reuse bool `json:"reuse,omitempty"`

	// Enabled network
	//
	// One of `tcp` `udp`.
	//
	// Both is enabled by default.
	Network NetworkList `json:"network,omitempty"`
}

type SnellOutboundOptions _SnellOutboundOptions

func (o *SnellOutboundOptions) UnmarshalJSON(content []byte) error {
	err := json.Unmarshal(content, (*_SnellOutboundOptions)(o))
	if err != nil {
		return err
	}
	var versionOptions any
	switch o.Version {
	case 4:
		versionOptions = &o.ObfsOptions
	case 6:
		versionOptions = &o.V6Options
	case 0:
		return E.New("snell: missing version")
	default:
		return E.New("snell: unsupported version: ", o.Version)
	}
	return badjson.UnmarshallExcluded(content, (*_SnellOutboundOptions)(o), versionOptions)
}

func (o SnellOutboundOptions) MarshalJSON() ([]byte, error) {
	var versionOptions any
	switch o.Version {
	case 4:
		versionOptions = o.ObfsOptions
	case 6:
		versionOptions = o.V6Options
	case 0:
		return nil, E.New("snell: missing version")
	default:
		return nil, E.New("snell: unsupported version: ", o.Version)
	}
	return badjson.MarshallObjects(_SnellOutboundOptions(o), versionOptions)
}

func (o SnellOutboundOptions) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return schema.DiscriminatedUnion(builder, "version", true, []schema.UnionVariant{
		{Value: 4, StructType: reflect.TypeFor[SnellObfsClientOptions]()},
		{Value: 6, StructType: reflect.TypeFor[SnellV6Options]()},
	}, func(variant *schema.Node) error {
		return builder.FlattenStruct(variant, reflect.TypeFor[AbstractSnellOutboundOptions]())
	})
}

type SnellObfsServerOptions struct {
	// HTTP obfuscation mode, one of `none` `http`.
	//
	// @Default: none
	//
	// @Indicator: Version 5 only
	ObfsMode string `json:"obfs_mode,omitempty" enum:"none,http,tls"`
}

type SnellUser struct {
	Name    string `json:"name,omitempty"`
	UserKey string `json:"userkey"`
}

type SnellObfsClientOptions struct {
	// HTTP obfuscation mode, one of `none` `http`.
	//
	// @Default: none
	//
	// @Indicator: Version 4 only
	ObfsMode string `json:"obfs_mode,omitempty" enum:"none,http,tls"`

	// The HTTP `Host` header sent when `obfs_mode` is `http`.
	//
	// @Default: bing.com
	//
	// @Indicator: Version 4 only
	ObfsHost string `json:"obfs_host,omitempty"`
}

type SnellV6Options struct {
	// Traffic shaping mode, one of `default` `unshaped` `unsafe-raw`.
	//
	// @Default: default
	//
	// @Indicator: Version 6 only
	Mode string `json:"mode,omitempty" enum:"default,unshaped,unsafe-raw"`
}
