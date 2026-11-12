### Structure

<!-- AUTO-GENERATED-EXAMPLE-START -->
```json
{
  "type": "tor",
  "tag": "tor-out",
  "executable_path": "/usr/bin/tor",
  "extra_args": [],
  "data_directory": "$HOME/.cache/tor",
  "torrc": {
    "ClientOnly": 1
  }

  ... // Dial Fields
}
```
<!-- AUTO-GENERATED-EXAMPLE-END -->

!!! info ""

    Embedded Tor is not included by default, see [Installation](/installation/build-from-source/#build-tags).

### Fields

<!-- AUTO-GENERATED-CONFIG-START -->
<!-- 此区域由脚本自动生成，切勿手动修改 -->

#### executable_path

The path to the Tor executable.

Embedded Tor will be ignored if set.

#### extra_args

List of extra arguments passed to the Tor instance when started.

#### data_directory

==Recommended==

The data directory of Tor.

Each start will be very slow if not specified.

#### torrc

Map of torrc options.

See [tor(1)](https://linux.die.net/man/1/tor) for details.

<!-- AUTO-GENERATED-CONFIG-END -->

### Dial Fields

<!-- AUTO-GENERATED-ORIGIN-START -->
See [Dial Fields](/configuration/shared/dial/) for details.
<!-- AUTO-GENERATED-ORIGIN-END -->
