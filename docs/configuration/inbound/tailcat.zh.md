---
icon: material/new-box
---

!!! question "自 sing-box 1.15.0 起"

### 结构

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "tailcat",
  "tag": "tailcat-in",
  "private_key": "",
  "pre_shared_key": "",
  "users": [
    {
      "name": "",
      "public_key": ""
    }
  ],
  "derp_map_url": "",
  "derp_region": 0,
  "derp_servers": [],
  "http_client": {},
  "bind_interface": "",
  "inet4_bind_address": "",
  "inet6_bind_address": "",
  "bind_address_no_port": false,
  "protect_path": "",
  "routing_mark": 0,
  "reuse_addr": false,
  "netns": "",
  "connect_timeout": "",
  "tcp_fast_open": false,
  "tcp_multi_path": false,
  "disable_tcp_keep_alive": false,
  "tcp_keep_alive": "",
  "tcp_keep_alive_interval": "",
  "udp_fragment": false,
  "domain_resolver": "",
  "network_strategy": "",
  "network_type": [],
  "fallback_network_type": [],
  "fallback_delay": "",

  // Deprecated
  "domain_strategy": ""
}
```
<!-- AUTO-GENERATED-EXAMPLE-END -->

### 字段

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### private_key

==必填==

私钥。

使用 `sing-box generate tailcat-keypair` 生成。

#### pre_shared_key

预共享密钥。

#### users

Tailcat 用户。

为空时不验证客户端。

如需使用带客户端验证的 DERP 服务器，在 [DERP 服务](/zh/configuration/service/derp/#verify_client_inbound)
中设置 `verify_client_inbound` 或 `verify_client_key`，且客户端必须使用固定的私钥。

#### users.public_key

==必填==

客户端公钥。

#### derp_map_url

[DERP map](https://pkg.go.dev/tailscale.com/tailcfg#DERPMap) 的 URL。

- 默认值：`https://tailcat.dev/derpmap.json`

#### derp_region

DERP map 中的 DERP 区域 ID。

与 `derp_servers` 冲突。

#### derp_servers

