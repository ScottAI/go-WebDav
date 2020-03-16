package entities

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

// https://tools.ietf.org/html/rfc4918#section-14.24
type Response struct {
	XMLName             xml.Name   `xml:"DAV: response"`
	Hrefs               []Href     `xml:"href"`
	Propstats           []Propstat `xml:"propstat,omitempty"`
	ResponseDescription string     `xml:"responsedescription,omitempty"`
	Status              *Status    `xml:"status,omitempty"`
	Error               *Error     `xml:"error,omitempty"`
	Location            *Location  `xml:"location,omitempty"`
}

func NewOKResponse(path string) *Response {
	href := Href{Path: path}
	return &Response{
		Hrefs:  []Href{href},
		Status: &Status{Code: http.StatusOK},
	}
}

func (resp *Response) Path() (string, error) {
	if err := resp.Status.Err(); err != nil {
		return "", err
	}
	if len(resp.Hrefs) != 1 {
		return "", fmt.Errorf("webdav: malformed response: expected exactly one href element, got %v", len(resp.Hrefs))
	}
	return resp.Hrefs[0].Path, nil
}

type missingPropError struct {
	XMLName xml.Name
}

func (err *missingPropError) Error() string {
	return fmt.Sprintf("webdav: missing prop %q %q", err.XMLName.Space, err.XMLName.Local)
}

func IsMissingProp(err error) bool {
	_, ok := err.(*missingPropError)
	return ok
}

func (resp *Response) DecodeProp(values ...interface{}) error {
	for _, v := range values {
		// TODO wrap errors with more context (XML name)
		name, err := valueXMLName(v)
		if err != nil {
			return err
		}
		if err := resp.Status.Err(); err != nil {
			return err
		}
		for _, propstat := range resp.Propstats {
			raw := propstat.Prop.Get(name)
			if raw == nil {
				continue
			}
			if err := propstat.Status.Err(); err != nil {
				return err
			}
			return raw.Decode(v)
		}
		return &missingPropError{name}
	}

	return nil
}

func (resp *Response) EncodeProp(code int, v interface{}) error {
	raw, err := EncodeOriginXMLElement(v)
	if err != nil {
		return err
	}

	for i := range resp.Propstats {
		propstat := &resp.Propstats[i]
		if propstat.Status.Code == code {
			propstat.Prop.Raw = append(propstat.Prop.Raw, *raw)
			return nil
		}
	}

	resp.Propstats = append(resp.Propstats, Propstat{
		Status: Status{Code: code},
		Prop:   Prop{Raw: []OriginXMLValue{*raw}},
	})
	return nil
}

// https://tools.ietf.org/html/rfc4918#section-14.9
type Location struct {
	XMLName xml.Name `xml:"DAV: location"`
	Href    Href     `xml:"href"`
}

