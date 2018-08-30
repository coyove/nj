//line .\parser.go.y:1
package parser

import __yyfmt__ "fmt"

//line .\parser.go.y:3
import (
	"fmt"
	"github.com/coyove/common/rand"
	"path/filepath"
)

//line .\parser.go.y:43
type yySymType struct {
	yys   int
	token Token
	expr  *Node
	str   string
}

const TAddressof = 57346
const TAssert = 57347
const TBreak = 57348
const TContinue = 57349
const TElse = 57350
const TFor = 57351
const TFunc = 57352
const TIf = 57353
const TLen = 57354
const TNew = 57355
const TNil = 57356
const TNot = 57357
const TReturn = 57358
const TRequire = 57359
const TTypeof = 57360
const TVar = 57361
const TWhile = 57362
const TYield = 57363
const TAddAdd = 57364
const TSubSub = 57365
const TEqeq = 57366
const TNeq = 57367
const TLsh = 57368
const TRsh = 57369
const TURsh = 57370
const TLte = 57371
const TGte = 57372
const TIdent = 57373
const TNumber = 57374
const TString = 57375
const FUN = 57376
const TOr = 57377
const TAnd = 57378
const UNARY = 57379
const TMinMin = 57380
const COPY = 57381

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"TAddressof",
	"TAssert",
	"TBreak",
	"TContinue",
	"TElse",
	"TFor",
	"TFunc",
	"TIf",
	"TLen",
	"TNew",
	"TNil",
	"TNot",
	"TReturn",
	"TRequire",
	"TTypeof",
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
	"'['",
	"'('",
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
	"COPY",
	"'}'",
	"'='",
	"']'",
	"'.'",
	"','",
	"':'",
	"')'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line .\parser.go.y:366

var _rand = rand.New()

