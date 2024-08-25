package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io"
	"log"
	"net/http"
	"os"
)

const myChatID = -1002241449422

type message struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

var (
	tg *tgbotapi.BotAPI
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Origin", "http://roshydrostandart.ru")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		// Отвечаем на предварительные запросы OPTIONS
		w.WriteHeader(http.StatusOK)
		return
	}

	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	fmt.Println("Полученные данные: ", string(body))

	var info message
	if err = json.Unmarshal(body, &info); err != nil {
		fmt.Println("Ошибка при декодировании данных: ", err)
		return
	}
	_, err = tg.Send(tgbotapi.NewMessage(myChatID, fmt.Sprintf("Новая заявка!\nИмя: %s\nТелефон: %s", info.Name, info.Phone)))
	if err != nil {
		fmt.Println("Ошибка при отправке сообщения в Telegram:", err)
		return
	}
	// Отправляем ответ клиенту
	w.WriteHeader(http.StatusOK)
}

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("Не указан токен Telegram бота")
	}

	port := ":8080"

	tgInstance, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Ошибка при создании бота: %s", err)
	}
	tg = tgInstance

	http.HandleFunc("/", handler)

	fmt.Printf("Сервер запущен на порту %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Ошибка при запуске сервера:", err)
	}
}
