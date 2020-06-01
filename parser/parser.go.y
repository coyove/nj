%{
package parser

import "strconv"
import "math/rand"

%}
%type<expr> stats
%type<expr> stat
%type<expr> declarator
%type<expr> declarator_list
%type<expr> ident_list
%type<expr> ident_dot_list
%type<expr> expr_list
%type<expr> expr_list_paren
%type<expr> expr_assign_list
%type<expr> expr
%type<expr> postfix_incdec
%type<atom> _postfix_assign
%type<expr> prefix_expr
%type<expr> assign_stat
%type<expr> for_stat
%type<expr> if_stat
%type<expr> elseif_stat
%type<expr> jmp_stat
%type<expr> flow_stat
%type<expr> func
%type<expr> func_stat
%type<expr> function
%type<expr> func_params_list
%type<expr> table_gen

%union {
  token Token
  expr  *Node
  atom  Symbol
}

/* Reserved words */
%token<token> TDo TIn TLocal TElseIf TThen TEnd TBreak TContinue TElse TFor TWhile TFunc TIf TLen TReturn TReturnVoid TImport TYield TYieldVoid TRepeat TUntil TNot

/* Literals */
%token<token> TEqeq TNeq TLte TGte TIdent TNumber TString '{' '[' '('
%token<token> TAddEq TSubEq TMulEq TDivEq TModEq
%token<token> TSquare TDotDotDot TDotDot TSet

/* Operators */
%right 'T'
%right TElse
%left ASSIGN
%right FUNC
%left TOr
%left TAnd
%left '>' '<' TGte TLte TEqeq TNeq
%left TDotDot
%left '+' '-' '^'
%left '*' '/' '%' 
%right UNARY /* not # -(unary) */
%right TImport

%% 

stats: 
        {
            $$ = __chain()
            if l, ok := yylex.(*Lexer); ok {
                l.Stmts = $$
            }
        } |
        stats stat {
            $$ = $1.CplAppend($2)
            if l, ok := yylex.(*Lexer); ok {
                l.Stmts = $$
            }
        }

stat:
        jmp_stat       { $$ = $1 } |
        flow_stat      { $$ = $1 } |
        assign_stat    { $$ = $1 } |
        TDo stats TEnd { $$ = __do($2) } |
        ';'            { $$ = emptyNode }

flow_stat:
        for_stat       { $$ = $1 } |
        if_stat        { $$ = $1 } |
        func_stat      { $$ = $1 }

_postfix_assign:
        TAddEq         { $$ = AAdd } |
        TSubEq         { $$ = ASub } |
        TMulEq         { $$ = AMul } |
        TDivEq         { $$ = ADiv } |
        TModEq         { $$ = AMod }

assign_stat:
        prefix_expr {
            $$ = $1
        } | 
        postfix_incdec {
            $$ = $1
        } |
        TLocal ident_list {
            $$ = __chain()
            for _, v := range $2.Cpl() {
                $$ = $$.CplAppend(__set(v, nilNode).pos0($1))
            }
        } |
        TLocal ident_list '=' expr_list {
            m, n := len($2.Cpl()), len($4.Cpl())
            for i := 0; i < m - n; i++ {
                $4.CplAppend(popvNode)
            }

            $$ = __chain()
            for i, v := range $2.Cpl() {
                $$ = $$.CplAppend(__set(v, $4.CplIndex(i)).pos0($1))
            }
        } |
        declarator_list '=' expr_list {
            nodes := $1.Cpl()
            m, n := len(nodes), len($3.Cpl())
            for i := 0; i < m - n; i++ {
                $3.CplAppend(popvNode)
            } 
             
            if head := nodes[0]; len(nodes) == 1 {
                // a0 = b0
                $$ = head.moveLoadStore(__move, $3.CplIndex(0)).pos0($1)
                if a, s := $3.CplIndex(0).isSimpleAddSub(); a != "" && a == head.Sym() {
                    $$ = __inc(head, Num(s)).pos0($1)
                }
            } else { 
                // a0, ..., an = b0, ..., bn
                $$ = __chain()
                names, retaddr := []*Node{}, Cpl(ARetAddr)
                for i := range nodes {
                    names = append(names, randomVarname())
                    retaddr.CplAppend(names[i])
                    $$.CplAppend(__set(names[i], $3.CplIndex(i)).pos0($1))
                }
                for i, v := range nodes {
                    $$.CplAppend(v.moveLoadStore(__move, names[i]).pos0($1))
                }
                $$.CplAppend(retaddr)
            }
        } 

