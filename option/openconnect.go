package option

import "github.com/sagernet/sing/common/json/badoption"

// @Tag oc-client
type OpenConnectEndpointOptions struct {
	DialerOptions

	// Use a system interface.
	//
	// Requires privilege and cannot conflict with existing system interfaces.
	//
	// If disabled, sing-box uses the internal network stack.
	System bool `json:"system,omitempty"`

	// Custom interface name for the system interface.
	//
	// An automatically generated `oc` interface name will be used by default.
	Name string `json:"name,omitempty"`

	// @Omit
	UDPTimeout badoption.Duration `json:"udp_timeout,omitempty"`

	// @Omit
	UDPMapping UDPNATBehavior `json:"udp_mapping,omitempty"`

	// @Omit
	UDPFiltering UDPNATBehavior `json:"udp_filtering,omitempty"`

	// @Omit
	UDPNATMax uint32 `json:"udp_nat_max,omitempty"`

	// OpenConnect VPN server HTTPS URL.
	//
	// The `https://` scheme is added if omitted. URL user information, queries, and fragments are not supported.
	//
	// @Example: "vpn.example.com"
	//
	// @Required
	Server string `json:"server"`

	// OpenConnect protocol flavor, one of `anyconnect`, `gp`, `fortinet`, `f5`, `pulse`, or `nc`.
	//
	// @Default: anyconnect
	//
	// @Example: "anyconnect"
	Flavor string `json:"flavor,omitempty" enum:"anyconnect,gp,fortinet,f5,pulse,nc"`

	// Username used to fill matching authentication form fields.
	Username string `json:"username,omitempty"`

	// Password used to fill matching authentication form fields.
	Password string `json:"password,omitempty"`

	// Authentication group used to preselect a matching group, realm, domain, or gateway choice when supported by the selected flavor.
	AuthGroup string `json:"auth_group,omitempty"`

	// Existing authentication session used to connect without first prompting for credentials.
	//
	// The accepted format depends on `flavor`:
	//
	// - `anyconnect`: A `webvpn` value, or a semicolon-separated cookie list containing `webvpn`.
	// - `gp`: The complete authenticated query string returned by GlobalProtect authentication.
	// - `nc`: A `DSID` value, or a semicolon-separated cookie list containing `DSID`.
	// - `pulse`: The raw Pulse authentication cookie value.
	// - `f5`: An `MRHSession` value, or a semicolon-separated cookie list containing `MRHSession` and optionally `F5_ST`.
	// - `fortinet`: An `SVPNCOOKIE` value, or a semicolon-separated cookie list containing `SVPNCOOKIE`.
	//
	// If the server rejects the supplied session, normal authentication is attempted.
	Cookie string `json:"cookie,omitempty"`

	// Token configuration for automatically answering matching token fields or HTTP Bearer authentication.
	//
	// One of `token.secret` or `token.secret_path` is required.
	//
	// @Example:
	//    {
	//      "mode": "",
	//      "secret": "",
	//      "secret_path": "",
	//      "pin": "",
	//      "password": "",
	//      "device_id": "",
	//      "counter": 0
	//    }
	Token *OpenConnectTokenOptions `json:"token,omitempty"`

	// Operating system identity reported to the VPN server when supported by the selected flavor.
	//
	// For `anyconnect`, `gp`, and `pulse`, the supported values are `linux`, `linux-64`, `win`, `mac-intel`, `android`, and `apple-ios`.
	//
	// The default is selected from the system platform: `win` on Windows, `mac-intel` on macOS, `android` on Android, `apple-ios` on iOS, and `linux-64` or `linux` on other 64-bit or 32-bit systems.
	ReportedOS string `json:"reported_os,omitempty"`

	// User agent reported to the VPN server when supported by the selected flavor.
	//
	// The default is flavor-specific. AnyConnect, Network Connect, Pulse, and F5 use `AnyConnect-compatible OpenConnect VPN Agent v9.21`; GlobalProtect uses `PAN GlobalProtect`; Fortinet uses `Mozilla/5.0 SV1`.
	UserAgent string `json:"user_agent,omitempty"`

	// Client version reported separately from `user_agent` when supported by the selected flavor.
	//
	// @Default: v9.21
	Version string `json:"version,omitempty"`

	// Local hostname reported to the VPN server when supported by the selected flavor.
	//
	// The system hostname will be used by default, or `localhost` if it is unavailable.
	LocalHostname string `json:"local_hostname,omitempty"`

	// AnyConnect mobile client identity. When configured, all three fields are required.
	//
	// @Example:
	//    {
	//      "platform_version": "",
	//      "device_type": "",
	//      "device_unique_id": ""
	//    }
	Mobile *OpenConnectMobileOptions `json:"mobile,omitempty"`

	// AnyConnect CSD/host scan compliance options.
	//
	// Built-in CSD handling will be used by default when requested by the server.
	//
	// @Example:
	//    {
	//      "wrapper_path": ""
	//    }
	CSD *OpenConnectCSDOptions `json:"csd,omitempty"`

	// GlobalProtect HIP check and report options.
	//
	// Built-in HIP reporting will be used by default when requested by the server.
	//
	// @Example:
	//    {
	//      "wrapper_path": ""
	//    }
	HIP *OpenConnectHIPOptions `json:"hip,omitempty"`

	// Network Connect TNCC compliance options.
	//
	// Built-in TNCC handling will be used by default when requested by the server.
	//
	// @Example:
	//    {
	//      "wrapper_path": "",
	//      "device_id": "",
	//      "user_agent": "",
	//      "machine_identification_enabled": false,
	//      "certificates": [
	//        {
	//          "certificate": [],
	//          "certificate_path": ""
	//        }
	//      ]
	//    }
	TNCC *OpenConnectTNCCOptions `json:"tncc,omitempty"`

	// Fortinet hostcheck result override, submitted when requested by the server.
	//
	// Disabled if `fortinet_host_check.hostcheck` is empty.
	//
	// Some Fortinet servers only request hostcheck from recognized FortiClient user agents. Configure `user_agent` when required by the server policy.
	//
	// @Example:
	//    {
	//      "hostcheck": "",
	//      "check_virtual_desktop": ""
	//    }
	FortinetHostCheck *OpenConnectFortinetHostCheckOptions `json:"fortinet_host_check,omitempty"`

	// Disable the DTLS or ESP secondary data channel and use the TLS data channel only.
	NoUDP bool `json:"no_udp,omitempty"`

	// Local UDP port used by the direct DTLS or ESP secondary data channel.
	//
	// An automatically selected ephemeral port will be used by default.
	DTLSLocalPort uint16 `json:"dtls_local_port,omitempty"`

	// Disable AnyConnect compression negotiation.
	//
	// By default, stateless `oc-lz4` and `lzs` compression is negotiated for CSTP and DTLS when supported by the server.
	//
	// Conflict with `compression_mode` set to `all`.
	CompressionDisabled bool `json:"compression_disabled,omitempty"`

	// AnyConnect compression mode, one of:
	//
	// - `stateless`: Advertise stateless `oc-lz4` and `lzs` compression.
	// - `all`: Additionally advertise stateful `deflate` compression for CSTP.
	//
	// @Default: stateless
	CompressionMode string `json:"compression_mode,omitempty" enum:"stateless,all"`

	// Disable requesting and using IPv6 tunnel configuration.
	IPv6Disabled bool `json:"ipv6_disabled,omitempty"`

	// Disable HTTP connection reuse during authentication and configuration requests.
	HTTPKeepAliveDisabled bool `json:"http_keepalive_disabled,omitempty"`

	// Disable AnyConnect XML POST authentication and start authentication with the legacy GET flow.
	XMLPostDisabled bool `json:"xml_post_disabled,omitempty"`

	// Disable external browser authentication such as SSO and SAML for AnyConnect, GlobalProtect, and Fortinet.
	ExternalAuthDisabled bool `json:"external_auth_disabled,omitempty"`

	// Abort AnyConnect authentication if the server returns a non-success authentication form, matching OpenConnect `--no-passwd` behavior.
	PasswordAuthenticationDisabled bool `json:"password_authentication_disabled,omitempty"`

	// Enable TCP keep alive for direct VPN server connections.
	//
	// Setting `tcp_keep_alive` or `tcp_keep_alive_interval` also enables it. When enabled without either duration, the operating system defaults are used.
	//
	// Conflict with `disable_tcp_keep_alive`.
	TCPKeepAliveEnabled bool `json:"tcp_keep_alive_enabled,omitempty"`

	// Require forward-secret TLS cipher suites for TLS 1.2 and earlier.
	PFS bool `json:"pfs,omitempty"`

	// Preferred tunnel MTU.
	//
	// The negotiated MTU is limited to this value.
	//
	// Non-zero values below `576` are treated as `576`. The maximum value is `65535`.
	MTU uint32 `json:"mtu,omitempty"`

	// Base path MTU used to calculate the AnyConnect, GlobalProtect, F5, and Fortinet tunnel MTU after outer IP, transport, and protocol overhead.
	//
	// These flavors treat values below `1280` as `1280`. The maximum value is `65535`.
	//
	// @Default: 1406
	BaseMTU uint32 `json:"base_mtu,omitempty"`

	// Override the Dead Peer Detection interval.
	//
	// The server-provided or flavor-specific interval will be used by default.
	//
	// Positive values below `2s` are treated as `2s`. The value must not be negative.
	DPDInterval badoption.Duration `json:"dpd_interval,omitempty"`

	// Maximum accumulated backoff time after failed reconnect attempts.
	//
	// The value must not be negative.
	//
	// @Default: 300s
	ReconnectTimeout badoption.Duration `json:"reconnect_timeout,omitempty"`

	// Override the interval between GlobalProtect HIP reports or Network Connect TNCC checks.
	//
	// The server-provided interval will be used by default. GlobalProtect uses `1h` when the server does not provide one.
	//
	// The value must not be negative.
	TrojanInterval badoption.Duration `json:"trojan_interval,omitempty"`

	// Inbound and outbound packet queue length between the VPN transport and the tunnel interface.
	//
	// @Default: 32
	QueueLength uint32 `json:"queue_length,omitempty"`

	// Enable weak TLS and DTLS cipher suites and TLS 1.0 compatibility required by legacy VPN servers.
	//
	// Disabled by default; TLS versions below 1.2 are otherwise rejected.
	AllowInsecureCrypto bool `json:"allow_insecure_crypto,omitempty"`

	// OpenConnect TLS configuration.
	//
	// @Example:
	//    {
	//      "insecure": false,
	//      "server_name": "",
	//      "peer_fingerprint": [],
	//      "system_trust_disabled": false,
	//      "certificate_authority": [],
	//      "certificate_authority_path": "",
	//      "client_certificate": [],
	//      "client_certificate_path": "",
	//      "client_key": [],
	//      "client_key_path": "",
	//      "client_key_password": "",
	//      "mca_certificate": [],
	//      "mca_certificate_path": "",
	//      "mca_key": [],
	//      "mca_key_path": "",
	//      "mca_key_password": ""
	//    }
	TLS OpenConnectTLSOptions `json:"tls,omitempty"`

	// Authentication form field overrides.
	//
	// An entry matches by `submission_key` when set, or by the combination of `form_id` and `name`. Later matching entries take precedence.
	//
	// @Example:
	//    [
	//      {
	//        "form_id": "",
	//        "submission_key": "",
	//        "name": "",
	//        "value": "",
	//        "promote": false
	//      }
	//    ]
	FormEntries []OpenConnectFormEntryOptions `json:"form_entries,omitempty"`

	// Allow the endpoint to be disconnected when necessary.
	//
	// @Since 1.15.0
	OnDemand bool `json:"on_demand,omitempty"`
}

