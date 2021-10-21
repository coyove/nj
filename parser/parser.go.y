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
%type<expr> call_expr
%type<expr> assign_stat
%type<expr> for_stat
%type<expr> if_stat
%type<expr> elseif_stat
%type<expr> jmp_stat
%type<expr> flow_stat
%type<expr> func_stat
%type<expr> func_expr
%type<expr> comma
%type<token> lparen

%union {
    token Token
    expr  Node
}

/* Reserved words */
%token<token> TDo TLocal TElseIf TThen TEnd TBreak TElse TFor TWhile TFunc TLambda TIf TReturn TReturnVoid TRepeat TUntil TNot TLabel TGoto TIn TNext TLsh TRsh TURsh TDotDotDot TLParen

/* Literals */
%token<token> TOr TAnd TEqeq TNeq TLte TGte TIdent TNumber TString 
%token<token> '{' '[' '(' '=' '>' '<' '+' '-' '*' '/' '%' '^' '#' '.' '&' '|' '~' '@' TIDiv

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
        {
            $$ = __chain()
        } |
        stats stat {
            $$ = $1.append($2)
        }

prog_stat:
        jmp_stat       { $$ = $1 } |
        flow_stat      { $$ = $1 } |
        assign_stat    { $$ = $1 } |
        func_stat      { $$ = $1 } |
        TDo stats TEnd { $$ = __do($2) } |
        ';'            { $$ = emptyNode }

stat:
        jmp_stat       { $$ = $1 } |
        flow_stat      { $$ = $1 } |
        assign_stat    { $$ = $1 } |
        TDo stats TEnd { $$ = __do($2) } |
        ';'            { $$ = emptyNode }

flow_stat:
        for_stat       { $$ = $1 } |
        if_stat        { $$ = $1 }

assign_stat:
        prefix_expr {
            $$ = $1
        } | 
        TLocal ident_list {
            $$ = __chain()
            for _, v := range $2.Nodes {
                $$ = $$.append(__set(v, NewSymbol(ANil)).SetPos($1.Pos))
            }
        } |
        TLocal ident_list '=' expr_list {
            if len($4.Nodes) == 1 && len($2.Nodes) > 1 {
                tmp := randomVarname()
                $$ = __chain(__local([]Node{tmp}, $4.Nodes, $1.Pos))
                for i, ident := range $2.Nodes {
                    $$ = $$.append(__local([]Node{ident}, []Node{__load(tmp, NewNumberFromInt(int64(i))).SetPos($1.Pos)}, $1.Pos))
                }
            } else {
                $$ = __local($2.Nodes, $4.Nodes, $1.Pos)
            }
        } |
        declarator_list '=' expr_list {
            if len($3.Nodes) == 1 && len($1.Nodes) > 1 {
                tmp := randomVarname()
                $$ = __chain(__local([]Node{tmp}, $3.Nodes, $2.Pos))
                for i, decl := range $1.Nodes {
                    x := decl.moveLoadStore(__move, __load(tmp, NewNumberFromInt(int64(i ))).SetPos($2.Pos)).SetPos($2.Pos)
                    $$ = $$.append(x)
                }
            } else {
                $$ = __moveMulti($1.Nodes, $3.Nodes, $2.Pos)
            }
        }

for_stat:
        TWhile expr TDo stats TEnd {
            $$ = __loop(__if($2, $4, breakNode).SetPos($1.Pos)).SetPos($1.Pos)
        } |
        TRepeat stats TUntil expr {
            $$ = __loop(
                __chain(
                    $2,
                    __if($4, breakNode, emptyNode).SetPos($1.Pos),
                ).SetPos($1.Pos),
            ).SetPos($1.Pos)
        } |
        TFor TIdent '=' expr ',' expr TDo stats TEnd {
            forVar, forEnd := NewSymbolFromToken($2), randomVarname()
            $$ = __do(
                    __set(forVar, $4).SetPos($1.Pos),
                    __set(forEnd, $6).SetPos($1.Pos),
                    __loop(
                        __if(
                            __less(forVar, forEnd),
                            __chain($8, __inc(forVar, oneNode).SetPos($1.Pos)),
                            breakNode,
                        ).SetPos($1.Pos),
                    ).SetPos($1.Pos),
                )
        } |
        TFor TIdent '=' expr ',' expr ',' expr TDo stats TEnd {
            forVar, forEnd, forStep := NewSymbolFromToken($2), randomVarname(), randomVarname()
            body := __chain($10, __inc(forVar, forStep))
            $$ = __do(
                __set(forVar, $4).SetPos($1.Pos),
                __set(forEnd, $6).SetPos($1.Pos),
                __set(forStep, $8).SetPos($1.Pos))

            if $8.IsNumber() { // step is a static number, easy case
                if $8.IsNegativeNumber() {
                    $$ = $$.append(__loop(__if(__less(forEnd, forVar), body, breakNode).SetPos($1.Pos)).SetPos($1.Pos))
                } else {
                    $$ = $$.append(__loop(__if(__less(forVar, forEnd), body, breakNode).SetPos($1.Pos)).SetPos($1.Pos))
                }
            } else { 
                $$ = $$.append(__loop(
                    __if(
                        __less(zeroNode, forStep).SetPos($1.Pos),
                        // +step
                        __if(__lessEq(forEnd, forVar), breakNode, body).SetPos($1.Pos),
                        // -step
                        __if(__lessEq(forVar, forEnd), breakNode, body).SetPos($1.Pos),
                    ).SetPos($1.Pos),
                ).SetPos($1.Pos))
            }
        } |
        TFor TIdent ',' TIdent TIn expr TDo stats TEnd {
            $$ = __forIn($2, $4, $6, $8, $1.Pos)
        } |
        TFor TIdent TIn expr TDo stats TEnd {
            $$ = __forIn($2, $1, $4, $6, $1.Pos)
        }


