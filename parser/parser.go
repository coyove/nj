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

//line parser.go.y:269

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 26,
	45, 88,
	64, 88,
	-2, 81,
	-1, 110,
	45, 89,
	64, 89,
	-2, 81,
}

const yyPrivate = 57344

const yyLast = 1358

var yyAct = [...]int{
	33, 34, 192, 18, 59, 148, 92, 178, 169, 45,
	177, 26, 217, 175, 203, 199, 197, 93, 172, 187,
	183, 94, 54, 51, 151, 57, 58, 156, 66, 29,
	185, 18, 51, 149, 50, 223, 202, 200, 152, 26,
	88, 95, 96, 97, 98, 99, 91, 100, 108, 174,
	123, 32, 55, 108, 171, 111, 104, 18, 27, 121,
	164, 170, 157, 110, 25, 26, 160, 125, 126, 127,
	128, 129, 130, 131, 132, 133, 134, 135, 136, 137,
	138, 139, 140, 141, 142, 143, 144, 145, 146, 119,
	90, 116, 107, 120, 52, 153, 122, 124, 150, 179,
	101, 103, 56, 180, 109, 61, 51, 153, 31, 155,
	158, 108, 114, 53, 162, 163, 87, 165, 147, 62,
	30, 18, 28, 87, 65, 168, 79, 61, 190, 26,
	49, 115, 173, 75, 76, 77, 78, 80, 208, 3,
	48, 62, 161, 46, 7, 10, 87, 166, 9, 18,
	125, 8, 60, 181, 20, 47, 79, 26, 125, 159,
	2, 1, 184, 18, 186, 77, 78, 80, 18, 0,
	63, 26, 4, 0, 0, 0, 26, 0, 0, 0,
	196, 201, 0, 198, 0, 0, 0, 0, 0, 0,
	205, 206, 204, 0, 0, 210, 0, 18, 195, 0,
	0, 214, 0, 216, 0, 26, 212, 0, 213, 18,
	0, 18, 0, 0, 0, 0, 0, 26, 0, 26,
	225, 207, 0, 209, 228, 0, 18, 0, 18, 18,
	232, 0, 0, 0, 26, 18, 26, 26, 0, 0,
	0, 0, 0, 26, 0, 0, 0, 0, 224, 0,
	226, 0, 227, 0, 0, 0, 84, 85, 86, 0,
	233, 87, 67, 68, 73, 74, 72, 71, 0, 0,
	0, 79, 0, 0, 0, 0, 69, 70, 75, 76,
	77, 78, 80, 83, 0, 0, 81, 82, 0, 0,
	0, 0, 0, 84, 85, 86, 0, 215, 87, 67,
	68, 73, 74, 72, 71, 0, 0, 0, 79, 0,
	0, 0, 0, 69, 70, 75, 76, 77, 78, 80,
	83, 0, 0, 81, 82, 0, 0, 0, 0, 0,
	84, 85, 86, 0, 182, 87, 67, 68, 73, 74,
	72, 71, 0, 0, 218, 79, 0, 0, 0, 0,
	69, 70, 75, 76, 77, 78, 80, 83, 0, 0,
	81, 82, 0, 0, 0, 0, 84, 85, 86, 0,
	154, 87, 67, 68, 73, 74, 72, 71, 0, 0,
	0, 79, 0, 0, 0, 0, 69, 70, 75, 76,
	77, 78, 80, 83, 0, 0, 81, 82, 0, 0,
	84, 85, 86, 0, 219, 87, 67, 68, 73, 74,
	72, 71, 0, 0, 0, 79, 0, 0, 0, 0,
	69, 70, 75, 76, 77, 78, 80, 83, 0, 0,
	81, 82, 0, 0, 0, 84, 85, 86, 0, 118,
	87, 67, 68, 73, 74, 72, 71, 0, 0, 0,
	79, 0, 0, 0, 0, 69, 70, 75, 76, 77,
	78, 80, 83, 0, 0, 81, 82, 0, 0, 6,
	19, 194, 0, 189, 12, 13, 193, 23, 21, 0,
	0, 24, 17, 16, 22, 0, 0, 15, 14, 0,
	35, 0, 0, 0, 0, 0, 42, 0, 0, 0,
	0, 0, 0, 27, 0, 0, 0, 0, 0, 25,
	0, 0, 0, 105, 36, 37, 0, 106, 38, 25,
	0, 0, 0, 44, 43, 0, 6, 19, 11, 41,
	234, 12, 13, 40, 23, 21, 0, 0, 24, 17,
	16, 22, 91, 0, 15, 14, 0, 0, 0, 6,
	19, 0, 0, 231, 12, 13, 0, 23, 21, 0,
	27, 24, 17, 16, 22, 0, 25, 15, 14, 230,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 27, 0, 11, 0, 0, 0, 25,
	0, 84, 85, 86, 0, 0, 87, 67, 68, 73,
	74, 72, 71, 0, 0, 0, 79, 0, 11, 0,
	0, 69, 70, 75, 76, 77, 78, 80, 83, 6,
	19, 81, 82, 229, 12, 13, 0, 23, 21, 0,
	0, 24, 17, 16, 22, 0, 0, 15, 14, 0,
	0, 222, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 27, 0, 0, 0, 0, 0, 25,
	84, 85, 86, 0, 0, 87, 67, 68, 73, 74,
	72, 71, 0, 0, 0, 79, 0, 0, 11, 0,
	69, 70, 75, 76, 77, 78, 80, 83, 6, 19,
	81, 82, 221, 12, 13, 0, 23, 21, 0, 0,
	24, 17, 16, 22, 0, 0, 15, 14, 220, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 27, 0, 0, 0, 0, 0, 25, 0,
	84, 85, 86, 0, 0, 87, 67, 68, 73, 74,
	72, 71, 0, 0, 0, 79, 0, 11, 0, 0,
	69, 70, 75, 76, 77, 78, 80, 83, 6, 19,
	81, 82, 211, 12, 13, 0, 23, 21, 0, 0,
	24, 17, 16, 22, 0, 0, 15, 14, 191, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 27, 0, 0, 0, 0, 0, 25, 0,
	84, 85, 86, 0, 0, 87, 67, 68, 73, 74,
	72, 71, 0, 0, 0, 79, 0, 11, 0, 0,
	69, 70, 75, 76, 77, 78, 80, 83, 6, 19,
	81, 82, 188, 12, 13, 0, 23, 21, 0, 0,
	24, 17, 16, 22, 0, 0, 15, 14, 0, 0,
	0, 6, 19, 0, 0, 176, 12, 13, 0, 23,
	21, 0, 27, 24, 17, 16, 22, 0, 25, 15,
	14, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 27, 0, 11, 6, 19,
	0, 25, 167, 12, 13, 0, 23, 21, 0, 0,
	24, 17, 16, 22, 0, 0, 15, 14, 0, 0,
	11, 117, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 27, 0, 0, 0, 0, 0, 25, 0,
	84, 85, 86, 0, 0, 87, 67, 68, 73, 74,
	72, 71, 0, 0, 0, 79, 0, 11, 0, 0,
	69, 70, 75, 76, 77, 78, 80, 83, 6, 19,
	81, 82, 0, 12, 13, 0, 23, 21, 0, 0,
	24, 17, 16, 22, 113, 0, 15, 14, 112, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 27, 0, 0, 0, 0, 0, 25, 0,
	84, 85, 86, 0, 0, 87, 67, 68, 73, 74,
	72, 71, 35, 0, 0, 79, 0, 11, 42, 0,
	69, 70, 75, 76, 77, 78, 80, 83, 0, 0,
	81, 82, 0, 0, 0, 27, 36, 37, 0, 39,
	38, 25, 0, 0, 0, 44, 43, 0, 0, 0,
	0, 41, 6, 19, 0, 40, 64, 12, 13, 0,
	23, 21, 0, 89, 24, 17, 16, 22, 0, 0,
	15, 14, 0, 0, 0, 6, 19, 0, 0, 0,
	12, 13, 0, 23, 21, 5, 27, 24, 17, 16,
	22, 0, 25, 15, 14, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 27,
	0, 11, 0, 0, 0, 25, 84, 85, 86, 0,
	0, 87, 67, 68, 73, 74, 72, 71, 0, 0,
	0, 79, 0, 0, 11, 0, 69, 70, 75, 76,
	77, 78, 80, 83, 6, 19, 81, 82, 35, 12,
	13, 0, 23, 21, 42, 0, 24, 17, 16, 22,
	0, 0, 15, 14, 0, 0, 0, 0, 0, 0,
	0, 105, 36, 37, 0, 106, 38, 25, 27, 0,
	0, 44, 43, 0, 25, 0, 0, 41, 0, 0,
	0, 40, 0, 0, 0, 0, 0, 0, 102, 84,
	85, 86, 0, 11, 87, 0, 68, 73, 74, 72,
	71, 0, 0, 0, 79, 0, 0, 0, 0, 69,
	70, 75, 76, 77, 78, 80, 83, 0, 0, 81,
	82, 84, 85, 86, 0, 0, 87, 0, 0, 73,
	74, 72, 71, 35, 0, 0, 79, 0, 0, 42,
	0, 69, 70, 75, 76, 77, 78, 80, 83, 0,
	0, 81, 82, 0, 0, 0, 27, 36, 37, 35,
	39, 38, 25, 0, 0, 42, 44, 43, 0, 0,
	0, 0, 41, 0, 0, 0, 40, 0, 0, 0,
	0, 0, 105, 36, 37, 0, 106, 38, 25, 84,
	85, 86, 44, 43, 87, 0, 0, 0, 41, 0,
	0, 0, 40, 0, 79, 0, 0, 0, 0, 0,
	0, 75, 76, 77, 78, 80, 83, 0, 0, 81,
	82, 84, 85, 86, 0, 0, 87, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 79, 0, 0, 0,
	0, 0, 0, 75, 76, 77, 78, 80,
}

