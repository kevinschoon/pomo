---
title: 'Installation'
date: 2019-02-11T19:27:37+10:00
weight: 2
---

## Binaries

Binaries are available for Linux and OSX platforms in the [releases section](https://github.com/kevinschoon/pomo/releases) on github.

## Installer Script

A bash script to download and verify the latest release for Linux and OSX platforms can be run
with the following command:

```bash
curl -L -s https://kevinschoon.github.io/pomo/install.sh | bash /dev/stdin
```

## Source

 ```bash
 go get github.com/kevinschoon/pomo
 pomo -v
 ```

## Migration

Once `pomo` is installed you need to initialize it's database.

``` bash
pomo init
```

