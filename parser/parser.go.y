%{
package parser

import "github.com/coyove/bracket/vm"
%}
%type<stmts> block
%type<stmt>  stat
%type<stmts> elseifs
%type<expr> var
%type<namelist> namelist
%type<exprlist> exprlist
%type<expr> expr
%type<expr> string
%type<expr> prefixexp
%type<expr> functioncall
%type<expr> afunctioncall
%type<exprlist> args
%type<expr> function

%union {
  token  Token

  stmts    *Node
  stmt     *Node

  funcname interface{}
  funcexpr interface{}

  exprlist *Node
  expr     *Node

  namelist *Node
}

/* Reserved words */
%token<token> TAnd TAssert TBreak TContinue TDo TElse TElseIf TEnd TFalse TFor TIf TIn TLambda TNil TNot TOr TReturn TSet TThen TTrue TWhile TXor

/* Literals */
%token<token> TEqeq TNeq TLsh TRsh TLte TGte TIdent TNumber TString '{' '(' '|' '&'

/* Operators */
%left TOr
%left TAnd
%left '>' '<' TGte TLte TEqeq TNeq
%left '+' '-'
%left '*' '/' '%'
%left '|' '&'
%left TLsh TRsh
%right UNARY /* not # -(unary) */
%right '^'

%% 

block: 
        {
            $$ = NewCompoundNode("chain")
            if l, ok := yylex.(*Lexer); ok {
                l.Stmts = $$
            }
        } |
        block stat {
            $1.Compound = append($1.Compound, $2)
            $$ = $1
            if l, ok := yylex.(*Lexer); ok {
                l.Stmts = $$
            }
        } | 
        block ';' {
            $$ = $1
            if l, ok := yylex.(*Lexer); ok {
                l.Stmts = $$
            }
        }

stat:
        var '=' expr {
            if len($1.Compound) > 0 && $1.Compound[0].Value.(string) == "load" {
                $$ = NewCompoundNode("store", $1.Compound[1], $1.Compound[2], $3)
            } else if len($1.Compound) > 0 && $1.Compound[0].Value.(string) == "safeload" {
                $$ = NewCompoundNode("safestore", $1.Compound[1], $1.Compound[2], $3)
            } else {
                $$ = NewCompoundNode("move", $1, $3)
            }
        } |
        /* 'stat = functioncal' causes a reduce/reduce conflict */
        prefixexp {
            // if _, ok := $1.(*FuncCallExpr); !ok {
            //    yylex.(*Lexer).Error("parse error")
            // } else {
            $$ = $1
            // }
        } |
        TWhile expr TDo block TEnd {
            $$ = NewCompoundNode("while", $2, $4)
        } |
        TIf expr TThen block elseifs TEnd {
            $$ = NewCompoundNode("if", $2, $4, NewCompoundNode())
            cur := $$
            for _, e := range $5.Compound {
                cur.Compound[3] = e
                cur = e
            }
        } |
        TIf expr TThen block elseifs TElse block TEnd {
            $$ = NewCompoundNode("if", $2, $4, NewCompoundNode())
            cur := $$
            for _, e := range $5.Compound {
                cur.Compound[3] = e
                cur = e
            }
            cur.Compound[3] = $7
        } |
        TSet namelist '=' exprlist {
            $$ = NewCompoundNode("chain")
            for i, name := range $2.Compound {
                var e *Node
                if i < len($4.Compound) {
                    e = $4.Compound[i]
                } else {
                    e = $4.Compound[len($4.Compound) - 1]
                }
                $$.Compound = append($$.Compound, NewCompoundNode("set", name, e))
            }
        } |
        TReturn {
            $$ = NewCompoundNode("ret")
        } |
        TReturn expr {
            $$ = NewCompoundNode("ret", $2)
        } |
        TBreak  {
            $$ = NewCompoundNode("break")
        } |
        TContinue  {
            $$ = NewCompoundNode("continue")
        } |
        TAssert expr {
            $$ = NewCompoundNode("assert", $2)
            $$.Compound[0].Pos = $2.Pos
        }

elseifs: 
        {
            $$ = NewCompoundNode()
        } | 
        elseifs TElseIf expr TThen block {
            $$.Compound = append($$.Compound, NewCompoundNode("if", $3, $5, NewCompoundNode()))
        }

var:
        TIdent {
            $$ = NewAtomNode($1)
            _, $$.LibWH = vm.LibLookup[$1.Str]
        } |
        prefixexp '[' expr ']' {
            $$ = NewCompoundNode("load", $1, $3)
        } | 
        prefixexp '{' expr '}' {
            $$ = NewCompoundNode("safeload", $1, $3)
        } | 
        prefixexp '.' TIdent {
            $$ = NewCompoundNode("load", $1, NewStringNode($3.Str))
        }

namelist:
        TIdent {
            $$ = NewCompoundNode($1.Str)
        } | 
        namelist ','  TIdent {
            $1.Compound = append($1.Compound, NewAtomNode($3))
            $$ = $1
        }

exprlist:
        expr {
            $$ = NewCompoundNode($1)
        } |
        exprlist ',' expr {
            $1.Compound = append($1.Compound, $3)
            $$ = $1
        }

