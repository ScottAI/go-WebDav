package go_WebDav

import (
	"encoding/xml"
	"fmt"
	"mime"
	"net/http"
)

type HttpError struct {
	Code int
	Err error
}

func HttpErrorTrans(err error) *HttpError {
	if err == nil {
		return nil
	}

	return &HttpError{http.StatusInternalServerError,err}

}

func IsNotFound(err error) bool {
	return HttpErrorTrans(err).Code == http.StatusNotFound
}

func HttpErrorf(code int,format string,a ...interface{}) *HttpError {
	return &HttpError{code,fmt.Errorf(format,a...)}
}

func (err *HttpError) Error() string {
	s := fmt.Sprintf("%v %v",err.Code,http.StatusText(err.Code))
	if err.Err != nil {
		return fmt.Sprintf("%v:%v",s,err.Err)
	}else {
		return s
	}
}

func ServeError(w http.ResponseWriter,err error) {
	code := http.StatusInternalServerError
	if httpErr,ok := err.(*HttpError);ok {
		code = httpErr.Code
	}
	http.Error(w,err.Error(),code)
}

func DecodeXMLRequest(r *http.Request,v interface{}) error {
	t,_,_ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if t != "application/xml" && t != "text/xml" {
		return HttpErrorf(http.StatusBadRequest,"webdav: expected application/xml request")
	}

	if err := xml.NewDecoder(r.Body).Decode(v); err != nil {
		return &HttpError{http.StatusBadRequest,err}
	}
	return nil
}

func ServeXML(w http.ResponseWriter) *xml.Encoder {
	w.Header().Add("Content-Type","text/xml;charset=\"utf-8\"")
	w.Write([]byte(xml.Header))
	return xml.NewEncoder(w)
}

func ServeMultistatus(w http.ResponseWriter,ms *Multistatus)  {
	
}
