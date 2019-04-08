package commands

type HelpHandler struct {}

func (b *HelpHandler) Handle(msg string) (string, error) {
	return "fuck you", nil
}
