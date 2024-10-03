package dto

import (
	"fmt"
)

type Error struct {
	Message  string `json:"message"`
	Details  any    `json:"details,omitempty"`
	Code     int    `json:"-"`
	DebugMsg any    `json:"-"`
	Location string `json:"-"`
}

func NewError(code int, message string, soruce string, details any, debugmsg any) *Error {
	return &Error{Message: message, Code: code, Details: details, DebugMsg: debugmsg}
}

func (e *Error) Error() string {
	var msg string

	msg += fmt.Sprintf("code=%d message=%s ", e.Code, e.Message)
	if e.Location != "" {
		msg += fmt.Sprintf("source=%v ", e.Details)
	}
	if e.Details != nil {
		msg += fmt.Sprintf("details=%v ", e.Details)
	}

	if e.DebugMsg != nil {
		msg += fmt.Sprintf("debug=%v", e.DebugMsg)
	}

	return msg
}
