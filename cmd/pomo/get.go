package main

import (
	"encoding/json"
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config"
	ttemplate "github.com/kevinschoon/pomo/pkg/display/template/task"
	"github.com/kevinschoon/pomo/pkg/display/tree"
	"github.com/kevinschoon/pomo/pkg/store"
	"github.com/kevinschoon/pomo/pkg/tags"
	"strconv"
)

func get(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] [ID]"
		cmd.LongDesc = `
Examples:

# Output all tasks as a tree
pomo get
`
		var (
			// search options
			taskIDStr  = cmd.StringArg("ID", "", "task id")
			parentID   = cmd.IntOpt("p parent", int(cfg.CurrentRoot), "show tasks below the specified parent")
			strFilters = cmd.StringsOpt("m message", []string{}, "string filters")
			tagFilters = cmd.StringsOpt("t tag", []string{}, "tag filters")
			// display options
			//flatten       = cmd.BoolOpt("f flatten", false, "flatten all projects to one level")
			showPomodoros = cmd.BoolOpt("P pomodoros", true, "show status of each pomodoro")
			//recent        = cmd.BoolOpt("r recent", true, "sort by most recently modified tasks")
			//ascend        = cmd.BoolOpt("a ascend", false, "sort from oldest to newest")
		)
		cmd.Action = func() {
			db, err := store.NewSQLiteStore(cfg.DBPath, -1)
			maybe(err)
			defer db.Close()
			tgs, err := tags.FromKVs(*tagFilters)
			maybe(err)
			var results []*pomo.Task
			maybe(db.With(func(db store.Store) error {
				if *taskIDStr != "" || len(os.Args) == 2 {
					var taskID int64
					if len(os.Args) == 2 {
						taskID = int64(0)
					} else {
						parsed, err := strconv.ParseInt(*taskIDStr, 0, 64)
						if err != nil {
							return err
						}
						taskID = parsed
					}
					root, err := db.ReadTask(taskID)
					if err != nil {
						return err
					}
					err = store.ReadAll(db, root)
					if err != nil {
						return err
					}
					results = append(results, root)
					return nil
				}
				options := &store.SearchOptions{
					ParentID: int64(*parentID),
					Messages: *strFilters,
					Tags:     tgs,
				}
				tasks, err := db.Search(*options)
				if err != nil {
					return err
				}
				for _, task := range tasks {
					err := store.ReadAll(db, task)
					if err != nil {
						return err
					}
					results = append(results, task)
				}

				return nil
			}))
			if cfg.JSON {
				maybe(json.NewEncoder(os.Stdout).Encode(results))
			} else {
				for _, result := range results {
					if pomo.Depth(*result) == 1 {
						templater := ttemplate.NewTemplater(ttemplate.Options{Template: ttemplate.DefaultTemplate})
						for _, task := range result.Tasks {
							if pomo.Complete(*task) {
								cfg.Colors.Tertiary.Println(templater(*task))
							} else {
								cfg.Colors.Primary.Println(templater(*task))
							}
						}
					} else {
						t := tree.New(*result, *showPomodoros)
						t.Colors = *cfg.Colors
						t.TaskTemplater = ttemplate.NewTemplater(ttemplate.Options{
							Template: ttemplate.DefaultTemplate,
						})
						fmt.Println(t.String())
					}
				}

			}
		}
	}
}
