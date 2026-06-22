package execution

type JavaScriptRunetime struct{}

func (JavaScriptRunetime) Image() string {
	return "node:22-alpine"
}

func (JavaScriptRunetime) Command() []string {
	return []string{"node"}
}
