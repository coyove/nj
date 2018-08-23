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
const TNot = 57354
const TReturn = 57355
const TRequire = 57356
const TNop = 57357
const TVar = 57358
const TWhile = 57359
const TYield = 57360
const TAddAdd = 57361
const TSubSub = 57362
const TEqeq = 57363
const TNeq = 57364
const TLsh = 57365
const TRsh = 57366
const TURsh = 57367
const TLte = 57368
const TGte = 57369
const TIdent = 57370
const TNumber = 57371
const TString = 57372
const FUN = 57373
const TOr = 57374
const TAnd = 57375
const UNARY = 57376
const TMinMin = 57377

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
	"TNot",
	"TReturn",
	"TRequire",
	"TNop",
	"TVar",
	"TWhile",
	"TYield",
	"TAddAdd",
	"TSubSub",
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
	"'_'",
	"'T'",
	"FUN",
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
	"'='",
	"'['",
	"']'",
	"'.'",
	"';'",
	"','",
	"')'",
	"'!'",
	"':'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.go.y:354

var typesLookup = map[string]string{
	"nil": "0", "number": "1", "string": "2", "map": "4", "closure": "6", "generic": "7",
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 857

var yyAct = [...]int{

	179, 98, 70, 37, 16, 23, 101, 105, 31, 152,
	177, 48, 50, 108, 67, 38, 130, 24, 113, 53,
	6, 173, 174, 23, 22, 71, 72, 153, 69, 159,
	158, 121, 152, 97, 118, 24, 55, 155, 185, 183,
	27, 93, 94, 95, 96, 101, 106, 122, 44, 100,
	21, 110, 46, 3, 25, 115, 23, 23, 112, 119,
	14, 100, 66, 123, 117, 57, 128, 129, 24, 24,
	99, 6, 25, 127, 132, 133, 134, 135, 136, 137,
	138, 139, 140, 141, 142, 143, 144, 145, 146, 147,
	148, 149, 150, 120, 109, 111, 5, 13, 65, 157,
	154, 21, 156, 58, 3, 25, 25, 83, 84, 85,
	17, 14, 92, 51, 56, 175, 151, 167, 23, 161,
	62, 164, 63, 168, 23, 171, 125, 61, 172, 59,
	24, 1, 28, 162, 166, 192, 24, 81, 82, 83,
	84, 85, 71, 72, 103, 36, 35, 5, 13, 52,
	64, 4, 15, 114, 176, 54, 39, 60, 104, 178,
	68, 180, 2, 0, 0, 23, 23, 25, 0, 23,
	0, 187, 186, 25, 0, 181, 182, 24, 24, 184,
	0, 24, 22, 0, 191, 0, 193, 0, 0, 0,
	0, 195, 0, 0, 29, 23, 23, 17, 27, 87,
	88, 89, 0, 0, 0, 196, 197, 24, 24, 0,
	0, 0, 0, 0, 25, 25, 0, 0, 25, 81,
	82, 83, 84, 85, 0, 0, 0, 0, 0, 0,
	0, 79, 80, 87, 88, 89, 78, 77, 0, 0,
	0, 0, 0, 0, 25, 25, 73, 74, 90, 91,
	86, 75, 76, 81, 82, 83, 84, 85, 79, 80,
	87, 88, 89, 78, 77, 169, 0, 0, 0, 0,
	0, 170, 0, 73, 74, 90, 91, 86, 75, 76,
	81, 82, 83, 84, 85, 79, 80, 87, 88, 89,
	78, 77, 0, 0, 0, 0, 0, 0, 190, 0,
	73, 74, 90, 91, 86, 75, 76, 81, 82, 83,
	84, 85, 0, 0, 10, 8, 9, 0, 19, 0,
	20, 0, 0, 11, 12, 160, 22, 18, 7, 0,
	0, 79, 80, 87, 88, 89, 78, 77, 29, 0,
	0, 17, 27, 0, 0, 0, 73, 74, 90, 91,
	86, 75, 76, 81, 82, 83, 84, 85, 79, 80,
	87, 88, 89, 78, 77, 0, 0, 0, 0, 0,
	0, 0, 0, 73, 74, 90, 91, 86, 75, 76,
	81, 82, 83, 84, 85, 79, 80, 87, 88, 89,
	78, 77, 0, 0, 0, 0, 165, 0, 0, 0,
	73, 74, 90, 91, 86, 75, 76, 81, 82, 83,
	84, 85, 79, 80, 87, 88, 89, 78, 77, 0,
	0, 0, 0, 131, 0, 0, 0, 73, 74, 90,
	91, 86, 75, 76, 81, 82, 83, 84, 85, 79,
	80, 87, 88, 89, 78, 77, 0, 0, 163, 0,
	0, 0, 0, 0, 73, 74, 90, 91, 86, 75,
	76, 81, 82, 83, 84, 85, 79, 80, 87, 88,
	89, 78, 77, 194, 26, 0, 32, 42, 0, 34,
	0, 73, 74, 90, 91, 86, 75, 76, 81, 82,
	83, 84, 85, 29, 33, 47, 45, 27, 0, 26,
	189, 32, 42, 0, 34, 0, 0, 0, 0, 40,
	0, 0, 0, 0, 41, 43, 0, 0, 29, 33,
	47, 45, 27, 0, 0, 26, 124, 32, 42, 0,
	34, 0, 0, 0, 40, 0, 0, 0, 0, 41,
	43, 0, 0, 0, 29, 33, 47, 45, 27, 126,
	0, 0, 0, 0, 0, 0, 26, 0, 32, 42,
	40, 34, 0, 0, 0, 41, 43, 26, 0, 32,
	42, 0, 34, 116, 0, 29, 33, 47, 45, 27,
	0, 0, 0, 0, 0, 0, 29, 33, 47, 45,
	27, 40, 0, 0, 0, 0, 41, 43, 0, 0,
	0, 0, 40, 49, 0, 0, 0, 41, 43, 0,
	26, 0, 32, 42, 30, 34, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 29,
	33, 47, 45, 27, 10, 8, 9, 0, 19, 26,
	20, 0, 0, 11, 12, 40, 22, 18, 7, 0,
	41, 43, 0, 0, 0, 0, 188, 0, 29, 0,
	0, 17, 27, 0, 0, 0, 0, 0, 79, 80,
	87, 88, 89, 78, 77, 0, 0, 0, 0, 0,
	0, 0, 107, 73, 74, 90, 91, 86, 75, 76,
	81, 82, 83, 84, 85, 79, 80, 87, 88, 89,
	78, 77, 0, 26, 0, 32, 42, 0, 34, 0,
	0, 74, 90, 91, 86, 75, 76, 81, 82, 83,
	84, 85, 29, 33, 47, 45, 27, 0, 0, 0,
	79, 80, 87, 88, 89, 78, 77, 0, 40, 0,
	0, 0, 0, 41, 43, 0, 102, 90, 91, 86,
	75, 76, 81, 82, 83, 84, 85, 26, 0, 32,
	42, 0, 34, 79, 80, 87, 88, 89, 78, 77,
	0, 0, 0, 0, 0, 0, 29, 33, 47, 45,
	27, 0, 0, 75, 76, 81, 82, 83, 84, 85,
	0, 0, 40, 0, 0, 0, 0, 41, 43, 10,
	8, 9, 0, 19, 26, 20, 0, 0, 11, 12,
	0, 22, 18, 7, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 29, 0, 0, 17, 27, 10, 8,
	9, 0, 19, 0, 20, 0, 0, 11, 12, 0,
	22, 18, 7, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 29, 0, 0, 17, 27,
}
var yyPact = [...]int{

	-1000, 795, -1000, -1000, -1000, -1000, -1000, 558, -1000, -1000,
	748, 547, 83, -1000, -1000, -1000, -1000, -1000, 748, 8,
	71, 101, 99, 66, -1000, 9, -46, 748, -1000, 123,
	-1000, 647, -1000, -1000, 82, -1000, -1000, 66, -1000, -1000,
	748, 748, 748, 748, 17, 694, -1000, -1000, 647, -1000,
	647, -1000, 630, 310, 516, 6, -23, -1000, 748, 29,
	-27, -6, 465, 98, -1000, 490, 748, -1000, -43, 364,
	-1000, -1000, -1000, 748, 748, 748, 748, 748, 748, 748,
	748, 748, 748, 748, 748, 748, 748, 748, 748, 748,
	748, 748, -1000, -1000, -1000, -1000, -1000, 79, -26, 748,
	-22, -1000, -1000, 47, -28, -29, 264, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, 166, 391, -1000, 93, -1000, 337,
	824, 89, 748, 210, 748, 123, -1000, -37, 647, 647,
	-1000, -1000, 674, 709, 176, 176, 176, 176, 176, 176,
	62, 62, -1000, -1000, -1000, 742, 94, 94, 94, 742,
	742, -1000, 87, 748, 647, -1000, -49, -1000, 748, 748,
	748, 824, 824, -1000, -14, 824, -1000, -15, 647, 123,
	601, 445, -1000, 748, -1000, -1000, 647, -1000, 237, 647,
	647, -1000, -1000, 748, 128, 748, -1000, 418, -1000, -1000,
	748, 310, 824, 647, -1000, 647, -1000, -1000,
}
var yyPgo = [...]int{

	0, 131, 18, 162, 52, 1, 7, 158, 157, 0,
	15, 2, 156, 3, 4, 94, 95, 155, 153, 58,
	13, 51, 152, 151, 48, 132, 150, 146, 33, 145,
	144,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 15, 3, 3, 3, 3, 20,
	20, 20, 20, 20, 23, 23, 23, 14, 14, 14,
	14, 11, 11, 10, 10, 10, 17, 17, 18, 18,
	16, 16, 16, 16, 19, 19, 24, 24, 22, 21,
	21, 21, 21, 21, 21, 21, 21, 4, 4, 4,
	4, 4, 4, 5, 5, 6, 6, 7, 7, 8,
	8, 8, 8, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 12, 13, 13, 13, 13, 25,
	26, 26, 27, 27, 27, 28, 28, 29, 29, 30,
	30, 30, 30,
}
var yyR2 = [...]int{

	0, 0, 2, 3, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 2, 1, 1,
	3, 1, 1, 2, 5, 4, 2, 1, 2, 1,
	3, 5, 5, 7, 5, 7, 1, 2, 4, 2,
	2, 1, 1, 2, 2, 2, 2, 1, 4, 6,
	5, 5, 3, 1, 3, 1, 3, 3, 5, 1,
	3, 5, 3, 1, 1, 2, 1, 1, 1, 1,
	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	2, 2, 2, 2, 1, 1, 3, 1, 3, 2,
	2, 3, 3, 4, 3, 2, 3, 2, 3, 1,
	2, 1, 2,
}
var yyChk = [...]int{

	-1000, -1, -3, -21, -23, -15, -2, 18, 5, 6,
	4, 13, 14, -16, -19, -22, -14, 31, 17, 8,
	10, -24, 16, -13, -10, -4, 9, 32, -25, 28,
	56, -9, 11, 29, 14, -27, -29, -13, -10, -12,
	44, 49, 12, 50, -24, 31, -4, 30, -9, 56,
	-9, 30, -1, -9, -17, 28, -15, 57, 32, 28,
	-8, 28, 54, 56, -26, 32, 53, 60, -25, -9,
	-11, 19, 20, 36, 37, 41, 42, 27, 26, 21,
	22, 43, 44, 45, 46, 47, 40, 23, 24, 25,
	38, 39, 30, -9, -9, -9, -9, -28, -5, 53,
	32, 28, 52, -30, -7, -6, -9, 52, -20, -15,
	-21, -16, -19, -2, -18, -9, 57, 58, 57, -9,
	-28, 58, 53, -9, 61, 28, 59, -6, -9, -9,
	59, 59, -9, -9, -9, -9, -9, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, -9, -9, -9, -9,
	-9, -2, 58, 53, -9, 59, -5, 52, 58, 58,
	61, -14, -2, 57, 28, 59, -20, 28, -9, 55,
	61, -9, -11, 58, 59, 28, -9, 59, -9, -9,
	-9, -20, -20, 53, -20, 53, -11, -9, 55, 55,
	61, -9, 7, -9, 55, -9, -20, -20,
}
var yyDef = [...]int{

	1, -2, 2, 5, 6, 7, 8, 0, 41, 42,
	0, 0, 0, 14, 15, 16, 4, 1, 0, 0,
	0, 0, 0, 18, 19, 95, 36, 0, 97, 47,
	39, 40, 63, 64, 0, 66, 67, 68, 69, 70,
	0, 0, 0, 0, 0, 0, 95, 94, 43, 44,
	45, 46, 0, 0, 0, 47, 0, 27, 0, 0,
	17, 59, 0, 0, 99, 0, 0, 37, 97, 0,
	23, 21, 22, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 65, 90, 91, 92, 93, 0, 0, 0,
	0, 53, 107, 0, 109, 111, 55, 3, 30, 9,
	10, 11, 12, 13, 0, 0, 29, 0, 26, 0,
	0, 0, 0, 0, 0, 52, 100, 0, 55, 20,
	96, 98, 71, 72, 73, 74, 75, 76, 77, 78,
	79, 80, 81, 82, 83, 84, 85, 86, 87, 88,
	89, 102, 0, 0, 104, 105, 0, 108, 110, 112,
	0, 0, 0, 28, 0, 0, 38, 62, 60, 48,
	0, 0, 25, 0, 101, 54, 103, 106, 0, 56,
	57, 31, 32, 0, 34, 0, 24, 0, 50, 51,
	0, 0, 0, 61, 49, 58, 33, 35,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 60, 3, 50, 3, 47, 39, 3,
	32, 59, 45, 43, 58, 44, 56, 46, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 61, 57,
	42, 53, 41, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 34, 3, 3, 3, 3, 3,
	3, 54, 3, 55, 40, 33, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 31, 38, 52, 49,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 35,
	36, 37, 48, 51,
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
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:91
		{
			if yyDollar[1].expr.isIsolatedCopy() {
				yyDollar[1].expr.Cx(2).C()[0] = NNode(0.0)
			}
			yyVAL.expr = yyDollar[1].expr
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:99
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
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
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:106
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
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
			yyVAL.expr = yyDollar[1].expr
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:112
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
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:117
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:118
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:119
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:120
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
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:140
		{
			yyVAL.expr = NNode(1.0)
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:141
		{
			yyVAL.expr = NNode(-1.0)
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:144
		{
			yyVAL.expr = CNode("inc", ANode(yyDollar[1].token).setPos(yyDollar[1].token), yyDollar[2].expr)
		}
	case 24:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:145
		{
			yyVAL.expr = CNode("store", yyDollar[1].expr, yyDollar[3].expr, CNode("+", CNode("load", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr), yyDollar[5].expr).setPos0(yyDollar[1].expr))
		}
	case 25:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:146
		{
			yyVAL.expr = CNode("store", yyDollar[1].expr, yyDollar[3].token, CNode("+", CNode("load", yyDollar[1].expr, yyDollar[3].token).setPos0(yyDollar[1].expr), yyDollar[4].expr).setPos0(yyDollar[1].expr))
		}
	case 26:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:149
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:150
		{
			yyVAL.expr = CNode("chain")
		}
	case 28:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:153
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:154
		{
			yyVAL.expr = NNode(1.0)
		}
	case 30:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:157
		{
			yyVAL.expr = CNode("for", yyDollar[2].expr, CNode(), yyDollar[3].expr).setPos0(yyDollar[1].token)
		}
	case 31:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:160
		{
			yyVAL.expr = yyDollar[2].expr.Cappend(CNode("for", yyDollar[3].expr, CNode("chain", yyDollar[4].expr), yyDollar[5].expr)).setPos0(yyDollar[1].token)
		}
	case 32:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:163
		{
			yyVAL.expr = yyDollar[2].expr.Cappend(CNode("for", yyDollar[3].expr, yyDollar[4].expr, yyDollar[5].expr)).setPos0(yyDollar[1].token)
		}
	case 33:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:166
		{
			yyVAL.expr = CNode("call", "copy", CNode(
				NNode(0.0),
				yyDollar[6].expr,
				CNode("func", "<anony-map-iter-callback>", CNode(yyDollar[2].token.Str, yyDollar[4].token.Str), yyDollar[7].expr),
			))
		}
	case 34:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:175
		{
			yyVAL.expr = CNode("if", yyDollar[3].expr, yyDollar[5].expr, CNode())
		}
	case 35:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:176
		{
			yyVAL.expr = CNode("if", yyDollar[3].expr, yyDollar[5].expr, yyDollar[7].expr)
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:179
		{
			yyVAL.str = "func"
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:180
		{
			yyVAL.str = "safefunc"
		}
	case 38:
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
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:195
		{
			yyVAL.expr = CNode("yield").setPos0(yyDollar[1].token)
		}
	case 40:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:196
		{
			yyVAL.expr = CNode("yield", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 41:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:197
		{
			yyVAL.expr = CNode("break").setPos0(yyDollar[1].token)
		}
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:198
		{
			yyVAL.expr = CNode("continue").setPos0(yyDollar[1].token)
		}
	case 43:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:199
		{
			yyVAL.expr = CNode("assert", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 44:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:200
		{
			yyVAL.expr = CNode("ret").setPos0(yyDollar[1].token)
		}
	case 45:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:201
		{
			if yyDollar[2].expr.isIsolatedCopy() && yyDollar[2].expr.Cx(2).Cx(2).N() == 1 {
				yyDollar[2].expr.Cx(2).C()[2] = NNode(2.0)
			}
			yyVAL.expr = CNode("ret", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 46:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:207
		{
			path := filepath.Join(filepath.Dir(yyDollar[1].token.Pos.Source), yyDollar[2].token.Str)
			yyVAL.expr = yylex.(*Lexer).loadFile(path)
		}
	case 47:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:213
		{
			yyVAL.expr = ANode(yyDollar[1].token).setPos(yyDollar[1].token)
		}
	case 48:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:214
		{
			yyVAL.expr = CNode("load", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 49:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:215
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, yyDollar[3].expr, yyDollar[5].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 50:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:216
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, yyDollar[3].expr, NNode("-1")).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 51:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:217
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, NNode("0"), yyDollar[4].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:218
		{
			yyVAL.expr = CNode("load", yyDollar[1].expr, SNode(yyDollar[3].token.Str)).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 53:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:221
		{
			yyVAL.expr = CNode(yyDollar[1].token.Str)
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:222
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(ANode(yyDollar[3].token))
		}
	case 55:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:225
		{
			yyVAL.expr = CNode(yyDollar[1].expr)
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:226
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[3].expr)
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:229
		{
			yyVAL.expr = CNode(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 58:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:230
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[3].expr).Cappend(yyDollar[5].expr)
		}
	case 59:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:233
		{
			yyVAL.expr = CNode("chain", CNode("set", ANode(yyDollar[1].token), NilNode()).setPos0(yyDollar[1].token))
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:234
		{
			yyVAL.expr = CNode("chain", CNode("set", ANode(yyDollar[1].token), yyDollar[3].expr).setPos0(yyDollar[1].token))
		}
	case 61:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:235
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(CNode("set", ANode(yyDollar[3].token), yyDollar[5].expr).setPos0(yyDollar[1].expr))
		}
	case 62:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:236
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(CNode("set", ANode(yyDollar[3].token), NilNode()).setPos0(yyDollar[1].expr))
		}
	case 63:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:239
		{
			yyVAL.expr = NilNode().SetPos(yyDollar[1].token)
		}
	case 64:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:240
		{
			yyVAL.expr = NNode(yyDollar[1].token.Str).SetPos(yyDollar[1].token)
		}
	case 65:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:241
		{
			yyVAL.expr = yylex.(*Lexer).loadFile(filepath.Join(filepath.Dir(yyDollar[1].token.Pos.Source), yyDollar[2].token.Str))
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:242
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 67:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:243
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 68:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:244
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 69:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:245
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 70:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:246
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:247
		{
			yyVAL.expr = CNode("or", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 72:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:248
		{
			yyVAL.expr = CNode("and", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:249
		{
			yyVAL.expr = CNode("<", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:250
		{
			yyVAL.expr = CNode("<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 75:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:251
		{
			yyVAL.expr = CNode("<=", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:252
		{
			yyVAL.expr = CNode("<=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 77:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:253
		{
			yyVAL.expr = CNode("==", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 78:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:254
		{
			yyVAL.expr = CNode("!=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 79:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:255
		{
			yyVAL.expr = CNode("+", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 80:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:256
		{
			yyVAL.expr = CNode("-", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 81:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:257
		{
			yyVAL.expr = CNode("*", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 82:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:258
		{
			yyVAL.expr = CNode("/", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 83:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:259
		{
			yyVAL.expr = CNode("%", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 84:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:260
		{
			yyVAL.expr = CNode("^", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:261
		{
			yyVAL.expr = CNode("<<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 86:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:262
		{
			yyVAL.expr = CNode(">>", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 87:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:263
		{
			yyVAL.expr = CNode(">>>", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 88:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:264
		{
			yyVAL.expr = CNode("|", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 89:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:265
		{
			yyVAL.expr = CNode("&", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 90:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:266
		{
			yyVAL.expr = CNode("-", NNode(0.0), yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 91:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:267
		{
			yyVAL.expr = CNode("~", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 92:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:268
		{
			yyVAL.expr = CNode("!", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 93:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:269
		{
			yyVAL.expr = CNode("#", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 94:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:272
		{
			yyVAL.expr = SNode(yyDollar[1].token.Str).SetPos(yyDollar[1].token)
		}
	case 95:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:275
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 96:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:276
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 97:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:277
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 98:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:278
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 99:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:281
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
						yylex.(*Lexer).Error("invalid argument for copy")
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
	case 100:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:332
		{
			yyVAL.expr = CNode()
		}
	case 101:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:333
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 102:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:336
		{
			yyVAL.expr = CNode(yyDollar[1].str, "<a>", yyDollar[2].expr, yyDollar[3].expr).setPos0(yyDollar[2].expr)
		}
	case 103:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:337
		{
			yyVAL.expr = CNode(yyDollar[1].str, "<a>", yyDollar[2].expr, CNode("chain", CNode("ret", yyDollar[4].expr))).setPos0(yyDollar[2].expr)
		}
	case 104:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:338
		{
			yyVAL.expr = CNode(yyDollar[1].str, "<a>", CNode(), CNode("chain", CNode("ret", yyDollar[3].expr))).setPos0(yyDollar[3].expr)
		}
	case 105:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:341
		{
			yyVAL.expr = CNode()
		}
	case 106:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:342
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 107:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:345
		{
			yyVAL.expr = CNode("map", CNode()).setPos0(yyDollar[1].token)
		}
	case 108:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:346
		{
			yyVAL.expr = yyDollar[2].expr.setPos0(yyDollar[1].token)
		}
	case 109:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:349
		{
			yyVAL.expr = CNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 110:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:350
		{
			yyVAL.expr = CNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 111:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:351
		{
			yyVAL.expr = CNode("array", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 112:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:352
		{
			yyVAL.expr = CNode("array", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	}
	goto yystack /* stack new state and value */
}
