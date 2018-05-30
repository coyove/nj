//line parser.go.y:2
package parser

import __yyfmt__ "fmt"

//line parser.go.y:2
import (
	"bytes"
	"io/ioutil"
	"path/filepath"
)

//line parser.go.y:28
type yySymType struct {
	yys   int
	token Token

	stmts *Node
	stmt  *Node

	funcname interface{}
	funcexpr interface{}

	exprlist *Node
	expr     *Node

	namelist *Node
}

const TAnd = 57346
const TAssert = 57347
const TBreak = 57348
const TContinue = 57349
const TDo = 57350
const TElse = 57351
const TElseIf = 57352
const TEnd = 57353
const TFalse = 57354
const TFor = 57355
const TIf = 57356
const TLambda = 57357
const TList = 57358
const TNil = 57359
const TNot = 57360
const TOr = 57361
const TReturn = 57362
const TRequire = 57363
const TSet = 57364
const TThen = 57365
const TTrue = 57366
const TYield = 57367
const TEqeq = 57368
const TNeq = 57369
const TLsh = 57370
const TRsh = 57371
const TLte = 57372
const TGte = 57373
const TIdent = 57374
const TNumber = 57375
const TString = 57376
const UNARY = 57377

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"TAnd",
	"TAssert",
	"TBreak",
	"TContinue",
	"TDo",
	"TElse",
	"TElseIf",
	"TEnd",
	"TFalse",
	"TFor",
	"TIf",
	"TLambda",
	"TList",
	"TNil",
	"TNot",
	"TOr",
	"TReturn",
	"TRequire",
	"TSet",
	"TThen",
	"TTrue",
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
	"';'",
	"'='",
	"','",
	"'['",
	"']'",
	"'.'",
	"')'",
	"'}'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.go.y:508

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
}

const yyPrivate = 57344

const yyLast = 739

