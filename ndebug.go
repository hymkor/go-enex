//go:build !debug

package enex

const debugBuild = false

func debug(v ...any) {}

func stackTrace(e error, _ ...string) error {
	return e
}
