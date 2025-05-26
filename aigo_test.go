package aigo_test

import (
	"context"
	"strings"
	"testing"

	"github.com/DeluxeOwl/aigo"
	"github.com/DeluxeOwl/aigo/provider"
	"github.com/stretchr/testify/require"
)

func TestOllama(t *testing.T) {
	ctx := t.Context()

	resp, err := aigo.GenText(ctx, &aigo.GenTextOptions{
		Provider: provider.NewOllama("qwen3:0.6B"),
		Message:  "/no_think What can you tell me about a fox?",
		Middleware: []aigo.GenTextMiddleware{
			aigo.GenTextMiddlewareFunc(func(ctx context.Context, options *aigo.GenTextOptions, next aigo.GenTextNextFn) (*aigo.GenTextResponse, error) {
				options.Message = strings.ReplaceAll(options.Message, "fox", "bear")

				return next(ctx, options)
			}),

			aigo.GenTextMiddlewareFunc(func(ctx context.Context, options *aigo.GenTextOptions, next aigo.GenTextNextFn) (*aigo.GenTextResponse, error) {
				res, err := next(ctx, options)

				if err == nil && res != nil {
					res.Text = strings.ToUpper(res.Text)
				}
				return res, err
			}),
		},
	})
	require.NoError(t, err)

	t.Log(resp.Text)
}
