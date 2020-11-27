package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMin(t *testing.T) {

	assert.Equal(t, 2, Min(2, 3))
	assert.Equal(t, 2, Min(2, 2))

}
