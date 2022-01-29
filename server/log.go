package server

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"kapokmq/model"
	"log"
	"os"
	"time"
)

/**
 * 日志记录
 */

var Loger *log.Logger

func InitLog() {

	file := "./log/kapokmq-" + time.Now().Format("2006-01-02") + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		//创建log目录
		err = os.Mkdir("./log", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	Loger = log.New(logFile, "", log.LstdFlags|log.Lshortfile|log.LUTC) // 将文件设置为loger作为输出
}

var isPersistent int
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
