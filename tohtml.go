package enex

import (
	"fmt"
	"net/url"
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

var rxLiDiv = regexp.MustCompile(`<li><div>([^<>]*)</div></li>`)

var rxBrSomething = regexp.MustCompile(`<br />(<(?:(?:div)|(?:ol)|(?:ul)))`)

func DefaultRenamer(imagePathHeader string) func(string, int) string {
	return func(baseName string, n int) string {
		return imagePathHeader + renameWithNumber(baseName, n)
	}
}

func (exp *Export) Html(imagePathHeader string) (html string, images map[string][]byte) {
	html, rsc := exp.HtmlAndImagesWithRenamer(DefaultRenamer(imagePathHeader))
	images = map[string][]byte{}
	for name, r := range rsc {
		images[name] = r.Data()
	}
	return html, images
}

func (exp *Export) HtmlAndImagesWithRenamer(renamer func(string, int) string) (html string, images map[string]*Resource) {

	html = exp.Content
	html = rxXml.ReplaceAllString(html, "")
	html = rxDocType.ReplaceAllString(html, "<!DOCTYPE html>")
	html = strings.ReplaceAll(html, "<en-note>",
		"<html><head><meta charset=\"utf-8\"></head><body>\n")
	html = strings.ReplaceAll(html, "</en-note>", "</body></html>\n")
	html = strings.ReplaceAll(html, "<div><br /></div>", "<br />")
	html = rxLiDiv.ReplaceAllString(html, `<li>${1}</li>`)
	html = rxBrSomething.ReplaceAllString(html, `${1}`)

	images = make(map[string]*Resource)

	html = convMediaTag(html, func(hash string) string {
		if rsc, ok := exp.Hash[hash]; ok {
			fname := renamer(rsc.FileName, rsc.index)
			images[fname] = rsc
			return fmt.Sprintf(
				`<img alt="%[1]s" src="%[1]s" width="%[2]d" height="%[3]d" />`,
				url.QueryEscape(fname),
				rsc.Width,
				rsc.Height)
		} else {
			return fmt.Sprintf(`<!-- Error: hash="%s" -->`, hash)
		}
	})
	return html, images
}
