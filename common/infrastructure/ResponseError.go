package infrastructure

import (
	"fmt"
	"errors"

	"UnpakSiamida/common/domain"
	"github.com/gofiber/fiber/v2"
)

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

func HandleError(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	// 1) coba unwrap ResponseError
	var respErr *ResponseError
	if errors.As(err, &respErr) {
		// map code -> http status jika perlu
		status := 400
		switch respErr.Code {
		case "AkurasiPenelitian.Validation":
			status = 400
		case "AkurasiPenelitian.NotFound":
			status = 404
		case "AkurasiPenelitian.Conflict":
			status = 409
		default:
			status = 400
		}
		return c.Status(status).JSON(respErr)
	}

	// 2) coba unwrap domain.Error
	var derr domain.Error
	if errors.As(err, &derr) {
		re := NewResponseError(derr.Code, derr.Description)
		status := 400
		switch derr.Type {
		case domain.NotFound:
			status = 404
		case domain.Conflict:
			status = 409
		case domain.Validation:
			status = 400
		default:
			status = 500
		}
		return c.Status(status).JSON(re)
	}

	// 3) fallback: unknown/wrapped error -> internal server error
	return c.Status(500).JSON(NewInternalError(err))
}