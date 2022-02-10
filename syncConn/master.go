package syncConn

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"kapokmq/config"
	"kapokmq/model"
	"kapokmq/mqLog"
	"sync"
)

var SendLock sync.Mutex

// Master 主节点同步协程
func Master(c *gin.Context) {

	//升级get请求为webSocket协议
	ws, err := UpGrader.Upgrade(c.Writer, c.Request, nil)

	defer func(ws *websocket.Conn) {
		err := ws.Close()
		Conn = nil
		if err != nil {
			mqLog.Loger.Println(err)
		}
	}(ws)

	Conn = ws

	//登录验证
	for {
		//连接成功，等待从节点输入访问密钥
		err = ws.WriteMessage(1, []byte("Please enter the secret key"))
		if err != nil {
			mqLog.Loger.Println(err)
			return
		}

		//读取ws中的数据，获取访问密钥
		_, sk, err := ws.ReadMessage()
		if err != nil {
			mqLog.Loger.Println(err)
			return
		}
		if string(sk) == config.Get.Mq.SecretKey {
			//访问密钥匹配成功
			err = ws.WriteMessage(1, []byte("Secret key matching succeeded"))
			if err != nil {
				mqLog.Loger.Println(err)
				return
			}
			break
		}

		//访问密钥匹配失败
		err = ws.WriteMessage(1, []byte("Secret key matching error"))
		if err != nil {
			mqLog.Loger.Println(err)
			return
		}
	}

	for {
		//读取从节点发来的心跳检测
		_, _, err = ws.ReadMessage()
		if err != nil {
			mqLog.Loger.Println(err)
			return
		}
	}
}

// SendMessage 主从同步，向从节点发送消息
func SendMessage(message model.Message) {

	SendLock.Lock()
	err := Conn.WriteJSON(message)
	SendLock.Unlock()

	if err != nil {
		mqLog.Loger.Println(err)
		return
	}
}
