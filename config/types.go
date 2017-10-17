package config

//-----------------------------------------------------------------------------

// Loader loads the config into a struct which is passed as a pointer
type Loader interface {
	Load(ptr interface{}, filePath ...string) error
}

//-----------------------------------------------------------------------------

// LoaderFunc is a function type that implements Loader
type LoaderFunc func(interface{}, ...string) error

// Load implements Loader
func (lf LoaderFunc) Load(ptr interface{}, filePath ...string) error {
	return lf(ptr, filePath...)
}

//-----------------------------------------------------------------------------
