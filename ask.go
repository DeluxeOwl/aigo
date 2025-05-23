package aigo

import (
	"context"
)

type Asker interface {
	Ask(ctx context.Context, message string) (string, error)
}

type AskOptions struct {
	Provider Asker
	Message  string
}

func Ask(ctx context.Context, options AskOptions) (string, error) {
	return options.Provider.Ask(ctx, options.Message)
}
