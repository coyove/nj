%{
package parser
%}
%type<stmts> chunk
%type<stmts> chunk1
%type<stmts> block
%type<stmt>  stat
%type<stmts> elseifs
%type<stmt>  laststat
%type<funcname> funcname
%type<funcname> funcname1
%type<exprlist> varlist
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
%type<expr> tableconstructor
%type<fieldlist> fieldlist
%type<field> field
%type<fieldsep> fieldsep

%union {
  token  Token

  stmts    []Stmt
  stmt     Stmt

  funcname *FuncName
  funcexpr *FunctionExpr

  exprlist []Expr
  expr   Expr

  fieldlist []*Field
  field     *Field
  fieldsep  string

  namelist []string
  parlist  *ParList
}

/* Reserved words */
%token<token> TAnd TBreak TContinue TDo TElse TElseIf TEnd TFalse TFor TFunction TIf TIn TLocal TNil TNot TOr TReturn TRepeat TThen TTrue TUntil TWhile 

/* Literals */
%token<token> TEqeq TNeq TLte TGte T2Comma T3Comma TIdent TNumber TString '{' '('

/* Operators */
%left TOr
%left TAnd
%left '>' '<' TGte TLte TEqeq TNeq
%right T2Comma
%left '+' '-'
%left '*' '/' '%'
%right UNARY /* not # -(unary) */
%right '^'

%%

chunk: 
        chunk1 {
            $$ = $1
            if l, ok := yylex.(*Lexer); ok {
                l.Stmts = $$
            }
        } |
        chunk1 laststat {
            $$ = append($1, $2)
            if l, ok := yylex.(*Lexer); ok {
                l.Stmts = $$
            }
        } | 
        chunk1 laststat ';' {
            $$ = append($1, $2)
            if l, ok := yylex.(*Lexer); ok {
                l.Stmts = $$
            }
        }

chunk1: 
        {
            $$ = []Stmt{}
        } |
        chunk1 stat {
            $$ = append($1, $2)
        } | 
        chunk1 ';' {
            $$ = $1
        }

block: 
        chunk {
            $$ = $1
        }

stat:
        varlist '=' exprlist {
            $$ = &AssignStmt{Lhs: $1, Rhs: $3}
            $$.SetLine($1[0].Line())
        } |
        /* 'stat = functioncal' causes a reduce/reduce conflict */
        prefixexp {
            if _, ok := $1.(*FuncCallExpr); !ok {
               yylex.(*Lexer).Error("parse error")
            } else {
              $$ = &FuncCallStmt{Expr: $1}
              $$.SetLine($1.Line())
            }
        } |
        TWhile expr TDo block TEnd {
            $$ = &WhileStmt{Condition: $2, Stmts: $4}
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($5.Pos.Line)
        } |
        TIf expr TThen block elseifs TEnd {
            $$ = &IfStmt{Condition: $2, Then: $4}
            cur := $$
            for _, elseif := range $5 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($6.Pos.Line)
        } |
        TIf expr TThen block elseifs TElse block TEnd {
            $$ = &IfStmt{Condition: $2, Then: $4}
            cur := $$
            for _, elseif := range $5 {
                cur.(*IfStmt).Else = []Stmt{elseif}
                cur = elseif
            }
            cur.(*IfStmt).Else = $7
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($8.Pos.Line)
        } |
        TFunction funcname funcbody {
            $$ = &FuncDefStmt{Name: $2, Func: $3}
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($3.LastLine())
        } |
        TLocal TFunction TIdent funcbody {
            $$ = &LocalAssignStmt{Names:[]string{$3.Str}, Exprs: []Expr{$4}}
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($4.LastLine())
        } | 
        TLocal namelist '=' exprlist {
            $$ = &LocalAssignStmt{Names: $2, Exprs:$4}
            $$.SetLine($1.Pos.Line)
        } |
        TLocal namelist {
            $$ = &LocalAssignStmt{Names: $2, Exprs:[]Expr{}}
            $$.SetLine($1.Pos.Line)
        }

