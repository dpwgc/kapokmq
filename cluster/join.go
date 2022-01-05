package cluster

import (
	"flag"
	"github.com/hashicorp/memberlist"
	"github.com/spf13/viper"
)

//集群节点列表
var list *memberlist.Memberlist

// InitCluster Gossip集群注册
func InitCluster() {

	//获取集群名称
	clusterName := viper.GetString("cluster.clusterName")

	//获取该节点的ip地址及端口号
	ip := viper.GetString("server.ip")
	port := viper.GetString("server.port")

	node := flag.String("node", ip+":"+port, "mq node")
	cluster := flag.String("cluster", clusterName, "add exist cluster")

	flag.Parse()
	conf := memberlist.DefaultLANConfig()
	conf.Name = *node
	conf.BindAddr = *node

	//创建一个节点
	list, err := memberlist.Create(conf)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	//将list加入到已存在的集群.
	if *cluster == "" {
		_, err := list.Join([]string{*cluster})
		if err != nil {
			panic("Failed to join cluster: " + err.Error())
		}
	}
}
