//line parser.go.y:1
package parser

import __yyfmt__ "fmt"

//line parser.go.y:3
import (
	"bytes"
	"io/ioutil"
	"path/filepath"
)

//line parser.go.y:38
type yySymType struct {
	yys   int
	token Token
	expr  *Node
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
	"'['",
	"']'",
	"'.'",
	"','",
	"':'",
	"'!'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.go.y:567

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
	-1, 99,
	44, 19,
	-2, 53,
}

const yyPrivate = 57344

const yyLast = 661

var yyAct = [...]int{

	140, 1, 148, 39, 4, 9, 152, 149, 87, 5,
	146, 34, 49, 143, 33, 50, 132, 142, 151, 106,
	17, 53, 150, 105, 16, 131, 57, 129, 174, 171,
	59, 60, 168, 94, 64, 30, 139, 31, 47, 63,
	8, 58, 83, 84, 85, 86, 29, 93, 162, 88,
	161, 127, 7, 136, 26, 96, 100, 98, 99, 103,
	19, 25, 130, 55, 20, 101, 107, 108, 109, 110,
	111, 112, 113, 114, 115, 116, 117, 118, 119, 120,
	121, 122, 123, 124, 170, 54, 92, 73, 74, 75,
	76, 77, 126, 8, 125, 75, 76, 77, 135, 7,
	71, 72, 79, 80, 70, 69, 51, 19, 138, 159,
	102, 20, 65, 66, 81, 82, 78, 67, 68, 73,
	74, 75, 76, 77, 61, 56, 28, 166, 141, 21,
	90, 38, 144, 37, 145, 160, 32, 154, 9, 6,
	157, 9, 18, 40, 4, 9, 27, 91, 156, 5,
	52, 128, 2, 0, 0, 0, 163, 9, 165, 9,
	17, 167, 0, 0, 16, 0, 0, 169, 9, 0,
	9, 0, 9, 8, 9, 9, 8, 172, 9, 173,
	8, 175, 176, 0, 0, 177, 0, 0, 0, 0,
	0, 0, 8, 0, 8, 0, 0, 0, 0, 0,
	0, 0, 0, 8, 0, 8, 0, 8, 0, 8,
	8, 0, 0, 8, 71, 72, 79, 80, 70, 69,
	0, 0, 0, 0, 0, 0, 65, 66, 81, 82,
	78, 67, 68, 73, 74, 75, 76, 77, 71, 72,
	79, 80, 70, 69, 0, 0, 0, 0, 0, 133,
	65, 66, 81, 82, 78, 67, 68, 73, 74, 75,
	76, 77, 71, 72, 79, 80, 70, 69, 0, 0,
	104, 0, 0, 0, 65, 66, 81, 82, 78, 67,
	68, 73, 74, 75, 76, 77, 71, 72, 79, 80,
	70, 69, 137, 0, 0, 0, 0, 0, 65, 66,
	81, 82, 78, 67, 68, 73, 74, 75, 76, 77,
	45, 0, 35, 79, 80, 7, 134, 0, 0, 0,
	0, 0, 0, 19, 36, 48, 46, 20, 0, 0,
	73, 74, 75, 76, 77, 0, 0, 41, 0, 0,
	0, 0, 42, 44, 0, 97, 71, 72, 79, 80,
	70, 69, 0, 43, 0, 0, 0, 0, 65, 66,
	81, 82, 78, 67, 68, 73, 74, 75, 76, 77,
	45, 0, 35, 0, 0, 0, 95, 45, 0, 35,
	0, 0, 0, 19, 36, 48, 46, 20, 0, 0,
	19, 36, 48, 46, 20, 0, 0, 41, 0, 0,
	0, 0, 42, 44, 41, 155, 0, 0, 0, 42,
	44, 89, 45, 43, 35, 0, 0, 0, 0, 45,
	43, 35, 0, 0, 0, 19, 36, 48, 46, 20,
	0, 0, 19, 36, 48, 46, 20, 0, 0, 41,
	0, 0, 0, 0, 42, 44, 41, 0, 0, 62,
	0, 42, 44, 0, 0, 43, 71, 72, 79, 80,
	70, 69, 43, 0, 0, 0, 0, 0, 65, 66,
	81, 82, 78, 67, 68, 73, 74, 75, 76, 77,
	0, 0, 0, 0, 164, 71, 72, 79, 80, 70,
	69, 0, 0, 0, 0, 0, 0, 65, 66, 81,
	82, 78, 67, 68, 73, 74, 75, 76, 77, 0,
	14, 12, 13, 153, 22, 24, 23, 0, 10, 15,
	7, 11, 0, 0, 0, 14, 12, 13, 19, 22,
	0, 23, 20, 10, 15, 7, 11, 0, 0, 0,
	0, 0, 0, 19, 0, 0, 126, 20, 0, 158,
	3, 0, 0, 0, 71, 72, 79, 80, 70, 69,
	0, 0, 0, 0, 0, 147, 65, 66, 81, 82,
	78, 67, 68, 73, 74, 75, 76, 77, 14, 12,
	13, 0, 22, 24, 23, 0, 10, 15, 7, 11,
	0, 0, 0, 0, 0, 0, 19, 0, 0, 0,
	20, 0, 0, 0, 0, 0, 71, 72, 79, 80,
	70, 69, 0, 0, 0, 0, 0, 0, 3, 66,
	81, 82, 78, 67, 68, 73, 74, 75, 76, 77,
	71, 72, 79, 80, 70, 69, 0, 71, 72, 79,
	80, 70, 69, 0, 81, 82, 78, 67, 68, 73,
	74, 75, 76, 77, 67, 68, 73, 74, 75, 76,
	77,
}
var yyPact = [...]int{

	-1000, 574, -1000, -1000, 17, 10, -1000, 104, 1, -12,
	410, 410, -1000, -1000, 410, 82, -1000, -1000, -1000, -1000,
	410, -1000, 59, 37, 103, -1000, -1000, -24, -4, 410,
	410, 102, -1000, 403, 538, -1000, -1000, -1000, -1000, -12,
	-1000, 410, 410, 410, 410, 23, 368, -1000, -1000, 538,
	538, -1000, -13, 330, 301, 410, 23, 88, 410, 538,
	222, -1000, -1000, -27, 538, 410, 410, 410, 410, 410,
	410, 410, 410, 410, 410, 410, 410, 410, 410, 410,
	410, 410, 410, -1000, -1000, -1000, -1000, 67, 5, -1000,
	19, -25, -34, 198, -1000, -1000, 270, 410, 9, -12,
	246, 67, -9, 538, -1000, 410, -1000, 590, 614, 295,
	295, 295, 295, 295, 295, 58, 58, -1000, -1000, -1000,
	621, 52, 52, 621, 621, -1000, -1000, -1000, -33, -1000,
	-1000, 410, 410, 410, 521, 469, 361, 521, -1000, 410,
	538, 506, 87, -1000, 84, 538, -1000, -1000, 6, 4,
	-1000, -1000, -1000, 85, 440, 85, 120, 538, -1000, -1000,
	410, -1000, -1000, -14, 38, -17, 521, 538, 521, -18,
	521, 521, -1000, -1000, 521, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 1, 6, 152, 38, 151, 39, 147, 146, 0,
	143, 3, 2, 22, 18, 10, 7, 142, 139, 129,
	136, 133, 8, 131, 130,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 3, 3, 3, 3, 15, 15,
	15, 15, 15, 15, 18, 18, 18, 12, 12, 12,
	13, 13, 13, 13, 13, 14, 14, 17, 16, 16,
	16, 16, 16, 16, 16, 16, 4, 4, 4, 5,
	5, 6, 6, 7, 7, 8, 8, 8, 8, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 10, 11, 11,
	11, 11, 19, 20, 20, 21, 22, 22, 23, 23,
	24, 24, 24, 24,
}
var yyR2 = [...]int{

	0, 0, 2, 3, 1, 2, 2, 1, 1, 2,
	2, 1, 1, 1, 1, 1, 1, 2, 3, 1,
	5, 8, 9, 8, 8, 5, 7, 4, 1, 2,
	1, 2, 1, 1, 2, 2, 1, 4, 3, 1,
	3, 1, 3, 3, 5, 1, 3, 5, 3, 1,
	1, 1, 1, 1, 1, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 2, 2, 2, 2, 1, 1, 3,
	1, 3, 2, 2, 3, 3, 2, 3, 2, 3,
	1, 2, 1, 2,
}
var yyChk = [...]int{

	-1000, -1, -3, 44, -12, -16, -18, 14, -4, -11,
	12, 15, 5, 6, 4, 13, -13, -14, -17, 22,
	26, -19, 8, 10, 9, 44, 44, -8, 22, 45,
	47, 49, -20, 26, -9, 11, 23, -21, -23, -11,
	-10, 36, 41, 52, 42, 9, 25, -4, 24, -9,
	-9, 24, -19, -9, 26, 26, 22, 50, 45, -9,
	-9, 22, 46, -6, -9, 28, 29, 33, 34, 21,
	20, 16, 17, 35, 36, 37, 38, 39, 32, 18,
	19, 30, 31, -9, -9, -9, -9, -22, 26, 43,
	-24, -7, -6, -9, 46, 46, -9, 44, -12, -11,
	-9, -22, 22, -9, 48, 50, 46, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, -9, -9, -9, -9,
	-9, -9, -9, -9, -9, -2, 25, 46, -5, 22,
	43, 50, 50, 51, 46, -9, 44, 46, -2, 45,
	-9, -1, 50, 46, -9, -9, -15, 44, -12, -16,
	-13, -14, -2, 44, -9, 44, -15, -9, 43, 22,
	51, 44, 44, -12, 44, -12, 7, -9, 46, -12,
	46, 46, -15, -15, 46, -15, -15, -15,
}
var yyDef = [...]int{

	1, -2, 2, 4, 0, 0, 7, 0, 78, 19,
	28, 30, 32, 33, 0, 0, 14, 15, 16, 36,
	0, 80, 0, 0, 0, 5, 6, 17, 45, 0,
	0, 0, 82, 0, 29, 49, 50, 51, 52, 53,
	54, 0, 0, 0, 0, 0, 0, 78, 77, 31,
	34, 35, 80, 0, 0, 0, 0, 0, 0, 18,
	0, 38, 83, 0, 41, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 73, 74, 75, 76, 0, 0, 88,
	0, 90, 92, 41, 79, 81, 0, 0, 0, -2,
	0, 0, 48, 46, 37, 0, 84, 55, 56, 57,
	58, 59, 60, 61, 62, 63, 64, 65, 66, 67,
	68, 69, 70, 71, 72, 85, 1, 86, 0, 39,
	89, 91, 93, 0, 0, 0, 0, 0, 27, 0,
	42, 0, 0, 87, 0, 43, 20, 8, 0, 0,
	11, 12, 13, 0, 0, 0, 25, 47, 3, 40,
	0, 9, 10, 0, 0, 0, 0, 44, 0, 0,
	0, 0, 26, 21, 0, 23, 24, 22,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 52, 3, 42, 3, 39, 31, 3,
	26, 46, 37, 35, 50, 36, 49, 38, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 51, 44,
	34, 45, 33, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 27, 3, 3, 3, 3, 3,
	3, 47, 3, 48, 32, 3, 3, 3, 3, 3,
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
		//line parser.go.y:67
		{
			yyVAL.expr = NewCompoundNode("chain")
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:73
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, yyDollar[2].expr)
			yyVAL.expr = yyDollar[1].expr
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 3:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:82
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:87
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:90
		{
			if yyDollar[1].expr.isIsolatedDupCall() {
				yyDollar[1].expr.Compound[2].Compound[0] = NewNumberNode("0")
			}
			yyVAL.expr = yyDollar[1].expr
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:96
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:99
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:104
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:107
		{
			if yyDollar[1].expr.isIsolatedDupCall() {
				yyDollar[1].expr.Compound[2].Compound[0] = NewNumberNode("0")
			}
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:113
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:116
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:119
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[1].expr)
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:122
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:127
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:130
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:133
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 17:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:138
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:141
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
		//line parser.go.y:161
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 20:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:166
		{
			yyVAL.expr = NewCompoundNode("for", yyDollar[3].expr, NewCompoundNode(), yyDollar[5].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 21:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.go.y:170
		{
			yyVAL.expr = NewCompoundNode("for", yyDollar[4].expr, NewCompoundNode("chain", yyDollar[6].expr), yyDollar[8].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 22:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.go.y:174
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[3].expr, NewCompoundNode("for", yyDollar[5].expr, NewCompoundNode("chain", yyDollar[7].expr), yyDollar[9].expr))
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 23:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.go.y:178
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[3].expr, NewCompoundNode("for", yyDollar[5].expr, NewCompoundNode(), yyDollar[8].expr))
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 24:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.go.y:182
		{
			yyVAL.expr = NewCompoundNode("chain", yyDollar[3].expr, NewCompoundNode("for", NewNumberNode("1"), NewCompoundNode("chain", yyDollar[6].expr), yyDollar[8].expr))
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 25:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:188
		{
			yyVAL.expr = NewCompoundNode("if", yyDollar[3].expr, yyDollar[5].expr, NewCompoundNode())
		}
	case 26:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:191
		{
			yyVAL.expr = NewCompoundNode("if", yyDollar[3].expr, yyDollar[5].expr, yyDollar[7].expr)
		}
	case 27:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:196
		{
			funcname := NewAtomNode(yyDollar[2].token)
			yyVAL.expr = NewCompoundNode(
				"chain",
				NewCompoundNode("set", funcname, NewNilNode()),
				NewCompoundNode("move", funcname, NewCompoundNode("lambda", yyDollar[3].expr, yyDollar[4].expr)))
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:205
		{
			yyVAL.expr = NewCompoundNode("ret")
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 29:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:209
		{
			if yyDollar[2].expr.isIsolatedDupCall() {
				if h, _ := yyDollar[2].expr.Compound[2].Compound[2].Value.(float64); h == 1 {
					yyDollar[2].expr.Compound[2].Compound[2] = NewNumberNode("2")
				}
			}
			yyVAL.expr = NewCompoundNode("ret", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:218
		{
			yyVAL.expr = NewCompoundNode("yield")
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 31:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:222
		{
			yyVAL.expr = NewCompoundNode("yield", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:226
		{
			yyVAL.expr = NewCompoundNode("break")
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:230
		{
			yyVAL.expr = NewCompoundNode("continue")
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 34:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:234
		{
			yyVAL.expr = NewCompoundNode("assert", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 35:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:238
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
			yyVAL.expr = NewCompoundNode("set", filename, call)
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:260
		{
			yyVAL.expr = NewAtomNode(yyDollar[1].token)
		}
	case 37:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:263
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, yyDollar[3].expr)
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:266
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, NewStringNode(yyDollar[3].token.Str))
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:271
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].token.Str)
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:274
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, NewAtomNode(yyDollar[3].token))
			yyVAL.expr = yyDollar[1].expr
		}
	case 41:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:280
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].expr)
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:283
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, yyDollar[3].expr)
			yyVAL.expr = yyDollar[1].expr
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:289
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 44:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:292
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, yyDollar[3].expr, yyDollar[5].expr)
			yyVAL.expr = yyDollar[1].expr
		}
	case 45:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:298
		{
			yyVAL.expr = NewCompoundNode("chain", NewCompoundNode("set", NewAtomNode(yyDollar[1].token), NewNilNode()))
		}
	case 46:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:301
		{
			yyVAL.expr = NewCompoundNode("chain", NewCompoundNode("set", NewAtomNode(yyDollar[1].token), yyDollar[3].expr))
		}
	case 47:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:304
		{
			yyDollar[1].expr.Compound = append(yyVAL.expr.Compound, NewCompoundNode("set", NewAtomNode(yyDollar[3].token), yyDollar[5].expr))
			yyVAL.expr = yyDollar[1].expr
		}
	case 48:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:308
		{
			yyDollar[1].expr.Compound = append(yyDollar[1].expr.Compound, NewCompoundNode("set", NewAtomNode(yyDollar[3].token), NewNilNode()))
			yyVAL.expr = yyDollar[1].expr
		}
	case 49:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:314
		{
			yyVAL.expr = NewNilNode()
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 50:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:318
		{
			yyVAL.expr = NewNumberNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:322
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 52:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:325
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 53:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:328
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:331
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:334
		{
			yyVAL.expr = NewCompoundNode("or", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:338
		{
			yyVAL.expr = NewCompoundNode("and", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:342
		{
			yyVAL.expr = NewCompoundNode("<", yyDollar[3].expr, yyDollar[1].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:346
		{
			yyVAL.expr = NewCompoundNode("<", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:350
		{
			yyVAL.expr = NewCompoundNode("<=", yyDollar[3].expr, yyDollar[1].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:354
		{
			yyVAL.expr = NewCompoundNode("<=", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:358
		{
			yyVAL.expr = NewCompoundNode("==", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 62:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:362
		{
			yyVAL.expr = NewCompoundNode("!=", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 63:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:366
		{
			yyVAL.expr = NewCompoundNode("+", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:370
		{
			yyVAL.expr = NewCompoundNode("-", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 65:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:374
		{
			yyVAL.expr = NewCompoundNode("*", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:378
		{
			yyVAL.expr = NewCompoundNode("/", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:382
		{
			yyVAL.expr = NewCompoundNode("%", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:386
		{
			yyVAL.expr = NewCompoundNode("^", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:390
		{
			yyVAL.expr = NewCompoundNode("<<", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 70:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:394
		{
			yyVAL.expr = NewCompoundNode(">>", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:398
		{
			yyVAL.expr = NewCompoundNode("|", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 72:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:402
		{
			yyVAL.expr = NewCompoundNode("&", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 73:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:406
		{
			yyVAL.expr = NewCompoundNode("-", NewNumberNode("0"), yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 74:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:410
		{
			yyVAL.expr = NewCompoundNode("~", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 75:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:414
		{
			yyVAL.expr = NewCompoundNode("!", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 76:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:418
		{
			yyVAL.expr = NewCompoundNode("#", yyDollar[2].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[2].expr.Pos
		}
	case 77:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:424
		{
			yyVAL.expr = NewStringNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 78:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:430
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 79:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:433
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 80:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:436
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 81:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:439
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 82:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:444
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
		}
	case 83:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:504
		{
			if yylex.(*Lexer).PNewLine {
				yylex.(*Lexer).TokenError(yyDollar[1].token, "ambiguous syntax (function call x new statement)")
			}
			yyVAL.expr = NewCompoundNode()
		}
	case 84:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:510
		{
			if yylex.(*Lexer).PNewLine {
				yylex.(*Lexer).TokenError(yyDollar[1].token, "ambiguous syntax (function call x new statement)")
			}
			yyVAL.expr = yyDollar[2].expr
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:518
		{
			yyVAL.expr = NewCompoundNode("lambda", yyDollar[2].expr, yyDollar[3].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 86:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:524
		{
			yyVAL.expr = NewCompoundNode()
		}
	case 87:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:527
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 88:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:532
		{
			yyVAL.expr = NewCompoundNode("map", NewCompoundNode())
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 89:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:536
		{
			yyVAL.expr = yyDollar[2].expr
			yyVAL.expr.Compound[0].Pos = yyDollar[1].token.Pos
		}
	case 90:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:542
		{
			yyVAL.expr = NewCompoundNode("map", yyDollar[1].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 91:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:546
		{
			yyVAL.expr = NewCompoundNode("map", yyDollar[1].expr)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 92:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:550
		{
			table := NewCompoundNode()
			for i, v := range yyDollar[1].expr.Compound {
				table.Compound = append(table.Compound, &Node{Type: NTNumber, Value: float64(i)}, v)
			}
			yyVAL.expr = NewCompoundNode("map", table)
			yyVAL.expr.Compound[0].Pos = yyDollar[1].expr.Pos
		}
	case 93:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:558
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
