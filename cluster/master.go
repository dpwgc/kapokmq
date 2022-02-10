package cluster

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"kapokmq/config"
	"kapokmq/model"
	"kapokmq/server"
	"time"
)

// Master 主节点同步协程
func Master(c *gin.Context) {

	defer func(SyncConn *websocket.Conn) {
		err := SyncConn.Close()
		if err != nil {
			server.Loger.Println(err)
		}
	}(SyncConn)

	//升级get请求为webSocket协议
	SyncConn, err := UpGrader.Upgrade(c.Writer, c.Request, nil)

	//登录验证
	for {
		//连接成功，等待从节点输入访问密钥
		err = SyncConn.WriteMessage(1, []byte("Please enter the secret key"))
		if err != nil {
			server.Loger.Println(err)
			return
		}

		//读取ws中的数据，获取访问密钥
		_, sk, err := SyncConn.ReadMessage()
		if err != nil {
			server.Loger.Println(err)
			return
		}
		if string(sk) == config.Get.Mq.SecretKey {
			//访问密钥匹配成功
			err = SyncConn.WriteMessage(1, []byte("Secret key matching succeeded"))
			if err != nil {
				server.Loger.Println(err)
				return
			}
			break
		}

		//访问密钥匹配失败
		err = SyncConn.WriteMessage(1, []byte("Secret key matching error"))
		if err != nil {
			server.Loger.Println(err)
			return
		}
	}

	//挂起连接
	for {
		time.Sleep(time.Second * 10)
	}
}

// SendMessage 主从同步，向从节点发送消息
func SendMessage(message model.Message) {
	err := SyncConn.WriteJSON(message)
	if err != nil {
		server.Loger.Println(err)
		return
	}
}
