package persistent

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"kapokmq/model"
	"kapokmq/server"
	"os"
)

/*
 * 持久化文件读写
 */

var wFile *os.File
var rFile *os.File
var dataFile string

//消息复制map的初始大小
var copyMapSize int

// InitFileRW 初始化文件读写模块
func InitFileRW() {

	//持久化文件名
	dataFile = viper.GetString("mq.persistentFile")

	//获取消息通道的缓冲区大小
	messageChanBuffer := viper.GetInt("mq.messageChanBuffer")
	//如果消息通道缓冲区大小小于等于1000
	if messageChanBuffer <= 1000 {
		//消息复制map初始大小设为1000
		copyMapSize = 1000
	} else {
		//消息复制map初始大小设为消息通道缓冲区大小的千分之一
		copyMapSize = messageChanBuffer / 1000
	}

	//判断持久化文件是否存在
	f, err := os.Open(dataFile)
	if err != nil {
		//创建持久化文件
		server.Loger.Println(err)
		server.Loger.Println(fmt.Sprintf("%s%s", "Create persistent file: ", dataFile))
		_, err = os.Create(dataFile)
		if err != nil {
			server.Loger.Println(err)
		}
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			server.Loger.Println(err)
		}
	}(f)
}

// Write 写入持久化文件
func Write() {

	var err error
	//写文件，设置为只写、覆盖，权限设置为777
	wFile, err = os.OpenFile(dataFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)

	writer := gob.NewEncoder(wFile)

	//复制消息列表
	copyMap := make(map[string]interface{}, copyMapSize)
	server.MessageList.Range(func(key, value interface{}) bool {
		copyMap[key.(string)] = value
		return true
	})

	//将消息列表拷贝转为[]Byte数据存入
	copyBytes, _ := json.Marshal(copyMap)
	err = writer.Encode(copyBytes)
	if err != nil {
		server.Loger.Println(err)
		return
	}
	//关闭文件流
	err = wFile.Close()
	if err != nil {
		server.Loger.Println(err)
	}
}

// Read 加载持久化文件内的数据到内存中
func Read() {
	var err error
	//读文件，设置为只读，权限设置为777
	rFile, err = os.OpenFile(dataFile, os.O_RDONLY, 0777)
	if err != nil {
		server.Loger.Println(err)
		return
	}

	//取出数据
	reader := gob.NewDecoder(rFile)
	var info []byte
	err = reader.Decode(&info)
	if err != nil {
		server.Loger.Println(err)
		return
	}
	if len(info) == 0 {
		server.Loger.Println("The file is empty")
		return
	}

	copyMap, err := jsonToMessage(string(info))
	if err != nil {
		server.Loger.Println(err)
		return
	}

	//将本地持久化文件数据（localMap）循环写入从节点map（DataMap）
	for key, value := range copyMap {
		server.MessageList.Store(key, value)
	}
	err = rFile.Close()
	if err != nil {
		server.Loger.Println(err)
	}
}

//json字符串转Data结构体
func jsonToMessage(jsonStr string) (map[string]model.Message, error) {
	m := make(map[string]model.Message)
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		server.Loger.Println(err)
		return nil, err
	}
	return m, nil
}
