// Code generated by goyacc -o parser.go parser.go.y. DO NOT EDIT.

//line parser.go.y:2
package parser

import __yyfmt__ "fmt"

//line parser.go.y:2

//line parser.go.y:24
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
const TLambda = 57357
const TIf = 57358
const TReturn = 57359
const TReturnVoid = 57360
const TRepeat = 57361
const TUntil = 57362
const TNot = 57363
const TLabel = 57364
const TGoto = 57365
const TIn = 57366
const TNext = 57367
const TLsh = 57368
const TRsh = 57369
const TURsh = 57370
const TDotDotDot = 57371
const TLParen = 57372
const TIs = 57373
const TOr = 57374
const TAnd = 57375
const TEqeq = 57376
const TNeq = 57377
const TLte = 57378
const TGte = 57379
const TIdent = 57380
const TNumber = 57381
const TString = 57382
const TIDiv = 57383
const ASSIGN = 57384
const FUNC = 57385
const UNARY = 57386

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
	"TLambda",
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
	"TLParen",
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

//line parser.go.y:264

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 27,
	45, 85,
	64, 85,
	-2, 79,
	-1, 109,
	45, 86,
	64, 86,
	-2, 79,
}

const yyPrivate = 57344

const yyLast = 1365

var yyAct = [...]int{
	34, 35, 91, 18, 92, 188, 174, 147, 93, 45,
	173, 27, 58, 209, 197, 193, 192, 183, 179, 50,
	106, 50, 53, 171, 150, 56, 155, 57, 65, 181,
	215, 51, 18, 90, 30, 148, 49, 196, 60, 107,
	27, 94, 95, 96, 97, 98, 122, 99, 120, 87,
	52, 103, 61, 33, 110, 169, 18, 54, 107, 170,
	194, 156, 109, 59, 27, 151, 124, 125, 126, 127,
	128, 129, 130, 131, 132, 133, 134, 135, 136, 137,
	138, 139, 140, 141, 142, 143, 144, 145, 163, 159,
	107, 168, 89, 118, 152, 60, 115, 121, 123, 149,
	175, 28, 102, 25, 176, 108, 152, 26, 119, 61,
	154, 157, 48, 161, 162, 100, 164, 113, 55, 50,
	18, 86, 146, 32, 86, 46, 31, 29, 27, 86,
	186, 78, 167, 64, 78, 202, 114, 47, 74, 75,
	76, 77, 79, 76, 77, 79, 160, 3, 18, 124,
	7, 165, 177, 10, 9, 36, 27, 124, 8, 182,
	158, 42, 18, 180, 62, 20, 4, 18, 2, 1,
	27, 0, 0, 0, 0, 27, 0, 195, 104, 37,
	25, 0, 105, 38, 26, 0, 199, 200, 44, 43,
	198, 204, 0, 18, 41, 206, 0, 208, 40, 0,
	0, 27, 191, 18, 0, 18, 0, 90, 0, 36,
	0, 27, 217, 27, 0, 42, 220, 0, 18, 0,
	18, 18, 201, 0, 203, 224, 27, 18, 27, 27,
	0, 0, 28, 37, 25, 27, 39, 38, 26, 0,
	0, 0, 44, 43, 0, 216, 0, 218, 41, 219,
	0, 0, 40, 83, 84, 85, 0, 225, 86, 66,
	67, 72, 73, 71, 70, 0, 0, 0, 78, 0,
	0, 0, 0, 68, 69, 74, 75, 76, 77, 79,
	82, 0, 0, 80, 81, 0, 0, 0, 0, 0,
	83, 84, 85, 0, 207, 86, 66, 67, 72, 73,
	71, 70, 0, 0, 0, 78, 0, 0, 0, 0,
	68, 69, 74, 75, 76, 77, 79, 82, 0, 0,
	80, 81, 0, 0, 0, 0, 0, 83, 84, 85,
	0, 178, 86, 66, 67, 72, 73, 71, 70, 0,
	0, 210, 78, 0, 0, 0, 0, 68, 69, 74,
	75, 76, 77, 79, 82, 0, 0, 80, 81, 0,
	0, 0, 0, 83, 84, 85, 0, 153, 86, 66,
	67, 72, 73, 71, 70, 0, 0, 0, 78, 0,
	0, 0, 0, 68, 69, 74, 75, 76, 77, 79,
	82, 0, 0, 80, 81, 0, 0, 83, 84, 85,
	0, 211, 86, 66, 67, 72, 73, 71, 70, 0,
	0, 0, 78, 0, 0, 0, 0, 68, 69, 74,
	75, 76, 77, 79, 82, 0, 0, 80, 81, 0,
	0, 0, 83, 84, 85, 0, 117, 86, 66, 67,
	72, 73, 71, 70, 0, 0, 0, 78, 0, 0,
	0, 0, 68, 69, 74, 75, 76, 77, 79, 82,
	0, 0, 80, 81, 0, 0, 6, 19, 190, 0,
	185, 12, 13, 189, 23, 21, 0, 0, 24, 17,
	16, 22, 0, 0, 15, 14, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	28, 0, 25, 6, 19, 0, 26, 226, 12, 13,
	0, 23, 21, 0, 0, 24, 17, 16, 22, 0,
	0, 15, 14, 0, 0, 11, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 28, 0, 25,
	6, 19, 0, 26, 223, 12, 13, 0, 23, 21,
	0, 0, 24, 17, 16, 22, 0, 0, 15, 14,
	0, 0, 11, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 28, 0, 25, 6, 19, 0,
	26, 221, 12, 13, 0, 23, 21, 0, 0, 24,
	17, 16, 22, 0, 0, 15, 14, 0, 0, 11,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 28, 0, 25, 6, 19, 0, 26, 213, 12,
	13, 0, 23, 21, 0, 0, 24, 17, 16, 22,
	0, 0, 15, 14, 0, 0, 11, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 28, 0,
	25, 6, 19, 0, 26, 205, 12, 13, 0, 23,
	21, 0, 0, 24, 17, 16, 22, 0, 0, 15,
	14, 0, 0, 11, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 28, 0, 25, 6, 19,
	0, 26, 184, 12, 13, 0, 23, 21, 0, 0,
	24, 17, 16, 22, 0, 0, 15, 14, 0, 0,
	11, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 28, 0, 25, 6, 19, 0, 26, 172,
	12, 13, 0, 23, 21, 0, 0, 24, 17, 16,
	22, 0, 0, 15, 14, 0, 0, 11, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 28,
	0, 25, 6, 19, 0, 26, 166, 12, 13, 0,
	23, 21, 0, 0, 24, 17, 16, 22, 0, 0,
	15, 14, 0, 0, 11, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 28, 0, 25, 6,
	19, 0, 26, 0, 12, 13, 0, 23, 21, 0,
	0, 24, 17, 16, 22, 112, 0, 15, 14, 0,
	0, 11, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 28, 0, 25, 6, 19, 0, 26,
	63, 12, 13, 0, 23, 21, 0, 0, 24, 17,
	16, 22, 0, 0, 15, 14, 0, 0, 11, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	28, 0, 25, 6, 19, 0, 26, 0, 12, 13,
	0, 23, 21, 5, 0, 24, 17, 16, 22, 0,
	0, 15, 14, 222, 0, 11, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 28, 0, 25,
	0, 0, 0, 26, 0, 83, 84, 85, 0, 0,
	86, 66, 67, 72, 73, 71, 70, 0, 214, 0,
	78, 0, 11, 0, 0, 68, 69, 74, 75, 76,
	77, 79, 82, 0, 0, 80, 81, 83, 84, 85,
	0, 0, 86, 66, 67, 72, 73, 71, 70, 0,
	0, 0, 78, 0, 0, 0, 0, 68, 69, 74,
	75, 76, 77, 79, 82, 6, 19, 80, 81, 0,
	12, 13, 0, 23, 21, 0, 0, 24, 17, 16,
	22, 0, 0, 15, 14, 212, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 28,
	0, 25, 0, 0, 0, 26, 0, 83, 84, 85,
	0, 0, 86, 66, 67, 72, 73, 71, 70, 187,
	0, 0, 78, 0, 11, 0, 0, 68, 69, 74,
	75, 76, 77, 79, 82, 0, 0, 80, 81, 0,
	0, 83, 84, 85, 0, 0, 86, 66, 67, 72,
	73, 71, 70, 0, 116, 0, 78, 0, 0, 0,
	0, 68, 69, 74, 75, 76, 77, 79, 82, 0,
	0, 80, 81, 83, 84, 85, 0, 0, 86, 66,
	67, 72, 73, 71, 70, 111, 0, 0, 78, 0,
	0, 0, 0, 68, 69, 74, 75, 76, 77, 79,
	82, 0, 0, 80, 81, 0, 0, 83, 84, 85,
	0, 0, 86, 66, 67, 72, 73, 71, 70, 36,
	0, 0, 78, 0, 0, 42, 0, 68, 69, 74,
	75, 76, 77, 79, 82, 0, 0, 80, 81, 0,
	0, 0, 28, 37, 25, 0, 39, 38, 26, 0,
	0, 0, 44, 43, 0, 0, 0, 0, 41, 0,
	0, 0, 40, 0, 0, 0, 83, 84, 85, 0,
	88, 86, 66, 67, 72, 73, 71, 70, 36, 0,
	0, 78, 0, 0, 42, 0, 68, 69, 74, 75,
	76, 77, 79, 82, 0, 0, 80, 81, 0, 0,
	0, 104, 37, 25, 0, 105, 38, 26, 0, 0,
	0, 44, 43, 0, 0, 0, 0, 41, 0, 0,
	0, 40, 83, 84, 85, 0, 0, 86, 101, 67,
	72, 73, 71, 70, 0, 0, 0, 78, 0, 0,
	0, 0, 68, 69, 74, 75, 76, 77, 79, 82,
	0, 0, 80, 81, 83, 84, 85, 0, 0, 86,
	0, 0, 72, 73, 71, 70, 36, 0, 0, 78,
	0, 0, 42, 0, 68, 69, 74, 75, 76, 77,
	79, 82, 0, 0, 80, 81, 0, 0, 0, 104,
	37, 25, 0, 105, 38, 26, 83, 84, 85, 44,
	43, 86, 0, 0, 0, 41, 0, 0, 0, 40,
	0, 78, 0, 0, 0, 0, 0, 0, 74, 75,
	76, 77, 79, 82, 0, 0, 80, 81, 83, 84,
	85, 0, 0, 86, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 78, 0, 0, 0, 0, 0, 0,
	74, 75, 76, 77, 79,
}

