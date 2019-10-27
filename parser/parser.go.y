%{
package parser

import (
    "path/filepath"
    "github.com/coyove/common/rand"
    "fmt"
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
%type<str>  _postfix_assign
%type<expr> string
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
%type<str>  _func
%type<expr> func_call
%type<expr> func_args
%type<expr> function
%type<expr> func_params_list
%type<expr> map_gen
%type<expr> _map_gen

%union {
  token Token
  expr  *Node
  str   string
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
            $$ = CNode("chain")
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
                $1.Cx(2).C()[0] = NNode(0.0)
            }
            $$ = $1
        }

stat:
        jmp_stat               { $$ = $1 } |
        flow_stat              { $$ = $1 } |
        assign_stat            { $$ = $1 } |
        block                  { $$ = $1 }

oneline_or_block:
        assign_stat            { $$ = CNode("chain", $1) } |
        jmp_stat               { $$ = CNode("chain", $1) } |
        for_stat               { $$ = CNode("chain", $1) } |
        if_stat                { $$ = CNode("chain", $1) } |
        switch_stat            { $$ = CNode("chain", $1) } |
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
            $$ = CNode("move", $1, $3).setPos0($1)
            if $1.Cn() > 0 && $1.Cx(0).S() == "load" {
                $$ = CNode("store", $1.Cx(1), $1.Cx(2), $3)
            }
            if c := $1.S(); c != "" && $1.Type == Natom {
                // For 'a = a +/- n', we will simplify it as 'inc a +/- n'
                if a, b, s := $3.isSimpleAddSub(); a == c {
                    $3.Cx(2).Value = $3.Cx(2).N() * s
                    $$ = CNode("inc", $1, $3.Cx(2))
                    $$.Cx(1).SetPos($1)
                } else if b == c {
                    $3.Cx(1).Value = $3.Cx(1).N() * s
                    $$ = CNode("inc", $1, $3.Cx(1))
                    $$.Cx(1).SetPos($1)
                }
            }
            $$.Cx(0).SetPos($1)
        }

_postfix_incdec:
        TAddAdd { $$ = NNode(1.0) } |
        TSubSub { $$ = NNode(-1.0) }

_postfix_assign:
        TAddEq  { $$ = "+" } |
        TSubEq  { $$ = "-" } |
        TMulEq  { $$ = "*" } |
        TDivEq  { $$ = "/" } |
        TModEq  { $$ = "%" } |
        TAndEq  { $$ = "&" } |
        TOrEq   { $$ = "|" } |
        TXorEq  { $$ = "^" } |
        TLshEq  { $$ = "<<" } |
        TRshEq  { $$ = ">>" } |
        TURshEq { $$ = ">>>" }

postfix_incdec:
        TIdent _postfix_incdec                    { $$ = CNode("inc", ANode($1).setPos($1), $2) } |
        TIdent _postfix_assign expr %prec ASSIGN  { $$ = CNode("move", ANode($1), CNode($2, ANode($1).setPos($1), $3)).setPos0($1) } |
        prefix_expr '[' expr ']' _postfix_incdec  { $$ = CNode("store", $1, $3, CNode("+", CNode("load", $1, $3).setPos0($1), $5).setPos0($1)) } |
        prefix_expr '.' TIdent   _postfix_incdec  { $$ = CNode("store", $1, $3, CNode("+", CNode("load", $1, $3).setPos0($1), $4).setPos0($1)) } |
        prefix_expr '[' expr ']' _postfix_assign expr %prec ASSIGN {
            $$ = CNode("store", $1, $3, CNode($5, CNode("load", $1, $3).setPos0($1), $6).setPos0($1))
        } |
        prefix_expr '.' TIdent   _postfix_assign expr %prec ASSIGN {
            $$ = CNode("store", $1, $3, CNode($4, CNode("load", $1, $3).setPos0($1), $5).setPos0($1))
        }

