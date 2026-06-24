package execution

import (
	"os"
	"path/filepath"
)

type workspace struct {
	path string
}

func createWorkspace(files []InputFile) (*workspace, error) {
	dir, err := os.MkdirTemp("", "agent-executor-*")

	if err != nil {
		return nil, err
	}

	ws := &workspace{path: dir}

	if err := ws.writeFiles(files); err != nil {
		_ = ws.cleanup()
		return nil, err
	}

	return ws, nil
}

func (w *workspace) writeFiles(files []InputFile) error {
	for _, file := range files {
		targetPath := filepath.Join(w.path, filepath.FromSlash(file.Path))

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		if err := os.WriteFile(targetPath, []byte(file.Content), 0644); err != nil {
			return err
		}
	}

	return nil
}

func (w *workspace) cleanup() error {
	return os.RemoveAll(w.path)
}
