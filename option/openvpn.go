package option

import (
	"net/netip"

	"github.com/sagernet/sing/common/auth"
	"github.com/sagernet/sing/common/json/badoption"
)

type OpenVPNEndpointOptions struct {
	// Use a system interface.
	//
	// Requires privilege and cannot conflict with existing system interfaces.
	//
	// If disabled, sing-box uses the internal network stack.
	//
	System bool `json:"system,omitempty"`

	// Custom interface name for the system interface.
	//
	// An automatically generated `ovpn` interface name will be used by default.
	//
	Name string `json:"name,omitempty"`

	// OpenVPN interface MTU.
	//
	// When empty, `1500` is used until a server-pushed MTU is received.
	//
	// @Example: 1500
	MTU uint32 `json:"mtu,omitempty"`

	// @Omit
	UDPMapping UDPNATBehavior `json:"udp_mapping,omitempty"`

	// @Omit
	UDPFiltering UDPNATBehavior `json:"udp_filtering,omitempty"`

	// @Omit
	UDPNATMax uint32 `json:"udp_nat_max,omitempty"`
}

// @Tag ovpn-client
type OpenVPNClientEndpointOptions struct {
	DialerOptions

	// @Override server:
	//    OpenVPN server address.
	//
	//    Either `server` or `servers` is required.
	//
	//    Conflict with `servers`.
	//
	//    @Example: "127.0.0.1"
	//
	// @Override server_port:
	//    OpenVPN server port.
	//
	//    Required when `server` is set.
	//
	//    @Example: 1194
	ServerOptions
	OpenVPNEndpointOptions

	// OpenVPN session mode, one of `tls` or `static_key`.
	//
	// `static_key` is a deprecated OpenVPN mode without a TLS control channel or
	// forward secrecy.
	//
	// @Default: tls
	//
	// @Example: "tls"
	Mode string `json:"mode,omitempty" enum:"tls,static_key"`

	// Default OpenVPN transport network, one of `udp` or `tcp`.
	//
	// This value applies to `server` and to `servers` entries without their own `network`.
	//
	// @Default: udp
	//
	// @Example: "udp"
	Network string `json:"network,omitempty" enum:"udp,udp4,udp6,tcp,tcp4,tcp6"`

	// List of OpenVPN servers.
	//
	// The client tries the servers in order and moves to the next server when a connection fails.
	//
	// Either `server` or `servers` is required.
	//
	// Conflict with `server`.
	//
	// @Example:
	//    [
	//      {
	//        "server": "127.0.0.1",
	//        "server_port": 1194,
	//        "network": "udp"
	//      }
	//    ]
	Servers []OpenVPNRemoteOptions `json:"servers,omitempty"`

	// Randomize the `servers` order before connecting.
	//
	// Disabled by default.
	RemoteRandom bool `json:"remote_random,omitempty"`

	// Local IPv4 and IPv6 tunnel prefixes.
	//
	// At least one address is required in `static_key` mode. In TLS mode these
	// addresses are optional and can be replaced by addresses pulled from the
	// server.
	Address badoption.Listable[netip.Prefix] `json:"address,omitempty"`

	// IPv4 tunnel peer address and VPN gateway.
	//
	// Required when an IPv4 `address` is configured in `static_key` mode.
	PeerAddress *badoption.Addr `json:"peer_address,omitempty"`

	// IPv6 tunnel peer address and VPN gateway.
	//
	// Required when an IPv6 `address` is configured in `static_key` mode.
	PeerAddressIPv6 *badoption.Addr `json:"peer_address_ipv6,omitempty"`

	// Tunnel topology, one of `net30`, `p2p`, or `subnet`.
	//
	// The topology pulled from the server is used when empty in TLS mode.
	Topology string `json:"topology,omitempty" enum:"net30,p2p,subnet"`

	// Username for OpenVPN username/password authentication.
	//
	// Only available in TLS mode.
	Username string `json:"username,omitempty"`

	// Password for OpenVPN username/password authentication.
	Password string `json:"password,omitempty"`

	// Behavior after username/password authentication fails, one of `none`, `nointeract`, or `interact`.
	//
	// `none` will be used by default and treats a permanent authentication failure as terminal.
	//
	// `nointeract` and `interact` allow authentication retries.
	//
	// @Example: "none"
	AuthRetry string `json:"auth_retry,omitempty" enum:"none,nointeract,interact"`

	// Static challenge text shown when requesting an authentication response.
	StaticChallenge string `json:"static_challenge,omitempty"`

	// Show the static challenge response as plain text.
	StaticChallengeEcho bool `json:"static_challenge_echo,omitempty"`

	// OpenVPN static key content.
	//
	// Required in `static_key` mode.
	//
	// Conflict with `static_key_path`.
	StaticKey badoption.Listable[string] `json:"static_key,omitempty"`

	// OpenVPN static key path.
	//
	// Required in `static_key` mode when `static_key` is not set.
	//
	// Conflict with `static_key`.
	StaticKeyPath string `json:"static_key_path,omitempty"`

	// Static key direction, one of `server` or `client`.
	//
	// The key is used bidirectionally if empty. Only available in `static_key` mode.
	KeyDirection string `json:"key_direction,omitempty" enum:"server,client"`

	// Required in TLS mode.
	//
	// OpenVPN control channel TLS configuration.
	//
	// @Example:
	//    {
	//      "server_name": "",
	//      "server_name_type": "name",
	//      "certificate": [],
	//      "certificate_path": "",
	//      "client_certificate": [],
	//      "client_certificate_path": "",
	//      "client_key": [],
	//      "client_key_path": "",
	//      "peer_fingerprint": [],
	//      "crl_path": "",
	//      "remote_certificate_ku": [],
	//      "remote_certificate_eku": "",
	//      "remote_certificate_tls": "",
	//      "certificate_profile": "",
	//      "ns_certificate_type": "",
	//      "version_min": "1.2",
	//      "version_max": "",
	//      "cipher": "",
	//      "groups": "",
	//      "control_wrap": {
	//        "type": "",
	//        "key": [],
	//        "key_path": "",
	//        "direction": ""
	//      }
	//    }
	TLS *OpenVPNOutboundTLSOptions `json:"tls,omitempty"`

	// Data-channel cipher used in `static_key` mode.
	//
	// The upstream static-key default `BF-CBC` is used when empty. Static-key ciphers include `BF-CBC`,
	// `CAST5-CBC`, `DES-CBC`, `DES-EDE-CBC`, `DES-EDE3-CBC`, the AES-CBC,
	// ARIA-CBC, and Camellia-CBC families, `SEED-CBC`, `SM4-CBC`, and `NONE`.
	//
	// Only available in `static_key` mode.
	Cipher string `json:"cipher,omitempty"`

	// Allowed OpenVPN data channel ciphers.
	//
	// Only available in TLS mode.
	//
	// `AES-256-GCM`, `AES-128-GCM`, and `CHACHA20-POLY1305` are used by default.
	//
	// The AES-GCM family includes `AES-192-GCM`. Retained ciphers include the CBC,
	// CFB, and OFB forms of AES, ARIA, Camellia, DES, Blowfish, and CAST5, the CBC,
	// CFB, and OFB forms of SEED and SM4, and `NONE`. CFB and OFB are available only
	// in TLS mode.
	DataCiphers LegacyListable[string] `json:"data_ciphers,omitempty"`

	// Data channel cipher for peers that do not support cipher negotiation.
	//
	// Disabled by default.
	//
	// Only available in TLS mode.
	DataCiphersFallback string `json:"data_ciphers_fallback,omitempty"`

	// OpenVPN data channel authentication digest.
	//
	// `SHA1` will be used by default. It only applies to non-AEAD data ciphers and `tls_auth`.
	Auth string `json:"auth,omitempty"`

	// Maximum OpenVPN UDP packet size used to clamp the MSS of TCP connections sent through the tunnel.
	//
	// When empty, the upstream OpenVPN default is used: `fragment` when configured,
	// otherwise `1492` for the default tunnel MTU or the configured tunnel MTU.
	MSSFix uint32 `json:"mss_fix,omitempty"`

	// Disable MSS clamping, including the default clamp.
	//
	// Conflict with `mss_fix` and `mss_fix_mode`.
	MSSFixDisabled bool `json:"mss_fix_disabled,omitempty"`

	// OpenVPN MSS calculation mode for an explicit `mss_fix`, one of `mtu` or `fixed`.
	//
	// An empty value uses the normal OpenVPN encapsulation-aware calculation. `mtu` also accounts for the outer IP and UDP/TCP transport headers. `fixed` treats `mss_fix` as an inner IPv4 packet size.
	//
	// Requires `mss_fix`.
	MSSFixMode string `json:"mss_fix_mode,omitempty" enum:"mtu,fixed"`

	// Maximum OpenVPN UDP packet size used for OpenVPN data channel fragmentation.
	//
	// Disabled when `0`. A non-zero value must be at least `68`.
	//
	// Conflict with TCP transport.
	Fragment uint32 `json:"fragment,omitempty"`

	// UDP data-channel replay window size. `64` will be used by default. The maximum is `65536`.
	ReplayWindow uint32 `json:"replay_window,omitempty"`

	// UDP data-channel replay window duration. `15s` will be used by default and the maximum is `10m`.
	//
	// The value must use whole seconds.
	ReplayWindowTime badoption.Duration `json:"replay_window_time,omitempty"`

	// OpenVPN `compress` framing mode, one of `none`, `no`, `lz4`, `lz4-v2`, `stub`, `stub-v2`, `disabled`, or `off`.
	//
	// Disabled by default.
	Compression string `json:"compression,omitempty" enum:"none,no,lz4,lz4-v2,stub,stub-v2,disabled,off"`

	// OpenVPN `comp-lzo` mode, one of `none`, `no`, `yes`, `adaptive`, `asym`, `disabled`, or `off`.
	//
	// Disabled by default.
	CompressionLZO string `json:"compression_lzo,omitempty" enum:"none,no,yes,adaptive,asym,disabled,off"`

	// Policy for compression pushed by the server, one of `no`, `asym`, or `yes`.
	//
	// `no` will be used by default and permits only compression stub framing. `asym` accepts compressed packets from the server but does not compress outgoing packets. `yes` is a legacy alias for `asym`.
	//
	// Conflict with non-stub compression enabled by `compression` or `compression_lzo` when set to `no`.
	//
	// @Example: "no"
	AllowCompression string `json:"allow_compression,omitempty" enum:"no,asym,yes"`

	// Ignore routes, DNS and DHCP settings, route metrics, `redirect-gateway`,
	// `redirect-private`, `block-ipv6`, and `block-outside-dns` pushed by the server.
	//
	// Disabled by default.
	RouteNoPull bool `json:"route_no_pull,omitempty"`

	// Ordered filters for options pushed by the server.
	//
	// The first filter whose `text` is a case-sensitive prefix of the complete pushed option is applied. Options that match no filter are accepted.
	//
	// @Example:
	//    [
	//      {
	//        "action": "ignore",
	//        "text": "route "
	//      }
	//    ]
	PullFilters []OpenVPNPullFilterOptions `json:"pull_filters,omitempty"`

	// IPv4 and IPv6 prefixes preferred by sing-box routing for this OpenVPN endpoint.
	//
	// These routes are used in addition to routes accepted from the server.
	Routes LegacyListable[netip.Prefix] `json:"routes,omitempty"`

	// IPv4 gateway for routes through the OpenVPN endpoint.
	//
	// When empty, the VPN gateway received from the server is used.
	RouteGateway *badoption.Addr `json:"route_gateway,omitempty"`

	// Default metric for routes through the OpenVPN endpoint.
	//
	// The platform default is used when `0`.
	RouteMetric int `json:"route_metric,omitempty"`

	// Prefer the OpenVPN endpoint for all IPv4 destinations in sing-box routing.
	//
	// Disabled by default.
	RedirectGateway bool `json:"redirect_gateway,omitempty"`

	// OpenVPN `redirect-gateway` flags.
	//
	// `!ipv4` disables IPv4 preference, `def1` represents it with two `/1`
	// prefixes, and `ipv6` also prefers the upstream-specific IPv6 prefixes.
	// `local`, `autolocal`, `bypass-dhcp` and `bypass-dns` have no effect.
	// `block-local` is not supported.
	//
	// Empty by default.
	RedirectGatewayFlags LegacyListable[string] `json:"redirect_gateway_flags,omitempty"`

	// Accept `redirect_gateway_flags` without adding a default-route preference.
	//
	// Disabled by default.
	RedirectPrivate bool `json:"redirect_private,omitempty"`

	// Reject IPv6 traffic locally instead of sending it through the VPN.
	//
	// Disabled by default.
	BlockIPv6 bool `json:"block_ipv6,omitempty"`

	// Interval after which the client sends a data-channel ping when no packet has been sent to the server.
	//
	// A server-pushed OpenVPN `ping` value overrides this value.
	//
	// The value must use whole seconds.
	//
	// Disabled by default.
	PingInterval badoption.Duration `json:"ping_interval,omitempty"`

	// Time without receiving a packet after which the client reconnects to the server.
	//
	// A server-pushed OpenVPN `ping-restart` value overrides this value.
	//
	// The value must use whole seconds.
	//
	// When empty, `120s` is used for UDP connections with pull enabled until the
	// server pushes another value. No default receive timeout is used for TCP.
	PingRestart badoption.Duration `json:"ping_restart,omitempty"`

	// Disable the initial `120s` UDP pull timeout and any locally configured ping restart timeout.
	//
	// Conflict with `ping_restart`.
	PingRestartDisabled bool `json:"ping_restart_disabled,omitempty"`

	// OpenVPN TLS renegotiation interval.
	//
	// When empty, the OpenVPN default `1h` is used.
	RenegotiateInterval badoption.Duration `json:"renegotiate_interval,omitempty"`

	// Disable time-based TLS renegotiation, including the default interval.
	//
	// Conflict with `renegotiate_interval`.
	RenegotiateDisabled bool `json:"renegotiate_disabled,omitempty"`

	// Renegotiate data-channel keys after this many bytes. `0` uses the cipher-dependent OpenVPN default.
	RenegotiateBytes uint64 `json:"renegotiate_bytes,omitempty"`

	// Renegotiate data-channel keys after this many packets. `0` uses the cipher-dependent OpenVPN default.
	RenegotiatePackets uint64 `json:"renegotiate_packets,omitempty"`

	// Initial retransmission timeout for TLS control packets. The OpenVPN default `2s` is used when empty.
	TLSTimeout badoption.Duration `json:"tls_timeout,omitempty"`

	// Maximum time allowed for the initial TLS handshake and each renegotiation. The OpenVPN default `1m` is used when empty.
	HandshakeWindow badoption.Duration `json:"handshake_window,omitempty"`

	// Number of OpenVPN exit notifications sent when closing a UDP connection.
	//
	// Notifications are sent one second apart. Disabled when `0`.
	ExplicitExitNotify uint32 `json:"explicit_exit_notify,omitempty"`

	// @Omit
	UDPTimeout UDPTimeoutCompat `json:"udp_timeout,omitempty"`

	// Allow the endpoint to be disconnected when necessary.
	//
	// @Since 1.15.0
	OnDemand bool `json:"on_demand,omitempty"`
}

