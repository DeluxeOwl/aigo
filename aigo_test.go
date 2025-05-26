package aigo_test

import (
	"context"
	"strings"
	"testing"

	"github.com/DeluxeOwl/aigo"
	"github.com/DeluxeOwl/aigo/provider"
	"github.com/DeluxeOwl/aigo/provider/schema"
	"github.com/stretchr/testify/require"
)

func TestOllama(t *testing.T) {
	ctx := t.Context()

	resp, err := aigo.GenText(ctx, &aigo.GenTextOptions{
		Provider: provider.NewOllama("qwen3:0.6B"),
		Messages: []schema.Message{
			schema.NewSystemMessage("TALK LIKE A PIRATE WITH EMOJIS"),
			schema.NewUserMessage([]schema.ContentPartUser{
				schema.NewTextPart("/no_think What can you tell me about a fox?"),
			}),
		},
		Middleware: []aigo.GenTextMiddleware{
			aigo.GenTextMiddlewareFunc(func(ctx context.Context, options *aigo.GenTextOptions, next aigo.GenTextNextFn) (*aigo.GenTextResponse, error) {
				for _, message := range options.Messages {
					if muser, ok := message.(*schema.UserMessage); ok {
						for _, c := range muser.Content {
							if text, ok := c.(*schema.TextPart); ok {
								text.Text = strings.ReplaceAll(text.Text, "fox", "bear")
							}
						}
					}
				}

				return next(ctx, options)
			}),

			aigo.GenTextMiddlewareFunc(func(ctx context.Context, options *aigo.GenTextOptions, next aigo.GenTextNextFn) (*aigo.GenTextResponse, error) {
				res, err := next(ctx, options)

				if err == nil && res != nil {
					res.Text = strings.ToUpper(res.Text)
				}
				return res, err
			}),
		},
	})
	require.NoError(t, err)

	t.Log(resp.Text)
}