var yyPact = [...]int{
	-1000, 1071, -1000, -1000, -1000, 84, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, 82, 70, -1000, 1228, 100, 68,
	49, 1228, -1000, 64, 1228, 1228, -1000, -1000, 97, 1048,
	-1000, 102, -36, 1090, 100, 75, -1000, -1000, 997, -21,
	1228, 1228, 1228, 1228, 1228, -1000, 1228, 62, -1000, 1133,
	47, -1000, 1228, 20, 974, 954, 67, 904, 374, -1000,
	55, -6, -15, -1000, -1000, -1000, 1228, 1228, 1228, 1228,
	1228, 1228, 1228, 1228, 1228, 1228, 1228, 1228, 1228, 1228,
	1228, 1228, 1228, 1228, 1228, 1228, 1228, 1228, -1000, -1000,
	-31, -1000, -40, -7, 1228, -1000, -1000, -1000, -1000, -1000,
	304, -1000, -1000, -2, -40, -7, 475, 1228, 28, -36,
	-1000, 100, -1000, 1228, 1228, 22, 1228, -1000, -1000, 884,
	75, 21, -11, 21, -16, 1090, 1173, 1205, 1273, 1273,
	1273, 1273, 1273, 1273, 115, 115, 92, 92, 92, 92,
	1305, 1305, 1305, 85, 85, 85, -1000, 847, -56, 1228,
	-60, 61, 1228, 267, -1000, -45, -34, 1254, -46, -36,
	-1000, 824, 1090, 409, 104, 774, 465, -1000, -1000, -1000,
	-1000, 21, -49, -1000, 21, -50, -1000, -1000, -1000, -8,
	1228, 1090, -9, -1000, -51, -1000, -40, -1000, -1000, 1228,
	1228, -1000, 130, -1000, 1228, 754, -1000, 21, -1000, 21,
	1228, 230, 1228, -1000, -53, 340, 704, 684, -1000, 1140,
	634, -1000, -1000, -1000, 1090, -10, 1090, -1000, -1000, 1228,
	-1000, -1000, -1000, 1228, 615, 565, 545, 465, 1090, -1000,
	-1000, -1000, -1000, 522, -1000,
}

