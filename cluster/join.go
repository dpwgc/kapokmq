package cluster

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"github.com/spf13/viper"
	"kapokmq/server"
)

//集群节点列表
var list *memberlist.Memberlist

// InitCluster 向Gossip集群的注册中心注册
func InitCluster() {

	//如果是单机部署模式
	if viper.GetInt("cluster.isCluster") != 1 {
		return
	}

	//获取该节点的地址
	addr := viper.GetString("server.addr")
	//获取该节点的Gin服务端口号
	port := viper.GetString("server.port")

	//获取设置的Gossip服务端口号
	gossipPort := viper.GetInt("cluster.gossipPort")

	registryAddr := viper.GetString("cluster.registryAddr")
	registryGossipPort := viper.GetString("cluster.registryGossipPort")

	//配置本节点信息
	conf := memberlist.DefaultLANConfig()
	//addr缺省，addr为空默认设为0.0.0.0
	if addr == "" {
		addr = "0.0.0.0"
	}
	//本节点名称
	conf.Name = fmt.Sprintf("%s%s%s%s", "mq:", addr, ":", port) //前缀r:表明这是注册中心，前缀mq-表明这是消息队列节点
	//本节点的地址
	conf.BindAddr = addr
	//本节点Gossip服务端口号
	conf.BindPort = gossipPort
	conf.AdvertisePort = gossipPort

	var err error

	//申请创建一个Gossip服务节点
	list, err = memberlist.Create(conf)
	if err != nil {
		server.Loger.Println("Failed to create memberlist: " + err.Error())
		panic(err)
		return
	}

	//将节点加入到已存在的集群（即注册中心所在集群）
	_, err = list.Join([]string{registryAddr + ":" + registryGossipPort})
	if err != nil {
		server.Loger.Println("Failed to join cluster: " + err.Error())
		panic(err)
		return
	}
}
