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
	"'@'",
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

//line parser.go.y:349

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 25,
	40, 47,
	61, 47,
	-2, 81,
	-1, 101,
	40, 48,
	61, 48,
	-2, 81,
}

const yyPrivate = 57344

const yyLast = 1297

var yyAct = [...]int{
	31, 32, 197, 17, 134, 172, 156, 30, 155, 109,
	207, 47, 166, 48, 48, 181, 178, 41, 161, 142,
	64, 163, 52, 137, 98, 55, 135, 107, 49, 17,
	203, 177, 151, 175, 85, 138, 99, 196, 91, 92,
	93, 182, 145, 94, 87, 99, 99, 165, 108, 50,
	103, 167, 102, 135, 97, 17, 147, 100, 144, 157,
	40, 104, 25, 158, 26, 112, 113, 114, 115, 116,
	117, 118, 119, 120, 121, 122, 123, 124, 125, 126,
	127, 128, 129, 130, 131, 132, 95, 44, 25, 42,
	46, 139, 54, 136, 51, 48, 27, 29, 28, 170,
	63, 43, 141, 187, 6, 15, 143, 149, 150, 38,
	152, 101, 14, 17, 25, 146, 59, 45, 4, 53,
	75, 76, 78, 88, 89, 34, 35, 36, 90, 33,
	77, 58, 19, 3, 39, 60, 112, 5, 56, 159,
	2, 1, 0, 37, 0, 0, 0, 162, 0, 0,
	17, 0, 0, 0, 86, 17, 0, 0, 111, 176,
	0, 73, 74, 75, 76, 78, 17, 0, 0, 38,
	184, 185, 25, 77, 0, 189, 190, 0, 192, 183,
	0, 0, 17, 0, 26, 34, 35, 36, 17, 33,
	17, 0, 0, 0, 39, 0, 17, 17, 0, 0,
	209, 0, 148, 37, 212, 0, 0, 153, 17, 25,
	17, 0, 17, 17, 25, 0, 17, 218, 0, 0,
	0, 0, 17, 0, 0, 25, 0, 61, 18, 174,
	0, 0, 9, 173, 23, 21, 0, 24, 13, 12,
	22, 25, 164, 11, 10, 0, 0, 25, 0, 25,
	0, 0, 0, 0, 0, 25, 25, 26, 0, 0,
	0, 0, 180, 0, 0, 0, 0, 25, 186, 25,
	188, 25, 25, 0, 0, 25, 0, 16, 194, 195,
	0, 25, 0, 62, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 206, 0, 208, 0, 210, 0, 211,
	0, 82, 83, 84, 214, 65, 66, 71, 72, 70,
	69, 0, 0, 220, 0, 0, 0, 0, 67, 68,
	73, 74, 75, 76, 78, 81, 0, 0, 79, 80,
	0, 0, 77, 0, 0, 0, 0, 82, 83, 84,
	191, 65, 66, 71, 72, 70, 69, 0, 0, 0,
	0, 0, 0, 0, 67, 68, 73, 74, 75, 76,
	78, 81, 0, 0, 79, 80, 0, 0, 77, 0,
	0, 0, 0, 82, 83, 84, 160, 65, 66, 71,
	72, 70, 69, 0, 0, 0, 0, 0, 0, 198,
	67, 68, 73, 74, 75, 76, 78, 81, 0, 0,
	79, 80, 0, 0, 77, 0, 0, 0, 0, 82,
	83, 84, 140, 65, 66, 71, 72, 70, 69, 0,
	0, 0, 0, 0, 0, 0, 67, 68, 73, 74,
	75, 76, 78, 81, 0, 0, 79, 80, 0, 0,
	77, 0, 0, 82, 83, 84, 199, 65, 66, 71,
	72, 70, 69, 0, 0, 0, 0, 0, 0, 0,
	67, 68, 73, 74, 75, 76, 78, 81, 0, 0,
	79, 80, 0, 0, 77, 0, 0, 0, 82, 83,
	84, 133, 65, 66, 71, 72, 70, 69, 0, 0,
	0, 0, 216, 0, 0, 67, 68, 73, 74, 75,
	76, 78, 81, 0, 0, 79, 80, 0, 0, 77,
	0, 0, 82, 83, 84, 169, 65, 66, 71, 72,
	70, 69, 0, 0, 0, 0, 0, 202, 0, 67,
	68, 73, 74, 75, 76, 78, 81, 0, 0, 79,
	80, 0, 0, 77, 82, 83, 84, 0, 65, 66,
	71, 72, 70, 69, 0, 0, 200, 0, 0, 0,
	0, 67, 68, 73, 74, 75, 76, 78, 81, 0,
	0, 79, 80, 0, 0, 77, 82, 83, 84, 0,
	65, 66, 71, 72, 70, 69, 0, 0, 171, 0,
	0, 0, 0, 67, 68, 73, 74, 75, 76, 78,
	81, 0, 0, 79, 80, 0, 0, 77, 82, 83,
	84, 0, 65, 66, 71, 72, 70, 69, 0, 0,
	0, 0, 0, 110, 0, 67, 68, 73, 74, 75,
	76, 78, 81, 0, 0, 79, 80, 0, 0, 77,
	82, 83, 84, 0, 65, 66, 71, 72, 70, 69,
	0, 0, 105, 0, 0, 0, 0, 67, 68, 73,
	74, 75, 76, 78, 81, 0, 0, 79, 80, 0,
	0, 77, 82, 83, 84, 0, 65, 66, 71, 72,
	70, 69, 0, 0, 0, 0, 0, 0, 0, 67,
	68, 73, 74, 75, 76, 78, 81, 61, 18, 79,
	80, 221, 9, 77, 23, 21, 0, 24, 13, 12,
	22, 0, 0, 11, 10, 0, 0, 0, 61, 18,
	0, 0, 219, 9, 0, 23, 21, 26, 24, 13,
	12, 22, 0, 0, 11, 10, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 16, 26, 0,
	0, 0, 0, 62, 0, 0, 0, 0, 82, 83,
	84, 0, 65, 66, 71, 72, 70, 69, 16, 0,
	0, 0, 0, 0, 62, 67, 68, 73, 74, 75,
	76, 78, 81, 61, 18, 79, 80, 217, 9, 77,
	23, 21, 0, 24, 13, 12, 22, 0, 0, 11,
	10, 0, 0, 0, 61, 18, 0, 0, 215, 9,
	0, 23, 21, 26, 24, 13, 12, 22, 0, 0,
	11, 10, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 16, 26, 0, 61, 18, 0, 62,
	213, 9, 0, 23, 21, 0, 24, 13, 12, 22,
	0, 0, 11, 10, 16, 0, 0, 61, 18, 0,
	62, 205, 9, 0, 23, 21, 26, 24, 13, 12,
	22, 0, 0, 11, 10, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 16, 26, 0, 61,
	18, 0, 62, 204, 9, 0, 23, 21, 0, 24,
	13, 12, 22, 0, 0, 11, 10, 16, 0, 0,
	61, 18, 0, 62, 201, 9, 0, 23, 21, 26,
	24, 13, 12, 22, 0, 0, 11, 10, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 16,
	26, 0, 61, 18, 0, 62, 193, 9, 0, 23,
	21, 0, 24, 13, 12, 22, 0, 0, 11, 10,
	16, 0, 0, 61, 18, 0, 62, 179, 9, 0,
	23, 21, 26, 24, 13, 12, 22, 0, 0, 11,
	10, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 16, 26, 0, 61, 18, 0, 62, 168,
	9, 0, 23, 21, 0, 24, 13, 12, 22, 0,
	0, 11, 10, 16, 0, 0, 61, 18, 0, 62,
	154, 9, 0, 23, 21, 26, 24, 13, 12, 22,
	0, 0, 11, 10, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 16, 26, 61, 18, 0,
	0, 62, 9, 0, 23, 21, 0, 24, 13, 12,
	22, 106, 0, 11, 10, 0, 16, 0, 0, 0,
	0, 0, 62, 0, 61, 18, 0, 26, 57, 9,
	0, 23, 21, 0, 24, 13, 12, 22, 0, 0,
	11, 10, 0, 0, 0, 0, 0, 16, 0, 0,
	0, 7, 18, 62, 26, 0, 9, 0, 23, 21,
	20, 24, 13, 12, 22, 0, 0, 11, 10, 0,
	0, 0, 0, 0, 16, 0, 0, 0, 61, 18,
	62, 26, 0, 9, 0, 23, 21, 0, 24, 13,
	12, 22, 0, 0, 11, 10, 0, 0, 0, 0,
	0, 16, 0, 0, 0, 0, 0, 8, 26, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 82, 83,
	84, 0, 0, 66, 71, 72, 70, 69, 16, 0,
	0, 0, 0, 0, 62, 67, 68, 73, 74, 75,
	76, 78, 81, 0, 0, 79, 80, 0, 0, 77,
	82, 83, 84, 0, 0, 0, 71, 72, 70, 69,
	0, 0, 0, 0, 0, 0, 0, 67, 68, 73,
	74, 75, 76, 78, 81, 38, 0, 79, 80, 0,
	0, 77, 0, 82, 83, 84, 0, 0, 0, 0,
	26, 34, 35, 36, 0, 33, 0, 0, 0, 0,
	39, 0, 73, 74, 75, 76, 78, 81, 0, 37,
	79, 80, 0, 0, 77, 82, 83, 84, 96, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 73, 74, 75, 76, 78, 0,
	0, 0, 0, 0, 0, 0, 77,
}