if_stat:
        TIf expr TThen stats elseif_stat TEnd %prec 'T' {
            $$ = __if($2, $4, $5).SetPos($1.Pos)
        }

elseif_stat:
        {
            $$ = NewComplex()
        } |
        TElse stats {
            $$ = $2
        } |
        TElseIf expr TThen stats elseif_stat {
            $$ = __if($2, $4, $5).SetPos($1.Pos)
        }

func_stat:
        TFunc TIdent lparen ')' stats TEnd                               { $$ = __func($2, emptyNode, "", $5) } | 
        TFunc TIdent lparen ident_list ')' stats TEnd                    { $$ = __func($2, $4, "", $6) } | 
        TFunc TIdent lparen ident_list TDotDotDot ')' stats TEnd         { $$ = __func($2, __dotdotdot($4), "", $7) } | 
        TFunc TIdent '.' TIdent lparen ')' stats TEnd                    {
            $$ = __store(NewSymbolFromToken($2), NewString($4.Str), __func(__markupFuncName($2, $4), emptyNode, "", $7)) 
        } | 
        TFunc TIdent '.' TIdent lparen ident_list ')' stats TEnd         {
            $$ = __store(NewSymbolFromToken($2), NewString($4.Str), __func(__markupFuncName($2, $4), $6, "", $8)) 
        } | 
        TFunc TIdent '.' TIdent lparen ident_list TDotDotDot ')' stats TEnd         {
            $$ = __store(NewSymbolFromToken($2), NewString($4.Str), __func(__markupFuncName($2, $4), __dotdotdot($6), "", $9)) 
        }

func_expr:
        TLambda lparen ')' stats TEnd                               { $$ = __func(__markupLambdaName($1), emptyNode, "", $4) } | 
        TLambda lparen ident_list ')' stats TEnd                    { $$ = __func(__markupLambdaName($1), $3, "", $5) } | 
        TLambda lparen ident_list TDotDotDot ')' stats TEnd         { $$ = __func(__markupLambdaName($1), __dotdotdot($3), "", $6) }

jmp_stat:
        TBreak {
            $$ = NewComplex(NewSymbol(ABreak)).SetPos($1.Pos) 
        } |
        TGoto TIdent {
            $$ = NewComplex(NewSymbol(AGoto), NewSymbolFromToken($2)).SetPos($1.Pos) 
        } |
        TLabel TIdent TLabel {
            $$ = NewComplex(NewSymbol(ALabel), NewSymbolFromToken($2)) 
        } |
        TReturnVoid {
            $$ = NewComplex(NewSymbol(AReturn), NewSymbol(ANil)).SetPos($1.Pos) 
        } |
        TReturn expr_list {
            if len($2.Nodes) == 1 {
                a := $2.Nodes[0]
                if len(a.Nodes) == 3 && a.Nodes[0].SymbolValue() == ACall { 
                    // return call(...) -> return tailcall(...)
                    a.Nodes[0].strSym = ATailCall
                }
                $$ = NewComplex(NewSymbol(AReturn), a).SetPos($1.Pos) 
            } else {
                $$ = NewComplex(NewSymbol(AReturn), NewComplex(NewSymbol(AArray), $2)).SetPos($1.Pos) 
            }
        }

declarator:
        TIdent {
            $$ = NewSymbolFromToken($1) 
        } |
        prefix_expr '[' expr ']' {
            $$ = __load($1, $3).SetPos($2.Pos) 
        } |
        prefix_expr '.' TIdent {
            $$ = __load($1, NewString($3.Str)).SetPos($2.Pos) 
        }

declarator_list:
        declarator {
            $$ = NewComplex($1) 
        } |
        declarator_list ',' declarator {
            $$ = $1.append($3) 
        }

ident_list:
        TIdent {
            $$ = NewComplex(NewSymbolFromToken($1)) 
        } | 
        ident_list ',' TIdent {
            $$ = $1.append(NewSymbolFromToken($3)) 
        }

