### Structure

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "vmess",
  "tag": "vmess-in",
  "users": [
    {
      "name": "sekai",
      "uuid": "bf000d23-0752-40b4-affe-68f7707a9661",
      "alterId": 0
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

!!! warning ""

    Legacy protocol support (VMess MD5 Authentication) is provided for compatibility purposes only, use of alterId > 1 is not recommended.

VMess users.

| Alter ID | Description             |
|----------|-------------------------|
| 0        | Disable legacy protocol |
| > 0      | Enable legacy protocol  |

#### tls

TLS configuration, see [TLS](/configuration/shared/tls/#inbound).

#### multiplex

See [Multiplex](/configuration/shared/multiplex#inbound) for details.

#### transport

V2Ray Transport configuration, see [V2Ray Transport](/configuration/shared/v2ray-transport/).

<!-- AUTO-GENERATED-CONFIG-END -->
