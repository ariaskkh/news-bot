package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	// Set up Logging
	logFile, err := os.OpenFile("app.log", os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	logWrapper := logWrapper(logger)
	// Initialize the bot
	telegramBotToken := GetTeleBotToken()
	chatID := GetChatID()
	teleBot := CreateTeleBot(telegramBotToken, chatID, logWrapper)
	if teleBot == nil {
		logger.Fatal("Failed to create TeleBot instance")
	}

	message := make(chan string)
	go func() {
		for msg := range message {
			teleBot.SendMessage(msg)
		}
	}()
	// Initialize the craler
	crawler := CreateYahooFinanceCrawler(logWrapper, message, &teleBot.keywords)
	if crawler == nil {
		logger.Fatal("Failed to create crawler instance")
	}

	// Example usage
	// crawler.AddKeyword("Korea")
	// crawler.AddKeyword("Trump")
	// crawler.AddKeyword("Dean")

	// TODO: 저장되어 있는 keywords set
	// teleBot.SetKeywords()
	
	go teleBot.Start()
	go func() {
		for {
			log.Println("crawling 시작")
			crawler.CrawlYahooNews()
			time.Sleep(3 * time.Second)
		}
	}()

	// Keep the main function running
	select {}
}

// 다양한 로그 원하는 경우 클래스로 만들기
func logWrapper(logger *log.Logger) func(string) {
	return func(message string) {
		logger.Println(message)
	}
}

func GetTeleBotToken() string {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is not set")
	}
	return token;
}

func GetChatID() int64 {
	chatIDStr := os.Getenv("CHAT_ID")
	if chatIDStr == "" {
		log.Fatal("CHAT_ID environment variable is not set")
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Failed to convert CHAT_ID to int64: %v", err)
	}
	return chatID
}
