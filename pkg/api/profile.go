package api

type Profile struct {
	ID       string   `json:"sub"`
	Email    string   `json:"email"`
	Username string   `json:"preferred_username"`
	Groups   []string `json:"groups"`
}

func (p *Profile) IsInGroup(groupName string) bool {
	for _, group := range p.Groups {
		if group == groupName {
			return true
		}
	}
	return false
}
