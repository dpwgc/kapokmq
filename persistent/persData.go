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

	//周期性全量持久化方式
	if isPersistent == 1 {
		go func() {
			server.Loger.Println("Start persistence")
			for {
				//将消息列表全量持久化到二进制文件
				Write()
				time.Sleep(time.Second * time.Duration(persistentTime))
			}
		}()
	}

	//周期性全量持久化与追加写入日志结合
	if isPersistent == 2 {
		go func() {
			server.Loger.Println("Start persistence")
			for {
				//将消息列表全量持久化到二进制文件
				Write()
				//清空两次全量持久化之间的WAL日志
				CleanWAL()
				time.Sleep(time.Second * time.Duration(persistentTime))
			}
		}()
	}
}
