package basecrm

import (
	"fmt"
)

type ErrorLinks struct {
	// An optional link to resources that can be helpful in solving the issue.
	MoreInfo string `json:"more_info"`
}

type Error struct {
	// The resource the error relates to.
	Resource string `json:"resource"`
	// The field of the resource the error relates to.
	Field string `json:"field"`
	// The error code.
	Code string `json:"code"`
	// Human readable error description in a language specified by the Content-Language header.
	Message string `json:"message"`
	// An optional detailed descriptive text targeted at the client developer in English.
	Details string `json:"details"`
}

type ErrorMeta struct {
	Type string `json:"type"`
	// An optional link to resources that can be helpful in solving the issue.
	Links *ErrorLinks `json:"links"`
}

type ErrorEnvelope struct {
	Error *Error     `json:"error"`
	Meta  *ErrorMeta `json:"meta"`
}

type ErrorsMeta struct {
	// Errors envelope type. Must be set to **errors**
	Type string `json:"type"`
	// HTTP status code of the response plus HTTP response status message.
	HttpStatus string `json:"http_status"`
	// Unique id of the request. This is the same value as the X-Request-Id header.
	Logref string `json:"logref"`
	// An optional link to resources that can be helpful in solving the issue.
	Links *ErrorLinks `json:"links"`
}

type ErrorsEnvelope struct {
	Errors []*ErrorEnvelope `json:"errors"`
	Meta   *ErrorsMeta      `json:"meta"`
}

func (envelope *ErrorsEnvelope) String() string {
	err := envelope.Errors[0].Error
	return fmt.Sprintf("resource=%v, field=%v, code=%v, message=%v, details=%v, logref=%v",
		err.Resource,
		err.Field,
		err.Code,
		err.Message,
		err.Details,
		envelope.Meta.Logref)
}

func (envelope *ErrorsEnvelope) Error() string {
	return envelope.String()
}
