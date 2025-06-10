package language

import (
	"errors"
	"path/filepath"

	"github.com/zgsm/mock-kbcenter/config"
)

func Detect(filePath string) (string, error) {
	cfg := config.GetConfig()
	ext := filepath.Ext(filePath)
	if lang, ok := cfg.LanguageMapping[ext]; ok {
		return lang, nil
	}
	return "", errors.New("not supported file type")
}
