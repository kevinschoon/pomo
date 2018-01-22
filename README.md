# 🍅 pomo

`pomo` is a simple CLI for using the [Pomodoro Technique](https://en.wikipedia.org/wiki/Pomodoro_Technique). There are [some](https://taskwarrior.org/) [amazing](https://todoist.com/) task management systems but `pomo` is more of a task *execution* or [timeboxing](https://en.wikipedia.org/wiki/Timeboxing) system. `pomo` helps you track what you did, how long it took you do it, and how long you expect it will take.

## Background

The Pomodoro Technique is simple and effective:

  * Decide on a task you want to accomplish
  * Break the task into timed intervals (pomodoros), [approx. 25 min]
  * After each pomodoro take a short break [approx. 3 - 5 min]
  * Once all pomodoros are completed take a longer break [approx 15 - 20 min]
  * Repeat

## Installation

`pomo` depends on the [libnotify](https://developer.gnome.org/libnotify/) client package, a notification [server](https://wiki.archlinux.org/index.php/Desktop_notifications#Notification_servers) (installed with most Linux desktop environments), and [SQLite](https://sqlite.org/).

Binaries are available in the [releases section](https://github.com/kevinschoon/pomo/releases) on github.

### Linux

#### Binaries

```
curl -L -o pomo https://github.com/kevinschoon/pomo/releases/download/0.1.0/pomo-0.1.0-linux 
echo f4587b566d135e05a6c1b1bec50fe3378f643f654319ca4662d5fe3aa590b8d2 pomo | sha256sum -c -
chmod +x pomo
./pomo -v
# Copy pomo to somewhere on your $PATH
```

#### Source

 ```
 go get github.com/kevinschoon/pomo
 cd $GOPATH/github.com/kevinschoon/pomo
 make
 ./bin/pomo
 ```

## Usage

```
# Initialize the SQLite database and state directory
pomo init
# Start a new task
# Add a tag "dev", allocate 2 pomodoros for 1 minute each
pomo start -t dev -p 2 -d 1m "Write Some Codes"
...
# List previous tasks
# pomo list
...
```

## Roadmap

  * Support OSX
  * Support Windows
  * Generate charts
  * Alternate notifiers
  * ??

## Credits

 * [pomodoro technique](https://cirillocompany.de/pages/pomodoro-technique/book/)
 * [logo by rones](https://openclipart.org/detail/262421/tomato-by-rones)
 * [website generate by hugo](http://gohugo.io/)
 * [theme by calintat](https://github.com/calintat/minimal)
