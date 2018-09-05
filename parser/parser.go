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
const TUse = 57359
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
	"TUse",
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

//line .\parser.go.y:353

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

const yyLast = 929

var yyAct = [...]int{

	184, 101, 108, 40, 69, 23, 134, 157, 30, 182,
	164, 51, 52, 41, 177, 24, 178, 158, 104, 55,
	157, 57, 119, 163, 123, 120, 104, 190, 68, 61,
	64, 103, 188, 166, 100, 92, 93, 118, 95, 6,
	49, 124, 25, 162, 96, 97, 98, 99, 160, 109,
	103, 102, 62, 112, 65, 47, 17, 21, 23, 23,
	91, 23, 125, 53, 180, 130, 131, 129, 24, 24,
	171, 24, 168, 136, 137, 138, 139, 140, 141, 142,
	143, 144, 145, 146, 147, 148, 149, 150, 151, 152,
	153, 154, 6, 122, 127, 25, 25, 115, 25, 3,
	28, 94, 117, 159, 14, 161, 116, 114, 13, 5,
	21, 121, 80, 81, 82, 83, 84, 82, 83, 84,
	167, 60, 58, 56, 1, 172, 23, 175, 67, 10,
	8, 9, 176, 19, 26, 20, 24, 132, 156, 155,
	11, 12, 54, 22, 18, 7, 133, 70, 71, 169,
	106, 39, 3, 38, 63, 29, 179, 14, 17, 181,
	27, 13, 5, 25, 183, 66, 185, 4, 15, 16,
	23, 42, 59, 23, 107, 192, 170, 2, 191, 111,
	24, 0, 0, 24, 0, 0, 0, 0, 197, 198,
	23, 199, 0, 86, 87, 88, 201, 0, 0, 0,
	24, 23, 23, 0, 205, 0, 0, 25, 0, 23,
	25, 24, 24, 80, 81, 82, 83, 84, 0, 24,
	186, 0, 0, 189, 0, 0, 0, 25, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 25, 25,
	196, 0, 0, 0, 0, 0, 25, 0, 0, 0,
	0, 202, 204, 0, 10, 8, 9, 0, 19, 206,
	20, 0, 0, 0, 0, 11, 12, 0, 22, 18,
	7, 0, 0, 78, 79, 86, 87, 88, 77, 76,
	29, 0, 0, 17, 0, 27, 0, 0, 72, 73,
	89, 90, 85, 74, 75, 80, 81, 82, 83, 84,
	10, 8, 9, 0, 19, 0, 20, 0, 203, 0,
	0, 11, 12, 0, 22, 18, 7, 0, 0, 78,
	79, 86, 87, 88, 77, 76, 29, 0, 0, 17,
	0, 27, 0, 0, 72, 73, 89, 90, 85, 74,
	75, 80, 81, 82, 83, 84, 10, 8, 113, 0,
	19, 0, 20, 0, 0, 0, 0, 11, 12, 0,
	22, 18, 7, 0, 0, 78, 79, 86, 87, 88,
	77, 76, 29, 0, 0, 17, 0, 27, 0, 0,
	72, 73, 89, 90, 85, 74, 75, 80, 81, 82,
	83, 84, 78, 79, 86, 87, 88, 77, 76, 0,
	0, 0, 0, 0, 0, 0, 0, 72, 73, 89,
	90, 85, 74, 75, 80, 81, 82, 83, 84, 0,
	0, 0, 0, 0, 0, 173, 0, 0, 174, 78,
	79, 86, 87, 88, 77, 76, 0, 0, 0, 0,
	0, 0, 0, 0, 72, 73, 89, 90, 85, 74,
	75, 80, 81, 82, 83, 84, 78, 79, 86, 87,
	88, 77, 76, 0, 0, 0, 135, 0, 0, 0,
	0, 72, 73, 89, 90, 85, 74, 75, 80, 81,
	82, 83, 84, 78, 79, 86, 87, 88, 77, 76,
	0, 0, 195, 0, 0, 0, 0, 0, 72, 73,
	89, 90, 85, 74, 75, 80, 81, 82, 83, 84,
	78, 79, 86, 87, 88, 77, 76, 0, 0, 165,
	0, 0, 0, 0, 0, 72, 73, 89, 90, 85,
	74, 75, 80, 81, 82, 83, 84, 0, 36, 0,
	0, 0, 0, 0, 26, 187, 35, 37, 31, 45,
	0, 33, 34, 0, 78, 79, 86, 87, 88, 77,
	76, 0, 0, 0, 0, 29, 32, 50, 48, 0,
	27, 89, 90, 85, 74, 75, 80, 81, 82, 83,
	84, 43, 0, 0, 0, 0, 44, 46, 78, 79,
	86, 87, 88, 77, 76, 128, 0, 0, 0, 0,
	0, 0, 0, 72, 73, 89, 90, 85, 74, 75,
	80, 81, 82, 83, 84, 0, 0, 0, 0, 0,
	0, 200, 78, 79, 86, 87, 88, 77, 76, 0,
	0, 0, 0, 0, 0, 0, 0, 72, 73, 89,
	90, 85, 74, 75, 80, 81, 82, 83, 84, 0,
	36, 0, 0, 0, 0, 194, 26, 0, 35, 37,
	31, 45, 36, 33, 34, 0, 0, 0, 26, 0,
	35, 37, 31, 45, 0, 33, 34, 29, 32, 50,
	48, 0, 27, 0, 0, 0, 0, 0, 0, 29,
	32, 50, 48, 43, 27, 0, 0, 0, 44, 46,
	0, 0, 0, 0, 0, 43, 126, 0, 36, 0,
	44, 46, 0, 0, 26, 193, 35, 37, 31, 45,
	0, 33, 34, 78, 79, 86, 87, 88, 77, 76,
	0, 0, 0, 0, 0, 29, 32, 50, 48, 0,
	27, 0, 0, 74, 75, 80, 81, 82, 83, 84,
	0, 43, 0, 0, 0, 0, 44, 46, 0, 105,
	78, 79, 86, 87, 88, 77, 76, 0, 0, 110,
	0, 0, 0, 0, 0, 72, 73, 89, 90, 85,
	74, 75, 80, 81, 82, 83, 84, 36, 0, 0,
	0, 0, 0, 26, 0, 35, 37, 31, 45, 0,
	33, 34, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 29, 32, 50, 48, 0, 27,
	0, 0, 78, 79, 86, 87, 88, 77, 76, 0,
	43, 0, 0, 0, 0, 44, 46, 72, 73, 89,
	90, 85, 74, 75, 80, 81, 82, 83, 84, 78,
	79, 86, 87, 88, 77, 76, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 73, 89, 90, 85, 74,
	75, 80, 81, 82, 83, 84, 10, 8, 9, 0,
	19, 26, 20, 0, 0, 0, 0, 11, 12, 0,
	22, 18, 7, 0, 0, 0, 0, 10, 8, 9,
	0, 19, 29, 20, 0, 17, 0, 27, 11, 12,
	0, 22, 18, 7, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 29, 0, 0, 17, 0, 27,
}
var yyPact = [...]int{

	-1000, 871, -1000, -1000, -1000, -1000, -1000, 783, -1000, -1000,
	783, 783, 30, -1000, -1000, -1000, -1000, -1000, 783, 92,
	783, 91, 90, -6, -1000, -2, -1000, 783, -1000, 125,
	798, -1000, -1000, 27, 783, 783, 70, 783, -1000, -1000,
	-6, -1000, -1000, 783, 783, 783, 783, -5, 704, -1000,
	-1000, 736, 798, -1000, 124, 341, -34, 295, 14, -35,
	-15, 646, 63, -1000, 534, 783, 127, -55, 405, -1000,
	-1000, -1000, 783, 783, 783, 783, 783, 783, 783, 783,
	783, 783, 783, 783, 783, 783, 783, 783, 783, 783,
	783, -1000, -1000, -1000, -1000, 122, -1000, -1000, -1000, -1000,
	22, -39, 783, -13, -1000, -1000, -12, -36, -49, 459,
	-1000, -1000, -1000, -23, -1000, -1000, -1000, -1000, -1000, 783,
	41, 141, 892, 39, 783, 368, 783, 125, -1000, -45,
	798, 798, -1000, -1000, -1000, -1000, 825, 530, 167, 167,
	167, 167, 167, 167, 69, 69, -1000, -1000, -1000, 699,
	66, 66, 66, 699, 699, 783, -1000, 33, 783, 798,
	-1000, -52, -1000, 783, 783, 783, 892, 486, -24, 892,
	-1000, -29, 798, 125, 658, 598, -1000, 783, -1000, -1000,
	-1000, 798, -1000, 432, 798, 798, 892, 783, 783, -1000,
	783, -1000, 564, -1000, -1000, 783, -1000, 249, 295, 798,
	-1000, 798, -1000, 783, -1000, 295, -1000,
}
var yyPgo = [...]int{

	0, 124, 37, 177, 40, 1, 2, 174, 172, 0,
	13, 4, 171, 3, 169, 107, 106, 102, 53, 97,
	168, 167, 55, 165, 100, 154, 153, 34, 151, 150,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 15, 3, 3, 3, 3, 18,
	18, 18, 18, 18, 21, 21, 21, 14, 14, 14,
	14, 11, 11, 10, 10, 10, 16, 16, 16, 16,
	16, 17, 17, 23, 23, 23, 22, 20, 19, 19,
	19, 19, 19, 19, 19, 4, 4, 4, 4, 4,
	4, 5, 5, 6, 6, 7, 7, 8, 8, 8,
	8, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 12, 13, 13,
	13, 13, 24, 25, 25, 26, 26, 26, 27, 27,
	28, 28, 29, 29, 29, 29,
}
var yyR2 = [...]int{

	0, 0, 2, 3, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 2, 1, 1,
	3, 1, 1, 2, 5, 4, 3, 6, 7, 9,
	7, 3, 5, 0, 2, 2, 2, 4, 2, 1,
	1, 2, 3, 2, 2, 1, 4, 6, 5, 5,
	3, 1, 3, 1, 3, 3, 5, 1, 3, 5,
	3, 1, 1, 2, 2, 2, 2, 2, 4, 1,
	1, 1, 1, 1, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 2, 2, 2, 2, 1, 1, 3,
	1, 3, 2, 2, 3, 3, 4, 3, 2, 3,
	2, 3, 1, 2, 1, 2,
}
var yyChk = [...]int{

	-1000, -1, -3, -19, -21, -15, -2, 21, 6, 7,
	5, 16, 17, -16, -17, -20, -14, 34, 20, 9,
	11, -22, 19, -13, -10, -4, 10, 36, -24, 31,
	-9, 14, 32, 17, 18, 12, 4, 13, -26, -28,
	-13, -10, -12, 47, 52, 15, 53, -22, 34, -4,
	33, -9, -9, 33, -1, -9, 31, -9, 31, -8,
	31, 35, 58, -25, 36, 56, -23, -24, -9, -11,
	22, 23, 39, 40, 44, 45, 30, 29, 24, 25,
	46, 47, 48, 49, 50, 43, 26, 27, 28, 41,
	42, 33, -9, -9, 31, -9, -9, -9, -9, -9,
	-27, -5, 56, 36, 31, 55, -29, -7, -6, -9,
	33, 55, -18, 7, -15, -19, -16, -17, -2, 56,
	59, -18, -27, 59, 56, -9, 60, 31, 61, -6,
	-9, -9, 10, 19, 61, 61, -9, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, -9, -9, -9, -9,
	-9, -9, -9, -9, -9, 17, -2, 59, 56, -9,
	61, -5, 55, 59, 59, 60, 56, -9, 31, 8,
	-18, 31, -9, 57, 60, -9, -11, 59, 61, -9,
	31, -9, 61, -9, -9, -9, -18, 59, 56, -18,
	56, -11, -9, 57, 57, 60, -18, -9, -9, -9,
	57, -9, -18, 59, -18, -9, -18,
}
var yyDef = [...]int{

	1, -2, 2, 5, 6, 7, 8, 0, 39, 40,
	0, 0, 0, 14, 15, 16, 4, 1, 0, 0,
	0, 0, 0, 18, 19, 98, 33, 0, 100, 45,
	38, 61, 62, 0, 0, 0, 0, 0, 69, 70,
	71, 72, 73, 0, 0, 0, 0, 0, 0, 98,
	97, 41, 43, 44, 0, 0, 0, 0, 0, 17,
	57, 0, 0, 102, 0, 0, 36, 100, 0, 23,
	21, 22, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 63, 64, 65, 66, 67, 93, 94, 95, 96,
	0, 0, 0, 0, 51, 110, 0, 112, 114, 53,
	42, 3, 26, 40, 9, 10, 11, 12, 13, 0,
	0, 31, 0, 0, 0, 0, 0, 50, 103, 0,
	53, 20, 34, 35, 99, 101, 74, 75, 76, 77,
	78, 79, 80, 81, 82, 83, 84, 85, 86, 87,
	88, 89, 90, 91, 92, 0, 105, 0, 0, 107,
	108, 0, 111, 113, 115, 0, 0, 0, 0, 0,
	37, 60, 58, 46, 0, 0, 25, 0, 104, 68,
	52, 106, 109, 0, 54, 55, 0, 0, 0, 32,
	0, 24, 0, 48, 49, 0, 27, 0, 0, 59,
	47, 56, 28, 0, 30, 0, 29,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 53, 3, 50, 42, 3,
	36, 61, 48, 46, 59, 47, 58, 49, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 60, 3,
	45, 56, 44, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 37, 3, 3, 3, 3, 3,
	3, 35, 3, 57, 43, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 34, 41, 55, 52,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 38, 39, 40, 51, 54,
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
		//line .\parser.go.y:76
		{
			yyVAL.expr = CNode("chain")
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:82
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[2].expr)
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 3:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:90
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:93
		{
			if yyDollar[1].expr.isIsolatedCopy() {
				yyDollar[1].expr.Cx(2).C()[0] = NNode(0.0)
			}
			yyVAL.expr = yyDollar[1].expr
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:101
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:102
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:103
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:104
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:107
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:108
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:109
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:110
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:111
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:114
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:115
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:116
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 17:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:119
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:120
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:121
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:122
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
		//line .\parser.go.y:142
		{
			yyVAL.expr = NNode(1.0)
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:143
		{
			yyVAL.expr = NNode(-1.0)
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:146
		{
			yyVAL.expr = CNode("inc", ANode(yyDollar[1].token).setPos(yyDollar[1].token), yyDollar[2].expr)
		}
	case 24:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:147
		{
			yyVAL.expr = CNode("store", yyDollar[1].expr, yyDollar[3].expr, CNode("+", CNode("load", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr), yyDollar[5].expr).setPos0(yyDollar[1].expr))
		}
	case 25:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:148
		{
			yyVAL.expr = CNode("store", yyDollar[1].expr, yyDollar[3].token, CNode("+", CNode("load", yyDollar[1].expr, yyDollar[3].token).setPos0(yyDollar[1].expr), yyDollar[4].expr).setPos0(yyDollar[1].expr))
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:151
		{
			yyVAL.expr = CNode("for", yyDollar[2].expr, CNode(), yyDollar[3].expr).setPos0(yyDollar[1].token)
		}
	case 27:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line .\parser.go.y:154
		{
			yyVAL.expr = CNode("for", yyDollar[2].expr, yyDollar[5].expr, yyDollar[6].expr).setPos0(yyDollar[1].token)
		}
	case 28:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line .\parser.go.y:157
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
		//line .\parser.go.y:171
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
		//line .\parser.go.y:211
		{
			yyVAL.expr = CNode("call", "copy", CNode(
				NNode(0),
				yyDollar[6].expr,
				CNode("func", "<anony-map-iter-callback>", CNode(yyDollar[2].token.Str, yyDollar[4].token.Str), yyDollar[7].expr),
			))
		}
	case 31:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:220
		{
			yyVAL.expr = CNode("if", yyDollar[2].expr, yyDollar[3].expr, CNode())
		}
	case 32:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:221
		{
			yyVAL.expr = CNode("if", yyDollar[2].expr, yyDollar[3].expr, yyDollar[5].expr)
		}
	case 33:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line .\parser.go.y:224
		{
			yyVAL.str = ""
		}
	case 34:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:225
		{
			yyVAL.str = yyDollar[1].str + ",safe"
		}
	case 35:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:226
		{
			yyVAL.str = yyDollar[1].str + ",var"
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:229
		{
			yyVAL.str = "func," + yyDollar[2].str
		}
	case 37:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:232
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
		//line .\parser.go.y:244
		{
			yyVAL.expr = CNode("yield", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:245
		{
			yyVAL.expr = CNode("break").setPos0(yyDollar[1].token)
		}
	case 40:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:246
		{
			yyVAL.expr = CNode("continue").setPos0(yyDollar[1].token)
		}
	case 41:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:247
		{
			yyVAL.expr = CNode("assert", yyDollar[2].expr, NilNode()).setPos0(yyDollar[1].token)
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:248
		{
			yyVAL.expr = CNode("assert", yyDollar[2].expr, SNode(yyDollar[3].token.Str)).setPos0(yyDollar[1].token)
		}
	case 43:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:249
		{
			yyVAL.expr = CNode("ret", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 44:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:250
		{
			yyVAL.expr = yylex.(*Lexer).loadFile(filepath.Join(filepath.Dir(yyDollar[1].token.Pos.Source), yyDollar[2].token.Str))
		}
	case 45:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:253
		{
			yyVAL.expr = ANode(yyDollar[1].token).setPos(yyDollar[1].token)
		}
	case 46:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:254
		{
			yyVAL.expr = CNode("load", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 47:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line .\parser.go.y:255
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, yyDollar[3].expr, yyDollar[5].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 48:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:256
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, yyDollar[3].expr, NNode("-1")).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 49:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:257
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, NNode("0"), yyDollar[4].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 50:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:258
		{
			yyVAL.expr = CNode("load", yyDollar[1].expr, SNode(yyDollar[3].token.Str)).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:261
		{
			yyVAL.expr = CNode(yyDollar[1].token.Str)
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:262
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(ANode(yyDollar[3].token))
		}
	case 53:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:265
		{
			yyVAL.expr = CNode(yyDollar[1].expr)
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:266
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[3].expr)
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:269
		{
			yyVAL.expr = CNode(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 56:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:270
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[3].expr).Cappend(yyDollar[5].expr)
		}
	case 57:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:273
		{
			yyVAL.expr = CNode("chain", CNode("set", ANode(yyDollar[1].token), NilNode()).setPos0(yyDollar[1].token))
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:274
		{
			yyVAL.expr = CNode("chain", CNode("set", ANode(yyDollar[1].token), yyDollar[3].expr).setPos0(yyDollar[1].token))
		}
	case 59:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:275
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(CNode("set", ANode(yyDollar[3].token), yyDollar[5].expr).setPos0(yyDollar[1].expr))
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:276
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(CNode("set", ANode(yyDollar[3].token), NilNode()).setPos0(yyDollar[1].expr))
		}
	case 61:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:279
		{
			yyVAL.expr = NilNode().SetPos(yyDollar[1].token)
		}
	case 62:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:280
		{
			yyVAL.expr = NNode(yyDollar[1].token.Str).SetPos(yyDollar[1].token)
		}
	case 63:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:281
		{
			yyVAL.expr = yylex.(*Lexer).loadFile(filepath.Join(filepath.Dir(yyDollar[1].token.Pos.Source), yyDollar[2].token.Str))
		}
	case 64:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:282
		{
			yyVAL.expr = CNode("typeof", yyDollar[2].expr)
		}
	case 65:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:283
		{
			yyVAL.expr = CNode("len", yyDollar[2].expr)
		}
	case 66:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:284
		{
			yyVAL.expr = CNode("call", "addressof", CNode(ANode(yyDollar[2].token)))
		}
	case 67:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:285
		{
			yyVAL.expr = CNode("call", "copy", CNode(NNode(1), yyDollar[2].expr, NilNode()))
		}
	case 68:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:286
		{
			yyVAL.expr = CNode("call", "copy", CNode(NNode(1), yyDollar[2].expr, yyDollar[4].expr))
		}
	case 69:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:287
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 70:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:288
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 71:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:289
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:290
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 73:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:291
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:292
		{
			yyVAL.expr = CNode("or", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 75:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:293
		{
			yyVAL.expr = CNode("and", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:294
		{
			yyVAL.expr = CNode("<", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 77:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:295
		{
			yyVAL.expr = CNode("<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 78:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:296
		{
			yyVAL.expr = CNode("<=", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 79:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:297
		{
			yyVAL.expr = CNode("<=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 80:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:298
		{
			yyVAL.expr = CNode("==", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 81:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:299
		{
			yyVAL.expr = CNode("!=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 82:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:300
		{
			yyVAL.expr = CNode("+", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 83:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:301
		{
			yyVAL.expr = CNode("-", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 84:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:302
		{
			yyVAL.expr = CNode("*", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:303
		{
			yyVAL.expr = CNode("/", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 86:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:304
		{
			yyVAL.expr = CNode("%", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 87:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:305
		{
			yyVAL.expr = CNode("^", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 88:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:306
		{
			yyVAL.expr = CNode("<<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 89:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:307
		{
			yyVAL.expr = CNode(">>", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 90:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:308
		{
			yyVAL.expr = CNode(">>>", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 91:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:309
		{
			yyVAL.expr = CNode("|", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 92:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:310
		{
			yyVAL.expr = CNode("&", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 93:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:311
		{
			yyVAL.expr = CNode("-", NNode(0.0), yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 94:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:312
		{
			yyVAL.expr = CNode("~", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 95:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:313
		{
			yyVAL.expr = CNode("!", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 96:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:314
		{
			yyVAL.expr = CNode("#", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 97:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:317
		{
			yyVAL.expr = SNode(yyDollar[1].token.Str).SetPos(yyDollar[1].token)
		}
	case 98:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:320
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 99:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:321
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 100:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:322
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 101:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:323
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 102:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:326
		{
			yyVAL.expr = CNode("call", yyDollar[1].expr, yyDollar[2].expr).setPos0(yyDollar[1].expr)
		}
	case 103:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:331
		{
			yyVAL.expr = CNode()
		}
	case 104:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:332
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 105:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:335
		{
			yyVAL.expr = CNode(yyDollar[1].str, "<a>", yyDollar[2].expr, yyDollar[3].expr).setPos0(yyDollar[2].expr)
		}
	case 106:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:336
		{
			yyVAL.expr = CNode(yyDollar[1].str, "<a>", yyDollar[2].expr, CNode("chain", CNode("ret", yyDollar[4].expr))).setPos0(yyDollar[2].expr)
		}
	case 107:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:337
		{
			yyVAL.expr = CNode(yyDollar[1].str, "<a>", CNode(), CNode("chain", CNode("ret", yyDollar[3].expr))).setPos0(yyDollar[3].expr)
		}
	case 108:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:340
		{
			yyVAL.expr = CNode()
		}
	case 109:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:341
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 110:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:344
		{
			yyVAL.expr = CNode("map", CNode()).setPos0(yyDollar[1].token)
		}
	case 111:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:345
		{
			yyVAL.expr = yyDollar[2].expr.setPos0(yyDollar[1].token)
		}
	case 112:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:348
		{
			yyVAL.expr = CNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 113:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:349
		{
			yyVAL.expr = CNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 114:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:350
		{
			yyVAL.expr = CNode("array", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 115:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:351
		{
			yyVAL.expr = CNode("array", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	}
	goto yystack /* stack new state and value */
}
