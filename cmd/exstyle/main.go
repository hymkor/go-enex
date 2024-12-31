package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type Html struct {
	XMLName xml.Name `xml:"html"`
	Style   string   `xml:"head>style"`
}

func mains(args []string) error {
	var in *os.File
	if len(args) > 0 {
		var err error
		in, err = os.Open(args[0])
		if err != nil {
			return err
		}
		defer in.Close()
	} else {
		in = os.Stdin
	}
	decoder := xml.NewDecoder(in)
	decoder.Strict = false
	decoder.AutoClose = xml.HTMLAutoClose
	decoder.Entity = xml.HTMLEntity
	decoder.Entity["nbsp"] = " "

	var htm Html
	err := decoder.Decode(&htm)
	if err != nil {
		return err
	}
	io.WriteString(os.Stdout, htm.Style)
	return nil
}

func main() {
	if err := mains(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
