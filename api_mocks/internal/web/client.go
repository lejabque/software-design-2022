package web

import "net/http"

//go:generate mockery --name=HTTPClient --exported --case underscore
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
