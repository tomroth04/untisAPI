package types

import (
	"log"
	"time"
)

type Homework struct {
	Id          int           `json:"id"`
	LessonId    int           `json:"lessonId"`
	Date        int           `json:"date"`
	DueDate     int           `json:"dueDate"`
	Text        string        `json:"text"`
	Remark      string        `json:"remark"`
	Completed   bool          `json:"completed"`
	Attachments []interface{} `json:"attachments"`
	LessonName  string
}

func (h Homework) GetDate() time.Time {
	d, err := ParseUntisDate(str(h.Date))
	if err != nil {
		log.Fatalln(err)
	}
	return d
}

func (h Homework) GetDueDate() time.Time {
	d, err := ParseUntisDate(str(h.DueDate))
	if err != nil {
		log.Fatalln(err)
	}
	return d
}

func (h Homework) String() string {
	return h.LessonName + ": " + h.Text +
		" (" + h.GetDate().Format("02.01.2006") +
		" - " + h.GetDueDate().Format("02.01.2006") + ")"
}
