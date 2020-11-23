package store

import (
	"encoding/json"
	"errors"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/utils"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"sort"
	"strings"
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
		Model(model.Message{}).
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

func (cs *ChatStore) CreateThread(thread model.Thread) error {
	return cs.db.Create(thread).Error
}

func (cs *ChatStore) SendMessageToThread(sendMessageRequest *chat.SendMessageToThreadRequest) *chat.SendMessageToThreadResponse {

	lastMsgChars := utils.FirstChars(sendMessageRequest.Text, 20)
	authorKey := sendMessageRequest.FromUser
	authorId := authorKey.String()
	sentAt := time.Now()

	thread, err := cs.GetThread(sendMessageRequest.ThreadKey)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		thread = &model.Thread{
			UserID:              sendMessageRequest.ThreadKey.UserKey.String(),
			TopicID:             sendMessageRequest.ThreadKey.TopicKey.String(),
			CreatedAt:           time.Now(),
			LastMessageAt:       time.Unix(0, 0),
			LastTimeRead:        time.Unix(0, 0),
			LastMessageChars:    "",
			LastMessageUserId:   "",
			LastMessageUserName: "",
		}
		err := cs.CreateThread(*thread)
		if err != nil {
			return &chat.SendMessageToThreadResponse{
				Error: err,
			}
		}
	} else if err != nil {
		return &chat.SendMessageToThreadResponse{
			Error: err,
		}
	}

	recipientThreadKey := thread.GetKey()
	isAuthorThread := authorKey == recipientThreadKey.UserKey

	// Updating the thread with the last message
	threadUpdateValues := map[string]interface{}{
		"last_message_at":        sentAt,
		"last_message_chars":     lastMsgChars,
		"last_message_user_id":   authorId,
		"last_message_user_name": sendMessageRequest.FromUserName,
	}
	if isAuthorThread {
		threadUpdateValues["last_time_read"] = sentAt
	}

	err = cs.db.Model(&model.Thread{}).
		Where("topic_id = ? AND user_id = ?", recipientThreadKey.TopicKey.String(), recipientThreadKey.UserKey.String()).
		Updates(threadUpdateValues).Error

	if err != nil {
		return &chat.SendMessageToThreadResponse{
			Error: err,
		}
	}

	blocks := sendMessageRequest.Blocks
	blocksJson, err := getBlocksJson(blocks)
	if err != nil {
		return &chat.SendMessageToThreadResponse{
			Error: err,
		}
	}

	attachments := sendMessageRequest.Attachments
	attachmentsJson, err := getAttachmentsJson(attachments)
	if err != nil {
		return &chat.SendMessageToThreadResponse{
			Error: err,
		}
	}

	// Building the message
	senderMessage := model.Message{
		ID:             uuid.NewV4(),
		UserID:         recipientThreadKey.UserKey.String(),
		TopicID:        recipientThreadKey.TopicKey.ID,
		SentAt:         sentAt,
		Blocks:         blocksJson,
		Attachments:    attachmentsJson,
		Text:           sendMessageRequest.Text,
		MessageType:    model.NormalMessage,
		MessageSubType: model.UserMessage,
		IsPersonal:     true,
		SentBy:         sendMessageRequest.FromUser.String(),
		SentByUsername: sendMessageRequest.FromUserName,
	}

	err = cs.db.Create(senderMessage).Error

	return &chat.SendMessageToThreadResponse{
		Error: err,
	}

}

