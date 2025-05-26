package llm

const DefaultBaseURLOllama = "http://localhost:11434/v1"

func NewOllama(model OpenAICompatibleModel) *OpenAICompatible {
	ollama := &OpenAICompatible{
		Base:  NewBaseConfig(BaseURL(DefaultBaseURLOllama)),
		model: model,
	}

	return ollama
}

func NewOllamaWithConfig(model OpenAICompatibleModel, cfg *OpenAICompatibleConfig) *OpenAICompatible {
	baseConfigOpts := []BaseOption{BaseURL(DefaultBaseURLOllama)}
	if cfg != nil {
		baseConfigOpts = append(baseConfigOpts, cfg.Base...)
	}

	ollama := &OpenAICompatible{
		Base:  NewBaseConfig(baseConfigOpts...),
		model: model,
	}

	if cfg != nil {
		for _, opt := range cfg.OpenAICompatible {
			opt(ollama)
		}
	}

	return ollama
}
