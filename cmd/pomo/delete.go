package main

import (
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/functional"
	"github.com/kevinschoon/pomo/pkg/store"
	"github.com/kevinschoon/pomo/pkg/tree"
)

func deleteTask(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] [ID]"
		var (
			taskID     = cmd.IntArg("ID", -1, "task to delete")
			filterArgs = cmd.StringsOpt("f filter", []string{}, "filters")
		)

		cmd.Action = func() {
			db, err := store.NewSQLiteStore(cfg.DBPath)
			maybe(err)
			defer db.Close()
			maybe(db.With(func(db store.Store) error {
				if *taskID > 0 {
					task := &pomo.Task{ID: int64(*taskID)}
					err := db.ReadTask(task)
					if err != nil {
						return err
					}
					return db.DeleteTask(int64(*taskID))
				}
				root := &pomo.Task{
					ID: int64(0),
				}
				err := db.ReadTask(root)
				if err != nil {
					return err
				}
				tasks := functional.FindMany(root.Tasks, functional.FiltersFromStrings(*filterArgs)...)
				fmt.Println("are you sure you want to delete the following tasks:")
				for _, subTask := range tasks {
					tree.Tree{Task: *subTask}.Write(os.Stdout, nil)
				}
				fmt.Println("type YES to confirm")
				err = promptConfirm("YES")
				if err != nil {
					return err
				}
				for _, subTask := range root.Tasks {
					err = db.DeleteTask(subTask.ID)
					if err != nil {
						return err
					}
				}
				return nil
			}))
		}
	}
}
