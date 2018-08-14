//line parser.go.y:2
package parser

import __yyfmt__ "fmt"

//line parser.go.y:2
import (
	"path/filepath"
)

//line parser.go.y:40
type yySymType struct {
	yys   int
	token Token
	expr  *Node
	str   string
}

const TAssert = 57346
const TBreak = 57347
const TContinue = 57348
const TElse = 57349
const TFor = 57350
const TFunc = 57351
const TGoto = 57352
const TIf = 57353
const TNil = 57354
const TReturn = 57355
const TRequire = 57356
const TVar = 57357
const TYield = 57358
const TEqeq = 57359
const TNeq = 57360
const TLsh = 57361
const TRsh = 57362
const TLte = 57363
const TGte = 57364
const TIdent = 57365
const TNumber = 57366
const TString = 57367
const TOr = 57368
const TAnd = 57369
const UNARY = 57370

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"TAssert",
	"TBreak",
	"TContinue",
	"TElse",
	"TFor",
	"TFunc",
	"TGoto",
	"TIf",
	"TNil",
	"TReturn",
	"TRequire",
	"TVar",
	"TYield",
	"TEqeq",
	"TNeq",
	"TLsh",
	"TRsh",
	"TLte",
	"TGte",
	"TIdent",
	"TNumber",
	"TString",
	"'{'",
	"'('",
	"'T'",
	"TOr",
	"TAnd",
	"'|'",
	"'&'",
	"'^'",
	"'>'",
	"'<'",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'%'",
	"UNARY",
	"'~'",
	"'#'",
	"'}'",
	"';'",
	"'='",
	"')'",
	"'!'",
	"'['",
	"']'",
	"':'",
	"'.'",
	"','",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.go.y:392

var typesLookup = map[string]string{
	"nil": "0", "number": "1", "string": "2", "map": "4", "closure": "6", "generic": "8",
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 95,
	45, 20,
	-2, 64,
}

const yyPrivate = 57344

const yyLast = 745