var yyAct = [...]int{

	56, 55, 89, 43, 91, 93, 130, 25, 118, 42,
	120, 45, 46, 119, 93, 47, 24, 20, 44, 2,
	49, 51, 52, 81, 1, 32, 125, 5, 39, 93,
	4, 48, 21, 94, 22, 77, 78, 79, 80, 85,
	84, 116, 88, 89, 67, 68, 69, 70, 71, 69,
	70, 71, 53, 44, 41, 138, 139, 137, 128, 18,
	97, 98, 99, 100, 101, 102, 103, 104, 105, 106,
	107, 108, 109, 110, 111, 112, 113, 114, 96, 50,
	86, 31, 95, 30, 5, 117, 23, 4, 17, 33,
	124, 83, 134, 0, 126, 0, 0, 0, 14, 12,
	13, 0, 0, 0, 144, 115, 6, 8, 7, 0,
	0, 122, 123, 10, 15, 9, 73, 74, 11, 131,
	0, 5, 132, 0, 4, 16, 0, 0, 0, 19,
	67, 68, 69, 70, 71, 0, 0, 141, 0, 0,
	143, 5, 3, 0, 4, 0, 0, 0, 5, 5,
	0, 4, 4, 135, 60, 0, 0, 0, 0, 0,
	0, 5, 0, 142, 4, 0, 0, 0, 5, 59,
	146, 4, 5, 0, 0, 4, 65, 66, 73, 74,
	64, 63, 0, 0, 60, 0, 0, 75, 76, 72,
	61, 62, 67, 68, 69, 70, 71, 0, 0, 59,
	0, 0, 0, 0, 0, 90, 65, 66, 73, 74,
	64, 63, 60, 0, 0, 0, 57, 75, 76, 72,
	61, 62, 67, 68, 69, 70, 71, 59, 0, 0,
	0, 0, 0, 92, 65, 66, 73, 74, 64, 63,
	60, 0, 0, 0, 0, 75, 76, 72, 61, 62,
	67, 68, 69, 70, 71, 59, 0, 0, 0, 58,
	0, 0, 65, 66, 73, 74, 64, 63, 60, 0,
	0, 0, 0, 75, 76, 72, 61, 62, 67, 68,
	69, 70, 71, 59, 0, 0, 136, 0, 0, 0,
	65, 66, 73, 74, 64, 63, 0, 0, 0, 0,
	0, 75, 76, 72, 61, 62, 67, 68, 69, 70,
	71, 0, 27, 0, 121, 37, 0, 26, 36, 0,
	0, 0, 0, 0, 28, 0, 0, 0, 0, 0,
	0, 0, 16, 29, 40, 38, 19, 0, 0, 27,
	0, 0, 37, 34, 26, 36, 0, 0, 35, 0,
	0, 28, 0, 0, 60, 0, 82, 0, 0, 16,
	29, 40, 38, 19, 0, 0, 0, 0, 0, 59,
	34, 0, 0, 145, 0, 35, 65, 66, 73, 74,
	64, 63, 54, 60, 0, 0, 0, 75, 76, 72,
	61, 62, 67, 68, 69, 70, 71, 0, 59, 0,
	0, 0, 87, 0, 0, 65, 66, 73, 74, 64,
	63, 60, 0, 0, 0, 0, 75, 76, 72, 61,
	62, 67, 68, 69, 70, 71, 59, 0, 0, 0,
	0, 0, 0, 65, 66, 73, 74, 64, 63, 0,
	0, 0, 0, 0, 75, 76, 72, 61, 62, 67,
	68, 69, 70, 71, 14, 12, 13, 0, 0, 0,
	140, 0, 6, 8, 7, 0, 0, 0, 0, 10,
	15, 9, 0, 0, 11, 0, 0, 0, 0, 0,
	0, 16, 0, 0, 0, 19, 14, 12, 13, 0,
	0, 0, 133, 0, 6, 8, 7, 0, 3, 0,
	0, 10, 15, 9, 0, 0, 11, 0, 0, 0,
	0, 0, 0, 16, 0, 0, 0, 19, 14, 12,
	13, 0, 0, 0, 129, 0, 6, 8, 7, 0,
	3, 60, 0, 10, 15, 9, 0, 0, 11, 0,
	0, 0, 0, 0, 0, 16, 0, 0, 0, 19,
	0, 0, 0, 65, 66, 73, 74, 64, 63, 0,
	0, 0, 3, 0, 75, 76, 72, 61, 62, 67,
	68, 69, 70, 71, 14, 12, 13, 0, 0, 0,
	127, 0, 6, 8, 7, 0, 0, 0, 0, 10,
	15, 9, 0, 0, 11, 0, 0, 0, 0, 0,
	0, 16, 0, 0, 0, 19, 14, 12, 13, 0,
	0, 0, 0, 0, 6, 8, 7, 0, 3, 0,
	0, 10, 15, 9, 0, 0, 11, 0, 0, 0,
	0, 0, 0, 16, 0, 0, 0, 19, 0, 0,
	65, 66, 73, 74, 64, 63, 0, 0, 0, 0,
	3, 75, 76, 72, 61, 62, 67, 68, 69, 70,
	71, 27, 0, 0, 37, 0, 26, 36, 0, 0,
	0, 0, 0, 28, 0, 0, 0, 0, 0, 0,
	0, 16, 29, 40, 38, 19, 65, 66, 73, 74,
	64, 63, 34, 0, 0, 0, 0, 35, 0, 0,
	61, 62, 67, 68, 69, 70, 71, 14, 12, 13,
	0, 0, 0, 0, 0, 6, 8, 7, 0, 0,
	0, 0, 10, 15, 9, 0, 0, 11, 0, 0,
	0, 0, 0, 0, 16, 0, 0, 0, 19,
}
var yyPact = [...]int{

	-1000, 601, -1000, -1000, -33, -20, 649, 22, 649, 21,
	649, 649, -1000, -1000, 649, -3, -1000, -1000, -1000, 649,
	649, 649, 20, -1000, 327, 208, -1000, -1000, -1000, -1000,
	-1000, -1000, -20, -1000, 649, 649, 649, -13, 300, -1000,
	-1000, -13, 379, -8, -1000, 407, 407, 407, -1000, 150,
	-51, 407, 180, -1000, -1000, -22, 407, -1000, 702, 649,
	649, 649, 649, 649, 649, 649, 649, 649, 649, 649,
	649, 649, 649, 649, 649, 649, 649, -1000, -1000, -1000,
	-1000, -14, -1000, -43, -46, 264, -1000, -1000, 649, -6,
	-1000, -1000, -1000, 649, -1000, 569, 50, 527, 614, 88,
	88, 88, 88, 88, 88, 5, 5, -1000, -1000, -1000,
	660, 2, 2, 660, 660, 513, -1000, -49, 649, -1000,
	-1000, 649, 481, 601, -37, -1000, 407, -1000, -1000, -1000,
	-1000, 236, 407, -1000, 46, 449, 649, -1000, -1000, 649,
	-1000, 407, 93, 350, -1000, -1000, 601,
}
var yyPgo = [...]int{

	0, 24, 19, 92, 28, 3, 1, 91, 0, 89,
	25, 59, 88, 86, 83, 38, 81,
}
var yyR1 = [...]int{

	0, 1, 1, 1, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
	3, 3, 4, 4, 4, 5, 5, 6, 6, 7,
	7, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	9, 10, 10, 10, 10, 12, 11, 13, 13, 14,
	15, 15, 16, 16, 16,
}
var yyR2 = [...]int{

	0, 0, 2, 2, 3, 1, 5, 7, 5, 6,
	8, 4, 1, 2, 1, 2, 1, 1, 2, 2,
	0, 5, 1, 4, 3, 1, 3, 1, 3, 3,
	5, 1, 1, 1, 1, 1, 1, 1, 1, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 2, 2, 2,
	1, 1, 1, 1, 3, 3, 2, 2, 3, 4,
	2, 3, 2, 3, 3,
}
var yyChk = [...]int{

	-1000, -1, -2, 49, -4, -10, 13, 15, 14, 22,
	20, 25, 6, 7, 5, 21, 32, -12, -11, 36,
	50, 52, 54, -13, 36, -8, 17, 12, 24, 33,
	-14, -16, -10, -9, 43, 48, 18, 15, 35, -4,
	34, 32, -8, -5, 32, -8, -8, -8, 34, -8,
	-11, -8, -8, 32, 55, -6, -8, 8, 51, 19,
	4, 40, 41, 31, 30, 26, 27, 42, 43, 44,
	45, 46, 39, 28, 29, 37, 38, -8, -8, -8,
	-15, 36, 56, -7, -6, -8, -15, 23, 50, 51,
	55, 55, 53, 51, 55, -1, -2, -8, -8, -8,
	-8, -8, -8, -8, -8, -8, -8, -8, -8, -8,
	-8, -8, -8, -8, -8, -1, 55, -5, 51, 56,
	56, 50, -1, -1, -6, 32, -8, 11, 8, 11,
	55, -8, -8, 11, -3, -1, 50, 11, 9, 10,
	11, -8, -1, -8, 11, 23, -1,
}
var yyDef = [...]int{

	1, -2, 2, 3, 61, 5, 0, 0, 0, 0,
	12, 14, 16, 17, 0, 0, 22, 62, 63, 0,
	0, 0, 0, 66, 0, 0, 31, 32, 33, 34,
	35, 36, 37, 38, 0, 0, 0, 0, 0, 61,
	60, 0, 0, 0, 25, 13, 15, 18, 19, 0,
	63, 4, 0, 24, 67, 0, 27, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 57, 58, 59,
	1, 0, 72, 0, 0, 27, 1, 1, 0, 0,
	64, 65, 23, 0, 68, 0, 0, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 0, 70, 0, 0, 73,
	74, 0, 0, 20, 11, 26, 28, 6, 1, 69,
	71, 0, 29, 8, 0, 0, 0, 9, 1, 0,
	7, 30, 0, 0, 10, 1, 21,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 46, 38, 3,
	36, 55, 44, 42, 51, 43, 54, 45, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 49,
	41, 50, 40, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 52, 3, 53, 39, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 35, 37, 56, 48,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 47,
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
		//line parser.go.y:63
		{
			yyVAL.stmts = NewCompoundNode("chain")
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.stmts
			}
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:69
		{
			if yyDollar[2].stmt.isIsolatedDupCall() {
				yyDollar[2].stmt.Compound[2].Compound[0] = NewNumberNode("0")
			}
			yyDollar[1].stmts.Compound = append(yyDollar[1].stmts.Compound, yyDollar[2].stmt)
			yyVAL.stmts = yyDollar[1].stmts
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.stmts
			}
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:79
		{
			yyVAL.stmts = yyDollar[1].stmts
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.stmts
			}
		}
	case 4:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:87
		{
			yyVAL.stmt = NewCompoundNode("move", yyDollar[1].expr, yyDollar[3].expr)
			if len(yyDollar[1].expr.Compound) > 0 {
				if c, _ := yyDollar[1].expr.Compound[0].Value.(string); c == "load" {
					yyVAL.stmt = NewCompoundNode("store", yyDollar[1].expr.Compound[1], yyDollar[1].expr.Compound[2], yyDollar[3].expr)
				}
			}
			if c, _ := yyDollar[1].expr.Value.(string); c != "" && yyDollar[1].expr.Type == NTAtom {
				if a, b, s := yyDollar[3].expr.isSimpleAddSub(); a == c {
					yyDollar[3].expr.Compound[2].Value = yyDollar[3].expr.Compound[2].Value.(float64) * s
					yyVAL.stmt = NewCompoundNode("inc", yyDollar[1].expr, yyDollar[3].expr.Compound[2])
					yyVAL.stmt.Compound[1].Pos = yyDollar[1].expr.Pos
				} else if b == c {
					yyDollar[3].expr.Compound[1].Value = yyDollar[3].expr.Compound[1].Value.(float64) * s
					yyVAL.stmt = NewCompoundNode("inc", yyDollar[1].expr, yyDollar[3].expr.Compound[1])
					yyVAL.stmt.Compound[1].Pos = yyDollar[1].expr.Pos
				}
			}
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:108
		{
			// if _, ok := $1.(*FuncCallExpr); !ok {
			//    yylex.(*Lexer).Error("parse error")
			// } else {
			yyVAL.stmt = yyDollar[1].expr
			// }
		}
	case 6:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:115
		{
			yyVAL.stmt = NewCompoundNode("for", yyDollar[2].expr, NewCompoundNode(), yyDollar[4].stmts)
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 7:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:119
		{
			yyVAL.stmt = NewCompoundNode("for", yyDollar[2].expr, NewCompoundNode("chain", yyDollar[4].stmt), yyDollar[6].stmts)
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 8:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:123
		{
			funcname := NewAtomNode(yyDollar[2].token)
			yyVAL.stmt = NewCompoundNode("chain", NewCompoundNode("set", funcname, NewCompoundNode("nil")), NewCompoundNode("move", funcname, NewCompoundNode("lambda", yyDollar[3].expr, yyDollar[4].stmts)))
		}
	case 9:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:127
		{
			yyVAL.stmt = NewCompoundNode("if", yyDollar[2].expr, yyDollar[4].stmts, NewCompoundNode())
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
			cur := yyVAL.stmt
			for _, e := range yyDollar[5].stmts.Compound {
				cur.Compound[3] = NewCompoundNode("chain", e)
				cur = e
			}
		}
	case 10:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.go.y:136
		{
			yyVAL.stmt = NewCompoundNode("if", yyDollar[2].expr, yyDollar[4].stmts, NewCompoundNode())
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
			cur := yyVAL.stmt
			for _, e := range yyDollar[5].stmts.Compound {
				cur.Compound[3] = NewCompoundNode("chain", e)
				cur = e
			}
			cur.Compound[3] = yyDollar[7].stmts
		}
	case 11:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:146
		{
			yyVAL.stmt = NewCompoundNode("chain")
			for i, name := range yyDollar[2].namelist.Compound {
				var e *Node
				if i < len(yyDollar[4].exprlist.Compound) {
					e = yyDollar[4].exprlist.Compound[i]
				} else {
					e = yyDollar[4].exprlist.Compound[len(yyDollar[4].exprlist.Compound)-1]
				}
				c := NewCompoundNode("set", name, e)
				name.Pos, e.Pos = yyDollar[1].token.Pos, yyDollar[1].token.Pos
				c.Compound[0].Pos = yyDollar[1].token.Pos
				yyVAL.stmt.Compound = append(yyVAL.stmt.Compound, c)
			}
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:161
		{
			yyVAL.stmt = NewCompoundNode("ret")
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 13:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:165
		{
			if yyDollar[2].expr.isIsolatedDupCall() {
				if h, _ := yyDollar[2].expr.Compound[2].Compound[2].Value.(float64); h == 1 {
					yyDollar[2].expr.Compound[2].Compound[2] = NewNumberNode("2")
				}
			}
			yyVAL.stmt = NewCompoundNode("ret", yyDollar[2].expr)
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:174
		{
			yyVAL.stmt = NewCompoundNode("yield")
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:178
		{
			yyVAL.stmt = NewCompoundNode("yield", yyDollar[2].expr)
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:182
		{
			yyVAL.stmt = NewCompoundNode("break")
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:186
		{
			yyVAL.stmt = NewCompoundNode("continue")
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 18:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:190
		{
			yyVAL.stmt = NewCompoundNode("assert", yyDollar[2].expr)
			yyVAL.stmt.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:194
		{
			path := filepath.Dir(yyDollar[1].token.Pos.Source)
			path = filepath.Join(path, yyDollar[2].token.Str)
			filename := filepath.Base(yyDollar[2].token.Str)
			filename = filename[:len(filename)-len(filepath.Ext(filename))]

			code, err := ioutil.ReadFile(path)
			if err != nil {
				yylex.(*Lexer).Error(err.Error())
			}
			n, err := Parse(bytes.NewReader(code), path)
			if err != nil {
				yylex.(*Lexer).Error(err.Error())
			}

			// now the required code is loaded, for naming scope we will wrap them into a closure
			cls := NewCompoundNode("lambda", NewCompoundNode(), n)
			call := NewCompoundNode("call", cls, NewCompoundNode())
			yyVAL.stmt = NewCompoundNode("set", filename, call)
		}
	case 20:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.go.y:216
		{
			yyVAL.stmts = NewCompoundNode()
		}
	case 21:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:219
		{
			yyVAL.stmts.Compound = append(yyDollar[1].stmts.Compound, NewCompoundNode("if", yyDollar[3].expr, yyDollar[5].stmts, NewCompoundNode()))
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:224
		{
			yyVAL.expr = NewAtomNode(yyDollar[1].token)
		}
	case 23:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:227
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, yyDollar[3].expr)
		}
	case 24:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:230
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, NewStringNode(yyDollar[3].token.Str))
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:235
		{
			yyVAL.namelist = NewCompoundNode(yyDollar[1].token.Str)
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:238
		{
			yyDollar[1].namelist.Compound = append(yyDollar[1].namelist.Compound, NewAtomNode(yyDollar[3].token))
			yyVAL.namelist = yyDollar[1].namelist
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:244
		{
			yyVAL.exprlist = NewCompoundNode(yyDollar[1].expr)
		}
	case 28:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:247
		{
			yyDollar[1].exprlist.Compound = append(yyDollar[1].exprlist.Compound, yyDollar[3].expr)
			yyVAL.exprlist = yyDollar[1].exprlist
		}
	case 29:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:253
		{
			yyVAL.exprlist = NewCompoundNode(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 30:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:256
		{
			yyDollar[1].exprlist.Compound = append(yyDollar[1].exprlist.Compound, yyDollar[3].expr, yyDollar[5].expr)
			yyVAL.exprlist = yyDollar[1].exprlist
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:262
		{
			yyVAL.expr = NewCompoundNode("nil")
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:266
		{
			yyVAL.expr = NewCompoundNode("false")
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:270
		{
			yyVAL.expr = NewCompoundNode("true")
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:274
		{
			yyVAL.expr = NewNumberNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:278
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:281
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:284
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:287
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:290
		{
			yyVAL.expr = NewCompoundNode("or", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:294
		{
			yyVAL.expr = NewCompoundNode("and", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 41:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:298
		{
			yyVAL.expr = NewCompoundNode("<", yyDollar[3].expr, yyDollar[1].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:302
		{
			yyVAL.expr = NewCompoundNode("<", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:306
		{
			yyVAL.expr = NewCompoundNode("<=", yyDollar[3].expr, yyDollar[1].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 44:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:310
		{
			yyVAL.expr = NewCompoundNode("<=", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:314
		{
			yyVAL.expr = NewCompoundNode("eq", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 46:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:318
		{
			yyVAL.expr = NewCompoundNode("neq", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:322
		{
			yyVAL.expr = NewCompoundNode("+", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 48:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:326
		{
			yyVAL.expr = NewCompoundNode("-", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 49:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:330
		{
			yyVAL.expr = NewCompoundNode("*", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 50:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:334
		{
			yyVAL.expr = NewCompoundNode("/", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:338
		{
			yyVAL.expr = NewCompoundNode("%", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:342
		{
			yyVAL.expr = NewCompoundNode("^", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:346
		{
			yyVAL.expr = NewCompoundNode("<<", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:350
		{
			yyVAL.expr = NewCompoundNode(">>", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:354
		{
			yyVAL.expr = NewCompoundNode("|", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:358
		{
			yyVAL.expr = NewCompoundNode("&", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 57:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:362
		{
			yyVAL.expr = NewCompoundNode("-", NewNumberNode("0"), yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 58:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:366
		{
			yyVAL.expr = NewCompoundNode("~", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 59:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:370
		{
			yyVAL.expr = NewCompoundNode("not", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 60:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:376
		{
			yyVAL.expr = NewStringNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 61:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:382
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 62:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:385
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 63:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:388
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:391
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 65:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:396
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 66:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:401
		{
			switch c, _ := yyDollar[1].expr.Value.(string); c {
			case "dup":
				switch len(yyDollar[2].exprlist.Compound) {
				case 0:
					yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(NewNumberNode("1"), NewNumberNode("1"), NewNumberNode("1")))
				case 1:
					yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(NewNumberNode("1"), yyDollar[2].exprlist.Compound[0], NewNumberNode("0")))
				default:
					p := yyDollar[2].exprlist.Compound[1]
					if p.Type != NTCompound && p.Type != NTAtom {
						yylex.(*Lexer).Error("the second argument of dup must be a closure")
					}
					yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(NewNumberNode("1"), yyDollar[2].exprlist.Compound[0], p))
				}
			case "error":
				if len(yyDollar[2].exprlist.Compound) == 0 {
					yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(NewCompoundNode("nil")))
				} else {
					yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, yyDollar[2].exprlist)
				}
			case "typeof":
				switch len(yyDollar[2].exprlist.Compound) {
				case 0:
					yylex.(*Lexer).Error("typeof takes at least 1 argument")
				case 1:
					yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].exprlist.Compound[0], NewNumberNode("255")))
				default:
					switch x, _ := yyDollar[2].exprlist.Compound[1].Value.(string); x {
					case "nil":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].exprlist.Compound[0], NewNumberNode("0")))
					case "number":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].exprlist.Compound[0], NewNumberNode("1")))
					case "string":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].exprlist.Compound[0], NewNumberNode("2")))
					case "bool":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].exprlist.Compound[0], NewNumberNode("3")))
					case "map":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].exprlist.Compound[0], NewNumberNode("4")))
					case "closure":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].exprlist.Compound[0], NewNumberNode("5")))
					case "generic":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].exprlist.Compound[0], NewNumberNode("6")))
					default:
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].exprlist.Compound[0], yyDollar[2].exprlist.Compound[1]))
					}
				}
			case "len":
				switch len(yyDollar[2].exprlist.Compound) {
				case 0:
					yylex.(*Lexer).Error("len takes 1 argument")
				default:
					yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, yyDollar[2].exprlist)
				}
			default:
				yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, yyDollar[2].exprlist)
			}
		}
	case 67:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:461
		{
			if yylex.(*Lexer).PNewLine {
				yylex.(*Lexer).TokenError(yyDollar[1].token, "ambiguous syntax (function call x new statement)")
			}
			yyVAL.exprlist = NewCompoundNode()
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:467
		{
			if yylex.(*Lexer).PNewLine {
				yylex.(*Lexer).TokenError(yyDollar[1].token, "ambiguous syntax (function call x new statement)")
			}
			yyVAL.exprlist = yyDollar[2].exprlist
		}
	case 69:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:475
		{
			yyVAL.expr = NewCompoundNode("lambda", yyDollar[2].expr, yyDollar[3].stmts)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 70:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:481
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:484
		{
			yyVAL.expr = yyDollar[2].namelist
		}
	case 72:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:489
		{
			yyVAL.expr = NewCompoundNode("map", NewCompoundNode())
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:493
		{
			yyVAL.expr = NewCompoundNode("map", yyDollar[2].exprlist)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:497
		{
			table := NewCompoundNode()
			for i, v := range yyDollar[2].exprlist.Compound {
				table.Compound = append(table.Compound,
					&Node{Type: NTNumber, Value: float64(i)},
					v)
			}
			yyVAL.expr = NewCompoundNode("map", table)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	}
	goto yystack /* stack new state and value */
}
