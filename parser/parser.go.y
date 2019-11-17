%{
package parser

import "strconv"

%}
%type<expr> stats
%type<expr> block
%type<expr> stat
%type<expr> declarator
%type<expr> declarator_list_assign
%type<expr> ident_list
%type<expr> expr_list
%type<expr> expr_assign_list
%type<expr> expr
%type<expr> postfix_incdec
%type<expr> _postfix_incdec
%type<atom> _postfix_assign
%type<expr> prefix_expr
%type<expr> assign_stat
%type<expr> for_stat
%type<expr> if_stat
%type<expr> oneline_or_block
%type<expr> jmp_stat
%type<expr> func_stat
%type<expr> flow_stat
%type<expr> function
%type<expr> func_params_list
%type<expr> struct_slice_gen

%union {
  token Token
  expr  *Node
  atom  Atom
}

/* Reserved words */
%token<token> TAssert TBreak TContinue TElse TFor TFunc TIf TLen TReturn TReturnVoid TImport TTypeof TYield TYieldVoid TStruct 

/* Literals */
%token<token> TAddAdd TSubSub TEqeq TNeq TLsh TRsh TURsh TLte TGte TIdent TNumber TString '{' '[' '('
%token<token> TAddEq TSubEq TMulEq TDivEq TModEq TBitAndEq TBitOrEq TXorEq TLshEq TRshEq TURshEq
%token<token> TSquare TDotDotDot TSet

/* Operators */
%right 'T'
%right TElse

%left ASSIGN
%right FUNC
%left TDotDotDot
%left TOr
%left TAnd
%left '>' '<' TGte TLte TEqeq TNeq
%left '+' '-' '|' '^'
%left '*' '/' '%' TLsh TRsh TURsh '&'
%right UNARY /* not # -(unary) */
%right '~'
%right '#'
%left TAddAdd TMinMin
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

block: 
        '{' stats '}'  { $$ = $2 }

stat:
        jmp_stat       { $$ = $1 } |
        flow_stat      { $$ = $1 } |
        assign_stat    { $$ = $1 } |
        block          { $$ = $1 } |
        ';'            { $$ = emptyNode }

oneline_or_block:
        assign_stat    { $$ = __chain($1) } |
        jmp_stat       { $$ = __chain($1) } |
        for_stat       { $$ = __chain($1) } |
        if_stat        { $$ = __chain($1) } |
        block          { $$ = $1 }

flow_stat:
        for_stat       { $$ = $1 } |
        if_stat        { $$ = $1 } |
        func_stat      { $$ = $1 }

_postfix_incdec:
        TAddAdd        { $$ = oneNode } |
        TSubSub        { $$ = moneNode }

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

            $$ = __chain()
            if len(nodes) == 2 && $2.Cn() == 1 { // a, b = c()
                y := __move
                if op.A() == ASet {
                    y = __set
                }
                if nodes[0].Type() == Natom && nodes[1].Type() == Natom {
                    $$ = __chain(
                        $2.Cx(0),
                        y(nodes[0], nilNode).pos0($1),
                        y(nodes[1], nilNode).pos0($1),
                        CompNode(ASetFromAB, nodes[0], nodes[1]),
                    )
                } else {
                    $$ = __chain(
                        $2.Cx(0),
                        __set("(2)a", nilNode).pos0($1),
                        __set("(2)b", nilNode).pos0($1),
                        CompNode(ASetFromAB, "(2)a", "(2)b"),
                    )
                    x := func(n *Node, src string) {
                        if n.Cn() > 0 && n.Cx(0).A() == ALoad {
                            $$.Cappend(__store(n.Cx(1), n.Cx(2), src).pos0(n))
                        } else {
                            $$.Cappend(y(n, src).pos0(n))
                        }
                    }
                    x(nodes[0], "(2)a")
                    x(nodes[1], "(2)b")
                }
            } else if op.A() == ASet { // a0, ..., an := b0, ..., bn
                for i, v := range nodes {
                    $$ = $$.Cappend(__set(v, $2.Cx(i)).pos0($1))
                }
            } else if head := nodes[0]; len(nodes) == 1 { // a0 = b0
                $$ = __move(head, $2.Cx(0)).pos0($1)
                if head.Cn() > 0 && head.Cx(0).A() == ALoad {
                    $$ = __store(head.Cx(1), head.Cx(2), $2.Cx(0)).pos0($1)
                }
                if a, s := $2.Cx(0).isSimpleAddSub(); a != "" && a == head.A() { // Note that a := a + v is different
                    $$ = __inc(head, NewNumberNode(s)).pos0($1)
                }
            } else { // a0, ..., an = b0, ..., bn
                for i := range nodes {
                    $$.Cappend(__set("(1)" + strconv.Itoa(i), $2.Cx(i)).pos0($1))
                }
                for i, v := range nodes {
                    if v.Cn() > 0 && v.Cx(0).A() == ALoad {
                        $$.Cappend(__store(v.Cx(1), v.Cx(2), "(1)" + strconv.Itoa(i)).pos0($1))
                    } else {
                        $$.Cappend(__move(v, "(1)" + strconv.Itoa(i)).pos0($1))
                    }
                }
            }
        } 

