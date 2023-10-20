package types

import "strconv"

type Holiday struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	LongName  string `json:"longName"`
	StartDate int    `json:"startDate"`
	EndDate   int    `json:"endDate"`
}

func (h Holiday) GetStartDate() Time {
	t, err := ParseUntisDate(strconv.Itoa(h.StartDate))
	if err != nil {
		return Time{}
	}
	return Time(t)
}

func (h Holiday) GetEndDate() Time {
	t, err := ParseUntisDate(strconv.Itoa(h.EndDate))
	if err != nil {
		return Time{}
	}
	return Time(t)
}
