//line parser.go.y:2
package parser

import __yyfmt__ "fmt"

//line parser.go.y:2
import (
	"bytes"
	"io/ioutil"
	"path/filepath"
)

//line parser.go.y:41
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

//line parser.go.y:605

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
	-1, 100,
	44, 19,
	-2, 57,
}

const yyPrivate = 57344

const yyLast = 621

var yyAct = [...]int{

	145, 1, 46, 40, 23, 9, 157, 154, 151, 5,
	135, 34, 50, 148, 156, 51, 16, 155, 147, 15,
	53, 109, 134, 58, 57, 168, 108, 96, 89, 33,
	60, 61, 144, 132, 65, 59, 29, 48, 167, 8,
	64, 166, 141, 85, 86, 87, 88, 153, 95, 4,
	26, 30, 25, 31, 7, 98, 103, 130, 100, 133,
	106, 90, 18, 76, 77, 78, 19, 110, 111, 112,
	113, 114, 115, 116, 117, 118, 119, 120, 121, 122,
	123, 124, 125, 126, 127, 104, 159, 55, 94, 54,
	129, 84, 8, 164, 105, 62, 128, 56, 28, 170,
	139, 92, 101, 72, 73, 80, 81, 71, 70, 20,
	39, 143, 37, 32, 6, 66, 67, 82, 83, 79,
	68, 69, 74, 75, 76, 77, 78, 17, 138, 52,
	99, 146, 41, 27, 93, 149, 131, 150, 2, 165,
	0, 9, 9, 0, 0, 162, 9, 0, 0, 23,
	9, 161, 0, 0, 5, 74, 75, 76, 77, 78,
	0, 16, 0, 9, 15, 0, 171, 0, 169, 0,
	80, 81, 9, 0, 9, 8, 8, 172, 0, 173,
	8, 0, 0, 0, 8, 0, 158, 74, 75, 76,
	77, 78, 0, 0, 4, 0, 0, 8, 72, 73,
	80, 81, 71, 70, 0, 0, 8, 0, 8, 0,
	66, 67, 82, 83, 79, 68, 69, 74, 75, 76,
	77, 78, 0, 72, 73, 80, 81, 71, 70, 0,
	0, 0, 0, 0, 136, 66, 67, 82, 83, 79,
	68, 69, 74, 75, 76, 77, 78, 72, 73, 80,
	81, 71, 70, 0, 0, 0, 107, 0, 0, 66,
	67, 82, 83, 79, 68, 69, 74, 75, 76, 77,
	78, 72, 73, 80, 81, 71, 70, 142, 0, 0,
	0, 0, 0, 66, 67, 82, 83, 79, 68, 69,
	74, 75, 76, 77, 78, 72, 73, 80, 81, 71,
	70, 137, 0, 0, 0, 0, 0, 66, 67, 82,
	83, 79, 68, 69, 74, 75, 76, 77, 78, 72,
	73, 80, 81, 71, 70, 97, 0, 0, 0, 0,
	0, 66, 67, 82, 83, 79, 68, 69, 74, 75,
	76, 77, 78, 0, 24, 0, 35, 160, 38, 7,
	0, 24, 0, 35, 0, 38, 0, 18, 36, 49,
	47, 19, 0, 0, 18, 36, 49, 47, 19, 0,
	0, 42, 0, 0, 0, 0, 43, 45, 42, 102,
	0, 0, 44, 43, 45, 0, 140, 0, 24, 44,
	35, 0, 38, 0, 0, 0, 24, 0, 35, 0,
	38, 18, 36, 49, 47, 19, 0, 0, 0, 18,
	36, 49, 47, 19, 0, 42, 0, 0, 0, 0,
	43, 45, 91, 42, 0, 0, 44, 0, 43, 45,
	0, 0, 0, 63, 44, 24, 0, 35, 0, 38,
	0, 0, 0, 0, 0, 0, 0, 0, 18, 36,
	49, 47, 19, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 42, 0, 0, 0, 0, 43, 45, 0,
	0, 0, 0, 44, 72, 73, 80, 81, 71, 70,
	0, 0, 0, 0, 0, 0, 66, 67, 82, 83,
	79, 68, 69, 74, 75, 76, 77, 78, 14, 12,
	13, 0, 21, 24, 22, 0, 10, 0, 7, 11,
	0, 14, 12, 13, 0, 21, 18, 22, 0, 10,
	19, 7, 11, 0, 0, 0, 0, 0, 0, 18,
	0, 0, 129, 19, 0, 0, 0, 163, 3, 72,
	73, 80, 81, 71, 70, 0, 0, 0, 0, 0,
	0, 152, 67, 82, 83, 79, 68, 69, 74, 75,
	76, 77, 78, 14, 12, 13, 0, 21, 24, 22,
	0, 10, 0, 7, 11, 0, 0, 0, 0, 0,
	0, 18, 0, 0, 0, 19, 0, 0, 0, 0,
	72, 73, 80, 81, 71, 70, 0, 72, 73, 80,
	81, 71, 70, 3, 82, 83, 79, 68, 69, 74,
	75, 76, 77, 78, 68, 69, 74, 75, 76, 77,
	78,
}
var yyPact = [...]int{

	-1000, 559, -1000, -1000, 8, 6, -1000, 76, -9, 3,
	426, 426, -1000, -1000, 426, -1000, -1000, -1000, -1000, 426,
	-1000, 63, 61, 75, -23, -1000, -1000, -28, -10, 426,
	426, 73, -1000, 387, 458, -1000, -1000, -1000, 67, -1000,
	3, -1000, 426, 426, 426, 426, 35, 379, -1000, -1000,
	458, 458, -19, 279, 335, 426, 35, -1000, 72, 426,
	458, 207, -1000, -1000, -25, 458, 426, 426, 426, 426,
	426, 426, 426, 426, 426, 426, 426, 426, 426, 426,
	426, 426, 426, 426, -1000, -1000, -1000, -1000, -1000, 65,
	11, -1000, 16, -29, -41, 182, -1000, -1000, 255, 342,
	3, -2, -1000, 231, 65, -13, 458, -1000, 426, -1000,
	523, 574, 152, 152, 152, 152, 152, 152, 26, 26,
	-1000, -1000, -1000, 581, 120, 120, 581, 581, -1000, -1000,
	-1000, -33, -1000, -1000, 426, 426, 426, 507, 40, 303,
	-1000, -1000, 507, -1000, 426, 458, 494, 71, -1000, 87,
	458, -1000, -1000, -3, -6, -1000, -1000, -1000, -21, 507,
	-1000, 92, 458, -1000, -1000, 426, -1000, -1000, 507, -1000,
	507, 458, -1000, -1000,
}
var yyPgo = [...]int{

	0, 1, 6, 138, 37, 136, 40, 134, 133, 0,
	132, 3, 47, 17, 130, 128, 14, 8, 7, 127,
	114, 2, 109, 113, 112, 28, 110, 101,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 3, 3, 3, 3, 17, 17,
	17, 17, 17, 17, 20, 20, 20, 12, 12, 12,
	14, 14, 15, 15, 13, 13, 13, 16, 16, 21,
	21, 19, 18, 18, 18, 18, 18, 18, 18, 4,
	4, 4, 5, 5, 6, 6, 7, 7, 8, 8,
	8, 8, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 10, 11, 11, 11, 11, 22, 23, 23, 24,
	25, 25, 26, 26, 27, 27, 27, 27,
}
var yyR2 = [...]int{

	0, 0, 2, 3, 1, 2, 2, 1, 1, 2,
	2, 1, 1, 1, 1, 1, 1, 2, 3, 1,
	2, 1, 2, 1, 5, 7, 6, 5, 7, 1,
	2, 4, 1, 2, 1, 2, 1, 1, 2, 1,
	4, 3, 1, 3, 1, 3, 3, 5, 1, 3,
	5, 3, 1, 1, 1, 2, 1, 1, 1, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 2, 2, 2,
	2, 1, 1, 3, 1, 3, 2, 2, 3, 3,
	2, 3, 2, 3, 1, 2, 1, 2,
}
var yyChk = [...]int{

	-1000, -1, -3, 44, -12, -18, -20, 14, -4, -11,
	12, 15, 5, 6, 4, -13, -16, -19, 22, 26,
	-22, 8, 10, -21, 9, 44, 44, -8, 22, 45,
	48, 50, -23, 26, -9, 11, 23, -24, 13, -26,
	-11, -10, 36, 41, 47, 42, -21, 25, -4, 24,
	-9, -9, -22, -9, 26, 26, 22, 47, 51, 45,
	-9, -9, 22, 46, -6, -9, 28, 29, 33, 34,
	21, 20, 16, 17, 35, 36, 37, 38, 39, 32,
	18, 19, 30, 31, 24, -9, -9, -9, -9, -25,
	26, 43, -27, -7, -6, -9, 46, 46, -9, -14,
	-11, -12, 44, -9, -25, 22, -9, 49, 51, 46,
	-9, -9, -9, -9, -9, -9, -9, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, -9, -9, -2, 25,
	46, -5, 22, 43, 51, 51, 52, 46, -15, -9,
	44, 44, 46, -2, 45, -9, -1, 51, 46, -9,
	-9, -17, 44, -12, -18, -13, -16, -2, -12, 46,
	44, -17, -9, 43, 22, 52, 44, 44, 46, -17,
	7, -9, -17, -17,
}
var yyDef = [...]int{

	1, -2, 2, 4, 0, 0, 7, 0, 82, 19,
	32, 34, 36, 37, 0, 14, 15, 16, 39, 0,
	84, 0, 0, 0, 29, 5, 6, 17, 48, 0,
	0, 0, 86, 0, 33, 52, 53, 54, 0, 56,
	57, 58, 0, 0, 0, 0, 0, 0, 82, 81,
	35, 38, 84, 0, 0, 0, 0, 30, 0, 0,
	18, 0, 41, 87, 0, 44, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 55, 77, 78, 79, 80, 0,
	0, 92, 0, 94, 96, 44, 83, 85, 0, 0,
	-2, 0, 21, 0, 0, 51, 49, 40, 0, 88,
	59, 60, 61, 62, 63, 64, 65, 66, 67, 68,
	69, 70, 71, 72, 73, 74, 75, 76, 89, 1,
	90, 0, 42, 93, 95, 97, 0, 0, 0, 0,
	23, 20, 0, 31, 0, 45, 0, 0, 91, 0,
	46, 24, 8, 0, 0, 11, 12, 13, 0, 0,
	22, 27, 50, 3, 43, 0, 9, 10, 0, 26,
	0, 47, 25, 28,
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
		//line parser.go.y:71
		{
			yyVAL.expr = NewCompoundNode("chain")
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:77
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, yyDollar[2].expr)
			yyVAL.expr = yyDollar[1].expr
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 3:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:86
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:91
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:94
		{
			if yyDollar[1].expr.isIsolatedDupCall() {
				yyDollar[1].expr.Compound[2].Compound[0] = NewNumberNode("0")
			}
			yyVAL.expr = yyDollar[1].expr
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:100
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:103
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:108
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:111
		{
			if yyDollar[1].expr.isIsolatedDupCall() {
				yyDollar[1].expr.Compound[2].Compound[0] = NewNumberNode("0")
			}
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:117
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:120
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:123
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:126
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:131
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:134
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:137
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 17:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:142
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:145
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
		//line parser.go.y:165
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 20:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:170
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:173
		{
			yyVAL.expr = NewCompoundNode("chain")
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:178
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 23:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:181
		{
			yyVAL.expr = NewNumberNode("1")
		}
	case 24:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:186
		{
			yyVAL.expr = NewCompoundNode("for", yyDollar[3].expr, NewCompoundNode(), yyDollar[5].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 25:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:190
		{
			yyVAL.expr = yyDollar[3].expr
			yyVAL.expr.Compound = append(yyVAL.expr.Compound, NewCompoundNode("for", yyDollar[4].expr, NewCompoundNode("chain", yyDollar[5].expr), yyDollar[7].expr))
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 26:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:195
		{
			yyVAL.expr = yyDollar[3].expr
			yyVAL.expr.Compound = append(yyVAL.expr.Compound, NewCompoundNode("for", yyDollar[4].expr, NewCompoundNode(), yyDollar[6].expr))
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 27:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:202
		{
			yyVAL.expr = NewCompoundNode("if", yyDollar[3].expr, yyDollar[5].expr, NewCompoundNode())
		}
	case 28:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:205
		{
			yyVAL.expr = NewCompoundNode("if", yyDollar[3].expr, yyDollar[5].expr, yyDollar[7].expr)
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:210
		{
			yyVAL.str = "func"
		}
	case 30:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:213
		{
			yyVAL.str = "safefunc"
		}
	case 31:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:218
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
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:230
		{
			yyVAL.expr = NewCompoundNode("ret")
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 33:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:234
		{
			if yyDollar[2].expr.isIsolatedDupCall() {
				if h, _ := yyDollar[2].expr.Compound[2].Compound[2].Value.(float64); h == 1 {
					yyDollar[2].expr.Compound[2].Compound[2] = NewNumberNode("2")
				}
			}
			yyVAL.expr = NewCompoundNode("ret", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:243
		{
			yyVAL.expr = NewCompoundNode("yield")
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 35:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:247
		{
			yyVAL.expr = NewCompoundNode("yield", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:251
		{
			yyVAL.expr = NewCompoundNode("break")
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:255
		{
			yyVAL.expr = NewCompoundNode("continue")
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 38:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:259
		{
			yyVAL.expr = NewCompoundNode("assert", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:265
		{
			yyVAL.expr = NewAtomNode(yyDollar[1].token)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 40:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:269
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 41:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:274
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, NewStringNode(yyDollar[3].token.Str))
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:281
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].token.Str)
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:284
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, NewAtomNode(yyDollar[3].token))
			yyVAL.expr = yyDollar[1].expr
		}
	case 44:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:290
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].expr)
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:293
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, yyDollar[3].expr)
			yyVAL.expr = yyDollar[1].expr
		}
	case 46:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:299
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 47:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:302
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, yyDollar[3].expr, yyDollar[5].expr)
			yyVAL.expr = yyDollar[1].expr
		}
	case 48:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:308
		{
			yyVAL.expr = NewCompoundNode("chain", NewCompoundNode("set", NewAtomNode(yyDollar[1].token), NewNilNode()))
			yyVAL.expr.Compound[1].Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 49:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:312
		{
			yyVAL.expr = NewCompoundNode("chain", NewCompoundNode("set", NewAtomNode(yyDollar[1].token), yyDollar[3].expr))
			yyVAL.expr.Compound[1].Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 50:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:316
		{
			x := NewCompoundNode("set", NewAtomNode(yyDollar[3].token), yyDollar[5].expr)
			x.Compound[0].Pos = yyDollar[1].expr.Pos
			yyDollar[1].expr.Compound = append(yyVAL.expr.Compound, x)
			yyVAL.expr = yyDollar[1].expr
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:322
		{
			x := NewCompoundNode("set", NewAtomNode(yyDollar[3].token), NewNilNode())
			x.Compound[0].Pos = yyDollar[1].expr.Pos
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, x)
			yyVAL.expr = yyDollar[1].expr
		}
	case 52:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:330
		{
			yyVAL.expr = NewNilNode()
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 53:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:334
		{
			yyVAL.expr = NewNumberNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:338
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 55:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:341
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
	case 56:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:358
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 57:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:361
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 58:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:364
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:367
		{
			yyVAL.expr = NewCompoundNode("or", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:371
		{
			yyVAL.expr = NewCompoundNode("and", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:375
		{
			yyVAL.expr = NewCompoundNode("<", yyDollar[3].expr, yyDollar[1].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 62:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:379
		{
			yyVAL.expr = NewCompoundNode("<", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 63:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:383
		{
			yyVAL.expr = NewCompoundNode("<=", yyDollar[3].expr, yyDollar[1].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:387
		{
			yyVAL.expr = NewCompoundNode("<=", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 65:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:391
		{
			yyVAL.expr = NewCompoundNode("==", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:395
		{
			yyVAL.expr = NewCompoundNode("!=", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:399
		{
			yyVAL.expr = NewCompoundNode("+", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:403
		{
			yyVAL.expr = NewCompoundNode("-", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:407
		{
			yyVAL.expr = NewCompoundNode("*", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 70:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:411
		{
			yyVAL.expr = NewCompoundNode("/", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:415
		{
			yyVAL.expr = NewCompoundNode("%", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 72:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:419
		{
			yyVAL.expr = NewCompoundNode("^", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:423
		{
			yyVAL.expr = NewCompoundNode("<<", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:427
		{
			yyVAL.expr = NewCompoundNode(">>", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 75:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:431
		{
			yyVAL.expr = NewCompoundNode("|", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:435
		{
			yyVAL.expr = NewCompoundNode("&", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 77:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:439
		{
			yyVAL.expr = NewCompoundNode("-", NewNumberNode("0"), yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 78:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:443
		{
			yyVAL.expr = NewCompoundNode("~", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 79:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:447
		{
			yyVAL.expr = NewCompoundNode("!", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 80:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:451
		{
			yyVAL.expr = NewCompoundNode("#", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 81:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:457
		{
			yyVAL.expr = NewStringNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 82:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:463
		{
			yyVAL.expr = yyDollar[1].expr
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 83:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:467
		{
			yyVAL.expr = yyDollar[2].expr
			yyVAL.expr.Pos = yyDollar[2].expr.Pos
		}
	case 84:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:471
		{
			yyVAL.expr = yyDollar[1].expr
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:475
		{
			yyVAL.expr = yyDollar[2].expr
			yyVAL.expr.Pos = yyDollar[2].expr.Pos
		}
	case 86:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:481
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
	case 87:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:542
		{
			if yylex.(*Lexer).PNewLine {
				yylex.(*Lexer).TokenError(yyDollar[1].token, "ambiguous syntax (function call x new statement)")
			}
			yyVAL.expr = NewCompoundNode()
		}
	case 88:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:548
		{
			if yylex.(*Lexer).PNewLine {
				yylex.(*Lexer).TokenError(yyDollar[1].token, "ambiguous syntax (function call x new statement)")
			}
			yyVAL.expr = yyDollar[2].expr
		}
	case 89:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:556
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].str, yyDollar[2].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 90:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:562
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 91:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:565
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 92:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:570
		{
			yyVAL.expr = NewCompoundNode("map", NewCompoundNode())
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 93:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:574
		{
			yyVAL.expr = yyDollar[2].expr
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 94:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:580
		{
			yyVAL.expr = NewCompoundNode("map", yyDollar[1].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 95:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:584
		{
			yyVAL.expr = NewCompoundNode("map", yyDollar[1].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 96:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:588
		{
			table := NewCompoundNode()
			for i, v := range yyDollar[1].expr.Compound {
				table.Compound = append(table.Compound, &Node{Type: NTNumber, Value: float64(i)}, v)
			}
			yyVAL.expr = NewCompoundNode("map", table)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 97:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:596
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
