package tfl

import (
	"encoding/xml"
)

type ArrayOfLineStatus struct {
	XMLName xml.Name     `xml:"ArrayOfLineStatus"`
	Xmlns   string       `xml:"xmlns,attr"`
	Xsd     string       `xml:"xsd,attr"`
	Xsi     string       `xml:"xsi,attr"`
	Lines   []LineStatus `xml:"LineStatus"`
}

type LineStatus struct {
	ID            int               `xml:"ID,attr"`
	StatusDetails string            `xml:"StatusDetails,attr"`
	Disruptions   BranchDisruptions `xml:"BranchDisruptions"`
	Line          Line              `xml:"Line"`
	Status        Status            `xml:"Status"`
}

type BranchDisruptions struct {
	Disruptions []BranchDisruption `xml:"BranchDisruption"`
}

type BranchDisruption struct {
	StationTo   Station `xml:"StationTo"`
	StationFrom Station `xml:"StationFrom"`
	Status      Status  `xml:"Status"`
}

type Station struct {
	ID   int    `xml:"ID,attr"`
	Name string `xml:"Name,attr"`
}

type Line struct {
	ID   int    `xml:"ID,attr"`
	Name string `xml:"Name,attr"`
}

type Status struct {
	ID          string     `xml:"ID,attr"`
	CssClass    string     `xml:"CssClass,attr"`
	Description string     `xml:"Description,attr"`
	IsActive    bool       `xml:"IsActive,attr"`
	StatusType  StatusType `xml:"StatusType"`
}

type StatusType struct {
	ID          int    `xml:"ID,attr"`
	Description string `xml:"Description,attr"`
}

type TFLParsed struct {
	Line   string
	Status string
}
