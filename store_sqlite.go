package main

import (
	"database/sql"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
)

// 2018-01-16 19:05:21.752851759+08:00
const datetimeFmt = "2006-01-02 15:04:05.999999999-07:00"

var _ Store = (*SQLiteStore)(nil)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &SQLiteStore{db: db}, nil
}

// With applies all of the given functions with
// a single transaction, rolling back on failure
// and commiting on success.
func (s SQLiteStore) With(fns ...func(tx *sql.Tx) error) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	for _, fn := range fns {
		err = fn(tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s SQLiteStore) CreateTask(tx *sql.Tx, task Task) (int64, error) {
	result, err := sq.
		Insert("task").
		Columns("message", "duration", "tags").
		Values(task.Message, task.Duration, strings.Join(task.Tags, ",")).
		RunWith(tx).Exec()
	if err != nil {
		return -1, err
	}
	taskId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return taskId, nil
}

func (s SQLiteStore) ReadTask(tx *sql.Tx, task *Task) error {
	var (
		tags string
	)
	err := sq.Select("task_id", "message", "duration", "tags").
		From("task").
		Where(sq.Eq{"task_id": task.ID}).
		RunWith(tx).
		QueryRow().
		Scan(&task.ID, &task.Message, &task.Duration, &tags)

	if err != nil {
		return err
	}

	// TODO: JSONB
	if tags != "" {
		task.Tags = strings.Split(tags, ",")
	}
	return nil
}

func (s SQLiteStore) UpdateTask(tx *sql.Tx, task Task) error {
	_, err := sq.
		Update("task").
		Set("duration", task.Duration).
		Set("message", task.Message).
		Set("tags", strings.Join(task.Tags, ",")).
		RunWith(tx).Exec()
	return err
}

func (s SQLiteStore) ReadTasks(tx *sql.Tx) ([]*Task, error) {
	var tasks []*Task
	rows, err := sq.
		Select("task_id", "message", "duration", "tags").
		From("task").
		RunWith(tx).Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var (
			tags string
		)
		task := &Task{Pomodoros: []*Pomodoro{}}
		err = rows.Scan(&task.ID, &task.Message, &task.Duration, &tags)
		if err != nil {
			return nil, err
		}
		if tags != "" {
			task.Tags = strings.Split(tags, ",")
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s SQLiteStore) DeleteTask(tx *sql.Tx, taskID int64) error {
	_, err := sq.
		Delete("task").
		Where(sq.Eq{"task_id": taskID}).
		RunWith(tx).Exec()
	return err
}

func (s SQLiteStore) CreatePomodoro(tx *sql.Tx, pomodoro Pomodoro) (int64, error) {
	result, err := sq.
		Insert("pomodoro").
		Columns("task_id", "start", "run_time", "pause_time").
		Values(pomodoro.TaskID, pomodoro.Start, pomodoro.RunTime, pomodoro.PauseTime).
		RunWith(tx).Exec()
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (s SQLiteStore) UpdatePomodoro(tx *sql.Tx, pomodoro Pomodoro) error {
	_, err := sq.
		Update("pomodoro").
		Set("start", pomodoro.Start).
		Set("run_time", pomodoro.RunTime).
		Set("pause_time", pomodoro.PauseTime).
		RunWith(tx).Exec()
	return err
}

// ReadPomodoros returns all pomodoros optionally matching
// taskID and pomodoroID if their IDs are greater than zero.
// To return all pomodoros:
// s.ReadPomodoros(tx, -1, -1)
func (s SQLiteStore) ReadPomodoros(tx *sql.Tx, taskID, pomodoroID int64) ([]*Pomodoro, error) {
	var pomodoros []*Pomodoro
	query := sq.
		Select("pomodoro_id", "task_id", "start", "run_time", "pause_time").
		From("pomodoro")
	conditional := sq.Eq{}
	if taskID > 0 {
		conditional["task_id"] = taskID
	}
	if pomodoroID > 0 {
		conditional["pomodoro_id"] = pomodoroID
	}
	if len(conditional) > 0 {
		query = query.Where(conditional)
	}
	rows, err := query.RunWith(tx).Query()
	if err != nil {
		return nil, err
	}
	var datetimeStr string
	for rows.Next() {
		pomodoro := &Pomodoro{}
		err = rows.Scan(
			&pomodoro.ID,
			&pomodoro.TaskID,
			&datetimeStr,
			&pomodoro.RunTime,
			&pomodoro.PauseTime,
		)
		if err != nil {
			return nil, err
		}
		start, _ := time.Parse(datetimeFmt, datetimeStr)
		pomodoro.Start = start
		pomodoros = append(pomodoros, pomodoro)
	}
	return pomodoros, nil
}

func (s SQLiteStore) DeletePomodoros(tx *sql.Tx, taskID, pomodoroID int64) error {
	conditional := sq.Eq{
		"task_id": taskID,
	}
	if pomodoroID > 0 {
		conditional["pomodoro_id"] = pomodoroID
	}
	_, err := sq.
		Delete("pomodoro").
		Where(conditional).
		RunWith(tx).Exec()
	return err
}

func (s SQLiteStore) Close() error { return s.db.Close() }

func initDB(db *SQLiteStore) error {
	// TODO Migrate
	stmt := `
    CREATE TABLE task (
    task_id INTEGER PRIMARY KEY,
	message TEXT,
	duration INTEGER,
	tags TEXT
    );
    CREATE TABLE pomodoro (
    pomodoro_id INTEGER PRIMARY KEY,
	task_id INTEGER,
	start DATETTIME,
	run_time INTEGER,
    pause_time INTEGER,
    FOREIGN KEY(task_id) REFERENCES task(task_id)
    );
    PRAGMA foreign_keys = ON;
    `
	_, err := db.db.Exec(stmt)
	return err
}
