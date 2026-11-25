package serializables

import (
	"bytes"
	"encoding/xml"
	"os"

	"golang.org/x/net/html/charset"
)

type PresentationLinkbaseFile struct {
	XMLName  xml.Name   `xml:"linkbase"`
	XMLAttrs []xml.Attr `xml:",any,attr"`
	RoleRef  []struct {
		XMLName  xml.Name
		XMLAttrs []xml.Attr `xml:",any,attr"`
	} `xml:"roleRef"`
	PresentationLink []struct {
		XMLName  xml.Name
		XMLAttrs []xml.Attr `xml:",any,attr"`
		Loc      []struct {
			XMLName  xml.Name
			XMLAttrs []xml.Attr `xml:",any,attr"`
		} `xml:"loc"`
		PresentationArc []struct {
			XMLName  xml.Name
			XMLAttrs []xml.Attr `xml:",any,attr"`
		} `xml:"presentationArc"`
	} `xml:"presentationLink"`
}

func decodePresentationLinkbaseFile(xmlData []byte) (*PresentationLinkbaseFile, error) {
	reader := bytes.NewReader(xmlData)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	decoded := PresentationLinkbaseFile{}
	err := decoder.Decode(&decoded)
	if err != nil {
		return nil, err
	}
	return &decoded, nil
}

func readPresentationLinkbaseFile(filepath string) (*PresentationLinkbaseFile, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	decoded, err := decodePresentationLinkbaseFile(data)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}
