[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/barryorourke/check-tempager-3e-temperature)
![Go Test](https://github.com/barryorourke/check-tempager-3e-temperature/workflows/Go%20Test/badge.svg)
![goreleaser](https://github.com/barryorourke/check-tempager-3e-temperature/workflows/goreleaser/badge.svg)

# Check Plugin Template

## Overview
# check-tempager-3e-temperature

## Table of Contents
- [Overview](#overview)
- [Files](#files)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Check definition](#check-definition)
- [Installation from source](#installation-from-source)
- [Additional notes](#additional-notes)
- [Contributing](#contributing)

## Overview

check-tempager-3e-temperature is a [Sensu Check][6] that reports on the temperature gathered via SNMP
from a Tempager 3e temperature monitor.

Like all Sensu plugins this one will run on anything that supports Nagios plugin output.

I'm still pretty new to go, so please go easy on me in the Issues and Pull requests.

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add barryorourke/check-tempager-3e-temperature
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][https://bonsai.sensu.io/assets/barryorourke/check-tempager-3e-temperature].

### Check definition

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: check-tempager-3e-temperature
  namespace: default
spec:
  command: check-tempager-3e-temperature -t <ip address of tempager> -c <community string>
  subscriptions:
  - system
  runtime_assets:
  - barryorourke/check-tempager-3e-temperature
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the check-tempager-3e-temperature repository:

```
go build
```

## Additional notes

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
[2]: https://github.com/sensu-community/sensu-plugin-sdk
[3]: https://github.com/sensu-plugins/community/blob/master/PLUGIN_STYLEGUIDE.md
[4]: https://github.com/sensu-community/check-plugin-template/blob/master/.github/workflows/release.yml
[5]: https://github.com/sensu-community/check-plugin-template/actions
[6]: https://docs.sensu.io/sensu-go/latest/reference/checks/
[7]: https://github.com/sensu-community/check-plugin-template/blob/master/main.go
[8]: https://bonsai.sensu.io/
[9]: https://github.com/sensu-community/sensu-plugin-tool
[10]: https://docs.sensu.io/sensu-go/latest/reference/assets/
