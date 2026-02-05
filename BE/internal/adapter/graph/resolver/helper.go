package resolver

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func safelyToFloat64(v interface{}) (float64, error) {
	switch val := v.(type) {
	case json.Number:
		return val.Float64()
	case float64:
		return val, nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	case float32:
		return float64(val), nil
	default:
		return 0, fmt.Errorf("unknown type")
	}
}
