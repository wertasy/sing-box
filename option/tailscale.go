package option

import (
	"net/netip"
	"net/url"
	"reflect"

	"github.com/sagernet/sing-box/schema"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badjson"
	"github.com/sagernet/sing/common/json/badoption"
	M "github.com/sagernet/sing/common/metadata"
)

type TailscaleEndpointOptions struct {
	DialerOptions
	// The directory where the Tailscale state is stored.
	//
	// `tailscale` is used by default.
	//
	// Example: `$HOME/.tailscale`
	StateDirectory string `json:"state_directory,omitempty"`
	// @note
	//
	//     Auth key is not required. By default, sing-box will log the login URL.
	//
	// The auth key to create the node. If the node is already created (from state previously stored), then this field is not
	// used.
	AuthKey string `json:"auth_key,omitempty"`
	// The coordination server URL.
	//
	// `https://controlplane.tailscale.com` is used by default.
	ControlURL string `json:"control_url,omitempty"`
	// Indicates whether the instance should register as an Ephemeral node (https://tailscale.com/s/ephemeral-nodes).
	Ephemeral bool `json:"ephemeral,omitempty"`
	// The hostname of the node.
	//
	// System hostname is used by default.
	//
	// @since 1.14.0
	//
	//     On iOS, tvOS and Android, the device name is used by default.
	//
	// Example: `localhost`
	Hostname string `json:"hostname,omitempty"`
	// Indicates whether the node should accept routes advertised by other nodes.
	AcceptRoutes bool `json:"accept_routes,omitempty"`
	// The exit node name or IP address to use.
	ExitNode string `json:"exit_node,omitempty"`
	// @note
	//
	//     When the exit node does not have a corresponding advertised route, private traffics cannot be routed to the exit node even if `exit_node_allow_lan_access is` set.
	//
	// Indicates whether locally accessible subnets should be routed directly or via the exit node.
	ExitNodeAllowLANAccess bool `json:"exit_node_allow_lan_access,omitempty"`
	// CIDR prefixes to advertise into the Tailscale network as reachable through the current node.
	//
	// Example: `["192.168.1.1/24"]`
	AdvertiseRoutes []netip.Prefix `json:"advertise_routes,omitempty"`
	// Indicates whether the node should advertise itself as an exit node.
	AdvertiseExitNode bool `json:"advertise_exit_node,omitempty"`
	// @since 1.13.0
	//
	// Tags to advertise for this node, for ACL enforcement purposes.
	//
	// Example: `["tag:server"]`
	AdvertiseTags LegacyListable[string] `json:"advertise_tags,omitempty"`
	// @since 1.14.0
	//
	// The UDP port to listen on for WireGuard and peer-to-peer traffic.
	//
	// A port is automatically selected by default.
	ListenPort uint16 `json:"listen_port,omitempty"`
	// @since 1.13.0
	//
	// The port to listen on for incoming relay connections from other Tailscale nodes.
	RelayServerPort *uint16 `json:"relay_server_port,omitempty"`
	// @since 1.13.0
	//
	// Static endpoints to advertise for the relay server.
	RelayServerStaticEndpoints []netip.AddrPort `json:"relay_server_static_endpoints,omitempty"`
	// @since 1.13.0
	//
	// Create a system TUN interface for Tailscale.
	SystemInterface bool `json:"system_interface,omitempty"`
	// @since 1.13.0
	//
	// Custom TUN interface name. By default, `tailscale` (or `utun` on macOS) will be used.
	SystemInterfaceName string `json:"system_interface_name,omitempty"`
	// @since 1.13.0
	//
	// Override the TUN MTU. By default, Tailscale's own MTU is used.
	SystemInterfaceMTU uint32 `json:"system_interface_mtu,omitempty"`
	// UDP NAT expiration time.
	//
	// `5m` will be used by default.
	UDPTimeout UDPTimeoutCompat `json:"udp_timeout,omitempty"`
	// @since 1.14.0
	//
	// Run a Tailscale SSH server on tailnet port 22.
	//
	// Access is controlled by the SSH ACL in the Tailscale admin console, which maps each connection to a local user. How that user is resolved, and which users are allowed, depends on the platform:
	//
	// - **Linux** and **macOS**: switching to a user other than the one sing-box runs as requires running as root.
	// - **Windows**: in the command line client, sessions mapped to a local account other than the one sing-box runs as are refused. In the graphical client, there is no such restriction.
	// - **Android**: `root` is the superuser (UID 0) and `shell` is the ADB shell user (UID 2000); every other name is resolved as the package name of an installed application, running as that application's UID with its data directory as the home directory. `termux` is a shortcut for `com.termux`, and `sing-box` for the app's own package name; when Termux is installed, the `root` and `termux` users load the Termux environment. Running as the sing-box application itself requires no root, while any other user requires granted root access.
	// - **macOS**: the SSH server is only available in the standalone version and requires the Helper Service; the App Store version is not supported.
	// - **iOS**: the SSH server is only available in the jailbreak build; the App Store and TestFlight versions are not supported.
	// - **tvOS**: not yet supported.
	//
	// Object format:
	//
	// ```json
	// {
	//   "enabled": true,
	//   "disable_pty": false,
	//   "disable_sftp": false,
	//   "disable_forwarding": false
	// }
	// ```
	//
	// Setting `ssh_server` value to `true` is equivalent to `{ "enabled": true }`.
	SSHServer *TailscaleSSHServerOptions `json:"ssh_server,omitempty"`
	// @since 1.14.0
	//
	// The directory where files received from tailnet peers are stored.
	//
	// Relative paths are resolved against the working directory, as [state_directory](#state_directory)
	// is.
	//
	// `Taildrop` is used by default.
	TaildropDirectory string `json:"taildrop_directory,omitempty"`
	// @since 1.15.0
	//
	// Allow the endpoint to be disconnected when necessary.
	OnDemand bool `json:"on_demand,omitempty"`
}