expr:
        TNil {
            $$ = NewCompoundNode("nil")
            $$.Compound[0].Pos = $1.Pos
        } | 
        TFalse {
            $$ = NewCompoundNode("false")
            $$.Compound[0].Pos = $1.Pos
        } | 
        TTrue {
            $$ = NewCompoundNode("true")
            $$.Compound[0].Pos = $1.Pos
        } | 
        TNumber {
            $$ = NewNumberNode($1.Str)
            $$.Pos = $1.Pos
        } |
        function {
            $$ = $1
        } | 
        prefixexp {
            $$ = $1
        } |
        string {
            $$ = $1
        } |
        expr TOr expr {
            $$ = NewCompoundNode("or", $1,$3)
            $$.Compound[0].Pos = $1.Pos
        } |
        expr TAnd expr {
            $$ = NewCompoundNode("and", $1,$3)
            $$.Compound[0].Pos = $1.Pos
        } |
        expr TXor expr {
            $$ = NewCompoundNode("xor", $1,$3)
            $$.Compound[0].Pos = $1.Pos
        } |
        expr '>' expr {
            $$ = NewCompoundNode(">", $1,$3)
            $$.Compound[0].Pos = $1.Pos
        } |
        expr '<' expr {
            $$ = NewCompoundNode("<", $1,$3)
            $$.Compound[0].Pos = $1.Pos
        } |
        expr TGte expr {
            $$ = NewCompoundNode(">=", $1,$3)
            $$.Compound[0].Pos = $1.Pos
        } |
        expr TLte expr {
            $$ = NewCompoundNode("<=", $1,$3)
            $$.Compound[0].Pos = $1.Pos
        } |
        expr TEqeq expr {
            $$ = NewCompoundNode("eq", $1,$3)
            $$.Compound[0].Pos = $1.Pos
        } |
        expr TNeq expr {
            $$ = NewCompoundNode("neq", $1,$3)
            $$.Compound[0].Pos = $1.Pos
        } |
        expr '+' expr {
            $$ = NewCompoundNode("+", $1,$3)
            $$.Compound[0].Pos = $1.Pos
        } |
        expr '-' expr {
            $$ = NewCompoundNode("-", $1,$3)
            $$.Compound[0].Pos = $1.Pos
        } |
        expr '*' expr {
            $$ = NewCompoundNode("*", $1,$3)
            $$.Compound[0].Pos = $1.Pos
        } |
        expr '/' expr {
            $$ = NewCompoundNode("/", $1,$3)
            $$.Compound[0].Pos = $1.Pos
        } |
        expr '%' expr {
            $$ = NewCompoundNode("%", $1,$3)
            $$.Compound[0].Pos = $1.Pos
        } |
        expr '^' expr {
            $$ = NewCompoundNode("^", $1,$3)
            $$.Compound[0].Pos = $1.Pos
        } |
        expr TLsh expr {
            $$ = NewCompoundNode("<<", $1,$3)
            $$.Compound[0].Pos = $1.Pos
        } |
        expr TRsh expr {
            $$ = NewCompoundNode(">>", $1,$3)
            $$.Compound[0].Pos = $1.Pos
        } |
        expr '|' expr {
            $$ = NewCompoundNode("|", $1,$3)
            $$.Compound[0].Pos = $1.Pos
        } |
        expr '&' expr {
            $$ = NewCompoundNode("&", $1,$3)
            $$.Compound[0].Pos = $1.Pos
        } |
        '-' expr %prec UNARY {
            $$ = NewCompoundNode("-", NewNumberNode("0"), $2)
            $$.Compound[0].Pos = $2.Pos
        } |
        '~' expr %prec UNARY {
            $$ = NewCompoundNode("~", $2)
            $$.Compound[0].Pos = $2.Pos
        } |
        TNot expr %prec UNARY {
            $$ = NewCompoundNode("not", $2)
            $$.Compound[0].Pos = $2.Pos
        }

string: 
        TString {
            $$ = NewStringNode($1.Str)
            $$.Pos = $1.Pos
        } 

prefixexp:
        var {
            $$ = $1
        } |
        afunctioncall {
            $$ = $1
        } |
        functioncall {
            $$ = $1
        } |
        '(' expr ')' {
            $$ = $2
        }

afunctioncall:
        '(' functioncall ')' {
            $$ = $2
        }

functioncall:
        prefixexp args {
            $$ = NewCompoundNode("call", $1, $2)
        }

args:
        '(' ')' {
            if yylex.(*Lexer).PNewLine {
               yylex.(*Lexer).TokenError($1, "ambiguous syntax (function call x new statement)")
            }
            $$ = NewCompoundNode()
        } |
        '(' exprlist ')' {
            if yylex.(*Lexer).PNewLine {
               yylex.(*Lexer).TokenError($1, "ambiguous syntax (function call x new statement)")
            }
            $$ = $2
        }

function:
        TLambda '(' ')' block TEnd {
            $$ = NewCompoundNode("lambda", NewCompoundNode(), $4)
            $$.Compound[0].Pos = $1.Pos
        } |
        TLambda '(' namelist ')' block TEnd {
            $$ = NewCompoundNode("lambda", $3, $5)
            $$.Compound[0].Pos = $1.Pos
        }

%%

func TokenName(c int) string {
	if c >= TAnd && c-TAnd < len(yyToknames) {
		if yyToknames[c-TAnd] != "" {
			return yyToknames[c-TAnd]
		}
	}
    return string([]byte{byte(c)})
}

