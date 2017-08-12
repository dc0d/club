package tune

// import (
// 	"bufio"
// 	"bytes"
// 	"log"
// 	"sync"

// 	"github.com/comail/colog"
// 	"github.com/fatih/color"
// )

// var bufferPool = sync.Pool{
// 	New: func() interface{} {
// 		return &bytes.Buffer{}
// 	},
// }

// func getBuffer() *bytes.Buffer {
// 	buff := bufferPool.Get().(*bytes.Buffer)
// 	buff.Reset()
// 	return buff
// }

// func putBuffer(buff *bytes.Buffer) {
// 	bufferPool.Put(buff)
// }

// // DevFormatter .
// type DevFormatter struct {
// 	flags int
// }

// // NewDevFormatter .
// func NewDevFormatter() (res *DevFormatter) {
// 	res = new(DevFormatter)
// 	return
// }

// // Flags .
// func (df *DevFormatter) Flags() int {
// 	return df.flags
// }

// type clset struct {
// 	headChar, level, message, fields, atTime, line, file []color.Attribute
// }

// // Format .
// func (df *DevFormatter) Format(e *colog.Entry) (res []byte, ferr error) {
// 	if e == nil {
// 		return nil, nil
// 	}

// 	buf := getBuffer()
// 	defer putBuffer(buf)
// 	defer func() { res = buf.Bytes() }()

// 	w := bufio.NewWriter(buf)
// 	defer w.Flush()

// 	var head clset
// 	head.headChar = []color.Attribute{color.BgBlack}
// 	head.level = []color.Attribute{color.FgWhite}
// 	head.message = []color.Attribute{color.FgWhite}
// 	head.atTime = []color.Attribute{color.Italic, color.FgHiMagenta}
// 	head.line = []color.Attribute{color.Italic, color.FgHiMagenta}
// 	head.file = []color.Attribute{color.Italic, color.FgHiMagenta}

// 	const headCharChar = "  "

// 	switch e.Level {
// 	case colog.LTrace:
// 		head.headChar = []color.Attribute{color.BgHiGreen}
// 		head.level = []color.Attribute{color.FgHiGreen, color.Underline}
// 	case colog.LDebug:
// 		head.headChar = []color.Attribute{color.BgBlack}
// 		head.level = []color.Attribute{color.FgWhite, color.Underline}
// 	case colog.LInfo:
// 		head.headChar = []color.Attribute{color.BgBlue}
// 		head.level = []color.Attribute{color.FgBlue, color.Underline}
// 	case colog.LWarning:
// 		head.headChar = []color.Attribute{color.BgYellow}
// 		head.level = []color.Attribute{color.FgHiYellow, color.Underline}
// 	case colog.LError:
// 		head.headChar = []color.Attribute{color.BgHiRed}
// 		head.level = []color.Attribute{color.FgHiRed, color.Underline}
// 	case colog.LAlert:
// 		head.headChar = []color.Attribute{color.BgRed}
// 		head.level = []color.Attribute{color.FgRed, color.Underline}
// 	}

// 	color.New(head.headChar...).Fprintf(w, headCharChar)
// 	color.New(head.level...).Fprintf(w, "+%-80v\n", e.Level)

// 	if len(e.Message) > 0 {
// 		msg := bytes.TrimSpace(e.Message)
// 		color.New(head.headChar...).Fprintf(w, headCharChar)
// 		color.New(head.message...).Fprintf(w, "%s\n", msg)
// 	}

// 	for k, v := range e.Fields {
// 		color.New(head.headChar...).Fprintf(w, headCharChar)
// 		color.New(head.fields...).Fprintf(w, "%s=%s\n", k, v)
// 	}

// 	if df.flags&(log.Ldate|log.Ltime|log.Lmicroseconds) != 0 {
// 		color.New(head.headChar...).Fprintf(w, headCharChar)
// 		dtFormat := ""
// 		dt := e.Time
// 		if df.flags&log.LUTC != 0 {
// 			dt = dt.UTC()
// 		}
// 		if df.flags&log.Ldate != 0 {
// 			dtFormat += "2006/01/02 "
// 		}
// 		if df.flags&log.Ltime != 0 {
// 			dtFormat += "15:04:05"
// 		} else if df.flags&log.Lmicroseconds != 0 {
// 			dtFormat += "15:04:05.000000"
// 		}
// 		color.New(head.atTime...).Fprintf(w, "@ %v\n", e.Time.Format(dtFormat))
// 	}

// 	if df.flags&(log.Lshortfile|log.Llongfile) != 0 {
// 		fp := e.File
// 		if df.flags&log.Lshortfile != 0 {
// 			short := fp
// 			for i := len(fp) - 1; i > 0; i-- {
// 				if fp[i] == '/' {
// 					short = fp[i+1:]
// 					break
// 				}
// 			}
// 			fp = short
// 		}
// 		color.New(head.headChar...).Fprintf(w, headCharChar)
// 		color.New(head.line...).Fprintf(w, "@ %v %v\n", e.Line, fp)
// 	}

// 	return
// }

// // SetFlags .
// func (df *DevFormatter) SetFlags(flags int) { df.flags = flags }
