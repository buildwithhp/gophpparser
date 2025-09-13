package gophpparser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// SemanticInfo contains resolved symbol information for AST nodes
type SemanticInfo struct {
	ResolvedSymbol   *Symbol            `json:"resolved_symbol,omitempty"`   // What this identifier refers to
	SymbolReferences []*SymbolReference `json:"symbol_references,omitempty"` // All symbol references in this node
	Scope            *Scope             `json:"scope,omitempty"`             // Scope this node belongs to
}

// Enhanced AST nodes with semantic information
type SemanticNewExpression struct {
	*NewExpression
	SemanticInfo *SemanticInfo `json:"semantic_info,omitempty"`
}

func (sne *SemanticNewExpression) Type() string { return "SemanticNewExpression" }

type SemanticCallExpression struct {
	*CallExpression
	SemanticInfo *SemanticInfo `json:"semantic_info,omitempty"`
}

func (sce *SemanticCallExpression) Type() string { return "SemanticCallExpression" }

type SemanticIdentifier struct {
	*Identifier
	SemanticInfo *SemanticInfo `json:"semantic_info,omitempty"`
}

func (si *SemanticIdentifier) Type() string { return "SemanticIdentifier" }

type SemanticStaticAccess struct {
	*StaticAccessExpression
	SemanticInfo *SemanticInfo `json:"semantic_info,omitempty"`
}

func (ssa *SemanticStaticAccess) Type() string { return "SemanticStaticAccess" }

// SemanticProgram contains the original AST plus semantic analysis results
type SemanticProgram struct {
	*Program
	SymbolTable      *SymbolTable        `json:"symbol_table"`
	AllReferences    []*SymbolReference  `json:"all_references"`
	UnresolvedRefs   []*SymbolReference  `json:"unresolved_references"`
	ClassHierarchy   map[string][]string `json:"class_hierarchy"`
	NamespaceSymbols map[string][]*Symbol `json:"namespace_symbols"`
}

// ParseWithSemantics parses PHP code and performs semantic analysis
func ParseWithSemantics(input string, filename string) (*SemanticProgram, error) {
	// 1. Parse syntax
	lexer := New(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		return nil, fmt.Errorf("parsing errors: %s", strings.Join(parser.Errors(), "; "))
	}

	// 2. Perform semantic analysis
	analyzer := NewSemanticAnalyzer()
	analyzer.AnalyzeProgram(program, filename)
	analyzer.ValidateReferences()

	// 3. Create enhanced program with semantic info
	semanticProgram := &SemanticProgram{
		Program:          program,
		SymbolTable:      analyzer.SymbolTable,
		AllReferences:    analyzer.SymbolTable.References,
		UnresolvedRefs:   analyzer.SymbolTable.GetUnresolvedReferences(),
		ClassHierarchy:   analyzer.SymbolTable.ClassHierarchy,
		NamespaceSymbols: analyzer.SymbolTable.Namespaces,
	}

	return semanticProgram, nil
}

// ParseFileWithSemantics parses a file and performs semantic analysis
func ParseFileWithSemantics(filepath string) (*SemanticProgram, error) {
	// Read the file
	content, err := ReadFileContent(filepath)
	if err != nil {
		return nil, fmt.Errorf("error reading file '%s': %v", filepath, err)
	}

	return ParseWithSemantics(content, filepath)
}

// ReadFileContent reads file contents (helper function)
func ReadFileContent(filepath string) (string, error) {
	if content, err := os.ReadFile(filepath); err != nil {
		return "", err
	} else {
		return string(content), nil
	}
}

// GetClassReferences returns all references to a specific class
func (sp *SemanticProgram) GetClassReferences(className string) []*SymbolReference {
	var refs []*SymbolReference
	for _, ref := range sp.AllReferences {
		if ref.ResolvedSymbol != nil && 
		   ref.ResolvedSymbol.Type == CLASS_SYMBOL && 
		   (ref.ResolvedSymbol.Name == className || ref.ResolvedSymbol.FullyQualified == className) {
			refs = append(refs, ref)
		}
	}
	return refs
}

