package llm

import "net/http"

type BaseConfig struct {
	client  *http.Client
	url     string `exhaustruct:"optional"`
	baseURL string `exhaustruct:"optional"`
	path    string `exhaustruct:"optional"`
}

type BaseConfigOption func(*BaseConfig)

func HTTPClient(client *http.Client) BaseConfigOption {
	return func(c *BaseConfig) {
		c.client = client
	}
}

func BaseURL(baseURL string) BaseConfigOption {
	return func(c *BaseConfig) {
		c.baseURL = baseURL
	}
}

func CompletionsPath(path string) BaseConfigOption {
	return func(c *BaseConfig) {
		c.path = path
	}
}

func NewBaseConfig(options ...BaseConfigOption) *BaseConfig {
	c := &BaseConfig{
		client: NewDefaultHTTPClient(),
		path:   DefaultChatCompletionsPath,
	}

	for _, opt := range options {
		opt(c)
	}

	c.url = BuildChatCompletionsURL(c.baseURL, c.path)

	return c
}
