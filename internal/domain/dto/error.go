package dto

import "fmt"

type Error struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("message: %s, details: %v", e.Message, e.Details)
}
