---
icon: material/new-box
---

!!! question "自 sing-box 1.12.0 起"

# DNS over TLS (DoT)

### 结构

```json
{
  "dns": {
    "servers": [
      {
        "type": "tls",
        "tag": "",

        "server": "",
        "server_port": 853,

        "tls": {},

        // 拨号字段
      }
    ]
  }
}
```

!!! info "与旧版 TLS 服务器的区别"

    * 旧服务器默认使用默认出站，除非指定了绕行；新服务器像出站一样使用拨号器，相当于默认使用空的直连出站。
    * 旧服务器使用 `address_resolver` 和 `address_strategy` 来解析服务器中的域名；新服务器改用 [拨号字段](/zh/configuration/shared/dial/) 中的 `domain_resolver` 和 `domain_strategy`。

### 字段

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### server

==必填==

DNS 服务器的地址。

如果使用域名，还必须设置 `domain_resolver` 来解析 IP 地址。

#### server_port

DNS 服务器的端口。

默认使用 `853`。

- 默认值：`53`

#### tls

TLS 配置，参阅 [TLS](/zh/configuration/shared/tls/#出站)。

<!-- AUTO-GENERATED-CONFIG-END -->

### 拨号字段

<!-- AUTO-GENERATED-ORIGIN-START -->
参阅 [拨号字段](/zh/configuration/shared/dial/) 了解详情。
<!-- AUTO-GENERATED-ORIGIN-END -->
