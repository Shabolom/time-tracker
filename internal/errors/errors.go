package errors

// Err Used for custom errors
type Err struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (err Err) Error() string {
	return err.Message
}

const ErrResourceUnavailable = "This resource is unavailable"
