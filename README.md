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

### Source

 ```bash
 git clone git@github.com:kevinschoon/pomo.git
 cd pomo
 make
 # copy pomo somewhere on your $PATH
 cp bin/pomo ~/bin/
 ```

### Package Managers

On Arch Pomo is available on the [aur](https://aur.archlinux.org/packages/pomo).

On macOS, `pomo` can be installed via [MacPorts](https://ports.macports.org/port/pomo/).

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

Pomo has a few configuration options which can be read from a JSON file in Pomo's config directory `~/.config/pomo/config.json`.

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

### Execute command on state change

Pomo will execute an arbitrary command specified in the array argument `onEvent`
when the state changes.  The first element of this array should be the
executable to run while the remaining elements are space delimited arguments.
The new state will be exported as an environment variable `POMO_STATE` for this
command.  Possible state values are `RUNNING`, `PAUSED`, `BREAKING`, or
`COMPLETE`.

For example, to trigger a terminal bell when a session completes, add the
following to `config.json`:
```json
...
"onEvent": ["/bin/sh", "/path/to/script/my_script.sh"],
...
```
where the contents of `my_script.sh` are
```bash
#!/bin/sh

if [ "$POMO_STATE" == "COMPLETE" ] ; then
   echo -e '\a'
fi
```

See the `contrib` directory for user contributed scripts for use with `onEvent`.

## Integrations

By default pomo will setup a Unix socket and serve it's status there.

```bash
echo | socat stdio UNIX-CONNECT:$HOME/.pomo/pomo.sock | jq .
{
  "state": 1,
  "remaining": 1492000000000,
  "count": 0,
  "n_pomodoros": 4
}
```

Alternately by setting the `publish` flag to `true` it will publish it's status
to an existing socket.

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

#### [luastatus](https://github.com/shdown/luastatus)

Configured this bar by setting `publish` to `true`.

```lua
widget = {
	plugin = "unixsock",
	opts = {
		path = "pomo.sock",
		timeout = 2,
	},
	cb = function(t)
		local full_text
		local foreground = ""
		local background = ""
		if t.what == "line" then
			if string.match(t.line, "R") then
				-- green
				foreground = "#ffffff"
				background = "#307335"
			end
			if string.match(t.line, "B") or string.match(t.line, "P") or string.match(t.line, "C") then
				-- red
				foreground = "#ffffff"
				background = "ff8080"
			end
			return { full_text = t.line, background = background, foreground = foreground }
		elseif t.what == "timeout" then
			return { full_text = "-" }
		elseif t.what == "hello" then
			return { full_text = "-" }
		end
	end,
}

```


## Roadmap

  * Generate charts/burn down
  * ??

## Credits

 * [pomodoro technique](https://cirillocompany.de/pages/pomodoro-technique/book/)
 * [logo by rones](https://openclipart.org/detail/262421/tomato-by-rones)
 * [website generate by hugo](http://gohugo.io/)
 * [theme by calintat](https://github.com/calintat/minimal)
