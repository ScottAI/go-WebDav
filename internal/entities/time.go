package entities


type Time time.Time

func (t *Time) UnmarshalText(b []byte) error {
	tt, err := http.ParseTime(string(b))
	if err != nil {
		return err
	}
	*t = Time(tt)
	return nil
}

func (t *Time) MarshalText() ([]byte, error) {
	s := time.Time(*t).Format(time.RFC1123Z)
	return []byte(s), nil
}

// https://tools.ietf.org/html/rfc4918#section-15.7
type GetLastModified struct {
	XMLName      xml.Name `xml:"DAV: getlastmodified"`
	LastModified Time     `xml:",chardata"`
}
