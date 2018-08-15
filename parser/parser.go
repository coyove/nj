//line parser.go.y:1
package parser

import __yyfmt__ "fmt"

//line parser.go.y:3
import (
	"path/filepath"
)

//line parser.go.y:55
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
const TAddAdd = 57359
const TMinMin = 57360
const TEqeq = 57361
const TNeq = 57362
const TLsh = 57363
const TRsh = 57364
const TLte = 57365
const TGte = 57366
const TIdent = 57367
const TNumber = 57368
const TString = 57369
const TOr = 57370
const TAnd = 57371
const UNARY = 57372

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
	"TAddAdd",
	"TMinMin",
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
	"')'",
	"'!'",
	"','",
	"':'",
	"'='",
	"'['",
	"']'",
	"'.'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.go.y:452

var typesLookup = map[string]string{
	"nil": "0", "number": "1", "string": "2", "map": "4", "closure": "6", "generic": "7",
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 60,
	47, 40,
	-2, 58,
}

const yyPrivate = 57344

const yyLast = 550

var yyAct = [...]int{

	79, 177, 27, 57, 49, 17, 51, 181, 2, 39,
	27, 7, 26, 27, 27, 23, 50, 48, 29, 35,
	47, 55, 65, 174, 58, 59, 62, 19, 45, 18,
	84, 77, 85, 86, 87, 88, 54, 169, 72, 54,
	168, 78, 27, 195, 54, 185, 72, 122, 189, 73,
	190, 160, 54, 93, 81, 27, 82, 180, 186, 172,
	54, 171, 127, 126, 111, 69, 128, 27, 162, 27,
	54, 156, 121, 54, 27, 155, 54, 154, 115, 138,
	120, 54, 27, 53, 108, 64, 54, 52, 110, 124,
	125, 123, 94, 130, 117, 103, 104, 135, 92, 136,
	129, 114, 142, 143, 144, 145, 150, 151, 152, 80,
	148, 149, 153, 140, 141, 139, 134, 27, 146, 147,
	27, 137, 89, 17, 112, 105, 106, 107, 158, 72,
	67, 66, 27, 100, 99, 164, 63, 68, 61, 165,
	166, 167, 63, 170, 67, 109, 97, 98, 60, 173,
	163, 132, 71, 101, 102, 95, 96, 27, 27, 90,
	91, 192, 75, 27, 183, 24, 17, 119, 1, 178,
	27, 182, 176, 42, 6, 83, 37, 5, 4, 16,
	27, 187, 157, 116, 166, 15, 8, 30, 36, 40,
	70, 27, 76, 27, 183, 176, 17, 78, 193, 194,
	197, 182, 113, 191, 0, 0, 176, 0, 12, 10,
	11, 0, 20, 22, 196, 21, 43, 13, 14, 25,
	9, 0, 0, 0, 0, 0, 0, 0, 0, 38,
	44, 46, 28, 41, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 31, 0, 0, 0, 0, 32, 34,
	161, 3, 0, 33, 12, 10, 11, 0, 20, 22,
	0, 21, 43, 13, 14, 25, 9, 0, 0, 0,
	0, 0, 0, 0, 0, 38, 44, 46, 184, 41,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 31,
	0, 0, 0, 0, 32, 34, 0, 3, 0, 33,
	12, 10, 11, 0, 20, 22, 0, 21, 43, 13,
	14, 25, 9, 0, 22, 0, 0, 43, 0, 56,
	0, 38, 44, 46, 28, 41, 0, 0, 0, 0,
	38, 44, 46, 28, 41, 31, 0, 0, 0, 0,
	32, 34, 0, 3, 31, 33, 0, 0, 22, 32,
	34, 43, 0, 56, 33, 0, 0, 22, 0, 188,
	43, 0, 56, 25, 38, 44, 46, 28, 41, 0,
	0, 0, 0, 38, 44, 46, 28, 41, 31, 0,
	0, 0, 0, 32, 34, 0, 22, 31, 33, 43,
	131, 56, 32, 34, 0, 118, 0, 33, 0, 0,
	0, 0, 38, 44, 46, 67, 41, 0, 0, 0,
	0, 0, 22, 0, 0, 43, 31, 56, 0, 0,
	22, 32, 34, 43, 175, 56, 33, 0, 38, 44,
	46, 28, 41, 0, 0, 0, 38, 44, 46, 28,
	41, 0, 31, 0, 0, 0, 0, 32, 34, 0,
	31, 179, 33, 0, 22, 32, 34, 43, 159, 56,
	33, 0, 22, 0, 0, 43, 0, 56, 0, 0,
	38, 44, 46, 0, 41, 0, 0, 0, 38, 44,
	46, 28, 41, 22, 31, 0, 43, 0, 56, 32,
	34, 74, 31, 0, 33, 0, 0, 32, 34, 38,
	44, 46, 33, 41, 0, 0, 0, 0, 0, 22,
	0, 0, 43, 31, 56, 0, 0, 0, 32, 34,
	0, 0, 133, 33, 0, 38, 44, 46, 0, 41,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 31,
	0, 0, 0, 0, 32, 34, 0, 0, 0, 33,
}
var yyPact = [...]int{

	-1000, 296, -1000, -1000, 40, -1000, -1000, 36, -1000, 453,
	-1000, -1000, 453, 453, 121, -1000, -1000, 113, -1000, 38,
	102, 108, 16, -1000, -1000, 127, 98, -3, 445, 77,
	1, 500, 500, 500, 500, 89, -1000, -1000, 142, 63,
	-1000, 453, -1000, -1000, -1000, 58, -1000, 136, 110, 132,
	57, 85, -1000, -1000, 453, 26, 118, 107, 26, 26,
	-1000, 107, 116, 76, -1000, -1000, 348, -1000, 453, -1000,
	22, -5, 500, 500, -1000, 44, 13, 12, 15, -1000,
	500, 339, 126, -1000, 474, -1000, -1000, -1000, -1000, 500,
	-1000, -1000, 500, 31, 500, 500, 500, 500, 500, 500,
	500, 500, 500, 500, 500, 500, 500, 500, -1000, -1000,
	116, -1000, -1000, 27, -1000, 23, 411, 4, -1000, 204,
	20, 125, 500, 77, -1000, -1000, 500, 500, 500, 89,
	-14, 453, -1000, -1000, 11, 98, 63, 58, -1000, 136,
	110, 110, 132, 132, 132, 132, 57, 57, 85, 85,
	-1000, -1000, -1000, -1000, 124, -1000, 377, 403, 10, -1000,
	-1000, -1000, 250, -7, 98, 7, 98, 98, -1000, 305,
	-6, 500, -1000, -1000, -1000, -1000, -1000, -1000, 2, 377,
	-1000, 154, -1000, -1000, 445, 500, 500, -11, -1000, -1000,
	377, -1000, 250, 98, 98, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 167, 1, 8, 7, 202, 31, 192, 190, 11,
	189, 188, 187, 0, 6, 16, 4, 17, 20, 28,
	9, 19, 18, 12, 15, 29, 27, 186, 185, 183,
	182, 179, 23, 178, 177, 3, 176, 175, 26, 174,
	173, 165, 162,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 27, 3, 3, 3, 3, 3,
	3, 4, 4, 32, 32, 32, 34, 34, 26, 29,
	29, 30, 30, 28, 28, 28, 28, 31, 31, 35,
	35, 39, 40, 33, 33, 33, 33, 33, 33, 33,
	33, 5, 5, 6, 6, 7, 7, 8, 8, 8,
	8, 10, 11, 11, 11, 11, 11, 11, 11, 12,
	12, 12, 12, 12, 12, 12, 12, 12, 13, 13,
	13, 13, 13, 14, 14, 14, 14, 15, 15, 15,
	16, 16, 16, 17, 17, 17, 17, 17, 18, 18,
	18, 19, 19, 20, 20, 21, 21, 22, 22, 23,
	23, 24, 24, 25, 25, 9, 9, 36, 37, 37,
	38, 38, 41, 41, 42, 42, 42, 42,
}
var yyR2 = [...]int{

	0, 0, 2, 3, 2, 1, 2, 1, 1, 2,
	1, 1, 1, 1, 1, 1, 1, 1, 2, 2,
	1, 2, 1, 2, 5, 7, 6, 5, 7, 1,
	2, 4, 3, 1, 2, 1, 1, 2, 1, 2,
	2, 1, 3, 1, 3, 3, 5, 1, 3, 5,
	3, 1, 1, 1, 3, 1, 1, 1, 2, 1,
	4, 6, 5, 5, 3, 1, 2, 2, 1, 2,
	2, 2, 2, 1, 3, 3, 3, 1, 3, 3,
	1, 3, 3, 1, 3, 3, 3, 3, 1, 3,
	3, 1, 3, 1, 3, 1, 3, 1, 3, 1,
	3, 1, 3, 1, 1, 1, 3, 2, 2, 3,
	2, 3, 2, 3, 1, 2, 1, 2,
}
var yyChk = [...]int{

	-1000, -1, -3, 47, -33, -34, -39, -9, -27, 16,
	5, 6, 4, 13, 14, -28, -31, -35, -25, -26,
	8, 11, 9, -24, -41, 15, -23, -13, 28, -22,
	-12, 39, 44, 49, 45, -21, -11, -36, 25, -20,
	-10, 29, -40, 12, 26, -19, 27, -18, -17, -16,
	-15, -14, 47, 47, 50, -9, 14, -35, -9, -9,
	27, 25, -38, 29, 47, -2, 29, 28, 29, 49,
	-8, 25, 31, 52, 46, -42, -7, -6, -23, -13,
	32, 53, 55, -37, 29, -13, -13, -13, -13, 33,
	17, 18, 35, -9, 34, 19, 20, 36, 37, 24,
	23, 21, 22, 38, 39, 40, 41, 42, -25, 27,
	-38, -2, 48, -5, 25, -9, -29, -26, 47, -1,
	-9, 50, 52, -22, -24, 46, 50, 50, 51, -21,
	-9, 51, 25, 48, -6, -23, -20, -19, 48, -18,
	-17, -17, -16, -16, -16, -16, -15, -15, -14, -14,
	-13, -13, -13, -2, 50, 48, 48, -30, -9, 47,
	47, 46, 48, 25, -23, -23, -23, -23, 54, 51,
	-9, 50, 48, 25, -32, 47, -24, -2, -9, 48,
	47, -4, -3, -2, 28, 52, 51, -9, 54, 54,
	48, -32, 7, -23, -23, 54, -32, -4,
}
var yyDef = [...]int{

	1, -2, 2, 5, 0, 7, 8, 0, 10, 33,
	35, 36, 0, 38, 0, 16, 17, 0, 105, 0,
	0, 0, 29, 103, 104, 0, 101, 73, 0, 99,
	68, 0, 0, 0, 0, 97, 59, 65, 52, 95,
	53, 0, 55, 56, 57, 93, 51, 91, 88, 83,
	80, 77, 6, 9, 0, 34, 0, 0, 37, 39,
	-2, 0, 0, 0, 4, 23, 0, 1, 0, 30,
	18, 47, 0, 0, 112, 0, 114, 116, 43, 73,
	0, 0, 0, 107, 0, 69, 70, 71, 72, 0,
	66, 67, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 106, 58,
	0, 32, 110, 0, 41, 0, 0, 0, 20, 0,
	0, 0, 0, 100, 102, 113, 115, 117, 0, 98,
	0, 0, 64, 108, 0, 43, 96, 94, 54, 92,
	89, 90, 84, 85, 86, 87, 81, 82, 78, 79,
	74, 75, 76, 31, 0, 111, 0, 0, 0, 22,
	19, 3, 0, 50, 48, 0, 44, 45, 60, 0,
	0, 0, 109, 42, 24, 13, 14, 15, 0, 0,
	21, 27, 11, 12, 1, 0, 0, 0, 62, 63,
	0, 26, 0, 49, 46, 61, 25, 28,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 49, 3, 45, 3, 42, 34, 3,
	29, 48, 40, 38, 50, 39, 55, 41, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 51, 47,
	37, 52, 36, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 30, 3, 3, 3, 3, 3,
	3, 53, 3, 54, 35, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 28, 33, 46, 44,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 31, 32, 43,
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
		//line parser.go.y:86
		{
			yyVAL.expr = NewCompoundNode("chain")
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:92
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, yyDollar[2].expr)
			yyVAL.expr = yyDollar[1].expr
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 3:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:101
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:104
		{
			if yyDollar[1].expr.isIsolatedDupCall() {
				yyDollar[1].expr.Compound[2].Compound[0] = NewNumberNode("0")
			}
			yyVAL.expr = yyDollar[1].expr
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:112
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:113
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:114
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:115
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:116
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:117
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:120
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:121
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:124
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:125
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:126
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:129
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:130
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 18:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:133
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:138
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 20:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:139
		{
			yyVAL.expr = NewCompoundNode("chain")
		}
	case 21:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:142
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:143
		{
			yyVAL.expr = NewNumberNode("1")
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:146
		{
			yyVAL.expr = NewCompoundNode("for", NewNumberNode("1"), NewCompoundNode(), yyDollar[2].expr).setPos0(yyDollar[1].token.Pos)
		}
	case 24:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:149
		{
			yyVAL.expr = NewCompoundNode("for", yyDollar[3].expr, NewCompoundNode(), yyDollar[5].expr).setPos0(yyDollar[1].token.Pos)
		}
	case 25:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:152
		{
			yyVAL.expr = yyDollar[3].expr
			yyVAL.expr.Compound = append(yyVAL.expr.Compound, NewCompoundNode("for", yyDollar[4].expr, NewCompoundNode("chain", yyDollar[5].expr), yyDollar[7].expr))
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 26:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:157
		{
			yyVAL.expr = yyDollar[3].expr
			yyVAL.expr.Compound = append(yyVAL.expr.Compound, NewCompoundNode("for", yyDollar[4].expr, NewCompoundNode(), yyDollar[6].expr))
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 27:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:164
		{
			yyVAL.expr = NewCompoundNode("if", yyDollar[3].expr, yyDollar[5].expr, NewCompoundNode())
		}
	case 28:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:165
		{
			yyVAL.expr = NewCompoundNode("if", yyDollar[3].expr, yyDollar[5].expr, yyDollar[7].expr)
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:168
		{
			yyVAL.str = "func"
		}
	case 30:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:169
		{
			yyVAL.str = "safefunc"
		}
	case 31:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:172
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
	case 32:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:184
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].str, "<a>", yyDollar[2].expr, yyDollar[3].expr).setPos0(yyDollar[2].expr.Pos)
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:187
		{
			yyVAL.expr = NewCompoundNode("yield").setPos0(yyDollar[1].token.Pos)
		}
	case 34:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:188
		{
			yyVAL.expr = NewCompoundNode("yield", yyDollar[2].expr).setPos0(yyDollar[1].token.Pos)
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:189
		{
			yyVAL.expr = NewCompoundNode("break").setPos0(yyDollar[1].token.Pos)
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:190
		{
			yyVAL.expr = NewCompoundNode("continue").setPos0(yyDollar[1].token.Pos)
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:191
		{
			yyVAL.expr = NewCompoundNode("assert", yyDollar[2].expr).setPos0(yyDollar[1].token.Pos)
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:192
		{
			yyVAL.expr = NewCompoundNode("ret").setPos0(yyDollar[1].token.Pos)
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:193
		{
			if yyDollar[2].expr.isIsolatedDupCall() {
				if h, _ := yyDollar[2].expr.Compound[2].Compound[2].Value.(float64); h == 1 {
					yyDollar[2].expr.Compound[2].Compound[2] = NewNumberNode("2")
				}
			}
			yyVAL.expr = NewCompoundNode("ret", yyDollar[2].expr).setPos0(yyDollar[1].token.Pos)
		}
	case 40:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:201
		{
			path := filepath.Join(filepath.Dir(yyDollar[1].token.Pos.Source), yyDollar[2].token.Str)
			yyVAL.expr = yylex.(*Lexer).loadFile(path)
		}
	case 41:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:207
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].token.Str)
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:210
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, NewAtomNode(yyDollar[3].token))
			yyVAL.expr = yyDollar[1].expr
		}
	case 43:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:216
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].expr)
		}
	case 44:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:219
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, yyDollar[3].expr)
			yyVAL.expr = yyDollar[1].expr
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:225
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 46:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:228
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, yyDollar[3].expr, yyDollar[5].expr)
			yyVAL.expr = yyDollar[1].expr
		}
	case 47:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:234
		{
			yyVAL.expr = NewCompoundNode("chain", NewCompoundNode("set", NewAtomNode(yyDollar[1].token), NewNilNode()))
			yyVAL.expr.Compound[1].Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 48:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:238
		{
			yyVAL.expr = NewCompoundNode("chain", NewCompoundNode("set", NewAtomNode(yyDollar[1].token), yyDollar[3].expr))
			yyVAL.expr.Compound[1].Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 49:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:242
		{
			x := NewCompoundNode("set", NewAtomNode(yyDollar[3].token), yyDollar[5].expr).setPos0(yyDollar[1].expr.Pos)
			yyDollar[1].expr.Compound = append(yyVAL.expr.Compound, x)
			yyVAL.expr = yyDollar[1].expr
		}
	case 50:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:247
		{
			x := NewCompoundNode("set", NewAtomNode(yyDollar[3].token), NewNilNode()).setPos0(yyDollar[1].expr.Pos)
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, x)
			yyVAL.expr = yyDollar[1].expr
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:254
		{
			yyVAL.expr = NewStringNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 52:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:260
		{
			yyVAL.expr = NewAtomNode(yyDollar[1].token).setPos(yyDollar[1].token.Pos)
		}
	case 53:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:261
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:262
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 55:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:263
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 56:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:264
		{
			yyVAL.expr = NewNilNode()
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 57:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:268
		{
			yyVAL.expr = NewNumberNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 58:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:272
		{
			path := filepath.Join(filepath.Dir(yyDollar[1].token.Pos.Source), yyDollar[2].token.Str)
			yyVAL.expr = yylex.(*Lexer).loadFile(path)
		}
	case 59:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:278
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 60:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:279
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos).setPos(yyDollar[1].expr.Pos)
		}
	case 61:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:280
		{
			yyVAL.expr = NewCompoundNode("slice", yyDollar[1].expr, yyDollar[3].expr, yyDollar[5].expr).setPos0(yyDollar[1].expr.Pos).setPos(yyDollar[1].expr.Pos)
		}
	case 62:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:281
		{
			yyVAL.expr = NewCompoundNode("slice", yyDollar[1].expr, yyDollar[3].expr, NewNumberNode("-1")).setPos0(yyDollar[1].expr.Pos).setPos(yyDollar[1].expr.Pos)
		}
	case 63:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:282
		{
			yyVAL.expr = NewCompoundNode("slice", yyDollar[1].expr, NewNumberNode("0"), yyDollar[4].expr).setPos0(yyDollar[1].expr.Pos).setPos(yyDollar[1].expr.Pos)
		}
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:283
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, NewStringNode(yyDollar[3].token.Str)).setPos0(yyDollar[1].expr.Pos).setPos(yyDollar[1].expr.Pos)
		}
	case 65:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:284
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 66:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:285
		{
			yyVAL.expr = NewCompoundNode("inc", NewAtomNode(yyDollar[1].token), NewNumberNode("1")).setPos0(yyDollar[1].token.Pos)
		}
	case 67:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:286
		{
			yyVAL.expr = NewCompoundNode("inc", NewAtomNode(yyDollar[1].token), NewNumberNode("-1")).setPos0(yyDollar[1].token.Pos)
		}
	case 68:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:289
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 69:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:290
		{
			yyVAL.expr = NewCompoundNode("-", NewNumberNode("0"), yyDollar[2].expr).setPos0(yyDollar[2].expr.Pos)
		}
	case 70:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:291
		{
			yyVAL.expr = NewCompoundNode("~", yyDollar[2].expr).setPos0(yyDollar[2].expr.Pos)
		}
	case 71:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:292
		{
			yyVAL.expr = NewCompoundNode("!", yyDollar[2].expr).setPos0(yyDollar[2].expr.Pos)
		}
	case 72:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:293
		{
			yyVAL.expr = NewCompoundNode("#", yyDollar[2].expr).setPos0(yyDollar[2].expr.Pos)
		}
	case 73:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:296
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:297
		{
			yyVAL.expr = NewCompoundNode("*", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 75:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:298
		{
			yyVAL.expr = NewCompoundNode("/", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:299
		{
			yyVAL.expr = NewCompoundNode("%", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 77:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:302
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 78:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:303
		{
			yyVAL.expr = NewCompoundNode("+", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 79:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:304
		{
			yyVAL.expr = NewCompoundNode("-", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 80:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:307
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 81:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:308
		{
			yyVAL.expr = NewCompoundNode("<<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 82:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:309
		{
			yyVAL.expr = NewCompoundNode(">>", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 83:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:312
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 84:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:313
		{
			yyVAL.expr = NewCompoundNode("<", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:314
		{
			yyVAL.expr = NewCompoundNode("<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 86:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:315
		{
			yyVAL.expr = NewCompoundNode("<=", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 87:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:316
		{
			yyVAL.expr = NewCompoundNode("<=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 88:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:319
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 89:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:320
		{
			yyVAL.expr = NewCompoundNode("==", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 90:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:321
		{
			yyVAL.expr = NewCompoundNode("!=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 91:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:324
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 92:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:325
		{
			yyVAL.expr = NewCompoundNode("&", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 93:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:328
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 94:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:329
		{
			yyVAL.expr = NewCompoundNode("^", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 95:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:332
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 96:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:333
		{
			yyVAL.expr = NewCompoundNode("|", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 97:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:336
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 98:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:337
		{
			yyVAL.expr = NewCompoundNode("and", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 99:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:340
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 100:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:341
		{
			yyVAL.expr = NewCompoundNode("or", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 101:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:344
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 102:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:345
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
	case 103:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:367
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 104:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:368
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 105:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:371
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 106:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:374
		{
			if yyDollar[1].expr.Type != NTCompound || len(yyDollar[1].expr.Compound) == 0 || yyDollar[1].expr.Compound[0].Value.(string) != "chain" {
				yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
			} else {
				yyVAL.expr = yyDollar[1].expr
			}
			yyVAL.expr.Compound = append(yyVAL.expr.Compound, yyDollar[3].expr)
		}
	case 107:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:384
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
	case 108:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:435
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 109:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:436
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 110:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:439
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 111:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:440
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 112:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:443
		{
			yyVAL.expr = NewCompoundNode("map", NewCompoundNode()).setPos0(yyDollar[1].token.Pos)
		}
	case 113:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:444
		{
			yyVAL.expr = yyDollar[2].expr.setPos0(yyDollar[1].token.Pos)
		}
	case 114:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:447
		{
			yyVAL.expr = NewCompoundNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 115:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:448
		{
			yyVAL.expr = NewCompoundNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 116:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:449
		{
			yyVAL.expr = NewCompoundNode("array", yyDollar[1].expr).setPos0(yyDollar[1].expr.Pos)
		}
	case 117:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:450
		{
			yyVAL.expr = NewCompoundNode("array", yyDollar[1].expr).setPos0(yyDollar[1].expr.Pos)
		}
	}
	goto yystack /* stack new state and value */
}
