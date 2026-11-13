`http` outbound is a HTTP CONNECT proxy client.

### Structure

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "http",
  "tag": "http-out",
  "server": "127.0.0.1",
  "server_port": 1080,
  "username": "sekai",
  "password": "admin",
  "tls": {},
  "path": "",
  "headers": {},
  "version": 0,
  "disable_version_fallback": false

  ... // Dial Fields
  ... // HTTP2 Fields / QUIC Fields
}
```
<!-- AUTO-GENERATED-EXAMPLE-END -->

### Fields

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### server

==Required==

The server address.

#### server_port

==Required==

The server port.

#### username

Basic authorization username.

#### password

Basic authorization password.

#### tls

TLS configuration, see [TLS](/configuration/shared/tls/#outbound).

#### path

Path of HTTP request.

#### headers

Extra headers of HTTP request.

#### version

!!! question "Since sing-box 1.15.0"


HTTP version.

Available values: `1`, `2`, `3`.

`2` will be used by default, or `1` if `path` or the `Host` header is set.

`path` and the `Host` header are only available for `1`.

When `3`, [HTTP2 Fields](#http2-fields) are replaced by [QUIC Fields](#quic-fields).

#### disable_version_fallback

!!! question "Since sing-box 1.15.0"


Disable automatic fallback to lower HTTP version.

<!-- AUTO-GENERATED-CONFIG-END -->

### HTTP2 Fields

!!! question "Since sing-box 1.15.0"

When `version` is `2` (default).

See [HTTP2 Fields](/configuration/shared/http2/) for details.

### QUIC Fields

!!! question "Since sing-box 1.15.0"

When `version` is `3`.

See [QUIC Fields](/configuration/shared/quic/) for details.

### Dial Fields

<!-- AUTO-GENERATED-ORIGIN-START -->
See [Dial Fields](/configuration/shared/dial/) for details.
<!-- AUTO-GENERATED-ORIGIN-END -->
