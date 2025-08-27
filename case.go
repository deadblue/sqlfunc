package sqlfunc

import (
	"strings"
	"unicode"
)

func pascalToSnake(name string) string {
	sb := strings.Builder{}
	for i, ch := range name {
		if unicode.IsUpper(ch) {
			if i != 0 {
				sb.WriteByte('_')
			}
			sb.WriteRune(unicode.ToLower(ch))
		} else {
			sb.WriteRune(ch)
		}
	}
	return sb.String()
}
