package services

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"net/url"
	"sort"
	"strings"
	"tg_pay_gate/internal/models"
	"tg_pay_gate/internal/utils/config"
	"tg_pay_gate/internal/utils/db"
	"tg_pay_gate/internal/utils/tg_bot"
	_type "tg_pay_gate/internal/utils/type"
	"time"
)

func GetSuccessOrderByTgID(tgID int64) (*models.Order, error) {
	var order *models.Order
	result := db.DB.Model(models.Order{}).Where("tg_id=? and status=?", tgID, _type.OrderStatusSuccess).Find(&order)
	if result.Error != nil {
		return nil, errors.New("查询订单失败")
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("没有找到订单")
	}
	return order, nil
}

func EpayUrl(orderID string, price decimal.Decimal, productName string, epayConfig config.EpayConfigStruct) string {
	siteConfig := config.GetSiteConfig()
	notifyHost := siteConfig.Host
	submitData := map[string]string{
		"pid":          epayConfig.Pid,
		"type":         epayConfig.PayType,
		"out_trade_no": orderID,
		"notify_url":   fmt.Sprintf("%s%s", notifyHost, epayConfig.NotifyUrl),
		"return_url":   fmt.Sprintf("https://web.telegram.org"),
		"name":         productName,
		"money":        price.String(),
	}
	submitData["sign"] = EpaySign(submitData, epayConfig.Key)
	submitData["sign_type"] = "MD5"

	//生成url
	values := url.Values{}
	for key, value := range submitData {
		values.Add(key, value)
	}
	payUrl := fmt.Sprintf("%s?%s", epayConfig.Url, values.Encode())
	return payUrl
}

func EpaySign(mapInput map[string]string, epayKey string) string {
	//排序key获取排序后的key列表
	var keys []string
	for k := range mapInput {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var queryParts []string
	for _, key := range keys {
		key := key
		value := mapInput[key]
		if value == "" || key == "sign" || key == "sign_type" {
			continue
		}
		queryParts = append(queryParts, key+"="+value)
	}

	stringToSign := strings.Join(queryParts, "&")
	stringToSign = stringToSign + epayKey
	stringToSign, _ = url.QueryUnescape(stringToSign)

	//md5
	inputBytes := []byte(stringToSign)
	md5Hash := md5.Sum(inputBytes)
	md5String := fmt.Sprintf("%x", md5Hash)

	return md5String
}

func GetOrderByOrderID(orderID uuid.UUID) (*models.Order, error) {
	var order *models.Order
	if result := db.DB.Model(models.Order{}).Where("id = ?", orderID).Find(&order); result.RowsAffected == 0 {
		return nil, errors.New("没有找到订单")
	}
	return order, nil
}

func EpayNotify(order *models.Order, epayKey string, c *gin.Context) error {
	notifyData := map[string]string{
		"pid":          c.Query("pid"),
		"trade_no":     c.Query("trade_no"),
		"out_trade_no": c.Query("out_trade_no"),
		"type":         c.Query("type"),
		"name":         c.Query("name"),
		"money":        c.Query("money"),
		"trade_status": c.Query("trade_status"),
		"sign":         c.Query("sign"),
		"sign_type":    c.Query("sign_type"),
	}
	inputSign := c.Query("sign")
	calculateSign := EpaySign(notifyData, epayKey)

	if inputSign != calculateSign {
		return errors.New("签名错误")
	}

	// 更新订单状态,拉用户进群
	err := OrderSuccess(order)
	if err != nil {
		return err
	}

	return nil

}

func OrderSuccess(order *models.Order) error {
	result := db.DB.Model(models.Order{}).Where("id=? and status=?", order.ID, _type.OrderStatusPending).Updates(map[string]interface{}{
		"status": _type.OrderStatusSuccess,
	})

	if result.Error != nil {
		return errors.New("更新订单状态错误")
	}

	if result.RowsAffected == 0 {
		return errors.New("没有该订单")
	}

	siteConfig := config.GetSiteConfig()
	// 订单完成前，尝试解封用户，防止用户被踢出后无法使用邀请链接，tg踢出用户默认行为是封禁用户
	_ = tg_bot.UnbanUser(siteConfig.GroupID, order.TgID)

	// 发送成功消息
	msg := tgbotapi.NewMessage(order.TgID, "支付成功")
	_, _ = tg_bot.Bot.Send(msg)

	// 拉用户进群
	err := tg_bot.SendInviteJoinGroup(order.TgID)
	if err != nil {
		return err
	}

	return nil
}

func ClearPendingOrder(timeThreshold time.Duration) error {
	nowTimestamp := time.Now().Unix()
	thresholdTime := nowTimestamp - int64(timeThreshold.Seconds())

	// 查询符合条件的订单
	var orders []models.Order
	result := db.DB.Where("create_time < ? AND status=?", thresholdTime, _type.OrderStatusPending).Find(&orders)
	if result.Error != nil {
		return errors.New("查询进行订单失败")
	}

	// 如果有符合条件的订单，执行删除操作
	if len(orders) > 0 {
		deleteResult := db.DB.Delete(&orders)
		if deleteResult.Error != nil {
			return errors.New("清理订单失败")
		}
	}

	return nil
}
