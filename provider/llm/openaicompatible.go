package llm

import (
	"bytes"
	"context"
	"encoding/json"
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
func (o *OpenAICompatible) Gen(ctx context.Context, messages []schema.Message) (*aigo.GenResponse, error) {
	body := schema.Request{
		Model:    string(o.model),
		Messages: messages,
	}

	if o.baseOnBeforeRequestMarshal != nil {
		o.baseOnBeforeRequestMarshal(&body)
	}

	jsonb, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("gen: %w", err)
	}

	if o.baseOnBeforeRequestBody != nil {
		o.baseOnBeforeRequestBody(jsonb)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.url, bytes.NewBuffer(jsonb))
	if err != nil {
		return nil, fmt.Errorf("gen: %w", err)
	}

	if o.baseOnBeforeRequestSend != nil {
		o.baseOnBeforeRequestSend(req)
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gen: %w", err)
	}
	defer resp.Body.Close()

	if o.baseOnBeforeResponseRead != nil {
		o.baseOnBeforeResponseRead(resp)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gen: %w", err)
	}

	// TODO: unmarshal into a custom struct and stop? Return it as a method on GenResponse?
	if o.baseOnBeforeResponseUnmarshal != nil {
		o.baseOnBeforeResponseUnmarshal(b)
	}

	var u schema.Response
	err = json.Unmarshal(b, &u)
	if err != nil {
		return nil, fmt.Errorf("gen: %w", err)
	}

	return &aigo.GenResponse{
		Response: u,
	}, nil
}
