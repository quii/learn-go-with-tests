package mysolution

import "strings"

func Repeat(s string) string {
	var output strings.Builder

	for _ = range 5 {
		output.WriteString(s)
	}
	return output.String()
}
