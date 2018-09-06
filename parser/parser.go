//line .\parser.go.y:1
package parser

import __yyfmt__ "fmt"

//line .\parser.go.y:3
import (
	"fmt"
	"github.com/coyove/common/rand"
	"path/filepath"
)

//line .\parser.go.y:44
type yySymType struct {
	yys   int
	token Token
	expr  *Node
	str   string
}

const TAddressof = 57346
const TAssert = 57347
const TBreak = 57348
const TCase = 57349
const TContinue = 57350
const TElse = 57351
const TFor = 57352
const TFunc = 57353
const TIf = 57354
const TLen = 57355
const TNew = 57356
const TNil = 57357
const TNot = 57358
const TReturn = 57359
const TUse = 57360
const TSwitch = 57361
const TTypeof = 57362
const TVar = 57363
const TWhile = 57364
const TYield = 57365
const TAddAdd = 57366
const TSubSub = 57367
const TEqeq = 57368
const TNeq = 57369
const TLsh = 57370
const TRsh = 57371
const TURsh = 57372
const TLte = 57373
const TGte = 57374
const TIdent = 57375
const TNumber = 57376
const TString = 57377
const TAddEq = 57378
const TSubEq = 57379
const TMulEq = 57380
const TDivEq = 57381
const TModEq = 57382
const TAndEq = 57383
const TOrEq = 57384
const TXorEq = 57385
const TLshEq = 57386
const TRshEq = 57387
const TURshEq = 57388
const ASSIGN = 57389
const FUN = 57390
const TOr = 57391
const TAnd = 57392
const UNARY = 57393
const TMinMin = 57394

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"TAddressof",
	"TAssert",
	"TBreak",
	"TCase",
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
	"TSwitch",
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
	"TAddEq",
	"TSubEq",
	"TMulEq",
	"TDivEq",
	"TModEq",
	"TAndEq",
	"TOrEq",
	"TXorEq",
	"TLshEq",
	"TRshEq",
	"TURshEq",
	"'T'",
	"ASSIGN",
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

//line .\parser.go.y:387

var _rand = rand.New()

func randomName() string {
	return fmt.Sprintf("%x", _rand.Fetch(16))
}

