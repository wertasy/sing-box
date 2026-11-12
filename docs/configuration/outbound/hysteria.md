---
icon: material/new-box
---

!!! quote "Changes in sing-box 1.12.0"

    :material-plus: [server_ports](#server_ports)  
    :material-plus: [hop_interval](#hop_interval)

### Structure

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "hysteria",
  "tag": "hysteria-out",
  "server": "127.0.0.1",
  "server_port": 1080,
  "server_ports": [
    "2080:3000"
  ],
  "hop_interval": "",
  "up": "100 Mbps",
  "up_mbps": 100,
  "down": "100 Mbps",
  "down_mbps": 100,
  "obfs": "fuck me till the daylight",
  "auth": "",
  "auth_str": "password",
  "network": "",
  "tls": {},

  ... // Dial Fields
  ... // QUIC Fields

  // Deprecated
  "recv_window_conn": 0,
  "recv_window": 0,
  "disable_mtu_discovery": false
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

#### server_ports

!!! question "Since sing-box 1.12.0"

Server port range list.

Conflicts with `server_port`.

#### hop_interval

!!! question "Since sing-box 1.12.0"

Port hopping interval.

- Default: `30s`

#### up, down

==Required==

Format: `[Integer] [Unit]` e.g. `100 Mbps, 640 KBps, 2 Gbps`

Supported units (case sensitive, b = bits, B = bytes, 8b=1B):

    bps (bits per second)
    Bps (bytes per second)
    Kbps (kilobits per second)
    KBps (kilobytes per second)
    Mbps (megabits per second)
    MBps (megabytes per second)
    Gbps (gigabits per second)
    GBps (gigabytes per second)
    Tbps (terabits per second)
    TBps (terabytes per second)

#### up_mbps, down_mbps

==Required==

`up, down` in Mbps.

#### obfs

Obfuscated password.

#### auth

Authentication password, in base64.

#### auth_str

Authentication password.

#### network

Enabled network

One of `tcp` `udp`.

Both is enabled by default.

#### tls

==Required==

TLS configuration, see [TLS](/configuration/shared/tls/#outbound).

#### recv_window_conn


- Deprecated.

#### recv_window


- Deprecated.

#### disable_mtu_discovery


- Deprecated.

<!-- AUTO-GENERATED-CONFIG-END -->

### QUIC Fields

See [QUIC Fields](/configuration/shared/quic/) for details.

### Dial Fields

<!-- AUTO-GENERATED-ORIGIN-START -->
See [Dial Fields](/configuration/shared/dial/) for details.
<!-- AUTO-GENERATED-ORIGIN-END -->

### Deprecated Fields

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### recv_window_conn

!!! failure "Deprecated in sing-box 1.14.0"

    Use QUIC fields `stream_receive_window` instead.

#### recv_window

!!! failure "Deprecated in sing-box 1.14.0"

    Use QUIC fields `connection_receive_window` instead.

#### disable_mtu_discovery

!!! failure "Deprecated in sing-box 1.14.0"

    Use QUIC fields `disable_path_mtu_discovery` instead.

<!-- AUTO-GENERATED-CONFIG-END -->
