package api

import (
	"errors"
	"fmt"
	"path"
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
)

func validateExecutionRequest(req ExecutionRequest, cfg Config) error {
	if req.Language == "" {
		return ErrLanguageRequired
	}

	if req.Code == "" {
		return ErrCodeRequired
	}
	if len(req.Code) > cfg.MaxCodeSize {
		return fmt.Errorf("%w: got %d bytes, max %d", ErrCodeTooLarge, len(req.Code), cfg.MaxCodeSize)
	}

	if err := validateInputFiles(req.Files, cfg); err != nil {
		return err
	}

	return nil
}

func validateInputFiles(files []InputFile, cfg Config) error {
	if len(files) > cfg.MaxFileCount {
		return fmt.Errorf("%w: got %d, max %d", ErrTooManyFiles, len(files), cfg.MaxFileCount)
	}

	totalSize := 0

	for _, file := range files {
		if err := validateInputFilePath(file.Path); err != nil {
			return err
		}

		size := len([]byte(file.Content))

		if size > cfg.MaxFileSizeBytes {
			return fmt.Errorf("%w: %s is %d bytes, max %d", ErrFileTooLarge, file.Path, size, cfg.MaxFileSizeBytes)
		}
		totalSize += size
		if totalSize > cfg.MaxTotalFileSize {
			return fmt.Errorf("%w: got %d bytes, max %d", ErrTotalFilesTooLarge, totalSize, cfg.MaxTotalFileSize)
		}
	}

	return nil
}

func validateInputFilePath(filePath string) error {
	if strings.TrimSpace(filePath) == "" {
		return fmt.Errorf("%w: path is empty", ErrInvalidFilePath)
	}

	if strings.HasPrefix(filePath, "/") {
		return fmt.Errorf("%w: absolute paths are not allowed: %s", ErrInvalidFilePath, filePath)
	}

	if strings.Contains(filePath, "\\") {
		return fmt.Errorf("%w: backslashes are not allowed: %s", ErrInvalidFilePath, filePath)
	}

	cleaned := path.Clean(filePath)

	if cleaned == "." {
		return fmt.Errorf("%w: path is empty", ErrInvalidFilePath)
	}

	if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return fmt.Errorf("%w: path traversal is not allowed: %s", ErrInvalidFilePath, filePath)
	}

	return nil
}