// @Tag ovpn-server
type OpenVPNServerEndpointOptions struct {
	ListenOptions

	// @Override system:
	//    Use system interface.
	//
	//    Requires privilege and cannot conflict with existing system interfaces.
	//
	//    If disabled, sing-box uses the internal network stack.
	//
	// @Override name:
	//    Custom interface name for system interface.
	//
	//    An automatically generated `ovpn` interface name will be used by default.
	//
	// @Override mtu:
	//    OpenVPN interface MTU.
	//
	//    `1500` will be used by default.
	OpenVPNEndpointOptions

	// OpenVPN session mode, one of `tls` or `static_key`.
	//
	// `static_key` serves one peer without a TLS control channel or forward secrecy.
	//
	// @Default: tls
	//
	// @Example: "tls"
	Mode string `json:"mode,omitempty" enum:"tls,static_key"`

	// OpenVPN transport network, one of `udp` or `tcp`.
	//
	// Only one transport network is served per endpoint; to serve both TCP and UDP,
	// configure two endpoints with separate `address` subnets.
	//
	// @Default: udp
	//
	// @Example: "udp"
	Network string `json:"network,omitempty" enum:"tcp,udp"`

	// Fixed remote peer address for a UDP `static_key` server.
	//
	// Required with `remote_port` in UDP `static_key` mode.
	Remote string `json:"remote,omitempty"`

	// Fixed remote peer port for a UDP `static_key` server.
	//
	// Required with `remote` in UDP `static_key` mode.
	RemotePort uint16 `json:"remote_port,omitempty"`

	// Maximum number of established and pending TLS client sessions.
	//
	// `1024` will be used by default. The maximum value is `16777215`.
	//
	// Must be `0` or `1` in `static_key` mode.
	//
	// @Example: 1024
	MaxClients int `json:"max_clients,omitempty"`

	// List of OpenVPN server address prefixes.
	//
	// At most one IPv4 prefix and one IPv6 prefix are supported.
	//
	// The prefix address is assigned to the server interface. The masked prefix is used as the client address pool and route.
	//
	// The first IPv4 and IPv6 prefix addresses are used as the endpoint's local addresses.
	//
	// In `static_key` mode these are the local tunnel prefixes rather than address pools.
	//
	// @Required
	Address badoption.Listable[netip.Prefix] `json:"address"`

	// IPv4 tunnel peer address.
	//
	// Required when an IPv4 `address` is configured in `static_key` mode.
	PeerAddress *badoption.Addr `json:"peer_address,omitempty"`

	// IPv6 tunnel peer address.
	//
	// Required when an IPv6 `address` is configured in `static_key` mode.
	PeerAddressIPv6 *badoption.Addr `json:"peer_address_ipv6,omitempty"`

	// OpenVPN topology pushed to clients, one of `subnet`, `p2p` or `net30`.
	//
	// `subnet` will be used by default in TLS mode. `p2p` will be used by default in
	// `static_key` mode.
	//
	// @Example: "subnet"
	Topology string `json:"topology,omitempty" enum:"net30,p2p,subnet"`

	// Allow multiple active clients with the same authenticated certificate common name or username.
	//
	// When disabled, a newly authenticated session replaces the existing session with the same identity and reuses its tunnel address when available.
	//
	// Disabled by default.
	//
	// Only available in TLS mode.
	DuplicateCN bool `json:"duplicate_cn,omitempty"`

	// List of OpenVPN username/password users.
	//
	// If set, clients must pass username/password authentication in addition to any certificate policy configured by `tls.verify_client_certificate`.
	//
	// Only available in TLS mode.
	//
	// @Example:
	//    [
	//      {
	//        "username": "",
	//        "password": ""
	//      }
	//    ]
	// @Doc users.username:
	//
	//	Username.
	//
	// @Doc users.password:
	//
	//	Password.
	Users []auth.User `json:"users,omitempty"`

	// OpenVPN static key content.
	//
	// Required in `static_key` mode.
	//
	// Conflict with `static_key_path`.
	StaticKey badoption.Listable[string] `json:"static_key,omitempty"`

	// OpenVPN static key path.
	//
	// Required in `static_key` mode when `static_key` is not set.
	//
	// Conflict with `static_key`.
	StaticKeyPath string `json:"static_key_path,omitempty"`

	// Static key direction, one of `server` or `client`.
	//
	// The key is used bidirectionally if empty. Conventionally the server uses
	// `server` and the peer uses `client`.
	//
	// Only available in `static_key` mode.
	KeyDirection string `json:"key_direction,omitempty" enum:"server,client"`

	// Required in TLS mode.
	//
	// OpenVPN control channel TLS configuration.
	//
	// @Example:
	//    {
	//      "certificate": [],
	//      "certificate_path": "",
	//      "key": [],
	//      "key_path": "",
	//      "client_certificate": [],
	//      "client_certificate_path": "",
	//      "verify_client_certificate": "require",
	//      "client_name": "",
	//      "client_name_type": "name",
	//      "peer_fingerprint": [],
	//      "crl_path": "",
	//      "remote_certificate_ku": [],
	//      "remote_certificate_eku": "",
	//      "remote_certificate_tls": "",
	//      "certificate_profile": "",
	//      "ns_certificate_type": "",
	//      "version_min": "1.2",
	//      "version_max": "",
	//      "cipher": "",
	//      "groups": "",
	//      "control_wrap": {
	//        "type": "tls_crypt",
	//        "key": [],
	//        "key_path": "",
	//        "direction": "",
	//        "force_cookie": false
	//      }
	//    }
	TLS *OpenVPNInboundTLSOptions `json:"tls,omitempty"`

	// Data-channel cipher used in `static_key` mode.
	//
	// The upstream static-key default `BF-CBC` is used when empty. Supported
	// static-key ciphers are the AES-CBC, ARIA-CBC, Camellia-CBC, DES-CBC,
	// Blowfish-CBC, CAST5-CBC families, `SEED-CBC`, `SM4-CBC`, and `NONE`.
	//
	// Only available in `static_key` mode.
	Cipher string `json:"cipher,omitempty"`

	// Allowed OpenVPN data channel ciphers.
	//
	// `AES-256-GCM`, `AES-128-GCM` and `CHACHA20-POLY1305` are used by default.
	//
	// The AES-GCM family includes `AES-192-GCM`. Retained ciphers include the CBC,
	// CFB, and OFB forms of AES, ARIA, Camellia, DES, Blowfish, and CAST5, the CBC,
	// CFB, and OFB forms of SEED and SM4, and `NONE`. CFB and OFB are available only
	// in TLS mode.
	//
	// Only available in TLS mode.
	DataCiphers LegacyListable[string] `json:"data_ciphers,omitempty"`

	// OpenVPN data channel cipher for legacy clients that do not support cipher negotiation.
	//
	// Equivalent to OpenVPN `data-ciphers-fallback`.
	//
	// Disabled by default.
	//
	// Only available in TLS mode.
	DataCiphersFallback string `json:"data_ciphers_fallback,omitempty"`

	// OpenVPN data channel authentication digest.
	//
	// `SHA1` will be used by default. It only applies to non-AEAD data ciphers and `tls_auth`.
	Auth string `json:"auth,omitempty"`

	// Maximum encapsulated packet size used to clamp TCP MSS. The upstream default calculation uses `1492` with the default MTU.
	MSSFix uint32 `json:"mss_fix,omitempty"`

	// Disable MSS clamping, including the default clamp.
	MSSFixDisabled bool `json:"mss_fix_disabled,omitempty"`

	// Calculation mode for an explicit `mss_fix`, one of `mtu` or `fixed`. Requires `mss_fix`.
	MSSFixMode string `json:"mss_fix_mode,omitempty" enum:"mtu,fixed"`

	// UDP data-channel replay window size. `64` will be used by default.
	ReplayWindow uint32 `json:"replay_window,omitempty"`

	// UDP replay window duration. `15s` will be used by default. The value must use whole seconds.
	ReplayWindowTime badoption.Duration `json:"replay_window_time,omitempty"`

	// Options pushed to clients.
	//
	// @Example:
	//    {
	//      "routes": [],
	//      "dns": [],
	//      "dns_servers": [],
	//      "search_domains": [],
	//      "dhcp_options": [],
	//      "redirect_gateway": false,
	//      "redirect_gateway_flags": [],
	//      "block_outside_dns": false,
	//      "ping_interval": "",
	//      "ping_restart": ""
	//    }
	Push *OpenVPNPushOptions `json:"push,omitempty"`

	// Interval after which the server sends a data-channel ping when no packet has been sent to a client.
	//
	// This value applies to the server. Use `push.ping_interval` to configure clients.
	//
	// The value must use whole seconds.
	//
	// Disabled by default.
	PingInterval badoption.Duration `json:"ping_interval,omitempty"`

	// Time without receiving a packet after which the server closes the client session.
	//
	// This value applies to the server. Use `push.ping_restart` to configure clients.
	//
	// The value must use whole seconds.
	//
	// Disabled by default.
	PingRestart badoption.Duration `json:"ping_restart,omitempty"`

	// OpenVPN TLS renegotiation interval.
	//
	// When empty, the OpenVPN default `1h` is used.
	//
	// Only available in TLS mode.
	RenegotiateInterval badoption.Duration `json:"renegotiate_interval,omitempty"`

	// Disable time-based TLS renegotiation, including the default interval.
	//
	// Only available in TLS mode.
	RenegotiateDisabled bool `json:"renegotiate_disabled,omitempty"`

	// Renegotiate data-channel keys after this many bytes. `0` uses the cipher-dependent OpenVPN default.
	//
	// Only available in TLS mode.
	RenegotiateBytes uint64 `json:"renegotiate_bytes,omitempty"`

	// Renegotiate data-channel keys after this many packets. `0` uses the cipher-dependent OpenVPN default.
	//
	// Only available in TLS mode.
	RenegotiatePackets uint64 `json:"renegotiate_packets,omitempty"`

	// Maximum time allowed for the initial TLS handshake and each TLS renegotiation.
	//
	// Only available in TLS mode.
	//
	// @Default: 1m
	//
	// @Example: "1m"
	HandshakeWindow badoption.Duration `json:"handshake_window,omitempty"`
}

