package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/lejabque/software-design-2022/mvc_todo/database"
	"github.com/lejabque/software-design-2022/mvc_todo/views"
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
	if task.Folder == "" {
		return NewUserError("folder is required")
	}
	return nil
}

func (c *TaskController) createTask(w http.ResponseWriter, r *http.Request) error {
	var task database.Task
	err := json.Unmarshal([]byte(r.FormValue("task")), &task)
	if err != nil {
		return NewUserError("invalid task: %s", err)
	}
	err = c.validateTask(&task)
	if err != nil {
		return err
	}
	err = c.tasks.CreateTask(r.Context(), &task)
	if err != nil {
		return err
	}
	// TODO: render ?
	return c.views.New.Render(w, nil)
}

func (c *TaskController) CreateTask(w http.ResponseWriter, r *http.Request) {
	err := c.createTask(w, r)
	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
	}
}

func (c *TaskController) index(w http.ResponseWriter, r *http.Request) error {
	tasks, err := c.tasks.GetFolderTasks(r.Context(), r.FormValue("folder"))
	if err != nil {
		return err
	}
	return c.views.Index.Render(w, tasks)
}

func (c *TaskController) Index(w http.ResponseWriter, r *http.Request) {
	err := c.index(w, r)
	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
	}
}
