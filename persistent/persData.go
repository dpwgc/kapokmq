package persistent

import (
	"DPMQ/models"
	"DPMQ/servers"
	"github.com/spf13/viper"
	"time"
)

/*
 * 持久化数据到硬盘
 */

func InitPers() {

	//是否开启持久化
	isPersistent := viper.GetInt("mq.isPersistent")
	if isPersistent == 0 {
		return
	}

	//两次持久化操作的间隔时间
	persistentTime := viper.GetInt("mq.persistentTime")

	go func() {
		servers.Loger.Println("Start persistence")
		for {
			//复制消息列表
			copyData()
			//持久化写入
			Write()
			time.Sleep(time.Second * time.Duration(persistentTime))
		}
	}()
}

//消息列表拷贝
var copyMessageList []models.Message

//复制该主节点的数据
func copyData() {

	copyMessageList = servers.MessageList
}
