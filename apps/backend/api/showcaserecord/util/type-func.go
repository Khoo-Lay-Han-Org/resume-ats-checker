package showcaserecord_util

import (
	"encoding/json"

	"gorm.io/datatypes"
)

func ToStringSlice(v any) []string {
	if v == nil {
		return []string{}
	}
	s, ok := v.([]any)
	if !ok {
		return []string{}
	}
	result := make([]string, len(s))
	for i, item := range s {
		result[i], _ = item.(string)
	}
	return result
}

func ToJSON(v any) datatypes.JSON {
	if v == nil {
		return datatypes.JSON("null")
	}
	b, err := json.Marshal(v)
	if err != nil {
		return datatypes.JSON("null")
	}
	return datatypes.JSON(b)
}
