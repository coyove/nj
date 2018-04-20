%{
package parser
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
%token<token> TAnd TAssert TBreak TContinue TDo TElse TElseIf TEnd TFalse TFor TFunction TIf TIn TNil TNot TOr TReturn TSet TThen TTrue TTypeIs TWhile TXor

/* Literals */
%token<token> TEqeq TNeq TLte TGte TIdent TNumber TString '{' '('

/* Operators */
%left TOr
%left TAnd
%left '>' '<' TGte TLte TEqeq TNeq
%left '+' '-'
%left '*' '/' '%'
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
        TSet TIdent '=' expr {
            $$ = NewCompoundNode("set", $2.Str, $4)
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
            $$.Pos = $2.Pos
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
        } |
        prefixexp '[' expr ']' {
            $$ = NewCompoundNode("load", $1, $3)
        } | 
        prefixexp '.' TIdent {
            $$ = NewCompoundNode("load", $1, $3.Str)
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
            $$.Pos = $1.Pos
        } |
        exprlist ',' expr {
            $1.Compound = append($1.Compound, $3)
            $$ = $1
        }

expr:
        TNil {
            $$ = NewCompoundNode("nil")
            $$.Pos = $1.Pos
        } | 
        TFalse {
            $$ = NewCompoundNode("false")
            $$.Pos = $1.Pos
        } | 
        TTrue {
            $$ = NewCompoundNode("true")
            $$.Pos = $1.Pos
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
        expr TTypeIs TIdent {
            $$ = NewCompoundNode("typeof", $1, $3)
            $$.Pos = $1.Pos
        } |
        expr TOr expr {
            $$ = NewCompoundNode("or", $1,$3)
            $$.Pos = $1.Pos
        } |
        expr TAnd expr {
            $$ = NewCompoundNode("and", $1,$3)
            $$.Pos = $1.Pos
        } |
        expr TXor expr {
            $$ = NewCompoundNode("xor", $1,$3)
            $$.Pos = $1.Pos
        } |
        expr '>' expr {
            $$ = NewCompoundNode(">", $1,$3)
            $$.Pos = $1.Pos
        } |
        expr '<' expr {
            $$ = NewCompoundNode("<", $1,$3)
            $$.Pos = $1.Pos
        } |
        expr TGte expr {
            $$ = NewCompoundNode(">=", $1,$3)
            $$.Pos = $1.Pos
        } |
        expr TLte expr {
            $$ = NewCompoundNode("<=", $1,$3)
            $$.Pos = $1.Pos
        } |
        expr TEqeq expr {
            $$ = NewCompoundNode("eq", $1,$3)
            $$.Pos = $1.Pos
        } |
        expr TNeq expr {
            $$ = NewCompoundNode("neq", $1,$3)
            $$.Pos = $1.Pos
        } |
        expr '+' expr {
            $$ = NewCompoundNode("+", $1,$3)
            $$.Pos = $1.Pos
        } |
        expr '-' expr {
            $$ = NewCompoundNode("-", $1,$3)
            $$.Pos = $1.Pos
        } |
        expr '*' expr {
            $$ = NewCompoundNode("*", $1,$3)
            $$.Pos = $1.Pos
        } |
        expr '/' expr {
            $$ = NewCompoundNode("/", $1,$3)
            $$.Pos = $1.Pos
        } |
        expr '%' expr {
            $$ = NewCompoundNode("%", $1,$3)
            $$.Pos = $1.Pos
        } |
        expr '^' expr {
            $$ = NewCompoundNode("^", $1,$3)
            $$.Pos = $1.Pos
        } |
        '-' expr %prec UNARY {
            $$ = NewCompoundNode("-", $2)
            $$.Pos = $2.Pos
        } |
        '~' expr %prec UNARY {
            $$ = NewCompoundNode("~", $2)
            $$.Pos = $2.Pos
        } |
        TNot expr %prec UNARY {
            $$ = NewCompoundNode("not", $2)
            $$.Pos = $2.Pos
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
        TFunction '(' ')' block TEnd {
            $$ = NewCompoundNode("lambda", NewCompoundNode(), $4)
            $$.Pos = $1.Pos
        } |
        TFunction '(' namelist ')' block TEnd {
            $$ = NewCompoundNode("lambda", $3, $5)
            $$.Pos = $1.Pos
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