type OpenVPNRemoteOptions struct {
	// @Override server:
	//    @Required
	//
	//    OpenVPN server address.
	//
	// @Override server_port:
	//    @Required
	//
	//    OpenVPN server port.
	//
	ServerOptions

	// OpenVPN transport network for this server, one of `udp` or `tcp`.
	//
	// The top-level `network` will be used by default.
	Network string `json:"network,omitempty" enum:"udp,udp4,udp6,tcp,tcp4,tcp6"`
}

type OpenVPNPullFilterOptions struct {
	// Filter action, one of `accept`, `ignore`, or `reject`.
	//
	// `accept` applies the option, `ignore` discards it, and `reject` terminates the connection.
	//
	// @Required
	Action string `json:"action" enum:"accept,ignore,reject"`

	// Case-sensitive prefix to match against the pushed option name and value.
	//
	// For example, `route ` matches pushed IPv4 route options without matching `route-gateway`.
	//
	// @Required
	Text string `json:"text"`
}

type OpenVPNOutboundTLSOptions struct {
	// Expected server certificate name.
	//
	// Certificate name verification is disabled if empty.
	ServerName string `json:"server_name,omitempty"`

	// Certificate field matched by `tls.server_name`, one of `subject`, `name`, or `name-prefix`.
	//
	// `name` will be used by default when `tls.server_name` is set.
	//
	// `subject` matches the full certificate subject, `name` matches the common name exactly, and `name-prefix` matches a common name prefix.
	ServerNameType string `json:"server_name_type,omitempty" enum:"subject,name,name-prefix"`

	// Trusted CA certificate content.
	//
	// One of `tls.certificate`, `tls.certificate_path`, or `tls.peer_fingerprint` is required.
	//
	// Conflict with `tls.certificate_path`.
	Certificate badoption.Listable[string] `json:"certificate,omitempty"`

	// Trusted CA certificate path.
	//
	// One of `tls.certificate`, `tls.certificate_path`, or `tls.peer_fingerprint` is required.
	//
	// Conflict with `tls.certificate`.
	CertificatePath string `json:"certificate_path,omitempty"`

	// Client certificate content.
	//
	// Conflict with `tls.client_certificate_path`.
	ClientCertificate badoption.Listable[string] `json:"client_certificate,omitempty"`

	// Client certificate path.
	//
	// Conflict with `tls.client_certificate`.
	ClientCertificatePath string `json:"client_certificate_path,omitempty"`

	// Client private key content.
	//
	// Conflict with `tls.client_key_path`.
	ClientKey badoption.Listable[string] `json:"client_key,omitempty"`

	// Client private key path.
	//
	// Conflict with `tls.client_key`.
	//
	// The client certificate and key must both be set or both be empty.
	ClientKeyPath string `json:"client_key_path,omitempty"`

	// Allowed SHA-256 fingerprints of the server leaf certificate.
	//
	// Each fingerprint must be 64 lowercase hexadecimal characters without separators.
	PeerFingerprint badoption.Listable[string] `json:"peer_fingerprint,omitempty"`

	// Path to a PEM or DER certificate revocation list used to reject revoked server certificates.
	//
	// Disabled by default.
	CRLPath string `json:"crl_path,omitempty"`

	// Required server certificate key usage masks, written as hexadecimal values in OpenVPN `remote-cert-ku` format.
	//
	// The certificate must contain all bits from at least one configured mask.
	//
	// Disabled by default.
	RemoteCertificateKU badoption.Listable[string] `json:"remote_certificate_ku,omitempty"`

	// Required server certificate extended key usage.
	//
	// OpenSSL names, object identifiers, and the aliases `server` and `client` are accepted.
	//
	// When set, this field replaces the default `tls.remote_certificate_tls` check.
	//
	// Conflict with an explicitly configured `tls.remote_certificate_tls`.
	RemoteCertificateEKU string `json:"remote_certificate_eku,omitempty"`

	// Peer certificate purpose check, one of `server`, `client`, or `none`.
	//
	// `none` disables the certificate purpose check.
	//
	// Conflict with `tls.remote_certificate_eku`.
	//
	// @Default: server
	RemoteCertificateTLS string `json:"remote_certificate_tls,omitempty" enum:"server,client,none"`

	// Certificate profile, one of `insecure`, `legacy`, `preferred`, or `suiteb`.
	//
	// `insecure` accepts MD5- and SHA-1-signed certificate chains and smaller legacy
	// keys for compatibility with immutable peers. Use it only when the peer cannot
	// be upgraded. `legacy` accepts SHA-1 but rejects MD5 signatures; `preferred`
	// requires stronger signatures and keys.
	//
	// When `suiteb` is selected, the TLS 1.2 cipher list defaults to the Suite B ECDHE-ECDSA AES-GCM suites when `tls.cipher` is empty, and the TLS groups default to `secp256r1` and `secp384r1` when `tls.groups` is empty.
	//
	// @Default: legacy
	CertificateProfile string `json:"certificate_profile,omitempty" enum:"legacy,preferred,insecure,suiteb"`

	// Deprecated Netscape certificate type check, one of `server` or `client`.
	//
	// Disabled by default. Prefer `tls.remote_certificate_tls`.
	NSCertificateType string `json:"ns_certificate_type,omitempty" enum:"server,client"`

	// Minimum TLS version, one of `1.0`, `1.1`, `1.2`, or `1.3`.
	//
	// @Default: 1.2
	VersionMin string `json:"version_min,omitempty" enum:"1.0,1.1,1.2,1.3"`

	// Maximum TLS version, one of `1.0`, `1.1`, `1.2`, or `1.3`.
	//
	// The maximum supported version will be used by default.
	//
	// The value cannot be lower than `tls.version_min`.
	VersionMax string `json:"version_max,omitempty" enum:"1.0,1.1,1.2,1.3"`

	// Colon-separated OpenSSL or IANA cipher suite names allowed for TLS 1.2 and earlier, such as `ECDHE-ECDSA-AES256-GCM-SHA384` or `TLS-ECDHE-ECDSA-WITH-AES-256-GCM-SHA384`. The order of the list is ignored.
	//
	// The default TLS cipher suites are used when empty.
	Cipher string `json:"cipher,omitempty"`

	// Colon-separated TLS key exchange groups. The order of the list is ignored.
	//
	// Supported groups are `X25519`, `X25519MLKEM768`, `secp256r1`, `secp384r1`, and `secp521r1`, including their common OpenSSL and NIST aliases. Names are case-insensitive.
	//
	// The default TLS groups are used when empty.
	Groups string `json:"groups,omitempty"`

	// OpenVPN control channel wrapping.
	//
	// Equivalent to OpenVPN `tls-auth`, `tls-crypt`, and `tls-crypt-v2`.
	//
	// Disabled if empty.
	ControlWrap *OpenVPNControlWrapOptions `json:"control_wrap,omitempty"`
}

