%{
package parser

import "github.com/coyove/nj/typ"

func ss(yylex yyLexer) *Lexer { return yylex.(*Lexer) }
%}
%type<expr> prog
%type<expr> stats
%type<expr> declarator
%type<expr> declarator_list
%type<expr> ident_list
%type<expr> expr
%type<expr> expr_list
%type<expr> expr_assign_list
%type<expr> prefix_expr
%type<expr> assign_stat
%type<expr> for_stat
%type<expr> if_stat
%type<expr> elseif_stat
%type<expr> jmp_stat
%type<expr> func_stat
%type<expr> func_params
%type<token> comma

%union {
    token  Token
    expr   Node
}

/* Reserved words */
%token<token> TDo TLocal TElseIf TThen TEnd TBreak TContinue TElse TFor TWhile TFunc TIf TReturn TReturnVoid TRepeat TUntil TNot TLabel TGoto TIn TLsh TRsh TURsh TDotDotDot TLParen TLBracket TIs

/* Literals */
%token<token> TOr TAnd TEqeq TNeq TLte TGte TIdent TNumber TString TIDiv TInv
%token<token> TAddEq TSubEq TMulEq TDivEq TIDivEq TModEq TBitAndEq TBitOrEq TBitXorEq TBitLshEq TBitRshEq TBitURshEq
%token<token> '{' '[' '(' '=' '>' '<' '+' '-' '*' '/' '%' '^' '#' '.' '&' '|' '~' ':' ')' ','

/* Operators */
%right 'T'
%right TElse
%left ASSIGN
%right FUNC
%left TOr
%left TAnd
%left '>' '<' TGte TLte TEqeq TNeq
%left '+' '-' '|' '^'
%left '*' '/' '%' TIDiv TLsh TRsh TURsh '&'
%left TIs
%right UNARY /* not # -(unary) */

%% 

prog: stats { ss(yylex).Stmts = $1 }

stats: 
                                    { $$ = &Prog{} } |
    stats func_stat                 { $$ = $1.(*Prog).Append($2) } |
    stats TDo stats TEnd            { $3.(*Prog).DoBlock = true; $$ = $1.(*Prog).Append($3) } |
    stats jmp_stat                  { $$ = $1.(*Prog).Append($2) } |
    stats assign_stat               { $$ = $1.(*Prog).Append($2) } |
    stats for_stat                  { $$ = $1.(*Prog).Append($2) } |
    stats if_stat                   { $$ = $1.(*Prog).Append($2) } |
    stats ';'                       { $$ = $1 }

assign_stat:
    expr                            { $$ = $1 } | 
    TLocal ident_list               { $$ = ss(yylex).pDeclareAssign([]Node($2.(IdentList)), nil, false, $1) } |
    TLocal ident_list '=' expr_list { $$ = ss(yylex).pDeclareAssign([]Node($2.(IdentList)), $4.(ExprList), false, $1) } |
    declarator_list '=' expr_list   { $$ = ss(yylex).pDeclareAssign([]Node($1.(DeclList)), $3.(ExprList), true, $2) } |
    declarator TAddEq expr          { $$ = assignLoadStore($1, ss(yylex).pBinary(typ.OpAdd, $1, $3, $2), $2) } |
    declarator TSubEq expr          { $$ = assignLoadStore($1, ss(yylex).pBinary(typ.OpSub, $1, $3, $2), $2) } |
    declarator TMulEq expr          { $$ = assignLoadStore($1, ss(yylex).pBinary(typ.OpMul, $1, $3, $2), $2) } |
    declarator TDivEq expr          { $$ = assignLoadStore($1, ss(yylex).pBinary(typ.OpDiv, $1, $3, $2), $2) } |
    declarator TIDivEq expr         { $$ = assignLoadStore($1, ss(yylex).pBinary(typ.OpIDiv, $1, $3, $2), $2) } |
    declarator TModEq expr          { $$ = assignLoadStore($1, ss(yylex).pBinary(typ.OpMod, $1, $3, $2), $2) } |
    declarator TBitAndEq expr       { $$ = assignLoadStore($1, ss(yylex).pBitwise("and", $1, $3, $2), $2) } |
    declarator TBitOrEq expr        { $$ = assignLoadStore($1, ss(yylex).pBitwise("or", $1, $3, $2), $2) } |
    declarator TBitXorEq expr       { $$ = assignLoadStore($1, ss(yylex).pBitwise("xor", $1, $3, $2), $2) } |
    declarator TBitLshEq expr       { $$ = assignLoadStore($1, ss(yylex).pBitwise("lsh", $1, $3, $2), $2) } |
    declarator TBitRshEq expr       { $$ = assignLoadStore($1, ss(yylex).pBitwise("rsh", $1, $3, $2), $2) } |
    declarator TBitURshEq expr      { $$ = assignLoadStore($1, ss(yylex).pBitwise("ursh", $1, $3, $2), $2) }

