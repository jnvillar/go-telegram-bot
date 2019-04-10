package commands

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"go-telegram-bot/password"
	"go-telegram-bot/validator"
	"strings"
)

type CommandHandler struct {
	PasswordManager *password.PassWordManager
}

func New(pm *password.PassWordManager) *CommandHandler {
	return &CommandHandler{PasswordManager: pm}
}

func (b *CommandHandler) Master(msg *tgbotapi.Message) string {
	params := strings.Fields(msg.Text)
	if len(params) < 2 {
		return "Faltan parámetros: <masterPass>"
	}
	master := params[1]
	v, _ := validator.MinLength(master, 5)
	if !v {
		return "Contraseña maestra muy corta"
	}
	b.PasswordManager.StoreMasterPassword(msg.From.ID, msg.Text)
	return "Contraseña maestra seteada"
}

func (b *CommandHandler) Store(msg *tgbotapi.Message) string {
	params := strings.Fields(msg.Text)
	v, err := validator.LengthOfParameters(params)
	if !v {
		return err.Error()
	}
	if len(params) < 3 {
		return "Faltan parámetros: <nombre> <pass>"
	}
	err = b.PasswordManager.StorePassword(msg.From.ID, strings.ToLower(params[1]), params[2])
	if err != nil {
		return err.Error()
	}
	return "Contraseña guardada"
}

func (b *CommandHandler) Load(msg *tgbotapi.Message) string {
	params := strings.Fields(msg.Text)
	v, err := validator.LengthOfParameters(params)
	if !v {
		return err.Error()
	}
	if len(params) < 2 {
		return "Faltan parámetros: <nombre>"
	}
	pass, found, err := b.PasswordManager.LoadPassword(msg.From.ID, strings.ToLower(params[1]))
	if err != nil {
		return err.Error()
	}
	if !found {
		return "Contraseña no encontrada"
	}
	return fmt.Sprintf("Contraseña: %s", pass)
}
