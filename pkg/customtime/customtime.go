package customtime

import (
	"strings"
	"time"
)

type Time struct {
	time.Time
}

func (ct *Time) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	ct.Time, _ = time.Parse("2006-01-02T15:04:05Z0700", s)

	return nil
}