func randomName() string {
	return fmt.Sprintf("%x", _rand.Fetch(16))
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 983

var yyAct = [...]int{

	185, 106, 70, 41, 110, 23, 103, 158, 30, 183,
	135, 52, 53, 119, 178, 6, 179, 159, 165, 56,
	158, 58, 120, 164, 124, 121, 191, 163, 69, 106,
	62, 65, 161, 189, 105, 93, 94, 42, 96, 24,
	102, 167, 156, 125, 113, 98, 99, 100, 101, 66,
	111, 50, 105, 25, 63, 104, 83, 84, 85, 23,
	23, 17, 23, 126, 92, 54, 131, 132, 181, 6,
	130, 172, 169, 128, 137, 138, 139, 140, 141, 142,
	143, 144, 145, 146, 147, 148, 149, 150, 151, 152,
	153, 154, 155, 24, 24, 48, 24, 21, 97, 116,
	123, 3, 118, 122, 14, 160, 95, 25, 25, 61,
	25, 59, 162, 57, 28, 170, 157, 1, 10, 8,
	9, 168, 19, 26, 20, 108, 173, 23, 176, 11,
	12, 177, 22, 18, 7, 55, 40, 117, 115, 13,
	5, 133, 68, 39, 29, 71, 72, 17, 64, 27,
	134, 21, 67, 4, 15, 3, 16, 180, 14, 43,
	182, 24, 87, 88, 89, 184, 60, 186, 171, 112,
	109, 23, 2, 0, 23, 25, 193, 192, 0, 0,
	0, 0, 81, 82, 83, 84, 85, 0, 0, 199,
	200, 23, 201, 13, 5, 0, 0, 203, 204, 0,
	0, 0, 0, 23, 23, 24, 0, 209, 24, 210,
	0, 0, 187, 23, 23, 190, 0, 0, 0, 25,
	0, 0, 25, 0, 0, 24, 81, 82, 83, 84,
	85, 0, 198, 0, 0, 0, 0, 24, 24, 25,
	0, 0, 0, 0, 205, 207, 0, 24, 24, 0,
	0, 25, 25, 0, 211, 212, 0, 0, 0, 0,
	0, 25, 25, 10, 8, 9, 0, 19, 0, 20,
	0, 0, 0, 0, 11, 12, 0, 22, 18, 7,
	0, 0, 79, 80, 87, 88, 89, 78, 77, 29,
	0, 0, 17, 0, 27, 0, 0, 73, 74, 90,
	91, 86, 75, 76, 81, 82, 83, 84, 85, 10,
	8, 9, 0, 19, 0, 20, 0, 0, 206, 0,
	11, 12, 0, 22, 18, 7, 0, 0, 79, 80,
	87, 88, 89, 78, 77, 29, 0, 0, 17, 0,
	27, 0, 0, 73, 74, 90, 91, 86, 75, 76,
	81, 82, 83, 84, 85, 79, 80, 87, 88, 89,
	78, 77, 0, 0, 0, 0, 0, 0, 0, 0,
	73, 74, 90, 91, 86, 75, 76, 81, 82, 83,
	84, 85, 79, 80, 87, 88, 89, 78, 77, 174,
	0, 0, 175, 0, 0, 0, 0, 73, 74, 90,
	91, 86, 75, 76, 81, 82, 83, 84, 85, 0,
	10, 8, 114, 0, 19, 0, 20, 0, 0, 0,
	136, 11, 12, 0, 22, 18, 7, 0, 0, 79,
	80, 87, 88, 89, 78, 77, 29, 0, 0, 17,
	0, 27, 0, 0, 73, 74, 90, 91, 86, 75,
	76, 81, 82, 83, 84, 85, 79, 80, 87, 88,
	89, 78, 77, 0, 0, 0, 0, 0, 0, 0,
	0, 73, 74, 90, 91, 86, 75, 76, 81, 82,
	83, 84, 85, 79, 80, 87, 88, 89, 78, 77,
	0, 0, 0, 197, 0, 0, 0, 0, 73, 74,
	90, 91, 86, 75, 76, 81, 82, 83, 84, 85,
	79, 80, 87, 88, 89, 78, 77, 0, 0, 0,
	166, 0, 0, 0, 0, 73, 74, 90, 91, 86,
	75, 76, 81, 82, 83, 84, 85, 79, 80, 87,
	88, 89, 78, 77, 0, 0, 208, 0, 0, 0,
	0, 0, 73, 74, 90, 91, 86, 75, 76, 81,
	82, 83, 84, 85, 79, 80, 87, 88, 89, 78,
	77, 0, 0, 196, 0, 0, 0, 0, 0, 73,
	74, 90, 91, 86, 75, 76, 81, 82, 83, 84,
	85, 0, 0, 36, 0, 0, 0, 0, 0, 26,
	188, 35, 38, 31, 46, 0, 33, 34, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 36, 0,
	29, 32, 51, 49, 26, 27, 35, 38, 31, 46,
	0, 33, 34, 0, 0, 37, 44, 0, 0, 0,
	0, 45, 47, 0, 0, 29, 32, 51, 49, 0,
	27, 129, 0, 0, 0, 0, 0, 0, 0, 0,
	37, 44, 0, 0, 0, 0, 45, 47, 79, 80,
	87, 88, 89, 78, 77, 127, 0, 0, 0, 0,
	0, 0, 0, 73, 74, 90, 91, 86, 75, 76,
	81, 82, 83, 84, 85, 79, 80, 87, 88, 89,
	78, 77, 202, 0, 0, 0, 0, 0, 0, 0,
	73, 74, 90, 91, 86, 75, 76, 81, 82, 83,
	84, 85, 36, 0, 0, 0, 0, 0, 26, 195,
	35, 38, 31, 46, 0, 33, 34, 79, 80, 87,
	88, 89, 78, 77, 0, 0, 0, 0, 0, 29,
	32, 51, 49, 0, 27, 0, 0, 75, 76, 81,
	82, 83, 84, 85, 37, 44, 0, 0, 0, 36,
	45, 47, 0, 0, 0, 26, 194, 35, 38, 31,
	46, 0, 33, 34, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 29, 32, 51, 49,
	0, 27, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 37, 44, 0, 36, 0, 0, 45, 47, 0,
	26, 107, 35, 38, 31, 46, 0, 33, 34, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 29, 32, 51, 49, 0, 27, 0, 0, 79,
	80, 87, 88, 89, 78, 77, 37, 44, 0, 0,
	0, 0, 45, 47, 73, 74, 90, 91, 86, 75,
	76, 81, 82, 83, 84, 85, 79, 80, 87, 88,
	89, 78, 77, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 74, 90, 91, 86, 75, 76, 81, 82,
	83, 84, 85, 79, 80, 87, 88, 89, 78, 77,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	90, 91, 86, 75, 76, 81, 82, 83, 84, 85,
	10, 8, 9, 0, 19, 26, 20, 0, 0, 0,
	0, 11, 12, 0, 22, 18, 7, 0, 0, 0,
	0, 10, 8, 9, 0, 19, 29, 20, 0, 17,
	0, 27, 11, 12, 0, 22, 18, 7, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 29, 0, 0,
	17, 0, 27,
}
var yyPact = [...]int{

	-1000, 925, -1000, -1000, -1000, -1000, -1000, 810, -1000, -1000,
	810, 810, 32, -1000, -1000, -1000, -1000, -1000, 810, 82,
	810, 80, 78, -5, -1000, -8, -1000, 810, -1000, 123,
	825, -1000, -1000, 31, 810, 810, 75, 810, 67, -1000,
	-1000, -5, -1000, -1000, 810, 810, 810, 810, -2, 765,
	-1000, -1000, 825, 825, -1000, 113, 405, -35, 304, 16,
	-36, -14, 614, 42, -1000, 589, 810, 131, -52, 358,
	-1000, -1000, -1000, 810, 810, 810, 810, 810, 810, 810,
	810, 810, 810, 810, 810, 810, 810, 810, 810, 810,
	810, 810, -1000, -1000, -1000, -1000, -1000, -15, -1000, -1000,
	-1000, -1000, 27, -40, 810, -30, -1000, -1000, -29, -37,
	-42, 459, -1000, -1000, -16, -1000, -1000, -1000, -1000, -1000,
	810, 41, 107, 946, 40, 810, 331, 810, 123, -1000,
	-46, 825, 825, -1000, -1000, -1000, -1000, 852, 879, 136,
	136, 136, 136, 136, 136, 8, 8, -1000, -1000, -1000,
	713, 180, 180, 180, 713, 713, 810, -1000, 37, 810,
	825, -1000, -53, -1000, 810, 810, 810, 946, 540, -24,
	946, -1000, -31, 825, 123, 718, 671, -1000, 810, -1000,
	513, -1000, 825, -1000, 432, 825, 825, 946, 810, 810,
	-1000, 810, -1000, 644, -1000, -1000, 810, 810, -1000, 258,
	304, 825, -1000, 486, 825, -1000, 810, -1000, 810, 304,
	304, -1000, -1000,
}
var yyPgo = [...]int{

	0, 117, 13, 172, 51, 6, 4, 170, 166, 0,
	37, 2, 159, 3, 156, 138, 137, 102, 44, 99,
	154, 153, 95, 152, 114, 148, 143, 40, 136, 125,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 15, 3, 3, 3, 3, 18,
	18, 18, 18, 18, 21, 21, 21, 14, 14, 14,
	14, 11, 11, 10, 10, 10, 16, 16, 16, 16,
	16, 17, 17, 23, 23, 23, 22, 20, 19, 19,
	19, 19, 19, 19, 4, 4, 4, 4, 4, 4,
	5, 5, 6, 6, 7, 7, 8, 8, 8, 8,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 12, 13, 13, 13,
	13, 24, 25, 25, 26, 26, 26, 27, 27, 28,
	28, 29, 29, 29, 29,
}
var yyR2 = [...]int{

	0, 0, 2, 3, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 2, 1, 1,
	3, 1, 1, 2, 5, 4, 3, 6, 7, 9,
	7, 3, 5, 0, 2, 2, 2, 4, 2, 1,
	1, 2, 2, 2, 1, 4, 6, 5, 5, 3,
	1, 3, 1, 3, 3, 5, 1, 3, 5, 3,
	1, 1, 2, 2, 2, 2, 2, 9, 1, 1,
	1, 1, 1, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 2, 2, 2, 2, 1, 1, 3, 1,
	3, 2, 2, 3, 3, 4, 3, 2, 3, 2,
	3, 1, 2, 1, 2,
}
var yyChk = [...]int{

	-1000, -1, -3, -19, -21, -15, -2, 21, 6, 7,
	5, 16, 17, -16, -17, -20, -14, 34, 20, 9,
	11, -22, 19, -13, -10, -4, 10, 36, -24, 31,
	-9, 14, 32, 17, 18, 12, 4, 46, 13, -26,
	-28, -13, -10, -12, 47, 52, 15, 53, -22, 34,
	-4, 33, -9, -9, 33, -1, -9, 31, -9, 31,
	-8, 31, 35, 59, -25, 36, 57, -23, -24, -9,
	-11, 22, 23, 39, 40, 44, 45, 30, 29, 24,
	25, 46, 47, 48, 49, 50, 43, 26, 27, 28,
	41, 42, 33, -9, -9, 31, -9, 31, -9, -9,
	-9, -9, -27, -5, 57, 36, 31, 56, -29, -7,
	-6, -9, 56, -18, 7, -15, -19, -16, -17, -2,
	57, 60, -18, -27, 60, 57, -9, 61, 31, 62,
	-6, -9, -9, 10, 19, 62, 62, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, -9, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, 57, -2, 60, 57,
	-9, 62, -5, 56, 60, 60, 61, 57, -9, 31,
	8, -18, 31, -9, 58, 61, -9, -11, 60, 62,
	-9, 31, -9, 62, -9, -9, -9, -18, 60, 57,
	-18, 57, -11, -9, 58, 58, 60, 61, -18, -9,
	-9, -9, 58, -9, -9, -18, 60, -18, 60, -9,
	-9, -18, -18,
}
var yyDef = [...]int{

	1, -2, 2, 5, 6, 7, 8, 0, 39, 40,
	0, 0, 0, 14, 15, 16, 4, 1, 0, 0,
	0, 0, 0, 18, 19, 97, 33, 0, 99, 44,
	38, 60, 61, 0, 0, 0, 0, 0, 0, 68,
	69, 70, 71, 72, 0, 0, 0, 0, 0, 0,
	97, 96, 41, 42, 43, 0, 0, 0, 0, 0,
	17, 56, 0, 0, 101, 0, 0, 36, 99, 0,
	23, 21, 22, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 62, 63, 64, 65, 66, 0, 92, 93,
	94, 95, 0, 0, 0, 0, 50, 109, 0, 111,
	113, 52, 3, 26, 40, 9, 10, 11, 12, 13,
	0, 0, 31, 0, 0, 0, 0, 0, 49, 102,
	0, 52, 20, 34, 35, 98, 100, 73, 74, 75,
	76, 77, 78, 79, 80, 81, 82, 83, 84, 85,
	86, 87, 88, 89, 90, 91, 0, 104, 0, 0,
	106, 107, 0, 110, 112, 114, 0, 0, 0, 0,
	0, 37, 59, 57, 45, 0, 0, 25, 0, 103,
	0, 51, 105, 108, 0, 53, 54, 0, 0, 0,
	32, 0, 24, 0, 47, 48, 0, 0, 27, 0,
	0, 58, 46, 0, 55, 28, 0, 30, 0, 0,
	0, 29, 67,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 53, 3, 50, 42, 3,
	36, 62, 48, 46, 60, 47, 59, 49, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 61, 3,
	45, 57, 44, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 37, 3, 3, 3, 3, 3,
	3, 35, 3, 58, 43, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 34, 41, 56, 52,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 38, 39, 40, 51, 54, 55,
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
		//line .\parser.go.y:77
		{
			yyVAL.expr = CNode("chain")
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:83
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[2].expr)
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 3:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:91
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:94
		{
			if yyDollar[1].expr.isIsolatedCopy() {
				yyDollar[1].expr.Cx(2).C()[0] = NNode(0.0)
			}
			yyVAL.expr = yyDollar[1].expr
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:102
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:103
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:104
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:105
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:108
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:109
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:110
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:111
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:112
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:115
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:116
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:117
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 17:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:120
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:121
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:122
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:123
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
		//line .\parser.go.y:143
		{
			yyVAL.expr = NNode(1.0)
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:144
		{
			yyVAL.expr = NNode(-1.0)
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:147
		{
			yyVAL.expr = CNode("inc", ANode(yyDollar[1].token).setPos(yyDollar[1].token), yyDollar[2].expr)
		}
	case 24:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:148
		{
			yyVAL.expr = CNode("store", yyDollar[1].expr, yyDollar[3].expr, CNode("+", CNode("load", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr), yyDollar[5].expr).setPos0(yyDollar[1].expr))
		}
	case 25:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:149
		{
			yyVAL.expr = CNode("store", yyDollar[1].expr, yyDollar[3].token, CNode("+", CNode("load", yyDollar[1].expr, yyDollar[3].token).setPos0(yyDollar[1].expr), yyDollar[4].expr).setPos0(yyDollar[1].expr))
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:152
		{
			yyVAL.expr = CNode("for", yyDollar[2].expr, CNode(), yyDollar[3].expr).setPos0(yyDollar[1].token)
		}
	case 27:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line .\parser.go.y:155
		{
			yyVAL.expr = CNode("for", yyDollar[2].expr, yyDollar[5].expr, yyDollar[6].expr).setPos0(yyDollar[1].token)
		}
	case 28:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line .\parser.go.y:158
		{
			vname, ename := ANode(yyDollar[2].token), ANodeS(yyDollar[2].token.Str+randomName())
			yyVAL.expr = CNode("chain",
				CNode("set", vname, yyDollar[4].expr),
				CNode("set", ename, yyDollar[6].expr),
				CNode("for",
					CNode("<", vname, ename).setPos0(yyDollar[1].token),
					CNode("chain",
						CNode("inc", vname, NNode(1.0)).setPos0(yyDollar[1].token),
					),
					yyDollar[7].expr,
				).setPos0(yyDollar[1].token),
			)
		}
	case 29:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line .\parser.go.y:172
		{
			vname, sname, ename := ANode(yyDollar[2].token), ANodeS(yyDollar[2].token.Str+randomName()), ANodeS(yyDollar[2].token.Str+randomName())
			if yyDollar[6].expr.Type == Nnumber {
				// easy case
				chain, cmp := CNode("chain", CNode("inc", vname, yyDollar[6].expr).setPos0(yyDollar[1].token)), "<="
				if yyDollar[6].expr.N() < 0 {
					cmp = ">="
				}
				yyVAL.expr = CNode("chain",
					CNode("set", vname, yyDollar[4].expr),
					CNode("set", ename, yyDollar[8].expr),
					CNode("for", CNode(cmp, vname, ename), chain, yyDollar[9].expr).setPos0(yyDollar[1].token),
				)
			} else {
				bname := ANodeS(yyDollar[2].token.Str + randomName())
				yyVAL.expr = CNode("chain",
					CNode("set", vname, yyDollar[4].expr),
					CNode("set", bname, yyDollar[4].expr),
					CNode("set", sname, yyDollar[6].expr),
					CNode("set", ename, yyDollar[8].expr),
					CNode("if", CNode("<=", NNode(0.0), CNode("*", CNode("-", ename, vname).setPos0(yyDollar[1].token), sname).setPos0(yyDollar[1].token)),
						CNode("chain",
							CNode("for",
								CNode("<=",
									CNode("*",
										CNode("-", vname, bname).setPos0(yyDollar[1].token),
										CNode("-", vname, ename).setPos0(yyDollar[1].token),
									),
									NNode(0.0),
								),
								CNode("chain", CNode("move", vname, CNode("+", vname, sname).setPos0(yyDollar[1].token))),
								yyDollar[9].expr,
							),
						),
						CNode("chain"),
					),
				)
			}

		}
	case 30:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line .\parser.go.y:212
		{
			yyVAL.expr = CNode("call", "copy", CNode(
				NNode(0),
				yyDollar[6].expr,
				CNode("func", "<anony-map-iter-callback>", CNode(yyDollar[2].token.Str, yyDollar[4].token.Str), yyDollar[7].expr),
			))
		}
	case 31:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:221
		{
			yyVAL.expr = CNode("if", yyDollar[2].expr, yyDollar[3].expr, CNode())
		}
	case 32:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:222
		{
			yyVAL.expr = CNode("if", yyDollar[2].expr, yyDollar[3].expr, yyDollar[5].expr)
		}
	case 33:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line .\parser.go.y:225
		{
			yyVAL.str = ""
		}
	case 34:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:226
		{
			yyVAL.str = yyDollar[1].str + ",safe"
		}
	case 35:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:227
		{
			yyVAL.str = yyDollar[1].str + ",var"
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:230
		{
			yyVAL.str = "func," + yyDollar[2].str
		}
	case 37:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:233
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
	case 38:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:245
		{
			yyVAL.expr = CNode("yield", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:246
		{
			yyVAL.expr = CNode("break").setPos0(yyDollar[1].token)
		}
	case 40:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:247
		{
			yyVAL.expr = CNode("continue").setPos0(yyDollar[1].token)
		}
	case 41:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:248
		{
			yyVAL.expr = CNode("assert", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 42:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:249
		{
			yyVAL.expr = CNode("ret", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 43:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:250
		{
			yyVAL.expr = yylex.(*Lexer).loadFile(filepath.Join(filepath.Dir(yyDollar[1].token.Pos.Source), yyDollar[2].token.Str))
		}
	case 44:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:253
		{
			yyVAL.expr = ANode(yyDollar[1].token).setPos(yyDollar[1].token)
		}
	case 45:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:254
		{
			yyVAL.expr = CNode("load", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 46:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line .\parser.go.y:255
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, yyDollar[3].expr, yyDollar[5].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 47:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:256
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, yyDollar[3].expr, NNode("-1")).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 48:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:257
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, NNode("0"), yyDollar[4].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 49:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:258
		{
			yyVAL.expr = CNode("load", yyDollar[1].expr, SNode(yyDollar[3].token.Str)).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 50:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:261
		{
			yyVAL.expr = CNode(yyDollar[1].token.Str)
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:262
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(ANode(yyDollar[3].token))
		}
	case 52:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:265
		{
			yyVAL.expr = CNode(yyDollar[1].expr)
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:266
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[3].expr)
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:269
		{
			yyVAL.expr = CNode(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 55:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:270
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[3].expr).Cappend(yyDollar[5].expr)
		}
	case 56:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:273
		{
			yyVAL.expr = CNode("chain", CNode("set", ANode(yyDollar[1].token), NilNode()).setPos0(yyDollar[1].token))
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:274
		{
			yyVAL.expr = CNode("chain", CNode("set", ANode(yyDollar[1].token), yyDollar[3].expr).setPos0(yyDollar[1].token))
		}
	case 58:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:275
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(CNode("set", ANode(yyDollar[3].token), yyDollar[5].expr).setPos0(yyDollar[1].expr))
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:276
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(CNode("set", ANode(yyDollar[3].token), NilNode()).setPos0(yyDollar[1].expr))
		}
	case 60:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:279
		{
			yyVAL.expr = NilNode().SetPos(yyDollar[1].token)
		}
	case 61:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:280
		{
			yyVAL.expr = NNode(yyDollar[1].token.Str).SetPos(yyDollar[1].token)
		}
	case 62:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:281
		{
			yyVAL.expr = yylex.(*Lexer).loadFile(filepath.Join(filepath.Dir(yyDollar[1].token.Pos.Source), yyDollar[2].token.Str))
		}
	case 63:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:282
		{
			yyVAL.expr = CNode("typeof", yyDollar[2].expr)
		}
	case 64:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:283
		{
			yyVAL.expr = CNode("len", yyDollar[2].expr)
		}
	case 65:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:284
		{
			yyVAL.expr = CNode("call", "addressof", CNode(ANode(yyDollar[2].token)))
		}
	case 66:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:285
		{
			yyVAL.expr = CNode("call", "copy", CNode(NNode(1), yyDollar[2].expr, NilNode()))
		}
	case 67:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line .\parser.go.y:286
		{
		}
	case 68:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:287
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 69:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:288
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 70:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:289
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 71:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:290
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:291
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:292
		{
			yyVAL.expr = CNode("or", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:293
		{
			yyVAL.expr = CNode("and", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 75:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:294
		{
			yyVAL.expr = CNode("<", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:295
		{
			yyVAL.expr = CNode("<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 77:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:296
		{
			yyVAL.expr = CNode("<=", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 78:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:297
		{
			yyVAL.expr = CNode("<=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 79:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:298
		{
			yyVAL.expr = CNode("==", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 80:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:299
		{
			yyVAL.expr = CNode("!=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 81:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:300
		{
			yyVAL.expr = CNode("+", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 82:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:301
		{
			yyVAL.expr = CNode("-", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 83:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:302
		{
			yyVAL.expr = CNode("*", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 84:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:303
		{
			yyVAL.expr = CNode("/", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:304
		{
			yyVAL.expr = CNode("%", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 86:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:305
		{
			yyVAL.expr = CNode("^", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 87:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:306
		{
			yyVAL.expr = CNode("<<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 88:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:307
		{
			yyVAL.expr = CNode(">>", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 89:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:308
		{
			yyVAL.expr = CNode(">>>", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 90:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:309
		{
			yyVAL.expr = CNode("|", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 91:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:310
		{
			yyVAL.expr = CNode("&", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 92:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:311
		{
			yyVAL.expr = CNode("-", NNode(0.0), yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 93:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:312
		{
			yyVAL.expr = CNode("~", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 94:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:313
		{
			yyVAL.expr = CNode("!", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 95:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:314
		{
			yyVAL.expr = CNode("#", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 96:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:317
		{
			yyVAL.expr = SNode(yyDollar[1].token.Str).SetPos(yyDollar[1].token)
		}
	case 97:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:320
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 98:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:321
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 99:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:322
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 100:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:323
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 101:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:326
		{
			switch yyDollar[1].expr.S() {
			case "copy":
				switch yyDollar[2].expr.Cn() {
				case 0:
					yylex.(*Lexer).Error("copy takes at least 1 argument")
				case 1:
					yyVAL.expr = CNode("call", yyDollar[1].expr, CNode(NNode(1), yyDollar[2].expr.Cx(0), NilNode()))
				default:
					yyVAL.expr = CNode("call", yyDollar[1].expr, CNode(NNode(1), yyDollar[2].expr.Cx(0), yyDollar[2].expr.Cx(1)))
				}
			default:
				yyVAL.expr = CNode("call", yyDollar[1].expr, yyDollar[2].expr)
			}
			yyVAL.expr.Cx(0).SetPos(yyDollar[1].expr)
		}
	case 102:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:344
		{
			yyVAL.expr = CNode()
		}
	case 103:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:345
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 104:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:348
		{
			yyVAL.expr = CNode(yyDollar[1].str, "<a>", yyDollar[2].expr, yyDollar[3].expr).setPos0(yyDollar[2].expr)
		}
	case 105:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:349
		{
			yyVAL.expr = CNode(yyDollar[1].str, "<a>", yyDollar[2].expr, CNode("chain", CNode("ret", yyDollar[4].expr))).setPos0(yyDollar[2].expr)
		}
	case 106:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:350
		{
			yyVAL.expr = CNode(yyDollar[1].str, "<a>", CNode(), CNode("chain", CNode("ret", yyDollar[3].expr))).setPos0(yyDollar[3].expr)
		}
	case 107:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:353
		{
			yyVAL.expr = CNode()
		}
	case 108:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:354
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 109:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:357
		{
			yyVAL.expr = CNode("map", CNode()).setPos0(yyDollar[1].token)
		}
	case 110:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:358
		{
			yyVAL.expr = yyDollar[2].expr.setPos0(yyDollar[1].token)
		}
	case 111:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:361
		{
			yyVAL.expr = CNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 112:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:362
		{
			yyVAL.expr = CNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 113:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:363
		{
			yyVAL.expr = CNode("array", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 114:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:364
		{
			yyVAL.expr = CNode("array", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	}
	goto yystack /* stack new state and value */
}