func (cs *ChatStore) GetOrCreateConversationTopic(request *chat.GetOrCreateConversationTopicRequest) *chat.GetOrCreateConversationTopicResponse {

	var shortUids []string
	for _, key := range request.ParticipantList.Items {
		sid, err := utils.ShortUuidFromStr(key.String())
		if err != nil {
			return &chat.GetOrCreateConversationTopicResponse{
				Error: err,
			}
		}
		shortUids = append(shortUids, sid)
	}

	sort.Strings(shortUids)
	convoId := strings.Join(shortUids, "-")
	topicKey := model.NewTopicKey(convoId)

	var topic model.Topic
	found := cs.db.Model(model.Topic{}).Where("id = ?", topicKey.String()).First(&topic)
	if errors.Is(found.Error, gorm.ErrRecordNotFound) {
		topic := model.Topic{
			ID:    topicKey.String(),
			Title: "",
		}
		err := cs.db.Model(model.Topic{}).Create(topic).Error
		if err != nil {
			return &chat.GetOrCreateConversationTopicResponse{
				Error: err,
			}
		}

		var threads []model.Thread
		for _, key := range request.ParticipantList.Items {
			thread := model.Thread{
				UserID:              key.String(),
				TopicID:             topicKey.String(),
				CreatedAt:           time.Now(),
				LastMessageAt:       time.Unix(0, 0),
				LastTimeRead:        time.Unix(0, 0),
				LastMessageChars:    "",
				LastMessageUserId:   "",
				LastMessageUserName: "",
			}
			threads = append(threads, thread)
		}

		err = cs.db.Model(model.Thread{}).Create(threads).Error
		if err != nil {
			return &chat.GetOrCreateConversationTopicResponse{
				Error: err,
			}
		}

	}

	return &chat.GetOrCreateConversationTopicResponse{
		TopicKey:        topicKey,
		ParticipantList: request.ParticipantList,
		Error:           nil,
	}

}

func (cs *ChatStore) SendMessage(sendMessageRequest *chat.SendMessageRequest) *chat.SendMessageResponse {

	lastMsgChars := utils.FirstChars(sendMessageRequest.Text, 20)
	topicKey := sendMessageRequest.TopicKey
	topic := topicKey.String()
	authorKey := sendMessageRequest.FromUser
	authorId := authorKey.String()
	sentAt := time.Now()

	err := cs.db.Transaction(func(tx *gorm.DB) error {

		// Finding threads attached to topic
		var recipientThreadsForTopic []model.Thread
		err := tx.
			Where("topic_id = ?", topic).
			Find(&recipientThreadsForTopic).
			Error
		if err != nil {
			return err
		}

		blocks := sendMessageRequest.Blocks
		blocksJson, err := getBlocksJson(blocks)
		if err != nil {
			return err
		}

		attachments := sendMessageRequest.Attachments
		attachmentsJson, err := getAttachmentsJson(attachments)
		if err != nil {
			return err
		}

		// Looping through each thread
		for _, recipientThreadForTopic := range recipientThreadsForTopic {

			recipientThreadKey := recipientThreadForTopic.GetKey()
			isAuthorThread := authorKey == recipientThreadKey.UserKey

			// Updating the thread with the last message
			threadUpdateValues := map[string]interface{}{
				"last_message_at":        sentAt,
				"last_message_chars":     lastMsgChars,
				"last_message_user_id":   authorId,
				"last_message_user_name": sendMessageRequest.FromUserName,
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

			// Building the message
			senderMessage := model.Message{
				ID:             uuid.NewV4(),
				UserID:         recipientThreadKey.UserKey.String(),
				TopicID:        recipientThreadKey.TopicKey.ID,
				SentAt:         sentAt,
				Blocks:         blocksJson,
				Attachments:    attachmentsJson,
				Text:           sendMessageRequest.Text,
				MessageType:    model.NormalMessage,
				MessageSubType: model.UserMessage,
				IsPersonal:     false,
				SentBy:         sendMessageRequest.FromUser.String(),
				SentByUsername: sendMessageRequest.FromUserName,
			}

			err = tx.Create(senderMessage).Error
			if err != nil {
				return err
			}

		}

		return nil

	})

	return &chat.SendMessageResponse{
		Error: err,
	}

}

func getAttachmentsJson(attachments []model.Attachment) (string, error) {
	attachmentsJson := "[]"
	if attachments != nil {
		attachmentBytes, err := json.Marshal(attachments)
		if err != nil {
			return "", err
		}
		attachmentsJson = string(attachmentBytes)
	}
	return attachmentsJson, nil
}

func getBlocksJson(blocks []model.Block) (string, error) {
	blocksJson := "[]"
	if blocks != nil {
		blocksBytes, err := json.Marshal(blocks)
		if err != nil {
			return "", err
		}
		blocksJson = string(blocksBytes)
	}
	return blocksJson, nil
}

func (cs *ChatStore) GetTopic(key model.TopicKey) (*model.Topic, error) {
	var topic model.Topic
	err := cs.db.Where("id = ?", key.ID).First(&topic).Error
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
				TopicId:    uuid.NewV4().String(),
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
