package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"os/user"
)

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

func (s Store) AddTask(task Task) (int, error) {
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

func (s Store) AddRecord(taskID int, record Record) error {
	_, err := s.db.Exec(
		`INSERT INTO record (task_id, start, end) VALUES ($1, $2, $3)`,
		taskID,
		record.Start,
		record.End,
	)
	return err
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
