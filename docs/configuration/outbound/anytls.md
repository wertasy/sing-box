---
icon: material/new-box
---

!!! question "Since sing-box 1.12.0"

### Structure

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "anytls",
  "tag": "anytls-out",
  "server": "127.0.0.1",
  "server_port": 1080,
  "tls": {},
  "password": "8JCsPssfgS8tiRwiMlhARg==",
  "idle_session_check_interval": "30s",
  "idle_session_timeout": "30s",
  "min_idle_session": 5,
  "client_metadata": ""

  ... // Dial Fields
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

#### tls

==Required==

TLS configuration, see [TLS](/configuration/shared/tls/#outbound).

#### password

==Required==

The AnyTLS password.

#### idle_session_check_interval

Interval checking for idle sessions. Default: 30s.

#### idle_session_timeout

In the check, close sessions that have been idle for longer than this. Default: 30s.

#### min_idle_session

In the check, at least the first `n` idle sessions are kept open. Default value: `n`=0

#### client_metadata

!!! question "Since sing-box 1.13.16"

    Check [AnyTLS client metadata](/manual/misc/anytls-client-metadata/).

<!-- AUTO-GENERATED-CONFIG-END -->

### Dial Fields

<!-- AUTO-GENERATED-ORIGIN-START -->
See [Dial Fields](/configuration/shared/dial/) for details.
<!-- AUTO-GENERATED-ORIGIN-END -->
