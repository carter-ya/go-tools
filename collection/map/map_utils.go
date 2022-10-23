package _map

import (
	"fmt"
	"strings"
)

// MapString returns a string representation of the map.
func MapString[K comparable, V any](m Map[K, V]) string {
	sb := strings.Builder{}
	sb.WriteString("{")
	separator := ""
	m.ForEach(func(key K, value V) {
		sb.WriteString(separator)
		sb.WriteString(fmt.Sprint(key))
		sb.WriteString(": ")
		sb.WriteString(fmt.Sprint(value))
		separator = ", "
	})
	sb.WriteString("}")
	return sb.String()
}
