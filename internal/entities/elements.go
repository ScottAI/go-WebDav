package entities

import "encoding/xml"

// https://tools.ietf.org/html/rfc4918#section-15.2
type DisplayName struct {
	XMLName xml.Name `xml:"DAV: displayname"`
	Name    string   `xml:",chardata"`
}

// https://tools.ietf.org/html/rfc5397#section-3
type CurrentUserPrincipal struct {
	XMLName         xml.Name  `xml:"DAV: current-user-principal"`
	Href            Href      `xml:"href,omitempty"`
	Unauthenticated *struct{} `xml:"unauthenticated,omitempty"`
}

// https://tools.ietf.org/html/rfc4918#section-14.19
type Propertyupdate struct {
	XMLName xml.Name `xml:"DAV: propertyupdate"`
	Remove  []Remove `xml:"remove"`
	Set     []Set    `xml:"set"`
}

// https://tools.ietf.org/html/rfc4918#section-14.23
type Remove struct {
	XMLName xml.Name `xml:"DAV: remove"`
	Prop    Prop     `xml:"prop"`
}

// https://tools.ietf.org/html/rfc4918#section-14.26
type Set struct {
	XMLName xml.Name `xml:"DAV: set"`
	Prop    Prop     `xml:"prop"`
}
