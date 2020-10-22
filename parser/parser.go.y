%{
package parser
%}
%type<expr> stats
%type<expr> stat
%type<expr> declarator
%type<expr> declarator_list
%type<expr> ident_list
%type<expr> expr_list
%type<expr> expr_list_paren
%type<expr> expr
%type<expr> postfix_incdec
%type<expr> _postfix_assign
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

%union {
  token Token
  expr  Node
  atom  Symbol
}

/* Reserved words */
%token<token> TDo TIn TLocal TElseIf TThen TEnd TBreak TElse TFor TWhile TFunc TIf TLen TReturn TReturnVoid TImport TYield TYieldVoid TRepeat TUntil TNot TLabel TGoto

/* Literals */
%token<token> TOr TAnd TEqeq TNeq TLte TGte TIdent TNumber TString 
%token<token> '{' '[' '(' '=' '>' '<' '+' '-' '*' '/' '%' '^' '#' '.' '&'
%token<token> TAddEq TSubEq TMulEq TDivEq TModEq
%token<token> TSquare TDotDot 

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
        TAddEq         { $$ = Node{AAdd.SetPos($1.Pos)} } |
        TSubEq         { $$ = Node{ASub.SetPos($1.Pos)} } |
        TMulEq         { $$ = Node{AMul.SetPos($1.Pos)} } |
        TDivEq         { $$ = Node{ADiv.SetPos($1.Pos)} } |
        TModEq         { $$ = Node{AMod.SetPos($1.Pos)} }

assign_stat:
        prefix_expr {
            if $1.isCallStat() {
                // Single call statement, clear env.V to avoid side effects
                $$ = __chain($1, popvClearNode)
            } else {
                $$ = $1
            }
        } | 
        postfix_incdec {
            $$ = $1
        } |
        TLocal ident_list {
            $$ = __chain()
            for _, v := range $2.Cpl() {
                $$ = $$.CplAppend(__set(v, Node{ANil}).SetPos($1.Pos))
            }
        } |
        TLocal ident_list '=' expr_list {
            m, n := len($2.Cpl()), len($4.Cpl())
            for i, count := 0, m - n; i < count; i++ {
                if i == count - 1 {
                    $4 = $4.CplAppend(__chain(popvNode, popvClearNode))
		} else {
		  $4 = $4.CplAppend(popvNode)
		}
            }

            $$ = __chain()
            for i, v := range $2.Cpl() {
                if v.SymDDD() { 
                    $$ = $$.CplAppend(__set(v, __popvAll(i, $4.CplIndex(i))).SetPos($1.Pos))
                } else {
                    $$ = $$.CplAppend(__set(v, $4.CplIndex(i)).SetPos($1.Pos))
                }
            }

            if m == 1 && n == 1 && $4.CplIndex(0).isCallStat() {
                // Single call statement with single assignment, clear env.V to avoid side effects
                $$ = $$.CplAppend(popvClearNode)
            }
        } |
        declarator_list '=' expr_list {
            nodes := $1.Cpl()
            m, n := len(nodes), len($3.Cpl())
            for i, count := 0, m - n; i < count; i++ {
                if i == count - 1 {
		    $3 = $3.CplAppend(__chain(popvNode, popvClearNode))
		} else {
		    $3 = $3.CplAppend(popvNode)
		}
            } 
             
	    if head := nodes[0]; len(nodes) == 1 && !nodes[0].SymDDD() {
                // a0 = b0
                // if a, s, ok := $3.CplIndex(0).isSimpleAddSub(); ok && a.Equals(head.Sym()) {
                //    $$ = __inc(head, Num(s)).SetPos($2.Pos)
                // } else {
                    $$ = head.moveLoadStore(__move, $3.CplIndex(0)).SetPos($2.Pos)
                // }
            } else { 
                // a0, ..., an = b0, ..., bn
                $$ = __chain()
                names, retaddr := []Node{}, Cpl(Node{ARetAddr})
                for i := range nodes {
                    names = append(names, randomVarname())
                    retaddr = retaddr.CplAppend(names[i])
                    if nodes[i].SymDDD() {
                        $$ = $$.CplAppend(__set(names[i], __popvAll(i, $3.CplIndex(i))).SetPos($2.Pos))
                    } else {
                        $$ = $$.CplAppend(__set(names[i], $3.CplIndex(i)).SetPos($2.Pos))
                    }
                }
                for i, v := range nodes {
                    $$ = $$.CplAppend(v.moveLoadStore(__move, names[i]).SetPos($2.Pos))
                }
                $$ = $$.CplAppend(retaddr)
            }

            if m == 1 && n == 1 && $3.CplIndex(0).isCallStat() {
                // Single call statement with single assignment, clear env.V to avoid side effects
                $$ = __chain($$, popvClearNode)
            }
        } 