var yyPact = [...]int{
	-1000, 869, -1000, -1000, -1000, 89, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, 88, 85, -1000, 194, 82, 81,
	-14, 194, -1000, 80, 194, -1000, 194, -1000, -1000, 8,
	832, -1000, 111, -36, 1150, 82, 65, -1000, 1114, -34,
	194, 194, 194, 194, 194, -1000, 194, 77, 1173, -25,
	-1000, 194, 63, 1091, 795, 72, 1057, 371, -1000, 70,
	-17, -19, -1000, -1000, -1000, 194, 194, 194, 194, 194,
	194, 194, 194, 194, 194, 194, 194, 194, 194, 194,
	194, 194, 194, 194, 194, 194, 194, -1000, -1000, -29,
	-1000, -40, 20, 194, -1000, -1000, -1000, -1000, -1000, 301,
	-1000, -1000, -3, -40, 20, 140, 194, 51, -36, -1000,
	82, -1000, 194, 194, 50, 194, -1000, -1000, 758, 65,
	-1000, 26, -1000, -6, 1150, 1206, 1238, 1280, 1280, 1280,
	1280, 1280, 1280, 93, 93, 98, 98, 98, 98, 1312,
	1312, 1312, 90, 90, 90, -1000, 721, -56, 194, -61,
	62, 194, 264, -1000, -47, -35, 1261, -48, -36, -1000,
	684, 1150, 406, 106, 1025, 462, -1000, -1000, -1000, -49,
	-1000, -50, -1000, -1000, -1000, 15, 194, 1150, -8, -1000,
	-51, -1000, -40, -1000, -1000, 194, 194, -1000, 127, -1000,
	194, 647, -1000, -1000, 194, 227, 194, -1000, -52, 337,
	991, 610, -1000, 971, 921, -1000, 1150, -15, 1150, -1000,
	-1000, 194, -1000, -1000, -1000, 194, 573, 889, 536, 462,
	1150, -1000, -1000, -1000, -1000, 499, -1000,
}

