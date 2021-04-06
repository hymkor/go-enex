package main

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"os"
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

type Enex struct {
	Content  string
	Resource map[string][][]byte
}

func ReadEnex(in io.Reader) (*Enex, error) {
	data, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}

	var xml1 _EnexXML
	err = xml.Unmarshal(data, &xml1)
	if err != nil {
		return nil, err
	}
	resource := make(map[string][][]byte)
	for _, rsc := range xml1.Resource {
		strReader := strings.NewReader(rsc.Data)
		binReader := base64.NewDecoder(base64.StdEncoding, strReader)
		var buffer bytes.Buffer
		io.Copy(&buffer, binReader)
		resource[rsc.FileName] = append(resource[rsc.FileName], buffer.Bytes())
	}
	return &Enex{Content: xml1.Content, Resource: resource}, nil
}

func mains() error {
	enex, err := ReadEnex(os.Stdin)
	if err != nil {
		return err
	}
	for fname, bins := range enex.Resource {
		for i, bin := range bins {
			name := fname
			if i > 0 {
				name = fmt.Sprintf("%d-%s", i, fname)
			}
			os.WriteFile(name, bin, 0666)
		}
	}
	return nil
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
