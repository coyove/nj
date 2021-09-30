// Code generated by goyacc -o parser.go parser.go.y. DO NOT EDIT.

//line parser.go.y:2
package parser

import __yyfmt__ "fmt"

//line parser.go.y:2

//line parser.go.y:25
type yySymType struct {
	yys   int
	token Token
	expr  Node
}

const TDo = 57346
const TLocal = 57347
const TElseIf = 57348
const TThen = 57349
const TEnd = 57350
const TBreak = 57351
const TElse = 57352
const TFor = 57353
const TWhile = 57354
const TFunc = 57355
const TIf = 57356
const TReturn = 57357
const TReturnVoid = 57358
const TRepeat = 57359
const TUntil = 57360
const TNot = 57361
const TLabel = 57362
const TGoto = 57363
const TIn = 57364
const TNext = 57365
const TLsh = 57366
const TRsh = 57367
const TURsh = 57368
const TDotDotDot = 57369
const TOr = 57370
const TAnd = 57371
const TEqeq = 57372
const TNeq = 57373
const TLte = 57374
const TGte = 57375
const TIdent = 57376
const TNumber = 57377
const TString = 57378
const TIDiv = 57379
const ASSIGN = 57380
const FUNC = 57381
const UNARY = 57382

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"TDo",
	"TLocal",
	"TElseIf",
	"TThen",
	"TEnd",
	"TBreak",
	"TElse",
	"TFor",
	"TWhile",
	"TFunc",
	"TIf",
	"TReturn",
	"TReturnVoid",
	"TRepeat",
	"TUntil",
	"TNot",
	"TLabel",
	"TGoto",
	"TIn",
	"TNext",
	"TLsh",
	"TRsh",
	"TURsh",
	"TDotDotDot",
	"TOr",
	"TAnd",
	"TEqeq",
	"TNeq",
	"TLte",
	"TGte",
	"TIdent",
	"TNumber",
	"TString",
	"'{'",
	"'['",
	"'('",
	"'='",
	"'>'",
	"'<'",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'%'",
	"'^'",
	"'#'",
	"'.'",
	"'&'",
	"'|'",
	"'~'",
	"TIDiv",
	"'T'",
	"ASSIGN",
	"FUNC",
	"UNARY",
	"';'",
	"','",
	"')'",
	"']'",
	"'}'",
}

var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.go.y:350

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 24,
	40, 46,
	60, 46,
	-2, 80,
	-1, 99,
	40, 47,
	60, 47,
	-2, 80,
}

const yyPrivate = 57344

const yyLast = 1162

