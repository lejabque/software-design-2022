package vkstats

import (
	"context"
	"github.com/lejabque/software-design-2022/api_mocks/internal/vkapi"
	"time"
)

//go:generate mockery --name=vkClient  --exported --case underscore
type vkClient interface {
	SearchPosts(ctx context.Context, queryText string, startTime time.Time, endTime time.Time, startFrom string, count int) (vkapi.PostsResponse, error)
}

type Counter struct {
	client vkClient
}

func (c Counter) CountPostsDistribution(ctx context.Context, queryText string, hoursCnt int) ([]int64, error) {
	out := make([]int64, hoursCnt)
	now := time.Now()
	startFrom := ""
	for i := 0; i < hoursCnt; i++ {
		hourNum := hoursCnt - i
		startTime := now.Add(time.Hour * time.Duration(-hourNum))
		endTime := now.Add(time.Hour * time.Duration(-hourNum+1))
		posts, err := c.client.SearchPosts(ctx, queryText, startTime, endTime, startFrom, 0)
		if err != nil {
			return nil, err
		}
		out[i] = int64(posts.TotalCount)
		startFrom = posts.NextFrom
	}
	return out, nil
}
