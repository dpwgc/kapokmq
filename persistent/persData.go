package persistent

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"kapokmq/model"
	"kapokmq/server"
	"log"
	"os"
	"time"
)

/*
 * 持久化数据到硬盘
 */

var isPersistent int

// InitPers 加载数据持久化模块
func InitPers() {

	//是否开启持久化
	isPersistent = viper.GetInt("mq.isPersistent")
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

var WAL *log.Logger

// InitWAL WAL持久化日志
func InitWAL() {

	//是否开启WAL持久化
	isPersistent = viper.GetInt("mq.isPersistent")
	if isPersistent != 2 {
		return
	}

	file := "./WAL.log"
	logFile, _ := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	WAL = log.New(logFile, "", log.LstdFlags|log.Lshortfile|log.LUTC) // 将文件设置为loger作为输出
}

// SetWAL 消息追加写入日志
func SetWAL(message model.Message) {

	//如果没有开启WAL持久化
	if isPersistent != 2 {
		return
	}

	//追加写入
	jsonStr, _ := json.Marshal(message)
	WAL.Println(fmt.Sprintf("%s%s", "\t\\|SET|\\\t", string(jsonStr)))
}
