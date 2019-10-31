%{
package parser

import (
    "path/filepath"
)

%}
%type<expr> stats
%type<expr> block
%type<expr> stat
%type<expr> declarator
%type<expr> ident_list
%type<expr> expr_list
%type<expr> expr_assign_list
%type<expr> expr
%type<expr> postfix_incdec
%type<expr> _postfix_incdec
%type<atom> _postfix_assign
%type<expr> prefix_expr
%type<expr> assign_stat
%type<expr> for_stat
%type<expr> if_stat
%type<expr> switch_stat
%type<expr> switch_body
%type<expr> oneline_or_block
%type<expr> jmp_stat
%type<expr> func_stat
%type<expr> flow_stat
%type<expr> function
%type<expr> func_params_list
%type<expr> map_gen

%union {
  token Token
  expr  *Node
  atom  Atom
}

/* Reserved words */
%token<token> TAssert TBreak TCase TContinue TElse TFor TFunc TIf TLen TReturn TReturnNil TUse TSwitch TTypeof TYield TYieldNil

/* Literals */
%token<token> TAddAdd TSubSub TEqeq TNeq TLsh TRsh TURsh TLte TGte TIdent TNumber TString '{' '[' '('
%token<token> TAddEq TSubEq TMulEq TDivEq TModEq TAndEq TOrEq TXorEq TLshEq TRshEq TURshEq
%token<token> TSquare

/* Operators */
%right 'T'
%right TElse

%left ASSIGN
%right FUN
%left TOr
%left TAnd
%left '>' '<' TGte TLte TEqeq TNeq
%left '+' '-' '|' '^'
%left '*' '/' '%' TLsh TRsh TURsh '&'
%right UNARY /* not # -(unary) */
%right '~'
%right '#'
%left TAddAdd TMinMin
%right TTypeof, TLen, TUse

%% 

stats: 
        {
            $$ = __chain()
            if l, ok := yylex.(*Lexer); ok {
                l.Stmts = $$
            }
        } |
        stats stat {
            $$ = $1.Cappend($2)
            if l, ok := yylex.(*Lexer); ok {
                l.Stmts = $$
            }
        }

block: 
        '{' stats '}' { $$ = $2 }

stat:
        jmp_stat               { $$ = $1 } |
        flow_stat              { $$ = $1 } |
        assign_stat            { $$ = $1 } |
        block                  { $$ = $1 } |
        ';'                    { $$ = emptyNode }

oneline_or_block:
        assign_stat            { $$ = __chain($1) } |
        jmp_stat               { $$ = __chain($1) } |
        for_stat               { $$ = __chain($1) } |
        if_stat                { $$ = __chain($1) } |
        switch_stat            { $$ = __chain($1) } |
        block                  { $$ = $1 }

flow_stat:
        for_stat               { $$ = $1 } |
        if_stat                { $$ = $1 } |
        switch_stat            { $$ = $1 } |
        func_stat              { $$ = $1 }

assign_stat:
        prefix_expr            { $$ = $1 } |
        postfix_incdec         { $$ = $1 } |
        declarator '=' expr {
            $$ = __move($1, $3).pos0($1)
            if $1.Cn() > 0 && $1.Cx(0).A() == ALoad {
                $$ = __store($1.Cx(1), $1.Cx(2), $3).pos0($1)
            }
            if c := $1.A(); c != "" && $1.Type() == Natom {
                // For 'a = a +/- n', we will simplify it as 'inc a +/- n'
                if a, b, s := $3.isSimpleAddSub(); a == c {
                    $3.Cx(2).Value = $3.Cx(2).N() * s
                    $$ = __inc($1, $3.Cx(2)).pos0($1)
                } else if b == c {
                    $3.Cx(1).Value = $3.Cx(1).N() * s
                    $$ = __inc($1, $3.Cx(1)).pos0($1)
                }
            }
        }

_postfix_incdec:
        TAddAdd { $$ = oneNode } |
        TSubSub { $$ = moneNode }

_postfix_assign:
        TAddEq  { $$ = AAdd } |
        TSubEq  { $$ = ASub } |
        TMulEq  { $$ = AMul } |
        TDivEq  { $$ = ADiv } |
        TModEq  { $$ = AMod } |
        TAndEq  { $$ = ABitAnd } |
        TOrEq   { $$ = ABitOr } |
        TXorEq  { $$ = ABitXor } |
        TLshEq  { $$ = ABitLsh } |
        TRshEq  { $$ = ABitRsh } |
        TURshEq { $$ = ABitURsh }

