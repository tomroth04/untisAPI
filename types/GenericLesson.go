package types

import (
	"github.com/tidwall/gjson"
	"log"
	"strconv"
	"time"
)

type GenericLesson struct {
	R gjson.Result
}

// IsCancelled check if some hours are cancelled
func (g GenericLesson) IsCancelled() bool {
	return g.R.Get("code").Exists()
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
	if !g.R.Get("su").Exists() {
		log.Println("error getting subject of: " + g.R.String())
	}

	return g.R.Get("su").Array()[0].Get("longname").String()
}

// GetDate gets the date of the lesson
func (g GenericLesson) GetDate() time.Time {
	t, err := ParseUntisDate(strconv.Itoa(int(g.R.Get("date").Int())))
	if err != nil {
		log.Println(err)
		return time.Time{}
	}
	return t
}

// GetDateFormatted gets the date formatted
func (g GenericLesson) GetDateFormatted() string {
	t, err := ParseUntisDate(strconv.Itoa(int(g.R.Get("date").Int())))
	if err != nil {
		log.Println(err)
		return ""
	}
	return t.Format("02 January 2006")
}

// IsReplaced checks if the lesson has a replacement teacher
func (g GenericLesson) IsReplaced() bool {
	if !g.R.Get("te").Exists() {
		log.Println("No teacher: " + g.R.String())
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

func (g GenericLesson) IsEqual(b GenericLesson) bool {
	return g.R.String() == b.R.String()
}