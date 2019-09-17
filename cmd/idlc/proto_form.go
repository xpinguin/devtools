// proto_form.go: the syntactic aspect of the protobuf IDL
package main

import (
	"fmt"
	"strconv"
	"strings"
)

//go:generate go get -v golang.org/x/tools/cmd/stringer

///////////////////////////
///////////////////////////
//go:generate stringer -type=NameAfx
type NameAfx int

const (
	NULL NameAfx = iota

	REPEATED NameAfx = 1 << iota
	ONE_OF
	DIRECTIVE
	SERVICE
	MESSAGE
	ENUM
	FIELD
	VARIANT
	RPC

	METHOD = RPC
	OPTION = DIRECTIVE
)

///////////////////////////
//go:generate stringer -type=TypeAfx
type TypeAfx int

const (
	NONE TypeAfx = iota
	BOOL TypeAfx = 1 << iota
	STRING
	BYTES
	INT32
	INT64
	INT
	FLOAT
	PB_ANY
	PB_TIMESTAMP
	PB_DURATION
)

///////////////////////////
const hiddenTypes = NONE /* ... | INT | PB_ANY | ... */
const hiddenFlags = NULL | DIRECTIVE | FIELD | VARIANT

///////////////////////////
/** REM it is somewhat helpful to think about "terminal" & "non-terminal" fields */
type Field struct {
	flags NameAfx
	typ   TypeAfx
	name  string
	tag   int
}

///////////////////////////
type Message struct {
	/*head*/ *Field // NB. untagged field semantically
	/*body*/ ents []*Field
}

///////////////////////////
type Method struct {
	/*label*/ *Field
	arg, ret *Message
}

type Service struct {
	/*label*/ *Field
	meths []*Method
}

///////////////////////////
///////////////////////////
type ProtoPkg struct {
	name, qualName string

	svcs []*Service
	msgs []*Message

	opts map[string]string
}

///////////////////////////
///////////////////////////
func (f *Field) Name() (res string) {
	flags := f.flags.String()
	////
	if f.typ == NONE {
		return strings.Join([]string{flags, f.name}, " ")
	}
	////
	typ := f.typ.String()
	return strings.Join([]string{flags, typ, f.name}, " ")
}

func (msg *Message) Name() string {
	return msg.Field.name
}

///////////////////////////
///////////////////////////
func (f *Field) String() string {
	parts := []string{}
	////
	flags := f.flags & ^hiddenFlags
	typ := f.typ & ^hiddenTypes
	////
	if flags != NameAfx(0) {
		for i := 1; i <= int(flags); i = i << 1 {
			if flags&NameAfx(i) != 0 {
				parts = append(parts, strings.ToLower(NameAfx(i).String()))
			}
		}
	}
	if typ != TypeAfx(0) {
		for i := 1; i <= int(typ); i = i << 1 {
			if typ&TypeAfx(i) != 0 {
				parts = append(parts, strings.ToLower(TypeAfx(i).String()))
			}
		}
	}
	parts = append(parts, f.name)
	////
	if f.flags&FIELD != 0 {
		parts = append(parts, "=", strconv.FormatInt(int64(f.tag), 10))
	}
	return strings.Join(parts, " ")
}

func (meth *Method) String() (res string) {
	res = meth.Field.String()
	if meth.arg != nil {
		res += fmt.Sprintf("(%s)", meth.arg.Name())
	}
	if meth.ret != nil {
		res += fmt.Sprintf(" returns (%s)", meth.ret.Name())
	}
	return res
}

func (svc *Service) String() (res string) {
	res = svc.Field.String() + " {\n"
	for _, m := range svc.meths {
		res += "    " + m.String() + ";\n"
	}
	return res + "}\n"
}

func (msg *Message) String() (res string) {
	res = msg.Field.String() + " {\n"
	for i, ent := range msg.ents {
		if ent.tag == 0 {
			ent.tag = i*10 + ent.tag%10
		}
		res += "    " + ent.String() + ";\n"
	}
	return res + "}\n"
}

///////////////////////////
/////////////////////////// TEST CODE
func main() {
	//////////////
	msgs := []*Message{
		&Message{
			Field: &Field{name: "SearchByRegionRequest", flags: MESSAGE},
			ents: []*Field{
				&Field{typ: INT32, name: "region_id", tag: 1, flags: FIELD},
				&Field{name: "SearchRequest q", tag: 2, flags: FIELD},
			}},

		&Message{
			Field: &Field{name: "SearchByHotelsRequest", flags: MESSAGE},
			ents: []*Field{
				&Field{typ: STRING, name: "hotels_ids", tag: 1, flags: FIELD | REPEATED},
				&Field{name: "SearchRequest q", tag: 2, flags: FIELD},
			}},

		&Message{
			Field: &Field{name: "SearchResponse", flags: MESSAGE},
			ents: []*Field{
				&Field{name: "Component components", tag: 1, flags: FIELD | REPEATED},
				&Field{name: "Accommodation components_metadata", tag: 2, flags: FIELD | REPEATED},
				&Field{typ: STRING, name: "error", tag: 10, flags: FIELD},
			}},
	}

	//////
	SearchByRegion := &Method{
		Field: &Field{name: "SearchByRegion", flags: METHOD},
		arg:   msgs[0],
		ret:   &Message{Field: &Field{name: "SearchResponse", flags: MESSAGE}}, // NB. that's won't work...
	}
	SearchByHotels := &Method{
		Field: &Field{name: "SearchByHotels", flags: METHOD},
		arg:   msgs[1],
		ret:   &Message{Field: &Field{name: "SearchResponse", flags: MESSAGE}},
	}

	svc1 := &Service{
		Field: &Field{name: "Accommodations", flags: SERVICE},
		meths: []*Method{SearchByRegion, SearchByHotels},
	}
	////
	fmt.Printf("-----\n%v\n---------\n", svc1.String())

	for _, m := range msgs {
		fmt.Println(m.String())
	}
}
