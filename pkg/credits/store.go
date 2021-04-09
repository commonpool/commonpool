package credits

import (
	"cp/pkg/api"
	"cp/pkg/utils"
	"gorm.io/gorm"
)

type Store interface {
	Create(credits *api.Credits) error
	GetAllForGroup(groupID string) ([]*api.Credits, error)
}

type CreditStore struct {
	db *gorm.DB
}

func NewCreditStore(db *gorm.DB) *CreditStore {
	return &CreditStore{db: db}
}

func (s *CreditStore) Create(credits *api.Credits) error {
	return s.db.Create(credits).Error
}

func (s *CreditStore) GetAllForGroup(groupID string) ([]*api.Credits, error) {
	var results []*api.Credits
	if err := s.db.Model(&api.Credits{}).Find(&results, "group_id = ?", groupID).Error; err != nil {
		return nil, err
	}
	var allTargets []*api.Target
	for _, result := range results {
		allTargets = append(allTargets, result.SentTo)
		allTargets = append(allTargets, result.SentBy)
	}
	if err := utils.PopulateTargets(s.db, allTargets);err != nil {
		return nil, err
	}
	return results, nil
}