var yyPgo = [...]int{
	0, 169, 34, 168, 164, 9, 165, 36, 0, 53,
	2, 1, 158, 154, 153, 5, 150, 147, 12, 7,
}

var yyR1 = [...]int{
	0, 1, 1, 2, 2, 3, 3, 4, 4, 4,
	4, 4, 4, 12, 12, 12, 12, 13, 13, 13,
	13, 13, 13, 14, 15, 15, 15, 17, 17, 18,
	18, 18, 18, 18, 18, 16, 16, 16, 16, 16,
	16, 5, 5, 5, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 11, 11, 11,
	11, 11, 11, 11, 11, 6, 6, 7, 7, 9,
	9, 10, 10, 10, 10, 19, 19,
}

var yyR2 = [...]int{
	0, 0, 2, 0, 2, 1, 1, 3, 1, 1,
	1, 1, 1, 1, 2, 4, 3, 5, 4, 9,
	11, 9, 7, 6, 0, 2, 5, 5, 7, 2,
	3, 4, 2, 3, 4, 1, 1, 2, 3, 1,
	2, 1, 4, 3, 1, 4, 1, 2, 2, 4,
	4, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 2, 2, 2, 2, 2, 1, 3, 1,
	3, 5, 5, 6, 7, 1, 3, 1, 3, 1,
	3, 3, 5, 5, 7, 0, 1,
}

