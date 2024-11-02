package main

import (
	"log"
	"os"
	"strconv"
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
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := GetChatID()
	teleBot := CreateTeleBot(telegramBotToken, chatID, logWrapper)
	if teleBot == nil {
		logger.Fatal("Failed to create TeleBot instance")
	}

	// Initialize the craler
	crawler := CreateYahooFinanceCrawler(logWrapper)

	// Example usage
	crawler.AddKeyword("Korea")
	crawler.AddKeyword("Trump")
	crawler.AddKeyword("Dean")
	// crawler.CrawlYahooNews()

	teleBot.SetKeywords(crawler.keywords)
	teleBot.Start()

	// Keep the main function running
	select {}
}

func logWrapper(logger *log.Logger) func(string) {
	return func(message string) {
		logger.Println(message)
	}
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
