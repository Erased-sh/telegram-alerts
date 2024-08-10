package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io"
	"log"
	"net/http"
)

const (
	token    = "7422382792:AAH06eSEfSufB3DIIpWjFWDn8Ejt5YTpSoE"
	myChatID = -1002241449422
)

type message struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

var (
	tg *tgbotapi.BotAPI
)

// handler функция для обработки HTTP-запросов с поддержкой CORS
func handler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовки CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")                   // Разрешаем все источники (или укажите конкретный домен)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Разрешаем методы
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")       // Разрешаем заголовки

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
	// Определите порт для сервера
	port := ":8080"

	tgInstance, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Ошибка при создании бота: %s", err)
	}
	tg = tgInstance
	// Установите маршрут и обработчик
	http.HandleFunc("/", handler)

	// Запустите сервер
	fmt.Printf("Сервер запущен на порту %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Ошибка при запуске сервера:", err)
	}
}
