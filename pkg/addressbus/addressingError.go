package addressbus

import "fmt"

type addressingError struct {
	Address uint16
	Op      string
}

func (e *addressingError) Error() string {
	return e.Op + " " + fmt.Sprintf("%x", e.Address)
}
