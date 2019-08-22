package main

import (
	"fmt"
	"strconv"
)

type IntSliceFlag []int

func (is *IntSliceFlag) String() string {
	return fmt.Sprintf("%#v", is)
}

func (is *IntSliceFlag) Set(v string) error {
	i, err := strconv.Atoi(v)
	if err != nil {
		return err
	}
	*is = append(*is, i)
	return nil
}
