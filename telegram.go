package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"time"
)

func getBot() *tgbotapi.BotAPI {
	token := getConfigValue("telegram", "token")
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		fmt.Println(err)
	}

	bot.Debug = true
	return bot
}

func getMessage(text string) tgbotapi.MessageConfig {
	chatID, err := strconv.Atoi(getConfigValue("telegram", "chat_id"))

	if err != nil {
		panic(err)
	}

	message := tgbotapi.NewMessage(int64(chatID), text)
	return message
}

func getButtonsByMode(mode string) []tgbotapi.KeyboardButton {
	buttons := make([]tgbotapi.KeyboardButton, 0)
	switch mode {
	case "awaiting":
		buttons = append(buttons, tgbotapi.NewKeyboardButton("drink"))
		buttons = append(buttons, tgbotapi.NewKeyboardButton("p"))
		buttons = append(buttons, tgbotapi.NewKeyboardButton("agg"))
	case "choose_drink":
		for _, drink := range getDrinks() {
			buttons = append(buttons, tgbotapi.NewKeyboardButton(drink))
		}
		buttons = append(buttons, tgbotapi.NewKeyboardButton("another"))
	case "agg":
		buttons = append(buttons, tgbotapi.NewKeyboardButton("p"))
		buttons = append(buttons, tgbotapi.NewKeyboardButton("drink"))
	case "agg.p":
		buttons = append(buttons, tgbotapi.NewKeyboardButton("freq"))
		buttons = append(buttons, tgbotapi.NewKeyboardButton("period"))
	}
	if mode != "awaiting" {
		buttons = append(buttons, tgbotapi.NewKeyboardButton("cancel"))
	}
	return buttons
}

func getMsgByMode(mode string) string {
	switch mode {
	case "put_drink":
		return "enter drink"
	case "choose_drink":
		return "choose drink"
	case "volume":
		return "enter volume"
	case "agg":
		return "choose agg event type"
	case "agg.p":
		return "choose agg type"
	case "agg.p.freq":
		return getPFreq()
	case "agg.p.period":
		return getPPeriod()
	case "agg.drink":
		return getRate()
	default:
		return "default"
	}
}

func processInp(inp string, currentMode string, event *Event) string {
	if inp == "cancel" {
		return "awaiting"
	}

	switch currentMode {
	case "awaiting":
		switch inp {
		case "drink":
			event.Type = EventType{Name: "drink"}
			event.Time = time.Now()
			return "choose_drink"
		case "p":
			event.Type = EventType{Name: "p"}
			event.ID = getLastId("ns") + 1
			event.Time = time.Now()
			postEvent(event)
			return "awaiting"
		case "agg":
			return "agg"
		}
	case "choose_drink":
		if inp == "another" {
			return "put_drink"
		} else {
			event.Drink = inp
			return "volume"
		}
	case "put_drink":
		putDrink(inp)
		event.Drink = inp
		return "volume"
	case "volume":
		volume, err := strconv.Atoi(inp)
		treatErr(err)
		event.Volume = volume
		event.ID = getLastId("ns") + 1
		postEvent(event)
		return "awaiting"
	case "agg":
		switch inp {
		case "p":
			return "agg.p"
		case "drink":
			return "agg.drink"
		}
	case "agg.p":
		switch inp {
		case "freq":
			return "agg.p.freq"
		case "period":
			return "agg.p.period"
		}
	case "agg.p.freq":
	case "agg.p.period":
	case "agg.drink":
		return "awaiting"
	}
	return "awaiting"
}

func runBot() {
	bot := getBot()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	mode := "awaiting"

	event := Event{}
	msg := getMessage(getMsgByMode(mode))
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(getButtonsByMode(mode))
	bot.Send(msg)
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		mode = processInp(update.Message.Text, mode, &event)
		msg := getMessage(getMsgByMode(mode))
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(getButtonsByMode(mode))
		bot.Send(msg)
	}
}
