package errgo

//-----------------------------------------------------------------------------

const markerFormat = "%s/%02d:%s()"

// Marker prepend err with markerFormat
type Marker struct {
	name string
	err  error
}

func (m *Marker) Error() string {
	var cause string
	if m.err != nil {
		cause = m.err.Error()
	} else {
		cause = "CAUSE:N/A"
	}
	return m.name + ": " + cause
}

// Cause implements Causer interface
func (m *Marker) Cause() error { return m.err }

// NewMarker whereever it be called, will prepend the other error with:
// markerFormat
func NewMarker(cause error) error {
	if cause == nil {
		return nil
	}
	name := makeMarker(markerFormat)
	return &Marker{
		name: name,
		err:  cause,
	}
}

//-----------------------------------------------------------------------------
