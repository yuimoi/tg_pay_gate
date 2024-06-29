package config

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"os"
	"sync"
	_type "tg_pay_gate/internal/utils/type"
)

type SiteConfigStruct struct {
	TgBotToken string `json:"tg_bot_token" desc:"进群机器人的token"`
	GroupID    int64  `json:"group_id" desc:"需要进群的群组id（需要将机器人邀请进群并设置为管理员）"`

	Price decimal.Decimal `json:"price" desc:"进群费用"`
	Host  string          `json:"host" desc:"http域名,用于支付回调"`

	Proxy _type.Proxy `json:"proxy" desc:"网络代理"`
}

func LoadSiteConfig() {
	SiteConfigLock.Lock()         // 在读取和处理配置之前加锁
	defer SiteConfigLock.Unlock() // 确保函数退出前释放锁

	path := configBaseDir + "/config.json"
	config := new(SiteConfigStruct)

	// 读取JSON文件
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// 反序列化JSON到config
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}

	SiteConfig = config

	fmt.Printf("Config loaded: %+v\n", config)
}

func GetSiteConfig() *SiteConfigStruct {
	SiteConfigLock.RLock()
	defer SiteConfigLock.RUnlock()
	return SiteConfig
}

var SiteConfigLock = &sync.RWMutex{}

var SiteConfig *SiteConfigStruct
