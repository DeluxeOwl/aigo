package aigo_test

import (
	"testing"

	"github.com/DeluxeOwl/aigo/provider"
	"github.com/stretchr/testify/require"
)

func TestOllama(t *testing.T) {
	p := provider.NewProviderOllama("qwen3:0.6B")

	resp, err := p.Ask(t.Context(), "/no_think How are you?")
	require.NoError(t, err)

	t.Log(resp)
}
