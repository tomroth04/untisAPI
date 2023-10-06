package types

import (
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"testing"
)

func TestGenericLesson_GetSubject(t *testing.T) {
	assert.Equal(t, "Mathematik", GenericLesson{R: gjson.Parse(`{"su":[{"longname":"Mathematik"}]}`)}.GetSubject())
}

func TestGenericLesson_GetDate(t *testing.T) {
	d := GenericLesson{R: gjson.Parse(`{"date":20221202}`)}.GetDate()
	assert.Equal(t, 2022, d.Year())
	assert.Equal(t, 12, int(d.Month()))
	assert.Equal(t, 2, d.Day())
}
