package utils

func Map[T, U any](s []T, f func(T) U) []U {
	result := make([]U, len(s))
	for i, item := range s {
		result[i] = f(item)
	}
	return result
}

func Filter[T any](s []T, f func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range s {
		if f(item) {
			result = append(result, item)
		}
	}
	return result
}
