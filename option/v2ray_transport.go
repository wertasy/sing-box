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

type _V2RayTransportOptions struct {
	Type               string                  `json:"type" enum:"http,ws,quic,grpc,httpupgrade"`
	HTTPOptions        V2RayHTTPOptions        `json:"-"`
	WebsocketOptions   V2RayWebsocketOptions   `json:"-"`
	QUICOptions        V2RayQUICOptions        `json:"-"`
	GRPCOptions        V2RayGRPCOptions        `json:"-"`
	HTTPUpgradeOptions V2RayHTTPUpgradeOptions `json:"-"`
}

type V2RayTransportOptions _V2RayTransportOptions

func (o V2RayTransportOptions) MarshalJSON() ([]byte, error) {
	var v any
	switch o.Type {
	case C.V2RayTransportTypeHTTP:
		v = o.HTTPOptions
	case C.V2RayTransportTypeWebsocket:
		v = o.WebsocketOptions
	case C.V2RayTransportTypeQUIC:
		v = o.QUICOptions
	case C.V2RayTransportTypeGRPC:
		v = o.GRPCOptions
	case C.V2RayTransportTypeHTTPUpgrade:
		v = o.HTTPUpgradeOptions
	case "":
		return nil, E.New("missing transport type")
	default:
		return nil, E.New("unknown transport type: " + o.Type)
	}
	return badjson.MarshallObjects(_V2RayTransportOptions(o), v)
}

func (o *V2RayTransportOptions) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, (*_V2RayTransportOptions)(o))
	if err != nil {
		return err
	}
	var v any
	switch o.Type {
	case C.V2RayTransportTypeHTTP:
		v = &o.HTTPOptions
	case C.V2RayTransportTypeWebsocket:
		v = &o.WebsocketOptions
	case C.V2RayTransportTypeQUIC:
		v = &o.QUICOptions
	case C.V2RayTransportTypeGRPC:
		v = &o.GRPCOptions
	case C.V2RayTransportTypeHTTPUpgrade:
		v = &o.HTTPUpgradeOptions
	default:
		return E.New("unknown transport type: " + o.Type)
	}
	err = badjson.UnmarshallExcluded(bytes, (*_V2RayTransportOptions)(o), v)
	if err != nil {
		return err
	}
	return nil
}

func (o V2RayTransportOptions) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return builder.Define("V2RayTransport", func() (*schema.Node, error) {
		return schema.DiscriminatedUnion(builder, "type", true, []schema.UnionVariant{
			{Value: C.V2RayTransportTypeHTTP, StructType: reflect.TypeFor[V2RayHTTPOptions]()},
			{Value: C.V2RayTransportTypeWebsocket, StructType: reflect.TypeFor[V2RayWebsocketOptions]()},
			{Value: C.V2RayTransportTypeQUIC, StructType: reflect.TypeFor[V2RayQUICOptions]()},
			{Value: C.V2RayTransportTypeGRPC, StructType: reflect.TypeFor[V2RayGRPCOptions]()},
			{Value: C.V2RayTransportTypeHTTPUpgrade, StructType: reflect.TypeFor[V2RayHTTPUpgradeOptions]()},
		}, nil)
	})
}

type V2RayHTTPOptions struct {
	// List of host domain.
	//
	// The client will choose randomly and the server will verify if not empty.
	Host badoption.Listable[string] `json:"host,omitempty"`
	// @warning
	//
	//     V2Ray's documentation says that the path between the server and the client must be consistent,
	//     but the actual code allows the client to add any suffix to the path.
	//     sing-box uses the same behavior as V2Ray, but note that the behavior does not exist in `WebSocket` and `HTTPUpgrade` transport.
	//
	// Path of HTTP request.
	//
	// The server will verify.
	Path string `json:"path,omitempty"`
	// Method of HTTP request.
	//
	// The server will verify if not empty.
	Method string `json:"method,omitempty"`
	// Extra headers of HTTP request.
	//
	// The server will write in response if not empty.
	Headers badoption.HTTPHeader `json:"headers,omitempty"`
	// In HTTP2 server:
	//
	// Specifies the time until idle clients should be closed with a GOAWAY frame. PING frames are not considered as activity.
	//
	// In HTTP2 client:
	//
	// Specifies the period of time after which a health check will be performed using a ping frame if no frames have been
	// received on the connection.Please note that a ping response is considered a received frame, so if there is no other
	// traffic on the connection, the health check will be executed every interval. If the value is zero, no health check will
	// be performed.
	//
	// Zero is used by default.
	IdleTimeout badoption.Duration `json:"idle_timeout,omitempty"`
	// In HTTP2 client:
	//
	// Specifies the timeout duration after sending a PING frame, within which a response must be received.
	// If a response to the PING frame is not received within the specified timeout duration, the connection will be closed.
	// The default timeout duration is 15 seconds.
	PingTimeout badoption.Duration `json:"ping_timeout,omitempty"`
}

type V2RayWebsocketOptions struct {
	// Path of HTTP request.
	//
	// The server will verify.
	Path string `json:"path,omitempty"`
	// Extra headers of HTTP request.
	//
	// The server will write in response if not empty.
	Headers badoption.HTTPHeader `json:"headers,omitempty"`
	// Allowed payload size is in the request. Enabled if not zero.
	MaxEarlyData uint32 `json:"max_early_data,omitempty"`
	// Early data is sent in path instead of header by default.
	//
	// To be compatible with Xray-core, set this to `Sec-WebSocket-Protocol`.
	//
	// It needs to be consistent with the server.
	EarlyDataHeaderName string `json:"early_data_header_name,omitempty"`
}

type V2RayQUICOptions struct{}

type V2RayGRPCOptions struct {
	// Service name of gRPC.
	ServiceName string `json:"service_name,omitempty"`
	// In standard gRPC server/client:
	//
	// If the transport doesn't see any activity after a duration of this time,
	// it pings the client to check if the connection is still active.
	//
	// In default gRPC server/client:
	//
	// It has the same behavior as the corresponding setting in HTTP transport.
	IdleTimeout badoption.Duration `json:"idle_timeout,omitempty"`
	// In standard gRPC server/client:
	//
	// The timeout that after performing a keepalive check, the client will wait for activity.
	// If no activity is detected, the connection will be closed.
	//
	// In default gRPC server/client:
	//
	// It has the same behavior as the corresponding setting in HTTP transport.
	PingTimeout badoption.Duration `json:"ping_timeout,omitempty"`
	// In standard gRPC client:
	//
	// If enabled, the client transport sends keepalive pings even with no active connections.
	// If disabled, when there are no active connections, `idle_timeout` and `ping_timeout` will be ignored and no keepalive
	// pings will be sent.
	//
	// Disabled by default.
	PermitWithoutStream bool `json:"permit_without_stream,omitempty"`
	ForceLite           bool `json:"-"` // for test
}

type V2RayHTTPUpgradeOptions struct {
	// Host domain.
	//
	// The server will verify if not empty.
	Host string `json:"host,omitempty"`
	// Path of HTTP request.
	//
	// The server will verify.
	Path string `json:"path,omitempty"`
	// Extra headers of HTTP request.
	//
	// The server will write in response if not empty.
	Headers badoption.HTTPHeader `json:"headers,omitempty"`
}
