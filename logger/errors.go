package logger

import "fmt"

type RetroError struct {
	Err     string
	Details []error
}

func (e RetroError) Error() string {
	ers := e.Err
	for _, detail := range e.Details {
		ers += " " + detail.Error()
	}
	return ers
}

func GError(parts ...any) RetroError {
	return RetroError{
		Err:     fmt.Sprint(parts...), // Fix is here
		Details: nil,
	}
}
