package search

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/lejabque/software-design-2022/actors/api"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func TestSucessRequest(t *testing.T) {
	endpoints := RunStubServers()
	searchConfig := SearchConfig{Endpoints: endpoints}
	logger := zap.NewExample()
	defer logger.Sync()

	server := NewServer(searchConfig, logger)
	ctx := context.Background()
	req := &api.SearchRequest{Query: "test"}
	resp, err := server.Search(ctx, req)

	assert.Nil(t, err)
	assert.Equal(t, 15, len(resp.Results))
}

type errorServer struct {
	api.UnimplementedSearchServiceServer
}

func (s *errorServer) Search(ctx context.Context, req *api.SearchRequest) (*api.SearchResponse, error) {
	return &api.SearchResponse{}, fmt.Errorf("error")
}

type sleepingServer struct {
	api.UnimplementedSearchServiceServer
}

func (s *sleepingServer) Search(ctx context.Context, req *api.SearchRequest) (*api.SearchResponse, error) {
	time.Sleep(3 * ReceiveTimeout)
	return &api.SearchResponse{}, nil
}

func runServer(port int, server api.SearchServiceServer) string {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	api.RegisterSearchServiceServer(s, server)
	reflection.Register(s)
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	return listener.Addr().String()
}

func TestAllFailedRequest(t *testing.T) {
	endpoints := []string{runServer(8081, &errorServer{}), runServer(8082, &errorServer{})}
	searchConfig := SearchConfig{Endpoints: endpoints}
	logger := zap.NewExample()
	defer logger.Sync()

	server := NewServer(searchConfig, logger)
	req := &api.SearchRequest{Query: "fail"}
	resp, err := server.Search(context.Background(), req)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 0, len(resp.Results))
}

func TestOneFailedRequest(t *testing.T) {
	endpoints := []string{runServer(8083, &errorServer{}), runServer(8084, &stubServer{response: makeStubResponse(5, "localhost:8084")})}
	searchConfig := SearchConfig{Endpoints: endpoints}
	logger := zap.NewExample()
	defer logger.Sync()

	server := NewServer(searchConfig, logger)
	req := &api.SearchRequest{Query: "fail"}
	resp, err := server.Search(context.Background(), req)

	assert.Nil(t, err)
	assert.Equal(t, 5, len(resp.Results))
	assert.Equal(t, "[::]:8084", resp.Results[0].Source)
}

func TestTimeoutRequest(t *testing.T) {
	endpoints := []string{runServer(8085, &sleepingServer{}), runServer(8086, &sleepingServer{})}
	searchConfig := SearchConfig{Endpoints: endpoints}
	logger := zap.NewExample()
	defer logger.Sync()

	server := NewServer(searchConfig, logger)
	ctx := context.Background()
	req := &api.SearchRequest{Query: "fail"}
	runTime := time.Now()
	resp, err := server.Search(ctx, req)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(resp.Results))
	assert.True(t, time.Since(runTime) < 2*ReceiveTimeout)
}

func TestOneTimeoutRequest(t *testing.T) {
	endpoints := []string{runServer(8087, &sleepingServer{}), runServer(8088, &stubServer{response: makeStubResponse(5, "localhost:8088")})}
	searchConfig := SearchConfig{Endpoints: endpoints}
	logger := zap.NewExample()
	defer logger.Sync()

	server := NewServer(searchConfig, logger)
	req := &api.SearchRequest{Query: "fail"}
	resp, err := server.Search(context.Background(), req)

	assert.Nil(t, err)
	assert.Equal(t, 5, len(resp.Results))
	assert.Equal(t, "[::]:8088", resp.Results[0].Source)
}
