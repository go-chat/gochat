package apperror

type AppError struct {
	StatusCode int
	Error      error
	Message    string
}

func NewAppError(err error, message string, statusCode int) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Error:      err,
		Message:    message,
	}
}
