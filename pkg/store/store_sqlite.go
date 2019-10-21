package store

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/url"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/tags"
)

// 2018-01-16 19:05:21.752851759+08:00
const datetimeFmt = "2006-01-02 15:04:05.999999999-07:00"

var _ Store = (*SQLiteStore)(nil)

// SQLiteStore implements a Pomo store
// backed by SQLite
type SQLiteStore struct {
	db        *sql.DB
	tx        *sql.Tx
	snapshots int
}

// NewSQLiteStore returns a new SQLiteStore
func NewSQLiteStore(path string, history int) (*SQLiteStore, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	qs := &url.Values{}
	qs.Add("_fk", "yes")
	u.RawQuery = qs.Encode()
	db, err := sql.Open("sqlite3", u.String())
	if err != nil {
		return nil, err
	}
	return &SQLiteStore{db: db, snapshots: history}, nil
}

// Close closes the underlying SQLite connection
func (s *SQLiteStore) Close() error { return s.db.Close() }

// Init initalizes the SQLite database
func (s *SQLiteStore) Init() error {
	// TODO Migrate
	stmt := `
    CREATE TABLE task (
    task_id INTEGER PRIMARY KEY,
    parent_id INTEGER,
	message TEXT,
    notes TEXT,
	duration INTEGER,
    FOREIGN KEY(parent_id) REFERENCES task(task_id) ON DELETE CASCADE
    );
    CREATE TABLE pomodoro (
    pomodoro_id INTEGER PRIMARY KEY,
	task_id INTEGER,
	start DATETTIME,
	run_time INTEGER,
    pause_time INTEGER,
    FOREIGN KEY(task_id) REFERENCES task(task_id) ON DELETE CASCADE
    );
    CREATE TABLE tag (
    project_id INTEGER,
    task_id INTEGER,
    key TEXT,
    value TEXT,
    FOREIGN KEY(task_id) REFERENCES task(task_id) ON DELETE CASCADE
    );
    CREATE TABLE snapshot (
    snapshot_id INTEGER PRIMARY KEY,
    data JSON
    );
    PRAGMA foreign_keys = ON;
    INSERT INTO task (task_id, message, notes, duration) VALUES (0, "root", "", 0) ON CONFLICT(task_id) DO UPDATE SET task_id = task_id;
    `
	_, err := s.db.Exec(stmt)
	return err
}

// With executes a StoreFunc in the context
// of a single transaction
func (s *SQLiteStore) With(fn func(Store) error) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	s.tx = tx
	err = fn(s)
	if err != nil {
		tx.Rollback()
		s.tx = nil
		return err
	}
	s.tx = nil
	return tx.Commit()
}

// Reset completely empties the database an all
// associated data within it aside from the root
// task created on initialization

func (s *SQLiteStore) Reset() error {
	_, err := sq.Delete("task").
		RunWith(s.tx).
		Exec()
	if err != nil {
		return err
	}
	_, err = sq.
		Insert("task").
		Columns("task_id", "message", "duration").
		Values(0, "root", 0).
		Suffix("ON CONFLICT(task_id) DO UPDATE SET task_id = task_id").
		RunWith(s.tx).
		Exec()
	return err
}

// Snapshot saves the entire task tree in JSON format
// in the snapshot table.
// TODO: This currently duplicates the state everytime
// it is called which means the database file size
// will double on each successful transaction!
func (s *SQLiteStore) Snapshot() error {
	// snapshots are disabled
	if s.snapshots == -1 {
		return nil
	}
	// limit stored number of snapshots
	if s.snapshots > 0 {
		// TODO: believe this can be implemented with squirrel but
		// not immediately sure how to accomplish that
		_, err := s.tx.Exec(
			"delete from snapshot where snapshot_id = (select min(snapshot_id) from snapshot) and (select count(*) from snapshot) = ?;", s.snapshots)
		if err != nil {
			return err
		}
	}
	root := pomo.NewTask()
	root.ID = int64(0)
	err := ReadAll(s, root)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(nil)
	err = json.NewEncoder(buf).Encode(root)
	if err != nil {
		return err
	}
	_, err = sq.
		Insert("snapshot").
		Columns("data").
		Values(buf.Bytes()).
		RunWith(s.tx).Exec()
	return err
}

