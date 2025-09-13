# Semantic Analysis Guide

This guide explains how to use the semantic analysis features of the PHP parser to resolve symbol references and understand code structure.

## Overview

The semantic analyzer adds **symbol resolution** capabilities to the PHP parser, allowing you to:

- ✅ **Resolve class instantiations** to their fully qualified names
- ✅ **Track function calls** to their declarations  
- ✅ **Handle namespace imports** and aliases
- ✅ **Detect undefined symbols**
- ✅ **Analyze code dependencies**
- ✅ **Generate usage statistics**

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    parser "github.com/buildwithhp/gophpparser"
)

func main() {
    phpCode := `<?php
namespace App\Services;

use App\Models\User;

class UserService {
    public function createUser($name) {
        $user = new User($name);  // This will resolve to App\Models\User
        return $user;
    }
}

$service = new UserService();
?>`

    // Parse with semantic analysis
    semanticProgram, err := parser.ParseWithSemantics(phpCode, "example.php")
    if err != nil {
        log.Fatal(err)
    }

    // Show resolved class references
    fmt.Println("Class Instantiations:")
    for _, ref := range semanticProgram.AllReferences {
        if ref.ResolvedSymbol != nil && ref.ResolvedSymbol.Type == parser.CLASS_SYMBOL {
            fmt.Printf("  Line %d: 'new %s()' -> %s\n", 
                ref.Line, ref.Name, ref.ResolvedSymbol.FullyQualified)
        }
    }
}
```

**Output:**
```
Class Instantiations:
  Line 8: 'new User()' -> App\Models\User
  Line 12: 'new UserService()' -> App\Services\UserService
```

### Parse File with Semantics

```go
// Parse a PHP file directly
semanticProgram, err := parser.ParseFileWithSemantics("path/to/file.php")
if err != nil {
    log.Fatal(err)
}

// Access semantic information
fmt.Printf("Total symbols: %d\n", len(semanticProgram.SymbolTable.AllSymbols))
fmt.Printf("Total references: %d\n", len(semanticProgram.AllReferences))
fmt.Printf("Unresolved references: %d\n", len(semanticProgram.UnresolvedRefs))
```

## Core Concepts

### Symbol Types

The analyzer tracks different types of symbols:

```go
const (
    CLASS_SYMBOL     SymbolType = iota  // Classes
    FUNCTION_SYMBOL                     // Functions and methods
    VARIABLE_SYMBOL                     // Variables and properties
    CONSTANT_SYMBOL                     // Constants
    INTERFACE_SYMBOL                    // Interfaces
    TRAIT_SYMBOL                        // Traits
)
```

### Symbol Resolution Rules

The analyzer resolves symbols using PHP's resolution rules:

1. **Absolute references** (`\ClassName`) - Global namespace
2. **Import aliases** (`use App\User as AppUser`)
3. **Current scope** (local variables, methods)
4. **Current namespace** (`namespace App; new User()` -> `App\User`)
5. **Global namespace** (fallback)

### Scopes

The analyzer tracks nested scopes:

- **Global scope** - Top level
- **Namespace scope** - Within namespace declarations
- **Class scope** - Within class definitions
- **Function/Method scope** - Within function bodies

## Practical Examples

### Example 1: Resolving Conflicting Class Names

```php
<?php
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
    public function createUsers() {
        $hrUser = new HRUser();     // Resolves to HR\User
        $payrollUser = new User();  // Resolves to Payroll\User
        return [$hrUser, $payrollUser];
    }
}
?>
```

```go
semanticProgram, _ := parser.ParseWithSemantics(phpCode, "example.php")

// Find all class references
for _, ref := range semanticProgram.AllReferences {
    if ref.ResolvedSymbol != nil && ref.ResolvedSymbol.Type == parser.CLASS_SYMBOL {
        fmt.Printf("Line %d: '%s' -> %s\n", 
            ref.Line, ref.Name, ref.ResolvedSymbol.FullyQualified)
    }
}
```

**Output:**
```
Line 16: 'HRUser' -> HR\User
Line 17: 'User' -> Payroll\User
```

### Example 2: Function Call Resolution

```php
<?php
namespace Utils;