var yyPact = [...]int{
	-1000, 1097, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	64, 63, -1000, 150, -1000, -1000, 150, 51, 61, -12,
	60, 150, -1000, 58, 150, -1000, -1000, 1070, -1000, 80,
	-41, 734, 51, 150, -1000, -1000, 90, 150, 150, 150,
	-1000, 734, 150, 52, -1000, -1000, 1206, -16, -1000, 150,
	30, 11, 648, 1043, -13, 616, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, 150, 150, 150, 150, 150, 150,
	150, 150, 150, 150, 150, 150, 150, 150, 150, 150,
	150, 150, 150, 150, 150, 419, -1000, -35, -38, -5,
	150, -1000, -1000, -1000, 349, -1000, -1000, -8, 150, 24,
	-41, -1000, 51, -20, 22, -1000, 150, 150, -2, 150,
	-1000, 1012, 734, 1144, 1176, 1209, 1209, 1209, 1209, 1209,
	1209, 75, 75, -1000, -1000, -1000, -1000, 1241, 1241, 1241,
	118, 118, 118, -1000, -56, 150, -58, 25, 150, 313,
	-1000, -44, -40, -41, -1000, -1000, -15, 12, 991, 734,
	454, 77, 584, 223, -1000, -1000, -1000, -7, 150, 734,
	-9, -1000, -46, -1000, 959, -1000, -47, -21, -1000, 150,
	150, -1000, 95, -1000, 150, 150, 277, 150, -1000, -1000,
	938, -1000, -1000, -25, 385, 552, 906, -1000, 1124, 520,
	734, -10, 734, -1000, 885, 853, -1000, -52, -1000, 150,
	-1000, -1000, -1000, 150, -1000, -1000, 832, -1000, 800, 488,
	779, 223, 734, -1000, 714, -1000, -1000, -1000, -1000, -1000,
	693, -1000,
}

