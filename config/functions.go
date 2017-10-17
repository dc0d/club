package config

import (
	"os"
	"path/filepath"
)

//-----------------------------------------------------------------------------

// RelativeResource returns the file path 1 - if it's inside the directory of the program
// 2 - if it's inside the working directory 3 - returns os.ErrNotExist error
func RelativeResource(fileName string) (filePath string, funcErr error) {
	appDir, err := os.Executable()
	if err != nil {
		funcErr = err
		return
	}
	appDir, err = filepath.Abs(filepath.Dir(appDir))
	if err != nil {
		funcErr = err
		return
	}
	p := filepath.Join(appDir, fileName)
	if _, err := os.Stat(p); err == nil {
		filePath = p
		return
	}
	wd, err := os.Getwd()
	if err != nil {
		funcErr = err
		return
	}
	p = filepath.Join(wd, fileName)
	if _, err := os.Stat(p); err == nil {
		filePath = p
		return
	}
	funcErr = os.ErrNotExist
	return
}

//-----------------------------------------------------------------------------
