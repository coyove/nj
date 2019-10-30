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
%type<atom>  _postfix_assign
%type<expr> prefix_expr
%type<expr> _assign_stat
%type<expr> assign_stat
%type<expr> for_stat
%type<expr> if_stat
%type<expr> switch_stat
%type<expr> switch_body
%type<expr> oneline_or_block
%type<expr> jmp_stat
%type<expr> func_stat
%type<expr> flow_stat
%type<expr> func
%type<atom>  _func
%type<expr> func_call
%type<expr> func_args
%type<expr> function
%type<expr> func_params_list
%type<expr> map_gen
%type<expr> _map_gen

%union {
  token Token
  expr  *Node
  atom  Atom
}

/* Reserved words */
%token<token> TAssert TBreak TCase TContinue TElse TFor TFunc TIf TLen TNot TReturn TReturnNil TUse TSwitch TTypeof TVar TYield TYieldNil

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
            $$ = CompNode(AChain)
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

assign_stat:
        _assign_stat {
            if $1.isIsolatedCopy() {
                $1.Cx(2).C()[0] = zeroNode
            }
            $$ = $1
        }

stat:
        jmp_stat               { $$ = $1 } |
        flow_stat              { $$ = $1 } |
        assign_stat            { $$ = $1 } |
        block                  { $$ = $1 }

oneline_or_block:
        assign_stat            { $$ = ChainNode($1) } |
        jmp_stat               { $$ = ChainNode($1) } |
        for_stat               { $$ = ChainNode($1) } |
        if_stat                { $$ = ChainNode($1) } |
        switch_stat            { $$ = ChainNode($1) } |
        block                  { $$ = $1 }

flow_stat:
        for_stat               { $$ = $1 } |
        if_stat                { $$ = $1 } |
        switch_stat            { $$ = $1 } |
        func_stat              { $$ = $1 }

_assign_stat:
        prefix_expr            { $$ = $1 } |
        postfix_incdec         { $$ = $1 } |
        declarator '=' expr {
            $$ = CompNode("move", $1, $3).setPos0($1)
            if $1.Cn() > 0 && $1.Cx(0).A() == ALoad {
                $$ = CompNode("store", $1.Cx(1), $1.Cx(2), $3)
            }
            if c := $1.A(); c != "" && $1.Type() == Natom {
                // For 'a = a +/- n', we will simplify it as 'inc a +/- n'
                if a, b, s := $3.isSimpleAddSub(); a == c {
                    $3.Cx(2).Value = $3.Cx(2).N() * s
                    $$ = CompNode("inc", $1, $3.Cx(2))
                    $$.Cx(1).SetPos($1)
                } else if b == c {
                    $3.Cx(1).Value = $3.Cx(1).N() * s
                    $$ = CompNode("inc", $1, $3.Cx(1))
                    $$.Cx(1).SetPos($1)
                }
            }
            $$.Cx(0).SetPos($1)
        }

_postfix_incdec:
        TAddAdd { $$ = oneNode } |
        TSubSub { $$ = moneNode }

_postfix_assign:
        TAddEq  { $$ = AAdd } |
        TSubEq  { $$ = ASub } |
        TMulEq  { $$ = AMul } |
        TDivEq  { $$ = "/" } |
        TModEq  { $$ = "%" } |
        TAndEq  { $$ = "&" } |
        TOrEq   { $$ = "|" } |
        TXorEq  { $$ = "^" } |
        TLshEq  { $$ = "<<" } |
        TRshEq  { $$ = ">>" } |
        TURshEq { $$ = ">>>" }

postfix_incdec:
        TIdent _postfix_incdec                    { $$ = CompNode("inc", ANode($1).setPos($1), $2) } |
        TIdent _postfix_assign expr %prec ASSIGN  { $$ = CompNode("move", ANode($1), CompNode($2, ANode($1).setPos($1), $3)).setPos0($1) } |
        prefix_expr '[' expr ']' _postfix_incdec  { $$ = CompNode("store", $1, $3, CompNode(AAdd, CompNode(ALoad, $1, $3).setPos0($1), $5).setPos0($1)) } |
        prefix_expr '.' TIdent   _postfix_incdec  { $$ = CompNode("store", $1, $3, CompNode(AAdd, CompNode(ALoad, $1, $3).setPos0($1), $4).setPos0($1)) } |
        prefix_expr '[' expr ']' _postfix_assign expr %prec ASSIGN {
            $$ = CompNode("store", $1, $3, CompNode($5, CompNode(ALoad, $1, $3).setPos0($1), $6).setPos0($1))
        } |
        prefix_expr '.' TIdent   _postfix_assign expr %prec ASSIGN {
            $$ = CompNode("store", $1, $3, CompNode($4, CompNode(ALoad, $1, $3).setPos0($1), $5).setPos0($1))
        }

