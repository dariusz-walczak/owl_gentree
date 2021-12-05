package main

import (
	"fmt"
)

const (
	appCode = "GenTree"
)

/* The error message constants are used to avoid potential, subtle differences in feedback messages
   allowing potential sniffing of the server implementation details */
const (
	/* The internalErrorMsg text is returned as the 'message' field of the JSON error response to
	   indicate an internal server error. */
	internalErrorMsg = "Unexpected error occurred"
	/* The payloadErrorMsg text is returned as the 'message' field of the JSON error response to
	   indicate the JSON payload schema validation error. */
	payloadErrorMsg = "Payload validation error"
)

type AppError struct {
	Code int
	msg  string
}

func (e AppError) Error() string {
	return fmt.Sprintf("%s#%d: %s", appCode, e.Code, e.msg)
}
