package main

import (
	"flag"
	"fmt"

	"github.com/json5/json5-go"
	"gopkg.in/yaml.v3"
)

type (
	_ = yaml.Encoder
	_ = json5.Unmarshaler
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

	fmt.Println(specPath, ":", dataSamplePath)
}
