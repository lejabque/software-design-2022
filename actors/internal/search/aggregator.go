package search

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lejabque/software-design-2022/actors/api"
	"go.uber.org/zap"
)

const ReceiveTimeout = 5 * time.Second

type searchAggregatorActor struct {
	searchConfig SearchConfig
	logger       *zap.Logger
	results      []*searchChildResponse
	errors       []*api.SearchError
	sender       *actor.PID
}

func NewSearchAggregatorActor(config SearchConfig, logger *zap.Logger) actor.Actor {
	return &searchAggregatorActor{searchConfig: config, logger: logger}
}

func (s *searchAggregatorActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *api.SearchRequest:
		s.handleSearchRequest(ctx, msg)
	case *searchChildResponse:
		s.handleSearchResponse(ctx, msg)
	case *actor.ReceiveTimeout:
		s.logger.Sugar().Infof("timeout received")
		s.respondAndStop(ctx)
	}
}

func (s *searchAggregatorActor) stopChildren(ctx actor.Context) {
	for _, child := range ctx.Children() {
		ctx.Stop(child)
	}
	ctx.Poison(ctx.Self())
}

func (s *searchAggregatorActor) respondAndStop(ctx actor.Context) {
	defer s.stopChildren(ctx)
	resp := &api.SearchResponse{}
	for _, result := range s.results {
		resp.Results = append(resp.Results, result.Results...)
	}
	ctx.Send(s.sender, resp)
}

func (s *searchAggregatorActor) handleSearchRequest(ctx actor.Context, msg *api.SearchRequest) {
	if s.sender == nil {
		s.sender = ctx.Sender()
		ctx.SetReceiveTimeout(ReceiveTimeout)
		s.logger.Sugar().Infof("handling search request: %s", msg.Query)
		for _, endpoint := range s.searchConfig.Endpoints {
			props := actor.PropsFromFunc(newSearchChildActor(endpoint, s.logger).Receive)
			pid := ctx.Spawn(props)
			ctx.Forward(pid)
		}
	} else {
		s.logger.Sugar().Infof("ignoring search request because it's not the first one")
	}
}

func (s *searchAggregatorActor) handleSearchResponse(ctx actor.Context, msg *searchChildResponse) {
	s.results = append(s.results, msg)
	if len(s.results) == len(s.searchConfig.Endpoints) {
		s.respondAndStop(ctx)
	}
}
