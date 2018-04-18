//line parser.go.y:2
package parser

import __yyfmt__ "fmt"

//line parser.go.y:2
//line parser.go.y:22
type yySymType struct {
	yys   int
	token Token

	stmts []interface{}
	stmt  interface{}

	funcname interface{}
	funcexpr interface{}

	exprlist []interface{}
	expr     interface{}

	namelist []interface{}
	parlist  interface{}
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
const TEqeq = 57366
const TNeq = 57367
const TLte = 57368
const TGte = 57369
const TIdent = 57370
const TNumber = 57371
const TString = 57372
const UNARY = 57373

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

//line parser.go.y:304

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

const yyLast = 392

var yyAct = [...]int{

	37, 27, 32, 5, 4, 47, 94, 21, 34, 89,
	66, 70, 68, 96, 66, 1, 38, 40, 41, 46,
	56, 57, 58, 65, 59, 52, 53, 51, 50, 87,
	60, 61, 16, 59, 48, 49, 54, 55, 56, 57,
	58, 63, 59, 100, 42, 35, 69, 72, 73, 74,
	75, 76, 77, 78, 79, 80, 81, 82, 83, 84,
	85, 71, 86, 36, 20, 14, 91, 92, 10, 11,
	102, 103, 101, 5, 4, 62, 17, 7, 18, 26,
	90, 39, 9, 8, 44, 19, 6, 13, 28, 88,
	97, 12, 5, 4, 2, 15, 0, 5, 4, 47,
	5, 4, 0, 95, 106, 3, 0, 5, 4, 0,
	98, 5, 4, 46, 0, 0, 0, 0, 105, 52,
	53, 51, 50, 0, 109, 0, 0, 0, 48, 49,
	54, 55, 56, 57, 58, 23, 59, 31, 0, 0,
	22, 30, 0, 67, 47, 0, 24, 0, 0, 0,
	0, 0, 12, 25, 33, 0, 15, 0, 46, 0,
	29, 108, 0, 0, 52, 53, 51, 50, 47, 0,
	0, 0, 43, 48, 49, 54, 55, 56, 57, 58,
	0, 59, 46, 0, 0, 64, 0, 0, 52, 53,
	51, 50, 47, 0, 0, 45, 0, 48, 49, 54,
	55, 56, 57, 58, 0, 59, 46, 0, 0, 0,
	0, 0, 52, 53, 51, 50, 47, 0, 0, 0,
	0, 48, 49, 54, 55, 56, 57, 58, 0, 59,
	46, 0, 0, 0, 47, 0, 52, 53, 51, 50,
	0, 0, 0, 0, 0, 48, 49, 54, 55, 56,
	57, 58, 0, 59, 52, 53, 51, 50, 0, 0,
	0, 0, 0, 48, 49, 54, 55, 56, 57, 58,
	0, 59, 52, 53, 51, 50, 0, 0, 0, 0,
	0, 48, 49, 54, 55, 56, 57, 58, 0, 59,
	10, 11, 10, 11, 0, 107, 0, 104, 0, 7,
	0, 7, 0, 0, 9, 8, 9, 8, 6, 0,
	6, 0, 0, 12, 0, 12, 0, 15, 0, 15,
	0, 10, 11, 10, 11, 0, 99, 3, 93, 3,
	7, 0, 7, 0, 0, 9, 8, 9, 8, 6,
	0, 6, 0, 0, 12, 0, 12, 0, 15, 0,
	15, 0, 0, 0, 0, 0, 0, 0, 3, 0,
	3, 54, 55, 56, 57, 58, 23, 59, 31, 0,
	0, 22, 30, 0, 0, 0, 0, 24, 0, 0,
	0, 0, 0, 12, 25, 33, 0, 15, 0, 0,
	0, 29,
}
var yyPact = [...]int{

	-1000, 63, -1000, -1000, -11, 32, 355, 355, 17, 355,
	-1000, -1000, -1000, -1000, -1000, 355, 355, 355, 16, -1000,
	124, 188, -1000, -1000, -1000, -1000, -1000, 32, -1000, 355,
	355, 9, -1000, -1000, 164, -20, -33, 212, 95, -36,
	212, 1, -1000, -1000, -37, -1000, 355, 355, 355, 355,
	355, 355, 355, 355, 355, 355, 355, 355, 355, 355,
	-8, -8, -1000, -19, -1000, 355, 355, -1000, -1000, -1000,
	-1000, 318, 230, 248, 326, 326, 326, 326, 326, 326,
	-17, -17, -8, -8, -8, -8, -42, -1000, -34, -1000,
	63, 212, 212, -1000, -1000, 316, 15, 62, 287, -1000,
	-1000, -1000, -1000, 355, -1000, 285, 140, -1000, -1000, 63,
}
var yyPgo = [...]int{

	0, 15, 94, 90, 90, 90, 2, 89, 63, 0,
	88, 1, 65, 87, 85, 79, 75, 62,
}
var yyR1 = [...]int{

	0, 1, 1, 1, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 3, 3, 4, 5, 6, 6,
	6, 7, 7, 8, 8, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 10, 11,
	11, 11, 11, 13, 12, 14, 14, 15, 16, 16,
	17,
}
var yyR2 = [...]int{

	0, 0, 2, 2, 3, 1, 5, 6, 8, 4,
	1, 2, 1, 1, 0, 5, 1, 1, 1, 4,
	3, 1, 3, 1, 3, 1, 1, 1, 1, 1,
	1, 1, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 2, 2, 1, 1,
	1, 1, 3, 3, 2, 2, 3, 2, 5, 4,
	1,
}
var yyChk = [...]int{

	-1000, -1, -2, 42, -6, -11, 23, 14, 20, 19,
	5, 6, 28, -13, -12, 32, 43, 44, 46, -14,
	32, -9, 16, 11, 22, 29, -15, -11, -10, 36,
	17, 13, -6, 30, -9, 28, -8, -9, -9, -12,
	-9, -9, 28, 48, -8, 7, 18, 4, 33, 34,
	27, 26, 24, 25, 35, 36, 37, 38, 39, 41,
	-9, -9, -16, 32, 21, 43, 47, 48, 48, 45,
	48, -1, -9, -9, -9, -9, -9, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, -17, 48, -7, 28,
	-1, -9, -9, 10, 48, -1, 47, -3, -1, 10,
	28, 10, 8, 9, 10, -1, -9, 10, 21, -1,
}
var yyDef = [...]int{

	1, -2, 2, 3, 49, 5, 0, 0, 0, 10,
	12, 13, 18, 50, 51, 0, 0, 0, 0, 54,
	0, 0, 25, 26, 27, 28, 29, 30, 31, 0,
	0, 0, 49, 48, 0, 0, 11, 23, 0, 51,
	4, 0, 20, 55, 0, 1, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	46, 47, 57, 0, 1, 0, 0, 52, 53, 19,
	56, 0, 32, 33, 34, 35, 36, 37, 38, 39,
	40, 41, 42, 43, 44, 45, 0, 1, 60, 21,
	14, 9, 24, 6, 1, 0, 0, 0, 0, 59,
	22, 7, 1, 0, 58, 0, 0, 8, 1, 15,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 39, 3, 3,
	32, 48, 37, 35, 47, 36, 46, 38, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 42,
	34, 43, 33, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 44, 3, 45, 41, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 31,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 40,
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
		//line parser.go.y:56
		{
			yyVAL.stmts = []interface{}{"chain"}
			if l, ok := yylex.(*Lexer); ok {
				l.Stmts = yyVAL.stmts
			}
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:62
		{
			yyVAL.stmts = append(yyDollar[1].stmts, yyDollar[2].stmt)
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
			yyVAL.stmt = []interface{}{"set", yyDollar[1].expr, yyDollar[3].expr}
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
			yyVAL.stmt = []interface{}{"while", yyDollar[2].expr, yyDollar[4].stmts}
		}
	case 7:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:90
		{
			yyVAL.stmt = []interface{}{"if", yyDollar[2].expr, yyDollar[4].stmts, nil}
			if len(yyDollar[5].stmts) > 0 {
				cur := yyVAL.stmt
				for _, e := range yyDollar[5].stmts {
					cur.([]interface{})[3] = e
					cur = e
				}
			}
		}
	case 8:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.go.y:100
		{
			yyVAL.stmt = []interface{}{"if", yyDollar[2].expr, yyDollar[4].stmts, nil}
			cur := yyVAL.stmt
			if len(yyDollar[5].stmts) > 0 {
				for _, e := range yyDollar[5].stmts {
					cur.([]interface{})[3] = e
					cur = e
				}
			}
			cur.([]interface{})[3] = yyDollar[7].stmts
		}
	case 9:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:111
		{
			yyVAL.stmt = []interface{}{"var", yyDollar[2].token.Str, yyDollar[4].expr}
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:114
		{
			yyVAL.stmt = []interface{}{"ret"}
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:117
		{
			yyVAL.stmt = []interface{}{"ret", yyDollar[2].exprlist}
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:120
		{
			yyVAL.stmt = "break"
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:123
		{
			yyVAL.stmt = "continue"
		}
	case 14:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.go.y:128
		{
			yyVAL.stmts = []interface{}{}
		}
	case 15:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:131
		{
			yyVAL.stmts = append(yyVAL.stmts, []interface{}{"if", yyDollar[3].expr, yyDollar[5].stmts, nil})
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:136
		{
			yyVAL.funcname = yyDollar[1].funcname
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:141
		{
			yyVAL.funcname = yyDollar[1].token.Str
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:146
		{
			yyVAL.expr = yyDollar[1].token.Str
		}
	case 19:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:149
		{
			yyVAL.expr = []interface{}{yyDollar[1].expr, ":", yyDollar[3].expr}
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:152
		{
			yyVAL.expr = []interface{}{yyDollar[1].expr, ":", yyDollar[3].token.Str}
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:157
		{
			yyVAL.namelist = []interface{}{yyDollar[1].token.Str}
		}
	case 22:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:160
		{
			yyVAL.namelist = append(yyDollar[1].namelist, yyDollar[3].token.Str)
		}
	case 23:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:165
		{
			yyVAL.exprlist = []interface{}{yyDollar[1].expr}
		}
	case 24:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:168
		{
			yyVAL.exprlist = append(yyDollar[1].exprlist, yyDollar[3].expr)
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:173
		{
			yyVAL.expr = "nil"
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:176
		{
			yyVAL.expr = "false"
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:179
		{
			yyVAL.expr = "true"
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:182
		{
			yyVAL.expr = yyDollar[1].token.Str
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
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:191
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 32:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:194
		{
			yyVAL.expr = []interface{}{"or", yyDollar[1].expr, yyDollar[3].expr}
		}
	case 33:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:197
		{
			yyVAL.expr = []interface{}{"and", yyDollar[1].expr, yyDollar[3].expr}
		}
	case 34:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:200
		{
			yyVAL.expr = []interface{}{">", yyDollar[1].expr, yyDollar[3].expr}
		}
	case 35:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:203
		{
			yyVAL.expr = []interface{}{"<", yyDollar[1].expr, yyDollar[3].expr}
		}
	case 36:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:206
		{
			yyVAL.expr = []interface{}{">=", yyDollar[1].expr, yyDollar[3].expr}
		}
	case 37:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:209
		{
			yyVAL.expr = []interface{}{"<=", yyDollar[1].expr, yyDollar[3].expr}
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:212
		{
			yyVAL.expr = []interface{}{"==", yyDollar[1].expr, yyDollar[3].expr}
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:215
		{
			yyVAL.expr = []interface{}{"!=", yyDollar[1].expr, yyDollar[3].expr}
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:218
		{
			yyVAL.expr = []interface{}{"+", yyDollar[1].expr, yyDollar[3].expr}
		}
	case 41:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:221
		{
			yyVAL.expr = []interface{}{"-", yyDollar[1].expr, yyDollar[3].expr}
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:224
		{
			yyVAL.expr = []interface{}{"*", yyDollar[1].expr, yyDollar[3].expr}
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:227
		{
			yyVAL.expr = []interface{}{"/", yyDollar[1].expr, yyDollar[3].expr}
		}
	case 44:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:230
		{
			yyVAL.expr = []interface{}{"%", yyDollar[1].expr, yyDollar[3].expr}
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:233
		{
			yyVAL.expr = []interface{}{"^", yyDollar[1].expr, yyDollar[3].expr}
		}
	case 46:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:236
		{
			yyVAL.expr = []interface{}{"-", yyDollar[2].expr}
		}
	case 47:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:239
		{
			yyVAL.expr = []interface{}{"not", yyDollar[2].expr}
		}
	case 48:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:244
		{
			yyVAL.expr = yyDollar[1].token.Str
		}
	case 49:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:249
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 50:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:252
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:255
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:258
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:263
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 54:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:268
		{
			yyVAL.expr = []interface{}{"call", yyDollar[1].expr, yyDollar[2].exprlist}
		}
	case 55:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:273
		{
			if yylex.(*Lexer).PNewLine {
				yylex.(*Lexer).TokenError(yyDollar[1].token, "ambiguous syntax (function call x new statement)")
			}
			yyVAL.exprlist = []interface{}{}
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:279
		{
			if yylex.(*Lexer).PNewLine {
				yylex.(*Lexer).TokenError(yyDollar[1].token, "ambiguous syntax (function call x new statement)")
			}
			yyVAL.exprlist = yyDollar[2].exprlist
		}
	case 57:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:287
		{
			yyVAL.expr = []interface{}{"lambda", yyDollar[2].funcexpr.([]interface{})[0], yyDollar[2].funcexpr.([]interface{})[1]}
		}
	case 58:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:292
		{
			yyVAL.funcexpr = []interface{}{yyDollar[2].parlist, yyDollar[4].stmts}
		}
	case 59:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:295
		{
			yyVAL.funcexpr = []interface{}{nil, yyDollar[3].stmts}
		}
	case 60:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:300
		{
			yyVAL.parlist = yyDollar[1].namelist
		}
	}
	goto yystack /* stack new state and value */
}
