// Code generated by goyacc -o parser.go parser.go.y. DO NOT EDIT.

//line parser.go.y:2
package parser

import __yyfmt__ "fmt"

//line parser.go.y:2

import "github.com/coyove/nj/typ"

func ss(yylex yyLexer) *Lexer { return yylex.(*Lexer) }

//line parser.go.y:26
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
const TContinue = 57352
const TElse = 57353
const TFor = 57354
const TWhile = 57355
const TFunc = 57356
const TIf = 57357
const TReturn = 57358
const TReturnVoid = 57359
const TRepeat = 57360
const TUntil = 57361
const TNot = 57362
const TLabel = 57363
const TGoto = 57364
const TIn = 57365
const TLsh = 57366
const TRsh = 57367
const TURsh = 57368
const TDotDotDot = 57369
const TLParen = 57370
const TLBracket = 57371
const TIs = 57372
const TOr = 57373
const TAnd = 57374
const TEqeq = 57375
const TNeq = 57376
const TLte = 57377
const TGte = 57378
const TIdent = 57379
const TNumber = 57380
const TString = 57381
const TIDiv = 57382
const TInv = 57383
const TAddEq = 57384
const TSubEq = 57385
const TMulEq = 57386
const TDivEq = 57387
const TIDivEq = 57388
const TModEq = 57389
const TBitAndEq = 57390
const TBitOrEq = 57391
const TBitXorEq = 57392
const TBitLshEq = 57393
const TBitRshEq = 57394
const TBitURshEq = 57395
const ASSIGN = 57396
const FUNC = 57397
const UNARY = 57398

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
	"TContinue",
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
	"TLsh",
	"TRsh",
	"TURsh",
	"TDotDotDot",
	"TLParen",
	"TLBracket",
	"TIs",
	"TOr",
	"TAnd",
	"TEqeq",
	"TNeq",
	"TLte",
	"TGte",
	"TIdent",
	"TNumber",
	"TString",
	"TIDiv",
	"TInv",
	"TAddEq",
	"TSubEq",
	"TMulEq",
	"TDivEq",
	"TIDivEq",
	"TModEq",
	"TBitAndEq",
	"TBitOrEq",
	"TBitXorEq",
	"TBitLshEq",
	"TBitRshEq",
	"TBitURshEq",
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
	"':'",
	"')'",
	"','",
	"'T'",
	"ASSIGN",
	"FUNC",
	"UNARY",
	"';'",
	"']'",
	"'}'",
}

var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.go.y:214

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 20,
	57, 100,
	73, 100,
	-2, 85,
	-1, 141,
	57, 101,
	73, 101,
	-2, 85,
}

const yyPrivate = 57344

const yyLast = 1885

