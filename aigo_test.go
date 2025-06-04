package aigo_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/DeluxeOwl/aigo"
	"github.com/DeluxeOwl/aigo/provider"
	"github.com/DeluxeOwl/aigo/provider/llm"
	"github.com/DeluxeOwl/aigo/provider/schema"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
)

func TestOllama(t *testing.T) {
	ctx := t.Context()

	resp, err := aigo.Gen(ctx, &aigo.GenOptions{
		Provider: provider.NewOpenAICompatibleWithConfig("google/gemini-2.5-flash-preview-05-20", llm.OpenAICompatibleConfig{
			Base: []llm.BaseOption{
				llm.APIBearerKey(os.Getenv("OPENWEB_UI_KEY")),
				llm.BaseURL("http://localhost:3000"),
				llm.CompletionsPath("/api/chat/completions"),
			},
		}),
		Messages: []schema.Message{
			schema.NewSystemMessage("TALK LIKE A PIRATE WITH EMOJIS, also, answer in 100 words"),
			schema.NewUserMessage([]schema.ContentPartUser{
				schema.NewTextPart("/no_think What can you tell me about a fox?"),
			}),
		},
		Middleware: []aigo.GenMiddleware{
			aigo.GenMiddlewareFunc(func(ctx context.Context, options *aigo.GenOptions, next aigo.GenNextFn) (*aigo.GenResponse, error) {
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

			aigo.GenMiddlewareFunc(func(ctx context.Context, options *aigo.GenOptions, next aigo.GenNextFn) (*aigo.GenResponse, error) {
				res, err := next(ctx, options)

				if err == nil && res != nil {
					am, err := res.GetLastAssistantMessage()
					if err != nil {
						return res, err
					}

					am.RunIfText(func(message aigo.AssistantText) {
						message.SetText(strings.ToUpper(message.Text))
					})

				}
				return res, err
			}),
		},
	})
	require.NoError(t, err)

	am, err := resp.GetLastAssistantMessage()
	require.NoError(t, err)

	am.RunIfText(func(message aigo.AssistantText) {
		t.Log(message.Text)
	})
}
