package aigo

import (
	"context"
)

type Asker interface {
	Ask(ctx context.Context, message string) (string, error)
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
	AfterAsk(ctx context.Context, res string, err error) (string, error)
}
type AfterAsk func(ctx context.Context, res string, err error) (string, error)

func (f AfterAsk) AfterAsk(ctx context.Context, response string, err error) (string, error) {
	return f(ctx, response, err)
}

// TODO: what about hooks that always run?
func Ask(ctx context.Context, options *AskOptions) (string, error) {
	if options.Hooks != nil && len(options.Hooks.BeforeAsk) > 0 {
		for _, h := range options.Hooks.BeforeAsk {
			err := h.BeforeAsk(ctx, options)
			if err != nil {
				return "", err
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
