package untisAPI

import (
	b64 "encoding/base64"
	"encoding/json"
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
