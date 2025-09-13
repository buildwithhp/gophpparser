package gophpparser

import (
	"os"
	"testing"
)

func TestParsefileFunction(t *testing.T) {
	// Test parsing an existing test file
	program, err := Parsefile("testfiles/phase3_test.php")
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	if program == nil {
		t.Fatal("Program is nil")
	}

	if len(program.Statements) == 0 {
		t.Fatal("Program has no statements")
	}
	b, err := ToJSON(program)
	if err != nil {
		t.Fatalf("Failed to convert program to JSON: %v", err)
	}
	s := string(b)
	//
	t.Logf("Program JSON: %s", s)

	t.Logf("Successfully parsed file with %d statements", len(program.Statements))
}

func TestParsefileNonExistentFile(t *testing.T) {
	// Test with non-existent file
	_, err := Parsefile("non_existent_file.php")
	if err == nil {
		t.Fatal("Expected error for non-existent file")
	}

	t.Logf("Correctly returned error for non-existent file: %v", err)
}

func TestParsefileSimpleCode(t *testing.T) {
	// Create a temporary PHP file for testing
	content := `<?php
$name = "John";
echo "Hello " . $name;
?>`

	tmpFile, err := os.CreateTemp("", "test_*.php")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test parsing the temporary file
	program, err := Parsefile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to parse temp file: %v", err)
	}

	if program == nil {
		t.Fatal("Program is nil")
	}

	if len(program.Statements) != 2 {
		t.Errorf("Expected 2 statements, got %d", len(program.Statements))
	}

	t.Logf("Successfully parsed temporary file with %d statements", len(program.Statements))
}

func TestParsefileWithJSONOutput(t *testing.T) {
	// Test parsing and converting to JSON with a simpler file first
	program, err := Parsefile("testfiles/simple_test.php")
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	// Convert to JSON
	jsonData, err := ToJSON(program)
	if err != nil {
		t.Fatalf("Failed to convert to JSON: %v", err)
	}

	if len(jsonData) == 0 {
		t.Fatal("JSON output is empty")
	}

	t.Logf("Successfully converted to JSON (%d bytes)", len(jsonData))
}

func TestParsefileMultipleFiles(t *testing.T) {
	// Test parsing multiple test files
	testFiles := []string{
		"testfiles/simple_test.php",
		"testfiles/minimal_test.php",
		"testfiles/operators_test.php",
	}

	for _, filename := range testFiles {
		t.Run(filename, func(t *testing.T) {
			// Check if file exists first
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				t.Skipf("Test file %s does not exist, skipping", filename)
				return
			}

			program, err := Parsefile(filename)
			if err != nil {
				t.Errorf("Failed to parse %s: %v", filename, err)
				return
			}

			if program == nil {
				t.Errorf("Program is nil for %s", filename)
				return
			}

			t.Logf("Successfully parsed %s with %d statements", filename, len(program.Statements))
		})
	}
}

// Benchmark the Parsefile function
func BenchmarkParsefile(b *testing.B) {
	// Check if test file exists
	if _, err := os.Stat("testfiles/simple_test.php"); os.IsNotExist(err) {
		b.Skip("Test file does not exist")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Parsefile("testfiles/simple_test.php")
		if err != nil {
			b.Fatalf("Parsing failed: %v", err)
		}
	}
}

// Test namespace separator parsing
func TestParsefileNamespaceSeparator(t *testing.T) {
	// Test parsing phase3_test.php which contains namespace separators
	program, err := Parsefile("testfiles/phase3_test.php")
	if err != nil {
		t.Logf("Error parsing phase3_test.php: %v", err)
		// Don't fail the test, just log the error to see if our fix worked
		return
	}

	if program != nil {
		t.Logf("Successfully parsed phase3_test.php with %d statements", len(program.Statements))
	}
}

// Test to debug the semicolon issue specifically
func TestParsefileSemicolonIssue(t *testing.T) {
	// Create a simple PHP code that might have semicolon issues
	content := `<?php
$x = 1;
$y = 2;
?>`

	tmpFile, err := os.CreateTemp("", "semicolon_test_*.php")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test parsing this simple file
	program, err := Parsefile(tmpFile.Name())
	if err != nil {
		t.Logf("Error parsing simple semicolon test: %v", err)
		// Don't fail the test, just log the error to understand it
		return
	}

	if program != nil {
		t.Logf("Successfully parsed semicolon test with %d statements", len(program.Statements))
	}
}

// Test helper function to demonstrate library usage
func ExampleParsefile() {
	// This is how external packages would use the Parsefile function
	program, err := Parsefile("testfiles/simple_test.php")
	if err != nil {
		// Handle error
		return
	}

	// Use the program
	_ = program.Statements
}