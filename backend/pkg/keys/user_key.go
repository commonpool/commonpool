package keys

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap/zapcore"
)

type UserKey struct {
	subject string
}

func (k UserKey) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("user_id", k.subject)
	return nil
}

var _ zapcore.ObjectMarshaler = &UserKey{}

func NewUserKey(subject string) UserKey {
	return UserKey{subject: subject}
}

func (k UserKey) String() string {
	return k.subject
}

func (k UserKey) GetExchangeName() string {
	return "users." + k.String()
}

func (k UserKey) GetFrontendLink() string {
	return fmt.Sprintf("<commonpool-user id='%s'></commonpool-user>", k.String())
}

func (k UserKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.subject)
}

func (k *UserKey) UnmarshalJSON(data []byte) error {
	var uid string
	if err := json.Unmarshal(data, &uid); err != nil {
		return err
	}
	k.subject = uid
	return nil
}
