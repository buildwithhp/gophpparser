package gophpparser

import (
	"fmt"
	"strings"
)

// SymbolType represents the type of symbol
type SymbolType int

const (
	CLASS_SYMBOL SymbolType = iota
	FUNCTION_SYMBOL
	VARIABLE_SYMBOL
	CONSTANT_SYMBOL
	INTERFACE_SYMBOL
	TRAIT_SYMBOL
)

func (st SymbolType) String() string {
	switch st {
	case CLASS_SYMBOL:
		return "class"
	case FUNCTION_SYMBOL:
		return "function"
	case VARIABLE_SYMBOL:
		return "variable"
	case CONSTANT_SYMBOL:
		return "constant"
	case INTERFACE_SYMBOL:
		return "interface"
	case TRAIT_SYMBOL:
		return "trait"
	default:
		return "unknown"
	}
}

// Symbol represents a declared symbol with its fully qualified name
type Symbol struct {
	Name         string     `json:"name"`           // Local name (e.g., "User")
	FullyQualified string   `json:"fully_qualified"` // Full name (e.g., "HR\\User")
	Type         SymbolType `json:"type"`           // Symbol type
	Namespace    string     `json:"namespace"`      // Declaring namespace
	File         string     `json:"file,omitempty"` // Source file
	Line         int        `json:"line,omitempty"` // Line number
}

// SymbolReference represents a reference to a symbol with resolved information
type SymbolReference struct {
	Name           string  `json:"name"`             // Used name (e.g., "User")
	ResolvedSymbol *Symbol `json:"resolved_symbol"`  // What it actually refers to
	Line           int     `json:"line,omitempty"`   // Where it's used
	Column         int     `json:"column,omitempty"` // Column position
}

// Scope represents a lexical scope (global, namespace, class, function)
type Scope struct {
	Type      string             `json:"type"`      // "global", "namespace", "class", "function"
	Name      string             `json:"name"`      // Scope identifier
	Parent    *Scope             `json:"-"`         // Parent scope
	Symbols   map[string]*Symbol `json:"symbols"`   // Symbols declared in this scope
	Children  []*Scope           `json:"children"`  // Child scopes
	Namespace string             `json:"namespace"` // Current namespace
	Imports   map[string]string  `json:"imports"`   // use statements (alias -> fully qualified)
}

// SymbolTable manages all symbols and scopes
type SymbolTable struct {
	GlobalScope   *Scope                        `json:"global_scope"`
	CurrentScope  *Scope                        `json:"-"`
	AllSymbols    map[string]*Symbol            `json:"all_symbols"`    // All symbols by fully qualified name
	References    []*SymbolReference            `json:"references"`     // All symbol references
	Namespaces    map[string][]*Symbol          `json:"namespaces"`     // Symbols grouped by namespace
	ClassHierarchy map[string][]string          `json:"class_hierarchy"` // class -> [parent, interfaces...]
}

// NewSymbolTable creates a new symbol table
func NewSymbolTable() *SymbolTable {
	globalScope := &Scope{
		Type:      "global",
		Name:      "global",
		Symbols:   make(map[string]*Symbol),
		Children:  []*Scope{},
		Namespace: "",
		Imports:   make(map[string]string),
	}

	return &SymbolTable{
		GlobalScope:    globalScope,
		CurrentScope:   globalScope,
		AllSymbols:     make(map[string]*Symbol),
		References:     []*SymbolReference{},
		Namespaces:     make(map[string][]*Symbol),
		ClassHierarchy: make(map[string][]string),
	}
}

// EnterScope creates a new child scope
func (st *SymbolTable) EnterScope(scopeType, name string) {
	newScope := &Scope{
		Type:      scopeType,
		Name:      name,
		Parent:    st.CurrentScope,
		Symbols:   make(map[string]*Symbol),
		Children:  []*Scope{},
		Namespace: st.CurrentScope.Namespace, // Inherit namespace
		Imports:   make(map[string]string),   // Copy imports from parent
	}

	// Copy imports from parent
	for alias, fqn := range st.CurrentScope.Imports {
		newScope.Imports[alias] = fqn
	}

	st.CurrentScope.Children = append(st.CurrentScope.Children, newScope)
	st.CurrentScope = newScope
}