type OpenVPNInboundTLSOptions struct {
	// TLS server certificate content.
	//
	// Either `tls.certificate` or `tls.certificate_path` is required.
	//
	// Conflict with `tls.certificate_path`.
	Certificate badoption.Listable[string] `json:"certificate,omitempty"`

	// TLS server certificate path.
	//
	// Either `tls.certificate` or `tls.certificate_path` is required.
	//
	// Conflict with `tls.certificate`.
	CertificatePath string `json:"certificate_path,omitempty"`

	// TLS server private key content.
	//
	// Either `tls.key` or `tls.key_path` is required.
	//
	// Conflict with `tls.key_path`.
	Key badoption.Listable[string] `json:"key,omitempty"`

	// TLS server private key path.
	//
	// Either `tls.key` or `tls.key_path` is required.
	//
	// Conflict with `tls.key`.
	KeyPath string `json:"key_path,omitempty"`

	// TLS CA certificate content, used to verify client certificates.
	//
	// One of `tls.client_certificate`, `tls.client_certificate_path`, or `tls.peer_fingerprint` is required when `tls.verify_client_certificate` is `require` or `optional`.
	//
	// Conflict with `tls.client_certificate_path`.
	ClientCertificate badoption.Listable[string] `json:"client_certificate,omitempty"`

	// TLS CA certificate path, used to verify client certificates.
	//
	// One of `tls.client_certificate`, `tls.client_certificate_path`, or `tls.peer_fingerprint` is required when `tls.verify_client_certificate` is `require` or `optional`.
	//
	// Conflict with `tls.client_certificate`.
	ClientCertificatePath string `json:"client_certificate_path,omitempty"`

	// OpenVPN client certificate policy, one of `require`, `optional` or `none`.
	//
	// If set to `optional`, a client certificate is verified when provided, but clients without a certificate are allowed.
	//
	// If set to `none`, client certificates are not requested.
	//
	// @Default: require
	VerifyClientCertificate string `json:"verify_client_certificate,omitempty" enum:"require,optional,none"`

	// Expected client certificate name. Disabled when empty.
	ClientName string `json:"client_name,omitempty"`

	// Certificate field matched by `tls.client_name`, one of `subject`, `name`, or `name-prefix`.
	//
	// `name` will be used by default when `tls.client_name` is configured.
	ClientNameType string `json:"client_name_type,omitempty" enum:"subject,name,name-prefix"`

	// Allowed SHA-256 fingerprints of client leaf certificates. Fingerprint-only verification can be used without a client CA.
	PeerFingerprint badoption.Listable[string] `json:"peer_fingerprint,omitempty"`

	// Path to a certificate revocation list used to reject revoked client certificates.
	CRLPath string `json:"crl_path,omitempty"`

	// Required client certificate key usage masks in OpenVPN `remote-cert-ku` format.
	RemoteCertificateKU badoption.Listable[string] `json:"remote_certificate_ku,omitempty"`

	// Required client certificate extended key usage. Conflict with an explicitly configured `tls.remote_certificate_tls`.
	RemoteCertificateEKU string `json:"remote_certificate_eku,omitempty"`

	// Client certificate purpose check, one of `server`, `client`, or `none`. `client` will be used by default.
	RemoteCertificateTLS string `json:"remote_certificate_tls,omitempty" enum:"server,client,none"`

	// Certificate profile, one of `insecure`, `legacy`, `preferred`, or `suiteb`.
	//
	// `insecure` accepts MD5- and SHA-1-signed certificate chains and smaller legacy
	// keys for compatibility with immutable peers. Use it only when the peer cannot
	// be upgraded. `legacy` accepts SHA-1 but rejects MD5 signatures; `preferred`
	// requires stronger signatures and keys.
	//
	// When `suiteb` is selected, the TLS 1.2 cipher list defaults to the Suite B ECDHE-ECDSA AES-GCM suites when `tls.cipher` is empty, and the TLS groups default to `secp256r1` and `secp384r1` when `tls.groups` is empty.
	//
	// @Default: legacy
	CertificateProfile string `json:"certificate_profile,omitempty" enum:"legacy,preferred,insecure,suiteb"`

	// Deprecated Netscape certificate type check, one of `server` or `client`.
	NSCertificateType string `json:"ns_certificate_type,omitempty" enum:"server,client"`

	// Minimum TLS version. `1.2` will be used by default.
	VersionMin string `json:"version_min,omitempty" enum:"1.0,1.1,1.2,1.3"`

	// Maximum TLS version. The maximum supported version will be used by default.
	VersionMax string `json:"version_max,omitempty" enum:"1.0,1.1,1.2,1.3"`

	// Colon-separated OpenSSL or IANA cipher suite names allowed for TLS 1.2 and earlier, such as `ECDHE-ECDSA-AES256-GCM-SHA384` or `TLS-ECDHE-ECDSA-WITH-AES-256-GCM-SHA384`. The order of the list is ignored.
	//
	// The default TLS cipher suites are used when empty.
	Cipher string `json:"cipher,omitempty"`

	// Colon-separated TLS key exchange groups. The order of the list is ignored.
	//
	// Supported groups are `X25519`, `X25519MLKEM768`, `secp256r1`, `secp384r1`, and `secp521r1`, including their common OpenSSL and NIST aliases. Names are case-insensitive.
	//
	// The default TLS groups are used when empty.
	Groups string `json:"groups,omitempty"`

	// OpenVPN control channel wrapping.
	//
	// Equivalent to OpenVPN `tls-auth`, `tls-crypt` and `tls-crypt-v2`.
	//
	// Disabled by default.
	ControlWrap *OpenVPNInboundControlWrapOptions `json:"control_wrap,omitempty"`
}

