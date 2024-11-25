package gtml

import "strings"

func For[T any](slice []T, callback func(T) string) string {
	var builder strings.Builder
	for _, item := range slice {
		builder.WriteString(callback(item))
	}
	return builder.String()
}

func If(condition bool, fn func() string) string {
	if condition {
		return fn()
	}
	return ""
}

func Else(condition bool, fn func() string) string {
	if !condition {
		return fn()
	}
	return ""
}

func Slot(contentFunc func() string) string {
	return contentFunc()
}