postfix_incdec:
        TIdent _postfix_incdec                    { $$ = __inc(ANode($1), $2).pos0($1) } |
        TIdent _postfix_assign expr %prec ASSIGN  { $$ = __move(ANode($1), CompNode($2, ANode($1).setPos($1), $3)).pos0($1) } |
        prefix_expr '[' expr ']' _postfix_incdec  { $$ = __store($1, $3, CompNode(AAdd, __load($1, $3).pos0($1), $5).pos0($1)) } |
        prefix_expr '.' TIdent   _postfix_incdec  { $$ = __store($1, $3, CompNode(AAdd, __load($1, $3).pos0($1), $4).pos0($1)) } |
        prefix_expr '[' expr ']' _postfix_assign expr %prec ASSIGN {
            $$ = __store($1, $3, CompNode($5, __load($1, $3).pos0($1), $6).pos0($1))
        } |
        prefix_expr '.' TIdent _postfix_assign expr %prec ASSIGN {
            $$ = __store($1, $3, CompNode($4, __load($1, $3).pos0($1), $5).pos0($1))
        }

for_stat:
        TFor expr oneline_or_block {
            $$ = __for($2).__continue(emptyNode).__body($3).pos0($1)
        } |
        TFor ';' expr ';' oneline_or_block oneline_or_block {
            $$ = __for($3).__continue($5).__body($6).pos0($1)
        } |
        TFor expr ';' expr ';' oneline_or_block oneline_or_block {
            $$ = __chain(
                $2,
                __for($4).__continue($6).__body($7).pos0($1),
            )
        } |
        TFor TIdent '=' expr ',' expr oneline_or_block {
            forVar, forEnd := ANode($2), ANodeS($2.Str + "_end")
            $$ = __chain(
                __move(forVar, $4).pos0($1),
                __move(forEnd, $6).pos0($1),
                __for(
                    CompNode(ALess, forVar, forEnd).pos0($1),
                ).
                __continue(
                    __chain(__inc(forVar, oneNode).pos0($1)),
                ).
                __body($7).pos0($1),
            )
        } |
        TFor TIdent '=' expr ',' expr ',' expr oneline_or_block {
            forVar, forEnd := ANode($2), ANodeS($2.Str + "_end") 
            if $6.Type() == Nnumber { // easy case
                var cond *Node
                if $6.N() < 0 {
                    cond = __lessEq(forEnd, forVar)
                } else {
                    cond = __lessEq(forVar, forEnd)
                }
                $$ = __chain(
                    __move(forVar, $4).pos0($1),
                    __move(forEnd, $8).pos0($1),
                    __for(cond).
                    __continue(__chain(__inc(forVar, $6).pos0($1))).
                    __body($9).pos0($1),
                )
            } else {
                forStep := ANodeS($2.Str + "_step")
                forBegin := ANodeS($2.Str + "_begin")
                $$ = __chain(
                    __move(forVar, $4).pos0($1),
                    __move(forBegin, $4).pos0($1),
                    __move(forStep, $6).pos0($1),
                    __move(forEnd, $8).pos0($1),
                    __if(
                        __lessEq(
                            zeroNode,
                            __mul(
                                __sub(forEnd, forVar).pos0($1),
                                forStep,
                            ).pos0($1),
                        ).pos0($1),
                    ).
                    __then(
                        __chain(
                            __for(
                                __lessEq(
                                    __mul(
                                        __sub(forVar, forBegin).pos0($1), 
                                        __sub(forVar, forEnd).pos0($1),
                                    ),
                                    zeroNode,
                                ).pos0($1), // (forVar - forBegin) * (forVar - forEnd) <= 0
                            ).
                            __continue(
                                __chain(__inc(forVar, forStep).pos0($1)),
                            ).
                            __body($9).pos0($1),
                        ),
                    ).
                    __else(
                        emptyNode,
                    ).pos0($1),
                )
            }
            
        } |
        TFor expr ',' expr {
            $$ = CompNode(AForeach, $2, $4).pos0($1)
        } 

