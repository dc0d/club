package config

// Loader loads the config into a struct which is passed as a pointer
type Loader interface {
	Load(ptr interface{}, filePath ...string) error
}
