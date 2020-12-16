package chat

type BlockElement struct {
	Type           ElementType    `json:"type,omitempty"`
	Text           *BlockElement  `json:"text,omitempty"`
	ActionID       *string        `json:"actionId,omitempty"`
	URL            *string        `json:"url,omitempty"`
	Value          *string        `json:"value,omitempty"`
	Style          *ButtonStyle   `json:"style,omitempty"`
	Confirm        *ConfirmObject `json:"confirm,omitempty"`
	Options        []OptionObject `json:"options,omitempty"`
	InitialOptions []OptionObject `json:"initialOptions,omitempty"`
	Placeholder    *BlockElement  `json:"placeholder,omitempty"`
	InitialDate    *string        `json:"initialDate,omitempty"`
	ImageURL       *string        `json:"imageUrl,omitempty"`
	AltText        *string        `json:"altText,omitempty"`
	Emoji          *bool          `json:"emoji,omitempty"`
}

func NewMarkdownObject(text string) BlockElement {
	return BlockElement{
		Type:  MarkdownTextType,
		Value: &text,
	}
}

func NewPlainTextObject(text string) BlockElement {
	return BlockElement{
		Type:  PlainTextType,
		Value: &text,
	}
}

func NewButtonElement(text BlockElement, style *ButtonStyle, actionId *string, url *string, value *string, confirm *ConfirmObject) *BlockElement {
	return &BlockElement{
		Type:     ButtonElement,
		Text:     &text,
		ActionID: actionId,
		URL:      url,
		Value:    value,
		Style:    style,
		Confirm:  confirm,
	}
}

func NewCheckboxesElement(text BlockElement, options []OptionObject, initialOptions []OptionObject, actionId *string, confirm *ConfirmObject) BlockElement {
	return BlockElement{
		Type:           CheckboxesElement,
		Text:           &text,
		ActionID:       actionId,
		Confirm:        confirm,
		InitialOptions: initialOptions,
		Options:        options,
	}
}

func NewDatePickerElement(actionId *string, placeholder *BlockElement, initialDate *string, confirm *ConfirmObject) BlockElement {
	return BlockElement{
		Type:        DatepickerElement,
		ActionID:    actionId,
		Confirm:     confirm,
		Placeholder: placeholder,
		InitialDate: initialDate,
	}
}

func NewImageElement(imageUrl string, altText string) BlockElement {
	return BlockElement{
		Type:     ImageElement,
		ImageURL: &imageUrl,
		AltText:  &altText,
	}
}

func NewRadioButtonsElement(options []OptionObject, initialOptions []OptionObject, actionId *string, confirm *ConfirmObject) BlockElement {
	return BlockElement{
		Type:           RadioButtonsElement,
		Options:        options,
		InitialOptions: initialOptions,
		ActionID:       actionId,
		Confirm:        confirm,
	}
}