type _TailscaleSSHServerOptions struct {
	// Enable the SSH server.
	Enabled bool `json:"enabled,omitempty"`
	// Refuse PTY allocation requests.
	DisablePTY bool `json:"disable_pty,omitempty"`
	// Refuse the SFTP subsystem.
	DisableSFTP bool `json:"disable_sftp,omitempty"`
	// Refuse local and remote TCP and Unix-socket forwarding, including SSH agent forwarding.
	DisableForwarding bool `json:"disable_forwarding,omitempty"`
}

type TailscaleSSHServerOptions _TailscaleSSHServerOptions

func (o TailscaleSSHServerOptions) MarshalJSON() ([]byte, error) {
	if !o.DisablePTY && !o.DisableSFTP && !o.DisableForwarding {
		return json.Marshal(o.Enabled)
	}
	return json.Marshal(_TailscaleSSHServerOptions(o))
}

func (o *TailscaleSSHServerOptions) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, &o.Enabled)
	if err == nil {
		return nil
	}
	return json.UnmarshalDisallowUnknownFields(bytes, (*_TailscaleSSHServerOptions)(o))
}

func (o TailscaleSSHServerOptions) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	objectForm := schema.StrictObject()
	err := builder.FlattenStruct(objectForm, reflect.TypeFor[TailscaleSSHServerOptions]())
	if err != nil {
		return nil, err
	}
	return schema.AnyOf(schema.BooleanNode(), objectForm), nil
}

type TailscaleDNSServerOptions struct {
	// @required
	//
	// The tag of the [Tailscale Endpoint](/configuration/endpoint/tailscale).
	Endpoint string `json:"endpoint,omitempty"`
	// Indicates whether default DNS resolvers should be accepted for fallback queries in addition to MagicDNS。
	//
	// if not enabled, `NXDOMAIN` will be returned for non-Tailscale domain queries.
	AcceptDefaultResolvers bool `json:"accept_default_resolvers,omitempty"`
	// @since 1.14.0
	//
	// When enabled, single-label queries (e.g. `my-device`) are retried against each Tailscale search domain until one resolves.
	AcceptSearchDomain bool `json:"accept_search_domain,omitempty"`
}

type TailscaleCertificateProviderOptions struct {
	// @required
	//
	// The tag of the [Tailscale endpoint](/configuration/endpoint/tailscale/) to reuse.
	//
	// [MagicDNS and HTTPS](https://tailscale.com/kb/1153/enabling-https) must be enabled in the Tailscale admin console.
	Endpoint string `json:"endpoint,omitempty"`
}

