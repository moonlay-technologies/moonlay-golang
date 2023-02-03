package telegram

import (
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"order-service/global/utils/helper"
	"os"
)

type TelegramBotInterface interface {
	SetChatID(ChatID string)
	SetToken(token string)
	SendMessage(messages string, result chan error)
}

func InitTelegramBot(token string) TelegramBotInterface {
	return &telegramBot{
		Token: token,
	}
}

type telegramBot struct {
	Token  string `json:"token"`
	ChatID string `json:"group_id"`
}

func (g *telegramBot) SetChatID(ChatID string) {
	g.ChatID = ChatID
}

func (g *telegramBot) SetToken(token string) {
	g.Token = token
}

func (g *telegramBot) SendMessage(messages string, result chan error) {

	if len(g.ChatID) == 0 {
		newErr := helper.NewError("Group ID is required")
		result <- newErr
		return
	}

	if len(g.Token) == 0 {
		newErr := helper.NewError("Token is required")
		result <- newErr
		return
	}

	telegramData := map[string]interface{}{
		"chat_id": g.ChatID,
		"text":    messages,
	}

	telegramDataJson, err := json.Marshal(telegramData)

	if err != nil {
		errStr := fmt.Sprintf("Error Encode Json Body Telegram %s", err.Error())
		helper.SetSentryError(err, errStr, sentry.LevelError)
		result <- err
		return
	}

	httpRequestOption := helper.Options{
		Method:      "POST",
		Body:        telegramDataJson,
		ContentType: "application/json",
		URL:         fmt.Sprintf("%s/bot%s/sendMessage", os.Getenv("TELEGRAM_URL"), g.Token),
	}

	telegramRest := helper.POST(&httpRequestOption)

	if telegramRest.StatusCode == 200 {
		result <- nil
		return
	}

	result <- telegramRest.Error
	return
}
