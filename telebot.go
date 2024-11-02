package main

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NewsKeyword struct {
	Text string
}

type TeleBot struct {
	bot      *tgbotapi.BotAPI
	chatID   int64
	keywords []NewsKeyword
	log      func(string)
}

func CreateTeleBot(token string, chatID int64, log func(string)) *TeleBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log("Failed to initialize bot: " + err.Error())
		return nil
	}
	return &TeleBot{
		bot:      bot,
		chatID:   chatID,
		keywords: []NewsKeyword{},
		log:      log,
	}
}

func (t *TeleBot) SetKeywords(keywords []NewsKeyword) {
	t.keywords = keywords
}

func (t *TeleBot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// updates := t.bot.GetUpdatesChan(u)
	t.log("Bot started successfully")

	// msg := tgbotapi.NewMessage(t.chatID, "hihi")
	// t.bot.Send(msg)

	// for update := range updates {
	// 	t.log(update.Message.Command())
	// 	if update.Message != nil && update.Message.IsCommand() {
	// 		switch update.Message.Command() {
	// 		case "Keywords":
	// 			t.SendKeywords(update.Message.Chat.ID)
	// 		}
	// 	}
	// }
}

func (t *TeleBot) SendKeywords(chatID int64) {
	var keywordTexts []string
	for _, keyword := range t.keywords {
		keywordTexts = append(keywordTexts, "'"+keyword.Text+"'")
	}
	msg := tgbotapi.NewMessage(chatID, "등록된 keyword: ["+strings.Join(keywordTexts, ", ")+"]")
	t.bot.Send(msg)
}
