package types

import (
	"strconv"
	"time"
)

// necessary for unmarshalling time from json, that isn't in the default go time format

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

func (t Time) ToTime() time.Time {
	return time.Time(t)
}

func (t LessonTime) String() string {
	return t.toTime().Format("15:04")
}
