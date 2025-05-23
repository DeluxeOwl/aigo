package llm

import "net/http"

type Base struct {
	client  *http.Client
	url     string `exhaustruct:"optional"`
	baseURL string `exhaustruct:"optional"`
	path    string `exhaustruct:"optional"`
}

type BaseOption func(*Base)

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

	return c
}
