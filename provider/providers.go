package provider

import (
	"github.com/DeluxeOwl/aigo/provider/ollama"
)

func NewProviderOllama(model ollama.Model, options ...ollama.Option) *ollama.Ollama {
	return ollama.New(model, options...)
}
