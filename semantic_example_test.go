package gophpparser

import (
	"fmt"
	"testing"
)

func TestSemanticAnalysisExample(t *testing.T) {
	// Example PHP code with multiple namespaces and class references
	phpCode := `<?php
namespace Payroll;

use HR\User;
use Finance\Calculator as Calc;

class PayrollService {
    private $calculator;
    private $users = [];
    
    public function __construct() {
        $this->calculator = new Calc();
    }
    
    public function processPayroll($userId) {
        $user = new User($userId);  // This refers to HR\User
        $salary = $this->calculator->calculate($user->getSalary());
        return $salary;
    }
    
    public function createLocalUser($name) {
        $localUser = new PayrollUser($name);  // This refers to Payroll\PayrollUser
        return $localUser;
    }
}

class PayrollUser {
    private $name;
    
    public function __construct($name) {
        $this->name = $name;
    }
}

function calculateTax($amount) {
    return $amount * 0.3;
}

// Usage examples
$service = new PayrollService();
$result = $service->processPayroll(123);
$tax = calculateTax($result);
?>`

	// Parse with semantic analysis
	semanticProgram, err := ParseWithSemantics(phpCode, "payroll.php")
	if err != nil {
		t.Fatalf("Failed to parse with semantics: %v", err)
	}

	// Test 1: Verify namespace resolution
	t.Run("NamespaceResolution", func(t *testing.T) {
		// Check that PayrollService is in Payroll namespace
		payrollService := semanticProgram.GetSymbolByFullyQualifiedName("Payroll\\PayrollService")
		if payrollService == nil {
			t.Error("PayrollService not found in symbol table")
		} else {
			if payrollService.Namespace != "Payroll" {
				t.Errorf("Expected namespace 'Payroll', got '%s'", payrollService.Namespace)
			}
			t.Logf("✓ PayrollService found: %s (namespace: %s)", payrollService.FullyQualified, payrollService.Namespace)
		}

		// Check PayrollUser is also in Payroll namespace
		payrollUser := semanticProgram.GetSymbolByFullyQualifiedName("Payroll\\PayrollUser")
		if payrollUser == nil {
			t.Error("PayrollUser not found in symbol table")
		} else {
			t.Logf("✓ PayrollUser found: %s", payrollUser.FullyQualified)
		}
	})

	// Test 2: Verify class instantiation resolution
	t.Run("ClassInstantiationResolution", func(t *testing.T) {
		// Find all class references
		classRefs := make(map[string][]*SymbolReference)
		for _, ref := range semanticProgram.AllReferences {
			if ref.ResolvedSymbol != nil && ref.ResolvedSymbol.Type == CLASS_SYMBOL {
				className := ref.ResolvedSymbol.FullyQualified
				classRefs[className] = append(classRefs[className], ref)
			}
		}

		t.Logf("Found class references:")
		for className, refs := range classRefs {
			t.Logf("  %s: %d references", className, len(refs))
			for _, ref := range refs {
				t.Logf("    - Line %d: '%s' -> %s", ref.Line, ref.Name, ref.ResolvedSymbol.FullyQualified)
			}
		}

		// Verify specific resolutions
		foundHRUser := false
		foundPayrollUser := false
		foundCalc := false

		for _, ref := range semanticProgram.AllReferences {
			if ref.ResolvedSymbol != nil && ref.ResolvedSymbol.Type == CLASS_SYMBOL {
				switch ref.ResolvedSymbol.FullyQualified {
				case "HR\\User":
					if ref.Name == "User" {
						foundHRUser = true
						t.Logf("✓ 'User' correctly resolved to HR\\User at line %d", ref.Line)
					}
				case "Payroll\\PayrollUser":
					if ref.Name == "PayrollUser" {
						foundPayrollUser = true
						t.Logf("✓ 'PayrollUser' correctly resolved to Payroll\\PayrollUser at line %d", ref.Line)
					}
				case "Finance\\Calculator":
					if ref.Name == "Calc" {
						foundCalc = true
						t.Logf("✓ 'Calc' correctly resolved to Finance\\Calculator at line %d", ref.Line)
					}
				}
			}
		}

		if !foundHRUser {
			t.Error("❌ Failed to resolve 'User' to HR\\User")
		}
		if !foundPayrollUser {
			t.Error("❌ Failed to resolve 'PayrollUser' to Payroll\\PayrollUser")
		}
		if !foundCalc {
			t.Error("❌ Failed to resolve 'Calc' to Finance\\Calculator")
		}
	})

	// Test 3: Function resolution
	t.Run("FunctionResolution", func(t *testing.T) {
		// Find calculateTax function
		calculateTaxSymbol := semanticProgram.GetSymbolByFullyQualifiedName("Payroll\\calculateTax")
		if calculateTaxSymbol == nil {
			t.Error("calculateTax function not found")
		} else {
			t.Logf("✓ calculateTax function found: %s", calculateTaxSymbol.FullyQualified)
		}

		// Find function calls
		functionRefs := semanticProgram.GetFunctionReferences("calculateTax")
		if len(functionRefs) == 0 {
			t.Error("No references to calculateTax found")
		} else {
			t.Logf("✓ Found %d references to calculateTax", len(functionRefs))
			for _, ref := range functionRefs {
				t.Logf("  - Line %d: %s -> %s", ref.Line, ref.Name, ref.ResolvedSymbol.FullyQualified)
			}
		}
	})

	// Test 4: Symbol statistics
	t.Run("SymbolStatistics", func(t *testing.T) {
		report := semanticProgram.GenerateReferenceReport()
		t.Logf("Reference Report:")
		
		if summary, ok := report["summary"].(map[string]any); ok {
			t.Logf("  Total symbols: %v", summary["total_symbols"])
			t.Logf("  Total references: %v", summary["total_references"])
			t.Logf("  Unresolved references: %v", summary["unresolved_references"])
			t.Logf("  Resolution rate: %.1f%%", summary["resolution_rate"])
		}

		if byType, ok := report["by_symbol_type"].(map[string]map[string]int); ok {
			t.Logf("  By symbol type:")
			for symbolType, counts := range byType {
				t.Logf("    %s: %d declared, %d referenced", symbolType, counts["declared"], counts["referenced"])
			}
		}
	})

	// Test 5: Unresolved references
	t.Run("UnresolvedReferences", func(t *testing.T) {
		if len(semanticProgram.UnresolvedRefs) > 0 {
			t.Logf("Unresolved references:")
			for _, ref := range semanticProgram.UnresolvedRefs {
				t.Logf("  - Line %d: '%s' (could not resolve)", ref.Line, ref.Name)
			}
		} else {
			t.Logf("✓ All references resolved successfully!")
		}
	})

	// Test 6: Generate JSON output with semantic info
	t.Run("SemanticJSON", func(t *testing.T) {
		jsonData, err := semanticProgram.SemanticJSON()
		if err != nil {
			t.Errorf("Failed to generate semantic JSON: %v", err)
		} else {
			t.Logf("✓ Generated semantic JSON (%d bytes)", len(jsonData))
			// Optionally print first 500 characters for inspection
			if len(jsonData) > 500 {
				t.Logf("JSON preview: %s...", string(jsonData[:500]))
			} else {
				t.Logf("JSON output: %s", string(jsonData))
			}
		}
	})
}

