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
	"github.com/DeluxeOwl/aigo/provider/schema"
)

type OpenAICompatible struct {
	*Base
	model OpenAICompatibleModel
}

type OpenAICompatibleOption func(*OpenAICompatible)

type OpenAICompatibleModel string

type OpenAICompatibleConfig struct {
	Base             []BaseOption
	OpenAICompatible []OpenAICompatibleOption
}

func NewOpenAICompatibleWithConfig(model OpenAICompatibleModel, cfg *OpenAICompatibleConfig) *OpenAICompatible {
	baseConfigOpts := []BaseOption{}
	if cfg != nil {
		baseConfigOpts = append(baseConfigOpts, cfg.Base...)
	}

	openAICompatible := &OpenAICompatible{
		Base:  NewBaseConfig(baseConfigOpts...),
		model: model,
	}

	if cfg != nil {
		for _, opt := range cfg.OpenAICompatible {
			opt(openAICompatible)
		}
	}

	return openAICompatible
}

// TODO: errors
func (o *OpenAICompatible) GenText(ctx context.Context, message string) (*aigo.GenTextResponse, error) {
	body := schema.Request{
		Model: string(o.model),
		Messages: []schema.Message{
			{
				Role:    schema.MessageRoleUser,
				Content: message,
			},
		},
	}

	jsonb, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("gen text: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.url, bytes.NewBuffer(jsonb))
	if err != nil {
		return nil, fmt.Errorf("gen text: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gen text: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gen text: %w", err)
	}

	var u schema.Response
	err = json.Unmarshal(b, &u)
	if err != nil {
		return nil, fmt.Errorf("gen text: %w", err)
	}

	if len(u.Choices) == 0 {
		return nil, errors.New("gen text: empty response")
	}

	return &aigo.GenTextResponse{
		Text: u.Choices[0].Message.Content,
	}, nil
}
