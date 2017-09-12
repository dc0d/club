package club

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteSize(t *testing.T) {
	var n ByteSize = 2048
	assert.Equal(t, "2.00 KB", fmt.Sprint(n))
}