// ExitScope returns to parent scope
func (st *SymbolTable) ExitScope() {
	if st.CurrentScope.Parent != nil {
		st.CurrentScope = st.CurrentScope.Parent
	}
}

// SetNamespace sets the current namespace
func (st *SymbolTable) SetNamespace(namespace string) {
	st.CurrentScope.Namespace = namespace
}

// AddImport adds a use statement
func (st *SymbolTable) AddImport(fullyQualified, alias string) {
	if alias == "" {
		// Extract class name from fully qualified name
		parts := strings.Split(fullyQualified, "\\")
		alias = parts[len(parts)-1]
	}
	st.CurrentScope.Imports[alias] = fullyQualified
}

// DeclareSymbol declares a new symbol in current scope
func (st *SymbolTable) DeclareSymbol(name string, symbolType SymbolType, file string, line int) *Symbol {
	// Create fully qualified name
	fqn := st.makeFullyQualified(name)

	symbol := &Symbol{
		Name:           name,
		FullyQualified: fqn,
		Type:           symbolType,
		Namespace:      st.CurrentScope.Namespace,
		File:           file,
		Line:           line,
	}

	// Add to current scope
	st.CurrentScope.Symbols[name] = symbol

	// Add to global registry
	st.AllSymbols[fqn] = symbol

	// Add to namespace registry
	if st.Namespaces[symbol.Namespace] == nil {
		st.Namespaces[symbol.Namespace] = []*Symbol{}
	}
	st.Namespaces[symbol.Namespace] = append(st.Namespaces[symbol.Namespace], symbol)

	return symbol
}

// ResolveSymbol resolves a symbol reference
func (st *SymbolTable) ResolveSymbol(name string, symbolType SymbolType) *Symbol {
	// 1. Check if it's an absolute reference (starts with \)
	if strings.HasPrefix(name, "\\") {
		if symbol, exists := st.AllSymbols[name]; exists && symbol.Type == symbolType {
			return symbol
		}
		return nil
	}

	// 2. Check imports/aliases first
	if fqn, exists := st.CurrentScope.Imports[name]; exists {
		if symbol, exists := st.AllSymbols[fqn]; exists && symbol.Type == symbolType {
			return symbol
		}
	}

	// 3. Check current scope and parent scopes
	scope := st.CurrentScope
	for scope != nil {
		if symbol, exists := scope.Symbols[name]; exists && symbol.Type == symbolType {
			return symbol
		}
		scope = scope.Parent
	}

	// 4. Check current namespace
	currentNamespace := st.CurrentScope.Namespace
	if currentNamespace != "" {
		fqn := currentNamespace + "\\" + name
		if symbol, exists := st.AllSymbols[fqn]; exists && symbol.Type == symbolType {
			return symbol
		}
	}

	// 5. Check global namespace
	if symbol, exists := st.AllSymbols[name]; exists && symbol.Type == symbolType {
		return symbol
	}

	return nil
}

// AddReference adds a symbol reference
func (st *SymbolTable) AddReference(name string, symbolType SymbolType, line, column int) *SymbolReference {
	resolvedSymbol := st.ResolveSymbol(name, symbolType)

	ref := &SymbolReference{
		Name:           name,
		ResolvedSymbol: resolvedSymbol,
		Line:           line,
		Column:         column,
	}

	st.References = append(st.References, ref)
	return ref
}

// AddClassHierarchy adds inheritance information
func (st *SymbolTable) AddClassHierarchy(className string, extends string, implements []string) {
	hierarchy := []string{}
	if extends != "" {
		hierarchy = append(hierarchy, extends)
	}
	hierarchy = append(hierarchy, implements...)
	st.ClassHierarchy[className] = hierarchy
}

// makeFullyQualified creates a fully qualified name
func (st *SymbolTable) makeFullyQualified(name string) string {
	if strings.HasPrefix(name, "\\") {
		return name // Already fully qualified
	}

	if st.CurrentScope.Namespace == "" {
		return name // Global namespace
	}

	return st.CurrentScope.Namespace + "\\" + name
}

// GetClassHierarchy returns the inheritance chain for a class
func (st *SymbolTable) GetClassHierarchy(className string) []string {
	return st.ClassHierarchy[className]
}

// FindSymbolsInNamespace returns all symbols in a given namespace
func (st *SymbolTable) FindSymbolsInNamespace(namespace string) []*Symbol {
	return st.Namespaces[namespace]
}

