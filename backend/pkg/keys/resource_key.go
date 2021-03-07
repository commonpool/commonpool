package keys

import (
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/satori/go.uuid"
)

type ResourceKey struct {
	ID uuid.UUID
}

func NewResourceKey(id uuid.UUID) ResourceKey {
	return ResourceKey{
		ID: id,
	}
}
func GenerateResourceKey() ResourceKey {
	return ResourceKey{
		ID: uuid.NewV4(),
	}
}

func ParseResourceKey(key string) (ResourceKey, error) {
	resourceUuid, err := uuid.FromString(key)
	if err != nil {
		response := exceptions.ErrInvalidResourceKey(key)
		return ResourceKey{}, &response
	}
	resourceKey := ResourceKey{
		ID: resourceUuid,
	}
	return resourceKey, nil
}

func (r ResourceKey) GetUUID() uuid.UUID {
	return r.ID
}

func (r ResourceKey) String() string {
	return r.ID.String()
}

func (r ResourceKey) GetFrontendLink() string {
	return fmt.Sprintf("<commonpool-resource id='%s'><commonpool-resource>", r.String())
}

func (k ResourceKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.ID.String())
}

func (k *ResourceKey) UnmarshalJSON(data []byte) error {
	var uid string
	if err := json.Unmarshal(data, &uid); err != nil {
		return err
	}
	id, err := uuid.FromString(uid)
	if err != nil {
		return err
	}
	k.ID = id
	return nil
}