// Revert reverts the database to the given snapshot_id
func (s *SQLiteStore) Revert(id int, task *pomo.Task) error {
	buf := bytes.NewBuffer(nil)
	switch {
	// if id is zero return the most recent snapshot
	case id == 0:
		var data string
		err := sq.
			Select("snapshot_id", "data").
			From("snapshot").
			OrderBy("snapshot_id desc").
			Limit(1).
			RunWith(s.tx).
			QueryRow().
			Scan(&sql.NullInt64{}, &data)
		if err != nil {
			return err
		}
		buf.WriteString(data)
	// if id is less than zero take offset from the tail
	case id < 0:
		var data string
		err := sq.
			Select("snapshot_id", "data").
			From("snapshot").
			OrderBy("snapshot_id desc").
			Limit(1).
			Offset(uint64(-id)).
			RunWith(s.tx).
			QueryRow().
			Scan(&sql.NullInt64{}, &data)
		if err != nil {
			return err
		}
		buf.WriteString(data)
	// if id is greater than zero take the reset by index
	case id > 0:
		var data string
		err := sq.
			Select("snapshot_id", "data").
			From("snapshot").
			Where(sq.Eq{"snapshot_id": id}).
			RunWith(s.tx).
			QueryRow().
			Scan(&sql.NullInt64{}, &data)
		if err != nil {
			return err
		}
		buf.WriteString(data)
	}
	return json.NewDecoder(buf).Decode(task)
}

func (s *SQLiteStore) Search(opts SearchOptions) ([]*pomo.Task, error) {
	var (
		// find many
		statement = sq.
				Select("task.task_id").
				From("task").
				LeftJoin("tag on task.task_id = tag.task_id")

		conditions []sq.Sqlizer
		results    []*pomo.Task
	)
	if opts.ParentID > 0 {
		conditions = append(conditions, sq.Eq{"parent_id": opts.ParentID})
	}

	for _, value := range opts.Messages {
		conditions = append(conditions, sq.Like{"message": value})
	}

	for _, value := range opts.Notes {
		conditions = append(conditions, sq.Like{"notes": value})
	}

	if opts.Tags != nil && opts.Tags.Len() > 0 {
		keys := opts.Tags.Keys()
		var values []string
		conditions = append(conditions, sq.Eq{"key": keys})
		for _, key := range keys {
			value := opts.Tags.Get(key)
			if value != "" {
				values = append(values, value)
			}
		}
		if len(values) > 0 {
			conditions = append(conditions, sq.Eq{"value": values})
		}
	}

	if len(conditions) > 0 {
		if opts.MatchAny {
			statement = statement.Where(sq.Or(conditions))
		} else {
			statement = statement.Where(sq.And(conditions))
		}
	}

	rows, err := statement.
		RunWith(s.tx).
		Query()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var taskID int64
		err = rows.Scan(&taskID)
		if err != nil {
			return nil, err
		}
		task, err := s.ReadTask(taskID)
		if err != nil {
			return nil, err
		}
		results = append(results, task)
	}

	return results, nil
}

func (s *SQLiteStore) ReadTask(taskID int64) (*pomo.Task, error) {
	task := &pomo.Task{}

	notesValue := sql.NullString{}
	parentIDValue := sql.NullInt64{}
	err := sq.
		Select("task_id", "parent_id", "message", "notes", "duration").
		From("task").
		Where(sq.Eq{"task_id": taskID}).
		RunWith(s.tx).
		QueryRow().
		Scan(&task.ID, &parentIDValue, &task.Message, &notesValue, &task.Duration)
	if err != nil {
		return nil, err
	}
	task.ParentID = parentIDValue.Int64
	task.Notes = notesValue.String
	return task, nil
}

func (s *SQLiteStore) ReadTasks(taskID int64, parentID int64) ([]*pomo.Task, error) {
	var (
		statement = sq.
				Select("task_id", "parent_id", "message", "notes", "duration").
				From("task")
		conditions []sq.Sqlizer
		results    []*pomo.Task
	)

	if taskID >= 0 {
		conditions = append(conditions, sq.Eq{"task_id": taskID})
	}

	if parentID >= 0 {
		conditions = append(conditions, sq.Eq{"parent_id": parentID})
	}

	if len(conditions) > 0 {
		statement = statement.Where(sq.And(conditions))
	}

	rows, err := statement.
		RunWith(s.tx).
		Query()

	if err != nil {
		return nil, err
	}

	parentIDValue := sql.NullInt64{}
	notesValue := sql.NullString{}

	for rows.Next() {
		task := &pomo.Task{}
		err := rows.Scan(&task.ID, &parentIDValue, &task.Message, &notesValue, &task.Duration)
		if err != nil {
			return nil, err
		}
		task.ParentID = parentIDValue.Int64
		task.Notes = notesValue.String
		results = append(results, task)
	}

	return results, nil
}

// DeleteTask deletes a task with the given ID
func (s *SQLiteStore) DeleteTask(taskID int64) error {
	_, err := sq.
		Delete("task").
		Where(sq.Eq{"task_id": taskID}).
		RunWith(s.tx).Exec()
	return err
}

