---
icon: material/new-box
---

!!! question "Since sing-box 1.14.0"

# Default

Attach to an existing network namespace.

### Structure

```json
{
  "network_namespaces": [
    {
      "type": "default", // optional
      "tag": "",
      "path": ""
    }
  ]
}
```

### Fields

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### path

==Required==

Name or path of the network namespace, for example `sing` or `/run/netns/sing`.

<!-- AUTO-GENERATED-CONFIG-END -->
