package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	cli "github.com/jawher/mow.cli"

	pomo "github.com/kevinschoon/pomo/pkg/internal"
)

func maybe(err error) {
	if err != nil {
		fmt.Printf("Error:\n%s\n", err)
		os.Exit(1)
	}
}

func defaultConfigPath() string {
	u, err := user.Current()
	maybe(err)
	return path.Join(u.HomeDir, "/.pomo/config.json")
}

func parseRange(arg string) (int, int, error) {
	if strings.Contains(arg, ":") {
		split := strings.Split(arg, ":")
		start, err := strconv.ParseInt(split[0], 0, 64)
		if err != nil {
			return -1, -1, err
		}
		end, err := strconv.ParseInt(split[1], 0, 64)
		if err != nil {
			return -1, -1, err
		}
		return int(start), int(end), nil
	}
	n, err := strconv.ParseInt(arg, 0, 64)
	if err != nil {
		return -1, -1, err
	}
	return int(n), int(n), err
}

func start(config *pomo.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] MESSAGE"
		var (
			duration  = cmd.StringOpt("d duration", "25m", "duration of each stent")
			pomodoros = cmd.IntOpt("p pomodoros", 4, "number of pomodoros")
			message   = cmd.StringArg("MESSAGE", "", "descriptive name of the given task")
			tags      = cmd.StringsOpt("t tag", []string{}, "tags associated with this task")
		)
		cmd.Action = func() {
			parsed, err := time.ParseDuration(*duration)
			maybe(err)
			db, err := pomo.NewStore(config.DBPath)
			maybe(err)
			defer db.Close()
			task := &pomo.Task{
				Message:    *message,
				Tags:       *tags,
				NPomodoros: *pomodoros,
				Duration:   parsed,
			}
			maybe(db.With(func(tx *sql.Tx) error {
				id, err := db.CreateTask(tx, *task)
				if err != nil {
					return err
				}
				task.ID = id
				return nil
			}))
			runner, err := pomo.NewTaskRunner(task, config)
			maybe(err)
			server, err := pomo.NewServer(runner, config)
			maybe(err)
			server.Start()
			defer server.Stop()
			runner.Start()
			pomo.StartUI(runner)
		}
	}
}

func create(config *pomo.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] MESSAGE"
		var (
			duration  = cmd.StringOpt("d duration", "25m", "duration of each stent")
			pomodoros = cmd.IntOpt("p pomodoros", 4, "number of pomodoros")
			message   = cmd.StringArg("MESSAGE", "", "descriptive name of the given task")
			tags      = cmd.StringsOpt("t tag", []string{}, "tags associated with this task")
		)
		cmd.Action = func() {
			parsed, err := time.ParseDuration(*duration)
			maybe(err)
			db, err := pomo.NewStore(config.DBPath)
			maybe(err)
			defer db.Close()
			task := &pomo.Task{
				Message:    *message,
				Tags:       *tags,
				NPomodoros: *pomodoros,
				Duration:   parsed,
			}
			maybe(db.With(func(tx *sql.Tx) error {
				taskId, err := db.CreateTask(tx, *task)
				if err != nil {
					return err
				}
				fmt.Println(taskId)
				return nil
			}))
		}
	}
}

func begin(config *pomo.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] TASK_ID"
		var (
			taskId = cmd.IntArg("TASK_ID", -1, "ID of Pomodoro to begin")
		)

		cmd.Action = func() {
			db, err := pomo.NewStore(config.DBPath)
			maybe(err)
			defer db.Close()
			var task *pomo.Task
			maybe(db.With(func(tx *sql.Tx) error {
				read, err := db.ReadTask(tx, *taskId)
				if err != nil {
					return err
				}
				task = read
				err = db.DeletePomodoros(tx, *taskId)
				if err != nil {
					return err
				}
				task.Pomodoros = []*pomo.Pomodoro{}
				return nil
			}))
			runner, err := pomo.NewTaskRunner(task, config)
			maybe(err)
			server, err := pomo.NewServer(runner, config)
			maybe(err)
			server.Start()
			defer server.Stop()
			runner.Start()
			pomo.StartUI(runner)
		}
	}
}

