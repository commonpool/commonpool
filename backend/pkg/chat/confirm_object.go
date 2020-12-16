package chat

type ConfirmObject struct {
	Title   BlockElement `json:"title,omitempty"`
	Text    BlockElement `json:"text,omitempty"`
	Confirm BlockElement `json:"confirm,omitempty"`
	Deny    BlockElement `json:"deny,omitempty"`
	Style   ButtonStyle  `json:"style,omitempty"`
}
