package model

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
