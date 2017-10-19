package sentinel

//-----------------------------------------------------------------------------

// ErrString a string that inplements the error interface
type ErrString string

func (v ErrString) Error() string { return string(v) }

//-----------------------------------------------------------------------------

// ErrCode an int that inplements the error interface
type ErrCode int

func (v ErrCode) Error() string { return string(v) }

//-----------------------------------------------------------------------------
