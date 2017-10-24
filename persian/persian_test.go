package persian

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToPersianText(t *testing.T) {
	assert := assert.New(t)

	txt := ToPersianText(0)
	assert.Equal("۰", txt)

	txt = ToPersianText(9)
	assert.Equal("۹", txt)

	txt = ToPersianText(1234567890)
	assert.Equal("۱۲۳۴۵۶۷۸۹۰", txt)

	txt = ToPersianText(48)
	assert.Equal("۴۸", txt)

	txt = ToPersianText(110)
	assert.Equal("۱۱۰", txt)
}
