package persistent

import (
	"DPMQ/server"
	"encoding/csv"
	"encoding/json"
	"github.com/spf13/viper"
	"os"
)

/*
 * 持久化文件读写
 */

var wFile *os.File
var rFile *os.File
var path string

//初始化文件
func InitFileRW() {

	path = viper.GetString("mq.persistentPath")
	//判断持久化文件是否存在
	_, err := os.Stat(path)
	if err != nil {
		//创建持久化文件
		server.Loger.Println(err)
		server.Loger.Println("Create persistent file: " + path)
		_, err = os.Create(path)
		if err != nil {
			server.Loger.Println(err)
		}
	}
}

//写入持久化文件
func Write() {
	var err error
	//写文件，设置为只写、覆盖，权限设置为777
	wFile, err = os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)

	writer := csv.NewWriter(wFile)

	//将复制得到的消息列表转为[]Byte数据
	copyBytes, _ := json.Marshal(copyMessageList)
	//将数据以json字符串形式存入持久化文件
	jsonStr := string(copyBytes)
	err = writer.Write([]string{jsonStr})
	if err != nil {
		server.Loger.Println(err)
		return
	}
	writer.Flush()
	//关闭文件流
	err = wFile.Close()
	if err != nil {
		server.Loger.Println(err)
	}
}

//加载持久化文件内的数据到内存中
func Read() {
	var err error
	//读文件，设置为只读，权限设置为777
	rFile, err = os.OpenFile(path, os.O_RDONLY, 0777)
	if err != nil {
		server.Loger.Println(err)
		return
	}
	reader := csv.NewReader(rFile)
	reader.FieldsPerRecord = -1
	record, err := reader.ReadAll()
	if err != nil {
		server.Loger.Println(err)
		return
	}
	if len(record) == 0 {
		server.Loger.Println("The file is empty")
		return
	}
	//解析本地持久化文件的数据到MessageList
	err = json.Unmarshal([]byte(record[0][0]), &server.MessageList)
	if err != nil {
		server.Loger.Println(err)
		return
	}

	//将未消费的消息插入消息队列
	for _, m := range server.MessageList {
		if m.Status == 0 {
			server.MessageChan <- m
		}
	}

	err = rFile.Close()
	if err != nil {
		server.Loger.Println(err)
	}
}