elseifs: 
        {
            $$ = []Stmt{}
        } | 
        elseifs TElseIf expr TThen block {
            $$ = append($1, &IfStmt{Condition: $3, Then: $5})
            $$[len($$)-1].SetLine($2.Pos.Line)
        }

laststat:
        TReturn {
            $$ = &ReturnStmt{Exprs:nil}
            $$.SetLine($1.Pos.Line)
        } |
        TReturn exprlist {
            $$ = &ReturnStmt{Exprs:$2}
            $$.SetLine($1.Pos.Line)
        } |
        TBreak  {
            $$ = &BreakStmt{}
            $$.SetLine($1.Pos.Line)
        } |
        TContinue  {
            $$ = &ContinueStmt{}
            $$.SetLine($1.Pos.Line)
        }

funcname: 
        funcname1 {
            $$ = $1
        } |
        funcname1 ':' TIdent {
            $$ = &FuncName{Func:nil, Receiver:$1.Func, Method: $3.Str}
        }

funcname1:
        TIdent {
            $$ = &FuncName{Func: &IdentExpr{Value:$1.Str}}
            $$.Func.SetLine($1.Pos.Line)
        } | 
        funcname1 '.' TIdent {
            key:= &StringExpr{Value:$3.Str}
            key.SetLine($3.Pos.Line)
            fn := &AttrGetExpr{Object: $1.Func, Key: key}
            fn.SetLine($3.Pos.Line)
            $$ = &FuncName{Func: fn}
        }

varlist:
        var {
            $$ = []Expr{$1}
        } | 
        varlist ',' var {
            $$ = append($1, $3)
        }

var:
        TIdent {
            $$ = &IdentExpr{Value:$1.Str}
            $$.SetLine($1.Pos.Line)
        } |
        prefixexp '[' expr ']' {
            $$ = &AttrGetExpr{Object: $1, Key: $3}
            $$.SetLine($1.Line())
        } | 
        prefixexp '.' TIdent {
            key := &StringExpr{Value:$3.Str}
            key.SetLine($3.Pos.Line)
            $$ = &AttrGetExpr{Object: $1, Key: key}
            $$.SetLine($1.Line())
        }

namelist:
        TIdent {
            $$ = []string{$1.Str}
        } | 
        namelist ','  TIdent {
            $$ = append($1, $3.Str)
        }

exprlist:
        expr {
            $$ = []Expr{$1}
        } |
        exprlist ',' expr {
            $$ = append($1, $3)
        }

