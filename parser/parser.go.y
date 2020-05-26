%{
package parser

import "strconv"
import "math/rand"

%}
%type<expr> stats
%type<expr> stat
%type<expr> declarator
%type<expr> declarator_list_assign
%type<expr> _declarator_list_assign
%type<expr> ident_list
%type<expr> expr_list
%type<expr> expr_list_paren
%type<expr> expr_assign_list
%type<expr> expr
%type<expr> postfix_incdec
%type<atom> _postfix_assign
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
%type<expr> table_gen

%union {
  token Token
  expr  *Node
  atom  Atom
}

/* Reserved words */
%token<token> TDo TIn TLocal TElseIf TThen TEnd TBreak TContinue TElse TFor TWhile TFunc TIf TLen TReturn TReturnVoid TImport TYield TYieldVoid TRepeat TUntil TNot

/* Literals */
%token<token> TEqeq TNeq TLsh TRsh TURsh TLte TGte TIdent TNumber TString '{' '[' '('
%token<token> TAddEq TSubEq TMulEq TDivEq TModEq TBitAndEq TBitOrEq TXorEq TLshEq TRshEq TURshEq
%token<token> TSquare TDotDotDot TDotDot TSet

/* Operators */
%right 'T'
%right TElse
%left ASSIGN
%right FUNC
%left TDotDotDot
%left TOr
%left TAnd
%left '>' '<' TGte TLte TEqeq TNeq
%left '+' '-' '|' '^' TDotDot
%left '*' '/' '%' TLsh TRsh TURsh '&'
%right UNARY /* not # -(unary) */
%right '#'
%right TTypeof, TLen, TImport

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
        TAddEq         { $$ = AAdd } |
        TSubEq         { $$ = ASub } |
        TMulEq         { $$ = AMul } |
        TDivEq         { $$ = ADiv } |
        TModEq         { $$ = AMod } |
        TBitAndEq      { $$ = ABitAnd } |
        TBitOrEq       { $$ = ABitOr } |
        TXorEq         { $$ = ABitXor } |
        TLshEq         { $$ = ABitLsh } |
        TRshEq         { $$ = ABitRsh } |
        TURshEq        { $$ = ABitURsh }

assign_stat:
        prefix_expr {
            $$ = $1
        } | 
        postfix_incdec {
            $$ = $1
        } |
        declarator_list_assign expr_list {
            nodes := $1.C()
            op := nodes[len(nodes) - 1]
            nodes = nodes[:len(nodes) - 1]

            if len(nodes) == 2 && $2.Cn() == 1 { // local? a, b = c()
                bb := CompNode(AGetB)
                if nodes[0].Type() == Natom && nodes[1].Type() == Natom {
                    $$ = __chain(
                        opSetMove(op)(nodes[0], $2.Cx(0)).pos0($1),
                        opSetMove(op)(nodes[1], bb).pos0($1),
                    )
                } else {
                    $$ = __do(
                        __set("(1)a", $2.Cx(0)).pos0($1),
                        __set("(1)b", bb).pos0($1),
                        nodes[0].moveLoadStore(__move, ANodeS("(1)a")).pos0(nodes[0]),
                        nodes[1].moveLoadStore(__move, ANodeS("(1)b")).pos0(nodes[1]),
                    )
                }
            } else if len(nodes) != $2.Cn() {
                panic(&Error{Pos: $2.Position, Message: "unmatched assignments", Token: string(op.A())})
            } else if op.A() == ASet { // local a0, ..., an = b0, ..., bn
                $$ = __chain()
                for i, v := range nodes {
                    $$ = $$.Cappend(__set(v, $2.Cx(i)).pos0($1))
                }
            } else if head := nodes[0]; len(nodes) == 1 { // a0 = b0
                $$ = head.moveLoadStore(__move, $2.Cx(0)).pos0($1)
                if a, s := $2.Cx(0).isSimpleAddSub(); a != "" && a == head.A() {
                    // Note that a := a + v is different
                    $$ = __inc(head, NewNumberNode(s)).pos0($1)
                }
            } else { // a0, ..., an = b0, ..., bn
                $$ = __chain()
                names := []*Node{}
                for i := range nodes {
                    names = append(names, ANodeS("(1)a" + strconv.Itoa(i)))
                    $$.Cappend(__set(names[i], $2.Cx(i)).pos0($1))
                }
                for i, v := range nodes {
                    $$.Cappend(v.moveLoadStore(__move, names[i]).pos0($1))
                }
            }
        } 

