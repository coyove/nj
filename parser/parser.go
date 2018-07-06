//line parser.go.y:1
package parser

import __yyfmt__ "fmt"

//line parser.go.y:3
import (
	"bytes"
	"io/ioutil"
	"path/filepath"
)

//line parser.go.y:39
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
	"'.'",
	"','",
	"':'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.go.y:593

func TokenName(c int) string {
	if c >= TAnd && c-TAnd < len(yyToknames) {
		if yyToknames[c-TAnd] != "" {
			return yyToknames[c-TAnd]
		}
	}
	return string([]byte{byte(c)})
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 101,
	44, 19,
	-2, 55,
}

const yyPrivate = 57344

const yyLast = 668

var yyAct = [...]int{

	142, 1, 150, 40, 4, 9, 154, 46, 89, 23,
	148, 34, 50, 145, 108, 51, 134, 33, 144, 107,
	53, 151, 153, 5, 16, 133, 58, 152, 131, 15,
	60, 61, 57, 176, 65, 173, 170, 64, 48, 30,
	8, 31, 96, 85, 86, 87, 88, 141, 95, 59,
	29, 164, 129, 163, 7, 98, 102, 100, 101, 138,
	105, 26, 18, 25, 132, 103, 19, 109, 110, 111,
	112, 113, 114, 115, 116, 117, 118, 119, 120, 121,
	122, 123, 124, 125, 126, 94, 172, 74, 75, 76,
	77, 78, 90, 8, 55, 54, 127, 76, 77, 78,
	137, 72, 73, 80, 81, 71, 70, 128, 80, 81,
	140, 84, 168, 66, 67, 82, 83, 79, 68, 69,
	74, 75, 76, 77, 78, 74, 75, 76, 77, 78,
	143, 161, 104, 7, 146, 62, 147, 162, 56, 156,
	9, 18, 159, 9, 28, 19, 4, 9, 20, 92,
	158, 23, 39, 37, 32, 6, 17, 41, 165, 9,
	167, 9, 27, 169, 93, 5, 16, 130, 52, 171,
	9, 15, 9, 2, 9, 8, 9, 9, 8, 174,
	9, 175, 8, 177, 178, 0, 0, 179, 0, 0,
	0, 0, 0, 0, 8, 0, 8, 0, 0, 0,
	0, 0, 0, 0, 0, 8, 0, 8, 0, 8,
	0, 8, 8, 0, 0, 8, 72, 73, 80, 81,
	71, 70, 0, 0, 0, 0, 0, 0, 66, 67,
	82, 83, 79, 68, 69, 74, 75, 76, 77, 78,
	0, 72, 73, 80, 81, 71, 70, 0, 0, 0,
	0, 0, 135, 66, 67, 82, 83, 79, 68, 69,
	74, 75, 76, 77, 78, 72, 73, 80, 81, 71,
	70, 0, 0, 0, 106, 0, 0, 66, 67, 82,
	83, 79, 68, 69, 74, 75, 76, 77, 78, 72,
	73, 80, 81, 71, 70, 139, 0, 0, 0, 0,
	0, 66, 67, 82, 83, 79, 68, 69, 74, 75,
	76, 77, 78, 72, 73, 80, 81, 71, 70, 136,
	0, 0, 0, 0, 0, 66, 67, 82, 83, 79,
	68, 69, 74, 75, 76, 77, 78, 72, 73, 80,
	81, 71, 70, 97, 0, 0, 0, 0, 0, 66,
	67, 82, 83, 79, 68, 69, 74, 75, 76, 77,
	78, 0, 0, 0, 0, 166, 72, 73, 80, 81,
	71, 70, 0, 0, 0, 0, 0, 0, 66, 67,
	82, 83, 79, 68, 69, 74, 75, 76, 77, 78,
	0, 24, 0, 35, 155, 38, 7, 0, 24, 0,
	35, 0, 38, 0, 18, 36, 49, 47, 19, 0,
	0, 18, 36, 49, 47, 19, 0, 0, 42, 0,
	0, 0, 0, 43, 45, 42, 99, 0, 0, 44,
	43, 45, 0, 157, 0, 24, 44, 35, 0, 38,
	0, 0, 0, 24, 0, 35, 0, 38, 18, 36,
	49, 47, 19, 0, 0, 0, 18, 36, 49, 47,
	19, 0, 42, 0, 0, 0, 0, 43, 45, 91,
	42, 0, 0, 44, 0, 43, 45, 0, 0, 0,
	63, 44, 24, 0, 35, 0, 38, 0, 0, 0,
	0, 0, 0, 0, 0, 18, 36, 49, 47, 19,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 42,
	0, 0, 0, 0, 43, 45, 0, 0, 0, 0,
	44, 72, 73, 80, 81, 71, 70, 0, 0, 0,
	0, 0, 0, 66, 67, 82, 83, 79, 68, 69,
	74, 75, 76, 77, 78, 14, 12, 13, 0, 21,
	24, 22, 0, 10, 0, 7, 11, 0, 14, 12,
	13, 0, 21, 18, 22, 0, 10, 19, 7, 11,
	0, 0, 0, 0, 0, 0, 18, 0, 0, 128,
	19, 0, 0, 0, 160, 3, 72, 73, 80, 81,
	71, 70, 0, 0, 0, 0, 0, 0, 149, 67,
	82, 83, 79, 68, 69, 74, 75, 76, 77, 78,
	14, 12, 13, 0, 21, 24, 22, 0, 10, 0,
	7, 11, 0, 0, 0, 0, 0, 0, 18, 0,
	0, 0, 19, 0, 0, 0, 0, 72, 73, 80,
	81, 71, 70, 0, 72, 73, 80, 81, 71, 70,
	3, 82, 83, 79, 68, 69, 74, 75, 76, 77,
	78, 68, 69, 74, 75, 76, 77, 78,
}
var yyPact = [...]int{

	-1000, 606, -1000, -1000, 19, 17, -1000, 122, 5, -9,
	473, 473, -1000, -1000, 473, -1000, -1000, -1000, -1000, 473,
	-1000, 69, 68, 116, -15, -1000, -1000, -25, 4, 473,
	473, 113, -1000, 434, 505, -1000, -1000, -1000, 87, -1000,
	-9, -1000, 473, 473, 473, 473, 66, 426, -1000, -1000,
	505, 505, -4, 297, 382, 473, 66, -1000, 110, 473,
	505, 225, -1000, -1000, -32, 505, 473, 473, 473, 473,
	473, 473, 473, 473, 473, 473, 473, 473, 473, 473,
	473, 473, 473, 473, -1000, -1000, -1000, -1000, -1000, 82,
	6, -1000, 21, -26, -35, 200, -1000, -1000, 273, 473,
	15, -9, 249, 82, 2, 505, -1000, 473, -1000, 570,
	621, 90, 90, 90, 90, 90, 90, 60, 60, -1000,
	-1000, -1000, 628, 52, 52, 628, 628, -1000, -1000, -1000,
	-33, -1000, -1000, 473, 473, 473, 554, 350, 389, 554,
	-1000, 473, 505, 541, 109, -1000, 85, 505, -1000, -1000,
	9, 7, -1000, -1000, -1000, 119, 321, 119, 105, 505,
	-1000, -1000, 473, -1000, -1000, -10, 40, -11, 554, 505,
	554, -13, 554, 554, -1000, -1000, 554, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 1, 6, 173, 38, 167, 37, 164, 162, 0,
	157, 3, 2, 27, 22, 10, 21, 156, 155, 7,
	148, 154, 153, 8, 152, 149,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 3, 3, 3, 3, 15, 15,
	15, 15, 15, 15, 18, 18, 18, 12, 12, 12,
	13, 13, 13, 13, 13, 14, 14, 19, 19, 17,
	16, 16, 16, 16, 16, 16, 16, 4, 4, 4,
	5, 5, 6, 6, 7, 7, 8, 8, 8, 8,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 10,
	11, 11, 11, 11, 20, 21, 21, 22, 23, 23,
	24, 24, 25, 25, 25, 25,
}
var yyR2 = [...]int{

	0, 0, 2, 3, 1, 2, 2, 1, 1, 2,
	2, 1, 1, 1, 1, 1, 1, 2, 3, 1,
	5, 8, 9, 8, 8, 5, 7, 1, 2, 4,
	1, 2, 1, 2, 1, 1, 2, 1, 4, 3,
	1, 3, 1, 3, 3, 5, 1, 3, 5, 3,
	1, 1, 1, 2, 1, 1, 1, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 2, 2, 2, 2, 1,
	1, 3, 1, 3, 2, 2, 3, 3, 2, 3,
	2, 3, 1, 2, 1, 2,
}
var yyChk = [...]int{

	-1000, -1, -3, 44, -12, -16, -18, 14, -4, -11,
	12, 15, 5, 6, 4, -13, -14, -17, 22, 26,
	-20, 8, 10, -19, 9, 44, 44, -8, 22, 45,
	48, 50, -21, 26, -9, 11, 23, -22, 13, -24,
	-11, -10, 36, 41, 47, 42, -19, 25, -4, 24,
	-9, -9, -20, -9, 26, 26, 22, 47, 51, 45,
	-9, -9, 22, 46, -6, -9, 28, 29, 33, 34,
	21, 20, 16, 17, 35, 36, 37, 38, 39, 32,
	18, 19, 30, 31, 24, -9, -9, -9, -9, -23,
	26, 43, -25, -7, -6, -9, 46, 46, -9, 44,
	-12, -11, -9, -23, 22, -9, 49, 51, 46, -9,
	-9, -9, -9, -9, -9, -9, -9, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, -9, -2, 25, 46,
	-5, 22, 43, 51, 51, 52, 46, -9, 44, 46,
	-2, 45, -9, -1, 51, 46, -9, -9, -15, 44,
	-12, -16, -13, -14, -2, 44, -9, 44, -15, -9,
	43, 22, 52, 44, 44, -12, 44, -12, 7, -9,
	46, -12, 46, 46, -15, -15, 46, -15, -15, -15,
}
var yyDef = [...]int{

	1, -2, 2, 4, 0, 0, 7, 0, 80, 19,
	30, 32, 34, 35, 0, 14, 15, 16, 37, 0,
	82, 0, 0, 0, 27, 5, 6, 17, 46, 0,
	0, 0, 84, 0, 31, 50, 51, 52, 0, 54,
	55, 56, 0, 0, 0, 0, 0, 0, 80, 79,
	33, 36, 82, 0, 0, 0, 0, 28, 0, 0,
	18, 0, 39, 85, 0, 42, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 53, 75, 76, 77, 78, 0,
	0, 90, 0, 92, 94, 42, 81, 83, 0, 0,
	0, -2, 0, 0, 49, 47, 38, 0, 86, 57,
	58, 59, 60, 61, 62, 63, 64, 65, 66, 67,
	68, 69, 70, 71, 72, 73, 74, 87, 1, 88,
	0, 40, 91, 93, 95, 0, 0, 0, 0, 0,
	29, 0, 43, 0, 0, 89, 0, 44, 20, 8,
	0, 0, 11, 12, 13, 0, 0, 0, 25, 48,
	3, 41, 0, 9, 10, 0, 0, 0, 0, 45,
	0, 0, 0, 0, 26, 21, 0, 23, 24, 22,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 47, 3, 42, 3, 39, 31, 3,
	26, 46, 37, 35, 51, 36, 50, 38, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 52, 44,
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
		//line parser.go.y:69
		{
			yyVAL.expr = NewCompoundNode("chain")
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:75
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, yyDollar[2].expr)
			yyVAL.expr = yyDollar[1].expr
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 3:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:84
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:89
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:92
		{
			if yyDollar[1].expr.isIsolatedDupCall() {
				yyDollar[1].expr.Compound[2].Compound[0] = NewNumberNode("0")
			}
			yyVAL.expr = yyDollar[1].expr
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:98
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:101
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:106
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:109
		{
			if yyDollar[1].expr.isIsolatedDupCall() {
				yyDollar[1].expr.Compound[2].Compound[0] = NewNumberNode("0")
			}
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:115
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:118
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:121
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:124
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:129
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:132
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:135
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 17:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:140
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:143
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
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:163
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 20:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:168
		{
			yyVAL.expr = NewCompoundNode("for", yyDollar[3].expr, NewCompoundNode(), yyDollar[5].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 21:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.go.y:172
		{
			yyVAL.expr = NewCompoundNode("for", yyDollar[4].expr, NewCompoundNode("chain", yyDollar[6].expr), yyDollar[8].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 22:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.go.y:176
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[3].expr, NewCompoundNode("for", yyDollar[5].expr, NewCompoundNode("chain", yyDollar[7].expr), yyDollar[9].expr))
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 23:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.go.y:180
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[3].expr, NewCompoundNode("for", yyDollar[5].expr, NewCompoundNode(), yyDollar[8].expr))
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 24:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.go.y:184
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[3].expr, NewCompoundNode("for", NewNumberNode("1"), NewCompoundNode("chain", yyDollar[6].expr), yyDollar[8].expr))
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 25:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:190
		{
			yyVAL.expr = NewCompoundNode("if", yyDollar[3].expr, yyDollar[5].expr, NewCompoundNode())
		}
	case 26:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:193
		{
			yyVAL.expr = NewCompoundNode("if", yyDollar[3].expr, yyDollar[5].expr, yyDollar[7].expr)
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:198
		{
			yyVAL.str = "func"
		}
	case 28:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:201
		{
			yyVAL.str = "safefunc"
		}
	case 29:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:206
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
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:218
		{
			yyVAL.expr = NewCompoundNode("ret")
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 31:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:222
		{
			if yyDollar[2].expr.isIsolatedDupCall() {
				if h, _ := yyDollar[2].expr.Compound[2].Compound[2].Value.(float64); h == 1 {
					yyDollar[2].expr.Compound[2].Compound[2] = NewNumberNode("2")
				}
			}
			yyVAL.expr = NewCompoundNode("ret", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:231
		{
			yyVAL.expr = NewCompoundNode("yield")
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 33:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:235
		{
			yyVAL.expr = NewCompoundNode("yield", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:239
		{
			yyVAL.expr = NewCompoundNode("break")
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:243
		{
			yyVAL.expr = NewCompoundNode("continue")
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:247
		{
			yyVAL.expr = NewCompoundNode("assert", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:253
		{
			yyVAL.expr = NewAtomNode(yyDollar[1].token)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 38:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:257
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:262
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, NewStringNode(yyDollar[3].token.Str))
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 40:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:269
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].token.Str)
		}
	case 41:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:272
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, NewAtomNode(yyDollar[3].token))
			yyVAL.expr = yyDollar[1].expr
		}
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:278
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].expr)
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:281
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, yyDollar[3].expr)
			yyVAL.expr = yyDollar[1].expr
		}
	case 44:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:287
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 45:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:290
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, yyDollar[3].expr, yyDollar[5].expr)
			yyVAL.expr = yyDollar[1].expr
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:296
		{
			yyVAL.expr = NewCompoundNode("chain", NewCompoundNode("set", NewAtomNode(yyDollar[1].token), NewNilNode()))
			yyVAL.expr.Compound[1].Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:300
		{
			yyVAL.expr = NewCompoundNode("chain", NewCompoundNode("set", NewAtomNode(yyDollar[1].token), yyDollar[3].expr))
			yyVAL.expr.Compound[1].Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 48:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:304
		{
			x := NewCompoundNode("set", NewAtomNode(yyDollar[3].token), yyDollar[5].expr)
			x.Compound[0].Pos = yyDollar[1].expr.Pos
			yyDollar[1].expr.Compound = append(yyVAL.expr.Compound, x)
			yyVAL.expr = yyDollar[1].expr
		}
	case 49:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:310
		{
			x := NewCompoundNode("set", NewAtomNode(yyDollar[3].token), NewNilNode())
			x.Compound[0].Pos = yyDollar[1].expr.Pos
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, x)
			yyVAL.expr = yyDollar[1].expr
		}
	case 50:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:318
		{
			yyVAL.expr = NewNilNode()
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:322
		{
			yyVAL.expr = NewNumberNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 52:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:326
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 53:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:329
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
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:346
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 55:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:349
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 56:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:352
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:355
		{
			yyVAL.expr = NewCompoundNode("or", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:359
		{
			yyVAL.expr = NewCompoundNode("and", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:363
		{
			yyVAL.expr = NewCompoundNode("<", yyDollar[3].expr, yyDollar[1].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:367
		{
			yyVAL.expr = NewCompoundNode("<", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:371
		{
			yyVAL.expr = NewCompoundNode("<=", yyDollar[3].expr, yyDollar[1].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 62:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:375
		{
			yyVAL.expr = NewCompoundNode("<=", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 63:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:379
		{
			yyVAL.expr = NewCompoundNode("==", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:383
		{
			yyVAL.expr = NewCompoundNode("!=", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 65:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:387
		{
			yyVAL.expr = NewCompoundNode("+", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:391
		{
			yyVAL.expr = NewCompoundNode("-", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:395
		{
			yyVAL.expr = NewCompoundNode("*", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:399
		{
			yyVAL.expr = NewCompoundNode("/", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:403
		{
			yyVAL.expr = NewCompoundNode("%", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 70:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:407
		{
			yyVAL.expr = NewCompoundNode("^", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:411
		{
			yyVAL.expr = NewCompoundNode("<<", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 72:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:415
		{
			yyVAL.expr = NewCompoundNode(">>", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:419
		{
			yyVAL.expr = NewCompoundNode("|", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:423
		{
			yyVAL.expr = NewCompoundNode("&", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 75:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:427
		{
			yyVAL.expr = NewCompoundNode("-", NewNumberNode("0"), yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 76:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:431
		{
			yyVAL.expr = NewCompoundNode("~", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 77:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:435
		{
			yyVAL.expr = NewCompoundNode("!", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 78:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:439
		{
			yyVAL.expr = NewCompoundNode("#", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 79:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:445
		{
			yyVAL.expr = NewStringNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 80:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:451
		{
			yyVAL.expr = yyDollar[1].expr
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 81:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:455
		{
			yyVAL.expr = yyDollar[2].expr
			yyVAL.expr.Pos = yyDollar[2].expr.Pos
		}
	case 82:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:459
		{
			yyVAL.expr = yyDollar[1].expr
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 83:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:463
		{
			yyVAL.expr = yyDollar[2].expr
			yyVAL.expr.Pos = yyDollar[2].expr.Pos
		}
	case 84:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:469
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
					switch x, _ := yyDollar[2].expr.Compound[1].Value.(string); x {
					case "nil":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].expr.Compound[0], NewNumberNode("0")))
					case "number":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].expr.Compound[0], NewNumberNode("1")))
					case "string":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].expr.Compound[0], NewNumberNode("2")))
					case "map":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].expr.Compound[0], NewNumberNode("3")))
					case "closure":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].expr.Compound[0], NewNumberNode("4")))
					case "generic":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].expr.Compound[0], NewNumberNode("5")))
					default:
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].expr.Compound[0], yyDollar[2].expr.Compound[1]))
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
	case 85:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:530
		{
			if yylex.(*Lexer).PNewLine {
				yylex.(*Lexer).TokenError(yyDollar[1].token, "ambiguous syntax (function call x new statement)")
			}
			yyVAL.expr = NewCompoundNode()
		}
	case 86:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:536
		{
			if yylex.(*Lexer).PNewLine {
				yylex.(*Lexer).TokenError(yyDollar[1].token, "ambiguous syntax (function call x new statement)")
			}
			yyVAL.expr = yyDollar[2].expr
		}
	case 87:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:544
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].str, yyDollar[2].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 88:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:550
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 89:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:553
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 90:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:558
		{
			yyVAL.expr = NewCompoundNode("map", NewCompoundNode())
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 91:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:562
		{
			yyVAL.expr = yyDollar[2].expr
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 92:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:568
		{
			yyVAL.expr = NewCompoundNode("map", yyDollar[1].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 93:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:572
		{
			yyVAL.expr = NewCompoundNode("map", yyDollar[1].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 94:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:576
		{
			table := NewCompoundNode()
			for i, v := range yyDollar[1].expr.Compound {
				table.Compound = append(table.Compound, &Node{Type: NTNumber, Value: float64(i)}, v)
			}
			yyVAL.expr = NewCompoundNode("map", table)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 95:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:584
		{
			table := NewCompoundNode()
			for i, v := range yyDollar[1].expr.Compound {
				table.Compound = append(table.Compound, &Node{Type: NTNumber, Value: float64(i)}, v)
			}
			yyVAL.expr = NewCompoundNode("map", table)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	}
	goto yystack /* stack new state and value */
}