var yyAct = [...]int{
	17, 168, 217, 69, 38, 200, 199, 231, 139, 230,
	139, 212, 139, 209, 139, 226, 138, 44, 115, 171,
	169, 71, 85, 240, 181, 88, 224, 25, 93, 94,
	95, 96, 139, 227, 97, 44, 103, 72, 43, 172,
	190, 178, 105, 109, 112, 186, 70, 175, 198, 116,
	117, 118, 119, 120, 121, 122, 123, 124, 125, 126,
	127, 128, 129, 130, 131, 132, 133, 134, 135, 180,
	139, 164, 44, 99, 143, 144, 145, 146, 147, 148,
	149, 150, 151, 152, 153, 154, 177, 139, 70, 45,
	161, 162, 20, 44, 169, 70, 136, 39, 159, 2,
	142, 92, 90, 170, 36, 65, 66, 67, 87, 113,
	140, 68, 89, 42, 179, 41, 183, 182, 110, 70,
	68, 60, 86, 210, 39, 40, 20, 39, 70, 215,
	207, 166, 157, 37, 114, 235, 106, 3, 107, 44,
	91, 56, 57, 58, 59, 61, 64, 5, 158, 62,
	63, 8, 40, 7, 111, 40, 6, 188, 189, 101,
	191, 19, 141, 108, 196, 184, 1, 0, 197, 0,
	183, 0, 202, 203, 204, 0, 20, 185, 0, 0,
	206, 0, 208, 0, 0, 211, 0, 0, 0, 0,
	0, 0, 0, 0, 220, 0, 221, 20, 0, 0,
	225, 0, 0, 0, 0, 174, 0, 47, 46, 0,
	0, 0, 0, 30, 0, 232, 233, 0, 0, 0,
	237, 0, 0, 0, 0, 0, 0, 0, 241, 242,
	102, 26, 32, 0, 29, 0, 65, 66, 67, 249,
	0, 0, 68, 0, 47, 46, 251, 35, 34, 33,
	137, 0, 60, 0, 0, 187, 258, 0, 0, 28,
	192, 0, 0, 27, 20, 0, 0, 31, 0, 32,
	0, 0, 0, 100, 58, 59, 61, 20, 0, 0,
	62, 0, 20, 0, 35, 34, 33, 0, 0, 0,
	0, 0, 0, 65, 66, 67, 0, 0, 0, 68,
	48, 49, 54, 55, 53, 52, 229, 0, 0, 60,
	0, 0, 0, 0, 0, 0, 234, 0, 236, 20,
	0, 0, 0, 0, 20, 0, 20, 50, 51, 56,
	57, 58, 59, 61, 64, 0, 0, 62, 63, 0,
	20, 254, 20, 20, 250, 0, 252, 0, 253, 20,
	0, 0, 4, 18, 219, 0, 259, 11, 12, 218,
	23, 21, 10, 24, 16, 15, 22, 0, 30, 14,
	13, 73, 74, 75, 76, 77, 78, 79, 80, 81,
	82, 83, 84, 0, 0, 31, 26, 32, 0, 29,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 35, 34, 33, 0, 0, 0, 0, 47,
	46, 4, 18, 0, 28, 260, 11, 12, 27, 23,
	21, 10, 24, 16, 15, 22, 9, 30, 14, 13,
	0, 0, 31, 0, 32, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 31, 26, 32, 0, 29, 35,
	34, 33, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 35, 34, 33, 0, 0, 0, 0, 0, 0,
	4, 18, 0, 28, 257, 11, 12, 27, 23, 21,
	10, 24, 16, 15, 22, 9, 30, 14, 13, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 31, 26, 32, 0, 29, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	35, 34, 33, 0, 0, 0, 0, 0, 0, 4,
	18, 0, 28, 255, 11, 12, 27, 23, 21, 10,
	24, 16, 15, 22, 9, 30, 14, 13, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 31, 26, 32, 0, 29, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 35,
	34, 33, 0, 0, 0, 0, 0, 0, 4, 18,
	0, 28, 247, 11, 12, 27, 23, 21, 10, 24,
	16, 15, 22, 9, 30, 14, 13, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 31, 26, 32, 0, 29, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 35, 34,
	33, 0, 0, 0, 0, 0, 0, 4, 18, 0,
	28, 243, 11, 12, 27, 23, 21, 10, 24, 16,
	15, 22, 9, 30, 14, 13, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	31, 26, 32, 0, 29, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 35, 34, 33,
	0, 0, 0, 0, 0, 0, 4, 18, 0, 28,
	213, 11, 12, 27, 23, 21, 10, 24, 16, 15,
	22, 9, 30, 14, 13, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 31,
	26, 32, 0, 29, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 35, 34, 33, 0,
	0, 0, 0, 0, 0, 4, 18, 0, 28, 205,
	11, 12, 27, 23, 21, 10, 24, 16, 15, 22,
	9, 30, 14, 13, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 31, 26,
	32, 0, 29, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 35, 34, 33, 0, 0,
	0, 0, 0, 0, 4, 18, 0, 28, 176, 11,
	12, 27, 23, 21, 10, 24, 16, 15, 22, 9,
	30, 14, 13, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 31, 26, 32,
	0, 29, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 35, 34, 33, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 28, 0, 0, 0,
	27, 4, 18, 0, 0, 0, 11, 12, 9, 23,
	21, 10, 24, 16, 15, 22, 156, 30, 14, 13,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 31, 26, 32, 0, 29, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 35, 34, 33, 0, 0, 0, 0, 0, 0,
	4, 18, 0, 28, 104, 11, 12, 27, 23, 21,
	10, 24, 16, 15, 22, 9, 30, 14, 13, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 31, 26, 32, 0, 29, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	35, 34, 33, 0, 0, 0, 0, 0, 0, 4,
	18, 0, 28, 0, 11, 12, 27, 23, 21, 10,
	24, 16, 15, 22, 9, 30, 14, 13, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 31, 26, 32, 0, 29, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 65, 66, 67, 35,
	34, 33, 68, 48, 49, 54, 55, 53, 52, 0,
	0, 28, 60, 0, 0, 27, 0, 0, 0, 0,
	0, 0, 0, 9, 0, 0, 0, 0, 0, 0,
	50, 51, 56, 57, 58, 59, 61, 64, 0, 0,
	62, 63, 0, 195, 0, 0, 65, 66, 67, 0,
	0, 194, 68, 48, 49, 54, 55, 53, 52, 65,
	66, 67, 60, 0, 0, 68, 48, 49, 54, 55,
	53, 52, 0, 0, 0, 60, 0, 0, 0, 0,
	50, 51, 56, 57, 58, 59, 61, 64, 0, 244,
	62, 63, 0, 50, 51, 56, 57, 58, 59, 61,
	64, 239, 0, 62, 63, 0, 0, 0, 0, 65,
	66, 67, 0, 0, 223, 68, 48, 49, 54, 55,
	53, 52, 0, 0, 0, 60, 65, 66, 67, 0,
	0, 0, 68, 48, 49, 54, 55, 53, 52, 0,
	0, 0, 60, 50, 51, 56, 57, 58, 59, 61,
	64, 0, 0, 62, 63, 0, 0, 0, 245, 0,
	50, 51, 56, 57, 58, 59, 61, 64, 0, 0,
	62, 63, 65, 66, 67, 238, 0, 0, 68, 48,
	49, 54, 55, 53, 52, 0, 0, 0, 60, 65,
	66, 67, 0, 0, 0, 68, 48, 49, 54, 55,
	53, 52, 0, 0, 0, 60, 50, 51, 56, 57,
	58, 59, 61, 64, 0, 0, 62, 63, 0, 0,
	0, 214, 0, 50, 51, 56, 57, 58, 59, 61,
	64, 0, 0, 62, 63, 65, 66, 67, 193, 0,
	0, 68, 48, 49, 54, 55, 53, 52, 0, 0,
	0, 60, 65, 66, 67, 0, 0, 0, 68, 48,
	49, 54, 55, 53, 52, 0, 0, 0, 60, 50,
	51, 56, 57, 58, 59, 61, 64, 0, 0, 62,
	63, 0, 0, 167, 0, 0, 50, 51, 56, 57,
	58, 59, 61, 64, 0, 0, 62, 63, 0, 228,
	65, 66, 67, 0, 0, 0, 68, 48, 49, 54,
	55, 53, 52, 0, 0, 0, 60, 0, 0, 0,
	0, 0, 0, 0, 0, 47, 46, 0, 0, 0,
	0, 30, 0, 0, 50, 51, 56, 57, 58, 59,
	61, 64, 0, 0, 62, 63, 0, 173, 31, 26,
	32, 0, 29, 0, 0, 47, 46, 0, 0, 0,
	0, 30, 0, 0, 0, 35, 34, 33, 0, 0,
	0, 0, 256, 0, 0, 0, 0, 28, 31, 26,
	32, 27, 29, 0, 0, 0, 0, 0, 0, 0,
	222, 0, 65, 66, 67, 35, 34, 33, 68, 48,
	49, 54, 55, 53, 52, 0, 0, 28, 60, 0,
	0, 27, 0, 0, 0, 0, 0, 0, 0, 0,
	98, 248, 0, 0, 0, 0, 50, 51, 56, 57,
	58, 59, 61, 64, 0, 246, 62, 63, 65, 66,
	67, 0, 0, 0, 68, 48, 49, 54, 55, 53,
	52, 0, 0, 0, 60, 65, 66, 67, 0, 0,
	0, 68, 48, 49, 54, 55, 53, 52, 0, 0,
	0, 60, 50, 51, 56, 57, 58, 59, 61, 64,
	0, 216, 62, 63, 0, 0, 0, 0, 0, 50,
	51, 56, 57, 58, 59, 61, 64, 160, 0, 62,
	63, 65, 66, 67, 0, 0, 0, 68, 48, 49,
	54, 55, 53, 52, 65, 66, 67, 60, 0, 0,
	68, 48, 49, 54, 55, 53, 52, 0, 0, 0,
	60, 0, 0, 0, 0, 50, 51, 56, 57, 58,
	59, 61, 64, 155, 0, 62, 63, 0, 50, 51,
	56, 57, 58, 59, 61, 64, 0, 0, 62, 63,
	0, 0, 0, 65, 66, 67, 0, 0, 0, 68,
	48, 49, 54, 55, 53, 52, 65, 66, 67, 60,
	0, 0, 68, 48, 49, 54, 55, 53, 52, 0,
	0, 0, 60, 0, 0, 0, 0, 50, 51, 56,
	57, 58, 59, 61, 64, 0, 0, 62, 63, 0,
	50, 51, 56, 57, 58, 59, 61, 64, 0, 0,
	62, 63, 65, 66, 67, 0, 0, 0, 68, 0,
	49, 54, 55, 53, 52, 65, 66, 67, 60, 0,
	0, 68, 0, 0, 54, 55, 53, 52, 0, 0,
	0, 60, 0, 0, 0, 0, 50, 51, 56, 57,
	58, 59, 61, 64, 0, 0, 62, 63, 0, 50,
	51, 56, 57, 58, 59, 61, 64, 47, 46, 62,
	63, 0, 0, 30, 47, 46, 0, 0, 0, 0,
	30, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	31, 26, 32, 0, 29, 0, 0, 31, 26, 32,
	0, 29, 0, 0, 0, 0, 0, 35, 34, 33,
	0, 0, 0, 0, 35, 34, 33, 0, 0, 28,
	0, 47, 46, 27, 0, 165, 28, 30, 47, 46,
	27, 163, 0, 0, 30, 89, 0, 0, 0, 0,
	0, 0, 0, 0, 31, 26, 32, 0, 29, 0,
	0, 31, 26, 32, 0, 29, 0, 0, 47, 46,
	0, 35, 34, 33, 30, 0, 0, 0, 35, 34,
	33, 0, 0, 28, 0, 0, 0, 27, 0, 0,
	28, 201, 26, 32, 27, 29, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 35, 34,
	33, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	28, 0, 0, 0, 27,
}

