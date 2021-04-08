package enex_test

import (
	"strings"
	"testing"

	"github.com/zetamatta/go-enex"
)

func TestShrinkMarkdown(t *testing.T) {
	source := `* A

* B

* C

* D

aaaa
bbbb
cccc
`

	expect := `* A
* B
* C
* D

aaaa
bbbb
cccc
`
	var buffer strings.Builder
	enex.ShrinkMarkdown(strings.NewReader(source), &buffer)
	actual := buffer.String()
	if actual != expect {
		t.Fatalf("expect `%s` but `%s`\n", expect, actual)
	}
}
