package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/DeluxeOwl/aigo"
)

const DefaultBaseURLOllama = "http://localhost:11434/v1"

type Ollama struct {
	*Base
	model OllamaModel
}

type OllamaOption func(*Ollama)

type OllamaModel string

type OllamaConfig struct {
	Base   []BaseOption
	Ollama []OllamaOption
}

func NewOllama(model OllamaModel) *Ollama {
	ollama := &Ollama{
		Base:  NewBaseConfig(BaseURL(DefaultBaseURLOllama)),
		model: model,
	}

	return ollama
}

func NewOllamaWithConfig(model OllamaModel, cfg *OllamaConfig) *Ollama {
	baseConfigOpts := []BaseOption{BaseURL(DefaultBaseURLOllama)}
	if cfg != nil {
		baseConfigOpts = append(baseConfigOpts, cfg.Base...)
	}

	ollama := &Ollama{
		Base:  NewBaseConfig(baseConfigOpts...),
		model: model,
	}

	if cfg != nil {
		for _, opt := range cfg.Ollama {
			opt(ollama)
		}
	}

	return ollama
}

// TODO: errors
func (o *Ollama) Ask(ctx context.Context, message string) (*aigo.AskResponse, error) {
	body := Request{
		Model: string(o.model),
		Messages: []Message{
			{
				Role:    MessageRoleUser,
				Content: message,
			},
		},
	}

	jsonb, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("ask: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.url, bytes.NewBuffer(jsonb))
	if err != nil {
		return nil, fmt.Errorf("ask: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ask: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ask: %w", err)
	}

	var u Response
	err = json.Unmarshal(b, &u)
	if err != nil {
		return nil, fmt.Errorf("ask: %w", err)
	}

	if len(u.Choices) == 0 {
		return nil, errors.New("ask: empty response")
	}

	return &aigo.AskResponse{
		Text: u.Choices[0].Message.Content,
	}, nil
}
