package collection

import (
	"github.com/carter-ya/go-tools/stream"
	"strings"
)

// String returns a string representation of the collection.
func String[E comparable](c Collection[E]) string {
	sb := strings.Builder{}
	sb.WriteString("[")
	itemsString := c.Stream().Collect(
		stream.JoiningSupplier[E](),
		stream.JoiningAccumulator[E](", "),
	).(string)
	sb.WriteString(itemsString)
	sb.WriteString("]")
	return sb.String()
}
