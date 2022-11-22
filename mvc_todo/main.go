package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/lejabque/software-design-2022/mvc_todo/controllers"
	"github.com/lejabque/software-design-2022/mvc_todo/database"
	"github.com/lejabque/software-design-2022/mvc_todo/views"
)

type mvcConfig struct {
	Ydb database.YdbConfig `json:"ydb"`
}

func readConfig(path string, cfg interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	return decoder.Decode(cfg)
}

func main() {
	saKeyPath := os.Getenv("SA_KEY_PATH")
	if saKeyPath == "" {
		panic("SA_KEY_PATH is not set")
	}
	// read YdbConfig from json file in path from cmd argument
	var cfg mvcConfig
	cfgPath := os.Args[1]
	if err := readConfig(cfgPath, &cfg); err != nil {
		panic(err)
	}

	ydb := database.NewYdbClient(cfg.Ydb, saKeyPath)
	tasksRepo := database.NewTaskRepo(ydb)
	foldersRepo := database.NewFoldersRepo(ydb)
	views := views.NewView("views/layouts")
	tasksController := controllers.NewTaskController(tasksRepo, views)
	foldersController := controllers.NewFolderController(foldersRepo, views)

	port := os.Args[2]
	router := httprouter.New()
	router.GET("/", foldersController.ListFolders)
	router.GET("/folders", foldersController.ListFolders)
	router.POST("/folders/:name/create", foldersController.CreateFolder)
	router.POST("/folders/:name/delete", foldersController.DeleteFolder)

	router.GET("/folders/:name/tasks", tasksController.ListTasks)
	router.GET("/folders/:name/tasks/:id/create", tasksController.CreateTask)
	router.GET("/folders/:name/tasks/:id/complete", tasksController.CompleteTask)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		panic(err)
	}

	// http.HandleFunc("/", foldersController.Index)
	// http.HandleFunc("/index", foldersController.Index)
	// http.HandleFunc("/create", controller.CreateTask)
	// if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
	// 	panic(err)
	// }
}