var yyAct = [...]int{
	30, 31, 37, 16, 132, 170, 29, 154, 153, 205,
	179, 45, 195, 46, 164, 176, 46, 25, 33, 34,
	35, 50, 32, 159, 53, 140, 62, 38, 16, 96,
	161, 135, 47, 83, 133, 201, 36, 89, 90, 91,
	180, 92, 85, 143, 94, 97, 194, 97, 163, 97,
	100, 95, 48, 16, 98, 175, 173, 136, 133, 165,
	39, 149, 24, 110, 111, 112, 113, 114, 115, 116,
	117, 118, 119, 120, 121, 122, 123, 124, 125, 126,
	127, 128, 129, 130, 107, 145, 155, 24, 101, 137,
	156, 134, 73, 74, 76, 142, 26, 25, 93, 102,
	139, 75, 105, 141, 52, 147, 148, 49, 150, 99,
	46, 16, 24, 144, 28, 27, 59, 17, 51, 168,
	219, 9, 106, 22, 20, 61, 23, 13, 12, 21,
	185, 6, 11, 10, 110, 15, 14, 157, 71, 72,
	73, 74, 76, 43, 86, 160, 25, 18, 16, 75,
	54, 2, 42, 16, 40, 44, 109, 174, 57, 56,
	4, 3, 1, 58, 16, 5, 41, 37, 182, 183,
	24, 60, 0, 187, 188, 0, 190, 181, 0, 0,
	16, 0, 25, 33, 34, 35, 16, 32, 16, 0,
	0, 0, 38, 0, 16, 16, 0, 0, 207, 0,
	146, 36, 210, 0, 0, 151, 16, 24, 16, 0,
	16, 16, 24, 0, 16, 216, 0, 0, 0, 0,
	16, 0, 0, 24, 0, 59, 17, 172, 0, 0,
	9, 171, 22, 20, 0, 23, 13, 12, 21, 24,
	162, 11, 10, 0, 0, 24, 0, 24, 0, 0,
	0, 0, 0, 24, 24, 25, 0, 0, 0, 0,
	178, 0, 80, 81, 82, 24, 184, 24, 186, 24,
	24, 0, 0, 24, 0, 0, 192, 193, 0, 24,
	60, 71, 72, 73, 74, 76, 79, 0, 0, 77,
	78, 204, 75, 206, 0, 208, 0, 209, 0, 80,
	81, 82, 212, 63, 64, 69, 70, 68, 67, 0,
	0, 218, 0, 0, 0, 0, 65, 66, 71, 72,
	73, 74, 76, 79, 0, 0, 77, 78, 0, 75,
	0, 0, 0, 0, 80, 81, 82, 189, 63, 64,
	69, 70, 68, 67, 0, 0, 0, 0, 0, 0,
	0, 65, 66, 71, 72, 73, 74, 76, 79, 0,
	0, 77, 78, 0, 75, 0, 0, 0, 0, 80,
	81, 82, 158, 63, 64, 69, 70, 68, 67, 0,
	0, 0, 0, 0, 196, 0, 65, 66, 71, 72,
	73, 74, 76, 79, 0, 0, 77, 78, 0, 75,
	0, 0, 0, 0, 80, 81, 82, 138, 63, 64,
	69, 70, 68, 67, 0, 0, 0, 0, 0, 0,
	0, 65, 66, 71, 72, 73, 74, 76, 79, 0,
	0, 77, 78, 0, 75, 0, 0, 80, 81, 82,
	197, 63, 64, 69, 70, 68, 67, 0, 0, 0,
	0, 0, 0, 0, 65, 66, 71, 72, 73, 74,
	76, 79, 0, 0, 77, 78, 0, 75, 0, 0,
	0, 80, 81, 82, 131, 63, 64, 69, 70, 68,
	67, 0, 0, 0, 214, 0, 0, 0, 65, 66,
	71, 72, 73, 74, 76, 79, 0, 0, 77, 78,
	0, 75, 0, 0, 80, 81, 82, 167, 63, 64,
	69, 70, 68, 67, 0, 0, 0, 0, 200, 0,
	0, 65, 66, 71, 72, 73, 74, 76, 79, 0,
	0, 77, 78, 0, 75, 80, 81, 82, 0, 63,
	64, 69, 70, 68, 67, 0, 198, 0, 0, 0,
	0, 0, 65, 66, 71, 72, 73, 74, 76, 79,
	0, 0, 77, 78, 0, 75, 80, 81, 82, 0,
	63, 64, 69, 70, 68, 67, 0, 169, 0, 0,
	0, 0, 0, 65, 66, 71, 72, 73, 74, 76,
	79, 0, 0, 77, 78, 0, 75, 80, 81, 82,
	0, 63, 64, 69, 70, 68, 67, 0, 0, 0,
	0, 108, 0, 0, 65, 66, 71, 72, 73, 74,
	76, 79, 0, 0, 77, 78, 0, 75, 80, 81,
	82, 0, 63, 64, 69, 70, 68, 67, 0, 103,
	0, 0, 0, 0, 0, 65, 66, 71, 72, 73,
	74, 76, 79, 0, 0, 77, 78, 0, 75, 80,
	81, 82, 0, 63, 64, 69, 70, 68, 67, 0,
	0, 0, 0, 0, 0, 0, 65, 66, 71, 72,
	73, 74, 76, 79, 0, 0, 77, 78, 0, 75,
	80, 81, 82, 0, 63, 64, 69, 70, 68, 67,
	0, 0, 0, 0, 0, 0, 0, 65, 66, 71,
	72, 73, 74, 76, 79, 0, 0, 77, 78, 0,
	75, 59, 17, 0, 0, 217, 9, 0, 22, 20,
	0, 23, 13, 12, 21, 0, 0, 11, 10, 0,
	0, 0, 59, 17, 0, 0, 215, 9, 0, 22,
	20, 25, 23, 13, 12, 21, 0, 0, 11, 10,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	59, 17, 25, 0, 213, 9, 60, 22, 20, 0,
	23, 13, 12, 21, 0, 0, 11, 10, 0, 0,
	0, 59, 17, 0, 0, 211, 9, 60, 22, 20,
	25, 23, 13, 12, 21, 0, 0, 11, 10, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 59,
	17, 25, 0, 203, 9, 60, 22, 20, 0, 23,
	13, 12, 21, 0, 0, 11, 10, 0, 0, 0,
	59, 17, 0, 0, 202, 9, 60, 22, 20, 25,
	23, 13, 12, 21, 0, 0, 11, 10, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 59, 17,
	25, 0, 199, 9, 60, 22, 20, 0, 23, 13,
	12, 21, 0, 0, 11, 10, 0, 0, 0, 59,
	17, 0, 0, 191, 9, 60, 22, 20, 25, 23,
	13, 12, 21, 0, 0, 11, 10, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 59, 17, 25,
	0, 177, 9, 60, 22, 20, 0, 23, 13, 12,
	21, 0, 0, 11, 10, 0, 0, 0, 59, 17,
	0, 0, 166, 9, 60, 22, 20, 25, 23, 13,
	12, 21, 0, 0, 11, 10, 80, 81, 82, 0,
	0, 64, 69, 70, 68, 67, 0, 0, 25, 0,
	0, 0, 60, 65, 66, 71, 72, 73, 74, 76,
	79, 0, 0, 77, 78, 0, 75, 59, 17, 0,
	0, 152, 9, 60, 22, 20, 0, 23, 13, 12,
	21, 0, 0, 11, 10, 0, 0, 0, 59, 17,
	0, 37, 0, 9, 0, 22, 20, 25, 23, 13,
	12, 21, 104, 0, 11, 10, 87, 33, 34, 35,
	88, 32, 0, 0, 0, 0, 38, 0, 25, 0,
	0, 0, 60, 0, 0, 36, 0, 0, 0, 59,
	17, 0, 0, 55, 9, 84, 22, 20, 0, 23,
	13, 12, 21, 60, 0, 11, 10, 0, 0, 0,
	0, 0, 0, 7, 17, 0, 0, 0, 9, 25,
	22, 20, 19, 23, 13, 12, 21, 0, 0, 11,
	10, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	59, 17, 0, 25, 60, 9, 0, 22, 20, 0,
	23, 13, 12, 21, 0, 0, 11, 10, 80, 81,
	82, 0, 0, 0, 69, 70, 68, 67, 8, 0,
	25, 80, 81, 82, 0, 65, 66, 71, 72, 73,
	74, 76, 79, 0, 0, 77, 78, 0, 75, 0,
	71, 72, 73, 74, 76, 60, 0, 0, 0, 0,
	0, 75,
}

