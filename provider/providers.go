package provider

import "github.com/DeluxeOwl/aigo/provider/llm"

func NewOllama(model llm.OpenAICompatibleModel) *llm.OpenAICompatible {
	return llm.NewOllama(model)
}

func NewOllamaWithConfig(model llm.OpenAICompatibleModel, cfg llm.OpenAICompatibleConfig) *llm.OpenAICompatible {
	return llm.NewOllamaWithConfig(model, &cfg)
}

func NewOpenAICompatibleWithConfig(model llm.OpenAICompatibleModel, cfg llm.OpenAICompatibleConfig) *llm.OpenAICompatible {
	return llm.NewOpenAICompatibleWithConfig(model, &cfg)
}