postfix_incdec:
        TIdent _postfix_assign expr %prec ASSIGN  {
            $$ = __move(ANode($1), CompNode($2, ANode($1).setPos($1), $3)).pos0($1)
        } |
        prefix_expr '[' expr ']' _postfix_assign expr %prec ASSIGN {
            $$ = __store($1, $3, CompNode($5, __load($1, $3).pos0($1), $6).pos0($1))
        } |
        prefix_expr '.' TIdent _postfix_assign expr %prec ASSIGN {
            $$ = __store($1, NewNode($3.Str), CompNode($4, __load($1, NewNode($3.Str)).pos0($1), $5).pos0($1))
        }

for_stat:
        TWhile expr TDo stats TEnd {
            $$ = __for(
                __chain(
                    __if($2).__then($4).__else(breakNode).pos0($1),
                ).pos0($1),
            ).pos0($1)
        } |
        TRepeat stats TUntil expr {
            $$ = __for(
                __chain(
                    $2,
                    __if($4).__then(breakNode).__else(emptyNode).pos0($1),
                ).pos0($1),
            ).pos0($1)
        } |
        TFor TIdent TIn expr TDo stats TEnd {
            iter := randomVarname()
            $$ = __chain(
                __set(iter, $4).pos0($1),
                __for(
                    __chain(
                        __set($2, __call(iter, emptyNode).pos0($1)).pos0($1),
                        __if($2).__then($6).__else(breakNode).pos0($1),
                    ),
                ).pos0($1),
            )
        } |
        TFor TIdent ',' TIdent TIn expr TDo stats TEnd {
            iter := randomVarname()
            $$ = __chain(
                __set(iter, $6).pos0($1),
                __for(
                    __chain(
                        __set($2, __call(iter, emptyNode).pos0($1)).pos0($1),
                        __set($4, CompNode(AGetB).pos0($1)).pos0($1),
                        __if($2).__then($8).__else(breakNode).pos0($1),
                    ),
                ).pos0($1),
            )
        } |
        TFor TIdent '=' expr ',' expr TDo stats TEnd {
            forVar, forEnd := ANode($2), randomVarname()
            $$ = __do(
                    __set(forVar, $4).pos0($1),
                    __set(forEnd, $6).pos0($1),
                    __for(
                        __chain(
                            __if(__lessEq(forVar, forEnd)).
                            __then(
                                __chain(
                                    $8,
                                    __inc(forVar, oneNode),
                                ),
                            ).
                            __else(CompNode(ABreak).pos0($1)).pos0($1),
                        ),
                    ).pos0($1),
                )
        } |
        TFor TIdent '=' expr ',' expr ',' expr TDo stats TEnd {
            forVar, forEnd := ANode($2), randomVarname()
            if $8.Type() == Nnumber { // step is a static number, easy case
                var cond *Node
                if $8.N() < 0 {
                    cond = __lessEq(forEnd, forVar)
                } else {
                    cond = __lessEq(forVar, forEnd)
                }
                $$ = __do(
                    __set(forVar, $4).pos0($1),
                    __set(forEnd, $6).pos0($1),
                    __for(
                        __chain(
                            __if(cond).
                            __then(
                                __chain($10, __inc(forVar, $8)),
                            ).
                            __else(breakNode).pos0($1),
                        ),
                    ).pos0($1),
                )
            } else { 
                forStep := randomVarname()
                $$ = __do(
                    __set(forVar, $4).pos0($1),
                    __set(forEnd, $6).pos0($1),
                    __set(forStep, $8).pos0($1),
                    __for(
                        __chain(
                            __if(__less(zeroNode, forStep)).
                            __then( // +step
                                __if(__less(forEnd, forVar)).__then(breakNode).__else(emptyNode).pos0($1),
                            ).
                            __else( // -step
                                __if(__less(forVar, forEnd)).__then(breakNode).__else(emptyNode).pos0($1),
                            ).pos0($1),
                            $10,
                            __inc(forVar, forStep),
                        ),
                    ).pos0($1),
                )
            }
            
        } 

if_stat:
        TIf expr TThen stats elseif_stat TEnd %prec 'T' {
            $$ = __if($2).__then($4).__else($5).pos0($1)
        }

elseif_stat:
        {
            $$ = CompNode()
        } |
        TElse stats {
            $$ = $2
        } |
        TElseIf expr TThen stats elseif_stat {
            $$ = __if($2).__then($4).__else($5).pos0($1)
        }

