package llm

import (
	"fmt"
	"net/http"
	"time"
)

const (
	DefaultHTTPTimeout         = 30 * time.Second
	DefaultChatCompletionsPath = "/chat/completions"
)

func NewDefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   DefaultHTTPTimeout,
		Transport: http.DefaultTransport,
	}
}

func BuildChatCompletionsURL(baseURL, path string) string {
	return fmt.Sprintf("%s%s", baseURL, path)
}
