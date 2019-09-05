package main

import (
	"fmt"
)

type Definition struct {
	*Schema
	Name string
}

func sqlCreate(tbls chan *Definition) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		////
		for t := range tbls {
			switch t.Type {
			case "object":
				go func() {
					for stmt := range sqlCreate(tbls) {
						out <- stmt
					}
				}()
				////////
				for name, tbl := range t.Properties {
					tbls <- &Definition{Schema: tbl, Name: name}
				}

			case "array", "list":

			default:
				out <- fmt.Sprint(t.Name, t.Type)
			}
		}
	}()
	return out
}

func GoSrcDump(def *OpenDef) []byte {
	return nil
}