自定义 DERP 服务器，为 [DERPNode](https://pkg.go.dev/tailscale.com/tailcfg#DERPNode) 格式，字段名使用 snake_case。

与 `derp_map_url` 和 `derp_region` 冲突。

将数组值设置为字符串 `__HOST__` 等同于配置：

```json
{ "host": __HOST__ }
```

#### http_client

用于获取 DERP map 的 HTTP 客户端。

参阅 [HTTP 客户端](/zh/configuration/shared/http-client/) 了解详情。

#### bind_interface

The network interface to bind to.

#### inet4_bind_address

The IPv4 address to bind to.

#### inet6_bind_address

The IPv6 address to bind to.

#### bind_address_no_port

!!! question "Since sing-box 1.13.0"

!!! quote ""

    Only supported on Linux.

Do not reserve a port when binding to a source address.

This allows reusing the same source port for multiple connections if the full 4-tuple (source IP, source port, destination IP, destination port) remains unique.

#### protect_path


#### routing_mark

!!! quote ""

    Only supported on Linux.

Set netfilter routing mark.

Integers (e.g. `1234`) and string hexadecimals (e.g. `"0x1234"`) are supported.

#### reuse_addr

Reuse listener address.

#### netns

!!! question "Since sing-box 1.12.0"

!!! quote ""

    Only supported on Linux.

Set network namespace, name or path.

Since sing-box 1.14.0, the tag of a [network namespace](/configuration/network-namespace/) can also be used.
Referencing an `unshare` network namespace should be avoided.

#### connect_timeout

Connect timeout, in golang's Duration format.

A duration string is a possibly signed sequence of
decimal numbers, each with optional fraction and a unit suffix,
such as "300ms", "-1.5h" or "2h45m".
Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".

#### tcp_fast_open

Enable TCP Fast Open.

#### tcp_multi_path

!!! warning ""

    Go 1.21 required.

Enable TCP Multi Path.

#### disable_tcp_keep_alive

!!! question "Since sing-box 1.13.0"

Disable TCP keep alive.

#### tcp_keep_alive

!!! question "Since sing-box 1.13.0"

    Default value changed from `10m` to `5m`.

TCP keep alive initial period.

- 默认值：`5m`

#### tcp_keep_alive_interval

!!! question "Since sing-box 1.13.0"

TCP keep alive interval.

- 默认值：`75s`

#### udp_fragment

Enable UDP fragmentation.

#### domain_resolver

!!! info ""

    `domain_resolver` or `route.default_domain_resolver` is optional when only one DNS server is configured.

!!! warning ""

    `outbound` DNS rule items are deprecated and will be removed in sing-box 1.14.0, so this item will be required for outbound/endpoints using domain name in server address since sing-box 1.14.0.

Set domain resolver to use for resolving domain names.

This option uses the same format as the [route DNS rule action](/configuration/dns/rule_action/#route) without the `action` field.

Setting this option directly to a string is equivalent to setting `server` of this options.

| Outbound/Endpoints | Effected domains         |
|--------------------|--------------------------|
| `direct`           | Domain in request        |
| others             | Domain in server address |

#### network_strategy

!!! question "Since sing-box 1.11.0"

!!! quote ""

    Only supported in graphical clients on Android and Apple platforms with `auto_detect_interface` enabled.

Strategy for selecting network interfaces.

Available values:

- `default` (default): Connect to default network or networks specified in `network_type` sequentially.
- `hybrid`: Connect to all networks or networks specified in `network_type` concurrently.
- `fallback`: Connect to default network or preferred networks specified in `network_type` concurrently, and try fallback networks when unavailable or timeout.

For fallback, when preferred interfaces fails or times out,
it will enter a 15s fast fallback state (Connect to all preferred and fallback networks concurrently),
and exit immediately if preferred networks recover.

Conflicts with `bind_interface`, `inet4_bind_address` and `inet6_bind_address`.

#### network_type

!!! question "Since sing-box 1.11.0"

!!! quote ""

    Only supported in graphical clients on Android and Apple platforms with `auto_detect_interface` enabled.

Network types to use when using `default` or `hybrid` network strategy or
preferred network types to use when using `fallback` network strategy.

Available values: `wifi`, `cellular`, `ethernet`, `other`.

Device's default network is used by default.

#### fallback_network_type

!!! question "Since sing-box 1.11.0"

!!! quote ""

    Only supported in graphical clients on Android and Apple platforms with `auto_detect_interface` enabled.

Fallback network types when preferred networks are unavailable or timeout when using `fallback` network strategy.

All other networks expect preferred are used by default.

#### fallback_delay

!!! question "Since sing-box 1.11.0"

!!! quote ""

    Only supported in graphical clients on Android and Apple platforms with `auto_detect_interface` enabled.

The length of time to wait before spawning a RFC 6555 Fast Fallback connection.

For `domain_strategy`, is the amount of time to wait for connection to succeed before assuming
that IPv4/IPv6 is misconfigured and falling back to other type of addresses.

For `network_strategy`, is the amount of time to wait for connection to succeed before falling
back to other interfaces.

Only take effect when `domain_strategy` or `network_strategy` is set.

- 默认值：`300ms`

#### domain_strategy

!!! failure "已在 sing-box 1.12.0 废弃"

    `domain_strategy` 已废弃且将在 sing-box 1.14.0 中被移除，参阅 [迁移指南](/zh/migration/#迁移出站域名策略选项到域名解析器)。

可选值：`prefer_ipv4` `prefer_ipv6` `ipv4_only` `ipv6_only`。

如果设置，域名将在请求发出之前解析为 IP。

| 出站       | 受影响的域名    | 默认回退值                     |
|----------|-----------|---------------------------|
| `direct` | 请求中的域名    | `inbound.domain_strategy` | 
| others   | 服务器地址中的域名 | /                         |

- 已废弃。

<!-- AUTO-GENERATED-CONFIG-END -->

### 拨号字段

!!! note

    Tailcat 入站中的拨号字段仅控制它如何连接到 DERP 服务器，与实际连接无关。

参阅 [拨号字段](/zh/configuration/shared/dial/) 了解详情。
