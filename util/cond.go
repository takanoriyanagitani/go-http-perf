package util

func Select[T any](f T, t T, cond bool) T {
	switch cond {
	case true:
		return t
	default:
		return f
	}
}
