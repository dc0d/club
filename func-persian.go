package club

import (
	"strings"
	"time"
)

//-----------------------------------------------------------------------------

// PolishYeKaf replaces characters by their persian equivalent (for keyboards or OSs that suck)
func PolishYeKaf(s string) (res string) {
	res = s
	res = strings.Replace(
		res,
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

// IranTime .
func IranTime(source time.Time) time.Time {
	var dest time.Time
	loc, err := time.LoadLocation(`Asia/Tehran`)
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
