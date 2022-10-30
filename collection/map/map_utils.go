package _map

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
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

func MarshalJSON[K comparable, V any](m Map[K, V]) (bz []byte, err error) {
	if m == nil {
		return []byte("null"), nil
	}
	buf := bytes.NewBuffer(bz)
	buf.WriteByte('{')
	m.ForEachIndexed(func(index int, k K, v V) (stop bool) {
		if index > 0 {
			buf.WriteByte(',')
		}

		// write key
		isString := reflect.ValueOf(k).Kind() == reflect.String
		if !isString {
			buf.WriteByte('"')
		}
		var keyBz []byte
		keyBz, err = json.Marshal(k)
		if err != nil {
			return true
		}
		buf.Write(keyBz)
		if !isString {
			buf.WriteByte('"')
		}

		// write separator
		buf.WriteByte(':')

		// write value
		var valueBz []byte
		valueBz, err = json.Marshal(v)
		if err != nil {
			return true
		}
		buf.Write(valueBz)
		return false
	})
	if err != nil {
		return nil, err
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func UnmarshalJSON[K comparable, V any](m Map[K, V], data []byte) error {
	items := make(map[K]V, 0)
	err := json.Unmarshal(data, &items)
	if err != nil {
		return err
	}
	index := make(map[K]int, len(items))
	keys := make([]K, 0, len(items))
	for key := range items {
		keys = append(keys, key)
		keyBz, err := json.Marshal(key)
		if err != nil {
			return err
		}
		index[key] = bytes.Index(data, keyBz)
		if index[key] == -1 {
			return fmt.Errorf("key %v not found in data", key)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		return index[keys[i]] < index[keys[j]]
	})
	for _, key := range keys {
		m.Put(key, items[key])
	}
	return nil
}
