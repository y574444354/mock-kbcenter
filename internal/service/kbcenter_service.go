package service

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zgsm/mock-kbcenter/i18n"
	"github.com/zgsm/mock-kbcenter/pkg/language"
)

type KBCenterMockService struct {
	baseDir string
}

func NewKBCenterMockService(baseDir string) *KBCenterMockService {
	return &KBCenterMockService{
		baseDir: baseDir,
	}
}

// GetFileContent reads file content and returns lines between startLine and endLine (inclusive)
func (s *KBCenterMockService) GetFileContent(ctx context.Context, filePath string, startLine, endLine int) ([]byte, error) {
	fullPath := filepath.Join(s.baseDir, filePath)

	content, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%s", i18n.Translate("kbcenter.file_not_found", "", map[string]interface{}{"path": fullPath}))
		}
		return nil, fmt.Errorf("%s", i18n.Translate("kbcenter.read_file_failed", "", map[string]interface{}{"error": err.Error()}))
	}

	// Split content into lines
	lines := bytes.Split(content, []byte{'\n'})

	// Validate line numbers
	if startLine < 1 || startLine > len(lines) {
		return nil, fmt.Errorf("%s", i18n.Translate("kbcenter.invalid_start_line", "", map[string]interface{}{
			"startLine": startLine,
			"maxLine":   len(lines),
		}))
	}
	if endLine > len(lines) {
		return nil, fmt.Errorf("%s", i18n.Translate("kbcenter.invalid_end_line", "", map[string]interface{}{
			"endLine": endLine,
			"maxLine": len(lines),
		}))
	}
	if endLine < 1 {
		endLine = len(lines)
	}
	if startLine > endLine {
		return nil, fmt.Errorf("%s", i18n.Translate("kbcenter.invalid_line_range", "", map[string]interface{}{
			"startLine": startLine,
			"endLine":   endLine,
		}))
	}

	// Extract requested lines
	var result [][]byte
	for i := startLine - 1; i < endLine; i++ {
		result = append(result, lines[i])
	}

	// Join lines with newlines
	return bytes.Join(result, []byte{'\n'}), nil
}

type DirectoryNode struct {
	Name     string          `json:"name"`
	Type     string          `json:"type"`
	Path     string          `json:"path"`
	Children []DirectoryNode `json:"children,omitempty"`
}

func (s *KBCenterMockService) buildDirectoryTree(basePath, relativePath string, depth, currentDepth int, includeFiles bool) (DirectoryNode, error) {
	node := DirectoryNode{
		Name: filepath.Base(basePath),
		Type: "directory",
		Path: relativePath,
	}

	if currentDepth >= depth {
		return node, nil
	}

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return node, fmt.Errorf("%s", i18n.Translate("kbcenter.read_dir_failed", "", map[string]interface{}{"path": basePath, "error": err.Error()}))
	}

	for _, entry := range entries {
		// Skip files and directories starting with '.'
		if entry.Name()[0] == '.' {
			continue
		}

		// Skip node_models directory
		if entry.IsDir() && entry.Name() == "node_modules" {
			continue
		}

		entryPath := filepath.Join(basePath, entry.Name())
		entryRelativePath := filepath.Join(relativePath, entry.Name())

		if entry.IsDir() {
			child, err := s.buildDirectoryTree(entryPath, entryRelativePath, depth, currentDepth+1, includeFiles)
			if err != nil {
				return node, err
			}
			node.Children = append(node.Children, child)
		} else if includeFiles {
			node.Children = append(node.Children, DirectoryNode{
				Name: entry.Name(),
				Type: "file",
				Path: entryRelativePath,
			})
		}
	}

	return node, nil
}

func (s *KBCenterMockService) GetFileStructure(ctx context.Context, filePath string) ([]language.FunctionInfo, error) {
	fullPath := filepath.Join(s.baseDir, filePath)

	content, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%s", i18n.Translate("kbcenter.file_not_found", "", map[string]interface{}{"path": fullPath}))
		}
		return nil, fmt.Errorf("%s", i18n.Translate("kbcenter.read_file_failed", "", map[string]interface{}{"error": err.Error()}))
	}

	// Detect language from file extension
	lang, err := language.Detect(filePath)
	if err != nil {
		return nil, fmt.Errorf("%s", i18n.Translate("language.unsupported_file_type", "", map[string]interface{}{
			"file": filePath,
		}))
	}

	return language.ExtractFunctions(lang, string(content))
}

func (s *KBCenterMockService) GetDirectoryTree(ctx context.Context, clientId, projectPath, subDir string, depth int, includeFiles bool) (interface{}, error) {
	basePath := filepath.Join(s.baseDir, subDir)

	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s", i18n.Translate("kbcenter.dir_not_found", "", map[string]interface{}{"path": basePath}))
	}

	var result struct {
		CodebaseId    string        `json:"codebaseId"`
		Name          string        `json:"name"`
		RootPath      string        `json:"rootPath"`
		DirectoryTree DirectoryNode `json:"directoryTree"`
	}

	rootNode, err := s.buildDirectoryTree(basePath, subDir, depth, 0, includeFiles)
	if err != nil {
		return nil, err
	}

	result.CodebaseId = clientId
	result.Name = projectPath
	result.RootPath = basePath
	result.DirectoryTree = rootNode

	return &result, nil
}
