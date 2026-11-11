---
icon: material/new-box
---

!!! question "自 sing-box 1.12.0 起"

# Fake IP

### 结构

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

### 字段

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### inet4_range

用于 FakeIP 的 IPv4 地址范围。

#### inet6_range

用于 FakeIP 的 IPv6 地址范围。

<!-- AUTO-GENERATED-CONFIG-END -->
