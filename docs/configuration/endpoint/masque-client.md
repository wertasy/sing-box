# MASQUE Client

!!! question "Since sing-box 1.15.0"

`masque-client` endpoint is an IP proxying over HTTP ([RFC 9484](https://datatracker.ietf.org/doc/html/rfc9484), CONNECT-IP) client.

### Structure

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "masque-client",
  "tag": "masque-client",
  "server": "127.0.0.1",
  "server_port": 443,
  "system": false,
  "name": "",
  "mtu": 1280,
  "username": "",
  "password": "",
  "tls": {},
  "path": "",
  "headers": {},
  "version": 0,
  "disable_version_fallback": false,
  "advertise_routes": [],
  "on_demand": false

  ... // Dial Fields
  ... // HTTP2 Fields / QUIC Fields
  ... // UDP NAT Fields
}
```
<!-- AUTO-GENERATED-EXAMPLE-END -->

!!! note ""

    You can ignore the JSON Array [] tag when the content is only one item

### Fields

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### server

==Required==

The server address.

#### server_port

==Required==

The server port.

#### system

Use system interface.

Requires privilege and cannot conflict with existing system interfaces.

If disabled, sing-box uses the internal network stack.

#### name

Custom interface name for system interface.

An automatically generated `masque` interface name will be used by default.

#### mtu

Tunnel MTU.

`1280` will be used by default.

#### username

Basic authorization username.

#### password

Basic authorization password.

#### tls

TLS configuration, see [TLS](/configuration/shared/tls/#outbound).

Required for HTTP/3.

#### path

URI template path of the IP proxying resource, may contain the `target` and `ipproto` variables.

`/.well-known/masque/ip/{target}/{ipproto}/` will be used by default.

#### headers

Extra headers of HTTP request.

#### version

HTTP version.

Available values: `1`, `2`, `3`.

When `2`, [QUIC Fields](#quic-fields) are replaced by [HTTP2 Fields](#http2-fields).

`3` will be used by default.

#### disable_version_fallback

Disable automatic fallback to lower HTTP version.

#### advertise_routes

List of IP prefixes to advertise to the server.

The server will route traffic for these prefixes into this endpoint, where it is handled as inbound traffic.

#### on_demand

Allow the endpoint to be disconnected when necessary.

<!-- AUTO-GENERATED-CONFIG-END -->

### HTTP2 Fields

When `version` is `2`.

See [HTTP2 Fields](/configuration/shared/http2/) for details.

### QUIC Fields

When `version` is `3` (default).

See [QUIC Fields](/configuration/shared/quic/) for details.

### UDP NAT Fields

See [UDP NAT Fields](/configuration/shared/udp-nat/) for details.

### Dial Fields

<!-- AUTO-GENERATED-ORIGIN-START -->
See [Dial Fields](/configuration/shared/dial/) for details.
<!-- AUTO-GENERATED-ORIGIN-END -->