expr:
        TNil {
            $$ = &NilExpr{}
            $$.SetLine($1.Pos.Line)
        } | 
        TFalse {
            $$ = &FalseExpr{}
            $$.SetLine($1.Pos.Line)
        } | 
        TTrue {
            $$ = &TrueExpr{}
            $$.SetLine($1.Pos.Line)
        } | 
        TNumber {
            $$ = &NumberExpr{Value: $1.Str}
            $$.SetLine($1.Pos.Line)
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
            $$ = &LogicalOpExpr{Lhs: $1, Operator: "or", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr TAnd expr {
            $$ = &LogicalOpExpr{Lhs: $1, Operator: "and", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr '>' expr {
            $$ = &RelationalOpExpr{Lhs: $1, Operator: ">", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr '<' expr {
            $$ = &RelationalOpExpr{Lhs: $1, Operator: "<", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr TGte expr {
            $$ = &RelationalOpExpr{Lhs: $1, Operator: ">=", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr TLte expr {
            $$ = &RelationalOpExpr{Lhs: $1, Operator: "<=", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr TEqeq expr {
            $$ = &RelationalOpExpr{Lhs: $1, Operator: "==", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr TNeq expr {
            $$ = &RelationalOpExpr{Lhs: $1, Operator: "~=", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr T2Comma expr {
            $$ = &StringConcatOpExpr{Lhs: $1, Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr '+' expr {
            $$ = &ArithmeticOpExpr{Lhs: $1, Operator: "+", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr '-' expr {
            $$ = &ArithmeticOpExpr{Lhs: $1, Operator: "-", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr '*' expr {
            $$ = &ArithmeticOpExpr{Lhs: $1, Operator: "*", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr '/' expr {
            $$ = &ArithmeticOpExpr{Lhs: $1, Operator: "/", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr '%' expr {
            $$ = &ArithmeticOpExpr{Lhs: $1, Operator: "%", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        expr '^' expr {
            $$ = &ArithmeticOpExpr{Lhs: $1, Operator: "^", Rhs: $3}
            $$.SetLine($1.Line())
        } |
        '-' expr %prec UNARY {
            $$ = &UnaryMinusOpExpr{Expr: $2}
            $$.SetLine($2.Line())
        } |
        TNot expr %prec UNARY {
            $$ = &UnaryNotOpExpr{Expr: $2}
            $$.SetLine($2.Line())
        }

string: 
        TString {
            $$ = &StringExpr{Value: $1.Str}
            $$.SetLine($1.Pos.Line)
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
            $$.SetLine($1.Pos.Line)
        }

afunctioncall:
        '(' functioncall ')' {
            $2.(*FuncCallExpr).AdjustRet = true
            $$ = $2
        }

functioncall:
        prefixexp args {
            $$ = &FuncCallExpr{Func: $1, Args: $2}
            $$.SetLine($1.Line())
        } |
        prefixexp ':' TIdent args {
            $$ = &FuncCallExpr{Method: $3.Str, Receiver: $1, Args: $4}
            $$.SetLine($1.Line())
        }

args:
        '(' ')' {
            if yylex.(*Lexer).PNewLine {
               yylex.(*Lexer).TokenError($1, "ambiguous syntax (function call x new statement)")
            }
            $$ = []Expr{}
        } |
        '(' exprlist ')' {
            if yylex.(*Lexer).PNewLine {
               yylex.(*Lexer).TokenError($1, "ambiguous syntax (function call x new statement)")
            }
            $$ = $2
        } |
        tableconstructor {
            $$ = []Expr{$1}
        } | 
        string {
            $$ = []Expr{$1}
        }

function:
        TFunction funcbody {
            $$ = &FunctionExpr{ParList:$2.ParList, Stmts: $2.Stmts}
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($2.LastLine())
        }

funcbody:
        '(' parlist ')' block TEnd {
            $$ = &FunctionExpr{ParList: $2, Stmts: $4}
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($5.Pos.Line)
        } | 
        '(' ')' block TEnd {
            $$ = &FunctionExpr{ParList: &ParList{HasVargs: false, Names: []string{}}, Stmts: $3}
            $$.SetLine($1.Pos.Line)
            $$.SetLastLine($4.Pos.Line)
        }

parlist:
        namelist {
          $$ = &ParList{HasVargs: false, Names: []string{}}
          $$.Names = append($$.Names, $1...)
        }

tableconstructor:
        '{' '}' {
            $$ = &TableExpr{Fields: []*Field{}}
            $$.SetLine($1.Pos.Line)
        } |
        '{' fieldlist '}' {
            $$ = &TableExpr{Fields: $2}
            $$.SetLine($1.Pos.Line)
        }


fieldlist:
        field {
            $$ = []*Field{$1}
        } | 
        fieldlist fieldsep field {
            $$ = append($1, $3)
        } | 
        fieldlist fieldsep {
            $$ = $1
        }

field:
        TIdent '=' expr {
            $$ = &Field{Key: &StringExpr{Value:$1.Str}, Value: $3}
            $$.Key.SetLine($1.Pos.Line)
        } | 
        '[' expr ']' '=' expr {
            $$ = &Field{Key: $2, Value: $5}
        } |
        expr {
            $$ = &Field{Value: $1}
        }

fieldsep:
        ',' {
            $$ = ","
        } | 
        ';' {
            $$ = ";"
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

