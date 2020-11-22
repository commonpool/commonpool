package store

import (
	"errors"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/utils"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type ChatStore struct {
	db *gorm.DB
}

var _ chat.Store = &ChatStore{}

func (cs *ChatStore) GetLatestThreads(userKey model.UserKey, take int, skip int) ([]model.Thread, error) {
	var threads []model.Thread
	err := cs.db.
		Where("user_id = ?", userKey.String()).
		Order("last_message_at desc").
		Offset(skip).
		Limit(take).
		Find(&threads).
		Error

	if err != nil {
		return nil, err
	}
	return threads, nil
}

func (cs *ChatStore) GetThreadMessages(threadKey model.ThreadKey, take int, skip int) ([]model.Message, error) {
	var messages []model.Message
	err := cs.db.
		Where("topic_id = ? AND user_id = ?", threadKey.TopicKey.String(), threadKey.UserKey.String()).
		Order("sent_at desc").
		Offset(skip).
		Limit(take).
		Find(&messages).
		Error
	if err != nil {
		return nil, err
	}

	if len(messages) > 0 {
		lastMessageTs := messages[0].SentAt
		err = cs.db.Model(&model.Thread{}).
			Where("topic_id = ? AND user_id = ?", threadKey.TopicKey.String(), threadKey.UserKey.String(), lastMessageTs).
			Update("last_time_read", lastMessageTs).
			Error
		if err != nil {
			return nil, err
		}
	}

	return messages, nil
}

func (cs *ChatStore) GetThread(threadKey model.ThreadKey) (*model.Thread, error) {
	thread := &model.Thread{}
	err := cs.db.First(thread, "topic_id = ? and user_id = ?", threadKey.TopicKey.String(), threadKey.UserKey.String()).Error
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (cs *ChatStore) SendMessage(author model.UserKey, authorUserName string, topic model.TopicKey, content string) error {

	lastMsgChars := utils.FirstChars(content, 20)

	return cs.db.Transaction(func(tx *gorm.DB) error {

		sentAt := time.Now()

		var recipientThreadsForTopic []model.Thread
		err := tx.
			Where("topic_id = ?", topic.ID.String()).
			Find(&recipientThreadsForTopic).
			Error

		if err != nil {
			return err
		}

		for _, recipientThreadForTopic := range recipientThreadsForTopic {

			recipientThreadKey := recipientThreadForTopic.GetKey()
			isAuthorThread := author == recipientThreadKey.UserKey

			threadUpdateValues := map[string]interface{}{
				"last_message_at":        sentAt,
				"last_message_chars":     lastMsgChars,
				"last_message_user_id":   author.String(),
				"last_message_user_name": authorUserName,
			}
			if isAuthorThread {
				threadUpdateValues["last_time_read"] = sentAt
			}

			err = tx.Model(&model.Thread{}).
				Where("topic_id = ? AND user_id = ?", recipientThreadKey.TopicKey.String(), recipientThreadKey.UserKey.String()).
				Updates(threadUpdateValues).Error

			if err != nil {
				return err
			}

			senderMessage := model.Message{
				ID:       uuid.NewV4(),
				UserID:   recipientThreadKey.UserKey.String(),
				TopicId:  recipientThreadKey.TopicKey.ID,
				SentAt:   sentAt,
				Content:  content,
				AuthorID: author.String(),
			}

			err = tx.Create(senderMessage).Error
			if err != nil {
				return err
			}

		}

		return nil

	})
}

func (cs *ChatStore) GetTopic(key model.TopicKey) (*model.Topic, error) {
	var topic model.Topic
	err := cs.db.Where("id = ?", key.ID.String()).First(&topic).Error
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

func (cs *ChatStore) GetOrCreateResourceTopicMapping(rk model.ResourceKey, uk model.UserKey, rs resource.Store) (*model.ResourceTopic, error) {

	var rt model.ResourceTopic

	err := cs.db.First(&rt, "resource_id = ? AND user_id = ?", rk.String(), uk.String()).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {

		err = cs.db.Transaction(func(tx *gorm.DB) error {

			getResourceByKeyResponse := rs.GetByKey(resource.NewGetResourceByKeyQuery(rk))
			if getResourceByKeyResponse.Error != nil {
				return getResourceByKeyResponse.Error
			}
			res := getResourceByKeyResponse.Resource

			newResourceTopic := model.ResourceTopic{
				TopicId:    uuid.NewV4(),
				UserId:     uk.String(),
				ResourceId: rk.GetUUID(),
			}

			err = cs.db.Create(newResourceTopic).Error
			if err != nil {
				return err
			}

			newTopic := model.Topic{
				ID:    newResourceTopic.TopicId,
				Title: "About " + res.Summary,
			}
			err := cs.db.Create(newTopic).Error
			if err != nil {
				return err
			}

			createdAt := time.Now()

			inquirerThread := model.Thread{
				TopicID:   newResourceTopic.TopicId,
				UserID:    uk.String(),
				CreatedAt: createdAt,
			}

			ownerKey := res.GetUserKey()
			ownerThread := model.Thread{
				TopicID:   newResourceTopic.TopicId,
				UserID:    ownerKey.String(),
				CreatedAt: createdAt,
			}

			err = cs.db.Create(inquirerThread).Error
			if err != nil {
				return err
			}

			err = cs.db.Create(ownerThread).Error
			if err != nil {
				return err
			}

			rt = newResourceTopic

			return nil

		})

		return &rt, err

	}

	if err != nil {
		return nil, err
	}

	return &rt, nil
}

func NewChatStore(db *gorm.DB) *ChatStore {
	return &ChatStore{
		db: db,
	}
}
