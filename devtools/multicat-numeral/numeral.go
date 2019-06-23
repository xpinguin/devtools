/*compound numeral (\in N{0..sup}:N{1..sup}) to explicitly label: data types+instances (eg. map,json,struct,term,tree(_AST_,XML,...)) and computation (eg. SSA-instruction, _AST_, principal term)"
[master 7316826] init[nocode,task,todo]: compound numeral (\in N{0..sup}:N{1..sup}) to explicitly label: data types+instances (eg. map,json,struct,term,tree(_AST_,XML,...)) and computation (eg. SSA-instruction, _AST_, principal term)*/
package numeral

import (
	"fmt"
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
func StripStructSrc(structSrc string) (src string, indentCnt map[string]int) {
	structSrc = strings.TrimSpace(structSrc)
	structSrc = strings.ReplaceAll(structSrc, "\r\n", "\n")
	structLines := strings.Split(structSrc, "\n")

	indentCnt = map[string]int{}
	for _, line := range structLines {
		stripped := srcStructLn.ReplaceAllString(line, "$1$4")
		src += strings.TrimRight(stripped, " }\t{") + "\n"

		indent := whitespacePfxLn.ReplaceAllString(stripped, "$1")
		indentCnt[indent] += 1

	}

	return src, indentCnt
}

//////////////////////
/// ?HINT: numerator/denominator <=> rational
/////////////////////
func NumerateStructSrc(structSrc string) (nsrc string) {
	structSrc, indents := StripStructSrc(structSrc)
	structLines := strings.Split(structSrc, "\n")

	/// FIXME
	lineTypeOffs := 0
	if strings.HasPrefix(strings.TrimSpace(structLines[0]), "type") {
		lineTypeOffs = 1
	}
	////
	delete(indents, "") // module-level
	////
	for pfx, c := range indents {
		if strings.ContainsAny(pfx, " ") { // NB. space and \t could repr different cardinality, eg ' ' = \t div 4 (or 2, or 8, you see, don't you?)
			nsrc += fmt.Sprintf("//!!! XXX: NOT IMPLEMENTED: Rational sub-cardinality: '%s'. Opaque cardinality: %d\n", pfx, c)
		}
	}
	////
	if nsrc != "" {
		nsrc += "\n\n"
	}
	////

	// or 'denominator' instead???
	var lastDepth int
	depthFieldPos := map[int]uint{}
	/// stupid pre-init, just to be sure
	for d := 0; d <= len(indents); d++ {
		depthFieldPos[d] = 1
	}
	///

	for numerator, line := range structLines[lineTypeOffs:] {
		var nline string

		for pfx, cardinality := range indents {
			if constrs := strings.TrimPrefix(line, pfx); constrs == strings.TrimLeft(line, "\t") {
				depth := len(pfx)
				if depth < lastDepth {
					depthFieldPos[lastDepth] = 1
				}
				if depth != lastDepth {
					lastDepth = depth
				}
				/////
				nline = fmt.Sprintf("%d:%d\t//%d\t|%s\t**%d", len(pfx), depthFieldPos[depth], numerator+1, constrs, cardinality)
				/////
				depthFieldPos[depth] += 1
				break
			}
		}

		if nline != line {
			nsrc += nline + "\n"
		} else if line != "" {
			nsrc += "%% " + nline + "\n"
		}
	}
	return nsrc
}
