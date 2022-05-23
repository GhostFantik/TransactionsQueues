package pkg

import "fmt"

type Error struct {
	Code        int
	Description string
}

func NewError(code int, description string) Error {
	return Error{
		Code:        code,
		Description: description,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("status %d: err %v", e.Code, e.Description)
}
