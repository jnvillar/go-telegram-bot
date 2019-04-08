package main

import (
	"bot/commands"
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var CommandHandlers map[string]commands.CommandHandler

func main() {
	CommandHandlers = loadHandlers()
	var token = os.Getenv("TELEGRAMTOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		logMessage(update.Message)

		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			msg := handleCommand(update.Message)
			bot.Send(msg)
		} 
	}
}

func logMessage(message *tgbotapi.Message){
	log.Printf("[%s] %s", message.From.UserName, message.Text)
}

func handleCommand(message *tgbotapi.Message) *tgbotapi.MessageConfig{
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	command := message.Command()
	handler, ok := CommandHandlers[command]
	if !ok{
		msg.Text = "I don't know that command"
	}else{
		msg.Text, _ = handler.Handle(message.Text)
	}
	return &msg
}

func loadHandlers() map[string]commands.CommandHandler{
	return map[string]commands.CommandHandler{
		"help": &commands.HelpHandler{},
	}
}