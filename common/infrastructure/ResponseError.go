package infrastructure

import "fmt"

type ResponseError struct {
    Code    string      `json:"code"`
    Message interface{} `json:"message"` // <-- bisa string atau object (map)
}

// implement error interface
func (e *ResponseError) Error() string {
    // jika message berupa string, kembalikan langsung; jika object, kita format ringkas
    switch m := e.Message.(type) {
    case string:
        return m
    default:
        return fmt.Sprintf("%s: %v", e.Code, m)
    }
}

func NewResponseError(code string, message interface{}) *ResponseError {
    return &ResponseError{
        Code:    code,
        Message: message,
    }
}

func NewInternalError(err error) *ResponseError {
    return &ResponseError{
        Code:    "INTERNAL_SERVER_ERROR",
        Message: err.Error(),
    }
}
