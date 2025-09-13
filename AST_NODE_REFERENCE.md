# PHP Parser AST Node Reference

This document provides a comprehensive reference for all Abstract Syntax Tree (AST) node types in the Go PHP Parser.

## Table of Contents

- [Base Interfaces](#base-interfaces)
- [Core Program Structure](#core-program-structure)
- [Literal Expressions](#literal-expressions)
- [Variables and Identifiers](#variables-and-identifiers)
- [Operators and Expressions](#operators-and-expressions)
- [Control Flow Statements](#control-flow-statements)
- [Function and Method Declarations](#function-and-method-declarations)
- [Object-Oriented Programming](#object-oriented-programming)
- [Arrays and Collections](#arrays-and-collections)
- [Namespace and Import System](#namespace-and-import-system)
- [Exception Handling](#exception-handling)
- [Advanced Features](#advanced-features)
- [JSON Serialization](#json-serialization)

## Base Interfaces

### Node Interface
```go
type Node interface {
    String() string        // Human-readable string representation
    TokenLiteral() string // Source token literal
    Type() string         // Node type identifier for JSON
}
```

### Statement Interface
```go
type Statement interface {
    Node
    statementNode()
}
```

### Expression Interface  
```go
type Expression interface {
    Node
    expressionNode()
}
```

## Core Program Structure

### Program
**Type:** Root node  
**Description:** Top-level container for all PHP statements  

```go
type Program struct {
    Statements []Statement `json:"statements"`
}
```

**PHP Example:**
```php
<?php
$x = 1;
echo $x;
?>
```

**JSON Output:**
```json
{
  "type": "Program",
  "statements": [...]
}
```

---

## Literal Expressions

### IntegerLiteral
**Type:** Expression  
**Description:** Integer numeric literals  

```go
type IntegerLiteral struct {
    Token Token `json:"token"`
    Value int64 `json:"value"`
}
```

**PHP Examples:**
```php
42
0x1A    // Hexadecimal
0755    // Octal
```

### FloatLiteral  
**Type:** Expression  
**Description:** Floating-point numeric literals  

```go
type FloatLiteral struct {
    Token Token   `json:"token"`
    Value float64 `json:"value"`
}
```

**PHP Examples:**
```php
3.14
1.5e10
.5
```

### StringLiteral
**Type:** Expression  
**Description:** String literals (quoted strings)  

```go
type StringLiteral struct {
    Token Token  `json:"token"`
    Value string `json:"value"`
}
```

**PHP Examples:**
```php
"Hello World"
'Single quotes'
```

### BooleanLiteral
**Type:** Expression  
**Description:** Boolean true/false values  

```go
type BooleanLiteral struct {
    Token Token `json:"token"`
    Value bool  `json:"value"`
}
```

**PHP Examples:**
```php
true
false
```

---

## Variables and Identifiers

### Variable
**Type:** Expression  
**Description:** PHP variables (prefixed with $)  

```go
type Variable struct {
    Token Token  `json:"token"`
    Name  string `json:"name"`  // Without the $ prefix
}
```

**PHP Examples:**
```php
$userName
$_POST
$GLOBALS
```

### Identifier
**Type:** Expression  
**Description:** Function names, class names, constants  

```go
type Identifier struct {
    Token Token  `json:"token"`
    Value string `json:"value"`
}
```

**PHP Examples:**
```php
myFunction
ClassName
CONSTANT_NAME
```

### NamespacedIdentifier
**Type:** Expression  
**Description:** Fully qualified or namespaced identifiers  

```go
type NamespacedIdentifier struct {
    Token     Token         `json:"token"`
    Namespace []*Identifier `json:"namespace"`
    Name      *Identifier   `json:"name"`
}
```

**PHP Examples:**
```php
\Exception                    // Global namespace
App\Controllers\UserController
```

---

## Operators and Expressions

### AssignmentExpression
**Type:** Expression  
**Description:** Variable assignment operations  

```go
type AssignmentExpression struct {
    Token Token      `json:"token"`
    Name  *Variable  `json:"name"`
    Value Expression `json:"value"`
}
```

**PHP Examples:**
```php
$x = 5
$name = "John"
$result = $a + $b
```

### InfixExpression
**Type:** Expression  
**Description:** Binary operations (arithmetic, comparison, logical)  

```go
type InfixExpression struct {
    Token    Token      `json:"token"`
    Left     Expression `json:"left"`
    Operator string     `json:"operator"`
    Right    Expression `json:"right"`
}
```

**PHP Examples:**
```php
$a + $b          // Arithmetic
$x == $y         // Comparison  
$flag && $other  // Logical
$str . $suffix   // String concatenation
```

### PrefixExpression
**Type:** Expression  
**Description:** Unary operations (negation, increment, etc.)  

```go
type PrefixExpression struct {
    Token    Token      `json:"token"`
    Operator string     `json:"operator"`
    Right    Expression `json:"right"`
}
```

**PHP Examples:**
```php
-$value      // Negation
!$flag       // Logical NOT
++$counter   // Pre-increment
--$index     // Pre-decrement
```

### PostfixExpression
**Type:** Expression  
**Description:** Post-increment and post-decrement operations  

```go
type PostfixExpression struct {
    Token    Token      `json:"token"`
    Left     Expression `json:"left"`
    Operator string     `json:"operator"`
}
```

**PHP Examples:**
```php
$counter++   // Post-increment
$index--     // Post-decrement
```

### TernaryExpression
**Type:** Expression  
**Description:** Conditional ternary operator  

```go
type TernaryExpression struct {
    Token      Token      `json:"token"`
    Condition  Expression `json:"condition"`
    TrueValue  Expression `json:"true_value"`
    FalseValue Expression `json:"false_value"`
}
```

**PHP Examples:**
```php
$result = $x > 0 ? "positive" : "negative"
$value = isset($data) ? $data : "default"
```

---

## Control Flow Statements

### ExpressionStatement
**Type:** Statement  
**Description:** Wrapper for expressions used as statements  

```go
type ExpressionStatement struct {
    Token      Token      `json:"token"`
    Expression Expression `json:"expression"`
}
```

**PHP Examples:**
```php
$x = 5;        // Assignment expression as statement
functionCall(); // Function call as statement
```

### IfStatement
**Type:** Statement  
**Description:** Conditional execution (if/else)  

```go
type IfStatement struct {
    Token       Token           `json:"token"`
    Condition   Expression      `json:"condition"`
    Consequence *BlockStatement `json:"consequence"`
    Alternative *BlockStatement `json:"alternative"`
}
```

**PHP Examples:**
```php
if ($x > 0) {
    echo "positive";
} else {
    echo "not positive";
}
```

### ForStatement
**Type:** Statement  
**Description:** Traditional for loops  

```go
type ForStatement struct {
    Token     Token           `json:"token"`
    Init      Expression      `json:"init"`
    Condition Expression      `json:"condition"`
    Update    Expression      `json:"update"`
    Body      *BlockStatement `json:"body"`
}
```

**PHP Examples:**
```php
for ($i = 0; $i < 10; $i++) {
    echo $i;
}
```

### WhileStatement
**Type:** Statement  
**Description:** While loop construct  

```go
type WhileStatement struct {
    Token     Token           `json:"token"`
    Condition Expression      `json:"condition"`
    Body      *BlockStatement `json:"body"`
}
```

**PHP Examples:**
```php
while ($i < 10) {
    echo $i;
    $i++;
}
```

### ForeachStatement
**Type:** Statement  
**Description:** Foreach loop for array iteration  

```go
type ForeachStatement struct {
    Token Token           `json:"token"`
    Array Expression      `json:"array"`
    Key   *Variable       `json:"key"`
    Value *Variable       `json:"value"`
    Body  *BlockStatement `json:"body"`
}
```

**PHP Examples:**
```php
foreach ($array as $value) {
    echo $value;
}

foreach ($array as $key => $value) {
    echo "$key: $value";
}
```

### BreakStatement
**Type:** Statement  
**Description:** Break from loops with optional level  

```go
type BreakStatement struct {
    Token Token      `json:"token"`
    Level Expression `json:"level,omitempty"`
}
```

**PHP Examples:**
```php
break;      // Break current loop
break 2;    // Break out of 2 nested loops
```

### ContinueStatement
**Type:** Statement  
**Description:** Continue to next iteration with optional level  

```go
type ContinueStatement struct {
    Token Token      `json:"token"`
    Level Expression `json:"level,omitempty"`
}
```

**PHP Examples:**
```php
continue;   // Continue current loop
continue 2; // Continue outer loop
```

### BlockStatement
**Type:** Statement  
**Description:** Block of statements enclosed in braces  

```go
type BlockStatement struct {
    Token      Token       `json:"token"`
    Statements []Statement `json:"statements"`
}
```

**PHP Examples:**
```php
{
    $x = 1;
    echo $x;
    return $x;
}
```

---

## Function and Method Declarations

### FunctionDeclaration
**Type:** Statement  
**Description:** Function definition  

```go
type FunctionDeclaration struct {
    Token      Token           `json:"token"`
    Name       *Identifier     `json:"name"`
    Parameters []*Variable     `json:"parameters"`
    Body       *BlockStatement `json:"body"`
}
```

**PHP Examples:**
```php
function greet($name) {
    echo "Hello " . $name;
}

function add($a, $b) {
    return $a + $b;
}
```

### AnonymousFunction
**Type:** Expression  
**Description:** Anonymous function/closure  

```go
type AnonymousFunction struct {
    Token      Token           `json:"token"`
    Parameters []*Variable     `json:"parameters"`
    UseClause  []*Variable     `json:"use_clause,omitempty"`
    Body       *BlockStatement `json:"body"`
}
```

**PHP Examples:**
```php
$closure = function($x) {
    return $x * 2;
};

$callback = function($data) use ($multiplier) {
    return $data * $multiplier;
};
```

### CallExpression
**Type:** Expression  
**Description:** Function or method call  

```go
type CallExpression struct {
    Token     Token        `json:"token"`
    Function  Expression   `json:"function"`
    Arguments []Expression `json:"arguments"`
}
```

**PHP Examples:**
```php
strlen($string)
$object->method($arg1, $arg2)
MyClass::staticMethod()
```

### ReturnStatement
**Type:** Statement  
**Description:** Return from function  

```go
type ReturnStatement struct {
    Token       Token      `json:"token"`
    ReturnValue Expression `json:"return_value"`
}
```

**PHP Examples:**
```php
return;           // Return without value
return $result;   // Return with value
return $a + $b;   // Return expression
```

---

## Object-Oriented Programming

### ClassDeclaration
**Type:** Statement  
**Description:** Class definition  

```go
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
```

**PHP Examples:**
```php
class User extends Person implements UserInterface {
    use LoggerTrait;
    
    public const STATUS_ACTIVE = 1;
    private $name;
    protected static $count = 0;
    
    public function __construct($name) {
        $this->name = $name;
    }
    
    public function getName() {
        return $this->name;
    }
}
```

### PropertyDeclaration
**Type:** Statement  
**Description:** Class property definition  

```go
type PropertyDeclaration struct {
    Token      Token      `json:"token"`
    Visibility string     `json:"visibility"`  // public, private, protected
    Static     bool       `json:"static"`
    Name       *Variable  `json:"name"`
    Value      Expression `json:"value,omitempty"`
}
```

**PHP Examples:**
```php
public $name;
private static $instance = null;
protected $data = [];
```

### MethodDeclaration
**Type:** Statement  
**Description:** Class method definition  

```go
type MethodDeclaration struct {
    Token      Token           `json:"token"`
    Visibility string          `json:"visibility"`
    Static     bool            `json:"static"`
    Name       *Identifier     `json:"name"`
    Parameters []*Variable     `json:"parameters"`
    Body       *BlockStatement `json:"body"`
}
```

**PHP Examples:**
```php
public function getName() {
    return $this->name;
}

private static function getInstance() {
    return self::$instance;
}
```

### InterfaceDeclaration
**Type:** Statement  
**Description:** Interface definition  

```go
type InterfaceDeclaration struct {
    Token   Token              `json:"token"`
    Name    *Identifier        `json:"name"`
    Methods []*InterfaceMethod `json:"methods"`
}
```

**PHP Examples:**
```php
interface UserInterface {
    public function getName();
    public function setName($name);
}
```

### InterfaceMethod
**Type:** Statement  
**Description:** Interface method signature  

```go
type InterfaceMethod struct {
    Token      Token       `json:"token"`
    Visibility string      `json:"visibility"`
    Name       *Identifier `json:"name"`
    Parameters []*Variable `json:"parameters"`
}
```

### TraitDeclaration
**Type:** Statement  
**Description:** Trait definition  

```go
type TraitDeclaration struct {
    Token      Token                  `json:"token"`
    Name       *Identifier            `json:"name"`
    Properties []*PropertyDeclaration `json:"properties"`
    Methods    []*MethodDeclaration   `json:"methods"`
}
```

**PHP Examples:**
```php
trait LoggerTrait {
    protected $logLevel = 'info';
    
    public function log($message) {
        echo "[{$this->logLevel}] $message";
    }
}
```

### TraitUse
**Type:** Statement  
**Description:** Use trait in class  

```go
type TraitUse struct {
    Token  Token         `json:"token"`
    Traits []*Identifier `json:"traits"`
}
```

**PHP Examples:**
```php
use LoggerTrait;
use LoggerTrait, CacheTrait;
```

### ConstantDeclaration
**Type:** Statement  
**Description:** Class constant definition  

```go
type ConstantDeclaration struct {
    Token      Token       `json:"token"`
    Visibility string      `json:"visibility"`
    Name       *Identifier `json:"name"`
    Value      Expression  `json:"value"`
}
```

**PHP Examples:**
```php
const STATUS_ACTIVE = 1;
public const MAX_SIZE = 1024;
private const SECRET_KEY = 'abc123';
```

### NewExpression
**Type:** Expression  
**Description:** Object instantiation  

```go
type NewExpression struct {
    Token     Token        `json:"token"`
    ClassName *Identifier  `json:"class_name"`
    Arguments []Expression `json:"arguments"`
}
```

**PHP Examples:**
```php
new User()
new Database($host, $user, $pass)
new \DateTime('now')
```

### ObjectAccessExpression
**Type:** Expression  
**Description:** Object property/method access (->)  

```go
type ObjectAccessExpression struct {
    Token    Token      `json:"token"`
    Object   Expression `json:"object"`
    Property Expression `json:"property"`
}
```

**PHP Examples:**
```php
$user->name
$user->getName()
$this->property
```

### StaticAccessExpression
**Type:** Expression  
**Description:** Static property/method access (::)  

```go
type StaticAccessExpression struct {
    Token    Token      `json:"token"`
    Class    Expression `json:"class"`
    Property Expression `json:"property"`
}
```

**PHP Examples:**
```php
User::$count
MyClass::CONSTANT
self::getInstance()
parent::__construct()
```

---

## Arrays and Collections

### ArrayLiteral
**Type:** Expression  
**Description:** Indexed array literal  

```go
type ArrayLiteral struct {
    Token    Token        `json:"token"`
    Elements []Expression `json:"elements"`
}
```

**PHP Examples:**
```php
[1, 2, 3]
["apple", "banana", "cherry"]
[$var1, $var2, $var3]
```

### AssociativeArrayLiteral
**Type:** Expression  
**Description:** Associative array (key-value pairs)  

```go
type AssociativeArrayLiteral struct {
    Token Token       `json:"token"`
    Pairs []ArrayPair `json:"pairs"`
}

type ArrayPair struct {
    Key   Expression `json:"key"`
    Value Expression `json:"value"`
}
```

**PHP Examples:**
```php
["name" => "John", "age" => 30]
[0 => "first", 1 => "second"]
[$key1 => $value1, $key2 => $value2]
```

### IndexExpression
**Type:** Expression  
**Description:** Array element access  

```go
type IndexExpression struct {
    Token Token      `json:"token"`
    Left  Expression `json:"left"`
    Index Expression `json:"index"`
}
```

**PHP Examples:**
```php
$array[0]
$data["key"]
$matrix[$i][$j]
```

---

## Namespace and Import System

### NamespaceDeclaration
**Type:** Statement  
**Description:** Namespace declaration  

```go
type NamespaceDeclaration struct {
    Token Token       `json:"token"`
    Name  *Identifier `json:"name"`
}
```

**PHP Examples:**
```php
namespace App\Controllers;
namespace MyProject\Utils;
```

### UseStatement
**Type:** Statement  
**Description:** Import statement with optional alias  

```go
type UseStatement struct {
    Token     Token       `json:"token"`
    Namespace *Identifier `json:"namespace"`
    Alias     *Identifier `json:"alias,omitempty"`
}
```

**PHP Examples:**
```php
use App\Models\User;
use Very\Long\Namespace\ClassName as ShortName;
use PDO;
```

---

## Exception Handling

### TryStatement
**Type:** Statement  
**Description:** Try-catch-finally block  

```go
type TryStatement struct {
    Token   Token           `json:"token"`
    Body    *BlockStatement `json:"body"`
    Catches []*CatchClause  `json:"catches"`
    Finally *BlockStatement `json:"finally,omitempty"`
}
```

**PHP Examples:**
```php
try {
    $result = riskyOperation();
} catch (DatabaseException $e) {
    logError($e);
} catch (Exception $e) {
    handleGenericError($e);
} finally {
    cleanup();
}
```

### CatchClause
**Type:** Statement  
**Description:** Catch block for exception handling  

```go
type CatchClause struct {
    Token         Token           `json:"token"`
    ExceptionType *Identifier     `json:"exception_type"`
    Variable      *Variable       `json:"variable"`
    Body          *BlockStatement `json:"body"`
}
```

**PHP Examples:**
```php
catch (Exception $e) {
    echo $e->getMessage();
}

catch (DatabaseException $dbError) {
    rollback();
}
```

### ThrowStatement
**Type:** Statement  
**Description:** Throw exception  

```go
type ThrowStatement struct {
    Token      Token      `json:"token"`
    Expression Expression `json:"expression"`
}
```

**PHP Examples:**
```php
throw new Exception("Error message");
throw $customException;
throw new InvalidArgumentException("Invalid input: " . $input);
```

---

## Advanced Features

### EchoStatement
**Type:** Statement  
**Description:** Output statement  

```go
type EchoStatement struct {
    Token  Token        `json:"token"`
    Values []Expression `json:"values"`
}
```

**PHP Examples:**
```php
echo "Hello World";
echo $name, " is ", $age, " years old";
echo $result;
```

### YieldExpression
**Type:** Expression  
**Description:** Generator yield expression  

```go
type YieldExpression struct {
    Token Token      `json:"token"`
    Key   Expression `json:"key,omitempty"`
    Value Expression `json:"value,omitempty"`
}
```

**PHP Examples:**
```php
yield $value;              // Yield value only
yield $key => $value;      // Yield key-value pair
yield;                     // Yield without value
```

### InterpolatedString
**Type:** Expression  
**Description:** String with embedded variables  

```go
type InterpolatedString struct {
    Token Token        `json:"token"`
    Parts []Expression `json:"parts"`
}
```

**PHP Examples:**
```php
"Hello $name, you are $age years old"
"The result is: {$calculation}"
```

---

## JSON Serialization

The parser provides a `ToJSON()` function that converts any AST node to JSON format:

```go
func ToJSON(node Node) ([]byte, error)
```

### Usage Example:
```go
// Parse PHP code
lexer := New(phpCode)
parser := NewParser(lexer)
program := parser.ParseProgram()

// Convert to JSON
jsonData, err := ToJSON(program)
if err != nil {
    log.Fatal(err)
}

fmt.Println(string(jsonData))
```

### Sample JSON Output:
```json
{
  "type": "Program",
  "statements": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "name": {
          "type": "Variable",
          "name": "greeting"
        },
        "value": {
          "type": "StringLiteral",
          "value": "Hello World"
        }
      }
    }
  ]
}
```

---

## Node Hierarchy Summary

```
Node (interface)
├── Statement (interface)
│   ├── ExpressionStatement
│   ├── FunctionDeclaration
│   ├── ReturnStatement
│   ├── BlockStatement
│   ├── IfStatement
│   ├── ForStatement
│   ├── WhileStatement
│   ├── ForeachStatement
│   ├── BreakStatement
│   ├── ContinueStatement
│   ├── EchoStatement
│   ├── ClassDeclaration
│   ├── PropertyDeclaration
│   ├── MethodDeclaration
│   ├── InterfaceDeclaration
│   ├── InterfaceMethod
│   ├── TraitDeclaration
│   ├── TraitUse
│   ├── ConstantDeclaration
│   ├── NamespaceDeclaration
│   ├── UseStatement
│   ├── TryStatement
│   ├── CatchClause
│   └── ThrowStatement
├── Expression (interface)
│   ├── Identifier
│   ├── Variable
│   ├── IntegerLiteral
│   ├── FloatLiteral
│   ├── StringLiteral
│   ├── BooleanLiteral
│   ├── AssignmentExpression
│   ├── InfixExpression
│   ├── PrefixExpression
│   ├── PostfixExpression
│   ├── TernaryExpression
│   ├── CallExpression
│   ├── ArrayLiteral
│   ├── AssociativeArrayLiteral
│   ├── IndexExpression
│   ├── NewExpression
│   ├── ObjectAccessExpression
│   ├── StaticAccessExpression
│   ├── AnonymousFunction
│   ├── NamespacedIdentifier
│   ├── YieldExpression
│   └── InterpolatedString
└── Program (root node)
```

This comprehensive reference covers all AST node types available in the Go PHP Parser, with their structure, usage examples, and relationships.