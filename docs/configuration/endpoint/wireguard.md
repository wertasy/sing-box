!!! quote "Changes in sing-box 1.15.0"

    :material-plus: [on_demand](#on_demand)

!!! quote "Changes in sing-box 1.14.0"

    :material-plus: [udp_mapping](/configuration/shared/udp-nat/#udp_mapping)  
    :material-plus: [udp_filtering](/configuration/shared/udp-nat/#udp_filtering)  
    :material-plus: [udp_nat_max](/configuration/shared/udp-nat/#udp_nat_max)

!!! question "Since sing-box 1.11.0"

### Structure

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "wireguard",
  "tag": "wg-ep",
  "system": false,
  "name": "",
  "mtu": 1408,
  "address": [],
  "private_key": "",
  "listen_port": 10000,
  "peers": [
    {
      "address": "127.0.0.1",
      "port": 10001,
      "public_key": "",
      "pre_shared_key": "",
      "allowed_ips": [],
      "persistent_keepalive_interval": 0,
      "reserved": [0, 0, 0]
    }
  ],
  "udp_timeout": "",
  "udp_mapping": "",
  "udp_filtering": "",
  "udp_nat_max": 0,
  "workers": 0,
  "on_demand": false

  ... // Dial Fields
}
```
<!-- AUTO-GENERATED-EXAMPLE-END -->

!!! note ""

    You can ignore the JSON Array [] tag when the content is only one item

### Fields

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### system

Use system interface.

Requires privilege and cannot conflict with exists system interfaces.

#### name

Custom interface name for system interface.

#### mtu

WireGuard MTU.

- Default: `1408`

#### address

==Required==

List of IP (v4 or v6) address prefixes to be assigned to the interface.

#### private_key

==Required==

WireGuard requires base64-encoded public and private keys. These can be generated using the wg(8) utility:

```shell
wg genkey
echo "private key" || wg pubkey
```

or `sing-box generate wg-keypair`.

#### peers

==Required==

List of WireGuard peers.

#### peers.address

WireGuard peer address.

#### peers.port

WireGuard peer port.

#### peers.public_key

==Required==

WireGuard peer public key.

#### peers.pre_shared_key

WireGuard peer pre-shared key.

#### peers.allowed_ips

==Required==

WireGuard allowed IPs.

#### peers.persistent_keepalive_interval

WireGuard persistent keepalive interval, in seconds.

Disabled by default.

#### peers.reserved

WireGuard reserved field bytes.

#### workers

WireGuard worker count.

CPU count is used by default.

#### on_demand

!!! question "Since sing-box 1.15.0"

Allow the endpoint to be disconnected when necessary.

#### listen_port

Listen port for incoming connections.

<!-- AUTO-GENERATED-CONFIG-END -->

### UDP NAT Fields

See [UDP NAT Fields](/configuration/shared/udp-nat/) for details.

### Dial Fields

<!-- AUTO-GENERATED-ORIGIN-START -->
See [Dial Fields](/configuration/shared/dial/) for details.
<!-- AUTO-GENERATED-ORIGIN-END -->
