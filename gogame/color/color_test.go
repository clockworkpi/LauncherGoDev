package color

// > go test color.go color_test.go  -v
// Or
// > go test -v
// to test all test files

import "testing"

func TestColor(t *testing.T) {
	c := &Color{244,124,244,0}
	t.Logf("%x", c.ToHex())
}

