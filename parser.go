package enex

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"io"
	"net/url"
	"strings"
)

type _EnexResource struct {
	XMLName   xml.Name `xml:"resource"`
	Data      string   `xml:"data"`
	Mime      string   `xml:"mime"`
	FileName  string   `xml:"resource-attributes>file-name"`
	SourceUrl string   `xml:"resource-attributes>source-url"`
}

type _EnexXML struct {
	XMLName  xml.Name         `xml:"en-export"`
	Content  string           `xml:"note>content"`
	Resource []*_EnexResource `xml:"note>resource"`
}

type Resource struct {
	Data      []byte
	Mime      string
	SourceUrl string
	Hash      string
	Index     int
	FileName  string
}

type Enex struct {
	Content  string
	Resource map[string][]*Resource
	Hash     map[string]*Resource
}

func Parse(in io.Reader) (*Enex, error) {
	data, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}

	var xml1 _EnexXML
	err = xml.Unmarshal(data, &xml1)
	if err != nil {
		return nil, err
	}
	resource := make(map[string][]*Resource)
	hash := make(map[string]*Resource)
	for i, rsc := range xml1.Resource {
		strReader := strings.NewReader(rsc.Data)
		binReader := base64.NewDecoder(base64.StdEncoding, strReader)
		var buffer bytes.Buffer
		io.Copy(&buffer, binReader)

		r := &Resource{
			Data:     buffer.Bytes(),
			Mime:     strings.TrimSpace(rsc.Mime),
			Index:    i,
			FileName: rsc.FileName,
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

	return &Enex{
		Content:  strings.TrimSpace(xml1.Content),
		Resource: resource,
		Hash:     hash,
	}, nil
}
