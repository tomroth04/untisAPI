package types

import (
	"github.com/tidwall/gjson"
	"log/slog"
	"strconv"
	"time"
)

// GenericLesson is a generic lesson
// uses gjson.Result as data source as i didn't know the exact structure of the data
// and i didn't want the code to break if the structure changes or I missed something
type GenericLesson struct {
	R gjson.Result
}

// IsCancelled check if a lesson is marked as cancelled
func (g GenericLesson) IsCancelled() bool {
	if g.R.Get("code").Exists() && g.R.Get("code").String() == "cancelled" {
		return true
	}
	return false
}

// IsIrregular checks if a lesson is marked as irregular
func (g GenericLesson) IsIrregular() bool {
	if g.R.Get("code").Exists() && g.R.Get("code").String() == "irregular" {
		return true
	}
	return false
}

// GetLessonId gets the id of the lesson
func (g GenericLesson) GetLessonId() int {
	if !g.R.Get("id").Exists() {
		return 0
	}

	return int(g.R.Get("id").Int())
}

// GetSubject gets the subject of the lesson
func (g GenericLesson) GetSubject() string {
	su := g.R.Get("su")
	if !su.Exists() || !su.IsArray() || len(su.Array()) == 0 {
		return ""
	}

	return g.R.Get("su").Array()[0].Get("longname").String()
}

// GetDate gets the date of the lesson
func (g GenericLesson) GetDate() time.Time {
	t, err := ParseUntisDate(strconv.Itoa(int(g.R.Get("date").Int())))
	if err != nil {
		slog.Error("Error parsing date", "error", err)
		return time.Time{}
	}
	return t
}

// GetDateFormatted gets the date formatted
func (g GenericLesson) GetDateFormatted() string {
	return g.GetDate().Format("Monday, 02 January 2006")
}

// IsReplaced checks if the lesson has a substitute teacher
func (g GenericLesson) IsReplaced() bool {
	if !g.R.Get("te").Exists() {
		slog.Error("No teacher", "data", g.R.String())
		return false
	}

	for _, teacher := range g.R.Get("te").Array() {
		if teacher.Get("orgname").Exists() {
			return true
		}
	}

	return false
}

// GetStartTimeFormatted gets the start Time formatted
func (g GenericLesson) GetStartTimeFormatted() string {
	return getLessonTimeFromInteger(int(g.R.Get("startTime").Int()))
}

// GetLessonInfo gets the lesson information
func (g GenericLesson) GetLessonInfo() string {
	if !g.R.Get("lstext").Exists() {
		return ""
	}
	return g.R.Get("lstext").String()
}

func (g GenericLesson) GetSubstituteText() string {
	if !g.R.Get("substText").Exists() {
		return ""
	}
	return g.R.Get("substText").String()
}

func (g GenericLesson) IsEqual(b GenericLesson) bool {
	return g.R.String() == b.R.String()
}

func (g GenericLesson) GetActivityType() string {
	if !g.R.Get("activityType").Exists() {
		return ""
	}
	return g.R.Get("activityType").String()
}
