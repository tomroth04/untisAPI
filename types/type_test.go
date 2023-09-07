package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClass_String(t *testing.T) {
	assert.Equal(t, "2CD1 (2CD1) Id: 1 teacher1:2", Class{
		Id:       1,
		Name:     "2CD1",
		LongName: "2CD1",
		Active:   true,
		Teacher1: 2,
	}.String())
}
