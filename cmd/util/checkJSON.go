package util

import (
	"encoding/json"
	"errors"
	"fmt"
)

func CheckJSON(decoder *json.Decoder) (json.Delim, error) {

	// check if the first token is a JSON object or array
	token, err := decoder.Token()
	if err != nil {
		return 0, err
	}

	delim, ok := token.(json.Delim)
	if !ok {
		return 0, errors.New("top-level JSON is not an object or array")
	}

	if delim != '{' && delim != '[' {
		return 0, fmt.Errorf("unsupported JSON delimiter: %v", delim)
	}

	return delim, nil
}
