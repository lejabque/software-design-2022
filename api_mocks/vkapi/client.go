package vkapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lejabque/software-design-2022/api_mocks/internal/web"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	token   string
	baseURL string
	client  web.HTTPClient
}

type PostsResponse struct {
	Items []struct {
		ID       int    `json:"id"`
		Date     int    `json:"date"`
		OwnerId  int    `json:"owner_id"`
		FromId   int    `json:"from_id"`
		PostType string `json:"post_type"`
		Text     string `json:"text"`
	} `json:"items"`
	NextFrom   string `json:"next_from"`
	TotalCount int    `json:"total_count"`
}

func NewVKClient(token string, baseURL string, client web.HTTPClient) Client {
	return Client{
		token:   token,
		baseURL: baseURL,
		client:  client,
	}
}

func (c Client) SearchPosts(ctx context.Context, queryText string, startTime time.Time, endTime time.Time, startFrom string, count int) (PostsResponse, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/method/newsfeed.search", c.baseURL), nil)
	if err != nil {
		return PostsResponse{}, err
	}
	params := url.Values{}
	params.Add("access_token", c.token)
	params.Add("v", "5.120")
	params.Add("q", queryText)
	params.Add("start_time", fmt.Sprintf("%d", startTime.Unix()))
	params.Add("end_time", fmt.Sprintf("%d", endTime.Unix()))
	params.Add("start_from", startFrom)
	params.Add("count", strconv.Itoa(count))
	r.URL.RawQuery = params.Encode()
	resp, err := c.client.Do(r)
	if err != nil {
		return PostsResponse{}, err
	}
	defer resp.Body.Close()
	var parsed struct {
		Response PostsResponse `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return PostsResponse{}, err
	}
	return parsed.Response, nil
}
