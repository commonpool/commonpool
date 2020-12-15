package model

type Groups struct {
	Items []*Group
}

func NewGroups(groups []*Group) *Groups {
	return &Groups{
		Items: groups,
	}
}
