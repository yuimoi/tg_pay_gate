package tg_bot

import (
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
	"runtime"
	"tg_pay_gate/internal/models"
	"tg_pay_gate/internal/utils/config"
	"time"
)

var Bot *tgbotapi.BotAPI

func InitTGBot() {
	client := &http.Client{}

	myProxy := config.GetSiteConfig().Proxy
	if myProxy.EnableProxy == true && runtime.GOOS == "windows" {
		tgProxyURL, err := url.Parse(fmt.Sprintf("%s://%s:%d", myProxy.Protocol, myProxy.Host, myProxy.Port))
		if err != nil {
			panic(fmt.Sprintf("Failed to parse proxy: %s\n", err))
		}

		tgDialer, err := proxy.FromURL(tgProxyURL, proxy.Direct)
		if err != nil {
			panic(fmt.Sprintf("Failed to obtain proxy dialer: %s\n", err))
		}
		tgTransport := &http.Transport{
			Dial: tgDialer.Dial,
		}
		client.Transport = tgTransport
	}

	fmt.Println("正在连接TG")
	var err error
	Bot, err = tgbotapi.NewBotAPIWithClient(config.GetSiteConfig().TgBotToken, "https://api.telegram.org/bot%s/%s", client)
	if err != nil {
		panic(err)
	}
	fmt.Println("TG连接成功")

	fmt.Println("正在校验管理员身份")
	isAdmin, err := isBotAdmin()
	if err != nil {
		panic(err)
	}
	if !isAdmin {
		panic("机器人不是群组管理员")
	}
	fmt.Println("校验管理员身份成功")

	// 检查成员

	//Bot.Debug = config.GetSiteConfig().EnableTGBotDebug
}

func SendTgMsg(msgText string, currentUser *models.User) error {
	msg := tgbotapi.NewMessage(currentUser.TgID, msgText)
	msg.DisableWebPagePreview = true
	msg.ParseMode = "HTML"
	_, err := Bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func DeleteMsg(chatID int64, msgID int) error {
	deleteConfig := tgbotapi.DeleteMessageConfig{
		ChatID:    chatID,
		MessageID: msgID,
	}
	if _, err := Bot.Request(deleteConfig); err != nil {
		return err
	}
	return nil
}

// 检查机器人是否是群组的管理员
func isBotAdmin() (bool, error) {
	chatID := config.GetSiteConfig().GroupID

	// 获取群组管理员列表
	admins, err := Bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: chatID}})
	if err != nil {
		return false, err
	}

	// 检查机器人是否在管理员列表中
	for _, admin := range admins {
		if admin.User.ID == Bot.Self.ID { // 使用全局变量 Bot.Self.ID
			if admin.CanInviteUsers {
				return true, nil
			} else {
				return false, errors.New("机器人没有邀请权限")
			}
		}
	}
	return false, errors.New("机器人不是群组管理员")
}

func CreateInviteLink(inviteTgID int64, duration time.Duration) (string, error) {
	// 假设要邀请的群组 ID
	var resMap map[string]interface{}

	// 生成邀请链接
	inviteLinkConfig := tgbotapi.CreateChatInviteLinkConfig{
		ChatConfig:         tgbotapi.ChatConfig{ChatID: inviteTgID},
		Name:               "邀请进群",
		ExpireDate:         int(time.Now().Add(duration).Unix()),
		MemberLimit:        1,
		CreatesJoinRequest: false, //是否需要管理员同意
	}
	response, _ := Bot.Request(inviteLinkConfig)
	err := json.Unmarshal(response.Result, &resMap)
	if err != nil {
		return "", errors.New("创建邀请链接失败")
	}

	return resMap["invite_link"].(string), nil
}

func SendInviteJoinGroup(userTgID int64) error {
	// 假设要邀请的群组 ID
	groupID := config.GetSiteConfig().GroupID
	inviteLinkDuration := time.Minute * 1
	inviteLink, err := CreateInviteLink(groupID, inviteLinkDuration)
	if err != nil {
		return err
		//msg := tgbotapi.NewMessage(userTgID, err.Error())
		//_, err = Bot.Send(msg)
	}

	sendText := fmt.Sprintf("<a href=\"%s\">点击进群(60秒内有效)</a>", inviteLink)

	// 发送邀请链接给用户
	msg := tgbotapi.NewMessage(userTgID, sendText)
	msg.DisableWebPagePreview = true
	msg.ParseMode = "HTML"
	_, err = Bot.Send(msg)

	return nil
}

func UnbanUser(chatID int64, userID int64) error {
	unbanConfig := tgbotapi.UnbanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		OnlyIfBanned: true, // 仅当用户被封禁时才解除封禁
	}
	_, err := Bot.Request(unbanConfig)
	return err
}
