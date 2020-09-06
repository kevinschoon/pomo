<p align="center"><img src="https://raw.githubusercontent.com/kevinschoon/pomo/master/www/static/demo.gif" alt="demo"/></p>

# 🍅 pomo

![pomo](https://github.com/kevinschoon/pomo/workflows/pomo/badge.svg)

`pomo` is a simple CLI for using the [Pomodoro Technique](https://en.wikipedia.org/wiki/Pomodoro_Technique). There are [some](https://taskwarrior.org/) [amazing](https://todoist.com/) task management systems but `pomo` is more of a task *execution* or [timeboxing](https://en.wikipedia.org/wiki/Timeboxing) system. `pomo` helps you track what you did, how long it took you to do it, and how much effort you expect it to take.

## Background

The Pomodoro Technique is simple and effective:

  * Decide on a task you want to accomplish
  * Break the task into timed intervals (pomodoros), [approx. 25 min]
  * After each pomodoro take a short break [approx. 3 - 5 min]
  * Once all pomodoros are completed take a longer break [approx 15 - 20 min]
  * Repeat

## Installation

### Binaries

Binaries are available for Linux and OSX platforms in the [releases section](https://github.com/kevinschoon/pomo/releases) on github.

### Installer Script

A bash script to download and verify the latest release for Linux and OSX platforms can be run
with the following command:

```bash
curl -L -s https://kevinschoon.github.io/pomo/install.sh | bash /dev/stdin
```

### Source

 ```bash
 go get github.com/kevinschoon/pomo
 pomo -v
 ```

## Usage

Once `pomo` is installed you need to initialize it's database.

``` bash
pomo init
```

Start a 4 pomodoro session at 25 minute intervals:
```bash
pomo start -t my-project "write some codes"
```

## Configuration

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

## Integrations

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


## Roadmap

  * Generate charts/burn down
  * ??

## Credits

 * [pomodoro technique](https://cirillocompany.de/pages/pomodoro-technique/book/)
 * [logo by rones](https://openclipart.org/detail/262421/tomato-by-rones)
 * [website generate by hugo](http://gohugo.io/)
 * [theme by calintat](https://github.com/calintat/minimal)
