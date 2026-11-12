# OpenConnect 客户端

!!! question "自 sing-box 1.14.0 起"

!!! quote "sing-box 1.15.0 中的更改"

    :material-plus: [on_demand](#on_demand)

==仅客户端==

### 结构

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "openconnect",
  "tag": "oc-client",
  "system": false,
  "name": "",
  "server": "vpn.example.com",
  "flavor": "anyconnect",
  "username": "",
  "password": "",
  "auth_group": "",
  "cookie": "",
  "token": {
    "mode": "",
    "secret": "",
    "secret_path": "",
    "pin": "",
    "password": "",
    "device_id": "",
    "counter": 0
  },
  "reported_os": "",
  "user_agent": "",
  "version": "",
  "local_hostname": "",
  "mobile": {
    "platform_version": "",
    "device_type": "",
    "device_unique_id": ""
  },
  "csd": {
    "wrapper_path": ""
  },
  "hip": {
    "wrapper_path": ""
  },
  "tncc": {
    "wrapper_path": "",
    "device_id": "",
    "user_agent": "",
    "machine_identification_enabled": false,
    "certificates": [
      {
        "certificate": [],
        "certificate_path": ""
      }
    ]
  },
  "fortinet_host_check": {
    "hostcheck": "",
    "check_virtual_desktop": ""
  },
  "no_udp": false,
  "dtls_local_port": 0,
  "compression_disabled": false,
  "compression_mode": "",
  "ipv6_disabled": false,
  "http_keepalive_disabled": false,
  "xml_post_disabled": false,
  "external_auth_disabled": false,
  "password_authentication_disabled": false,
  "tcp_keep_alive_enabled": false,
  "pfs": false,
  "mtu": 0,
  "base_mtu": 0,
  "dpd_interval": "",
  "reconnect_timeout": "",
  "trojan_interval": "",
  "queue_length": 0,
  "allow_insecure_crypto": false,
  "tls": {
    "insecure": false,
    "server_name": "",
    "peer_fingerprint": [],
    "system_trust_disabled": false,
    "certificate_authority": [],
    "certificate_authority_path": "",
    "client_certificate": [],
    "client_certificate_path": "",
    "client_key": [],
    "client_key_path": "",
    "client_key_password": "",
    "mca_certificate": [],
    "mca_certificate_path": "",
    "mca_key": [],
    "mca_key_path": "",
    "mca_key_password": ""
  },
  "form_entries": [
    {
      "form_id": "",
      "submission_key": "",
      "name": "",
      "value": "",
      "promote": false
    }
  ],
  "on_demand": false

  ... // 拨号字段
  ... // UDP NAT 字段
}
```
<!-- AUTO-GENERATED-EXAMPLE-END -->

!!! note ""

    当内容只有一项时，可以忽略 JSON 数组 [] 标签。

### 字段

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### system

使用系统接口。

需要权限，且不能与现有系统接口冲突。

禁用时，sing-box 使用内部网络栈。

#### name

系统接口的自定义接口名称。

默认使用自动生成的 `oc` 接口名称。

#### server

==必填==

OpenConnect VPN 服务器 HTTPS URL。

省略协议时会添加 `https://`。不支持 URL 用户信息、查询和片段。

#### flavor

OpenConnect 协议 flavor，可选值为 `anyconnect`、`gp`、`fortinet`、`f5`、`pulse` 或 `nc`。

- 默认值：`anyconnect`

#### username

用于填充匹配认证表单字段的用户名。

#### password

用于填充匹配认证表单字段的密码。

#### auth_group

认证组，用于在所选 flavor 支持时预选匹配的组、realm、domain 或 gateway 选项。

#### cookie

用于跳过凭据提示并直接连接的现有认证会话。

接受的格式取决于 `flavor`：

- `anyconnect`：`webvpn` 值，或包含 `webvpn` 的分号分隔 cookie 列表。
- `gp`：GlobalProtect 认证返回的完整 authenticated query string。
- `nc`：`DSID` 值，或包含 `DSID` 的分号分隔 cookie 列表。
- `pulse`：原始 Pulse 认证 cookie 值。
- `f5`：`MRHSession` 值，或包含 `MRHSession` 及可选 `F5_ST` 的分号分隔 cookie 列表。
- `fortinet`：`SVPNCOOKIE` 值，或包含 `SVPNCOOKIE` 的分号分隔 cookie 列表。

如果服务器拒绝提供的会话，将尝试正常认证。

#### token

用于自动回答匹配 token 字段或进行 HTTP Bearer 认证的 token 配置。

必须设置 `token.secret` 或 `token.secret_path` 之一。

#### reported_os

所选 flavor 支持时向 VPN 服务器报告的操作系统标识。

对于 `anyconnect`、`gp` 和 `pulse`，支持的值为 `linux`、`linux-64`、`win`、`mac-intel`、`android` 和 `apple-ios`。

默认值根据系统平台选择：Windows 使用 `win`，macOS 使用 `mac-intel`，Android 使用 `android`，iOS 使用 `apple-ios`，其他 64 位或 32 位系统使用 `linux-64` 或 `linux`。

#### user_agent

所选 flavor 支持时向 VPN 服务器报告的 User-Agent。

默认值由 flavor 决定。AnyConnect、Network Connect、Pulse 和 F5 使用 `AnyConnect-compatible OpenConnect VPN Agent v9.21`；GlobalProtect 使用 `PAN GlobalProtect`；Fortinet 使用 `Mozilla/5.0 SV1`。

