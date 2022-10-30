package collector

import (
	"fmt"
	"strings"
	"sync"
)

// NewJoiningCollector returns a collector that concatenates the input elements,
// separated by the specified delimiter, in encounter order.
//
// Note: This collector is not routine-safe.
func NewJoiningCollector[T any](separator string) Collector[T, *strings.Builder, string] {
	return NewBaseCollector[T, *strings.Builder, string](
		func() *strings.Builder {
			return &strings.Builder{}
		},
		func(container *strings.Builder, item T) {
			if container.Len() > 0 {
				container.WriteString(separator)
			}
			container.WriteString(fmt.Sprint(item))
		},
		func(container *strings.Builder) string {
			return container.String()
		},
	)
}

// NewJoiningCollectorInParallel returns a collector that concatenates the input elements,
// separated by the specified delimiter, in encounter order.
//
// Note: This collector is routine-safe.
func NewJoiningCollectorInParallel[T any](separator string) Collector[T, *strings.Builder, string] {
	mu := &sync.Mutex{}
	return NewBaseCollector[T, *strings.Builder, string](
		func() *strings.Builder {
			return &strings.Builder{}
		},
		func(container *strings.Builder, item T) {
			mu.Lock()
			defer mu.Unlock()

			if container.Len() > 0 {
				container.WriteString(separator)
			}
			container.WriteString(fmt.Sprint(item))
		},
		func(container *strings.Builder) string {
			return container.String()
		},
	)
}