var yyPact = [...]int{
	-1000, 1069, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	81, 80, -1000, 148, -1000, -1000, 116, 76, -8, 73,
	148, -1000, 70, 148, -1000, -1000, 1045, -1000, 105, -34,
	666, 116, 148, -1000, -1000, 992, 148, 148, 148, -1000,
	148, 64, -1000, -1000, -17, -11, -1000, 148, 63, 49,
	635, 1004, 62, 604, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, 148, 148, 148, 148, 148, 148, 148, 148,
	148, 148, 148, 148, 148, 148, 148, 148, 148, 148,
	148, 148, 148, 413, -1000, -26, -29, 17, 148, -1000,
	-1000, -1000, 345, -1000, -1000, -2, 148, 61, -34, -1000,
	116, -18, 51, -1000, 148, 148, 27, 148, -1000, 983,
	666, 932, 1094, 238, 238, 238, 238, 238, 238, 47,
	47, -1000, -1000, -1000, -1000, 1107, 1107, 1107, 95, 95,
	95, -1000, -55, 148, -56, 52, 148, 310, -1000, -38,
	-30, -34, -1000, -1000, -13, 20, 934, 666, 447, 97,
	573, 221, -1000, -1000, -1000, 16, 148, 666, 15, -1000,
	-46, -1000, 913, -1000, -51, -21, -1000, 148, 148, -1000,
	122, -1000, 148, 148, 275, 148, -1000, -1000, 885, -1000,
	-1000, -15, 380, 542, 864, -1000, 1096, 511, 666, -5,
	666, -1000, 836, 815, -1000, -52, -1000, 148, -1000, -1000,
	-1000, 148, -1000, -1000, 787, -1000, 766, 480, 738, 221,
	666, -1000, 717, -1000, -1000, -1000, -1000, -1000, 112, -1000,
}

