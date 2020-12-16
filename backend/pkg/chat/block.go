package chat

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