postfix_incdec:
        TIdent _postfix_incdec {
            $$ = __inc(ANode($1), $2).pos0($1)
        } |
        TIdent _postfix_assign expr %prec ASSIGN  {
            $$ = __move(ANode($1), CompNode($2, ANode($1).setPos($1), $3)).pos0($1)
        } |
        prefix_expr '[' expr ']' _postfix_incdec  {
            $$ = __store($1, $3, CompNode(AAdd, __load($1, $3).pos0($1), $5).pos0($1))
        } |
        prefix_expr '.' TIdent   _postfix_incdec  {
            $$ = __store($1, __hash($3.Str), CompNode(AAdd, __load($1, __hash($3.Str)).pos0($1), $4).pos0($1)) 
        } |
        prefix_expr '[' expr ']' _postfix_assign expr %prec ASSIGN {
            $$ = __store($1, $3, CompNode($5, __load($1, $3).pos0($1), $6).pos0($1))
        } |
        prefix_expr '.' TIdent _postfix_assign expr %prec ASSIGN {
            $$ = __store($1, __hash($3.Str), CompNode($4, __load($1, __hash($3.Str)).pos0($1), $5).pos0($1))
        }

for_stat:
        TFor expr oneline_or_block {
            $$ = __for($2).__continue(emptyNode).__body($3).pos0($1)
        } |
        TFor ';' expr ';' oneline_or_block oneline_or_block {
            $$ = __for($3).__continue($5).__body($6).pos0($1)
        } |
        TFor expr ';' expr ';' oneline_or_block oneline_or_block {
            $$ = __chain(
                $2,
                __for($4).__continue($6).__body($7).pos0($1),
            )
        } |
        TFor TIdent '=' expr oneline_or_block {
            forVar, forEnd := ANode($2), ANodeS($2.Str + "_end")
            $$ = __chain(
                __move(forVar, NewNumberNode(0)).pos0($1),
                __move(forEnd, CompNode(ALen, $4).pos0($1)).pos0($1),
                __for(
                    CompNode(ALess, forVar, forEnd).pos0($1),
                ).
                __continue(
                    __chain(__inc(forVar, oneNode).pos0($1)),
                ).
                __body($5).pos0($1),
            )
        } |
        TFor TIdent '=' expr ',' expr oneline_or_block {
            forVar, forEnd := ANode($2), ANodeS($2.Str + "_end")
            $$ = __chain(
                __move(forVar, $4).pos0($1),
                __move(forEnd, $6).pos0($1),
                __for(
                    CompNode(ALess, forVar, forEnd).pos0($1),
                ).
                __continue(
                    __chain(__inc(forVar, oneNode).pos0($1)),
                ).
                __body($7).pos0($1),
            )
        } |
        TFor TIdent '=' expr ',' expr ',' expr oneline_or_block {
            forVar, forEnd := ANode($2), ANodeS($2.Str + "_end") 
            if $8.Type() == Nnumber { // easy case
                var cond *Node
                if $8.N() < 0 {
                    cond = __lessEq(forEnd, forVar)
                } else {
                    cond = __lessEq(forVar, forEnd)
                }
                $$ = __chain(
                    __move(forVar, $4).pos0($1),
                    __move(forEnd, $6).pos0($1),
                    __for(cond).
                    __continue(__chain(__inc(forVar, $8).pos0($1))).
                    __body($9).pos0($1),
                )
            } else {
                forStep := ANodeS($2.Str + "_step")
                forBegin := ANodeS($2.Str + "_begin")
                $$ = __chain(
                    __move(forVar, $4).pos0($1),
                    __move(forBegin, $4).pos0($1),
                    __move(forEnd, $6).pos0($1),
                    __move(forStep, $8).pos0($1),
                    __if(
                        __lessEq(
                            zeroNode,
                            __mul(
                                __sub(forEnd, forVar).pos0($1),
                                forStep,
                            ).pos0($1),
                        ).pos0($1),
                    ).
                    __then(
                        __chain(
                            __for(
                                __lessEq(
                                    __mul(
                                        __sub(forVar, forBegin).pos0($1), 
                                        __sub(forVar, forEnd).pos0($1),
                                    ),
                                    zeroNode,
                                ).pos0($1), // (forVar - forBegin) * (forVar - forEnd) <= 0
                            ).
                            __continue(
                                __chain(__inc(forVar, forStep).pos0($1)),
                            ).
                            __body($9).pos0($1),
                        ),
                    ).
                    __else(
                        emptyNode,
                    ).pos0($1),
                )
            }
            
        } 

