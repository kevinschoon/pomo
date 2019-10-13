package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	cli "github.com/jawher/mow.cli"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/store"
	"github.com/kevinschoon/pomo/pkg/tags"
)

type sliceValue struct {
	start int
	end   int
}

func parseSliceString(str string) (*sliceValue, error) {
	split := strings.Split(str, ":")
	if !(len(split) == 1 || len(split) == 2) {
		return nil, fmt.Errorf("bad slice: %s", str)
	}
	sv := &sliceValue{
		end: -1,
	}
	start, err := strconv.ParseInt(split[0], 0, 16)
	if err != nil {
		return nil, err
	}
	if start < 0 {
		return nil, fmt.Errorf("bad slice: %s", str)
	}
	sv.start = int(start)
	if len(split) == 2 {
		end, err := strconv.ParseInt(split[1], 0, 16)
		if err != nil {
			return nil, err
		}
		if end < 0 {
			return nil, fmt.Errorf("bad slice: %s", str)
		}
		sv.end = int(end)
	}
	return sv, nil
}

func edit(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] ID"
		cmd.LongDesc = `
Edit an existing task

Examples:

# Remove pomodoros 2-4
pomo edit --remove 2:3
        `
		var (
			taskID       = cmd.IntArg("ID", -1, "task identifier")
			parentID     = cmd.IntOpt("p parent", -1, "parent id")
			durationStr  = cmd.StringOpt("d duration", "", "pomodoro duration")
			addPomodoros = cmd.IntOpt("a add", 0, "add n pomodoros")
			message      = cmd.StringOpt("m message", "", "modify the task message")
			rmPomodoros  = cmd.StringOpt("r remove", "", "remove a subset of pomodoros between start:end")
			truncate     = cmd.BoolOpt("t truncate", false, "truncate the task to it's current runtime")
			done         = cmd.BoolOpt("D done", false, "mark the task as completed")
			addTags      = cmd.StringsOpt("T tag", []string{}, "add or modify an existing tag")
			rmTags       = cmd.StringsOpt("R removeTag", []string{}, "remove existing tags")
		)
		cmd.Action = func() {
			db, err := store.NewSQLiteStore(cfg.DBPath, cfg.Snapshots)
			maybe(err)
			defer db.Close()
			maybe(db.With(func(db store.Store) error {
				err := db.Snapshot()
				if err != nil {
					return err
				}
				task, err := db.ReadTask(int64(*taskID))
				if err != nil {
					return err
				}
				err = store.ReadAll(db, task)
				if err != nil {
					return err
				}
				var taskUpdated bool
				if *parentID != -1 {
					task.ParentID = int64(*parentID)
					taskUpdated = true
				}
				if *durationStr != "" {
					duration, err := time.ParseDuration(*durationStr)
					if err != nil {
						return err
					}
					task.Duration = duration
					taskUpdated = true
				}
				if *message != "" {
					task.Message = *message
					taskUpdated = true
				}
				if taskUpdated {
					err = db.UpdateTask(task)
					if err != nil {
						return err
					}
				}
				if *rmPomodoros != "" {
					sv, err := parseSliceString(*rmPomodoros)
					if err != nil {
						return err
					}
					subset := map[int64]bool{}
					for _, pomodoro := range task.Pomodoros[sv.start:sv.end] {
						subset[pomodoro.ID] = true
					}
					for _, pomodoro := range task.Pomodoros {
						if _, ok := subset[pomodoro.ID]; !ok {
							err := db.DeletePomodoro(pomodoro.ID)
							if err != nil {
								return err
							}
						}
					}
				}
				if *addPomodoros > 0 {
					for _, pomodoro := range pomo.NewPomodoros(*addPomodoros) {
						pomodoro.TaskID = task.ID
						_, err := db.WritePomodoro(pomodoro)
						if err != nil {
							return err
						}
					}
				}
				if *truncate {
					task.Duration = pomo.Truncate(*task)
					err := db.UpdateTask(task)
					if err != nil {
						return err
					}
					for _, pomodoro := range task.Pomodoros {
						pomodoro.Start = time.Time{}
						pomodoro.RunTime = task.Duration
						err := db.UpdatePomodoro(pomodoro)
						if err != nil {
							return err
						}
					}
				}
				if *done {
					for _, pomodoro := range task.Pomodoros {
						if pomodoro.Start.IsZero() {
							pomodoro.Start = time.Now()
						}
						pomodoro.RunTime += (task.Duration - pomodoro.RunTime)
						err := db.UpdatePomodoro(pomodoro)
						if err != nil {
							return err
						}
					}
				}
				if len(*addTags) > 0 {
					tgs, err := tags.FromKVs(*addTags)
					if err != nil {
						return err
					}
					err = db.WriteTags(tags.Merge(task.Tags, tgs))
					if err != nil {
						return err
					}
				}
				if len(*rmTags) > 0 {
					for _, key := range *rmTags {
						task.Tags.Delete(key)
					}
					err := db.WriteTags(task.Tags)
					if err != nil {
						return err
					}
				}
				return nil
			}))

		}
	}
}
