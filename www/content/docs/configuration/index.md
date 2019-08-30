---
title: 'Configuration'
date: 2019-02-11T19:27:37+10:00
weight: 4
---

Pomo has a few configuration options which can be read from a JSON file in Pomo's state directory `~/.pomo/config.json`.

### colors

You can map colors to specific tags in the `colors` field.

Example:
```json
{
    "colors": {
        "my-project": "hiyellow",
        "another-project": "green"
    }
}
```

