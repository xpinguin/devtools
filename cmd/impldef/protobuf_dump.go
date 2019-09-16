package main

import (
	"crypto/rand"
	"fmt"
	"reflect"
)

//////////////////////////
type genFunc func(x interface{}) string

type genHandlers map[reflect.Type]genFunc

type genResults map[interface{}]string

///////////////////////////
type swaggerToPB struct{}

func (h *swaggerToPB) OpenDef(def *OpenDef) string {
	// TODO
}

func (h *swaggerToPB) Schema(s *Schema) string {
	// TODO
}

func (h *swaggerToPB) Param(s *Param) string {
	// TODO
}

func (h *swaggerToPB) Response(s *Response) string {
	// TODO
}

func (h *swaggerToPB) MethodInfo(info *MethodInfo) string {
	// TODO
}

///////////////////////////
type pbGen struct {
	handlers  genHandlers
	generated genResults
}

func NewProtobufGen(hnds genHandlers) *pbGen {
	return &pbGen{handlers: genHandlers, generated: genResults{}}
}

func (st *pbGen) Dump(def *OpenDef) (lines []string) {
	declPfx := "" // TODO pfx := []string{}
_Again:
	switch s.Type {
	case "", "$ref":
		if s.s != nil {
			s = s.s
			s.Type = "$ref"
			goto _Again
		}

	case "array":
		if s.Items == nil {
			panic("empty items of an array")
		}
		s = s.Items
		declPfx = "repeated " + declPfx
		goto _Again

	case "object":
		if s.Properties == nil {
			panic("empty properties of an object")
		}
		for name, ss := range s.Properties {
			//fmt.Println("message", name, ss.ProtobufSrc(name))
		}

	default:
		pbs = append(pbs, fmt.Sprintf("%s %s = %d;",
			protobufType(s.Type), name, rand.Int()%100)) // FIXME
	}
	return pbs
}

/////////////////
func (d *OpenDef) ProtobufSrc() (pbs []string) {
	/// Data model
	for sname, s := range d.Definitions {
		pbs = append(pbs, s.ProtobufSrc(sname)...)
	}

	for _, m := range d.PathsMethods {
		for _, info := range []*MethodInfo{m.Get, m.Post} {
			if info == nil {
				continue
			}
			/// Arguments
			for _, arg := range info.Args {
				pbs = append(pbs, arg.ProtobufSrc()...)
			}
			/// Return types
			for _, ret := range info.CodeRets {
				s := ret.Schema
				if s == nil {
					pbs = append(pbs, fmt.Sprintf("// ERR: nil Schema in: %v", ret))
					continue
				}
				if _, err := s.LinkRef(d); err != nil {
					pbs = append(pbs, fmt.Sprintf("// ERR: %v", err))
					continue
				}
			}
			for code, ret := range info.CodeRets {
				if s := ret.Schema; s != nil {
					pbs = append(pbs,
						s.ProtobufSrc(fmt.Sprintf("%sResponse_%s", info.OperationId, code))...)
				}
			}
		}
	}
	return
}

///////////////////

func protobufType(oapiType string) string {
	switch oapiType {
	case "number":
		return "float"
	case "integer":
		return "int"
	case "array":
		return "repeated"
	case "object":
		return "message"
	case "boolean":
		return "bool"
	}
	return oapiType
}

/***************************************************************
		X          					{"type": "object",
			     										"properties": {
		X.code     					{"type": "integer"},
		X.data:    					{"type": "object",
			            				    			"properties": {
        X.data.flights: 			{"type": "array",
		X.data.flights.[]
        data.flights.[]Flig	ht		   				    "items": {"$ref": "#/definitions/Flight"}

		////////////////////////////////////////////////////////////////////
		X  						    {"type": "object",
				  			  	    					"properties": {
        X.Documents 				{"type": "object",
				  						               "properties": {
        X.Documents.Adt 			{"type": "object",
				  			  	            			"properties": {
	    X.Documents.Adt.Other 		{"type": "array",
	 // X.Documents.Adt.Other.[]
		X.Documents.Adt.Other.[]string	                "items": {"type": "string"}
*****************************************************************/
