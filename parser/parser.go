//line parser.go.y:2
package parser

import __yyfmt__ "fmt"

//line parser.go.y:2
import (
	"bytes"
	"io/ioutil"
	"path/filepath"
)

//line parser.go.y:42
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
const TIf = 57352
const TNil = 57353
const TReturn = 57354
const TRequire = 57355
const TVar = 57356
const TYield = 57357
const TEqeq = 57358
const TNeq = 57359
const TLsh = 57360
const TRsh = 57361
const TLte = 57362
const TGte = 57363
const TIdent = 57364
const TNumber = 57365
const TString = 57366
const TOr = 57367
const TAnd = 57368
const UNARY = 57369

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

//line parser.go.y:415

var typesLookup = map[string]string{
	"nil": "0", "number": "1", "string": "2", "map": "3", "closure": "4", "generic": "5",
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 93,
	44, 20,
	-2, 63,
}

const yyPrivate = 57344

const yyLast = 713

var yyAct = [...]int{

	154, 33, 89, 21, 162, 15, 156, 134, 27, 57,
	133, 43, 44, 41, 39, 20, 18, 159, 161, 4,
	13, 46, 160, 152, 12, 60, 158, 150, 6, 151,
	99, 54, 84, 149, 55, 131, 80, 81, 82, 83,
	58, 90, 1, 176, 175, 108, 168, 100, 91, 93,
	97, 53, 174, 94, 101, 102, 140, 132, 107, 129,
	106, 20, 110, 111, 112, 113, 114, 115, 116, 117,
	118, 119, 120, 121, 122, 123, 124, 125, 126, 127,
	45, 26, 85, 98, 10, 8, 9, 49, 16, 128,
	17, 96, 11, 138, 19, 7, 48, 47, 21, 48,
	79, 145, 23, 143, 148, 48, 24, 172, 144, 19,
	20, 18, 75, 76, 4, 13, 178, 23, 104, 12,
	48, 24, 25, 6, 157, 71, 72, 73, 52, 69,
	70, 71, 72, 73, 153, 50, 155, 87, 21, 21,
	32, 165, 164, 163, 21, 31, 56, 59, 169, 167,
	20, 20, 5, 14, 137, 92, 20, 69, 70, 71,
	72, 73, 34, 51, 88, 130, 2, 21, 0, 179,
	0, 0, 177, 0, 181, 0, 0, 21, 21, 20,
	21, 0, 182, 183, 0, 184, 0, 0, 0, 20,
	20, 0, 20, 67, 68, 75, 76, 66, 65, 0,
	0, 0, 0, 0, 0, 61, 62, 77, 78, 74,
	63, 64, 69, 70, 71, 72, 73, 67, 68, 75,
	76, 66, 65, 0, 0, 0, 146, 147, 0, 61,
	62, 77, 78, 74, 63, 64, 69, 70, 71, 72,
	73, 67, 68, 75, 76, 66, 65, 0, 0, 0,
	0, 173, 0, 61, 62, 77, 78, 74, 63, 64,
	69, 70, 71, 72, 73, 67, 68, 75, 76, 66,
	65, 0, 0, 0, 0, 135, 0, 61, 62, 77,
	78, 74, 63, 64, 69, 70, 71, 72, 73, 67,
	68, 75, 76, 66, 65, 0, 0, 0, 180, 0,
	0, 61, 62, 77, 78, 74, 63, 64, 69, 70,
	71, 72, 73, 67, 68, 75, 76, 66, 65, 0,
	0, 0, 171, 0, 0, 61, 62, 77, 78, 74,
	63, 64, 69, 70, 71, 72, 73, 67, 68, 75,
	76, 66, 65, 142, 0, 0, 0, 0, 0, 61,
	62, 77, 78, 74, 63, 64, 69, 70, 71, 72,
	73, 67, 68, 75, 76, 66, 65, 136, 0, 0,
	0, 0, 0, 61, 62, 77, 78, 74, 63, 64,
	69, 70, 71, 72, 73, 67, 68, 75, 76, 66,
	65, 109, 0, 0, 0, 0, 0, 61, 62, 77,
	78, 74, 63, 64, 69, 70, 71, 72, 73, 0,
	22, 0, 28, 166, 30, 0, 0, 22, 0, 28,
	0, 30, 0, 23, 29, 42, 40, 24, 0, 0,
	23, 29, 42, 40, 24, 0, 0, 35, 0, 0,
	0, 0, 36, 38, 35, 0, 0, 0, 37, 36,
	38, 103, 0, 0, 22, 37, 28, 170, 30, 19,
	0, 22, 0, 28, 0, 30, 0, 23, 29, 42,
	40, 24, 0, 0, 23, 29, 42, 40, 24, 0,
	0, 35, 0, 0, 0, 0, 36, 38, 35, 95,
	0, 0, 37, 36, 38, 0, 139, 0, 22, 37,
	28, 0, 30, 0, 0, 22, 0, 28, 0, 30,
	0, 23, 29, 42, 40, 24, 0, 0, 23, 29,
	42, 40, 24, 0, 0, 35, 0, 0, 0, 0,
	36, 38, 35, 0, 0, 105, 37, 36, 38, 86,
	22, 0, 28, 37, 30, 0, 0, 0, 0, 0,
	0, 0, 0, 23, 29, 42, 40, 24, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 35, 0, 0,
	0, 0, 36, 38, 0, 0, 0, 0, 37, 67,
	68, 75, 76, 66, 65, 0, 0, 0, 0, 0,
	0, 61, 62, 77, 78, 74, 63, 64, 69, 70,
	71, 72, 73, 10, 8, 9, 0, 16, 22, 17,
	0, 11, 0, 19, 7, 0, 0, 0, 0, 0,
	0, 23, 0, 0, 0, 24, 0, 0, 0, 0,
	0, 67, 68, 75, 76, 66, 65, 0, 0, 0,
	0, 0, 141, 3, 62, 77, 78, 74, 63, 64,
	69, 70, 71, 72, 73, 10, 8, 9, 0, 16,
	22, 17, 0, 11, 0, 19, 7, 0, 0, 0,
	0, 0, 0, 23, 0, 0, 0, 24, 0, 0,
	0, 0, 67, 68, 75, 76, 66, 65, 0, 67,
	68, 75, 76, 66, 65, 3, 77, 78, 74, 63,
	64, 69, 70, 71, 72, 73, 63, 64, 69, 70,
	71, 72, 73,
}
var yyPact = [...]int{

	-1000, 651, -1000, -1000, 37, -1000, -1000, 531, -1000, -1000,
	531, 531, -1000, -1000, -1000, 36, 71, 61, 113, 106,
	6, -17, -7, -1000, 531, -1000, -1000, 563, -1000, -1000,
	76, -1000, -1000, -17, -1000, 531, 531, 531, 531, 56,
	496, -1000, -1000, 563, 563, -1000, -1000, 445, -1000, 531,
	56, -22, 2, 531, 401, 96, -1000, 489, -1000, -1,
	345, 531, 531, 531, 531, 531, 531, 531, 531, 531,
	531, 531, 531, 531, 531, 531, 531, 531, 531, -1000,
	-1000, -1000, -1000, -1000, 74, 13, -1000, 14, -42, -45,
	225, 321, 452, -17, 12, -1000, 599, 297, 74, 86,
	531, 563, 177, 531, -1000, -1000, -19, 563, -1000, -1000,
	615, 666, 94, 94, 94, 94, 94, 94, 88, 88,
	-1000, -1000, -1000, 673, 122, 122, 673, 673, -1000, -1000,
	-23, -1000, -1000, 531, 531, 531, 80, 95, 369, -1000,
	-1000, -1000, 80, -1000, 1, 563, -1000, 408, 273, 531,
	-1000, 85, -1000, 201, 563, 563, -1000, -1000, -1000, 8,
	-1000, -1000, -1000, -2, -3, 80, -1000, 109, 531, 249,
	-1000, -1000, -1000, 531, -1000, 80, 80, -1000, 80, 563,
	-1000, 563, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 42, 4, 166, 13, 165, 2, 164, 163, 0,
	162, 1, 5, 26, 22, 155, 154, 18, 6, 17,
	153, 152, 14, 122, 146, 145, 32, 140, 137,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 13, 3, 3, 3, 3, 18,
	18, 18, 18, 18, 18, 21, 21, 21, 12, 12,
	12, 15, 15, 16, 16, 14, 14, 14, 14, 14,
	17, 17, 22, 22, 20, 19, 19, 19, 19, 19,
	19, 19, 4, 4, 4, 4, 4, 4, 5, 5,
	6, 6, 7, 7, 8, 8, 8, 8, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 10, 11, 11,
	11, 11, 23, 24, 24, 25, 26, 26, 27, 27,
	28, 28, 28, 28,
}
var yyR2 = [...]int{

	0, 0, 2, 3, 2, 1, 2, 1, 1, 1,
	1, 2, 1, 1, 1, 1, 1, 1, 2, 3,
	1, 2, 1, 2, 1, 2, 5, 7, 7, 6,
	5, 7, 1, 2, 4, 1, 2, 1, 1, 2,
	1, 2, 1, 4, 6, 5, 5, 3, 1, 3,
	1, 3, 3, 5, 1, 3, 5, 3, 1, 1,
	2, 1, 1, 1, 1, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 2, 2, 2, 2, 1, 1, 3,
	1, 3, 2, 2, 3, 3, 2, 3, 2, 3,
	1, 2, 1, 2,
}
var yyChk = [...]int{

	-1000, -1, -3, 44, -19, -21, -13, 15, 5, 6,
	4, 12, -14, -17, -20, -12, 8, 10, -22, 14,
	-4, -11, 9, 22, 26, -23, 44, -9, 11, 23,
	13, -25, -27, -11, -10, 36, 41, 47, 42, -22,
	25, -4, 24, -9, -9, 44, -2, 26, 25, 26,
	22, -8, 22, 45, 48, 51, -24, 26, 47, -23,
	-9, 28, 29, 33, 34, 21, 20, 16, 17, 35,
	36, 37, 38, 39, 32, 18, 19, 30, 31, 24,
	-9, -9, -9, -9, -26, 26, 43, -28, -7, -6,
	-9, -9, -15, -11, -12, 44, -1, -9, -26, 52,
	45, -9, -9, 50, 22, 46, -6, -9, 46, 46,
	-9, -9, -9, -9, -9, -9, -9, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, -9, -9, -2, 46,
	-5, 22, 43, 52, 52, 50, 46, -16, -9, 44,
	44, 43, 46, -2, 22, -9, 49, 50, -9, 52,
	46, 52, 46, -9, -9, -9, -18, 44, -13, -19,
	-14, -17, -2, -12, -2, 46, 44, -18, 45, -9,
	49, 49, 22, 50, 44, 46, 46, -18, 7, -9,
	49, -9, -18, -18, -18,
}
var yyDef = [...]int{

	1, -2, 2, 5, 0, 7, 8, 35, 37, 38,
	0, 40, 15, 16, 17, 0, 0, 0, 0, 0,
	88, 20, 32, 42, 0, 90, 6, 36, 58, 59,
	0, 61, 62, 63, 64, 0, 0, 0, 0, 0,
	0, 88, 87, 39, 41, 4, 25, 0, 1, 0,
	0, 18, 54, 0, 0, 0, 92, 0, 33, 90,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 60,
	83, 84, 85, 86, 0, 0, 98, 0, 100, 102,
	50, 0, 0, -2, 0, 22, 0, 0, 0, 0,
	0, 19, 0, 0, 47, 93, 0, 50, 89, 91,
	65, 66, 67, 68, 69, 70, 71, 72, 73, 74,
	75, 76, 77, 78, 79, 80, 81, 82, 95, 96,
	0, 48, 99, 101, 103, 0, 0, 0, 0, 24,
	21, 3, 0, 34, 57, 55, 43, 0, 0, 0,
	94, 0, 97, 0, 51, 52, 26, 9, 10, 0,
	12, 13, 14, 0, 0, 0, 23, 30, 0, 0,
	45, 46, 49, 0, 11, 0, 0, 29, 0, 56,
	44, 53, 27, 28, 31,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 47, 3, 42, 3, 39, 31, 3,
	26, 46, 37, 35, 52, 36, 51, 38, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 50, 44,
	34, 45, 33, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 27, 3, 3, 3, 3, 3,
	3, 48, 3, 49, 32, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 25, 30, 43, 41,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 28, 29, 40,
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
		//line parser.go.y:72
		{
			yyVAL.expr = NewCompoundNode("chain")
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:78
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, yyDollar[2].expr)
			yyVAL.expr = yyDollar[1].expr
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 3:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:87
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:90
		{
			if yyDollar[1].expr.isIsolatedDupCall() {
				yyDollar[1].expr.Compound[2].Compound[0] = NewNumberNode("0")
			}
			yyVAL.expr = yyDollar[1].expr
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:98
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:99
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:100
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:101
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:104
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:105
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:106
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:107
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:108
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:109
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:112
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:113
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:114
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 18:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:117
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:120
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
		//line parser.go.y:140
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 21:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:145
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:146
		{
			yyVAL.expr = NewCompoundNode("chain")
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:149
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:150
		{
			yyVAL.expr = NewNumberNode("1")
		}
	case 25:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:153
		{
			yyVAL.expr = NewCompoundNode("for", NewNumberNode("1"), NewCompoundNode(), yyDollar[2].expr).setPos0(yyDollar[1].token.Pos)
		}
	case 26:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:156
		{
			yyVAL.expr = NewCompoundNode("for", yyDollar[3].expr, NewCompoundNode(), yyDollar[5].expr).setPos0(yyDollar[1].token.Pos)
		}
	case 27:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:159
		{
			yyVAL.expr = yyDollar[3].expr
			yyVAL.expr.Compound = append(yyVAL.expr.Compound, NewCompoundNode("for", yyDollar[4].expr, NewCompoundNode("chain", yyDollar[5].expr), yyDollar[7].expr))
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 28:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:164
		{
			yyVAL.expr = yyDollar[3].expr
			yyVAL.expr.Compound = append(yyVAL.expr.Compound, NewCompoundNode("for", yyDollar[4].expr, yyDollar[5].expr, yyDollar[7].expr))
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 29:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:169
		{
			yyVAL.expr = yyDollar[3].expr
			yyVAL.expr.Compound = append(yyVAL.expr.Compound, NewCompoundNode("for", yyDollar[4].expr, NewCompoundNode(), yyDollar[6].expr))
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 30:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:176
		{
			yyVAL.expr = NewCompoundNode("if", yyDollar[3].expr, yyDollar[5].expr, NewCompoundNode())
		}
	case 31:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:177
		{
			yyVAL.expr = NewCompoundNode("if", yyDollar[3].expr, yyDollar[5].expr, yyDollar[7].expr)
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:180
		{
			yyVAL.str = "func"
		}
	case 33:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:181
		{
			yyVAL.str = "safefunc"
		}
	case 34:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:184
		{
			funcname := NewAtomNode(yyDollar[2].token)
			yyVAL.expr = NewCompoundNode(
				"chain",
				NewCompoundNode("set", funcname, NewNilNode()),
				NewCompoundNode("move", funcname, NewCompoundNode(yyDollar[1].str, yyDollar[3].expr, yyDollar[4].expr)))
			yyVAL.expr.Compound[1].Compound[0].Pos = yyDollar[2].token.Pos
			yyVAL.expr.Compound[2].Compound[0].Pos = yyDollar[2].token.Pos
			yyVAL.expr.Compound[2].Compound[2].Compound[0].Pos = yyDollar[2].token.Pos
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:196
		{
			yyVAL.expr = NewCompoundNode("yield").setPos0(yyDollar[1].token.Pos)
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:197
		{
			yyVAL.expr = NewCompoundNode("yield", yyDollar[2].expr).setPos0(yyDollar[1].token.Pos)
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:198
		{
			yyVAL.expr = NewCompoundNode("break").setPos0(yyDollar[1].token.Pos)
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:199
		{
			yyVAL.expr = NewCompoundNode("continue").setPos0(yyDollar[1].token.Pos)
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:200
		{
			yyVAL.expr = NewCompoundNode("assert", yyDollar[2].expr).setPos0(yyDollar[1].token.Pos)
		}
	case 40:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:201
		{
			yyVAL.expr = NewCompoundNode("ret").setPos0(yyDollar[1].token.Pos)
		}
	case 41:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:202
		{
			if yyDollar[2].expr.isIsolatedDupCall() {
				if h, _ := yyDollar[2].expr.Compound[2].Compound[2].Value.(float64); h == 1 {
					yyDollar[2].expr.Compound[2].Compound[2] = NewNumberNode("2")
				}
			}
			yyVAL.expr = NewCompoundNode("ret", yyDollar[2].expr).setPos0(yyDollar[1].token.Pos)
		}
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:212
		{
			yyVAL.expr = NewAtomNode(yyDollar[1].token).setPos(yyDollar[1].token.Pos)
		}
	case 43:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:213
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos).setPos(yyDollar[1].expr.Pos)
		}
	case 44:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:214
		{
			yyVAL.expr = NewCompoundNode("slice", yyDollar[1].expr, yyDollar[3].expr, yyDollar[5].expr).setPos0(yyDollar[1].expr.Pos).setPos(yyDollar[1].expr.Pos)
		}
	case 45:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:215
		{
			yyVAL.expr = NewCompoundNode("slice", yyDollar[1].expr, yyDollar[3].expr, NewNumberNode("-1")).setPos0(yyDollar[1].expr.Pos).setPos(yyDollar[1].expr.Pos)
		}
	case 46:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:216
		{
			yyVAL.expr = NewCompoundNode("slice", yyDollar[1].expr, NewNumberNode("0"), yyDollar[4].expr).setPos0(yyDollar[1].expr.Pos).setPos(yyDollar[1].expr.Pos)
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:217
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, NewStringNode(yyDollar[3].token.Str)).setPos0(yyDollar[1].expr.Pos).setPos(yyDollar[1].expr.Pos)
		}
	case 48:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:220
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].token.Str)
		}
	case 49:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:223
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, NewAtomNode(yyDollar[3].token))
			yyVAL.expr = yyDollar[1].expr
		}
	case 50:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:229
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].expr)
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:232
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, yyDollar[3].expr)
			yyVAL.expr = yyDollar[1].expr
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:238
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 53:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:241
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, yyDollar[3].expr, yyDollar[5].expr)
			yyVAL.expr = yyDollar[1].expr
		}
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:247
		{
			yyVAL.expr = NewCompoundNode("chain", NewCompoundNode("set", NewAtomNode(yyDollar[1].token), NewNilNode()))
			yyVAL.expr.Compound[1].Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:251
		{
			yyVAL.expr = NewCompoundNode("chain", NewCompoundNode("set", NewAtomNode(yyDollar[1].token), yyDollar[3].expr))
			yyVAL.expr.Compound[1].Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 56:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:255
		{
			x := NewCompoundNode("set", NewAtomNode(yyDollar[3].token), yyDollar[5].expr).setPos0(yyDollar[1].expr.Pos)
			yyDollar[1].expr.Compound = append(yyVAL.expr.Compound, x)
			yyVAL.expr = yyDollar[1].expr
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:260
		{
			x := NewCompoundNode("set", NewAtomNode(yyDollar[3].token), NewNilNode()).setPos0(yyDollar[1].expr.Pos)
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, x)
			yyVAL.expr = yyDollar[1].expr
		}
	case 58:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:267
		{
			yyVAL.expr = NewNilNode()
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 59:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:271
		{
			yyVAL.expr = NewNumberNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 60:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:275
		{
			path := filepath.Dir(yyDollar[1].token.Pos.Source)
			path = filepath.Join(path, yyDollar[2].token.Str)

			code, err := ioutil.ReadFile(path)
			if err != nil {
				yylex.(*Lexer).Error(err.Error())
			}
			n, err := Parse(bytes.NewReader(code), path)
			if err != nil {
				yylex.(*Lexer).Error(err.Error())
			}

			// now the required code is loaded, for naming scope we will wrap them into a closure
			cls := NewCompoundNode("func", NewCompoundNode(), n)
			yyVAL.expr = NewCompoundNode("call", cls, NewCompoundNode())
		}
	case 61:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:292
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 62:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:293
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 63:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:294
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 64:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:295
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 65:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:296
		{
			yyVAL.expr = NewCompoundNode("or", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:297
		{
			yyVAL.expr = NewCompoundNode("and", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:298
		{
			yyVAL.expr = NewCompoundNode("<", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:299
		{
			yyVAL.expr = NewCompoundNode("<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:300
		{
			yyVAL.expr = NewCompoundNode("<=", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 70:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:301
		{
			yyVAL.expr = NewCompoundNode("<=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:302
		{
			yyVAL.expr = NewCompoundNode("==", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 72:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:303
		{
			yyVAL.expr = NewCompoundNode("!=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:304
		{
			yyVAL.expr = NewCompoundNode("+", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:305
		{
			yyVAL.expr = NewCompoundNode("-", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 75:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:306
		{
			yyVAL.expr = NewCompoundNode("*", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:307
		{
			yyVAL.expr = NewCompoundNode("/", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 77:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:308
		{
			yyVAL.expr = NewCompoundNode("%", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 78:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:309
		{
			yyVAL.expr = NewCompoundNode("^", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 79:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:310
		{
			yyVAL.expr = NewCompoundNode("<<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 80:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:311
		{
			yyVAL.expr = NewCompoundNode(">>", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 81:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:312
		{
			yyVAL.expr = NewCompoundNode("|", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 82:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:313
		{
			yyVAL.expr = NewCompoundNode("&", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 83:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:314
		{
			yyVAL.expr = NewCompoundNode("-", NewNumberNode("0"), yyDollar[2].expr).setPos0(yyDollar[2].expr.Pos)
		}
	case 84:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:315
		{
			yyVAL.expr = NewCompoundNode("~", yyDollar[2].expr).setPos0(yyDollar[2].expr.Pos)
		}
	case 85:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:316
		{
			yyVAL.expr = NewCompoundNode("!", yyDollar[2].expr).setPos0(yyDollar[2].expr.Pos)
		}
	case 86:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:317
		{
			yyVAL.expr = NewCompoundNode("#", yyDollar[2].expr).setPos0(yyDollar[2].expr.Pos)
		}
	case 87:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:320
		{
			yyVAL.expr = NewStringNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 88:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:326
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 89:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:327
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 90:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:328
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 91:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:329
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 92:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:332
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
	case 93:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:383
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 94:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:384
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 95:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:387
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].str, yyDollar[2].expr, yyDollar[3].expr).setPos0(yyDollar[2].expr.Pos)
		}
	case 96:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:390
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 97:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:391
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 98:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:394
		{
			yyVAL.expr = NewCompoundNode("map", NewCompoundNode()).setPos0(yyDollar[1].token.Pos)
		}
	case 99:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:395
		{
			yyVAL.expr = yyDollar[2].expr.setPos0(yyDollar[1].token.Pos)
		}
	case 100:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:398
		{
			yyVAL.expr = NewCompoundNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 101:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:399
		{
			yyVAL.expr = NewCompoundNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 102:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:400
		{
			table := NewCompoundNode()
			for i, v := range yyDollar[1].expr.Compound {
				table.Compound = append(table.Compound, &Node{Type: NTNumber, Value: float64(i)}, v)
			}
			yyVAL.expr = NewCompoundNode("map", table).setPos0(yyDollar[1].expr.Pos)
		}
	case 103:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:407
		{
			table := NewCompoundNode()
			for i, v := range yyDollar[1].expr.Compound {
				table.Compound = append(table.Compound, &Node{Type: NTNumber, Value: float64(i)}, v)
			}
			yyVAL.expr = NewCompoundNode("map", table).setPos0(yyDollar[1].expr.Pos)
		}
	}
	goto yystack /* stack new state and value */
}
