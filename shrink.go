package enex

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

const (
	_blank = iota
	_item
	_asis
)

var (
	rxBlank = regexp.MustCompile(`^\s*$`)
	rxUl    = regexp.MustCompile(`^\s*\*\s`)
	rxOl    = regexp.MustCompile(`^\s*\d+\.`)
	rxPre   = regexp.MustCompile("^```")
)

func ShrinkMarkdown(r io.Reader, w io.Writer) {
	sc := bufio.NewScanner(r)
	pre := false
	lastLine := _asis
	pendingBlank := false
	for sc.Scan() {
		line := sc.Text()
		if rxPre.MatchString(line) {
			pre = !pre
			lastLine = _asis
			if pendingBlank {
				fmt.Fprintln(w)
				pendingBlank = false
			}
		} else if pre {
			lastLine = _asis
		} else if rxBlank.MatchString(line) {
			switch lastLine {
			case _blank:
				continue
			case _item:
				pendingBlank = true
				continue
			}
			lastLine = _blank
		} else if rxUl.MatchString(line) || rxOl.MatchString(line) {
			lastLine = _item
			pendingBlank = false
		} else {
			if pendingBlank {
				fmt.Fprintln(w)
				pendingBlank = false
			}
			lastLine = _asis
		}
		fmt.Fprintln(w, line)
	}
}
