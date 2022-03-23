# sensu-check-dns-server

## Table of Contents
- [sensu-check-dns-server](#sensu-check-dns-server)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Usage examples](#usage-examples)
    - [Help output](#help-output)
  - [Configuration](#configuration)
    - [Asset registration](#asset-registration)
    - [Check definition](#check-definition)
  - [Installation from source](#installation-from-source)
  - [Contributing](#contributing)

## Overview

The check-dns-server is a [Sensu Check][6] that checks if a DNS server is working.  
It checks if the specified server responds to UPD, TCP or DoT requests.

## Usage examples

### Help output
```
Check DNS Server functionality

Usage:
  sensu-dns-server-check [flags]
  sensu-dns-server-check [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -h, --help                 help for sensu-dns-server-check
  -p, --port int             Port to check (default 53)
  -P, --protocol string      DNS Protocol to check (udp, tcp, dot) (default "udp")
  -r, --record string        DNS Record to check (default "sensu.io")
  -s, --server string        DNS Server to check
  -n, --server-name string   Hostname for DoT

Use "sensu-dns-server-check [command] --help" for more information about a command.
```

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add larsl-net/sensu-check-dns-server
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][https://bonsai.sensu.io/assets/larsl-net/sensu-check-dns-server].

### Check definition

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: sensu-check-dns-server
  namespace: default
spec:
  command: sensu-dns-server-check -s 1.1.1.1 -p 853 -P dot -r sensu.io -n cloudflare-dns.com
  subscriptions:
  - system
  runtime_assets:
  - larsl-net/sensu-check-dns-server
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the sensu-check-dns-server repository:

```
go build
```


## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