if_stat:
        TIf expr oneline_or_block %prec 'T'              { $$ = __if($2).__then($3).__else(emptyNode).pos0($1) } |
        TIf expr oneline_or_block TElse oneline_or_block { $$ = __if($2).__then($3).__else($5).pos0($1) }

switch_stat:
        TSwitch expr '{' switch_body '}'         { $$ = expandSwitch($1, $2, $4.C()) }

switch_body:
        TCase expr ':' stats             { $$ = CompNode($2, $4).pos0($1) } |
        TCase TElse ':' stats            { $$ = CompNode(ANode($2), $4).pos0($1) } |
        switch_body TCase expr ':' stats { $$ = $1.Cappend($3, $5) } |
        switch_body TCase TElse ':' stats{ $$ = $1.Cappend(ANode($3), $5) }

func_stat:
        TFunc TIdent func_params_list oneline_or_block {
            funcname := ANode($2)
            $$ = __chain(
                __set(funcname, nilNode).pos0($2), 
                __move(funcname, CompNode(AFunc, funcname, $3, $4).pos0($2)).pos0($2),
            )
        }

jmp_stat:
        TYield expr           { $$ = CompNode(AYield, $2).pos0($1) } |
        TYieldNil             { $$ = CompNode(AYield, CompNode(APop, nilNode).pos0($1)).pos0($1) } |
        TBreak                { $$ = CompNode(ABreak).pos0($1) } |
        TContinue             { $$ = CompNode(AContinue).pos0($1) } |
        TAssert expr          { $$ = CompNode(AAssert, $2, nilNode).pos0($1) } |
        TAssert expr TString  { $$ = CompNode(AAssert, $2, NewNode($3.Str)).pos0($1) } |
        TReturn expr          { $$ = __return($2).pos0($1) } |
        TReturnNil            { $$ = __return(CompNode(APop, nilNode).pos0($1)).pos0($1) } |
        TUse TString          { $$ = yylex.(*Lexer).loadFile(filepath.Join(filepath.Dir($1.Pos.Source), $2.Str), $1) }

declarator:
        TIdent                                { $$ = ANode($1).setPos($1) } |
        TIdent TSquare                        { $$ = __load(nilNode, $1).pos0($1) } |
        prefix_expr '[' expr ']'              { $$ = __load($1, $3).pos0($3).setPos($3) } |
        prefix_expr '.' TIdent                { $$ = __load($1, NewNode($3.Str)).pos0($3).setPos($3) } |
        prefix_expr '[' expr ':' expr ']'     { $$ = CompNode(ASlice, $1, $3, $5).pos0($3).setPos($3) } |
        prefix_expr '[' expr ':' ']'          { $$ = CompNode(ASlice, $1, $3, moneNode).pos0($3).setPos($3) } |
        prefix_expr '[' ':' expr ']'          { $$ = CompNode(ASlice, $1, zeroNode, $4).pos0($4).setPos($4) }

ident_list:
        TIdent                                { $$ = CompNode($1.Str) } | 
        ident_list ',' TIdent                 { $$ = $1.Cappend(ANode($3)) }

