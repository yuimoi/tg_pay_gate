package config

import (
	"math/rand"
	"os"
)

var SiteSecret string

func GenerateRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func LoadSiteSecret() {
	secretFilePath := configBaseDir + "/.secret"
	if _, err := os.Stat(secretFilePath); os.IsNotExist(err) {
		secret := GenerateRandomString(32)
		err = os.WriteFile(secretFilePath, []byte(secret), 0600)
		if err != nil {
			panic(err)
		}
	}

	secret, err := os.ReadFile(secretFilePath)
	if err != nil {
		panic(err)
	}
	SiteSecret = string(secret)
}
