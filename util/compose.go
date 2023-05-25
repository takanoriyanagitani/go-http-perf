package util

func ComposeErr[T, U, V any](f func(T) (U, error), g func(U) (V, error)) func(T) (V, error) {
	return func(t T) (v V, e error) {
		u, e := f(t)
		return Select(
			func() (V, error) { return v, e },
			func() (V, error) { return g(u) },
			nil == e,
		)()
	}
}
