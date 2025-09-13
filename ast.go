package gophpparser

import (
	"encoding/json"
)

type Node interface {
	String() string
	TokenLiteral() string
	Type() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement `json:"statements"`
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	out := ""
	for _, s := range p.Statements {
		out += s.String()
	}
	return out
}

func (p *Program) Type() string {
	return "Program"
}

type Identifier struct {
	Token Token  `json:"token"`
	Value string `json:"value"`
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }
func (i *Identifier) Type() string         { return "Identifier" }

type Variable struct {
	Token Token  `json:"token"`
	Name  string `json:"name"`
}

func (v *Variable) expressionNode()      {}
func (v *Variable) TokenLiteral() string { return v.Token.Literal }
func (v *Variable) String() string       { return "$" + v.Name }
func (v *Variable) Type() string         { return "Variable" }

type IntegerLiteral struct {
	Token Token `json:"token"`
	Value int64 `json:"value"`
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }
func (il *IntegerLiteral) Type() string         { return "IntegerLiteral" }

type FloatLiteral struct {
	Token Token   `json:"token"`
	Value float64 `json:"value"`
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }
func (fl *FloatLiteral) Type() string         { return "FloatLiteral" }

type StringLiteral struct {
	Token Token  `json:"token"`
	Value string `json:"value"`
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }
func (sl *StringLiteral) Type() string         { return "StringLiteral" }

type BooleanLiteral struct {
	Token Token `json:"token"`
	Value bool  `json:"value"`
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }
func (bl *BooleanLiteral) Type() string         { return "BooleanLiteral" }

type MagicConstant struct {
	Token Token  `json:"token"`
	Value string `json:"value"`
}

func (mc *MagicConstant) expressionNode()      {}
func (mc *MagicConstant) TokenLiteral() string { return mc.Token.Literal }
func (mc *MagicConstant) String() string       { return mc.Token.Literal }
func (mc *MagicConstant) Type() string         { return "MagicConstant" }

