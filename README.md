<p align="center"><img src="https://raw.githubusercontent.com/kevinschoon/pomo/master/www/static/demo.gif" alt="demo"/></p>

# 🍅 pomo

`pomo` is a simple CLI for using the [Pomodoro Technique](https://en.wikipedia.org/wiki/Pomodoro_Technique). There are [some](https://taskwarrior.org/) [amazing](https://todoist.com/) task management systems but `pomo` is more of a task *execution* or [timeboxing](https://en.wikipedia.org/wiki/Timeboxing) system. `pomo` helps you track what you did, how long it took you do it, and how much effort you expect it to take.

## Background

The Pomodoro Technique is simple and effective:

  * Decide on a task you want to accomplish
  * Break the task into timed intervals (pomodoros), [approx. 25 min]
  * After each pomodoro take a short break [approx. 3 - 5 min]
  * Once all pomodoros are completed take a longer break [approx 15 - 20 min]
  * Repeat

## Installation

### Dependencies

`pomo` depends on the [SQLite](https://sqlite.org) runtime, make sure it is installed before continuing. If running on Linux `pomo` will attempt to send notifications via [libnotify](https://developer.gnome.org/libnotify). 


### Binaries

Binaries are available in the [releases section](https://github.com/kevinschoon/pomo/releases) on github.

```
curl -L -o pomo https://github.com/kevinschoon/pomo/releases/download/0.2.0/pomo-0.2.0-linux 
echo f74156fdcdaec9574d64891b8fc76b7e5b0b2c00a446eb58e0fe4341fc83e9ad pomo | sha256sum -c -
chmod +x pomo
./pomo -v
# Copy pomo to somewhere on your $PATH
```

### Source

 ```
 go get github.com/kevinschoon/pomo
 cd $GOPATH/github.com/kevinschoon/pomo
 make
 ./bin/pomo
 ```

## Usage

Once `pomo` is installed you need to initialize it's database.

```
pomo init
```

Start a 4 pomodoro session at 25 minute intervals:
```
pomo start -t my-project "write some codes"
```

## Roadmap

  * Generate charts
  * Notifiers for OSX and Windows(?)
  * ??

## Credits

 * [pomodoro technique](https://cirillocompany.de/pages/pomodoro-technique/book/)
 * [logo by rones](https://openclipart.org/detail/262421/tomato-by-rones)
 * [website generate by hugo](http://gohugo.io/)
 * [theme by calintat](https://github.com/calintat/minimal)