var yyAct = [...]int{

	156, 34, 91, 22, 164, 16, 158, 59, 28, 136,
	135, 44, 45, 42, 40, 21, 19, 161, 163, 4,
	14, 162, 48, 13, 154, 160, 62, 6, 152, 56,
	153, 101, 57, 86, 151, 133, 1, 82, 83, 84,
	85, 60, 92, 178, 177, 110, 170, 102, 55, 176,
	93, 95, 99, 142, 47, 96, 103, 104, 27, 131,
	109, 134, 108, 21, 112, 113, 114, 115, 116, 117,
	118, 119, 120, 121, 122, 123, 124, 125, 126, 127,
	128, 129, 73, 74, 75, 87, 100, 98, 50, 49,
	51, 130, 50, 23, 81, 140, 29, 26, 31, 46,
	22, 174, 146, 147, 106, 145, 150, 24, 30, 43,
	41, 25, 21, 19, 54, 52, 4, 14, 180, 89,
	13, 36, 20, 61, 6, 33, 37, 39, 32, 58,
	24, 5, 38, 50, 25, 105, 155, 15, 157, 139,
	22, 22, 77, 78, 166, 165, 22, 94, 35, 53,
	171, 169, 21, 21, 167, 90, 132, 2, 21, 71,
	72, 73, 74, 75, 71, 72, 73, 74, 75, 22,
	0, 181, 0, 0, 179, 0, 183, 0, 0, 22,
	22, 21, 22, 0, 184, 185, 0, 186, 0, 0,
	0, 21, 21, 0, 21, 69, 70, 77, 78, 68,
	67, 0, 0, 0, 0, 0, 0, 63, 64, 79,
	80, 76, 65, 66, 71, 72, 73, 74, 75, 69,
	70, 77, 78, 68, 67, 0, 0, 0, 148, 149,
	0, 63, 64, 79, 80, 76, 65, 66, 71, 72,
	73, 74, 75, 69, 70, 77, 78, 68, 67, 0,
	0, 0, 0, 175, 0, 63, 64, 79, 80, 76,
	65, 66, 71, 72, 73, 74, 75, 69, 70, 77,
	78, 68, 67, 0, 0, 0, 0, 137, 0, 63,
	64, 79, 80, 76, 65, 66, 71, 72, 73, 74,
	75, 69, 70, 77, 78, 68, 67, 0, 0, 0,
	182, 0, 0, 63, 64, 79, 80, 76, 65, 66,
	71, 72, 73, 74, 75, 69, 70, 77, 78, 68,
	67, 0, 0, 0, 173, 0, 0, 63, 64, 79,
	80, 76, 65, 66, 71, 72, 73, 74, 75, 69,
	70, 77, 78, 68, 67, 144, 0, 0, 0, 0,
	0, 63, 64, 79, 80, 76, 65, 66, 71, 72,
	73, 74, 75, 69, 70, 77, 78, 68, 67, 138,
	0, 0, 0, 0, 0, 63, 64, 79, 80, 76,
	65, 66, 71, 72, 73, 74, 75, 69, 70, 77,
	78, 68, 67, 111, 0, 0, 0, 0, 0, 63,
	64, 79, 80, 76, 65, 66, 71, 72, 73, 74,
	75, 23, 0, 0, 29, 168, 31, 0, 0, 0,
	0, 0, 0, 0, 0, 24, 30, 43, 41, 25,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 36,
	0, 0, 0, 23, 37, 39, 29, 0, 31, 20,
	38, 0, 172, 0, 0, 0, 0, 24, 30, 43,
	41, 25, 0, 0, 0, 0, 0, 0, 0, 23,
	0, 36, 29, 0, 31, 0, 37, 39, 0, 97,
	0, 0, 38, 24, 30, 43, 41, 25, 0, 0,
	0, 0, 0, 0, 0, 23, 0, 36, 29, 0,
	31, 0, 37, 39, 0, 141, 0, 0, 38, 24,
	30, 43, 41, 25, 0, 0, 0, 0, 0, 0,
	0, 23, 0, 36, 29, 0, 31, 0, 37, 39,
	0, 0, 0, 107, 38, 24, 30, 43, 41, 25,
	0, 0, 0, 0, 0, 0, 0, 23, 0, 36,
	29, 0, 31, 0, 37, 39, 88, 0, 0, 0,
	38, 24, 30, 43, 41, 25, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 36, 0, 0, 0, 0,
	37, 39, 0, 10, 8, 9, 38, 17, 23, 0,
	18, 0, 11, 12, 20, 7, 0, 0, 0, 0,
	0, 0, 24, 0, 0, 0, 25, 0, 0, 0,
	0, 0, 0, 69, 70, 77, 78, 68, 67, 0,
	0, 0, 0, 143, 3, 63, 64, 79, 80, 76,
	65, 66, 71, 72, 73, 74, 75, 10, 8, 9,
	0, 17, 0, 0, 18, 0, 11, 12, 20, 7,
	0, 0, 0, 0, 0, 0, 24, 0, 0, 50,
	25, 10, 8, 9, 0, 17, 23, 0, 18, 0,
	11, 12, 20, 7, 0, 0, 0, 0, 159, 0,
	24, 0, 0, 0, 25, 0, 0, 0, 0, 0,
	69, 70, 77, 78, 68, 67, 0, 0, 0, 0,
	0, 0, 3, 64, 79, 80, 76, 65, 66, 71,
	72, 73, 74, 75, 69, 70, 77, 78, 68, 67,
	0, 69, 70, 77, 78, 68, 67, 0, 79, 80,
	76, 65, 66, 71, 72, 73, 74, 75, 65, 66,
	71, 72, 73, 74, 75,
}
var yyPact = [...]int{

	-1000, 657, -1000, -1000, 13, -1000, -1000, 538, -1000, -1000,
	538, 538, 74, -1000, -1000, -1000, 9, 62, 63, 92,
	91, 2, -20, -7, -1000, 538, -1000, -1000, 596, -1000,
	-1000, 69, -1000, -1000, -20, -1000, 538, 538, 538, 538,
	58, 512, -1000, -1000, 596, 596, -1000, -1000, -1000, 434,
	-1000, 538, 58, -22, 1, 538, 84, 81, -1000, 486,
	-1000, -2, 346, 538, 538, 538, 538, 538, 538, 538,
	538, 538, 538, 538, 538, 538, 538, 538, 538, 538,
	538, -1000, -1000, -1000, -1000, -1000, 66, 12, -1000, 17,
	-43, -44, 226, 322, 460, -20, 8, -1000, 579, 298,
	66, 79, 538, 596, 178, 538, -1000, -1000, -19, 596,
	-1000, -1000, 673, 697, 123, 123, 123, 123, 123, 123,
	44, 44, -1000, -1000, -1000, 704, 128, 128, 704, 704,
	-1000, -1000, -23, -1000, -1000, 538, 538, 538, 633, 107,
	370, -1000, -1000, -1000, 633, -1000, 0, 596, -1000, 402,
	274, 538, -1000, 78, -1000, 202, 596, 596, -1000, -1000,
	-1000, 4, -1000, -1000, -1000, -3, -4, 633, -1000, 111,
	538, 250, -1000, -1000, -1000, 538, -1000, 633, 633, -1000,
	633, 596, -1000, 596, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 36, 4, 157, 13, 156, 2, 155, 149, 0,
	148, 1, 5, 25, 21, 147, 139, 18, 6, 17,
	137, 131, 14, 97, 129, 128, 33, 125, 119,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 13, 3, 3, 3, 3, 18,
	18, 18, 18, 18, 18, 21, 21, 21, 12, 12,
	12, 15, 15, 16, 16, 14, 14, 14, 14, 14,
	17, 17, 22, 22, 20, 19, 19, 19, 19, 19,
	19, 19, 19, 4, 4, 4, 4, 4, 4, 5,
	5, 6, 6, 7, 7, 8, 8, 8, 8, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 10, 11,
	11, 11, 11, 23, 24, 24, 25, 26, 26, 27,
	27, 28, 28, 28, 28,
}
var yyR2 = [...]int{

	0, 0, 2, 3, 2, 1, 2, 1, 1, 1,
	1, 2, 1, 1, 1, 1, 1, 1, 2, 3,
	1, 2, 1, 2, 1, 2, 5, 7, 7, 6,
	5, 7, 1, 2, 4, 1, 2, 1, 1, 2,
	1, 2, 2, 1, 4, 6, 5, 5, 3, 1,
	3, 1, 3, 3, 5, 1, 3, 5, 3, 1,
	1, 2, 1, 1, 1, 1, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 2, 2, 2, 2, 1, 1,
	3, 1, 3, 2, 2, 3, 3, 2, 3, 2,
	3, 1, 2, 1, 2,
}
var yyChk = [...]int{

	-1000, -1, -3, 45, -19, -21, -13, 16, 5, 6,
	4, 13, 14, -14, -17, -20, -12, 8, 11, -22,
	15, -4, -11, 9, 23, 27, -23, 45, -9, 12,
	24, 14, -25, -27, -11, -10, 37, 42, 48, 43,
	-22, 26, -4, 25, -9, -9, 25, 45, -2, 27,
	26, 27, 23, -8, 23, 46, 49, 52, -24, 27,
	48, -23, -9, 29, 30, 34, 35, 22, 21, 17,
	18, 36, 37, 38, 39, 40, 33, 19, 20, 31,
	32, 25, -9, -9, -9, -9, -26, 27, 44, -28,
	-7, -6, -9, -9, -15, -11, -12, 45, -1, -9,
	-26, 53, 46, -9, -9, 51, 23, 47, -6, -9,
	47, 47, -9, -9, -9, -9, -9, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, -9, -9, -9, -9,
	-2, 47, -5, 23, 44, 53, 53, 51, 47, -16,
	-9, 45, 45, 44, 47, -2, 23, -9, 50, 51,
	-9, 53, 47, 53, 47, -9, -9, -9, -18, 45,
	-13, -19, -14, -17, -2, -12, -2, 47, 45, -18,
	46, -9, 50, 50, 23, 51, 45, 47, 47, -18,
	7, -9, 50, -9, -18, -18, -18,
}
var yyDef = [...]int{

	1, -2, 2, 5, 0, 7, 8, 35, 37, 38,
	0, 40, 0, 15, 16, 17, 0, 0, 0, 0,
	0, 89, 20, 32, 43, 0, 91, 6, 36, 59,
	60, 0, 62, 63, 64, 65, 0, 0, 0, 0,
	0, 0, 89, 88, 39, 41, 42, 4, 25, 0,
	1, 0, 0, 18, 55, 0, 0, 0, 93, 0,
	33, 91, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 61, 84, 85, 86, 87, 0, 0, 99, 0,
	101, 103, 51, 0, 0, -2, 0, 22, 0, 0,
	0, 0, 0, 19, 0, 0, 48, 94, 0, 51,
	90, 92, 66, 67, 68, 69, 70, 71, 72, 73,
	74, 75, 76, 77, 78, 79, 80, 81, 82, 83,
	96, 97, 0, 49, 100, 102, 104, 0, 0, 0,
	0, 24, 21, 3, 0, 34, 58, 56, 44, 0,
	0, 0, 95, 0, 98, 0, 52, 53, 26, 9,
	10, 0, 12, 13, 14, 0, 0, 0, 23, 30,
	0, 0, 46, 47, 50, 0, 11, 0, 0, 29,
	0, 57, 45, 54, 27, 28, 31,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 48, 3, 43, 3, 40, 32, 3,
	27, 47, 38, 36, 53, 37, 52, 39, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 51, 45,
	35, 46, 34, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 28, 3, 3, 3, 3, 3,
	3, 49, 3, 50, 33, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 26, 31, 44, 42,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 29, 30, 41,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.go.y:70
		{
			yyVAL.expr = NewCompoundNode("chain")
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:76
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, yyDollar[2].expr)
			yyVAL.expr = yyDollar[1].expr
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 3:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:85
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:88
		{
			if yyDollar[1].expr.isIsolatedDupCall() {
				yyDollar[1].expr.Compound[2].Compound[0] = NewNumberNode("0")
			}
			yyVAL.expr = yyDollar[1].expr
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:96
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:97
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:98
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:99
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:102
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:103
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:104
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:105
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:106
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:107
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:110
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:111
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:112
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 18:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:115
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:118
		{
			yyVAL.expr = NewCompoundNode("move", yyDollar[1].expr, yyDollar[3].expr)
			if len(yyDollar[1].expr.Compound) > 0 {
				if c, _ := yyDollar[1].expr.Compound[0].Value.(string); c == "load" {
					yyVAL.expr = NewCompoundNode("store", yyDollar[1].expr.Compound[1], yyDollar[1].expr.Compound[2], yyDollar[3].expr)
				}
			}
			if c, _ := yyDollar[1].expr.Value.(string); c != "" && yyDollar[1].expr.Type == NTAtom {
				if a, b, s := yyDollar[3].expr.isSimpleAddSub(); a == c {
					yyDollar[3].expr.Compound[2].Value = yyDollar[3].expr.Compound[2].Value.(float64) * s
					yyVAL.expr = NewCompoundNode("inc", yyDollar[1].expr, yyDollar[3].expr.Compound[2])
					yyVAL.expr.Compound[1].Pos = yyDollar[1].expr.Pos
				} else if b == c {
					yyDollar[3].expr.Compound[1].Value = yyDollar[3].expr.Compound[1].Value.(float64) * s
					yyVAL.expr = NewCompoundNode("inc", yyDollar[1].expr, yyDollar[3].expr.Compound[1])
					yyVAL.expr.Compound[1].Pos = yyDollar[1].expr.Pos
				}
			}
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 20:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:138
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 21:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:143
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:144
		{
			yyVAL.expr = NewCompoundNode("chain")
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:147
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:148
		{
			yyVAL.expr = NewNumberNode("1")
		}
	case 25:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:151
		{
			yyVAL.expr = NewCompoundNode("for", NewNumberNode("1"), NewCompoundNode(), yyDollar[2].expr).setPos0(yyDollar[1].token.Pos)
		}
	case 26:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:154
		{
			yyVAL.expr = NewCompoundNode("for", yyDollar[3].expr, NewCompoundNode(), yyDollar[5].expr).setPos0(yyDollar[1].token.Pos)
		}
	case 27:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:157
		{
			yyVAL.expr = yyDollar[3].expr
			yyVAL.expr.Compound = append(yyVAL.expr.Compound, NewCompoundNode("for", yyDollar[4].expr, NewCompoundNode("chain", yyDollar[5].expr), yyDollar[7].expr))
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 28:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:162
		{
			yyVAL.expr = yyDollar[3].expr
			yyVAL.expr.Compound = append(yyVAL.expr.Compound, NewCompoundNode("for", yyDollar[4].expr, yyDollar[5].expr, yyDollar[7].expr))
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 29:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:167
		{
			yyVAL.expr = yyDollar[3].expr
			yyVAL.expr.Compound = append(yyVAL.expr.Compound, NewCompoundNode("for", yyDollar[4].expr, NewCompoundNode(), yyDollar[6].expr))
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 30:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:174
		{
			yyVAL.expr = NewCompoundNode("if", yyDollar[3].expr, yyDollar[5].expr, NewCompoundNode())
		}
	case 31:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:175
		{
			yyVAL.expr = NewCompoundNode("if", yyDollar[3].expr, yyDollar[5].expr, yyDollar[7].expr)
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:178
		{
			yyVAL.str = "func"
		}
	case 33:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:179
		{
			yyVAL.str = "safefunc"
		}
	case 34:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:182
		{
			funcname := NewAtomNode(yyDollar[2].token)
			yyVAL.expr = NewCompoundNode(
				"chain",
				NewCompoundNode("set", funcname, NewNilNode()),
				NewCompoundNode("move", funcname, NewCompoundNode(yyDollar[1].str, funcname, yyDollar[3].expr, yyDollar[4].expr)))
			yyVAL.expr.Compound[1].Compound[0].Pos = yyDollar[2].token.Pos
			yyVAL.expr.Compound[2].Compound[0].Pos = yyDollar[2].token.Pos
			yyVAL.expr.Compound[2].Compound[2].Compound[0].Pos = yyDollar[2].token.Pos
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:194
		{
			yyVAL.expr = NewCompoundNode("yield").setPos0(yyDollar[1].token.Pos)
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:195
		{
			yyVAL.expr = NewCompoundNode("yield", yyDollar[2].expr).setPos0(yyDollar[1].token.Pos)
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:196
		{
			yyVAL.expr = NewCompoundNode("break").setPos0(yyDollar[1].token.Pos)
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:197
		{
			yyVAL.expr = NewCompoundNode("continue").setPos0(yyDollar[1].token.Pos)
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:198
		{
			yyVAL.expr = NewCompoundNode("assert", yyDollar[2].expr).setPos0(yyDollar[1].token.Pos)
		}
	case 40:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:199
		{
			yyVAL.expr = NewCompoundNode("ret").setPos0(yyDollar[1].token.Pos)
		}
	case 41:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:200
		{
			if yyDollar[2].expr.isIsolatedDupCall() {
				if h, _ := yyDollar[2].expr.Compound[2].Compound[2].Value.(float64); h == 1 {
					yyDollar[2].expr.Compound[2].Compound[2] = NewNumberNode("2")
				}
			}
			yyVAL.expr = NewCompoundNode("ret", yyDollar[2].expr).setPos0(yyDollar[1].token.Pos)
		}
	case 42:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:208
		{
			path := filepath.Join(filepath.Dir(yyDollar[1].token.Pos.Source), yyDollar[2].token.Str)
			yyVAL.expr = yylex.(*Lexer).loadFile(path)
		}
	case 43:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:214
		{
			yyVAL.expr = NewAtomNode(yyDollar[1].token).setPos(yyDollar[1].token.Pos)
		}
	case 44:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:215
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos).setPos(yyDollar[1].expr.Pos)
		}
	case 45:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:216
		{
			yyVAL.expr = NewCompoundNode("slice", yyDollar[1].expr, yyDollar[3].expr, yyDollar[5].expr).setPos0(yyDollar[1].expr.Pos).setPos(yyDollar[1].expr.Pos)
		}
	case 46:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:217
		{
			yyVAL.expr = NewCompoundNode("slice", yyDollar[1].expr, yyDollar[3].expr, NewNumberNode("-1")).setPos0(yyDollar[1].expr.Pos).setPos(yyDollar[1].expr.Pos)
		}
	case 47:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:218
		{
			yyVAL.expr = NewCompoundNode("slice", yyDollar[1].expr, NewNumberNode("0"), yyDollar[4].expr).setPos0(yyDollar[1].expr.Pos).setPos(yyDollar[1].expr.Pos)
		}
	case 48:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:219
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, NewStringNode(yyDollar[3].token.Str)).setPos0(yyDollar[1].expr.Pos).setPos(yyDollar[1].expr.Pos)
		}
	case 49:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:222
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].token.Str)
		}
	case 50:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:225
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, NewAtomNode(yyDollar[3].token))
			yyVAL.expr = yyDollar[1].expr
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:231
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].expr)
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:234
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, yyDollar[3].expr)
			yyVAL.expr = yyDollar[1].expr
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:240
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 54:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:243
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, yyDollar[3].expr, yyDollar[5].expr)
			yyVAL.expr = yyDollar[1].expr
		}
	case 55:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:249
		{
			yyVAL.expr = NewCompoundNode("chain", NewCompoundNode("set", NewAtomNode(yyDollar[1].token), NewNilNode()))
			yyVAL.expr.Compound[1].Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:253
		{
			yyVAL.expr = NewCompoundNode("chain", NewCompoundNode("set", NewAtomNode(yyDollar[1].token), yyDollar[3].expr))
			yyVAL.expr.Compound[1].Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 57:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:257
		{
			x := NewCompoundNode("set", NewAtomNode(yyDollar[3].token), yyDollar[5].expr).setPos0(yyDollar[1].expr.Pos)
			yyDollar[1].expr.Compound = append(yyVAL.expr.Compound, x)
			yyVAL.expr = yyDollar[1].expr
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:262
		{
			x := NewCompoundNode("set", NewAtomNode(yyDollar[3].token), NewNilNode()).setPos0(yyDollar[1].expr.Pos)
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, x)
			yyVAL.expr = yyDollar[1].expr
		}
	case 59:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:269
		{
			yyVAL.expr = NewNilNode()
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 60:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:273
		{
			yyVAL.expr = NewNumberNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 61:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:277
		{
			path := filepath.Join(filepath.Dir(yyDollar[1].token.Pos.Source), yyDollar[2].token.Str)
			yyVAL.expr = yylex.(*Lexer).loadFile(path)
		}
	case 62:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:281
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 63:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:282
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 64:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:283
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 65:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:284
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:285
		{
			yyVAL.expr = NewCompoundNode("or", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:286
		{
			yyVAL.expr = NewCompoundNode("and", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:287
		{
			yyVAL.expr = NewCompoundNode("<", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:288
		{
			yyVAL.expr = NewCompoundNode("<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 70:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:289
		{
			yyVAL.expr = NewCompoundNode("<=", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:290
		{
			yyVAL.expr = NewCompoundNode("<=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 72:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:291
		{
			yyVAL.expr = NewCompoundNode("==", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:292
		{
			yyVAL.expr = NewCompoundNode("!=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:293
		{
			yyVAL.expr = NewCompoundNode("+", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 75:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:294
		{
			yyVAL.expr = NewCompoundNode("-", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:295
		{
			yyVAL.expr = NewCompoundNode("*", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 77:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:296
		{
			yyVAL.expr = NewCompoundNode("/", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 78:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:297
		{
			yyVAL.expr = NewCompoundNode("%", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 79:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:298
		{
			yyVAL.expr = NewCompoundNode("^", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 80:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:299
		{
			yyVAL.expr = NewCompoundNode("<<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 81:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:300
		{
			yyVAL.expr = NewCompoundNode(">>", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 82:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:301
		{
			yyVAL.expr = NewCompoundNode("|", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 83:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:302
		{
			yyVAL.expr = NewCompoundNode("&", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 84:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:303
		{
			yyVAL.expr = NewCompoundNode("-", NewNumberNode("0"), yyDollar[2].expr).setPos0(yyDollar[2].expr.Pos)
		}
	case 85:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:304
		{
			yyVAL.expr = NewCompoundNode("~", yyDollar[2].expr).setPos0(yyDollar[2].expr.Pos)
		}
	case 86:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:305
		{
			yyVAL.expr = NewCompoundNode("!", yyDollar[2].expr).setPos0(yyDollar[2].expr.Pos)
		}
	case 87:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:306
		{
			yyVAL.expr = NewCompoundNode("#", yyDollar[2].expr).setPos0(yyDollar[2].expr.Pos)
		}
	case 88:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:309
		{
			yyVAL.expr = NewStringNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 89:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:315
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 90:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:316
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 91:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:317
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 92:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:318
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 93:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:321
		{
			switch c, _ := yyDollar[1].expr.Value.(string); c {
			case "copy":
				switch len(yyDollar[2].expr.Compound) {
				case 0:
					yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(NewNumberNode("1"), NewNumberNode("1"), NewNumberNode("1")))
				case 1:
					yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(NewNumberNode("1"), yyDollar[2].expr.Compound[0], NewNumberNode("0")))
				default:
					p := yyDollar[2].expr.Compound[1]
					if p.Type != NTCompound && p.Type != NTAtom {
						yylex.(*Lexer).Error("invalid argument for S")
					}
					yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(NewNumberNode("1"), yyDollar[2].expr.Compound[0], p))
				}
			case "typeof":
				switch len(yyDollar[2].expr.Compound) {
				case 0:
					yylex.(*Lexer).Error("typeof takes at least 1 argument")
				case 1:
					yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].expr.Compound[0], NewNumberNode("255")))
				default:
					x, _ := yyDollar[2].expr.Compound[1].Value.(string)
					if ti, ok := typesLookup[x]; ok {
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].expr.Compound[0], NewNumberNode(ti)))
					} else {
						yylex.(*Lexer).Error("invalid typename in typeof")
					}
				}
			case "addressof":
				if len(yyDollar[2].expr.Compound) != 1 {
					yylex.(*Lexer).Error("addressof takes 1 argument")
				}
				if yyDollar[2].expr.Compound[0].Type != NTAtom {
					yylex.(*Lexer).Error("addressof can only get the address of a variable")
				}
				yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, yyDollar[2].expr)
			case "len":
				switch len(yyDollar[2].expr.Compound) {
				case 0:
					yylex.(*Lexer).Error("len takes 1 argument")
				default:
					yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, yyDollar[2].expr)
				}
			default:
				yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, yyDollar[2].expr)
			}
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 94:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:372
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 95:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:373
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 96:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:376
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].str, "<a>", yyDollar[2].expr, yyDollar[3].expr).setPos0(yyDollar[2].expr.Pos)
		}
	case 97:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:379
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 98:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:380
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 99:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:383
		{
			yyVAL.expr = NewCompoundNode("map", NewCompoundNode()).setPos0(yyDollar[1].token.Pos)
		}
	case 100:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:384
		{
			yyVAL.expr = yyDollar[2].expr.setPos0(yyDollar[1].token.Pos)
		}
	case 101:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:387
		{
			yyVAL.expr = NewCompoundNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 102:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:388
		{
			yyVAL.expr = NewCompoundNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 103:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:389
		{
			yyVAL.expr = NewCompoundNode("array", yyDollar[1].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 104:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:390
		{
			yyVAL.expr = NewCompoundNode("array", yyDollar[1].expr).setPos0(yyDollar[1].expr.Pos)
		}
	}
	goto yystack /* stack new state and value */
}
