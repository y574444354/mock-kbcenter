package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zgsm/mock-kbcenter/i18n"
)

type KBCenterMockService struct {
	baseDir string
}

func NewKBCenterMockService(baseDir string) *KBCenterMockService {
	return &KBCenterMockService{
		baseDir: baseDir,
	}
}

func (s *KBCenterMockService) GetFileContent(ctx context.Context, filePath string, startLine, endLine int) ([]byte, error) {
	fullPath := filepath.Join(s.baseDir, filePath)

	content, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%s", i18n.Translate("kbcenter.file_not_found", "", map[string]interface{}{"path": fullPath}))
		}
		return nil, fmt.Errorf("%s", i18n.Translate("kbcenter.read_file_failed", "", map[string]interface{}{"error": err.Error()}))
	}

	return content, nil
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

func (s *KBCenterMockService) GetDirectoryTree(ctx context.Context, clientId, projectPath, subDir string, depth int, includeFiles bool) (interface{}, error) {
	basePath := filepath.Join(s.baseDir, subDir)

	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s", i18n.Translate("kbcenter.dir_not_found", "", map[string]interface{}{"path": basePath}))
	}

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			CodebaseId    string        `json:"codebaseId"`
			Name          string        `json:"name"`
			RootPath      string        `json:"rootPath"`
			DirectoryTree DirectoryNode `json:"directoryTree"`
		} `json:"data"`
	}

	rootNode, err := s.buildDirectoryTree(basePath, subDir, depth, 0, includeFiles)
	if err != nil {
		return nil, err
	}

	result.Code = 0
	result.Message = "success"
	result.Data.CodebaseId = clientId
	result.Data.Name = projectPath
	result.Data.RootPath = basePath
	result.Data.DirectoryTree = rootNode

	return &result, nil
}