postfix_incdec:
        TIdent _postfix_assign expr %prec ASSIGN  {
            $$ = __move(SymTok($1), Cpl($2, SymTok($1), $3)).SetPos($2.Pos())
        } |
        prefix_expr '[' expr ']' _postfix_assign expr %prec ASSIGN {
            $$ = __store($1, $3, Cpl($5, __load($1, $3), $6).SetPos($5.Pos()))
        } |
        prefix_expr '.' TIdent _postfix_assign expr %prec ASSIGN {
            i := Node{$3.Str}
            $$ = __store($1, i, Cpl($4, __load($1, i), $5).SetPos($4.Pos()))
        }

for_stat:
        TWhile expr TDo stats TEnd {
            $$ = __loop(__if($2, $4, breakNode).SetPos($1.Pos)).SetPos($1.Pos)
        } |
        TRepeat stats TUntil expr {
            $$ = __loop(
                __chain(
                    $2,
                    __if($4, breakNode, emptyNode).SetPos($1.Pos),
                ).SetPos($1.Pos),
            ).SetPos($1.Pos)
        } |
        TFor ident_list TIn expr_list TDo stats TEnd {
            $$ = forLoop($1.Pos, $2.Cpl(), $4.Cpl(), $6)
        } |
        TFor TIdent '=' expr ',' expr TDo stats TEnd {
            forVar, forEnd := SymTok($2), randomVarname()
            $$ = __do(
                    __set(forVar, $4).SetPos($1.Pos),
                    __set(forEnd, $6).SetPos($1.Pos),
                    __loop(
                        __if(
                            __lessEq(forVar, forEnd),
                            __chain($8, __inc(forVar, oneNode).SetPos($1.Pos)),
                            breakNode,
                        ).SetPos($1.Pos),
                    ).SetPos($1.Pos),
                )
        } |
        TFor TIdent '=' expr ',' expr ',' expr TDo stats TEnd {
            forVar, forEnd := SymTok($2), randomVarname()
            if $8.Type() == NUM { // step is a static number, easy case
                var cond Node
                if f, i := $8.Num(); f < 0 || i < 0 {
                    cond = __lessEq(forEnd, forVar)
                } else {
                    cond = __lessEq(forVar, forEnd)
                }
                $$ = __do(
                    __set(forVar, $4).SetPos($1.Pos),
                    __set(forEnd, $6).SetPos($1.Pos),
                    __loop(
                        __chain(
                            __if(
                                cond,
                                __chain($10, __inc(forVar, $8)),
                                breakNode,
                            ).SetPos($1.Pos),
                        ),
                    ).SetPos($1.Pos),
                )
            } else { 
                forStep := randomVarname()
                $$ = __do(
                    __set(forVar, $4).SetPos($1.Pos),
                    __set(forEnd, $6).SetPos($1.Pos),
                    __set(forStep, $8).SetPos($1.Pos),
                    __loop(
                        __chain(
                            __if(
                                __less(zeroNode, forStep).SetPos($1.Pos),
                                // +step
                                __if(__less(forEnd, forVar), breakNode, emptyNode).SetPos($1.Pos),
                                // -step
                                __if(__less(forVar, forEnd), breakNode, emptyNode).SetPos($1.Pos),
                            ).SetPos($1.Pos),
                            $10,
                            __inc(forVar, forStep),
                        ),
                    ).SetPos($1.Pos),
                )
            }
            
        } 

if_stat:
        TIf expr TThen stats elseif_stat TEnd %prec 'T' {
            $$ = __if($2, $4, $5).SetPos($1.Pos)
        }

elseif_stat:
        {
            $$ = Cpl()
        } |
        TElse stats {
            $$ = $2
        } |
        TElseIf expr TThen stats elseif_stat {
            $$ = __if($2, $4, $5).SetPos($1.Pos)
        }

func:
        TFunc        { $$ = Node{AMove}.SetPos($1.Pos) } |
        TLocal TFunc { $$ = Node{ASet}.SetPos($1.Pos) }

func_stat:
        func TIdent func_params_list stats TEnd {
            funcname := SymTok($2)
            x := __move
            if $1.Sym().Equals(ASet) {
                x = __set
            }
            $$ = __chain(
                x(funcname, Node{ANil}).SetPos($1.Pos()), 
                __move(funcname, __func(funcname, $3, $4).SetPos($1.Pos())).SetPos($1.Pos()),
            )
        }

function:
        func func_params_list stats TEnd %prec FUNC {
	    $$ = __func(emptyNode, $2, $3).SetPos($1.Pos()).SetPos($1.Pos()) 
        }

func_params_list:
        '(' ')'                           { $$ = Cpl() } |
        '(' ident_list ')'                { $$ = $2 }

