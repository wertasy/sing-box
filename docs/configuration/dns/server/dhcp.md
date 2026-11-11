---
icon: material/new-box
---

!!! question "Since sing-box 1.12.0"

# DHCP

### Structure

```json
{
  "dns": {
    "servers": [
      {
        "type": "dhcp",
        "tag": "",

        "interface": "",
        
        // Dial Fields
      }
    ]
  }
}
```

### Fields

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### interface

Interface name to listen on.

Tge default interface will be used by default.

<!-- AUTO-GENERATED-CONFIG-END -->

### Dial Fields

See [Dial Fields](/configuration/shared/dial/) for details. 
