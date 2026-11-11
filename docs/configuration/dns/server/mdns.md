---
icon: material/new-box
---

!!! question "Since sing-box 1.14.0"

# mDNS

### Structure

```json
{
  "dns": {
    "servers": [
      {
        "type": "mdns",
        "tag": "",

        "interface": [],

        // Dial Fields
      }
    ]
  }
}
```

!!! info ""

    `*.local.` and IPv4/IPv6 link-local reverse zones are also resolved by the [Local](./local/) server.

### Fields

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### interface

List of network interface names to send mDNS queries on.

When empty, all interfaces that are up, multicast-capable, and non-loopback are used.

<!-- AUTO-GENERATED-CONFIG-END -->

### Dial Fields

See [Dial Fields](/configuration/shared/dial/) for details.
