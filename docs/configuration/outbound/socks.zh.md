`socks` 出站是 socks4/socks4a/socks5 客户端

### 结构

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "socks",
  "tag": "socks-out",
  "server": "",
  "server_port": 0,
  "version": "5",
  "username": "sekai",
  "password": "admin",
  "network": "udp",
  "udp_over_tcp": false | {}

  ... // 拨号字段
}
```
<!-- AUTO-GENERATED-EXAMPLE-END -->

### 字段

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### server

==必填==

服务器地址。

#### server_port

==必填==

服务器端口。

#### version

SOCKS 版本, 可为 `4` `4a` `5`.

默认使用 SOCKS5。

#### username

SOCKS 用户名。

#### password

SOCKS5 密码。

#### network

启用的网络协议

`tcp` 或 `udp`。

默认所有。

#### udp_over_tcp

UDP over TCP 配置。

参阅 [UDP Over TCP](/zh/configuration/shared/udp-over-tcp/)。

<!-- AUTO-GENERATED-CONFIG-END -->

### 拨号字段

<!-- AUTO-GENERATED-ORIGIN-START -->
参阅 [拨号字段](/zh/configuration/shared/dial/) 了解详情。
<!-- AUTO-GENERATED-ORIGIN-END -->
