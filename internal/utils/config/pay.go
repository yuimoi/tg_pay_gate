package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type EpayConfigStruct struct {
	Pid     string `json:"pid"`
	Key     string `json:"key"`
	Url     string `json:"url"`
	PayType string `json:"pay_type"`

	NotifyUrl string `json:"notify_url"`
}

func LoadPayConfig() {
	path := configBaseDir + "/epay_config.json"
	config := new(EpayConfigStruct)

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

	EpayConfig = config

	fmt.Printf("付费计划信息: %+v\n", config)
}

var EpayConfig *EpayConfigStruct