// GetUnresolvedReferences returns references that couldn't be resolved
func (st *SymbolTable) GetUnresolvedReferences() []*SymbolReference {
	var unresolved []*SymbolReference
	for _, ref := range st.References {
		if ref.ResolvedSymbol == nil {
			unresolved = append(unresolved, ref)
		}
	}
	return unresolved
}

// SemanticAnalyzer performs semantic analysis on AST
type SemanticAnalyzer struct {
	SymbolTable *SymbolTable
	CurrentFile string
	Errors      []string
}

// NewSemanticAnalyzer creates a new semantic analyzer
func NewSemanticAnalyzer() *SemanticAnalyzer {
	return &SemanticAnalyzer{
		SymbolTable: NewSymbolTable(),
		Errors:      []string{},
	}
}

// AnalyzeProgram performs semantic analysis on a program
func (sa *SemanticAnalyzer) AnalyzeProgram(program *Program, filename string) {
	sa.CurrentFile = filename
	sa.visitProgram(program)
}

// visitProgram visits program node
func (sa *SemanticAnalyzer) visitProgram(program *Program) {
	for _, stmt := range program.Statements {
		sa.visitStatement(stmt)
	}
}

// visitStatement visits statement nodes
func (sa *SemanticAnalyzer) visitStatement(stmt Statement) {
	switch s := stmt.(type) {
	case *NamespaceDeclaration:
		sa.visitNamespaceDeclaration(s)
	case *UseStatement:
		sa.visitUseStatement(s)
	case *ClassDeclaration:
		sa.visitClassDeclaration(s)
	case *InterfaceDeclaration:
		sa.visitInterfaceDeclaration(s)
	case *TraitDeclaration:
		sa.visitTraitDeclaration(s)
	case *FunctionDeclaration:
		sa.visitFunctionDeclaration(s)
	case *ExpressionStatement:
		sa.visitExpression(s.Expression)
	case *BlockStatement:
		sa.visitBlockStatement(s)
	case *IfStatement:
		sa.visitIfStatement(s)
	case *ForStatement:
		sa.visitForStatement(s)
	case *WhileStatement:
		sa.visitWhileStatement(s)
	case *ForeachStatement:
		sa.visitForeachStatement(s)
	case *ReturnStatement:
		sa.visitReturnStatement(s)
	case *EchoStatement:
		sa.visitEchoStatement(s)
	case *TryStatement:
		sa.visitTryStatement(s)
	case *ThrowStatement:
		sa.visitThrowStatement(s)
	}
}

// visitExpression visits expression nodes
func (sa *SemanticAnalyzer) visitExpression(expr Expression) {
	switch e := expr.(type) {
	case *NewExpression:
		sa.visitNewExpression(e)
	case *CallExpression:
		sa.visitCallExpression(e)
	case *ObjectAccessExpression:
		sa.visitObjectAccessExpression(e)
	case *StaticAccessExpression:
		sa.visitStaticAccessExpression(e)
	case *AssignmentExpression:
		sa.visitAssignmentExpression(e)
	case *InfixExpression:
		sa.visitInfixExpression(e)
	case *PrefixExpression:
		sa.visitPrefixExpression(e)
	case *PostfixExpression:
		sa.visitPostfixExpression(e)
	case *ArrayLiteral:
		sa.visitArrayLiteral(e)
	case *AssociativeArrayLiteral:
		sa.visitAssociativeArrayLiteral(e)
	case *IndexExpression:
		sa.visitIndexExpression(e)
	case *AnonymousFunction:
		sa.visitAnonymousFunction(e)
	case *YieldExpression:
		sa.visitYieldExpression(e)
	case *TernaryExpression:
		sa.visitTernaryExpression(e)
	case *Identifier:
		// This might be a function call or constant reference
		sa.addIdentifierReference(e)
	}
}

// Specific visit methods for each node type
func (sa *SemanticAnalyzer) visitNamespaceDeclaration(stmt *NamespaceDeclaration) {
	sa.SymbolTable.SetNamespace(stmt.Name.Value)
}

func (sa *SemanticAnalyzer) visitUseStatement(stmt *UseStatement) {
	alias := ""
	if stmt.Alias != nil {
		alias = stmt.Alias.Value
	}
	sa.SymbolTable.AddImport(stmt.Namespace.Value, alias)
}