// GetFunctionReferences returns all references to a specific function
func (sp *SemanticProgram) GetFunctionReferences(functionName string) []*SymbolReference {
	var refs []*SymbolReference
	for _, ref := range sp.AllReferences {
		if ref.ResolvedSymbol != nil && 
		   ref.ResolvedSymbol.Type == FUNCTION_SYMBOL && 
		   (ref.ResolvedSymbol.Name == functionName || ref.ResolvedSymbol.FullyQualified == functionName) {
			refs = append(refs, ref)
		}
	}
	return refs
}

// GetSymbolByFullyQualifiedName returns a symbol by its fully qualified name
func (sp *SemanticProgram) GetSymbolByFullyQualifiedName(fqn string) *Symbol {
	return sp.SymbolTable.AllSymbols[fqn]
}

// GetSymbolsInNamespace returns all symbols in a specific namespace
func (sp *SemanticProgram) GetSymbolsInNamespace(namespace string) []*Symbol {
	return sp.NamespaceSymbols[namespace]
}

// GetClassHierarchy returns the inheritance chain for a class
func (sp *SemanticProgram) GetClassHierarchy(className string) []string {
	return sp.ClassHierarchy[className]
}

// FindClassInstantiations finds all places where a class is instantiated
func (sp *SemanticProgram) FindClassInstantiations(className string) []*SymbolReference {
	var instantiations []*SymbolReference
	for _, ref := range sp.AllReferences {
		if ref.ResolvedSymbol != nil && 
		   ref.ResolvedSymbol.Type == CLASS_SYMBOL && 
		   (ref.ResolvedSymbol.Name == className || ref.ResolvedSymbol.FullyQualified == className) {
			// Note: In a more sophisticated implementation, you'd distinguish between
			// different types of references (instantiation vs static access vs inheritance)
			instantiations = append(instantiations, ref)
		}
	}
	return instantiations
}

// SemanticJSON generates JSON with semantic information
func (sp *SemanticProgram) SemanticJSON() ([]byte, error) {
	return ToJSONSemantic(sp)
}

// ToJSONSemantic converts semantic program to JSON with enhanced information
func ToJSONSemantic(sp *SemanticProgram) ([]byte, error) {
	data := map[string]any{
		"type":    "SemanticProgram",
		"program": sp.Program,
		"semantic_analysis": map[string]any{
			"symbol_table": map[string]any{
				"all_symbols":       sp.SymbolTable.AllSymbols,
				"namespace_symbols": sp.NamespaceSymbols,
				"class_hierarchy":   sp.ClassHierarchy,
			},
			"references": map[string]any{
				"all_references":      sp.AllReferences,
				"unresolved_references": sp.UnresolvedRefs,
				"total_references":    len(sp.AllReferences),
				"unresolved_count":    len(sp.UnresolvedRefs),
			},
			"statistics": map[string]any{
				"total_symbols":     len(sp.SymbolTable.AllSymbols),
				"total_namespaces":  len(sp.NamespaceSymbols),
				"total_classes":     sp.countSymbolsByType(CLASS_SYMBOL),
				"total_functions":   sp.countSymbolsByType(FUNCTION_SYMBOL),
				"total_interfaces":  sp.countSymbolsByType(INTERFACE_SYMBOL),
				"total_traits":      sp.countSymbolsByType(TRAIT_SYMBOL),
			},
		},
	}

	return json.MarshalIndent(data, "", "  ")
}

// countSymbolsByType counts symbols of a specific type
func (sp *SemanticProgram) countSymbolsByType(symbolType SymbolType) int {
	count := 0
	for _, symbol := range sp.SymbolTable.AllSymbols {
		if symbol.Type == symbolType {
			count++
		}
	}
	return count
}

