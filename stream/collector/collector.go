package collector

// Identify returns the given item
func Identify[V any]() func(v V) V {
	return func(v V) V {
		return v
	}
}

// Collector is a mutable reduction operation that
// accumulates input elements into a mutable result container,
// optionally transforming the accumulated result into a final representation
// after all input elements have been processed.
// Reduction operations can be performed either sequentially or in parallel.
type Collector[T any, A any, R any] interface {
	// Supplier returns a function that returns a new, mutable result container.
	Supplier() func() A
	// Accumulator returns a function that folds a value into a mutable result container.
	Accumulator() func(container A, item T)
	// Finisher returns a function that transforms the final accumulator into the result of the reduction.
	Finisher() func(container A) R
}

var _ Collector[any, any, any] = (*BaseCollector[any, any, any])(nil)

type BaseCollector[T any, A any, R any] struct {
	supplier    func() A
	accumulator func(container A, item T)
	finisher    func(container A) R
}

func NewBaseCollector[T any, A any, R any](
	supplier func() A,
	accumulator func(container A, item T),
	finisher func(container A) R,
) *BaseCollector[T, A, R] {
	return &BaseCollector[T, A, R]{
		supplier:    supplier,
		accumulator: accumulator,
		finisher:    finisher,
	}
}

func (bc *BaseCollector[T, A, R]) Supplier() func() A {
	return bc.supplier
}

func (bc *BaseCollector[T, A, R]) Accumulator() func(container A, item T) {
	return bc.accumulator
}

func (bc *BaseCollector[T, A, R]) Finisher() func(container A) R {
	return bc.finisher
}
