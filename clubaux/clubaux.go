package clubaux

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/dc0d/club/errors"
	"github.com/hashicorp/hcl"
)

//-----------------------------------------------------------------------------

// OnSignal runs function on receiving the OS signal
func OnSignal(f func(), sig ...os.Signal) {
	if f == nil {
		return
	}
	sigc := make(chan os.Signal, 1)
	if len(sig) > 0 {
		signal.Notify(sigc, sig...)
	} else {
		signal.Notify(sigc,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
			syscall.SIGSTOP,
			syscall.SIGABRT,
			syscall.SIGTSTP,
			syscall.SIGKILL)
	}
	go func() {
		<-sigc
		f()
	}()
}

var (
	errNotAvailable = errors.Errorf("N/A")
)

//-----------------------------------------------------------------------------

var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

// GetBuffer .
func GetBuffer() *bytes.Buffer {
	buff := bufferPool.Get().(*bytes.Buffer)
	return buff
}

// PutBuffer .
func PutBuffer(buff *bytes.Buffer) {
	buff.Reset() // ooch!
	bufferPool.Put(buff)
}

//-----------------------------------------------------------------------------

// LoadHCL loads hcl conf file. default conf file names (if filePath not provided)
// in the same directory are <appname>.conf and if not fount app.conf
func LoadHCL(ptr interface{}, filePath ...string) error {
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

// Chain .
func Chain(steps ...func() error) (chainerr error) {
	for _, v := range steps {
		chainerr = v()
		if chainerr != nil {
			return
		}
	}
	return
}

//-----------------------------------------------------------------------------
