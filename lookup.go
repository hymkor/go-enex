package enex

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
)

func dumpAttr(attr map[string]string, w io.Writer) {
	for key, val := range attr {
		fmt.Fprintf(w, "%q=%q\n", key, val)
	}
}

var (
	rxStyleWidth  = regexp.MustCompile(`--en-naturalWidth:(\d+)`)
	rxStyleHeight = regexp.MustCompile(`--en-naturalHeight:(\d+)`)
)

func (note *Note) LookupWithStyle(style string, log io.Writer) *Resource {
	w, h := -1, -1
	if s := rxStyleWidth.FindStringSubmatch(style); s != nil {
		if v, err := strconv.Atoi(s[1]); err == nil {
			w = v
		}
	}
	if s := rxStyleHeight.FindStringSubmatch(style); s != nil {
		if v, err := strconv.Atoi(s[1]); err == nil {
			h = v
		}
	}
	var lastFound *Resource
	defer func() {
		if lastFound != nil {
			fmt.Fprintf(log, "Falling back to dimension-based matching (width=%d, height=%d).\n", w, h)
		}
	}()
	for _, resouces := range note.Resource {
		for _, r := range resouces {
			if r.Width == w && r.Height == h {
				if !r.touch {
					lastFound = r
					r.touch = true
					return r
				}
				lastFound = r
			}
		}
	}
	return lastFound
}

func (note *Note) LookupWithHash(hash string, log io.Writer) *Resource {
	rsc, ok := note.Hash[hash]
	if !ok {
		fmt.Fprintf(log, "No resource matched the hash (%q).\n", hash)
		return nil
	}
	return rsc
}

func (note *Note) Lookup(
	attr map[string]string,
	noFallback bool,
	log io.Writer) (*Resource, string) {

	var reason string
	if hash, ok := attr["hash"]; ok {
		if rsc := note.LookupWithHash(hash, log); rsc != nil {
			return rsc, ""
		}
		reason = fmt.Sprintf(`<!-- Error: No resource matched the hash (%q). -->`, hash)
	} else {
		reason = `<!-- Error: no hash specified in en-media -->`
		fmt.Fprintln(log, "No hash specified in en-media.")
	}
	if !noFallback {
		if style, ok := attr["style"]; ok {
			if r := note.LookupWithStyle(style, log); r != nil {
				return r, ""
			}
		}
	}
	return nil, reason
}
