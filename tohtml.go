package enex

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var (
	rxXml         = regexp.MustCompile(`\s*<\?xml[^>]*>\s*`)
	rxDocType     = regexp.MustCompile(`(?s)\s*<!DOCTYPE[^>]*>\s*`)
	rxDivBrDiv    = regexp.MustCompile(`(?s)<div>\s*<br\s*/>\s*</div>`)
	rxDivBrDiv2   = regexp.MustCompile(`(?s)</div>\s*<br\s*/>\s*<div>`)
	rxLiDiv       = regexp.MustCompile(`(?s)<li>\s*<div>([^<>]*)</div>\s*</li>`)
	rxBrSomething = regexp.MustCompile(`(?s)<br\s*/>\s*(<(?:(?:div)|(?:ol)|(?:ul)))`)
	rxMedia       = regexp.MustCompile(`(?s)\s*<en-media[^>]*hash="([^"]*)"[^>]*/>\s*`)
	rxEnds        = regexp.MustCompile(`(?s)</(?:(?:div)|(?:p))>`)
)

var ToSafe = strings.NewReplacer(
	`<`, `＜`,
	`>`, `＞`,
	`"`, `”`,
	`/`, `／`,
	`\`, `＼`,
	`|`, `｜`,
	`?`, `？`,
	`*`, `＊`,
	`:`, `：`,
	`(`, `（`,
	`)`, `）`,
	` `, `_`,
)

type SerialNo map[string][]int

func (s SerialNo) ToUniqName(name string, index int) string {
	indexList, ok := s[name]
	if !ok {
		// New table and no need to rename
		s[name] = []int{index}
		return name
	}
	serial := slices.Index(indexList, index)
	if serial == 0 {
		// Modifying and count-up are not needed
		return name
	}
	if serial < 0 {
		// Count-up and update table is required
		indexList = append(indexList, index)
		serial = len(indexList)
		s[name] = indexList
	}
	ext := path.Ext(name)
	base := name[:len(name)-len(ext)]
	return fmt.Sprintf("%s (%d)%s", base, serial, ext)
}

type ImgSrc struct {
	Images    map[string]*Resource
	baseName  string
	Dir       string
	dirEscape string
	serialNo  SerialNo
}

func NewImgSrc(note *Export) *ImgSrc {
	baseName := ToSafe.Replace(note.Title)
	dir := baseName + ".files"
	dirEscape := url.PathEscape(dir)
	return &ImgSrc{
		Images:    make(map[string]*Resource),
		baseName:  baseName,
		Dir:       dir,
		dirEscape: dirEscape,
		serialNo:  make(SerialNo),
	}
}

func (imgSrc *ImgSrc) Make(rsc *Resource) string {
	name := ToSafe.Replace(imgSrc.serialNo.ToUniqName(rsc.FileName, rsc.Index))
	imgSrc.Images[filepath.Join(imgSrc.Dir, name)] = rsc
	return path.Join(imgSrc.dirEscape, url.PathEscape(name))
}

func (exp *Export) ToHtml(imgSrc interface{ Make(*Resource) string }) string {
	html := exp.Content
	html = rxXml.ReplaceAllString(html, "")
	html = rxDocType.ReplaceAllString(html, "<!DOCTYPE html>")
	html = strings.ReplaceAll(html, "<en-note>",
		"<html><head><meta charset=\"utf-8\"></head><body>\n")
	html = strings.ReplaceAll(html, "</en-note>", "</body></html>\n")
	html = rxDivBrDiv.ReplaceAllString(html, "<br/>\n")
	html = rxDivBrDiv2.ReplaceAllString(html, "</div><div>")
	html = rxLiDiv.ReplaceAllString(html, "<li>${1}</li>\n")
	html = rxBrSomething.ReplaceAllString(html, `${1}`)
	html = rxEnds.ReplaceAllString(html, "${0}\n")

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
			imgsrc1 := imgSrc.Make(rsc)
			switch strings.ToUpper(filepath.Ext(imgsrc1)) {
			case ".JPG", ".JPEG", ".PNG", ".GIF":
				fmt.Fprintf(&buffer,
					`<img src="%s" width="%d" height="%d" />`,
					imgsrc1,
					rsc.Width,
					rsc.Height)
			default:
				fmt.Fprintf(&buffer, `<a href="%s">%s</a>`,
					imgsrc1,
					filepath.Base(imgsrc1))
			}
		} else {
			fmt.Fprintf(&buffer, `<!-- Error: hash="%s" -->`, hash)
		}
		html = html[m[1]:]
	}
	return buffer.String()
}

func (exp *Export) HtmlAndDir() (string, *ImgSrc) {
	imgSrc := NewImgSrc(exp)
	html := exp.ToHtml(imgSrc)
	return html, imgSrc
}
