!!! quote "sing-box 1.13.0 中的更改"

    :material-plus: [quic_congestion_control](#quic_congestion_control)

### 结构

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "naive",
  "tag": "naive-in",
  "users": [
    {
      "username": "sekai",
      "password": "password"
    }
  ],
  "network": "udp",
  "quic_congestion_control": "",
  "tls": {}

  ... // 监听字段
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

#### network

监听的网络协议，`tcp` `udp` 之一。

默认所有。

#### users

==必填==

Naive 用户。

#### quic_congestion_control

!!! question "Since sing-box 1.13.0"

QUIC 拥塞控制算法。

| 算法             | 描述                 |
|----------------|--------------------|
| `bbr`          | BBR                |
| `cubic`        | CUBIC              |
| `reno`         | New Reno           |

- 默认值：`bbr`

#### tls

TLS 配置, 参阅 [TLS](/zh/configuration/shared/tls/#入站)。

<!-- AUTO-GENERATED-CONFIG-END -->