for_stat:
        TFor expr oneline_or_block {
            $$ = CNode("for", $2, CNode(), $3).setPos0($1)
        } |
        TFor ';' expr ';' oneline_or_block oneline_or_block {
            $$ = CNode("for", $3, $5, $6).setPos0($1)
        } |
        TFor expr ';' expr ';' oneline_or_block oneline_or_block {
            $$ = CNode("chain",
                $2,
                CNode("for", $4, $6, $7).setPos0($1),
            )
        } |
        TFor TIdent '=' expr ',' expr oneline_or_block {
            vname, ename := ANode($2), ANodeS($2.Str + randomName())
            $$ = CNode("chain",
                CNode("set", vname, $4).setPos0($1),
                CNode("set", ename, $6).setPos0($1),
                CNode("for", 
                    CNode("<", vname, ename).setPos0($1), 
                    CNode("chain", 
                        CNode("inc", vname, NNode(1.0)).setPos0($1),
                    ), 
                    $7,
                ).setPos0($1),
            )
        } |
        TFor TIdent '=' expr ',' expr ',' expr oneline_or_block {
            vname, sname, ename := ANode($2), ANodeS($2.Str + randomName()), ANodeS($2.Str + randomName()) 
            if $6.Type == Nnumber {
                // easy case
                chain := CNode("chain", CNode("inc", vname, $6).setPos0($1))
                var cond *Node
                if $6.N() < 0 {
                    cond = CNode("<=", ename, vname)
                } else {
                    cond = CNode("<=", vname, ename)
                }
                $$ = CNode("chain",
                    CNode("set", vname, $4).setPos0($1),
                    CNode("set", ename, $8).setPos0($1),
                    CNode("for", cond, chain, $9).setPos0($1),
                )
            } else {
                bname := ANodeS($2.Str + randomName())
                $$ = CNode("chain", 
                    CNode("set", vname, $4).setPos0($1),
                    CNode("set", bname, $4).setPos0($1),
                    CNode("set", sname, $6).setPos0($1),
                    CNode("set", ename, $8).setPos0($1),
                    CNode("if", CNode("<=", NNode(0.0), CNode("*", CNode("-", ename, vname).setPos0($1), sname).setPos0($1)),
                        CNode("chain",
                            CNode("for",
                                CNode("<=",
                                    CNode("*",
                                        CNode("-", vname, bname).setPos0($1), 
                                        CNode("-", vname, ename).setPos0($1),
                                    ),
                                    NNode(0.0),
                                ),
                                CNode("chain", 
                                    CNode("move", vname, CNode("+", vname, sname).setPos0($1),
                                ).setPos0($1)),
                                $9,
                            ).setPos0($1),
                        ),
                        CNode("chain"),
                    ).setPos0($1),
                )
            }
            
        } |
        TFor expr ',' expr {
            $$ = CNode("foreach", $2, $4).setPos0($1)
        } 

if_stat:
        TIf expr oneline_or_block %prec 'T'              { $$ = CNode("if", $2, $3, CNode()).setPos0($1) } |
        TIf expr oneline_or_block TElse oneline_or_block { $$ = CNode("if", $2, $3, $5).setPos0($1) }

switch_stat:
        TSwitch expr '{' switch_body '}'         { $$ = expandSwitch($1, $2, $4.C()) }

switch_body:
        TCase expr ':' stats             { $$ = CNode($2, $4).setPos0($1) } |
        TCase TElse ':' stats            { $$ = CNode(ANode($2), $4).setPos0($1) } |
        switch_body TCase expr ':' stats { $$ = $1.Cappend($3, $5) } |
        switch_body TCase TElse ':' stats{ $$ = $1.Cappend(ANode($3), $5) }

_func:
        { $$ = "" } |
        _func TFunc  { $$ = $1 + ",safe" } |
        _func TVar   { $$ = $1 + ",var" }

func:
        TFunc _func { $$ = ANodeS("func," + $2).setPos($1) }

