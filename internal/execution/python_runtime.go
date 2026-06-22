package execution

type PythonRuntime struct{}

func (PythonRuntime) Image() string {
	return "python:3.12-alpine"
}

func (PythonRuntime) Command() []string {
	return []string{"python", "-"}
}
