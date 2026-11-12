# NTP

内建的 NTP 客户端服务。

如果启用，它将为像 TLS/Shadowsocks/VMess 这样的协议提供时间，这对于无法进行时间同步的环境很有用。

### 结构

```json
{
  "ntp": {
    "enabled": false,
    "server": "time.apple.com",
    "server_port": 123,
    "interval": "30m",
    
    ... // 拨号字段
  }
}

```

### 字段

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### enabled

启用 NTP 服务。

#### server

==必填==

NTP 服务器地址。

#### server_port

NTP 服务器端口。

默认使用 123。

#### interval

时间同步间隔。

- 默认值：`30 minutes`

#### write_to_system

Update system time after synchronization.

<!-- AUTO-GENERATED-CONFIG-END -->

### 拨号字段

<!-- AUTO-GENERATED-ORIGIN-START -->
参阅 [拨号字段](/zh/configuration/shared/dial/) 了解详情。
<!-- AUTO-GENERATED-ORIGIN-END -->
