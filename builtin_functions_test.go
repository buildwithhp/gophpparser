package gophpparser

import (
	"testing"
)

func TestBuiltinFunctions(t *testing.T) {
	phpCode := `<?php
$path = dirname(__FILE__);
$base = basename($path);
$exists = file_exists($path);
$length = strlen("hello");
$merged = array_merge($arr1, $arr2);
$parts = explode(",", $str);
$joined = implode(" ", $parts);
$trimmed = trim($data);
$replaced = str_replace("old", "new", $text);
$json = json_encode($data);
$decoded = json_decode($json);
?>`

	t.Logf("=== Built-in Functions Test ===")
	
	// Try semantic parsing
	semanticProgram, err := ParseWithSemantics(phpCode, "builtin_test.php")
	if err != nil {
		t.Logf("❌ Parse error: %v", err)
		
		// Debug what's failing
		debug := DebugParsePHP(phpCode)
		t.Logf("Parsing errors: %d", len(debug.ParsingErrors))
		if len(debug.ParsingErrors) > 0 {
			t.Logf("First few errors:")
			for i := 0; i < 3 && i < len(debug.ParsingErrors); i++ {
				t.Logf("  %d. %s", i+1, debug.ParsingErrors[i])
			}
		}
		return
	}
	
	t.Logf("✅ Successfully parsed built-in functions!")
	t.Logf("   Symbols found: %d", len(semanticProgram.SymbolTable.AllSymbols))
	t.Logf("   References: %d", len(semanticProgram.AllReferences))
	t.Logf("   Unresolved: %d", len(semanticProgram.UnresolvedRefs))
	
	// Count function calls
	functionCalls := 0
	for _, stmt := range semanticProgram.Program.Statements {
		if exprStmt, ok := stmt.(*ExpressionStatement); ok {
			if assignExpr, ok := exprStmt.Expression.(*AssignmentExpression); ok {
				if callExpr, ok := assignExpr.Value.(*CallExpression); ok {
					functionCalls++
					if ident, ok := callExpr.Function.(*Identifier); ok {
						t.Logf("   Found function call: %s", ident.Value)
					}
				}
			}
		}
	}
	
	if functionCalls == 0 {
		t.Error("❌ No function calls found in AST")
	} else {
		t.Logf("✅ Found %d function calls", functionCalls)
	}
}

func TestSpecificBuiltinFunctions(t *testing.T) {
	tests := []struct {
		name     string
		phpCode  string
		function string
	}{
		{
			name:     "dirname function",
			phpCode:  `<?php $path = dirname(__FILE__); ?>`,
			function: "dirname",
		},
		{
			name:     "basename function", 
			phpCode:  `<?php $base = basename($path); ?>`,
			function: "basename",
		},
		{
			name:     "file_exists function",
			phpCode:  `<?php $exists = file_exists($file); ?>`,
			function: "file_exists",
		},
		{
			name:     "strlen function",
			phpCode:  `<?php $len = strlen($str); ?>`,
			function: "strlen",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			semanticProgram, err := ParseWithSemantics(tt.phpCode, "test.php")
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", tt.name, err)
			}

			// Look for function call in the AST
			found := false
			for _, stmt := range semanticProgram.Program.Statements {
				if exprStmt, ok := stmt.(*ExpressionStatement); ok {
					if assignExpr, ok := exprStmt.Expression.(*AssignmentExpression); ok {
						if callExpr, ok := assignExpr.Value.(*CallExpression); ok {
							if ident, ok := callExpr.Function.(*Identifier); ok {
								if ident.Value == tt.function {
									found = true
									t.Logf("✅ Found expected function: %s", ident.Value)
								}
							}
						}
					}
				}
			}

			if !found {
				t.Errorf("❌ Expected function %s not found", tt.function)
			}
		})
	}
}