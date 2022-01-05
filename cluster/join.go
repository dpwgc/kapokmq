package cluster

import (
	"flag"
	"github.com/hashicorp/memberlist"
)

//集群节点列表
var list *memberlist.Memberlist

// InitCluster Gossip集群注册
func InitCluster() {

	node := flag.String("node", "127.0.0.1", "mq node")
	cluster := flag.String("cluster", "1.2.3.4", "add exist cluster")
	flag.Parse()
	conf := memberlist.DefaultLANConfig()
	conf.Name = *node
	conf.BindAddr = *node

	//创建一个节点
	list, err := memberlist.Create(conf)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	// 将list加入到已存在的集群.
	if *cluster == "" {
		_, err := list.Join([]string{*cluster})
		if err != nil {
			panic("Failed to join cluster: " + err.Error())
		}
	}
}
