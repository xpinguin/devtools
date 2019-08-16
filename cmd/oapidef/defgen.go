package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/davecgh/go-spew/spew"

	"github.com/gertd/go-pluralize"
	"github.com/hjson/hjson-go"
	"gopkg.in/yaml.v3"
)

type (
	_ = yaml.Encoder
)

type (
	JsonObj  = map[string]interface{}
	JsonList = []interface{}
)

var (
	plz *pluralize.Client

	go2jsonTypeNames = map[string]string{
		"float": "number",
		"int":   "integer",
		"bool":  "boolean",
	}
)

func normalizeTypeName(t string) (n string) {
	for goType, jsonType := range go2jsonTypeNames {
		if strings.HasPrefix(t, goType) {
			return jsonType
		}
	}
	return t
}

func normalizeModelName(m string) string {
	return strcase.ToCamel(m)
}

func oapiDefinition(name string, obj JsonObj) JsonObj {
	name = normalizeModelName(name)

	def := make(JsonObj)
	props := make(JsonObj)
	def["type"] = "object"
	def["properties"] = props

	for fldName, fld := range obj {
		fldName = normalizeModelName(fldName)
		///
		var propDef JsonObj
		///
		switch x := fld.(type) {
		case JsonObj:
			for _, prop := range oapiDefinition(fldName, x) {
				propDef = prop.(JsonObj)
				break
			}
		case JsonList:
			itemsDef := make(JsonObj)
			for _, v := range x {
				if _, ok := v.(JsonList); ok {
					fmt.Fprintln(os.Stderr, "IDKWTF")
				} else if obj, ok := v.(JsonObj); ok {
					for _, item := range oapiDefinition(plz.Singular(fldName), obj) {
						itemsDef = item.(JsonObj)
						break
					}
				} else if v != nil {
					itemsDef = JsonObj{
						"type": normalizeTypeName(reflect.TypeOf(v).Name()),
					}
				}
				break
			}
			if itemsDef == nil {
				continue
			}
			propDef = JsonObj{
				"type":  "array",
				"items": itemsDef,
			}
		default:
			if fld == nil {
				fmt.Fprintln(os.Stderr, "WARN:", "nil field:", fldName)
				continue
			}
			propDef = JsonObj{
				"type": normalizeTypeName(reflect.TypeOf(fld).Name()),
			}
		}
		props[fldName] = propDef
	}

	return JsonObj{name: def}
}

func init() {
	plz = pluralize.NewClient()
}

func main() {
	var dataSamplePath string
	///
	flag.StringVar(&dataSamplePath, "sample", "", "Path to the sample API response")
	flag.Parse()
	///
	switch "" {
	case dataSamplePath:
		fmt.Fprintln(os.Stderr, "Assuming STDIN for the sample response; see help")
	}

	///
	dataSampleRaw, err := ioutil.ReadFile(dataSamplePath)
	if err != nil {
		fmt.Println("ERR:", err)
		return
	}

	///
	var dataSample map[string]interface{}
	if err := hjson.Unmarshal(dataSampleRaw, &dataSample); err != nil {
		fmt.Println("Failed to parse sample data:", err)
		return
	}
	spew.Fdump(os.Stderr, dataSample)

	///
	flightSample := dataSample["data"].(JsonObj)["flights"].(JsonList)[0].(JsonObj)
	def := oapiDefinition("flight", flightSample)
	{
		enc := yaml.NewEncoder(os.Stdout)
		defer enc.Close()
		enc.SetIndent(2)
		if err := enc.Encode(JsonObj{"definitions": def}); err != nil {
			fmt.Println("Failed to marshal sample data into YAML:", err)
			return
		}
	}
}