func_stat:
        func TIdent func_params_list oneline_or_block {
            funcname := ANode($2)
            $$ = CNode(
                "chain", 
                CNode("set", funcname, ANodeS("nil")).setPos0($2), 
                CNode("move", funcname, 
                    CNode($1, funcname, $3, $4).setPos0($2),
                ).setPos0($2),
            )
        }

jmp_stat:
        TYield expr           { $$ = CNode("yield", $2).setPos0($1) } |
        TYieldNil             { $$ = CNode("yield", ANodeS("nil")).setPos0($1) } |
        TBreak                { $$ = CNode("break").setPos0($1) } |
        TContinue             { $$ = CNode("continue").setPos0($1) } |
        TAssert expr          { $$ = CNode("assert", $2, ANodeS("nil")).setPos0($1) } |
        TAssert expr TString  { $$ = CNode("assert", $2, SNode($3.Str)).setPos0($1) } |
        TReturn expr          { $$ = CNode("ret", $2).setPos0($1) } |
        TReturnNil            { $$ = CNode("ret", ANodeS("nil")).setPos0($1) } |
        TUse TString          { $$ = yylex.(*Lexer).loadFile(filepath.Join(filepath.Dir($1.Pos.Source), $2.Str), $1) }

declarator:
        TIdent                                { $$ = ANode($1).setPos($1) } |
        TIdent TSquare                        { $$ = CNode("load", ANodeS("nil"), $1.Str).setPos0($1) } |
        prefix_expr '[' expr ']'              { $$ = CNode("load", $1, $3).setPos0($3).setPos($3) } |
        prefix_expr '[' expr ':' expr ']'     { $$ = CNode("slice", $1, $3, $5).setPos0($3).setPos($3) } |
        prefix_expr '[' expr ':' ']'          { $$ = CNode("slice", $1, $3, NNode("-1")).setPos0($3).setPos($3) } |
        prefix_expr '[' ':' expr ']'          { $$ = CNode("slice", $1, NNode("0"), $4).setPos0($4).setPos($4) } |
        prefix_expr '.' TIdent                { $$ = CNode("load", $1, SNode($3.Str)).setPos0($3).setPos($3) }

ident_list:
        TIdent                                { $$ = CNode($1.Str) } | 
        ident_list ',' TIdent                 { $$ = $1.Cappend(ANode($3)) }

expr_list:
        expr                                  { $$ = CNode($1) } |
        expr_list ',' expr                    { $$ = $1.Cappend($3) }

expr_assign_list:
        expr ':' expr                         { $$ = CNode($1, $3) } |
        expr_assign_list ',' expr ':' expr    { $$ = $1.Cappend($3).Cappend($5) }

