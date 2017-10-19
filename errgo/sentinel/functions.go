package sentinel

import "fmt"

//-----------------------------------------------------------------------------

// Errorf value type (string) error
func Errorf(format string, a ...interface{}) error {
	return ErrString(fmt.Sprintf(format, a...))
}

//-----------------------------------------------------------------------------
