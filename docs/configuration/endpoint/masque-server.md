# MASQUE Server

!!! question "Since sing-box 1.15.0"

`masque-server` endpoint is an IP proxying over HTTP ([RFC 9484](https://datatracker.ietf.org/doc/html/rfc9484), CONNECT-IP) server.

### Structure

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "masque-server",
  "tag": "masque-server",
  "system": false,
  "name": "",
  "mtu": 1280,
  "users": [
    {
      "username": "",
      "password": ""
    }
  ],
  "version": [],
  "tls": {},
  "path": "",
  "address": [],
  "advertise_routes": []

  ... // Listen Fields
  ... // HTTP2 Fields / QUIC Fields
  ... // UDP NAT Fields
}
```
<!-- AUTO-GENERATED-EXAMPLE-END -->

!!! note ""

    You can ignore the JSON Array [] tag when the content is only one item

### Listen Fields

See [Listen Fields](/configuration/shared/listen/) for details. `udp_timeout` is part of the [UDP NAT Fields](#udp-nat-fields) below.

### Fields

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

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

#### users

HTTP users, verified by the `Authorization` header.

No authentication required if empty.

#### version

List of HTTP versions to serve.

Available values: `1`, `2`, `3`.

All versions are used by default.

TLS is required for `3`.

#### tls

TLS configuration, see [TLS](/configuration/shared/tls/#inbound).

IP proxying must be operated over TLS or QUIC. Leave it disabled only when the server is placed behind an HTTP intermediary that terminates TLS.

#### path

URI template path of the IP proxying resource, may contain the `target` and `ipproto` variables.

`/.well-known/masque/ip/{target}/{ipproto}/` will be used by default.

#### address

==Required==

List of IP prefixes of the tunnel network, at most one for each IP version.

The address of the prefix is used by the server itself, other addresses in the prefix are assigned to clients.

#### advertise_routes

List of IP prefixes to advertise to clients, in addition to the tunnel network.

Traffic from clients to other destinations is rejected.

All addresses are advertised by default.

<!-- AUTO-GENERATED-CONFIG-END -->

### HTTP2 Fields

When `version` contains `2`.

See [HTTP2 Fields](/configuration/shared/http2/) for details.

### QUIC Fields

When `version` contains `3` (default), [HTTP2 Fields](#http2-fields) are replaced by QUIC Fields.

See [QUIC Fields](/configuration/shared/quic/) for details.

### UDP NAT Fields

See [UDP NAT Fields](/configuration/shared/udp-nat/) for details.
