package untisApi

import (
	b64 "encoding/base64"
	"encoding/json"
	"strconv"
	"time"
)

func toJsonStr(data any) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func toJson(data any) []byte {
	jsonData, _ := json.Marshal(data)
	return jsonData
}

func GetDateUntisFormat(date time.Time) string {
	return date.Format("20060102")
}

func ParseUntisDate(date string) (time.Time, error) {
	return time.Parse("20060102", date)
}

func MustParseUntisDate(date string) time.Time {
	t, _ := ParseUntisDate(date)
	return t
}

func ToBase64(str string) string {
	return b64.StdEncoding.EncodeToString([]byte(str))
}

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
