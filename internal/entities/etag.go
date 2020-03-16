package entities

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

type ETag string

// https://tools.ietf.org/html/rfc4918#section-15.6
type GetETag struct {
	XMLName xml.Name `xml:"DAV: getetag"`
	ETag    ETag     `xml:",chardata"`
}

func (etag *ETag) UnmarshalText(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return fmt.Errorf("webdav: failed to unquote ETag: %v", err)
	}
	*etag = ETag(s)
	return nil
}

func (etag ETag) MarshalText() ([]byte, error) {
	return []byte(etag.String()), nil
}

func (etag ETag) String() string {
	return fmt.Sprintf("%q", string(etag))
}