function calculateTax($amount) {
    return $amount * 0.3;
}

class Calculator {
    public function process($amount) {
        $tax = calculateTax($amount);  // Resolves to Utils\calculateTax
        return $amount - $tax;
    }
}
?>
```

```go
// Find function references
functionRefs := semanticProgram.GetFunctionReferences("calculateTax")
for _, ref := range functionRefs {
    fmt.Printf("Function call at line %d resolves to: %s\n", 
        ref.Line, ref.ResolvedSymbol.FullyQualified)
}
```

### Example 3: Detecting Undefined Symbols

```php
<?php
namespace App;

class Service {
    public function test() {
        $user = new UndefinedClass();  // This won't resolve
        return $user;
    }
}
?>
```

```go
semanticProgram, _ := parser.ParseWithSemantics(phpCode, "example.php")

// Check for unresolved references
if len(semanticProgram.UnresolvedRefs) > 0 {
    fmt.Println("Unresolved references:")
    for _, ref := range semanticProgram.UnresolvedRefs {
        fmt.Printf("  Line %d: '%s' could not be resolved\n", ref.Line, ref.Name)
    }
}
```

## Advanced Features

### 1. Usage Statistics

```go
stats := semanticProgram.GetUsageStatistics()

// Most used classes
if mostUsed, ok := stats["most_used_classes"].([]map[string]any); ok {
    fmt.Println("Most used classes:")
    for _, item := range mostUsed {
        symbol := item["symbol"].(*parser.Symbol)
        count := item["usage_count"].(int)
        fmt.Printf("  %s: %d uses\n", symbol.FullyQualified, count)
    }
}

// Unused symbols
if unused, ok := stats["unused_symbols"].([]*parser.Symbol); ok {
    fmt.Printf("Unused symbols: %d\n", len(unused))
    for _, symbol := range unused {
        fmt.Printf("  %s (declared but never used)\n", symbol.FullyQualified)
    }
}
```

### 2. Class Hierarchy Analysis

```go
// Get inheritance information
hierarchy := semanticProgram.GetClassHierarchy("App\\UserService")
if len(hierarchy) > 0 {
    fmt.Printf("App\\UserService extends/implements: %v\n", hierarchy)
}

// Find all classes in a namespace
symbols := semanticProgram.GetSymbolsInNamespace("App\\Models")
fmt.Printf("Classes in App\\Models: %d\n", len(symbols))
for _, symbol := range symbols {
    if symbol.Type == parser.CLASS_SYMBOL {
        fmt.Printf("  - %s\n", symbol.Name)
    }
}
```

### 3. Reference Report Generation

```go
report := semanticProgram.GenerateReferenceReport()

fmt.Printf("Analysis Report:\n")
fmt.Printf("  Resolution Rate: %.1f%%\n", report["summary"].(map[string]any)["resolution_rate"])

// Show statistics by symbol type
if byType, ok := report["by_symbol_type"].(map[string]map[string]int); ok {
    for symbolType, counts := range byType {
        declared := counts["declared"]
        referenced := counts["referenced"]
        fmt.Printf("  %s: %d declared, %d referenced\n", symbolType, declared, referenced)
    }
}
```

### 4. JSON Export with Semantic Information

```go
// Generate JSON with full semantic analysis
jsonData, err := semanticProgram.SemanticJSON()
if err != nil {
    log.Fatal(err)
}

