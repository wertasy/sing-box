# OpenConnect Client

!!! question "Since sing-box 1.14.0"

!!! quote "Changes in sing-box 1.15.0"

    :material-plus: [on_demand](#on_demand)

==Client only==

### Structure

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "openconnect",
  "tag": "oc-client",
  "system": false,
  "name": "",
  "server": "vpn.example.com",
  "flavor": "anyconnect",
  "username": "",
  "password": "",
  "auth_group": "",
  "cookie": "",
  "token": {
    "mode": "",
    "secret": "",
    "secret_path": "",
    "pin": "",
    "password": "",
    "device_id": "",
    "counter": 0
  },
  "reported_os": "",
  "user_agent": "",
  "version": "",
  "local_hostname": "",
  "mobile": {
    "platform_version": "",
    "device_type": "",
    "device_unique_id": ""
  },
  "csd": {
    "wrapper_path": ""
  },
  "hip": {
    "wrapper_path": ""
  },
  "tncc": {
    "wrapper_path": "",
    "device_id": "",
    "user_agent": "",
    "machine_identification_enabled": false,
    "certificates": [
      {
        "certificate": [],
        "certificate_path": ""
      }
    ]
  },
  "fortinet_host_check": {
    "hostcheck": "",
    "check_virtual_desktop": ""
  },
  "no_udp": false,
  "dtls_local_port": 0,
  "compression_disabled": false,
  "compression_mode": "",
  "ipv6_disabled": false,
  "http_keepalive_disabled": false,
  "xml_post_disabled": false,
  "external_auth_disabled": false,
  "password_authentication_disabled": false,
  "tcp_keep_alive_enabled": false,
  "pfs": false,
  "mtu": 0,
  "base_mtu": 0,
  "dpd_interval": "",
  "reconnect_timeout": "",
  "trojan_interval": "",
  "queue_length": 0,
  "allow_insecure_crypto": false,
  "tls": {
    "insecure": false,
    "server_name": "",
    "peer_fingerprint": [],
    "system_trust_disabled": false,
    "certificate_authority": [],
    "certificate_authority_path": "",
    "client_certificate": [],
    "client_certificate_path": "",
    "client_key": [],
    "client_key_path": "",
    "client_key_password": "",
    "mca_certificate": [],
    "mca_certificate_path": "",
    "mca_key": [],
    "mca_key_path": "",
    "mca_key_password": ""
  },
  "form_entries": [
    {
      "form_id": "",
      "submission_key": "",
      "name": "",
      "value": "",
      "promote": false
    }
  ],
  "on_demand": false

  ... // Dial Fields
  ... // UDP NAT Fields
}
```
<!-- AUTO-GENERATED-EXAMPLE-END -->

!!! note ""

    You can ignore the JSON Array [] tag when the content is only one item.

### Fields

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### system

Use a system interface.

Requires privilege and cannot conflict with existing system interfaces.

If disabled, sing-box uses the internal network stack.

#### name

Custom interface name for the system interface.

An automatically generated `oc` interface name is used by default.

#### server

==Required==

OpenConnect VPN server HTTPS URL.

The `https://` scheme is added if omitted. URL user information, queries, and fragments are not supported.

#### flavor

OpenConnect protocol flavor, one of `anyconnect`, `gp`, `fortinet`, `f5`, `pulse`, or `nc`.

- Default: `anyconnect`

#### username

Username used to fill matching authentication form fields.

#### password

Password used to fill matching authentication form fields.

#### auth_group

Authentication group used to preselect a matching group, realm, domain, or gateway choice when supported by the selected flavor.

#### cookie

Existing authentication session used to connect without first prompting for credentials.

The accepted format depends on `flavor`:

- `anyconnect`: A `webvpn` value, or a semicolon-separated cookie list containing `webvpn`.
- `gp`: The complete authenticated query string returned by GlobalProtect authentication.
- `nc`: A `DSID` value, or a semicolon-separated cookie list containing `DSID`.
- `pulse`: The raw Pulse authentication cookie value.
- `f5`: An `MRHSession` value, or a semicolon-separated cookie list containing `MRHSession` and optionally `F5_ST`.
- `fortinet`: An `SVPNCOOKIE` value, or a semicolon-separated cookie list containing `SVPNCOOKIE`.

If the server rejects the supplied session, normal authentication is attempted.

#### token

Token configuration for automatically answering matching token fields or HTTP Bearer authentication.

One of `token.secret` or `token.secret_path` is required.

#### reported_os

Operating system identity reported to the VPN server when supported by the selected flavor.

For `anyconnect`, `gp`, and `pulse`, the supported values are `linux`, `linux-64`, `win`, `mac-intel`, `android`, and `apple-ios`.

The default is selected from the system platform: `win` on Windows, `mac-intel` on macOS, `android` on Android, `apple-ios` on iOS, and `linux-64` or `linux` on other 64-bit or 32-bit systems.

#### user_agent

User agent reported to the VPN server when supported by the selected flavor.

