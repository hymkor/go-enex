package enex

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
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
	XMLName xml.Name   `xml:"en-export"`
	Note    []*xmlNote `xml:"note"`
}

type xmlNote struct {
	XMLName  xml.Name       `xml:"note"`
	Title    string         `xml:"title"`
	Content  string         `xml:"content"`
	Resource []*xmlResource `xml:"resource"`
}

type Resource struct {
	data        string
	Mime        string
	SourceUrl   string
	Hash        string
	FileName    string
	Width       int
	Height      int
	NewFileName string
}

func (rsc *Resource) DataBeforeDecoded() string {
	return rsc.data
}

func (rsc *Resource) WriteTo(w io.Writer) (int64, error) {
	strReader := strings.NewReader(strings.TrimSpace(rsc.data))
	binReader := base64.NewDecoder(base64.StdEncoding, strReader)
	return io.Copy(w, binReader)
}

func (rsc *Resource) Data() []byte {
	var buffer bytes.Buffer
	rsc.WriteTo(&buffer)
	return buffer.Bytes()
}

type Export struct {
	Title    string
	Content  string
	Resource map[string][]*Resource // filename to the multi resources
	Hash     map[string]*Resource   // hash to the one resource
	ExHeader string
}

func Parse(data []byte) (*Export, error) {
	return ParseVerbose(data, io.Discard)
}

func ParseVerbose(data []byte, warn io.Writer) (*Export, error) {
	exports, err := ParseMulti(data, warn)
	if err != nil {
		return nil, err
	}
	if len(exports) >= 2 {
		return nil, errors.New("ParseVerbose: not support multi notes")
	}
	if len(exports) <= 0 {
		return nil, errors.New("ParseVerbose: zero notes")
	}
	return exports[0], nil
}

func ParseMulti(data []byte, warn io.Writer) ([]*Export, error) {
	var theXml xmlEnExport
	err := xml.Unmarshal(data, &theXml)
	if err != nil {
		return nil, err
	}
	exports := make([]*Export, 0, len(theXml.Note))
	for _, note := range theXml.Note {
		resource := make(map[string][]*Resource)
		hash := make(map[string]*Resource)
		for _, rsc := range note.Resource {
			r := &Resource{
				data:     rsc.Data,
				Mime:     strings.TrimSpace(rsc.Mime),
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
		exports = append(exports, &Export{
			Title:    note.Title,
			Content:  strings.TrimSpace(note.Content),
			Resource: resource,
			Hash:     hash,
		})
	}
	return exports, nil
}
