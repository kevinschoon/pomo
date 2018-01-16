package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"os/user"
	"time"
)

// 2018-01-16 19:05:21.752851759+08:00
const datetimeFmt = "2006-01-02 15:04:05.999999999-07:00"

func defaultDBPath() string {
	u, err := user.Current()
	maybe(err)
	return u.HomeDir + "/.pomo"
}

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
	_, err = tx.Exec("INSERT INTO task (name) VALUES ($1)", task.Name)
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

func (s Store) CreateRecord(taskID int, record Record) error {
	_, err := s.db.Exec(
		`INSERT INTO record (task_id, start, end) VALUES ($1, $2, $3)`,
		taskID,
		record.Start,
		record.End,
	)
	return err
}

func (s Store) ReadTasks() ([]*Task, error) {
	rows, err := s.db.Query(`SELECT rowid,name FROM task`)
	if err != nil {
		return nil, err
	}
	tasks := []*Task{}
	for rows.Next() {
		task := &Task{Records: []*Record{}}
		err = rows.Scan(&task.ID, &task.Name)
		if err != nil {
			return nil, err
		}
		records, err := s.ReadRecords(task.ID)
		if err != nil {
			return nil, err
		}
		for _, record := range records {
			task.Records = append(task.Records, record)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s Store) ReadRecords(taskID int) ([]*Record, error) {
	rows, err := s.db.Query(`SELECT start,end FROM record WHERE task_id = $1`, &taskID)
	if err != nil {
		return nil, err
	}
	records := []*Record{}
	for rows.Next() {
		var (
			startStr string
			endStr   string
		)
		record := &Record{}
		err = rows.Scan(&startStr, &endStr)
		if err != nil {
			return nil, err
		}
		start, _ := time.Parse(datetimeFmt, startStr)
		end, _ := time.Parse(datetimeFmt, endStr)
		record.Start = start
		record.End = end
		records = append(records, record)
	}
	return records, nil
}

func (s Store) Close() error { return s.db.Close() }

func initDB(db *Store) error {
	stmt := `
    CREATE TABLE task (
	name TEXT
    );
    CREATE TABLE record (
	task_id INTEGER,
	start DATETTIME,
	end DATETTIME
    );
    `
	_, err := db.db.Exec(stmt)
	return err
}