postfix_incdec:
        TIdent _postfix_assign expr %prec ASSIGN  {
            $$ = __move(SymTok($1), Cpl($2, SymTok($1).pos0($1), $3)).pos0($1)
        } |
        prefix_expr '[' expr ']' _postfix_assign expr %prec ASSIGN {
            $$ = __store($1, $3, Cpl($5, __load($1, $3).pos0($1), $6).pos0($1))
        } |
        prefix_expr '.' TIdent _postfix_assign expr %prec ASSIGN {
            $$ = __store($1, Nod($3.Str), Cpl($4, __load($1, Nod($3.Str)).pos0($1), $5).pos0($1))
        }

for_stat:
        TWhile expr TDo stats TEnd {
            $$ = __loop(__if($2, $4, breakNode).pos0($1)).pos0($1)
        } |
        TRepeat stats TUntil expr {
            $$ = __loop(
                __chain(
                    $2,
                    __if($4, breakNode, emptyNode).pos0($1),
                ).pos0($1),
            ).pos0($1)
        } |
        TFor TIdent TIn expr TDo stats TEnd {
            iter := randomVarname()
            $$ = __do(
                __set(iter, $4).pos0($1),
                __loop(
                    __chain(
                        __set($2, __call(iter, emptyNode).pos0($1)).pos0($1),
                        __if($2, $6, breakNode).pos0($1),
                    ),
                ).pos0($1),
            )
        } |
        TFor TIdent ',' TIdent TIn expr TDo stats TEnd {
            iter := randomVarname()
            $$ = __do(
                __set(iter, $6).pos0($1),
                __loop(
                    __chain(
                        __set($2, __call(iter, emptyNode).pos0($1)).pos0($1),
                        __set($4, popvNode).pos0($1),
                        __if($2, $8, breakNode).pos0($1),
                    ),
                ).pos0($1),
            )
        } |
        TFor TIdent '=' expr ',' expr TDo stats TEnd {
            forVar, forEnd := SymTok($2), randomVarname()
            $$ = __do(
                    __set(forVar, $4).pos0($1),
                    __set(forEnd, $6).pos0($1),
                    __loop(
                        __if(
                            __lessEq(forVar, forEnd),
                            __chain($8, __inc(forVar, oneNode).pos0($1)),
                            breakNode,
                        ).pos0($1),
                    ).pos0($1),
                )
        } |
        TFor TIdent '=' expr ',' expr ',' expr TDo stats TEnd {
            forVar, forEnd := SymTok($2), randomVarname()
            if $8.Type() == NUM { // step is a static number, easy case
                var cond *Node
                if $8.Num() < 0 {
                    cond = __lessEq(forEnd, forVar)
                } else {
                    cond = __lessEq(forVar, forEnd)
                }
                $$ = __do(
                    __set(forVar, $4).pos0($1),
                    __set(forEnd, $6).pos0($1),
                    __loop(
                        __chain(
                            __if(
                                cond,
                                __chain($10, __inc(forVar, $8)),
                                breakNode,
                            ).pos0($1),
                        ),
                    ).pos0($1),
                )
            } else { 
                forStep := randomVarname()
                $$ = __do(
                    __set(forVar, $4).pos0($1),
                    __set(forEnd, $6).pos0($1),
                    __set(forStep, $8).pos0($1),
                    __loop(
                        __chain(
                            __if(
                                __less(zeroNode, forStep).pos0($1),
                                // +step
                                __if(__less(forEnd, forVar), breakNode, emptyNode).pos0($1),
                                // -step
                                __if(__less(forVar, forEnd), breakNode, emptyNode).pos0($1),
                            ).pos0($1),
                            $10,
                            __inc(forVar, forStep),
                        ),
                    ).pos0($1),
                )
            }
            
        } 

if_stat:
        TIf expr TThen stats elseif_stat TEnd %prec 'T' {
            $$ = __if($2, $4, $5).pos0($1)
        }

