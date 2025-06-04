package aigo

import (
	"context"

	"github.com/DeluxeOwl/aigo/provider/schema"
)

type Generator interface {
	Gen(ctx context.Context, messages []schema.Message) (*GenResponse, error)
}

type GenNextFn func(ctx context.Context, options *GenOptions) (*GenResponse, error)

type GenMiddleware interface {
	Process(ctx context.Context, options *GenOptions, next GenNextFn) (*GenResponse, error)
}

type GenMiddlewareFunc func(ctx context.Context, options *GenOptions, next GenNextFn) (*GenResponse, error)

func (f GenMiddlewareFunc) Process(ctx context.Context, options *GenOptions, next GenNextFn) (*GenResponse, error) {
	return f(ctx, options, next)
}

type GenOptions struct {
	Provider   Generator
	Messages   []schema.Message
	Middleware []GenMiddleware `exhaustruct:"optional"`
}

func Gen(ctx context.Context, options *GenOptions) (*GenResponse, error) {
	coreOperation := func(currentCtx context.Context, currentOpts *GenOptions) (*GenResponse, error) {
		return currentOpts.Provider.Gen(currentCtx, currentOpts.Messages)
	}

	chainedHandler := GenNextFn(coreOperation)

	if options.Middleware != nil {
		for i := len(options.Middleware) - 1; i >= 0; i-- {
			mw := options.Middleware[i]
			nextInChain := chainedHandler

			chainedHandler = func(c context.Context, o *GenOptions) (*GenResponse, error) {
				return mw.Process(c, o, nextInChain)
			}
		}
	}

	return chainedHandler(ctx, options)
}
