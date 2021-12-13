package config

import (
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.AddConfigPath("./config")
	viper.SetConfigFile("application") // 指定配置文件路径
	viper.SetConfigName("application") //配置文件名
	viper.SetConfigType("yaml")        //配置文件类型
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
