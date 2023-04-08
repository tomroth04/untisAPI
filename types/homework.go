package types

import (
	"github.com/rs/zerolog/log"
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
		log.Error().Err(err).Timestamp().Caller(0).
			Msg("error parsing date")
	}
	return d
}

func (h Homework) GetDueDate() time.Time {
	d, err := ParseUntisDate(str(h.DueDate))
	if err != nil {
		log.Error().Err(err).Caller(0).Timestamp().
			Msg("error parsing due date")
	}
	return d
}

func (h Homework) String() string {
	return h.LessonName + ": " + h.Text +
		" (" + h.GetDate().Format("02.01.2006") +
		" - " + h.GetDueDate().Format("02.01.2006") + ")"
}