var yyPgo = [...]int{
	0, 161, 29, 160, 170, 9, 154, 34, 0, 51,
	6, 1, 151, 148, 145, 2, 144, 139, 4, 8,
	5,
}

var yyR1 = [...]int{
	0, 1, 1, 2, 2, 3, 3, 4, 4, 4,
	4, 4, 4, 12, 12, 12, 12, 13, 13, 13,
	13, 13, 13, 14, 15, 15, 15, 17, 17, 18,
	18, 18, 18, 18, 18, 19, 19, 16, 16, 16,
	16, 16, 16, 5, 5, 5, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	11, 11, 11, 11, 11, 11, 11, 11, 6, 6,
	7, 7, 9, 9, 10, 10, 10, 10, 20, 20,
}

var yyR2 = [...]int{
	0, 0, 2, 0, 2, 1, 1, 3, 1, 1,
	1, 1, 1, 1, 2, 4, 3, 5, 4, 9,
	11, 9, 7, 6, 0, 2, 5, 5, 7, 3,
	4, 5, 3, 4, 5, 0, 1, 1, 1, 2,
	3, 1, 2, 1, 4, 3, 1, 4, 1, 1,
	2, 2, 4, 4, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 2, 2, 2, 2, 2,
	3, 1, 2, 3, 5, 5, 6, 7, 1, 3,
	1, 3, 1, 3, 3, 5, 5, 7, 0, 1,
}

