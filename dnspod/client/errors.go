package client

import (
	"fmt"
)

type InvalidClientFieldError string

func (icfe InvalidClientFieldError) Error() string {
	return "Invalid client field: " + string(icfe)
}

type BadStatusCodeError struct {
	Code    string
	Message string
}

func (bsce *BadStatusCodeError) Error() string {
	return fmt.Sprintf("Bad StatusCode %s: %s", bsce.Code, bsce.Message)
}
