package types

import (
	b64 "encoding/base64"
	"encoding/json"
	"github.com/tidwall/gjson"
	"strconv"
	"time"
)

// ToJsonStr convert anything/interface{} to  string
func ToJsonStr(data any) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

// GetDateUntisFormat formats date to string
func GetDateUntisFormat(date time.Time) string {
	return date.Format("20060102")
}

// ParseUntisDate parses a string date
func ParseUntisDate(date string) (time.Time, error) {
	return time.Parse("20060102", date)
}

// ToBase64 transform to base64
func ToBase64(str string) string {
	return b64.StdEncoding.EncodeToString([]byte(str))
}

// parse the hourly time from a lesson
func getLessonTimeFromInteger(i int) string {
	n := strconv.Itoa(i)
	if len(n) == 3 {
		// parse time like 800
		h, _ := strconv.Atoi(n[0:1])
		m, _ := strconv.Atoi(n[1:3])
		return time.Date(0, 0, 0, h, m, 0, 0, time.UTC).Format("15:04")
	} else if len(n) == 4 {
		h, _ := strconv.Atoi(n[0:2])
		m, _ := strconv.Atoi(n[2:4])
		return time.Date(0, 0, 0, h, m, 0, 0, time.UTC).Format("15:04")
	}
	return ""
}

func str(n int) string {
	return strconv.Itoa(n)
}

// TransformResultLesson converts gjson.Result to GenericLesson
func TransformResultLesson(res []gjson.Result) []GenericLesson {
	result := make([]GenericLesson, len(res))
	for i := 0; i < len(res); i++ {
		result[i].R = res[i]
	}
	return result
}

// GetLessonMap returns a map of lesson id and lesson subject
func GetLessonMap(subjs []LessonWithSubj) (map[int]string, error) {
	lessons := make(map[int]string)

	for _, subj := range subjs {
		lessons[subj.Id] = subj.Subject
	}

	return lessons, nil
}
