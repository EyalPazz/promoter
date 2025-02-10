package utils

func PtrTo[T any](a T) *T {
	return &a
}