var yyPact = [...]int{
	-1000, -1000, 1005, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	96, -1000, -1000, 78, 76, -1000, 1784, 1612, 9, -36,
	329, 1784, -1000, 71, 1777, 73, -1000, 1784, 1784, 1784,
	1784, -1000, -1000, 1784, 1401, 193, 946, 69, -1000, 91,
	82, -1000, 113, -55, 1612, -1000, 84, 99, 1784, 1784,
	1784, 1784, 1784, 1784, 1784, 1784, 1784, 1784, 1784, 1784,
	1784, 1784, 1784, 1784, 1784, 1784, 1784, 1784, 230, -41,
	-1000, 1784, 395, 1784, 1784, 1784, 1784, 1784, 1784, 1784,
	1784, 1784, 1784, 1784, 1784, 1599, 887, 75, 1550, 1784,
	1730, 34, 1723, -1000, -1000, -1000, -1000, 1271, -1000, -53,
	-1000, -54, -18, 1336, -1000, -1000, 10, 820, -1000, 14,
	9, -1000, -3, 9, -1000, 1784, 1658, 1671, 81, 81,
	81, 81, 81, 81, 212, 212, 90, 90, 90, 90,
	90, 212, 212, 90, 90, 90, 73, 395, 1784, 8,
	-55, -1000, 73, 1612, 1612, 1612, 1612, 1612, 1612, 1612,
	1612, 1612, 1612, 1612, 1612, -1000, 1784, 1784, 3, 1784,
	-1000, 1225, 1032, 1784, -1000, -1000, 21, -1000, -73, 1784,
	-75, 1814, 1784, 1784, 761, 99, -1000, -1000, 58, -59,
	-1000, 51, -61, 1612, 73, -55, -1000, 702, 1612, 1208,
	106, 1537, 348, 1784, -1000, 1371, 1095, -46, -58, -1000,
	-1000, -24, 1288, 1612, 1612, -1000, -1000, -1000, -63, -1000,
	-1000, -65, -1000, -1000, 1784, 1784, -1000, 127, -1000, 1784,
	1162, 1082, -1000, -1000, -1000, -49, -1000, 1784, 1784, 643,
	-1000, -1000, 1145, 1491, 584, -1000, 1005, 1474, 1784, -1000,
	-1000, 1612, 1612, -1000, -1000, 1784, -1000, -1000, -1000, 269,
	525, 1428, 466, 348, -1000, -1000, -1000, -1000, -1000, 407,
	-1000,
}

