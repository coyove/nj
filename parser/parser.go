//line parser.go.y:1
package parser

import __yyfmt__ "fmt"

//line parser.go.y:3
import "github.com/coyove/bracket/vm"

//line parser.go.y:20
type yySymType struct {
	yys   int
	token Token

	stmts *Node
	stmt  *Node

	funcname interface{}
	funcexpr interface{}

	exprlist *Node
	expr     *Node

	namelist *Node
}

const TAnd = 57346
const TAssert = 57347
const TBreak = 57348
const TContinue = 57349
const TDo = 57350
const TElse = 57351
const TElseIf = 57352
const TEnd = 57353
const TFalse = 57354
const TFor = 57355
const TIf = 57356
const TIn = 57357
const TLambda = 57358
const TNil = 57359
const TNot = 57360
const TOr = 57361
const TReturn = 57362
const TSet = 57363
const TThen = 57364
const TTrue = 57365
const TWhile = 57366
const TXor = 57367
const TEqeq = 57368
const TNeq = 57369
const TLsh = 57370
const TRsh = 57371
const TLte = 57372
const TGte = 57373
const TIdent = 57374
const TNumber = 57375
const TString = 57376
const UNARY = 57377

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"TAnd",
	"TAssert",
	"TBreak",
	"TContinue",
	"TDo",
	"TElse",
	"TElseIf",
	"TEnd",
	"TFalse",
	"TFor",
	"TIf",
	"TIn",
	"TLambda",
	"TNil",
	"TNot",
	"TOr",
	"TReturn",
	"TSet",
	"TThen",
	"TTrue",
	"TWhile",
	"TXor",
	"TEqeq",
	"TNeq",
	"TLsh",
	"TRsh",
	"TLte",
	"TGte",
	"TIdent",
	"TNumber",
	"TString",
	"'{'",
	"'('",
	"'|'",
	"'&'",
	"'>'",
	"'<'",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'%'",
	"UNARY",
	"'^'",
	"';'",
	"'='",
	"'['",
	"']'",
	"'.'",
	"','",
	"'~'",
	"')'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.go.y:348