var yyChk = [...]int{
	-1000, -1, -3, -17, -4, 14, 4, -16, -12, -13,
	-14, 63, 9, 10, 23, 22, 18, 17, -11, 5,
	-6, 13, 19, 12, 16, 40, 44, -5, 38, 38,
	-2, 38, 38, -9, -8, -11, 15, 39, 43, 42,
	58, 54, 21, 49, 48, -5, 43, 55, 30, -7,
	38, 45, 64, -8, -2, 38, -8, -8, -18, 55,
	30, 44, -4, 8, 22, 64, 32, 33, 46, 47,
	37, 36, 34, 35, 48, 49, 50, 51, 41, 52,
	56, 57, 53, 26, 27, 28, 31, -18, 66, -9,
	67, -10, 38, 42, -8, -8, -8, -8, -8, -8,
	38, 65, -9, -10, 38, 42, 45, 64, -9, -5,
	-11, 4, 20, 45, 64, 24, 7, 65, -2, 38,
	65, -7, 65, -7, -8, -8, -8, -8, -8, -8,
	-8, -8, -8, -8, -8, -8, -8, -8, -8, -8,
	-8, -8, -8, -8, -8, -8, -2, -19, 64, -19,
	64, 45, -8, 66, -19, 29, 64, -19, -9, 38,
	-2, -8, -8, 38, -8, -2, 8, -18, 65, 29,
	65, 29, 8, 66, 67, 38, 42, -8, 67, 65,
	-19, 64, -10, 65, 8, 64, 24, 4, -15, 11,
	6, -2, 65, 65, 45, -8, 45, 65, -19, -8,
	-8, -2, 8, -2, -8, 8, -8, 67, -8, 65,
	4, 64, 4, 8, 7, 45, -2, -8, -2, -2,
	-8, 8, 4, 8, -15, -2, 8,
}

