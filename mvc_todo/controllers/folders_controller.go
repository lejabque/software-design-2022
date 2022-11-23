package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lejabque/software-design-2022/mvc_todo/database"
	"github.com/lejabque/software-design-2022/mvc_todo/views"
	"golang.org/x/exp/slices"
)

type foldersRepo interface {
	CreateFolder(ctx context.Context, folder *database.Folder) error
	ListFolders(ctx context.Context) ([]*database.Folder, error)
	DeleteFolder(ctx context.Context, folder string) error
}

type FolderController struct {
	folders foldersRepo
	views   views.View
}

func NewFolderController(folders foldersRepo, views views.View) *FolderController {
	return &FolderController{folders: folders, views: views}
}

func (c *FolderController) createFolder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	folder := ps.ByName("name")
	if folder == "" {
		return NewUserError("folder name is required")
	}
	err := c.folders.CreateFolder(r.Context(), &database.Folder{Name: folder})
	if err != nil {
		return err
	}
	return nil
}

func (c *FolderController) CreateFolder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := c.createFolder(w, r, ps)
	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
	} else {
		http.Redirect(w, r, fmt.Sprintf("/folders/%s/tasks", ps.ByName("name")), http.StatusSeeOther)
	}
}

func (c *FolderController) deleteFolder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	folder := ps.ByName("name")
	if folder == "" {
		return NewUserError("folder name is required")
	}
	err := c.folders.DeleteFolder(r.Context(), folder)
	if err != nil {
		return err
	}
	return nil
}

func (c *FolderController) DeleteFolder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	WrapHandler(c.deleteFolder)(w, r, ps)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *FolderController) listFolders(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	folders, err := c.folders.ListFolders(r.Context())
	if err != nil {
		return err
	}
	slices.SortFunc(folders, func(left, right *database.Folder) bool {
		return left.Name < right.Name
	})
	type RenderFolder struct {
		Name string
	}
	var data struct {
		Folders []*RenderFolder
	}
	for _, f := range folders {
		data.Folders = append(data.Folders, &RenderFolder{
			Name: f.Name,
		})
	}
	return c.views.Folders.Render(w, data)
}

func (c *FolderController) ListFolders(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	WrapHandler(c.listFolders)(w, r, ps)
}
