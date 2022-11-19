package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/lejabque/software-design-2022/mvc_todo/database"
	"github.com/lejabque/software-design-2022/mvc_todo/views"
	"golang.org/x/exp/slices"
)

type tasksRepo interface {
	CreateTask(ctx context.Context, task *database.Task) error
	GetFolderTasks(ctx context.Context, folder string) ([]*database.Task, error)
	GetTask(ctx context.Context, folder string, id uint64) (*database.Task, error)
	UpdateTask(ctx context.Context, task *database.Task) error
	DeleteTask(ctx context.Context, folder string, id uint64) error
}

type TaskController struct {
	tasks tasksRepo
	views views.View
}

func NewTaskController(tasks tasksRepo, views views.View) *TaskController {
	return &TaskController{tasks: tasks, views: views}
}

func (c *TaskController) validateTask(task *database.Task) error {
	if task.Title == "" {
		return NewUserError("title is required")
	}
	// todo: implement tasks folders
	// if task.Folder == "" {
	// 	return NewUserError("folder is required")
	// }
	return nil
}

func (c *TaskController) taskFromForm(r *http.Request) (*database.Task, error) {
	var deadline time.Time
	if deadlineStr := r.FormValue("deadline"); deadlineStr != "" {
		var err error
		deadline, err = time.Parse("2006-01-02", deadlineStr)
		if err != nil {
			return nil, err
		}
	}
	return &database.Task{
		Folder:      r.FormValue("folder"),
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Priority:    database.PriorityFromString(r.FormValue("priority")),
		Deadline:    deadline,
	}, nil
}

func (c *TaskController) createTask(w http.ResponseWriter, r *http.Request) error {
	task, err := c.taskFromForm(r)
	if err != nil {
		return err
	}
	err = c.validateTask(task)
	if err != nil {
		return err
	}
	err = c.tasks.CreateTask(r.Context(), task)
	if err != nil {
		return err
	}
	return nil
}

func (c *TaskController) CreateTask(w http.ResponseWriter, r *http.Request) {
	err := c.createTask(w, r)
	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
	}
	http.Redirect(w, r, "/index", http.StatusSeeOther)
}

func (c *TaskController) index(w http.ResponseWriter, r *http.Request) error {
	tasks, err := c.tasks.GetFolderTasks(r.Context(), r.FormValue("folder"))
	if err != nil {
		return err
	}
	type RenderTask struct {
		Title       string
		Description string
		Priority    string
	}
	var data struct {
		Tasks []*RenderTask
	}
	// sort tasks by priority, asc
	slices.SortFunc(tasks, func(left, right *database.Task) bool {
		return left.Priority > right.Priority
	})
	for _, task := range tasks {
		data.Tasks = append(data.Tasks, &RenderTask{
			Title:       task.Title,
			Description: task.Description,
			Priority:    database.PriorityToString(task.Priority),
		})
	}
	return c.views.Index.Render(w, data)
}

func (c *TaskController) Index(w http.ResponseWriter, r *http.Request) {
	err := c.index(w, r)
	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
	}
}
