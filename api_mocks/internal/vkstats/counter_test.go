package vkstats

import (
	"context"
	"github.com/lejabque/software-design-2022/api_mocks/internal/vkapi"
	"github.com/lejabque/software-design-2022/api_mocks/internal/vkstats/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CounterTestSuite struct {
	suite.Suite
	vkClient *mocks.VkClient
	counter  Counter
}

func (s *CounterTestSuite) SetupTest() {
	s.vkClient = &mocks.VkClient{}
	s.counter = Counter{s.vkClient}
}

func (s *CounterTestSuite) TestCountPostsDistribution() {
	s.vkClient.On("SearchPosts", mock.Anything, "#test_tag", mock.Anything, mock.Anything, "", 0).
		Return(vkapi.PostsResponse{
			TotalCount: 10,
			NextFrom:   "first_req",
		}, nil).Once()
	s.vkClient.On("SearchPosts", mock.Anything, "#test_tag", mock.Anything, mock.Anything, "first_req", 0).
		Return(vkapi.PostsResponse{
			TotalCount: 20,
			NextFrom:   "second_req",
		}, nil).Once()
	s.vkClient.On("SearchPosts", mock.Anything, "#test_tag", mock.Anything, mock.Anything, "second_req", 0).
		Return(vkapi.PostsResponse{
			TotalCount: 30,
			NextFrom:   "third_req",
		}, nil).Once()

	res, err := s.counter.CountPostsDistribution(context.Background(), "#test_tag", 3)
	s.NoError(err)
	if s.Len(res, 3) {
		s.Equal(int64(10), res[0])
		s.Equal(int64(20), res[1])
		s.Equal(int64(30), res[2])
	}
}

func TestCounterTestSuite(t *testing.T) {
	suite.Run(t, new(CounterTestSuite))
}
