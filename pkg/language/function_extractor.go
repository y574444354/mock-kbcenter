package language

import (
	"context"
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	tree_sitter_c "github.com/smacker/go-tree-sitter/c"
	tree_sitter_cpp "github.com/smacker/go-tree-sitter/cpp"
	tree_sitter_go "github.com/smacker/go-tree-sitter/golang"
	tree_sitter_java "github.com/smacker/go-tree-sitter/java"
	tree_sitter_javascript "github.com/smacker/go-tree-sitter/javascript"
	tree_sitter_php "github.com/smacker/go-tree-sitter/php"
	tree_sitter_python "github.com/smacker/go-tree-sitter/python"
	tree_sitter_ruby "github.com/smacker/go-tree-sitter/ruby"
	tree_sitter_tsx "github.com/smacker/go-tree-sitter/typescript/tsx"
	tree_sitter_typescript "github.com/smacker/go-tree-sitter/typescript/typescript"
	"github.com/zgsm/mock-kbcenter/config"
	"github.com/zgsm/mock-kbcenter/i18n"
)

func getLanguage(lang string) (*sitter.Language, error) {
	// TODO: Can add logic to load languages from .so files in the future
	switch strings.ToLower(lang) {
	case "go":
		return tree_sitter_go.GetLanguage(), nil
	case "javascript", "js":
		return tree_sitter_javascript.GetLanguage(), nil
	case "tsx":
		return tree_sitter_tsx.GetLanguage(), nil
	case "typescript", "ts":
		return tree_sitter_typescript.GetLanguage(), nil
	case "python", "py":
		return tree_sitter_python.GetLanguage(), nil
	case "java":
		return tree_sitter_java.GetLanguage(), nil
	case "php":
		return tree_sitter_php.GetLanguage(), nil
	case "ruby", "rb":
		return tree_sitter_ruby.GetLanguage(), nil
	case "c":
		return tree_sitter_c.GetLanguage(), nil
	case "cpp", "cxx", "cc", "hpp":
		return tree_sitter_cpp.GetLanguage(), nil
	default:
		return nil, fmt.Errorf("%s", i18n.Translate("language.unsupported", "", map[string]interface{}{
			"language": lang,
		}))
	}
}

// FunctionInfo contains function code and its location information
type FunctionInfo struct {
	Code      string // Function code content
	StartLine int    // Start line number
	EndLine   int    // End line number
}

// GetFunctionName extracts function name from code using tree-sitter
func GetFunctionName(lang string, code string) (string, error) {
	parser := sitter.NewParser()
	defer parser.Close()

	language, err := getLanguage(lang)
	if err != nil {
		return "", err
	}

	parser.SetLanguage(language)
	tree, err := parser.ParseCtx(context.Background(), nil, []byte(code))
	if err != nil {
		return "", err
	}

	rootNode := tree.RootNode()

	// Query pattern depends on language
	var queryPattern string
	switch strings.ToLower(lang) {
	case "go":
		queryPattern = "(function_declaration name: (identifier) @name)"
	case "javascript", "typescript":
		queryPattern = "(function_declaration name: (identifier) @name)"
	case "python":
		queryPattern = "(function_definition name: (identifier) @name)"
	default:
		return "", fmt.Errorf("unsupported language: %s", lang)
	}

	query, err := sitter.NewQuery([]byte(queryPattern), language)
	if err != nil {
		return "", err
	}
	defer query.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()
	qc.Exec(query, rootNode)

	match, ok := qc.NextMatch()
	if !ok {
		return "", fmt.Errorf("no function name found in code")
	}

	for _, capture := range match.Captures {
		if capture.Node.Type() == "identifier" {
			return capture.Node.Content([]byte(code)), nil
		}
	}

	return "", fmt.Errorf("no function name found in code")
}

// FormatCodeWithLineNumbers formats code segment with line numbers
// Parameters:
//   - code: Code content
//   - startLine: Start line number (for validation)
//   - endLine: End line number (for validation)
//   - relativeStart: Relative start line number (as offset for the whole code segment)
//
// Returns:
//   - Formatted code text with line numbers
func FormatCodeWithLineNumbers(code string, startLine, endLine, relativeStart int) string {
	lines := strings.Split(code, "\n")
	var builder strings.Builder

	// Calculate max digits for line number alignment
	maxDigits := len(fmt.Sprintf("%d", endLine+relativeStart-1))

	for i, line := range lines {
		// Line number = relative start + current line index
		lineNum := relativeStart + startLine + i - 1
		// Format as "line number | code"
		lineStr := fmt.Sprintf("%*d | %s\n", maxDigits, lineNum, line)
		builder.WriteString(lineStr)
	}

	return builder.String()
}

// ExtractFunctions extracts function definitions from source code
// Parameters:
//   - lang: Language type (go, javascript, python)
//   - content: Source code content
//
// Returns:
//   - Slice of function info (containing code content and line range)
//   - Error information
func ExtractFunctions(lang string, content string) ([]FunctionInfo, error) {
	parser := sitter.NewParser()
	defer parser.Close()

	language, err := getLanguage(lang)
	if err != nil {
		return nil, err
	}

	if language == nil {
		return nil, fmt.Errorf("%s", i18n.Translate("language.invalid_language", "", map[string]interface{}{
			"language": lang,
		}))
	}
	parser.SetLanguage(language)
	tree, err := parser.ParseCtx(context.Background(), nil, []byte(content))
	if err != nil {
		return nil, err
	}

	rootNode := tree.RootNode()
	var functions []FunctionInfo

	// Get query pattern from config
	cfg := config.GetConfig()

	queryPattern := cfg.LanguageQueries[strings.ToLower(lang)]

	// Also used in unit tests
	if queryPattern == "" {
		// If not defined in config, use default query pattern
		switch strings.ToLower(lang) {
		case "go":
			queryPattern = "(function_declaration) @func\n(method_declaration) @func"
		case "javascript", "typescript", "tsx":
			queryPattern = "(function_declaration) @func\n(arrow_function) @func"
		case "python":
			queryPattern = "(function_definition) @func"
		default:
			queryPattern = "(function_declaration) @func"
		}
	}

	query, err := sitter.NewQuery([]byte(queryPattern), language)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", i18n.Translate("language.query_error", "", nil), err)
	}
	defer query.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()
	qc.Exec(query, rootNode)

	for {
		match, ok := qc.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			node := capture.Node
			start := node.StartByte()
			end := node.EndByte()
			startPoint := node.StartPoint()
			endPoint := node.EndPoint()
			functions = append(functions, FunctionInfo{
				Code:      string(content[start:end]),
				StartLine: int(startPoint.Row) + 1, // Line numbers start from 1
				EndLine:   int(endPoint.Row) + 1,
			})
		}
	}

	return functions, nil
}
