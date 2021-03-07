package domain

import "github.com/commonpool/backend/pkg/keys"

type ResourceSharing struct {
	GroupKey keys.GroupKey `json:"group_key"`
}
