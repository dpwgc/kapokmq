package persistent

import (
	"github.com/spf13/viper"
	"kapokmq/model"
	"kapokmq/server"
)

/*
 * 主节点重启后进行数据恢复
 */

// InitRecovery 数据恢复到内存
func InitRecovery() {

	// 数据恢复策略
	recoveryStrategy := viper.GetInt("mq.recoveryStrategy")

	//从本地持久化文件中获取数据
	if recoveryStrategy == 1 {
		//本地恢复数据
		Read()
		//将未消费的信息重新导入消息队列
		server.MessageList.Range(func(key, value interface{}) bool {
			msg := value.(model.Message)
			if msg.Status == -1 {
				server.MessageChan <- msg
			}
			return true
		})
		server.Loger.Println("Recovery from local")
	}
	//recoveryStrategy为其他数值时不进行数据恢复操作
}
