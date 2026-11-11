### Server Requirements

* Linux
* `brutal` congestion control algorithm kernel module installed

See [tcp-brutal](https://github.com/apernet/tcp-brutal) for details.

### Structure

```json
{
  "enabled": true,
  "up_mbps": 100,
  "down_mbps": 100
}
```

### Fields

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### enabled

Enable TCP Brutal congestion control algorithm。

#### up_mbps, down_mbps

==Required==

Upload and download bandwidth, in Mbps.

<!-- AUTO-GENERATED-CONFIG-END -->