type OpenConnectTokenOptions struct {
	// Token mode, one of:
	//
	// - `totp`: Time-based One-Time Password.
	// - `hotp`: HMAC-based One-Time Password.
	// - `stoken`: RSA SecurID software token.
	// - `oidc`: OIDC access token used for HTTP Bearer authentication.
	//
	// @Required
	Mode string `json:"mode,omitempty" enum:"totp,hotp,stoken,oidc"`

	// Software token secret.
	//
	// For `totp` and `hotp`, this can be a Base32 secret, a `base32:`-prefixed secret, or an `otpauth://` URI of the matching type.
	//
	// For `stoken`, this is the encoded RSA SecurID CTF token content.
	//
	// For `oidc`, this is the access token value.
	//
	// Conflict with `token.secret_path`.
	Secret string `json:"secret,omitempty"`

	// Path to the software token secret or OIDC access token.
	//
	// Conflict with `token.secret`.
	SecretPath string `json:"secret_path,omitempty"`

	// RSA SecurID PIN for `stoken` mode.
	PIN string `json:"pin,omitempty"`

	// Password for decrypting a password-protected RSA SecurID token in `stoken` mode.
	Password string `json:"password,omitempty"`

	// Device ID for decrypting a device-bound RSA SecurID token in `stoken` mode.
	DeviceID string `json:"device_id,omitempty"`

	// Initial counter for `hotp` mode.
	//
	// If zero, the counter from an `otpauth://` URI is used when present; otherwise the counter starts at zero.
	Counter uint64 `json:"counter,omitempty"`
}