func:
        TFunc {
            $$ = NewNode(AMove).SetPos($1)
        } |
        TLocal TFunc {
            $$ = NewNode(ASet).SetPos($1)
        }

func_stat:
        func TIdent func_params_list stats TEnd {
            funcname := ANode($2)
            $$ = __chain(
                opSetMove($1)(funcname, nilNode).pos0($2), 
                __move(funcname, __func($3).__body($4).pos0($2)).pos0($2),
            )
        } |
        func TIdent '.' TIdent func_params_list stats TEnd {
            $$ = __store(
                ANode($2), NewNode($4.Str), __func($5).__body($6).pos0($2),
            ).pos0($2) 
        } |
        func TIdent ':' TIdent func_params_list stats TEnd {
            paramlist := $5.Cprepend(ANodeS("self"))
            $$ = __store(
                ANode($2), NewNode($4.Str), __func(paramlist).__body($6).pos0($2),
            ).pos0($2) 
        }

function:
        func func_params_list stats TEnd %prec FUNC {
            $$ = __func($2).__body($3).pos0($1).SetPos($1) 
        }

func_params_list:
        '(' ')'                           { $$ = CompNode() } |
        '(' ident_list ')'                { $$ = $2 }

jmp_stat:
        TYield expr                       { $$ = CompNode(AYield, $2).pos0($1) } |
        TYield expr ',' expr              { $$ = __chain(CompNode(ASetB, $4).pos0($1), CompNode(AYield, $2).pos0($1)) } |
        TYieldVoid                        { $$ = CompNode(AYield, nilNode).pos0($1) } |
        TBreak                            { $$ = CompNode(ABreak).pos0($1) } |
        TContinue                         { $$ = CompNode(AContinue).pos0($1) } |
        TReturn expr                      { $$ = __return($2).pos0($1) } |
        TReturn expr ',' expr             { $$ = __chain(CompNode(ASetB, $4).pos0($1), __return($2).pos0($1)) } |
        TReturnVoid                       { $$ = __return(nilNode).pos0($1) } |
        TImport TString                   { $$ = yylex.(*Lexer).loadFile(joinSourcePath($1.Pos.Source, $2.Str), $1) }

declarator:
        TIdent                            { $$ = ANode($1).setPos($1) } |
        TIdent TSquare                    { $$ = __load(nilNode, $1).pos0($1) } |
        prefix_expr '[' expr ']'          { $$ = __load($1, $3).pos0($3).setPos($3) } |
        prefix_expr '.' TIdent            { $$ = __load($1, NewNode($3.Str)).pos0($3).setPos($3) }

declarator_list_assign:
        TLocal _declarator_list_assign    { a := $2.Value.([]*Node); a[len(a)-1] = NewNode(ASet); $$ = $2 } |
        _declarator_list_assign           { $$ = $1 }

_declarator_list_assign:
        declarator '='                    { $$ = CompNode($1, NewNode(AMove)) } |
        declarator ',' declarator_list_assign { $$ = $3.Cprepend($1) }

ident_list:
        TIdent                            { $$ = CompNode($1.Str) } | 
        TDotDotDot                        { $$ = CompNode("...") } | 
        ident_list ',' TIdent             { $$ = $1.Cappend(ANode($3)) } |
        ident_list ',' TDotDotDot         { $$ = $1.Cappend(ANodeS("...").SetPos($3)) }

