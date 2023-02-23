package search

import (
	"context"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lejabque/software-design-2022/actors/api"
	"go.uber.org/zap"
)

type SearchConfig struct {
	Endpoints []string
}

type Server struct {
	api.UnimplementedSearchServiceServer
	actors *actor.ActorSystem
	config SearchConfig
	logger *zap.Logger
}

func NewServer(config SearchConfig, logger *zap.Logger) *Server {
	return &Server{
		actors: actor.NewActorSystem(),
		config: config,
		logger: logger,
	}
}

func (s *Server) Search(ctx context.Context, req *api.SearchRequest) (*api.SearchResponse, error) {
	props := actor.PropsFromFunc(NewSearchAggregatorActor(s.config, s.logger).Receive)
	pid := s.actors.Root.Spawn(props)
	timeout := 1 * time.Minute
	if d, ok := ctx.Deadline(); ok {
		timeout = d.Sub(time.Now())
	}
	resp, err := s.actors.Root.RequestFuture(pid, req, timeout).Result()
	if err != nil {
		return nil, err
	}
	return resp.(*api.SearchResponse), nil
}
