package numeral

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

//////////////////////
//////////////////////
const goSrc2 = `struct
	string
	int32
	int32
	int32
	int32
	int32
	int32
	string
	*struct
		int32
		string
		*struct
			int32
			string
			string
			string
			string
			[]int32
		*struct
			string
			string
	*struct
		bool
		bool
		bool
		bool
		bool
		struct
	[]uint8
`

//////////////////////
//////////////////////
func TestStripGosrc(t *testing.T) {
	src, indents := StripStructSrc(goSrc2)
	fmt.Print(src, "\n\n///////////////\n")
	spew.Dump(indents)
}

//////////////////////
func TestStructNumbering(t *testing.T) {
	nsrc := NumerateStructSrc(goSrc2)
	fmt.Println(nsrc)
}
