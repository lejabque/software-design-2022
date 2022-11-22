package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/lejabque/software-design-2022/mvc_todo/database"
	"github.com/lejabque/software-design-2022/mvc_todo/views"
	"golang.org/x/exp/slices"
)

type tasksRepo interface {
	CreateTask(ctx context.Context, task *database.Task) error
	GetFolderTasks(ctx context.Context, folder string) ([]*database.Task, error)
	GetTask(ctx context.Context, folder string, id uint64) (*database.Task, error)
	UpdateTask(ctx context.Context, task *database.Task) error
	CompleteTask(ctx context.Context, folder string, id uint64) error
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
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Priority:    database.PriorityFromString(r.FormValue("priority")),
		Deadline:    deadline,
	}, nil
}

func (c *TaskController) createTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	folder := ps.ByName("name")
	task, err := c.taskFromForm(r)
	if err != nil {
		return err
	}
	task.Folder = folder
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

func (c *TaskController) CreateTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	WrapHandler(c.createTask)(w, r, ps)
	// TODO: redirect?
}

func (c *TaskController) completeTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	folder := ps.ByName("name")
	if folder == "" {
		return NewUserError("folder name is required")
	}
	id := ps.ByName("id")
	if id == "" {
		return NewUserError("task id is required")
	}
	taskId, err := strconv.Atoi(id)
	if err != nil {
		return NewUserError("invalid id")
	}
	return c.tasks.CompleteTask(r.Context(), folder, uint64(taskId))
}

func (c *TaskController) CompleteTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	WrapHandler(c.completeTask)(w, r, ps)
}

func (c *TaskController) listTasks(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	tasks, err := c.tasks.GetFolderTasks(r.Context(), ps.ByName("name"))
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

func (c *TaskController) ListTasks(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	WrapHandler(c.listTasks)(w, r, ps)
}
