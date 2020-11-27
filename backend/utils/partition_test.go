package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPartition(t *testing.T) {

	slice := []int{1, 2, 3, 4, 5}

	i := -1
	err := Partition(len(slice), 2, func(i1 int, i2 int) error {
		i++
		if i == 0 {
			assert.Equal(t, 0, i1)
			assert.Equal(t, 1, i2)
		}
		if i == 1 {
			assert.Equal(t, 2, i1)
			assert.Equal(t, 3, i2)
		}
		if i == 2 {
			assert.Equal(t, 4, i1)
			assert.Equal(t, 4, i2)
		}
		assert.NotEqual(t, 3, i)

		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, i)

}
