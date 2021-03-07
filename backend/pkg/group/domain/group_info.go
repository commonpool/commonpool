package domain

type GroupInfo struct {
	Name        string
	Description string
}

func NewGroupInfo() GroupInfo {
	return GroupInfo{}
}

func (g GroupInfo) WithName(name string) GroupInfo {
	return GroupInfo{
		Name:        name,
		Description: g.Description,
	}
}

func (g GroupInfo) WithDescription(description string) GroupInfo {
	return GroupInfo{
		Name:        g.Name,
		Description: description,
	}
}
