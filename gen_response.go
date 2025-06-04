package aigo

import (
	"errors"
	"fmt"

	"github.com/DeluxeOwl/aigo/provider/schema"
)

type GenResponse struct {
	Response *schema.Response
}

func (gr *GenResponse) GetLastAssistantMessage() (*AssistantMessage, error) {
	if len(gr.Response.Choices) == 0 {
		return nil, errors.New("GetLastAssistantMessage: empty response choices")
	}

	lastChoice := gr.Response.Choices[len(gr.Response.Choices)-1]
	assistantMsg, ok := lastChoice.Message.(*schema.AssistantMessage)
	if !ok {
		return nil, fmt.Errorf("GetLastAssistantMessage: last message in choices is not an *schema.AssistantMessage, but %T", lastChoice.Message)
	}

	return &AssistantMessage{Message: assistantMsg}, nil
}

type AssistantMessage struct {
	Message *schema.AssistantMessage // A direct pointer to the message in the GenResponse
}

type AssistantText struct {
	Text       string
	SetText    func(text string)
	SetContent func(content schema.ContentPartAssistant)
}

func (h *AssistantMessage) RunIfText(cb func(message AssistantText)) bool {
	if h.Message == nil {
		return false
	}
	stringPart, ok := h.Message.Content.(schema.StringPart)
	if !ok {
		return false
	}

	cb(AssistantText{
		Text: stringPart.String(),
		SetText: func(text string) {
			h.Message.Content = schema.StringPart(text)
		},
		SetContent: func(content schema.ContentPartAssistant) {
			h.Message.Content = content
		},
	})

	return true
}
