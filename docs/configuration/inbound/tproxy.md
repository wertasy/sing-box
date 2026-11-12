!!! quote "Changes in sing-box 1.14.0"

    :material-plus: [udp_mapping](/configuration/shared/udp-nat/#udp_mapping)  
    :material-plus: [udp_filtering](/configuration/shared/udp-nat/#udp_filtering)  
    :material-plus: [udp_nat_max](/configuration/shared/udp-nat/#udp_nat_max)

!!! quote ""

    Only supported on Linux.

### Structure

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "tproxy",
  "tag": "tproxy-in",
  "network": "udp"

  ... // Listen Fields
  ... // UDP NAT Fields
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

<!-- AUTO-GENERATED-CONFIG-END -->

### UDP NAT Fields

See [UDP NAT Fields](/configuration/shared/udp-nat/) for details.
