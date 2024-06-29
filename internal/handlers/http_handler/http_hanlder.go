package http_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tg_pay_gate/internal/services"
	"tg_pay_gate/internal/utils/config"
	"tg_pay_gate/internal/utils/functions"
)

func EpayNotify(c *gin.Context) {
	orderIDString := c.Query("out_trade_no")
	orderID, err := functions.ParseUUID(orderIDString)
	if err != nil {
		c.String(http.StatusBadRequest, "id格式错误")
		return
	}

	order, err := services.GetOrderByOrderID(orderID)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	epayConfig := config.EpayConfig
	epayKey := epayConfig.Key
	if err := services.EpayNotify(order, epayKey, c); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK, "success")
}
