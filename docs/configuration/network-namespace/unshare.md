---
icon: material/new-box
---

!!! question "Since sing-box 1.14.0"

# Unshare

Create a new network namespace, without root privilege.

!!! info ""

    Rootless operation requires the kernel to allow unprivileged user namespace creation.

### Structure

```json
{
  "network_namespaces": [
    {
      "type": "unshare",
      "tag": "",
      "pid_file": ""
    }
  ]
}
```

### Fields

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### pid_file

If set, the PID of the process holding the namespace open is written to this path.

The namespace can be entered with `nsenter -t <pid> -n` when sing-box is run as root,
or `nsenter -t <pid> -U --preserve-credentials -n` otherwise.

<!-- AUTO-GENERATED-CONFIG-END -->
