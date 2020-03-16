package entities

import (
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strings"
)

//implements xml.Unmarshaler and xml.Marshaler
//延迟XML解码、预处理
type OriginXMLValue struct {
	token xml.Token//确保不是endelement
	children []OriginXMLValue

	//缓存输出数据
	out interface{}
}

//CreateOriginXMLValue:为元素创建OriginXMLValue
func CreateOriginXMLValue(name xml.Name,attrs []xml.Attr,children []OriginXMLValue) *OriginXMLValue{
	return &OriginXMLValue{token:xml.StartElement{name,attrs},children:children}
}

func EncodeOriginXMLElement(value interface{}) (*OriginXMLValue,error)  {
	return &OriginXMLValue{out:value},nil
}

// UnmarshalXML implements xml.Unmarshaler.
func (val *OriginXMLValue) UnmarshalXML(d *xml.Decoder,start xml.StartElement) error  {
	val.token = start
	val.children = nil
	val.out = nil

	for{
		token,err := d.Token()
		if err != nil{
			return err
		}
		switch token := token.(type) {
		case xml.StartElement:
			child := OriginXMLValue{}
			if err := child.UnmarshalXML(d,token);err != nil {
				return err
			}
			val.children = append(val.children,child)
		case xml.EndElement:
			return nil
		default:
			val.children = append(val.children,OriginXMLValue{token:xml.CopyToken(token)})
		}
	}
}

// MarshalXML implements xml.Marshaler.
func (val *OriginXMLValue) MarshalXML(e *xml.Encoder,start xml.StartElement) error  {
	if val.out != nil{
		return e.Encode(val.out)
	}
	switch token := val.token.(type) {
	case xml.StartElement:
		if err := e.EncodeToken(token);err != nil {
			return err
		}
		for _,child := range val.children {
			if err := child.MarshalXML(e,xml.StartElement{}); err != nil {
				return err
			}
		}
		return e.EncodeToken(token.End())
	case xml.EndElement:
		panic("Unexpected End Element")
	default:
		return e.EncodeToken(token)
	}
}

type originXMLValueReader struct {
	val *OriginXMLValue
	start,end bool
	child int
	childReader xml.TokenReader
}

//TokenReader
func (val *OriginXMLValue) TokenReader() xml.TokenReader {
	if val.out != nil {
		panic("webdav:failed when call OriginXMLValue.TokenReader on a marshal-only XML")
	}
	return originXMLValueReader{val:val}.childReader//需要测试
}

func (val *OriginXMLValue) Decode(v interface{}) error {
	return xml.NewTokenDecoder(val.TokenReader()).Decode(&v)
}

func (val *OriginXMLValue) XMLName() (name xml.Name,ok bool)  {
	if start,ok := val.token.(xml.StartElement);ok {
		return start.Name,true
	}
	return xml.Name{},false
}

func (reader *originXMLValueReader) Token() (xml.Token,error)  {
	if reader.end{
		return nil,io.EOF
	}

	start,ok := reader.val.token.(xml.StartElement)
	if !ok {
		reader.end = true
		return reader.val.token,nil
	}

	if !reader.start{
		reader.start = true
		return start,nil
	}

	for reader.child < len(reader.val.children){
		if reader.childReader == nil {
			reader.childReader = reader.val.children[reader.child].TokenReader()
		}
		token,err := reader.childReader.Token()
		if err == io.EOF{
			reader.childReader = nil
			reader.child++
		}else {
				return token,err
		}
	}
	reader.end = true
	return start.End(),nil
}

func valueXMLName(v interface{}) (xml.Name,error)  {
	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Ptr{
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct{
		return xml.Name{},fmt.Errorf("webdav:%T is not a struct",v)
	}
	nameField,ok := t.FieldByName("XMLName")
	if !ok{
		return xml.Name{},fmt.Errorf("webdav: %T is missing an XMLName struct field",v)
	}
	if nameField.Type != reflect.TypeOf(xml.Name{}){
		return xml.Name{},fmt.Errorf("webdav:%T.XMLName is not an xml.Name",v)
	}
	tag := nameField.Tag.Get("xml")
	if tag == ""{
		return xml.Name{},fmt.Errorf(`webdav:%T.XMLName is missing an "xml" tag`,v)
	}
	name := strings.Split(tag,",")[0]
	nameParts := strings.Split(name," ")
	if len(nameParts) != 2 {
		return xml.Name{},fmt.Errorf("webdav: expected a namespace and local name in %T.XMLName's xml tag",v)
	}
	return xml.Name{nameParts[0],nameParts[1]},nil
}
