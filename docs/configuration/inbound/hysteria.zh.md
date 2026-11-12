### 结构

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "hysteria",
  "tag": "hysteria-in",
  "up": "100 Mbps",
  "up_mbps": 100,
  "down": "100 Mbps",
  "down_mbps": 100,
  "obfs": "fuck me till the daylight",
  "users": [
    {
      "name": "sekai",
      "auth": "",
      "auth_str": "password"
    }
  ],
  "tls": {},

  ... // 监听字段
  ... // QUIC 字段

  // Deprecated
  "recv_window_conn": 0,
  "recv_window_client": 0,
  "max_conn_client": 0,
  "disable_mtu_discovery": false
}
```
<!-- AUTO-GENERATED-EXAMPLE-END -->

### 监听字段

<!-- AUTO-GENERATED-ORIGIN-START -->
参阅 [监听字段](/zh/configuration/shared/listen/) 了解详情。
<!-- AUTO-GENERATED-ORIGIN-END -->

### 字段

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### up, down

==必填==

格式: `[Integer] [Unit]` 例如： `100 Mbps, 640 KBps, 2 Gbps`

支持的单位 (大小写敏感, b = bits, B = bytes, 8b=1B)：

    bps (bits per second)
    Bps (bytes per second)
    Kbps (kilobits per second)
    KBps (kilobytes per second)
    Mbps (megabits per second)
    MBps (megabytes per second)
    Gbps (gigabits per second)
    GBps (gigabytes per second)
    Tbps (terabits per second)
    TBps (terabytes per second)

#### up_mbps, down_mbps

==必填==

以 Mbps 为单位的 `up, down`。

#### obfs

混淆密码。

#### users

Hysteria 用户

#### users.auth

base64 编码的认证密码。

#### users.auth_str

认证密码。

#### tls

==必填==

TLS 配置, 参阅 [TLS](/zh/configuration/shared/tls/#入站)。

#### users.name


#### recv_window_conn

!!! failure "已在 sing-box 1.14.0 废弃"

    请使用 QUIC 字段 `stream_receive_window` 代替。

- 已废弃。

#### recv_window_client

!!! failure "已在 sing-box 1.14.0 废弃"

    请使用 QUIC 字段 `connection_receive_window` 代替。

- 已废弃。

#### max_conn_client

!!! failure "已在 sing-box 1.14.0 废弃"

    请使用 QUIC 字段 `max_concurrent_streams` 代替。

- 已废弃。

#### disable_mtu_discovery

!!! failure "已在 sing-box 1.14.0 废弃"

    请使用 QUIC 字段 `disable_path_mtu_discovery` 代替。

- 已废弃。

<!-- AUTO-GENERATED-CONFIG-END -->

### QUIC 字段

参阅 [QUIC 字段](/zh/configuration/shared/quic/) 了解详情。

### 废弃字段

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### recv_window_conn

!!! failure "已在 sing-box 1.14.0 废弃"

    请使用 QUIC 字段 `stream_receive_window` 代替。

#### recv_window_client

!!! failure "已在 sing-box 1.14.0 废弃"

    请使用 QUIC 字段 `connection_receive_window` 代替。

#### max_conn_client

!!! failure "已在 sing-box 1.14.0 废弃"

    请使用 QUIC 字段 `max_concurrent_streams` 代替。

#### disable_mtu_discovery

!!! failure "已在 sing-box 1.14.0 废弃"

    请使用 QUIC 字段 `disable_path_mtu_discovery` 代替。

<!-- AUTO-GENERATED-CONFIG-END -->
