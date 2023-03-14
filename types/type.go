package types

import (
	"fmt"
	"strconv"
)

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

func (c Class) String() string {
	return fmt.Sprintf(
		"%s (%s) Id: %s teacher1:%s",
		c.Name, c.LongName, strconv.Itoa(c.Id), strconv.Itoa(c.Teacher1),
	)
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

type LessonWithSubj struct {
	Id         int    `json:"id"`
	Subject    string `json:"subject"`
	LessonType string `json:"lessonType"`
}
