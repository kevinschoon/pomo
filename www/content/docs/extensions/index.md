---
title: 'Extensions'
date: 2019-02-11T19:27:37+10:00
weight: 5
---

### Status Bars

The Pomo CLI can output the current state of a running task session via the `pomo status`
making it easy to script and embed it's output in various Linux status bars.

#### [Polybar](https://github.com/jaagr/polybar)

You can create a module with the `custom/script` type and 
embed Pomo's status output in your Polybar:

```ini
[module/pomo]
type = custom/script
interval = 1
exec = pomo status
```
