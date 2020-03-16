package entities

import "encoding/xml"

//Resource Type Collection
type ResourceTypeCollection struct {
	XMLName xml.Name `xml:"collection"`
}
//Calendar Resource Type
type ResourceTypeCalendar struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:caldav calendar"`
}
//Type of Resource
type ResourceType struct {
	XMLName xml.Name `xml:"resourcetype"`
	Collection *ResourceTypeCollection `xml:"omitempty"`
	Calendar *ResourceTypeCalendar `xml:"omitempty"`
}

// https://tools.ietf.org/html/rfc4918#section-14.22
type Propstat struct {
	XMLName             xml.Name `xml:"DAV: propstat"`
	Prop                Prop     `xml:"prop"`
	Status              Status   `xml:"status"`
	ResponseDescription string   `xml:"responsedescription,omitempty"`
	Error               *Error   `xml:"error,omitempty"`
}

// https://tools.ietf.org/html/rfc4918#section-14.18
type Prop struct {
	XMLName xml.Name      `xml:"DAV: prop"`
	Raw     []OriginXMLValue `xml:",any"`
}

func EncodeProp(values ...interface{}) (*Prop, error) {
	l := make([]OriginXMLValue, len(values))
	for i, v := range values {
		raw, err := EncodeOriginXMLElement(v)
		if err != nil {
			return nil, err
		}
		l[i] = *raw
	}
	return &Prop{Raw: l}, nil
}

func (p *Prop) Get(name xml.Name) *OriginXMLValue {
	for i := range p.Raw {
		raw := &p.Raw[i]
		if n, ok := raw.XMLName(); ok && name == n {
			return raw
		}
	}
	return nil
}

func (p *Prop) Decode(v interface{}) error {
	name, err := valueXMLName(v)
	if err != nil {
		return err
	}

	raw := p.Get(name)
	if raw == nil {
		return &missingPropError{name}
	}

	return raw.Decode(v)
}
