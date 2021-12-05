package main

import (
	"fmt"
)

const (
	appCode = "GenTree"
)

/* The error message constants are used to avoid potential, subtle differences in feedback messages
   allowing potential sniffing of the server implementation details. They are returned as the
   'message' field of the JSON error response. */
const (
	internalErrorMsg = "Unexpected error occurred"
	payloadErrorMsg = "Payload validation error"
)

type AppError struct {
	Code int
	msg  string
}

func (e AppError) Error() string {
	return fmt.Sprintf("%s#%d: %s", appCode, e.Code, e.msg)
}