func (sa *SemanticAnalyzer) visitClassDeclaration(stmt *ClassDeclaration) {
	// Declare the class
	symbol := sa.SymbolTable.DeclareSymbol(stmt.Name.Value, CLASS_SYMBOL, sa.CurrentFile, stmt.Token.Line)

	// Add inheritance information
	extends := ""
	if stmt.SuperClass != nil {
		extends = stmt.SuperClass.Value
	}
	
	implements := []string{}
	for _, iface := range stmt.Interfaces {
		implements = append(implements, iface.Value)
	}
	
	sa.SymbolTable.AddClassHierarchy(symbol.FullyQualified, extends, implements)

	// Enter class scope
	sa.SymbolTable.EnterScope("class", stmt.Name.Value)

	// Visit class members
	for _, constant := range stmt.Constants {
		sa.visitConstantDeclaration(constant)
	}
	for _, property := range stmt.Properties {
		sa.visitPropertyDeclaration(property)
	}
	for _, method := range stmt.Methods {
		sa.visitMethodDeclaration(method)
	}

	// Exit class scope
	sa.SymbolTable.ExitScope()
}

func (sa *SemanticAnalyzer) visitInterfaceDeclaration(stmt *InterfaceDeclaration) {
	sa.SymbolTable.DeclareSymbol(stmt.Name.Value, INTERFACE_SYMBOL, sa.CurrentFile, stmt.Token.Line)

	sa.SymbolTable.EnterScope("interface", stmt.Name.Value)
	for _, method := range stmt.Methods {
		sa.visitInterfaceMethod(method)
	}
	sa.SymbolTable.ExitScope()
}

func (sa *SemanticAnalyzer) visitTraitDeclaration(stmt *TraitDeclaration) {
	sa.SymbolTable.DeclareSymbol(stmt.Name.Value, TRAIT_SYMBOL, sa.CurrentFile, stmt.Token.Line)

	sa.SymbolTable.EnterScope("trait", stmt.Name.Value)
	for _, property := range stmt.Properties {
		sa.visitPropertyDeclaration(property)
	}
	for _, method := range stmt.Methods {
		sa.visitMethodDeclaration(method)
	}
	sa.SymbolTable.ExitScope()
}

func (sa *SemanticAnalyzer) visitFunctionDeclaration(stmt *FunctionDeclaration) {
	sa.SymbolTable.DeclareSymbol(stmt.Name.Value, FUNCTION_SYMBOL, sa.CurrentFile, stmt.Token.Line)

	sa.SymbolTable.EnterScope("function", stmt.Name.Value)
	for _, param := range stmt.Parameters {
		sa.SymbolTable.DeclareSymbol(param.Name, VARIABLE_SYMBOL, sa.CurrentFile, param.Token.Line)
	}
	sa.visitBlockStatement(stmt.Body)
	sa.SymbolTable.ExitScope()
}

func (sa *SemanticAnalyzer) visitNewExpression(expr *NewExpression) {
	// Add reference to the class being instantiated
	_ = sa.SymbolTable.AddReference(expr.ClassName.Value, CLASS_SYMBOL, expr.Token.Line, 0)
	
	// Visit constructor arguments
	for _, arg := range expr.Arguments {
		sa.visitExpression(arg)
	}
}

func (sa *SemanticAnalyzer) visitCallExpression(expr *CallExpression) {
	// If it's a simple function call (Identifier), add reference
	if identifier, ok := expr.Function.(*Identifier); ok {
		sa.SymbolTable.AddReference(identifier.Value, FUNCTION_SYMBOL, expr.Token.Line, 0)
	} else {
		// Visit the function expression (could be method call, etc.)
		sa.visitExpression(expr.Function)
	}

	// Visit arguments
	for _, arg := range expr.Arguments {
		sa.visitExpression(arg)
	}
}

func (sa *SemanticAnalyzer) visitObjectAccessExpression(expr *ObjectAccessExpression) {
	sa.visitExpression(expr.Object)
	sa.visitExpression(expr.Property)
}

func (sa *SemanticAnalyzer) visitStaticAccessExpression(expr *StaticAccessExpression) {
	// Add reference to the class
	if identifier, ok := expr.Class.(*Identifier); ok {
		sa.SymbolTable.AddReference(identifier.Value, CLASS_SYMBOL, expr.Token.Line, 0)
	} else {
		sa.visitExpression(expr.Class)
	}
	sa.visitExpression(expr.Property)
}

