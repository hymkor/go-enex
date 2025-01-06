package enex

import (
	"strings"
	"testing"
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
	shrinkMarkdown(strings.NewReader(source), &buffer)
	actual := buffer.String()
	if actual != expect {
		t.Fatalf("expect `%s` but `%s`\n", expect, actual)
	}
}
