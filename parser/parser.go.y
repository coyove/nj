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
%type<expr> expr_declare_list
%type<expr> expr
%type<expr> postfix_incdec
%type<expr> _postfix_incdec
%type<expr> string

%type<expr> prefix_expr
%type<expr> _assign_stat
%type<expr> assign_stat
%type<expr> for_stat
%type<expr> if_stat
%type<expr> oneline_or_block
%type<expr> jmp_stat
%type<expr> func_stat
%type<expr> flow_stat

%type<str>  func
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
%token<token> TAssert TBreak TContinue TElse TFor TFunc TIf TNil TNot TReturn TRequire TNop TVar TWhile TYield

/* Literals */
%token<token> TAddAdd TSubSub TEqeq TNeq TLsh TRsh TURsh TLte TGte TIdent TNumber TString '{' '[' '('

/* Operators */
%right 'T'
%right TElse

%right FUN
%left TOr
%left TAnd
%left '|' '&' '^'
%left '>' '<' TGte TLte TEqeq TNeq
%left TLsh TRsh TURsh
%left '+' '-'
%left '*' '/' '%'
%right UNARY /* not # -(unary) */
%right '~'
%right '#'
%left TAddAdd TMinMin

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
        jmp_stat     { $$ = $1 } |
        flow_stat        { $$ = $1 } |
        assign_stat { $$ = $1 } |
        block {$$ = $1 }

oneline_or_block:
        assign_stat            { $$ = CNode("chain", $1) } |
        jmp_stat               { $$ = CNode("chain", $1) } |
        for_stat               { $$ = CNode("chain", $1) } |
        if_stat                { $$ = CNode("chain", $1) } |
        block                  { $$ = $1 }

flow_stat:
        for_stat               { $$ = $1 } |
        if_stat                { $$ = $1 } |
        func_stat              { $$ = $1 }

_assign_stat:
        TVar expr_declare_list { $$ = $2 } |
        prefix_expr            { $$ = $1 } |
        postfix_incdec         { $$ = $1 } |
        declarator '=' expr {
            $$ = CNode("move", $1, $3)
            if $1.Cn() > 0 && $1.Cx(0).S() == "load" {
                $$ = CNode("store", $1.Cx(1), $1.Cx(2), $3)
            }
            if c := $1.S(); c != "" && $1.Type == Natom {
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

postfix_incdec:
        TIdent _postfix_incdec                    { $$ = CNode("inc", ANode($1).setPos($1), $2) } |
        prefix_expr '[' expr ']' _postfix_incdec  { $$ = CNode("store", $1, $3, CNode("+", CNode("load", $1, $3).setPos0($1), $5).setPos0($1)) } |
        prefix_expr '.' TIdent   _postfix_incdec  { $$ = CNode("store", $1, $3, CNode("+", CNode("load", $1, $3).setPos0($1), $4).setPos0($1)) }

for_stat:
        TWhile expr oneline_or_block {
            $$ = CNode("for", $2, CNode(), $3).setPos0($1)
        } |
        TWhile expr TContinue '=' oneline_or_block oneline_or_block {
            $$ = CNode("for", $2, $5, $6).setPos0($1)
        } |
        TFor TIdent '=' expr ',' expr oneline_or_block {
            vname, ename := ANode($2), ANodeS($2.Str + randomName())
            $$ = CNode("chain",
                CNode("set", vname, $4),
                CNode("set", ename, $6),
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
                chain, cmp := CNode("chain", CNode("inc", vname, $6).setPos0($1)), "<="
                if $6.N() < 0 {
                    cmp = ">="
                }
                $$ = CNode("chain",
                    CNode("set", vname, $4),
                    CNode("set", ename, $8),
                    CNode("for", CNode(cmp, vname, ename), chain, $9).setPos0($1),
                )
            } else {
                bname := ANodeS($2.Str + randomName())
                $$ = CNode("chain", 
                    CNode("set", vname, $4),
                    CNode("set", bname, $4),
                    CNode("set", sname, $6),
                    CNode("set", ename, $8),
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
                                CNode("chain", CNode("move", vname, CNode("+", vname, sname).setPos0($1))),
                                $9,
                            ),
                        ),
                        CNode("chain"),
                    ),
                )
            }
            
        } |
        TFor TIdent ',' TIdent '=' expr oneline_or_block {
            $$ = CNode("call", "copy", CNode(
               NNode(0.0),
               $6,
               CNode("func", "<anony-map-iter-callback>", CNode($2.Str, $4.Str), $7),
            ))
        }

if_stat:
        TIf expr oneline_or_block %prec 'T'     { $$ = CNode("if", $2, $3, CNode()) } |
        TIf expr oneline_or_block TElse oneline_or_block { $$ = CNode("if", $2, $3, $5) }

func:
        TFunc     { $$ = "func" } |
        TFunc '!' { $$ = "safefunc" }

func_stat:
        func TIdent func_params_list oneline_or_block {
            funcname := ANode($2)
            $$ = CNode(
                "chain", 
                CNode("set", funcname, NilNode()).setPos0($2), 
                CNode("move", funcname, 
                    CNode($1, funcname, $3, $4).setPos0($2),
                ).setPos0($2),
            )
        }

