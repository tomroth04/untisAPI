package types

import (
	"log/slog"
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

// GetDate returns the inscription date of the homework as a time.Time object
func (h Homework) GetDate() time.Time {
	d, err := ParseUntisDate(str(h.Date))
	if err != nil {
		slog.Error("Error parsing date", "error", err)
	}
	return d
}

// GetDueDate returns the due date of the homework as a time.Time object
func (h Homework) GetDueDate() time.Time {
	d, err := ParseUntisDate(str(h.DueDate))
	if err != nil {
		slog.Error("Error parsing due date", "error", err)
	}
	return d
}
