package model

type MessageType string

const (
	NormalMessage MessageType = "message"
)

type MessageSubType string

const (
	UserMessage MessageSubType = "user"
	BotMessage  MessageSubType = "bot"
)

type SurfaceType string

const (
	Modals   SurfaceType = "modals"
	Messages SurfaceType = "messages"
	Home     SurfaceType = "home"
	Tabs     SurfaceType = "tabs"
)

type BlockType string

const (
	Actions BlockType = "actions"
	Context BlockType = "context"
	Divider BlockType = "divider"
	File    BlockType = "file"
	Header  BlockType = "header"
	Image   BlockType = "image"
	Input   BlockType = "input"
	Section BlockType = "section"
)

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

type ButtonStyle string

const (
	Primary ButtonStyle= "primary"
	Danger  ButtonStyle= "danger"
)

type FileSource string

const (
	Remote FileSource = "remote"
)

type ConfirmObject struct {
	Title   BlockElement `json:"title,omitempty"`
	Text    BlockElement `json:"text,omitempty"`
	Confirm BlockElement `json:"confirm,omitempty"`
	Deny    BlockElement `json:"deny,omitempty"`
	Style   ButtonStyle  `json:"style,omitempty"`
}

type OptionObject struct {
	Text        BlockElement  `json:"text,omitempty"`
	Value       string        `json:"value,omitempty"`
	Description *BlockElement `json:"description,omitempty"`
}

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

type Block struct {
	Type       BlockType      `json:"type,omitempty"`
	Text       *BlockElement  `json:"text,omitempty"`
	Elements   []BlockElement `json:"elements,omitempty"`
	ImageURL   *string        `json:"imageUrl,omitempty"`
	AltText    *string        `json:"altText,omitempty"`
	Title      *BlockElement  `json:"textObject,omitempty"`
	Fields     []BlockElement `json:"fields,omitempty"`
	Accessory  *BlockElement  `json:"accessory,omitempty"`
	BlockID    *string        `json:"blockId,omitempty"`
	ExternalID *string        `json:"externalId,omitempty"`
	Source     *FileSource    `json:"fileSource,omitempty"`
}

type Attachment struct {
	Color  string  `json:"color,omitempty"`
	Blocks []Block `json:"blocks"`
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

func NewButtonElement(text BlockElement, style *ButtonStyle, actionId *string, url *string, value *string, confirm *ConfirmObject) BlockElement {
	return BlockElement{
		Type:     ButtonElement,
		Text:     &text,
		ActionID: actionId,
		URL:      url,
		Value:    value,
		Style:    style,
		Confirm:  confirm,
	}
}

func NewOptionObject(text BlockElement, value string, description *BlockElement) OptionObject {
	return OptionObject{
		Text:        text,
		Value:       value,
		Description: description,
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

func NewActionBlock(elements []BlockElement, blockID *string) *Block {
	return &Block{
		Type:     Actions,
		Elements: elements,
		BlockID:  blockID,
	}
}

func NewContextBlock(elements []BlockElement, blockID *string) *Block {
	return &Block{
		Type:     Context,
		Elements: elements,
		BlockID:  blockID,
	}
}

func NewDividerBlock() *Block {
	return &Block{
		Type: Divider,
	}
}

func NewFileBlock(externalId string, source FileSource, blockID *string) *Block {
	return &Block{
		Type:       File,
		ExternalID: &externalId,
		Source:     &source,
		BlockID:    blockID,
	}
}

func NewHeaderBlock(text BlockElement, blockID *string) *Block {
	return &Block{
		Type:    Header,
		Text:    &text,
		BlockID: blockID,
	}
}

func NewImageBlock(imageURL string, altText string, title BlockElement, blockID *string) *Block {
	return &Block{
		Type:     Image,
		ImageURL: &imageURL,
		AltText:  &altText,
		Title:    &title,
		BlockID:  blockID,
	}
}

func NewSectionBlock(text BlockElement, fields []BlockElement, accessory *BlockElement, blockID *string) *Block {
	return &Block{
		Type:      Section,
		Text:      &text,
		Fields:    fields,
		Accessory: accessory,
		BlockID:   blockID,
	}
}
