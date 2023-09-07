package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestToJsonStr(t *testing.T) {
	assert.Equal(t, "{\"a\":1}", ToJsonStr(map[string]int{"a": 1}))
}

func TestGetDateUntisFormat(t *testing.T) {
	assert.Equal(t, "20221202", GetDateUntisFormat(time.Date(2022, 12, 2, 0, 0, 0, 0, time.Now().Location())))
}

func TestParseUntisDate(t *testing.T) {
	if date, err := ParseUntisDate("20221202"); err != nil {
		assert.NoError(t, err)
	} else {
		assert.Equal(t, time.Date(2022, 12, 2, 0, 0, 0, 0, time.UTC), date)
	}
}

func TestBase64(t *testing.T) {
	assert.Equal(t, "bG1ybA==", ToBase64("lmrl"))
}

func TestGetLessonTimeFromInteger(t *testing.T) {
	assert.Equal(t, "08:00", getLessonTimeFromInteger(800))
	assert.Equal(t, "12:00", getLessonTimeFromInteger(1200))
	assert.Equal(t, "", getLessonTimeFromInteger(0))
	assert.Equal(t, "", getLessonTimeFromInteger(10))
}

func TestStr(t *testing.T) {
	assert.Equal(t, "1", str(1))
}
