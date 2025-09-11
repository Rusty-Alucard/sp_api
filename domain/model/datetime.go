// Package model provides domain models and value objects for the application.
package model

import (
	"fmt"
	"io"
	"time"
)

// DateTime represents a custom date time type
type DateTime time.Time

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (d *DateTime) UnmarshalGQL(v interface{}) error {
	switch v := v.(type) {
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return err
		}
		*d = DateTime(t)
		return nil
	case time.Time:
		*d = DateTime(v)
		return nil
	default:
		return fmt.Errorf("DateTime is invalid")
	}
}

// MarshalGQL implements the graphql.Marshaler interface
func (d DateTime) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(`"` + time.Time(d).Format(time.RFC3339) + `"`))
}