type OpenVPNControlWrapOptions struct {
	// Control channel wrapping type, one of `tls_auth`, `tls_crypt`, or `tls_crypt_v2`.
	Type string `json:"type,omitempty" enum:"tls_auth,tls_crypt,tls_crypt_v2"`

	// Control channel wrapping key content.
	//
	// Conflict with `tls.control_wrap.key_path`.
	Key badoption.Listable[string] `json:"key,omitempty"`

	// Control channel wrapping key path.
	//
	// Conflict with `tls.control_wrap.key`.
	KeyPath string `json:"key_path,omitempty"`

	// `tls-auth` key direction, one of `server` or `client`.
	//
	// Only available when `tls.control_wrap.type` is `tls_auth`. The key is used bidirectionally if empty.
	Direction string `json:"direction,omitempty" enum:"server,client"`
}

type OpenVPNInboundControlWrapOptions struct {
	// Control channel wrapping type, one of `tls_auth`, `tls_crypt` or `tls_crypt_v2`.
	//
	// For `tls_crypt_v2`, the key is the server key.
	//
	// @Required
	Type string `json:"type,omitempty" enum:"tls_auth,tls_crypt,tls_crypt_v2"`

	// Control channel wrapping key content.
	//
	// Either `tls.control_wrap.key` or `tls.control_wrap.key_path` is required.
	//
	// Conflict with `tls.control_wrap.key_path`.
	Key badoption.Listable[string] `json:"key,omitempty"`

	// Control channel wrapping key path.
	//
	// Either `tls.control_wrap.key` or `tls.control_wrap.key_path` is required.
	//
	// Conflict with `tls.control_wrap.key`.
	KeyPath string `json:"key_path,omitempty"`

	// OpenVPN `tls-auth` key direction, one of `server` or `client`.
	//
	// Only available when `tls.control_wrap.type` is `tls_auth`.
	//
	// `server` maps to OpenVPN key direction `0`, and `client` maps to `1`; by convention servers use `0` and clients use `1`.
	//
	// If empty, the key is used bidirectionally, matching an omitted `key-direction` on both peers.
	Direction string `json:"direction,omitempty" enum:"server,client"`

	// Require `tls-crypt-v2` clients over UDP to support stateless session cookies.
	//
	// Only available when `tls.control_wrap.type` is `tls_crypt_v2`. When disabled,
	// clients without cookie support are accepted using the upstream `allow-noncookie` behavior.
	//
	// Disabled by default.
	ForceCookie bool `json:"force_cookie,omitempty"`
}

