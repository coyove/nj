%{
package parser

import (
    "bytes"
    "io/ioutil"
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
%type<expr> expr_declare_list
%type<expr> expr
%type<expr> string

%type<expr> prefix_expr
%type<expr> assign_stat
%type<expr> assign_stat_semi
%type<expr> for_stat
%type<expr> for_stat1
%type<expr> for_stat2
%type<expr> if_stat
%type<expr> if_body
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
%token<token> TAssert TBreak TContinue TElse TFor TFunc TGoto TIf TNil TReturn TRequire TVar TYield

/* Literals */
%token<token> TEqeq TNeq TLsh TRsh TLte TGte TIdent TNumber TString '{' '('

/* Operators */
%right 'T'
%right TElse

%left TOr
%left TAnd
%left '|' '&' '^'
%left '>' '<' TGte TLte TEqeq TNeq
%left TLsh TRsh
%left '+' '-'
%left '*' '/' '%'
%right UNARY /* not # -(unary) */
%right '~'
%right '#'

%% 

stats: 
        {
            $$ = NewCompoundNode("chain")
            if l, ok := yylex.(*Lexer); ok {
                l.Stmts = $$
            }
        } |
        stats stat {
            $1.Compound = append($1.Compound, $2)
            $$ = $1
            if l, ok := yylex.(*Lexer); ok {
                l.Stmts = $$
            }
        }

block: 
        '{' stats '}' { $$ = $2 }

assign_stat_semi:
        assign_stat ';' {
            if $1.isIsolatedDupCall() {
                $1.Compound[2].Compound[0] = NewNumberNode("0")
            }
            $$ = $1
        }

stat:
        ';'              { $$ = NewCompoundNode() } |
        jmp_stat ';'     { $$ = $1 } |
        flow_stat        { $$ = $1 } |
        assign_stat_semi { $$ = $1 } |
        TIdent ':'       { $$ = NewCompoundNode("label", NewAtomNode($1)) }

if_body:
        ';'              { $$ = NewCompoundNode() } |
        assign_stat_semi { $$ = NewCompoundNode("chain", $1) } |
        jmp_stat ';'     { $$ = NewCompoundNode("chain", $1) } |
        for_stat         { $$ = NewCompoundNode("chain", $1) } |
        if_stat          { $$ = NewCompoundNode("chain", $1) } |
        block            { $$ = $1 }

flow_stat:
        for_stat  { $$ = $1 } |
        if_stat   { $$ = $1 } |
        func_stat { $$ = $1 }

assign_stat:
        TVar expr_declare_list {
            $$ = $2
        } |
        declarator '=' expr {
            $$ = NewCompoundNode("move", $1, $3)
            if len($1.Compound) > 0 {
                if c, _ := $1.Compound[0].Value.(string); c == "load" {
                    $$ = NewCompoundNode("store", $1.Compound[1], $1.Compound[2], $3)
                }
            }
            if c, _ := $1.Value.(string); c != "" && $1.Type == NTAtom {
                if a, b, s := $3.isSimpleAddSub(); a == c {
                    $3.Compound[2].Value = $3.Compound[2].Value.(float64) * s
                    $$ = NewCompoundNode("inc", $1, $3.Compound[2])
                    $$.Compound[1].Pos = $1.Pos
                } else if b == c {
                    $3.Compound[1].Value = $3.Compound[1].Value.(float64) * s
                    $$ = NewCompoundNode("inc", $1, $3.Compound[1])
                    $$.Compound[1].Pos = $1.Pos
                }
            }
            $$.Compound[0].Pos = $1.Pos
        } |
        prefix_expr {
            $$ = $1
        }

for_stat1:
        assign_stat ';' { $$ = NewCompoundNode("chain", $1) } |
        ';'             { $$ = NewCompoundNode("chain") }

for_stat2:
        expr ';' { $$ = $1 } |
        ';'      { $$ = NewNumberNode("1") }

for_stat:
        TFor block {
            $$ = NewCompoundNode("for", NewNumberNode("1"), NewCompoundNode(), $2).setPos0($1.Pos)
        } |
        TFor '(' expr ')' if_body {
            $$ = NewCompoundNode("for", $3, NewCompoundNode(), $5).setPos0($1.Pos)
        } |
        TFor '(' for_stat1 for_stat2 assign_stat ')' if_body {
            $$ = $3
            $$.Compound = append($$.Compound, NewCompoundNode("for", $4, NewCompoundNode("chain", $5), $7))
            $$.Compound[0].Pos = $1.Pos
        } |
        TFor '(' for_stat1 for_stat2 block ')' if_body {
            $$ = $3
            $$.Compound = append($$.Compound, NewCompoundNode("for", $4, $5, $7))
            $$.Compound[0].Pos = $1.Pos
        } |
        TFor '(' for_stat1 for_stat2 ')' if_body {
            $$ = $3
            $$.Compound = append($$.Compound, NewCompoundNode("for", $4, NewCompoundNode(), $6))
            $$.Compound[0].Pos = $1.Pos
        }

if_stat:
        TIf '(' expr ')' if_body %prec 'T'     { $$ = NewCompoundNode("if", $3, $5, NewCompoundNode()) } |
        TIf '(' expr ')' if_body TElse if_body { $$ = NewCompoundNode("if", $3, $5, $7) }

func:
        TFunc     { $$ = "func" } |
        TFunc '!' { $$ = "safefunc" }

func_stat:
        func TIdent func_params_list block {
            funcname := NewAtomNode($2)
            $$ = NewCompoundNode(
                "chain", 
                NewCompoundNode("set", funcname, NewNilNode()), 
                NewCompoundNode("move", funcname, NewCompoundNode($1, $3, $4)))
            $$.Compound[1].Compound[0].Pos = $2.Pos
            $$.Compound[2].Compound[0].Pos = $2.Pos
            $$.Compound[2].Compound[2].Compound[0].Pos = $2.Pos
        }

jmp_stat:
        TGoto TIdent { $$ = NewCompoundNode("goto", NewAtomNode($2)).setPos0($1.Pos) } |
        TYield       { $$ = NewCompoundNode("yield").setPos0($1.Pos) } |
        TYield expr  { $$ = NewCompoundNode("yield", $2).setPos0($1.Pos) } |
        TBreak       { $$ = NewCompoundNode("break").setPos0($1.Pos) } |
        TContinue    { $$ = NewCompoundNode("continue").setPos0($1.Pos) } |
        TAssert expr { $$ = NewCompoundNode("assert", $2).setPos0($1.Pos) } |
        TReturn      { $$ = NewCompoundNode("ret").setPos0($1.Pos) } |
        TReturn expr {
            if $2.isIsolatedDupCall() {
                if h, _ := $2.Compound[2].Compound[2].Value.(float64); h == 1 {
                    $2.Compound[2].Compound[2] = NewNumberNode("2")
                }
            }
            $$ = NewCompoundNode("ret", $2).setPos0($1.Pos)
        }

declarator:
        TIdent                            { $$ = NewAtomNode($1).setPos($1.Pos) } |
        prefix_expr '[' expr ']'          { $$ = NewCompoundNode("load", $1, $3).setPos0($1.Pos).setPos($1.Pos) } |
        prefix_expr '[' expr ':' expr ']' { $$ = NewCompoundNode("slice", $1, $3, $5).setPos0($1.Pos).setPos($1.Pos) } |
        prefix_expr '[' expr ':' ']'      { $$ = NewCompoundNode("slice", $1, $3, NewNumberNode("-1")).setPos0($1.Pos).setPos($1.Pos) } |
        prefix_expr '[' ':' expr ']'      { $$ = NewCompoundNode("slice", $1, NewNumberNode("0"), $4).setPos0($1.Pos).setPos($1.Pos) } |
        prefix_expr '.' TIdent            { $$ = NewCompoundNode("load", $1, NewStringNode($3.Str)).setPos0($1.Pos).setPos($1.Pos) }

ident_list:
        TIdent {
            $$ = NewCompoundNode($1.Str)
        } | 
        ident_list ',' TIdent {
            $1.Compound = append($1.Compound, NewAtomNode($3))
            $$ = $1
        }

expr_list:
        expr {
            $$ = NewCompoundNode($1)
        } |
        expr_list ',' expr {
            $1.Compound = append($1.Compound, $3)
            $$ = $1
        }

expr_assign_list:
        expr ':' expr {
            $$ = NewCompoundNode($1, $3)
        } |
        expr_assign_list ',' expr ':' expr {
            $1.Compound = append($1.Compound, $3, $5)
            $$ = $1
        }

expr_declare_list:
        TIdent {
            $$ = NewCompoundNode("chain", NewCompoundNode("set", NewAtomNode($1), NewNilNode()))
            $$.Compound[1].Compound[0].Pos = $1.Pos
        } |
        TIdent '=' expr {
            $$ = NewCompoundNode("chain", NewCompoundNode("set", NewAtomNode($1), $3))
            $$.Compound[1].Compound[0].Pos = $1.Pos
        } |
        expr_declare_list ',' TIdent '=' expr {
            x := NewCompoundNode("set", NewAtomNode($3), $5).setPos0($1.Pos)
            $1.Compound = append($$.Compound, x)
            $$ = $1
        } |
        expr_declare_list ',' TIdent {
            x := NewCompoundNode("set", NewAtomNode($3), NewNilNode()).setPos0($1.Pos)
            $1.Compound = append($1.Compound, x)
            $$ = $1
        }

expr:
        TNil {
            $$ = NewNilNode()
            $$.Pos = $1.Pos
        } |
        TNumber {
            $$ = NewNumberNode($1.Str)
            $$.Pos = $1.Pos
        } |
        TRequire TString {
            path := filepath.Dir($1.Pos.Source)
            path = filepath.Join(path, $2.Str)

            code, err := ioutil.ReadFile(path)
            if err != nil {
                yylex.(*Lexer).Error(err.Error())
            }
            n, err := Parse(bytes.NewReader(code), path)
            if err != nil {
                yylex.(*Lexer).Error(err.Error())
            }

            // now the required code is loaded, for naming scope we will wrap them into a closure
            cls := NewCompoundNode("func", NewCompoundNode(), n)
            $$ = NewCompoundNode("call", cls, NewCompoundNode())
        } |
        function             { $$ = $1 } |
        map_gen              { $$ = $1 } |
        prefix_expr          { $$ = $1 } |
        string               { $$ = $1 } |
        expr TOr expr        { $$ = NewCompoundNode("or", $1,$3).setPos0($1.Pos) } |
        expr TAnd expr       { $$ = NewCompoundNode("and", $1,$3).setPos0($1.Pos) } |
        expr '>' expr        { $$ = NewCompoundNode("<", $3,$1).setPos0($1.Pos) } |
        expr '<' expr        { $$ = NewCompoundNode("<", $1,$3).setPos0($1.Pos) } |
        expr TGte expr       { $$ = NewCompoundNode("<=", $3,$1).setPos0($1.Pos) } |
        expr TLte expr       { $$ = NewCompoundNode("<=", $1,$3).setPos0($1.Pos) } |
        expr TEqeq expr      { $$ = NewCompoundNode("==", $1,$3).setPos0($1.Pos) } |
        expr TNeq expr       { $$ = NewCompoundNode("!=", $1,$3).setPos0($1.Pos) } |
        expr '+' expr        { $$ = NewCompoundNode("+", $1,$3).setPos0($1.Pos) } |
        expr '-' expr        { $$ = NewCompoundNode("-", $1,$3).setPos0($1.Pos) } |
        expr '*' expr        { $$ = NewCompoundNode("*", $1,$3).setPos0($1.Pos) } |
        expr '/' expr        { $$ = NewCompoundNode("/", $1,$3).setPos0($1.Pos) } |
        expr '%' expr        { $$ = NewCompoundNode("%", $1,$3).setPos0($1.Pos) } |
        expr '^' expr        { $$ = NewCompoundNode("^", $1,$3).setPos0($1.Pos) } |
        expr TLsh expr       { $$ = NewCompoundNode("<<", $1,$3).setPos0($1.Pos) } |
        expr TRsh expr       { $$ = NewCompoundNode(">>", $1,$3).setPos0($1.Pos) } |
        expr '|' expr        { $$ = NewCompoundNode("|", $1,$3).setPos0($1.Pos) } |
        expr '&' expr        { $$ = NewCompoundNode("&", $1,$3).setPos0($1.Pos) } |
        '-' expr %prec UNARY { $$ = NewCompoundNode("-", NewNumberNode("0"), $2).setPos0($2.Pos) } |
        '~' expr %prec UNARY { $$ = NewCompoundNode("~", $2).setPos0($2.Pos) } |
        '!' expr %prec UNARY { $$ = NewCompoundNode("!", $2).setPos0($2.Pos) } |
        '#' expr %prec UNARY { $$ = NewCompoundNode("#", $2).setPos0($2.Pos) }

string: 
        TString {
            $$ = NewStringNode($1.Str)
            $$.Pos = $1.Pos
        } 

prefix_expr:
        declarator        { $$ = $1 } |
        '(' func_call ')' { $$ = $2 } |
        func_call         { $$ = $1 } |
        '(' expr ')'      { $$ = $2 }

func_call:
        prefix_expr func_args {
            switch c, _ := $1.Value.(string); c {
            case "copy":
                switch len($2.Compound) {
                case 0:
                    $$ = NewCompoundNode("call", $1, NewCompoundNode(NewNumberNode("1"), NewNumberNode("1"), NewNumberNode("1")))
                case 1:
                    $$ = NewCompoundNode("call", $1, NewCompoundNode(NewNumberNode("1"), $2.Compound[0], NewNumberNode("0")))
                default:
                    p := $2.Compound[1]
                    if p.Type != NTCompound && p.Type != NTAtom {
                        yylex.(*Lexer).Error("invalid argument for S")
                    }
                    $$ = NewCompoundNode("call", $1, NewCompoundNode(NewNumberNode("1"), $2.Compound[0], p))
                }
            case "typeof":
                switch len($2.Compound) {
                case 0:
                    yylex.(*Lexer).Error("typeof takes at least 1 argument")
                case 1:
                    $$ = NewCompoundNode("call", $1, NewCompoundNode($2.Compound[0], NewNumberNode("255")))
                default:
                    x, _ := $2.Compound[1].Value.(string);
                    if ti, ok := typesLookup[x]; ok {
                        $$ = NewCompoundNode("call", $1, NewCompoundNode($2.Compound[0], NewNumberNode(ti)))
                    } else {
                        yylex.(*Lexer).Error("invalid typename in typeof")
                    }
                }
            case "addressof":
                if len($2.Compound) != 1 {
                    yylex.(*Lexer).Error("addressof takes 1 argument")
                }
                if $2.Compound[0].Type != NTAtom {
                    yylex.(*Lexer).Error("addressof can only get the address of a variable")
                }
                $$ = NewCompoundNode("call", $1, $2)
            case "len":
                switch len($2.Compound) {
                case 0:
                    yylex.(*Lexer).Error("len takes 1 argument")
                default:
                    $$ = NewCompoundNode("call", $1, $2)
                }
            default:
                $$ = NewCompoundNode("call", $1, $2)
            }
            $$.Compound[0].Pos = $1.Pos
        }

func_args:
        '(' ')'           { $$ = NewCompoundNode() } |
        '(' expr_list ')' { $$ = $2 }

function:
        func func_params_list block { $$ = NewCompoundNode($1, $2, $3).setPos0($2.Pos) }

func_params_list:
        '(' ')'            { $$ = NewCompoundNode() } |
        '(' ident_list ')' { $$ = $2 }

map_gen:
        '{' '}'          { $$ = NewCompoundNode("map", NewCompoundNode()).setPos0($1.Pos) } |
        '{' _map_gen '}' { $$ = $2.setPos0($1.Pos) }

_map_gen:
        expr_assign_list     { $$ = NewCompoundNode("map", $1).setPos0($1.Pos) } |
        expr_assign_list ',' { $$ = NewCompoundNode("map", $1).setPos0($1.Pos) } |
        expr_list {
            table := NewCompoundNode()
            for i, v := range $1.Compound {
                table.Compound = append(table.Compound, &Node{ Type:  NTNumber, Value: float64(i) }, v)
            }
            $$ = NewCompoundNode("map", table).setPos0($1.Pos)
        } |
        expr_list ',' {
            table := NewCompoundNode()
            for i, v := range $1.Compound {
                table.Compound = append(table.Compound, &Node{ Type:  NTNumber, Value: float64(i) }, v)
            }
            $$ = NewCompoundNode("map", table).setPos0($1.Pos)
        }

%%

var typesLookup = map[string]string {
    "nil": "0", "number": "1", "string": "2", "map": "3", "closure": "4", "generic": "5",
}