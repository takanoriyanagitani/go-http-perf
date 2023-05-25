package util

func CurryErr[T, U, V any](f func(T, U) (V, error)) func(T) func(U) (V, error) {
	return func(t T) func(U) (V, error) {
		return func(u U) (V, error) {
			return f(t, u)
		}
	}
}

func CurryErrIII[T, U, V, W any](f func(T, U, V) (W, error)) func(T) func(U) func(V) (W, error) {
	return func(t T) func(U) func(V) (W, error) {
		return func(u U) func(V) (W, error) {
			return func(v V) (W, error) {
				return f(t, u, v)
			}
		}
	}
}
