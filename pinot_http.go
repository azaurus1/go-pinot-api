package goPinotAPI

import (
	"net/http"
	"net/url"
)

type pinotHttp struct {
	httpClient         *http.Client
	pinotControllerUrl *url.URL
	httpAuthWriter     httpAuthWriter
}

type httpAuthWriter func(*http.Request)

func (p *pinotHttp) Do(req *http.Request) (*http.Response, error) {
	p.httpAuthWriter(req)
	return p.httpClient.Do(req)
}