type OpenConnectMobileOptions struct {
	// Mobile operating system version reported to the AnyConnect server.
	PlatformVersion string `json:"platform_version"`

	// Mobile device model or type reported to the AnyConnect server.
	DeviceType string `json:"device_type"`

	// Mobile device identifier reported to the AnyConnect server.
	DeviceUniqueID string `json:"device_unique_id"`
}

type OpenConnectCSDOptions struct {
	// Path to an external AnyConnect CSD wrapper executable.
	//
	// Built-in CSD handling is used if empty.
	WrapperPath string `json:"wrapper_path,omitempty"`
}

type OpenConnectHIPOptions struct {
	// Path to an external GlobalProtect HIP report wrapper executable.
	//
	// Built-in HIP reporting is used if empty.
	WrapperPath string `json:"wrapper_path,omitempty"`
}

type OpenConnectTNCCOptions struct {
	// Path to an external Network Connect TNCC wrapper executable.
	//
	// Built-in TNCC handling is used if empty.
	//
	// Conflict with `tncc.device_id`, `tncc.user_agent`, `tncc.machine_identification_enabled`, and `tncc.certificates`.
	WrapperPath string `json:"wrapper_path,omitempty"`

	// Device ID reported by the built-in TNCC handler.
	//
	// Conflict with `tncc.wrapper_path`.
	DeviceID string `json:"device_id,omitempty"`

	// User agent used by the built-in TNCC handler.
	//
	// `Neoteris HC Http` will be used by default.
	//
	// Conflict with `tncc.wrapper_path`.
	UserAgent string `json:"user_agent,omitempty"`

	// Enable built-in TNCC machine identification, including the platform, hostname, and observed MAC addresses.
	//
	// Conflict with `tncc.wrapper_path`.
	MachineIdentificationEnabled bool `json:"machine_identification_enabled,omitempty"`

	// Machine certificates used by the built-in TNCC handler to answer certificate requests.
	//
	// Requires `tncc.machine_identification_enabled`.
	//
	// Conflict with `tncc.wrapper_path`.
	Certificates []OpenConnectTNCCCertificateOptions `json:"certificates,omitempty"`
}