func TokenName(c int) string {
	if c >= TAnd && c-TAnd < len(yyToknames) {
		if yyToknames[c-TAnd] != "" {
			return yyToknames[c-TAnd]
		}
	}
	return string([]byte{byte(c)})
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 520

var yyAct = [...]int{

	48, 28, 34, 5, 4, 37, 47, 22, 36, 75,
	39, 109, 79, 40, 80, 77, 79, 41, 43, 44,
	1, 17, 52, 57, 58, 65, 66, 56, 55, 72,
	105, 69, 70, 71, 67, 68, 53, 54, 59, 60,
	61, 62, 63, 74, 64, 45, 38, 75, 27, 38,
	20, 82, 83, 84, 85, 86, 87, 88, 89, 90,
	91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
	81, 24, 101, 21, 14, 33, 23, 32, 102, 29,
	106, 104, 25, 5, 4, 52, 15, 18, 110, 19,
	52, 13, 26, 35, 103, 16, 52, 67, 68, 65,
	66, 30, 2, 42, 0, 5, 4, 64, 67, 68,
	5, 4, 64, 31, 5, 4, 118, 0, 64, 5,
	4, 51, 108, 5, 4, 114, 115, 113, 0, 0,
	112, 0, 0, 0, 0, 117, 50, 0, 0, 0,
	0, 121, 52, 57, 58, 65, 66, 56, 55, 0,
	0, 0, 51, 0, 67, 68, 53, 54, 59, 60,
	61, 62, 63, 0, 64, 0, 0, 50, 0, 0,
	0, 0, 76, 52, 57, 58, 65, 66, 56, 55,
	0, 0, 51, 0, 0, 67, 68, 53, 54, 59,
	60, 61, 62, 63, 0, 64, 0, 50, 0, 78,
	120, 0, 0, 52, 57, 58, 65, 66, 56, 55,
	0, 51, 0, 0, 0, 67, 68, 53, 54, 59,
	60, 61, 62, 63, 0, 64, 50, 0, 0, 73,
	0, 0, 52, 57, 58, 65, 66, 56, 55, 0,
	0, 0, 0, 0, 67, 68, 53, 54, 59, 60,
	61, 62, 63, 51, 64, 0, 0, 49, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 50, 0,
	0, 0, 0, 0, 52, 57, 58, 65, 66, 56,
	55, 0, 0, 0, 0, 0, 67, 68, 53, 54,
	59, 60, 61, 62, 63, 24, 64, 0, 0, 33,
	23, 32, 0, 0, 0, 0, 25, 51, 0, 0,
	0, 0, 0, 0, 0, 13, 26, 35, 0, 16,
	0, 0, 50, 0, 0, 30, 0, 0, 52, 57,
	58, 65, 66, 56, 55, 51, 0, 31, 46, 0,
	67, 68, 53, 54, 59, 60, 61, 62, 63, 0,
	64, 0, 0, 0, 0, 0, 52, 57, 58, 65,
	66, 56, 55, 0, 0, 0, 0, 0, 67, 68,
	53, 54, 59, 60, 61, 62, 63, 0, 64, 12,
	10, 11, 0, 0, 0, 119, 0, 0, 7, 0,
	0, 0, 0, 0, 9, 8, 0, 0, 6, 12,
	10, 11, 0, 0, 0, 116, 13, 0, 7, 0,
	16, 0, 0, 0, 9, 8, 0, 0, 6, 12,
	10, 11, 3, 0, 0, 111, 13, 0, 7, 0,
	16, 0, 0, 0, 9, 8, 0, 0, 6, 12,
	10, 11, 3, 0, 0, 107, 13, 0, 7, 0,
	16, 0, 0, 0, 9, 8, 0, 0, 6, 12,
	10, 11, 3, 0, 0, 0, 13, 0, 7, 0,
	16, 0, 0, 0, 9, 8, 52, 0, 6, 65,
	66, 0, 3, 0, 0, 0, 13, 0, 67, 68,
	16, 0, 59, 60, 61, 62, 63, 52, 64, 0,
	65, 66, 3, 0, 0, 0, 0, 0, 0, 67,
	68, 0, 0, 0, 0, 61, 62, 63, 0, 64,
}
var yyPact = [...]int{

	-1000, 454, -1000, -1000, -28, 37, 59, 59, 14, 59,
	-1000, -1000, 59, -1000, -1000, -1000, 59, 59, 59, 13,
	-1000, 283, 249, -1000, -1000, -1000, -1000, -1000, 37, -1000,
	59, 59, 59, -7, -1000, -1000, 207, -6, -1000, 303,
	303, 117, -40, 303, 148, -1000, -1000, -41, 303, -1000,
	59, 59, 59, 59, 59, 59, 59, 59, 59, 59,
	59, 59, 59, 59, 59, 59, 59, 59, 59, 65,
	65, 65, 17, -1000, 59, -2, -1000, -1000, -1000, 59,
	-1000, 434, 331, -3, 303, 451, 451, 451, 451, 451,
	451, 472, 472, 71, 71, 71, 65, 60, 60, 65,
	65, -1000, -44, 454, -37, -1000, 303, -1000, 414, -1000,
	116, -1000, 394, -1000, -1000, 59, -1000, 374, 178, -1000,
	-1000, 454,
}
var yyPgo = [...]int{

	0, 20, 102, 88, 2, 5, 6, 0, 79, 1,
	86, 74, 50, 48,
}
var yyR1 = [...]int{

	0, 1, 1, 1, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 3, 3, 4, 4, 4,
	5, 5, 6, 6, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 8, 9, 9, 9, 9, 11, 10,
	12, 12, 13, 13,
}
var yyR2 = [...]int{

	0, 0, 2, 2, 3, 1, 5, 6, 8, 4,
	1, 2, 1, 1, 2, 0, 5, 1, 4, 3,
	1, 3, 1, 3, 1, 1, 1, 1, 1, 1,
	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	2, 2, 2, 1, 1, 1, 1, 3, 3, 2,
	2, 3, 5, 6,
}
var yyChk = [...]int{

	-1000, -1, -2, 48, -4, -9, 24, 14, 21, 20,
	6, 7, 5, 32, -11, -10, 36, 49, 50, 52,
	-12, 36, -7, 17, 12, 23, 33, -13, -9, -8,
	42, 54, 18, 16, -4, 34, -7, -5, 32, -7,
	-7, -7, -10, -7, -7, 32, 55, -6, -7, 8,
	19, 4, 25, 39, 40, 31, 30, 26, 27, 41,
	42, 43, 44, 45, 47, 28, 29, 37, 38, -7,
	-7, -7, 36, 22, 49, 53, 55, 55, 51, 53,
	55, -1, -7, -7, -7, -7, -7, -7, -7, -7,
	-7, -7, -7, -7, -7, -7, -7, -7, -7, -7,
	-7, 55, -5, -1, -6, 32, -7, 11, -1, 55,
	-3, 11, -1, 11, 9, 10, 11, -1, -7, 11,
	22, -1,
}
var yyDef = [...]int{

	1, -2, 2, 3, 54, 5, 0, 0, 0, 10,
	12, 13, 0, 17, 55, 56, 0, 0, 0, 0,
	59, 0, 0, 24, 25, 26, 27, 28, 29, 30,
	0, 0, 0, 0, 54, 53, 0, 0, 20, 11,
	14, 0, 56, 4, 0, 19, 60, 0, 22, 1,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 50,
	51, 52, 0, 1, 0, 0, 57, 58, 18, 0,
	61, 0, 31, 32, 33, 34, 35, 36, 37, 38,
	39, 40, 41, 42, 43, 44, 45, 46, 47, 48,
	49, 1, 0, 15, 9, 21, 23, 6, 0, 1,
	0, 62, 0, 7, 1, 0, 63, 0, 0, 8,
	1, 16,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 45, 38, 3,
	36, 55, 43, 41, 53, 42, 52, 44, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 48,
	40, 49, 39, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 50, 3, 51, 47, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 35, 37, 3, 54,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 46,
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
			yyVAL.stmts = NewCompoundNode("chain")
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.stmts
			}
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:61
		{
			yyDollar[1].stmts.Compound = append(yyDollar[1].stmts.Compound, yyDollar[2].stmt)
			yyVAL.stmts = yyDollar[1].stmts
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.stmts
			}
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:68
		{
			yyVAL.stmts = yyDollar[1].stmts
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.stmts
			}
		}
	case 4:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:76
		{
			if len(yyDollar[1].expr.Compound) > 0 && yyDollar[1].expr.Compound[0].Value.(string) == "load" {
				yyVAL.stmt = NewCompoundNode("store", yyDollar[1].expr.Compound[1], yyDollar[1].expr.Compound[2], yyDollar[3].expr)
			} else {
				yyVAL.stmt = NewCompoundNode("move", yyDollar[1].expr, yyDollar[3].expr)
			}
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:84
		{
			// if _, ok := $1.(*FuncCallExpr); !ok {
			//    yylex.(*Lexer).Error("parse error")
			// } else {
			yyVAL.stmt = yyDollar[1].expr
			// }
		}
	case 6:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:91
		{
			yyVAL.stmt = NewCompoundNode("while", yyDollar[2].expr, yyDollar[4].stmts)
		}
	case 7:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:94
		{
			yyVAL.stmt = NewCompoundNode("if", yyDollar[2].expr, yyDollar[4].stmts, NewCompoundNode())
			cur := yyVAL.stmt
			for _, e := range yyDollar[5].stmts.Compound {
				cur.Compound[3] = e
				cur = e
			}
		}
	case 8:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.go.y:102
		{
			yyVAL.stmt = NewCompoundNode("if", yyDollar[2].expr, yyDollar[4].stmts, NewCompoundNode())
			cur := yyVAL.stmt
			for _, e := range yyDollar[5].stmts.Compound {
				cur.Compound[3] = e
				cur = e
			}
			cur.Compound[3] = yyDollar[7].stmts
		}
	case 9:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:111
		{
			yyVAL.stmt = NewCompoundNode("chain")
			for i, name := range yyDollar[2].namelist.Compound {
				var e *Node
				if i < len(yyDollar[4].exprlist.Compound) {
					e = yyDollar[4].exprlist.Compound[i]
				} else {
					e = yyDollar[4].exprlist.Compound[len(yyDollar[4].exprlist.Compound)-1]
				}
				yyVAL.stmt.Compound = append(yyVAL.stmt.Compound, NewCompoundNode("set", name, e))
			}
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:123
		{
			yyVAL.stmt = NewCompoundNode("ret")
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:126
		{
			yyVAL.stmt = NewCompoundNode("ret", yyDollar[2].expr)
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:129
		{
			yyVAL.stmt = NewCompoundNode("break")
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:132
		{
			yyVAL.stmt = NewCompoundNode("continue")
		}
	case 14:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:135
		{
			yyVAL.stmt = NewCompoundNode("assert", yyDollar[2].expr)
			yyVAL.stmt.Pos = yyDollar[2].expr.Pos
		}
	case 15:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.go.y:141
		{
			yyVAL.stmts = NewCompoundNode()
		}
	case 16:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:144
		{
			yyVAL.stmts.Compound = append(yyVAL.stmts.Compound, NewCompoundNode("if", yyDollar[3].expr, yyDollar[5].stmts, NewCompoundNode()))
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:149
		{
			yyVAL.expr = NewAtomNode(yyDollar[1].token)
			_, yyVAL.expr.LibWH = vm.LibLookup[yyDollar[1].token.Str]
		}
	case 18:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:153
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, yyDollar[3].expr)
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:156
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, NewStringNode(yyDollar[3].token.Str))
		}
	case 20:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:161
		{
			yyVAL.namelist = NewCompoundNode(yyDollar[1].token.Str)
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:164
		{
			yyDollar[1].namelist.Compound = append(yyDollar[1].namelist.Compound, NewAtomNode(yyDollar[3].token))
			yyVAL.namelist = yyDollar[1].namelist
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:170
		{
			yyVAL.exprlist = NewCompoundNode(yyDollar[1].expr)
			yyVAL.exprlist.Pos = yyDollar[1].expr.Pos
		}
	case 23:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:174
		{
			yyDollar[1].exprlist.Compound = append(yyDollar[1].exprlist.Compound, yyDollar[3].expr)
			yyVAL.exprlist = yyDollar[1].exprlist
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:180
		{
			yyVAL.expr = NewCompoundNode("nil")
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:184
		{
			yyVAL.expr = NewCompoundNode("false")
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:188
		{
			yyVAL.expr = NewCompoundNode("true")
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:192
		{
			yyVAL.expr = NewNumberNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:196
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:199
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:202
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 31:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:205
		{
			yyVAL.expr = NewCompoundNode("or", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 32:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:209
		{
			yyVAL.expr = NewCompoundNode("and", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 33:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:213
		{
			yyVAL.expr = NewCompoundNode("xor", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 34:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:217
		{
			yyVAL.expr = NewCompoundNode(">", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 35:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:221
		{
			yyVAL.expr = NewCompoundNode("<", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 36:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:225
		{
			yyVAL.expr = NewCompoundNode(">=", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 37:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:229
		{
			yyVAL.expr = NewCompoundNode("<=", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:233
		{
			yyVAL.expr = NewCompoundNode("eq", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:237
		{
			yyVAL.expr = NewCompoundNode("neq", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:241
		{
			yyVAL.expr = NewCompoundNode("+", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 41:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:245
		{
			yyVAL.expr = NewCompoundNode("-", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:249
		{
			yyVAL.expr = NewCompoundNode("*", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:253
		{
			yyVAL.expr = NewCompoundNode("/", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 44:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:257
		{
			yyVAL.expr = NewCompoundNode("%", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:261
		{
			yyVAL.expr = NewCompoundNode("^", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 46:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:265
		{
			yyVAL.expr = NewCompoundNode("<<", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:269
		{
			yyVAL.expr = NewCompoundNode(">>", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 48:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:273
		{
			yyVAL.expr = NewCompoundNode("|", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 49:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:277
		{
			yyVAL.expr = NewCompoundNode("&", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 50:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:281
		{
			yyVAL.expr = NewCompoundNode("-", NewNumberNode("0"), yyDollar[2].expr)
			yyVAL.expr.Pos = yyDollar[2].expr.Pos
		}
	case 51:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:285
		{
			yyVAL.expr = NewCompoundNode("~", yyDollar[2].expr)
			yyVAL.expr.Pos = yyDollar[2].expr.Pos
		}
	case 52:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:289
		{
			yyVAL.expr = NewCompoundNode("not", yyDollar[2].expr)
			yyVAL.expr.Pos = yyDollar[2].expr.Pos
		}
	case 53:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:295
		{
			yyVAL.expr = NewStringNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:301
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 55:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:304
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 56:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:307
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:310
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:315
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 59:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:320
		{
			yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, yyDollar[2].exprlist)
		}
	case 60:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:325
		{
			if yylex.(*Lexer).PNewLine {
				yylex.(*Lexer).TokenError(yyDollar[1].token, "ambiguous syntax (function call x new statement)")
			}
			yyVAL.exprlist = NewCompoundNode()
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:331
		{
			if yylex.(*Lexer).PNewLine {
				yylex.(*Lexer).TokenError(yyDollar[1].token, "ambiguous syntax (function call x new statement)")
			}
			yyVAL.exprlist = yyDollar[2].exprlist
		}
	case 62:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:339
		{
			yyVAL.expr = NewCompoundNode("lambda", NewCompoundNode(), yyDollar[4].stmts)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 63:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:343
		{
			yyVAL.expr = NewCompoundNode("lambda", yyDollar[3].namelist, yyDollar[5].stmts)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	}
	goto yystack /* stack new state and value */
}
