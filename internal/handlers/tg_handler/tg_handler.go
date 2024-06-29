package tg_handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg_pay_gate/internal/models"
	"tg_pay_gate/internal/services"
	"tg_pay_gate/internal/utils/config"
	"tg_pay_gate/internal/utils/db"
	"tg_pay_gate/internal/utils/tg_bot"
	"time"
)

var PayCommandString = "/pay"
var StartCommandString = "/start"
var JoinGroupCommandString = "/join_group"

func StartCommand(update tgbotapi.Update) {
	msgText := fmt.Sprintf(`欢迎使用进群机器人
付费进群: %s
已付费, 点击进群: %s
`, PayCommandString, JoinGroupCommandString)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	_, _ = tg_bot.Bot.Send(msg)
}

func PayCommand(update tgbotapi.Update) {
	tgID := update.Message.Chat.ID
	successOrder, _ := services.GetSuccessOrderByTgID(tgID)
	if successOrder != nil {
		msg := tgbotapi.NewMessage(tgID, "已付款,请点击进群 "+JoinGroupCommandString)
		_, _ = tg_bot.Bot.Send(msg)
		return
	}

	// 清理太久未支付的订单
	_ = services.ClearPendingOrder(time.Hour * 24)

	siteConfig := config.GetSiteConfig()
	// 创建订单
	newOrder := models.NewOrder(siteConfig.Price, tgID)
	result := db.DB.Create(&newOrder)
	if result.Error != nil {
		msg := tgbotapi.NewMessage(tgID, "创建订单失败")
		_, _ = tg_bot.Bot.Send(msg)
		return
	}

	// 构建订单url
	epayConfig := *config.EpayConfig
	payUrl := services.EpayUrl(newOrder.ID.String(), siteConfig.Price, "付费进群", epayConfig)

	sendText := fmt.Sprintf("<a href=\"%s\">点击支付</a>", payUrl)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, sendText)
	msg.DisableWebPagePreview = true
	msg.ParseMode = "HTML"
	_, _ = tg_bot.Bot.Send(msg)
}

func JoinGroupCommand(update tgbotapi.Update) {
	userTgID := update.Message.Chat.ID
	successOrder, _ := services.GetSuccessOrderByTgID(userTgID)
	if successOrder == nil {
		msg := tgbotapi.NewMessage(userTgID, "未付款,请先付款后再进群 "+PayCommandString)
		_, _ = tg_bot.Bot.Send(msg)
		return
	}

	// 假设要邀请的群组 ID
	err := tg_bot.SendInviteJoinGroup(userTgID)
	if err != nil {
		msg := tgbotapi.NewMessage(userTgID, err.Error())
		_, err = tg_bot.Bot.Send(msg)
	}

}

func MemberJoin(update tgbotapi.Update) {
	siteConfig := config.GetSiteConfig()
	newChatMembers := update.Message.NewChatMembers
	for _, newChatMember := range newChatMembers {
		successOrder, _ := services.GetSuccessOrderByTgID(newChatMember.ID)
		if successOrder == nil {
			// 这个命令会使当前所有邀请链接对该用户失效
			kickConfig := tgbotapi.KickChatMemberConfig{
				ChatMemberConfig: tgbotapi.ChatMemberConfig{
					ChatID: siteConfig.GroupID,
					UserID: newChatMember.ID,
				},
			}
			_, err := tg_bot.Bot.Request(kickConfig)
			if err != nil {

			}
		}
	}

}
