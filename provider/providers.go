package provider

import "github.com/DeluxeOwl/aigo/provider/llm"

func NewOllama(model llm.OllamaModel, options ...llm.OllamaOption) *llm.Ollama {
	return llm.NewOllama(model, options...)
}
