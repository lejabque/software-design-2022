package main

import (
	"log"
	"net"

	"github.com/lejabque/software-design-2022/actors/api"
	"github.com/lejabque/software-design-2022/actors/internal/search"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	endpoints := search.RunStubServers()
	searchConfig := search.SearchConfig{Endpoints: endpoints}
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	logger := zap.NewExample()
	defer logger.Sync()

	api.RegisterSearchServiceServer(s, search.NewServer(searchConfig, logger))
	reflection.Register(s)
	logger.Sugar().Infof("starting server on %s", listener.Addr().String())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