type OpenConnectFortinetHostCheckOptions struct {
	// Fortinet hostcheck result string.
	//
	// The conventional format is `<security-status>,<os-version>`, for example `0100,10.0.19042`. `security-status` contains four `0` or `1` characters representing, in order, third-party firewall, third-party antivirus, FortiClient firewall, and FortiClient antivirus.
	HostCheck string `json:"hostcheck,omitempty"`

	// Fortinet virtual desktop check result string.
	//
	// FortiClient conventionally sends colon-separated MAC addresses joined by `|`, for example `74:78:27:4d:81:93|84:1b:77:3a:95:84`.
	CheckVirtualDesktop string `json:"check_virtual_desktop,omitempty"`
}

type OpenConnectTNCCCertificateOptions struct {
	// TNCC machine certificate content in PEM format.
	//
	// Conflict with `tncc.certificates.certificate_path`.
	Certificate badoption.Listable[string] `json:"certificate,omitempty"`

	// TNCC machine certificate path in PEM format.
	//
	// Conflict with `tncc.certificates.certificate`.
	CertificatePath string `json:"certificate_path,omitempty"`
}

type OpenConnectTLSOptions struct {
	// Disable verification of the VPN server certificate and hostname.
	//
	// Disabled by default. Prefer `tls.certificate_authority` or `tls.peer_fingerprint` when possible.
	Insecure bool `json:"insecure,omitempty"`

	// Server name used for TLS SNI and certificate hostname verification.
	//
	// The hostname from `server` will be used by default.
	ServerName string `json:"server_name,omitempty"`

	// Allowed server certificate fingerprints. A single string or a list can be specified.
	//
	// Supported formats:
	//
	// - An unprefixed SHA-1 certificate fingerprint compatible with OpenConnect `--servercert`.
	// - `sha1:<hex>`: SHA-1 SPKI fingerprint.
	// - `sha256:<hex>`: SHA-256 SPKI fingerprint.
	// - `pin-sha256:<base64>`: Base64-encoded SHA-256 SPKI pin.
	//
	// The encoded fingerprint in every format can be abbreviated to a prefix of at least four characters. When configured, the peer certificate must match one of these fingerprints; a match can authorize a certificate that is not otherwise trusted.
	PeerFingerprint badoption.Listable[string] `json:"peer_fingerprint,omitempty"`

	// Disable the system CA certificate pool.
	//
	// Use `tls.certificate_authority` or `tls.peer_fingerprint` to establish trust when enabled.
	SystemTrustDisabled bool `json:"system_trust_disabled,omitempty"`

	// Additional trusted CA certificate content in PEM format.
	//
	// The certificates are added to the system certificate pool.
	//
	// Conflict with `tls.certificate_authority_path`.
	CertificateAuthority badoption.Listable[string] `json:"certificate_authority,omitempty"`

	// Path to additional trusted CA certificates in PEM format.
	//
	// The certificates are added to the system certificate pool.
	//
	// Conflict with `tls.certificate_authority`.
	CertificateAuthorityPath string `json:"certificate_authority_path,omitempty"`

	// Client certificate chain content in PEM format.
	//
	// Conflict with `tls.client_certificate_path`.
	ClientCertificate badoption.Listable[string] `json:"client_certificate,omitempty"`

	// Client certificate chain path in PEM format.
	//
	// Conflict with `tls.client_certificate`.
	ClientCertificatePath string `json:"client_certificate_path,omitempty"`

	// Client private key content in PEM format.
	//
	// Conflict with `tls.client_key_path`.
	ClientKey badoption.Listable[string] `json:"client_key,omitempty"`

	// Client private key path in PEM format.
	//
	// Conflict with `tls.client_key`.
	//
	// The client certificate and key must both be set or both be empty.
	ClientKeyPath string `json:"client_key_path,omitempty"`

	// Password for the encrypted client private key.
	ClientKeyPassword string `json:"client_key_password,omitempty"`

	// AnyConnect multiple-certificate authentication (MCA) certificate chain content in PEM format.
	//
	// Conflict with `tls.mca_certificate_path`.
	MCACertificate badoption.Listable[string] `json:"mca_certificate,omitempty"`

	// AnyConnect multiple-certificate authentication (MCA) certificate chain path in PEM format.
	//
	// Conflict with `tls.mca_certificate`.
	MCACertificatePath string `json:"mca_certificate_path,omitempty"`

	// AnyConnect multiple-certificate authentication (MCA) private key content in PEM format.
	//
	// Conflict with `tls.mca_key_path`.
	MCAKey badoption.Listable[string] `json:"mca_key,omitempty"`

	// AnyConnect multiple-certificate authentication (MCA) private key path in PEM format.
	//
	// Conflict with `tls.mca_key`.
	//
	// The MCA certificate and key must both be set or both be empty.
	MCAKeyPath string `json:"mca_key_path,omitempty"`

	// Password for the encrypted MCA private key.
	MCAKeyPassword string `json:"mca_key_password,omitempty"`
}

