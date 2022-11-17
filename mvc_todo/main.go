package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
	repo := database.NewTaskRepo(ydb)
	controller := controllers.NewTaskController(repo, views.NewView("views/layouts"))

	port := os.Args[2]
	http.HandleFunc("/", controller.Index)
	http.HandleFunc("/index", controller.Index)
	http.HandleFunc("/create", controller.CreateTask)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		panic(err)
	}
}