expr:
        TNumber                           { $$ = NewNumberNode($1.Str).SetPos($1) } |
        TImport TString                   { $$ = yylex.(*Lexer).loadFile(joinSourcePath($1.Pos.Source, $2.Str), $1) } |
        '#' expr                          { $$ = CompNode(ALen, $2) } |
        function                          { $$ = $1 } |
        table_gen                  { $$ = $1 } |
        prefix_expr                       { $$ = $1 } |
        TString                           { $$ = NewNode($1.Str).SetPos($1) }  |
        expr TOr expr                     { $$ = CompNode(AOr, $1,$3).pos0($1) } |
        expr TAnd expr                    { $$ = CompNode(AAnd, $1,$3).pos0($1) } |
        expr '>' expr                     { $$ = CompNode(ALess, $3,$1).pos0($1) } |
        expr '<' expr                     { $$ = CompNode(ALess, $1,$3).pos0($1) } |
        expr TGte expr                    { $$ = CompNode(ALessEq, $3,$1).pos0($1) } |
        expr TLte expr                    { $$ = CompNode(ALessEq, $1,$3).pos0($1) } |
        expr TEqeq expr                   { $$ = CompNode(AEq, $1,$3).pos0($1) } |
        expr TNeq expr                    { $$ = CompNode(ANeq, $1,$3).pos0($1) } |
        expr '+' expr                     { $$ = CompNode(AAdd, $1,$3).pos0($1) } |
        expr TDotDot expr                 { $$ = CompNode(AConcat, $1,$3).pos0($1) } |
        expr '-' expr                     { $$ = CompNode(ASub, $1,$3).pos0($1) } |
        expr '*' expr                     { $$ = CompNode(AMul, $1,$3).pos0($1) } |
        expr '/' expr                     { $$ = CompNode(ADiv, $1,$3).pos0($1) } |
        expr '%' expr                     { $$ = CompNode(AMod, $1,$3).pos0($1) } |
        expr '^' expr                     { $$ = CompNode(ABitXor, $1,$3).pos0($1) } |
        expr TLsh expr                    { $$ = CompNode(ABitLsh, $1,$3).pos0($1) } |
        expr TRsh expr                    { $$ = CompNode(ABitRsh, $1,$3).pos0($1) } |
        expr TURsh expr                   { $$ = CompNode(ABitURsh, $1,$3).pos0($1) } |
        expr '|' expr                     { $$ = CompNode(ABitOr, $1,$3).pos0($1) } |
        expr '&' expr                     { $$ = CompNode(ABitAnd, $1,$3).pos0($1) } |
        '^' expr %prec UNARY              { $$ = CompNode(ABitXor, $2, max32Node).pos0($2) } |
        '-' expr %prec UNARY              { $$ = CompNode(ASub, zeroNode, $2).pos0($2) } |
        TNot expr %prec UNARY             { $$ = CompNode(ANot, $2).pos0($2) } |
        '&' TIdent %prec UNARY            { $$ = CompNode(AAddrOf, ANode($2)).pos0($2) }

prefix_expr:
        declarator                        { $$ = $1 } |
        prefix_expr TString               { $$ = __call($1, CompNode(NewNode($2.Str))).pos0($1) } |
        TIdent ':' TIdent expr_list_paren { $$ = __call(__load($1, NewNode($3.Str)).pos0($1), $4.Cprepend(ANode($1))).pos0($1) } |
        prefix_expr expr_list_paren       { $$ = __call($1, $2).pos0($1) } |
        '(' expr ')'                      { $$ = $2 } // shift/reduce conflict

expr_list:
        expr                              { $$ = CompNode($1) } |
        expr_list ',' expr                { $$ = $1.Cappend($3) }

expr_list_paren:
        '(' ')'                           { $$ = CompNode() } |
        '{' '}'                           { $$ = CompNode() } |
        '(' expr_list ')'                 { $$ = $2 } |
        '{' expr_list '}'                 { $$ = $2 }

expr_assign_list:
        TIdent '=' expr                            { $$ = CompNode(NewNode($1.Str), $3) } |
        '[' expr ']' '=' expr                      { $$ = CompNode($2, $5) } |
        expr_assign_list ',' TIdent '=' expr       { $$ = $1.Cappend(NewNode($3.Str), $5) } |
        expr_assign_list ',' '[' expr ']' '=' expr { $$ = $1.Cappend($4, $7) }

table_gen:
        '{' '}'                                    { $$ = CompNode(AArray, emptyNode).pos0($1) } |
        '{' expr_assign_list     '}'               { $$ = CompNode(AHash, $2).pos0($1) } |
        '{' expr_assign_list ',' '}'               { $$ = CompNode(AHash, $2).pos0($1) } |
        '{' expr_assign_list ';' expr_list '}'     { $$ = CompNode(AHashArray, $2, $4).pos0($1) } |
        '{' expr_assign_list ';' expr_list ',' '}' { $$ = CompNode(AHashArray, $2, $4).pos0($1) } |
        '{' expr_list            '}'               { $$ = CompNode(AArray, $2).pos0($1) } |
        '{' expr_list ';' expr_assign_list '}'     { $$ = CompNode(AHashArray, $4, $2).pos0($1) } |
        '{' expr_list ';' expr_assign_list ',' '}' { $$ = CompNode(AHashArray, $4, $2).pos0($1) } |
        '{' expr_list ','        '}'               { $$ = CompNode(AArray, $2).pos0($1) }

%%

func opSetMove(op *Node) func(dest, src interface{}) *Node {
    if op.A() == ASet {
        return __set
    }
    return __move
}

func randomVarname() *Node {
    return ANodeS("v" + strconv.FormatInt(rand.Int63(), 10))
}
