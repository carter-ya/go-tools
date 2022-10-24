package collection

import (
	"bytes"
	"encoding/json"
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

func MarshalJSON[E comparable](c Collection[E]) (bz []byte, err error) {
	if c == nil {
		return []byte("null"), nil
	}
	buf := bytes.NewBuffer(bz)
	buf.WriteByte('[')
	var separator []byte
	c.ForEachIndexed(func(_ int, e E) (stop bool) {
		buf.Write(separator)
		var elementBz []byte
		elementBz, err = json.Marshal(e)
		if err != nil {
			return true
		}
		buf.Write(elementBz)
		separator = []byte(",")
		return false
	})

	if err != nil {
		return nil, err
	}

	buf.WriteByte(']')
	return buf.Bytes(), nil
}

func UnmarshalJSON[E comparable](c Collection[E], data []byte) error {
	items := make([]E, 0)
	err := json.Unmarshal(data, &items)
	if err != nil {
		return err
	}
	for _, item := range items {
		c.Add(item)
	}
	return nil
}
