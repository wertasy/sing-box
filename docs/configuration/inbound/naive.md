!!! quote "Changes in sing-box 1.13.0"

    :material-plus: [quic_congestion_control](#quic_congestion_control)

### Structure

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "naive",
  "tag": "naive-in",
  "users": [
    {
      "username": "sekai",
      "password": "password"
    }
  ],
  "network": "udp",
  "quic_congestion_control": "",
  "tls": {}

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

#### network

Listen network, one of `tcp` `udp`.

Both if empty.

#### users

==Required==

Naive users.

#### quic_congestion_control

!!! question "Since sing-box 1.13.0"


QUIC congestion control algorithm.

| Algorithm      | Description                     |
|----------------|---------------------------------|
| `bbr`          | BBR                             |
| `cubic`        | CUBIC                           |
| `reno`         | New Reno                        |

`bbr` will be used by default.

#### tls

TLS configuration, see [TLS](/configuration/shared/tls/#inbound).

<!-- AUTO-GENERATED-CONFIG-END -->
