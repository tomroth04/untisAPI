package types

import (
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"testing"
)

type Test struct {
	name string
	json string
	want interface{}
}

func RunTests(t *testing.T, tests []Test, assertionFunc func(t *testing.T, tt Test, lesson GenericLesson)) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lesson := GenericLesson{R: gjson.Parse(tt.json)}
			assertionFunc(t, tt, lesson)
		})
	}
}

func TestGenericLesson_GetSubject(t *testing.T) {
	tests := []Test{
		{name: "Lesson has subject", json: `{"su": [{"longname": "Mathematics"}]}`, want: "Mathematics"},
		{name: "Lesson does not have subject", json: `{"su": []}`, want: ""},
	}

	RunTests(t, tests, func(t *testing.T, tt Test, lesson GenericLesson) {
		assert.Equal(t, tt.want, lesson.GetSubject())
	})
}

func TestGenericLesson_IsCancelled(t *testing.T) {
	tests := []Test{
		{name: "Lesson is cancelled", json: `{"code": "cancelled"}`, want: true},
		{name: "Lesson is not cancelled", json: `{}`, want: false},
	}

	RunTests(t, tests, func(t *testing.T, tt Test, lesson GenericLesson) {
		assert.Equal(t, tt.want, lesson.IsCancelled())
	})
}

func TestGenericLesson_IsIrregular(t *testing.T) {
	tests := []Test{
		{name: "Lesson is irregular", json: `{"code": "irregular"}`, want: true},
		{name: "Lesson is not irregular", json: `{}`, want: false},
	}

	RunTests(t, tests, func(t *testing.T, tt Test, lesson GenericLesson) {
		assert.Equal(t, tt.want, lesson.IsIrregular())
	})
}

func TestGenericLesson_GetDate(t *testing.T) {
	tests := []Test{
		{name: "Lesson has date", json: `{"date": 20200101}`, want: "2020-01-01 01:00:00 +0100 CET"},
	}

	RunTests(t, tests, func(t *testing.T, tt Test, lesson GenericLesson) {
		assert.Equal(t, tt.want, lesson.GetDate().String())
	})
}

func TestGenericLesson_GetLessonId(t *testing.T) {
	tests := []Test{
		{name: "Lesson has id", json: `{"id": 123}`, want: 123},
		{name: "Lesson does not have id", json: `{}`, want: 0},
	}

	RunTests(t, tests, func(t *testing.T, tt Test, lesson GenericLesson) {
		assert.Equal(t, tt.want, lesson.GetLessonId())
	})
}