expr:
        prefix_expr                       { $$ = $1 } |
        TNumber                           { $$ = NewNumberFromString($1.Str) } |
        TString                           { $$ = NewString($1.Str) } |
        '{' '}'                           { $$ = NewComplex(NewSymbol(AArrayMap), emptyNode).SetPos($1.Pos) } |
        '{' expr_list comma '}'           { $$ = NewComplex(NewSymbol(AArray), $2).SetPos($1.Pos) } |
        '{' expr_assign_list comma'}'     { $$ = NewComplex(NewSymbol(AArrayMap), $2).SetPos($1.Pos) } |
        expr TOr expr                     { $$ = NewComplex(NewSymbol(AOr), $1,$3).SetPos($2.Pos) } |
        expr TAnd expr                    { $$ = NewComplex(NewSymbol(AAnd), $1,$3).SetPos($2.Pos) } |
        expr '>' expr                     { $$ = NewComplex(NewSymbol(ALess), $3,$1).SetPos($2.Pos) } |
        expr '<' expr                     { $$ = NewComplex(NewSymbol(ALess), $1,$3).SetPos($2.Pos) } |
        expr TGte expr                    { $$ = NewComplex(NewSymbol(ALessEq), $3,$1).SetPos($2.Pos) } |
        expr TLte expr                    { $$ = NewComplex(NewSymbol(ALessEq), $1,$3).SetPos($2.Pos) } |
        expr TEqeq expr                   { $$ = NewComplex(NewSymbol(AEq), $1,$3).SetPos($2.Pos) } |
        expr TNeq expr                    { $$ = NewComplex(NewSymbol(ANeq), $1,$3).SetPos($2.Pos) } |
        expr '+' expr                     { $$ = NewComplex(NewSymbol(AAdd), $1,$3).SetPos($2.Pos) } |
        expr '-' expr                     { $$ = NewComplex(NewSymbol(ASub), $1,$3).SetPos($2.Pos) } |
        expr '*' expr                     { $$ = NewComplex(NewSymbol(AMul), $1,$3).SetPos($2.Pos) } |
        expr '/' expr                     { $$ = NewComplex(NewSymbol(ADiv), $1,$3).SetPos($2.Pos) } |
        expr TIDiv expr                   { $$ = NewComplex(NewSymbol(AIDiv), $1,$3).SetPos($2.Pos) } |
        expr '%' expr                     { $$ = NewComplex(NewSymbol(AMod), $1,$3).SetPos($2.Pos) } |
        expr '&' expr                     { $$ = NewComplex(NewSymbol(ABitAnd), $1,$3).SetPos($2.Pos) } |
        expr '|' expr                     { $$ = NewComplex(NewSymbol(ABitOr), $1,$3).SetPos($2.Pos) } |
        expr '^' expr                     { $$ = NewComplex(NewSymbol(ABitXor), $1,$3).SetPos($2.Pos) } |
        expr TLsh expr                    { $$ = NewComplex(NewSymbol(ABitLsh), $1,$3).SetPos($2.Pos) } |
        expr TRsh expr                    { $$ = NewComplex(NewSymbol(ABitRsh), $1,$3).SetPos($2.Pos) } |
        expr TURsh expr                   { $$ = NewComplex(NewSymbol(ABitURsh), $1,$3).SetPos($2.Pos) } |
        '~' expr %prec UNARY              { $$ = NewComplex(NewSymbol(ABitNot), $2).SetPos($1.Pos) } |
        TNot expr %prec UNARY             { $$ = NewComplex(NewSymbol(ANot), $2).SetPos($1.Pos) } |
        '-' expr %prec UNARY              { $$ = NewComplex(NewSymbol(ASub), zeroNode, $2).SetPos($1.Pos) } |
        '+' expr %prec UNARY              { $$ = NewComplex(NewSymbol(AAdd), zeroNode, $2).SetPos($1.Pos) } |
        func_expr                         {
            $$ = $1 
        }

prefix_expr:
        '(' expr ')'                      { $$ = $2 } |
        declarator {
            $$ = $1 
        } |
        prefix_expr TString {
            $$ = __call($1, NewComplex(NewString($2.Str))).SetPos($1.Pos()) 
        } |
        prefix_expr call_expr {
            $2.Nodes[1] = $1
            $$ = $2
        }

call_expr:
        TLParen ')' {
            $$ = __call(emptyNode, emptyNode).SetPos($1.Pos) 
        } |
        TLParen expr_list comma ')' {
            $$ = __call(emptyNode, $2).SetPos($1.Pos) 
        } |
        TLParen expr_list TDotDotDot comma ')' {
            $$ = __call(emptyNode, __dotdotdot($2)).SetPos($1.Pos) 
        }

expr_list:
        expr {
            $$ = NewComplex($1) 
        } |
        expr_list ',' expr {
            $$ = $1.append($3) 
        }

expr_assign_list:
        TIdent '=' expr {
            $$ = NewComplex(NewString($1.Str), $3) 
        } |
        '[' expr ']' '=' expr {
            $$ = NewComplex($2, $5) 
        } |
        expr_assign_list ',' TIdent '=' expr {
            $$ = $1.append(NewString($3.Str)).append($5) 
        } |
        expr_assign_list ',' '[' expr ']' '=' expr {
            $$ = $1.append($4).append($7) 
        }

comma: { $$ = emptyNode } | ',' { $$ = emptyNode }
lparen:
        TLParen { $$ = $1 } |
        '(' { $$ = $1 }

%%