jmp_stat:
         TYield       '.'  { $$ = CNode("yield").setPos0($1) } |
         TYield expr    { $$ = CNode("yield", $2).setPos0($1) } |
         TBreak         { $$ = CNode("break").setPos0($1) } |
        TContinue      { $$ = CNode("continue").setPos0($1) } |
         TAssert expr   { $$ = CNode("assert", $2).setPos0($1) } |
         TReturn      '.'  { $$ = CNode("ret").setPos0($1) } |
         TReturn expr   {
            if $2.isIsolatedCopy() && $2.Cx(2).Cx(2).N() == 1 {
                $2.Cx(2).C()[2] = NNode(2.0)
            }
            $$ = CNode("ret", $2).setPos0($1)
        } |
        TRequire TString {
            path := filepath.Join(filepath.Dir($1.Pos.Source), $2.Str)
            $$ = yylex.(*Lexer).loadFile(path)
        }

declarator:
        TIdent                            { $$ = ANode($1).setPos($1) } |
        prefix_expr '[' expr ']'          { $$ = CNode("load", $1, $3).setPos0($1).setPos($1) } |
        prefix_expr '[' expr ':' expr ']' { $$ = CNode("slice", $1, $3, $5).setPos0($1).setPos($1) } |
        prefix_expr '[' expr ':' ']'      { $$ = CNode("slice", $1, $3, NNode("-1")).setPos0($1).setPos($1) } |
        prefix_expr '[' ':' expr ']'      { $$ = CNode("slice", $1, NNode("0"), $4).setPos0($1).setPos($1) } |
        prefix_expr '.' TIdent            { $$ = CNode("load", $1, SNode($3.Str)).setPos0($1).setPos($1) }

ident_list:
        TIdent                { $$ = CNode($1.Str) } | 
        ident_list ',' TIdent { $$ = $1.Cappend(ANode($3)) }

expr_list:
        expr               { $$ = CNode($1) } |
        expr_list ',' expr { $$ = $1.Cappend($3) }

expr_assign_list:
        expr ':' expr                      { $$ = CNode($1, $3) } |
        expr_assign_list ',' expr ':' expr { $$ = $1.Cappend($3).Cappend($5) }

expr_declare_list:
        TIdent                                { $$ = CNode("chain", CNode("set", ANode($1), NilNode()).setPos0($1)) } |
        TIdent '=' expr                       { $$ = CNode("chain", CNode("set", ANode($1), $3).setPos0($1)) } |
        expr_declare_list ',' TIdent '=' expr { $$ = $1.Cappend(CNode("set", ANode($3), $5).setPos0($1)) } |
        expr_declare_list ',' TIdent          { $$ = $1.Cappend(CNode("set", ANode($3), NilNode()).setPos0($1)) }

expr:
        TNil                 { $$ = NilNode().SetPos($1) } |
        TNumber              { $$ = NNode($1.Str).SetPos($1) } |
        TRequire TString     { $$ = yylex.(*Lexer).loadFile(filepath.Join(filepath.Dir($1.Pos.Source), $2.Str)) } |
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
        '#' expr %prec UNARY { $$ = CNode("#", $2).setPos0($2) }

string: 
        TString { $$ = SNode($1.Str).SetPos($1) } 

prefix_expr:
        declarator        { $$ = $1 } |
        '(' func_call ')' { $$ = $2 } |
        func_call         { $$ = $1 } |
        '(' expr ')'      { $$ = $2 }

func_call:
        prefix_expr func_args {
            switch $1.S() {
            case "copy":
                switch $2.Cn() {
                case 0:
                    $$ = CNode("call", $1, CNode(NNode("1"), NNode("1"), NNode("1")))
                case 1:
                    $$ = CNode("call", $1, CNode(NNode("1"), $2.Cx(0), NNode("0")))
                default:
                    p := $2.Cx(1)
                    if p.Type != Ncompound && p.Type != Natom {
                        yylex.(*Lexer).Error("invalid argument for copy")
                    }
                    $$ = CNode("call", $1, CNode(NNode("1"), $2.Cx(0), p))
                }
            case "typeof":
                switch $2.Cn() {
                case 0:
                    yylex.(*Lexer).Error("typeof takes at least 1 argument")
                case 1:
                    $$ = CNode("call", $1, CNode($2.Cx(0), NNode("255")))
                default:
                    x, _ := $2.Cx(1).Value.(string);
                    if ti, ok := typesLookup[x]; ok {
                        $$ = CNode("call", $1, CNode($2.Cx(0), NNode(ti)))
                    } else {
                        yylex.(*Lexer).Error("invalid typename in typeof")
                    }
                }
            case "addressof":
                if $2.Cn() != 1 {
                    yylex.(*Lexer).Error("addressof takes 1 argument")
                }
                if $2.Cx(0).Type != Natom {
                    yylex.(*Lexer).Error("addressof can only get the address of a variable")
                }
                $$ = CNode("call", $1, $2)
            case "len":
                switch $2.Cn() {
                case 0:
                    yylex.(*Lexer).Error("len takes 1 argument")
                default:
                    $$ = CNode("call", $1, $2)
                }
            default:
                $$ = CNode("call", $1, $2)
            }
            $$.Cx(0).SetPos($1)
        }

func_args:
        '(' ')'           { $$ = CNode() } |
        '(' expr_list ')' { $$ = $2 }

function:
        func func_params_list block %prec FUN{ $$ = CNode($1, "<a>", $2, $3).setPos0($2) } |
        func ident_list '=' expr %prec FUN  { $$ = CNode($1, "<a>", $2, CNode("chain", CNode("ret", $4))).setPos0($2) } |
        func '=' expr  %prec FUN { $$ = CNode($1, "<a>", CNode(), CNode("chain", CNode("ret", $3))).setPos0($3) }

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

var typesLookup = map[string]string {
    "nil": "0", "number": "1", "string": "2", "map": "4", "closure": "6", "generic": "7",
}

var _rand = rand.New()

func randomName() string {
    return fmt.Sprintf("%x", _rand.Fetch(16))
}