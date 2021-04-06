package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var rxMedia = regexp.MustCompile(`<en-media[^>]*hash="([^"]*)"[^>]*/>`)

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

func mains() error {
	enex, err := ReadEnex(os.Stdin)
	if err != nil {
		return err
	}
	c := strings.ReplaceAll(enex.Content, "</div>", "</div>\n")
	c = strings.ReplaceAll(c, "</ul>", "</ul>\n")
	c = strings.ReplaceAll(c, "</li>", "</li>\n")
	c = convMediaTag(c, func(hash string) string {
		if rsc, ok := enex.Hash[hash]; ok {
			return fmt.Sprintf(`<img alt="%[1]s" src="%[2]s" />`,
				hash,
				renameWithNumber(rsc.FileName, rsc.Index))
		} else {
			return "<!-- Error -->"
		}
	})

	io.WriteString(os.Stdout, c)
	for fname, rscs := range enex.Resource {
		for i, rsc := range rscs {
			name := renameWithNumber(fname, i)
			fmt.Fprintf(os.Stderr, "%s: %s\n", name, rsc.Hash)
			os.WriteFile(name, rsc.Data, 0666)
		}
	}
	return nil
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
