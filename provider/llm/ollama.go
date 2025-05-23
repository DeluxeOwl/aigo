package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const OllamaDefaultBaseURL = "http://localhost:11434/v1"

type Ollama struct {
	client  *http.Client
	baseURL string
	path    string
	model   OllamaModel

	url string `exhaustruct:"optional"`
}

type OllamaOption func(*Ollama)

func OllamaHTTPClient(client *http.Client) OllamaOption {
	return func(o *Ollama) {
		o.client = client
	}
}

func OllamaBaseURL(baseURL string) OllamaOption {
	return func(o *Ollama) {
		o.baseURL = baseURL
	}
}

func OllamaCompletionsPath(path string) OllamaOption {
	return func(o *Ollama) {
		o.path = path
	}
}

type OllamaModel string

func NewOllama(model OllamaModel, options ...OllamaOption) *Ollama {
	ollama := &Ollama{
		client:  NewDefaultHTTPClient(),
		baseURL: OllamaDefaultBaseURL,
		path:    DefaultChatCompletionsPath,
		model:   model,
	}

	for _, opt := range options {
		opt(ollama)
	}

	ollama.url = BuildChatCompletionsURL(ollama.baseURL, ollama.path)

	return ollama
}

// TODO: errors
func (o *Ollama) Ask(ctx context.Context, message string) (string, error) {
	body := fmt.Sprintf(`
{
	"model": "%s",
	"messages": [
		{
			"role": "user",
			"content": "%s"
		}
	]
}
`, o.model, message)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.url, strings.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("ask: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ask: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ask: %w", err)
	}

	var u Response
	err = json.Unmarshal(b, &u)
	if err != nil {
		return "", fmt.Errorf("ask: %w", err)
	}

	if len(u.Choices) == 0 {
		return "", errors.New("ask: empty response")
	}

	return u.Choices[0].Message.Content, nil
}
