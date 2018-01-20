package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
	"time"
)

// 2018-01-16 19:05:21.752851759+08:00
const datetimeFmt = "2006-01-02 15:04:05.999999999-07:00"

type Store struct {
	db *sql.DB
}

func NewStore(path string) (*Store, error) {
	os.Mkdir(path, 0755)
	db, err := sql.Open("sqlite3", path+"/pomo.db")
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s Store) CreateTask(task Task) (int, error) {
	var taskID int
	tx, err := s.db.Begin()
	if err != nil {
		return -1, err
	}
	_, err = tx.Exec("INSERT INTO task (message,tags) VALUES ($1,$2)", task.Message, strings.Join(task.Tags, ","))
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	err = tx.QueryRow("SELECT last_insert_rowid() FROM task").Scan(&taskID)
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	return taskID, tx.Commit()
}

func (s Store) CreatePomodoro(taskID int, pomodoro Pomodoro) error {
	_, err := s.db.Exec(
		`INSERT INTO pomodoro (task_id, start, end) VALUES ($1, $2, $3)`,
		taskID,
		pomodoro.Start,
		pomodoro.End,
	)
	return err
}

func (s Store) ReadTasks() ([]*Task, error) {
	rows, err := s.db.Query(`SELECT rowid,message,tags FROM task`)
	if err != nil {
		return nil, err
	}
	tasks := []*Task{}
	for rows.Next() {
		var tags string
		task := &Task{Pomodoros: []*Pomodoro{}}
		err = rows.Scan(&task.ID, &task.Message, &tags)
		if err != nil {
			return nil, err
		}
		if tags != "" {
			task.Tags = strings.Split(tags, ",")
		}
		pomodoros, err := s.ReadPomodoros(task.ID)
		if err != nil {
			return nil, err
		}
		for _, pomodoro := range pomodoros {
			task.Pomodoros = append(task.Pomodoros, pomodoro)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s Store) ReadPomodoros(taskID int) ([]*Pomodoro, error) {
	rows, err := s.db.Query(`SELECT start,end FROM pomodoro WHERE task_id = $1`, &taskID)
	if err != nil {
		return nil, err
	}
	pomodoros := []*Pomodoro{}
	for rows.Next() {
		var (
			startStr string
			endStr   string
		)
		pomodoro := &Pomodoro{}
		err = rows.Scan(&startStr, &endStr)
		if err != nil {
			return nil, err
		}
		start, _ := time.Parse(datetimeFmt, startStr)
		end, _ := time.Parse(datetimeFmt, endStr)
		pomodoro.Start = start
		pomodoro.End = end
		pomodoros = append(pomodoros, pomodoro)
	}
	return pomodoros, nil
}

func (s Store) DeleteTask(taskID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM task WHERE rowid = $1", &taskID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("DELETE FROM record WHERE task_id = $1", &taskID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s Store) Close() error { return s.db.Close() }

func initDB(db *Store) error {
	stmt := `
    CREATE TABLE task (
	message TEXT,
	tags TEXT
    );
    CREATE TABLE pomodoro (
	task_id INTEGER,
	start DATETTIME,
	end DATETTIME
    );
    `
	_, err := db.db.Exec(stmt)
	return err
}
