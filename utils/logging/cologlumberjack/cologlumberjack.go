package cologlumberjack

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/comail/colog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Hook struct {
	levels []colog.Level
	logger *lumberjack.Logger
}

func New(logger *lumberjack.Logger, levels ...colog.Level) *Hook {
	if len(levels) == 0 {
		levels = []colog.Level{
			colog.LTrace,
			colog.LDebug,
			colog.LInfo,
			colog.LWarning,
			colog.LError,
			colog.LAlert,
		}
	}

	return &Hook{
		levels: levels,
		logger: logger,
	}
}

func (h *Hook) Levels() []colog.Level {
	return h.levels
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func getBuffer() *bytes.Buffer {
	buff := bufferPool.Get().(*bytes.Buffer)
	buff.Reset()
	return buff
}

func putBuffer(buff *bytes.Buffer) {
	bufferPool.Put(buff)
}

/*
	Time    time.Time // time of the event
*/

func (h *Hook) Fire(e *colog.Entry) error {
	buf := getBuffer()
	defer putBuffer(buf)

	buf.Write([]byte("level="))
	buf.Write([]byte(e.Level.String()))
	buf.WriteByte(' ')

	buf.Write([]byte("time="))
	buf.WriteString(e.Time.Format(time.RFC3339))
	buf.WriteByte(' ')

	if e.Line >= 0 {
		buf.Write([]byte("line="))
		buf.WriteString(strconv.Itoa(e.Line))
		buf.WriteByte(' ')
	}

	if e.File != "" {
		buf.Write([]byte("file="))
		buf.WriteString(fmt.Sprintf("%q", filepath.Base(e.File)))
		buf.WriteByte(' ')
	}

	for k, v := range e.Fields {
		buf.Write([]byte(k))
		buf.Write([]byte("="))
		buf.WriteString(fmt.Sprintf("%+q", v))
		buf.WriteByte(' ')
	}

	buf.Write([]byte("message="))
	buf.WriteString(fmt.Sprintf("%+q", e.Message))

	buf.WriteByte('\n')
	_, err := h.logger.Write(buf.Bytes())
	return err
}
