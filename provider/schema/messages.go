package schema

type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
	MessageRoleSystem    MessageRole = "system"
	MessageRoleTool      MessageRole = "tool"
)

//sumtype:decl
type Message interface {
	isMessage()
}

type SystemMessage struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`
}

func (SystemMessage) isMessage() {}

func NewSystemMessage(message string) *SystemMessage {
	return &SystemMessage{
		Role:    MessageRoleSystem,
		Content: message,
	}
}

type UserMessage struct {
	Role    MessageRole       `json:"role"`
	Content []ContentPartUser `json:"content"`
}

func (UserMessage) isMessage() {}

func NewUserMessage(content []ContentPartUser) *UserMessage {
	return &UserMessage{
		Role:    MessageRoleUser,
		Content: content,
	}
}

// ....

type AssistantMessage struct {
	Role    MessageRole          `json:"role"`
	Content ContentPartAssistant `json:"content"`
}

func (*AssistantMessage) isMessage() {}
