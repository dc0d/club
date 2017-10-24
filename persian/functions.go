package persian

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

//-----------------------------------------------------------------------------

// PolishYeKaf replaces characters by their persian equivalent (for keyboards or OSs that suck)
func PolishYeKaf(s string) (res string) {
	res = strings.Replace(
		s,
		"ي",
		"ی",
		-1)
	res = strings.Replace(
		res,
		"ك",
		"ک",
		-1)

	return
}

//-----------------------------------------------------------------------------

// IranTime .
func IranTime(source time.Time) time.Time {
	var dest time.Time
	loc, err := time.LoadLocation("Asia/Tehran")
	if err == nil {
		dest = source.In(loc)
	} else {
		dest = source
	}
	return dest
}

// IranNow .
func IranNow() time.Time {
	return IranTime(time.Now())
}

//-----------------------------------------------------------------------------

const (
	persianNumbers = "۰۱۲۳۴۵۶۷۸۹"
	latinNumbers   = "0123456789"
)

var p2l = make(map[rune]rune)
var l2p = make(map[rune]rune)

func init() {
	var p = []rune(persianNumbers)
	var l = []rune(latinNumbers)
	if len(l) != len(p) {
		panic(errors.Errorf("BOTH []rune MUST BE OF EQUAL LENGTH"))
	}
	for i := 0; i < 10; i++ {
		p2l[p[i]] = l[i]
		l2p[l[i]] = p[i]
	}
}

// ToPersianText .
func ToPersianText(n int) string {
	q := []rune(fmt.Sprintf("%d", n))
	for k, v := range q {
		q[k] = l2p[v]
	}
	return string(q)
}

//-----------------------------------------------------------------------------
