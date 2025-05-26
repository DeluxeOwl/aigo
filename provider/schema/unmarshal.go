package schema

import (
	"encoding/json"
	"errors"
	"fmt"
)

type tempMessageRolePeek struct {
	Role MessageRole `json:"role"`
}

func (c *Choice) UnmarshalJSON(data []byte) error {
	type Alias Choice
	aux := &struct {
		Message json.RawMessage `json:"message"`
		*Alias
	}{
		Alias:   (*Alias)(c),
		Message: json.RawMessage{},
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("choice unmarshal: failed to unmarshal to auxiliary struct: %w", err)
	}

	if aux.Message == nil {
		return errors.New("choice unmarshal: 'message' field is missing or null")
	}

	var rolePeek tempMessageRolePeek
	if err := json.Unmarshal(aux.Message, &rolePeek); err != nil {
		return fmt.Errorf("choice unmarshal: failed to peek message role from raw message %s: %w", string(aux.Message), err)
	}

	switch rolePeek.Role {
	case MessageRoleSystem:
		var sm SystemMessage
		if err := json.Unmarshal(aux.Message, &sm); err != nil {
			return fmt.Errorf("choice unmarshal: failed to unmarshal SystemMessage from %s: %w", string(aux.Message), err)
		}
		c.Message = &sm
	case MessageRoleUser:
		var um UserMessage

		if err := json.Unmarshal(aux.Message, &um); err != nil {
			return fmt.Errorf("choice unmarshal: failed to unmarshal UserMessage from %s: %w", string(aux.Message), err)
		}
		c.Message = &um
	case MessageRoleAssistant:
		var am AssistantMessage
		if err := json.Unmarshal(aux.Message, &am); err != nil {
			return fmt.Errorf("choice unmarshal: failed to unmarshal AssistantMessage from %s: %w", string(aux.Message), err)
		}
		c.Message = &am
	case MessageRoleTool:
		return fmt.Errorf("choice unmarshal: MessageRoleTool deserialization not yet implemented for Choice.Message. Raw: %s", string(aux.Message))
	default:
		return fmt.Errorf("choice unmarshal: unknown or unhandled message role '%s' in raw message: %s", rolePeek.Role, string(aux.Message))
	}
	return nil
}

func (am *AssistantMessage) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Role    MessageRole     `json:"role"`
		Content json.RawMessage `json:"content"`
	}{
		Role:    "",
		Content: json.RawMessage{},
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("assistantmessage unmarshal: failed to unmarshal to auxiliary struct: %w", err)
	}

	am.Role = aux.Role

	if len(aux.Content) == 0 || string(aux.Content) == "null" {
		am.Content = StringPart("")
		return nil
	}

	var strContent string
	err := json.Unmarshal(aux.Content, &strContent)

	if err == nil {
		am.Content = StringPart(strContent)
		return nil
	}

	return fmt.Errorf("assistantmessage unmarshal: 'content' field was not a JSON string and could not be parsed as StringPart (e.g. for text content). Unmarshal error: %w. Raw 'content' data: %s", err, string(aux.Content))
}
