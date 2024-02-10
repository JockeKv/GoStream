package stream

type Stream[T any] chan T

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

// Returns a Stream of type T with len buflen created from function f.
// The function get the channel as an argument and should put the items in it.
// The channel is closed when the function returns.
func Func[T any](f func(c chan T), buflen int) Stream[T] {
	c := make(chan T, buflen)
	go func(ch chan T) {
		f(ch)
		defer close(c)
	}((c))
	return c
}

// Returns a Stream of type T and length len(s) from Slice s
func Slice[S []T, T any](s S) Stream[T] {
	c := make(chan T, len(s))
	go func() {
		defer close(c)
		for _, i := range s {
			c <- i
		}
	}()
	return c
}

// Returns a Stream of type Pair[K, V] with len(s) from Map s
func Map[K comparable, V any](s map[K]V) Stream[Pair[K, V]] {
	c := make(chan Pair[K, V], len(s))
	go func() {
		defer close(c)
		for k, v := range s {
			c <- Pair[K, V]{Key: k, Value: v}
		}
	}()
	return c
}

// Returns a Stream of type T from the parameters s with len(s)
func Of[T any](s ...T) Stream[T] {
	c := make(chan T, len(s))
	go func() {
		defer close(c)
		for _, i := range s {
			c <- i
		}
	}()
	return c
}

func (t Stream[T]) Filter(pred func(e T) bool) Stream[T] {
	c := make(chan T, len(t))
	go func() {
		defer close(c)
		for i := range t {
			if pred(i) {
				c <- i
			}
		}
	}()
	return c
}

func (t Stream[T]) Map(pred func(e T) T) Stream[T] {
	c := make(chan T, len(t))
	go func() {
		defer close(c)
		for i := range t {
			c <- pred(i)
		}
	}()
	return c
}

func (t Stream[T]) Reduce(pred func(a T, b T) T) T {
	var result T
	first := true
	for i := range t {
		if first {
			result = i
			first = false
		} else {
			result = pred(result, i)
		}
	}
	return result
}

func (t Stream[T]) ForEach(pred func(e T)) {
	for i := range t {
		pred(i)
	}
}

func (t Stream[T]) Collect() []T {
	r := []T{}
	for i := range t {
		r = append(r, i)
	}
	return r
}