for_stat:
        TFor expr oneline_or_block {
            $$ = CompNode("for", $2, emptyNode, $3).setPos0($1)
        } |
        TFor ';' expr ';' oneline_or_block oneline_or_block {
            $$ = CompNode("for", $3, $5, $6).setPos0($1)
        } |
        TFor expr ';' expr ';' oneline_or_block oneline_or_block {
            $$ = ChainNode(
                $2,
                CompNode("for", $4, $6, $7).setPos0($1),
            )
        } |
        TFor TIdent '=' expr ',' expr oneline_or_block {
            vname, ename := ANode($2), ANodeS($2.Str + "_end")
            $$ = ChainNode(
                CompNode("move", vname, $4).setPos0($1),
                CompNode("move", ename, $6).setPos0($1),
                CompNode("for", 
                    CompNode(ALess, vname, ename).setPos0($1), 
                    ChainNode(
                        CompNode("inc", vname, oneNode).setPos0($1),
                    ), 
                    $7,
                ).setPos0($1),
            )
        } |
        TFor TIdent '=' expr ',' expr ',' expr oneline_or_block {
            vname, sname, ename := ANode($2), ANodeS($2.Str + "_step"), ANodeS($2.Str + "_end") 
            if $6.Type() == Nnumber {
                // easy case
                chain := ChainNode(CompNode("inc", vname, $6).setPos0($1))
                var cond *Node
                if $6.N() < 0 {
                    cond = CompNode("<=", ename, vname)
                } else {
                    cond = CompNode("<=", vname, ename)
                }
                $$ = ChainNode(
                    CompNode("move", vname, $4).setPos0($1),
                    CompNode("move", ename, $8).setPos0($1),
                    CompNode("for", cond, chain, $9).setPos0($1),
                )
            } else {
                bname := ANodeS($2.Str + "_begin")
                $$ = ChainNode(
                    CompNode("move", vname, $4).setPos0($1),
                    CompNode("move", bname, $4).setPos0($1),
                    CompNode("move", sname, $6).setPos0($1),
                    CompNode("move", ename, $8).setPos0($1),
                    CompNode("if", 
                        CompNode("<=",
                            zeroNode,
                            CompNode(AMul,
                                CompNode(ASub, ename, vname).setPos0($1),
                                sname,
                            ).setPos0($1),
                        ).setPos0($1),
                        ChainNode(
                            CompNode("for",
                                CompNode("<=",
                                    CompNode(AMul,
                                        CompNode(ASub, vname, bname).setPos0($1), 
                                        CompNode(ASub, vname, ename).setPos0($1),
                                    ),
                                    zeroNode,
                                ).setPos0($1),
                                ChainNode(
                                    CompNode("move", vname, CompNode(AAdd, vname, sname).setPos0($1),
                                ).setPos0($1)),
                                $9,
                            ).setPos0($1),
                        ),
                        CompNode(AChain),
                    ).setPos0($1),
                )
            }
            
        } |
        TFor expr ',' expr {
            $$ = CompNode("foreach", $2, $4).setPos0($1)
        } 

if_stat:
        TIf expr oneline_or_block %prec 'T'              { $$ = CompNode("if", $2, $3, emptyNode).setPos0($1) } |
        TIf expr oneline_or_block TElse oneline_or_block { $$ = CompNode("if", $2, $3, $5).setPos0($1) }

switch_stat:
        TSwitch expr '{' switch_body '}'         { $$ = expandSwitch($1, $2, $4.C()) }

