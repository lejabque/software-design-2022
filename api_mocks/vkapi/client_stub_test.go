package vkapi

import (
	"context"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type ClientStubTestSuite struct {
	suite.Suite

	stubServer *httptest.Server
	apiClient  Client
}

func (s *ClientStubTestSuite) SetupTest() {
	s.stubServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		s.Equal("GET", req.Method)
		s.Equal("/method/newsfeed.search", req.URL.Path)
		s.Equal("token", req.URL.Query().Get("access_token"))
		s.Equal("5.120", req.URL.Query().Get("v"))
		s.Equal("1", req.URL.Query().Get("count"))
		if req.URL.Query().Get("q") == "forbidden" {
			w.WriteHeader(403)
			w.Write([]byte("forbidden"))
			return
		}
		if req.URL.Query().Get("q") == "unavailable" {
			w.WriteHeader(500)
			w.Write([]byte("Internal server error"))
			return
		}

		s.Equal("#test_tag", req.URL.Query().Get("q"))
		s.Equal("123", req.URL.Query().Get("start_time"))
		s.Equal("321", req.URL.Query().Get("end_time"))
		w.WriteHeader(200)
		switch req.URL.Query().Get("start_from") {
		case "":
			w.Write(makeFakeResponse(12314324323, "10/-77262232_4158", "#test_tag #first_post"))
		case "10/-77262232_4158":
			w.Write(makeFakeResponse(12314324324, "10/-77262232_4159", "#test_tag #first_post"))
		default:
			s.Fail("unexpected start_from")
		}
	}))

	s.apiClient = NewVKClient("token", s.stubServer.URL, &http.Client{})
}

func (s *ClientStubTestSuite) TearDownTest() {
	s.stubServer.Close()
}

func (s *ClientStubTestSuite) TestSearchPostsErrors() {
	_, err := s.apiClient.SearchPosts(context.Background(), "forbidden",
		time.Unix(123, 0), time.Unix(321, 0), "", 1)
	s.Error(err)
	_, err = s.apiClient.SearchPosts(context.Background(), "unavailable",
		time.Unix(123, 0), time.Unix(321, 0), "", 1)
	s.Error(err)
}

func (s *ClientStubTestSuite) TestSearchPostsPaging() {
	first, err := s.apiClient.SearchPosts(context.Background(), "#test_tag",
		time.Unix(123, 0), time.Unix(321, 0), "", 1)
	s.NoError(err)
	if s.Equal(1, len(first.Items)) {
		s.Equal(12314324323, first.Items[0].ID)
	}

	second, err := s.apiClient.SearchPosts(context.Background(), "#test_tag",
		time.Unix(123, 0), time.Unix(321, 0), first.NextFrom, 1)
	s.NoError(err)
	if s.Equal(1, len(first.Items)) {
		s.Equal(12314324324, second.Items[0].ID)
	}
}

func TestClientStubTestSuite(t *testing.T) {
	suite.Run(t, new(ClientStubTestSuite))
}
