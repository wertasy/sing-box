---
icon: material/new-box
---

!!! question "Since sing-box 1.12.0"

# Resolved

Resolved service is a fake systemd-resolved DBUS service to receive DNS settings from other programs
(e.g. NetworkManager) and provide DNS resolution.

See also: [Resolved DNS Server](/configuration/dns/server/resolved/)

### Structure

```json
{
  "type": "resolved",
  
  ... // Listen Fields
}
```

### Listen Fields

<!-- AUTO-GENERATED-ORIGIN-START -->
See [Listen Fields](/configuration/shared/listen/) for details.
<!-- AUTO-GENERATED-ORIGIN-END -->

### Fields

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### listen

==Required==

Listen address.

`127.0.0.53` will be used by default.

#### listen_port

==Required==

Listen port.

`53` will be used by default.

<!-- AUTO-GENERATED-CONFIG-END -->