switch_body:
        TCase expr ':' stats             { $$ = CompNode($2, $4).setPos0($1) } |
        TCase TElse ':' stats            { $$ = CompNode(ANode($2), $4).setPos0($1) } |
        switch_body TCase expr ':' stats { $$ = $1.Cappend($3, $5) } |
        switch_body TCase TElse ':' stats{ $$ = $1.Cappend(ANode($3), $5) }

_func:
        { $$ = Atom("") } |
        _func TFunc  { $$ = $1 + ",safe" } |
        _func TVar   { $$ = $1 + ",var" }

func:
        TFunc _func { $$ = ANodeS("func," + string($2)).setPos($1) }

func_stat:
        func TIdent func_params_list oneline_or_block {
            funcname := ANode($2)
            $$ = CompNode(
                AChain, 
                CompNode("set", funcname, nilNode).setPos0($2), 
                CompNode("move", funcname, 
                    CompNode($1, funcname, $3, $4).setPos0($2),
                ).setPos0($2),
            )
        }

jmp_stat:
        TYield expr           { $$ = CompNode("yield", $2).setPos0($1) } |
        TYieldNil             { $$ = CompNode("yield", CompNode("#", nilNode).setPos0($1)).setPos0($1) } |
        TBreak                { $$ = CompNode("break").setPos0($1) } |
        TContinue             { $$ = CompNode("continue").setPos0($1) } |
        TAssert expr          { $$ = CompNode("assert", $2, nilNode).setPos0($1) } |
        TAssert expr TString  { $$ = CompNode("assert", $2, NewNode($3.Str)).setPos0($1) } |
        TReturn expr          { $$ = CompNode("ret", $2).setPos0($1) } |
        TReturnNil            { $$ = CompNode("ret", CompNode("#", nilNode).setPos0($1)).setPos0($1) } |
        TUse TString          { $$ = yylex.(*Lexer).loadFile(filepath.Join(filepath.Dir($1.Pos.Source), $2.Str), $1) }

declarator:
        TIdent                                { $$ = ANode($1).setPos($1) } |
        TIdent TSquare                        { $$ = CompNode(ALoad, nilNode, $1.Str).setPos0($1) } |
        prefix_expr '[' expr ']'              { $$ = CompNode(ALoad, $1, $3).setPos0($3).setPos($3) } |
        prefix_expr '[' expr ':' expr ']'     { $$ = CompNode(ASlice, $1, $3, $5).setPos0($3).setPos($3) } |
        prefix_expr '[' expr ':' ']'          { $$ = CompNode(ASlice, $1, $3, moneNode).setPos0($3).setPos($3) } |
        prefix_expr '[' ':' expr ']'          { $$ = CompNode(ASlice, $1, zeroNode, $4).setPos0($4).setPos($4) } |
        prefix_expr '.' TIdent                { $$ = CompNode(ALoad, $1, NewNode($3.Str)).setPos0($3).setPos($3) }

ident_list:
        TIdent                                { $$ = CompNode($1.Str) } | 
        ident_list ',' TIdent                 { $$ = $1.Cappend(ANode($3)) }

expr_list:
        expr                                  { $$ = CompNode($1) } |
        expr_list ',' expr                    { $$ = $1.Cappend($3) }

expr_assign_list:
        expr ':' expr                         { $$ = CompNode($1, $3) } |
        expr_assign_list ',' expr ':' expr    { $$ = $1.Cappend($3).Cappend($5) }

