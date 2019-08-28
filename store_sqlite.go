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
	tx *sql.Tx
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Close() error { return s.db.Close() }

func (s *SQLiteStore) Init() error {
	// TODO Migrate
	stmt := `
    CREATE TABLE project (
    project_id INTEGER PRIMARY KEY,
    parent_id INTEGER,
    title TEXT,
    FOREIGN KEY(parent_id) REFERENCES project(project_id) ON DELETE CASCADE
    );
    CREATE TABLE task (
    task_id INTEGER PRIMARY KEY,
    project_id INTEGER,
	message TEXT,
	duration INTEGER,
	tags TEXT,
    FOREIGN KEY(project_id) REFERENCES project(project_id) ON DELETE CASCADE
    );
    CREATE TABLE pomodoro (
    pomodoro_id INTEGER PRIMARY KEY,
	task_id INTEGER,
	start DATETTIME,
	run_time INTEGER,
    pause_time INTEGER,
    FOREIGN KEY(task_id) REFERENCES task(task_id) ON DELETE CASCADE
    );
    PRAGMA foreign_keys = ON;
    INSERT INTO project (project_id, title) VALUES (0, "root") ON CONFLICT(project_id) DO UPDATE SET project_id = project_id;
    `
	_, err := s.db.Exec(stmt)
	return err
}

func (s *SQLiteStore) With(fn func(Store) error) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	s.tx = tx
	err = fn(s)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// CreateProject creates all projects and child projects recursively, each task
// and pomodoro are also created.
func (s *SQLiteStore) CreateProject(project *Project) error {
	result, err := sq.
		Insert("project").
		Columns("parent_id", "title").
		Values(project.ParentID, project.Title).
		RunWith(s.tx).Exec()
	if err != nil {
		return err
	}
	projectID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	project.ID = projectID
	for _, task := range project.Tasks {
		task.ProjectID = project.ID
		err = s.CreateTask(task)
		if err != nil {
			return err
		}
	}
	for _, child := range project.Children {
		child.ParentID = project.ID
		err = s.CreateProject(child)
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadProject returns the associated project, child projects,
// tasks, and pomodoros recursively.
func (s *SQLiteStore) ReadProject(project *Project) error {

	// special case when requesting the root project which
	// has no parentID and returns null, otherwise there are
	// no orphans.
	parentID := sql.NullInt64{}
	err := sq.
		Select("project_id", "parent_id", "title").
		From("project").
		Where(sq.Eq{"project_id": project.ID}).
		RunWith(s.tx).
		QueryRow().
		Scan(&project.ID, &parentID, &project.Title)

	if err != nil {
		return err
	}

	project.ParentID = parentID.Int64
	tasks, err := s.ReadTasks(project.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			tasks = []*Task{}
		} else {
			return err
		}
	}
	project.Tasks = tasks
	children, err := s.ReadProjects(project.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			children = []*Project{}
		} else {
			return err
		}
	}
	project.Children = children
	return nil
}

func (s *SQLiteStore) ReadProjects(parentID int64) ([]*Project, error) {
	var projects []*Project
	rows, err := sq.
		Select("project_id", "parent_id", "title").
		From("project").
		Where(sq.Eq{"parent_id": parentID}).
		RunWith(s.tx).Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		project := &Project{}
		err = rows.Scan(&project.ID, &project.ParentID, &project.Title)
		if err != nil {
			return nil, err
		}
		tasks, err := s.ReadTasks(project.ID)
		if err != nil {
			return nil, err
		}
		project.Tasks = tasks
		projects = append(projects, project)
	}
	for _, project := range projects {
		children, err := s.ReadProjects(project.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			return nil, err
		}
		project.Children = children
	}
	return projects, nil
}

// UpdateProject updates the title and parent association of the project
// it does not modify tasks, pomodoros, or child projects.
func (s *SQLiteStore) UpdateProject(project *Project) error {
	_, err := sq.
		Update("project").
		Set("title", project.Title).
		Set("parent_id", project.ParentID).
		RunWith(s.tx).Exec()
	return err
}

// DeleteProject deletes the given project ID which causes all
// decendant projects, tasks, and pomodoros to be deleted.
func (s *SQLiteStore) DeleteProject(projectID int64) error {
	_, err := sq.
		Delete("project").
		Where(sq.Eq{"project_id": projectID}).
		RunWith(s.tx).Exec()
	return err
}

func (s *SQLiteStore) CreateTask(task *Task) error {
	result, err := sq.
		Insert("task").
		Columns("project_id", "message", "duration", "tags").
		Values(task.ProjectID, task.Message, task.Duration, strings.Join(task.Tags, ",")).
		RunWith(s.tx).Exec()
	if err != nil {
		return err
	}
	taskId, err := result.LastInsertId()
	if err != nil {
		return err
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

func (s *SQLiteStore) ReadTask(task *Task) error {
	var (
		tags string
	)
	err := sq.Select("task_id", "project_id", "message", "duration", "tags").
		From("task").
		Where(sq.Eq{"task_id": task.ID}).
		RunWith(s.tx).
		QueryRow().
		Scan(&task.ID, &task.ProjectID, &task.Message, &task.Duration, &tags)

	if err != nil {
		return err
	}

	pomodoros, err := s.ReadPomodoros(task.ID, -1)
	if err != nil {
		return err
	}

	task.Pomodoros = pomodoros

	// TODO: JSONB
	if tags != "" {
		task.Tags = strings.Split(tags, ",")
	}
	return nil
}

func (s *SQLiteStore) ReadTasks(projectID int64) ([]*Task, error) {
	var tasks []*Task
	query := sq.
		Select("task_id", "project_id", "message", "duration", "tags").
		From("task")
	if projectID >= 0 {
		query = query.
			Where(sq.Eq{"project_id": projectID})
	}
	rows, err := query.
		RunWith(s.tx).
		Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var (
			tags string
		)
		task := &Task{Pomodoros: []*Pomodoro{}}
		err = rows.Scan(&task.ID, &task.ProjectID, &task.Message, &task.Duration, &tags)
		if err != nil {
			return nil, err
		}
		if tags != "" {
			task.Tags = strings.Split(tags, ",")
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
	return tasks, nil
}

func (s *SQLiteStore) UpdateTask(task *Task) error {
	_, err := sq.
		Update("task").
		Set("duration", task.Duration).
		Set("message", task.Message).
		Set("tags", strings.Join(task.Tags, ",")).
		RunWith(s.tx).Exec()
	return err
}

func (s *SQLiteStore) DeleteTask(taskID int64) error {
	_, err := sq.
		Delete("task").
		Where(sq.Eq{"task_id": taskID}).
		RunWith(s.tx).Exec()
	return err
}

func (s *SQLiteStore) CreatePomodoro(pomodoro *Pomodoro) error {
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

func (s *SQLiteStore) UpdatePomodoro(pomodoro *Pomodoro) error {
	_, err := sq.
		Update("pomodoro").
		Set("start", pomodoro.Start).
		Set("run_time", pomodoro.RunTime).
		Set("pause_time", pomodoro.PauseTime).
		RunWith(s.tx).Exec()
	return err
}

// ReadPomodoros returns all pomodoros optionally matching
// taskID and pomodoroID if their IDs are greater than zero.
// To return all pomodoros:
// s.ReadPomodoros(tx, -1, -1)
func (s *SQLiteStore) ReadPomodoros(taskID, pomodoroID int64) ([]*Pomodoro, error) {
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
	rows, err := query.RunWith(s.tx).Query()
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