var yyPgo = [...]int{
	0, 166, 99, 89, 161, 3, 0, 38, 159, 27,
	156, 153, 151, 2, 147, 137, 4, 1,
}

var yyR1 = [...]int{
	0, 1, 2, 2, 2, 2, 2, 2, 2, 2,
	10, 10, 10, 10, 10, 10, 10, 10, 10, 10,
	10, 10, 10, 10, 10, 10, 11, 11, 11, 11,
	11, 11, 12, 13, 13, 13, 15, 15, 16, 16,
	16, 16, 16, 16, 16, 16, 16, 16, 14, 14,
	14, 14, 14, 14, 3, 3, 3, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	4, 4, 5, 5, 7, 7, 8, 8, 8, 8,
	17, 17,
}

var yyR2 = [...]int{
	0, 1, 0, 2, 4, 2, 2, 2, 2, 2,
	1, 2, 4, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 5, 4, 9, 11,
	9, 7, 6, 0, 2, 5, 5, 7, 2, 3,
	4, 4, 5, 2, 3, 4, 4, 5, 1, 1,
	2, 3, 1, 2, 1, 4, 3, 1, 1, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	4, 2, 2, 2, 2, 1, 8, 4, 1, 3,
	2, 2, 4, 4, 6, 5, 5, 3, 5, 6,
	1, 3, 1, 3, 1, 3, 3, 3, 5, 5,
	0, 1,
}

var yyChk = [...]int{
	-1000, -1, -2, -15, 4, -14, -10, -11, -12, 78,
	14, 9, 10, 22, 21, 17, 16, -6, 5, -4,
	-3, 13, 18, 12, 15, -9, 38, 70, 66, 41,
	20, 37, 39, 56, 55, 54, -2, 37, -16, 28,
	56, 37, 37, -7, -6, -3, 15, 14, 31, 32,
	58, 59, 36, 35, 33, 34, 60, 61, 62, 63,
	40, 64, 68, 69, 65, 24, 25, 26, 30, -5,
	37, 57, 73, 42, 43, 44, 45, 46, 47, 48,
	49, 50, 51, 52, 53, -6, -2, 37, -6, 28,
	29, 67, 28, -6, -6, -6, -6, -6, 79, -7,
	80, -8, 37, -6, 8, -16, 67, -2, 72, -5,
	27, 72, -5, 27, 21, 73, -6, -6, -6, -6,
	-6, -6, -6, -6, -6, -6, -6, -6, -6, -6,
	-6, -6, -6, -6, -6, -6, -9, 20, 57, 73,
	-7, -3, -9, -6, -6, -6, -6, -6, -6, -6,
	-6, -6, -6, -6, -6, 4, 19, 57, 73, 23,
	7, -6, -6, 71, 37, 72, -7, 72, -17, 73,
	-17, 73, 57, 71, -2, 37, 8, 72, 27, -5,
	72, 27, -5, -6, -9, -7, 37, -2, -6, -6,
	37, -6, -2, 73, 79, 71, -6, -17, 27, 79,
	80, 37, -6, -6, -6, 8, -16, 72, -5, 72,
	72, -5, 72, 8, 73, 23, 4, -13, 11, 6,
	-6, -6, 79, 79, 72, -17, 73, 57, 71, -2,
	72, 72, -6, -6, -2, 8, -2, -6, 73, 79,
	72, -6, -6, 8, 4, 73, 4, 8, 7, -6,
	-2, -6, -2, -2, 72, 8, 4, 8, -13, -2,
	8,
}

