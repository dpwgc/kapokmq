package persistent

import (
	"DPMQ/model"
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

	//复制消息列表
	copyMap := make(map[string]interface{})
	server.MessageList.Range(func(key, value interface{}) bool {
		copyMap[key.(string)] = value
		return true
	})

	//将消息列表拷贝转为[]Byte数据
	copyBytes, _ := json.Marshal(copyMap)
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

	copyMap, err := jsonToMessage(record[0][0])
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
