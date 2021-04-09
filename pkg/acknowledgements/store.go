package acknowledgements

import (
	"cp/pkg/api"
	"cp/pkg/utils"
	"gorm.io/gorm"
)

type Store interface {
	Save(acknowledgement *api.Acknowledgement) error
	GetForUser(userID string) ([]*api.Acknowledgement, error)
	GetForGroup(groupID string) ([]*api.Acknowledgement, error)
	GetAllInGroup(groupID string) ([]*api.Acknowledgement, error)
}

type AcknowledgementStore struct {
	db *gorm.DB
}

func NewAcknowledgementStore(db *gorm.DB) *AcknowledgementStore {
	return &AcknowledgementStore{db: db}
}

func (s *AcknowledgementStore) Save(acknowledgement *api.Acknowledgement) error {
	return s.db.Create(acknowledgement).Error
}

func (s *AcknowledgementStore) GetForUser(userID string) ([]*api.Acknowledgement, error) {
	var acknowledgements []*api.Acknowledgement
	if err := s.db.
		Model(&api.Acknowledgement{}).
		Find(&acknowledgements, "sent_to_user_id = ?", userID).
		Error; err != nil {
		return nil, err
	}

	var allTargets []*api.Target
	for _, acknowledgement := range acknowledgements {
		allTargets = append(allTargets, acknowledgement.SentBy)
		allTargets = append(allTargets, acknowledgement.SentTo)
	}
	if err := utils.PopulateTargets(s.db, allTargets); err != nil {
		return nil, err
	}
	return acknowledgements, nil
}

func (s *AcknowledgementStore) GetForGroup(groupID string) ([]*api.Acknowledgement, error) {
	var acknowledgements []*api.Acknowledgement
	if err := s.db.
		Model(&api.Acknowledgement{}).
		Find(&acknowledgements, "sent_to_group_id = ?", groupID).
		Error; err != nil {
		return nil, err
	}

	var allTargets []*api.Target
	for _, acknowledgement := range acknowledgements {
		allTargets = append(allTargets, acknowledgement.SentBy)
		allTargets = append(allTargets, acknowledgement.SentTo)
	}
	if err := utils.PopulateTargets(s.db, allTargets); err != nil {
		return nil, err
	}

	return acknowledgements, nil
}

func (s *AcknowledgementStore) GetAllInGroup(groupID string) ([]*api.Acknowledgement, error) {
	var acknowledgements []*api.Acknowledgement
	if err := s.db.
		Model(&api.Acknowledgement{}).
		Find(&acknowledgements, "group_id = ?", groupID).
		Error; err != nil {
		return nil, err
	}

	var allTargets []*api.Target
	for _, acknowledgement := range acknowledgements {
		allTargets = append(allTargets, acknowledgement.SentBy)
		allTargets = append(allTargets, acknowledgement.SentTo)
	}
	if err := utils.PopulateTargets(s.db, allTargets); err != nil {
		return nil, err
	}

	return acknowledgements, nil
}
