package option

import (
	"reflect"

	"github.com/sagernet/sing-box/schema"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badjson"
)

const (
	USBIPProviderDefault = "default"
	USBIPProviderDynamic = "dynamic"
)

type _USBIPServerServiceOptions struct {
	ListenOptions

	// The device source provider.
	//
	// - `default`: Exports the local devices matched by `devices`. The default value.
	// - `dynamic`: Devices are provided at runtime through a [sing-box API](/configuration/service/api/)
	//   client instead of from configuration, on supported platforms: the sing-box graphical clients on
	//   [macOS](/clients/apple/) and [Android](/clients/android/), and Chromium-based browsers with
	//   [sing-box Dashboard](https://github.com/SagerNet/sing-box-dashboard).
	//
	// @Quote:
	//
	//     The `default` provider is only supported when running directly via the CLI on Linux, Windows,
	//     and macOS, and requires elevated privileges.
	Provider string `json:"provider,omitempty" enum:"default,dynamic"`
	Options  any    `json:"-"`
}

type USBIPServerServiceOptions _USBIPServerServiceOptions

func (o USBIPServerServiceOptions) MarshalJSON() ([]byte, error) {
	if o.Provider == USBIPProviderDefault {
		o.Provider = ""
	}
	return badjson.MarshallObjects(_USBIPServerServiceOptions(o), o.Options)
}

func (o *USBIPServerServiceOptions) UnmarshalJSON(content []byte) error {
	err := json.Unmarshal(content, (*_USBIPServerServiceOptions)(o))
	if err != nil {
		return err
	}
	var options any
	switch o.Provider {
	case "", USBIPProviderDefault:
		o.Provider = USBIPProviderDefault
		options = new(USBIPDefaultProviderOptions)
	case USBIPProviderDynamic:
		options = new(USBIPDynamicProviderOptions)
	default:
		return E.New("unknown usbip provider type: ", o.Provider)
	}
	err = badjson.UnmarshallExcluded(content, (*_USBIPServerServiceOptions)(o), options)
	if err != nil {
		return err
	}
	o.Options = options
	return nil
}

func (o USBIPServerServiceOptions) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	return schema.DiscriminatedUnion(builder, "provider", false, []schema.UnionVariant{
		{Value: USBIPProviderDefault, StructType: reflect.TypeFor[USBIPDefaultProviderOptions](), TypeOptional: true},
		{Value: USBIPProviderDynamic, StructType: reflect.TypeFor[USBIPDynamicProviderOptions]()},
	}, func(variant *schema.Node) error {
		return builder.FlattenStruct(variant, reflect.TypeFor[ListenOptions]())
	})
}

type USBIPClientServiceOptions struct {
	DialerOptions

	// @Override server:
	//    @Required
	//
	//    The remote `usbip-server` address.
	//
	//    @Example: ""
	//
	// @Override server_port:
	//    The remote `usbip-server` port. Defaults to `3240`.
	//
	//    @Example: 0
	ServerOptions

	// List of device matches selecting which remote devices to import. If empty, all exported devices
	// are imported.
	//
	// Object format:
	//
	// ```json
	// {
	//   "bus_id": "",
	//   "vendor_id": 0,
	//   "product_id": 0,
	//   "serial": ""
	// }
	// ```
	//
	// Object fields:
	//
	// - `bus_id`: USB bus ID, e.g. `1-2`.
	// - `vendor_id`: USB vendor ID, as a number.
	// - `product_id`: USB product ID, as a number.
	// - `serial`: Device serial number.
	//
	// Within one object, all specified fields must match; multiple objects are combined as a union. At
	// least one field is required.
	Devices []USBIPDeviceMatch `json:"devices,omitempty"`
}

type USBIPDeviceMatch struct {
	BusID     string `json:"bus_id,omitempty"`
	VendorID  uint16 `json:"vendor_id,omitempty"`
	ProductID uint16 `json:"product_id,omitempty"`
	Serial    string `json:"serial,omitempty"`
}

type USBIPDefaultProviderOptions struct {
	// List of device matches selecting which local USB devices to export.
	//
	// Object format:
	//
	// ```json
	// {
	//   "bus_id": "",
	//   "vendor_id": 0,
	//   "product_id": 0,
	//   "serial": ""
	// }
	// ```
	//
	// Object fields:
	//
	// - `bus_id`: USB bus ID, e.g. `1-2`.
	// - `vendor_id`: USB vendor ID, as a number.
	// - `product_id`: USB product ID, as a number.
	// - `serial`: Device serial number.
	//
	// Within one object, all specified fields must match; multiple objects are combined as a union. At
	// least one field is required.
	//
	// @Required: with the `default` provider.
	Devices []USBIPDeviceMatch `json:"devices,omitempty"`
}

type USBIPDynamicProviderOptions struct{}
