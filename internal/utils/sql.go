package utils

import (
	"database/sql"
	"encoding/json"

	"github.com/sqlc-dev/pqtype"
)

func MarshalToNullRawMessage(value interface{}) pqtype.NullRawMessage {
	if value == nil {
		return pqtype.NullRawMessage{}
	}
	data, err := json.Marshal(value)
	return pqtype.NullRawMessage{RawMessage: data, Valid: err == nil}
}

func StringToNullString(value string) sql.NullString {
	return sql.NullString{String: value, Valid: value != ""}
}

func Int64ToNullInt64(value int64) sql.NullInt64 {
	return sql.NullInt64{Int64: value, Valid: value != 0}
}
