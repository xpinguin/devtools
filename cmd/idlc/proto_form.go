// proto_form.go: the syntactic aspect of the protobuf IDL
package main

import (
	"fmt"
)

//go:generate go get -v golang.org/x/tools/cmd/stringer

///////////////////////////
///////////////////////////
type Document struct {
	name string
	lns  int
}

///////////////////////////
//go:generate stringer -type=NameAfx
type NameAfx int

const (
	NULL NameAfx = iota

	REPEATED NameAfx = 1 << iota
	ONE_OF

	ENUM
	VARIANT

	DIRECTIVE
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
type Head struct {
	flags NameAfx
	name  string
}

type Field struct {
	*Head
	typ TypeAfx
	tag int
}

type Block struct {
	*Head
	fields []Field
}

///////////////////////////
type Message struct {
	*Block
}

type Service struct {
	*Block
}

type Enum struct {
	*Block
}

///// TEST CODE
func main() {
	fmt.Println(NameAfx(2), TypeAfx(4))
}