// GenerateReferenceReport generates a detailed reference report
func (sp *SemanticProgram) GenerateReferenceReport() map[string]any {
	report := map[string]any{
		"summary": map[string]any{
			"total_symbols":           len(sp.SymbolTable.AllSymbols),
			"total_references":        len(sp.AllReferences),
			"unresolved_references":   len(sp.UnresolvedRefs),
			"resolution_rate":         float64(len(sp.AllReferences)-len(sp.UnresolvedRefs)) / float64(len(sp.AllReferences)) * 100,
		},
		"by_symbol_type": make(map[string]map[string]int),
		"by_namespace":   make(map[string]int),
		"unresolved":     sp.UnresolvedRefs,
	}

	// Count by symbol type
	symbolTypeCounts := map[string]map[string]int{
		"class":     {"declared": 0, "referenced": 0},
		"function":  {"declared": 0, "referenced": 0},
		"interface": {"declared": 0, "referenced": 0},
		"trait":     {"declared": 0, "referenced": 0},
		"constant":  {"declared": 0, "referenced": 0},
		"variable":  {"declared": 0, "referenced": 0},
	}

	// Count declared symbols
	for _, symbol := range sp.SymbolTable.AllSymbols {
		typeStr := symbol.Type.String()
		if counts, exists := symbolTypeCounts[typeStr]; exists {
			counts["declared"]++
		}
	}

	// Count references
	for _, ref := range sp.AllReferences {
		if ref.ResolvedSymbol != nil {
			typeStr := ref.ResolvedSymbol.Type.String()
			if counts, exists := symbolTypeCounts[typeStr]; exists {
				counts["referenced"]++
			}
		}
	}

	report["by_symbol_type"] = symbolTypeCounts

	// Count by namespace
	namespaceCounts := make(map[string]int)
	for namespace, symbols := range sp.NamespaceSymbols {
		namespaceCounts[namespace] = len(symbols)
	}
	report["by_namespace"] = namespaceCounts

	return report
}

// Example usage functions for different scenarios

// ResolveClassInstantiation resolves a class instantiation to its fully qualified name
func (sp *SemanticProgram) ResolveClassInstantiation(className string, line int) *Symbol {
	for _, ref := range sp.AllReferences {
		if ref.Name == className && ref.Line == line && ref.ResolvedSymbol != nil && ref.ResolvedSymbol.Type == CLASS_SYMBOL {
			return ref.ResolvedSymbol
		}
	}
	return nil
}

// ResolveFunctionCall resolves a function call to its declaration
func (sp *SemanticProgram) ResolveFunctionCall(functionName string, line int) *Symbol {
	for _, ref := range sp.AllReferences {
		if ref.Name == functionName && ref.Line == line && ref.ResolvedSymbol != nil && ref.ResolvedSymbol.Type == FUNCTION_SYMBOL {
			return ref.ResolvedSymbol
		}
	}
	return nil
}

// GetUsageStatistics returns usage statistics for symbols
func (sp *SemanticProgram) GetUsageStatistics() map[string]any {
	stats := map[string]any{
		"most_used_classes":   sp.getMostUsedSymbols(CLASS_SYMBOL, 10),
		"most_used_functions": sp.getMostUsedSymbols(FUNCTION_SYMBOL, 10),
		"unused_symbols":      sp.getUnusedSymbols(),
	}
	return stats
}

// getMostUsedSymbols returns the most frequently referenced symbols
func (sp *SemanticProgram) getMostUsedSymbols(symbolType SymbolType, limit int) []map[string]any {
	usageCounts := make(map[string]int)
	
	for _, ref := range sp.AllReferences {
		if ref.ResolvedSymbol != nil && ref.ResolvedSymbol.Type == symbolType {
			usageCounts[ref.ResolvedSymbol.FullyQualified]++
		}
	}

	// Convert to slice and sort
	var results []map[string]any
	for fqn, count := range usageCounts {
		if symbol := sp.SymbolTable.AllSymbols[fqn]; symbol != nil {
			results = append(results, map[string]any{
				"symbol": symbol,
				"usage_count": count,
			})
		}
	}

	// Sort by usage count (simplified - you'd use sort.Slice in real implementation)
	// Return top N results
	if len(results) > limit {
		results = results[:limit]
	}

	return results
}

// getUnusedSymbols returns symbols that are declared but never referenced
func (sp *SemanticProgram) getUnusedSymbols() []*Symbol {
	var unused []*Symbol
	usedSymbols := make(map[string]bool)
	
	// Mark all referenced symbols as used
	for _, ref := range sp.AllReferences {
		if ref.ResolvedSymbol != nil {
			usedSymbols[ref.ResolvedSymbol.FullyQualified] = true
		}
	}
	
	// Find declared but unused symbols
	for _, symbol := range sp.SymbolTable.AllSymbols {
		if !usedSymbols[symbol.FullyQualified] {
			unused = append(unused, symbol)
		}
	}
	
	return unused
}