var yyPgo = [...]int{
	0, 162, 96, 151, 150, 60, 147, 11, 0, 6,
	144, 1, 143, 163, 136, 135, 5, 159, 158, 131,
	4,
}

var yyR1 = [...]int{
	0, 1, 1, 2, 2, 3, 3, 3, 3, 3,
	3, 4, 4, 4, 4, 4, 18, 18, 13, 13,
	13, 13, 14, 14, 14, 14, 14, 14, 15, 16,
	16, 16, 19, 19, 19, 19, 19, 19, 17, 17,
	17, 17, 17, 5, 5, 5, 6, 6, 7, 7,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	11, 11, 11, 12, 12, 12, 9, 9, 10, 10,
	10, 10, 20, 20,
}

var yyR2 = [...]int{
	0, 0, 2, 0, 2, 1, 1, 1, 1, 3,
	1, 1, 1, 1, 3, 1, 1, 1, 1, 2,
	4, 3, 5, 4, 9, 11, 9, 7, 6, 0,
	2, 5, 6, 7, 8, 8, 9, 10, 1, 2,
	3, 1, 2, 1, 4, 3, 1, 3, 1, 3,
	1, 3, 1, 1, 2, 4, 4, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 2, 2, 2,
	1, 2, 2, 2, 4, 5, 1, 3, 3, 5,
	5, 7, 0, 1,
}

var yyChk = [...]int{
	-1000, -1, -3, -17, -18, -13, -19, 4, 59, 9,
	21, 20, 16, 15, -14, -15, -11, 5, -6, 13,
	12, 17, 11, 14, -5, 34, -2, 34, 34, -9,
	-8, -11, 39, 35, 36, 37, 53, 19, 44, -5,
	38, 50, 36, -12, 39, -7, 34, 40, 60, 34,
	-8, -2, 34, -8, -4, 8, -17, -18, -13, 4,
	59, 20, 60, 28, 29, 41, 42, 33, 32, 30,
	31, 43, 44, 45, 46, 54, 47, 51, 52, 48,
	24, 25, 26, -8, 63, -9, -10, 34, 38, -8,
	-8, -8, -8, 34, 61, -9, 40, 60, -9, -5,
	-11, 39, 50, 4, 18, 40, 60, 22, 7, -2,
	-8, -8, -8, -8, -8, -8, -8, -8, -8, -8,
	-8, -8, -8, -8, -8, -8, -8, -8, -8, -8,
	-8, 61, -20, 60, -20, 60, 40, -8, 62, -20,
	27, -9, 34, 61, -7, 34, -2, -8, -8, 34,
	-8, -2, 8, 63, 63, 34, 38, -8, 62, 61,
	-20, 60, -2, 61, 27, 39, 8, 60, 22, 4,
	-16, 10, 6, 40, -8, 40, 61, 8, -2, 61,
	61, -7, -8, -8, -2, 8, -2, -8, -8, 62,
	-8, 8, -2, -2, 61, 27, 4, 60, 4, 8,
	7, 40, 8, 8, -2, 61, -2, -8, -2, -2,
	-8, 8, -2, 8, 4, 8, -16, 8, -2, 8,
}

var yyDef = [...]int{
	1, -2, 2, 5, 6, 7, 8, 3, 10, 38,
	0, 0, 41, 0, 16, 17, 18, 0, 0, 0,
	0, 3, 0, 0, -2, 43, 0, 39, 0, 42,
	86, 50, 0, 52, 53, 0, 0, 0, 0, 80,
	0, 0, 81, 82, 0, 19, 48, 0, 0, 0,
	0, 0, 0, 0, 4, 9, 11, 12, 13, 3,
	15, 40, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 54, 92, 92, 43, 0, 77,
	78, 79, 0, 45, 83, 92, 0, 0, 21, -2,
	0, 0, 0, 3, 0, 0, 0, 0, 3, 0,
	87, 57, 58, 59, 60, 61, 62, 63, 64, 65,
	66, 67, 68, 69, 70, 71, 72, 73, 74, 75,
	76, 51, 0, 93, 0, 93, 0, 0, 44, 0,
	92, 20, 49, 3, 0, 0, 0, 23, 0, 0,
	0, 29, 14, 55, 56, 0, 0, 88, 0, 84,
	0, 93, 0, 3, 0, 0, 22, 0, 0, 3,
	0, 3, 0, 0, 0, 0, 85, 32, 0, 3,
	3, 0, 0, 0, 0, 28, 30, 0, 90, 0,
	89, 33, 0, 0, 3, 0, 3, 0, 3, 27,
	3, 0, 34, 35, 0, 3, 0, 0, 0, 29,
	91, 36, 0, 24, 3, 26, 31, 37, 0, 25,
}

