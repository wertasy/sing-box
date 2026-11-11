---
icon: material/new-box
---

!!! question "Since sing-box 1.12.0"

# Fake IP

### Structure

```json
{
  "dns": {
    "servers": [
      {
        "type": "fakeip",
        "tag": "",

        "inet4_range": "198.18.0.0/15",
        "inet6_range": "2001:2::/48"
      }
    ]
  }
}
```

### Fields

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### inet4_range

IPv4 address range for FakeIP.

#### inet6_range

IPv6 address range for FakeIP.

<!-- AUTO-GENERATED-CONFIG-END -->
