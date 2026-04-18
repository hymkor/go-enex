//go:build debug

package enex

import (
	"fmt"
	"strings"

	"github.com/nyaosorg/go-windows-dbg"
)

const debugBuild = true

func debug(v ...any) {
	dbg.Println(v...)
}

func stackTrace(e error, s ...string) error {
	if e == nil {
		return nil
	}
	return fmt.Errorf("%w\nat %s", e, strings.Join(s, " "))
}
