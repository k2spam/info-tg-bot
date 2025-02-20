package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Bot struct {
	Users string
	Url   string
}

type RequestMessage struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Item     string `json:"item"`
	Type     string `json:"itemtype"`
	Problems string `json:"problems"`
	Repair   string `json:"repair"`
}

type UpdateData struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		Chat struct {
			ID int `json:"id"`
		} `json:"chat"`
		Text string `json:"text"`
	} `json:"message"`
}

type ResultData struct {
	Result []UpdateData `json:"result"`
}

func (bot Bot) loadChatIDs() ([]int, error) {
	data, err := os.ReadFile(bot.Users)
	if err != nil {
		if os.IsNotExist(err) {
			return []int{}, nil
		}
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var ids []int
	for _, line := range lines {
		if line == "" {
			continue
		}
		id, err := strconv.Atoi(line)
		if err == nil {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func (bot Bot) subscribeHandler(chatId int) error {
	ids, _ := bot.loadChatIDs()
	for _, id := range ids {
		if id == chatId {
			return nil
		}
	}

	f, err := os.OpenFile(bot.Users, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(strconv.Itoa(chatId) + "\n")
	return err
}

func (bot Bot) unsubscribeHandler(chatId int) error {
	ids, _ := bot.loadChatIDs()
	var newIds string

	for _, v := range ids {
		if v != chatId {
			newIds = fmt.Sprintf("%s%d\n", newIds, v)
		}
	}

	f, err := os.OpenFile(bot.Users, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(newIds)
	return err
}

func (bot Bot) getUpdates(offset int) ([]UpdateData, error) {
	res, err := http.Get(fmt.Sprintf("%s/getUpdates?offset=%d&timeout=60", bot.Url, offset))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result ResultData

	err = json.NewDecoder(res.Body).Decode(&result)
	return result.Result, err
}

func (bot Bot) sendMessage(chatId int, text string) error {
	data := map[string]interface{}{
		"chat_id": chatId,
		"text":    text,
	}
	jsonData, _ := json.Marshal(data)

	resp, err := http.Post(bot.Url+"/sendMessage", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (bot Bot) handleUpdate(update UpdateData) {
	chatId := update.Message.Chat.ID
	text := update.Message.Text

	if text == "/subscribe" {
		bot.subscribeHandler(chatId)
		bot.sendMessage(chatId, "✅ Ты подписался на уведомления!")
	}
	if text == "/unsubscribe" {
		bot.unsubscribeHandler(chatId)
		bot.sendMessage(chatId, "❌ Подписка отменена")
	}
}

func (bot Bot) startBot() {
	offset := 0
	for {
		updates, err := bot.getUpdates(offset)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		for _, update := range updates {
			bot.handleUpdate(update)
			offset = update.UpdateID + 1
		}
	}
}

func (bot Bot) getParamString(params RequestMessage) string {
	name := ""
	phone := ""
	item := ""
	itemType := ""
	problems := ""
	repair := ""

	if params.Name == "" {
		name = "Имя: Не указано\n"
	}
	if params.Phone != "" {
		phone = fmt.Sprintf("Телефон: %s\n", params.Phone)
	}
	if params.Item != "" {
		item = fmt.Sprintf("Техника: %s\n", params.Item)
	}
	if params.Type != "" {
		item = fmt.Sprintf("Техника(доп): %s\n", params.Type)
	}
	if params.Problems != "" {
		item = fmt.Sprintf("Поломки: %s\n", params.Problems)
	}
	if params.Repair != "" {
		item = fmt.Sprintf("Требуемый ремонт: %s\n", params.Repair)
	}
	return fmt.Sprintf(`%s%s%s%s%s%s`, name, phone, item, itemType, problems, repair)
}

func (bot Bot) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	var request RequestMessage
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&request)
	message := bot.getParamString(request)

	chatIDs, _ := bot.loadChatIDs()
	for _, chatID := range chatIDs {
		bot.sendMessage(chatID, message)
	}
}

func main() {
	godotenv.Load(".env")

	bot := Bot{
		Url:   "https://api.telegram.org/bot" + os.Getenv("TOKEN"),
		Users: os.Getenv("USERS"),
	}

	go bot.startBot()
	http.HandleFunc("/order", bot.HTTPHandler)
	http.ListenAndServe(":8081", nil)
}