var yyDef = [...]int{
	2, -2, 1, 3, 2, 5, 6, 7, 8, 9,
	0, 48, 49, 0, 0, 52, 0, 10, 0, 0,
	-2, 0, 2, 0, 0, 57, 58, 0, 0, 0,
	0, 54, 88, 0, 0, 0, 0, 0, 2, 0,
	0, 50, 0, 53, 104, 85, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 11,
	102, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 81, 82, 83, 84, 0, 90, 110,
	91, 110, 54, 0, 4, 2, 0, 0, 38, 0,
	0, 43, 0, 0, 51, 0, 59, 60, 61, 62,
	63, 64, 65, 66, 67, 68, 69, 70, 71, 72,
	73, 74, 75, 76, 77, 78, 79, 0, 0, 0,
	13, -2, 0, 14, 15, 16, 17, 18, 19, 20,
	21, 22, 23, 24, 25, 2, 0, 0, 0, 0,
	2, 0, 0, 0, 56, 97, 110, 89, 0, 111,
	0, 111, 0, 0, 0, 0, 87, 39, 0, 0,
	44, 0, 0, 105, 80, 12, 103, 0, 27, 0,
	0, 0, 33, 0, 55, 0, 0, 0, 110, 92,
	93, 54, 0, 106, 107, 36, 2, 40, 0, 41,
	45, 0, 46, 26, 0, 0, 2, 0, 2, 0,
	0, 0, 96, 95, 98, 0, 111, 0, 0, 0,
	42, 47, 0, 0, 0, 32, 34, 0, 0, 94,
	99, 108, 109, 37, 2, 0, 2, 31, 2, 0,
	0, 0, 0, 33, 86, 28, 2, 30, 35, 0,
	29,
}

var yyTok1 = [...]int{
	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 66, 3, 64, 68, 3,
	56, 72, 62, 60, 73, 61, 67, 63, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 71, 78,
	59, 57, 58, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 74, 3, 3, 3, 3, 3,
	3, 55, 3, 79, 65, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 54, 69, 80, 70,
}

