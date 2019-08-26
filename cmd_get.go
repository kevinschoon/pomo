package main

import (
	"encoding/json"
	"os"
	"sort"

	cli "github.com/jawher/mow.cli"
)

func get(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		var (
			ascend = cmd.BoolOpt("a ascend", false, "sort tasks in ascending order")
			limit  = cmd.IntOpt("l limit", 0, "limit returned tasks")
			asJson = cmd.BoolOpt("json", false, "write as json")
		)
		cmd.Action = func() {
			store, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer store.Close()
			var tasks []*Task
			maybe(store.With(func(s Store) error {
				t, err := s.ReadTasks(-1)
				if err != nil {
					return err
				}
				tasks = t
				return nil
			}))
			if *ascend {
				sort.Sort(sort.Reverse(ByID(tasks)))
			}
			if *limit > 0 && (len(tasks) > *limit) {
				tasks = tasks[0:*limit]
			}
			if *asJson {
				maybe(json.NewEncoder(os.Stdout).Encode(tasks))
			} else {
				summerizeTasks(config, tasks)
			}
		}
	}
}
