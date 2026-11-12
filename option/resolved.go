package option

import (
	"context"
	"net/netip"
	"reflect"

	"github.com/sagernet/sing-box/schema"
	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badoption"
)

type _ResolvedServiceOptions struct {
	// @Override listen:
	//    @Required
	//
	//    Listen address.
	//
	//    `127.0.0.53` will be used by default.
	//
	// @Override listen_port:
	//    @Required
	//
	//    Listen port.
	//
	//    `53` will be used by default.
	ListenOptions
}

type ResolvedServiceOptions _ResolvedServiceOptions

func (r ResolvedServiceOptions) MarshalJSONContext(ctx context.Context) ([]byte, error) {
	if r.Listen != nil && netip.Addr(*r.Listen) == netip.AddrFrom4([4]byte{127, 0, 0, 53}) {
		r.Listen = nil
	}
	if r.ListenPort == 53 {
		r.ListenPort = 0
	}
	return json.MarshalContext(ctx, (*_ResolvedServiceOptions)(&r))
}

func (r *ResolvedServiceOptions) UnmarshalJSONContext(ctx context.Context, bytes []byte) error {
	err := json.UnmarshalContextDisallowUnknownFields(ctx, bytes, (*_ResolvedServiceOptions)(r))
	if err != nil {
		return err
	}
	if r.Listen == nil {
		r.Listen = (*badoption.Addr)(common.Ptr(netip.AddrFrom4([4]byte{127, 0, 0, 53})))
	}
	if r.ListenPort == 0 {
		r.ListenPort = 53
	}
	return nil
}

func (r ResolvedServiceOptions) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	node := schema.StrictObject()
	err := builder.FlattenStruct(node, reflect.TypeFor[ResolvedServiceOptions]())
	if err != nil {
		return nil, err
	}
	return node, nil
}

type ResolvedDNSServerOptions struct {
	// The tag of the [Resolved Service](/configuration/service/resolved).
	//
	// @Required
	Service string `json:"service"`
	// Indicates whether the default DNS resolvers should be accepted for fallback queries in addition to matching domains.
	//
	// Specifically, default DNS resolvers are DNS servers that have `SetLinkDefaultRoute` or `SetLinkDomains ~.` set.
	//
	// If not enabled, `NXDOMAIN` will be returned for requests that do not match search or match domains.
	AcceptDefaultResolvers bool `json:"accept_default_resolvers,omitempty"`
	// Timeout                badoption.Duration `json:"timeout,omitempty"`
	// Attempts               int                `json:"attempts,omitempty"`
	// Rotate                 bool               `json:"rotate,omitempty"`
}
