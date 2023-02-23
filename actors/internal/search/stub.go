package search

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/lejabque/software-design-2022/actors/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type stubServer struct {
	api.UnimplementedSearchServiceServer
	response *api.SearchResponse
}

func (s *stubServer) Search(ctx context.Context, req *api.SearchRequest) (*api.SearchResponse, error) {
	return s.response, nil
}

func NewStubServer(response *api.SearchResponse) *stubServer {
	return &stubServer{response: response}
}

func makeStubResponse(count int, endpoint string) *api.SearchResponse {
	var results []*api.SearchResult
	for i := 0; i < 5; i++ {
		results = append(results, &api.SearchResult{
			Url:   fmt.Sprintf("https://example%d.com", i),
			Title: fmt.Sprintf("Example %s", endpoint),
		})
	}
	return &api.SearchResponse{Results: results}
}

func RunStubServers() []string {
	var endpoints []string
	for i := 0; i < 3; i++ {
		port := 10080 + i
		endpoint := fmt.Sprintf("localhost:%d", port)
		endpoints = append(endpoints, endpoint)
		go func() {
			listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
			if err != nil {
				panic(err)
			}

			s := grpc.NewServer()

			api.RegisterSearchServiceServer(s, NewStubServer(makeStubResponse(5, endpoint)))
			reflection.Register(s)
			if err := s.Serve(listener); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
		}()
	}
	return endpoints
}
