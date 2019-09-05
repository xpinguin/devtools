package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/davecgh/go-spew/spew"

	yaml "gopkg.in/yaml.v3"
)

type (
	_ = yaml.Encoder
)

type OpenDef struct {
	DefVersion string  `json:"swagger" yaml:"swagger,omitempty"`
	Info       APIInfo `json:"info" yaml:"info,omitempty"`
	Host       string  `json:"host" yaml:"host,omitempty"`
	BasePath   string  `json:"basePath" yaml:"basePath,omitempty"`

	PathsMethods map[string]HTTPMethod `json:"paths" yaml:"paths,omitempty"`
	Definitions  map[string]*Schema    `json:"definitions,omitempty" yaml:"definitions,omitempty"`
}

type APIInfo struct {
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	Title   string `json:"title,omitempty" yaml:"title,omitempty"`
}

type Paths map[string]HTTPMethod

type HTTPMethod struct {
	Get  *MethodInfo `json:"get,omitempty" yaml:"get,omitempty"`
	Post *MethodInfo `json:"post,omitempty" yaml:"post,omitempty"`
}

type MethodInfo struct {
	OperationId string              `json:"operationId" yaml:"operationId,omitempty"`
	Summary     string              `json:"summary,omitempty" yaml:"summary,omitempty"`
	Args        []Param             `json:"parameters" yaml:"parameters,omitempty"`
	CodeRets    map[string]Response `json:"responses" yaml:"responses,omitempty"`
}

type Response struct {
	Description string  `json:"description" yaml:"description,omitempty"`
	Schema      *Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
}

type Schema struct {
	Ref  string `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	Type string `json:"type,omitempty" yaml:"type,omitempty"`

	// object
	Required   []string           `json:"required,omitempty" yaml:"required,omitempty"`
	Properties map[string]*Schema `json:"properties,omitempty" yaml:"properties,omitempty"`

	// array
	Items *Schema `json:"items,omitempty" yaml:"items,omitempty"`
}

type Param struct {
	Type     string      `json:"type,omitempty" yaml:"type,omitempty"`
	Name     string      `json:"name,omitempty" yaml:"name,omitempty"`
	In       string      `json:"in,omitempty" yaml:"in,omitempty"`
	Min      int         `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	Max      int         `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	Default  interface{} `json:"default,omitempty" yaml:"default,omitempty"`
	Required bool        `json:"required,omitempty" yaml:"required,omitempty"`
	Format   string      `json:"format,omitempty" yaml:"format,omitempty"`

	// array
	Items *Param `json:"items,omitempty" yaml:"items,omitempty"`
}

func readOpenDef(defPath string) (def *OpenDef, err error) {
	defData, err := ioutil.ReadFile(defPath)
	if err != nil {
		return nil, err
	}

	def = &OpenDef{}
	switch ext := strings.ToLower(path.Ext(defPath)); ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(defData, def); err != nil {
			return nil, err
		}
	case ".json", ".js":
		/*if err := hjson.Unmarshal(defData, def); err != nil {
		log.Print("WARN: hjson-go: failed to unmarshal: %v; trying `encoding/json`...", err)*/
		if err := json.Unmarshal(defData, def); err != nil {
			return nil, err
		}
		//}
	default:
		log.Fatalf("ERR: parser not implemented for: %s (%s)", ext, defPath)
	}
	return
}

func main() {
	////////
	var inputDefPath string
	var outFmt string
	////////
	flag.StringVar(&inputDefPath, "in", "", "OpenAPI yaml/json definition to convert")
	flag.StringVar(&outFmt, "out", "json", "Output format: json, yaml, spew, printf, ...")
	flag.Parse()

	////////
	def, err := readOpenDef(inputDefPath)
	if err != nil {
		log.Fatal("ERR:", err)
		return
	}

	////
	spewCfg := spew.NewDefaultConfig()
	spewCfg.SortKeys = true

	var defDump []byte
	err = nil

	switch outFmt {
	case "json":
		defDump, err = json.MarshalIndent(def, "", "  ")
	case "yaml":
		defDump, err = yaml.Marshal(def)
	case "printf":
		defDump = []byte(fmt.Sprintf("%#v\n", def))
	case "spew":
		defDump = []byte(spewCfg.Sdump(def))
	case "go":
		defDump = GoSrcDump(def)
	default:
		err = fmt.Errorf("Unknown output format: %s", outFmt)
	}

	////
	if err != nil {
		log.Fatal("ERR:", err)
		return
	}
	os.Stdout.Write(defDump)
	//os.Stdout.Write([]byte{'\n'})
}
