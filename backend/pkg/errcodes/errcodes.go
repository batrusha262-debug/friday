package errcodes

import "git.appkode.ru/pub/go/failure"

var (
	ValidationError     = failure.NewErrorCode("VALIDATION_ERROR")
	NotFound            = failure.NewErrorCode("NOT_FOUND")
	Forbidden           = failure.NewErrorCode("FORBIDDEN")
	InternalServerError = failure.NewErrorCode("INTERNAL_SERVER_ERROR")
)
