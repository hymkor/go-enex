package enex

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var (
	rxXml         = regexp.MustCompile(`\s*<\?xml[^>]*>\s*`)
	rxDocType     = regexp.MustCompile(`(?s)\s*<!DOCTYPE[^>]*>\s*`)
	rxDivBrDiv    = regexp.MustCompile(`(?s)<div>\s*<br\s*/>\s*</div>`)
	rxDivBrDiv2   = regexp.MustCompile(`(?s)</div>\s*<br\s*/>\s*<div>`)
	rxLiDiv       = regexp.MustCompile(`(?s)<li>\s*<div>([^<>]*)</div>\s*</li>`)
	rxBrSomething = regexp.MustCompile(`(?s)<br\s*/>\s*(<(?:(?:div)|(?:ol)|(?:ul)))`)
	rxMedia       = regexp.MustCompile(`(?s)\s*<en-media([^>]*)>\s*`)
	rxEnds        = regexp.MustCompile(`(?s)</(?:(?:div)|(?:p))>`)
)

func parseEnMediaAttr(s string) map[string]string {
	result := map[string]string{}
	for len(s) > 0 {
		// skip spaces
		c, siz := utf8.DecodeRuneInString(s)
		s = s[siz:]
		if strings.ContainsRune(" \v\t\r\n", c) {
			continue
		}
		var name strings.Builder
		for {
			if c == '=' {
				var value strings.Builder
				q := false
				for {
					if len(s) <= 0 {
						result[name.String()] = value.String()
						break
					}
					c, siz = utf8.DecodeRuneInString(s)
					s = s[siz:]
					if c == '"' {
						q = !q
					} else if !q && strings.ContainsRune(" \v\t\r\n", c) {
						result[name.String()] = value.String()
						break
					} else {
						value.WriteRune(c)
					}
				}
				break
			}
			if strings.ContainsRune(" \v\t\r\n", c) {
				result[name.String()] = ""
				break
			}
			name.WriteRune(c)
			if len(s) <= 0 {
				result[name.String()] = ""
				return result
			}
			c, siz = utf8.DecodeRuneInString(s)
			s = s[siz:]
		}
	}
	return result
}

func parseEnMediaStyle(style string) map[string]string {
	result := map[string]string{}
	for {
		var eq string
		var ok bool
		eq, style, ok = strings.Cut(style, ";")
		if eq != "" {
			name, value, _ := strings.Cut(eq, ":")
			result[strings.TrimSpace(name)] = strings.TrimSpace(value)
		}
		if !ok {
			return result
		}
	}
}

func findWidth(attr, style map[string]string) string {
	if value, ok := attr["width"]; ok {
		return value
	}
	if value, ok := style["--en-naturalWidth"]; ok {
		return value
	}
	return ""
}

func findHeight(attr, style map[string]string) string {
	if value, ok := attr["height"]; ok {
		return value
	}
	if value, ok := style["--en-naturalHeight"]; ok {
		return value
	}
	return ""
}

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
		attr := parseEnMediaAttr(html[m[2]:m[3]])
		if hash, ok := attr["hash"]; ok {
			if rsc, ok := exp.Hash[hash]; ok {
				imgsrc1 := imgSrc.Make(rsc)
				switch strings.ToUpper(filepath.Ext(imgsrc1)) {
				case ".JPG", ".JPEG", ".PNG", ".GIF":
					fmt.Fprintf(&buffer, `<a href="%[1]s"><img src="%[1]s" border="0"`, imgsrc1)
					style := parseEnMediaStyle(attr["style"])
					if w := findWidth(attr, style); w != "" {
						fmt.Fprintf(&buffer, ` width="%s"`, w)
					} else {
						fmt.Fprintf(&buffer, ` width="%d"`, rsc.Width)
					}
					if h := findHeight(attr, style); h != "" {
						fmt.Fprintf(&buffer, ` height="%s"`, h)
					} else {
						fmt.Fprintf(&buffer, ` height="%d"`, rsc.Height)
					}
					fmt.Fprintf(&buffer, ` /></a>`)
				default:
					fmt.Fprintf(&buffer, `<a href="%s">%s</a>`,
						imgsrc1,
						filepath.Base(imgsrc1))
				}
			} else {
				fmt.Fprintf(&buffer, `<!-- Error: hash="%s" -->`, hash)
			}
		} else {
			fmt.Fprintf(&buffer, `<!-- Error: hash not found -->`)
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