func (sa *SemanticAnalyzer) visitAssignmentExpression(expr *AssignmentExpression) {
	// Declare variable if it's new
	sa.SymbolTable.DeclareSymbol(expr.Name.Name, VARIABLE_SYMBOL, sa.CurrentFile, expr.Token.Line)
	sa.visitExpression(expr.Value)
}

// Helper methods
func (sa *SemanticAnalyzer) visitBlockStatement(stmt *BlockStatement) {
	for _, s := range stmt.Statements {
		sa.visitStatement(s)
	}
}

func (sa *SemanticAnalyzer) visitIfStatement(stmt *IfStatement) {
	sa.visitExpression(stmt.Condition)
	sa.visitBlockStatement(stmt.Consequence)
	if stmt.Alternative != nil {
		sa.visitBlockStatement(stmt.Alternative)
	}
}

func (sa *SemanticAnalyzer) visitForStatement(stmt *ForStatement) {
	sa.visitExpression(stmt.Init)
	sa.visitExpression(stmt.Condition)
	sa.visitExpression(stmt.Update)
	sa.visitBlockStatement(stmt.Body)
}

func (sa *SemanticAnalyzer) visitWhileStatement(stmt *WhileStatement) {
	sa.visitExpression(stmt.Condition)
	sa.visitBlockStatement(stmt.Body)
}

func (sa *SemanticAnalyzer) visitForeachStatement(stmt *ForeachStatement) {
	sa.visitExpression(stmt.Array)
	if stmt.Key != nil {
		sa.SymbolTable.DeclareSymbol(stmt.Key.Name, VARIABLE_SYMBOL, sa.CurrentFile, stmt.Token.Line)
	}
	sa.SymbolTable.DeclareSymbol(stmt.Value.Name, VARIABLE_SYMBOL, sa.CurrentFile, stmt.Token.Line)
	sa.visitBlockStatement(stmt.Body)
}

func (sa *SemanticAnalyzer) visitReturnStatement(stmt *ReturnStatement) {
	if stmt.ReturnValue != nil {
		sa.visitExpression(stmt.ReturnValue)
	}
}

func (sa *SemanticAnalyzer) visitEchoStatement(stmt *EchoStatement) {
	for _, value := range stmt.Values {
		sa.visitExpression(value)
	}
}

func (sa *SemanticAnalyzer) visitTryStatement(stmt *TryStatement) {
	sa.visitBlockStatement(stmt.Body)
	for _, catchClause := range stmt.Catches {
		sa.visitCatchClause(catchClause)
	}
	if stmt.Finally != nil {
		sa.visitBlockStatement(stmt.Finally)
	}
}

func (sa *SemanticAnalyzer) visitThrowStatement(stmt *ThrowStatement) {
	sa.visitExpression(stmt.Expression)
}

func (sa *SemanticAnalyzer) visitCatchClause(clause *CatchClause) {
	if clause.ExceptionType != nil {
		sa.SymbolTable.AddReference(clause.ExceptionType.Value, CLASS_SYMBOL, clause.Token.Line, 0)
	}
	sa.SymbolTable.DeclareSymbol(clause.Variable.Name, VARIABLE_SYMBOL, sa.CurrentFile, clause.Token.Line)
	sa.visitBlockStatement(clause.Body)
}

func (sa *SemanticAnalyzer) visitInfixExpression(expr *InfixExpression) {
	sa.visitExpression(expr.Left)
	sa.visitExpression(expr.Right)
}

func (sa *SemanticAnalyzer) visitPrefixExpression(expr *PrefixExpression) {
	sa.visitExpression(expr.Right)
}

func (sa *SemanticAnalyzer) visitPostfixExpression(expr *PostfixExpression) {
	sa.visitExpression(expr.Left)
}

func (sa *SemanticAnalyzer) visitArrayLiteral(expr *ArrayLiteral) {
	for _, element := range expr.Elements {
		sa.visitExpression(element)
	}
}

func (sa *SemanticAnalyzer) visitAssociativeArrayLiteral(expr *AssociativeArrayLiteral) {
	for _, pair := range expr.Pairs {
		sa.visitExpression(pair.Key)
		sa.visitExpression(pair.Value)
	}
}

