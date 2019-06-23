/*compound numeral (\in N{0..sup}:N{1..sup}) to explicitly label: data types+instances (eg. map,json,struct,term,tree(_AST_,XML,...)) and computation (eg. SSA-instruction, _AST_, principal term)"
[master 7316826] init[nocode,task,todo]: compound numeral (\in N{0..sup}:N{1..sup}) to explicitly label: data types+instances (eg. map,json,struct,term,tree(_AST_,XML,...)) and computation (eg. SSA-instruction, _AST_, principal term)*/
package numeral

import (
	"regexp"
	"strings"
)

//////////////////////
//////////////////////
var whitespacePfxLn = regexp.MustCompile(
	`^(\s*)` +
		`(?:\S+)?(?:.+)?$`,
)

var srcStructLn = regexp.MustCompile(
	`^(\s+)` +
		`((?:[A-Za-z_])[\w\d_]*)` +
		`(\s+)` +
		`((?:\*|(?:\[\]))?[A-Za-z_0-9]+(?:\s*[{][^/]+)?)` +
		`((?:\s+)?/[*/].*)?$`,
)

//////////////////////
func stripStructSrc(structSrc string) (src string, indentCnt map[string]int) {
	structSrc = strings.TrimSpace(structSrc)
	structSrc = strings.ReplaceAll(structSrc, "\r\n", "\n")
	structLines := strings.Split(structSrc, "\n")

	indentCnt = map[string]int{}
	for _, line := range structLines {
		stripped := srcStructLn.ReplaceAllString(line, "$1$4")
		src += stripped + "\n"

		indent := whitespacePfxLn.ReplaceAllString(stripped, "$1")
		indentCnt[indent] += 1

	}

	return src, indentCnt
}
