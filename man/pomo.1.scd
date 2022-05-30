pomo(1)

# NAME

*Pomo* is a simple CLI for using the Pomodoro Technique.

# SYNOPSIS

*pomo* [OPTIONS] COMMAND [arg...]

# DESCRIPTION

*pomo* helps you track what you did, how long it took you to do it, 
and how much effort you expect it to take.

The Pomodoro Technique is simple and effective:

	** Decide on a task you want to accomplish
	** Break the task into timed intervals (pomodoros), [approx. 25 min]
	** After each pomodoro take a short break [approx. 3 - 5 min]
	** Once all pomodoros are completed take a longer break [approx 15 - 20 min]
	** Repeat

# SUBCOMMANDS 

See --help for the complete command usage

```
  start, s        start a new task
  init            initialize the sqlite database
  config, cf      display the current configuration
  create, c       create a new task without starting
  begin, b        begin requested pomodoro
  list, l         list historical tasks
  delete, d       delete a stored task
  status, st      output the current status

```

# CONFIGURATION

Pomo has a configuration file that is stored in *~/.config/pomo/config.json*.

```
{
	"colors": null,
	"dateTimeFmt": "2006-01-02 15:04",
	"publish": false,
	"publishJson": false,
	"publishSocketPath": ""
}
```

# EXAMPLES

## GETTING STARTED

```
# ensure your database has been initialized
pomo init
# run a new pomodoro
pomo start -t my-project "write some code"
# once finished view previously completed pomodoros
pomo list
```

# SEE ALSO

See the pomo source repository on Github at https://github.com/kevinschoon/pomo for complete documentation.

# AUTHORS

Written by Kevin Schoon <me@kevinschoon.com> with help from the open source
community.
