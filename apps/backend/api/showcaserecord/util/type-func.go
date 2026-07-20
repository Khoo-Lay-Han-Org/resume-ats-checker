package showcaserecord_util

import "encoding/json"

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

func ToJSON(v any) any {
	if v == nil {
		return "null"
	}
	b, err := json.Marshal(v)
	if err != nil {
		return "null"
	}
	return string(b)
}