var yyPgo = [...]int{
	0, 141, 96, 140, 138, 60, 132, 11, 0, 7,
	123, 1, 117, 135, 112, 105, 5, 131, 116, 104,
	4,
}

var yyR1 = [...]int{
	0, 1, 1, 2, 2, 3, 3, 3, 3, 3,
	3, 4, 4, 4, 4, 4, 18, 18, 13, 13,
	13, 13, 13, 14, 14, 14, 14, 14, 14, 15,
	16, 16, 16, 19, 19, 19, 19, 19, 19, 17,
	17, 17, 17, 17, 5, 5, 5, 6, 6, 7,
	7, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 11, 11, 11, 12, 12, 12, 9, 9, 10,
	10, 10, 10, 20, 20,
}

var yyR2 = [...]int{
	0, 0, 2, 0, 2, 1, 1, 1, 1, 3,
	1, 1, 1, 1, 3, 1, 1, 1, 2, 1,
	2, 4, 3, 5, 4, 9, 11, 9, 7, 6,
	0, 2, 5, 6, 7, 8, 8, 9, 10, 1,
	2, 3, 1, 2, 1, 4, 3, 1, 3, 1,
	3, 1, 3, 1, 1, 2, 4, 4, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 2, 2,
	2, 1, 2, 2, 2, 4, 5, 1, 3, 3,
	5, 5, 7, 0, 1,
}

