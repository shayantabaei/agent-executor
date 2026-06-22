package execution

import (
	"fmt"
)

type UnsupportedLanguageError struct {
	Language string
}

func (e UnsupportedLanguageError) Error() string {
	return fmt.Sprintf("unsupported language: %s", e.Language)
}

type Runtime interface {
	Image() string
	Command() []string
}

func runtimeForLanguage(language string) (Runtime, error) {
	switch language {
	case "python":
		return PythonRuntime{}, nil
	case "javascript":
		return JavaScriptRunetime{}, nil
	default:
		return nil, UnsupportedLanguageError{Language: language}
	}
}
