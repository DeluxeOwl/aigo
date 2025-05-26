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
	Base             []BaseOption             `exhaustruct:"optional"`
	OpenAICompatible []OpenAICompatibleOption `exhaustruct:"optional"`
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
func (o *OpenAICompatible) GenText(ctx context.Context, messages []schema.Message) (*aigo.GenTextResponse, error) {
	body := schema.Request{
		Model:    string(o.model),
		Messages: messages,
	}

	jsonb, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("gen text: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.url, bytes.NewBuffer(jsonb))
	if err != nil {
		return nil, fmt.Errorf("gen text: %w", err)
	}

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
		return nil, errors.New("gen text: empty response: " + u.Detail)
	}

	assistantMessage, ok := u.Choices[0].Message.(*schema.AssistantMessage)
	if !ok {
		return nil, errors.New("gen text: missing assistant message")
	}

	text, ok := assistantMessage.Content.(schema.StringPart)
	if !ok {
		return nil, errors.New("gen text: missing assistant string")
	}

	return &aigo.GenTextResponse{
		Text: text.String(),
	}, nil
}
