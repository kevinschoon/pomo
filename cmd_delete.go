package main

import (
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
)

func deleteProject(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[--tag] [ID]"
		cmd.LongDesc = `
Delete a project by ID or when all tags are matched
        `
		var (
			projectID = cmd.IntArg("ID", -1, "project to delete")
			tags      = cmd.StringsOpt("tag", []string{}, "delete projects with matching tags")
		)
		cmd.Action = func() {
			if *projectID == 0 {
				maybe(fmt.Errorf("cannot delete root project"))
			}
			// TODO
			kv, err := parseTags(*tags)
			maybe(err)
			if len(kv) == 0 && *projectID == -1 {
				cmd.PrintLongHelp()
				os.Exit(1)
			}
			store, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer store.Close()
			maybe(store.With(func(s Store) error {
				project := &Project{ID: int64(*projectID)}
				err := s.ReadProject(project)
				if err != nil {
					return err
				}
				return s.DeleteProject(int64(*projectID))
			}))
		}
	}
}

func deleteTask(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[--tag] [ID]"
		var (
			taskID = cmd.IntArg("ID", -1, "task to delete")
			tags   = cmd.StringsOpt("tag", []string{}, "delete projects with matching tags")
		)
		cmd.Action = func() {
			// TODO
			kv, err := parseTags(*tags)
			maybe(err)
			if len(kv) == 0 && *taskID == -1 {
				cmd.PrintLongHelp()
				os.Exit(1)
			}
			store, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer store.Close()
			maybe(store.With(func(s Store) error {
				task := &Task{ID: int64(*taskID)}
				err := s.ReadTask(task)
				if err != nil {
					return err
				}
				return s.DeleteTask(int64(*taskID))
			}))
		}
	}
}

func deletePomodoro(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "ID"
		cmd.Action = func() {
			// TODO: need to update pomodoro store function
		}
	}
}

func _delete(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[-f --force] [ [PROJECT] | [PROJECT TASK] | [PROJECT TASK POMODORO]]"
		cmd.LongDesc = `
Delete one or more resources. Resources can be specified by relative index
position (index starts from 0), or the case of tasks or pomodoros, e.g.

pomo delete PROJECT_ID | PROJECT_ID [TASK_N] | PROJECT_ID TASK_N [POMODORO_N]

or by unique resource identifier e.g.

pomo delete task 2

        `

		cmd.Command("project t", "delete a project", deleteProject(config))
		cmd.Command("task t", "delete a task", deleteTask(config))
		cmd.Command("pomodoro po", "delete a pomodoro", deletePomodoro(config))

		var (
			force         = cmd.BoolOpt("force f", false, "do not prompt before deleting")
			projectID     = cmd.IntArg("PROJECT", 0, "project identifier")
			taskIndex     = cmd.IntArg("TASK", -1, "task index")
			pomodoroIndex = cmd.IntArg("POMODORO", -1, "pomodoro index")
		)

		cmd.Action = func() {
			if *projectID == 0 {
				maybe(fmt.Errorf("cannot delete root project"))
			}
			store, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer store.Close()
			maybe(store.With(func(s Store) error {

				project := &Project{
					ID: int64(*projectID),
				}
				err := s.ReadProject(project)
				if err != nil {
					return err
				}
				if *taskIndex == -1 && *pomodoroIndex == -1 {
					if !*force {
						fmt.Println("Are you sure you want to delete the following project: ")
						Tree(*project).Write(os.Stdout, 0, Tree(*project).MaxDepth() == 0)
						fmt.Println("Type YES to continue: ")
						err = promptConfirm("YES")
						if err != nil {
							return err
						}
						return s.DeleteProject(project.ID)
					}
				} else if *taskIndex >= 0 && *pomodoroIndex == -1 {
					if *taskIndex > len(project.Tasks)-1 {
						maybe(fmt.Errorf("no task available at index %d", *taskIndex))
					}
					fmt.Println("Are you sure you want to delete the following task: ")
					fmt.Println(project.Tasks[*taskIndex].Info())
					fmt.Println("Type YES to continue: ")
					err = promptConfirm("YES")
					if err != nil {
						return err
					}
					return s.DeleteTask(project.Tasks[*taskIndex].ID)

				} else if *taskIndex >= 0 && *pomodoroIndex >= 0 {
					if *taskIndex > len(project.Tasks)-1 {
						maybe(fmt.Errorf("no task available at index %d", *taskIndex))
					}
					if *pomodoroIndex > len(project.Tasks[*taskIndex].Pomodoros) {
						maybe(fmt.Errorf("no pomodoro available at index %d", *pomodoroIndex))
					}
					fmt.Println("Are you sure you want to delete the following pomodoro: ")
					fmt.Println(project.Tasks[*taskIndex].Pomodoros[*pomodoroIndex].Info(project.Tasks[*taskIndex].Duration))
					err = promptConfirm("YES")
					if err != nil {
						return err
					}
					return s.DeletePomodoros(project.Tasks[*taskIndex].ID, project.Tasks[*taskIndex].Pomodoros[*pomodoroIndex].ID)
				} else {
					maybe(fmt.Errorf("bad options: %d %d %d", *projectID, *taskIndex, *pomodoroIndex))
				}

				return nil

			}))
		}
	}
}
