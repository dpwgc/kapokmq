package persistent

import (
	"DPMQ/server"
	"github.com/spf13/viper"
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
		server.Loger.Println("Recovery from local")
	}
	//recoveryStrategy为其他数值时不进行数据恢复操作
}
