package tune

// import (
// 	"runtime"
// )

// // Preempt helper that calls runtime.Gosched(), for using in for loops.
// // By default returns true unless the result is explicitly set.
// func Preempt(res ...bool) bool {
// 	runtime.Gosched()
// 	if len(res) > 0 {
// 		return res[0]
// 	}
// 	return true
// }
