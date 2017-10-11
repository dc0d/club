package hclconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dc0d/club/config"
	"github.com/hashicorp/hcl"
)

//-----------------------------------------------------------------------------

type loaderFunc func(interface{}, ...string) error

func (lf loaderFunc) Load(ptr interface{}, filePath ...string) error {
	return lf(ptr, filePath...)
}

//-----------------------------------------------------------------------------

// New returns a hcl config.Loader that loads hcl conf file. default conf file names
// (if filePath not provided) in the same directory are <appname>.conf and if
// not fount app.conf
func New() config.Loader {
	return loaderFunc(loadHCL)
}

//-----------------------------------------------------------------------------

func loadHCL(ptr interface{}, filePath ...string) error {
	var fp string
	if len(filePath) > 0 {
		fp = filePath[0]
	}
	if fp == "" {
		fp = _confFilePath()
	}
	cn, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}
	err = hcl.Unmarshal(cn, ptr)
	if err != nil {
		return err
	}

	return nil
}

func _confFilePath() string {
	appName := filepath.Base(os.Args[0])
	appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {

	}
	appConfName := fmt.Sprintf("%s.conf", appName)
	genericConfName := "app.conf"

	for _, vn := range []string{appConfName, genericConfName} {
		currentPath := filepath.Join(appDir, vn)
		if _, err := os.Stat(currentPath); err == nil {
			return currentPath
		}
	}

	for _, vn := range []string{appConfName, genericConfName} {
		wd, err := os.Getwd()
		if err != nil {
			continue
		}
		currentPath := filepath.Join(wd, vn)
		if _, err := os.Stat(currentPath); err == nil {
			return currentPath
		}
	}

	if _, err := os.Stat(appConfName); err == nil {
		return appConfName
	}

	return genericConfName
}

//-----------------------------------------------------------------------------
