package model

type Target struct {
	UserKey  *UserKey
	GroupKey *GroupKey
	Type     OfferItemTargetType
}

func (t Target) Equals(target *Target) bool {

	if t.Type != target.Type {
		return false
	}

	if t.Type == GroupTarget {
		return *t.GroupKey == *target.GroupKey
	}

	return *t.UserKey == *target.UserKey
}

func (t Target) IsForGroup() bool {
	return t.Type == GroupTarget
}

func (t Target) IsForUser() bool {
	return t.Type == UserTarget
}

func (t Target) GetGroupKey() GroupKey {
	return *t.GroupKey
}
func (t Target) GetUserKey() UserKey {
	return *t.UserKey
}

type Targets struct {
	Items []*Target
}

func NewTargets(items []*Target) *Targets {
	copied := make([]*Target, len(items))
	copy(copied, items)
	return &Targets{
		Items: copied,
	}
}

func NewEmptyTargets() *Targets {
	return &Targets{
		Items: []*Target{},
	}
}

type OfferItemTargetType string

const (
	UserTarget  OfferItemTargetType = "user"
	GroupTarget OfferItemTargetType = "group"
)
