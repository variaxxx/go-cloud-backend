package core_http_request

import (
	"encoding/json"

	"github.com/google/uuid"
)

type NullableUUIDField struct {
	Present bool
	Value   *uuid.UUID
}

func (f *NullableUUIDField) UnmarshalJSON(data []byte) error {
	f.Present = true

	if string(data) == "null" {
		f.Value = nil
		return nil
	}

	var value uuid.UUID
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	f.Value = &value
	return nil
}
