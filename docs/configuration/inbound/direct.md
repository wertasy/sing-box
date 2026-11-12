`direct` inbound is a tunnel server.

### Structure

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "direct",
  "tag": "direct-in",
  "network": "udp",
  "override_address": "1.0.0.1",
  "override_port": 53

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

#### override_address

Override the connection destination address.

#### override_port

Override the connection destination port.

<!-- AUTO-GENERATED-CONFIG-END -->
