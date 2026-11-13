!!! question "Since sing-box 1.8.0"

!!! quote "Changes in sing-box 1.15.0"

    :material-plus: [buffer_size](#buffer_size)  
    :material-plus: [flush_interval](#flush_interval)

!!! quote "Changes in sing-box 1.14.0"

    :material-delete-clock: [store_rdrc](#store_rdrc)  
    :material-delete-clock: [rdrc_timeout](#rdrc_timeout)  
    :material-plus: [store_dns](#store_dns)

!!! quote "Changes in sing-box 1.9.0"

    :material-plus: [store_rdrc](#store_rdrc)  
    :material-plus: [rdrc_timeout](#rdrc_timeout)

### Structure

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "enabled": true,
  "path": "",
  "cache_id": "",
  "store_fakeip": false,
  "store_dns": false,
  "buffer_size": "",
  "flush_interval": "",

  // Deprecated
  "store_rdrc": false,
  "rdrc_timeout": ""
}
```
<!-- AUTO-GENERATED-EXAMPLE-END -->

### Fields

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### enabled

Enable cache file.

#### path

Path to the cache file.

`cache.db` will be used by default.

#### cache_id

Identifier in the cache file

If not empty, configuration specified data will use a separate store keyed by it.

#### store_fakeip

Store fakeip in the cache file

#### store_rdrc

!!! failure "Deprecated in sing-box 1.14.0"

    `store_rdrc` is deprecated and will be removed in sing-box 1.16.0, check [Migration](/migration/#migrate-store-rdrc).

Store rejected DNS response cache in the cache file

The check results of [Legacy Address Filter Fields](/configuration/dns/rule/#legacy-address-filter-fields)
will be cached until expiration.

- Deprecated.

#### rdrc_timeout

!!! failure "Deprecated in sing-box 1.14.0"

    `rdrc_timeout` is deprecated and will be removed in sing-box 1.16.0, check [Migration](/migration/#migrate-store-rdrc).

Timeout of rejected DNS response cache.

`7d` will be used by default.
- Deprecated.

#### store_dns

!!! question "Since sing-box 1.14.0"


Store DNS cache in the cache file.

#### buffer_size

!!! question "Since sing-box 1.15.0"


Size of the write buffer.

`1MB` will be used by default.

#### flush_interval

!!! question "Since sing-box 1.15.0"


Interval for flushing the write buffer automatically.

Disabled by default.

<!-- AUTO-GENERATED-CONFIG-END -->
