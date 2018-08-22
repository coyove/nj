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
const TIf = 57352
const TNil = 57353
const TReturn = 57354
const TRequire = 57355
const TVar = 57356
const TYield = 57357
const TAddAdd = 57358
const TSubSub = 57359
const TAddEq = 57360
const TSubEq = 57361
const TEqeq = 57362
const TNeq = 57363
const TLsh = 57364
const TRsh = 57365
const TURsh = 57366
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
	"TURsh",
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
	"','",
	"'!'",
	"':'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.go.y:359

var typesLookup = map[string]string{
	"nil": "0", "number": "1", "string": "2", "map": "4", "closure": "6", "generic": "7",
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 102,
	50, 19,
	-2, 71,
	-1, 103,
	50, 20,
	-2, 72,
}

const yyPrivate = 57344

const yyLast = 783

var yyAct = [...]int{

	167, 35, 65, 21, 16, 175, 97, 169, 29, 66,
	67, 46, 47, 42, 62, 19, 36, 172, 22, 4,
	174, 145, 14, 50, 144, 173, 64, 13, 171, 142,
	6, 165, 164, 163, 162, 109, 44, 191, 23, 88,
	89, 90, 91, 60, 98, 190, 118, 1, 92, 151,
	193, 182, 99, 102, 107, 110, 104, 140, 111, 61,
	189, 116, 117, 152, 57, 49, 58, 115, 103, 120,
	121, 122, 123, 124, 125, 126, 127, 128, 129, 130,
	131, 132, 133, 134, 135, 136, 137, 138, 23, 20,
	28, 143, 76, 77, 78, 79, 80, 93, 139, 53,
	106, 149, 27, 108, 52, 52, 25, 87, 21, 52,
	51, 157, 48, 160, 155, 187, 161, 180, 156, 113,
	19, 194, 56, 22, 4, 54, 26, 14, 66, 67,
	178, 95, 13, 34, 33, 6, 78, 79, 80, 59,
	5, 15, 148, 23, 100, 166, 37, 168, 55, 21,
	21, 96, 63, 176, 177, 141, 21, 82, 83, 84,
	184, 183, 181, 2, 22, 22, 0, 0, 0, 0,
	0, 22, 0, 0, 0, 76, 77, 78, 79, 80,
	21, 0, 0, 195, 23, 23, 192, 0, 0, 197,
	0, 23, 21, 21, 200, 22, 21, 0, 198, 199,
	0, 0, 201, 0, 21, 0, 0, 22, 22, 0,
	203, 22, 0, 0, 0, 23, 0, 0, 0, 22,
	0, 0, 0, 0, 0, 0, 0, 23, 23, 0,
	0, 23, 74, 75, 82, 83, 84, 73, 72, 23,
	0, 0, 0, 0, 0, 68, 69, 85, 86, 81,
	70, 71, 76, 77, 78, 79, 80, 0, 74, 75,
	82, 83, 84, 73, 72, 158, 0, 0, 0, 0,
	159, 68, 69, 85, 86, 81, 70, 71, 76, 77,
	78, 79, 80, 0, 74, 75, 82, 83, 84, 73,
	72, 0, 0, 0, 0, 0, 188, 68, 69, 85,
	86, 81, 70, 71, 76, 77, 78, 79, 80, 0,
	74, 75, 82, 83, 84, 73, 72, 0, 0, 0,
	0, 0, 146, 68, 69, 85, 86, 81, 70, 71,
	76, 77, 78, 79, 80, 74, 75, 82, 83, 84,
	73, 72, 0, 0, 0, 202, 0, 0, 68, 69,
	85, 86, 81, 70, 71, 76, 77, 78, 79, 80,
	74, 75, 82, 83, 84, 73, 72, 0, 0, 0,
	154, 0, 0, 68, 69, 85, 86, 81, 70, 71,
	76, 77, 78, 79, 80, 74, 75, 82, 83, 84,
	73, 72, 0, 0, 0, 147, 0, 0, 68, 69,
	85, 86, 81, 70, 71, 76, 77, 78, 79, 80,
	74, 75, 82, 83, 84, 73, 72, 0, 0, 0,
	119, 0, 0, 68, 69, 85, 86, 81, 70, 71,
	76, 77, 78, 79, 80, 74, 75, 82, 83, 84,
	73, 72, 24, 196, 30, 0, 32, 0, 68, 69,
	85, 86, 81, 70, 71, 76, 77, 78, 79, 80,
	27, 31, 45, 43, 25, 24, 0, 30, 186, 32,
	20, 0, 0, 0, 38, 0, 0, 0, 0, 39,
	41, 0, 0, 101, 31, 45, 43, 25, 0, 0,
	40, 112, 0, 0, 0, 0, 24, 38, 30, 0,
	32, 0, 39, 41, 0, 24, 105, 30, 0, 32,
	0, 0, 0, 40, 27, 31, 45, 43, 25, 0,
	0, 0, 0, 27, 31, 45, 43, 25, 38, 0,
	0, 0, 0, 39, 41, 0, 0, 38, 0, 0,
	185, 0, 39, 41, 40, 74, 75, 82, 83, 84,
	73, 72, 24, 40, 30, 0, 32, 0, 68, 69,
	85, 86, 81, 70, 71, 76, 77, 78, 79, 80,
	27, 31, 45, 43, 25, 179, 24, 0, 30, 0,
	32, 0, 0, 24, 38, 30, 0, 32, 0, 39,
	41, 0, 0, 150, 27, 31, 45, 43, 25, 0,
	40, 27, 31, 45, 43, 25, 0, 0, 38, 0,
	0, 0, 0, 39, 41, 38, 0, 0, 0, 0,
	39, 41, 114, 94, 40, 0, 0, 0, 10, 8,
	9, 40, 17, 24, 18, 0, 11, 12, 20, 7,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 27, 10, 8, 9, 25, 17, 0, 18, 0,
	11, 12, 20, 7, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 153, 3, 27, 0, 0, 52, 25,
	0, 0, 0, 0, 0, 0, 74, 75, 82, 83,
	84, 73, 72, 0, 0, 0, 0, 0, 170, 68,
	69, 85, 86, 81, 70, 71, 76, 77, 78, 79,
	80, 10, 8, 9, 0, 17, 24, 18, 0, 11,
	12, 20, 7, 0, 0, 74, 75, 82, 83, 84,
	73, 72, 0, 0, 27, 0, 0, 0, 25, 69,
	85, 86, 81, 70, 71, 76, 77, 78, 79, 80,
	74, 75, 82, 83, 84, 73, 72, 3, 74, 75,
	82, 83, 84, 73, 72, 85, 86, 81, 70, 71,
	76, 77, 78, 79, 80, 0, 70, 71, 76, 77,
	78, 79, 80,
}
var yyPact = [...]int{

	-1000, 707, -1000, -1000, 40, -1000, -1000, 496, -1000, -1000,
	496, 496, 83, -1000, -1000, -1000, 15, 79, 68, 98,
	95, 12, -1000, 8, -43, 496, -1000, 112, -1000, 666,
	-1000, -1000, 78, -1000, -1000, 12, -1000, -1000, 496, 496,
	496, 496, 66, 574, -1000, -1000, 666, 666, -1000, -1000,
	-1000, 456, -1000, 496, 66, -21, 4, 433, 92, -1000,
	567, 496, -1000, -9, 365, -1000, -1000, -1000, 496, 496,
	496, 496, 496, 496, 496, 496, 496, 496, 496, 496,
	496, 496, 496, 496, 496, 496, 496, -1000, -1000, -1000,
	-1000, -1000, 74, 2, -1000, 42, -32, -35, 264, 340,
	543, -7, 12, -1000, 13, -1000, 624, 315, 74, 91,
	496, 212, 496, 112, -1000, -22, 666, 666, -1000, -1000,
	705, 730, 135, 135, 135, 135, 135, 135, 94, 94,
	-1000, -1000, -1000, 738, 52, 52, 52, 738, 738, -1000,
	-1000, -24, -1000, -1000, 496, 496, 496, 648, 75, 525,
	-1000, 90, -1000, -1000, 648, -1000, 0, 666, 112, 487,
	415, -1000, 496, -1000, 88, -1000, 238, 666, 666, -1000,
	-1000, -1000, 10, -1000, -1000, -1000, -10, -18, 648, -1000,
	-1, 114, 496, -1000, 390, -1000, -1000, -1000, 496, -1000,
	648, 648, -1000, 496, 648, 666, -1000, 666, -1000, -1000,
	290, -1000, 648, -1000,
}
var yyPgo = [...]int{

	0, 47, 5, 163, 36, 155, 6, 151, 148, 0,
	16, 2, 146, 1, 4, 28, 25, 144, 142, 20,
	7, 17, 141, 140, 13, 126, 139, 134, 48, 133,
	131,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 15, 3, 3, 3, 3, 20,
	20, 20, 20, 20, 20, 23, 23, 23, 14, 14,
	14, 14, 11, 11, 10, 10, 10, 17, 17, 18,
	18, 16, 16, 16, 16, 16, 16, 19, 19, 24,
	24, 22, 21, 21, 21, 21, 21, 21, 21, 21,
	4, 4, 4, 4, 4, 4, 5, 5, 6, 6,
	7, 7, 8, 8, 8, 8, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 12, 13, 13,
	13, 13, 25, 26, 26, 27, 28, 28, 29, 29,
	30, 30, 30, 30,
}
var yyR2 = [...]int{

	0, 0, 2, 3, 2, 1, 2, 1, 1, 1,
	1, 2, 1, 1, 1, 1, 1, 1, 2, 1,
	1, 3, 1, 1, 2, 5, 4, 2, 1, 2,
	1, 2, 5, 7, 7, 6, 9, 5, 7, 1,
	2, 4, 1, 2, 1, 1, 2, 1, 2, 2,
	1, 4, 6, 5, 5, 3, 1, 3, 1, 3,
	3, 5, 1, 3, 5, 3, 1, 1, 2, 1,
	1, 1, 1, 1, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 2, 2, 2, 2, 1, 1, 3,
	1, 3, 2, 2, 3, 3, 2, 3, 2, 3,
	1, 2, 1, 2,
}
var yyChk = [...]int{

	-1000, -1, -3, 50, -21, -23, -15, 15, 5, 6,
	4, 12, 13, -16, -19, -22, -14, 8, 10, -24,
	14, -13, -10, -4, 9, 31, -25, 27, 50, -9,
	11, 28, 13, -27, -29, -13, -10, -12, 41, 46,
	57, 47, -24, 30, -4, 29, -9, -9, 29, 50,
	-2, 31, 30, 31, 27, -8, 27, 52, 54, -26,
	31, 51, 57, -25, -9, -11, 16, 17, 33, 34,
	38, 39, 26, 25, 20, 21, 40, 41, 42, 43,
	44, 37, 22, 23, 24, 35, 36, 29, -9, -9,
	-9, -9, -28, 31, 49, -30, -7, -6, -9, -9,
	-17, 27, -13, -10, -14, 50, -1, -9, -28, 56,
	51, -9, 58, 27, 55, -6, -9, -9, 55, 55,
	-9, -9, -9, -9, -9, -9, -9, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, -9, -9, -9, -2,
	55, -5, 27, 49, 56, 56, 58, 55, -18, -9,
	50, 56, 50, 49, 55, -2, 27, -9, 53, 58,
	-9, -11, 56, 55, 56, 55, -9, -9, -9, -20,
	50, -15, -21, -16, -19, -2, -14, -2, 55, 50,
	27, -20, 51, -11, -9, 53, 53, 27, 58, 50,
	55, 55, -20, 51, 7, -9, 53, -9, -20, -20,
	-9, -20, 55, -20,
}
var yyDef = [...]int{

	1, -2, 2, 5, 0, 7, 8, 42, 44, 45,
	0, 47, 0, 15, 16, 17, 0, 0, 0, 0,
	0, 19, 20, 98, 39, 0, 100, 50, 6, 43,
	66, 67, 0, 69, 70, 71, 72, 73, 0, 0,
	0, 0, 0, 0, 98, 97, 46, 48, 49, 4,
	31, 0, 1, 0, 0, 18, 62, 0, 0, 102,
	0, 0, 40, 100, 0, 24, 22, 23, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 68, 93, 94,
	95, 96, 0, 0, 108, 0, 110, 112, 58, 0,
	0, 50, -2, -2, 0, 28, 0, 0, 0, 0,
	0, 0, 0, 55, 103, 0, 58, 21, 99, 101,
	74, 75, 76, 77, 78, 79, 80, 81, 82, 83,
	84, 85, 86, 87, 88, 89, 90, 91, 92, 105,
	106, 0, 56, 109, 111, 113, 0, 0, 0, 0,
	30, 0, 27, 3, 0, 41, 65, 63, 51, 0,
	0, 26, 0, 104, 0, 107, 0, 59, 60, 32,
	9, 10, 0, 12, 13, 14, 0, 0, 0, 29,
	0, 37, 0, 25, 0, 53, 54, 57, 0, 11,
	0, 0, 35, 0, 0, 64, 52, 61, 33, 34,
	0, 38, 0, 36,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 57, 3, 47, 3, 44, 36, 3,
	31, 55, 42, 40, 56, 41, 54, 43, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 58, 50,
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
			if c := yyDollar[1].expr.S(); c != "" && yyDollar[1].expr.Type == Natom {
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
			yyVAL.expr = NNode(1.0)
		}
	case 23:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:142
		{
			yyVAL.expr = NNode(-1.0)
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
			yyVAL.expr = NNode(1.0)
		}
	case 31:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:158
		{
			yyVAL.expr = CNode("for", NNode(1.0), CNode(), yyDollar[2].expr).setPos0(yyDollar[1].token)
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
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.go.y:173
		{
			yyVAL.expr = CNode("call", "copy", CNode(
				NNode(0.0),
				yyDollar[7].expr,
				CNode("func", "<anony-map-iter-callback>", CNode(yyDollar[3].token.Str, yyDollar[5].token.Str), yyDollar[9].expr),
			))
		}
	case 37:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:182
		{
			yyVAL.expr = CNode("if", yyDollar[3].expr, yyDollar[5].expr, CNode())
		}
	case 38:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:183
		{
			yyVAL.expr = CNode("if", yyDollar[3].expr, yyDollar[5].expr, yyDollar[7].expr)
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:186
		{
			yyVAL.str = "func"
		}
	case 40:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:187
		{
			yyVAL.str = "safefunc"
		}
	case 41:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:190
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
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:202
		{
			yyVAL.expr = CNode("yield").setPos0(yyDollar[1].token)
		}
	case 43:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:203
		{
			yyVAL.expr = CNode("yield", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 44:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:204
		{
			yyVAL.expr = CNode("break").setPos0(yyDollar[1].token)
		}
	case 45:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:205
		{
			yyVAL.expr = CNode("continue").setPos0(yyDollar[1].token)
		}
	case 46:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:206
		{
			yyVAL.expr = CNode("assert", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 47:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:207
		{
			yyVAL.expr = CNode("ret").setPos0(yyDollar[1].token)
		}
	case 48:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:208
		{
			if yyDollar[2].expr.isIsolatedDupCall() && yyDollar[2].expr.Cx(2).Cx(2).N() == 1 {
				yyDollar[2].expr.Cx(2).C()[2] = NNode(2.0)
			}
			yyVAL.expr = CNode("ret", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:214
		{
			path := filepath.Join(filepath.Dir(yyDollar[1].token.Pos.Source), yyDollar[2].token.Str)
			yyVAL.expr = yylex.(*Lexer).loadFile(path)
		}
	case 50:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:220
		{
			yyVAL.expr = ANode(yyDollar[1].token).setPos(yyDollar[1].token)
		}
	case 51:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:221
		{
			yyVAL.expr = CNode("load", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 52:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:222
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, yyDollar[3].expr, yyDollar[5].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 53:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:223
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, yyDollar[3].expr, NNode("-1")).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 54:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:224
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, NNode("0"), yyDollar[4].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:225
		{
			yyVAL.expr = CNode("load", yyDollar[1].expr, SNode(yyDollar[3].token.Str)).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 56:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:228
		{
			yyVAL.expr = CNode(yyDollar[1].token.Str)
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:229
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(ANode(yyDollar[3].token))
		}
	case 58:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:232
		{
			yyVAL.expr = CNode(yyDollar[1].expr)
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:233
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[3].expr)
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:236
		{
			yyVAL.expr = CNode(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 61:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:237
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[3].expr).Cappend(yyDollar[5].expr)
		}
	case 62:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:240
		{
			yyVAL.expr = CNode("chain", CNode("set", ANode(yyDollar[1].token), NilNode()).setPos0(yyDollar[1].token))
		}
	case 63:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:241
		{
			yyVAL.expr = CNode("chain", CNode("set", ANode(yyDollar[1].token), yyDollar[3].expr).setPos0(yyDollar[1].token))
		}
	case 64:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:242
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(CNode("set", ANode(yyDollar[3].token), yyDollar[5].expr).setPos0(yyDollar[1].expr))
		}
	case 65:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:243
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(CNode("set", ANode(yyDollar[3].token), NilNode()).setPos0(yyDollar[1].expr))
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:246
		{
			yyVAL.expr = NilNode().SetPos(yyDollar[1].token)
		}
	case 67:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:247
		{
			yyVAL.expr = NNode(yyDollar[1].token.Str).SetPos(yyDollar[1].token)
		}
	case 68:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:248
		{
			yyVAL.expr = yylex.(*Lexer).loadFile(filepath.Join(filepath.Dir(yyDollar[1].token.Pos.Source), yyDollar[2].token.Str))
		}
	case 69:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:249
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 70:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:250
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 71:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:251
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:252
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 73:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:253
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:254
		{
			yyVAL.expr = CNode("or", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 75:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:255
		{
			yyVAL.expr = CNode("and", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:256
		{
			yyVAL.expr = CNode("<", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 77:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:257
		{
			yyVAL.expr = CNode("<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 78:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:258
		{
			yyVAL.expr = CNode("<=", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 79:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:259
		{
			yyVAL.expr = CNode("<=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 80:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:260
		{
			yyVAL.expr = CNode("==", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 81:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:261
		{
			yyVAL.expr = CNode("!=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 82:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:262
		{
			yyVAL.expr = CNode("+", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 83:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:263
		{
			yyVAL.expr = CNode("-", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 84:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:264
		{
			yyVAL.expr = CNode("*", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:265
		{
			yyVAL.expr = CNode("/", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 86:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:266
		{
			yyVAL.expr = CNode("%", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 87:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:267
		{
			yyVAL.expr = CNode("^", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 88:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:268
		{
			yyVAL.expr = CNode("<<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 89:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:269
		{
			yyVAL.expr = CNode(">>", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 90:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:270
		{
			yyVAL.expr = CNode(">>>", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 91:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:271
		{
			yyVAL.expr = CNode("|", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 92:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:272
		{
			yyVAL.expr = CNode("&", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 93:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:273
		{
			yyVAL.expr = CNode("-", NNode(0.0), yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 94:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:274
		{
			yyVAL.expr = CNode("~", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 95:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:275
		{
			yyVAL.expr = CNode("!", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 96:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:276
		{
			yyVAL.expr = CNode("#", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 97:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:279
		{
			yyVAL.expr = SNode(yyDollar[1].token.Str).SetPos(yyDollar[1].token)
		}
	case 98:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:282
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 99:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:283
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 100:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:284
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 101:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:285
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 102:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:288
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
					if p.Type != Ncompound && p.Type != Natom {
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
				if yyDollar[2].expr.Cx(0).Type != Natom {
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
	case 103:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:339
		{
			yyVAL.expr = CNode()
		}
	case 104:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:340
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 105:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:343
		{
			yyVAL.expr = CNode(yyDollar[1].str, "<a>", yyDollar[2].expr, yyDollar[3].expr).setPos0(yyDollar[2].expr)
		}
	case 106:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:346
		{
			yyVAL.expr = CNode()
		}
	case 107:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:347
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 108:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:350
		{
			yyVAL.expr = CNode("map", CNode()).setPos0(yyDollar[1].token)
		}
	case 109:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:351
		{
			yyVAL.expr = yyDollar[2].expr.setPos0(yyDollar[1].token)
		}
	case 110:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:354
		{
			yyVAL.expr = CNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 111:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:355
		{
			yyVAL.expr = CNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 112:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:356
		{
			yyVAL.expr = CNode("array", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 113:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:357
		{
			yyVAL.expr = CNode("array", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	}
	goto yystack /* stack new state and value */
}
