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
		Hooks: &aigo.GenTextHooks{
			BeforeGenText: []aigo.BeforeGenTexter{
				aigo.BeforeGenText(func(_ context.Context, options *aigo.GenTextOptions) error {
					options.Message = strings.ReplaceAll(options.Message, "fox", "bear")
					return nil
				}),
			},
			AfterGenText: []aigo.AfterGenTexter{
				aigo.AfterGenText(func(_ context.Context, res *aigo.GenTextResponse, err error) (*aigo.GenTextResponse, error) {
					res.Text = strings.ToUpper(res.Text)
					return res, err
				}),
			},
		},
	})
	require.NoError(t, err)

	t.Log(resp)
}
