package enex

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"io"
	"net/url"
	"strings"
)

type xmlResource struct {
	XMLName   xml.Name `xml:"resource"`
	Data      string   `xml:"data"`
	Mime      string   `xml:"mime"`
	Width     int      `xml:"width"`
	Height    int      `xml:"height"`
	FileName  string   `xml:"resource-attributes>file-name"`
	SourceUrl string   `xml:"resource-attributes>source-url"`
}

type xmlEnExport struct {
	XMLName  xml.Name       `xml:"en-export"`
	Content  string         `xml:"note>content"`
	Resource []*xmlResource `xml:"note>resource"`
}

type Resource struct {
	Data      []byte
	Mime      string
	SourceUrl string
	Hash      string
	index     int
	FileName  string
	Width     int
	Height    int
}

type Export struct {
	Content  string
	Resource map[string][]*Resource // filename to the multi resources
	Hash     map[string]*Resource   // hash to the one resource
}

func Parse(data []byte) (*Export, error) {
	var theXml xmlEnExport
	err := xml.Unmarshal(data, &theXml)
	if err != nil {
		return nil, err
	}
	resource := make(map[string][]*Resource)
	hash := make(map[string]*Resource)
	for i, rsc := range theXml.Resource {
		strReader := strings.NewReader(rsc.Data)
		binReader := base64.NewDecoder(base64.StdEncoding, strReader)
		var buffer bytes.Buffer
		io.Copy(&buffer, binReader)

		r := &Resource{
			Data:     buffer.Bytes(),
			Mime:     strings.TrimSpace(rsc.Mime),
			index:    i,
			FileName: rsc.FileName,
			Width:    rsc.Width,
			Height:   rsc.Height,
		}
		sourceUrl := strings.TrimSpace(rsc.SourceUrl)
		if u, err := url.QueryUnescape(sourceUrl); err == nil {
			r.SourceUrl = u
			if field := strings.Fields(u); len(field) >= 4 {
				r.Hash = field[2]
				hash[field[2]] = r
			}
		}
		resource[rsc.FileName] = append(resource[rsc.FileName], r)
	}

	return &Export{
		Content:  strings.TrimSpace(theXml.Content),
		Resource: resource,
		Hash:     hash,
	}, nil
}
