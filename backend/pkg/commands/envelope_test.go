package commands

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type omit *struct{}

type TestPayload struct {
	Bla int `json:"bla"`
}

type TestCommand struct {
	CommandEnvelope
	TestPayload `json:"payload"`
}

func TestEnvelope(t *testing.T) {

	cmd := &TestCommand{
		CommandEnvelope: CommandEnvelope{
			CommandTime:   time.Now(),
			CommandType:   "bla",
			AggregateID:   "123",
			AggregateType: "test",
		},
		TestPayload: TestPayload{
			Bla: 2,
		},
	}

	bytes, err := json.MarshalIndent(cmd, "", "  ")
	if !assert.NoError(t, err) {
		return
	}

	t.Log("\n" + string(bytes))

}