for_stat:
    TWhile expr TDo stats TEnd                            { $$ = ss(yylex).pLoop(&If{$2, $4, emptyBreak}) } |
    TRepeat stats TUntil expr                             { $$ = ss(yylex).pLoop($2, &If{$4, emptyBreak, emptyProg}) } |
    TFor TIdent '=' expr ',' expr TDo stats TEnd          { $$ = ss(yylex).pForRange($2, $4, $6, one, $8, $1) } |
    TFor TIdent '=' expr ',' expr ',' expr TDo stats TEnd { $$ = ss(yylex).pForRange($2, $4, $6, $8, $10, $1) } |
    TFor TIdent ',' TIdent TIn expr TDo stats TEnd        { $$ = ss(yylex).pForIn($2, $4, $6, $8, $1) } |
    TFor TIdent TIn expr TDo stats TEnd                   { $$ = ss(yylex).pForIn($2, $1, $4, $6, $1) }

if_stat:
    TIf expr TThen stats elseif_stat TEnd %prec 'T' { $$ = &If{$2, $4, $5} }

elseif_stat:
                                                    { $$ = nil } |
    TElse stats                                     { $$ = $2 } |
    TElseIf expr TThen stats elseif_stat            { $$ = &If{$2, $4, $5} }

func_stat:
    TFunc TIdent func_params stats TEnd {
        $$ = ss(yylex).pFunc(false, $2, $3, $4, $1)
    } | 
    TFunc TIdent '.' TIdent func_params stats TEnd {
        m := ss(yylex).pFunc(true, __markupFuncName($2, $4), $5, $6, $1)
        $$ = &Tenary{typ.OpStore, Sym($2), ss(yylex).Str($4.Str), m, $1.Line()}
    }

func_params:
    TLParen ')'                                   { $$ = (IdentList)(nil) } | 
    TLParen ident_list ')'                        { $$ = $2 } |
    TLParen ident_list TDotDotDot ')'             { $$ = IdentVarargList{$2.(IdentList)} } |
    TLParen TDotDotDot ident_list ')'             { $$ = IdentVarargExpandList{nil, $3.(IdentList)} } |
    TLParen ident_list TDotDotDot ident_list ')'  { $$ = IdentVarargExpandList{$2.(IdentList), $4.(IdentList)} } |
    '(' ')'                                       { $$ = (IdentList)(nil) } | 
    '(' ident_list ')'                            { $$ = $2 } |
    '(' ident_list TDotDotDot ')'                 { $$ = IdentVarargList{$2.(IdentList)} } |
    '(' TDotDotDot ident_list ')'                 { $$ = IdentVarargExpandList{nil, $3.(IdentList)} } |
    '(' ident_list TDotDotDot ident_list ')'      { $$ = IdentVarargExpandList{$2.(IdentList), $4.(IdentList)} }

jmp_stat:
    TBreak               { $$ = &BreakContinue{true, $1.Line()} } |
    TContinue            { $$ = &BreakContinue{false, $1.Line()} } |
    TGoto TIdent         { $$ = &GotoLabel{$2.Str, true, $1.Line()} } |
    TLabel TIdent TLabel { $$ = &GotoLabel{$2.Str, false, $1.Line()} } |
    TReturnVoid          { $$ = &Unary{typ.OpRet, SNil, $1.Line()} } |
    TReturn expr_list {
        if el := $2.(ExprList); len(el) == 1 {
            ss(yylex).pFindTailCall(el[0])
            $$ = &Unary{typ.OpRet, el[0], $1.Line()}
        } else {
            $$ = &Unary{typ.OpRet, $2, $1.Line()}
        }
    }

declarator:
    TIdent {
        $$ = Sym($1)
    } |
    prefix_expr TLBracket expr ']' {
        $$ = &Tenary{typ.OpLoad, $1, $3, Address(typ.RegA), $2.Line()}
    } |
    prefix_expr '.' TIdent {
        $$ = &Tenary{typ.OpLoad, $1, ss(yylex).Str($3.Str), Address(typ.RegA), $2.Line()}
    } 

