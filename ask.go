package aigo

import (
	"context"
)

type Asker interface {
	Ask(ctx context.Context, message string) (*AskResponse, error)
}

type AskOptions struct {
	Provider Asker
	Message  string
	Hooks    *AskHooks `exhaustruct:"optional"`
}

type AskHooks struct {
	BeforeAsk []BeforeAsker `exhaustruct:"optional"`
	AfterAsk  []AfterAsker  `exhaustruct:"optional"`
}

type BeforeAsker interface {
	BeforeAsk(ctx context.Context, options *AskOptions) error
}
type BeforeAsk func(ctx context.Context, options *AskOptions) error

func (f BeforeAsk) BeforeAsk(ctx context.Context, options *AskOptions) error {
	return f(ctx, options)
}

type AfterAsker interface {
	AfterAsk(ctx context.Context, res *AskResponse, err error) (*AskResponse, error)
}
type AfterAsk func(ctx context.Context, res *AskResponse, err error) (*AskResponse, error)

func (f AfterAsk) AfterAsk(ctx context.Context, res *AskResponse, err error) (*AskResponse, error) {
	return f(ctx, res, err)
}

type AskResponse struct {
	Text string `json:"text"`
}

// TODO: what about hooks that always run?
func Ask(ctx context.Context, options *AskOptions) (*AskResponse, error) {
	if options.Hooks != nil && len(options.Hooks.BeforeAsk) > 0 {
		for _, h := range options.Hooks.BeforeAsk {
			err := h.BeforeAsk(ctx, options)
			if err != nil {
				return nil, err
			}
		}
	}

	res, err := options.Provider.Ask(ctx, options.Message)

	if options.Hooks != nil && len(options.Hooks.AfterAsk) > 0 {
		for _, h := range options.Hooks.AfterAsk {
			res, err = h.AfterAsk(ctx, res, err)
		}
	}

	return res, err
}
