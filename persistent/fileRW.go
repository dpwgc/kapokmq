package persistent

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"kapokmq/config"
	"kapokmq/memory"
	"kapokmq/model"
	"kapokmq/mqLog"
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

	//是否开启持久化
	isPersistent := config.Get.Mq.IsPersistent
	if isPersistent == 0 {
		return
	}

	//持久化文件名
	dataFile = config.Get.Mq.PersistentFile

	//获取消息通道的缓冲区大小
	messageChanBuffer := config.Get.Mq.MessageChanBuffer

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
		mqLog.Loger.Println(err)
		mqLog.Loger.Println(fmt.Sprintf("%s%s", "Create persistent file: ", dataFile))
		_, err = os.Create(dataFile)
		if err != nil {
			mqLog.Loger.Println(err)
		}
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			mqLog.Loger.Println(err)
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
	memory.MessageList.Range(func(key, value interface{}) bool {
		//如果是已经消费的消息
		if value.(model.Message).Status == 1 {
			return true
		}
		copyMap[key.(string)] = value
		return true
	})

	//将消息列表拷贝转为[]Byte数据存入
	copyBytes, _ := json.Marshal(copyMap)
	err = writer.Encode(copyBytes)
	if err != nil {
		mqLog.Loger.Println(err)
		return
	}
	//关闭文件流
	err = wFile.Close()
	if err != nil {
		mqLog.Loger.Println(err)
	}
}

// Read 加载持久化文件内的数据到内存中
func Read() {
	var err error
	//读文件，设置为只读，权限设置为777
	rFile, err = os.OpenFile(dataFile, os.O_RDONLY, 0777)
	if err != nil {
		mqLog.Loger.Println(err)
		return
	}

	//取出数据
	reader := gob.NewDecoder(rFile)
	var info []byte
	err = reader.Decode(&info)
	if err != nil {
		mqLog.Loger.Println(err)
		return
	}
	if len(info) == 0 {
		mqLog.Loger.Println("The file is empty")
		return
	}

	copyMap := jsonToMessageMap(string(info))

	//将本地持久化文件数据（copyMap）循环写入消息列表（MessageList）
	for key, value := range copyMap {
		memory.MessageList.Store(key, value)
	}
	err = rFile.Close()
	if err != nil {
		mqLog.Loger.Println(err)
	}
}

// ReadWAL 读取WAL日志文件
func ReadWAL() {
	f, err := os.Open("./WAL.log")
	//如果文件不存在，则直接返回
	if err != nil {
		mqLog.Loger.Println("WAL.log does not exist")
		return
	}
	r := bufio.NewReader(f)
	for {
		// 读取文件(行读取)
		slice, err := r.ReadString('\n')

		//将行字符串解析为Message结构体
		message := jsonToMessage(slice)
		//将读取到的消息更新到消息列表
		memory.MessageList.Store(message.MessageCode, message)

		//如果读取到文件末尾
		if err == io.EOF {
			break
		}
	}
	err = f.Close()
	if err != nil {
		mqLog.Loger.Println(err)
		return
	}
}

// CleanWAL 清空WAL日志文件
func CleanWAL() {

	//关闭WAL日志文件
	err := mqLog.WALFile.Close()
	if err != nil {
		mqLog.Loger.Println(err)
		return
	}

	//删除WAL日志文件
	err = os.Remove("WAL.log")
	if err != nil {
		mqLog.Loger.Println(err)
		return
	}

	//重新创建WAL日志文件
	mqLog.InitWAL()
}

//json字符串转Message Map
func jsonToMessageMap(jsonStr string) map[string]model.Message {
	m := make(map[string]model.Message, copyMapSize)
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		mqLog.Loger.Println(err)
		return nil
	}
	return m
}

//json字符串转Message
func jsonToMessage(jsonStr string) model.Message {
	m := model.Message{}
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		mqLog.Loger.Println(err)
		return m
	}
	return m
}
