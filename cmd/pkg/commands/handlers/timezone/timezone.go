package timezone

import (
	"errors"
)

func Get(args map[string]any) (map[string]any, error) {
	return map[string]any{
		"timezone": "Europe/Moscow",
	}, nil
}

func Set(args map[string]any) (map[string]any, error) {
	newTimezone, ok := args["timezone"].(string)
	if !ok {
		return nil, errors.New("invalid timezone arg")
	}
	return map[string]any{
		"timezone": newTimezone,
	}, nil
}
