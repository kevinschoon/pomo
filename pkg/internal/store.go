package pomo

import (
	"database/sql"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// 2018-01-16 19:05:21.752851759+08:00
const datetimeFmt = "2006-01-02 15:04:05.999999999-07:00"

type StoreFunc func(tx *sql.Tx) error

type Store struct {
	db *sql.DB
}

func NewStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

// With applies all of the given functions with
// a single transaction, rolling back on failure
// and commiting on success.
func (s Store) With(fns ...func(tx *sql.Tx) error) error {
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

func (s Store) CreateTask(tx *sql.Tx, task Task) (int, error) {
	var taskID int
	_, err := tx.Exec(
		"INSERT INTO task (message,pomodoros,duration,tags) VALUES ($1,$2,$3,$4)",
		task.Message, task.NPomodoros, task.Duration.String(), strings.Join(task.Tags, ","))
	if err != nil {
		return -1, err
	}
	err = tx.QueryRow("SELECT last_insert_rowid() FROM task").Scan(&taskID)
	if err != nil {
		return -1, err
	}
	err = tx.QueryRow("SELECT last_insert_rowid() FROM task").Scan(&taskID)
	if err != nil {
		return -1, err
	}
	return taskID, nil
}

func (s Store) ReadTasks(tx *sql.Tx) ([]*Task, error) {
	rows, err := tx.Query(`SELECT rowid,message,pomodoros,duration,tags FROM task`)
	if err != nil {
		return nil, err
	}
	tasks := []*Task{}
	for rows.Next() {
		var (
			tags        string
			strDuration string
		)
		task := &Task{Pomodoros: []*Pomodoro{}}
		err = rows.Scan(&task.ID, &task.Message, &task.NPomodoros, &strDuration, &tags)
		if err != nil {
			return nil, err
		}
		duration, _ := time.ParseDuration(strDuration)
		task.Duration = duration
		if tags != "" {
			task.Tags = strings.Split(tags, ",")
		}
		pomodoros, err := s.ReadPomodoros(tx, task.ID)
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

func (s Store) DeleteTask(tx *sql.Tx, taskID int) error {
	_, err := tx.Exec("DELETE FROM task WHERE rowid = $1", &taskID)
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM pomodoro WHERE task_id = $1", &taskID)
	if err != nil {
		return err
	}
	return nil
}

func (s Store) ReadTask(tx *sql.Tx, taskID int) (*Task, error) {
	task := &Task{}
	var (
		tags        string
		strDuration string
	)
	err := tx.QueryRow(`SELECT rowid,message,pomodoros,duration,tags FROM task WHERE rowid = $1`, &taskID).
		Scan(&task.ID, &task.Message, &task.NPomodoros, &strDuration, &tags)
	if err != nil {
		return nil, err
	}
	duration, _ := time.ParseDuration(strDuration)
	task.Duration = duration
	if tags != "" {
		task.Tags = strings.Split(tags, ",")
	}
	return task, nil
}

func (s Store) CreatePomodoro(tx *sql.Tx, taskID int, pomodoro Pomodoro) error {
	_, err := tx.Exec(
		`INSERT INTO pomodoro (task_id, start, end) VALUES ($1, $2, $3)`,
		taskID,
		pomodoro.Start,
		pomodoro.End,
	)
	return err
}

func (s Store) ReadPomodoros(tx *sql.Tx, taskID int) ([]*Pomodoro, error) {
	rows, err := tx.Query(`SELECT start,end FROM pomodoro WHERE task_id = $1`, &taskID)
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

func (s Store) DeletePomodoros(tx *sql.Tx, taskID int) error {
	_, err := tx.Exec("DELETE FROM pomodoro WHERE task_id = $1", &taskID)
	return err
}

func (s Store) Close() error { return s.db.Close() }

func InitDB(db *Store) error {
	stmt := `
    CREATE TABLE task (
	message TEXT,
	pomodoros INTEGER,
	duration TEXT,
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
