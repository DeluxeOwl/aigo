package llm

import "net/http"

type Base struct {
	client  *http.Client
	url     string `exhaustruct:"optional"`
	baseURL string `exhaustruct:"optional"`
	path    string `exhaustruct:"optional"`
	apiKey  string `exhaustruct:"optional"`
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
