package llm

import (
	"encoding/json"
	"net/http"

	"github.com/DeluxeOwl/aigo/provider/schema"
)

type Base struct {
	client  *http.Client
	url     string `exhaustruct:"optional"`
	baseURL string `exhaustruct:"optional"`
	path    string `exhaustruct:"optional"`
	apiKey  string `exhaustruct:"optional"`

	baseOnBeforeRequestMarshal    func(req *schema.Request)  `exhaustruct:"optional"`
	baseOnBeforeRequestBody       func(body json.RawMessage) `exhaustruct:"optional"`
	baseOnBeforeRequestSend       func(req *http.Request)    `exhaustruct:"optional"`
	baseOnBeforeResponseRead      func(req *http.Response)   `exhaustruct:"optional"`
	baseOnBeforeResponseUnmarshal func(body []byte)          `exhaustruct:"optional"`
}

type BaseOption func(*Base)

func APIBearerKey(key string) BaseOption {
	return func(c *Base) {
		c.apiKey = key
	}
}

func HTTPClient(client *http.Client) BaseOption {
	return func(c *Base) {
		c.client = client
	}
}

func BaseURL(baseURL string) BaseOption {
	return func(c *Base) {
		c.baseURL = baseURL
	}
}

func CompletionsPath(path string) BaseOption {
	return func(c *Base) {
		c.path = path
	}
}

func OnBeforeRequestMarshal(cb func(req *schema.Request)) BaseOption {
	return func(c *Base) {
		c.baseOnBeforeRequestMarshal = cb
	}
}

func OnBeforeRequestBody(cb func(body json.RawMessage)) BaseOption {
	return func(c *Base) {
		c.baseOnBeforeRequestBody = cb
	}
}

func OnBeforeRequestSend(cb func(*http.Request)) BaseOption {
	return func(c *Base) {
		c.baseOnBeforeRequestSend = cb
	}
}

func OnBeforeResponseRead(cb func(*http.Response)) BaseOption {
	return func(c *Base) {
		c.baseOnBeforeResponseRead = cb
	}
}

func OnBeforeResponseUnmarshal(cb func(body []byte)) BaseOption {
	return func(c *Base) {
		c.baseOnBeforeResponseUnmarshal = cb
	}
}

func NewBaseConfig(options ...BaseOption) *Base {
	c := &Base{
		client: NewDefaultHTTPClient(),
		path:   DefaultChatCompletionsPath,
	}

	for _, opt := range options {
		opt(c)
	}

	c.url = BuildChatCompletionsURL(c.baseURL, c.path)

	// TODO: is this good here?
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	if c.apiKey != "" {
		headers["Authorization"] = "Bearer " + c.apiKey
	}

	c.client.Transport = &withHeadersTransport{
		Headers:   headers,
		Transport: c.client.Transport,
	}

	return c
}

type withHeadersTransport struct {
	Transport http.RoundTripper
	Headers   map[string]string
}

func (t *withHeadersTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	outReq := req.Clone(req.Context())

	for k, v := range t.Headers {
		outReq.Header.Set(k, v)
	}

	return t.Transport.RoundTrip(outReq)
}