func initialize(config *pomo.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		cmd.Action = func() {
			db, err := pomo.NewStore(config.DBPath)
			maybe(err)
			defer db.Close()
			maybe(pomo.InitDB(db))
		}
	}
}

func list(config *pomo.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		var (
			asJSON   = cmd.BoolOpt("json", false, "output task history as JSON")
			assend   = cmd.BoolOpt("assend", false, "sort tasks assending in age")
			all      = cmd.BoolOpt("a all", true, "output all tasks")
			limit    = cmd.IntOpt("n limit", 0, "limit the number of results by n")
			duration = cmd.StringOpt("d duration", "24h", "show tasks within this duration")
		)
		cmd.Action = func() {
			duration, err := time.ParseDuration(*duration)
			maybe(err)
			db, err := pomo.NewStore(config.DBPath)
			maybe(err)
			defer db.Close()
			maybe(db.With(func(tx *sql.Tx) error {
				tasks, err := db.ReadTasks(tx)
				maybe(err)
				if *assend {
					sort.Sort(sort.Reverse(pomo.ByID(tasks)))
				}
				if !*all {
					tasks = pomo.After(time.Now().Add(-duration), tasks)
				}
				if *limit > 0 && (len(tasks) > *limit) {
					tasks = tasks[0:*limit]
				}
				if *asJSON {
					maybe(json.NewEncoder(os.Stdout).Encode(tasks))
					return nil
				}
				maybe(err)
				pomo.SummerizeTasks(config, tasks)
				return nil
			}))
		}
	}
}

func _delete(config *pomo.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] [TASK_ID...]"
		cmd.LongDesc = `
delete one or more tasks by ID

## Examples:
# delete a single task
pomo delete 1
# delete a range of tasks (1 - 10)
pomo delete 1:10
# delete multiple tasks 5, 10, and 20
pomo delete 5 10 20
`
		var taskIDs = cmd.StringsArg("TASK_ID", nil, "task to delete")
		cmd.Action = func() {

			db, err := pomo.NewStore(config.DBPath)
			maybe(err)
			defer db.Close()
			maybe(db.With(func(tx *sql.Tx) error {
				for _, expr := range *taskIDs {
					start, end, err := parseRange(expr)
					if err != nil {
						return err
					}
					for i := start; i <= end; i++ {
						err := db.DeleteTask(tx, i)
						if err != nil {
							return err
						}
						fmt.Printf("deleted task %d\n", i)
					}
				}

				return nil
			}))
		}
	}
}

func _status(config *pomo.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		cmd.Action = func() {
			client, err := pomo.NewClient(config.SocketPath)
			if err != nil {
				fmt.Println(pomo.FormatStatus(pomo.Status{}))
				return
			}
			defer client.Close()
			status, err := client.Status()
			maybe(err)
			fmt.Println(pomo.FormatStatus(*status))
		}
	}
}

func _config(config *pomo.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(config))
		}
	}
}

func New(config *pomo.Config) *cli.Cli {
	app := cli.App("pomo", "Pomodoro CLI")
	app.LongDesc = "Pomo helps you track what you did, how long it took you to do it, and how much effort you expect it to take."
	app.Spec = "[OPTIONS]"
	var (
		path = app.StringOpt("p path", defaultConfigPath(), "path to the pomo config directory")
	)
	app.Before = func() {
		maybe(pomo.LoadConfig(*path, config))
	}
	app.Version("v version", pomo.Version)
	app.Command("start s", "start a new task", start(config))
	app.Command("init", "initialize the sqlite database", initialize(config))
	app.Command("config cf", "display the current configuration", _config(config))
	app.Command("create c", "create a new task without starting", create(config))
	app.Command("begin b", "begin requested pomodoro", begin(config))
	app.Command("list l", "list historical tasks", list(config))
	app.Command("delete d", "delete a stored task", _delete(config))
	app.Command("status st", "output the current status", _status(config))
	return app
}

func Run() { New(&pomo.Config{}).Run(os.Args) }
