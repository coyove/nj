//line parser.go.y:2
package parser

import __yyfmt__ "fmt"

//line parser.go.y:2
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
const TAssert = 57347
const TBreak = 57348
const TContinue = 57349
const TDo = 57350
const TElse = 57351
const TElseIf = 57352
const TEnd = 57353
const TFalse = 57354
const TFor = 57355
const TFunction = 57356
const TIf = 57357
const TIn = 57358
const TNil = 57359
const TNot = 57360
const TOr = 57361
const TReturn = 57362
const TSet = 57363
const TThen = 57364
const TTrue = 57365
const TTypeIs = 57366
const TWhile = 57367
const TXor = 57368
const TEqeq = 57369
const TNeq = 57370
const TLte = 57371
const TGte = 57372
const TIdent = 57373
const TNumber = 57374
const TString = 57375
const UNARY = 57376

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
	"TTypeIs",
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
	"'~'",
	"')'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.go.y:322

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

const yyLast = 450

var yyAct = [...]int{

	22, 28, 34, 5, 4, 101, 72, 102, 36, 95,
	38, 70, 74, 39, 75, 17, 68, 40, 42, 43,
	105, 1, 47, 49, 77, 52, 57, 58, 56, 55,
	93, 65, 66, 67, 44, 53, 54, 59, 60, 61,
	62, 63, 37, 64, 108, 109, 107, 27, 20, 14,
	29, 78, 79, 80, 81, 82, 83, 84, 85, 86,
	87, 88, 89, 90, 91, 92, 12, 10, 11, 46,
	76, 97, 113, 21, 94, 98, 7, 15, 5, 4,
	103, 9, 8, 2, 0, 18, 6, 19, 0, 0,
	0, 96, 13, 0, 41, 0, 16, 0, 5, 4,
	0, 0, 5, 4, 0, 49, 3, 52, 5, 4,
	112, 0, 0, 5, 4, 100, 0, 5, 4, 51,
	0, 61, 62, 63, 106, 64, 0, 49, 0, 52,
	111, 0, 0, 0, 50, 0, 115, 0, 0, 49,
	0, 52, 57, 58, 56, 55, 51, 64, 0, 0,
	0, 53, 54, 59, 60, 61, 62, 63, 0, 64,
	0, 50, 0, 0, 0, 0, 49, 71, 52, 57,
	58, 56, 55, 0, 0, 0, 0, 0, 53, 54,
	59, 60, 61, 62, 63, 0, 64, 24, 0, 33,
	73, 0, 23, 32, 0, 0, 0, 0, 25, 0,
	0, 0, 0, 0, 0, 0, 13, 26, 35, 0,
	16, 24, 0, 33, 30, 0, 23, 32, 0, 0,
	0, 51, 25, 0, 0, 0, 31, 45, 0, 0,
	13, 26, 35, 0, 16, 0, 50, 0, 30, 114,
	0, 49, 0, 52, 57, 58, 56, 55, 51, 0,
	31, 0, 0, 53, 54, 59, 60, 61, 62, 63,
	0, 64, 0, 50, 0, 0, 69, 0, 49, 0,
	52, 57, 58, 56, 55, 51, 0, 0, 0, 48,
	53, 54, 59, 60, 61, 62, 63, 0, 64, 0,
	50, 0, 0, 0, 0, 49, 0, 52, 57, 58,
	56, 55, 51, 0, 0, 0, 0, 53, 54, 59,
	60, 61, 62, 63, 0, 64, 0, 50, 0, 0,
	0, 51, 49, 0, 52, 57, 58, 56, 55, 0,
	0, 0, 0, 0, 53, 54, 59, 60, 61, 62,
	63, 49, 64, 52, 57, 58, 56, 55, 0, 0,
	0, 0, 0, 53, 54, 59, 60, 61, 62, 63,
	0, 64, 12, 10, 11, 12, 10, 11, 110, 0,
	0, 104, 7, 0, 0, 7, 0, 9, 8, 0,
	9, 8, 6, 0, 0, 6, 0, 0, 13, 0,
	0, 13, 16, 0, 0, 16, 0, 12, 10, 11,
	0, 0, 3, 99, 0, 3, 0, 7, 0, 12,
	10, 11, 9, 8, 0, 0, 0, 6, 0, 7,
	0, 0, 0, 13, 9, 8, 49, 16, 52, 6,
	0, 0, 0, 0, 0, 13, 0, 3, 0, 16,
	59, 60, 61, 62, 63, 0, 64, 0, 0, 3,
}
var yyPact = [...]int{

	-1000, 404, -1000, -1000, -31, 38, 199, 199, 11, 199,
	-1000, -1000, 199, -1000, -1000, -1000, 199, 199, 199, 3,
	-1000, 175, 271, -1000, -1000, -1000, -1000, -1000, 38, -1000,
	199, 199, 199, -19, -1000, -1000, 244, -35, 298, 298,
	115, -46, 298, 142, -1000, -1000, -38, 298, -1000, -7,
	199, 199, 199, 199, 199, 199, 199, 199, 199, 199,
	199, 199, 199, 199, 199, 103, 103, 103, -22, -1000,
	199, -1000, -1000, -1000, 199, -1000, 392, -1000, 317, -1,
	298, 402, 402, 402, 402, 402, 402, 81, 81, 103,
	103, 103, 103, -1000, -45, -1000, 404, 298, 298, -1000,
	360, -11, -1000, 35, -1000, -1000, 357, -1000, -1000, 199,
	-1000, 61, 217, -1000, -1000, 404,
}
var yyPgo = [...]int{

	0, 21, 83, 80, 2, 74, 69, 0, 50, 1,
	77, 49, 48, 47,
}
var yyR1 = [...]int{

	0, 1, 1, 1, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 3, 3, 4, 4, 4,
	5, 5, 6, 6, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	8, 9, 9, 9, 9, 11, 10, 12, 12, 13,
	13,
}
var yyR2 = [...]int{

	0, 0, 2, 2, 3, 1, 5, 6, 8, 4,
	1, 2, 1, 1, 2, 0, 5, 1, 4, 3,
	1, 3, 1, 3, 1, 1, 1, 1, 1, 1,
	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 2, 2, 2,
	1, 1, 1, 1, 3, 3, 2, 2, 3, 5,
	6,
}
var yyChk = [...]int{

	-1000, -1, -2, 45, -4, -9, 25, 15, 21, 20,
	6, 7, 5, 31, -11, -10, 35, 46, 47, 49,
	-12, 35, -7, 17, 12, 23, 32, -13, -9, -8,
	39, 51, 18, 14, -4, 33, -7, 31, -7, -7,
	-7, -10, -7, -7, 31, 52, -6, -7, 8, 24,
	19, 4, 26, 36, 37, 30, 29, 27, 28, 38,
	39, 40, 41, 42, 44, -7, -7, -7, 35, 22,
	46, 52, 52, 48, 50, 52, -1, 31, -7, -7,
	-7, -7, -7, -7, -7, -7, -7, -7, -7, -7,
	-7, -7, -7, 52, -5, 31, -1, -7, -7, 11,
	-1, 50, 52, -3, 11, 31, -1, 11, 9, 10,
	11, -1, -7, 11, 22, -1,
}
var yyDef = [...]int{

	1, -2, 2, 3, 51, 5, 0, 0, 0, 10,
	12, 13, 0, 17, 52, 53, 0, 0, 0, 0,
	56, 0, 0, 24, 25, 26, 27, 28, 29, 30,
	0, 0, 0, 0, 51, 50, 0, 0, 11, 14,
	0, 53, 4, 0, 19, 57, 0, 22, 1, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 47, 48, 49, 0, 1,
	0, 54, 55, 18, 0, 58, 0, 31, 32, 33,
	34, 35, 36, 37, 38, 39, 40, 41, 42, 43,
	44, 45, 46, 1, 0, 20, 15, 9, 23, 6,
	0, 0, 1, 0, 59, 21, 0, 7, 1, 0,
	60, 0, 0, 8, 1, 16,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 42, 3, 3,
	35, 52, 40, 38, 50, 39, 49, 41, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 45,
	37, 46, 36, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 47, 3, 48, 44, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 34, 3, 3, 51,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 43,
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
			if len(yyDollar[1].expr.Compound) > 0 && yyDollar[1].expr.Compound[0].Value.(string) == "load" {
				yyVAL.stmt = NewCompoundNode("store", yyDollar[1].expr.Compound[1], yyDollar[1].expr.Compound[2], yyDollar[3].expr)
			} else {
				yyVAL.stmt = NewCompoundNode("move", yyDollar[1].expr, yyDollar[3].expr)
			}
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:80
		{
			// if _, ok := $1.(*FuncCallExpr); !ok {
			//    yylex.(*Lexer).Error("parse error")
			// } else {
			yyVAL.stmt = yyDollar[1].expr
			// }
		}
	case 6:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:87
		{
			yyVAL.stmt = NewCompoundNode("while", yyDollar[2].expr, yyDollar[4].stmts)
		}
	case 7:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:90
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
		//line parser.go.y:98
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
		//line parser.go.y:107
		{
			yyVAL.stmt = NewCompoundNode("set", yyDollar[2].token.Str, yyDollar[4].expr)
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:110
		{
			yyVAL.stmt = NewCompoundNode("ret")
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:113
		{
			yyVAL.stmt = NewCompoundNode("ret", yyDollar[2].expr)
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:116
		{
			yyVAL.stmt = NewCompoundNode("break")
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:119
		{
			yyVAL.stmt = NewCompoundNode("continue")
		}
	case 14:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:122
		{
			yyVAL.stmt = NewCompoundNode("assert", yyDollar[2].expr)
			yyVAL.stmt.Pos = yyDollar[2].expr.Pos
		}
	case 15:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.go.y:128
		{
			yyVAL.stmts = NewCompoundNode()
		}
	case 16:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:131
		{
			yyVAL.stmts.Compound = append(yyVAL.stmts.Compound, NewCompoundNode("if", yyDollar[3].expr, yyDollar[5].stmts, NewCompoundNode()))
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:136
		{
			yyVAL.expr = NewAtomNode(yyDollar[1].token)
		}
	case 18:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:139
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, yyDollar[3].expr)
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:142
		{
			yyVAL.expr = NewCompoundNode("load", yyDollar[1].expr, yyDollar[3].token.Str)
		}
	case 20:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:147
		{
			yyVAL.namelist = NewCompoundNode(yyDollar[1].token.Str)
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:150
		{
			yyDollar[1].namelist.Compound = append(yyDollar[1].namelist.Compound, NewAtomNode(yyDollar[3].token))
			yyVAL.namelist = yyDollar[1].namelist
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:156
		{
			yyVAL.exprlist = NewCompoundNode(yyDollar[1].expr)
			yyVAL.exprlist.Pos = yyDollar[1].expr.Pos
		}
	case 23:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:160
		{
			yyDollar[1].exprlist.Compound = append(yyDollar[1].exprlist.Compound, yyDollar[3].expr)
			yyVAL.exprlist = yyDollar[1].exprlist
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:166
		{
			yyVAL.expr = NewCompoundNode("nil")
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:170
		{
			yyVAL.expr = NewCompoundNode("false")
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:174
		{
			yyVAL.expr = NewCompoundNode("true")
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:178
		{
			yyVAL.expr = NewNumberNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:182
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:185
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:188
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 31:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:191
		{
			yyVAL.expr = NewCompoundNode("typeof", yyDollar[1].expr, yyDollar[3].token)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 32:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:195
		{
			yyVAL.expr = NewCompoundNode("or", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 33:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:199
		{
			yyVAL.expr = NewCompoundNode("and", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 34:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:203
		{
			yyVAL.expr = NewCompoundNode("xor", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 35:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:207
		{
			yyVAL.expr = NewCompoundNode(">", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 36:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:211
		{
			yyVAL.expr = NewCompoundNode("<", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 37:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:215
		{
			yyVAL.expr = NewCompoundNode(">=", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:219
		{
			yyVAL.expr = NewCompoundNode("<=", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:223
		{
			yyVAL.expr = NewCompoundNode("eq", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:227
		{
			yyVAL.expr = NewCompoundNode("neq", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 41:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:231
		{
			yyVAL.expr = NewCompoundNode("+", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:235
		{
			yyVAL.expr = NewCompoundNode("-", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:239
		{
			yyVAL.expr = NewCompoundNode("*", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 44:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:243
		{
			yyVAL.expr = NewCompoundNode("/", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:247
		{
			yyVAL.expr = NewCompoundNode("%", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 46:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:251
		{
			yyVAL.expr = NewCompoundNode("^", yyDollar[1].expr, yyDollar[3].expr)
			yyVAL.expr.Pos = yyDollar[1].expr.Pos
		}
	case 47:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:255
		{
			yyVAL.expr = NewCompoundNode("-", yyDollar[2].expr)
			yyVAL.expr.Pos = yyDollar[2].expr.Pos
		}
	case 48:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:259
		{
			yyVAL.expr = NewCompoundNode("~", yyDollar[2].expr)
			yyVAL.expr.Pos = yyDollar[2].expr.Pos
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:263
		{
			yyVAL.expr = NewCompoundNode("not", yyDollar[2].expr)
			yyVAL.expr.Pos = yyDollar[2].expr.Pos
		}
	case 50:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:269
		{
			yyVAL.expr = NewStringNode(yyDollar[1].token.Str)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:275
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 52:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:278
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 53:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:281
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:284
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:289
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 56:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:294
		{
			yyVAL.expr = NewCompoundNode("call", yyDollar[1].expr, yyDollar[2].exprlist)
		}
	case 57:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:299
		{
			if yylex.(*Lexer).PNewLine {
				yylex.(*Lexer).TokenError(yyDollar[1].token, "ambiguous syntax (function call x new statement)")
			}
			yyVAL.exprlist = NewCompoundNode()
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:305
		{
			if yylex.(*Lexer).PNewLine {
				yylex.(*Lexer).TokenError(yyDollar[1].token, "ambiguous syntax (function call x new statement)")
			}
			yyVAL.exprlist = yyDollar[2].exprlist
		}
	case 59:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:313
		{
			yyVAL.expr = NewCompoundNode("lambda", NewCompoundNode(), yyDollar[4].stmts)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	case 60:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:317
		{
			yyVAL.expr = NewCompoundNode("lambda", yyDollar[3].namelist, yyDollar[5].stmts)
			yyVAL.expr.Pos = yyDollar[1].token.Pos
		}
	}
	goto yystack /* stack new state and value */
}
