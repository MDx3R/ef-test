package dto

import (
	"encoding/json"
	"strings"
	"time"
)

const monthYearLayout = "01-2006"

type MonthYear string

func (my MonthYear) ToTime() time.Time {
	t, err := my.Parse()
	if err != nil || t == nil {
		return time.Time{}
	}
	return *t
}

func (my MonthYear) ToTimePtr() *time.Time {
	t, _ := my.Parse()
	return t
}

func (my MonthYear) MarshalJSON() ([]byte, error) {
	if err := my.Validate(); err != nil {
		return nil, err
	}

	s := strings.Trim(string(my), `"`)
	if s == "" || s == "null" {
		return []byte(`"null"`), nil
	}

	return json.Marshal(s)
}

func (my *MonthYear) UnmarshalJSON(b []byte) error {
	if string(b) == "null" || len(b) == 0 {
		*my = MonthYear("null")
		return nil
	}

	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	return my.UnmarshalText([]byte(s))
}

func (my *MonthYear) UnmarshalText(b []byte) error {
	s := string(b)
	if s == "" || s == "null" {
		*my = MonthYear("null")
		return nil
	}

	*my = MonthYear(s)
	return my.Validate()
}

func (my MonthYear) Validate() error {
	_, err := my.Parse()
	return err
}

func (my *MonthYear) Parse() (*time.Time, error) {
	if my == nil {
		return nil, nil
	}
	s := strings.Trim(string(*my), `"`)
	if s == "" || s == "null" {
		return nil, nil
	}

	t, err := time.Parse(monthYearLayout, s)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func FromTime(t *time.Time) MonthYear {
	if t == nil || t.IsZero() {
		return "null"
	}
	return MonthYear(t.Format(monthYearLayout))
}
