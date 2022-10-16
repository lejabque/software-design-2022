package vkapi

import (
	"bytes"
	"context"
	"fmt"
	webmocks "github.com/lejabque/software-design-2022/api_mocks/internal/web/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type ClientMockTestSuite struct {
	suite.Suite
	httpClient webmocks.HTTPClient
	apiClient  Client
}

func makeFakeResponse(id int, nextFrom string, text string) []byte {
	return []byte(fmt.Sprintf(`{
	"response": {
		"items": [
			{
				"id": %d,
				"date": 123,
				"owner_id": 1234556789,
				"from_id": 1234556789,
				"post_type": "post",
				"text": "%s"
			}
		],
		"next_from": "%s",
		"total_count": 39512
	}
}`, id, text, nextFrom))
}

func (s *ClientMockTestSuite) SetupTest() {
	s.httpClient = webmocks.HTTPClient{}
	s.apiClient = NewVKClient("token", "https://api.vk.com", &s.httpClient)
}

func (s *ClientMockTestSuite) TestSearchPostsErrors() {
	s.httpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Query().Get("access_token") == "token" && req.URL.Query().Get("q") == "forbidden"
	})).Return(&http.Response{
		StatusCode: 403,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("forbidden"))),
	}, nil).Once()
	s.httpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Query().Get("access_token") == "token" && req.URL.Query().Get("q") == "unavailable"
	})).Return(&http.Response{
		StatusCode: 500,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("Internal server error"))),
	}, nil).Once()

	_, err := s.apiClient.SearchPosts(context.Background(), "forbidden",
		time.Unix(123, 0), time.Unix(321, 0), "", 1)
	s.Error(err)
	_, err = s.apiClient.SearchPosts(context.Background(), "unavailable",
		time.Unix(123, 0), time.Unix(321, 0), "", 1)
	s.Error(err)
}

func (s *ClientMockTestSuite) TestSearchPostsPaging() {
	s.httpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		s.Equal("GET", req.Method)
		s.Equal("api.vk.com", req.Host)
		s.Equal("/method/newsfeed.search", req.URL.Path)
		s.Equal("token", req.URL.Query().Get("access_token"))
		s.Equal("5.120", req.URL.Query().Get("v"))
		s.Equal("1", req.URL.Query().Get("count"))
		s.Equal("#test_tag", req.URL.Query().Get("q"))
		s.Equal("123", req.URL.Query().Get("start_time"))
		s.Equal("321", req.URL.Query().Get("end_time"))
		return req.URL.Query().Get("next_from") == ""
	})).Return(&http.Response{
		StatusCode: 200,
		Body: ioutil.NopCloser(bytes.NewReader(
			makeFakeResponse(12314324323, "10/-77262232_4158", "#test_tag #first_post"))),
	}, nil).Once()

	s.httpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		s.Equal("GET", req.Method)
		s.Equal("api.vk.com", req.Host)
		s.Equal("/method/newsfeed.search", req.URL.Path)
		s.Equal("token", req.URL.Query().Get("access_token"))
		s.Equal("5.120", req.URL.Query().Get("v"))
		s.Equal("1", req.URL.Query().Get("count"))
		s.Equal("#test_tag", req.URL.Query().Get("q"))
		s.Equal("123", req.URL.Query().Get("start_time"))
		s.Equal("321", req.URL.Query().Get("end_time"))
		startFrom := req.URL.Query().Get("start_from")
		return startFrom == "10/-77262232_4158"
	})).Return(&http.Response{
		StatusCode: 200,
		Body: ioutil.NopCloser(bytes.NewReader(
			makeFakeResponse(12314324324, "10/-77262232_4159", "#test_tag #second_post"))),
	}, nil).Once()

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

func TestClientMockTestSuite(t *testing.T) {
	suite.Run(t, new(ClientMockTestSuite))
}
