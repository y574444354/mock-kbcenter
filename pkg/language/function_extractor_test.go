package language

import (
	"strings"
	"testing"
)

func TestGetFunctionName(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		code     string
		expected string
		wantErr  bool
	}{
		{
			name: "Go function",
			lang: "go",
			code: `package main

func helloWorld() {
	fmt.Println("Hello, World!")
}`,
			expected: "helloWorld",
			wantErr:  false,
		},
		{
			name: "JavaScript function",
			lang: "javascript",
			code: `function greet(name) {
	return "Hello, " + name;
}`,
			expected: "greet",
			wantErr:  false,
		},
		{
			name: "Python function",
			lang: "python",
			code: `def calculate_sum(a, b):
	return a + b`,
			expected: "calculate_sum",
			wantErr:  false,
		},
		{
			name:    "Unsupported language",
			lang:    "ruby",
			code:    `def hello; end`,
			wantErr: true,
		},
		{
			name:    "Invalid code",
			lang:    "go",
			code:    `invalid code`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFunctionName(tt.lang, tt.code)
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if got != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestExtractFunctions_Go(t *testing.T) {
	code := `package main

func add(a, b int) int {
	return a + b
}

func (p *Person) greet() string {
	return "Hello, " + p.Name
}`

	functions, err := ExtractFunctions("go", code)
	if err != nil {
		t.Fatalf("ExtractFunctions failed: %v", err)
	}

	if len(functions) != 2 {
		t.Fatalf("Expected 2 functions, got %d", len(functions))
	}

	// Verify first function
	if !strings.Contains(functions[0].Code, "func add(a, b int) int") {
		t.Error("First function code mismatch")
	}
	if functions[0].StartLine != 3 || functions[0].EndLine != 5 {
		t.Errorf("First function line range mismatch: got %d-%d", functions[0].StartLine, functions[0].EndLine)
	}

	// Verify second function
	if !strings.Contains(functions[1].Code, "func (p *Person) greet() string") {
		t.Error("Second function code mismatch")
	}
	if functions[1].StartLine != 7 || functions[1].EndLine != 9 {
		t.Errorf("Second function line range mismatch: got %d-%d", functions[1].StartLine, functions[1].EndLine)
	}
}

func TestExtractFunctions_JavaScript(t *testing.T) {
	code := `function add(a, b) {
	return a + b;
}

const multiply = (a, b) => {
	return a * b;
}`

	functions, err := ExtractFunctions("javascript", code)
	if err != nil {
		t.Fatalf("ExtractFunctions failed: %v", err)
	}

	if len(functions) != 2 {
		t.Fatalf("Expected 2 functions, got %d", len(functions))
	}
}

func TestExtractFunctions_Python(t *testing.T) {
	code := `def add(a, b):
	return a + b

class Calculator:
	def multiply(self, a, b):
		return a * b`

	functions, err := ExtractFunctions("python", code)
	if err != nil {
		t.Fatalf("ExtractFunctions failed: %v", err)
	}

	if len(functions) != 2 {
		t.Fatalf("Expected 2 functions, got %d", len(functions))
	}
}

func TestExtractFunctions_UnsupportedLanguage(t *testing.T) {
	_, err := ExtractFunctions("unknown", "some code")
	if err == nil {
		t.Error("Expected error for unsupported language")
	}
}

func TestExtractFunctions_EmptyCode(t *testing.T) {
	functions, err := ExtractFunctions("go", "")
	if err != nil {
		t.Fatalf("ExtractFunctions failed: %v", err)
	}
	if len(functions) != 0 {
		t.Errorf("Expected 0 functions, got %d", len(functions))
	}
}
