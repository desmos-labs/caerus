package database

import (
	"database/sql"
	"time"
)

func Uint64toNullInt(value uint64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: int64(value),
		Valid: true,
	}
}

func NullIntToUint64(value sql.NullInt64) uint64 {
	if value.Valid {
		return uint64(value.Int64)
	}
	return 0
}

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