The default is flavor-specific. AnyConnect, Network Connect, Pulse, and F5 use `AnyConnect-compatible OpenConnect VPN Agent v9.21`; GlobalProtect uses `PAN GlobalProtect`; Fortinet uses `Mozilla/5.0 SV1`.

#### version

Client version reported separately from `user_agent` when supported by the selected flavor.

- Default: `v9.21`

#### local_hostname

Local hostname reported to the VPN server when supported by the selected flavor.

The system hostname is used by default, or `localhost` if it is unavailable.

#### mobile

AnyConnect mobile client identity. When configured, all three fields are required.

#### csd

AnyConnect CSD/host scan compliance options.

Built-in CSD handling is used by default when requested by the server.

#### hip

GlobalProtect HIP check and report options.

Built-in HIP reporting is used by default when requested by the server.

#### tncc

Network Connect TNCC compliance options.

Built-in TNCC handling is used by default when requested by the server.

#### fortinet_host_check

Fortinet hostcheck result override, submitted when requested by the server.

Disabled if `fortinet_host_check.hostcheck` is empty.

Some Fortinet servers only request hostcheck from recognized FortiClient user agents. Configure `user_agent` when required by the server policy.

#### no_udp

Disable the DTLS or ESP secondary data channel and use the TLS data channel only.

#### dtls_local_port

Local UDP port used by the direct DTLS or ESP secondary data channel.

An automatically selected ephemeral port is used by default.

#### compression_disabled

Disable AnyConnect compression negotiation.

By default, stateless `oc-lz4` and `lzs` compression is negotiated for CSTP and DTLS when supported by the server.

Conflict with `compression_mode` set to `all`.

#### compression_mode

AnyConnect compression mode, one of:

- `stateless`: Advertise stateless `oc-lz4` and `lzs` compression.
- `all`: Additionally advertise stateful `deflate` compression for CSTP.

- Default: `stateless`

#### ipv6_disabled

Disable requesting and using IPv6 tunnel configuration.

#### http_keepalive_disabled

Disable HTTP connection reuse during authentication and configuration requests.

#### xml_post_disabled

Disable AnyConnect XML POST authentication and start authentication with the legacy GET flow.

#### external_auth_disabled

Disable external browser authentication such as SSO and SAML for AnyConnect, GlobalProtect, and Fortinet.

#### password_authentication_disabled

Abort AnyConnect authentication if the server returns a non-success authentication form, matching OpenConnect `--no-passwd` behavior.

#### tcp_keep_alive_enabled

Enable TCP keep alive for direct VPN server connections.

Setting `tcp_keep_alive` or `tcp_keep_alive_interval` also enables it. When enabled without either duration, the operating system defaults are used.

Conflict with `disable_tcp_keep_alive`.

#### pfs

Require forward-secret TLS cipher suites for TLS 1.2 and earlier.

#### mtu

Preferred tunnel MTU.

The negotiated MTU is limited to this value.

Non-zero values below `576` are treated as `576`. The maximum value is `65535`.

#### base_mtu

Base path MTU used to calculate the AnyConnect, GlobalProtect, F5, and Fortinet tunnel MTU after outer IP, transport, and protocol overhead.

These flavors treat values below `1280` as `1280`. The maximum value is `65535`.

- Default: `1406`

#### dpd_interval

Override the Dead Peer Detection interval.

The server-provided or flavor-specific interval is used by default.

Positive values below `2s` are treated as `2s`. The value must not be negative.

#### reconnect_timeout

Maximum accumulated backoff time after failed reconnect attempts.

The value must not be negative.

- Default: `300s`

#### trojan_interval

Override the interval between GlobalProtect HIP reports or Network Connect TNCC checks.

The server-provided interval is used by default. GlobalProtect uses `1h` when the server does not provide one.

The value must not be negative.

#### queue_length

Inbound and outbound packet queue length between the VPN transport and the tunnel interface.

- Default: `32`

#### allow_insecure_crypto

Enable weak TLS and DTLS cipher suites and TLS 1.0 compatibility required by legacy VPN servers.

Disabled by default; TLS versions below 1.2 are otherwise rejected.

#### tls

OpenConnect TLS configuration.

#### form_entries

Authentication form field overrides.

An entry matches by `submission_key` when set, or by the combination of `form_id` and `name`. Later matching entries take precedence.

#### on_demand

!!! question "Since sing-box 1.15.0"

Allow the endpoint to be disconnected when necessary.

<!-- AUTO-GENERATED-CONFIG-END -->

### UDP NAT Fields

See [UDP NAT Fields](/configuration/shared/udp-nat/) for details.

### Dial Fields

<!-- AUTO-GENERATED-ORIGIN-START -->
See [Dial Fields](/configuration/shared/dial/) for details.
<!-- AUTO-GENERATED-ORIGIN-END -->

## Interactive authentication

Use `Tools` > `Endpoints` in the sing-box dashboard or any sing-box graphical client to authenticate and manage the endpoint.

## DNS

Configure an [OpenConnect DNS server](/configuration/dns/server/openconnect/) to use the pushed DNS settings.