expr:
        TNumber              { $$ = NNode($1.Str).SetPos($1) } |
        TUse TString         { $$ = yylex.(*Lexer).loadFile(filepath.Join(filepath.Dir($1.Pos.Source), $2.Str), $1) } |
        TTypeof expr         { $$ = CNode("typeof", $2) } |
        TLen expr            { $$ = CNode("len", $2) } |
        function             { $$ = $1 } |
        map_gen              { $$ = $1 } |
        prefix_expr          { $$ = $1 } |
        postfix_incdec       { $$ = $1 } |
        string               { $$ = $1 } |
        expr TOr expr        { $$ = CNode("or", $1,$3).setPos0($1) } |
        expr TAnd expr       { $$ = CNode("and", $1,$3).setPos0($1) } |
        expr '>' expr        { $$ = CNode("<", $3,$1).setPos0($1) } |
        expr '<' expr        { $$ = CNode("<", $1,$3).setPos0($1) } |
        expr TGte expr       { $$ = CNode("<=", $3,$1).setPos0($1) } |
        expr TLte expr       { $$ = CNode("<=", $1,$3).setPos0($1) } |
        expr TEqeq expr      { $$ = CNode("==", $1,$3).setPos0($1) } |
        expr TNeq expr       { $$ = CNode("!=", $1,$3).setPos0($1) } |
        expr '+' expr        { $$ = CNode("+", $1,$3).setPos0($1) } |
        expr '-' expr        { $$ = CNode("-", $1,$3).setPos0($1) } |
        expr '*' expr        { $$ = CNode("*", $1,$3).setPos0($1) } |
        expr '/' expr        { $$ = CNode("/", $1,$3).setPos0($1) } |
        expr '%' expr        { $$ = CNode("%", $1,$3).setPos0($1) } |
        expr '^' expr        { $$ = CNode("^", $1,$3).setPos0($1) } |
        expr TLsh expr       { $$ = CNode("<<", $1,$3).setPos0($1) } |
        expr TRsh expr       { $$ = CNode(">>", $1,$3).setPos0($1) } |
        expr TURsh expr      { $$ = CNode(">>>", $1,$3).setPos0($1) } |
        expr '|' expr        { $$ = CNode("|", $1,$3).setPos0($1) } |
        expr '&' expr        { $$ = CNode("&", $1,$3).setPos0($1) } |
        '-' expr %prec UNARY { $$ = CNode("-", NNode(0.0), $2).setPos0($2) } |
        '~' expr %prec UNARY { $$ = CNode("~", $2).setPos0($2) } |
        TNot expr %prec UNARY { $$ = CNode("!", $2).setPos0($2) } |
        '#' expr %prec UNARY { $$ = CNode("#", $2).setPos0($2) } |
        '&' TIdent %prec UNARY    { $$ = CNode("addressof", ANode($2)).setPos0($2) }

string: 
        TString { $$ = SNode($1.Str).SetPos($1) } 

prefix_expr:
        declarator        { $$ = $1 } |
        '(' func_call ')' { $$ = $2 } |
        func_call         { $$ = $1 } |
        '(' expr ')'      { $$ = $2 }

func_call:
        prefix_expr func_args {
            $$ = CNode("call", $1, $2).setPos0($1)
        }

func_args:
        '(' ')'           { $$ = CNode() } |
        '(' expr_list ')' { $$ = $2 }

function:
        func func_params_list block %prec FUN { $$ = CNode($1, "<a>", $2, $3).setPos0($1) } |
        func ident_list '=' expr %prec FUN    { $$ = CNode($1, "<a>", $2, CNode("chain", CNode("ret", $4).setPos0($1))).setPos0($1) } |
        func '=' expr %prec FUN               { $$ = CNode($1, "<a>", CNode(), CNode("chain", CNode("ret", $3).setPos0($1))).setPos0($1) }

func_params_list:
        '(' ')'            { $$ = CNode() } |
        '(' ident_list ')' { $$ = $2 }

map_gen:
        '{' '}'          { $$ = CNode("map", CNode()).setPos0($1) } |
        '{' _map_gen '}' { $$ = $2.setPos0($1) }

_map_gen:
        expr_assign_list     { $$ = CNode("map", $1).setPos0($1) } |
        expr_assign_list ',' { $$ = CNode("map", $1).setPos0($1) } |
        expr_list            { $$ = CNode("array", $1).setPos0($1) } |
        expr_list ','        { $$ = CNode("array", $1).setPos0($1) }

%%

var _rand = rand.New()

func randomName() string {
    return fmt.Sprintf("%x", _rand.Fetch(16))
}

func expandSwitch(switchTok Token, sub *Node, cases []*Node) *Node {
    subject := ANodeS("switch" + randomName())
    ret := CNode("chain", CNode("set", subject, sub).setPos0(switchTok))

    var lastif, root *Node
    var defaultCase *Node
    
    for i := 0; i < len(cases); i+=2 {
        if cases[i].S() == "else" {
            defaultCase = cases[i + 1]
            continue
        }

        casestat := CNode("if", CNode("==", subject, cases[i]), cases[i + 1]).setPos0(cases[i])
        if lastif != nil {
            lastif.Cappend(CNode("chain", casestat))
        } else {
            root = casestat
        }
        lastif = casestat
    }

    if defaultCase == nil {
        lastif.Cappend(CNode("chain"))
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
