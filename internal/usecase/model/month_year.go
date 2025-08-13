package model

import (
	"fmt"
	"strings"
	"time"
)

type MonthYear struct {
	time.Time
}

func NewMonthYear(t time.Time) MonthYear {
	return MonthYear{Time: t}
}

func NewMonthYearFromPtr(t *time.Time) *MonthYear {
	if t == nil {
		return nil
	}
	return &MonthYear{Time: *t}
}

const monthYearLayout = "01-2006"

func (m *MonthYear) ToTime() time.Time {
	if m == nil {
		return time.Time{}
	}
	return m.Time
}

func (m *MonthYear) ToTimePtr() *time.Time {
	if m == nil {
		return nil
	}
	return &m.Time
}

func (my MonthYear) MarshalJSON() ([]byte, error) {
	if my.Time.IsZero() {
		return []byte(`null`), nil
	}
	return fmt.Appendf([]byte{}, `"%s"`, my.Time.Format(monthYearLayout)), nil
}

func (my *MonthYear) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), `"`)
	return my.UnmarshalText([]byte(str))
}

func (my *MonthYear) UnmarshalText(b []byte) error {
	s := string(b)
	if s == `` || s == `null` || s == `""` {
		my.Time = time.Time{}
		return nil
	}

	t, err := time.Parse(monthYearLayout, s)
	if err != nil {
		return err
	}

	my.Time = t
	return nil
}
