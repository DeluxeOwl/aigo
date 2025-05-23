package provider

import "github.com/DeluxeOwl/aigo/provider/llm"

func NewOllama(model llm.OllamaModel) *llm.Ollama {
	return llm.NewOllama(model)
}

func NewOllamaWithConfig(model llm.OllamaModel, cfg llm.OllamaConfig) *llm.Ollama {
	return llm.NewOllamaWithConfig(model, &cfg)
}
