package transaction

type Entries struct {
	Items []*Entry
}

func NewEntries(items []*Entry) *Entries {
	return &Entries{
		Items: items,
	}
}
