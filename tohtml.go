package enex

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	rxMedia   = regexp.MustCompile(`\s*<en-media[^>]*hash="([^"]*)"[^>]*/>\s*`)
	rxXml     = regexp.MustCompile(`\s*<\?xml[^>]*>\s*`)
	rxDocType = regexp.MustCompile(`\s*<!DOCTYPE[^>]*>\s*`)
)

func convMediaTag(s string, hash2tag func(string) string) string {
	var buffer strings.Builder
	for {
		m := rxMedia.FindStringSubmatchIndex(s)
		if m == nil {
			buffer.WriteString(s)
			break
		}
		buffer.WriteString(s[:m[0]])
		buffer.WriteString(hash2tag(s[m[2]:m[3]]))
		s = s[m[1]:]
	}
	return buffer.String()
}

func renameWithNumber(fname string, n int) string {
	if n <= 0 {
		return fname
	}
	ext := filepath.Ext(fname)
	base := fname[:len(fname)-len(ext)]
	return fmt.Sprintf("%s_%d%s", base, n, ext)
}

func (enex *Enex) Html(prefix string) (string, map[string][]byte) {
	resources := map[string][]byte{}
	c := enex.Content
	c = rxXml.ReplaceAllString(c, "")
	c = rxDocType.ReplaceAllString(c, "<!DOCTYPE html>")
	c = strings.ReplaceAll(c, "<en-note>",
		"<html><head><meta charset=\"utf-8\"></head><body>\n")
	c = strings.ReplaceAll(c, "</en-note>", "</body></html>\n")
	c = strings.ReplaceAll(c, "</div>", "</div>\n")
	c = strings.ReplaceAll(c, "</ul>", "</ul>\n")
	c = strings.ReplaceAll(c, "</li>", "</li>\n")
	c = convMediaTag(c, func(hash string) string {
		if rsc, ok := enex.Hash[hash]; ok {
			fname := prefix + renameWithNumber(rsc.FileName, rsc.Index)
			resources[fname] = rsc.Data
			return fmt.Sprintf(`<img src="%s" /><!-- hash="%s" -->`, fname, hash)
		} else {
			return fmt.Sprintf(`<!-- Error: hash="%s" -->`, hash)
		}
	})
	return c, resources
}
