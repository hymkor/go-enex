package enex

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"

	"strings"
)

type xmlResource struct {
	XMLName     xml.Name `xml:"resource"`
	Data        string   `xml:"data"`
	Mime        string   `xml:"mime"`
	Width       int      `xml:"width"`
	Height      int      `xml:"height"`
	FileName    string   `xml:"resource-attributes>file-name"`
	SourceUrl   string   `xml:"resource-attributes>source-url"`
	Recognition []byte   `xml:"recognition"`
}

type xmlRecoIndex struct {
	XMLName xml.Name `xml:"recoIndex"`
	ObjID   string   `xml:"objID,attr"`
}

type xmlEnExport struct {
	XMLName  xml.Name       `xml:"en-export"`
	Content  string         `xml:"note>content"`
	Resource []*xmlResource `xml:"note>resource"`
}

type Resource struct {
	data      string
	Mime      string
	SourceUrl string
	Hash      string
	index     int
	FileName  string
	Width     int
	Height    int
}

func (rsc *Resource) DataBeforeDecoded() string {
	return rsc.data
}

func (rsc *Resource) WriteTo(w io.Writer) (int64, error) {
	strReader := strings.NewReader(rsc.data)
	binReader := base64.NewDecoder(base64.StdEncoding, strReader)
	return io.Copy(w, binReader)
}

func (rsc *Resource) Data() []byte {
	var buffer bytes.Buffer
	rsc.WriteTo(&buffer)
	return buffer.Bytes()
}

type Export struct {
	Content  string
	Resource map[string][]*Resource // filename to the multi resources
	Hash     map[string]*Resource   // hash to the one resource
}

func Parse(data []byte) (*Export, error) {
	return ParseVerbose(data, io.Discard)
}

func ParseVerbose(data []byte, warn io.Writer) (*Export, error) {
	var theXml xmlEnExport
	err := xml.Unmarshal(data, &theXml)
	if err != nil {
		return nil, err
	}
	resource := make(map[string][]*Resource)
	hash := make(map[string]*Resource)
	for i, rsc := range theXml.Resource {
		r := &Resource{
			data:     rsc.Data,
			Mime:     strings.TrimSpace(rsc.Mime),
			index:    i,
			FileName: rsc.FileName,
			Width:    rsc.Width,
			Height:   rsc.Height,
		}
		fmt.Fprintln(warn, "Filename:", rsc.FileName)
		if len(rsc.Recognition) > 0 {
			var recoIndex xmlRecoIndex

			err = xml.Unmarshal(rsc.Recognition, &recoIndex)
			if err == nil && recoIndex.ObjID != "" {
				fmt.Fprintln(warn, "objID:", recoIndex.ObjID)
				r.Hash = recoIndex.ObjID
				hash[recoIndex.ObjID] = r
				resource[rsc.FileName] = append(resource[rsc.FileName], r)
				continue
			} else if err != nil {
				fmt.Fprintln(warn, err.Error())
			}
		}
		sourceUrl := strings.TrimSpace(rsc.SourceUrl)
		if u, err := url.QueryUnescape(sourceUrl); err == nil {
			fmt.Fprintln(warn, "Found SourceURL:", u)
			r.SourceUrl = u
			if field := strings.Fields(u); len(field) >= 4 {
				r.Hash = field[2]
				hash[field[2]] = r
			}
		} else {
			fmt.Fprintln(warn, "Can not Unescape SourceURL:", sourceUrl)
		}
		resource[rsc.FileName] = append(resource[rsc.FileName], r)
	}

	return &Export{
		Content:  strings.TrimSpace(theXml.Content),
		Resource: resource,
		Hash:     hash,
	}, nil
}
