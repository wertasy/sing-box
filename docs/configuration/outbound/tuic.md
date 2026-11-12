### Structure

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "tuic",
  "tag": "tuic-out",
  "server": "127.0.0.1",
  "server_port": 1080,
  "uuid": "2DD61D93-75D8-4DA4-AC0E-6AECE7EAC365",
  "password": "hello",
  "congestion_control": "cubic",
  "udp_relay_mode": "native",
  "udp_over_stream": false,
  "zero_rtt_handshake": false,
  "heartbeat": "10s",
  "network": "tcp",
  "tls": {}

  ... // Dial Fields
  ... // QUIC Fields
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

#### uuid

==Required==

TUIC user uuid

#### password

TUIC user password

#### congestion_control

QUIC congestion control algorithm

One of: `cubic`, `new_reno`, `bbr`

- Default: `cubic`

#### udp_relay_mode

UDP packet relay mode

| Mode   | Description                                                              |
|:-------|:-------------------------------------------------------------------------|
| native | native UDP characteristics                                               |
| quic   | lossless UDP relay using QUIC streams, additional overhead is introduced |

Conflict with `udp_over_stream`.

- Default: `native`

#### udp_over_stream

This is the TUIC port of the [UDP over TCP protocol](/configuration/shared/udp-over-tcp/), designed to provide a QUIC
stream based UDP relay mode that TUIC does not provide. Since it is an add-on protocol, you will need to use sing-box or
another program compatible with the protocol as a server.

This mode has no positive effect in a proper UDP proxy scenario and should only be applied to relay streaming UDP
traffic (basically QUIC streams).

Conflict with `udp_relay_mode`.

#### network

Enabled network

One of `tcp` `udp`.

Both is enabled by default.

#### tls

==Required==

TLS configuration, see [TLS](/configuration/shared/tls/#outbound).

#### zero_rtt_handshake

!!! warning ""

    Disabling this feature is highly recommended, as it is vulnerable to replay attacks.
    See [Attack of the clones](https://blog.cloudflare.com/even-faster-connection-establishment-with-quic-0-rtt-resumption/#attack-of-the-clones)

Enable 0-RTT QUIC connection handshake on the client.
This has little performance impact, because the protocol is fully multiplexed.

#### heartbeat

Heartbeat interval for keeping the connection alive

<!-- AUTO-GENERATED-CONFIG-END -->

### QUIC Fields

See [QUIC Fields](/configuration/shared/quic/) for details.

### Dial Fields

<!-- AUTO-GENERATED-ORIGIN-START -->
See [Dial Fields](/configuration/shared/dial/) for details.
<!-- AUTO-GENERATED-ORIGIN-END -->
