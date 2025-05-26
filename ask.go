package aigo

import (
	"context"
)

type GenTexter interface {
	GenText(ctx context.Context, message string) (*GenTextResponse, error)
}

type GenTextOptions struct {
	Provider GenTexter
	Message  string
	Hooks    *GenTextHooks `exhaustruct:"optional"`
}

type GenTextHooks struct {
	BeforeGenText []BeforeGenTexter `exhaustruct:"optional"`
	AfterGenText  []AfterGenTexter  `exhaustruct:"optional"`
}

type BeforeGenTexter interface {
	BeforeGenText(ctx context.Context, options *GenTextOptions) error
}
type BeforeGenText func(ctx context.Context, options *GenTextOptions) error

func (f BeforeGenText) BeforeGenText(ctx context.Context, options *GenTextOptions) error {
	return f(ctx, options)
}

type AfterGenTexter interface {
	AfterGenText(ctx context.Context, res *GenTextResponse, err error) (*GenTextResponse, error)
}
type AfterGenText func(ctx context.Context, res *GenTextResponse, err error) (*GenTextResponse, error)

func (f AfterGenText) AfterGenText(ctx context.Context, res *GenTextResponse, err error) (*GenTextResponse, error) {
	return f(ctx, res, err)
}

type GenTextResponse struct {
	Text string `json:"text"`
}

func GenText(ctx context.Context, options *GenTextOptions) (*GenTextResponse, error) {
	if options.Hooks != nil && len(options.Hooks.BeforeGenText) > 0 {
		for _, h := range options.Hooks.BeforeGenText {
			err := h.BeforeGenText(ctx, options)
			if err != nil {
				return nil, err
			}
		}
	}

	res, err := options.Provider.GenText(ctx, options.Message)

	if options.Hooks != nil && len(options.Hooks.AfterGenText) > 0 {
		for _, h := range options.Hooks.AfterGenText {
			res, err = h.AfterGenText(ctx, res, err)
		}
	}

	return res, err
}