var yyChk = [...]int{
	-1000, -1, -3, -17, -4, 14, 4, -16, -12, -13,
	-14, 63, 9, 10, 23, 22, 18, 17, -11, 5,
	-6, 13, 19, 12, 16, 44, -5, 38, 38, -2,
	38, 38, -9, -8, -11, 15, 39, 40, 43, 42,
	58, 54, 21, 49, 48, -5, 43, 55, 40, 30,
	-7, 38, 45, 64, -8, -2, 38, -8, -8, -18,
	55, 30, 44, -4, 8, 22, 64, 32, 33, 46,
	47, 37, 36, 34, 35, 48, 49, 50, 51, 41,
	52, 56, 57, 53, 26, 27, 28, 31, -18, 66,
	-9, 67, -10, 38, 42, -8, -8, -8, -8, -8,
	-8, 38, 65, -9, -10, 38, 42, 45, 64, -9,
	-5, -11, 4, 20, 45, 64, 24, 7, 65, -2,
	38, 65, -7, 65, -7, -8, -8, -8, -8, -8,
	-8, -8, -8, -8, -8, -8, -8, -8, -8, -8,
	-8, -8, -8, -8, -8, -8, -8, -2, -20, 64,
	-20, 64, 45, -8, 66, -20, 29, 64, -20, -9,
	38, -2, -8, -8, 38, -8, -2, 8, -18, -19,
	40, 65, 29, -19, 65, 29, 8, 66, 67, 38,
	42, -8, 67, 65, -20, 64, -10, 65, 8, 64,
	24, 4, -15, 11, 6, -2, -19, 65, -19, 65,
	45, -8, 45, 65, -20, -8, -8, -2, 8, -2,
	-8, 8, -19, -19, -8, 67, -8, 65, 4, 64,
	4, 8, 7, 45, -2, -8, -2, -2, -8, 8,
	4, 8, -15, -2, 8,
}

