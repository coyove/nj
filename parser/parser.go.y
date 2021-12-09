%{
package parser
%}
%type<expr> prog
%type<expr> stats
%type<expr> prog_stat
%type<expr> stat
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
%type<expr> func_docstring
%type<expr> comma

%union {
    token Token
    expr  Node
}

/* Reserved words */
%token<token> TDo TLocal TElseIf TThen TEnd TBreak TContinue TElse TFor TWhile TFunc TLambda TIf TReturn TReturnVoid TRepeat TUntil TNot TLabel TGoto TIn TNext TLsh TRsh TURsh TDotDotDot TLParen

/* Literals */
%token<token> TOr TAnd TEqeq TNeq TLte TGte TIdent TNumber TString TIDiv
%token<token> '{' '[' '(' '=' '>' '<' '+' '-' '*' '/' '%' '^' '#' '.' '&' '|' '~'

/* Operators */
%right 'T'
%right TElse
%left ASSIGN
%right FUNC
%left TOr
%left TAnd
%left '>' '<' TGte TLte TEqeq TNeq
%left '&' '|' '^'
%left TLsh TRsh TURsh 
%left '+' '-' 
%left '*' '/' '%' TIDiv
%right UNARY /* not # -(unary) */

%% 

prog: 
    {
        $$ = __chain()
        yylex.(*Lexer).Stmts = $$
    } |
    prog prog_stat {
        $$ = $1.append($2)
        yylex.(*Lexer).Stmts = $$
    }

stats: 
    { $$ = __chain() } | stats stat { $$ = $1.append($2) }

prog_stat:
    func_stat      { $$ = $1 } |
    stat           { $$ = $1 }

stat:
    TDo stats TEnd { $$ = __do($2) } |
    jmp_stat       { $$ = $1 } |
    assign_stat    { $$ = $1 } |
    for_stat       { $$ = $1 } |
    if_stat        { $$ = $1 } |
    ';'            { $$ = emptyNode }

assign_stat:
    prefix_expr {
        $$ = $1
    } | 
    TLocal ident_list {
        $$ = __chain()
        for _, v := range $2.Nodes() {
            $$ = $$.append(__set(v, SNil).At($1))
        }
    } |
    TLocal ident_list '=' expr_list {
        if len($4.Nodes()) == 1 && len($2.Nodes()) > 1 {
            tmp := randomVarname()
            $$ = __chain(__set(tmp, $4.Nodes()[0]).At($1))
            for i, ident := range $2.Nodes() {
                $$ = $$.append(__set(ident, __load(tmp, Int(int64(i))).At($1)).At($1))
            }
        } else {
            $$ = __local($2.Nodes(), $4.Nodes(), $1)
        }
    } |
    declarator_list '=' expr_list {
        if len($3.Nodes()) == 1 && len($1.Nodes()) > 1 {
            tmp := randomVarname()
            $$ = __chain(__set(tmp, $3.Nodes()[0]).At($2))
            for i, decl := range $1.Nodes() {
                x := decl.moveLoadStore(__move, __load(tmp, Int(int64(i))).At($2)).At($2)
                $$ = $$.append(x)
            }
        } else {
            $$ = __moveMulti($1.Nodes(), $3.Nodes(), $2)
        }
    }

