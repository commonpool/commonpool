package notifications

import (
	"cp/pkg/api"
	"gorm.io/gorm"
)

type Store interface {
	GetNotifications(userID string) ([]*api.Notification, error)
	ClearNotifications(userID string) error
	AddNotification(notification *api.Notification) error
	AddNotifications(notifications []*api.Notification) error
	GetUnreadCount(userID string) (int, error)
}

type NotificationStore struct {
	db *gorm.DB
}

func NewNotificationStore(db *gorm.DB) *NotificationStore {
	return &NotificationStore{db: db}
}

func (n *NotificationStore) GetNotifications(userID string) ([]*api.Notification, error) {
	var result []*api.Notification
	if err := n.db.Model(&api.Notification{}).
		Order("created_at desc").
		Find(&result, "user_id = ?", userID).
		Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (n *NotificationStore) ClearNotifications(userID string) error {
	return n.db.Where("user_id = ?").Delete(&api.Notification{}).Error
}

func (n *NotificationStore) AddNotification(notification *api.Notification) error {
	return n.db.Create(notification).Error
}

func (n *NotificationStore) AddNotifications(notifications []*api.Notification) error {
	if len(notifications) == 0 {
		return nil
	}
	return n.db.Create(notifications).Error
}

func (n *NotificationStore) GetUnreadCount(userID string) (int, error) {
	var result int64
	if err := n.db.Model(&api.Notification{}).Where("user_id = ?", userID).Count(&result).Error; err != nil {
		return 0, err
	}
	return int(result), nil
}
