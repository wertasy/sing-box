### Structure

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "selector",
  "tag": "select",
  "outbounds": [
    "proxy-a",
    "proxy-b",
    "proxy-c"
  ],
  "default": "proxy-c",
  "interrupt_exist_connections": false
}
```
<!-- AUTO-GENERATED-EXAMPLE-END -->

!!! quote ""

    The selector can only be controlled through the [Clash API](/configuration/experimental#clash-api-fields) currently.

### Fields

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### outbounds

==Required==

List of outbound tags to select.

#### default

The default outbound tag. The first outbound will be used if empty.

#### interrupt_exist_connections

Interrupt existing connections when the selected outbound has changed.

Only inbound connections are affected by this setting, internal connections will always be interrupted.

<!-- AUTO-GENERATED-CONFIG-END -->
