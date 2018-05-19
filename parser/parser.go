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
const TFalse = 57354
const TIf = 57355
const TLambda = 57356
const TList = 57357
const TNil = 57358
const TNot = 57359
const TMap = 57360
const TOr = 57361
const TReturn = 57362
const TRequire = 57363
const TSet = 57364
const TThen = 57365
const TTrue = 57366
const TWhile = 57367
const TXor = 57368
const TYield = 57369
const TEqeq = 57370
const TNeq = 57371
const TLsh = 57372
const TRsh = 57373
const TLte = 57374
const TGte = 57375
const TIdent = 57376
const TNumber = 57377
const TString = 57378
const UNARY = 57379

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
	"TIf",
	"TLambda",
	"TList",
	"TNil",
	"TNot",
	"TMap",
	"TOr",
	"TReturn",
	"TRequire",
	"TSet",
	"TThen",
	"TTrue",
	"TWhile",
	"TXor",
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
	"'['",
	"']'",
	"'}'",
	"'.'",
	"','",
	"')'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.go.y:502

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

const yyLast = 821

var yyAct = [...]int{

	60, 46, 59, 95, 138, 128, 47, 26, 126, 45,
	97, 48, 49, 100, 101, 50, 22, 25, 100, 94,
	52, 54, 55, 56, 95, 1, 34, 2, 5, 42,
	124, 4, 21, 20, 86, 23, 51, 82, 83, 84,
	65, 85, 91, 88, 78, 79, 133, 57, 65, 47,
	44, 127, 136, 33, 100, 146, 147, 145, 72, 73,
	74, 75, 76, 18, 104, 105, 106, 107, 108, 109,
	110, 111, 112, 113, 114, 115, 116, 117, 118, 119,
	120, 121, 122, 53, 65, 32, 92, 102, 125, 5,
	103, 31, 4, 24, 65, 17, 35, 132, 64, 90,
	142, 134, 72, 73, 74, 75, 76, 0, 0, 0,
	0, 123, 0, 63, 74, 75, 76, 0, 130, 131,
	65, 0, 70, 71, 78, 79, 69, 68, 139, 5,
	140, 0, 4, 80, 81, 77, 66, 67, 72, 73,
	74, 75, 76, 0, 0, 149, 0, 0, 151, 0,
	5, 0, 96, 4, 0, 0, 0, 5, 5, 0,
	4, 4, 143, 64, 0, 0, 0, 0, 0, 0,
	5, 0, 150, 4, 0, 0, 0, 5, 63, 154,
	4, 5, 0, 0, 4, 65, 0, 70, 71, 78,
	79, 69, 68, 64, 0, 0, 0, 0, 80, 81,
	77, 66, 67, 72, 73, 74, 75, 76, 63, 0,
	0, 0, 0, 0, 99, 65, 0, 70, 71, 78,
	79, 69, 68, 64, 0, 0, 0, 0, 80, 81,
	77, 66, 67, 72, 73, 74, 75, 76, 63, 0,
	0, 0, 0, 98, 0, 65, 0, 70, 71, 78,
	79, 69, 68, 64, 0, 0, 0, 0, 80, 81,
	77, 66, 67, 72, 73, 74, 75, 76, 63, 0,
	0, 144, 0, 0, 0, 65, 0, 70, 71, 78,
	79, 69, 68, 0, 0, 0, 0, 0, 80, 81,
	77, 66, 67, 72, 73, 74, 75, 76, 0, 0,
	28, 129, 39, 40, 27, 38, 41, 0, 0, 0,
	0, 0, 29, 0, 0, 64, 0, 0, 0, 61,
	0, 0, 16, 30, 43, 0, 19, 0, 0, 0,
	63, 0, 0, 36, 62, 0, 0, 65, 37, 70,
	71, 78, 79, 69, 68, 64, 58, 0, 0, 0,
	80, 81, 77, 66, 67, 72, 73, 74, 75, 76,
	63, 0, 0, 0, 153, 0, 0, 65, 0, 70,
	71, 78, 79, 69, 68, 64, 0, 0, 0, 0,
	80, 81, 77, 66, 67, 72, 73, 74, 75, 76,
	63, 0, 0, 0, 93, 0, 0, 65, 0, 70,
	71, 78, 79, 69, 68, 64, 0, 0, 0, 0,
	80, 81, 77, 66, 67, 72, 73, 74, 75, 76,
	63, 0, 0, 0, 0, 0, 0, 65, 64, 70,
	71, 78, 79, 69, 68, 0, 0, 0, 0, 0,
	80, 81, 77, 66, 67, 72, 73, 74, 75, 76,
	65, 0, 70, 71, 78, 79, 69, 68, 0, 0,
	0, 0, 0, 80, 81, 77, 66, 67, 72, 73,
	74, 75, 76, 14, 12, 13, 0, 0, 0, 152,
	0, 8, 7, 0, 0, 0, 0, 0, 10, 15,
	9, 0, 0, 6, 0, 11, 0, 0, 0, 0,
	0, 0, 16, 0, 0, 0, 19, 14, 12, 13,
	0, 0, 0, 148, 0, 8, 7, 0, 0, 3,
	0, 0, 10, 15, 9, 0, 0, 6, 0, 11,
	0, 0, 0, 0, 0, 0, 16, 0, 0, 0,
	19, 14, 12, 13, 0, 0, 0, 141, 0, 8,
	7, 0, 0, 3, 0, 0, 10, 15, 9, 0,
	0, 6, 0, 11, 0, 0, 0, 0, 0, 0,
	16, 0, 0, 0, 19, 14, 12, 13, 0, 0,
	0, 137, 0, 8, 7, 0, 0, 3, 0, 0,
	10, 15, 9, 0, 0, 6, 0, 11, 0, 0,
	0, 0, 0, 0, 16, 0, 0, 0, 19, 65,
	0, 70, 71, 78, 79, 69, 68, 0, 0, 0,
	0, 3, 80, 81, 77, 66, 67, 72, 73, 74,
	75, 76, 14, 12, 13, 0, 0, 0, 135, 0,
	8, 7, 14, 12, 13, 0, 0, 10, 15, 9,
	8, 7, 6, 0, 11, 0, 0, 10, 15, 9,
	0, 16, 6, 0, 11, 19, 0, 0, 0, 0,
	0, 16, 0, 0, 0, 19, 89, 28, 3, 39,
	40, 27, 38, 41, 0, 0, 0, 0, 3, 29,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 16,
	30, 43, 0, 19, 0, 0, 0, 0, 0, 0,
	36, 0, 0, 87, 28, 37, 39, 40, 27, 38,
	41, 0, 0, 0, 0, 0, 29, 28, 0, 39,
	40, 27, 38, 41, 0, 0, 16, 30, 43, 29,
	19, 0, 0, 0, 0, 0, 0, 36, 0, 16,
	30, 43, 37, 19, 0, 0, 0, 0, 0, 0,
	36, 0, 0, 0, 65, 37, 70, 71, 78, 79,
	69, 68, 0, 0, 0, 0, 0, 0, 0, 0,
	66, 67, 72, 73, 74, 75, 76, 14, 12, 13,
	0, 0, 0, 0, 0, 8, 7, 0, 0, 0,
	0, 0, 10, 15, 9, 0, 0, 6, 0, 11,
	0, 0, 0, 0, 0, 0, 16, 0, 0, 0,
	19,
}
var yyPact = [...]int{

	-1000, 637, -1000, -1000, -19, -21, 715, 16, 715, 15,
	715, 715, -1000, -1000, 715, 0, -1000, -1000, -1000, 715,
	715, 715, 715, 13, -1000, 288, 311, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -21, -1000, 715, 715, 715, -4,
	702, 665, -1000, -1000, -4, 371, -33, -1000, 401, 401,
	401, -1000, 94, -48, 401, 189, 159, -1000, -1000, -44,
	401, -1000, 782, 715, 715, 715, 715, 715, 715, 715,
	715, 715, 715, 715, 715, 715, 715, 715, 715, 715,
	715, 715, 22, 22, 22, -1000, -28, -1000, -3, -1000,
	-6, 249, -1000, -1000, 715, 12, -1000, -1000, -1000, -1000,
	715, -1000, 627, 44, 424, 583, 401, 14, 14, 14,
	14, 14, 14, 68, 68, 22, 22, 22, 738, 58,
	58, 738, 738, 570, -1000, -54, -1000, 715, -1000, 715,
	536, 637, -39, -1000, 401, -1000, -1000, -1000, -1000, 219,
	401, -1000, 46, 502, 715, -1000, -1000, 715, -1000, 401,
	468, 341, -1000, -1000, 637,
}
var yyPgo = [...]int{

	0, 25, 27, 100, 29, 1, 2, 99, 0, 96,
	26, 63, 95, 93, 91, 41, 85, 53,
}
var yyR1 = [...]int{

	0, 1, 1, 1, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
	3, 3, 4, 4, 4, 4, 5, 5, 6, 6,
	7, 7, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 9, 10, 10, 10, 10, 12, 11,
	13, 13, 14, 15, 15, 16, 16, 17, 17,
}
var yyR2 = [...]int{

	0, 0, 2, 2, 3, 1, 5, 7, 5, 6,
	8, 4, 1, 2, 1, 2, 1, 1, 2, 2,
	0, 5, 1, 4, 4, 3, 1, 3, 1, 3,
	3, 5, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	2, 2, 2, 1, 1, 1, 1, 3, 3, 2,
	2, 3, 4, 2, 3, 2, 3, 2, 3,
}
var yyChk = [...]int{

	-1000, -1, -2, 51, -4, -10, 25, 14, 13, 22,
	20, 27, 6, 7, 5, 21, 34, -12, -11, 38,
	52, 53, 37, 56, -13, 38, -8, 16, 12, 24,
	35, -14, -16, -17, -10, -9, 45, 50, 17, 14,
	15, 18, -4, 36, 34, -8, -5, 34, -8, -8,
	-8, 36, -8, -11, -8, -8, -8, 34, 58, -6,
	-8, 8, 23, 19, 4, 26, 42, 43, 33, 32,
	28, 29, 44, 45, 46, 47, 48, 41, 30, 31,
	39, 40, -8, -8, -8, -15, 38, 11, -6, 11,
	-7, -8, -15, 23, 52, 57, 58, 58, 54, 55,
	57, 58, -1, -2, -8, -8, -8, -8, -8, -8,
	-8, -8, -8, -8, -8, -8, -8, -8, -8, -8,
	-8, -8, -8, -1, 58, -5, 11, 57, 11, 52,
	-1, -1, -6, 34, -8, 11, 8, 11, 58, -8,
	-8, 11, -3, -1, 52, 11, 9, 10, 11, -8,
	-1, -8, 11, 23, -1,
}
var yyDef = [...]int{

	1, -2, 2, 3, 64, 5, 0, 0, 0, 0,
	12, 14, 16, 17, 0, 0, 22, 65, 66, 0,
	0, 0, 0, 0, 69, 0, 0, 32, 33, 34,
	35, 36, 37, 38, 39, 40, 0, 0, 0, 0,
	0, 0, 64, 63, 0, 0, 0, 26, 13, 15,
	18, 19, 0, 66, 4, 0, 0, 25, 70, 0,
	28, 1, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 60, 61, 62, 1, 0, 75, 0, 77,
	0, 0, 1, 1, 0, 0, 67, 68, 23, 24,
	0, 71, 0, 0, 41, 42, 43, 44, 45, 46,
	47, 48, 49, 50, 51, 52, 53, 54, 55, 56,
	57, 58, 59, 0, 73, 0, 76, 0, 78, 0,
	0, 20, 11, 27, 29, 6, 1, 72, 74, 0,
	30, 8, 0, 0, 0, 9, 1, 0, 7, 31,
	0, 0, 10, 1, 21,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 48, 40, 3,
	38, 58, 46, 44, 57, 45, 56, 47, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 51,
	43, 52, 42, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 53, 3, 54, 41, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 37, 39, 55, 50,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 49,
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
			if yyDollar[2].stmt.IsIsolatedDupCall() {
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
			if len(yyDollar[1].expr.Compound) > 0 {
				switch c, _ := yyDollar[1].expr.Compound[0].Value.(string); c {
				case "load":
					yyVAL.stmt = NewCompoundNode("store", yyDollar[1].expr.Compound[1], yyDollar[1].expr.Compound[2], yyDollar[3].expr)
				case "rload":
					yyVAL.stmt = NewCompoundNode("rstore", yyDollar[1].expr.Compound[1], yyDollar[1].expr.Compound[2], yyDollar[3].expr)
				case "safeload":
					yyVAL.stmt = NewCompoundNode("safestore", yyDollar[1].expr.Compound[1], yyDollar[1].expr.Compound[2], yyDollar[3].expr)
				}
			} else {
				yyVAL.stmt = NewCompoundNode("move", yyDollar[1].expr, yyDollar[3].expr)
			}
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:104
		{
			// if _, ok := $1.(*FuncCallExpr); !ok {
			//    yylex.(*Lexer).Error("parse error")
			// } else {
			yyVAL.stmt = yyDollar[1].expr
			// }
		}
	case 6:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:111
		{
			yyVAL.stmt = NewCompoundNode("while", yyDollar[2].expr, yyDollar[4].stmts)
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 7:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:115
		{
			yyDollar[6].stmts.Compound = append(yyDollar[6].stmts.Compound, yyDollar[4].stmt)
			yyVAL.stmt = NewCompoundNode("while", yyDollar[2].expr, yyDollar[6].stmts)
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 8:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:120
		{
			funcname := NewAtomNode(yyDollar[2].token)
			yyVAL.stmt = NewCompoundNode("chain", NewCompoundNode("set", funcname, NewCompoundNode("nil")), NewCompoundNode("move", funcname, NewCompoundNode("lambda", yyDollar[3].expr, yyDollar[4].stmts)))
		}
	case 9:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:124
		{
			yyVAL.stmt = NewCompoundNode("if", yyDollar[2].expr, yyDollar[4].stmts, NewCompoundNode())
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
			cur := yyVAL.stmt
			for _, e := range yyDollar[5].stmts.Compound {
				cur.Compound[3] = e
				cur = e
			}
		}
	case 10:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.go.y:133
		{
			yyVAL.stmt = NewCompoundNode("if", yyDollar[2].expr, yyDollar[4].stmts, NewCompoundNode())
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
			cur := yyVAL.stmt
			for _, e := range yyDollar[5].stmts.Compound {
				cur.Compound[3] = e
				cur = e
			}
			cur.Compound[3] = yyDollar[7].stmts
		}
	case 11:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:143
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
				c.Compound[0].Pos = yyDollar[1].token.Pos
				yyVAL.stmt.Compound = append(yyVAL.stmt.Compound, c)
			}
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:157
		{
			yyVAL.stmt = NewCompoundNode("ret")
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 13:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:161
		{
			yyVAL.stmt = NewCompoundNode("ret", yyDollar[2].expr)
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:165
		{
			yyVAL.stmt = NewCompoundNode("yield")
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:169
		{
			yyVAL.stmt = NewCompoundNode("yield", yyDollar[2].expr)
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:173
		{
			yyVAL.stmt = NewCompoundNode("break")
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:177
		{
			yyVAL.stmt = NewCompoundNode("continue")
			yyVAL.stmt.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 18:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:181
		{
			yyVAL.stmt = NewCompoundNode("assert", yyDollar[2].expr)
			yyVAL.stmt.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:185
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
		//line parser.go.y:207
		{
			yyVAL.stmts = NewCompoundNode()
		}
	case 21:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:210
		{
			yyVAL.stmts.Compound = append(yyVAL.stmts.Compound, NewCompoundNode("if", yyDollar[3].expr, yyDollar[5].stmts, NewCompoundNode()))
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:215
		{
			yyVAL.expr = NewAtomNode(yyDollar[1].token)
		}
	case 23:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:218
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, yyDollar[3].expr)
		}
	case 24:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:221
		{
			yyVAL.expr = NewCompoundNode("safeload", yyDollar[1].expr, yyDollar[3].expr)
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:224
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, NewStringNode(yyDollar[3].token.Str))
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:229
		{
			yyVAL.namelist = NewCompoundNode(yyDollar[1].token.Str)
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:232
		{
			yyDollar[1].namelist.Compound = append(yyDollar[1].namelist.Compound, NewAtomNode(yyDollar[3].token))
			yyVAL.namelist = yyDollar[1].namelist
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:238
		{
			yyVAL.exprlist = NewCompoundNode(yyDollar[1].expr)
		}
	case 29:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:241
		{
			yyDollar[1].exprlist.Compound = append(yyDollar[1].exprlist.Compound, yyDollar[3].expr)
			yyVAL.exprlist = yyDollar[1].exprlist
		}
	case 30:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:247
		{
			yyVAL.exprlist = NewCompoundNode(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 31:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:250
		{
			yyDollar[1].exprlist.Compound = append(yyDollar[1].exprlist.Compound, yyDollar[3].expr, yyDollar[5].expr)
			yyVAL.exprlist = yyDollar[1].exprlist
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:256
		{
			yyVAL.expr = NewCompoundNode("nil")
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:260
		{
			yyVAL.expr = NewCompoundNode("false")
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:264
		{
			yyVAL.expr = NewCompoundNode("true")
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:268
		{
			yyVAL.expr = NewNumberNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:272
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:275
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:278
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:281
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 40:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:284
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 41:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:287
		{
			yyVAL.expr = NewCompoundNode("or", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:291
		{
			yyVAL.expr = NewCompoundNode("and", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:295
		{
			yyVAL.expr = NewCompoundNode("xor", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 44:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:299
		{
			yyVAL.expr = NewCompoundNode(">", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:303
		{
			yyVAL.expr = NewCompoundNode("<", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 46:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:307
		{
			yyVAL.expr = NewCompoundNode(">=", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:311
		{
			yyVAL.expr = NewCompoundNode("<=", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 48:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:315
		{
			yyVAL.expr = NewCompoundNode("eq", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 49:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:319
		{
			yyVAL.expr = NewCompoundNode("neq", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 50:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:323
		{
			yyVAL.expr = NewCompoundNode("+", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:327
		{
			yyVAL.expr = NewCompoundNode("-", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:331
		{
			yyVAL.expr = NewCompoundNode("*", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:335
		{
			yyVAL.expr = NewCompoundNode("/", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:339
		{
			yyVAL.expr = NewCompoundNode("%", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:343
		{
			yyVAL.expr = NewCompoundNode("^", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:347
		{
			yyVAL.expr = NewCompoundNode("<<", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:351
		{
			yyVAL.expr = NewCompoundNode(">>", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:355
		{
			yyVAL.expr = NewCompoundNode("|", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:359
		{
			yyVAL.expr = NewCompoundNode("&", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 60:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:363
		{
			yyVAL.expr = NewCompoundNode("-", NewNumberNode("0"), yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 61:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:367
		{
			yyVAL.expr = NewCompoundNode("~", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 62:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:371
		{
			yyVAL.expr = NewCompoundNode("not", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 63:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:377
		{
			yyVAL.expr = NewStringNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 64:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:383
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 65:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:386
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:389
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:392
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:397
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 69:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:402
		{
			switch c, _ := yyDollar[1].expr.Value.(string); c {
			case "dup":
				switch len(yyDollar[2].exprlist.Compound) {
				case 0:
					yylex.(*Lexer).Error("dup takes at least 1 argument")
				case 1:
					yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(NewNumberNode("1"), yyDollar[2].exprlist.Compound[0], NewNumberNode("0")))
				default:
					yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(NewNumberNode("1"), yyDollar[2].exprlist.Compound[0], yyDollar[2].exprlist.Compound[1]))
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
					case "list":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].exprlist.Compound[0], NewNumberNode("4")))
					case "bytes":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].exprlist.Compound[0], NewNumberNode("5")))
					case "map":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].exprlist.Compound[0], NewNumberNode("6")))
					case "closure":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].exprlist.Compound[0], NewNumberNode("7")))
					case "generic":
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].exprlist.Compound[0], NewNumberNode("8")))
					default:
						yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, NewCompoundNode(yyDollar[2].exprlist.Compound[0], yyDollar[2].exprlist.Compound[1]))
					}
				}
			default:
				yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, yyDollar[2].exprlist)
			}
		}
	case 70:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:455
		{
			if yylex.(*Lexer).PNewLine {
				yylex.(*Lexer).TokenError(yyDollar[1].token, "ambiguous syntax (function call x new statement)")
			}
			yyVAL.exprlist = NewCompoundNode()
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:461
		{
			if yylex.(*Lexer).PNewLine {
				yylex.(*Lexer).TokenError(yyDollar[1].token, "ambiguous syntax (function call x new statement)")
			}
			yyVAL.exprlist = yyDollar[2].exprlist
		}
	case 72:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:469
		{
			yyVAL.expr = NewCompoundNode("lambda", yyDollar[2].expr, yyDollar[3].stmts)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 73:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:475
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:478
		{
			yyVAL.expr = yyDollar[2].namelist
		}
	case 75:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:483
		{
			yyVAL.expr = NewCompoundNode("list", NewCompoundNode())
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:487
		{
			yyVAL.expr = NewCompoundNode("list", yyDollar[2].exprlist)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 77:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:493
		{
			yyVAL.expr = NewCompoundNode("map", NewCompoundNode())
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 78:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:497
		{
			yyVAL.expr = NewCompoundNode("map", yyDollar[2].exprlist)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	}
	goto yystack /* stack new state and value */
}
