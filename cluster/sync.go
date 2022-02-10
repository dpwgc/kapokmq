package cluster

import (
	"github.com/gorilla/websocket"
	"kapokmq/config"
	"net/http"
	"time"
)

// SyncConn 主从同步连接
var SyncConn *websocket.Conn

// UpGrader websocket跨域配置
var UpGrader = websocket.Upgrader{
	//跨域设置
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// InitSync 主从同步
func InitSync() {

	//是否开启主从同步功能
	if config.Get.Sync.IsSync != 1 {
		return
	}

	//如果是从节点
	if config.Get.Sync.IsSlave == 1 {

		//建立主从连接
		SlaveConn()
		//启动从节点同步协程
		go SlaveSync()

		//开启连接检查协程
		go func() {
			for {
				//每隔一秒进行一次连接检查
				time.Sleep(time.Second * 1)
				checkConn()
			}
		}()
	}
}
