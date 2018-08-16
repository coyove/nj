//line parser.go.y:1
package parser

import __yyfmt__ "fmt"

//line parser.go.y:3
import (
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
const TGoto = 57352
const TIf = 57353
const TNil = 57354
const TReturn = 57355
const TRequire = 57356
const TVar = 57357
const TYield = 57358
const TAddAdd = 57359
const TSubSub = 57360
const TAddEq = 57361
const TSubEq = 57362
const TEqeq = 57363
const TNeq = 57364
const TLsh = 57365
const TRsh = 57366
const TLte = 57367
const TGte = 57368
const TIdent = 57369
const TNumber = 57370
const TString = 57371
const TOr = 57372
const TAnd = 57373
const UNARY = 57374
const TMinMin = 57375

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
	"TSubSub",
	"TAddEq",
	"TSubEq",
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
	"TMinMin",
	"'}'",
	"';'",
	"'='",
	"'['",
	"']'",
	"'.'",
	"')'",
	"'!'",
	"':'",
	"','",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.go.y:351

var typesLookup = map[string]string{
	"nil": "0", "number": "1", "string": "2", "map": "4", "closure": "6", "generic": "7",
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 100,
	50, 19,
	-2, 70,
	-1, 101,
	50, 20,
	-2, 71,
}

const yyPrivate = 57344

const yyLast = 736

var yyAct = [...]int{

	163, 35, 65, 21, 16, 171, 96, 165, 29, 142,
	141, 46, 47, 42, 36, 19, 22, 168, 170, 4,
	14, 107, 169, 50, 13, 167, 64, 6, 161, 159,
	186, 160, 158, 62, 60, 44, 185, 23, 116, 87,
	88, 89, 90, 91, 97, 177, 1, 108, 20, 61,
	184, 148, 98, 100, 105, 57, 102, 58, 109, 49,
	27, 114, 115, 52, 25, 139, 101, 113, 28, 118,
	119, 120, 121, 122, 123, 124, 125, 126, 127, 128,
	129, 130, 131, 132, 133, 134, 135, 23, 174, 140,
	78, 79, 80, 137, 52, 51, 92, 136, 106, 104,
	146, 76, 77, 78, 79, 80, 21, 53, 52, 153,
	86, 156, 151, 48, 157, 82, 83, 182, 19, 22,
	152, 111, 4, 14, 56, 54, 188, 13, 26, 94,
	6, 34, 76, 77, 78, 79, 80, 66, 67, 33,
	23, 59, 162, 5, 164, 15, 21, 21, 145, 99,
	172, 173, 21, 37, 63, 55, 179, 178, 176, 22,
	22, 95, 138, 2, 0, 22, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 21, 0, 189, 0,
	23, 23, 187, 0, 191, 0, 23, 21, 21, 22,
	21, 0, 0, 192, 193, 0, 194, 0, 0, 0,
	22, 22, 0, 22, 0, 0, 0, 0, 0, 0,
	23, 0, 0, 74, 75, 82, 83, 73, 72, 0,
	0, 23, 23, 0, 23, 68, 69, 84, 85, 81,
	70, 71, 76, 77, 78, 79, 80, 0, 74, 75,
	82, 83, 73, 72, 0, 154, 0, 0, 0, 155,
	68, 69, 84, 85, 81, 70, 71, 76, 77, 78,
	79, 80, 0, 74, 75, 82, 83, 73, 72, 0,
	0, 0, 0, 0, 183, 68, 69, 84, 85, 81,
	70, 71, 76, 77, 78, 79, 80, 0, 74, 75,
	82, 83, 73, 72, 0, 0, 0, 0, 0, 143,
	68, 69, 84, 85, 81, 70, 71, 76, 77, 78,
	79, 80, 74, 75, 82, 83, 73, 72, 0, 0,
	0, 0, 150, 0, 68, 69, 84, 85, 81, 70,
	71, 76, 77, 78, 79, 80, 74, 75, 82, 83,
	73, 72, 0, 0, 0, 0, 144, 0, 68, 69,
	84, 85, 81, 70, 71, 76, 77, 78, 79, 80,
	74, 75, 82, 83, 73, 72, 0, 0, 0, 0,
	117, 0, 68, 69, 84, 85, 81, 70, 71, 76,
	77, 78, 79, 80, 74, 75, 82, 83, 73, 72,
	24, 0, 190, 30, 0, 32, 68, 69, 84, 85,
	81, 70, 71, 76, 77, 78, 79, 80, 27, 31,
	45, 43, 25, 0, 24, 0, 181, 30, 0, 32,
	20, 24, 38, 0, 30, 0, 32, 39, 41, 0,
	0, 0, 27, 31, 45, 43, 25, 40, 110, 27,
	31, 45, 43, 25, 0, 0, 38, 0, 0, 0,
	0, 39, 41, 38, 0, 103, 0, 0, 39, 41,
	0, 40, 0, 0, 0, 180, 0, 0, 40, 74,
	75, 82, 83, 73, 72, 24, 0, 0, 30, 0,
	32, 68, 69, 84, 85, 81, 70, 71, 76, 77,
	78, 79, 80, 27, 31, 45, 43, 25, 175, 24,
	0, 0, 30, 0, 32, 0, 24, 38, 0, 30,
	0, 32, 39, 41, 0, 0, 147, 27, 31, 45,
	43, 25, 40, 0, 27, 31, 45, 43, 25, 24,
	0, 38, 30, 0, 32, 0, 39, 41, 38, 0,
	0, 0, 0, 39, 41, 112, 40, 27, 31, 45,
	43, 25, 0, 40, 0, 0, 0, 0, 0, 0,
	0, 38, 0, 0, 0, 0, 39, 41, 0, 93,
	10, 8, 9, 0, 17, 24, 40, 18, 0, 11,
	12, 20, 7, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 27, 10, 8, 9, 25, 17, 0,
	0, 18, 0, 11, 12, 20, 7, 0, 0, 0,
	0, 0, 0, 0, 0, 149, 3, 27, 0, 0,
	52, 25, 10, 8, 9, 0, 17, 24, 0, 18,
	0, 11, 12, 20, 7, 0, 0, 0, 0, 0,
	166, 0, 0, 0, 0, 27, 0, 0, 0, 25,
	0, 0, 0, 0, 0, 0, 0, 74, 75, 82,
	83, 73, 72, 0, 0, 0, 0, 0, 3, 68,
	69, 84, 85, 81, 70, 71, 76, 77, 78, 79,
	80, 74, 75, 82, 83, 73, 72, 0, 0, 0,
	0, 0, 0, 0, 69, 84, 85, 81, 70, 71,
	76, 77, 78, 79, 80, 74, 75, 82, 83, 73,
	72, 0, 74, 75, 82, 83, 73, 72, 0, 84,
	85, 81, 70, 71, 76, 77, 78, 79, 80, 70,
	71, 76, 77, 78, 79, 80,
}
var yyPact = [...]int{

	-1000, 618, -1000, -1000, 18, -1000, -1000, 497, -1000, -1000,
	497, 497, 84, -1000, -1000, -1000, 9, 64, 76, 98,
	97, 3, -1000, -2, -23, 497, -1000, 120, -1000, 636,
	-1000, -1000, 81, -1000, -1000, 3, -1000, -1000, 497, 497,
	497, 497, 65, 520, -1000, -1000, 636, 636, -1000, -1000,
	-1000, 405, -1000, 497, 65, -37, -4, 381, 94, -1000,
	490, 497, -1000, -17, 315, -1000, -1000, -1000, 497, 497,
	497, 497, 497, 497, 497, 497, 497, 497, 497, 497,
	497, 497, 497, 497, 497, 497, -1000, -1000, -1000, -1000,
	-1000, 78, 38, -1000, 40, -48, -49, 242, 291, 466,
	3, -1000, 1, -1000, 566, 267, 78, 93, 497, 192,
	497, 120, -1000, -26, 636, 636, -1000, -1000, 660, 684,
	92, 92, 92, 92, 92, 92, 48, 48, -1000, -1000,
	-1000, 691, 61, 61, 691, 691, -1000, -1000, -27, -1000,
	-1000, 497, 497, 497, 590, 33, 448, -1000, -1000, -1000,
	590, -1000, -6, 636, 120, 412, 363, -1000, 497, -1000,
	90, -1000, 217, 636, 636, -1000, -1000, -1000, 0, -1000,
	-1000, -1000, -19, -25, 590, -1000, 119, 497, -1000, 339,
	-1000, -1000, -1000, 497, -1000, 590, 590, -1000, 590, 636,
	-1000, 636, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 46, 5, 163, 35, 162, 6, 161, 155, 0,
	14, 2, 153, 1, 4, 25, 22, 149, 148, 18,
	7, 17, 145, 143, 13, 128, 141, 139, 43, 131,
	129,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 15, 3, 3, 3, 3, 20,
	20, 20, 20, 20, 20, 23, 23, 23, 14, 14,
	14, 14, 11, 11, 10, 10, 10, 17, 17, 18,
	18, 16, 16, 16, 16, 16, 19, 19, 24, 24,
	22, 21, 21, 21, 21, 21, 21, 21, 21, 4,
	4, 4, 4, 4, 4, 5, 5, 6, 6, 7,
	7, 8, 8, 8, 8, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 12, 13, 13, 13, 13,
	25, 26, 26, 27, 28, 28, 29, 29, 30, 30,
	30, 30,
}
var yyR2 = [...]int{

	0, 0, 2, 3, 2, 1, 2, 1, 1, 1,
	1, 2, 1, 1, 1, 1, 1, 1, 2, 1,
	1, 3, 1, 1, 2, 5, 4, 2, 1, 2,
	1, 2, 5, 7, 7, 6, 5, 7, 1, 2,
	4, 1, 2, 1, 1, 2, 1, 2, 2, 1,
	4, 6, 5, 5, 3, 1, 3, 1, 3, 3,
	5, 1, 3, 5, 3, 1, 1, 2, 1, 1,
	1, 1, 1, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 2, 2, 2, 2, 1, 1, 3, 1, 3,
	2, 2, 3, 3, 2, 3, 2, 3, 1, 2,
	1, 2,
}
var yyChk = [...]int{

	-1000, -1, -3, 50, -21, -23, -15, 16, 5, 6,
	4, 13, 14, -16, -19, -22, -14, 8, 11, -24,
	15, -13, -10, -4, 9, 31, -25, 27, 50, -9,
	12, 28, 14, -27, -29, -13, -10, -12, 41, 46,
	56, 47, -24, 30, -4, 29, -9, -9, 29, 50,
	-2, 31, 30, 31, 27, -8, 27, 52, 54, -26,
	31, 51, 56, -25, -9, -11, 17, 18, 33, 34,
	38, 39, 26, 25, 21, 22, 40, 41, 42, 43,
	44, 37, 23, 24, 35, 36, 29, -9, -9, -9,
	-9, -28, 31, 49, -30, -7, -6, -9, -9, -17,
	-13, -10, -14, 50, -1, -9, -28, 58, 51, -9,
	57, 27, 55, -6, -9, -9, 55, 55, -9, -9,
	-9, -9, -9, -9, -9, -9, -9, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, -2, 55, -5, 27,
	49, 58, 58, 57, 55, -18, -9, 50, 50, 49,
	55, -2, 27, -9, 53, 57, -9, -11, 58, 55,
	58, 55, -9, -9, -9, -20, 50, -15, -21, -16,
	-19, -2, -14, -2, 55, 50, -20, 51, -11, -9,
	53, 53, 27, 57, 50, 55, 55, -20, 7, -9,
	53, -9, -20, -20, -20,
}
var yyDef = [...]int{

	1, -2, 2, 5, 0, 7, 8, 41, 43, 44,
	0, 46, 0, 15, 16, 17, 0, 0, 0, 0,
	0, 19, 20, 96, 38, 0, 98, 49, 6, 42,
	65, 66, 0, 68, 69, 70, 71, 72, 0, 0,
	0, 0, 0, 0, 96, 95, 45, 47, 48, 4,
	31, 0, 1, 0, 0, 18, 61, 0, 0, 100,
	0, 0, 39, 98, 0, 24, 22, 23, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 67, 91, 92, 93,
	94, 0, 0, 106, 0, 108, 110, 57, 0, 0,
	-2, -2, 0, 28, 0, 0, 0, 0, 0, 0,
	0, 54, 101, 0, 57, 21, 97, 99, 73, 74,
	75, 76, 77, 78, 79, 80, 81, 82, 83, 84,
	85, 86, 87, 88, 89, 90, 103, 104, 0, 55,
	107, 109, 111, 0, 0, 0, 0, 30, 27, 3,
	0, 40, 64, 62, 50, 0, 0, 26, 0, 102,
	0, 105, 0, 58, 59, 32, 9, 10, 0, 12,
	13, 14, 0, 0, 0, 29, 36, 0, 25, 0,
	52, 53, 56, 0, 11, 0, 0, 35, 0, 63,
	51, 60, 33, 34, 37,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 56, 3, 47, 3, 44, 36, 3,
	31, 55, 42, 40, 58, 41, 54, 43, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 57, 50,
	39, 51, 38, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 32, 3, 3, 3, 3, 3,
	3, 52, 3, 53, 37, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 30, 35, 49, 46,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 33, 34,
	45, 48,
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
		//line parser.go.y:74
		{
			yyVAL.expr = CNode("chain")
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:80
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[2].expr)
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 3:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:88
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:91
		{
			if yyDollar[1].expr.isIsolatedDupCall() {
				yyDollar[1].expr.Cx(2).C()[0] = NNode(0.0)
			}
			yyVAL.expr = yyDollar[1].expr
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:99
		{
			yyVAL.expr = CNode()
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:100
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
		//line parser.go.y:102
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:105
		{
			yyVAL.expr = CNode()
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:106
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:107
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:108
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:109
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:110
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:113
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:114
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:115
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 18:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:118
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:119
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 20:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:120
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:121
		{
			yyVAL.expr = CNode("move", yyDollar[1].expr, yyDollar[3].expr)
			if yyDollar[1].expr.Cn() > 0 && yyDollar[1].expr.Cx(0).S() == "load" {
				yyVAL.expr = CNode("store", yyDollar[1].expr.Cx(1), yyDollar[1].expr.Cx(2), yyDollar[3].expr)
			}
			if c := yyDollar[1].expr.S(); c != "" && yyDollar[1].expr.Type == NTAtom {
				if a, b, s := yyDollar[3].expr.isSimpleAddSub(); a == c {
					yyDollar[3].expr.Cx(2).Value = yyDollar[3].expr.Cx(2).N() * s
					yyVAL.expr = CNode("inc", yyDollar[1].expr, yyDollar[3].expr.Cx(2))
					yyVAL.expr.Cx(1).SetPos(yyDollar[1].expr)
				} else if b == c {
					yyDollar[3].expr.Cx(1).Value = yyDollar[3].expr.Cx(1).N() * s
					yyVAL.expr = CNode("inc", yyDollar[1].expr, yyDollar[3].expr.Cx(1))
					yyVAL.expr.Cx(1).SetPos(yyDollar[1].expr)
				}
			}
			yyVAL.expr.Cx(0).SetPos(yyDollar[1].expr)
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:141
		{
			yyVAL.expr = NNode("1")
		}
	case 23:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:142
		{
			yyVAL.expr = NNode("-1")
		}
	case 24:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:145
		{
			yyVAL.expr = CNode("inc", ANode(yyDollar[1].token).setPos(yyDollar[1].token), yyDollar[2].expr)
		}
	case 25:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:146
		{
			yyVAL.expr = CNode("store", yyDollar[1].expr, yyDollar[3].expr, CNode("+", CNode("load", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr), yyDollar[5].expr).setPos0(yyDollar[1].expr))
		}
	case 26:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:147
		{
			yyVAL.expr = CNode("store", yyDollar[1].expr, yyDollar[3].token, CNode("+", CNode("load", yyDollar[1].expr, yyDollar[3].token).setPos0(yyDollar[1].expr), yyDollar[4].expr).setPos0(yyDollar[1].expr))
		}
	case 27:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:150
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:151
		{
			yyVAL.expr = CNode("chain")
		}
	case 29:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:154
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:155
		{
			yyVAL.expr = NNode("1")
		}
	case 31:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:158
		{
			yyVAL.expr = CNode("for", NNode("1"), CNode(), yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 32:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:161
		{
			yyVAL.expr = CNode("for", yyDollar[3].expr, CNode(), yyDollar[5].expr).setPos0(yyDollar[1].token)
		}
	case 33:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:164
		{
			yyVAL.expr = yyDollar[3].expr.Cappend(CNode("for", yyDollar[4].expr, CNode("chain", yyDollar[5].expr), yyDollar[7].expr)).setPos0(yyDollar[1].token)
		}
	case 34:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:167
		{
			yyVAL.expr = yyDollar[3].expr.Cappend(CNode("for", yyDollar[4].expr, yyDollar[5].expr, yyDollar[7].expr)).setPos0(yyDollar[1].token)
		}
	case 35:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:170
		{
			yyVAL.expr = yyDollar[3].expr.Cappend(CNode("for", yyDollar[4].expr, CNode(), yyDollar[6].expr)).setPos0(yyDollar[1].token)
		}
	case 36:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:175
		{
			yyVAL.expr = CNode("if", yyDollar[3].expr, yyDollar[5].expr, CNode())
		}
	case 37:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:176
		{
			yyVAL.expr = CNode("if", yyDollar[3].expr, yyDollar[5].expr, yyDollar[7].expr)
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:179
		{
			yyVAL.str = "func"
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:180
		{
			yyVAL.str = "safefunc"
		}
	case 40:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:183
		{
			funcname := ANode(yyDollar[2].token)
			yyVAL.expr = CNode(
				"chain",
				CNode("set", funcname, NilNode()).setPos0(yyDollar[2].token),
				CNode("move", funcname,
					CNode(yyDollar[1].str, funcname, yyDollar[3].expr, yyDollar[4].expr).setPos0(yyDollar[2].token),
				).setPos0(yyDollar[2].token),
			)
		}
	case 41:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:195
		{
			yyVAL.expr = CNode("yield").setPos0(yyDollar[1].token)
		}
	case 42:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:196
		{
			yyVAL.expr = CNode("yield", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 43:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:197
		{
			yyVAL.expr = CNode("break").setPos0(yyDollar[1].token)
		}
	case 44:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:198
		{
			yyVAL.expr = CNode("continue").setPos0(yyDollar[1].token)
		}
	case 45:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:199
		{
			yyVAL.expr = CNode("assert", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:200
		{
			yyVAL.expr = CNode("ret").setPos0(yyDollar[1].token)
		}
	case 47:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:201
		{
			if yyDollar[2].expr.isIsolatedDupCall() && yyDollar[2].expr.Cx(2).Cx(2).N() == 1 {
				yyDollar[2].expr.Cx(2).C()[2] = NNode(2.0)
			}
			yyVAL.expr = CNode("ret", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 48:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:207
		{
			path := filepath.Join(filepath.Dir(yyDollar[1].token.Pos.Source), yyDollar[2].token.Str)
			yyVAL.expr = yylex.(*Lexer).loadFile(path)
		}
	case 49:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:213
		{
			yyVAL.expr = ANode(yyDollar[1].token).setPos(yyDollar[1].token)
		}
	case 50:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:214
		{
			yyVAL.expr = CNode("load", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 51:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:215
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, yyDollar[3].expr, yyDollar[5].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 52:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:216
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, yyDollar[3].expr, NNode("-1")).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 53:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:217
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, NNode("0"), yyDollar[4].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:218
		{
			yyVAL.expr = CNode("load", yyDollar[1].expr, SNode(yyDollar[3].token.Str)).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 55:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:221
		{
			yyVAL.expr = CNode(yyDollar[1].token.Str)
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:222
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(ANode(yyDollar[3].token))
		}
	case 57:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:225
		{
			yyVAL.expr = CNode(yyDollar[1].expr)
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:226
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[3].expr)
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:229
		{
			yyVAL.expr = CNode(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 60:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:230
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[3].expr).Cappend(yyDollar[5].expr)
		}
	case 61:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:233
		{
			yyVAL.expr = CNode("chain", CNode("set", ANode(yyDollar[1].token), NilNode()).setPos0(yyDollar[1].token))
		}
	case 62:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:234
		{
			yyVAL.expr = CNode("chain", CNode("set", ANode(yyDollar[1].token), yyDollar[3].expr).setPos0(yyDollar[1].token))
		}
	case 63:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:235
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(CNode("set", ANode(yyDollar[3].token), yyDollar[5].expr).setPos0(yyDollar[1].expr))
		}
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:236
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(CNode("set", ANode(yyDollar[3].token), NilNode()).setPos0(yyDollar[1].expr))
		}
	case 65:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:239
		{
			yyVAL.expr = NilNode().SetPos(yyDollar[1].token)
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:240
		{
			yyVAL.expr = NNode(yyDollar[1].token.Str).SetPos(yyDollar[1].token)
		}
	case 67:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:241
		{
			yyVAL.expr = yylex.(*Lexer).loadFile(filepath.Join(filepath.Dir(yyDollar[1].token.Pos.Source), yyDollar[2].token.Str))
		}
	case 68:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:242
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 69:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:243
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 70:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:244
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 71:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:245
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:246
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:247
		{
			yyVAL.expr = CNode("or", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:248
		{
			yyVAL.expr = CNode("and", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 75:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:249
		{
			yyVAL.expr = CNode("<", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:250
		{
			yyVAL.expr = CNode("<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 77:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:251
		{
			yyVAL.expr = CNode("<=", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 78:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:252
		{
			yyVAL.expr = CNode("<=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 79:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:253
		{
			yyVAL.expr = CNode("==", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 80:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:254
		{
			yyVAL.expr = CNode("!=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 81:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:255
		{
			yyVAL.expr = CNode("+", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 82:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:256
		{
			yyVAL.expr = CNode("-", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 83:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:257
		{
			yyVAL.expr = CNode("*", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 84:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:258
		{
			yyVAL.expr = CNode("/", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:259
		{
			yyVAL.expr = CNode("%", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 86:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:260
		{
			yyVAL.expr = CNode("^", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 87:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:261
		{
			yyVAL.expr = CNode("<<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 88:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:262
		{
			yyVAL.expr = CNode(">>", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 89:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:263
		{
			yyVAL.expr = CNode("|", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 90:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:264
		{
			yyVAL.expr = CNode("&", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 91:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:265
		{
			yyVAL.expr = CNode("-", NNode(0.0), yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 92:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:266
		{
			yyVAL.expr = CNode("~", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 93:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:267
		{
			yyVAL.expr = CNode("!", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 94:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:268
		{
			yyVAL.expr = CNode("#", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 95:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:271
		{
			yyVAL.expr = SNode(yyDollar[1].token.Str).SetPos(yyDollar[1].token)
		}
	case 96:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:274
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 97:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:275
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 98:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:276
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 99:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:277
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 100:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:280
		{
			switch yyDollar[1].expr.S() {
			case "copy":
				switch yyDollar[2].expr.Cn() {
				case 0:
					yyVAL.expr = CNode("call", yyDollar[1].expr, CNode(NNode("1"), NNode("1"), NNode("1")))
				case 1:
					yyVAL.expr = CNode("call", yyDollar[1].expr, CNode(NNode("1"), yyDollar[2].expr.Cx(0), NNode("0")))
				default:
					p := yyDollar[2].expr.Cx(1)
					if p.Type != NTCompound && p.Type != NTAtom {
						yylex.(*Lexer).Error("invalid argument for S")
					}
					yyVAL.expr = CNode("call", yyDollar[1].expr, CNode(NNode("1"), yyDollar[2].expr.Cx(0), p))
				}
			case "typeof":
				switch yyDollar[2].expr.Cn() {
				case 0:
					yylex.(*Lexer).Error("typeof takes at least 1 argument")
				case 1:
					yyVAL.expr = CNode("call", yyDollar[1].expr, CNode(yyDollar[2].expr.Cx(0), NNode("255")))
				default:
					x, _ := yyDollar[2].expr.Cx(1).Value.(string)
					if ti, ok := typesLookup[x]; ok {
						yyVAL.expr = CNode("call", yyDollar[1].expr, CNode(yyDollar[2].expr.Cx(0), NNode(ti)))
					} else {
						yylex.(*Lexer).Error("invalid typename in typeof")
					}
				}
			case "addressof":
				if yyDollar[2].expr.Cn() != 1 {
					yylex.(*Lexer).Error("addressof takes 1 argument")
				}
				if yyDollar[2].expr.Cx(0).Type != NTAtom {
					yylex.(*Lexer).Error("addressof can only get the address of a variable")
				}
				yyVAL.expr = CNode("call", yyDollar[1].expr, yyDollar[2].expr)
			case "len":
				switch yyDollar[2].expr.Cn() {
				case 0:
					yylex.(*Lexer).Error("len takes 1 argument")
				default:
					yyVAL.expr = CNode("call", yyDollar[1].expr, yyDollar[2].expr)
				}
			default:
				yyVAL.expr = CNode("call", yyDollar[1].expr, yyDollar[2].expr)
			}
			yyVAL.expr.Cx(0).SetPos(yyDollar[1].expr)
		}
	case 101:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:331
		{
			yyVAL.expr = CNode()
		}
	case 102:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:332
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 103:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:335
		{
			yyVAL.expr = CNode(yyDollar[1].str, "<a>", yyDollar[2].expr, yyDollar[3].expr).setPos0(yyDollar[2].expr)
		}
	case 104:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:338
		{
			yyVAL.expr = CNode()
		}
	case 105:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:339
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 106:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:342
		{
			yyVAL.expr = CNode("map", CNode()).setPos0(yyDollar[1].token)
		}
	case 107:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:343
		{
			yyVAL.expr = yyDollar[2].expr.setPos0(yyDollar[1].token)
		}
	case 108:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:346
		{
			yyVAL.expr = CNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 109:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:347
		{
			yyVAL.expr = CNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 110:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:348
		{
			yyVAL.expr = CNode("array", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 111:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:349
		{
			yyVAL.expr = CNode("array", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	}
	goto yystack /* stack new state and value */
}
