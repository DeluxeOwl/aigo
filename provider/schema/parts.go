package schema

type ContentPartType string

const (
	ContentPartTypeText       ContentPartType = "text"
	ContentPartTypeData       ContentPartType = "data"
	ContentPartTypeURL        ContentPartType = "url"
	ContentPartTypeToolCall   ContentPartType = "tool-call"
	ContentPartTypeToolResult ContentPartType = "tool-result"
)

//sumtype:decl
type ContentPartUser interface {
	isContentPartUser()
}

type TextPart struct {
	Type ContentPartType `json:"type"`
	Text string          `json:"text"`
}

func (TextPart) isContentPartUser() {}
func NewTextPart(text string) *TextPart {
	return &TextPart{
		Type: ContentPartTypeText,
		Text: text,
	}
}

//sumtype:decl
type ContentPartAssistant interface {
	isContentPartAssistant()
}

type StringPart string

func (StringPart) isContentPartAssistant() {}
func (s StringPart) String() string        { return string(s) }
