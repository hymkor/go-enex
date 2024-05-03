package enex

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	rxMedia   = regexp.MustCompile(`\s*<en-media[^>]*hash="([^"]*)"[^>]*/>\s*`)
	rxXml     = regexp.MustCompile(`\s*<\?xml[^>]*>\s*`)
	rxDocType = regexp.MustCompile(`\s*<!DOCTYPE[^>]*>\s*`)
)

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

// Deprecated: use ToHtml
func (exp *Export) Html(imagePathHeader string) (html string, images map[string][]byte) {
	html, rsc := exp.HtmlAndImagesWithRenamer(DefaultRenamer(imagePathHeader))
	images = map[string][]byte{}
	for name, r := range rsc {
		images[name] = r.Data()
	}
	return html, images
}

var rxUrl = regexp.MustCompile(`^\w\w+\:`)

// Deprecated: use ToHtml
func (exp *Export) HtmlAndImagesWithRenamer(_renamer func(string, int) string) (string, map[string]*Resource) {

	images := make(map[string]*Resource)
	html := exp.ToHtml(func(rsc *Resource) string {
		fname := _renamer(rsc.FileName, rsc.index)
		images[fname] = rsc
		if rxUrl.MatchString(fname) {
			return fname
		}
		dir := path.Dir(fname)
		base := path.Base(fname)
		if dir == "" || dir == "." {
			return url.PathEscape(fname)
		}
		return path.Join(url.PathEscape(dir), url.PathEscape(base))
	})
	return html, images
}

func (exp *Export) ToHtml(mkImgSrc func(*Resource) string) string {
	html := exp.Content
	html = rxXml.ReplaceAllString(html, "")
	html = rxDocType.ReplaceAllString(html, "<!DOCTYPE html>")
	html = strings.ReplaceAll(html, "<en-note>",
		"<html><head><meta charset=\"utf-8\"></head><body>\n")
	html = strings.ReplaceAll(html, "</en-note>", "</body></html>\n")
	html = strings.ReplaceAll(html, "<div><br /></div>", "<br />")
	html = rxLiDiv.ReplaceAllString(html, `<li>${1}</li>`)
	html = rxBrSomething.ReplaceAllString(html, `${1}`)

	var buffer strings.Builder
	for {
		m := rxMedia.FindStringSubmatchIndex(html)
		if m == nil {
			buffer.WriteString(html)
			break
		}
		buffer.WriteString(html[:m[0]])
		hash := html[m[2]:m[3]]

		if rsc, ok := exp.Hash[hash]; ok {
			imgsrc := mkImgSrc(rsc)
			fmt.Fprintf(&buffer,
				`<img src="%s" width="%d" height="%d" />`,
				imgsrc,
				rsc.Width,
				rsc.Height)
		} else {
			fmt.Fprintf(&buffer, `<!-- Error: hash="%s" -->`, hash)
		}
		html = html[m[1]:]
	}
	return buffer.String()
}
