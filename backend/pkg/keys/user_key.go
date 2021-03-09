package keys

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap/zapcore"
)

type UserKey struct {
	subject string
}

func (u UserKey) GetUserKey() UserKey {
	return u
}

type UserKeyGetter interface {
	GetUserKey() UserKey
}

func (k UserKey) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("user_id", k.subject)
	return nil
}

var _ zapcore.ObjectMarshaler = &UserKey{}

func NewUserKey(subject string) UserKey {
	return UserKey{subject: subject}
}

func GenerateUserKey() UserKey {
	return NewUserKey(uuid.NewV4().String())
}

func (k UserKey) String() string {
	return k.subject
}

func (k UserKey) Target() *Target {
	return NewUserTarget(k)
}

func (k UserKey) StreamKey() StreamKey {
	return NewStreamKey("user", k.String())
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

func (k *UserKey) Scan(value interface{}) error {
	keyValue, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal string value:", value))
	}
	*k = NewUserKey(keyValue)
	return nil
}

func (k UserKey) Value() (driver.Value, error) {
	return driver.String.ConvertValue(k.subject)
}

func (k UserKey) GormDataType() string {
	return "varchar(128)"
}