elseif_stat:
        {
            $$ = Cpl()
        } |
        TElse stats {
            $$ = $2
        } |
        TElseIf expr TThen stats elseif_stat {
            $$ = __if($2, $4, $5).pos0($1)
        }

func:
        TFunc {
            $$ = Nod(AMove).SetPos($1)
        } |
        TLocal TFunc {
            $$ = Nod(ASet).SetPos($1)
        }

func_stat:
        func TIdent func_params_list stats TEnd {
            funcname := SymTok($2)
            $$ = __chain(
                opSetMove($1)(funcname, nilNode).pos0($2), 
                __move(funcname, __func($3, $4).pos0($2)).pos0($2),
            )
        } |
        func ident_dot_list '.' TIdent func_params_list stats TEnd {
            $$ = __store($2, Nod($4.Str), __func($5, $6).pos0($4)).pos0($4) 
        } |
        func ident_dot_list ':' TIdent func_params_list stats TEnd {
            paramlist := $5.CplPrepend(Sym("self"))
            $$ = __store(
                $2, Nod($4.Str), __func(paramlist, $6).pos0($4),
            ).pos0($4) 
        }

function:
        func func_params_list stats TEnd %prec FUNC {
            $$ = __func($2, $3).pos0($1).SetPos($1) 
        }

func_params_list:
        '(' ')'                           { $$ = Cpl() } |
        '(' ident_list ')'                { $$ = $2 }

jmp_stat:
        TYield expr_list                  { $$ = Cpl(AYield, $2).pos0($1) } |
        TYieldVoid                        { $$ = Cpl(AYield, emptyNode).pos0($1) } |
        TBreak                            { $$ = Cpl(ABreak).pos0($1) } |
        TContinue                         { $$ = Cpl(AContinue).pos0($1) } |
        TReturn expr_list                 { $$ = Cpl(AReturn, $2).pos0($1) } |
        TReturnVoid                       { $$ = Cpl(AReturn, emptyNode).pos0($1) } |
        TImport TString                   { $$ = __move(Sym(moduleNameFromPath($2.Str)), yylex.(*Lexer).loadFile(joinSourcePath($1.Pos.Source, $2.Str), $1)).pos0($1) }

declarator:
        TIdent                            { $$ = SymTok($1).SetPos($1) } |
        prefix_expr '[' expr ']'          { $$ = __load($1, $3).pos0($3).SetPos($3) } |
        prefix_expr '.' TIdent            { $$ = __load($1, Nod($3.Str)).pos0($3).SetPos($3) }

declarator_list:
        declarator                        { $$ = Cpl($1) } |
        declarator_list ',' declarator    { $$ = $1.CplAppend($3) }

ident_list:
        TIdent                            { $$ = Cpl($1.Str) } | 
        TDotDotDot                        { $$ = Cpl("...") } | 
        ident_list ',' TIdent             { $$ = $1.CplAppend(SymTok($3)) } |
        ident_list ',' TDotDotDot         { $$ = $1.CplAppend(Sym("...").SetPos($3)) }

ident_dot_list:
        TIdent                            { $$ = SymTok($1) } | 
        ident_dot_list '.' TIdent         { $$ = __load($1, Nod($3.Str)).pos0($3) }

