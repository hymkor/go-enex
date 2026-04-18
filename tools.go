package enex

import (
	"bytes"
	"encoding/xml"
	"html"
	"regexp"
	"strconv"
	"unicode/utf8"
)

var entityTables = [][2][]byte{
	[2][]byte{[]byte("&lt;"), []byte("\uE001")},
	[2][]byte{[]byte("&gt;"), []byte("\uE002")},
	[2][]byte{[]byte("&amp;"), []byte("\uE003")},
	[2][]byte{[]byte("&quot;"), []byte("\uE004")},
	[2][]byte{[]byte("&apos;"), []byte("\uE005")},
}

var reEntity = regexp.MustCompile(`&[a-zA-z][a-zA-Z0-9]*;`)

func xmlUnmarshal(data []byte, v any) error {
	for _, v := range entityTables {
		data = bytes.ReplaceAll(data, v[0], v[1])
	}
	data = reEntity.ReplaceAllFunc(data, func(srcBin []byte) []byte {
		srcStr := string(srcBin)
		decoded := html.UnescapeString(srcStr)
		if decoded == srcStr || decoded == "" {
			return srcBin
		}
		r, siz := utf8.DecodeRuneInString(decoded)
		if r == utf8.RuneError || len(decoded) > siz {
			debug("failed:", srcStr)
			return srcBin
		}
		bin := make([]byte, 0, 10)
		bin = append(bin, '&', '#')
		bin = strconv.AppendInt(bin, int64(r), 10)
		bin = append(bin, ';')
		debug(srcStr, "->", string(bin))
		return bin
	})
	for _, v := range entityTables {
		data = bytes.ReplaceAll(data, v[1], v[0])
	}
	return stackTrace(xml.Unmarshal(data, v), "xmlUnmarshal")
}
