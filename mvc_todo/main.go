package main

import (
	"context"
	"os"

	"github.com/lejabque/software-design-2022/mvc_todo/database"
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

	err := repo.ResetTable(context.Background())
	if err != nil {
		panic(err)
	}
}