func (sa *SemanticAnalyzer) visitIndexExpression(expr *IndexExpression) {
	sa.visitExpression(expr.Left)
	sa.visitExpression(expr.Index)
}

func (sa *SemanticAnalyzer) visitAnonymousFunction(expr *AnonymousFunction) {
	sa.SymbolTable.EnterScope("function", "anonymous")
	for _, param := range expr.Parameters {
		sa.SymbolTable.DeclareSymbol(param.Name, VARIABLE_SYMBOL, sa.CurrentFile, param.Token.Line)
	}
	for _, useVar := range expr.UseClause {
		sa.SymbolTable.AddReference(useVar.Name, VARIABLE_SYMBOL, useVar.Token.Line, 0)
	}
	sa.visitBlockStatement(expr.Body)
	sa.SymbolTable.ExitScope()
}

func (sa *SemanticAnalyzer) visitYieldExpression(expr *YieldExpression) {
	if expr.Key != nil {
		sa.visitExpression(expr.Key)
	}
	if expr.Value != nil {
		sa.visitExpression(expr.Value)
	}
}

func (sa *SemanticAnalyzer) visitTernaryExpression(expr *TernaryExpression) {
	sa.visitExpression(expr.Condition)
	sa.visitExpression(expr.TrueValue)
	sa.visitExpression(expr.FalseValue)
}

func (sa *SemanticAnalyzer) visitConstantDeclaration(stmt *ConstantDeclaration) {
	sa.SymbolTable.DeclareSymbol(stmt.Name.Value, CONSTANT_SYMBOL, sa.CurrentFile, stmt.Token.Line)
	sa.visitExpression(stmt.Value)
}

func (sa *SemanticAnalyzer) visitPropertyDeclaration(stmt *PropertyDeclaration) {
	sa.SymbolTable.DeclareSymbol(stmt.Name.Name, VARIABLE_SYMBOL, sa.CurrentFile, stmt.Token.Line)
	if stmt.Value != nil {
		sa.visitExpression(stmt.Value)
	}
}

func (sa *SemanticAnalyzer) visitMethodDeclaration(stmt *MethodDeclaration) {
	sa.SymbolTable.DeclareSymbol(stmt.Name.Value, FUNCTION_SYMBOL, sa.CurrentFile, stmt.Token.Line)

	sa.SymbolTable.EnterScope("method", stmt.Name.Value)
	for _, param := range stmt.Parameters {
		sa.SymbolTable.DeclareSymbol(param.Name, VARIABLE_SYMBOL, sa.CurrentFile, param.Token.Line)
	}
	sa.visitBlockStatement(stmt.Body)
	sa.SymbolTable.ExitScope()
}

func (sa *SemanticAnalyzer) visitInterfaceMethod(stmt *InterfaceMethod) {
	sa.SymbolTable.DeclareSymbol(stmt.Name.Value, FUNCTION_SYMBOL, sa.CurrentFile, stmt.Token.Line)
}

func (sa *SemanticAnalyzer) addIdentifierReference(identifier *Identifier) {
	// This could be a function call or constant reference
	// Try to resolve as function first, then as constant
	ref := sa.SymbolTable.AddReference(identifier.Value, FUNCTION_SYMBOL, identifier.Token.Line, 0)
	if ref.ResolvedSymbol == nil {
		sa.SymbolTable.AddReference(identifier.Value, CONSTANT_SYMBOL, identifier.Token.Line, 0)
	}
}

// AddError adds a semantic error
func (sa *SemanticAnalyzer) AddError(message string) {
	sa.Errors = append(sa.Errors, message)
}

// GetErrors returns all semantic errors
func (sa *SemanticAnalyzer) GetErrors() []string {
	return sa.Errors
}

// ValidateReferences validates all symbol references and reports errors
func (sa *SemanticAnalyzer) ValidateReferences() {
	for _, ref := range sa.SymbolTable.References {
		if ref.ResolvedSymbol == nil {
			sa.AddError(fmt.Sprintf("Undefined %s '%s' at line %d", 
				getSymbolTypeString(ref), ref.Name, ref.Line))
		}
	}
}

func getSymbolTypeString(_ *SymbolReference) string {
	// This is a simplified approach - in reality you'd track the expected type
	// The ref parameter is not used in this simple implementation
	return "symbol"
}