var yyChk = [...]int{
	-1000, -1, -3, -17, -18, -13, -19, 4, 60, 9,
	21, 20, 16, 15, -14, -15, 54, -11, 5, -6,
	13, 12, 17, 11, 14, -5, 34, -2, 34, 34,
	-9, -8, -11, 39, 35, 36, 37, 53, 19, 44,
	-5, -8, 38, 50, 36, -12, 39, -7, 34, 40,
	61, 34, -8, -2, 34, -8, -4, 8, -17, -18,
	-13, 4, 60, 20, 61, 28, 29, 41, 42, 33,
	32, 30, 31, 43, 44, 45, 46, 55, 47, 51,
	52, 48, 24, 25, 26, -8, 64, -9, -10, 34,
	38, -8, -8, -8, -8, 34, 62, -9, 40, 61,
	-9, -5, -11, 39, 50, 4, 18, 40, 61, 22,
	7, -2, -8, -8, -8, -8, -8, -8, -8, -8,
	-8, -8, -8, -8, -8, -8, -8, -8, -8, -8,
	-8, -8, -8, 62, -20, 61, -20, 61, 40, -8,
	63, -20, 27, -9, 34, 62, -7, 34, -2, -8,
	-8, 34, -8, -2, 8, 64, 64, 34, 38, -8,
	63, 62, -20, 61, -2, 62, 27, 39, 8, 61,
	22, 4, -16, 10, 6, 40, -8, 40, 62, 8,
	-2, 62, 62, -7, -8, -8, -2, 8, -2, -8,
	-8, 63, -8, 8, -2, -2, 62, 27, 4, 61,
	4, 8, 7, 40, 8, 8, -2, 62, -2, -8,
	-2, -2, -8, 8, -2, 8, 4, 8, -16, 8,
	-2, 8,
}

var yyDef = [...]int{
	1, -2, 2, 5, 6, 7, 8, 3, 10, 39,
	0, 0, 42, 0, 16, 17, 0, 19, 0, 0,
	0, 0, 3, 0, 0, -2, 44, 0, 40, 0,
	43, 87, 51, 0, 53, 54, 0, 0, 0, 0,
	81, 18, 0, 0, 82, 83, 0, 20, 49, 0,
	0, 0, 0, 0, 0, 0, 4, 9, 11, 12,
	13, 3, 15, 41, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 55, 93, 93, 44,
	0, 78, 79, 80, 0, 46, 84, 93, 0, 0,
	22, -2, 0, 0, 0, 3, 0, 0, 0, 0,
	3, 0, 88, 58, 59, 60, 61, 62, 63, 64,
	65, 66, 67, 68, 69, 70, 71, 72, 73, 74,
	75, 76, 77, 52, 0, 94, 0, 94, 0, 0,
	45, 0, 93, 21, 50, 3, 0, 0, 0, 24,
	0, 0, 0, 30, 14, 56, 57, 0, 0, 89,
	0, 85, 0, 94, 0, 3, 0, 0, 23, 0,
	0, 3, 0, 3, 0, 0, 0, 0, 86, 33,
	0, 3, 3, 0, 0, 0, 0, 29, 31, 0,
	91, 0, 90, 34, 0, 0, 3, 0, 3, 0,
	3, 28, 3, 0, 35, 36, 0, 3, 0, 0,
	0, 30, 92, 37, 0, 25, 3, 27, 32, 38,
	0, 26,
}

var yyTok1 = [...]int{
	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 49, 3, 47, 51, 3,
	39, 62, 45, 43, 61, 44, 50, 46, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 60,
	42, 40, 41, 3, 54, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 56, 3, 3, 3, 3, 3,
	3, 38, 3, 63, 48, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 37, 52, 64, 53,
}