// Save to file or send to another system
err = os.WriteFile("semantic_analysis.json", jsonData, 0644)
if err != nil {
    log.Fatal(err)
}
```

**JSON Structure:**
```json
{
  "type": "SemanticProgram",
  "program": { /* Original AST */ },
  "semantic_analysis": {
    "symbol_table": {
      "all_symbols": { /* Symbol definitions */ },
      "namespace_symbols": { /* Grouped by namespace */ },
      "class_hierarchy": { /* Inheritance info */ }
    },
    "references": {
      "all_references": [ /* All symbol references */ ],
      "unresolved_references": [ /* Failed resolutions */ ],
      "total_references": 45,
      "unresolved_count": 2
    },
    "statistics": {
      "total_symbols": 12,
      "total_namespaces": 3,
      "total_classes": 5,
      "total_functions": 4,
      "total_interfaces": 2,
      "total_traits": 1
    }
  }
}
```

## Integration Patterns

### 1. IDE Integration

```go
// For IDE features like "Go to Definition"
func FindDefinition(file string, line int, symbol string) *parser.Symbol {
    semanticProgram, err := parser.ParseFileWithSemantics(file)
    if err != nil {
        return nil
    }
    
    for _, ref := range semanticProgram.AllReferences {
        if ref.Line == line && ref.Name == symbol {
            return ref.ResolvedSymbol
        }
    }
    return nil
}
```

### 2. Dependency Analysis

```go
// Find all dependencies of a class
func FindClassDependencies(className string) []string {
    var dependencies []string
    
    for _, ref := range semanticProgram.AllReferences {
        if ref.ResolvedSymbol != nil {
            // Check if this reference is inside the target class
            // and refers to another class
            if ref.ResolvedSymbol.Type == parser.CLASS_SYMBOL {
                dependencies = append(dependencies, ref.ResolvedSymbol.FullyQualified)
            }
        }
    }
    
    return dependencies
}
```

### 3. Refactoring Support

```go
// Find all references to a symbol for safe renaming
func FindAllReferences(symbolFQN string) []*parser.SymbolReference {
    var refs []*parser.SymbolReference
    
    for _, ref := range semanticProgram.AllReferences {
        if ref.ResolvedSymbol != nil && ref.ResolvedSymbol.FullyQualified == symbolFQN {
            refs = append(refs, ref)
        }
    }
    
    return refs
}
```

## Error Handling

### Common Issues and Solutions

1. **Unresolved References**
   ```go
   if len(semanticProgram.UnresolvedRefs) > 0 {
       for _, ref := range semanticProgram.UnresolvedRefs {
           // Log or handle unresolved symbols
           fmt.Printf("Warning: Undefined %s '%s' at line %d\n", 
               getSymbolTypeHint(ref), ref.Name, ref.Line)
       }
   }
   ```

2. **Missing Import Declarations**
   - The analyzer only knows about symbols declared in the current file
   - For multi-file analysis, you'd need to parse all files together

3. **Dynamic Class Names**
   ```php
   $className = "User";
   $obj = new $className();  // Cannot be resolved statically
   ```

## Best Practices

1. **Parse Complete Codebases**: For accurate resolution, parse all related files together
2. **Handle Unresolved References**: Always check for and handle unresolved symbols
3. **Use Fully Qualified Names**: When comparing symbols, use `FullyQualified` field
4. **Cache Results**: Semantic analysis can be expensive, consider caching results
5. **Validate Input**: Ensure PHP code is syntactically correct before semantic analysis

## Performance Considerations

- **Memory Usage**: Symbol tables can be large for big codebases
- **Processing Time**: O(n) where n is the number of AST nodes
- **Optimization**: Consider processing files in dependency order

## Limitations

1. **Single File Analysis**: Currently analyzes one file at a time
2. **No Include Resolution**: Doesn't follow `include`/`require` statements  
3. **Dynamic Resolution**: Cannot resolve dynamically constructed names
4. **Incomplete PHP**: Some advanced PHP features may not be fully supported

## Future Enhancements

- **Multi-file analysis** with dependency resolution
- **Type inference** for variables and expressions
- **Control flow analysis** for more sophisticated error detection
- **Performance optimizations** for large codebases

This semantic analysis system transforms the PHP parser from a simple syntax analyzer into a powerful code understanding tool, enabling advanced IDE features, static analysis, and code intelligence capabilities.