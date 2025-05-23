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

const DefaultBaseURLOllama = "http://localhost:11434/v1"

type Ollama struct {
	*BaseConfig
	model OllamaModel
}

type OllamaOption func(*Ollama)

type OllamaModel string

type OllamaConfig struct {
	Base   []BaseConfigOption
	Ollama []OllamaOption
}

func NewOllama(model OllamaModel) *Ollama {
	ollama := &Ollama{
		BaseConfig: NewBaseConfig(BaseURL(DefaultBaseURLOllama)),
		model:      model,
	}

	return ollama
}

func NewOllamaWithConfig(model OllamaModel, cfg *OllamaConfig) *Ollama {
	baseConfigOpts := []BaseConfigOption{BaseURL(DefaultBaseURLOllama)}
	if cfg != nil {
		baseConfigOpts = append(baseConfigOpts, cfg.Base...)
	}

	ollama := &Ollama{
		BaseConfig: NewBaseConfig(baseConfigOpts...),
		model:      model,
	}

	if cfg != nil {
		for _, opt := range cfg.Ollama {
			opt(ollama)
		}
	}

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
