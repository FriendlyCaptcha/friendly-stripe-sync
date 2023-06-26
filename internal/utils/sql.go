package utils

import (
	"database/sql"
	"encoding/json"

	"github.com/tabbed/pqtype"
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