type ExpressionStatement struct {
	Token      Token      `json:"token"`
	Expression Expression `json:"expression"`
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
func (es *ExpressionStatement) Type() string { return "ExpressionStatement" }

type AssignmentExpression struct {
	Token Token      `json:"token"`
	Name  *Variable  `json:"name"`
	Value Expression `json:"value"`
}

func (ae *AssignmentExpression) expressionNode()      {}
func (ae *AssignmentExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignmentExpression) String() string {
	return ae.Name.String() + " = " + ae.Value.String()
}
func (ae *AssignmentExpression) Type() string { return "AssignmentExpression" }

type InfixExpression struct {
	Token    Token      `json:"token"`
	Left     Expression `json:"left"`
	Operator string     `json:"operator"`
	Right    Expression `json:"right"`
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}
func (ie *InfixExpression) Type() string { return "InfixExpression" }

type PrefixExpression struct {
	Token    Token      `json:"token"`
	Operator string     `json:"operator"`
	Right    Expression `json:"right"`
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	return "(" + pe.Operator + pe.Right.String() + ")"
}
func (pe *PrefixExpression) Type() string { return "PrefixExpression" }

type FunctionDeclaration struct {
	Token      Token           `json:"token"`
	Name       *Identifier     `json:"name"`
	Parameters []*Variable     `json:"parameters"`
	Body       *BlockStatement `json:"body"`
}

func (fd *FunctionDeclaration) statementNode()       {}
func (fd *FunctionDeclaration) TokenLiteral() string { return fd.Token.Literal }
func (fd *FunctionDeclaration) String() string {
	params := ""
	for i, p := range fd.Parameters {
		if i > 0 {
			params += ", "
		}
		params += p.String()
	}
	return "function " + fd.Name.String() + "(" + params + ") " + fd.Body.String()
}
func (fd *FunctionDeclaration) Type() string { return "FunctionDeclaration" }

type ReturnStatement struct {
	Token       Token      `json:"token"`
	ReturnValue Expression `json:"return_value"`
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	if rs.ReturnValue != nil {
		return rs.TokenLiteral() + " " + rs.ReturnValue.String() + ";"
	}
	return rs.TokenLiteral() + ";"
}
func (rs *ReturnStatement) Type() string { return "ReturnStatement" }

type BlockStatement struct {
	Token      Token       `json:"token"`
	Statements []Statement `json:"statements"`
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	out := "{"
	for _, s := range bs.Statements {
		out += s.String()
	}
	out += "}"
	return out
}
func (bs *BlockStatement) Type() string { return "BlockStatement" }

type IfStatement struct {
	Token       Token           `json:"token"`
	Condition   Expression      `json:"condition"`
	Consequence *BlockStatement `json:"consequence"`
	Alternative *BlockStatement `json:"alternative"`
}

func (ifs *IfStatement) statementNode()       {}
func (ifs *IfStatement) TokenLiteral() string { return ifs.Token.Literal }
func (ifs *IfStatement) String() string {
	out := "if" + ifs.Condition.String() + " " + ifs.Consequence.String()
	if ifs.Alternative != nil {
		out += "else " + ifs.Alternative.String()
	}
	return out
}
func (ifs *IfStatement) Type() string { return "IfStatement" }

type EchoStatement struct {
	Token  Token        `json:"token"`
	Values []Expression `json:"values"`
}

func (es *EchoStatement) statementNode()       {}
func (es *EchoStatement) TokenLiteral() string { return es.Token.Literal }
func (es *EchoStatement) String() string {
	out := "echo "
	for i, v := range es.Values {
		if i > 0 {
			out += ", "
		}
		out += v.String()
	}
	return out + ";"
}
func (es *EchoStatement) Type() string { return "EchoStatement" }

type CallExpression struct {
	Token     Token        `json:"token"`
	Function  Expression   `json:"function"`
	Arguments []Expression `json:"arguments"`
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	args := ""
	for i, a := range ce.Arguments {
		if i > 0 {
			args += ", "
		}
		args += a.String()
	}
	return ce.Function.String() + "(" + args + ")"
}
func (ce *CallExpression) Type() string { return "CallExpression" }

type ArrayLiteral struct {
	Token    Token        `json:"token"`
	Elements []Expression `json:"elements"`
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	elements := ""
	for i, e := range al.Elements {
		if i > 0 {
			elements += ", "
		}
		elements += e.String()
	}
	return "[" + elements + "]"
}
func (al *ArrayLiteral) Type() string { return "ArrayLiteral" }

type ForStatement struct {
	Token     Token           `json:"token"`
	Init      Expression      `json:"init"`
	Condition Expression      `json:"condition"`
	Update    Expression      `json:"update"`
	Body      *BlockStatement `json:"body"`
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	return "for (" + fs.Init.String() + "; " + fs.Condition.String() + "; " + fs.Update.String() + ") " + fs.Body.String()
}
func (fs *ForStatement) Type() string { return "ForStatement" }

type IndexExpression struct {
	Token Token      `json:"token"`
	Left  Expression `json:"left"`
	Index Expression `json:"index"`
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	return "(" + ie.Left.String() + "[" + ie.Index.String() + "])"
}
func (ie *IndexExpression) Type() string { return "IndexExpression" }

type PostfixExpression struct {
	Token    Token      `json:"token"`
	Left     Expression `json:"left"`
	Operator string     `json:"operator"`
}

func (pe *PostfixExpression) expressionNode()      {}
func (pe *PostfixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PostfixExpression) String() string {
	return "(" + pe.Left.String() + pe.Operator + ")"
}
func (pe *PostfixExpression) Type() string { return "PostfixExpression" }

type WhileStatement struct {
	Token     Token           `json:"token"`
	Condition Expression      `json:"condition"`
	Body      *BlockStatement `json:"body"`
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	return "while (" + ws.Condition.String() + ") " + ws.Body.String()
}
func (ws *WhileStatement) Type() string { return "WhileStatement" }

type ForeachStatement struct {
	Token Token           `json:"token"`
	Array Expression      `json:"array"`
	Key   *Variable       `json:"key"`
	Value *Variable       `json:"value"`
	Body  *BlockStatement `json:"body"`
}

func (fs *ForeachStatement) statementNode()       {}
func (fs *ForeachStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForeachStatement) String() string {
	out := "foreach (" + fs.Array.String() + " as "
	if fs.Key != nil {
		out += fs.Key.String() + " => "
	}
	out += fs.Value.String() + ") " + fs.Body.String()
	return out
}
func (fs *ForeachStatement) Type() string { return "ForeachStatement" }

type BreakStatement struct {
	Token Token      `json:"token"`
	Level Expression `json:"level,omitempty"`
}

func (bs *BreakStatement) statementNode()       {}
func (bs *BreakStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BreakStatement) String() string {
	if bs.Level != nil {
		return "break " + bs.Level.String() + ";"
	}
	return "break;"
}
func (bs *BreakStatement) Type() string { return "BreakStatement" }

type ContinueStatement struct {
	Token Token      `json:"token"`
	Level Expression `json:"level,omitempty"`
}

func (cs *ContinueStatement) statementNode()       {}
func (cs *ContinueStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ContinueStatement) String() string {
	if cs.Level != nil {
		return "continue " + cs.Level.String() + ";"
	}
	return "continue;"
}
func (cs *ContinueStatement) Type() string { return "ContinueStatement" }

type AssociativeArrayLiteral struct {
	Token Token       `json:"token"`
	Pairs []ArrayPair `json:"pairs"`
}

type ArrayPair struct {
	Key   Expression `json:"key"`
	Value Expression `json:"value"`
}

func (aal *AssociativeArrayLiteral) expressionNode()      {}
func (aal *AssociativeArrayLiteral) TokenLiteral() string { return aal.Token.Literal }
func (aal *AssociativeArrayLiteral) String() string {
	pairs := ""
	for i, pair := range aal.Pairs {
		if i > 0 {
			pairs += ", "
		}
		pairs += pair.Key.String() + " => " + pair.Value.String()
	}
	return "[" + pairs + "]"
}
func (aal *AssociativeArrayLiteral) Type() string { return "AssociativeArrayLiteral" }

type InterpolatedString struct {
	Token Token        `json:"token"`
	Parts []Expression `json:"parts"`
}

func (is *InterpolatedString) expressionNode()      {}
func (is *InterpolatedString) TokenLiteral() string { return is.Token.Literal }
func (is *InterpolatedString) String() string {
	out := "\""
	for _, part := range is.Parts {
		out += part.String()
	}
	out += "\""
	return out
}
func (is *InterpolatedString) Type() string { return "InterpolatedString" }

type ClassDeclaration struct {
	Token      Token                  `json:"token"`
	Name       *Identifier            `json:"name"`
	SuperClass *Identifier            `json:"super_class,omitempty"`
	Interfaces []*Identifier          `json:"interfaces,omitempty"`
	TraitUses  []*TraitUse            `json:"trait_uses,omitempty"`
	Properties []*PropertyDeclaration `json:"properties"`
	Methods    []*MethodDeclaration   `json:"methods"`
	Constants  []*ConstantDeclaration `json:"constants,omitempty"`
}

func (cd *ClassDeclaration) statementNode()       {}
func (cd *ClassDeclaration) TokenLiteral() string { return cd.Token.Literal }
func (cd *ClassDeclaration) String() string {
	out := "class " + cd.Name.String()
	if cd.SuperClass != nil {
		out += " extends " + cd.SuperClass.String()
	}
	if len(cd.Interfaces) > 0 {
		out += " implements "
		for i, iface := range cd.Interfaces {
			if i > 0 {
				out += ", "
			}
			out += iface.String()
		}
	}
	out += " {"
	for _, traitUse := range cd.TraitUses {
		out += traitUse.String()
	}
	for _, constant := range cd.Constants {
		out += constant.String()
	}
	for _, prop := range cd.Properties {
		out += prop.String()
	}
	for _, method := range cd.Methods {
		out += method.String()
	}
	out += "}"
	return out
}
func (cd *ClassDeclaration) Type() string { return "ClassDeclaration" }

type PropertyDeclaration struct {
	Token      Token      `json:"token"`
	Visibility string     `json:"visibility"`
	Static     bool       `json:"static"`
	Name       *Variable  `json:"name"`
	Value      Expression `json:"value,omitempty"`
}

func (pd *PropertyDeclaration) statementNode()       {}
func (pd *PropertyDeclaration) TokenLiteral() string { return pd.Token.Literal }
func (pd *PropertyDeclaration) String() string {
	out := pd.Visibility
	if pd.Static {
		out += " static"
	}
	out += " " + pd.Name.String()
	if pd.Value != nil {
		out += " = " + pd.Value.String()
	}
	out += ";"
	return out
}
func (pd *PropertyDeclaration) Type() string { return "PropertyDeclaration" }

type MethodDeclaration struct {
	Token      Token           `json:"token"`
	Visibility string          `json:"visibility"`
	Static     bool            `json:"static"`
	Name       *Identifier     `json:"name"`
	Parameters []*Variable     `json:"parameters"`
	Body       *BlockStatement `json:"body"`
}

func (md *MethodDeclaration) statementNode()       {}
func (md *MethodDeclaration) TokenLiteral() string { return md.Token.Literal }
func (md *MethodDeclaration) String() string {
	out := md.Visibility
	if md.Static {
		out += " static"
	}
	out += " function " + md.Name.String() + "("
	params := ""
	for i, p := range md.Parameters {
		if i > 0 {
			params += ", "
		}
		params += p.String()
	}
	out += params + ") " + md.Body.String()
	return out
}
func (md *MethodDeclaration) Type() string { return "MethodDeclaration" }

type InterfaceDeclaration struct {
	Token   Token              `json:"token"`
	Name    *Identifier        `json:"name"`
	Methods []*InterfaceMethod `json:"methods"`
}

func (id *InterfaceDeclaration) statementNode()       {}
func (id *InterfaceDeclaration) TokenLiteral() string { return id.Token.Literal }
func (id *InterfaceDeclaration) String() string {
	out := "interface " + id.Name.String() + " {"
	for _, method := range id.Methods {
		out += method.String()
	}
	out += "}"
	return out
}
func (id *InterfaceDeclaration) Type() string { return "InterfaceDeclaration" }

type InterfaceMethod struct {
	Token      Token       `json:"token"`
	Visibility string      `json:"visibility"`
	Name       *Identifier `json:"name"`
	Parameters []*Variable `json:"parameters"`
}

func (im *InterfaceMethod) statementNode()       {}
func (im *InterfaceMethod) TokenLiteral() string { return im.Token.Literal }
func (im *InterfaceMethod) String() string {
	params := ""
	for i, p := range im.Parameters {
		if i > 0 {
			params += ", "
		}
		params += p.String()
	}
	return im.Visibility + " function " + im.Name.String() + "(" + params + ");"
}
func (im *InterfaceMethod) Type() string { return "InterfaceMethod" }

type TraitDeclaration struct {
	Token      Token                  `json:"token"`
	Name       *Identifier            `json:"name"`
	Properties []*PropertyDeclaration `json:"properties"`
	Methods    []*MethodDeclaration   `json:"methods"`
}

func (td *TraitDeclaration) statementNode()       {}
func (td *TraitDeclaration) TokenLiteral() string { return td.Token.Literal }
func (td *TraitDeclaration) String() string {
	out := "trait " + td.Name.String() + " {"
	for _, prop := range td.Properties {
		out += prop.String()
	}
	for _, method := range td.Methods {
		out += method.String()
	}
	out += "}"
	return out
}
func (td *TraitDeclaration) Type() string { return "TraitDeclaration" }

type TraitUse struct {
	Token  Token         `json:"token"`
	Traits []*Identifier `json:"traits"`
}

func (tu *TraitUse) statementNode()       {}
func (tu *TraitUse) TokenLiteral() string { return tu.Token.Literal }
func (tu *TraitUse) String() string {
	traits := ""
	for i, trait := range tu.Traits {
		if i > 0 {
			traits += ", "
		}
		traits += trait.String()
	}
	return "use " + traits + ";"
}
func (tu *TraitUse) Type() string { return "TraitUse" }

type ConstantDeclaration struct {
	Token      Token       `json:"token"`
	Visibility string      `json:"visibility"`
	Name       *Identifier `json:"name"`
	Value      Expression  `json:"value"`
}

func (cd *ConstantDeclaration) statementNode()       {}
func (cd *ConstantDeclaration) TokenLiteral() string { return cd.Token.Literal }
func (cd *ConstantDeclaration) String() string {
	out := cd.Visibility + " const " + cd.Name.String() + " = " + cd.Value.String() + ";"
	return out
}
func (cd *ConstantDeclaration) Type() string { return "ConstantDeclaration" }

type NewExpression struct {
	Token     Token        `json:"token"`
	ClassName *Identifier  `json:"class_name"`
	Arguments []Expression `json:"arguments"`
}

func (ne *NewExpression) expressionNode()      {}
func (ne *NewExpression) TokenLiteral() string { return ne.Token.Literal }
func (ne *NewExpression) String() string {
	args := ""
	for i, arg := range ne.Arguments {
		if i > 0 {
			args += ", "
		}
		args += arg.String()
	}
	return "new " + ne.ClassName.String() + "(" + args + ")"
}
func (ne *NewExpression) Type() string { return "NewExpression" }

type ObjectAccessExpression struct {
	Token    Token      `json:"token"`
	Object   Expression `json:"object"`
	Property Expression `json:"property"`
}

func (oae *ObjectAccessExpression) expressionNode()      {}
func (oae *ObjectAccessExpression) TokenLiteral() string { return oae.Token.Literal }
func (oae *ObjectAccessExpression) String() string {
	return oae.Object.String() + "->" + oae.Property.String()
}
func (oae *ObjectAccessExpression) Type() string { return "ObjectAccessExpression" }

type StaticAccessExpression struct {
	Token    Token      `json:"token"`
	Class    Expression `json:"class"`
	Property Expression `json:"property"`
}

func (sae *StaticAccessExpression) expressionNode()      {}
func (sae *StaticAccessExpression) TokenLiteral() string { return sae.Token.Literal }
func (sae *StaticAccessExpression) String() string {
	return sae.Class.String() + "::" + sae.Property.String()
}
func (sae *StaticAccessExpression) Type() string { return "StaticAccessExpression" }

type NamespaceDeclaration struct {
	Token Token       `json:"token"`
	Name  *Identifier `json:"name"`
}

func (nd *NamespaceDeclaration) statementNode()       {}
func (nd *NamespaceDeclaration) TokenLiteral() string { return nd.Token.Literal }
func (nd *NamespaceDeclaration) String() string {
	return "namespace " + nd.Name.String() + ";"
}
func (nd *NamespaceDeclaration) Type() string { return "NamespaceDeclaration" }

type UseStatement struct {
	Token     Token       `json:"token"`
	Namespace *Identifier `json:"namespace"`
	Alias     *Identifier `json:"alias,omitempty"`
}

func (us *UseStatement) statementNode()       {}
func (us *UseStatement) TokenLiteral() string { return us.Token.Literal }
func (us *UseStatement) String() string {
	out := "use " + us.Namespace.String()
	if us.Alias != nil {
		out += " as " + us.Alias.String()
	}
	out += ";"
	return out
}
func (us *UseStatement) Type() string { return "UseStatement" }

type TryStatement struct {
	Token   Token           `json:"token"`
	Body    *BlockStatement `json:"body"`
	Catches []*CatchClause  `json:"catches"`
	Finally *BlockStatement `json:"finally,omitempty"`
}

func (ts *TryStatement) statementNode()       {}
func (ts *TryStatement) TokenLiteral() string { return ts.Token.Literal }
func (ts *TryStatement) String() string {
	out := "try " + ts.Body.String()
	for _, catch := range ts.Catches {
		out += catch.String()
	}
	if ts.Finally != nil {
		out += " finally " + ts.Finally.String()
	}
	return out
}
func (ts *TryStatement) Type() string { return "TryStatement" }

type CatchClause struct {
	Token         Token           `json:"token"`
	ExceptionType *Identifier     `json:"exception_type"`
	Variable      *Variable       `json:"variable"`
	Body          *BlockStatement `json:"body"`
}

func (cc *CatchClause) statementNode()       {}
func (cc *CatchClause) TokenLiteral() string { return cc.Token.Literal }
func (cc *CatchClause) String() string {
	out := " catch ("
	if cc.ExceptionType != nil {
		out += cc.ExceptionType.String() + " "
	}
	out += cc.Variable.String() + ") " + cc.Body.String()
	return out
}
func (cc *CatchClause) Type() string { return "CatchClause" }

type ThrowStatement struct {
	Token      Token      `json:"token"`
	Expression Expression `json:"expression"`
}

func (ts *ThrowStatement) statementNode()       {}
func (ts *ThrowStatement) TokenLiteral() string { return ts.Token.Literal }
func (ts *ThrowStatement) String() string {
	return "throw " + ts.Expression.String() + ";"
}
func (ts *ThrowStatement) Type() string { return "ThrowStatement" }

type AnonymousFunction struct {
	Token      Token           `json:"token"`
	Parameters []*Variable     `json:"parameters"`
	UseClause  []*Variable     `json:"use_clause,omitempty"`
	Body       *BlockStatement `json:"body"`
}

func (af *AnonymousFunction) expressionNode()      {}
func (af *AnonymousFunction) TokenLiteral() string { return af.Token.Literal }
func (af *AnonymousFunction) String() string {
	params := ""
	for i, p := range af.Parameters {
		if i > 0 {
			params += ", "
		}
		params += p.String()
	}
	out := "function(" + params + ")"

	if len(af.UseClause) > 0 {
		uses := ""
		for i, u := range af.UseClause {
			if i > 0 {
				uses += ", "
			}
			uses += u.String()
		}
		out += " use (" + uses + ")"
	}

	out += " " + af.Body.String()
	return out
}
func (af *AnonymousFunction) Type() string { return "AnonymousFunction" }

type NamespacedIdentifier struct {
	Token     Token         `json:"token"`
	Namespace []*Identifier `json:"namespace"`
	Name      *Identifier   `json:"name"`
}

func (ni *NamespacedIdentifier) expressionNode()      {}
func (ni *NamespacedIdentifier) TokenLiteral() string { return ni.Token.Literal }
func (ni *NamespacedIdentifier) String() string {
	out := ""
	for i, ns := range ni.Namespace {
		if i > 0 || len(ni.Namespace) > 0 {
			out += "\\"
		}
		out += ns.String()
	}
	if len(ni.Namespace) > 0 {
		out += "\\"
	}
	out += ni.Name.String()
	return out
}
func (ni *NamespacedIdentifier) Type() string { return "NamespacedIdentifier" }

type YieldExpression struct {
	Token Token      `json:"token"`
	Key   Expression `json:"key,omitempty"`
	Value Expression `json:"value,omitempty"`
}

func (ye *YieldExpression) expressionNode()      {}
func (ye *YieldExpression) TokenLiteral() string { return ye.Token.Literal }
func (ye *YieldExpression) String() string {
	if ye.Key != nil && ye.Value != nil {
		return "yield " + ye.Key.String() + " => " + ye.Value.String()
	} else if ye.Value != nil {
		return "yield " + ye.Value.String()
	}
	return "yield"
}
func (ye *YieldExpression) Type() string { return "YieldExpression" }

type TernaryExpression struct {
	Token      Token      `json:"token"`
	Condition  Expression `json:"condition"`
	TrueValue  Expression `json:"true_value"`
	FalseValue Expression `json:"false_value"`
}

func (te *TernaryExpression) expressionNode()      {}
func (te *TernaryExpression) TokenLiteral() string { return te.Token.Literal }
func (te *TernaryExpression) String() string {
	return "(" + te.Condition.String() + " ? " + te.TrueValue.String() + " : " + te.FalseValue.String() + ")"
}
func (te *TernaryExpression) Type() string { return "TernaryExpression" }

type DeclareStatement struct {
	Token      Token                    `json:"token"`
	Directives map[string]Expression    `json:"directives"`
	Body       *BlockStatement          `json:"body,omitempty"`
}

func (ds *DeclareStatement) statementNode()       {}
func (ds *DeclareStatement) TokenLiteral() string { return ds.Token.Literal }
func (ds *DeclareStatement) String() string {
	out := "declare("
	first := true
	for key, value := range ds.Directives {
		if !first {
			out += ", "
		}
		out += key + "=" + value.String()
		first = false
	}
	out += ")"
	if ds.Body != nil {
		out += " " + ds.Body.String()
	} else {
		out += ";"
	}
	return out
}
func (ds *DeclareStatement) Type() string { return "DeclareStatement" }

func ToJSON(node Node) ([]byte, error) {
	data := map[string]any{
		"type": node.Type(),
	}

	switch n := node.(type) {
	case *Program:
		data["statements"] = n.Statements
	case *Identifier:
		data["value"] = n.Value
	case *Variable:
		data["name"] = n.Name
	case *IntegerLiteral:
		data["value"] = n.Value
	case *FloatLiteral:
		data["value"] = n.Value
	case *StringLiteral:
		data["value"] = n.Value
	case *BooleanLiteral:
		data["value"] = n.Value
	case *ExpressionStatement:
		data["expression"] = n.Expression
	case *AssignmentExpression:
		data["name"] = n.Name
		data["value"] = n.Value
	case *InfixExpression:
		data["left"] = n.Left
		data["operator"] = n.Operator
		data["right"] = n.Right
	case *PrefixExpression:
		data["operator"] = n.Operator
		data["right"] = n.Right
	case *FunctionDeclaration:
		data["name"] = n.Name
		data["parameters"] = n.Parameters
		data["body"] = n.Body
	case *ReturnStatement:
		data["return_value"] = n.ReturnValue
	case *BlockStatement:
		data["statements"] = n.Statements
	case *IfStatement:
		data["condition"] = n.Condition
		data["consequence"] = n.Consequence
		if n.Alternative != nil {
			data["alternative"] = n.Alternative
		}
	case *EchoStatement:
		data["values"] = n.Values
	case *CallExpression:
		data["function"] = n.Function
		data["arguments"] = n.Arguments
	case *ArrayLiteral:
		data["elements"] = n.Elements
	case *ForStatement:
		data["init"] = n.Init
		data["condition"] = n.Condition
		data["update"] = n.Update
		data["body"] = n.Body
	case *IndexExpression:
		data["left"] = n.Left
		data["index"] = n.Index
	case *PostfixExpression:
		data["left"] = n.Left
		data["operator"] = n.Operator
	case *WhileStatement:
		data["condition"] = n.Condition
		data["body"] = n.Body
	case *ForeachStatement:
		data["array"] = n.Array
		if n.Key != nil {
			data["key"] = n.Key
		}
		data["value"] = n.Value
		data["body"] = n.Body
	case *BreakStatement:
		if n.Level != nil {
			data["level"] = n.Level
		}
	case *ContinueStatement:
		if n.Level != nil {
			data["level"] = n.Level
		}
	case *AssociativeArrayLiteral:
		data["pairs"] = n.Pairs
	case *InterpolatedString:
		data["parts"] = n.Parts
	case *ClassDeclaration:
		data["name"] = n.Name
		if n.SuperClass != nil {
			data["super_class"] = n.SuperClass
		}
		if len(n.Interfaces) > 0 {
			data["interfaces"] = n.Interfaces
		}
		if len(n.TraitUses) > 0 {
			data["trait_uses"] = n.TraitUses
		}
		data["properties"] = n.Properties
		data["methods"] = n.Methods
		if len(n.Constants) > 0 {
			data["constants"] = n.Constants
		}
	case *PropertyDeclaration:
		data["visibility"] = n.Visibility
		data["static"] = n.Static
		data["name"] = n.Name
		if n.Value != nil {
			data["value"] = n.Value
		}
	case *MethodDeclaration:
		data["visibility"] = n.Visibility
		data["static"] = n.Static
		data["name"] = n.Name
		data["parameters"] = n.Parameters
		data["body"] = n.Body
	case *NewExpression:
		data["class_name"] = n.ClassName
		data["arguments"] = n.Arguments
	case *ObjectAccessExpression:
		data["object"] = n.Object
		data["property"] = n.Property
	case *StaticAccessExpression:
		data["class"] = n.Class
		data["property"] = n.Property
	case *NamespaceDeclaration:
		data["name"] = n.Name
	case *UseStatement:
		data["namespace"] = n.Namespace
		if n.Alias != nil {
			data["alias"] = n.Alias
		}
	case *TryStatement:
		data["body"] = n.Body
		data["catches"] = n.Catches
		if n.Finally != nil {
			data["finally"] = n.Finally
		}
	case *CatchClause:
		if n.ExceptionType != nil {
			data["exception_type"] = n.ExceptionType
		}
		data["variable"] = n.Variable
		data["body"] = n.Body
	case *ThrowStatement:
		data["expression"] = n.Expression
	case *AnonymousFunction:
		data["parameters"] = n.Parameters
		if len(n.UseClause) > 0 {
			data["use_clause"] = n.UseClause
		}
		data["body"] = n.Body
	case *NamespacedIdentifier:
		data["namespace"] = n.Namespace
		data["name"] = n.Name
	case *YieldExpression:
		if n.Key != nil {
			data["key"] = n.Key
		}
		if n.Value != nil {
			data["value"] = n.Value
		}
	case *InterfaceDeclaration:
		data["name"] = n.Name
		data["methods"] = n.Methods
	case *InterfaceMethod:
		data["visibility"] = n.Visibility
		data["name"] = n.Name
		data["parameters"] = n.Parameters
	case *TraitDeclaration:
		data["name"] = n.Name
		data["properties"] = n.Properties
		data["methods"] = n.Methods
	case *TraitUse:
		data["traits"] = n.Traits
	case *ConstantDeclaration:
		data["visibility"] = n.Visibility
		data["name"] = n.Name
		data["value"] = n.Value
	case *TernaryExpression:
		data["condition"] = n.Condition
		data["true_value"] = n.TrueValue
		data["false_value"] = n.FalseValue
	case *DeclareStatement:
		data["directives"] = n.Directives
		if n.Body != nil {
			data["body"] = n.Body
		}
	}

	return json.MarshalIndent(data, "", "  ")
}