expr:
        TNumber                { $$ = NewNumberNode($1.Str).SetPos($1) } |
        TUse TString           { $$ = yylex.(*Lexer).loadFile(filepath.Join(filepath.Dir($1.Pos.Source), $2.Str), $1) } |
        TTypeof expr           { $$ = CompNode("typeof", $2) } |
        TLen expr              { $$ = CompNode("len", $2) } |
        function               { $$ = $1 } |
        map_gen                { $$ = $1 } |
        prefix_expr            { $$ = $1 } |
        postfix_incdec         { $$ = $1 } |
        TString                { $$ = NewNode($1.Str).SetPos($1) }  |
        expr TOr expr          { $$ = CompNode(AOr, $1,$3).pos0($1) } |
        expr TAnd expr         { $$ = CompNode(AAnd, $1,$3).pos0($1) } |
        expr '>' expr          { $$ = CompNode(ALess, $3,$1).pos0($1) } |
        expr '<' expr          { $$ = CompNode(ALess, $1,$3).pos0($1) } |
        expr TGte expr         { $$ = CompNode(ALessEq, $3,$1).pos0($1) } |
        expr TLte expr         { $$ = CompNode(ALessEq, $1,$3).pos0($1) } |
        expr TEqeq expr        { $$ = CompNode(AEq, $1,$3).pos0($1) } |
        expr TNeq expr         { $$ = CompNode(ANeq, $1,$3).pos0($1) } |
        expr '+' expr          { $$ = CompNode(AAdd, $1,$3).pos0($1) } |
        expr '-' expr          { $$ = CompNode(ASub, $1,$3).pos0($1) } |
        expr '*' expr          { $$ = CompNode(AMul, $1,$3).pos0($1) } |
        expr '/' expr          { $$ = CompNode(ADiv, $1,$3).pos0($1) } |
        expr '%' expr          { $$ = CompNode(AMod, $1,$3).pos0($1) } |
        expr '^' expr          { $$ = CompNode(ABitXor, $1,$3).pos0($1) } |
        expr TLsh expr         { $$ = CompNode(ABitLsh, $1,$3).pos0($1) } |
        expr TRsh expr         { $$ = CompNode(ABitRsh, $1,$3).pos0($1) } |
        expr TURsh expr        { $$ = CompNode(ABitURsh, $1,$3).pos0($1) } |
        expr '|' expr          { $$ = CompNode(ABitOr, $1,$3).pos0($1) } |
        expr '&' expr          { $$ = CompNode(ABitAnd, $1,$3).pos0($1) } |
        '~' expr %prec UNARY   { $$ = CompNode(ABitXor, $2, max32Node).pos0($2) } |
        '-' expr %prec UNARY   { $$ = CompNode(ASub, zeroNode, $2).pos0($2) } |
        '!' expr %prec UNARY   { $$ = CompNode(ANot, $2).pos0($2) } |
        '#' expr %prec UNARY   { $$ = CompNode(APop, $2).pos0($2) } |
        '&' TIdent %prec UNARY { $$ = CompNode(AAddrOf, ANode($2)).pos0($2) }

prefix_expr:
        declarator                    { $$ = $1 } |
        prefix_expr '(' ')'           { $$ = __call($1, emptyNode).pos0($1) } |
        prefix_expr '(' expr_list ')' { $$ = __call($1, $3).pos0($1) } |
        '(' expr ')'                  { $$ = $2 } // shift/reduce conflict

function:
        TFunc func_params_list block %prec FUN { $$ = CompNode(AFunc, "<a>", $2, $3).pos0($1).SetPos($1) } |
        TFunc ident_list '=' expr %prec FUN    { $$ = CompNode(AFunc, "<a>", $2, __chain(__return($4).pos0($1))).pos0($1).SetPos($1) } |
        TFunc '=' expr %prec FUN               { $$ = CompNode(AFunc, "<a>", emptyNode, __chain(__return($3).pos0($1))).pos0($1).SetPos($1) }

func_params_list:
        '(' ')'            { $$ = emptyNode } |
        '(' ident_list ')' { $$ = $2 }

expr_list:
        expr                                  { $$ = CompNode($1) } |
        expr_list ',' expr                    { $$ = $1.Cappend($3) }

expr_assign_list:
        expr ':' expr                         { $$ = CompNode($1, $3) } |
        expr_assign_list ',' expr ':' expr    { $$ = $1.Cappend($3).Cappend($5) }

map_gen:
        '{' '}'                      { $$ = CompNode(AMap, emptyNode).pos0($1) } |
        '{' expr_assign_list     '}' { $$ = CompNode(AMap, $2).pos0($2) } |
        '{' expr_assign_list ',' '}' { $$ = CompNode(AMap, $2).pos0($2) } |
        '{' expr_list            '}' { $$ = CompNode(AArray, $2).pos0($2) } |
        '{' expr_list ','        '}' { $$ = CompNode(AArray, $2).pos0($2) }

%%

func expandSwitch(switchTok Token, sub *Node, cases []*Node) *Node {
    subject := ANodeS("switch_tmp_var")
    ret := __chain(__set(subject, sub).pos0(switchTok))

    var lastif, root *Node
    var defaultCase *Node
    
    for i := 0; i < len(cases); i+=2 {
        if cases[i].A() == "else" {
            defaultCase = cases[i + 1]
            continue
        }

        casestat := CompNode(AIf, CompNode(AEq, subject, cases[i]), cases[i + 1]).pos0(cases[i])
        if lastif != nil {
            lastif.Cappend(__chain(casestat))
        } else {
            root = casestat
        }
        lastif = casestat
    }

    if defaultCase == nil {
        lastif.Cappend(CompNode(AChain))
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
