package vimnumbers

import "fmt"

func ConvertToHexColor(color int) *string {
	hexColor := fmt.Sprintf("#%06X", color)
	return &hexColor
}

func ForceInt(val interface{}) int {
	switch v := val.(type) {
	case uint64:
		return int(v)
	case int64:
		return int(v)
	default:
		panic(fmt.Sprintf("unexpected type: %T", v))
	}
}
