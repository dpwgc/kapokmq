package persistent

import (
	"github.com/spf13/viper"
	"kapokmq/server"
	"time"
)

/*
 * 持久化数据到硬盘
 */

// InitPers 加载数据持久化模块
func InitPers() {

	//是否开启持久化
	isPersistent := viper.GetInt("mq.isPersistent")
	if isPersistent == 0 {
		return
	}

	//两次持久化操作的间隔时间
	persistentTime := viper.GetInt("mq.persistentTime")

	go func() {
		server.Loger.Println("Start persistence")
		for {
			//将消息列表持久化写入文件
			Write()
			time.Sleep(time.Second * time.Duration(persistentTime))
		}
	}()
}
