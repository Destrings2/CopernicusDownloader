package api

import "encoding/xml"

type FeedEntry struct {
	Text  string `xml:",chardata"`
	Title string `xml:"title"`
	Link  []struct {
		Text string `xml:",chardata"`
		Href string `xml:"href,attr"`
		Rel  string `xml:"rel,attr"`
	} `xml:"link"`
	ID      string `xml:"id"`
	Summary string `xml:"summary"`
	Date    []struct {
		Text string `xml:",chardata"`
		Name string `xml:"name,attr"`
	} `xml:"date"`
	Int []struct {
		Text string `xml:",chardata"`
		Name string `xml:"name,attr"`
	} `xml:"int"`
	Str []struct {
		Text string `xml:",chardata"`
		Name string `xml:"name,attr"`
	} `xml:"str"`
}

type ODataEntry struct {
	XMLName xml.Name `xml:"entry"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	M       string   `xml:"m,attr"`
	D       string   `xml:"d,attr"`
	Base    string   `xml:"base,attr"`
	ID      string   `xml:"id"`
	Title   struct {
		Text string `xml:",chardata"`
		Type string `xml:"type,attr"`
	} `xml:"title"`
	Updated  string `xml:"updated"`
	Category struct {
		Text   string `xml:",chardata"`
		Term   string `xml:"term,attr"`
		Scheme string `xml:"scheme,attr"`
	} `xml:"category"`
	Link []struct {
		Text  string `xml:",chardata"`
		Href  string `xml:"href,attr"`
		Rel   string `xml:"rel,attr"`
		Title string `xml:"title,attr"`
		Type  string `xml:"type,attr"`
	} `xml:"link"`
	Content struct {
		Text string `xml:",chardata"`
		Type string `xml:"type,attr"`
		Src  string `xml:"src,attr"`
	} `xml:"content"`
	Properties struct {
		Text           string `xml:",chardata"`
		ID             string `xml:"Id"`
		Name           string `xml:"Name"`
		ContentType    string `xml:"ContentType"`
		ContentLength  string `xml:"ContentLength"`
		ChildrenNumber string `xml:"ChildrenNumber"`
		Value          struct {
			Text string `xml:",chardata"`
			Null string `xml:"null,attr"`
		} `xml:"Value"`
		CreationDate     string `xml:"CreationDate"`
		IngestionDate    string `xml:"IngestionDate"`
		ModificationDate string `xml:"ModificationDate"`
		EvictionDate     struct {
			Text string `xml:",chardata"`
			Null string `xml:"null,attr"`
		} `xml:"EvictionDate"`
		Online      bool `xml:"Online"`
		OnDemand    bool `xml:"OnDemand"`
		ContentDate struct {
			Text  string `xml:",chardata"`
			Type  string `xml:"type,attr"`
			Start string `xml:"Start"`
			End   string `xml:"End"`
		} `xml:"ContentDate"`
		Checksum struct {
			Text      string `xml:",chardata"`
			Type      string `xml:"type,attr"`
			Algorithm string `xml:"Algorithm"`
			Value     string `xml:"Value"`
		} `xml:"Checksum"`
	} `xml:"properties"`
}
