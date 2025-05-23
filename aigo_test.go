package aigo_test

import (
	"testing"

	"github.com/DeluxeOwl/aigo"
	"github.com/DeluxeOwl/aigo/provider"
	"github.com/stretchr/testify/require"
)

func TestOllama(t *testing.T) {
	ctx := t.Context()

	resp, err := aigo.Ask(ctx, aigo.AskOptions{
		Provider: provider.NewOllama("qwen3:0.6B"),
		Message:  "/no_think How are you?",
	})
	require.NoError(t, err)

	t.Log(resp)
}
