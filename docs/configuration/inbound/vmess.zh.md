### 结构

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "vmess",
  "tag": "vmess-in",
  "users": [
    {
      "name": "sekai",
      "uuid": "bf000d23-0752-40b4-affe-68f7707a9661",
      "alterId": 0
    }
  ],
  "tls": {},
  "multiplex": {},
  "transport": {}

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

#### users

==必填==

!!! warning ""

    提供旧协议支持（VMess MD5 身份验证）仅出于兼容性目的，不建议使用 alterId > 1。

VMess 用户。

| Alter ID | 描述    |
|----------|-------|
| 0        | 禁用旧协议 |
| > 0      | 启用旧协议 |

#### tls

TLS 配置, 参阅 [TLS](/zh/configuration/shared/tls/#入站)。

#### multiplex

参阅 [多路复用](/zh/configuration/shared/multiplex#入站)。

#### transport

V2Ray 传输配置，参阅 [V2Ray 传输层](/zh/configuration/shared/v2ray-transport/)。

<!-- AUTO-GENERATED-CONFIG-END -->
