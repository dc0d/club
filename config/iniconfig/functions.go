package iniconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dc0d/club/config"
	"gopkg.in/ini.v1"
)

//-----------------------------------------------------------------------------

// New returns a ini config.Loader that loads ini conf file. default conf file names
// (if filePath not provided) in the same directory are <appname>.ini and if
// not fount app.ini
func New() config.Loader {
	return config.LoaderFunc(loadINI)
}

//-----------------------------------------------------------------------------

func loadINI(ptr interface{}, filePath ...string) error {
	fp, err := checkFilePath(filePath...)
	if err != nil {
		return err
	}
	cf, err := loadINIFile(fp)
	if err != nil {
		return err
	}
	return cf.StrictMapTo(ptr)
}

func loadINIFile(fp string) (cnf *ini.File, funcErr error) {
	cnf, funcErr = ini.Load(fp)
	return
}

func checkFilePath(filePath ...string) (fp string, funcErr error) {
	if len(filePath) > 0 {
		fp = filePath[0]
	}

	if fp == "" {
		var err error
		appName := filepath.Base(os.Args[0])
		appConfName := fmt.Sprintf("%s.ini", appName)
		if fp, err = config.RelativeResource(appConfName); err != nil {
			if fp, err = config.RelativeResource(genericConfName); err != nil {
				return "", err
			}
		}
	}
	if _, err := os.Stat(fp); err != nil {
		return "", err
	}
	return
}

//-----------------------------------------------------------------------------

const (
	genericConfName = "app.ini"
)

//-----------------------------------------------------------------------------
