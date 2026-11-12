### 结构

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "shadowsocks",
  "tag": "ss-in",
  "method": "2022-blake3-aes-128-gcm",
  "password": "8JCsPssfgS8tiRwiMlhARg==",
  "multiplex": {},
  "managed": false

  ... // 监听字段
}
```
<!-- AUTO-GENERATED-EXAMPLE-END -->

### 多用户结构

```json
{
  "method": "2022-blake3-aes-128-gcm",
  "password": "8JCsPssfgS8tiRwiMlhARg==",
  "users": [
    {
      "name": "sekai",
      "password": "PCD2Z4o12bKUoFa3cC97Hw=="
    }
  ],
  "multiplex": {}
}
```

### 中转结构

```json
{
  "type": "shadowsocks",
  "method": "2022-blake3-aes-128-gcm",
  "password": "8JCsPssfgS8tiRwiMlhARg==",
  "destinations": [
    {
      "name": "test",
      "server": "example.com",
      "server_port": 8080,
      "password": "PCD2Z4o12bKUoFa3cC97Hw=="
    }
  ],
  "multiplex": {}
}
```

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

#### method

==必填==

| 方法                            | 密钥长度 |
|-------------------------------|------|
| 2022-blake3-aes-128-gcm       | 16   |
| 2022-blake3-aes-256-gcm       | 32   |
| 2022-blake3-chacha20-poly1305 | 32   |
| none                          | /    |
| aes-128-gcm                   | /    |
| aes-192-gcm                   | /    |
| aes-256-gcm                   | /    |
| chacha20-ietf-poly1305        | /    |
| xchacha20-ietf-poly1305       | /    |

#### password

==必填==

| 方法            | 密码格式                                     |
|---------------|------------------------------------------|
| none          | /                                        |
| 2022 methods  | `sing-box generate rand --base64 <密钥长度>` |
| other methods | 任意字符串                                    |

#### managed

默认为 `false`。当该入站需要由 [SSM API](/zh/configuration/service/ssm-api) 管理用户时必须启用此字段。

#### multiplex

参阅 [多路复用](/zh/configuration/shared/multiplex#入站)。

<!-- AUTO-GENERATED-CONFIG-END -->
