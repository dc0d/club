package club

import "fmt"

//-----------------------------------------------------------------------------

// ByteSize represents size in bytes
type ByteSize float64

// ByteSize values
const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB

	// TB
	// PB
	// EB
	// ZB
	// YB
)

func (bs ByteSize) String() string {
	n := bs
	switch {
	case n > GB:
		return fmt.Sprintf("%.2f GB", n/GB)
	case n > MB:
		return fmt.Sprintf("%.2f MB", n/MB)
	case n > KB:
		return fmt.Sprintf("%.2f KB", n/KB)
	default:
		return fmt.Sprintf("%.2f B", n)
	}
}

//-----------------------------------------------------------------------------
