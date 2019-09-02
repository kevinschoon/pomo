package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
)

func deleteProject(config *Config) func(*cli.Cmd) {
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
			store, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer store.Close()
			maybe(store.With(func(s Store) error {
				if *projectID > 0 {
					project := &Project{ID: int64(*projectID)}
					err := s.ReadProject(project)
					if err != nil {
						return err
					}
					return s.DeleteProject(int64(*projectID))
				}
				root := &Project{
					ID: int64(0),
				}
				err := s.ReadProject(root)
				if err != nil {
					return err
				}
				projects := root.Children
				projects = FilterProjects(projects, ProjectFiltersFromStrings(*filterArgs)...)
				if len(projects) == 1 {
					fmt.Println("are you sure you want to delete the following project?")
					fmt.Println(projects[0].Info())
					for _, task := range projects[0].Tasks {
						fmt.Println(task.Info())
					}
					fmt.Println("type YES to confirm")
					maybe(promptConfirm("YES"))
					return s.DeleteProject(projects[0].ID)
				} else if len(projects) > 1 {
					return fmt.Errorf("too ambiguous, got %d results", len(projects))
				}
				return fmt.Errorf("no results")
			}))
		}
	}
}

func deleteTask(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] [ID]"
		var (
			taskID            = cmd.IntArg("ID", -1, "task to delete")
			projectFilterArgs = cmd.StringsOpt("p project", []string{}, "project filters")
			taskFilterArgs    = cmd.StringsOpt("t task", []string{}, "task filters")
		)

		cmd.Action = func() {
			store, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer store.Close()
			maybe(store.With(func(s Store) error {
				if *taskID > 0 {
					task := &Task{ID: int64(*taskID)}
					err := s.ReadTask(task)
					if err != nil {
						return err
					}
					return s.DeleteTask(int64(*taskID))
				}
				root := &Project{
					ID: int64(0),
				}
				err := s.ReadProject(root)
				if err != nil {
					return err
				}
				projects := root.Children
				projects = FilterProjects(projects, ProjectFiltersFromStrings(*projectFilterArgs)...)
				for _, project := range projects {
					ForEachMutate(project, func(p *Project) {
						p.Tasks = FilterTasks(p.Tasks, TaskFiltersFromStrings(*taskFilterArgs)...)
					})
				}
				projects = FilterProjects(projects, ProjectFilterSomeTasks())
				if len(projects) == 1 {
					if len(projects[0].Tasks) == 1 {
						fmt.Println("are you sure you want to delete the following task: ")
						fmt.Println(projects[0].Tasks[0].Info())
						fmt.Println("type YES to confirm")
						maybe(promptConfirm("YES"))
						return s.DeleteTask(projects[0].Tasks[0].ID)
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

func deletePomodoro(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
	}
}

func _delete(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Command("project t", "delete a project", deleteProject(config))
		cmd.Command("task t", "delete a task", deleteTask(config))
		cmd.Command("pomodoro po", "delete a pomodoro", deletePomodoro(config))
	}
}
