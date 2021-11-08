package api

import (
	"encoding/xml"
)

type Feed struct {
	XMLName    xml.Name `xml:"feed"`
	Text       string   `xml:",chardata"`
	OpenSearch string   `xml:"opensearch,attr"`
	Xmlns      string   `xml:"xmlns,attr"`
	Title      string   `xml:"title"`
	Subtitle   string   `xml:"subtitle"`
	Updated    string   `xml:"updated"`
	Author     struct {
		Text string `xml:",chardata"`
		Name string `xml:"name"`
	} `xml:"author"`
	ID           string `xml:"id"`
	TotalResults string `xml:"totalResults"`
	StartIndex   string `xml:"startIndex"`
	ItemsPerPage string `xml:"itemsPerPage"`
	Link         []struct {
		Href  string `xml:"href,attr"`
		Rel   string `xml:"rel,attr"`
		Title string `xml:"title,attr"`
		Type  string `xml:"type,attr"`
	} `xml:"link"`
	Entry []FeedEntry `xml:"entry"`
}

func (f *Feed) HasNextPage() bool {
	for _, v := range f.Link {
		if v.Rel == "next" {
			return true
		}
	}

	return false
}
