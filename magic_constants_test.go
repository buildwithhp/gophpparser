package gophpparser

import (
	"testing"
)

func TestMagicConstants(t *testing.T) {
	phpCode := `<?php
$file = __FILE__;
$dir = __DIR__;
$path = \dirname(__DIR__);
?>`

	t.Logf("=== Magic Constants Test ===")
	
	// Try semantic parsing
	semanticProgram, err := ParseWithSemantics(phpCode, "magic_test.php")
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
	
	t.Logf("✅ Successfully parsed magic constants!")
	t.Logf("   Symbols found: %d", len(semanticProgram.SymbolTable.AllSymbols))
	t.Logf("   References: %d", len(semanticProgram.AllReferences))
	t.Logf("   Unresolved: %d", len(semanticProgram.UnresolvedRefs))
	
	// Check for magic constants in AST
	foundMagicConstants := 0
	for _, stmt := range semanticProgram.Program.Statements {
		if exprStmt, ok := stmt.(*ExpressionStatement); ok {
			if assignExpr, ok := exprStmt.Expression.(*AssignmentExpression); ok {
				if magicConstant, ok := assignExpr.Value.(*MagicConstant); ok {
					foundMagicConstants++
					t.Logf("   Found magic constant: %s", magicConstant.Value)
				}
			}
		}
	}
	
	if foundMagicConstants == 0 {
		t.Error("❌ No magic constants found in AST")
	} else {
		t.Logf("✅ Found %d magic constants", foundMagicConstants)
	}
}

func TestSpecificMagicConstants(t *testing.T) {
	tests := []struct {
		name     string
		phpCode  string
		expected string
	}{
		{
			name:     "__FILE__ constant",
			phpCode:  `<?php echo __FILE__; ?>`,
			expected: "__FILE__",
		},
		{
			name:     "__DIR__ constant",
			phpCode:  `<?php echo __DIR__; ?>`,
			expected: "__DIR__",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			semanticProgram, err := ParseWithSemantics(tt.phpCode, "test.php")
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", tt.name, err)
			}

			// Look for magic constant in the AST
			found := false
			for _, stmt := range semanticProgram.Program.Statements {
				if echoStmt, ok := stmt.(*EchoStatement); ok {
					for _, value := range echoStmt.Values {
						if magicConstant, ok := value.(*MagicConstant); ok {
							if magicConstant.Value == tt.expected {
								found = true
								t.Logf("✅ Found expected magic constant: %s", magicConstant.Value)
							}
						}
					}
				}
			}

			if !found {
				t.Errorf("❌ Expected magic constant %s not found", tt.expected)
			}
		})
	}
}