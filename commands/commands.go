package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"go-telegram-bot/password"
	"go-telegram-bot/validator"
	"strings"
)


type CommandHandler struct {
	PasswordManager *password.PassWordManager
}

func New(pm *password.PassWordManager) *CommandHandler{
	return &CommandHandler{PasswordManager: pm}
}

func (b *CommandHandler) Master(msg *tgbotapi.Message) string {
	v, err := validator.Length(msg.Text)
	if !v {
		return err
	}
	b.PasswordManager.StoreMasterPassword(msg.From.ID, msg.Text)
	return "Contraseña mestra seteada"
}

func (b *CommandHandler) Store(msg *tgbotapi.Message) string {
	params := strings.Split(msg.Text, "")
	v, err := validator.LenghOfParameters(params)
	if !v{
		return err
	}
	if len(params) < 2 {
		return "Faltan parametros: <nombre> <pass>"
	}
	b.PasswordManager.StorePassword(msg.From.ID, params[0], params[1])
	return "Contraseña guardada"
}

func (b *CommandHandler) Load(msg *tgbotapi.Message) string {
	params := strings.Split(msg.Text, "")
	v, err := validator.LenghOfParameters(params)
	if !v{
		return err
	}
	if len(params) < 1 {
		return "Faltan parametros: <nombre>"
	}
	pass, ok := b.PasswordManager.LoadPassword(msg.From.ID, params[0])
	if !ok{
		return "Contraseña no encontrada"
	}
	return pass
}