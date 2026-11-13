# NTP

Built-in NTP client service.

If enabled, it will provide time for protocols like TLS/Shadowsocks/VMess, which is useful for environments where time
synchronization is not possible.

### Structure

```json
{
  "ntp": {
    "enabled": false,
    "server": "time.apple.com",
    "server_port": 123,
    "interval": "30m",
    
    ... // Dial Fields
  }
}

```

### Fields

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### enabled

Enable NTP service.

#### server

==Required==

NTP server address.

#### server_port

NTP server port.

123 will be used by default.

#### interval

Time synchronization interval.

`30m` will be used by default.

#### write_to_system

Update system time after synchronization.

<!-- AUTO-GENERATED-CONFIG-END -->

### Dial Fields

<!-- AUTO-GENERATED-ORIGIN-START -->
See [Dial Fields](/configuration/shared/dial/) for details.
<!-- AUTO-GENERATED-ORIGIN-END -->