expr:
        TDotDotDot                        { $$ = __call(Sym("unpack"), Cpl(Sym("arg"))).pos0($1) } |
        TNumber                           { $$ = Num($1.Str).SetPos($1) } |
        TImport TString                   { $$ = yylex.(*Lexer).loadFile(joinSourcePath($1.Pos.Source, $2.Str), $1) } |
        function                          { $$ = $1 } |
        table_gen                         { $$ = $1 } |
        prefix_expr                       { $$ = $1 } |
        TString                           { $$ = Nod($1.Str).SetPos($1) }  |
        expr TOr expr                     { $$ = Cpl(AOr, $1,$3).pos0($1) } |
        expr TAnd expr                    { $$ = Cpl(AAnd, $1,$3).pos0($1) } |
        expr '>' expr                     { $$ = Cpl(ALess, $3,$1).pos0($1) } |
        expr '<' expr                     { $$ = Cpl(ALess, $1,$3).pos0($1) } |
        expr TGte expr                    { $$ = Cpl(ALessEq, $3,$1).pos0($1) } |
        expr TLte expr                    { $$ = Cpl(ALessEq, $1,$3).pos0($1) } |
        expr TEqeq expr                   { $$ = Cpl(AEq, $1,$3).pos0($1) } |
        expr TNeq expr                    { $$ = Cpl(ANeq, $1,$3).pos0($1) } |
        expr '+' expr                     { $$ = Cpl(AAdd, $1,$3).pos0($1) } |
        expr TDotDot expr                 { $$ = Cpl(AConcat, $1,$3).pos0($1) } |
        expr '-' expr                     { $$ = Cpl(ASub, $1,$3).pos0($1) } |
        expr '*' expr                     { $$ = Cpl(AMul, $1,$3).pos0($1) } |
        expr '/' expr                     { $$ = Cpl(ADiv, $1,$3).pos0($1) } |
        expr '%' expr                     { $$ = Cpl(AMod, $1,$3).pos0($1) } |
        expr '^' expr                     { $$ = Cpl(APow, $1,$3).pos0($1) } |
        '-' expr %prec UNARY              { $$ = Cpl(ASub, zeroNode, $2).pos0($2) } |
        TNot expr %prec UNARY             { $$ = Cpl(ANot, $2).pos0($2) } |
        '&' TIdent %prec UNARY            { $$ = Cpl(AAddrOf, SymTok($2)).pos0($2) } |
        '#' expr %prec UNARY              { $$ = Cpl(ALen, $2) }

prefix_expr:
        declarator                        { $$ = $1 } |
        prefix_expr TString               { $$ = __call($1, Cpl(Nod($2.Str))).pos0($1) } |
        TIdent ':' TIdent expr_list_paren { $$ = __call(__load($1, Nod($3.Str)).pos0($1), $4.CplPrepend(SymTok($1))).pos0($1) } |
        prefix_expr expr_list_paren       { $$ = __call($1, $2).pos0($1) } |
        '(' expr ')'                      { $$ = $2 } // shift/reduce conflict

expr_list:
        expr                              { $$ = Cpl($1) } |
        expr_list ',' expr                { $$ = $1.CplAppend($3) }

expr_list_paren:
        '(' ')'                           { $$ = Cpl() } |
        '(' expr_list ')'                 { $$ = $2 } |
        table_gen                         { $$ = Cpl($1) }

expr_assign_list:
        TIdent '=' expr                            { $$ = Cpl(Nod($1.Str), $3) } |
        '[' expr ']' '=' expr                      { $$ = Cpl($2, $5) } |
        expr_assign_list ',' TIdent '=' expr       { $$ = $1.CplAppend(Nod($3.Str), $5) } |
        expr_assign_list ',' '[' expr ']' '=' expr { $$ = $1.CplAppend($4, $7) }

table_gen:
        '{' '}'                                    { $$ = Cpl(AArray, emptyNode).pos0($1) } |
        '{' expr_assign_list     '}'               { $$ = Cpl(AHash, $2).pos0($1) } |
        '{' expr_assign_list ',' '}'               { $$ = Cpl(AHash, $2).pos0($1) } |
        '{' expr_assign_list ';' expr_list '}'     { $$ = Cpl(AHashArray, $2, $4).pos0($1) } |
        '{' expr_assign_list ';' expr_list ',' '}' { $$ = Cpl(AHashArray, $2, $4).pos0($1) } |
        '{' expr_list            '}'               { $$ = Cpl(AArray, $2).pos0($1) } |
        '{' expr_list ';' expr_assign_list '}'     { $$ = Cpl(AHashArray, $4, $2).pos0($1) } |
        '{' expr_list ';' expr_assign_list ',' '}' { $$ = Cpl(AHashArray, $4, $2).pos0($1) } |
        '{' expr_list ','        '}'               { $$ = Cpl(AArray, $2).pos0($1) }

%%

func opSetMove(op *Node) func(dest, src interface{}) *Node {
    if op.Sym() == ASet {
        return __set
    }
    return __move
}

func randomVarname() *Node {
    return Sym("v" + strconv.FormatInt(rand.Int63(), 10))
}
