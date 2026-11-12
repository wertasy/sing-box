### Structure

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "tuic",
  "tag": "tuic-in",
  "users": [
    {
      "name": "sekai",
      "uuid": "059032A9-7D40-4A96-9BB1-36823D848068",
      "password": "hello"
    }
  ],
  "congestion_control": "cubic",
  "auth_timeout": "3s",
  "zero_rtt_handshake": false,
  "heartbeat": "10s",
  "tls": {}

  ... // Listen Fields
  ... // QUIC Fields
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

TUIC users

#### users.uuid

==Required==

TUIC user uuid

#### users.password

TUIC user password

#### congestion_control

QUIC congestion control algorithm

One of: `cubic`, `new_reno`, `bbr`

- Default: `cubic`

#### auth_timeout

How long the server should wait for the client to send the authentication command

- Default: `3s`

#### zero_rtt_handshake

!!! warning ""
    Disabling this is highly recommended, as it is vulnerable to replay attacks.
    See [Attack of the clones](https://blog.cloudflare.com/even-faster-connection-establishment-with-quic-0-rtt-resumption/#attack-of-the-clones)

Enable 0-RTT QUIC connection handshake on the client side
This is not impacting much on the performance, as the protocol is fully multiplexed

#### heartbeat

Interval for sending heartbeat packets for keeping the connection alive

- Default: `10s`

#### tls

==Required==

TLS configuration, see [TLS](/configuration/shared/tls/#inbound).

#### users.name


<!-- AUTO-GENERATED-CONFIG-END -->

### QUIC Fields

See [QUIC Fields](/configuration/shared/quic/) for details.