func TestMultipleNamespaceResolution(t *testing.T) {
	// Test case with conflicting class names in different namespaces
	phpCode := `<?php
namespace HR;

class User {
    public $department = "HR";
}

namespace Payroll;

use HR\User as HRUser;

class User {
    public $department = "Payroll";
}

class Service {
    public function test() {
        $hrUser = new HRUser();      // Should resolve to HR\User
        $payrollUser = new User();   // Should resolve to Payroll\User
        $globalUser = new \User();   // Should resolve to global \User (if it existed)
        return [$hrUser, $payrollUser];
    }
}
?>`

	semanticProgram, err := ParseWithSemantics(phpCode, "multi_namespace.php")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Verify we have both User classes
	hrUser := semanticProgram.GetSymbolByFullyQualifiedName("HR\\User")
	payrollUser := semanticProgram.GetSymbolByFullyQualifiedName("Payroll\\User")

	if hrUser == nil {
		t.Error("HR\\User not found")
	} else {
		t.Logf("✓ Found HR\\User: %s", hrUser.FullyQualified)
	}

	if payrollUser == nil {
		t.Error("Payroll\\User not found")
	} else {
		t.Logf("✓ Found Payroll\\User: %s", payrollUser.FullyQualified)
	}

	// Check that references are resolved correctly
	t.Logf("Class references:")
	for _, ref := range semanticProgram.AllReferences {
		if ref.ResolvedSymbol != nil && ref.ResolvedSymbol.Type == CLASS_SYMBOL {
			t.Logf("  Line %d: '%s' -> %s", ref.Line, ref.Name, ref.ResolvedSymbol.FullyQualified)
		}
	}
}

func ExampleSemanticAnalysis() {
	// Simple example showing how to use semantic analysis
	phpCode := `<?php
namespace App\Services;

use App\Models\User;
use Database\Connection;

class UserService {
    public function createUser($name) {
        $user = new User($name);
        return $user;
    }
    
    public function findUser($id) {
        $connection = new Connection();
        return $connection->findById($id);
    }
}

$service = new UserService();
$user = $service->createUser("John");
?>`

	// Parse with semantic analysis
	semanticProgram, err := ParseWithSemantics(phpCode, "example.php")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Show all class instantiations with their resolved references
	fmt.Println("Class Instantiations:")
	for _, ref := range semanticProgram.AllReferences {
		if ref.ResolvedSymbol != nil && ref.ResolvedSymbol.Type == CLASS_SYMBOL {
			fmt.Printf("  Line %d: 'new %s()' resolves to %s\n", 
				ref.Line, ref.Name, ref.ResolvedSymbol.FullyQualified)
		}
	}

	// Show symbol table summary
	fmt.Printf("\nSymbol Summary:\n")
	fmt.Printf("  Total symbols: %d\n", len(semanticProgram.SymbolTable.AllSymbols))
	fmt.Printf("  Total references: %d\n", len(semanticProgram.AllReferences))
	fmt.Printf("  Unresolved: %d\n", len(semanticProgram.UnresolvedRefs))

	// Output:
	// Class Instantiations:
	//   Line 8: 'new User()' resolves to App\Models\User
	//   Line 13: 'new Connection()' resolves to Database\Connection
	//   Line 17: 'new UserService()' resolves to App\Services\UserService
	//
	// Symbol Summary:
	//   Total symbols: 2
	//   Total references: 3
	//   Unresolved: 0
}