func (s *SQLiteStore) UpdateTask(task *pomo.Task) error {
	// upsert
	_, err := sq.
		Update("task").
		Where(sq.Eq{"task_id": task.ID}).
		Set("parent_id", task.ParentID).
		Set("duration", task.Duration).
		Set("message", task.Message).
		Set("notes", task.Notes).
		RunWith(s.tx).Exec()

	return err
}

func (s *SQLiteStore) WriteTask(task *pomo.Task) (int64, error) {

	result, err := sq.
		Insert("task").
		Columns("parent_id", "message", "notes", "duration").
		Values(task.ParentID, task.Message, task.Notes, task.Duration).
		RunWith(s.tx).
		Exec()
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	// TODO
	task.ID = id
	return id, nil
}

func (s *SQLiteStore) ReadTags(taskID int64) (*tags.Tags, error) {
	var (
		statement = sq.
				Select("key", "value").
				From("tag")
		conditions []sq.Sqlizer
		results    = tags.New()
	)

	results.TaskID = taskID

	if taskID > 0 {
		conditions = append(conditions, sq.Eq{"task_id": taskID})
	}

	if len(conditions) > 0 {
		statement = statement.Where(sq.And(conditions))
	}

	rows, err := statement.
		RunWith(s.tx).
		Query()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			key, value string
		)
		err = rows.Scan(&key, &value)
		if err != nil {
			return nil, err
		}
		results.Set(key, value)
	}

	return results, nil
}

func (s *SQLiteStore) WriteTags(kvs *tags.Tags) error {

	for _, key := range kvs.Keys() {
		_, err := sq.
			Insert("tag").
			Columns("task_id", "key", "value").
			Values(kvs.TaskID, key, kvs.Get(key)).
			RunWith(s.tx).Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SQLiteStore) DeleteTags(id int64) error {
	_, err := sq.
		Delete("tag").
		Where(sq.Eq{"task_id": id}).
		RunWith(s.tx).
		Exec()
	return err
}

func (s *SQLiteStore) ReadPomodoros(pomodoroID int64, taskID int64) ([]*pomo.Pomodoro, error) {

	var (
		statement = sq.
				Select("pomodoro_id", "task_id", "start", "run_time", "pause_time").
				From("pomodoro")

		conditions []sq.Sqlizer
		results    []*pomo.Pomodoro
	)

	if (pomodoroID) >= 0 {
		conditions = append(conditions, sq.Eq{"pomodoro_id": pomodoroID})
	}

	if (taskID) >= 0 {
		conditions = append(conditions, sq.Eq{"task_id": taskID})
	}

	if len(conditions) > 0 {
		statement = statement.Where(sq.And(conditions))
	}

	rows, err := statement.
		RunWith(s.tx).
		Query()

	if err != nil {
		return nil, err
	}

	for rows.Next() {

		pomodoro := &pomo.Pomodoro{}

		var dateTimeStr string

		err := rows.Scan(
			&pomodoro.ID,
			&pomodoro.TaskID,
			&dateTimeStr,
			&pomodoro.RunTime,
			&pomodoro.PauseTime,
		)

		if err != nil {
			return nil, err
		}

		// TODO: store in unix time

		dt, _ := time.Parse(datetimeFmt, dateTimeStr)
		pomodoro.Start = dt

		results = append(results, pomodoro)
	}

	return results, nil
}

func (s *SQLiteStore) UpdatePomodoro(pomodoro *pomo.Pomodoro) error {
	// upsert
	_, err := sq.
		Update("pomodoro").
		Set("start", pomodoro.Start).
		Set("run_time", pomodoro.RunTime).
		Set("pause_time", pomodoro.PauseTime).
		Where(sq.Eq{"pomodoro_id": pomodoro.ID}).
		RunWith(s.tx).Exec()
	return err
}

func (s *SQLiteStore) WritePomodoro(pomodoro *pomo.Pomodoro) (int64, error) {
	result, err := sq.
		Insert("pomodoro").
		Columns("task_id", "start", "run_time", "pause_time").
		Values(pomodoro.TaskID, pomodoro.Start, pomodoro.RunTime, pomodoro.PauseTime).
		RunWith(s.tx).Exec()
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	// TODO
	pomodoro.ID = id
	return id, nil
}

func (s *SQLiteStore) DeletePomodoro(id int64) error {
	_, err := sq.
		Delete("pomodoro").
		Where(sq.Eq{"pomodoro_id": id}).
		RunWith(s.tx).
		Exec()
	return err
}