var yyTok1 = [...]int{
	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 49, 3, 47, 51, 3,
	39, 61, 45, 43, 60, 44, 50, 46, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 59,
	42, 40, 41, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 55, 3, 3, 3, 3, 3,
	3, 38, 3, 62, 48, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 37, 52, 63, 53,
}

var yyTok2 = [...]int{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 54, 56, 57, 58,
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
//line parser.go.y:54
		{
			yyVAL.expr = __chain()
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:60
		{
			yyVAL.expr = yyDollar[1].expr.append(yyDollar[2].expr)
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 3:
		yyDollar = yyS[yypt-0 : yypt+1]
//line parser.go.y:68
		{
			yyVAL.expr = __chain()
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:71
		{
			yyVAL.expr = yyDollar[1].expr.append(yyDollar[2].expr)
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:76
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:77
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:78
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:79
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 9:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:80
		{
			yyVAL.expr = __do(yyDollar[2].expr)
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:81
		{
			yyVAL.expr = emptyNode
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:84
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:85
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:86
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 14:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:87
		{
			yyVAL.expr = __do(yyDollar[2].expr)
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:88
		{
			yyVAL.expr = emptyNode
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:91
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:92
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:95
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:98
		{
			yyVAL.expr = __chain()
			for _, v := range yyDollar[2].expr.Nodes {
				yyVAL.expr = yyVAL.expr.append(__set(v, NewSymbol(ANil)).SetPos(yyDollar[1].token.Pos))
			}
		}
	case 20:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:104
		{
			if len(yyDollar[4].expr.Nodes) == 1 && len(yyDollar[2].expr.Nodes) > 1 {
				tmp := randomVarname()
				yyVAL.expr = __chain(__local([]Node{tmp}, yyDollar[4].expr.Nodes, yyDollar[1].token.Pos))
				for i, ident := range yyDollar[2].expr.Nodes {
					yyVAL.expr = yyVAL.expr.append(__local([]Node{ident}, []Node{__load(tmp, NewNumberFromInt(int64(i))).SetPos(yyDollar[1].token.Pos)}, yyDollar[1].token.Pos))
				}
			} else {
				yyVAL.expr = __local(yyDollar[2].expr.Nodes, yyDollar[4].expr.Nodes, yyDollar[1].token.Pos)
			}
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:115
		{
			if len(yyDollar[3].expr.Nodes) == 1 && len(yyDollar[1].expr.Nodes) > 1 {
				tmp := randomVarname()
				yyVAL.expr = __chain(__local([]Node{tmp}, yyDollar[3].expr.Nodes, yyDollar[2].token.Pos))
				for i, decl := range yyDollar[1].expr.Nodes {
					x := decl.moveLoadStore(__move, __load(tmp, NewNumberFromInt(int64(i))).SetPos(yyDollar[2].token.Pos)).SetPos(yyDollar[2].token.Pos)
					yyVAL.expr = yyVAL.expr.append(x)
				}
			} else {
				yyVAL.expr = __moveMulti(yyDollar[1].expr.Nodes, yyDollar[3].expr.Nodes, yyDollar[2].token.Pos)
			}
		}
	case 22:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:129
		{
			yyVAL.expr = __loop(__if(yyDollar[2].expr, yyDollar[4].expr, breakNode).SetPos(yyDollar[1].token.Pos)).SetPos(yyDollar[1].token.Pos)
		}
	case 23:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:132
		{
			yyVAL.expr = __loop(
				__chain(
					yyDollar[2].expr,
					__if(yyDollar[4].expr, breakNode, emptyNode).SetPos(yyDollar[1].token.Pos),
				).SetPos(yyDollar[1].token.Pos),
			).SetPos(yyDollar[1].token.Pos)
		}
	case 24:
		yyDollar = yyS[yypt-9 : yypt+1]
//line parser.go.y:140
		{
			forVar, forEnd := NewSymbolFromToken(yyDollar[2].token), randomVarname()
			yyVAL.expr = __do(
				__set(forVar, yyDollar[4].expr).SetPos(yyDollar[1].token.Pos),
				__set(forEnd, yyDollar[6].expr).SetPos(yyDollar[1].token.Pos),
				__loop(
					__if(
						__less(forVar, forEnd),
						__chain(yyDollar[8].expr, __inc(forVar, oneNode).SetPos(yyDollar[1].token.Pos)),
						breakNode,
					).SetPos(yyDollar[1].token.Pos),
				).SetPos(yyDollar[1].token.Pos),
			)
		}
	case 25:
		yyDollar = yyS[yypt-11 : yypt+1]
//line parser.go.y:154
		{
			forVar, forEnd, forStep := NewSymbolFromToken(yyDollar[2].token), randomVarname(), randomVarname()
			body := __chain(yyDollar[10].expr, __inc(forVar, forStep))
			yyVAL.expr = __do(
				__set(forVar, yyDollar[4].expr).SetPos(yyDollar[1].token.Pos),
				__set(forEnd, yyDollar[6].expr).SetPos(yyDollar[1].token.Pos),
				__set(forStep, yyDollar[8].expr).SetPos(yyDollar[1].token.Pos))

			if yyDollar[8].expr.IsNumber() { // step is a static number, easy case
				if yyDollar[8].expr.IsNegativeNumber() {
					yyVAL.expr = yyVAL.expr.append(__loop(__if(__less(forEnd, forVar), body, breakNode).SetPos(yyDollar[1].token.Pos)).SetPos(yyDollar[1].token.Pos))
				} else {
					yyVAL.expr = yyVAL.expr.append(__loop(__if(__less(forVar, forEnd), body, breakNode).SetPos(yyDollar[1].token.Pos)).SetPos(yyDollar[1].token.Pos))
				}
			} else {
				yyVAL.expr = yyVAL.expr.append(__loop(
					__if(
						__less(zeroNode, forStep).SetPos(yyDollar[1].token.Pos),
						// +step
						__if(__lessEq(forEnd, forVar), breakNode, body).SetPos(yyDollar[1].token.Pos),
						// -step
						__if(__lessEq(forVar, forEnd), breakNode, body).SetPos(yyDollar[1].token.Pos),
					).SetPos(yyDollar[1].token.Pos),
				).SetPos(yyDollar[1].token.Pos))
			}
		}
	case 26:
		yyDollar = yyS[yypt-9 : yypt+1]
//line parser.go.y:180
		{
			yyVAL.expr = __forIn(yyDollar[2].token, yyDollar[4].token, yyDollar[6].expr, yyDollar[8].expr, yyDollar[1].token.Pos)
		}
	case 27:
		yyDollar = yyS[yypt-7 : yypt+1]
//line parser.go.y:183
		{
			yyVAL.expr = __forIn(yyDollar[2].token, yyDollar[1].token, yyDollar[4].expr, yyDollar[6].expr, yyDollar[1].token.Pos)
		}
	case 28:
		yyDollar = yyS[yypt-6 : yypt+1]
//line parser.go.y:189
		{
			yyVAL.expr = __if(yyDollar[2].expr, yyDollar[4].expr, yyDollar[5].expr).SetPos(yyDollar[1].token.Pos)
		}
	case 29:
		yyDollar = yyS[yypt-0 : yypt+1]
//line parser.go.y:194
		{
			yyVAL.expr = NewComplex()
		}
	case 30:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:197
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 31:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:200
		{
			yyVAL.expr = __if(yyDollar[2].expr, yyDollar[4].expr, yyDollar[5].expr).SetPos(yyDollar[1].token.Pos)
		}
	case 32:
		yyDollar = yyS[yypt-6 : yypt+1]
//line parser.go.y:205
		{
			yyVAL.expr = __func(yyDollar[2].token, emptyNode, "", yyDollar[5].expr)
		}
	case 33:
		yyDollar = yyS[yypt-7 : yypt+1]
//line parser.go.y:206
		{
			yyVAL.expr = __func(yyDollar[2].token, yyDollar[4].expr, "", yyDollar[6].expr)
		}
	case 34:
		yyDollar = yyS[yypt-8 : yypt+1]
//line parser.go.y:207
		{
			yyVAL.expr = __func(yyDollar[2].token, __dotdotdot(yyDollar[4].expr), "", yyDollar[7].expr)
		}
	case 35:
		yyDollar = yyS[yypt-8 : yypt+1]
//line parser.go.y:208
		{
			yyVAL.expr = __store(NewSymbolFromToken(yyDollar[2].token), NewString(yyDollar[4].token.Str), __func(__markupFuncName(yyDollar[2].token, yyDollar[4].token), emptyNode, "", yyDollar[7].expr))
		}
	case 36:
		yyDollar = yyS[yypt-9 : yypt+1]
//line parser.go.y:211
		{
			yyVAL.expr = __store(NewSymbolFromToken(yyDollar[2].token), NewString(yyDollar[4].token.Str), __func(__markupFuncName(yyDollar[2].token, yyDollar[4].token), yyDollar[6].expr, "", yyDollar[8].expr))
		}
	case 37:
		yyDollar = yyS[yypt-10 : yypt+1]
//line parser.go.y:214
		{
			yyVAL.expr = __store(NewSymbolFromToken(yyDollar[2].token), NewString(yyDollar[4].token.Str), __func(__markupFuncName(yyDollar[2].token, yyDollar[4].token), __dotdotdot(yyDollar[6].expr), "", yyDollar[9].expr))
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:219
		{
			yyVAL.expr = NewComplex(NewSymbol(ABreak)).SetPos(yyDollar[1].token.Pos)
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:222
		{
			yyVAL.expr = NewComplex(NewSymbol(AGoto), NewSymbolFromToken(yyDollar[2].token)).SetPos(yyDollar[1].token.Pos)
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:225
		{
			yyVAL.expr = NewComplex(NewSymbol(ALabel), NewSymbolFromToken(yyDollar[2].token))
		}
	case 41:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:228
		{
			yyVAL.expr = NewComplex(NewSymbol(AReturn), NewSymbol(ANil)).SetPos(yyDollar[1].token.Pos)
		}
	case 42:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:231
		{
			if len(yyDollar[2].expr.Nodes) == 1 {
				a := yyDollar[2].expr.Nodes[0]
				if len(a.Nodes) == 3 && a.Nodes[0].SymbolValue() == ACall {
					// return call(...) -> return tailcall(...)
					a.Nodes[0].strSym = ATailCall
				}
				yyVAL.expr = NewComplex(NewSymbol(AReturn), a).SetPos(yyDollar[1].token.Pos)
			} else {
				yyVAL.expr = NewComplex(NewSymbol(AReturn), NewComplex(NewSymbol(AArray), yyDollar[2].expr)).SetPos(yyDollar[1].token.Pos)
			}
		}
	case 43:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:245
		{
			yyVAL.expr = NewSymbolFromToken(yyDollar[1].token)
		}
	case 44:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:248
		{
			yyVAL.expr = __load(yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:251
		{
			yyVAL.expr = __load(yyDollar[1].expr, NewString(yyDollar[3].token.Str)).SetPos(yyDollar[2].token.Pos)
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:256
		{
			yyVAL.expr = NewComplex(yyDollar[1].expr)
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:259
		{
			yyVAL.expr = yyDollar[1].expr.append(yyDollar[3].expr)
		}
	case 48:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:264
		{
			yyVAL.expr = NewComplex(NewSymbolFromToken(yyDollar[1].token))
		}
	case 49:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:267
		{
			yyVAL.expr = yyDollar[1].expr.append(NewSymbolFromToken(yyDollar[3].token))
		}
	case 50:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:272
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:273
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 52:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:274
		{
			yyVAL.expr = NewNumberFromString(yyDollar[1].token.Str)
		}
	case 53:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:275
		{
			yyVAL.expr = NewString(yyDollar[1].token.Str)
		}
	case 54:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:276
		{
			yyVAL.expr = NewComplex(NewSymbol(AArrayMap), emptyNode).SetPos(yyDollar[1].token.Pos)
		}
	case 55:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:277
		{
			yyVAL.expr = NewComplex(NewSymbol(AArray), yyDollar[2].expr).SetPos(yyDollar[1].token.Pos)
		}
	case 56:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:278
		{
			yyVAL.expr = NewComplex(NewSymbol(AArrayMap), yyDollar[2].expr).SetPos(yyDollar[1].token.Pos)
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:279
		{
			yyVAL.expr = NewComplex(NewSymbol(AOr), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:280
		{
			yyVAL.expr = NewComplex(NewSymbol(AAnd), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:281
		{
			yyVAL.expr = NewComplex(NewSymbol(ALess), yyDollar[3].expr, yyDollar[1].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:282
		{
			yyVAL.expr = NewComplex(NewSymbol(ALess), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:283
		{
			yyVAL.expr = NewComplex(NewSymbol(ALessEq), yyDollar[3].expr, yyDollar[1].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 62:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:284
		{
			yyVAL.expr = NewComplex(NewSymbol(ALessEq), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 63:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:285
		{
			yyVAL.expr = NewComplex(NewSymbol(AEq), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:286
		{
			yyVAL.expr = NewComplex(NewSymbol(ANeq), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 65:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:287
		{
			yyVAL.expr = NewComplex(NewSymbol(AAdd), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:288
		{
			yyVAL.expr = NewComplex(NewSymbol(ASub), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:289
		{
			yyVAL.expr = NewComplex(NewSymbol(AMul), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:290
		{
			yyVAL.expr = NewComplex(NewSymbol(ADiv), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:291
		{
			yyVAL.expr = NewComplex(NewSymbol(AIDiv), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 70:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:292
		{
			yyVAL.expr = NewComplex(NewSymbol(AMod), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:293
		{
			yyVAL.expr = NewComplex(NewSymbol(ABitAnd), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 72:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:294
		{
			yyVAL.expr = NewComplex(NewSymbol(ABitOr), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:295
		{
			yyVAL.expr = NewComplex(NewSymbol(ABitXor), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:296
		{
			yyVAL.expr = NewComplex(NewSymbol(ABitLsh), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 75:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:297
		{
			yyVAL.expr = NewComplex(NewSymbol(ABitRsh), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:298
		{
			yyVAL.expr = NewComplex(NewSymbol(ABitURsh), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 77:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:299
		{
			yyVAL.expr = NewComplex(NewSymbol(ABitNot), yyDollar[2].expr).SetPos(yyDollar[1].token.Pos)
		}
	case 78:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:300
		{
			yyVAL.expr = NewComplex(NewSymbol(ANot), yyDollar[2].expr).SetPos(yyDollar[1].token.Pos)
		}
	case 79:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:301
		{
			yyVAL.expr = NewComplex(NewSymbol(ASub), zeroNode, yyDollar[2].expr).SetPos(yyDollar[1].token.Pos)
		}
	case 80:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:304
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 81:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:307
		{
			yyVAL.expr = __call(yyDollar[1].expr, NewComplex(NewString(yyDollar[2].token.Str))).SetPos(yyDollar[1].expr.Pos())
		}
	case 82:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:310
		{
			yyDollar[2].expr.Nodes[1] = yyDollar[1].expr
			yyVAL.expr = yyDollar[2].expr
		}
	case 83:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:316
		{
			yyVAL.expr = __call(emptyNode, emptyNode).SetPos(yyDollar[1].token.Pos)
		}
	case 84:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:319
		{
			yyVAL.expr = __call(emptyNode, yyDollar[2].expr).SetPos(yyDollar[1].token.Pos)
		}
	case 85:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:322
		{
			yyVAL.expr = __call(emptyNode, __dotdotdot(yyDollar[2].expr)).SetPos(yyDollar[1].token.Pos)
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:327
		{
			yyVAL.expr = NewComplex(yyDollar[1].expr)
		}
	case 87:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:330
		{
			yyVAL.expr = yyDollar[1].expr.append(yyDollar[3].expr)
		}
	case 88:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:335
		{
			yyVAL.expr = NewComplex(NewString(yyDollar[1].token.Str), yyDollar[3].expr)
		}
	case 89:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:338
		{
			yyVAL.expr = NewComplex(yyDollar[2].expr, yyDollar[5].expr)
		}
	case 90:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:341
		{
			yyVAL.expr = yyDollar[1].expr.append(NewString(yyDollar[3].token.Str)).append(yyDollar[5].expr)
		}
	case 91:
		yyDollar = yyS[yypt-7 : yypt+1]
//line parser.go.y:344
		{
			yyVAL.expr = yyDollar[1].expr.append(yyDollar[4].expr).append(yyDollar[7].expr)
		}
	case 92:
		yyDollar = yyS[yypt-0 : yypt+1]
//line parser.go.y:348
		{
			yyVAL.expr = emptyNode
		}
	case 93:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:348
		{
			yyVAL.expr = emptyNode
		}
	}
	goto yystack /* stack new state and value */
}
