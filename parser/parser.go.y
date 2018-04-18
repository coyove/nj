%{
package parser
%}
%type<stmts> block
%type<stmt>  stat
%type<stmts> elseifs
%type<funcname> funcname
%type<funcname> funcname1
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
%type<funcexpr> funcbody
%type<parlist> parlist

%union {
  token  Token

  stmts    []interface{}
  stmt     interface{}

  funcname interface{}
  funcexpr interface{}

  exprlist []interface{} 
  expr   interface{}

  namelist []interface{}
  parlist  interface{}
}

/* Reserved words */
%token<token> TAnd TBreak TContinue TDo TElse TElseIf TEnd TFalse TFor TFunction TIf TIn TNil TNot TOr TReturn TSet TThen TTrue TWhile

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
            $$ = []interface{}{"chain"}
            if l, ok := yylex.(*Lexer); ok {
                l.Stmts = $$
            }
        } |
        block stat {
            $$ = append($1, $2)
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
            $$ = []interface{}{"set", $1, $3}
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
            $$ = []interface{}{"while", $2, $4}
        } |
        TIf expr TThen block elseifs TEnd {
            $$ = []interface{}{"if", $2, $4, nil}
            if len($5) > 0 {
                cur := $$
                for _, e := range $5 {
                    cur.([]interface{})[3] = e
                    cur = e
                }
            }
        } |
        TIf expr TThen block elseifs TElse block TEnd {
            $$ = []interface{}{"if", $2, $4, nil}
            cur := $$
            if len($5) > 0 {
                for _, e := range $5 {
                    cur.([]interface{})[3] = e
                    cur = e
                }
            }
            cur.([]interface{})[3] = $7
        } |
        TSet TIdent '=' expr {
            $$ = []interface{}{"var", $2.Str, $4}
        } |
        TReturn {
            $$ = []interface{}{"ret"}
        } |
        TReturn exprlist {
            $$ = []interface{}{"ret", $2}
        } |
        TBreak  {
            $$ = "break"
        } |
        TContinue  {
            $$ = "continue"
        }

elseifs: 
        {
            $$ = []interface{}{}
        } | 
        elseifs TElseIf expr TThen block {
            $$ = append($$, []interface{}{"if", $3, $5, nil})
        }

funcname: 
        funcname1 {
            $$ = $1
        }

funcname1:
        TIdent {
            $$ = $1.Str
        }

var:
        TIdent {
            $$ = $1.Str
        } |
        prefixexp '[' expr ']' {
            $$ = []interface{}{$1, ":", $3}
        } | 
        prefixexp '.' TIdent {
            $$ = []interface{}{$1, ":", $3.Str}
        }

namelist:
        TIdent {
            $$ = []interface{}{$1.Str}
        } | 
        namelist ','  TIdent {
            $$ = append($1, $3.Str)
        }

exprlist:
        expr {
            $$ = []interface{}{$1}
        } |
        exprlist ',' expr {
            $$ = append($1, $3)
        }

expr:
        TNil {
            $$ = "nil"
        } | 
        TFalse {
            $$ = "false"
        } | 
        TTrue {
            $$ = "true"
        } | 
        TNumber {
            $$ = $1.Str
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
            $$ = []interface{}{"or", $1,$3}
        } |
        expr TAnd expr {
            $$ = []interface{}{"and", $1,$3}
        } |
        expr '>' expr {
            $$ = []interface{}{">", $1,$3}
        } |
        expr '<' expr {
            $$ = []interface{}{"<", $1,$3}
        } |
        expr TGte expr {
            $$ = []interface{}{">=", $1,$3}
        } |
        expr TLte expr {
            $$ = []interface{}{"<=", $1,$3}
        } |
        expr TEqeq expr {
            $$ = []interface{}{"==", $1,$3}
        } |
        expr TNeq expr {
            $$ = []interface{}{"!=", $1,$3}
        } |
        expr '+' expr {
            $$ = []interface{}{"+", $1,$3}
        } |
        expr '-' expr {
            $$ = []interface{}{"-", $1,$3}
        } |
        expr '*' expr {
            $$ = []interface{}{"*", $1,$3}
        } |
        expr '/' expr {
            $$ = []interface{}{"/", $1,$3}
        } |
        expr '%' expr {
            $$ = []interface{}{"%", $1,$3}
        } |
        expr '^' expr {
            $$ = []interface{}{"^", $1,$3}
        } |
        '-' expr %prec UNARY {
            $$ = []interface{}{"-", $2}
        } |
        TNot expr %prec UNARY {
            $$ = []interface{}{"not", $2}
        }

string: 
        TString {
            $$ = $1.Str
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
            $$ = []interface{}{"call", $1, $2}
        }

args:
        '(' ')' {
            if yylex.(*Lexer).PNewLine {
               yylex.(*Lexer).TokenError($1, "ambiguous syntax (function call x new statement)")
            }
            $$ = []interface{}{}
        } |
        '(' exprlist ')' {
            if yylex.(*Lexer).PNewLine {
               yylex.(*Lexer).TokenError($1, "ambiguous syntax (function call x new statement)")
            }
            $$ = $2
        }

function:
        TFunction funcbody {
            $$ = []interface{}{"lambda", $2.([]interface{})[0], $2.([]interface{})[1]}
        }

funcbody:
        '(' parlist ')' block TEnd {
            $$ = []interface{}{$2, $4}
        } | 
        '(' ')' block TEnd {
            $$ = []interface{}{nil, $3}
        }

parlist:
        namelist {
          $$ = $1
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

