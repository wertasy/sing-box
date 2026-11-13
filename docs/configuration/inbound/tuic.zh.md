### 结构

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "tuic",
  "tag": "tuic-in",
  "users": [
    {
      "name": "sekai",
      "uuid": "059032A9-7D40-4A96-9BB1-36823D848068",
      "password": "hello"
    }
  ],
  "congestion_control": "cubic",
  "auth_timeout": "3s",
  "zero_rtt_handshake": false,
  "heartbeat": "10s",
  "tls": {}

  ... // 监听字段
  ... // QUIC 字段
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

#### users

TUIC 用户

#### users.uuid

==必填==

TUIC 用户 UUID

#### users.password

TUIC 用户密码

#### congestion_control

QUIC 拥塞控制算法

可选值: `cubic`, `new_reno`, `bbr`

默认使用 `cubic`。

#### auth_timeout

服务器等待客户端发送认证命令的时间

默认使用 `3s`。

#### zero_rtt_handshake

!!! warning ""
  强烈建议禁用此功能，因为它容易受到重放攻击。
  请参阅 [Attack of the clones](https://blog.cloudflare.com/even-faster-connection-establishment-with-quic-0-rtt-resumption/#attack-of-the-clones)

在客户端启用 0-RTT QUIC 连接握手
这对性能影响不大，因为协议是完全复用的

#### heartbeat

发送心跳包以保持连接存活的时间间隔

默认使用 `10s`。

#### tls

==必填==

TLS 配置, 参阅 [TLS](/zh/configuration/shared/tls/#入站)。

<!-- AUTO-GENERATED-CONFIG-END -->

### QUIC 字段

参阅 [QUIC 字段](/zh/configuration/shared/quic/) 了解详情。
