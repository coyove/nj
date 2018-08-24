//line .\parser.go.y:1
package parser

import __yyfmt__ "fmt"

//line .\parser.go.y:3
import (
	"fmt"
	"github.com/coyove/common/rand"
	"path/filepath"
)

//line .\parser.go.y:42
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
	"'!'",
	"':'",
	"')'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line .\parser.go.y:397

var typesLookup = map[string]string{
	"nil": "0", "number": "1", "string": "2", "map": "4", "closure": "6", "generic": "7",
}

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

const yyLast = 857

var yyAct = [...]int{

	173, 95, 102, 37, 67, 23, 147, 167, 31, 171,
	168, 48, 50, 38, 125, 24, 64, 98, 154, 53,
	148, 55, 112, 147, 153, 113, 59, 62, 66, 116,
	179, 177, 156, 98, 117, 94, 63, 46, 97, 25,
	152, 90, 91, 92, 93, 111, 103, 6, 150, 60,
	105, 78, 79, 80, 81, 82, 23, 23, 96, 23,
	118, 97, 17, 123, 124, 122, 24, 24, 89, 24,
	51, 127, 128, 129, 130, 131, 132, 133, 134, 135,
	136, 137, 138, 139, 140, 141, 142, 143, 144, 145,
	25, 25, 115, 25, 44, 169, 21, 149, 6, 151,
	161, 158, 26, 120, 32, 42, 114, 34, 58, 56,
	84, 85, 86, 157, 108, 159, 3, 54, 162, 23,
	165, 29, 33, 47, 45, 166, 27, 28, 100, 24,
	78, 79, 80, 81, 82, 68, 69, 40, 36, 110,
	146, 14, 41, 43, 1, 35, 109, 21, 13, 170,
	61, 4, 121, 25, 172, 65, 174, 15, 16, 107,
	23, 5, 52, 23, 39, 181, 160, 3, 180, 57,
	24, 101, 2, 24, 80, 81, 82, 186, 187, 23,
	188, 0, 0, 0, 0, 190, 0, 0, 0, 24,
	23, 23, 14, 194, 25, 0, 0, 25, 23, 13,
	24, 24, 0, 0, 0, 0, 0, 175, 24, 0,
	178, 0, 5, 25, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 25, 25, 185, 0, 0, 0,
	0, 0, 25, 0, 0, 0, 0, 191, 193, 0,
	10, 8, 9, 0, 19, 195, 20, 0, 0, 11,
	12, 0, 22, 18, 7, 0, 0, 76, 77, 84,
	85, 86, 75, 74, 29, 0, 0, 17, 0, 27,
	0, 0, 70, 71, 87, 88, 83, 72, 73, 78,
	79, 80, 81, 82, 10, 8, 9, 0, 19, 0,
	20, 0, 192, 11, 12, 0, 22, 18, 7, 0,
	0, 76, 77, 84, 85, 86, 75, 74, 29, 0,
	0, 17, 0, 27, 0, 0, 70, 71, 87, 88,
	83, 72, 73, 78, 79, 80, 81, 82, 76, 77,
	84, 85, 86, 75, 74, 0, 0, 0, 0, 0,
	0, 0, 0, 70, 71, 87, 88, 83, 72, 73,
	78, 79, 80, 81, 82, 0, 0, 0, 0, 0,
	0, 163, 0, 0, 0, 164, 76, 77, 84, 85,
	86, 75, 74, 0, 0, 0, 0, 0, 0, 0,
	0, 70, 71, 87, 88, 83, 72, 73, 78, 79,
	80, 81, 82, 10, 8, 106, 0, 19, 0, 20,
	0, 0, 11, 12, 126, 22, 18, 7, 0, 0,
	76, 77, 84, 85, 86, 75, 74, 29, 0, 0,
	17, 0, 27, 0, 0, 70, 71, 87, 88, 83,
	72, 73, 78, 79, 80, 81, 82, 76, 77, 84,
	85, 86, 75, 74, 0, 0, 0, 0, 0, 0,
	0, 0, 70, 71, 87, 88, 83, 72, 73, 78,
	79, 80, 81, 82, 76, 77, 84, 85, 86, 75,
	74, 0, 0, 0, 184, 0, 0, 0, 0, 70,
	71, 87, 88, 83, 72, 73, 78, 79, 80, 81,
	82, 76, 77, 84, 85, 86, 75, 74, 0, 0,
	0, 155, 0, 0, 0, 0, 70, 71, 87, 88,
	83, 72, 73, 78, 79, 80, 81, 82, 76, 77,
	84, 85, 86, 75, 74, 0, 176, 0, 0, 0,
	0, 0, 0, 70, 71, 87, 88, 83, 72, 73,
	78, 79, 80, 81, 82, 0, 0, 0, 0, 0,
	0, 189, 76, 77, 84, 85, 86, 75, 74, 0,
	26, 0, 32, 42, 0, 34, 0, 70, 71, 87,
	88, 83, 72, 73, 78, 79, 80, 81, 82, 29,
	33, 47, 45, 0, 27, 183, 26, 0, 32, 42,
	0, 34, 0, 0, 26, 40, 32, 42, 0, 34,
	41, 43, 0, 0, 0, 29, 33, 47, 45, 119,
	27, 0, 0, 29, 33, 47, 45, 26, 27, 32,
	42, 40, 34, 0, 0, 0, 41, 43, 0, 40,
	0, 0, 49, 0, 41, 43, 29, 33, 47, 45,
	30, 27, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 40, 0, 0, 0, 0, 41, 43, 10,
	8, 9, 182, 19, 26, 20, 0, 0, 11, 12,
	0, 22, 18, 7, 76, 77, 84, 85, 86, 75,
	74, 0, 26, 29, 32, 42, 17, 34, 27, 70,
	71, 87, 88, 83, 72, 73, 78, 79, 80, 81,
	82, 29, 33, 47, 45, 0, 27, 104, 76, 77,
	84, 85, 86, 75, 74, 0, 26, 40, 32, 42,
	0, 34, 41, 43, 71, 87, 88, 83, 72, 73,
	78, 79, 80, 81, 82, 29, 33, 47, 45, 0,
	27, 0, 0, 76, 77, 84, 85, 86, 75, 74,
	0, 40, 0, 0, 0, 0, 41, 43, 0, 99,
	87, 88, 83, 72, 73, 78, 79, 80, 81, 82,
	76, 77, 84, 85, 86, 75, 74, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	72, 73, 78, 79, 80, 81, 82, 10, 8, 9,
	0, 19, 26, 20, 0, 0, 11, 12, 0, 22,
	18, 7, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 29, 0, 0, 17, 0, 27, 10, 8, 9,
	0, 19, 0, 20, 0, 0, 11, 12, 0, 22,
	18, 7, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 29, 0, 0, 17, 0, 27,
}
var yyPact = [...]int{

	-1000, 793, -1000, -1000, -1000, -1000, -1000, 585, -1000, -1000,
	673, 577, 40, -1000, -1000, -1000, -1000, -1000, 673, 89,
	673, 81, 80, -6, -1000, -17, -41, 673, -1000, 116,
	-1000, 653, -1000, -1000, 38, -1000, -1000, -6, -1000, -1000,
	673, 673, 673, 673, 5, 707, -1000, -1000, 653, -1000,
	653, -1000, 655, 389, -31, 280, 28, -27, -19, 551,
	75, -1000, 93, 673, -1000, -45, 345, -1000, -1000, -1000,
	673, 673, 673, 673, 673, 673, 673, 673, 673, 673,
	673, 673, 673, 673, 673, 673, 673, 673, 673, -1000,
	-1000, -1000, -1000, -1000, 31, -33, 673, -11, -1000, -1000,
	-12, -32, -38, 443, -1000, -1000, -21, -1000, -1000, -1000,
	-1000, -1000, 673, 73, 108, 823, 72, 673, 307, 673,
	116, -1000, -49, 653, 653, -1000, -1000, 687, 722, 87,
	87, 87, 87, 87, 87, 129, 129, -1000, -1000, -1000,
	749, 8, 8, 8, 749, 749, -1000, 67, 673, 653,
	-1000, -50, -1000, 673, 673, 673, 823, 470, -22, 823,
	-1000, -23, 653, 116, 608, 531, -1000, 673, -1000, -1000,
	653, -1000, 416, 653, 653, 823, 673, 673, -1000, 673,
	-1000, 497, -1000, -1000, 673, -1000, 236, 280, 653, -1000,
	653, -1000, 673, -1000, 280, -1000,
}
var yyPgo = [...]int{

	0, 144, 45, 172, 37, 1, 2, 171, 169, 0,
	13, 4, 164, 3, 158, 159, 146, 139, 50, 114,
	157, 151, 94, 127, 150, 145, 35, 138, 128,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 15, 3, 3, 3, 3, 18,
	18, 18, 18, 18, 21, 21, 21, 14, 14, 14,
	14, 11, 11, 10, 10, 10, 16, 16, 16, 16,
	16, 17, 17, 22, 22, 20, 19, 19, 19, 19,
	19, 19, 19, 19, 4, 4, 4, 4, 4, 4,
	5, 5, 6, 6, 7, 7, 8, 8, 8, 8,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 12, 13, 13, 13, 13, 23, 24, 24, 25,
	25, 25, 26, 26, 27, 27, 28, 28, 28, 28,
}
var yyR2 = [...]int{

	0, 0, 2, 3, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 2, 1, 1,
	3, 1, 1, 2, 5, 4, 3, 6, 7, 9,
	7, 3, 5, 1, 2, 4, 2, 2, 1, 1,
	2, 2, 2, 2, 1, 4, 6, 5, 5, 3,
	1, 3, 1, 3, 3, 5, 1, 3, 5, 3,
	1, 1, 2, 1, 1, 1, 1, 1, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 2, 2, 2,
	2, 1, 1, 3, 1, 3, 2, 2, 3, 3,
	4, 3, 2, 3, 2, 3, 1, 2, 1, 2,
}
var yyChk = [...]int{

	-1000, -1, -3, -19, -21, -15, -2, 18, 5, 6,
	4, 13, 14, -16, -17, -20, -14, 31, 17, 8,
	10, -22, 16, -13, -10, -4, 9, 33, -23, 28,
	55, -9, 11, 29, 14, -25, -27, -13, -10, -12,
	44, 49, 12, 50, -22, 31, -4, 30, -9, 55,
	-9, 30, -1, -9, 28, -9, 28, -8, 28, 32,
	55, -24, 33, 53, 57, -23, -9, -11, 19, 20,
	36, 37, 41, 42, 27, 26, 21, 22, 43, 44,
	45, 46, 47, 40, 23, 24, 25, 38, 39, 30,
	-9, -9, -9, -9, -26, -5, 53, 33, 28, 52,
	-28, -7, -6, -9, 52, -18, 6, -15, -19, -16,
	-17, -2, 53, 56, -18, -26, 56, 53, -9, 58,
	28, 59, -6, -9, -9, 59, 59, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, -9, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, -2, 56, 53, -9,
	59, -5, 52, 56, 56, 58, 53, -9, 28, 7,
	-18, 28, -9, 54, 58, -9, -11, 56, 59, 28,
	-9, 59, -9, -9, -9, -18, 56, 53, -18, 53,
	-11, -9, 54, 54, 58, -18, -9, -9, -9, 54,
	-9, -18, 56, -18, -9, -18,
}
var yyDef = [...]int{

	1, -2, 2, 5, 6, 7, 8, 0, 38, 39,
	0, 0, 0, 14, 15, 16, 4, 1, 0, 0,
	0, 0, 0, 18, 19, 92, 33, 0, 94, 44,
	36, 37, 60, 61, 0, 63, 64, 65, 66, 67,
	0, 0, 0, 0, 0, 0, 92, 91, 40, 41,
	42, 43, 0, 0, 0, 0, 0, 17, 56, 0,
	0, 96, 0, 0, 34, 94, 0, 23, 21, 22,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 62,
	87, 88, 89, 90, 0, 0, 0, 0, 50, 104,
	0, 106, 108, 52, 3, 26, 39, 9, 10, 11,
	12, 13, 0, 0, 31, 0, 0, 0, 0, 0,
	49, 97, 0, 52, 20, 93, 95, 68, 69, 70,
	71, 72, 73, 74, 75, 76, 77, 78, 79, 80,
	81, 82, 83, 84, 85, 86, 99, 0, 0, 101,
	102, 0, 105, 107, 109, 0, 0, 0, 0, 0,
	35, 59, 57, 45, 0, 0, 25, 0, 98, 51,
	100, 103, 0, 53, 54, 0, 0, 0, 32, 0,
	24, 0, 47, 48, 0, 27, 0, 0, 58, 46,
	55, 28, 0, 30, 0, 29,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 57, 3, 50, 3, 47, 39, 3,
	33, 59, 45, 43, 56, 44, 55, 46, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 58, 3,
	42, 53, 41, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 34, 3, 3, 3, 3, 3,
	3, 32, 3, 54, 40, 3, 3, 3, 3, 3,
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
		//line .\parser.go.y:74
		{
			yyVAL.expr = CNode("chain")
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:80
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[2].expr)
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 3:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:88
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:91
		{
			if yyDollar[1].expr.isIsolatedCopy() {
				yyDollar[1].expr.Cx(2).C()[0] = NNode(0.0)
			}
			yyVAL.expr = yyDollar[1].expr
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:99
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:100
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:101
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:102
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:105
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:106
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:107
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:108
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:109
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:112
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:113
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:114
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 17:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:117
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:118
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:119
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:120
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
		//line .\parser.go.y:140
		{
			yyVAL.expr = NNode(1.0)
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:141
		{
			yyVAL.expr = NNode(-1.0)
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:144
		{
			yyVAL.expr = CNode("inc", ANode(yyDollar[1].token).setPos(yyDollar[1].token), yyDollar[2].expr)
		}
	case 24:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:145
		{
			yyVAL.expr = CNode("store", yyDollar[1].expr, yyDollar[3].expr, CNode("+", CNode("load", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr), yyDollar[5].expr).setPos0(yyDollar[1].expr))
		}
	case 25:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:146
		{
			yyVAL.expr = CNode("store", yyDollar[1].expr, yyDollar[3].token, CNode("+", CNode("load", yyDollar[1].expr, yyDollar[3].token).setPos0(yyDollar[1].expr), yyDollar[4].expr).setPos0(yyDollar[1].expr))
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:149
		{
			yyVAL.expr = CNode("for", yyDollar[2].expr, CNode(), yyDollar[3].expr).setPos0(yyDollar[1].token)
		}
	case 27:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line .\parser.go.y:152
		{
			yyVAL.expr = CNode("for", yyDollar[2].expr, yyDollar[5].expr, yyDollar[6].expr).setPos0(yyDollar[1].token)
		}
	case 28:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line .\parser.go.y:155
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
		//line .\parser.go.y:169
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
		//line .\parser.go.y:209
		{
			yyVAL.expr = CNode("call", "copy", CNode(
				NNode(0.0),
				yyDollar[6].expr,
				CNode("func", "<anony-map-iter-callback>", CNode(yyDollar[2].token.Str, yyDollar[4].token.Str), yyDollar[7].expr),
			))
		}
	case 31:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:218
		{
			yyVAL.expr = CNode("if", yyDollar[2].expr, yyDollar[3].expr, CNode())
		}
	case 32:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:219
		{
			yyVAL.expr = CNode("if", yyDollar[2].expr, yyDollar[3].expr, yyDollar[5].expr)
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:222
		{
			yyVAL.str = "func"
		}
	case 34:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:223
		{
			yyVAL.str = "safefunc"
		}
	case 35:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:226
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
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:238
		{
			yyVAL.expr = CNode("yield").setPos0(yyDollar[1].token)
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:239
		{
			yyVAL.expr = CNode("yield", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:240
		{
			yyVAL.expr = CNode("break").setPos0(yyDollar[1].token)
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:241
		{
			yyVAL.expr = CNode("continue").setPos0(yyDollar[1].token)
		}
	case 40:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:242
		{
			yyVAL.expr = CNode("assert", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 41:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:243
		{
			yyVAL.expr = CNode("ret").setPos0(yyDollar[1].token)
		}
	case 42:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:244
		{
			if yyDollar[2].expr.isIsolatedCopy() && yyDollar[2].expr.Cx(2).Cx(2).N() == 1 {
				yyDollar[2].expr.Cx(2).C()[2] = NNode(2.0)
			}
			yyVAL.expr = CNode("ret", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 43:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:250
		{
			path := filepath.Join(filepath.Dir(yyDollar[1].token.Pos.Source), yyDollar[2].token.Str)
			yyVAL.expr = yylex.(*Lexer).loadFile(path)
		}
	case 44:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:256
		{
			yyVAL.expr = ANode(yyDollar[1].token).setPos(yyDollar[1].token)
		}
	case 45:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:257
		{
			yyVAL.expr = CNode("load", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 46:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line .\parser.go.y:258
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, yyDollar[3].expr, yyDollar[5].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 47:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:259
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, yyDollar[3].expr, NNode("-1")).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 48:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:260
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, NNode("0"), yyDollar[4].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 49:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:261
		{
			yyVAL.expr = CNode("load", yyDollar[1].expr, SNode(yyDollar[3].token.Str)).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 50:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:264
		{
			yyVAL.expr = CNode(yyDollar[1].token.Str)
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:265
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(ANode(yyDollar[3].token))
		}
	case 52:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:268
		{
			yyVAL.expr = CNode(yyDollar[1].expr)
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:269
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[3].expr)
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:272
		{
			yyVAL.expr = CNode(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 55:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:273
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[3].expr).Cappend(yyDollar[5].expr)
		}
	case 56:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:276
		{
			yyVAL.expr = CNode("chain", CNode("set", ANode(yyDollar[1].token), NilNode()).setPos0(yyDollar[1].token))
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:277
		{
			yyVAL.expr = CNode("chain", CNode("set", ANode(yyDollar[1].token), yyDollar[3].expr).setPos0(yyDollar[1].token))
		}
	case 58:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:278
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(CNode("set", ANode(yyDollar[3].token), yyDollar[5].expr).setPos0(yyDollar[1].expr))
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:279
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(CNode("set", ANode(yyDollar[3].token), NilNode()).setPos0(yyDollar[1].expr))
		}
	case 60:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:282
		{
			yyVAL.expr = NilNode().SetPos(yyDollar[1].token)
		}
	case 61:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:283
		{
			yyVAL.expr = NNode(yyDollar[1].token.Str).SetPos(yyDollar[1].token)
		}
	case 62:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:284
		{
			yyVAL.expr = yylex.(*Lexer).loadFile(filepath.Join(filepath.Dir(yyDollar[1].token.Pos.Source), yyDollar[2].token.Str))
		}
	case 63:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:285
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 64:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:286
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 65:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:287
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:288
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 67:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:289
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:290
		{
			yyVAL.expr = CNode("or", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:291
		{
			yyVAL.expr = CNode("and", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 70:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:292
		{
			yyVAL.expr = CNode("<", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:293
		{
			yyVAL.expr = CNode("<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 72:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:294
		{
			yyVAL.expr = CNode("<=", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:295
		{
			yyVAL.expr = CNode("<=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:296
		{
			yyVAL.expr = CNode("==", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 75:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:297
		{
			yyVAL.expr = CNode("!=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:298
		{
			yyVAL.expr = CNode("+", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 77:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:299
		{
			yyVAL.expr = CNode("-", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 78:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:300
		{
			yyVAL.expr = CNode("*", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 79:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:301
		{
			yyVAL.expr = CNode("/", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 80:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:302
		{
			yyVAL.expr = CNode("%", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 81:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:303
		{
			yyVAL.expr = CNode("^", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 82:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:304
		{
			yyVAL.expr = CNode("<<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 83:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:305
		{
			yyVAL.expr = CNode(">>", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 84:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:306
		{
			yyVAL.expr = CNode(">>>", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:307
		{
			yyVAL.expr = CNode("|", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 86:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:308
		{
			yyVAL.expr = CNode("&", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 87:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:309
		{
			yyVAL.expr = CNode("-", NNode(0.0), yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 88:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:310
		{
			yyVAL.expr = CNode("~", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 89:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:311
		{
			yyVAL.expr = CNode("!", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 90:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:312
		{
			yyVAL.expr = CNode("#", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 91:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:315
		{
			yyVAL.expr = SNode(yyDollar[1].token.Str).SetPos(yyDollar[1].token)
		}
	case 92:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:318
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 93:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:319
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 94:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:320
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 95:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:321
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 96:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:324
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
	case 97:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:375
		{
			yyVAL.expr = CNode()
		}
	case 98:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:376
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 99:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:379
		{
			yyVAL.expr = CNode(yyDollar[1].str, "<a>", yyDollar[2].expr, yyDollar[3].expr).setPos0(yyDollar[2].expr)
		}
	case 100:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:380
		{
			yyVAL.expr = CNode(yyDollar[1].str, "<a>", yyDollar[2].expr, CNode("chain", CNode("ret", yyDollar[4].expr))).setPos0(yyDollar[2].expr)
		}
	case 101:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:381
		{
			yyVAL.expr = CNode(yyDollar[1].str, "<a>", CNode(), CNode("chain", CNode("ret", yyDollar[3].expr))).setPos0(yyDollar[3].expr)
		}
	case 102:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:384
		{
			yyVAL.expr = CNode()
		}
	case 103:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:385
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 104:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:388
		{
			yyVAL.expr = CNode("map", CNode()).setPos0(yyDollar[1].token)
		}
	case 105:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:389
		{
			yyVAL.expr = yyDollar[2].expr.setPos0(yyDollar[1].token)
		}
	case 106:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:392
		{
			yyVAL.expr = CNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 107:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:393
		{
			yyVAL.expr = CNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 108:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:394
		{
			yyVAL.expr = CNode("array", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 109:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:395
		{
			yyVAL.expr = CNode("array", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	}
	goto yystack /* stack new state and value */
}