if_stat:
        TIf expr oneline_or_block %prec 'T' {
            $$ = __if($2).__then($3).__else(emptyNode).pos0($1)
        } |
        TIf expr oneline_or_block TElse oneline_or_block {
            $$ = __if($2).__then($3).__else($5).pos0($1)
        }

func_stat:
        TFunc TIdent func_params_list oneline_or_block {
            funcname := ANode($2)
            $$ = __chain(
                __set(funcname, nilNode).pos0($2), 
                __move(funcname, __func(funcname).__params($3).__body($4).pos0($2)).pos0($2),
            )
        }

function:
        TFunc func_params_list block %prec FUNC {
            $$ = __func("<a>").__params($2).__body($3).pos0($1).SetPos($1) 
        } |
        TFunc ident_list '=' expr %prec FUNC {
            $$ = __func("<a>").__params($2).__body(__chain(__return($4).pos0($1))).pos0($1).SetPos($1)
        } |
        TFunc '=' expr %prec FUNC {
            $$ = __func("<a>").__params(emptyNode).__body(__chain(__return($3).pos0($1))).pos0($1).SetPos($1)
        }

func_params_list:
        '(' ')'                           { $$ = emptyNode } |
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
        prefix_expr '.' TIdent            { $$ = __load($1, __hash($3.Str)).pos0($3).setPos($3) } |
        prefix_expr '[' expr ':' expr ']' { $$ = CompNode(ASlice, $1, $3, $5).pos0($3).setPos($3) } |
        prefix_expr '[' expr ':' ']'      { $$ = CompNode(ASlice, $1, $3, moneNode).pos0($3).setPos($3) } |
        prefix_expr '[' ':' expr ']'      { $$ = CompNode(ASlice, $1, zeroNode, $4).pos0($4).setPos($4) }

declarator_list_assign:
        declarator TSet                       { $$ = CompNode($1, ASet) } |
        declarator '='                        { $$ = CompNode($1, AMove) } |
        declarator ',' declarator_list_assign { $$ = $3.Cprepend($1) }

ident_list:
        TIdent                            { $$ = CompNode($1.Str) } | 
        TIdent TDotDotDot                 { $$ = CompNode($1.Str + "...") } | 
        ident_list ',' TIdent             { $$ = $1.Cappend(ANode($3)) } |
        ident_list ',' TIdent TDotDotDot  { $$ = $1.Cappend(ANodeS($3.Str + "...").SetPos($3)) }

