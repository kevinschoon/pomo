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
    INSERT INTO task (task_id, message, duration) VALUES (0, "root", 0) ON CONFLICT(task_id) DO UPDATE SET task_id = task_id;
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
		Where(sq.Eq{"parent_id": int64(0)}).
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
	err := s.ReadTask(root)
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

// CreateTask creates a new Task
func (s *SQLiteStore) CreateTask(task *pomo.Task) error {
	result, err := sq.
		Insert("task").
		Columns("parent_id", "message", "duration").
		Values(task.ParentID, task.Message, task.Duration).
		RunWith(s.tx).Exec()
	if err != nil {
		return err
	}
	taskId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	task.ID = taskId
	for _, key := range task.Tags.Keys() {
		_, err := sq.
			Insert("tag").
			Columns("task_id", "key", "value").
			Values(task.ID, key, task.Tags.Get(key)).
			RunWith(s.tx).Exec()
		if err != nil {
			return err
		}
	}
	for _, pomodoro := range task.Pomodoros {
		pomodoro.TaskID = taskId
		err = s.CreatePomodoro(pomodoro)
		if err != nil {
			return err
		}
	}
	task.ID = taskId
	return nil
}

// ReadTask reads a single task recursively updating
// and sibling tasks
func (s *SQLiteStore) ReadTask(task *pomo.Task) error {
	// special case when requesting the root task which
	// has no parentID and returns null, otherwise there are
	// no orphans.

	parentID := sql.NullInt64{}

	err := sq.Select("task_id", "parent_id", "message", "duration").
		From("task").
		Where(sq.Eq{"task_id": task.ID}).
		RunWith(s.tx).
		QueryRow().
		Scan(&task.ID, &parentID, &task.Message, &task.Duration)

	if err != nil {
		return err
	}

	task.ParentID = parentID.Int64

	task.Tags = tags.New()

	rows, err := sq.
		Select("key", "value").
		From("tag").
		Where(sq.Eq{"task_id": task.ID}).
		RunWith(s.tx).
		Query()

	if err != nil {
		return err
	}

	for rows.Next() {
		var (
			key, value string
		)
		err = rows.Scan(&key, &value)
		if err != nil {
			return err
		}
		task.Tags.Set(key, value)
	}

	pomodoros, err := s.ReadPomodoros(task.ID, -1)
	if err != nil {
		return err
	}

	task.Pomodoros = pomodoros

	subTasks, err := s.ReadTasks(task.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			subTasks = []*pomo.Task{}
		} else {
			return err
		}
	}

	task.Tasks = subTasks

	return nil
}

// ReadTasks reads multiple tasks matching
// the parentID recursively updating any
// sibling tasks
func (s *SQLiteStore) ReadTasks(parentID int64) ([]*pomo.Task, error) {
	var tasks []*pomo.Task
	query := sq.
		Select("task_id", "parent_id", "message", "duration").
		From("task").
		Where(sq.Eq{"parent_id": parentID})
	rows, err := query.
		RunWith(s.tx).
		Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		task := &pomo.Task{Pomodoros: []*pomo.Pomodoro{}}
		err = rows.Scan(&task.ID, &task.ParentID, &task.Message, &task.Duration)
		if err != nil {
			return nil, err
		}
		rows, err := sq.
			Select("key", "value").
			From("tag").
			Where(sq.Eq{"task_id": task.ID}).
			RunWith(s.tx).
			Query()

		if err != nil {
			return nil, err
		}

		task.Tags = tags.New()

		for rows.Next() {
			var (
				key, value string
			)
			err = rows.Scan(&key, &value)
			if err != nil {
				return nil, err
			}
			task.Tags.Set(key, value)
		}
		tasks = append(tasks, task)
	}
	for _, task := range tasks {
		pomodoros, err := s.ReadPomodoros(task.ID, -1)
		if err != nil {
			return nil, err
		}
		task.Pomodoros = pomodoros
	}
	for _, task := range tasks {
		subTasks, err := s.ReadTasks(task.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			return nil, err
		}
		task.Tasks = subTasks
	}
	return tasks, nil
}

// UpdateTask updates a single task
func (s *SQLiteStore) UpdateTask(task *pomo.Task) error {
	_, err := sq.
		Update("task").
		Where(sq.Eq{"task_id": task.ID}).
		Set("parent_id", task.ParentID).
		Set("duration", task.Duration).
		Set("message", task.Message).
		RunWith(s.tx).Exec()
	// TODO generalize tags
	_, err = sq.
		Delete("tag").
		Where(sq.Eq{"task_id": task.ID}).
		RunWith(s.tx).
		Exec()
	if err != nil {
		return err
	}
	for _, key := range task.Tags.Keys() {
		_, err := sq.
			Insert("tag").
			Columns("task_id", "key", "value").
			Values(task.ID, key, task.Tags.Get(key)).
			RunWith(s.tx).Exec()
		if err != nil {
			return err
		}
	}
	return err
}

// DeleteTask deletes a task with the given ID
func (s *SQLiteStore) DeleteTask(taskID int64) error {
	_, err := sq.
		Delete("task").
		Where(sq.Eq{"task_id": taskID}).
		RunWith(s.tx).Exec()
	return err
}

// CreatePomodoro creates a new pomodoro
func (s *SQLiteStore) CreatePomodoro(pomodoro *pomo.Pomodoro) error {
	result, err := sq.
		Insert("pomodoro").
		Columns("task_id", "start", "run_time", "pause_time").
		Values(pomodoro.TaskID, pomodoro.Start, pomodoro.RunTime, pomodoro.PauseTime).
		RunWith(s.tx).Exec()
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	pomodoro.ID = id
	return nil
}

// UpdatePomodoro updates a single Pomodoro
func (s *SQLiteStore) UpdatePomodoro(pomodoro *pomo.Pomodoro) error {
	_, err := sq.
		Update("pomodoro").
		Set("start", pomodoro.Start).
		Set("run_time", pomodoro.RunTime).
		Set("pause_time", pomodoro.PauseTime).
		Where(sq.Eq{"pomodoro_id": pomodoro.ID}).
		RunWith(s.tx).Exec()
	return err
}

// ReadPomodoros returns all pomodoros optionally matching
// taskID and pomodoroID if their IDs are greater than zero.
// To return all pomodoros:
// s.ReadPomodoros(tx, -1, -1)
func (s *SQLiteStore) ReadPomodoros(taskID, pomodoroID int64) ([]*pomo.Pomodoro, error) {
	var pomodoros []*pomo.Pomodoro
	query := sq.
		Select("pomodoro_id", "task_id", "start", "run_time", "pause_time").
		From("pomodoro")

	conditional := sq.Eq{"task_id": taskID}

	if pomodoroID > 0 {
		conditional["pomodoro_id"] = pomodoroID
	}
	if len(conditional) > 0 {
		query = query.Where(conditional)
	}
	rows, err := query.RunWith(s.tx).Query()
	if err != nil {
		return nil, err
	}
	var datetimeStr string
	for rows.Next() {
		pomodoro := &pomo.Pomodoro{}
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

// DeletePomodoros deletes all pomodoros associated with
// taskID or a single Pomodoro matching both the taskID
// and pomodoroID
func (s *SQLiteStore) DeletePomodoros(taskID, pomodoroID int64) error {
	conditional := sq.Eq{
		"task_id": taskID,
	}
	if pomodoroID > 0 {
		conditional["pomodoro_id"] = pomodoroID
	}
	_, err := sq.
		Delete("pomodoro").
		Where(conditional).
		RunWith(s.tx).Exec()
	return err
}
