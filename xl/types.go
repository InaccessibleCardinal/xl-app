package xl

import (
	"encoding/xml"
)

type Sheet struct {
	XMLName xml.Name `xml:"worksheet"`
	Body    HLinks   `xml:"hyperlinks"`
}

type HLinks struct {
	XMLName xml.Name `xml:"hyperlinks"`
	Links   []HLink  `xml:"hyperlink"`
}

type HLink struct {
	XMLName xml.Name `xml:"hyperlink"`
	Ref     string   `xml:"ref,attr"`
	Display string   `xml:"display,attr"`
}

type HLinksMap map[string]string

type Cell struct {
	Ref    string
	Value  string
	Header string
}

type CellMap map[string]Cell

type Row []Cell