var yyDef = [...]int{
	1, -2, 2, 5, 6, 0, 3, 8, 9, 10,
	11, 12, 37, 38, 0, 0, 41, 0, 13, 0,
	0, 0, 3, 0, 0, 0, -2, 43, 0, 0,
	39, 0, 42, 92, 46, 0, 48, 49, 0, 0,
	0, 0, 0, 0, 0, 81, 0, 0, 82, 0,
	14, 90, 0, 0, 0, 0, 0, 0, 0, 3,
	0, 0, 0, 4, 7, 40, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 3, 50,
	98, 51, 98, 0, 0, 75, 76, 77, 78, 79,
	0, 45, 83, 98, 98, 43, 0, 0, 0, 16,
	-2, 0, 3, 0, 0, 0, 0, 3, 80, 0,
	0, 35, 0, 35, 0, 93, 54, 55, 56, 57,
	58, 59, 60, 61, 62, 63, 64, 65, 66, 67,
	68, 69, 70, 71, 72, 73, 74, 0, 0, 99,
	0, 99, 0, 0, 44, 0, 98, 99, 0, 15,
	91, 0, 18, 0, 0, 0, 24, 27, 3, 29,
	36, 35, 0, 32, 35, 0, 47, 52, 53, 0,
	0, 94, 0, 84, 0, 99, 98, 85, 17, 0,
	0, 3, 0, 3, 0, 0, 30, 35, 33, 35,
	0, 0, 0, 86, 0, 0, 0, 0, 23, 25,
	0, 28, 31, 34, 96, 0, 95, 87, 3, 0,
	3, 22, 3, 0, 0, 0, 0, 24, 97, 19,
	3, 21, 26, 0, 20,
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
//line parser.go.y:55
		{
			yyVAL.expr = __chain()
			yylex.(*Lexer).Stmts = yyVAL.expr
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:59
		{
			yyVAL.expr = yyDollar[1].expr.append(yyDollar[2].expr)
			yylex.(*Lexer).Stmts = yyVAL.expr
		}
	case 3:
		yyDollar = yyS[yypt-0 : yypt+1]
//line parser.go.y:65
		{
			yyVAL.expr = __chain()
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:65
		{
			yyVAL.expr = yyDollar[1].expr.append(yyDollar[2].expr)
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:68
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:69
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 7:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:72
		{
			yyVAL.expr = __do(yyDollar[2].expr)
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:73
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:74
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:75
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:76
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:77
		{
			yyVAL.expr = emptyNode
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:80
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 14:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:83
		{
			yyVAL.expr = __chain()
			for _, v := range yyDollar[2].expr.Nodes() {
				yyVAL.expr = yyVAL.expr.append(__set(v, SNil).At(yyDollar[1].token))
			}
		}
	case 15:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:89
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
//line parser.go.y:100
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
//line parser.go.y:114
		{
			yyVAL.expr = __loop(emptyNode, __if(yyDollar[2].expr, yyDollar[4].expr, breakNode).At(yyDollar[1].token)).At(yyDollar[1].token)
		}
	case 18:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:117
		{
			yyVAL.expr = __loop(emptyNode, yyDollar[2].expr, __if(yyDollar[4].expr, breakNode, emptyNode).At(yyDollar[1].token)).At(yyDollar[1].token)
		}
	case 19:
		yyDollar = yyS[yypt-9 : yypt+1]
//line parser.go.y:120
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
//line parser.go.y:136
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
//line parser.go.y:161
		{
			yyVAL.expr = __forIn(yyDollar[2].token, yyDollar[4].token, yyDollar[6].expr, yyDollar[8].expr, yyDollar[1].token)
		}
	case 22:
		yyDollar = yyS[yypt-7 : yypt+1]
//line parser.go.y:162
		{
			yyVAL.expr = __forIn(yyDollar[2].token, yyDollar[1].token, yyDollar[4].expr, yyDollar[6].expr, yyDollar[1].token)
		}
	case 23:
		yyDollar = yyS[yypt-6 : yypt+1]
//line parser.go.y:165
		{
			yyVAL.expr = __if(yyDollar[2].expr, yyDollar[4].expr, yyDollar[5].expr).At(yyDollar[1].token)
		}
	case 24:
		yyDollar = yyS[yypt-0 : yypt+1]
//line parser.go.y:168
		{
			yyVAL.expr = Nodes()
		}
	case 25:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:168
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 26:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:168
		{
			yyVAL.expr = __if(yyDollar[2].expr, yyDollar[4].expr, yyDollar[5].expr).At(yyDollar[1].token)
		}
	case 27:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:171
		{
			yyVAL.expr = __func(yyDollar[2].token, yyDollar[3].expr, yyDollar[4].expr)
		}
	case 28:
		yyDollar = yyS[yypt-7 : yypt+1]
//line parser.go.y:172
		{
			yyVAL.expr = __store(Sym(yyDollar[2].token), Str(yyDollar[4].token.Str), __func(__markupFuncName(yyDollar[2].token, yyDollar[4].token), yyDollar[5].expr, yyDollar[6].expr))
		}
	case 29:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:175
		{
			yyVAL.expr = Nodes(emptyNode, yyDollar[3].expr)
		}
	case 30:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:176
		{
			yyVAL.expr = Nodes(yyDollar[2].expr, yyDollar[4].expr)
		}
	case 31:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:177
		{
			yyVAL.expr = Nodes(__dotdotdot(yyDollar[2].expr), yyDollar[5].expr)
		}
	case 32:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:178
		{
			yyVAL.expr = Nodes(emptyNode, yyDollar[3].expr)
		}
	case 33:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:179
		{
			yyVAL.expr = Nodes(yyDollar[2].expr, yyDollar[4].expr)
		}
	case 34:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:180
		{
			yyVAL.expr = Nodes(__dotdotdot(yyDollar[2].expr), yyDollar[5].expr)
		}
	case 35:
		yyDollar = yyS[yypt-0 : yypt+1]
//line parser.go.y:183
		{
			yyVAL.expr = nullStr
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:183
		{
			yyVAL.expr = Str(yyDollar[1].token.Str)
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:186
		{
			yyVAL.expr = Nodes(SBreak).At(yyDollar[1].token)
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:187
		{
			yyVAL.expr = Nodes(SContinue).At(yyDollar[1].token)
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:188
		{
			yyVAL.expr = __goto(Sym(yyDollar[2].token)).At(yyDollar[1].token)
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:189
		{
			yyVAL.expr = __label(Sym(yyDollar[2].token))
		}
	case 41:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:190
		{
			yyVAL.expr = __ret(SNil).At(yyDollar[1].token)
		}
	case 42:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:191
		{
			if len(yyDollar[2].expr.Nodes()) == 1 {
				__findTailCall(yyDollar[2].expr.Nodes())
				yyVAL.expr = __ret(yyDollar[2].expr.Nodes()[0]).At(yyDollar[1].token)
			} else {
				yyVAL.expr = __ret(Nodes(SArray, yyDollar[2].expr)).At(yyDollar[1].token)
			}
		}
	case 43:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:201
		{
			yyVAL.expr = Sym(yyDollar[1].token)
		}
	case 44:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:202
		{
			yyVAL.expr = __load(yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:203
		{
			yyVAL.expr = __load(yyDollar[1].expr, Str(yyDollar[3].token.Str)).At(yyDollar[2].token)
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:206
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 47:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:207
		{
			yyVAL.expr = __lambda(__markupLambdaName(yyDollar[1].token), yyDollar[2].expr, yyDollar[3].expr)
		}
	case 48:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:208
		{
			yyVAL.expr = Num(yyDollar[1].token.Str)
		}
	case 49:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:209
		{
			yyVAL.expr = Str(yyDollar[1].token.Str)
		}
	case 50:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:210
		{
			yyVAL.expr = Nodes(SArray, emptyNode).At(yyDollar[1].token)
		}
	case 51:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:211
		{
			yyVAL.expr = Nodes(SObject, emptyNode).At(yyDollar[1].token)
		}
	case 52:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:212
		{
			yyVAL.expr = Nodes(SArray, yyDollar[2].expr).At(yyDollar[1].token)
		}
	case 53:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:213
		{
			yyVAL.expr = Nodes(SObject, yyDollar[2].expr).At(yyDollar[1].token)
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:214
		{
			yyVAL.expr = Nodes((SOr), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:215
		{
			yyVAL.expr = Nodes((SAnd), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:216
		{
			yyVAL.expr = Nodes((SLess), yyDollar[3].expr, yyDollar[1].expr).At(yyDollar[2].token)
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:217
		{
			yyVAL.expr = Nodes((SLess), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:218
		{
			yyVAL.expr = Nodes((SLessEq), yyDollar[3].expr, yyDollar[1].expr).At(yyDollar[2].token)
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:219
		{
			yyVAL.expr = Nodes((SLessEq), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:220
		{
			yyVAL.expr = Nodes((SEq), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:221
		{
			yyVAL.expr = Nodes((SNeq), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 62:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:222
		{
			yyVAL.expr = Nodes((SAdd), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 63:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:223
		{
			yyVAL.expr = Nodes((SSub), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:224
		{
			yyVAL.expr = Nodes((SMul), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 65:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:225
		{
			yyVAL.expr = Nodes((SDiv), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:226
		{
			yyVAL.expr = Nodes((SIDiv), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:227
		{
			yyVAL.expr = Nodes((SMod), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:228
		{
			yyVAL.expr = Nodes((SBitAnd), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:229
		{
			yyVAL.expr = Nodes((SBitOr), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 70:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:230
		{
			yyVAL.expr = Nodes((SBitXor), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:231
		{
			yyVAL.expr = Nodes((SBitLsh), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 72:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:232
		{
			yyVAL.expr = Nodes((SBitRsh), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:233
		{
			yyVAL.expr = Nodes((SBitURsh), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:234
		{
			yyVAL.expr = Nodes((SIs), yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 75:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:235
		{
			yyVAL.expr = Nodes((SBitNot), yyDollar[2].expr).At(yyDollar[1].token)
		}
	case 76:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:236
		{
			yyVAL.expr = Nodes((SLen), yyDollar[2].expr).At(yyDollar[1].token)
		}
	case 77:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:237
		{
			yyVAL.expr = Nodes((SNot), yyDollar[2].expr).At(yyDollar[1].token)
		}
	case 78:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:238
		{
			yyVAL.expr = Nodes((SSub), zero, yyDollar[2].expr).At(yyDollar[1].token)
		}
	case 79:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:239
		{
			yyVAL.expr = Nodes((SAdd), zero, yyDollar[2].expr).At(yyDollar[1].token)
		}
	case 80:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:242
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 81:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:243
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 82:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:244
		{
			yyVAL.expr = __call(yyDollar[1].expr, Nodes(Str(yyDollar[2].token.Str))).At(yyDollar[2].token)
		}
	case 83:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:245
		{
			yyVAL.expr = __call(yyDollar[1].expr, emptyNode).At(yyDollar[2].token)
		}
	case 84:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:246
		{
			yyVAL.expr = __call(yyDollar[1].expr, yyDollar[3].expr).At(yyDollar[2].token)
		}
	case 85:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:247
		{
			yyVAL.expr = __call(yyDollar[1].expr, Nodes(Nodes(SObject, yyDollar[3].expr).At(yyDollar[2].token))).At(yyDollar[2].token)
		}
	case 86:
		yyDollar = yyS[yypt-6 : yypt+1]
//line parser.go.y:248
		{
			yyVAL.expr = __call(yyDollar[1].expr, __dotdotdot(yyDollar[3].expr)).At(yyDollar[2].token)
		}
	case 87:
		yyDollar = yyS[yypt-7 : yypt+1]
//line parser.go.y:249
		{
			yyVAL.expr = __call(yyDollar[1].expr, yyDollar[3].expr.append(Nodes(SObject, yyDollar[5].expr).At(yyDollar[2].token))).At(yyDollar[2].token)
		}
	case 88:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:252
		{
			yyVAL.expr = Nodes(yyDollar[1].expr)
		}
	case 89:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:252
		{
			yyVAL.expr = yyDollar[1].expr.append(yyDollar[3].expr)
		}
	case 90:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:255
		{
			yyVAL.expr = Nodes(Sym(yyDollar[1].token))
		}
	case 91:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:255
		{
			yyVAL.expr = yyDollar[1].expr.append(Sym(yyDollar[3].token))
		}
	case 92:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:258
		{
			yyVAL.expr = Nodes(yyDollar[1].expr)
		}
	case 93:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:258
		{
			yyVAL.expr = yyDollar[1].expr.append(yyDollar[3].expr)
		}
	case 94:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:261
		{
			yyVAL.expr = Nodes(Str(yyDollar[1].token.Str), yyDollar[3].expr)
		}
	case 95:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:262
		{
			yyVAL.expr = Nodes(yyDollar[2].expr, yyDollar[5].expr)
		}
	case 96:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:263
		{
			yyVAL.expr = yyDollar[1].expr.append(Str(yyDollar[3].token.Str)).append(yyDollar[5].expr)
		}
	case 97:
		yyDollar = yyS[yypt-7 : yypt+1]
//line parser.go.y:264
		{
			yyVAL.expr = yyDollar[1].expr.append(yyDollar[4].expr).append(yyDollar[7].expr)
		}
	case 98:
		yyDollar = yyS[yypt-0 : yypt+1]
//line parser.go.y:267
		{
			yyVAL.expr = emptyNode
		}
	case 99:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:267
		{
			yyVAL.expr = emptyNode
		}
	}
	goto yystack /* stack new state and value */
}
