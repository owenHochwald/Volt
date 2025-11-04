package utils

import "strings"

// TODO: add allow keyword support for headers

func ParseKeyValuePairs(input string) (map[string]string, []string) {
	result := make(map[string]string)
	var errors []string

	lines := strings.Split(input, ",")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}
		key, value, found := strings.Cut(trimmed, "=")
		if !found {
			errors = append(errors, "Invalid key-value pair: "+trimmed)
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if len(key) == 0 {
			errors = append(errors, "Invalid key: "+key)
			continue
		}
		if len(value) == 0 {
			errors = append(errors, "Invalid value: "+value)
			continue
		}
		result[key] = value

	}

	return result, errors
}
