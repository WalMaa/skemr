package ai

import "encoding/json"

// toolJSON marshals the given value into a JSON string. It returns an error if the marshaling fails.
func toolJSON(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// toolError creates a JSON string representing an error with the given code and message.
func toolError(code string, message string) string {
	b, _ := json.Marshal(map[string]any{
		"ok": false,
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})

	return string(b)
}
