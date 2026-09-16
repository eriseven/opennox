// Package localization selects language resources with a fallback to the original game data.
package localization

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opennox/libs/ifs"
	"github.com/opennox/libs/strman"
)

// Directory returns the selected locales/<language> directory, or an empty path
// when the original game language is selected.
func Directory(dataDir, language string) (string, error) {
	if language == "" {
		return "", nil
	}
	if language == "." || language == ".." || strings.ContainsAny(language, `/\:`) {
		return "", fmt.Errorf("invalid language %q: expected a folder name such as zh-CN", language)
	}
	dir := filepath.Join(dataDir, "locales", language)
	info, err := ifs.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("cannot load language %q: %w", language, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("language path is not a directory: %s", dir)
	}
	return dir, nil
}

// Resolve tries all supported filenames in the locale before the base directory.
// Only missing files fall back; unreadable files and directories are errors.
func Resolve(dataDir, localeDir string, names ...string) (string, error) {
	dirs := []string{dataDir}
	if localeDir != "" {
		dirs = append([]string{localeDir}, dirs...)
	}
	for _, dir := range dirs {
		for _, name := range names {
			path := filepath.Join(dir, name)
			info, err := ifs.Stat(path)
			if os.IsNotExist(err) {
				continue
			} else if err != nil {
				return "", err
			}
			if info.IsDir() {
				return "", fmt.Errorf("expected a resource file: %s", path)
			}
			return path, nil
		}
	}
	return "", fmt.Errorf("resources %v in %v: %w", names, dirs, os.ErrNotExist)
}

// LoadStrings selects a complete string table. It does not merge individual entries.
func LoadStrings(sm *strman.StringManager, dataDir, localeDir string) error {
	path, err := Resolve(dataDir, localeDir, "nox.csf.json", "nox.csf")
	if err != nil {
		return err
	}
	if strings.HasSuffix(path, ".json") {
		err = sm.ReadJSON(path)
	} else {
		err = sm.ReadCSF(path)
	}
	if err != nil {
		return fmt.Errorf("cannot load strings %q: %w", path, err)
	}
	return nil
}
