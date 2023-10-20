package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLessonTime_UnmarshalJSON(t *testing.T) {
	var l LessonTime
	err := l.UnmarshalJSON([]byte("800"))
	assert.NoError(t, err)
	assert.Equal(t, LessonTime(time.Date(0, 0, 0, 8, 0, 0, 0, time.UTC)), l)
	err = l.UnmarshalJSON([]byte("1200"))
	assert.NoError(t, err)
	assert.Equal(t, LessonTime(time.Date(0, 0, 0, 12, 0, 0, 0, time.UTC)), l)
	err = l.UnmarshalJSON([]byte("0230"))
	assert.NoError(t, err)
	assert.Equal(t, LessonTime(time.Date(0, 0, 0, 2, 30, 0, 0, time.UTC)), l)
}

func TestLessonTime_String(t *testing.T) {
	var l LessonTime
	err := l.UnmarshalJSON([]byte("800"))
	assert.NoError(t, err)
	assert.Equal(t, "08:00", l.String())
	err = l.UnmarshalJSON([]byte("1200"))
	assert.NoError(t, err)
	assert.Equal(t, "12:00", l.String())
	err = l.UnmarshalJSON([]byte("0230"))
	assert.NoError(t, err)
	assert.Equal(t, "02:30", l.String())
}
