package lib

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestJson2Dart(t *testing.T) {
	example := `
		{
			"hello": "world",
			"what is": 2
		}
	`

	err := Json2DartFile(example, "haha", "")
	assert.Nil(t, err)
}
