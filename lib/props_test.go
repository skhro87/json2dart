package lib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_childObjectClassNameFromPropName(t *testing.T) {
	cases := map[string]string{
		"home_town":            "HomeTown",
		"glasses":              "Glass",
		"example_descriptions": "ExampleDescription",
	}

	for in, exp := range cases {
		assert.Equal(t, exp, childObjectClassNameFromPropName(in))
	}
}
