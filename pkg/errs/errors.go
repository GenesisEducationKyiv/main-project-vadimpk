package errs

import "strings"

// Err implements the Error interface with error marshaling.
type Err struct {
	Message string            `json:"message"`
	Code    int               `json:"code"`
	Details map[string]string `json:"details"`
}

func New(message string, code int) *Err {
	return &Err{Message: message, Code: code}
}

func (e *Err) Error() string {
	return e.Message
}

// IsExpected finds Err{} inside passed error.
func IsExpected(err error) bool {
	_, ok := err.(*Err)
	return ok
}

// GetDetails returns details of given error or nil when error is not custom
func GetDetails(err error) map[string]string {
	v, ok := err.(*Err)
	if !ok {
		return nil
	}
	return v.Details
}

// GetCode returns code of given error or empty string if error is not custom
func GetCode(err error) int {
	v, ok := err.(*Err)
	if !ok {
		return 500
	}
	return v.Code
}

// HasAnyGivenMessage returns true if given error has any given message.
// You can pass infinite number of messages into this function, but at least one is always required
func HasAnyGivenMessage(err error, firstMsg string, otherMsgs ...string) bool {
	if strings.Contains(err.Error(), firstMsg) {
		return true
	}
	for _, msg := range otherMsgs {
		if strings.Contains(err.Error(), msg) {
			return true
		}
	}
	return false
}

// DeleteFromDetails removes a value from details map by given key.
// Returns nothing on success and failure.
func DeleteFromDetails(err error, key string) {
	v, ok := err.(*Err)
	if !ok {
		return
	}
	delete(v.Details, key)
}
