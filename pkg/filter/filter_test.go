package filter_test

import (
	"testing"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/filter"
	"github.com/kevinschoon/pomo/pkg/tags"
)

func getTestProject() *pomo.Project {
	return &pomo.Project{
		ID:   int64(1),
		Tags: tags.New(),
		Children: []*pomo.Project{
			&pomo.Project{
				ID:    int64(2),
				Title: "empty container",
				Tags:  tags.New(),
				Children: []*pomo.Project{
					&pomo.Project{
						ID:    int64(3),
						Title: "sub-project",
						Tasks: []*pomo.Task{
							&pomo.Task{
								ID:      int64(1),
								Message: "sub-sub-task",
								Tags:    tags.New(),
							},
						},
						Tags: tags.New(),
					},
				},
			},
			&pomo.Project{
				ID:    int64(4),
				Title: "Some Random Side Project",
				Tasks: []*pomo.Task{
					&pomo.Task{
						ID:      int64(2),
						Message: "write some code",
						Tags: tags.FromMap(
							[]string{"coding"},
							map[string]string{"coding": ""},
						),
					},
					&pomo.Task{
						ID:      int64(3),
						Message: "write some other code",
						Tags: tags.FromMap(
							[]string{"coding"},
							map[string]string{"coding": ""},
						),
					},
				},
				Tags: tags.FromMap(
					[]string{"miscProject"},
					map[string]string{"miscProject": "fuu"},
				),
			},
			&pomo.Project{
				ID:    int64(5),
				Title: "Some Random Research Project",
				Tasks: []*pomo.Task{
					&pomo.Task{
						ID:      int64(4),
						Message: "read a book",
						Tags: tags.FromMap(
							[]string{"fuu"},
							map[string]string{"fuu": ""},
						),
					},
					&pomo.Task{
						ID:      int64(5),
						Message: "read another book",
						Tags: tags.FromMap(
							[]string{"bar"},
							map[string]string{"bar": ""},
						),
					},
				},
				Tags: tags.FromMap(
					[]string{"research"},
					map[string]string{"research": ""},
				),
			},
		},
	}
}

func TestFilterFindOneProject(t *testing.T) {
	project := getTestProject()
	filters := filter.Filters{
		TaskFilters: nil,
		ProjectFilters: []filter.ProjectFilter{
			filter.ProjectFilterByName("Some Random Side Project"),
			filter.ProjectFilterByTag("miscProject", "fuu"),
		},
	}
	result := filter.FindOne(*project, filters)
	if err := result.Error(); err != nil {
		t.Fatal(err)
	}
	if result.Project() == nil {
		t.Fatal("should have returned a project")
	}
	if result.Project().Title != "Some Random Side Project" {
		t.Fatal("returned the wrong task")
	}
}

func TestFilterFindOneTask(t *testing.T) {

	project := getTestProject()

	filters := filter.Filters{
		TaskFilters: []filter.TaskFilter{
			filter.TaskFilterByName("read a book"),
			filter.TaskFilterByTag("fuu", ""),
		},
	}
	result := filter.FindOne(*project, filters)
	if err := result.Error(); err != nil {
		t.Fatal(err)
	}
	if result.Task() == nil {
		t.Fatal("should have returned a task")
	}
	if result.Task().Message != "read a book" {
		t.Fatal("returned the wrong task")
	}

	filters = filter.Filters{
		TaskFilters: []filter.TaskFilter{
			filter.TaskFilterByName("read"),   // ambiguous name
			filter.TaskFilterByTag("fuu", ""), // matching tag
		},
	}

	result = filter.FindOne(*project, filters)
	if err := result.Error(); err != nil {
		t.Fatal(err)
	}

	if result.Task() == nil {
		t.Fatal("should have returned a task")
	}

	if result.Task().Message != "read a book" {
		t.Fatal("returned the wrong task")
	}

	filters = filter.Filters{
		TaskFilters: []filter.TaskFilter{
			filter.TaskFilterByName("read"), // will return 2
		},
	}

	result = filter.FindOne(*project, filters)

	if result.Error() != filter.ErrTooManyResults {
		t.Fatal("should have too many results")
	}
}

func TestFilterReduce(t *testing.T) {
	project := getTestProject()
	filters := filter.Filters{
		TaskFilters: []filter.TaskFilter{
			filter.TaskFilterByTag("bar", ""), // matching tag
		},
	}
	project = filter.Reduce(project, filters)
	if len(project.Children) != 1 {
		t.Fatalf("should have one child project, got %d", len(project.Children))
	}
	if len(project.Children[0].Tasks) != 1 {
		t.Fatal("should have one matching task")
	}
	t.Log(project)
}
