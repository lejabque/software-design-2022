package search

import (
	"context"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lejabque/software-design-2022/actors/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type searchChildActor struct {
	endpoint string
	logger   *zap.Logger
}

type searchChildResponse struct {
	Results []*api.SearchResult
}

func (*searchChildResponse) NotInfluenceReceiveTimeout() {}

func newSearchChildActor(endpoint string, logger *zap.Logger) actor.Actor {
	return &searchChildActor{endpoint: endpoint, logger: logger}
}

func (s *searchChildActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *api.SearchRequest:
		s.logger.Sugar().Infof("searching on %s", s.endpoint)
		reqCtx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		conn, err := grpc.Dial(s.endpoint, grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
		client := api.NewSearchServiceClient(conn)
		resp, err := client.Search(reqCtx, msg)
		if err != nil {
			panic(err)
		}
		res := &searchChildResponse{}
		for _, result := range resp.Results {
			result.Source = s.endpoint
			res.Results = append(res.Results, result)
		}
		ctx.Request(ctx.Parent(), res)
	}
}
