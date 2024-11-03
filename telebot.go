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
	// 동작 안하는 듯..?
	cmdConfig := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command: "keywords",
			Description: "List all keywords",
		},
		tgbotapi.BotCommand{
			Command: "add",
			Description: "Add a new keyword",
		},
		tgbotapi.BotCommand{
			Command: "remove",
			Description: "Remove a new keyword",
		},
	)

	t.bot.Send(cmdConfig)
	// t.bot.Request(cmdConfig)
	
	
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.bot.GetUpdatesChan(u)
	t.log("Bot started successfully")
	t.HandleCommands(updates)
}

func (t *TeleBot) SendMessage(message string) {
	msg := tgbotapi.NewMessage(t.chatID, message)
	t.bot.Send(msg)
}

func (t *TeleBot) HandleCommands(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if !update.Message.IsCommand() {
			continue
		}

		switch update.Message.Command() {
		case "keywords":
			t.GetKeywords()
		case "add":
			t.AddKeyword(update.Message.CommandArguments())
		case "remove":
			t.RemoveKeyword(update.Message.CommandArguments())
		default:
			t.SendMessage("등록되지 않은 커맨드입니다.")
		}
	}
}

func (t *TeleBot) GetKeywords() {
	var keywordTexts []string
	for _, keyword := range t.keywords {
		keywordTexts = append(keywordTexts, "'"+keyword.Text+"'")
	}
	msgStr := "등록된 keyword: ["+strings.Join(keywordTexts, ", ")+"]"
	t.SendMessage(msgStr)
}

func (t *TeleBot) AddKeyword(keyword string) {
	if keyword == "" {
		t.SendMessage("키워드를 입력해주세요")
		return
	}
	for _, k := range t.keywords {
		if k.Text == keyword {
			t.SendMessage("이미 등록된 키워드입니다")
			return
		}
	}
	t.keywords = append(t.keywords, NewsKeyword{Text: keyword})
	t.SendMessage("키워드가 등록되었습니다")
}

func (t *TeleBot) RemoveKeyword(keyword string) {
	if keyword == "" {
		t.SendMessage("키워드를 입력해주세요")
		return
	}
	for _, k := range t.keywords {
		if k.Text == keyword {
			
			t.SendMessage("키워드가 제거되었습니다.")
			return
		}
	}
}