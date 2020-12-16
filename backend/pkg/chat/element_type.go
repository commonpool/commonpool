package chat

type ElementType string

const (
	PlainTextType         ElementType = "plain_text"
	MarkdownTextType      ElementType = "mrkdwn"
	ButtonElement         ElementType = "button"
	PlainTextInputElement ElementType = "plain_text_input"
	ImageElement          ElementType = "image"
	CheckboxesElement     ElementType = "checkboxes"
	DatepickerElement     ElementType = "datepicker"
	RadioButtonsElement   ElementType = "radio_buttons"
)