expr:
        TNumber                           { $$ = NewNumberNode($1.Str).SetPos($1) } |
        TImport TString                   { $$ = yylex.(*Lexer).loadFile(joinSourcePath($1.Pos.Source, $2.Str), $1) } |
        TTypeof expr                      { $$ = CompNode(ATypeOf, $2) } |
        TLen expr                         { $$ = CompNode(ALen, $2) } |
        function                          { $$ = $1 } |
        struct_slice_gen                  { $$ = $1 } |
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
        expr TDotDotDot                   { $$ = CompNode(ADDD, $1).pos0($1) } |
        '^' expr %prec UNARY              { $$ = CompNode(ABitXor, $2, max32Node).pos0($2) } |
        '-' expr %prec UNARY              { $$ = CompNode(ASub, zeroNode, $2).pos0($2) } |
        '!' expr %prec UNARY              { $$ = CompNode(ANot, $2).pos0($2) } |
        '&' TIdent %prec UNARY            { $$ = CompNode(AAddrOf, ANode($2)).pos0($2) }

prefix_expr:
        declarator                        { $$ = $1 } |
        prefix_expr '(' ')'               { $$ = __call($1, emptyNode).pos0($1) } |
        prefix_expr '(' expr_list ')'     { $$ = patchVarargCall($1, $3) } |
        '(' expr ')'                      { $$ = $2 } // shift/reduce conflict

expr_list:
        expr                              { $$ = CompNode($1) } |
        expr_list ',' expr                { $$ = $1.Cappend($3) }

expr_assign_list:
        TIdent ':' expr                   { $$ = CompNode(NewNode($1.Str), __hash($1.Str), $3) } |
        expr_assign_list ',' TIdent':'expr{ $$ = $1.Cappend(NewNode($3.Str), __hash($3.Str), $5) }

struct_slice_gen:
        '{' '}'                           { $$ = CompNode(AArray, emptyNode).pos0($1) } |
        '{' expr_assign_list     '}'      { $$ = patchStruct(false, $2).pos0($2) } |
        '{' expr_assign_list ',' '}'      { $$ = patchStruct(false, $2).pos0($2) } |
        TStruct'{'expr_assign_list     '}'{ $$ = patchStruct(true, $3).pos0($2) } |
        TStruct'{'expr_assign_list ',' '}'{ $$ = patchStruct(true, $3).pos0($2) } |
        '{' expr_list            '}'      { $$ = CompNode(AArray, $2).pos0($2) } |
        '{' expr_list ','        '}'      { $$ = CompNode(AArray, $2).pos0($2) }

%%

var FieldsField = HashString("__fields")

func patchStruct(named bool, c *Node) *Node {
    x := c.C()
    names, args := CompNode(), CompNode()
    for i := 0; i < len(x); i += 3 {
        args.Cappend(x[i + 1], x[i + 2])
        names.Cappend(x[i])
    }
    if !named {
        return CompNode(AMap, args)
    }
    args.Cappend(NewNumberNode(FieldsField), CompNode(AArray, names).pos0(args))
    return CompNode(AMap, args)
}

func patchVarargCall(callee interface{}, args *Node) *Node {
    ddd := false
    for _, a := range args.C() {
        if a.Type() == Ncompound && a.Cx(0).A() == ADDD {
            ddd = true
            break
        }
    }
    if !ddd {
        return __call(callee, args).pos0(callee)
    }

    if args.Cn() == 1 {
        return __chain(
            CompNode(ADDD, args.Cx(0).Cx(1)).pos0(callee),
            __call(callee, emptyNode).pos0(callee),
        )
    }

    varname := "...vararg"
    res := __chain(__set(varname, CompNode(AArray, emptyNode).pos0(callee)).pos0(callee))
    for _, a := range args.C() {
        if a.Type() == Ncompound && a.Cx(0).A() == ADDD {
            res.Cappend(__move(varname, CompNode(ABitLsh, varname, a.Cx(1)).pos0(callee)).pos0(callee))
        } else {
            res.Cappend(__store(varname, CompNode(ALen, varname).pos0(callee), a).pos0(callee))
        }
    }
    res.Cappend(CompNode(ADDD, varname).pos0(callee))
    res.Cappend(__call(callee, emptyNode).pos0(callee))
    return res
}
