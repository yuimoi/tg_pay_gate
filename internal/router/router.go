package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg_pay_gate/internal/handlers/http_handler"
	"tg_pay_gate/internal/handlers/tg_handler"
	"tg_pay_gate/internal/services"
	"tg_pay_gate/internal/utils/tg_bot"
)

func SetupGinRoutes() *gin.Engine {
	r := gin.Default()

	r.GET("/api/epay_notify", http_handler.EpayNotify)

	return r
}

func RunTgBot() {
	bot := tg_bot.Bot
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		go handleUpdate(update)
	}
}

func handleUpdate(update tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			msgText := fmt.Sprintf("机器人处理消息崩溃, Error: %v", r)
			services.HandlePanic(r, msgText, "tg_bot_panic")
		}
	}()

	if update.Message != nil && update.Message.Chat.IsPrivate() { // 只处理私聊信息
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				tg_handler.StartCommand(update)
			case "pay":
				tg_handler.PayCommand(update)
			case "join_group":
				tg_handler.JoinGroupCommand(update)
			default:
				tg_handler.StartCommand(update)
			}

		} else {
			// 任何文本都返回start界面
			tg_handler.StartCommand(update)
		}
	}

	// 校验新成员
	if update.Message.NewChatMembers != nil {
		tg_handler.MemberJoin(update)
	}

}