var yyDef = [...]int{
	1, -2, 2, 5, 6, 0, 3, 8, 9, 10,
	11, 12, 35, 36, 0, 0, 39, 0, 13, 0,
	0, 0, 3, 0, 0, 77, 0, -2, 41, 0,
	0, 37, 0, 40, 89, 44, 0, 46, 0, 0,
	0, 0, 0, 0, 0, 79, 0, 0, 0, 14,
	87, 0, 0, 0, 0, 0, 0, 0, 3, 0,
	0, 0, 4, 7, 38, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 3, 47, 95,
	48, 95, 0, 0, 72, 73, 74, 75, 76, 0,
	43, 80, 95, 95, 41, 0, 0, 0, 16, -2,
	0, 3, 0, 0, 0, 0, 3, 78, 0, 0,
	29, 0, 32, 0, 90, 51, 52, 53, 54, 55,
	56, 57, 58, 59, 60, 61, 62, 63, 64, 65,
	66, 67, 68, 69, 70, 71, 0, 0, 96, 0,
	96, 0, 0, 42, 0, 95, 96, 0, 15, 88,
	0, 18, 0, 0, 0, 24, 27, 3, 30, 0,
	33, 0, 45, 49, 50, 0, 0, 91, 0, 81,
	0, 96, 95, 82, 17, 0, 0, 3, 0, 3,
	0, 0, 31, 34, 0, 0, 0, 83, 0, 0,
	0, 0, 23, 25, 0, 28, 93, 0, 92, 84,
	3, 0, 3, 22, 3, 0, 0, 0, 0, 24,
	94, 19, 3, 21, 26, 0, 20,
}

var yyTok1 = [...]int{
	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 54, 3, 52, 56, 3,
	44, 65, 50, 48, 64, 49, 55, 51, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 63,
	47, 45, 46, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 59, 3, 3, 3, 3, 3,
	3, 43, 3, 66, 53, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 42, 57, 67, 58,
}

