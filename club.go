package club

import (
	"bytes"
	"sync"
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
