---
icon: material/new-box
---

!!! question "Since sing-box 1.12.0"

### Structure

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "anytls",
  "tag": "anytls-in",
  "tls": {},
  "users": [
    {
      "name": "sekai",
      "password": "8JCsPssfgS8tiRwiMlhARg=="
    }
  ],
  "padding_scheme": []

  ... // Listen Fields
}
```
<!-- AUTO-GENERATED-EXAMPLE-END -->

### Listen Fields

<!-- AUTO-GENERATED-ORIGIN-START -->
See [Listen Fields](/configuration/shared/listen/) for details.
<!-- AUTO-GENERATED-ORIGIN-END -->

### Fields

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### users

==Required==

AnyTLS users.

#### padding_scheme

AnyTLS padding scheme line array.

Default padding scheme:

```json
[
  "stop=8",
  "0=30-30",
  "1=100-400",
  "2=400-500,c,500-1000,c,500-1000,c,500-1000,c,500-1000",
  "3=9-9,500-1000",
  "4=500-1000",
  "5=500-1000",
  "6=500-1000",
  "7=500-1000"
]
```

#### tls

TLS configuration, see [TLS](/configuration/shared/tls/#inbound).

<!-- AUTO-GENERATED-CONFIG-END -->