expr:
        TNumber                { $$ = NewNumberNode($1.Str).SetPos($1) } |
        TUse TString           { $$ = yylex.(*Lexer).loadFile(filepath.Join(filepath.Dir($1.Pos.Source), $2.Str), $1) } |
        TTypeof expr           { $$ = CompNode("typeof", $2) } |
        TLen expr              { $$ = CompNode("len", $2) } |
        function               { $$ = $1 } |
        map_gen                { $$ = $1 } |
        prefix_expr            { $$ = $1 } |
        postfix_incdec         { $$ = $1 } |
        TString { $$ = NewNode($1.Str).SetPos($1) }  |
        expr TOr expr          { $$ = CompNode("or", $1,$3).setPos0($1) } |
        expr TAnd expr         { $$ = CompNode("and", $1,$3).setPos0($1) } |
        expr '>' expr          { $$ = CompNode(ALess, $3,$1).setPos0($1) } |
        expr '<' expr          { $$ = CompNode(ALess, $1,$3).setPos0($1) } |
        expr TGte expr         { $$ = CompNode("<=", $3,$1).setPos0($1) } |
        expr TLte expr         { $$ = CompNode("<=", $1,$3).setPos0($1) } |
        expr TEqeq expr        { $$ = CompNode("==", $1,$3).setPos0($1) } |
        expr TNeq expr         { $$ = CompNode("!=", $1,$3).setPos0($1) } |
        expr '+' expr          { $$ = CompNode(AAdd, $1,$3).setPos0($1) } |
        expr '-' expr          { $$ = CompNode(ASub, $1,$3).setPos0($1) } |
        expr '*' expr          { $$ = CompNode(AMul, $1,$3).setPos0($1) } |
        expr '/' expr          { $$ = CompNode("/", $1,$3).setPos0($1) } |
        expr '%' expr          { $$ = CompNode("%", $1,$3).setPos0($1) } |
        expr '^' expr          { $$ = CompNode("^", $1,$3).setPos0($1) } |
        expr TLsh expr         { $$ = CompNode("<<", $1,$3).setPos0($1) } |
        expr TRsh expr         { $$ = CompNode(">>", $1,$3).setPos0($1) } |
        expr TURsh expr        { $$ = CompNode(">>>", $1,$3).setPos0($1) } |
        expr '|' expr          { $$ = CompNode("|", $1,$3).setPos0($1) } |
        expr '&' expr          { $$ = CompNode("&", $1,$3).setPos0($1) } |
        '-' expr %prec UNARY   { $$ = CompNode(ASub, zeroNode, $2).setPos0($2) } |
        '~' expr %prec UNARY   { $$ = CompNode("^", $2, max32Node).setPos0($2) } |
        TNot expr %prec UNARY  { $$ = CompNode("!", $2).setPos0($2) } |
        '#' expr %prec UNARY   { $$ = CompNode("#", $2).setPos0($2) } |
        '&' TIdent %prec UNARY { $$ = CompNode("addressof", ANode($2)).setPos0($2) }

prefix_expr:
        declarator        { $$ = $1 } |
        '(' func_call ')' { $$ = $2 } |
        func_call         { $$ = $1 } |
        '(' expr ')'      { $$ = $2 }

func_call:
        prefix_expr func_args { $$ = CompNode("call", $1, $2).setPos0($1) }

func_args:
        '(' ')'           { $$ = emptyNode } |
        '(' expr_list ')' { $$ = $2 }

function:
        func func_params_list block %prec FUN { $$ = CompNode($1, "<a>", $2, $3).setPos0($1) } |
        func ident_list '=' expr %prec FUN    { $$ = CompNode($1, "<a>", $2, ChainNode(CompNode("ret", $4).setPos0($1))).setPos0($1) } |
        func '=' expr %prec FUN               { $$ = CompNode($1, "<a>", emptyNode, ChainNode(CompNode("ret", $3).setPos0($1))).setPos0($1) }

func_params_list:
        '(' ')'            { $$ = emptyNode } |
        '(' ident_list ')' { $$ = $2 }

map_gen:
        '{' '}'          { $$ = CompNode("map", emptyNode).setPos0($1) } |
        '{' _map_gen '}' { $$ = $2.setPos0($1) }

_map_gen:
        expr_assign_list     { $$ = CompNode("map", $1).setPos0($1) } |
        expr_assign_list ',' { $$ = CompNode("map", $1).setPos0($1) } |
        expr_list            { $$ = CompNode("array", $1).setPos0($1) } |
        expr_list ','        { $$ = CompNode("array", $1).setPos0($1) }

%%

func expandSwitch(switchTok Token, sub *Node, cases []*Node) *Node {
    subject := ANodeS("switch_tmp_var")
    ret := ChainNode(CompNode("set", subject, sub).setPos0(switchTok))

    var lastif, root *Node
    var defaultCase *Node
    
    for i := 0; i < len(cases); i+=2 {
        if cases[i].A() == "else" {
            defaultCase = cases[i + 1]
            continue
        }

        casestat := CompNode("if", CompNode("==", subject, cases[i]), cases[i + 1]).setPos0(cases[i])
        if lastif != nil {
            lastif.Cappend(ChainNode(casestat))
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
