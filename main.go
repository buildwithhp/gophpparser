package gophpparser

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Parsefile parses the given PHP file and returns the parsed program
// and any errors encountered during parsing. If the file does not
// exist, it returns an error with a message indicating the file
// does not exist.
func Parsefile(filepath string) (*Program, error) {
	// Check if the file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("File '%s' does not exist", filepath)
	}

	// Read the contents of the file
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("Error reading file '%s': %v", filepath, err)
	}

	// Convert the content to a string
	input := string(content)

	// Create a lexer with the input string
	lexer := New(input)

	// Create a parser with the lexer
	parser := NewParser(lexer)

	// Parse the program
	program := parser.ParseProgram()

	// Check if there are any parser errors
	if len(parser.Errors()) != 0 {
		// If there are errors, construct an error string
		errStr := fmt.Sprintf("Parser errors for file '%s':\n", filepath)
		return nil, fmt.Errorf("%s%s", errStr, strings.Join(parser.Errors(), "\n"))
	}

	// Return the parsed program and nil for the error
	return program, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: php-parser <file.php>")
		fmt.Println("       php-parser -h | --help")
		os.Exit(1)
	}

	arg := os.Args[1]

	if arg == "-h" || arg == "--help" {
		printHelp()
		return
	}

	filename := arg

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatalf("File '%s' does not exist", filename)
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file '%s': %v", filename, err)
	}

	input := string(content)

	lexer := New(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		fmt.Fprintf(os.Stderr, "Parser errors:\n")
		for _, err := range parser.Errors() {
			fmt.Fprintf(os.Stderr, "  - %s\n", err)
		}
		os.Exit(1)
	}

	jsonOutput, err := ToJSON(program)
	if err != nil {
		log.Fatalf("Error converting AST to JSON: %v", err)
	}

	fmt.Println(string(jsonOutput))
}

func printHelp() {
	fmt.Println("PHP Parser to JSON")
	fmt.Println("")
	fmt.Println("A command-line tool that parses PHP code and outputs the Abstract Syntax Tree (AST) in JSON format.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  php-parser <file.php>    Parse a PHP file and output JSON")
	fmt.Println("  php-parser -h, --help    Show this help message")
	fmt.Println("")
	fmt.Println("Example:")
	fmt.Println("  php-parser test.php")
	fmt.Println("")
	fmt.Println("The parser supports basic PHP constructs including:")
	fmt.Println("  - Variables ($var)")
	fmt.Println("  - Functions")
	fmt.Println("  - Control structures (if/else)")
	fmt.Println("  - Expressions and operators")
	fmt.Println("  - Echo statements")
	fmt.Println("  - String and numeric literals")
}