for_stat:
    TWhile expr TDo stats TEnd {
        $$ = __loop(emptyNode, __if($2, $4, breakNode).At($1)).At($1)
    } |
    TRepeat stats TUntil expr {
        $$ = __loop(emptyNode, $2, __if($4, breakNode, emptyNode).At($1)).At($1)
    } |
    TFor TIdent '=' expr ',' expr TDo stats TEnd {
        forVar, forEnd := Sym($2), randomVarname()
        cont := __inc(forVar, one).At($1)
        $$ = __do(
            __set(forVar, $4).At($1),
            __set(forEnd, $6).At($1),
            __loop(
                cont,
                __if(
                    __less(forVar, forEnd),
                    __chain($8, cont),
                    breakNode,
                ).At($1),
            ).At($1),
        )
    } |
    TFor TIdent '=' expr ',' expr ',' expr TDo stats TEnd {
        forVar, forEnd, forStep := Sym($2), randomVarname(), randomVarname()
        body := __chain($10, __inc(forVar, forStep))
        $$ = __do(
            __set(forVar, $4).At($1),
            __set(forEnd, $6).At($1),
            __set(forStep, $8).At($1),
        )
        if $8.IsNum() { // step is a static number, easy case
            if $8.IsNegativeNumber() {
                $$ = $$.append(__loop(__inc(forVar, forStep), __if(__less(forEnd, forVar), body, breakNode).At($1)).At($1))
            } else {
                $$ = $$.append(__loop(__inc(forVar, forStep), __if(__less(forVar, forEnd), body, breakNode).At($1)).At($1))
            }
        } else { 
            $$ = $$.append(__loop(
                __inc(forVar, forStep),
                __if(
                    __less(zero, forStep).At($1),
                    __if(__lessEq(forEnd, forVar), breakNode, body).At($1), // +step
                    __if(__lessEq(forVar, forEnd), breakNode, body).At($1), // -step
                ).At($1),
            ).At($1))
        }
    } |
    TFor TIdent ',' TIdent TIn expr TDo stats TEnd          { $$ = __forIn($2, $4, $6, $8, $1) } |
    TFor TIdent TIn expr TDo stats TEnd                     { $$ = __forIn($2, $1, $4, $6, $1) }

if_stat:
    TIf expr TThen stats elseif_stat TEnd %prec 'T' { $$ = __if($2, $4, $5).At($1) }

elseif_stat:
    { $$ = Nodes() } | TElse stats { $$ = $2 } | TElseIf expr TThen stats elseif_stat { $$ = __if($2, $4, $5).At($1) }

func_stat:
    TFunc TIdent func_params stats TEnd            { $$ = __func($2, $3, $4) } | 
    TFunc TIdent '.' TIdent func_params stats TEnd { $$ = __store(Sym($2), Str($4.Str), __func(__markupFuncName($2, $4), $5, $6)) }

func_params:
    TLParen ')' func_docstring                       { $$ = Nodes(emptyNode, $3) } | 
    TLParen ident_list ')' func_docstring            { $$ = Nodes($2, $4) } |
    TLParen ident_list TDotDotDot ')' func_docstring { $$ = Nodes(__dotdotdot($2), $5) } |
    '(' ')' func_docstring                           { $$ = Nodes(emptyNode, $3) } | 
    '(' ident_list ')' func_docstring                { $$ = Nodes($2, $4) } |
    '(' ident_list TDotDotDot ')' func_docstring     { $$ = Nodes(__dotdotdot($2), $5) }

func_docstring:
    { $$ = nullStr } | TString { $$ = Str($1.Str) }

jmp_stat:
    TBreak               { $$ = Nodes(SBreak).At($1) } |
    TContinue            { $$ = Nodes(SContinue).At($1) } |
    TGoto TIdent         { $$ = __goto(Sym($2)).At($1) } |
    TLabel TIdent TLabel { $$ = __label(Sym($2)) } |
    TReturnVoid          { $$ = __ret(SNil).At($1) } |
    TReturn expr_list {
        if len($2.Nodes()) == 1 {
            __findTailCall($2.Nodes())
            $$ = __ret($2.Nodes()[0]).At($1) 
        } else {
            $$ = __ret(Nodes(SArray, $2)).At($1) 
        }
    }

declarator:
    TIdent                   { $$ = Sym($1) } |
    prefix_expr '[' expr ']' { $$ = __load($1, $3).At($2) } |
    prefix_expr '.' TIdent   { $$ = __load($1, Str($3.Str)).At($2) } 

