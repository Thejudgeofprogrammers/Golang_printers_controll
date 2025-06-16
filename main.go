package main

import (
	"fmt"
	"log"
	"printers/internal/config"
	"printers/internal/interfaces"
	"printers/internal/service"
	"strconv"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg, err := config.LoadJSON()
	if err != nil {
		log.Fatal("Ошибка загрузки конфига:", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg["TELEGRAM_BOT_TOKEN"])
	if err != nil {
		log.Fatal(err)
	}

	var debug string = cfg["DEBUG"]
	if debug == "1" {
		bot.Debug = true
	} else {
		bot.Debug = false
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	fmt.Println("Программа запущена")
	for update := range updates {
		if update.Message != nil {
			text := update.Message.Text

			if text == "/start" {
				rows := [][]tgbotapi.KeyboardButton{}
				for i := 1; i <= 2; i++ {
					hostKey := fmt.Sprintf("HOST_MOON_%d", i)
					hostIP := cfg[hostKey]
					if hostIP == "" {
						break
					}
					btn := tgbotapi.NewKeyboardButton(fmt.Sprintf("Принтер %d", i))
					rows = append(rows, tgbotapi.NewKeyboardButtonRow(btn))
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выбери принтер:")
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(rows...)
				if _, err := bot.Send(msg); err != nil {
					log.Println("Ошибка отправки:", err)
				}
				return
			}

			if strings.HasPrefix(text, "Принтер ") {
				go func(chatID int64, printerName string) {
					numStr := strings.TrimPrefix(printerName, "Принтер ")
					num, err := strconv.Atoi(numStr)
					if err != nil {
						bot.Send(tgbotapi.NewMessage(chatID, "Неверный формат принтера"))
						return
					}

					hostKey := fmt.Sprintf("HOST_MOON_%d", num)
					hostIP := cfg[hostKey]

					if hostIP == "" {
						bot.Send(tgbotapi.NewMessage(chatID, "IP принтера не найден"))
						return
					}

					var wg sync.WaitGroup
					wg.Add(2)

					var info interfaces.PrinterInfo
					var infoErr error

					var photoPath string
					var photoErr error

					go func() {
						defer wg.Done()
						info, infoErr = service.GetPrintersInfo(hostIP)
					}()

					go func() {
						defer wg.Done()
						photoPath, photoErr = service.GetPhoto(hostIP)
					}()

					wg.Wait()

					if infoErr != nil {
						bot.Send(tgbotapi.NewMessage(chatID, "Ошибка получения информации о принтере"))
						return
					}

					if photoErr != nil {
						bot.Send(tgbotapi.NewMessage(chatID, "Ошибка получения фото с принтера"))
						return
					}

					infoText := fmt.Sprintf(
						"Прогресс: %s%%\nОкончание: %s\nОсталось: %s",
						info.Success, info.DateEnd, info.EstimatedTime,
					)

					photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(photoPath))
					if _, err := bot.Send(photo); err != nil {
						log.Println("Ошибка отправки фото:", err)
						return
					}

					if _, err := bot.Send(tgbotapi.NewMessage(chatID, infoText)); err != nil {
						log.Println("Ошибка отправки сообщения:", err)
					}
				}(update.Message.Chat.ID, text)
			}
		}
	}
}