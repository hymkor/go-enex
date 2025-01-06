package enex

import (
	"fmt"
	"html"
	"net/url"
	"path"
	"path/filepath"
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

type Attachments struct {
	Images    map[string]*Resource
	BaseName  string
	Dir       string
	dirEscape string
	serialNo  _SerialNo
	sanitizer func(string) string
}

func newAttachments(note *Export, sanitizer func(string) string) *Attachments {
	baseName := sanitizer(note.Title)
	dir := baseName + ".files"
	dirEscape := url.PathEscape(dir)
	return &Attachments{
		Images:    make(map[string]*Resource),
		BaseName:  baseName,
		Dir:       dir,
		dirEscape: dirEscape,
		serialNo:  make(_SerialNo),
		sanitizer: sanitizer,
	}
}

func (attach *Attachments) Make(rsc *Resource) string {
	name := attach.sanitizer(attach.serialNo.ToUniqName(rsc.Mime, rsc.FileName, rsc.Hash))
	rsc.NewFileName = name
	attach.Images[filepath.Join(attach.Dir, name)] = rsc
	return path.Join(attach.dirEscape, url.PathEscape(name))
}

type Option struct {
	ExHeader  string
	Sanitizer func(string) string
}

func (exp *Export) ToHtml(makeRscUrl func(*Resource) string, opt *Option) string {
	exHeader := ""
	if opt != nil {
		exHeader = opt.ExHeader
	}
	content := "<html><head><meta charset=\"utf-8\">" +
		exHeader +
		"</head><body>" +
		"<en-note class=\"peso\" style=\"white-space: inherit;\">\n" +
		`<h1 class="noteTitle html-note" style="font-family: Source Sans Pro,-apple-system,system-ui,Segoe UI,Roboto, Oxygen,Ubuntu,Cantarell,Fira Sans,Droid Sans,Helvetica Neue,sans-serif; margin-top: 21px; margin-bottom: 21px; font-size: 32px;"><b>` +
		html.EscapeString(exp.Title) +
		"</b></h1>\n" +
		exp.Content +
		"</en-note></body></html>\n"

	var buffer strings.Builder
	for {
		m := rxMedia.FindStringSubmatchIndex(content)
		if m == nil {
			buffer.WriteString(content)
			break
		}
		buffer.WriteString(content[:m[0]])
		attr := parseEnMediaAttr(content[m[2]:m[3]])
		if hash, ok := attr["hash"]; ok {
			if rsc, ok := exp.Hash[hash]; ok {
				rscUrl := makeRscUrl(rsc)
				if strings.HasPrefix(strings.ToLower(rsc.Mime), "image") {
					fmt.Fprintf(&buffer, `<span class="goenex-attachment-image"><a href="%[1]s"><img src="%[1]s" border="0"`, rscUrl)
					if w, ok := attr["width"]; ok {
						fmt.Fprintf(&buffer, ` width="%s"`, w)
					}
					if h, ok := attr["height"]; ok {
						fmt.Fprintf(&buffer, ` height="%s"`, h)
					}
					fmt.Fprintf(&buffer, ` /></a></span>`)
				} else {
					fmt.Fprintf(&buffer, `<div class="goenex-attachment-link"><a href="%s">%s</a></div>`,
						rscUrl,
						rsc.NewFileName)
				}
			} else {
				fmt.Fprintf(&buffer, `<!-- Error: hash="%s" -->`, hash)
			}
		} else {
			fmt.Fprintf(&buffer, `<!-- Error: hash not found -->`)
		}
		content = content[m[1]:]
	}
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

func (exp *Export) Extract(opt *Option) (string, *Attachments) {
	sanitizer := defaultSanitizer.Replace
	if opt != nil && opt.Sanitizer != nil {
		sanitizer = opt.Sanitizer
	}
	attach := newAttachments(exp, sanitizer)
	content := exp.ToHtml(attach.Make, opt)
	return content, attach
}
