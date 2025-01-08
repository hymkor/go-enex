package enex

import (
	"fmt"
	"path"
	"regexp"
	"strings"
	"unicode/utf8"
)

var rxMedia = regexp.MustCompile(`(?s)\s*<en-media([^>]*)>\s*`)

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

type _SerialNo map[string]map[string]int

func (s _SerialNo) ToUniqName(mime, name string, hash string) string {
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
	if !ok || len(indexList) <= 0 {
		// New table and no need to rename
		s[uname] = map[string]int{hash: 0}
		return name
	}
	serial, ok := indexList[hash]
	if !ok {
		// Count-up and update table is required
		serial = len(indexList)
		s[uname][hash] = serial
	}
	if serial == 0 {
		return name
	}
	ext := path.Ext(name)
	base := name[:len(name)-len(ext)]
	return fmt.Sprintf("%s (%d)%s", base, serial, ext)
}

type Option struct {
	ExHeader  string
	Sanitizer func(string) string
}

const (
	_BALLOT_BOX            = " \u2610 "
	_BALLOT_BOX_WITH_CHECK = " \u2611 "
)

var enTagReplacer = strings.NewReplacer(
	`<en-todo checked="false" />`, _BALLOT_BOX,
	`<en-todo checked="true" />`, _BALLOT_BOX_WITH_CHECK,
)

func (note *Note) extract(makeRscUrl func(*Resource) string, opt *Option) string {
	var buffer strings.Builder

	buffer.WriteString("<!DOCTYPE html><html><head><meta charset=\"utf-8\">")
	if opt != nil {
		buffer.WriteString(opt.ExHeader)
	}
	buffer.WriteString(`</head><body>
<en-note class="peso" style="white-space: inherit;">
<h1 class="noteTitle html-note" style="font-family: Source Sans Pro,-apple-system,system-ui,Segoe UI,Roboto, Oxygen,Ubuntu,Cantarell,Fira Sans,Droid Sans,Helvetica Neue,sans-serif; margin-top: 21px; margin-bottom: 21px; font-size: 32px;"><b>`)
	buffer.WriteString(note.Title)
	buffer.WriteString("</b></h1>\n")

	content := enTagReplacer.Replace(note.Content)
	buffer.WriteString(rxMedia.ReplaceAllStringFunc(content, func(tag string) string {
		attr := parseEnMediaAttr(tag)
		hash, ok := attr["hash"]
		if !ok {
			return `<!-- Error: hash not found -->`
		}
		rsc, ok := note.Hash[hash]
		if !ok {
			return fmt.Sprintf(`<!-- Error: hash="%s" -->`, hash)
		}
		rscUrl := makeRscUrl(rsc)
		typ, _, ok := strings.Cut(rsc.Mime, "/")
		if ok && strings.EqualFold(typ, "image") {
			// image
			var b strings.Builder
			fmt.Fprintf(&b, `<span class="goenex-attachment-image"><a href="%[1]s"><img src="%[1]s" border="0"`, rscUrl)
			if w, ok := attr["width"]; ok {
				fmt.Fprintf(&b, ` width="%s"`, w)
			}
			if h, ok := attr["height"]; ok {
				fmt.Fprintf(&b, ` height="%s"`, h)
			}
			b.WriteString(" /></a></span>")
			return b.String()
		}
		// non-image attachment
		return fmt.Sprintf(`<div class="goenex-attachment-link"><a href="%s">%s</a></div>`,
			rscUrl,
			rsc.NewFileName)

	}))
	buffer.WriteString("</en-note></body></html>\n")
	return buffer.String()
}

var defaultSanitizer = strings.NewReplacer(
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

func (note *Note) Extract(opt *Option) (html string, bundle *Bundle) {
	sanitizer := defaultSanitizer.Replace
	if opt != nil && opt.Sanitizer != nil {
		sanitizer = opt.Sanitizer
	}
	bundle = newBundle(note, sanitizer)
	html = note.extract(bundle.makeUrlFor, opt)
	return
}
