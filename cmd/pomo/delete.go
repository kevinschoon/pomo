package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/filter"
	"github.com/kevinschoon/pomo/pkg/store"
)

func deleteProject(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] [ID]"
		cmd.LongDesc = `
Delete a project by ID or when all tags are matched
        `
		var (
			projectID  = cmd.IntArg("ID", -1, "project to delete")
			filterArgs = cmd.StringsOpt("f filter", []string{}, "project filters")
		)
		cmd.Action = func() {

			if *projectID == 0 {
				maybe(fmt.Errorf("cannot delete root project"))
			}
			db, err := store.NewSQLiteStore(cfg.DBPath)
			maybe(err)
			defer db.Close()
			maybe(db.With(func(db store.Store) error {
				if *projectID > 0 {
					project := &pomo.Project{ID: int64(*projectID)}
					err := db.ReadProject(project)
					if err != nil {
						return err
					}
					return db.DeleteProject(int64(*projectID))
				}
				root := &pomo.Project{
					ID: int64(0),
				}
				err := db.ReadProject(root)
				if err != nil {
					return err
				}
				projects := root.Children
				projects = filter.FilterProjects(projects, filter.ProjectFiltersFromStrings(*filterArgs)...)
				if len(projects) == 1 {
					fmt.Println("are you sure you want to delete the following project?")
					fmt.Println(projects[0].Info())
					for _, task := range projects[0].Tasks {
						fmt.Println(task.Info())
					}
					fmt.Println("type YES to confirm")
					maybe(promptConfirm("YES"))
					return db.DeleteProject(projects[0].ID)
				} else if len(projects) > 1 {
					return fmt.Errorf("too ambiguous, got %d results", len(projects))
				}
				return fmt.Errorf("no results")
			}))
		}
	}
}

func deleteTask(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] [ID]"
		var (
			taskID            = cmd.IntArg("ID", -1, "task to delete")
			projectFilterArgs = cmd.StringsOpt("p project", []string{}, "project filters")
			taskFilterArgs    = cmd.StringsOpt("t task", []string{}, "task filters")
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
				root := &pomo.Project{
					ID: int64(0),
				}
				err := db.ReadProject(root)
				if err != nil {
					return err
				}
				projects := root.Children
				projects = filter.FilterProjects(projects, filter.ProjectFiltersFromStrings(*projectFilterArgs)...)
				for _, project := range projects {
					pomo.ForEachMutate(project, func(p *pomo.Project) {
						p.Tasks = filter.FilterTasks(p.Tasks, filter.TaskFiltersFromStrings(*taskFilterArgs)...)
					})
				}
				projects = filter.FilterProjects(projects, filter.ProjectFilterSomeTasks())
				if len(projects) == 1 {
					if len(projects[0].Tasks) == 1 {
						fmt.Println("are you sure you want to delete the following task: ")
						fmt.Println(projects[0].Tasks[0].Info())
						fmt.Println("type YES to confirm")
						maybe(promptConfirm("YES"))
						return db.DeleteTask(projects[0].Tasks[0].ID)
					} else if len(projects[0].Tasks) > 1 {
						return fmt.Errorf("too ambiguous, got %d tasks of project %s", len(projects[0].Tasks), projects[0].Title)
					}
				} else if len(projects) > 1 {
					return fmt.Errorf("too ambiguous, got %d projects", len(projects))
				}
				return fmt.Errorf("no results")
			}))
		}
	}
}

func deletePomodoro(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
	}
}

func _delete(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Command("project t", "delete a project", deleteProject(cfg))
		cmd.Command("task t", "delete a task", deleteTask(cfg))
		cmd.Command("pomodoro po", "delete a pomodoro", deletePomodoro(cfg))
	}
}
