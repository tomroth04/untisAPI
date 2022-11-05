package untisApi

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"
)

type Time time.Time

func (t *Time) UnmarshalJSON(s []byte) (err error) {
	r := string(s)
	d, err := ParseUntisDate(r)
	if err != nil {
		return err
	}
	*(*time.Time)(t) = d

	return nil
}

func (t *Time) Before(u Time) bool {
	return time.Time(*t).Before(time.Time(u))
}

func (t Time) String() string {
	return time.Time(t).Format("02 January 2006")
}

type SessionInformation struct {
	SessionId  string `json:"sessionId"`
	PersonType int    `json:"personType"`
	PersonId   int    `json:"personId"`
	ClassId    int    `json:"ClassId"`
}

type SchoolYear struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	StartDate Time   `json:"startDate"`
	EndDate   Time   `json:"endDate"`
}

type Class struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	LongName string `json:"longName"`
	Active   bool   `json:"active"`
	Teacher1 int    `json:"teacher1"`
}

type Absences struct {
	Absence []struct {
		Id            int           `json:"id"`
		StartDate     int           `json:"startDate"`
		EndDate       int           `json:"endDate"`
		StartTime     int           `json:"startTime"`
		EndTime       int           `json:"endTime"`
		CreateDate    int64         `json:"createDate"`
		LastUpdate    int64         `json:"lastUpdate"`
		CreatedUser   string        `json:"createdUser"`
		UpdatedUser   string        `json:"updatedUser"`
		ReasonId      int           `json:"reasonId"`
		Reason        string        `json:"reason"`
		Text          string        `json:"text"`
		Interruptions []interface{} `json:"interruptions"`
		CanEdit       bool          `json:"canEdit"`
		StudentName   string        `json:"studentName"`
		ExcuseStatus  string        `json:"excuseStatus"`
		IsExcused     bool          `json:"isExcused"`
		Excuse        struct {
			Id           int    `json:"id"`
			Text         string `json:"text"`
			ExcuseDate   int    `json:"excuseDate"`
			ExcuseStatus string `json:"excuseStatus"`
			IsExcused    bool   `json:"isExcused"`
			UserId       int    `json:"userId"`
			Username     string `json:"username"`
		} `json:"excuse"`
	} `json:"absences"`
	AbsenceReasons          []interface{} `json:"absenceReasons"`
	ExcuseStatuses          interface{}   `json:"excuseStatuses"`
	ShowAbsenceReasonChange bool          `json:"showAbsenceReasonChange"`
	ShowCreateAbsence       bool          `json:"showCreateAbsence"`
}

type Subject struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	LongName      string `json:"longName"`
	AlternateName string `json:"alternateName"`
	Active        bool   `json:"active"`
}

type TimeGridLesson struct {
	Day       int `json:"day"`
	TimeUnits []struct {
		Name      string     `json:"name"`
		StartTime LessonTime `json:"startTime"`
		EndTime   LessonTime `json:"endTime"`
	} `json:"timeUnits"`
}

type LessonTime time.Time

func (t *LessonTime) UnmarshalJSON(b []byte) error {
	n := string(b)

	if len(n) == 3 {
		// parse time like 800
		h, _ := strconv.Atoi(n[0:1])
		m, _ := strconv.Atoi(n[1:3])
		*(*time.Time)(t) = time.Date(0, 0, 0, h, m, 0, 0, time.UTC)
	} else if len(n) == 4 {
		h, _ := strconv.Atoi(n[0:2])
		m, _ := strconv.Atoi(n[2:4])
		*(*time.Time)(t) = time.Date(0, 0, 0, h, m, 0, 0, time.UTC)
	}
	return nil
}

func (t LessonTime) toTime() time.Time {
	return time.Time(t)
}

func (t LessonTime) String() string {
	return t.toTime().Format("15:04")
}

type Lesson map[string]any

func (l *Lesson) UnmarshalJSON(b []byte) error {
	var unmarshalled map[string]any

	if err := json.Unmarshal(b, &unmarshalled); err != nil {
		return err
	}

	*(*map[string]any)(l) = unmarshalled
	return nil
}

// Check if some hours are cancalled or something similar
func (l Lesson) IsCancelled() bool {
	_, ok := l["code"]
	return ok
}

func (l Lesson) GetLessonId() int {
	if reflect.TypeOf(l["id"]) == nil {
		fmt.Println(l)
	}
	return int(l["id"].(float64))
}

func (l Lesson) GetSubject() string {
	if reflect.TypeOf(l["su"]) == reflect.TypeOf([]interface{}{}) {
		if reflect.TypeOf(l["su"]) == reflect.TypeOf([]map[any]any{}) {
			return l["su"].([]map[string]any)[0]["longname"].(string)
		}

		return ""
	}
	return l["su"].(map[string]any)["longname"].(string)
}

func (l Lesson) GetDateFormatted() string {
	t, err := ParseUntisDate(strconv.Itoa(int(l["date"].(float64))))
	if err != nil {
		log.Println(err)
		return ""
	}
	return t.Format("02 January 2006")
}

func (l Lesson) GetStartTimeFormatted() string {
	return getLessonTimeFromInteger(int(l["startTime"].(float64)))
}
