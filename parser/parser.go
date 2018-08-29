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
const TNil = 57354
const TNot = 57355
const TReturn = 57356
const TRequire = 57357
const TTypeof = 57358
const TVar = 57359
const TWhile = 57360
const TYield = 57361
const TAddAdd = 57362
const TSubSub = 57363
const TEqeq = 57364
const TNeq = 57365
const TLsh = 57366
const TRsh = 57367
const TURsh = 57368
const TLte = 57369
const TGte = 57370
const TIdent = 57371
const TNumber = 57372
const TString = 57373
const FUN = 57374
const TOr = 57375
const TAnd = 57376
const UNARY = 57377
const TMinMin = 57378

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

//line .\parser.go.y:369

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

const yyLast = 905

var yyAct = [...]int{

	177, 97, 67, 38, 151, 23, 175, 104, 30, 129,
	100, 49, 50, 39, 171, 24, 172, 152, 158, 53,
	151, 55, 114, 59, 62, 115, 157, 118, 66, 100,
	183, 181, 160, 119, 99, 90, 63, 47, 96, 25,
	154, 156, 92, 93, 94, 95, 60, 105, 99, 113,
	107, 6, 17, 89, 98, 51, 23, 23, 173, 23,
	120, 165, 162, 125, 126, 122, 24, 24, 91, 24,
	124, 131, 132, 133, 134, 135, 136, 137, 138, 139,
	140, 141, 142, 143, 144, 145, 146, 147, 148, 149,
	25, 25, 45, 25, 21, 117, 110, 58, 3, 153,
	56, 155, 6, 112, 54, 14, 116, 163, 111, 109,
	13, 5, 84, 85, 86, 161, 80, 81, 82, 28,
	166, 23, 169, 68, 69, 170, 78, 79, 80, 81,
	82, 24, 78, 79, 80, 81, 82, 127, 102, 1,
	37, 36, 61, 64, 128, 21, 150, 65, 4, 3,
	15, 16, 40, 174, 57, 25, 14, 52, 176, 103,
	178, 13, 5, 2, 23, 0, 0, 23, 164, 185,
	184, 0, 0, 0, 24, 0, 0, 24, 0, 0,
	0, 190, 191, 23, 192, 0, 0, 0, 0, 194,
	0, 0, 0, 24, 23, 23, 0, 198, 25, 0,
	0, 25, 23, 0, 24, 24, 0, 0, 0, 0,
	0, 179, 24, 0, 182, 0, 0, 25, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 25, 25,
	189, 0, 0, 0, 0, 0, 25, 0, 0, 0,
	0, 195, 197, 0, 10, 8, 9, 0, 19, 199,
	20, 0, 0, 11, 12, 0, 22, 18, 7, 0,
	0, 76, 77, 84, 85, 86, 75, 74, 29, 0,
	0, 17, 0, 27, 0, 0, 70, 71, 87, 88,
	83, 72, 73, 78, 79, 80, 81, 82, 10, 8,
	9, 0, 19, 0, 20, 0, 196, 11, 12, 0,
	22, 18, 7, 0, 0, 76, 77, 84, 85, 86,
	75, 74, 29, 0, 0, 17, 0, 27, 0, 0,
	70, 71, 87, 88, 83, 72, 73, 78, 79, 80,
	81, 82, 10, 8, 108, 0, 19, 0, 20, 0,
	0, 11, 12, 0, 22, 18, 7, 0, 0, 76,
	77, 84, 85, 86, 75, 74, 29, 0, 0, 17,
	0, 27, 0, 0, 70, 71, 87, 88, 83, 72,
	73, 78, 79, 80, 81, 82, 76, 77, 84, 85,
	86, 75, 74, 0, 0, 0, 0, 0, 0, 0,
	0, 70, 71, 87, 88, 83, 72, 73, 78, 79,
	80, 81, 82, 0, 0, 0, 0, 0, 0, 167,
	0, 0, 168, 76, 77, 84, 85, 86, 75, 74,
	0, 0, 0, 0, 0, 0, 0, 0, 70, 71,
	87, 88, 83, 72, 73, 78, 79, 80, 81, 82,
	76, 77, 84, 85, 86, 75, 74, 0, 0, 0,
	130, 0, 0, 0, 0, 70, 71, 87, 88, 83,
	72, 73, 78, 79, 80, 81, 82, 76, 77, 84,
	85, 86, 75, 74, 0, 0, 188, 0, 0, 0,
	0, 0, 70, 71, 87, 88, 83, 72, 73, 78,
	79, 80, 81, 82, 76, 77, 84, 85, 86, 75,
	74, 0, 0, 159, 0, 0, 0, 0, 0, 70,
	71, 87, 88, 83, 72, 73, 78, 79, 80, 81,
	82, 76, 77, 84, 85, 86, 75, 74, 0, 180,
	0, 0, 0, 0, 0, 0, 70, 71, 87, 88,
	83, 72, 73, 78, 79, 80, 81, 82, 0, 0,
	0, 0, 0, 0, 193, 76, 77, 84, 85, 86,
	75, 74, 0, 0, 0, 0, 0, 0, 0, 0,
	70, 71, 87, 88, 83, 72, 73, 78, 79, 80,
	81, 82, 0, 35, 0, 0, 0, 0, 187, 26,
	35, 31, 43, 0, 33, 34, 26, 0, 31, 43,
	0, 33, 34, 0, 0, 0, 0, 0, 29, 32,
	48, 46, 0, 27, 0, 29, 32, 48, 46, 0,
	27, 0, 0, 0, 41, 0, 0, 0, 0, 42,
	44, 41, 0, 0, 35, 0, 42, 44, 123, 0,
	26, 0, 31, 43, 121, 33, 34, 0, 76, 77,
	84, 85, 86, 75, 74, 0, 0, 0, 0, 29,
	32, 48, 46, 0, 27, 87, 88, 83, 72, 73,
	78, 79, 80, 81, 82, 41, 0, 0, 0, 0,
	42, 44, 10, 8, 9, 186, 19, 26, 20, 0,
	0, 11, 12, 35, 22, 18, 7, 0, 0, 26,
	0, 31, 43, 0, 33, 34, 29, 0, 0, 17,
	0, 27, 0, 0, 0, 0, 0, 0, 29, 32,
	48, 46, 0, 27, 0, 0, 0, 0, 0, 0,
	106, 0, 0, 0, 41, 0, 0, 0, 0, 42,
	44, 0, 101, 76, 77, 84, 85, 86, 75, 74,
	0, 0, 0, 0, 0, 0, 0, 0, 70, 71,
	87, 88, 83, 72, 73, 78, 79, 80, 81, 82,
	76, 77, 84, 85, 86, 75, 74, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 71, 87, 88, 83,
	72, 73, 78, 79, 80, 81, 82, 35, 0, 0,
	0, 0, 0, 26, 0, 31, 43, 0, 33, 34,
	76, 77, 84, 85, 86, 75, 74, 0, 0, 0,
	0, 0, 29, 32, 48, 46, 0, 27, 0, 0,
	72, 73, 78, 79, 80, 81, 82, 0, 41, 0,
	0, 0, 0, 42, 44, 10, 8, 9, 0, 19,
	26, 20, 0, 0, 11, 12, 0, 22, 18, 7,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 29,
	0, 0, 17, 0, 27, 10, 8, 9, 0, 19,
	0, 20, 0, 0, 11, 12, 0, 22, 18, 7,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 29,
	0, 0, 17, 0, 27,
}
var yyPact = [...]int{

	-1000, 840, -1000, -1000, -1000, -1000, -1000, 793, -1000, -1000,
	793, 793, 24, -1000, -1000, -1000, -1000, -1000, 793, 75,
	793, 71, 68, -10, -1000, -18, -1000, 793, -1000, 103,
	721, -1000, -1000, 22, 793, 39, -1000, -1000, -10, -1000,
	-1000, 793, 793, 793, 793, 0, 689, -1000, -1000, 721,
	721, -1000, 677, 327, -32, 283, 14, -30, -21, 586,
	36, -1000, 579, 793, 127, -50, 391, -1000, -1000, -1000,
	793, 793, 793, 793, 793, 793, 793, 793, 793, 793,
	793, 793, 793, 793, 793, 793, 793, 793, 793, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 20, -37, 793, -19,
	-1000, -1000, -12, -31, -39, 445, -1000, -1000, -22, -1000,
	-1000, -1000, -1000, -1000, 793, 33, 99, 870, 32, 793,
	354, 793, 103, -1000, -43, 721, 721, -1000, -1000, -1000,
	-1000, 748, 626, 88, 88, 88, 88, 88, 88, 70,
	70, -1000, -1000, -1000, 788, 82, 82, 82, 788, 788,
	-1000, 29, 793, 721, -1000, -53, -1000, 793, 793, 793,
	870, 472, -23, 870, -1000, -24, 721, 103, 630, 533,
	-1000, 793, -1000, -1000, 721, -1000, 418, 721, 721, 870,
	793, 793, -1000, 793, -1000, 499, -1000, -1000, 793, -1000,
	239, 283, 721, -1000, 721, -1000, 793, -1000, 283, -1000,
}
var yyPgo = [...]int{

	0, 139, 49, 163, 37, 1, 7, 159, 154, 0,
	13, 2, 152, 3, 151, 109, 108, 103, 50, 96,
	150, 148, 92, 143, 119, 142, 141, 38, 140, 138,
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
	9, 9, 9, 12, 13, 13, 13, 13, 24, 25,
	25, 26, 26, 26, 27, 27, 28, 28, 29, 29,
	29, 29,
}
var yyR2 = [...]int{

	0, 0, 2, 3, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 2, 1, 1,
	3, 1, 1, 2, 5, 4, 3, 6, 7, 9,
	7, 3, 5, 0, 2, 2, 2, 4, 2, 1,
	1, 2, 2, 2, 1, 4, 6, 5, 5, 3,
	1, 3, 1, 3, 3, 5, 1, 3, 5, 3,
	1, 1, 2, 2, 2, 1, 1, 1, 1, 1,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 2,
	2, 2, 2, 1, 1, 3, 1, 3, 2, 2,
	3, 3, 4, 3, 2, 3, 2, 3, 1, 2,
	1, 2,
}
var yyChk = [...]int{

	-1000, -1, -3, -19, -21, -15, -2, 19, 6, 7,
	5, 14, 15, -16, -17, -20, -14, 32, 18, 9,
	11, -22, 17, -13, -10, -4, 10, 34, -24, 29,
	-9, 12, 30, 15, 16, 4, -26, -28, -13, -10,
	-12, 45, 50, 13, 51, -22, 32, -4, 31, -9,
	-9, 31, -1, -9, 29, -9, 29, -8, 29, 33,
	56, -25, 34, 54, -23, -24, -9, -11, 20, 21,
	37, 38, 42, 43, 28, 27, 22, 23, 44, 45,
	46, 47, 48, 41, 24, 25, 26, 39, 40, 31,
	-9, 29, -9, -9, -9, -9, -27, -5, 54, 34,
	29, 53, -29, -7, -6, -9, 53, -18, 7, -15,
	-19, -16, -17, -2, 54, 57, -18, -27, 57, 54,
	-9, 58, 29, 59, -6, -9, -9, 10, 17, 59,
	59, -9, -9, -9, -9, -9, -9, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, -9, -9, -9, -9,
	-2, 57, 54, -9, 59, -5, 53, 57, 57, 58,
	54, -9, 29, 8, -18, 29, -9, 55, 58, -9,
	-11, 57, 59, 29, -9, 59, -9, -9, -9, -18,
	57, 54, -18, 54, -11, -9, 55, 55, 58, -18,
	-9, -9, -9, 55, -9, -18, 57, -18, -9, -18,
}
var yyDef = [...]int{

	1, -2, 2, 5, 6, 7, 8, 0, 39, 40,
	0, 0, 0, 14, 15, 16, 4, 1, 0, 0,
	0, 0, 0, 18, 19, 94, 33, 0, 96, 44,
	38, 60, 61, 0, 0, 0, 65, 66, 67, 68,
	69, 0, 0, 0, 0, 0, 0, 94, 93, 41,
	42, 43, 0, 0, 0, 0, 0, 17, 56, 0,
	0, 98, 0, 0, 36, 96, 0, 23, 21, 22,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 62,
	63, 64, 89, 90, 91, 92, 0, 0, 0, 0,
	50, 106, 0, 108, 110, 52, 3, 26, 40, 9,
	10, 11, 12, 13, 0, 0, 31, 0, 0, 0,
	0, 0, 49, 99, 0, 52, 20, 34, 35, 95,
	97, 70, 71, 72, 73, 74, 75, 76, 77, 78,
	79, 80, 81, 82, 83, 84, 85, 86, 87, 88,
	101, 0, 0, 103, 104, 0, 107, 109, 111, 0,
	0, 0, 0, 0, 37, 59, 57, 45, 0, 0,
	25, 0, 100, 51, 102, 105, 0, 53, 54, 0,
	0, 0, 32, 0, 24, 0, 47, 48, 0, 27,
	0, 0, 58, 46, 55, 28, 0, 30, 0, 29,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 51, 3, 48, 40, 3,
	34, 59, 46, 44, 57, 45, 56, 47, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 58, 3,
	43, 54, 42, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 35, 3, 3, 3, 3, 3,
	3, 33, 3, 55, 41, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 32, 39, 53, 50,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	36, 37, 38, 49, 52,
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
			yyVAL.expr = CNode("assert", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 42:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:248
		{
			yyVAL.expr = CNode("ret", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 43:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:249
		{
			yyVAL.expr = yylex.(*Lexer).loadFile(filepath.Join(filepath.Dir(yyDollar[1].token.Pos.Source), yyDollar[2].token.Str))
		}
	case 44:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:252
		{
			yyVAL.expr = ANode(yyDollar[1].token).setPos(yyDollar[1].token)
		}
	case 45:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:253
		{
			yyVAL.expr = CNode("load", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 46:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line .\parser.go.y:254
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, yyDollar[3].expr, yyDollar[5].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 47:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:255
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, yyDollar[3].expr, NNode("-1")).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 48:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:256
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, NNode("0"), yyDollar[4].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 49:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:257
		{
			yyVAL.expr = CNode("load", yyDollar[1].expr, SNode(yyDollar[3].token.Str)).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 50:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:260
		{
			yyVAL.expr = CNode(yyDollar[1].token.Str)
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:261
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(ANode(yyDollar[3].token))
		}
	case 52:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:264
		{
			yyVAL.expr = CNode(yyDollar[1].expr)
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:265
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[3].expr)
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:268
		{
			yyVAL.expr = CNode(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 55:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:269
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[3].expr).Cappend(yyDollar[5].expr)
		}
	case 56:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:272
		{
			yyVAL.expr = CNode("chain", CNode("set", ANode(yyDollar[1].token), NilNode()).setPos0(yyDollar[1].token))
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:273
		{
			yyVAL.expr = CNode("chain", CNode("set", ANode(yyDollar[1].token), yyDollar[3].expr).setPos0(yyDollar[1].token))
		}
	case 58:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:274
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(CNode("set", ANode(yyDollar[3].token), yyDollar[5].expr).setPos0(yyDollar[1].expr))
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:275
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(CNode("set", ANode(yyDollar[3].token), NilNode()).setPos0(yyDollar[1].expr))
		}
	case 60:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:278
		{
			yyVAL.expr = NilNode().SetPos(yyDollar[1].token)
		}
	case 61:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:279
		{
			yyVAL.expr = NNode(yyDollar[1].token.Str).SetPos(yyDollar[1].token)
		}
	case 62:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:280
		{
			yyVAL.expr = yylex.(*Lexer).loadFile(filepath.Join(filepath.Dir(yyDollar[1].token.Pos.Source), yyDollar[2].token.Str))
		}
	case 63:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:281
		{
			yyVAL.expr = CNode("typeof", yyDollar[2].expr)
		}
	case 64:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:282
		{
			yyVAL.expr = CNode("call", "addressof", CNode(ANode(yyDollar[2].token)))
		}
	case 65:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:283
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:284
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 67:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:285
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 68:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:286
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 69:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:287
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 70:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:288
		{
			yyVAL.expr = CNode("or", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:289
		{
			yyVAL.expr = CNode("and", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 72:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:290
		{
			yyVAL.expr = CNode("<", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:291
		{
			yyVAL.expr = CNode("<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:292
		{
			yyVAL.expr = CNode("<=", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 75:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:293
		{
			yyVAL.expr = CNode("<=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:294
		{
			yyVAL.expr = CNode("==", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 77:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:295
		{
			yyVAL.expr = CNode("!=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 78:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:296
		{
			yyVAL.expr = CNode("+", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 79:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:297
		{
			yyVAL.expr = CNode("-", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 80:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:298
		{
			yyVAL.expr = CNode("*", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 81:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:299
		{
			yyVAL.expr = CNode("/", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 82:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:300
		{
			yyVAL.expr = CNode("%", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 83:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:301
		{
			yyVAL.expr = CNode("^", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 84:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:302
		{
			yyVAL.expr = CNode("<<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:303
		{
			yyVAL.expr = CNode(">>", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 86:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:304
		{
			yyVAL.expr = CNode(">>>", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 87:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:305
		{
			yyVAL.expr = CNode("|", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 88:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:306
		{
			yyVAL.expr = CNode("&", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 89:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:307
		{
			yyVAL.expr = CNode("-", NNode(0.0), yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 90:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:308
		{
			yyVAL.expr = CNode("~", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 91:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:309
		{
			yyVAL.expr = CNode("!", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 92:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:310
		{
			yyVAL.expr = CNode("#", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 93:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:313
		{
			yyVAL.expr = SNode(yyDollar[1].token.Str).SetPos(yyDollar[1].token)
		}
	case 94:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:316
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 95:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:317
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 96:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:318
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 97:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:319
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 98:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:322
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
	case 99:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:347
		{
			yyVAL.expr = CNode()
		}
	case 100:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:348
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 101:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:351
		{
			yyVAL.expr = CNode(yyDollar[1].str, "<a>", yyDollar[2].expr, yyDollar[3].expr).setPos0(yyDollar[2].expr)
		}
	case 102:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:352
		{
			yyVAL.expr = CNode(yyDollar[1].str, "<a>", yyDollar[2].expr, CNode("chain", CNode("ret", yyDollar[4].expr))).setPos0(yyDollar[2].expr)
		}
	case 103:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:353
		{
			yyVAL.expr = CNode(yyDollar[1].str, "<a>", CNode(), CNode("chain", CNode("ret", yyDollar[3].expr))).setPos0(yyDollar[3].expr)
		}
	case 104:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:356
		{
			yyVAL.expr = CNode()
		}
	case 105:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:357
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 106:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:360
		{
			yyVAL.expr = CNode("map", CNode()).setPos0(yyDollar[1].token)
		}
	case 107:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:361
		{
			yyVAL.expr = yyDollar[2].expr.setPos0(yyDollar[1].token)
		}
	case 108:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:364
		{
			yyVAL.expr = CNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 109:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:365
		{
			yyVAL.expr = CNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 110:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:366
		{
			yyVAL.expr = CNode("array", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 111:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:367
		{
			yyVAL.expr = CNode("array", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	}
	goto yystack /* stack new state and value */
}
