package aigo

import (
	"context"

	"github.com/DeluxeOwl/aigo/provider/schema"
)

type GenTexter interface {
	GenText(ctx context.Context, messages []schema.Message) (*GenTextResponse, error)
}

type GenTextResponse struct {
	Text string `json:"text"`
}

type GenTextNextFn func(ctx context.Context, options *GenTextOptions) (*GenTextResponse, error)

type GenTextMiddleware interface {
	Process(ctx context.Context, options *GenTextOptions, next GenTextNextFn) (*GenTextResponse, error)
}

type GenTextMiddlewareFunc func(ctx context.Context, options *GenTextOptions, next GenTextNextFn) (*GenTextResponse, error)

func (f GenTextMiddlewareFunc) Process(ctx context.Context, options *GenTextOptions, next GenTextNextFn) (*GenTextResponse, error) {
	return f(ctx, options, next)
}

type GenTextOptions struct {
	Provider   GenTexter
	Messages   []schema.Message
	Middleware []GenTextMiddleware `exhaustruct:"optional"`
}

func GenText(ctx context.Context, options *GenTextOptions) (*GenTextResponse, error) {
	coreOperation := func(currentCtx context.Context, currentOpts *GenTextOptions) (*GenTextResponse, error) {
		return currentOpts.Provider.GenText(currentCtx, currentOpts.Messages)
	}

	chainedHandler := GenTextNextFn(coreOperation)

	if options.Middleware != nil {
		for i := len(options.Middleware) - 1; i >= 0; i-- {
			mw := options.Middleware[i]
			nextInChain := chainedHandler

			chainedHandler = func(c context.Context, o *GenTextOptions) (*GenTextResponse, error) {
				return mw.Process(c, o, nextInChain)
			}
		}
	}

	return chainedHandler(ctx, options)
}
