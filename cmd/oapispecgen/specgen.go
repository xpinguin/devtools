package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/davecgh/go-spew/spew"

	"github.com/hjson/hjson-go"
	"gopkg.in/yaml.v3"
)

type (
	_ = yaml.Encoder
)

func main() {
	var specPath, dataSamplePath string
	///
	flag.StringVar(&specPath, "spec", "", "OpenAPI/Swagger YAML path")
	flag.StringVar(&dataSamplePath, "sample", "", "Path to the sample API response")
	flag.Parse()
	///
	switch "" {
	case specPath:
		fmt.Println("ERR:", "No OpenAPI/Swagger YAML path has been specified")
		flag.PrintDefaults()
		return
	case dataSamplePath:
		fmt.Println("Assuming STDIN for the sample response; see help")

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
	{
		cfg := spew.NewDefaultConfig()
		cfg.DisableCapacities = true
		cfg.DisablePointerAddresses = true
		cfg.SortKeys = true
		cfg.SpewKeys = true
		///
		cfg.Dump(dataSample)
	}

	//fmt.Println(specPath, ":", dataSamplePath)
}
