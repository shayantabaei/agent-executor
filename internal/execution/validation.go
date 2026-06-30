package execution

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

var (
	ErrLanguageRequired   = errors.New("language is required")
	ErrCodeRequired       = errors.New("code is required")
	ErrCodeTooLarge       = errors.New("code exceeds maximum size")
	ErrTooManyFiles       = errors.New("too many files")
	ErrFileTooLarge       = errors.New("file too large")
	ErrTotalFilesTooLarge = errors.New("total file size too large")
	ErrInvalidFilePath    = errors.New("invalid file path")
	ErrInvalidRequest     = errors.New("invalid execution request")
)

func validateRequest(req Request, config ServiceConfig) error {
	if strings.TrimSpace(req.Language) == "" {
		return fmt.Errorf("%w: language is required", ErrLanguageRequired)
	}

	if strings.TrimSpace(req.Code) == "" {
		return fmt.Errorf("%w: code is required", ErrCodeRequired)
	}

	if len(req.Code) > config.MaxCodeSize {
		return fmt.Errorf("%w: code exceeds maximum size", ErrCodeTooLarge)
	}

	if len(req.Files) > config.MaxFileCount {
		return fmt.Errorf("%w: too many files", ErrTooManyFiles)
	}

	totalFileSize := 0
	seenPaths := make(map[string]struct{})

	for _, file := range req.Files {
		if strings.TrimSpace(file.Path) == "" {
			return fmt.Errorf("%w: file path is required", ErrInvalidFilePath)
		}

		if strings.Contains(file.Path, "\\") {
			return fmt.Errorf("%w: file path must use forward slashes", ErrInvalidFilePath)
		}

		if filepath.IsAbs(file.Path) {
			return fmt.Errorf("%w: absolute file paths are not allowed", ErrInvalidFilePath)
		}

		cleanPath := filepath.Clean(file.Path)

		if cleanPath == "." || strings.HasPrefix(cleanPath, "..") {
			return fmt.Errorf("%w: file path traversal is not allowed", ErrInvalidFilePath)
		}

		if _, exists := seenPaths[cleanPath]; exists {
			return fmt.Errorf("%w: duplicate file path %q", ErrInvalidFilePath, cleanPath)
		}

		seenPaths[cleanPath] = struct{}{}

		fileSize := len(file.Content)

		if fileSize > config.MaxFileSizeBytes {
			return fmt.Errorf("%w: file %q exceeds maximum size", ErrFileTooLarge, cleanPath)
		}

		totalFileSize += fileSize

		if totalFileSize > config.MaxTotalFileSize {
			return fmt.Errorf("%w: total file size exceeds maximum size", ErrTotalFilesTooLarge)
		}
	}

	return nil
}

func IsValidationError(err error) bool {
	return errors.Is(err, ErrLanguageRequired) ||
		errors.Is(err, ErrCodeRequired) ||
		errors.Is(err, ErrCodeTooLarge) ||
		errors.Is(err, ErrTooManyFiles) ||
		errors.Is(err, ErrFileTooLarge) ||
		errors.Is(err, ErrTotalFilesTooLarge) ||
		errors.Is(err, ErrInvalidFilePath) ||
		errors.Is(err, ErrInvalidRequest)
}
