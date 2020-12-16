package chat

type OptionObject struct {
	Text        BlockElement  `json:"text,omitempty"`
	Value       string        `json:"value,omitempty"`
	Description *BlockElement `json:"description,omitempty"`
}

func NewOptionObject(text BlockElement, value string, description *BlockElement) OptionObject {
	return OptionObject{
		Text:        text,
		Value:       value,
		Description: description,
	}
}
