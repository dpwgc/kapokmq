package cluster

import (
	"fmt"
	"strconv"
)

// SyncPush 集群消息同步推送
func SyncPush() {
	// 获取当前集群的节点
	for _, member := range list.Members() {
		//推送消息到所有节点
		fmt.Printf("Member: %s %s %s\n", member.Name, member.Addr, strconv.Itoa(list.NumMembers()))
		//将消息发送到该节点
	}
}

// SyncGet 接收其他节点的同步消息
func SyncGet(message []byte) {

}
