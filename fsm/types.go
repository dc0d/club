package fsm

//-----------------------------------------------------------------------------

// State represents a state activity
type State func() (State, error)

// Activate activates the state and it's consecutive states until the next state
// is nil or encounters an error
func Activate(s State) (funcErr error) {
	next := s
	for next != nil && funcErr == nil {
		next, funcErr = next()
	}
	return
}

//-----------------------------------------------------------------------------