expr:
    prefix_expr                       { $$ = $1 } |
    TLambda func_params stats TEnd    { $$ = __lambda(__markupLambdaName($1), $2, $3) } | 
    TNumber                           { $$ = Num($1.Str) } |
    TString                           { $$ = Str($1.Str) } |
    '[' ']'                           { $$ = Nodes(SArray, emptyNode).At($1) } |
    '{' '}'                           { $$ = Nodes(SObject, emptyNode).At($1) } |
    '[' expr_list comma ']'           { $$ = Nodes(SArray, $2).At($1) } |
    '{' expr_assign_list comma'}'     { $$ = Nodes(SObject, $2).At($1) } |
    expr TOr expr                     { $$ = Nodes((SOr), $1,$3).At($2) } |
    expr TAnd expr                    { $$ = Nodes((SAnd), $1,$3).At($2) } |
    expr '>' expr                     { $$ = Nodes((SLess), $3,$1).At($2) } |
    expr '<' expr                     { $$ = Nodes((SLess), $1,$3).At($2) } |
    expr TGte expr                    { $$ = Nodes((SLessEq), $3,$1).At($2) } |
    expr TLte expr                    { $$ = Nodes((SLessEq), $1,$3).At($2) } |
    expr TEqeq expr                   { $$ = Nodes((SEq), $1,$3).At($2) } |
    expr TNeq expr                    { $$ = Nodes((SNeq), $1,$3).At($2) } |
    expr '+' expr                     { $$ = Nodes((SAdd), $1,$3).At($2) } |
    expr '-' expr                     { $$ = Nodes((SSub), $1,$3).At($2) } |
    expr '*' expr                     { $$ = Nodes((SMul), $1,$3).At($2) } |
    expr '/' expr                     { $$ = Nodes((SDiv), $1,$3).At($2) } |
    expr TIDiv expr                   { $$ = Nodes((SIDiv), $1,$3).At($2) } |
    expr '%' expr                     { $$ = Nodes((SMod), $1,$3).At($2) } |
    expr '&' expr                     { $$ = Nodes((SBitAnd), $1,$3).At($2) } |
    expr '|' expr                     { $$ = Nodes((SBitOr), $1,$3).At($2) } |
    expr '^' expr                     { $$ = Nodes((SBitXor), $1,$3).At($2) } |
    expr TLsh expr                    { $$ = Nodes((SBitLsh), $1,$3).At($2) } |
    expr TRsh expr                    { $$ = Nodes((SBitRsh), $1,$3).At($2) } |
    expr TURsh expr                   { $$ = Nodes((SBitURsh), $1,$3).At($2) } |
    '~' expr %prec UNARY              { $$ = Nodes((SBitNot), $2).At($1) } |
    '#' expr %prec UNARY              { $$ = Nodes((SLen), $2).At($1) } |
    TNot expr %prec UNARY             { $$ = Nodes((SNot), $2).At($1) } |
    '-' expr %prec UNARY              { $$ = Nodes((SSub), zero, $2).At($1) } |
    '+' expr %prec UNARY              { $$ = Nodes((SAdd), zero, $2).At($1) }

prefix_expr:
    '(' expr ')'                                       { $$ = $2 } |
    declarator                                         { $$ = $1 } |
    prefix_expr TString                                { $$ = __call($1, Nodes(Str($2.Str))).At($2) } |
    prefix_expr TLParen ')'                            { $$ = __call($1, emptyNode).At($2) } |
    prefix_expr TLParen expr_list comma ')'            { $$ = __call($1, $3).At($2) } |
    prefix_expr TLParen expr_assign_list comma ')'     { $$ = __call($1, Nodes(Nodes(SObject, $3).At($2))).At($2) } |
    prefix_expr TLParen expr_list TDotDotDot comma ')' { $$ = __call($1, __dotdotdot($3)).At($2) } |
    prefix_expr TLParen expr_list ',' expr_assign_list comma ')' { $$ = __call($1, $3.append(Nodes(SObject, $5).At($2))).At($2) }

declarator_list:
    declarator { $$ = Nodes($1) } | declarator_list ',' declarator { $$ = $1.append($3) }

ident_list:
    TIdent { $$ = Nodes(Sym($1)) } | ident_list ',' TIdent { $$ = $1.append(Sym($3)) }

expr_list:
    expr { $$ = Nodes($1) } | expr_list ',' expr { $$ = $1.append($3) }

expr_assign_list:
    TIdent '=' expr                            { $$ = Nodes(Str($1.Str), $3) } |
    '{' expr '}' '=' expr                      { $$ = Nodes($2, $5) } |
    expr_assign_list ',' TIdent '=' expr       { $$ = $1.append(Str($3.Str)).append($5) } |
    expr_assign_list ',' '{' expr '}' '=' expr { $$ = $1.append($4).append($7) }

comma: 
    { $$ = emptyNode } | ',' { $$ = emptyNode }

%%
