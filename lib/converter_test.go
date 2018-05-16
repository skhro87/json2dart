package lib

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
)

func TestJson2Dart(t *testing.T) {
	example, err := ioutil.ReadFile("../example.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, err = Json2Dart(string(example), "Root")
	assert.Nil(t, err)
}


