package main

import (
	"context"
	"fmt"
	"github.com/lejabque/software-design-2022/api_mocks/vkapi"
	"github.com/lejabque/software-design-2022/api_mocks/vkstats"
	"net/http"
	"os"
	"time"
)

func main() {
	token := os.Getenv("VK_TOKEN")
	if token == "" {
		panic("VK_TOKEN is not set")
	}
	client := vkstats.Counter{Client: vkapi.NewVKClient(token, "https://api.vk.com", &http.Client{})}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	query := "#test"
	res, err := client.CountPostsDistribution(ctx, query, 12)
	if err != nil {
		panic(err)
	}
	for i, v := range res {
		fmt.Printf("posts count for %d hour ago for query \"%s\": %d\n", i, query, v)
	}
}
