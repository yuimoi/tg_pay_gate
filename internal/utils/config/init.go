package config

import (
	"os"
	"path/filepath"
	"time"
)

func GetRootDir() string {
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	rootDir := filepath.Dir(executable)
	return rootDir
}
func GetEnvDir() string {
	envDir := filepath.Join(GetRootDir(), ".env")
	return envDir
}

func LoadAllConfig() {
	LoadSiteSecret()
	LoadSiteConfig()
	LoadPayConfig()

	Loc, _ = time.LoadLocation("Asia/Shanghai")
}

var configBaseDir = GetEnvDir()

var Loc *time.Location