var yyTok2 = [...]int{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 55, 57, 58, 59,
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
//line parser.go.y:67
		{
			yyVAL.expr = yyDollar[1].expr.append(yyDollar[2].expr)
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:72
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:73
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:74
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:75
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 9:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:76
		{
			yyVAL.expr = __do(yyDollar[2].expr)
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:77
		{
			yyVAL.expr = emptyNode
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:80
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:81
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:82
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 14:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:83
		{
			yyVAL.expr = __do(yyDollar[2].expr)
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:84
		{
			yyVAL.expr = emptyNode
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:87
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:88
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 18:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:91
		{
			yyVAL.expr = __move(NewSymbol("$a"), yyDollar[2].expr).SetPos(yyDollar[1].token.Pos)
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:94
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 20:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:97
		{
			yyVAL.expr = __chain()
			for _, v := range yyDollar[2].expr.Nodes {
				yyVAL.expr = yyVAL.expr.append(__set(v, NewSymbol(ANil)).SetPos(yyDollar[1].token.Pos))
			}
		}
	case 21:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:103
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
	case 22:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:114
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
	case 23:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:128
		{
			yyVAL.expr = __loop(__if(yyDollar[2].expr, yyDollar[4].expr, breakNode).SetPos(yyDollar[1].token.Pos)).SetPos(yyDollar[1].token.Pos)
		}
	case 24:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:131
		{
			yyVAL.expr = __loop(
				__chain(
					yyDollar[2].expr,
					__if(yyDollar[4].expr, breakNode, emptyNode).SetPos(yyDollar[1].token.Pos),
				).SetPos(yyDollar[1].token.Pos),
			).SetPos(yyDollar[1].token.Pos)
		}
	case 25:
		yyDollar = yyS[yypt-9 : yypt+1]
//line parser.go.y:139
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
	case 26:
		yyDollar = yyS[yypt-11 : yypt+1]
//line parser.go.y:153
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
	case 27:
		yyDollar = yyS[yypt-9 : yypt+1]
//line parser.go.y:179
		{
			yyVAL.expr = __forIn(yyDollar[2].token, yyDollar[4].token, yyDollar[6].expr, yyDollar[8].expr, yyDollar[1].token.Pos)
		}
	case 28:
		yyDollar = yyS[yypt-7 : yypt+1]
//line parser.go.y:182
		{
			yyVAL.expr = __forIn(yyDollar[2].token, yyDollar[1].token, yyDollar[4].expr, yyDollar[6].expr, yyDollar[1].token.Pos)
		}
	case 29:
		yyDollar = yyS[yypt-6 : yypt+1]
//line parser.go.y:188
		{
			yyVAL.expr = __if(yyDollar[2].expr, yyDollar[4].expr, yyDollar[5].expr).SetPos(yyDollar[1].token.Pos)
		}
	case 30:
		yyDollar = yyS[yypt-0 : yypt+1]
//line parser.go.y:193
		{
			yyVAL.expr = NewComplex()
		}
	case 31:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:196
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 32:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:199
		{
			yyVAL.expr = __if(yyDollar[2].expr, yyDollar[4].expr, yyDollar[5].expr).SetPos(yyDollar[1].token.Pos)
		}
	case 33:
		yyDollar = yyS[yypt-6 : yypt+1]
//line parser.go.y:204
		{
			yyVAL.expr = __func(yyDollar[2].token, emptyNode, "", yyDollar[5].expr)
		}
	case 34:
		yyDollar = yyS[yypt-7 : yypt+1]
//line parser.go.y:205
		{
			yyVAL.expr = __func(yyDollar[2].token, yyDollar[4].expr, "", yyDollar[6].expr)
		}
	case 35:
		yyDollar = yyS[yypt-8 : yypt+1]
//line parser.go.y:206
		{
			yyVAL.expr = __func(yyDollar[2].token, __dotdotdot(yyDollar[4].expr), "", yyDollar[7].expr)
		}
	case 36:
		yyDollar = yyS[yypt-8 : yypt+1]
//line parser.go.y:207
		{
			yyVAL.expr = __store(NewSymbolFromToken(yyDollar[2].token), NewString(yyDollar[4].token.Str), __func(__markupFuncName(yyDollar[2].token, yyDollar[4].token), emptyNode, "", yyDollar[7].expr))
		}
	case 37:
		yyDollar = yyS[yypt-9 : yypt+1]
//line parser.go.y:210
		{
			yyVAL.expr = __store(NewSymbolFromToken(yyDollar[2].token), NewString(yyDollar[4].token.Str), __func(__markupFuncName(yyDollar[2].token, yyDollar[4].token), yyDollar[6].expr, "", yyDollar[8].expr))
		}
	case 38:
		yyDollar = yyS[yypt-10 : yypt+1]
//line parser.go.y:213
		{
			yyVAL.expr = __store(NewSymbolFromToken(yyDollar[2].token), NewString(yyDollar[4].token.Str), __func(__markupFuncName(yyDollar[2].token, yyDollar[4].token), __dotdotdot(yyDollar[6].expr), "", yyDollar[9].expr))
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:218
		{
			yyVAL.expr = NewComplex(NewSymbol(ABreak)).SetPos(yyDollar[1].token.Pos)
		}
	case 40:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:221
		{
			yyVAL.expr = NewComplex(NewSymbol(AGoto), NewSymbolFromToken(yyDollar[2].token)).SetPos(yyDollar[1].token.Pos)
		}
	case 41:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:224
		{
			yyVAL.expr = NewComplex(NewSymbol(ALabel), NewSymbolFromToken(yyDollar[2].token))
		}
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:227
		{
			yyVAL.expr = NewComplex(NewSymbol(AReturn), NewSymbol(ANil)).SetPos(yyDollar[1].token.Pos)
		}
	case 43:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:230
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
	case 44:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:244
		{
			yyVAL.expr = NewSymbolFromToken(yyDollar[1].token)
		}
	case 45:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:247
		{
			yyVAL.expr = __load(yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 46:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:250
		{
			yyVAL.expr = __load(yyDollar[1].expr, NewString(yyDollar[3].token.Str)).SetPos(yyDollar[2].token.Pos)
		}
	case 47:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:255
		{
			yyVAL.expr = NewComplex(yyDollar[1].expr)
		}
	case 48:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:258
		{
			yyVAL.expr = yyDollar[1].expr.append(yyDollar[3].expr)
		}
	case 49:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:263
		{
			yyVAL.expr = NewComplex(NewSymbolFromToken(yyDollar[1].token))
		}
	case 50:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:266
		{
			yyVAL.expr = yyDollar[1].expr.append(NewSymbolFromToken(yyDollar[3].token))
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:271
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:272
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 53:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:273
		{
			yyVAL.expr = NewNumberFromString(yyDollar[1].token.Str)
		}
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:274
		{
			yyVAL.expr = NewString(yyDollar[1].token.Str)
		}
	case 55:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:275
		{
			yyVAL.expr = NewComplex(NewSymbol(AArrayMap), emptyNode).SetPos(yyDollar[1].token.Pos)
		}
	case 56:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:276
		{
			yyVAL.expr = NewComplex(NewSymbol(AArray), yyDollar[2].expr).SetPos(yyDollar[1].token.Pos)
		}
	case 57:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:277
		{
			yyVAL.expr = NewComplex(NewSymbol(AArrayMap), yyDollar[2].expr).SetPos(yyDollar[1].token.Pos)
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:278
		{
			yyVAL.expr = NewComplex(NewSymbol(AOr), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:279
		{
			yyVAL.expr = NewComplex(NewSymbol(AAnd), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:280
		{
			yyVAL.expr = NewComplex(NewSymbol(ALess), yyDollar[3].expr, yyDollar[1].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:281
		{
			yyVAL.expr = NewComplex(NewSymbol(ALess), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 62:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:282
		{
			yyVAL.expr = NewComplex(NewSymbol(ALessEq), yyDollar[3].expr, yyDollar[1].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 63:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:283
		{
			yyVAL.expr = NewComplex(NewSymbol(ALessEq), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:284
		{
			yyVAL.expr = NewComplex(NewSymbol(AEq), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 65:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:285
		{
			yyVAL.expr = NewComplex(NewSymbol(ANeq), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:286
		{
			yyVAL.expr = NewComplex(NewSymbol(AAdd), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:287
		{
			yyVAL.expr = NewComplex(NewSymbol(ASub), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:288
		{
			yyVAL.expr = NewComplex(NewSymbol(AMul), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:289
		{
			yyVAL.expr = NewComplex(NewSymbol(ADiv), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 70:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:290
		{
			yyVAL.expr = NewComplex(NewSymbol(AIDiv), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:291
		{
			yyVAL.expr = NewComplex(NewSymbol(AMod), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 72:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:292
		{
			yyVAL.expr = NewComplex(NewSymbol(ABitAnd), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:293
		{
			yyVAL.expr = NewComplex(NewSymbol(ABitOr), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:294
		{
			yyVAL.expr = NewComplex(NewSymbol(ABitXor), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 75:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:295
		{
			yyVAL.expr = NewComplex(NewSymbol(ABitLsh), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:296
		{
			yyVAL.expr = NewComplex(NewSymbol(ABitRsh), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 77:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:297
		{
			yyVAL.expr = NewComplex(NewSymbol(ABitURsh), yyDollar[1].expr, yyDollar[3].expr).SetPos(yyDollar[2].token.Pos)
		}
	case 78:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:298
		{
			yyVAL.expr = NewComplex(NewSymbol(ABitNot), yyDollar[2].expr).SetPos(yyDollar[1].token.Pos)
		}
	case 79:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:299
		{
			yyVAL.expr = NewComplex(NewSymbol(ANot), yyDollar[2].expr).SetPos(yyDollar[1].token.Pos)
		}
	case 80:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:300
		{
			yyVAL.expr = NewComplex(NewSymbol(ASub), zeroNode, yyDollar[2].expr).SetPos(yyDollar[1].token.Pos)
		}
	case 81:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:303
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 82:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:306
		{
			yyVAL.expr = __call(yyDollar[1].expr, NewComplex(NewString(yyDollar[2].token.Str))).SetPos(yyDollar[1].expr.Pos())
		}
	case 83:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:309
		{
			yyDollar[2].expr.Nodes[1] = yyDollar[1].expr
			yyVAL.expr = yyDollar[2].expr
		}
	case 84:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser.go.y:315
		{
			yyVAL.expr = __call(emptyNode, emptyNode).SetPos(yyDollar[1].token.Pos)
		}
	case 85:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser.go.y:318
		{
			yyVAL.expr = __call(emptyNode, yyDollar[2].expr).SetPos(yyDollar[1].token.Pos)
		}
	case 86:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:321
		{
			yyVAL.expr = __call(emptyNode, __dotdotdot(yyDollar[2].expr)).SetPos(yyDollar[1].token.Pos)
		}
	case 87:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:326
		{
			yyVAL.expr = NewComplex(yyDollar[1].expr)
		}
	case 88:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:329
		{
			yyVAL.expr = yyDollar[1].expr.append(yyDollar[3].expr)
		}
	case 89:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser.go.y:334
		{
			yyVAL.expr = NewComplex(NewString(yyDollar[1].token.Str), yyDollar[3].expr)
		}
	case 90:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:337
		{
			yyVAL.expr = NewComplex(yyDollar[2].expr, yyDollar[5].expr)
		}
	case 91:
		yyDollar = yyS[yypt-5 : yypt+1]
//line parser.go.y:340
		{
			yyVAL.expr = yyDollar[1].expr.append(NewString(yyDollar[3].token.Str)).append(yyDollar[5].expr)
		}
	case 92:
		yyDollar = yyS[yypt-7 : yypt+1]
//line parser.go.y:343
		{
			yyVAL.expr = yyDollar[1].expr.append(yyDollar[4].expr).append(yyDollar[7].expr)
		}
	case 93:
		yyDollar = yyS[yypt-0 : yypt+1]
//line parser.go.y:347
		{
			yyVAL.expr = emptyNode
		}
	case 94:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser.go.y:347
		{
			yyVAL.expr = emptyNode
		}
	}
	goto yystack /* stack new state and value */
}
