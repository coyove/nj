//line parser.go.y:2
package parser

import __yyfmt__ "fmt"

//line parser.go.y:2
import (
	"bytes"
	"io/ioutil"
	"path/filepath"
)

//line parser.go.y:29
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
const TFor = 57354
const TIf = 57355
const TLambda = 57356
const TList = 57357
const TNil = 57358
const TNot = 57359
const TOr = 57360
const TReturn = 57361
const TRequire = 57362
const TSet = 57363
const TThen = 57364
const TYield = 57365
const TEqeq = 57366
const TNeq = 57367
const TLsh = 57368
const TRsh = 57369
const TLte = 57370
const TGte = 57371
const TIdent = 57372
const TNumber = 57373
const TString = 57374
const UNARY = 57375

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

//line parser.go.y:519

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

const yyLast = 717

var yyAct = [...]int{

	54, 53, 88, 41, 117, 90, 129, 25, 42, 40,
	92, 43, 44, 92, 93, 45, 87, 88, 119, 2,
	47, 49, 50, 30, 24, 5, 1, 20, 37, 118,
	4, 115, 79, 75, 76, 77, 46, 84, 83, 78,
	21, 124, 22, 65, 66, 67, 68, 69, 67, 68,
	69, 51, 42, 39, 137, 138, 136, 18, 96, 97,
	98, 99, 100, 101, 102, 103, 104, 105, 106, 107,
	108, 109, 110, 111, 112, 113, 95, 48, 81, 85,
	5, 127, 94, 116, 35, 4, 26, 34, 29, 123,
	28, 23, 17, 125, 31, 82, 133, 0, 0, 0,
	16, 27, 38, 36, 19, 114, 0, 0, 0, 0,
	0, 32, 121, 122, 71, 72, 33, 0, 5, 130,
	125, 131, 0, 4, 80, 0, 0, 0, 65, 66,
	67, 68, 69, 0, 0, 0, 140, 0, 5, 142,
	0, 0, 0, 4, 0, 5, 5, 0, 0, 0,
	4, 4, 0, 0, 134, 0, 0, 0, 5, 58,
	0, 0, 0, 4, 141, 5, 0, 0, 0, 5,
	4, 145, 0, 57, 4, 0, 0, 0, 0, 63,
	64, 71, 72, 62, 61, 0, 0, 0, 0, 58,
	73, 74, 70, 59, 60, 65, 66, 67, 68, 69,
	0, 0, 0, 57, 0, 0, 0, 0, 89, 63,
	64, 71, 72, 62, 61, 0, 0, 0, 0, 0,
	73, 74, 70, 59, 60, 65, 66, 67, 68, 69,
	58, 0, 0, 0, 55, 0, 91, 0, 0, 0,
	0, 0, 0, 0, 57, 0, 0, 0, 0, 0,
	63, 64, 71, 72, 62, 61, 0, 58, 0, 0,
	0, 73, 74, 70, 59, 60, 65, 66, 67, 68,
	69, 57, 0, 0, 0, 56, 0, 63, 64, 71,
	72, 62, 61, 0, 58, 0, 0, 0, 73, 74,
	70, 59, 60, 65, 66, 67, 68, 69, 57, 0,
	0, 135, 0, 0, 63, 64, 71, 72, 62, 61,
	0, 0, 0, 0, 0, 73, 74, 70, 59, 60,
	65, 66, 67, 68, 69, 0, 0, 35, 120, 26,
	34, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 16, 27, 38, 36, 19, 0, 58,
	0, 0, 0, 0, 32, 0, 0, 0, 0, 33,
	0, 0, 0, 57, 0, 0, 52, 144, 0, 63,
	64, 71, 72, 62, 61, 0, 58, 0, 0, 0,
	73, 74, 70, 59, 60, 65, 66, 67, 68, 69,
	57, 0, 0, 0, 86, 0, 63, 64, 71, 72,
	62, 61, 0, 58, 0, 0, 0, 73, 74, 70,
	59, 60, 65, 66, 67, 68, 69, 57, 0, 0,
	0, 0, 0, 63, 64, 71, 72, 62, 61, 0,
	0, 0, 0, 0, 73, 74, 70, 59, 60, 65,
	66, 67, 68, 69, 14, 12, 13, 0, 0, 0,
	143, 6, 8, 7, 0, 0, 0, 0, 10, 15,
	9, 0, 11, 0, 0, 0, 0, 0, 0, 16,
	0, 0, 0, 19, 14, 12, 13, 0, 0, 0,
	139, 6, 8, 7, 0, 0, 3, 0, 10, 15,
	9, 0, 11, 0, 0, 0, 0, 0, 0, 16,
	0, 0, 0, 19, 14, 12, 13, 0, 0, 0,
	132, 6, 8, 7, 0, 0, 3, 0, 10, 15,
	9, 0, 11, 0, 0, 0, 0, 0, 0, 16,
	0, 0, 0, 19, 14, 12, 13, 0, 0, 0,
	128, 6, 8, 7, 0, 0, 3, 58, 10, 15,
	9, 0, 11, 0, 0, 0, 0, 0, 0, 16,
	0, 0, 0, 19, 0, 0, 0, 63, 64, 71,
	72, 62, 61, 0, 0, 0, 3, 0, 73, 74,
	70, 59, 60, 65, 66, 67, 68, 69, 14, 12,
	13, 0, 0, 0, 126, 6, 8, 7, 0, 0,
	0, 0, 10, 15, 9, 0, 11, 14, 12, 13,
	0, 0, 0, 16, 6, 8, 7, 19, 0, 0,
	0, 10, 15, 9, 0, 11, 0, 0, 0, 0,
	3, 0, 16, 0, 0, 0, 19, 0, 0, 63,
	64, 71, 72, 62, 61, 0, 0, 0, 0, 3,
	73, 74, 70, 59, 60, 65, 66, 67, 68, 69,
	63, 64, 71, 72, 62, 61, 35, 0, 26, 34,
	0, 0, 0, 0, 59, 60, 65, 66, 67, 68,
	69, 0, 16, 27, 38, 36, 19, 14, 12, 13,
	0, 0, 0, 32, 6, 8, 7, 0, 33, 0,
	0, 10, 15, 9, 0, 11, 0, 0, 0, 0,
	0, 0, 16, 0, 0, 0, 19,
}
var yyPact = [...]int{

	-1000, 602, -1000, -1000, -21, -10, 652, 23, 652, 22,
	652, 652, -1000, -1000, 652, 4, -1000, -1000, -1000, 652,
	652, 652, 21, -1000, 313, 226, -1000, -1000, -1000, -1000,
	-10, -1000, 652, 652, 652, -2, 70, -1000, -1000, -2,
	372, -32, -1000, 399, 399, 399, -1000, 155, -48, 399,
	185, -1000, -1000, -39, 399, -1000, 682, 652, 652, 652,
	652, 652, 652, 652, 652, 652, 652, 652, 652, 652,
	652, 652, 652, 652, 652, -1000, -1000, -1000, -1000, -22,
	-1000, -50, -20, -31, 280, -1000, -1000, 652, 11, -1000,
	-1000, -1000, 652, -1000, 583, 73, 543, 615, 88, 88,
	88, 88, 88, 88, 6, 6, -1000, -1000, -1000, 636,
	3, 3, 636, 636, 529, -1000, -47, -1000, 652, 652,
	652, 499, 602, -36, -1000, 399, -1000, -1000, -1000, -1000,
	253, 399, -1000, 45, 469, 652, -1000, -1000, 652, -1000,
	399, 439, 345, -1000, -1000, 602,
}
var yyPgo = [...]int{

	0, 26, 19, 96, 28, 3, 1, 95, 0, 94,
	23, 57, 92, 91, 90, 39, 88, 78,
}
var yyR1 = [...]int{

	0, 1, 1, 1, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
	3, 3, 4, 4, 4, 5, 5, 6, 6, 7,
	7, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 9, 10,
	10, 10, 10, 12, 11, 13, 13, 14, 15, 15,
	16, 16, 17, 17, 17, 17,
}
var yyR2 = [...]int{

	0, 0, 2, 2, 3, 1, 5, 7, 5, 6,
	8, 4, 1, 2, 1, 2, 1, 1, 2, 2,
	0, 5, 1, 4, 3, 1, 3, 1, 3, 3,
	5, 1, 1, 1, 1, 1, 1, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 2, 2, 2, 1, 1,
	1, 1, 3, 3, 2, 2, 3, 4, 2, 3,
	2, 3, 1, 2, 1, 2,
}
var yyChk = [...]int{

	-1000, -1, -2, 47, -4, -10, 12, 14, 13, 21,
	19, 23, 6, 7, 5, 20, 30, -12, -11, 34,
	48, 50, 52, -13, 34, -8, 16, 31, -14, -16,
	-10, -9, 41, 46, 17, 14, 33, -4, 32, 30,
	-8, -5, 30, -8, -8, -8, 32, -8, -11, -8,
	-8, 30, 53, -6, -8, 8, 49, 18, 4, 38,
	39, 29, 28, 24, 25, 40, 41, 42, 43, 44,
	37, 26, 27, 35, 36, -8, -8, -8, -15, 34,
	54, -17, -7, -6, -8, -15, 22, 48, 49, 53,
	53, 51, 49, 53, -1, -2, -8, -8, -8, -8,
	-8, -8, -8, -8, -8, -8, -8, -8, -8, -8,
	-8, -8, -8, -8, -1, 53, -5, 54, 49, 49,
	48, -1, -1, -6, 30, -8, 11, 8, 11, 53,
	-8, -8, 11, -3, -1, 48, 11, 9, 10, 11,
	-8, -1, -8, 11, 22, -1,
}
var yyDef = [...]int{

	1, -2, 2, 3, 59, 5, 0, 0, 0, 0,
	12, 14, 16, 17, 0, 0, 22, 60, 61, 0,
	0, 0, 0, 64, 0, 0, 31, 32, 33, 34,
	35, 36, 0, 0, 0, 0, 0, 59, 58, 0,
	0, 0, 25, 13, 15, 18, 19, 0, 61, 4,
	0, 24, 65, 0, 27, 1, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 55, 56, 57, 1, 0,
	70, 0, 72, 74, 27, 1, 1, 0, 0, 62,
	63, 23, 0, 66, 0, 0, 37, 38, 39, 40,
	41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
	51, 52, 53, 54, 0, 68, 0, 71, 73, 75,
	0, 0, 20, 11, 26, 28, 6, 1, 67, 69,
	0, 29, 8, 0, 0, 0, 9, 1, 0, 7,
	30, 0, 0, 10, 1, 21,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 44, 36, 3,
	34, 53, 42, 40, 49, 41, 52, 43, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 47,
	39, 48, 38, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 50, 3, 51, 37, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 33, 35, 54, 46,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 45,
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
		//line parser.go.y:64
		{
			yyVAL.stmts = NewCompoundNode("chain")
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.stmts
			}
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:70
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
		//line parser.go.y:80
		{
			yyVAL.stmts = yyDollar[1].stmts
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.stmts
			}
		}
	case 4:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:88
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
		//line parser.go.y:109
		{
			// if _, ok := $1.(*FuncCallExpr); !ok {
			//    yylex.(*Lexer).Error("parse error")
			// } else {
			yyVAL.stmt = yyDollar[1].expr
			// }
		}
	case 6:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:116
		{
			yyVAL.stmt = NewCompoundNode("for", yyDollar[2].expr, NewCompoundNode(), yyDollar[4].stmts)
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 7:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:120
		{
			yyVAL.stmt = NewCompoundNode("for", yyDollar[2].expr, NewCompoundNode("chain", yyDollar[4].stmt), yyDollar[6].stmts)
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 8:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:124
		{
			funcname := NewAtomNode(yyDollar[2].token)
			yyVAL.stmt = NewCompoundNode("chain", NewCompoundNode("set", funcname, NewNilNode()), NewCompoundNode("move", funcname, NewCompoundNode("lambda", yyDollar[3].expr, yyDollar[4].stmts)))
		}
	case 9:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:128
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
		//line parser.go.y:137
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
		//line parser.go.y:147
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
		//line parser.go.y:162
		{
			yyVAL.stmt = NewCompoundNode("ret")
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 13:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:166
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
		//line parser.go.y:175
		{
			yyVAL.stmt = NewCompoundNode("yield")
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:179
		{
			yyVAL.stmt = NewCompoundNode("yield", yyDollar[2].expr)
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:183
		{
			yyVAL.stmt = NewCompoundNode("break")
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:187
		{
			yyVAL.stmt = NewCompoundNode("continue")
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 18:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:191
		{
			yyVAL.stmt = NewCompoundNode("assert", yyDollar[2].expr)
			yyVAL.stmt.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:195
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
		//line parser.go.y:217
		{
			yyVAL.stmts = NewCompoundNode()
		}
	case 21:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:220
		{
			yyVAL.stmts.Compound = append(yyDollar[1].stmts.Compound, NewCompoundNode("if", yyDollar[3].expr, yyDollar[5].stmts, NewCompoundNode()))
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:225
		{
			yyVAL.expr = NewAtomNode(yyDollar[1].token)
		}
	case 23:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:228
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, yyDollar[3].expr)
		}
	case 24:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:231
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, NewStringNode(yyDollar[3].token.Str))
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:236
		{
			yyVAL.namelist = NewCompoundNode(yyDollar[1].token.Str)
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:239
		{
			yyDollar[1].namelist.Compound = append(yyDollar[1].namelist.Compound, NewAtomNode(yyDollar[3].token))
			yyVAL.namelist = yyDollar[1].namelist
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:245
		{
			yyVAL.exprlist = NewCompoundNode(yyDollar[1].expr)
		}
	case 28:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:248
		{
			yyDollar[1].exprlist.Compound = append(yyDollar[1].exprlist.Compound, yyDollar[3].expr)
			yyVAL.exprlist = yyDollar[1].exprlist
		}
	case 29:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:254
		{
			yyVAL.exprlist = NewCompoundNode(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 30:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:257
		{
			yyDollar[1].exprlist.Compound = append(yyDollar[1].exprlist.Compound, yyDollar[3].expr, yyDollar[5].expr)
			yyVAL.exprlist = yyDollar[1].exprlist
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:263
		{
			yyVAL.expr = NewNilNode()
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:267
		{
			yyVAL.expr = NewNumberNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:271
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:274
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:277
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:280
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 37:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:283
		{
			yyVAL.expr = NewCompoundNode("or", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:287
		{
			yyVAL.expr = NewCompoundNode("and", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:291
		{
			yyVAL.expr = NewCompoundNode("<", yyDollar[3].expr, yyDollar[1].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:295
		{
			yyVAL.expr = NewCompoundNode("<", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 41:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:299
		{
			yyVAL.expr = NewCompoundNode("<=", yyDollar[3].expr, yyDollar[1].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:303
		{
			yyVAL.expr = NewCompoundNode("<=", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:307
		{
			yyVAL.expr = NewCompoundNode("eq", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 44:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:311
		{
			yyVAL.expr = NewCompoundNode("neq", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:315
		{
			yyVAL.expr = NewCompoundNode("+", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 46:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:319
		{
			yyVAL.expr = NewCompoundNode("-", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:323
		{
			yyVAL.expr = NewCompoundNode("*", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 48:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:327
		{
			yyVAL.expr = NewCompoundNode("/", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 49:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:331
		{
			yyVAL.expr = NewCompoundNode("%", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 50:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:335
		{
			yyVAL.expr = NewCompoundNode("^", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:339
		{
			yyVAL.expr = NewCompoundNode("<<", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:343
		{
			yyVAL.expr = NewCompoundNode(">>", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:347
		{
			yyVAL.expr = NewCompoundNode("|", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:351
		{
			yyVAL.expr = NewCompoundNode("&", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 55:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:355
		{
			yyVAL.expr = NewCompoundNode("-", NewNumberNode("0"), yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 56:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:359
		{
			yyVAL.expr = NewCompoundNode("~", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 57:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:363
		{
			yyVAL.expr = NewCompoundNode("not", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 58:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:369
		{
			yyVAL.expr = NewStringNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 59:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:375
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 60:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:378
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 61:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:381
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 62:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:384
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 63:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:389
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 64:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:394
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
					yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(NewNilNode()))
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
					case "map":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].exprlist.Compound[0], NewNumberNode("3")))
					case "closure":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].exprlist.Compound[0], NewNumberNode("4")))
					case "generic":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].exprlist.Compound[0], NewNumberNode("5")))
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
	case 65:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:452
		{
			if yylex.(*Lexer).PNewLine {
				yylex.(*Lexer).TokenError(yyDollar[1].token, "ambiguous syntax (function call x new statement)")
			}
			yyVAL.exprlist = NewCompoundNode()
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:458
		{
			if yylex.(*Lexer).PNewLine {
				yylex.(*Lexer).TokenError(yyDollar[1].token, "ambiguous syntax (function call x new statement)")
			}
			yyVAL.exprlist = yyDollar[2].exprlist
		}
	case 67:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:466
		{
			yyVAL.expr = NewCompoundNode("lambda", yyDollar[2].expr, yyDollar[3].stmts)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 68:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:472
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:475
		{
			yyVAL.expr = yyDollar[2].namelist
		}
	case 70:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:480
		{
			yyVAL.expr = NewCompoundNode("map", NewCompoundNode())
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:484
		{
			yyVAL.expr = yyDollar[2].expr
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:490
		{
			yyVAL.expr = NewCompoundNode("map", yyDollar[1].exprlist)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].exprlist.Pos
		}
	case 73:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:494
		{
			yyVAL.expr = NewCompoundNode("map", yyDollar[1].exprlist)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].exprlist.Pos
		}
	case 74:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:498
		{
			table := NewCompoundNode()
			for i, v := range yyDollar[1].exprlist.Compound {
				table.Compound = append(table.Compound,
					&Node{Type: NTNumber, Value: float64(i)},
					v)
			}
			yyVAL.expr = NewCompoundNode("map", table)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].exprlist.Pos
		}
	case 75:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:508
		{
			table := NewCompoundNode()
			for i, v := range yyDollar[1].exprlist.Compound {
				table.Compound = append(table.Compound,
					&Node{Type: NTNumber, Value: float64(i)},
					v)
			}
			yyVAL.expr = NewCompoundNode("map", table)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].exprlist.Pos
		}
	}
	goto yystack /* stack new state and value */
}
