package line

import (
	"strings"
	"testing"
)

func TestIndexgoOfString(t *testing.T) {
	i := strings.Index("sss   ssss", "   ")
	if i == -1 {
		t.Fail()
	}
}
