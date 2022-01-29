package persistent

import (
	"github.com/spf13/viper"
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
		//先获取二进制文件中的周期性全量备份数据
		Read()
		//再读取WAL日志文件中的消息
		ReadWAL()
		server.Loger.Println("Recovery from local")
	}
	//recoveryStrategy为其他数值时不进行数据恢复操作
}
