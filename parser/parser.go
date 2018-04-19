//line parser.go.y:1
package parser

import __yyfmt__ "fmt"

//line parser.go.y:3
//line parser.go.y:18
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
const TBreak = 57347
const TContinue = 57348
const TDo = 57349
const TElse = 57350
const TElseIf = 57351
const TEnd = 57352
const TFalse = 57353
const TFor = 57354
const TFunction = 57355
const TIf = 57356
const TIn = 57357
const TNil = 57358
const TNot = 57359
const TOr = 57360
const TReturn = 57361
const TSet = 57362
const TThen = 57363
const TTrue = 57364
const TWhile = 57365
const TXor = 57366
const TEqeq = 57367
const TNeq = 57368
const TLte = 57369
const TGte = 57370
const TIdent = 57371
const TNumber = 57372
const TString = 57373
const UNARY = 57374

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"TAnd",
	"TBreak",
	"TContinue",
	"TDo",
	"TElse",
	"TElseIf",
	"TEnd",
	"TFalse",
	"TFor",
	"TFunction",
	"TIf",
	"TIn",
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
	"TLte",
	"TGte",
	"TIdent",
	"TNumber",
	"TString",
	"'{'",
	"'('",
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
	"')'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.go.y:306

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

const yyLast = 413

var yyAct = [...]int{

	21, 27, 32, 5, 4, 47, 95, 96, 34, 67,
	36, 65, 89, 69, 70, 16, 37, 39, 40, 46,
	63, 44, 1, 99, 48, 48, 53, 54, 52, 51,
	61, 62, 87, 41, 35, 49, 50, 55, 56, 57,
	58, 59, 60, 60, 26, 19, 13, 72, 73, 74,
	75, 76, 77, 78, 79, 80, 81, 82, 83, 84,
	85, 86, 10, 11, 28, 43, 91, 107, 71, 20,
	92, 7, 14, 5, 4, 88, 9, 8, 97, 2,
	6, 17, 0, 18, 0, 0, 12, 90, 38, 0,
	15, 0, 5, 4, 0, 0, 5, 4, 0, 48,
	3, 0, 5, 4, 106, 0, 47, 5, 4, 0,
	94, 5, 4, 57, 58, 59, 0, 60, 0, 100,
	46, 102, 103, 101, 0, 105, 48, 53, 54, 52,
	51, 109, 47, 0, 0, 0, 49, 50, 55, 56,
	57, 58, 59, 0, 60, 0, 46, 0, 0, 0,
	0, 66, 48, 53, 54, 52, 51, 0, 0, 0,
	0, 0, 49, 50, 55, 56, 57, 58, 59, 0,
	60, 23, 0, 31, 68, 0, 22, 30, 0, 0,
	0, 47, 24, 0, 0, 0, 0, 0, 0, 12,
	25, 33, 0, 15, 0, 46, 0, 29, 108, 0,
	0, 48, 53, 54, 52, 51, 47, 0, 0, 42,
	0, 49, 50, 55, 56, 57, 58, 59, 0, 60,
	46, 0, 0, 64, 0, 0, 48, 53, 54, 52,
	51, 47, 0, 0, 45, 0, 49, 50, 55, 56,
	57, 58, 59, 0, 60, 46, 0, 0, 0, 0,
	47, 48, 53, 54, 52, 51, 0, 0, 0, 0,
	0, 49, 50, 55, 56, 57, 58, 59, 0, 60,
	48, 53, 54, 52, 51, 0, 0, 0, 0, 0,
	49, 50, 55, 56, 57, 58, 59, 0, 60, 48,
	53, 54, 52, 51, 0, 0, 0, 0, 0, 49,
	50, 55, 56, 57, 58, 59, 0, 60, 10, 11,
	10, 11, 0, 104, 0, 98, 0, 7, 0, 7,
	0, 0, 9, 8, 9, 8, 6, 0, 6, 0,
	0, 0, 12, 0, 12, 0, 15, 0, 15, 0,
	10, 11, 10, 11, 0, 93, 3, 0, 3, 7,
	0, 7, 0, 0, 9, 8, 9, 8, 6, 0,
	6, 0, 0, 0, 12, 0, 12, 0, 15, 48,
	15, 0, 0, 0, 0, 0, 0, 0, 3, 0,
	3, 55, 56, 57, 58, 59, 23, 60, 31, 0,
	0, 22, 30, 0, 0, 0, 0, 24, 0, 0,
	0, 0, 0, 0, 12, 25, 33, 0, 15, 0,
	0, 0, 29,
}
var yyPact = [...]int{

	-1000, 337, -1000, -1000, -29, 36, 375, 375, 5, 375,
	-1000, -1000, -1000, -1000, -1000, 375, 375, 375, 4, -1000,
	160, 227, -1000, -1000, -1000, -1000, -1000, 36, -1000, 375,
	375, -13, -1000, -1000, 202, -33, 1, 102, -40, 1,
	128, -1000, -1000, -35, 1, -1000, 375, 375, 375, 375,
	375, 375, 375, 375, 375, 375, 375, 375, 375, 375,
	375, 0, 0, -17, -1000, 375, -1000, -1000, -1000, 375,
	-1000, 335, 246, 265, 1, 345, 345, 345, 345, 345,
	345, 75, 75, 0, 0, 0, 0, -1000, -42, -1000,
	337, 1, 1, -1000, 305, -6, -1000, 113, -1000, -1000,
	303, -1000, -1000, 375, -1000, 57, 177, -1000, -1000, 337,
}
var yyPgo = [...]int{

	0, 22, 79, 78, 2, 75, 65, 0, 64, 1,
	72, 46, 45, 44,
}
var yyR1 = [...]int{

	0, 1, 1, 1, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 3, 3, 4, 4, 4, 5,
	5, 6, 6, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 8, 9, 9,
	9, 9, 11, 10, 12, 12, 13, 13,
}
var yyR2 = [...]int{

	0, 0, 2, 2, 3, 1, 5, 6, 8, 4,
	1, 2, 1, 1, 0, 5, 1, 4, 3, 1,
	3, 1, 3, 1, 1, 1, 1, 1, 1, 1,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 2, 2, 1, 1, 1,
	1, 3, 3, 2, 2, 3, 5, 6,
}
var yyChk = [...]int{

	-1000, -1, -2, 43, -4, -9, 23, 14, 20, 19,
	5, 6, 29, -11, -10, 33, 44, 45, 47, -12,
	33, -7, 16, 11, 22, 30, -13, -9, -8, 37,
	17, 13, -4, 31, -7, 29, -7, -7, -10, -7,
	-7, 29, 49, -6, -7, 7, 18, 4, 24, 34,
	35, 28, 27, 25, 26, 36, 37, 38, 39, 40,
	42, -7, -7, 33, 21, 44, 49, 49, 46, 48,
	49, -1, -7, -7, -7, -7, -7, -7, -7, -7,
	-7, -7, -7, -7, -7, -7, -7, 49, -5, 29,
	-1, -7, -7, 10, -1, 48, 49, -3, 10, 29,
	-1, 10, 8, 9, 10, -1, -7, 10, 21, -1,
}
var yyDef = [...]int{

	1, -2, 2, 3, 48, 5, 0, 0, 0, 10,
	12, 13, 16, 49, 50, 0, 0, 0, 0, 53,
	0, 0, 23, 24, 25, 26, 27, 28, 29, 0,
	0, 0, 48, 47, 0, 0, 11, 0, 50, 4,
	0, 18, 54, 0, 21, 1, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 45, 46, 0, 1, 0, 51, 52, 17, 0,
	55, 0, 30, 31, 32, 33, 34, 35, 36, 37,
	38, 39, 40, 41, 42, 43, 44, 1, 0, 19,
	14, 9, 22, 6, 0, 0, 1, 0, 56, 20,
	0, 7, 1, 0, 57, 0, 0, 8, 1, 15,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 40, 3, 3,
	33, 49, 38, 36, 48, 37, 47, 39, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 43,
	35, 44, 34, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 45, 3, 46, 42, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 32,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	41,
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
		//line parser.go.y:51
		{
			yyVAL.stmts = NewCompoundNode("chain")
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.stmts
			}
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:57
		{
			yyDollar[1].stmts.Compound = append(yyDollar[1].stmts.Compound, yyDollar[2].stmt)
			yyVAL.stmts = yyDollar[1].stmts
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.stmts
			}
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:64
		{
			yyVAL.stmts = yyDollar[1].stmts
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.stmts
			}
		}
	case 4:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:72
		{
			yyVAL.stmt = NewCompoundNode("move", yyDollar[1].expr, yyDollar[3].expr)
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:76
		{
			// if _, ok := $1.(*FuncCallExpr); !ok {
			//    yylex.(*Lexer).Error("parse error")
			// } else {
			yyVAL.stmt = yyDollar[1].expr
			// }
		}
	case 6:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:83
		{
			yyVAL.stmt = NewCompoundNode("while", yyDollar[2].expr, yyDollar[4].stmts)
		}
	case 7:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:86
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
		//line parser.go.y:94
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
		//line parser.go.y:103
		{
			yyVAL.stmt = NewCompoundNode("set", yyDollar[2].token.Str, yyDollar[4].expr)
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:106
		{
			yyVAL.stmt = NewCompoundNode("ret")
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:109
		{
			yyVAL.stmt = NewCompoundNode("ret", yyDollar[2].expr)
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:112
		{
			yyVAL.stmt = NewCompoundNode("break")
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:115
		{
			yyVAL.stmt = NewCompoundNode("continue")
		}
	case 14:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.go.y:120
		{
			yyVAL.stmts = NewCompoundNode()
		}
	case 15:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:123
		{
			yyVAL.stmts.Compound = append(yyVAL.stmts.Compound, NewCompoundNode("if", yyDollar[3].expr, yyDollar[5].stmts, NewCompoundNode()))
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:128
		{
			yyVAL.expr = NewAtomNode(yyDollar[1].token)
		}
	case 17:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:131
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].expr, ":", yyDollar[3].expr)
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:134
		{
			yyVAL.expr = NewCompoundNode(yyDollar[1].expr, ":", yyDollar[3].token.Str)
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:139
		{
			yyVAL.namelist = NewCompoundNode(yyDollar[1].token.Str)
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:142
		{
			yyDollar[1].namelist.Compound = append(yyDollar[1].namelist.Compound, NewAtomNode(yyDollar[3].token))
			yyVAL.namelist = yyDollar[1].namelist
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:148
		{
			yyVAL.exprlist = NewCompoundNode(yyDollar[1].expr)
			yyVAL.exprlist.Pos = yyDollar[1].expr.Pos
		}
	case 22:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:152
		{
			yyDollar[1].exprlist.Compound = append(yyDollar[1].exprlist.Compound, yyDollar[3].expr)
			yyVAL.exprlist = yyDollar[1].exprlist
		}
	case 23:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:158
		{
			yyVAL.expr = NewCompoundNode("nil")
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:162
		{
			yyVAL.expr = NewCompoundNode("false")
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:166
		{
			yyVAL.expr = NewCompoundNode("true")
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:170
		{
			yyVAL.expr = NewNumberNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:174
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:177
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:180
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 30:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:183
		{
			yyVAL.expr = NewCompoundNode("or", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 31:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:187
		{
			yyVAL.expr = NewCompoundNode("and", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 32:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:191
		{
			yyVAL.expr = NewCompoundNode("xor", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 33:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:195
		{
			yyVAL.expr = NewCompoundNode(">", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 34:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:199
		{
			yyVAL.expr = NewCompoundNode("<", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 35:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:203
		{
			yyVAL.expr = NewCompoundNode(">=", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 36:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:207
		{
			yyVAL.expr = NewCompoundNode("<=", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 37:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:211
		{
			yyVAL.expr = NewCompoundNode("eq", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:215
		{
			yyVAL.expr = NewCompoundNode("neq", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:219
		{
			yyVAL.expr = NewCompoundNode("+", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:223
		{
			yyVAL.expr = NewCompoundNode("-", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 41:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:227
		{
			yyVAL.expr = NewCompoundNode("*", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:231
		{
			yyVAL.expr = NewCompoundNode("/", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:235
		{
			yyVAL.expr = NewCompoundNode("%", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 44:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:239
		{
			yyVAL.expr = NewCompoundNode("^", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 45:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:243
		{
			yyVAL.expr = NewCompoundNode("-", yyDollar[2].expr)
			yyVAL.expr.Pos = yyDollar[2].expr.Pos
		}
	case 46:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:247
		{
			yyVAL.expr = NewCompoundNode("not", yyDollar[2].expr)
			yyVAL.expr.Pos = yyDollar[2].expr.Pos
		}
	case 47:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:253
		{
			yyVAL.expr = NewStringNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 48:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:259
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 49:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:262
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 50:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:265
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:268
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:273
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 53:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:278
		{
			yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, yyDollar[2].exprlist)
		}
	case 54:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:283
		{
			if yylex.(*Lexer).PNewLine {
				yylex.(*Lexer).TokenError(yyDollar[1].token, "ambiguous syntax (function call x new statement)")
			}
			yyVAL.exprlist = NewCompoundNode()
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:289
		{
			if yylex.(*Lexer).PNewLine {
				yylex.(*Lexer).TokenError(yyDollar[1].token, "ambiguous syntax (function call x new statement)")
			}
			yyVAL.exprlist = yyDollar[2].exprlist
		}
	case 56:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:297
		{
			yyVAL.expr = NewCompoundNode("lambda", NewCompoundNode(), yyDollar[4].stmts)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 57:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:301
		{
			yyVAL.expr = NewCompoundNode("lambda", yyDollar[3].namelist, yyDollar[5].stmts)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	}
	goto yystack /* stack new state and value */
}
