package service

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/utils"
	"sort"
	"strings"
)

// GetConversationChannelKey Returns the id of a conversation between users
// Only a single conversation can exist between a group of people.
// There can only be one conversation with Joe, Dana and Mark.
// So the ID of the conversation is composed of the
// sorted IDs of its participants.
func (c ChatService) GetConversationChannelKey(ctx context.Context, participants *keys.UserKeys) (keys.ChannelKey, error) {

	if participants == nil || len(participants.Items) == 0 {
		err := fmt.Errorf("cannot get conversation channel for 0 participants")
		return keys.ChannelKey{}, err
	}

	var shortUids []string
	for _, participant := range participants.Items {
		sid, err := utils.ShortUuidFromStr(participant.String())
		if err != nil {
			return keys.ChannelKey{}, err
		}
		shortUids = append(shortUids, sid)
	}
	sort.Strings(shortUids)
	channelId := strings.Join(shortUids, "-")
	channelKey := keys.NewConversationKey(channelId)

	return channelKey, nil
}
