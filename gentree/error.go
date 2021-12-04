package main


import (
	"fmt"
)

const (
	appCode = "GenTree"
	/* The internalErrorMsg is returned as the 'message' field of the JSON response sent in the
	case of an internal server error. It is defined as a constant to avoid potential, subtle
	differences in the message allowing sniffing of the server implementation details */
	internalErrorMsg = "Unexpected error occurred"
)

type AppError struct {
	Code int
	msg string
}

func (e AppError) Error() string {
	return fmt.Sprintf("%s#%d: %s", appCode, e.Code, e.msg)
}
