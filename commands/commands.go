package commands

type CommandHandler interface {
	Handle(msg string) (string, error)
}