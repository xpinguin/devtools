// dot_dump.go
package main

import (
	"fmt"
	"strings"
)

func objRec(props map[string]*Schema, requiredNames []string) string {
	propStrs := []string{}
	required := map[string]struct{}{}
	for _, n := range requiredNames {
		required[n] = struct{}{}
	}
	for name, _ := range props { // FIXME: deep recursion
		if _, ok := required[name]; ok {
			name += "*"
		}
		propStrs = append(propStrs, name)
	}
	return fmt.Sprintf(`[shape=Mrecord label="%s"];`, strings.Join(propStrs, "|"))
}

func DotDump(def *OpenDef) []byte {
	dot := &strings.Builder{}
	dot.WriteString("digraph G {\n")
	//////
	for name, schema := range def.Definitions {
		switch schema.Type {
		case "array":
			// returns plain label
		case "object":
			dot.WriteRune('\t')
			dot.WriteString(name)
			dot.WriteRune(' ')
			dot.WriteString(objRec(schema.Properties, schema.Required))
			dot.WriteRune('\n')

		default:
			fmt.Println("UNK:", schema.Type)
		}
	}
	//////
	dot.WriteString("}")
	return []byte(dot.String())
}
