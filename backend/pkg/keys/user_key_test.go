package keys

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshal(t *testing.T) {

	var userKey = NewUserKey("abc")
	bytes, _ := json.Marshal(userKey)

	var decoded UserKey
	_ = json.Unmarshal(bytes, &decoded)

	assert.Equal(t, "abc", decoded.subject)

}
