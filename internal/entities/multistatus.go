package entities

import (
	"encoding/xml"
	"fmt"
	"path"
)

// https://tools.ietf.org/html/rfc4918#section-14.16
type Multistatus struct {
	XMLName             xml.Name   `xml:"DAV: multistatus"`
	Responses           []Response `xml:"response"`
	ResponseDescription string     `xml:"responsedescription,omitempty"`
}

func NewMultistatus(resps ...Response) *Multistatus {
	return &Multistatus{Responses: resps}
}

func (ms *Multistatus) Get(p string) (*Response, error) {
	// Clean the path to avoid issues with trailing slashes
	p = path.Clean(p)
	for i := range ms.Responses {
		resp := &ms.Responses[i]
		for _, h := range resp.Hrefs {
			if path.Clean(h.Path) == p {
				return resp, resp.Status.Err()
			}
		}
	}

	return nil, fmt.Errorf("webdav: missing response for path %q", p)
}
