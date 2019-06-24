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
		/////// HACK
		var indent string
		if typeLine := strings.TrimPrefix(line, "type"); typeLine != line {
			line = " " + strings.TrimLeft(typeLine, " \t")
			indent = "toplevel"
		}
		////////

		stripped := srcStructLn.ReplaceAllString(line, "$1$4")
		if strings.HasSuffix(stripped, "}") {
			continue
		}

		////// HACK
		switch indent {
		case "toplevel":
			indent = ""
			stripped = stripped[1:]
		default:
			indent = whitespacePfxLn.ReplaceAllString(stripped, "$1")
		}
		//////////

		src += strings.TrimRight(stripped, " }\t{") + "\n"
		indentCnt[indent] += 1

	}

	return strings.TrimSpace(src), indentCnt
}

//////////////////////
/// ?HINT: numerator/denominator <=> rational
/////////////////////
func NumerateStructSrc(structSrc string) (nsrc string) {
	structSrc, indents := StripStructSrc(structSrc)
	structLines := strings.Split(structSrc, "\n")

	/// FIXME
	lineTypeOffs := 0
	/*if strings.HasPrefix(strings.TrimSpace(structLines[0]), "type") {
		lineTypeOffs = 1
	}*/
	////
	//delete(indents, "") // module-level*/
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
	var lastDepth int = 0

	depthFieldPos := map[int]uint{}
	globalDepthFieldPos := map[int]uint{}
	var globalCard int
	for pfx, card := range indents {
		globalCard += card

		d := len(pfx)
		depthFieldPos[d] = 1
		globalDepthFieldPos[d] = 1
	}

	for numerator, line := range structLines[lineTypeOffs:] {
		var nline string

		for pfx, depthCard := range indents {
			if constrs := strings.TrimPrefix(line, pfx); constrs == strings.TrimLeft(line, "\t") {
				depth := len(pfx)
				posCardStrs := []string{}
				for idt := ""; len(idt) < depth; idt += "\t" {
					if card, ok := indents[idt]; ok {
						posCardStrs = append(posCardStrs,
							fmt.Sprintf("%d/%d", globalDepthFieldPos[len(idt)]-1, card))
					}
				}
				posCardStrs = append(posCardStrs, fmt.Sprintf("%d/%d", globalDepthFieldPos[depth], depthCard))
				for dd := lastDepth; dd > depth; dd-- {
					depthFieldPos[dd] = 1
				}
				if depth != lastDepth {
					lastDepth = depth
				}
				/////

				posCards := strings.Join(posCardStrs, " + ")
				nline = fmt.Sprintf("%d:(%s)/%d==%d/%d\t\t\t|%s", depth, posCards, globalCard, numerator+1, globalCard, constrs)
				/////
				depthFieldPos[depth] += 1
				globalDepthFieldPos[depth] += 1
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