jmp_stat:
        TYield expr_list                  { $$ = Cpl(Node{AYield}, $2).SetPos($1.Pos) } |
        TYieldVoid                        { $$ = Cpl(Node{AYield}, emptyNode).SetPos($1.Pos) } |
        TBreak                            { $$ = Cpl(Node{ABreak}).SetPos($1.Pos) } |
        TImport TString                   { $$ = yylex.(*Lexer).loadFile(joinSourcePath($1.Pos.Source, $2.Str)) } |
        TGoto TIdent                      { $$ = Cpl(Node{AGoto}, SymTok($2)).SetPos($1.Pos) } |
        TLabel TIdent TLabel              { $$ = Cpl(Node{ALabel}, SymTok($2)) } |
        TReturnVoid                       { $$ = Cpl(Node{AReturn}, emptyNode).SetPos($1.Pos) } |
        TReturn expr_list                 {
            if len($2.Cpl()) == 1 {
                x := $2.CplIndex(0)
                if len(x.Cpl()) == 3 && x.CplIndex(0).Sym().Equals(ACall) {
                    tc := x.CplIndex(0).Sym()
                    tc.Text = ATailCall.Text
                    x.Value.([]Node)[0] = Node{tc}
                }
            }
            $$ = Cpl(Node{AReturn}, $2).SetPos($1.Pos) 
        }

declarator:
        TIdent                            { $$ = SymTok($1) } |
        prefix_expr '[' expr ']'          { $$ = __load($1, $3).SetPos($2.Pos) /* (10)[0] is valid if number has metamethod */ } |
        prefix_expr '.' TIdent            { $$ = __load($1, Node{$3.Str}).SetPos($2.Pos) }

declarator_list:
        declarator                        { $$ = Cpl($1) } |
        declarator_list ',' declarator    { $$ = $1.CplAppend($3) }

ident_list:
        TIdent                            { $$ = Cpl(SymTok($1)) } | 
        ident_list ',' TIdent             { $$ = $1.CplAppend(SymTok($3)) }

expr:
        TNumber                           { $$ = Num($1.Str) } |
        function                          { $$ = $1 } |
        TString                           { $$ = Node{$1.Str} } |
	prefix_expr                       { $$ = $1 } |
        expr TOr expr                     { $$ = Cpl(Node{AOr}, $1,$3).SetPos($2.Pos) } |
        expr TAnd expr                    { $$ = Cpl(Node{AAnd}, $1,$3).SetPos($2.Pos) } |
        expr '>' expr                     { $$ = Cpl(Node{ALess}, $3,$1).SetPos($2.Pos) } |
        expr '<' expr                     { $$ = Cpl(Node{ALess}, $1,$3).SetPos($2.Pos) } |
        expr TGte expr                    { $$ = Cpl(Node{ALessEq}, $3,$1).SetPos($2.Pos) } |
        expr TLte expr                    { $$ = Cpl(Node{ALessEq}, $1,$3).SetPos($2.Pos) } |
        expr TEqeq expr                   { $$ = Cpl(Node{AEq}, $1,$3).SetPos($2.Pos) } |
        expr TNeq expr                    { $$ = Cpl(Node{ANeq}, $1,$3).SetPos($2.Pos) } |
        expr '+' expr                     { $$ = Cpl(Node{AAdd}, $1,$3).SetPos($2.Pos) } |
        expr TDotDot expr                 { $$ = Cpl(Node{AConcat}, $1,$3).SetPos($2.Pos) } |
        expr '-' expr                     { $$ = Cpl(Node{ASub}, $1,$3).SetPos($2.Pos) } |
        expr '*' expr                     { $$ = Cpl(Node{AMul}, $1,$3).SetPos($2.Pos) } |
        expr '/' expr                     { $$ = Cpl(Node{ADiv}, $1,$3).SetPos($2.Pos) } |
        expr '%' expr                     { $$ = Cpl(Node{AMod}, $1,$3).SetPos($2.Pos) } |
        expr '^' expr                     { $$ = Cpl(Node{APow}, $1,$3).SetPos($2.Pos) } |
        TNot expr %prec UNARY             { $$ = Cpl(Node{ANot}, $2).SetPos($1.Pos) } |
        '-' expr %prec UNARY              { $$ = Cpl(Node{ASub}, zeroNode, $2).SetPos($1.Pos) } |
        '#' expr %prec UNARY              { $$ = Cpl(Node{ALen}, $2).SetPos($1.Pos) }

prefix_expr:
        declarator                        { $$ = $1 } |
        prefix_expr TString               { $$ = __call($1, Cpl(Node{$2.Str})).SetPos($1.Pos()) } |
        prefix_expr expr_list_paren       { $$ = __call($1, $2).SetPos($1.Pos()) } |
        '(' expr ')'                      { $$ = $2 } // shift/reduce conflict

expr_list:
        expr                              { $$ = Cpl($1) } |
        expr_list ',' expr                { $$ = $1.CplAppend($3) }

expr_list_paren:
        '(' ')'                           { $$ = Cpl() } |
        '(' expr_list ')'                 { $$ = $2 }

%%