expr:
    prefix_expr                       { $$ = $1 } |
    TNumber                           { $$ = ss(yylex).Num($1.Str) } |
    expr TOr expr                     { $$ = &Or{$1, $3} } |
    expr TAnd expr                    { $$ = &And{$1, $3} } |
    expr '>' expr                     { $$ = ss(yylex).pBinary(typ.OpLess, $3, $1, $2) } |
    expr '<' expr                     { $$ = ss(yylex).pBinary(typ.OpLess, $1, $3, $2) } |
    expr TGte expr                    { $$ = ss(yylex).pBinary(typ.OpLessEq, $3, $1, $2) } |
    expr TLte expr                    { $$ = ss(yylex).pBinary(typ.OpLessEq, $1, $3, $2) } |
    expr TEqeq expr                   { $$ = ss(yylex).pBinary(typ.OpEq, $1, $3, $2) } |
    expr TNeq expr                    { $$ = ss(yylex).pBinary(typ.OpNeq, $1, $3, $2) } |
    expr '+' expr                     { $$ = ss(yylex).pBinary(typ.OpAdd, $1, $3, $2) } |
    expr '-' expr                     { $$ = ss(yylex).pBinary(typ.OpSub, $1, $3, $2) } |
    expr '*' expr                     { $$ = ss(yylex).pBinary(typ.OpMul, $1, $3, $2) } |
    expr '/' expr                     { $$ = ss(yylex).pBinary(typ.OpDiv, $1, $3, $2) } |
    expr TIDiv expr                   { $$ = ss(yylex).pBinary(typ.OpIDiv, $1, $3, $2) } |
    expr '%' expr                     { $$ = ss(yylex).pBinary(typ.OpMod, $1, $3, $2) } |
    expr '&' expr                     { $$ = ss(yylex).pBitwise("and", $1, $3 ,$2) } |
    expr '|' expr                     { $$ = ss(yylex).pBitwise("or", $1, $3, $2) } |
    expr '^' expr                     { $$ = ss(yylex).pBitwise("xor", $1, $3, $2) } |
    expr TLsh expr                    { $$ = ss(yylex).pBitwise("lsh", $1, $3, $2) } |
    expr TRsh expr                    { $$ = ss(yylex).pBitwise("rsh", $1, $3, $2) } |
    expr TURsh expr                   { $$ = ss(yylex).pBitwise("ursh", $1, $3, $2) } |
    expr TIs prefix_expr              { $$ = ss(yylex).pBinary(typ.OpIsProto, $1, $3, $2) } |
    expr TIs TNot prefix_expr         { $$ = pUnary(typ.OpNot, ss(yylex).pBinary(typ.OpIsProto, $1, $4, $2), $2) } |
    '~' expr %prec UNARY              { $$ = ss(yylex).pBitwise("xor", ss(yylex).Int(-1), $2, $1) } |
    '#' expr %prec UNARY              { $$ = pUnary(typ.OpLen, $2, $1) } |
    TInv expr %prec UNARY             { $$ = ss(yylex).pBinary(typ.OpSub, zero, $2, $1) } |
    TNot expr %prec UNARY             { $$ = pUnary(typ.OpNot, $2, $1) }

prefix_expr:
    declarator                                         { $$ = $1 } |
    TIf TLParen expr ',' expr ',' expr ')'             { $$ = &If{$3, &Assign{Sa, $5, $1.Line()}, &Assign{Sa, $7, $1.Line()}} } |
    TFunc func_params stats TEnd                       { $$ = ss(yylex).pFunc(false, __markupLambdaName($1), $2, $3, $1) } | 
    TString                                            { $$ = ss(yylex).Str($1.Str) } |
    '(' expr ')'                                       { $$ = $2 } |
    '[' ']'                                            { $$ = ss(yylex).pEmptyArray() } |
    '{' '}'                                            { $$ = ss(yylex).pEmptyObject() } |
    '[' expr_list comma ']'                            { $$ = $2 } |
    '{' expr_assign_list comma'}'                      { $$ = $2 } |
    prefix_expr TLBracket expr ':' expr ']'            { $$ = &Tenary{typ.OpSlice, $1, $3, $5, $2.Line()} } |
    prefix_expr TLBracket ':' expr ']'                 { $$ = &Tenary{typ.OpSlice, $1, zero, $4, $2.Line()} } |
    prefix_expr TLBracket expr ':' ']'                 { $$ = &Tenary{typ.OpSlice, $1, $3, ss(yylex).Int(-1), $2.Line()} } |
    prefix_expr TLParen ')'                            { $$ = &Call{typ.OpCall, $1, ExprList(nil), false, $2.Line() } } |
    prefix_expr TLParen expr_list comma ')'            { $$ = &Call{typ.OpCall, $1, $3.(ExprList), false, $2.Line() } } |
    prefix_expr TLParen expr_list TDotDotDot comma ')' { $$ = &Call{typ.OpCall, $1, $3.(ExprList), true, $2.Line() } }

declarator_list:
    declarator                     { $$ = DeclList{$1} } |
    declarator_list ',' declarator { $$ = append($1.(DeclList), $3) }

ident_list:
    TIdent                         { $$ = IdentList{Sym($1)} } |
    ident_list ',' TIdent          { $$ = append($1.(IdentList), Sym($3)) }

expr_list:
    expr                           { $$ = ss(yylex).pArray(nil, $1) } |
    expr_list ',' expr             { $$ = ss(yylex).pArray($1, $3) }

expr_assign_list:
    TIdent '=' expr                       { $$ = ss(yylex).pObject(nil, ss(yylex).Str($1.Str), $3) } |
    expr ':' expr                         { $$ = ss(yylex).pObject(nil, $1, $3) } |
    expr_assign_list ',' TIdent '=' expr  { $$ = ss(yylex).pObject($1, ss(yylex).Str($3.Str), $5) } |
    expr_assign_list ',' expr ':' expr    { $$ = ss(yylex).pObject($1, $3, $5) }

comma: {} | ',' {}

%%
