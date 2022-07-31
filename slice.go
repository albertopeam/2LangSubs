package main

func LastIndex[T any](a []T, b []T) T {
	if len(a) > len(b) {
		return a[len(a)-1]
	} else {
		return b[len(b)-1]
	}
}
