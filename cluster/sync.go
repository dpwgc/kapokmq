package cluster

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"kapokmq/server"
)

// SyncPush 集群消息同步推送
func SyncPush(message []byte) {
	// 获取当前集群的节点
	for _, member := range list.Members() {
		//协程并发推送消息到所有节点
		go func(member *memberlist.Node) {
			fmt.Printf("Member: %s %s\n", member.Name, member.Addr)
			//将消息发送到该节点
			err := list.SendReliable(member, message)
			if err != nil {
				server.Loger.Println(err)
				return
			}
		}(member)
	}
}

// SyncGet 接收其他节点的同步消息
func SyncGet(message []byte) {

}
