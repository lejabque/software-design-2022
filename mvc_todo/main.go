package main

import (
	"net/http"
	"os"

	"github.com/lejabque/software-design-2022/mvc_todo/controllers"
	"github.com/lejabque/software-design-2022/mvc_todo/database"
	"github.com/lejabque/software-design-2022/mvc_todo/views"
)

func main() {
	// TODO: add cmd params
	saKeyPath := os.Getenv("SA_KEY_PATH")
	if saKeyPath == "" {
		panic("SA_KEY_PATH is not set")
	}
	// TODO: move to json config
	cfg := database.YdbConfig{
		Endpoint: "ydb.serverless.yandexcloud.net:2135",
		Database: "/ru-central1/b1gtehrpagsafsspa168/etn09q1egvnpnl0s3vhu",
	}
	ydb := database.NewYdbClient(cfg, saKeyPath)
	repo := database.NewTaskRepo(ydb)
	controller := controllers.NewTaskController(repo, views.NewView("views/layouts"))

	// bootstrap, err := template.ParseFiles(layoutFiles()...)
	// if err != nil {
	// 	panic(err)
	// }

	http.HandleFunc("/", controller.Index)
	http.HandleFunc("/index", controller.Index)
	http.HandleFunc("/create", controller.CreateTask)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
