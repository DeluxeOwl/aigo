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

	resp, err := aigo.Ask(ctx, &aigo.AskOptions{
		Provider: provider.NewOllama("qwen3:0.6B"),
		Message:  "/no_think What can you tell me about a fox?",
		Hooks: &aigo.AskHooks{
			BeforeAsk: []aigo.BeforeAsker{
				aigo.BeforeAsk(func(_ context.Context, options *aigo.AskOptions) error {
					options.Message = strings.ReplaceAll(options.Message, "fox", "bear")
					return nil
				}),
			},
			AfterAsk: []aigo.AfterAsker{
				aigo.AfterAsk(func(_ context.Context, s string, err error) (string, error) {
					return strings.ToUpper(s), err
				}),
			},
		},
	})
	require.NoError(t, err)

	t.Log(resp)
}