#### version

所选 flavor 支持时，与 `user_agent` 分开报告的客户端版本。

- 默认值：`v9.21`

#### local_hostname

所选 flavor 支持时向 VPN 服务器报告的本地主机名。

默认使用系统主机名；无法获取时使用 `localhost`。

#### mobile

AnyConnect 移动客户端身份。配置时三个字段均为必填。

#### csd

AnyConnect CSD/host scan 合规性选项。

服务器请求 CSD 时，默认使用内置 CSD 处理。

#### hip

GlobalProtect HIP 检查和报告选项。

服务器请求 HIP 时，默认使用内置 HIP 报告。

#### tncc

Network Connect TNCC 合规性选项。

服务器请求 TNCC 时，默认使用内置 TNCC 处理。

#### fortinet_host_check

Fortinet hostcheck 结果覆盖选项，在服务器要求时提交。

`fortinet_host_check.hostcheck` 为空时禁用。

部分 Fortinet 服务器只会要求可识别的 FortiClient User-Agent 执行 hostcheck。服务器策略有要求时请配置 `user_agent`。

#### no_udp

禁用 DTLS 或 ESP 辅助数据通道，仅使用 TLS 数据通道。

#### dtls_local_port

直连 DTLS 或 ESP 辅助数据通道使用的本地 UDP 端口。

默认自动选择临时端口。

#### compression_disabled

禁用 AnyConnect 压缩协商。

默认情况下，当服务器支持时，CSTP 和 DTLS 会协商无状态 `oc-lz4` 和 `lzs` 压缩。

与设置为 `all` 的 `compression_mode` 冲突。

#### compression_mode

AnyConnect 压缩模式，可选值为：

- `stateless`：声明支持无状态 `oc-lz4` 和 `lzs` 压缩。
- `all`：额外声明支持 CSTP 有状态 `deflate` 压缩。

- 默认值：`stateless`

#### ipv6_disabled

禁用请求和使用 IPv6 隧道配置。

#### http_keepalive_disabled

在认证和配置请求中禁用 HTTP 连接复用。

#### xml_post_disabled

禁用 AnyConnect XML POST 认证，并直接使用旧版 GET 流程开始认证。

#### external_auth_disabled

禁用 AnyConnect、GlobalProtect 和 Fortinet 的 SSO、SAML 等外部浏览器认证。

#### password_authentication_disabled

如果服务器返回非成功的认证表单，则中止 AnyConnect 认证，与 OpenConnect `--no-passwd` 行为一致。

#### tcp_keep_alive_enabled

为直接 VPN 服务器连接启用 TCP keep alive。

设置 `tcp_keep_alive` 或 `tcp_keep_alive_interval` 也会启用。启用但未设置这两个时间值时，使用操作系统默认值。

与 `disable_tcp_keep_alive` 冲突。

#### pfs

要求 TLS 1.2 及更早版本使用具有前向保密性的 TLS 密码套件。

#### mtu

首选隧道 MTU。

协商的 MTU 不超过此值。

非零值小于 `576` 时按 `576` 处理。最大值为 `65535`。

#### base_mtu

扣除外层 IP、传输和协议开销后，用于计算 AnyConnect、GlobalProtect、F5 和 Fortinet 隧道 MTU 的基础路径 MTU。

这些 flavor 会将小于 `1280` 的值按 `1280` 处理。最大值为 `65535`。

- 默认值：`1406`

#### dpd_interval

覆盖 Dead Peer Detection 间隔。

默认使用服务器提供或 flavor 特定的间隔。

大于零且小于 `2s` 的值按 `2s` 处理。值不得为负数。

#### reconnect_timeout

重连尝试失败后允许累计使用的最大退避时间。

值不得为负数。

- 默认值：`300s`

#### trojan_interval

覆盖 GlobalProtect HIP report 或 Network Connect TNCC check 的执行间隔。

默认使用服务器提供的间隔。服务器未提供时，GlobalProtect 使用 `1h`。

值不得为负数。

#### queue_length

VPN transport 与隧道接口之间的入站和出站数据包队列长度。

- 默认值：`32`

#### allow_insecure_crypto

启用旧版 VPN 服务器所需的弱 TLS 和 DTLS 密码套件及 TLS 1.0 兼容性。

默认禁用；未启用时会拒绝低于 TLS 1.2 的版本。

#### tls

OpenConnect TLS 配置。

#### form_entries

认证表单字段覆盖。

设置 `submission_key` 时按该字段匹配，否则按 `form_id` 和 `name` 的组合匹配。后面的匹配项优先。

#### on_demand

!!! question "自 sing-box 1.15.0 起"

允许该 endpoint 在需要时断开连接。

<!-- AUTO-GENERATED-CONFIG-END -->

### UDP NAT 字段

参阅 [UDP NAT 字段](/zh/configuration/shared/udp-nat/)。

### 拨号字段

<!-- AUTO-GENERATED-ORIGIN-START -->
参阅 [拨号字段](/zh/configuration/shared/dial/) 了解详情。
<!-- AUTO-GENERATED-ORIGIN-END -->

## 交互式认证

在 sing-box dashboard 或任意 sing-box 图形客户端的 `工具` > `端点` 中认证和管理 endpoint。

## DNS

配置 [OpenConnect DNS 服务器](/zh/configuration/dns/server/openconnect/) 以使用推送的 DNS 设置。