func expandSwitch(sub *Node, cases []*Node) *Node {
	subject := ANodeS("switch" + randomName())
	ret := CNode("chain", CNode("set", subject, sub))

	var lastif, root *Node
	var defaultCase *Node

	for i := 0; i < len(cases); i += 2 {
		if cases[i].S() == "else" {
			defaultCase = cases[i+1]
			continue
		}

		casestat := CNode("if", CNode("==", subject, cases[i]), cases[i+1])
		if lastif != nil {
			lastif.Cappend(CNode("chain", casestat))
		} else {
			root = casestat
		}
		lastif = casestat
	}

	if defaultCase == nil {
		lastif.Cappend(CNode("chain"))
	} else {
		if root == nil {
			ret.Cappend(defaultCase)
			return ret
		}

		lastif.Cappend(defaultCase)
	}

	ret.Cappend(root)
	return ret
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 1250

var yyAct = [...]int{

	205, 116, 73, 42, 151, 25, 72, 175, 32, 203,
	238, 53, 54, 198, 134, 199, 6, 229, 176, 182,
	57, 175, 59, 60, 119, 135, 181, 43, 136, 26,
	71, 140, 64, 67, 215, 209, 1, 107, 108, 119,
	110, 184, 141, 68, 118, 212, 111, 112, 113, 114,
	180, 124, 51, 123, 27, 56, 118, 49, 18, 23,
	25, 25, 130, 25, 3, 142, 178, 65, 147, 148,
	106, 6, 115, 201, 153, 133, 117, 15, 55, 191,
	132, 186, 14, 144, 26, 26, 173, 26, 154, 155,
	156, 157, 158, 159, 160, 161, 162, 163, 164, 165,
	166, 167, 168, 169, 170, 171, 172, 211, 109, 27,
	27, 131, 27, 13, 23, 97, 98, 99, 177, 3,
	179, 146, 63, 127, 101, 102, 103, 61, 58, 129,
	174, 5, 15, 189, 139, 149, 185, 14, 187, 30,
	121, 41, 192, 25, 195, 150, 40, 197, 66, 69,
	4, 196, 16, 188, 17, 44, 95, 96, 97, 98,
	99, 95, 96, 97, 98, 99, 62, 26, 13, 70,
	122, 2, 0, 0, 200, 0, 0, 202, 0, 0,
	0, 0, 204, 137, 206, 0, 5, 0, 25, 0,
	213, 25, 27, 0, 0, 218, 217, 0, 221, 0,
	216, 0, 0, 0, 0, 0, 0, 0, 0, 224,
	225, 25, 26, 226, 0, 26, 230, 0, 231, 0,
	0, 0, 0, 233, 0, 0, 0, 0, 25, 25,
	0, 0, 0, 0, 0, 26, 241, 27, 0, 0,
	27, 0, 0, 25, 25, 25, 25, 25, 0, 0,
	0, 0, 26, 26, 6, 6, 0, 6, 6, 0,
	27, 0, 0, 190, 0, 239, 240, 26, 26, 26,
	26, 26, 0, 0, 242, 243, 0, 27, 27, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 27, 27, 27, 27, 27, 23, 23, 0,
	23, 23, 3, 3, 0, 3, 3, 0, 207, 0,
	0, 210, 0, 0, 0, 15, 15, 0, 15, 15,
	14, 14, 0, 14, 14, 0, 0, 0, 0, 10,
	8, 223, 9, 0, 20, 28, 21, 0, 0, 0,
	0, 11, 12, 22, 0, 24, 19, 7, 234, 236,
	0, 13, 13, 0, 13, 13, 0, 31, 0, 0,
	18, 0, 29, 0, 0, 244, 0, 0, 0, 5,
	5, 0, 5, 5, 10, 8, 0, 9, 0, 20,
	0, 21, 0, 0, 0, 0, 11, 12, 22, 0,
	24, 19, 7, 126, 0, 93, 94, 101, 102, 103,
	92, 91, 31, 0, 0, 18, 0, 29, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 87, 88, 104, 105, 100, 89, 90, 95,
	96, 97, 98, 99, 10, 8, 0, 9, 0, 20,
	0, 21, 235, 0, 0, 0, 11, 12, 22, 0,
	24, 19, 7, 0, 0, 93, 94, 101, 102, 103,
	92, 91, 31, 0, 0, 18, 0, 29, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 87, 88, 104, 105, 100, 89, 90, 95,
	96, 97, 98, 99, 10, 8, 0, 128, 0, 20,
	0, 21, 0, 0, 0, 0, 11, 12, 22, 0,
	24, 19, 7, 0, 0, 93, 94, 101, 102, 103,
	92, 91, 31, 0, 0, 18, 0, 29, 93, 94,
	101, 102, 103, 92, 91, 0, 0, 0, 0, 0,
	0, 0, 87, 88, 104, 105, 100, 89, 90, 95,
	96, 97, 98, 99, 0, 87, 88, 104, 105, 100,
	89, 90, 95, 96, 97, 98, 99, 0, 0, 74,
	75, 0, 0, 193, 0, 0, 194, 93, 94, 101,
	102, 103, 92, 91, 76, 77, 78, 79, 80, 81,
	82, 83, 84, 85, 86, 93, 94, 101, 102, 103,
	92, 91, 0, 0, 87, 88, 104, 105, 100, 89,
	90, 95, 96, 97, 98, 99, 93, 94, 101, 102,
	103, 92, 91, 0, 0, 0, 152, 89, 90, 95,
	96, 97, 98, 99, 0, 0, 93, 94, 101, 102,
	103, 92, 91, 87, 88, 104, 105, 100, 89, 90,
	95, 96, 97, 98, 99, 93, 94, 101, 102, 103,
	92, 91, 0, 0, 237, 104, 105, 100, 89, 90,
	95, 96, 97, 98, 99, 0, 0, 0, 0, 0,
	0, 0, 87, 88, 104, 105, 100, 89, 90, 95,
	96, 97, 98, 99, 93, 94, 101, 102, 103, 92,
	91, 0, 0, 228, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 87, 88, 104, 105, 100, 89, 90, 95, 96,
	97, 98, 99, 93, 94, 101, 102, 103, 92, 91,
	0, 0, 222, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	87, 88, 104, 105, 100, 89, 90, 95, 96, 97,
	98, 99, 93, 94, 101, 102, 103, 92, 91, 0,
	0, 183, 38, 0, 0, 0, 0, 0, 0, 28,
	0, 37, 39, 33, 47, 0, 35, 0, 36, 87,
	88, 104, 105, 100, 89, 90, 95, 96, 97, 98,
	99, 31, 34, 52, 50, 0, 29, 0, 0, 208,
	0, 0, 0, 0, 0, 0, 0, 93, 94, 101,
	102, 103, 92, 91, 0, 0, 0, 0, 0, 45,
	0, 0, 0, 0, 46, 48, 93, 94, 101, 102,
	103, 92, 91, 145, 87, 88, 104, 105, 100, 89,
	90, 95, 96, 97, 98, 99, 0, 0, 0, 0,
	0, 0, 232, 87, 88, 104, 105, 100, 89, 90,
	95, 96, 97, 98, 99, 38, 0, 0, 0, 0,
	0, 220, 28, 0, 37, 39, 33, 47, 0, 35,
	0, 36, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 31, 34, 52, 50, 0, 29,
	38, 0, 0, 0, 0, 0, 0, 28, 0, 37,
	39, 33, 47, 0, 35, 0, 36, 0, 0, 0,
	0, 0, 45, 0, 0, 0, 0, 46, 48, 31,
	34, 52, 50, 0, 29, 143, 0, 0, 0, 38,
	0, 0, 0, 0, 0, 0, 28, 0, 37, 39,
	33, 47, 0, 35, 0, 36, 0, 45, 0, 0,
	0, 0, 46, 48, 0, 0, 0, 219, 31, 34,
	52, 50, 0, 29, 0, 0, 38, 0, 0, 0,
	0, 227, 0, 28, 0, 37, 39, 33, 47, 0,
	35, 0, 36, 0, 0, 0, 45, 0, 0, 0,
	0, 46, 48, 0, 120, 31, 34, 52, 50, 0,
	29, 38, 0, 0, 0, 0, 214, 0, 28, 0,
	37, 39, 33, 47, 0, 35, 0, 36, 0, 0,
	0, 0, 0, 45, 0, 0, 0, 0, 46, 48,
	31, 34, 52, 50, 0, 29, 0, 0, 93, 94,
	101, 102, 103, 92, 91, 0, 0, 0, 138, 0,
	93, 94, 101, 102, 103, 92, 91, 0, 45, 125,
	0, 0, 0, 46, 48, 87, 88, 104, 105, 100,
	89, 90, 95, 96, 97, 98, 99, 87, 88, 104,
	105, 100, 89, 90, 95, 96, 97, 98, 99, 38,
	0, 0, 0, 0, 0, 0, 28, 0, 37, 39,
	33, 47, 0, 35, 0, 36, 93, 94, 101, 102,
	103, 92, 91, 0, 0, 0, 0, 0, 31, 34,
	52, 50, 0, 29, 0, 93, 94, 101, 102, 103,
	92, 91, 0, 87, 88, 104, 105, 100, 89, 90,
	95, 96, 97, 98, 99, 0, 45, 0, 0, 0,
	0, 46, 48, 88, 104, 105, 100, 89, 90, 95,
	96, 97, 98, 99, 10, 8, 0, 9, 0, 20,
	28, 21, 0, 0, 0, 0, 11, 12, 22, 0,
	24, 19, 7, 0, 0, 0, 10, 8, 0, 9,
	0, 20, 31, 21, 0, 18, 0, 29, 11, 12,
	22, 0, 24, 19, 7, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 31, 0, 0, 18, 0, 29,
}
var yyPact = [...]int{

	-1000, 1189, -1000, -1000, -1000, -1000, -1000, 1115, -1000, -1000,
	1115, 1115, 43, -1000, -1000, -1000, -1000, -1000, -1000, 1115,
	95, 1115, 1115, 94, 89, -5, -1000, -27, -1000, 1115,
	-1000, 545, 1110, -1000, -1000, 35, 1115, 1115, 75, 1115,
	-1000, -1000, -5, -1000, -1000, 1115, 1115, 1115, 1115, 6,
	955, -1000, -1000, 1054, 1110, -1000, 324, 489, -45, 429,
	1042, 18, -42, -28, 881, 50, -1000, 778, 1115, 124,
	-71, 551, -1000, 1115, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 1115, 1115, 1115,
	1115, 1115, 1115, 1115, 1115, 1115, 1115, 1115, 1115, 1115,
	1115, 1115, 1115, 1115, 1115, 1115, -1000, -1000, -1000, -1000,
	68, -1000, -1000, -1000, -1000, 22, -52, 1115, -9, -1000,
	-1000, -19, -47, -54, 707, -1000, -1000, -1000, -29, -1000,
	-1000, -1000, -1000, -1000, -1000, 1115, 48, 129, 126, 1211,
	46, 1115, 502, 1115, 545, -1000, -60, 1110, 1110, -1000,
	-1000, -1000, -1000, 1110, 1129, 610, 96, 96, 96, 96,
	96, 96, 53, 53, -1000, -1000, -1000, 569, 101, 101,
	101, 569, 569, 1115, -1000, 40, 1115, 1110, -1000, -66,
	-1000, 1115, 1115, 1115, 1211, 746, -35, 1211, 38, 1027,
	-1000, -36, 1110, 545, 916, 820, -1000, 1115, 1115, -1000,
	-1000, -1000, 1110, -1000, 668, 1110, 1110, 1211, 1115, 1115,
	-1000, -1000, 992, 629, -57, 1115, -1000, 1115, 801, -1000,
	-1000, 1110, 1115, -1000, 369, 429, 590, -64, -1000, -1000,
	1110, 1110, -1000, 1110, -1000, 1115, -1000, -1000, -1000, 1189,
	1189, 429, 1189, 1189, -1000,
}
var yyPgo = [...]int{

	0, 36, 14, 171, 52, 1, 53, 170, 166, 0,
	27, 6, 2, 155, 3, 154, 129, 111, 80, 75,
	153, 123, 62, 152, 150, 57, 149, 139, 148, 146,
	72, 141, 140,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 16, 3, 3, 3, 3, 21,
	21, 21, 21, 21, 21, 24, 24, 24, 24, 15,
	15, 15, 15, 11, 11, 12, 12, 12, 12, 12,
	12, 12, 12, 12, 12, 12, 10, 10, 10, 10,
	10, 10, 17, 17, 17, 17, 17, 18, 18, 19,
	20, 20, 20, 20, 26, 26, 26, 25, 23, 22,
	22, 22, 22, 22, 22, 22, 4, 4, 4, 4,
	4, 4, 5, 5, 6, 6, 7, 7, 8, 8,
	8, 8, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 13, 14,
	14, 14, 14, 27, 28, 28, 29, 29, 29, 30,
	30, 31, 31, 32, 32, 32, 32,
}
var yyR2 = [...]int{

	0, 0, 2, 3, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 2,
	1, 1, 3, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 2, 3, 5, 4,
	6, 5, 3, 6, 7, 9, 7, 3, 5, 5,
	4, 4, 5, 5, 0, 2, 2, 2, 4, 2,
	1, 1, 2, 3, 2, 2, 1, 4, 6, 5,
	5, 3, 1, 3, 1, 3, 3, 5, 1, 3,
	5, 3, 1, 1, 2, 2, 2, 2, 2, 4,
	1, 1, 1, 1, 1, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 2, 2, 2, 2, 1, 1,
	3, 1, 3, 2, 2, 3, 3, 4, 3, 2,
	3, 2, 3, 1, 2, 1, 2,
}
var yyChk = [...]int{

	-1000, -1, -3, -22, -24, -16, -2, 23, 6, 8,
	5, 17, 18, -17, -18, -19, -23, -15, 36, 22,
	10, 12, 19, -25, 21, -14, -10, -4, 11, 38,
	-27, 33, -9, 15, 34, 18, 20, 13, 4, 14,
	-29, -31, -14, -10, -13, 61, 66, 16, 67, -25,
	36, -4, 35, -9, -9, 35, -1, -9, 33, -9,
	-9, 33, -8, 33, 37, 72, -28, 38, 70, -26,
	-27, -9, -11, -12, 24, 25, 39, 40, 41, 42,
	43, 44, 45, 46, 47, 48, 49, 53, 54, 58,
	59, 32, 31, 26, 27, 60, 61, 62, 63, 64,
	57, 28, 29, 30, 55, 56, 35, -9, -9, 33,
	-9, -9, -9, -9, -9, -30, -5, 70, 38, 33,
	69, -32, -7, -6, -9, 35, 69, -21, 8, -16,
	-22, -17, -18, -19, -2, 70, 73, -21, 36, -30,
	73, 70, -9, 74, 33, 75, -6, -9, -9, 11,
	21, 75, 75, -9, -9, -9, -9, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, -9, -9, -9, -9,
	-9, -9, -9, 18, -2, 73, 70, -9, 75, -5,
	69, 73, 73, 74, 70, -9, 33, 9, -20, 7,
	-21, 33, -9, 71, 74, -9, -11, -12, 73, 75,
	-9, 33, -9, 75, -9, -9, -9, -21, 73, 70,
	-21, 69, 7, -9, 9, 70, -11, -12, -9, 71,
	71, -9, 74, -21, -9, -9, -9, 9, 74, 74,
	-9, -9, 71, -9, -21, 73, -21, 74, 74, -1,
	-1, -9, -1, -1, -21,
}
var yyDef = [...]int{

	1, -2, 2, 5, 6, 7, 8, 0, 60, 61,
	0, 0, 0, 15, 16, 17, 18, 4, 1, 0,
	0, 0, 0, 0, 0, 20, 21, 119, 54, 0,
	121, 66, 59, 82, 83, 0, 0, 0, 0, 0,
	90, 91, 92, 93, 94, 0, 0, 0, 0, 0,
	0, 119, 118, 62, 64, 65, 0, 0, 0, 0,
	0, 0, 19, 78, 0, 0, 123, 0, 0, 57,
	121, 0, 36, 0, 23, 24, 25, 26, 27, 28,
	29, 30, 31, 32, 33, 34, 35, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 84, 85, 86, 87,
	88, 114, 115, 116, 117, 0, 0, 0, 0, 72,
	131, 0, 133, 135, 74, 63, 3, 42, 61, 9,
	10, 11, 12, 13, 14, 0, 0, 47, 0, 0,
	0, 0, 0, 0, 71, 124, 0, 74, 22, 55,
	56, 120, 122, 37, 95, 96, 97, 98, 99, 100,
	101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
	111, 112, 113, 0, 126, 0, 0, 128, 129, 0,
	132, 134, 136, 0, 0, 0, 0, 0, 0, 0,
	58, 81, 79, 67, 0, 0, 39, 0, 0, 125,
	89, 73, 127, 130, 0, 75, 76, 0, 0, 0,
	48, 49, 0, 0, 0, 0, 38, 0, 0, 69,
	70, 41, 0, 43, 0, 0, 0, 0, 1, 1,
	80, 40, 68, 77, 44, 0, 46, 1, 1, 50,
	51, 0, 52, 53, 45,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 67, 3, 64, 56, 3,
	38, 75, 62, 60, 73, 61, 72, 63, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 74, 3,
	59, 70, 58, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 50, 3, 3, 3, 3, 3,
	3, 37, 3, 71, 57, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 36, 55, 69, 66,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 39, 40, 41, 42, 43, 44,
	45, 46, 47, 48, 49, 51, 52, 53, 54, 65,
	68,
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
		//line .\parser.go.y:79
		{
			yyVAL.expr = CNode("chain")
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:85
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[2].expr)
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.expr
			}
		}
	case 3:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:93
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:96
		{
			if yyDollar[1].expr.isIsolatedCopy() {
				yyDollar[1].expr.Cx(2).C()[0] = NNode(0.0)
			}
			yyVAL.expr = yyDollar[1].expr
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:104
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:105
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:106
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:107
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:110
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:111
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:112
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:113
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:114
		{
			yyVAL.expr = CNode("chain", yyDollar[1].expr)
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:115
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:118
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:119
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:120
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:121
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:124
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 20:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:125
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:126
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 22:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:127
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
	case 23:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:147
		{
			yyVAL.expr = NNode(1.0)
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:148
		{
			yyVAL.expr = NNode(-1.0)
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:151
		{
			yyVAL.str = "+"
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:152
		{
			yyVAL.str = "-"
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:153
		{
			yyVAL.str = "*"
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:154
		{
			yyVAL.str = "/"
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:155
		{
			yyVAL.str = "%"
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:156
		{
			yyVAL.str = "&"
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:157
		{
			yyVAL.str = "|"
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:158
		{
			yyVAL.str = "^"
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:159
		{
			yyVAL.str = "<<"
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:160
		{
			yyVAL.str = ">>"
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:161
		{
			yyVAL.str = ">>>"
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:164
		{
			yyVAL.expr = CNode("inc", ANode(yyDollar[1].token).setPos(yyDollar[1].token), yyDollar[2].expr)
		}
	case 37:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:165
		{
			yyVAL.expr = CNode("move", ANode(yyDollar[1].token), CNode(yyDollar[2].str, ANode(yyDollar[1].token).setPos(yyDollar[1].token), yyDollar[3].expr))
		}
	case 38:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:166
		{
			yyVAL.expr = CNode("store", yyDollar[1].expr, yyDollar[3].expr, CNode("+", CNode("load", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr), yyDollar[5].expr).setPos0(yyDollar[1].expr))
		}
	case 39:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:167
		{
			yyVAL.expr = CNode("store", yyDollar[1].expr, yyDollar[3].token, CNode("+", CNode("load", yyDollar[1].expr, yyDollar[3].token).setPos0(yyDollar[1].expr), yyDollar[4].expr).setPos0(yyDollar[1].expr))
		}
	case 40:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line .\parser.go.y:168
		{
			yyVAL.expr = CNode("store", yyDollar[1].expr, yyDollar[3].expr, CNode(yyDollar[5].str, CNode("load", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr), yyDollar[6].expr).setPos0(yyDollar[1].expr))
		}
	case 41:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:171
		{
			yyVAL.expr = CNode("store", yyDollar[1].expr, yyDollar[3].token, CNode(yyDollar[4].str, CNode("load", yyDollar[1].expr, yyDollar[3].token).setPos0(yyDollar[1].expr), yyDollar[5].expr).setPos0(yyDollar[1].expr))
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:176
		{
			yyVAL.expr = CNode("for", yyDollar[2].expr, CNode(), yyDollar[3].expr).setPos0(yyDollar[1].token)
		}
	case 43:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line .\parser.go.y:179
		{
			yyVAL.expr = CNode("for", yyDollar[2].expr, yyDollar[5].expr, yyDollar[6].expr).setPos0(yyDollar[1].token)
		}
	case 44:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line .\parser.go.y:182
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
	case 45:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line .\parser.go.y:196
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
	case 46:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line .\parser.go.y:236
		{
			yyVAL.expr = CNode("call", "copy", CNode(
				NNode(0),
				yyDollar[6].expr,
				CNode("func", "<anony-map-iter-callback>", CNode(yyDollar[2].token.Str, yyDollar[4].token.Str), yyDollar[7].expr),
			))
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:245
		{
			yyVAL.expr = CNode("if", yyDollar[2].expr, yyDollar[3].expr, CNode())
		}
	case 48:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:246
		{
			yyVAL.expr = CNode("if", yyDollar[2].expr, yyDollar[3].expr, yyDollar[5].expr)
		}
	case 49:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:249
		{
			yyVAL.expr = expandSwitch(yyDollar[2].expr, yyDollar[4].expr.C())
		}
	case 50:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:252
		{
			yyVAL.expr = CNode(yyDollar[2].expr, yyDollar[4].expr)
		}
	case 51:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:253
		{
			yyVAL.expr = CNode(ANode(yyDollar[2].token), yyDollar[4].expr)
		}
	case 52:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:254
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[3].expr, yyDollar[5].expr)
		}
	case 53:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:255
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(ANode(yyDollar[3].token), yyDollar[5].expr)
		}
	case 54:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line .\parser.go.y:258
		{
			yyVAL.str = ""
		}
	case 55:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:259
		{
			yyVAL.str = yyDollar[1].str + ",safe"
		}
	case 56:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:260
		{
			yyVAL.str = yyDollar[1].str + ",var"
		}
	case 57:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:263
		{
			yyVAL.str = "func," + yyDollar[2].str
		}
	case 58:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:266
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
	case 59:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:278
		{
			yyVAL.expr = CNode("yield", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 60:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:279
		{
			yyVAL.expr = CNode("break").setPos0(yyDollar[1].token)
		}
	case 61:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:280
		{
			yyVAL.expr = CNode("continue").setPos0(yyDollar[1].token)
		}
	case 62:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:281
		{
			yyVAL.expr = CNode("assert", yyDollar[2].expr, NilNode()).setPos0(yyDollar[1].token)
		}
	case 63:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:282
		{
			yyVAL.expr = CNode("assert", yyDollar[2].expr, SNode(yyDollar[3].token.Str)).setPos0(yyDollar[1].token)
		}
	case 64:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:283
		{
			yyVAL.expr = CNode("ret", yyDollar[2].expr).setPos0(yyDollar[1].token)
		}
	case 65:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:284
		{
			yyVAL.expr = yylex.(*Lexer).loadFile(filepath.Join(filepath.Dir(yyDollar[1].token.Pos.Source), yyDollar[2].token.Str))
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:287
		{
			yyVAL.expr = ANode(yyDollar[1].token).setPos(yyDollar[1].token)
		}
	case 67:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:288
		{
			yyVAL.expr = CNode("load", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 68:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line .\parser.go.y:289
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, yyDollar[3].expr, yyDollar[5].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 69:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:290
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, yyDollar[3].expr, NNode("-1")).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 70:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:291
		{
			yyVAL.expr = CNode("slice", yyDollar[1].expr, NNode("0"), yyDollar[4].expr).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:292
		{
			yyVAL.expr = CNode("load", yyDollar[1].expr, SNode(yyDollar[3].token.Str)).setPos0(yyDollar[1].expr).setPos(yyDollar[1].expr)
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:295
		{
			yyVAL.expr = CNode(yyDollar[1].token.Str)
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:296
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(ANode(yyDollar[3].token))
		}
	case 74:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:299
		{
			yyVAL.expr = CNode(yyDollar[1].expr)
		}
	case 75:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:300
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[3].expr)
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:303
		{
			yyVAL.expr = CNode(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 77:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:304
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(yyDollar[3].expr).Cappend(yyDollar[5].expr)
		}
	case 78:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:307
		{
			yyVAL.expr = CNode("chain", CNode("set", ANode(yyDollar[1].token), NilNode()).setPos0(yyDollar[1].token))
		}
	case 79:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:308
		{
			yyVAL.expr = CNode("chain", CNode("set", ANode(yyDollar[1].token), yyDollar[3].expr).setPos0(yyDollar[1].token))
		}
	case 80:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line .\parser.go.y:309
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(CNode("set", ANode(yyDollar[3].token), yyDollar[5].expr).setPos0(yyDollar[1].expr))
		}
	case 81:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:310
		{
			yyVAL.expr = yyDollar[1].expr.Cappend(CNode("set", ANode(yyDollar[3].token), NilNode()).setPos0(yyDollar[1].expr))
		}
	case 82:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:313
		{
			yyVAL.expr = NilNode().SetPos(yyDollar[1].token)
		}
	case 83:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:314
		{
			yyVAL.expr = NNode(yyDollar[1].token.Str).SetPos(yyDollar[1].token)
		}
	case 84:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:315
		{
			yyVAL.expr = yylex.(*Lexer).loadFile(filepath.Join(filepath.Dir(yyDollar[1].token.Pos.Source), yyDollar[2].token.Str))
		}
	case 85:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:316
		{
			yyVAL.expr = CNode("typeof", yyDollar[2].expr)
		}
	case 86:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:317
		{
			yyVAL.expr = CNode("len", yyDollar[2].expr)
		}
	case 87:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:318
		{
			yyVAL.expr = CNode("call", "addressof", CNode(ANode(yyDollar[2].token)))
		}
	case 88:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:319
		{
			yyVAL.expr = CNode("call", "copy", CNode(NNode(1), yyDollar[2].expr, NilNode()))
		}
	case 89:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:320
		{
			yyVAL.expr = CNode("call", "copy", CNode(NNode(1), yyDollar[2].expr, yyDollar[4].expr))
		}
	case 90:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:321
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 91:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:322
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 92:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:323
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 93:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:324
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 94:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:325
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 95:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:326
		{
			yyVAL.expr = CNode("or", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 96:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:327
		{
			yyVAL.expr = CNode("and", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 97:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:328
		{
			yyVAL.expr = CNode("<", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 98:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:329
		{
			yyVAL.expr = CNode("<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 99:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:330
		{
			yyVAL.expr = CNode("<=", yyDollar[3].expr, yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 100:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:331
		{
			yyVAL.expr = CNode("<=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 101:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:332
		{
			yyVAL.expr = CNode("==", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 102:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:333
		{
			yyVAL.expr = CNode("!=", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 103:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:334
		{
			yyVAL.expr = CNode("+", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 104:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:335
		{
			yyVAL.expr = CNode("-", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 105:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:336
		{
			yyVAL.expr = CNode("*", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 106:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:337
		{
			yyVAL.expr = CNode("/", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 107:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:338
		{
			yyVAL.expr = CNode("%", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 108:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:339
		{
			yyVAL.expr = CNode("^", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 109:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:340
		{
			yyVAL.expr = CNode("<<", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 110:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:341
		{
			yyVAL.expr = CNode(">>", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 111:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:342
		{
			yyVAL.expr = CNode(">>>", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 112:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:343
		{
			yyVAL.expr = CNode("|", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 113:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:344
		{
			yyVAL.expr = CNode("&", yyDollar[1].expr, yyDollar[3].expr).setPos0(yyDollar[1].expr)
		}
	case 114:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:345
		{
			yyVAL.expr = CNode("-", NNode(0.0), yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 115:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:346
		{
			yyVAL.expr = CNode("~", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 116:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:347
		{
			yyVAL.expr = CNode("!", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 117:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:348
		{
			yyVAL.expr = CNode("#", yyDollar[2].expr).setPos0(yyDollar[2].expr)
		}
	case 118:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:351
		{
			yyVAL.expr = SNode(yyDollar[1].token.Str).SetPos(yyDollar[1].token)
		}
	case 119:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:354
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 120:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:355
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 121:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:356
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 122:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:357
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 123:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:360
		{
			yyVAL.expr = CNode("call", yyDollar[1].expr, yyDollar[2].expr).setPos0(yyDollar[1].expr)
		}
	case 124:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:365
		{
			yyVAL.expr = CNode()
		}
	case 125:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:366
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 126:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:369
		{
			yyVAL.expr = CNode(yyDollar[1].str, "<a>", yyDollar[2].expr, yyDollar[3].expr).setPos0(yyDollar[2].expr)
		}
	case 127:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line .\parser.go.y:370
		{
			yyVAL.expr = CNode(yyDollar[1].str, "<a>", yyDollar[2].expr, CNode("chain", CNode("ret", yyDollar[4].expr))).setPos0(yyDollar[2].expr)
		}
	case 128:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:371
		{
			yyVAL.expr = CNode(yyDollar[1].str, "<a>", CNode(), CNode("chain", CNode("ret", yyDollar[3].expr))).setPos0(yyDollar[3].expr)
		}
	case 129:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:374
		{
			yyVAL.expr = CNode()
		}
	case 130:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:375
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 131:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:378
		{
			yyVAL.expr = CNode("map", CNode()).setPos0(yyDollar[1].token)
		}
	case 132:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line .\parser.go.y:379
		{
			yyVAL.expr = yyDollar[2].expr.setPos0(yyDollar[1].token)
		}
	case 133:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:382
		{
			yyVAL.expr = CNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 134:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:383
		{
			yyVAL.expr = CNode("map", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 135:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line .\parser.go.y:384
		{
			yyVAL.expr = CNode("array", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	case 136:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line .\parser.go.y:385
		{
			yyVAL.expr = CNode("array", yyDollar[1].expr).setPos0(yyDollar[1].expr)
		}
	}
	goto yystack /* stack new state and value */
}
