package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"512b.it/godss/src/chart"
	"512b.it/godss/src/dss"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type ReplyData struct {
	Head   string
	Values []string
	Repeat bool
}

var replyManager = make(map[int64]*ReplyData) // stores state for each user (chat_id)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Load environment variable for TOKEN
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN not found in environment")
	}

	fmt.Printf("Token %s...%s\n", token[:3], token[len(token)-3:])

	// Initialize the bot
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	// Set bot to poll updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// Handle each update (message, command, etc.)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		text := update.Message.Text

		switch {
		case strings.HasPrefix(text, "/help"):
			handleStart(bot, chatID)
		case strings.HasPrefix(text, "/start"):
			handleStart(bot, chatID)
		case strings.HasPrefix(text, "/pie"):
			handlePieStart(bot, chatID)
		case replyManager[chatID] != nil:
			handleReply(bot, chatID, text)
		}
	}
}

func handleStart(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Hello, I'm a @dsspiebot!\n"+
		"Send me a phrase and some different conclusion of that phrase and I will tell you the online popularity.\nExample:\n"+
		"/pie\n   I like\n   foods\n   cats")
	bot.Send(msg)
}

func handlePieStart(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Please send the statement.")
	bot.Send(msg)

	// Initialize state for this chat
	replyManager[chatID] = &ReplyData{
		Values: []string{},
		Repeat: false,
	}
}

func handleReply(bot *tgbotapi.BotAPI, chatID int64, text string) {
	data := replyManager[chatID]

	if data.Head == "" { // First reply
		data.Head = text
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Send me the first option for (<i>\"%s\"</i>)", data.Head))
		msg.ParseMode = "HTML"
		bot.Send(msg)
		return
	}

	if text != "/done" { // Collecting options
		data.Values = append(data.Values, text)
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Send me the next option for (<i>\"%s\"</i>)\nSend /done when you are done.", data.Head))
		msg.ParseMode = "HTML"
		bot.Send(msg)
	} else { // User finished input
		samples := "\n- " + strings.Join(data.Values, "\n- ")
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Elaborating plot for <i>\"%s\"</i>\n%s", data.Head, samples))
		msg.ParseMode = "HTML"
		bot.Send(msg)

		// Call DSS function (replace with actual logic)
		go generateChart(bot, chatID, data.Head, data.Values)
		delete(replyManager, chatID) // Remove state after completion
	}
}

func generateChart(bot *tgbotapi.BotAPI, chatID int64, head string, values []string) {
	var err error
	// Simulate chart creation and response (replace with real implementation)
	bot.Send(tgbotapi.NewChatAction(chatID, "upload_photo"))

	counter := dss.Dss{}

	println("Generating chart for", head, strings.Join(values, ", "))
	var results map[string]int
	if results, err = counter.CountEvents(head, values); err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Unable to complete the operation for <i>\"%s\"</i>.", head))
		msg.ParseMode = "HTML"
		bot.Send(msg)
		return
	}

	orderResults := []int{}
	for _, value := range values {
		orderResults = append(orderResults, results[value])
	}

	println("Creating chart for", head, strings.Join(values, ", "))
	// Create the pie chart
	var pie []byte
	if pie, err = chart.CreatePieChart(head, values, orderResults); err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Unable to complete the operation for <i>\"%s\"</i>.", head))
		msg.ParseMode = "HTML"
		bot.Send(msg)
		return
	}

	photo := tgbotapi.FileBytes{
		Name:  "chart.png",
		Bytes: pie,
	}

	println("Sending chart for", head, strings.Join(values, ", "))
	msg := tgbotapi.NewPhoto(chatID, photo)
	bot.Send(msg)
}