type OpenVPNPushOptions struct {
	// Routes to push to clients.
	//
	// IPv4 and IPv6 prefixes can be mixed.
	Routes LegacyListable[netip.Prefix] `json:"routes,omitempty"`

	// DNS server addresses to push to clients.
	//
	// Uses legacy `dhcp-option DNS`/`DNS6`. A pushed modern DNS server group overrides these addresses on compatible clients.
	DNS badoption.Listable[netip.Addr] `json:"dns,omitempty"`

	// Modern OpenVPN DNS server groups to push. Each entry contains `priority`, `addresses`, optional `resolve_domains`, `dnssec`, `transport`, and `sni`.
	//
	// Addresses accept an IP address or `IP:port` (IPv6 ports use `[IPv6]:port`). `transport` is one of `plain`, `dot`, or `doh`; `dnssec` is one of `yes`, `optional`, or `no`. OpenVPN clients apply only the group with the lowest priority number.
	DNSServers []OpenVPNPushDNSServerOptions `json:"dns_servers,omitempty"`

	// Modern OpenVPN search domains to push.
	SearchDomains LegacyListable[string] `json:"search_domains,omitempty"`

	// Additional legacy `dhcp-option` values to push, without the `dhcp-option` prefix.
	DHCPOptions LegacyListable[string] `json:"dhcp_options,omitempty"`

	// Push `redirect-gateway` to clients, which routes client traffic through the VPN according to `push.redirect_gateway_flags`.
	//
	// When `push.redirect_gateway_flags` is empty, `def1` will be used by default.
	RedirectGateway bool `json:"redirect_gateway,omitempty"`

	// OpenVPN `redirect-gateway` flags to push to clients.
	//
	// Only available when `push.redirect_gateway` is enabled.
	//
	// `def1` is added by default.
	RedirectGatewayFlags LegacyListable[string] `json:"redirect_gateway_flags,omitempty"`

	// Push `block-outside-dns` to clients, which blocks DNS queries outside the VPN on Windows clients.
	BlockOutsideDNS bool `json:"block_outside_dns,omitempty"`

	// OpenVPN `ping` interval pushed to clients.
	//
	// After the interval passes without sending a packet, the client sends a data-channel ping to the server.
	//
	// The value must use whole seconds.
	//
	// Disabled by default.
	PingInterval badoption.Duration `json:"ping_interval,omitempty"`

	// OpenVPN `ping-restart` timeout pushed to clients.
	//
	// After the timeout passes without receiving a packet, the client reconnects to the server.
	//
	// The value must use whole seconds.
	//
	// Disabled by default.
	PingRestart badoption.Duration `json:"ping_restart,omitempty"`
}