var yyTok2 = [...]int{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	60, 61, 62,
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
			yylex.(*Lexer).Stmts = yyVAL.expr
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:58
		{
			yyVAL.expr = yyDollar[1].expr.append(yyDollar[2].expr)
			yylex.(*Lexer).Stmts = yyVAL.expr
		}
	case 3:
		yyDollar = yyS[yypt-0 : yypt+1]
//line parser.go.y:64
		{
			yyVAL.expr = __chain()
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:64
		{
			yyVAL.expr = yyDollar[1].expr.append(yyDollar[2].expr)
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:67
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:68
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 7:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:71
		{
			yyVAL.expr = __do(yyDollar[2].expr)
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:72
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:73
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:74
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:75
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:76
		{
			yyVAL.expr = emptyNode
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:79
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 14:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:82
		{
			yyVAL.expr = __chain()
			for _, v := range yyDollar[2].expr.Nodes() {
				yyVAL.expr = yyVAL.expr.append(__set(v, SNil).At(yyDollar[1].token))
			}
		}
	case 15:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:88
		{
			if len(yyDollar[4].expr.Nodes()) == 1 && len(yyDollar[2].expr.Nodes()) > 1 {
				tmp := randomVarname()
				yyVAL.expr = __chain(__set(tmp, yyDollar[4].expr.Nodes()[0]).At(yyDollar[1].token))
				for i, ident := range yyDollar[2].expr.Nodes() {
					yyVAL.expr = yyVAL.expr.append(__set(ident, __load(tmp, Int(int64(i))).At(yyDollar[1].token)).At(yyDollar[1].token))
				}
			} else {
				yyVAL.expr = __local(yyDollar[2].expr.Nodes(), yyDollar[4].expr.Nodes(), yyDollar[1].token)
			}
		}
	case 16:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:99
		{
			if len(yyDollar[3].expr.Nodes()) == 1 && len(yyDollar[1].expr.Nodes()) > 1 {
				tmp := randomVarname()
				yyVAL.expr = __chain(__set(tmp, yyDollar[3].expr.Nodes()[0]).At(yyDollar[2].token))
				for i, decl := range yyDollar[1].expr.Nodes() {
					x := decl.moveLoadStore(__move, __load(tmp, Int(int64(i))).At(yyDollar[2].token)).At(yyDollar[2].token)
					yyVAL.expr = yyVAL.expr.append(x)
				}
			} else {
				yyVAL.expr = __moveMulti(yyDollar[1].expr.Nodes(), yyDollar[3].expr.Nodes(), yyDollar[2].token)
			}
		}
	case 17:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:113
		{
			yyVAL.expr = __loop(emptyNode, __if(yyDollar[2].expr, yyDollar[4].expr, breakNode).At(yyDollar[1].token)).At(yyDollar[1].token)
		}
	case 18:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:116
		{
			yyVAL.expr = __loop(emptyNode, yyDollar[2].expr, __if(yyDollar[4].expr, breakNode, emptyNode).At(yyDollar[1].token)).At(yyDollar[1].token)
		}
	case 19:
		yyDollar = yyS[yypt-9 : yypt+1]
//line parser.go.y:119
		{
			forVar, forEnd := Sym(yyDollar[2].token), randomVarname()
			cont := __inc(forVar, one).At(yyDollar[1].token)
			yyVAL.expr = __do(
				__set(forVar, yyDollar[4].expr).At(yyDollar[1].token),
				__set(forEnd, yyDollar[6].expr).At(yyDollar[1].token),
				__loop(
					cont,
					__if(
						__less(forVar, forEnd),
						__chain(yyDollar[8].expr, cont),
						breakNode,
					).At(yyDollar[1].token),
				).At(yyDollar[1].token),
			)
		}
	case 20:
		yyDollar = yyS[yypt-11 : yypt+1]
//line parser.go.y:135
		{
			forVar, forEnd, forStep := Sym(yyDollar[2].token), randomVarname(), randomVarname()
			body := __chain(yyDollar[10].expr, __inc(forVar, forStep))
			yyVAL.expr = __do(
				__set(forVar, yyDollar[4].expr).At(yyDollar[1].token),
				__set(forEnd, yyDollar[6].expr).At(yyDollar[1].token),
				__set(forStep, yyDollar[8].expr).At(yyDollar[1].token),
			)
			if yyDollar[8].expr.IsNum() { // step is a static number, easy case
				if yyDollar[8].expr.IsNegativeNumber() {
					yyVAL.expr = yyVAL.expr.append(__loop(__inc(forVar, forStep), __if(__less(forEnd, forVar), body, breakNode).At(yyDollar[1].token)).At(yyDollar[1].token))
				} else {
					yyVAL.expr = yyVAL.expr.append(__loop(__inc(forVar, forStep), __if(__less(forVar, forEnd), body, breakNode).At(yyDollar[1].token)).At(yyDollar[1].token))
				}
			} else {
				yyVAL.expr = yyVAL.expr.append(__loop(
					__inc(forVar, forStep),
					__if(
						__less(zero, forStep).At(yyDollar[1].token),
						__if(__lessEq(forEnd, forVar), breakNode, body).At(yyDollar[1].token), // +step
						__if(__lessEq(forVar, forEnd), breakNode, body).At(yyDollar[1].token), // -step
					).At(yyDollar[1].token),
				).At(yyDollar[1].token))
			}
		}
	case 21:
		yyDollar = yyS[yypt-9 : yypt+1]
//line parser.go.y:160
		{
			yyVAL.expr = __forIn(yyDollar[2].token, yyDollar[4].token, yyDollar[6].expr, yyDollar[8].expr, yyDollar[1].token)
		}
	case 22:
		yyDollar = yyS[yypt-7 : yypt+1]
//line parser.go.y:161
		{
			yyVAL.expr = __forIn(yyDollar[2].token, yyDollar[1].token, yyDollar[4].expr, yyDollar[6].expr, yyDollar[1].token)
		}
	case 23:
		yyDollar = yyS[yypt-6 : yypt+1]
//line parser.go.y:164
		{
			yyVAL.expr = __if(yyDollar[2].expr, yyDollar[4].expr, yyDollar[5].expr).At(yyDollar[1].token)
		}
	case 24:
		yyDollar = yyS[yypt-0 : yypt+1]
//line parser.go.y:167
		{
			yyVAL.expr = Nodes()
		}
	case 25:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:167
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 26:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:167
		{
			yyVAL.expr = __if(yyDollar[2].expr, yyDollar[4].expr, yyDollar[5].expr).At(yyDollar[1].token)
		}
	case 27:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:170
		{
			yyVAL.expr = __func(yyDollar[2].token, yyDollar[3].expr, yyDollar[4].expr)
		}
	case 28:
		yyDollar = yyS[yypt-7 : yypt+1]
//line parser.go.y:171
		{
			yyVAL.expr = __store(Sym(yyDollar[2].token), Str(yyDollar[4].token.Str), __func(__markupFuncName(yyDollar[2].token, yyDollar[4].token), yyDollar[5].expr, yyDollar[6].expr))
		}
	case 29:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:174
		{
			yyVAL.expr = emptyNode
		}
	case 30:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:175
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 31:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:176
		{
			yyVAL.expr = __dotdotdot(yyDollar[2].expr)
		}
	case 32:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:177
		{
			yyVAL.expr = emptyNode
		}
	case 33:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:178
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 34:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:179
		{
			yyVAL.expr = __dotdotdot(yyDollar[2].expr)
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:182
		{
			yyVAL.expr = Nodes(SBreak).At(yyDollar[1].token)
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:183
		{
			yyVAL.expr = Nodes(SContinue).At(yyDollar[1].token)
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:184
		{
			yyVAL.expr = __goto(Sym(yyDollar[2].token)).At(yyDollar[1].token)
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:185
		{
			yyVAL.expr = __label(Sym(yyDollar[2].token))
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:186
		{
			yyVAL.expr = __ret(SNil).At(yyDollar[1].token)
		}
	case 40:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:187
		{
			if len(yyDollar[2].expr.Nodes()) == 1 {
				__findTailCall(yyDollar[2].expr.Nodes())
				yyVAL.expr = __ret(yyDollar[2].expr.Nodes()[0]).At(yyDollar[1].token)
			} else {
				yyVAL.expr = __ret(Nodes(SArray, yyDollar[2].expr)).At(yyDollar[1].token)
			}
		}
	case 41:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:197
		{
			yyVAL.expr = Sym(yyDollar[1].token)
		}
	case 42:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:198
		{
			yyVAL.expr = __load(yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:199
		{
			yyVAL.expr = __load(yyDollar[1].expr, Str(yyDollar[3].token.Str)).At(yyDollar[2].token)
		}
	case 44:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:202
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 45:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:203
		{
			yyVAL.expr = __lambda(__markupLambdaName(yyDollar[1].token), yyDollar[2].expr, yyDollar[3].expr)
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:204
		{
			yyVAL.expr = Num(yyDollar[1].token.Str)
		}
	case 47:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:205
		{
			yyVAL.expr = Nodes(SArray, emptyNode).At(yyDollar[1].token)
		}
	case 48:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:206
		{
			yyVAL.expr = Nodes(SObject, emptyNode).At(yyDollar[1].token)
		}
	case 49:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:207
		{
			yyVAL.expr = Nodes(SArray, yyDollar[2].expr).At(yyDollar[1].token)
		}
	case 50:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:208
		{
			yyVAL.expr = Nodes(SObject, yyDollar[2].expr).At(yyDollar[1].token)
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:209
		{
			yyVAL.expr = Nodes((SOr), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:210
		{
			yyVAL.expr = Nodes((SAnd), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:211
		{
			yyVAL.expr = Nodes((SLess), yyDollar[3].expr, yyDollar[1].expr).At(yyDollar[2].token)
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:212
		{
			yyVAL.expr = Nodes((SLess), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:213
		{
			yyVAL.expr = Nodes((SLessEq), yyDollar[3].expr, yyDollar[1].expr).At(yyDollar[2].token)
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:214
		{
			yyVAL.expr = Nodes((SLessEq), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:215
		{
			yyVAL.expr = Nodes((SEq), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:216
		{
			yyVAL.expr = Nodes((SNeq), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:217
		{
			yyVAL.expr = Nodes((SAdd), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:218
		{
			yyVAL.expr = Nodes((SSub), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:219
		{
			yyVAL.expr = Nodes((SMul), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 62:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:220
		{
			yyVAL.expr = Nodes((SDiv), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 63:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:221
		{
			yyVAL.expr = Nodes((SIDiv), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:222
		{
			yyVAL.expr = Nodes((SMod), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 65:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:223
		{
			yyVAL.expr = Nodes((SBitAnd), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:224
		{
			yyVAL.expr = Nodes((SBitOr), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:225
		{
			yyVAL.expr = Nodes((SBitXor), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:226
		{
			yyVAL.expr = Nodes((SBitLsh), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:227
		{
			yyVAL.expr = Nodes((SBitRsh), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 70:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:228
		{
			yyVAL.expr = Nodes((SBitURsh), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:229
		{
			yyVAL.expr = Nodes((SIs), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 72:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:230
		{
			yyVAL.expr = Nodes((SBitNot), yyDollar[2].expr).At(yyDollar[1].token)
		}
	case 73:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:231
		{
			yyVAL.expr = Nodes((SLen), yyDollar[2].expr).At(yyDollar[1].token)
		}
	case 74:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:232
		{
			yyVAL.expr = Nodes((SNot), yyDollar[2].expr).At(yyDollar[1].token)
		}
	case 75:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:233
		{
			yyVAL.expr = Nodes((SSub), zero, yyDollar[2].expr).At(yyDollar[1].token)
		}
	case 76:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:234
		{
			yyVAL.expr = Nodes((SAdd), zero, yyDollar[2].expr).At(yyDollar[1].token)
		}
	case 77:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:237
		{
			yyVAL.expr = Str(yyDollar[1].token.Str)
		}
	case 78:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:238
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 79:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:239
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 80:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:240
		{
			yyVAL.expr = __call(yyDollar[1].expr, emptyNode).At(yyDollar[2].token)
		}
	case 81:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:241
		{
			yyVAL.expr = __call(yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 82:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:242
		{
			yyVAL.expr = __call(yyDollar[1].expr, Nodes(Nodes(SObject, yyDollar[3].expr).At(yyDollar[2].token))).At(yyDollar[2].token)
		}
	case 83:
		yyDollar = yyS[yypt-6 : yypt+1]
//line parser.go.y:243
		{
			yyVAL.expr = __call(yyDollar[1].expr, __dotdotdot(yyDollar[3].expr)).At(yyDollar[2].token)
		}
	case 84:
		yyDollar = yyS[yypt-7 : yypt+1]
//line parser.go.y:244
		{
			yyVAL.expr = __call(yyDollar[1].expr, yyDollar[3].expr.append(Nodes(SObject, yyDollar[5].expr).At(yyDollar[2].token))).At(yyDollar[2].token)
		}
	case 85:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:247
		{
			yyVAL.expr = Nodes(yyDollar[1].expr)
		}
	case 86:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:247
		{
			yyVAL.expr = yyDollar[1].expr.append(yyDollar[3].expr)
		}
	case 87:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:250
		{
			yyVAL.expr = Nodes(Sym(yyDollar[1].token))
		}
	case 88:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:250
		{
			yyVAL.expr = yyDollar[1].expr.append(Sym(yyDollar[3].token))
		}
	case 89:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:253
		{
			yyVAL.expr = Nodes(yyDollar[1].expr)
		}
	case 90:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:253
		{
			yyVAL.expr = yyDollar[1].expr.append(yyDollar[3].expr)
		}
	case 91:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:256
		{
			yyVAL.expr = Nodes(Str(yyDollar[1].token.Str), yyDollar[3].expr)
		}
	case 92:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:257
		{
			yyVAL.expr = Nodes(yyDollar[2].expr, yyDollar[5].expr)
		}
	case 93:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:258
		{
			yyVAL.expr = yyDollar[1].expr.append(Str(yyDollar[3].token.Str)).append(yyDollar[5].expr)
		}
	case 94:
		yyDollar = yyS[yypt-7 : yypt+1]
//line parser.go.y:259
		{
			yyVAL.expr = yyDollar[1].expr.append(yyDollar[4].expr).append(yyDollar[7].expr)
		}
	case 95:
		yyDollar = yyS[yypt-0 : yypt+1]
//line parser.go.y:262
		{
			yyVAL.expr = emptyNode
		}
	case 96:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:262
		{
			yyVAL.expr = emptyNode
		}
	}
	goto yystack /* stack new state and value */
}
