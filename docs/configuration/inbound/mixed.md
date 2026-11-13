`mixed` inbound is a socks4, socks4a, socks5 and http server.

### Structure

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "mixed",
  "tag": "mixed-in",
  "users": [
    {
      "username": "admin",
      "password": "admin"
    }
  ],
  "set_system_proxy": false

  ... // Listen Fields
}
```
<!-- AUTO-GENERATED-EXAMPLE-END -->

### Listen Fields

<!-- AUTO-GENERATED-ORIGIN-START -->
See [Listen Fields](/configuration/shared/listen/) for details.
<!-- AUTO-GENERATED-ORIGIN-END -->

### Fields

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### users

SOCKS and HTTP users.

No authentication required if empty.

#### set_system_proxy

!!! quote ""

    Only supported on Linux, Android, Windows, and macOS.
!!! warning ""

    To work on Android and Apple platforms without privileges, use tun.platform.http_proxy instead.

Automatically set system proxy configuration when start and clean up when stop.

#### tls

TLS configuration, see [TLS](/configuration/shared/tls/#inbound).

<!-- AUTO-GENERATED-CONFIG-END -->