type DERPServiceOptions struct {
	ListenOptions
	InboundTLSOptionsContainer
	// @required
	//
	// Derper configuration file path.
	//
	// Example: `derper.key`
	ConfigPath string `json:"config_path,omitempty"`
	// Tailscale endpoints tags to verify clients.
	VerifyClientEndpoint badoption.Listable[string] `json:"verify_client_endpoint,omitempty"`
	// URL to verify clients.
	//
	// Object format:
	//
	// ```json
	// {
	//   "url": "",
	//
	//   ... // HTTP Client Fields
	// }
	// ```
	//
	// Setting Array value to a string `__URL__` is equivalent to configuring:
	//
	// ```json
	// { "url": __URL__ }
	// ```
	VerifyClientURL badoption.Listable[*DERPVerifyClientURLOptions] `json:"verify_client_url,omitempty"`
	// @since 1.15.0
	//
	// Tailcat inbound tags to verify clients.
	VerifyClientInbound badoption.Listable[string] `json:"verify_client_inbound,omitempty"`
	// @since 1.15.0
	//
	// Tailcat public keys to verify clients.
	VerifyClientKey badoption.Listable[string] `json:"verify_client_key,omitempty"`
	// What to serve at the root path. It may be left empty (the default, for a default homepage), `blank` for a blank page, or a URL to redirect to
	Home string `json:"home,omitempty"`
	// Mesh with other DERP servers.
	//
	// Object format:
	//
	// ```json
	// {
	//   "server": "",
	//   "server_port": "",
	//   "host": "",
	//   "tls": {},
	//
	//   ... // Dial Fields
	// }
	// ```
	//
	// Object fields:
	//
	// - `server`: **Required** DERP server address.
	// - `server_port`: **Required** DERP server port.
	// - `host`: Custom DERP hostname.
	// - `tls`: [TLS](/configuration/shared/tls/#outbound)
	// - `Dial Fields`: [Dial Fields](/configuration/shared/dial/)
	MeshWith badoption.Listable[*DERPMeshOptions] `json:"mesh_with,omitempty"`
	// Pre-shared key for DERP mesh.
	MeshPSK string `json:"mesh_psk,omitempty"`
	// Pre-shared key file for DERP mesh.
	MeshPSKFile string `json:"mesh_psk_file,omitempty"`
	// STUN server listen options.
	//
	// Object format:
	//
	// ```json
	// {
	//   "enabled": true,
	//
	//   ... // Listen Fields
	// }
	// ```
	//
	// Object fields:
	//
	// - `enabled`: **Required** Enable STUN server.
	// - `listen`: **Required** STUN server listen address, default to `::`.
	// - `listen_port`: **Required** STUN server listen port, default to `3478`.
	// - `other Listen Fields`: [Listen Fields](/configuration/shared/listen/)
	//
	// Setting `stun` value to a number `__PORT__` is equivalent to configuring:
	//
	// ```json
	// { "enabled": true, "listen_port": __PORT__ }
	// ```
	STUN *DERPSTUNListenOptions `json:"stun,omitempty"`
}

type _DERPVerifyClientURLBase struct {
	URL string `json:"url,omitempty"`
}

type _DERPVerifyClientURLOptions struct {
	_DERPVerifyClientURLBase
	HTTPClientOptions
}

type DERPVerifyClientURLOptions _DERPVerifyClientURLOptions

func (d DERPVerifyClientURLOptions) ServerIsDomain() bool {
	verifyURL, err := url.Parse(d.URL)
	if err != nil {
		return false
	}
	return M.ParseSocksaddr(verifyURL.Hostname()).IsDomain()
}

func (d DERPVerifyClientURLOptions) MarshalJSON() ([]byte, error) {
	if d.URL != "" && d.HTTPClientOptions.IsEmpty() {
		return json.Marshal(d.URL)
	}
	return badjson.MarshallObjects(d._DERPVerifyClientURLBase, HTTPClient(d.HTTPClientOptions))
}

func (d *DERPVerifyClientURLOptions) UnmarshalJSON(bytes []byte) error {
	var stringValue string
	err := json.Unmarshal(bytes, &stringValue)
	if err == nil {
		*d = DERPVerifyClientURLOptions{
			_DERPVerifyClientURLBase: _DERPVerifyClientURLBase{URL: stringValue},
		}
		return nil
	}
	err = json.Unmarshal(bytes, &d._DERPVerifyClientURLBase)
	if err != nil {
		return err
	}
	var client HTTPClient
	err = badjson.UnmarshallExcluded(bytes, &d._DERPVerifyClientURLBase, &client)
	if err != nil {
		return err
	}
	d.HTTPClientOptions = HTTPClientOptions(client)
	return nil
}

func (d DERPVerifyClientURLOptions) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	objectForm, err := describeHTTPClientObject(builder)
	if err != nil {
		return nil, err
	}
	objectForm.Properties.Put("url", schema.StringNode())
	return schema.AnyOf(schema.StringNode(), objectForm), nil
}

type DERPMeshOptions struct {
	ServerOptions
	Host string `json:"host,omitempty"`
	OutboundTLSOptionsContainer
	DialerOptions
}

type _DERPSTUNListenOptions struct {
	Enabled bool `json:"enabled,omitempty"`
	ListenOptions
}

type DERPSTUNListenOptions _DERPSTUNListenOptions

func (d DERPSTUNListenOptions) MarshalJSON() ([]byte, error) {
	if d.Enabled && d.ListenPort != 0 && d.ListenOptions == (ListenOptions{ListenPort: d.ListenPort}) {
		return json.Marshal(d.ListenPort)
	}
	return json.Marshal(_DERPSTUNListenOptions(d))
}

func (d *DERPSTUNListenOptions) UnmarshalJSON(bytes []byte) error {
	var portValue uint16
	err := json.Unmarshal(bytes, &portValue)
	if err == nil {
		d.Enabled = true
		d.ListenPort = portValue
		return nil
	}
	return json.Unmarshal(bytes, (*_DERPSTUNListenOptions)(d))
}

func (d DERPSTUNListenOptions) DescribeSchema(builder schema.Builder) (*schema.Node, error) {
	objectForm := schema.StrictObject()
	err := builder.FlattenStruct(objectForm, reflect.TypeFor[DERPSTUNListenOptions]())
	if err != nil {
		return nil, err
	}
	return schema.AnyOf(schema.UnsignedNode(16), objectForm), nil
}