type OpenVPNPushDNSServerOptions struct {
	Priority       int                    `json:"priority,omitempty"`
	Addresses      LegacyListable[string] `json:"addresses"`
	ResolveDomains LegacyListable[string] `json:"resolve_domains,omitempty"`
	DNSSEC         string                 `json:"dnssec,omitempty" enum:"yes,optional,no"`
	Transport      string                 `json:"transport,omitempty" enum:"plain,dot,doh"`
	SNI            string                 `json:"sni,omitempty"`
}

type OpenVPNDNSServerOptions struct {
	// The tag of the [OpenVPN Client Endpoint](/configuration/endpoint/openvpn-client).
	//
	// DNS queries are sent through the endpoint to resolvers pushed by the OpenVPN server.
	//
	// Both modern `dns` options and legacy `dhcp-option DNS`/`DNS6` and `DOMAIN-ROUTE` options are supported. `dnssec yes` is not supported.
	//
	// @Required
	Endpoint string `json:"endpoint,omitempty"`

	// Use pushed resolvers for queries that do not match a pushed `resolve-domains`, `DOMAIN-ROUTE`, or search-domain suffix.
	//
	// When disabled, unmatched queries return `NXDOMAIN`.
	AcceptDefaultResolvers bool `json:"accept_default_resolvers,omitempty"`

	// When enabled and pushed search domains are available, single-label queries (for example, `intranet`) are retried with each search domain until one resolves.
	//
	// If no search domain is available, the original single-label query follows normal default-resolver behavior.
	AcceptSearchDomain bool `json:"accept_search_domain,omitempty"`
}