type OpenConnectFormEntryOptions struct {
	// Authentication form identifier used with `form_entries.name` when `form_entries.submission_key` is empty.
	FormID string `json:"form_id,omitempty"`

	// Authentication field submission key.
	//
	// Either `form_entries.submission_key` or both `form_entries.form_id` and `form_entries.name` are required.
	SubmissionKey string `json:"submission_key,omitempty"`

	// Authentication field name used with `form_entries.form_id` when `form_entries.submission_key` is empty.
	Name string `json:"name,omitempty"`

	// Value supplied automatically for the matching authentication field.
	//
	// Conflict with `form_entries.promote`.
	Value string `json:"value,omitempty"`

	// Ask for the matching authentication field interactively instead of supplying an automatic value.
	//
	// Conflict with `form_entries.value`.
	Promote bool `json:"promote,omitempty"`
}

type OpenConnectDNSServerOptions struct {
	// The tag of the [OpenConnect Endpoint](/configuration/endpoint/openconnect).
	//
	// DNS queries are sent to the resolvers pushed by the VPN server through the OpenConnect endpoint.
	//
	// @Required
	Endpoint string `json:"endpoint,omitempty"`

	// Accept the general resolvers pushed by the VPN server for unmatched queries.
	//
	// When enabled, the general resolvers are used as the default only if the server requests all DNS through the tunnel, or if it does not provide split-DNS rules or suffixes. Otherwise, unmatched queries return `NXDOMAIN`.
	AcceptDefaultResolvers bool `json:"accept_default_resolvers,omitempty"`

	// When enabled and pushed search domains are available, single-label queries (for example, `intranet`) are retried with each search domain until one resolves.
	//
	// If every search-domain expansion returns `NXDOMAIN`, the original unqualified name follows normal default-resolver behavior.
	AcceptSearchDomain bool `json:"accept_search_domain,omitempty"`
}
