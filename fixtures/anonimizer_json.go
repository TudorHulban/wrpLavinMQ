package fixtures

import (
	"encoding/json"
	"io"
	"strconv"
)

var FieldNameMap = make(map[string]string)
var fieldCounter = 1

func anonymizeValue(value any) any {
	switch v := value.(type) {
	case map[string]any:
		result := make(map[string]any)

		for key, val := range v {
			// Get or create consistent field name
			anonName, exists := FieldNameMap[key]
			if !exists {
				anonName = "Field" + strconv.Itoa(fieldCounter)
				FieldNameMap[key] = anonName
				fieldCounter++
			}

			result[anonName] = anonymizeValue(val)
		}

		return result

	case []any:
		result := make([]any, len(v))

		for i, val := range v {
			result[i] = anonymizeValue(val)
		}

		return result

	case string:
		return "some string"

	case int:
		return 7 // Integer value

	case float64:
		if v == float64(int(v)) {
			return 7
		}

		return 7.0 // Float value

	case bool:
		return true

	case nil:
		return nil

	default:
		return "unknown type"
	}
}

func AnonymizeJSON(r io.Reader, w io.Writer) error {
	var data any

	if errDecode := json.NewDecoder(r).Decode(&data); errDecode != nil {
		return errDecode
	}

	anonymized := anonymizeValue(data)
	encoder := json.NewEncoder(w)

	encoder.SetIndent("", "  ") // Optional: pretty print

	return encoder.Encode(anonymized)
}
