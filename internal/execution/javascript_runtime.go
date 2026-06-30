package execution

type JavaScriptRuntime struct{}

func (JavaScriptRuntime) Image() string {
	return "node:22-alpine"
}

func (JavaScriptRuntime) Command() []string {
	return []string{"node"}
}
