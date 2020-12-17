package keys

import (
	"fmt"
	"github.com/commonpool/backend/pkg/utils"
	"go.uber.org/zap/zapcore"
	"sort"
	"strings"
)

type ChannelKey struct {
	ID string
}

var _ zapcore.ObjectMarshaler = &ChannelKey{}

func NewChannelKey(id string) ChannelKey {
	return ChannelKey{
		ID: id,
	}
}

func (tk ChannelKey) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("channel_id", tk.String())
	return nil
}

func (tk *ChannelKey) String() string {
	return tk.ID
}

func NewConversationKey(key string) ChannelKey {
	return ChannelKey{
		ID: key,
	}
}

func GetChannelKey(userKeys *UserKeys) (ChannelKey, error) {

	if userKeys == nil || len(userKeys.Items) == 0 {
		err := fmt.Errorf("cannot get conversation channel for 0 participants")
		return ChannelKey{}, err
	}

	var shortUids []string
	for _, userKey := range userKeys.Items {
		sid, err := utils.ShortUuidFromStr(userKey.String())
		if err != nil {
			return ChannelKey{}, err
		}
		shortUids = append(shortUids, sid)
	}
	sort.Strings(shortUids)
	channelId := strings.Join(shortUids, "-")
	channelKey := NewConversationKey(channelId)

	return channelKey, nil

}

func GetChannelKeyForGroup(groupKey GroupKey) ChannelKey {
	shortUid := utils.ShortUuid(groupKey.ID)
	return ChannelKey{
		ID: shortUid,
	}
}
