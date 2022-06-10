package app

import (
	"fmt"
	"strings"
)

type ErrorBag struct {
	Errors []error
}

func NewErrorBag() *ErrorBag {
	return &ErrorBag{
		Errors: []error{},
	}
}

func (bag *ErrorBag) hasError() bool {
	return len(bag.Errors) > 0
}

func (bag *ErrorBag) Error(msg string) error {
	if !bag.hasError() {
		return nil
	}

	errorMessages := make([]string, len(bag.Errors))
	for i, err := range bag.Errors {
		errorMessages[i] = fmt.Sprintf("%d: %s", i, err)
	}

	return fmt.Errorf("%s:\n\n%s", msg, strings.Join(errorMessages, "\n"))
}

func (bag *ErrorBag) Add(err error) {
	bag.Errors = append(bag.Errors, err)
}
