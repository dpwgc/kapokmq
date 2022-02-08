package server

import (
	"encoding/json"
	"kapokmq/config"
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

var WAL *log.Logger
var WALFile *os.File

// InitWAL WAL写前日志
func InitWAL() {

	//是否开启写前日志
	if config.Get.Mq.IsPersistent != 2 {
		return
	}

	file := "./WAL.log"
	WALFile, _ = os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	WAL = log.New(WALFile, "", 0) // 将文件设置为loger作为输出
}

// SetWAL 写前日志
func SetWAL(message model.Message) {

	//追加写入
	jsonStr, _ := json.Marshal(message)
	WAL.Println(string(jsonStr))
}
