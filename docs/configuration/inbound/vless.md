### Structure

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "vless",
  "tag": "vless-in",
  "users": [
    {
      "name": "sekai",
      "uuid": "bf000d23-0752-40b4-affe-68f7707a9661",
      "flow": ""
    }
  ],
  "tls": {},
  "multiplex": {},
  "transport": {}

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

==Required==

VLESS users.

#### users.uuid

==Required==

VLESS user id.

#### users.flow

VLESS Sub-protocol.

Available values:

* `xtls-rprx-vision`

#### tls

TLS configuration, see [TLS](/configuration/shared/tls/#inbound).

#### multiplex

See [Multiplex](/configuration/shared/multiplex#inbound) for details.

#### transport

V2Ray Transport configuration, see [V2Ray Transport](/configuration/shared/v2ray-transport/).

<!-- AUTO-GENERATED-CONFIG-END -->
