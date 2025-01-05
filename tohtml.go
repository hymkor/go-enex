package enex

import (
	"fmt"
	htmlPkg "html"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var (
	rxXml     = regexp.MustCompile(`\s*<\?xml[^>]*>\s*`)
	rxDocType = regexp.MustCompile(`(?s)\s*<!DOCTYPE[^>]*>\s*`)
	rxMedia   = regexp.MustCompile(`(?s)\s*<en-media([^>]*)>\s*`)
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

func findWidth(attr map[string]string) string {
	if value, ok := attr["width"]; ok {
		return value
	}
	return ""
}

func findHeight(attr map[string]string) string {
	if value, ok := attr["height"]; ok {
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

func (s SerialNo) ToUniqName(mime, name string, index int) string {
	if name == "" {
		mainType, subType, _ := strings.Cut(mime, "/")
		if strings.EqualFold(mainType, "image") {
			if subType == "jpeg" {
				name = "image.jpg"
			} else {
				name = "image." + subType
			}
		} else {
			name = "Evernote"
		}
	}
	uname := strings.ToUpper(name)
	indexList, ok := s[uname]
	if !ok {
		// New table and no need to rename
		s[uname] = []int{index}
		return name
	}
	serial := slices.Index(indexList, index)
	if serial == 0 {
		// Modifying and count-up are not needed
		return name
	}
	if serial < 0 {
		// Count-up and update table is required
		serial = len(indexList)
		s[uname] = append(indexList, index)
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
	name := ToSafe.Replace(imgSrc.serialNo.ToUniqName(rsc.Mime, rsc.FileName, rsc.Index))
	rsc.NewFileName = name
	imgSrc.Images[filepath.Join(imgSrc.Dir, name)] = rsc
	return path.Join(imgSrc.dirEscape, url.PathEscape(name))
}

func (exp *Export) ToHtml(imgSrc interface{ Make(*Resource) string }) string {
	html := exp.Content
	html = rxXml.ReplaceAllString(html, "")
	html = rxDocType.ReplaceAllString(html, "<!DOCTYPE html>")
	html = strings.ReplaceAll(html, "<en-note>",
		"<html><head><meta charset=\"utf-8\">"+
			exp.ExHeader+
			"</head><body>"+
			"<en-note class=\"peso\" style=\"white-space: inherit;\">\n"+
			`<h1 class="noteTitle html-note" style="font-family: Source Sans Pro,-apple-system,system-ui,Segoe UI,Roboto, Oxygen,Ubuntu,Cantarell,Fira Sans,Droid Sans,Helvetica Neue,sans-serif; margin-top: 21px; margin-bottom: 21px; font-size: 32px;"><b>`+
			htmlPkg.EscapeString(exp.Title)+
			"</b></h1>\n")
	html = strings.ReplaceAll(html, "</en-note>", "</en-note></body></html>\n")

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
				if strings.HasPrefix(strings.ToLower(rsc.Mime), "image") {
					fmt.Fprintf(&buffer, `<span class="goenex-attachment-image"><a href="%[1]s"><img src="%[1]s" border="0"`, imgsrc1)
					if w := findWidth(attr); w != "" {
						fmt.Fprintf(&buffer, ` width="%s"`, w)
					}
					if h := findHeight(attr); h != "" {
						fmt.Fprintf(&buffer, ` height="%s"`, h)
					}
					fmt.Fprintf(&buffer, ` /></a></span>`)
				} else {
					fmt.Fprintf(&buffer, `<div class="goenex-attachment-link"><a href="%s">%s</a></div>`,
						imgsrc1,
						rsc.NewFileName)
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
