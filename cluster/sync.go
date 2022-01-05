package cluster

import (
	"fmt"
)

//集群消息同步
func sync() {
	//同步工作协程
	go func() {
		// 获取当前集群的节点
		for _, member := range list.Members() {
			fmt.Printf("Member: %s %s\n", member.Name, member.Addr)
		}
	}()
}
