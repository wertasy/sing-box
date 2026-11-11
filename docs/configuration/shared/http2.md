---
icon: material/new-box
---

!!! question "Since sing-box 1.14.0"

### Structure

```json
{
  "idle_timeout": "",
  "keep_alive_period": "",
  "stream_receive_window": "",
  "connection_receive_window": "",
  "max_concurrent_streams": 0
}
```

### Fields

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### connection_receive_window

HTTP2 connection-level flow-control receive window size.

Accepts memory size format, e.g. `"64 MB"`.

#### idle_timeout

Idle connection timeout, in golang's Duration format.

#### keep_alive_period

Keep alive period, in golang's Duration format.

#### max_concurrent_streams

Maximum concurrent streams per connection.

#### stream_receive_window

HTTP2 stream-level flow-control receive window size.

Accepts memory size format, e.g. `"64 MB"`.

<!-- AUTO-GENERATED-CONFIG-END -->
