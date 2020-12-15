package model

type Attachment struct {
	Color  string  `json:"color,omitempty"`
	Blocks []Block `json:"blocks"`
}
