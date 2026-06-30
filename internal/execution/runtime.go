package execution

import (
	"fmt"
	"strings"
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

func SupportedLanguages() []string {
	return []string{"python", "javascript"}
}

func runtimeForLanguage(language string) (Runtime, error) {
	language = strings.ToLower(strings.TrimSpace(language))

	switch language {
	case "python":
		return PythonRuntime{}, nil
	case "javascript":
		return JavaScriptRuntime{}, nil
	default:
		return nil, UnsupportedLanguageError{Language: language}
	}
}
