package syncConn

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"kapokmq/config"
	"kapokmq/memory"
	"kapokmq/model"
	"kapokmq/mqLog"
	"log"
)

// StopPush 是否停止推送（true：是，false：否），仅在从节点生效
var StopPush bool

// SlaveConn 从节点连接到主节点
func SlaveConn() {

	var err error

	//获取主节点的地址
	masterProtocol := config.Get.Sync.MasterProtocol
	masterAddr := config.Get.Sync.MasterAddr
	masterPort := config.Get.Sync.MasterPort

	//与主节点建立websocket连接
	wsUrl := fmt.Sprintf("%s://%s:%s%s", masterProtocol, masterAddr, masterPort, "/Sync/Conn")
	Conn, _, err = websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		mqLog.Loger.Println(err)
		panic(err)
	}
}

// SlaveSync 从节点同步协程
func SlaveSync() {

	defer func(Conn *websocket.Conn) {
		err := Conn.Close()
		if err != nil {
			mqLog.Loger.Println(err)
		}
	}(Conn)

	//验证密钥
	for {
		//读取主节点发送过来的提示
		_, data, err := Conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
			return
		}

		//请输入访问密钥
		if string(data) == "Please enter the secret key" {

			//发送密钥
			err = Conn.WriteMessage(1, []byte(config.Get.Mq.SecretKey))
			if err != nil {
				log.Fatal(err)
				return
			}
		}

		//访问密钥错误
		if string(data) == "Secret key matching error" {
			log.Fatal("Secret key matching error")
		}

		//访问密钥正确
		if string(data) == "Secret key matching succeeded" {
			//关闭从节点的推送功能
			StopPush = true
			break
		}
	}

	//开始监听数据
	for {
		_, data, err := Conn.ReadMessage()
		if err != nil {
			mqLog.Loger.Println(err)
			return
		}

		//解析json字符串，生成消息模板
		message := model.Message{}
		err = json.Unmarshal(data, &message)
		if err != nil {
			mqLog.Loger.Println(err)
			return
		}

		//如果开启了WAL写前日志
		if config.Get.Mq.IsPersistent == 2 {
			mqLog.SetWAL(message)
		}

		//将消息更新到从节点消息列表
		memory.MessageList.Store(message.MessageCode, message)
	}
}

//检查与重连
func checkConn() {

	//向主节点发送心跳检测
	err := Conn.WriteMessage(1, []byte("hi"))
	//如果发送出错，则表明连接已断开
	if err != nil {

		//获取主节点的地址
		masterProtocol := config.Get.Sync.MasterProtocol
		masterAddr := config.Get.Sync.MasterAddr
		masterPort := config.Get.Sync.MasterPort

		//尝试与主节点建立websocket连接
		wsUrl := fmt.Sprintf("%s://%s:%s%s", masterProtocol, masterAddr, masterPort, "/Sync/Conn")
		Conn, _, err = websocket.DefaultDialer.Dial(wsUrl, nil)
		//如果依旧无法连接，则判定主节点宕机
		if err != nil {
			//开启从节点的消息推送功能
			StopPush = false
			mqLog.Loger.Println(err)
			mqLog.Loger.Println("Start push")
		} else {
			//重连成功，关闭从节点的消息推送功能
			StopPush = true
			//重新启动从节点同步协程
			go SlaveSync()
			mqLog.Loger.Println("Stop push")
		}
		return
	}
	StopPush = true
}
