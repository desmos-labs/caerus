package database

import (
	"database/sql"
	"time"
)

// StringToNullString allows to convert the given string to a sql.NullString
func StringToNullString(value string) sql.NullString {
	return sql.NullString{
		String: value,
		Valid:  value != "",
	}
}

// NullStringToString allows to convert the given sql.NullString to a string
func NullStringToString(value sql.NullString) string {
	if value.Valid {
		return value.String
	}
	return ""
}

// TimeToNullTime allows to convert the given time.Time to a sql.NullTime
func TimeToNullTime(value *time.Time) sql.NullTime {
	if value != nil {
		return sql.NullTime{
			Time:  *value,
			Valid: true,
		}
	}
	return sql.NullTime{
		Valid: false,
	}
}

// NullTimeToTime allows to convert the given sql.NullTime to a time.Time
func NullTimeToTime(value sql.NullTime) *time.Time {
	if value.Valid {
		return &value.Time
	}
	return nil
}