var yyTok2 = [...]int{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 75, 76, 77,
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
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:54
		{
			ss(yylex).Stmts = yyDollar[1].expr
		}
	case 2:
		yyDollar = yyS[yypt-0 : yypt+1]
//line parser.go.y:57
		{
			yyVAL.expr = &Prog{}
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:58
		{
			yyVAL.expr = yyDollar[1].expr.(*Prog).Append(yyDollar[2].expr)
		}
	case 4:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:59
		{
			yyDollar[3].expr.(*Prog).DoBlock = true
			yyVAL.expr = yyDollar[1].expr.(*Prog).Append(yyDollar[3].expr)
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:60
		{
			yyVAL.expr = yyDollar[1].expr.(*Prog).Append(yyDollar[2].expr)
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:61
		{
			yyVAL.expr = yyDollar[1].expr.(*Prog).Append(yyDollar[2].expr)
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:62
		{
			yyVAL.expr = yyDollar[1].expr.(*Prog).Append(yyDollar[2].expr)
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:63
		{
			yyVAL.expr = yyDollar[1].expr.(*Prog).Append(yyDollar[2].expr)
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:64
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:67
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:68
		{
			yyVAL.expr = ss(yylex).pDeclareAssign([]Node(yyDollar[2].expr.(IdentList)), nil, false, yyDollar[1].token)
		}
	case 12:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:69
		{
			yyVAL.expr = ss(yylex).pDeclareAssign([]Node(yyDollar[2].expr.(IdentList)), yyDollar[4].expr.(ExprList), false, yyDollar[1].token)
		}
	case 13:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:70
		{
			yyVAL.expr = ss(yylex).pDeclareAssign([]Node(yyDollar[1].expr.(DeclList)), yyDollar[3].expr.(ExprList), true, yyDollar[2].token)
		}
	case 14:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:71
		{
			yyVAL.expr = assignLoadStore(yyDollar[1].expr, ss(yylex).pBinary(typ.OpAdd, yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token), yyDollar[2].token)
		}
	case 15:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:72
		{
			yyVAL.expr = assignLoadStore(yyDollar[1].expr, ss(yylex).pBinary(typ.OpSub, yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token), yyDollar[2].token)
		}
	case 16:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:73
		{
			yyVAL.expr = assignLoadStore(yyDollar[1].expr, ss(yylex).pBinary(typ.OpMul, yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token), yyDollar[2].token)
		}
	case 17:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:74
		{
			yyVAL.expr = assignLoadStore(yyDollar[1].expr, ss(yylex).pBinary(typ.OpDiv, yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token), yyDollar[2].token)
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:75
		{
			yyVAL.expr = assignLoadStore(yyDollar[1].expr, ss(yylex).pBinary(typ.OpIDiv, yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token), yyDollar[2].token)
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:76
		{
			yyVAL.expr = assignLoadStore(yyDollar[1].expr, ss(yylex).pBinary(typ.OpMod, yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token), yyDollar[2].token)
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:77
		{
			yyVAL.expr = assignLoadStore(yyDollar[1].expr, ss(yylex).pBitwise("and", yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token), yyDollar[2].token)
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:78
		{
			yyVAL.expr = assignLoadStore(yyDollar[1].expr, ss(yylex).pBitwise("or", yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token), yyDollar[2].token)
		}
	case 22:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:79
		{
			yyVAL.expr = assignLoadStore(yyDollar[1].expr, ss(yylex).pBitwise("xor", yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token), yyDollar[2].token)
		}
	case 23:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:80
		{
			yyVAL.expr = assignLoadStore(yyDollar[1].expr, ss(yylex).pBitwise("lsh", yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token), yyDollar[2].token)
		}
	case 24:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:81
		{
			yyVAL.expr = assignLoadStore(yyDollar[1].expr, ss(yylex).pBitwise("rsh", yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token), yyDollar[2].token)
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:82
		{
			yyVAL.expr = assignLoadStore(yyDollar[1].expr, ss(yylex).pBitwise("ursh", yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token), yyDollar[2].token)
		}
	case 26:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:85
		{
			yyVAL.expr = ss(yylex).pLoop(&If{yyDollar[2].expr, yyDollar[4].expr, emptyBreak})
		}
	case 27:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:86
		{
			yyVAL.expr = ss(yylex).pLoop(yyDollar[2].expr, &If{yyDollar[4].expr, emptyBreak, emptyProg})
		}
	case 28:
		yyDollar = yyS[yypt-9 : yypt+1]
//line parser.go.y:87
		{
			yyVAL.expr = ss(yylex).pForRange(yyDollar[2].token, yyDollar[4].expr, yyDollar[6].expr, one, yyDollar[8].expr, yyDollar[1].token)
		}
	case 29:
		yyDollar = yyS[yypt-11 : yypt+1]
//line parser.go.y:88
		{
			yyVAL.expr = ss(yylex).pForRange(yyDollar[2].token, yyDollar[4].expr, yyDollar[6].expr, yyDollar[8].expr, yyDollar[10].expr, yyDollar[1].token)
		}
	case 30:
		yyDollar = yyS[yypt-9 : yypt+1]
//line parser.go.y:89
		{
			yyVAL.expr = ss(yylex).pForIn(yyDollar[2].token, yyDollar[4].token, yyDollar[6].expr, yyDollar[8].expr, yyDollar[1].token)
		}
	case 31:
		yyDollar = yyS[yypt-7 : yypt+1]
//line parser.go.y:90
		{
			yyVAL.expr = ss(yylex).pForIn(yyDollar[2].token, yyDollar[1].token, yyDollar[4].expr, yyDollar[6].expr, yyDollar[1].token)
		}
	case 32:
		yyDollar = yyS[yypt-6 : yypt+1]
//line parser.go.y:93
		{
			yyVAL.expr = &If{yyDollar[2].expr, yyDollar[4].expr, yyDollar[5].expr}
		}
	case 33:
		yyDollar = yyS[yypt-0 : yypt+1]
//line parser.go.y:96
		{
			yyVAL.expr = nil
		}
	case 34:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:97
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 35:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:98
		{
			yyVAL.expr = &If{yyDollar[2].expr, yyDollar[4].expr, yyDollar[5].expr}
		}
	case 36:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:101
		{
			yyVAL.expr = ss(yylex).pFunc(false, yyDollar[2].token, yyDollar[3].expr, yyDollar[4].expr, yyDollar[1].token)
		}
	case 37:
		yyDollar = yyS[yypt-7 : yypt+1]
//line parser.go.y:104
		{
			m := ss(yylex).pFunc(true, __markupFuncName(yyDollar[2].token, yyDollar[4].token), yyDollar[5].expr, yyDollar[6].expr, yyDollar[1].token)
			yyVAL.expr = &Tenary{typ.OpStore, Sym(yyDollar[2].token), ss(yylex).Str(yyDollar[4].token.Str), m, yyDollar[1].token.Line()}
		}
	case 38:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:110
		{
			yyVAL.expr = (IdentList)(nil)
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:111
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 40:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:112
		{
			yyVAL.expr = IdentVarargList{yyDollar[2].expr.(IdentList)}
		}
	case 41:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:113
		{
			yyVAL.expr = IdentVarargExpandList{nil, yyDollar[3].expr.(IdentList)}
		}
	case 42:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:114
		{
			yyVAL.expr = IdentVarargExpandList{yyDollar[2].expr.(IdentList), yyDollar[4].expr.(IdentList)}
		}
	case 43:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:115
		{
			yyVAL.expr = (IdentList)(nil)
		}
	case 44:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:116
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 45:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:117
		{
			yyVAL.expr = IdentVarargList{yyDollar[2].expr.(IdentList)}
		}
	case 46:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:118
		{
			yyVAL.expr = IdentVarargExpandList{nil, yyDollar[3].expr.(IdentList)}
		}
	case 47:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:119
		{
			yyVAL.expr = IdentVarargExpandList{yyDollar[2].expr.(IdentList), yyDollar[4].expr.(IdentList)}
		}
	case 48:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:122
		{
			yyVAL.expr = &BreakContinue{true, yyDollar[1].token.Line()}
		}
	case 49:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:123
		{
			yyVAL.expr = &BreakContinue{false, yyDollar[1].token.Line()}
		}
	case 50:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:124
		{
			yyVAL.expr = &GotoLabel{yyDollar[2].token.Str, true, yyDollar[1].token.Line()}
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:125
		{
			yyVAL.expr = &GotoLabel{yyDollar[2].token.Str, false, yyDollar[1].token.Line()}
		}
	case 52:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:126
		{
			yyVAL.expr = &Unary{typ.OpRet, SNil, yyDollar[1].token.Line()}
		}
	case 53:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:127
		{
			if el := yyDollar[2].expr.(ExprList); len(el) == 1 {
				ss(yylex).pFindTailCall(el[0])
				yyVAL.expr = &Unary{typ.OpRet, el[0], yyDollar[1].token.Line()}
			} else {
				yyVAL.expr = &Unary{typ.OpRet, yyDollar[2].expr, yyDollar[1].token.Line()}
			}
		}
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:137
		{
			yyVAL.expr = Sym(yyDollar[1].token)
		}
	case 55:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:140
		{
			yyVAL.expr = &Tenary{typ.OpLoad, yyDollar[1].expr, yyDollar[3].expr, Address(typ.RegA), yyDollar[2].token.Line()}
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:143
		{
			yyVAL.expr = &Tenary{typ.OpLoad, yyDollar[1].expr, ss(yylex).Str(yyDollar[3].token.Str), Address(typ.RegA), yyDollar[2].token.Line()}
		}
	case 57:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:148
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 58:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:149
		{
			yyVAL.expr = ss(yylex).Num(yyDollar[1].token.Str)
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:150
		{
			yyVAL.expr = &Or{yyDollar[1].expr, yyDollar[3].expr}
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:151
		{
			yyVAL.expr = &And{yyDollar[1].expr, yyDollar[3].expr}
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:152
		{
			yyVAL.expr = ss(yylex).pBinary(typ.OpLess, yyDollar[3].expr, yyDollar[1].expr, yyDollar[2].token)
		}
	case 62:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:153
		{
			yyVAL.expr = ss(yylex).pBinary(typ.OpLess, yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token)
		}
	case 63:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:154
		{
			yyVAL.expr = ss(yylex).pBinary(typ.OpLessEq, yyDollar[3].expr, yyDollar[1].expr, yyDollar[2].token)
		}
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:155
		{
			yyVAL.expr = ss(yylex).pBinary(typ.OpLessEq, yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token)
		}
	case 65:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:156
		{
			yyVAL.expr = ss(yylex).pBinary(typ.OpEq, yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token)
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:157
		{
			yyVAL.expr = ss(yylex).pBinary(typ.OpNeq, yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token)
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:158
		{
			yyVAL.expr = ss(yylex).pBinary(typ.OpAdd, yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token)
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:159
		{
			yyVAL.expr = ss(yylex).pBinary(typ.OpSub, yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token)
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:160
		{
			yyVAL.expr = ss(yylex).pBinary(typ.OpMul, yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token)
		}
	case 70:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:161
		{
			yyVAL.expr = ss(yylex).pBinary(typ.OpDiv, yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token)
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:162
		{
			yyVAL.expr = ss(yylex).pBinary(typ.OpIDiv, yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token)
		}
	case 72:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:163
		{
			yyVAL.expr = ss(yylex).pBinary(typ.OpMod, yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token)
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:164
		{
			yyVAL.expr = ss(yylex).pBitwise("and", yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token)
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:165
		{
			yyVAL.expr = ss(yylex).pBitwise("or", yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token)
		}
	case 75:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:166
		{
			yyVAL.expr = ss(yylex).pBitwise("xor", yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token)
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:167
		{
			yyVAL.expr = ss(yylex).pBitwise("lsh", yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token)
		}
	case 77:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:168
		{
			yyVAL.expr = ss(yylex).pBitwise("rsh", yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token)
		}
	case 78:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:169
		{
			yyVAL.expr = ss(yylex).pBitwise("ursh", yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token)
		}
	case 79:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:170
		{
			yyVAL.expr = ss(yylex).pBinary(typ.OpIsProto, yyDollar[1].expr, yyDollar[3].expr, yyDollar[2].token)
		}
	case 80:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:171
		{
			yyVAL.expr = pUnary(typ.OpNot, ss(yylex).pBinary(typ.OpIsProto, yyDollar[1].expr, yyDollar[4].expr, yyDollar[2].token), yyDollar[2].token)
		}
	case 81:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:172
		{
			yyVAL.expr = ss(yylex).pBitwise("xor", ss(yylex).Int(-1), yyDollar[2].expr, yyDollar[1].token)
		}
	case 82:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:173
		{
			yyVAL.expr = pUnary(typ.OpLen, yyDollar[2].expr, yyDollar[1].token)
		}
	case 83:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:174
		{
			yyVAL.expr = ss(yylex).pBinary(typ.OpSub, zero, yyDollar[2].expr, yyDollar[1].token)
		}
	case 84:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:175
		{
			yyVAL.expr = pUnary(typ.OpNot, yyDollar[2].expr, yyDollar[1].token)
		}
	case 85:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:178
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 86:
		yyDollar = yyS[yypt-8 : yypt+1]
//line parser.go.y:179
		{
			yyVAL.expr = &If{yyDollar[3].expr, &Assign{Sa, yyDollar[5].expr, yyDollar[1].token.Line()}, &Assign{Sa, yyDollar[7].expr, yyDollar[1].token.Line()}}
		}
	case 87:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:180
		{
			yyVAL.expr = ss(yylex).pFunc(false, __markupLambdaName(yyDollar[1].token), yyDollar[2].expr, yyDollar[3].expr, yyDollar[1].token)
		}
	case 88:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:181
		{
			yyVAL.expr = ss(yylex).Str(yyDollar[1].token.Str)
		}
	case 89:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:182
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 90:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:183
		{
			yyVAL.expr = ss(yylex).pEmptyArray()
		}
	case 91:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:184
		{
			yyVAL.expr = ss(yylex).pEmptyObject()
		}
	case 92:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:185
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 93:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:186
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 94:
		yyDollar = yyS[yypt-6 : yypt+1]
//line parser.go.y:187
		{
			yyVAL.expr = &Tenary{typ.OpSlice, yyDollar[1].expr, yyDollar[3].expr, yyDollar[5].expr, yyDollar[2].token.Line()}
		}
	case 95:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:188
		{
			yyVAL.expr = &Tenary{typ.OpSlice, yyDollar[1].expr, zero, yyDollar[4].expr, yyDollar[2].token.Line()}
		}
	case 96:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:189
		{
			yyVAL.expr = &Tenary{typ.OpSlice, yyDollar[1].expr, yyDollar[3].expr, ss(yylex).Int(-1), yyDollar[2].token.Line()}
		}
	case 97:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:190
		{
			yyVAL.expr = &Call{typ.OpCall, yyDollar[1].expr, ExprList(nil), false, yyDollar[2].token.Line()}
		}
	case 98:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:191
		{
			yyVAL.expr = &Call{typ.OpCall, yyDollar[1].expr, yyDollar[3].expr.(ExprList), false, yyDollar[2].token.Line()}
		}
	case 99:
		yyDollar = yyS[yypt-6 : yypt+1]
//line parser.go.y:192
		{
			yyVAL.expr = &Call{typ.OpCall, yyDollar[1].expr, yyDollar[3].expr.(ExprList), true, yyDollar[2].token.Line()}
		}
	case 100:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:195
		{
			yyVAL.expr = DeclList{yyDollar[1].expr}
		}
	case 101:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:196
		{
			yyVAL.expr = append(yyDollar[1].expr.(DeclList), yyDollar[3].expr)
		}
	case 102:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:199
		{
			yyVAL.expr = IdentList{Sym(yyDollar[1].token)}
		}
	case 103:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:200
		{
			yyVAL.expr = append(yyDollar[1].expr.(IdentList), Sym(yyDollar[3].token))
		}
	case 104:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:203
		{
			yyVAL.expr = ss(yylex).pArray(nil, yyDollar[1].expr)
		}
	case 105:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:204
		{
			yyVAL.expr = ss(yylex).pArray(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 106:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:207
		{
			yyVAL.expr = ss(yylex).pObject(nil, ss(yylex).Str(yyDollar[1].token.Str), yyDollar[3].expr)
		}
	case 107:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:208
		{
			yyVAL.expr = ss(yylex).pObject(nil, yyDollar[1].expr, yyDollar[3].expr)
		}
	case 108:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:209
		{
			yyVAL.expr = ss(yylex).pObject(yyDollar[1].expr, ss(yylex).Str(yyDollar[3].token.Str), yyDollar[5].expr)
		}
	case 109:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:210
		{
			yyVAL.expr = ss(yylex).pObject(yyDollar[1].expr, yyDollar[3].expr, yyDollar[5].expr)
		}
	case 110:
		yyDollar = yyS[yypt-0 : yypt+1]
//line parser.go.y:212
		{
		}
	case 111:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:212
		{
		}
	}
	goto yystack /* stack new state and value */
}
