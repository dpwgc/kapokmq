package cluster

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/memberlist"
	"kapokmq/config"
	"kapokmq/model"
	"kapokmq/mqLog"
	"strings"
)

//集群节点列表
var list *memberlist.Memberlist

// InitCluster 向Gossip集群的注册中心注册
func InitCluster() {

	//如果是单机部署模式
	if config.Get.Cluster.IsCluster != 1 {
		return
	}

	//获取该节点的地址
	addr := config.Get.Server.Addr
	//获取该节点的Gin服务端口号
	port := config.Get.Server.Port

	//获取设置的Gossip服务端口号
	gossipPort := config.Get.Cluster.GossipPort

	registryAddr := config.Get.Cluster.RegistryAddr
	registryGossipPort := config.Get.Cluster.RegistryGossipPort

	//配置本节点信息
	conf := memberlist.DefaultLANConfig()
	//addr缺省，addr为空默认设为127.0.0.1
	if addr == "" {
		addr = "127.0.0.1"
	}

	//本节点的名称（例：mq:0.0.0.0:8011）
	conf.Name = fmt.Sprintf("%s%s%s%s", "mq:", addr, ":", port) //前缀r:表明这是注册中心，前缀mq:表明这是消息队列节点

	//Bind：Gossip服务内部注册地址（0.0.0.0:gossipPort）
	conf.BindPort = gossipPort

	//本节点对外暴露的地址（公网IP，用于在公网环境下连接注册中心）
	conf.AdvertiseAddr = addr
	conf.AdvertisePort = gossipPort

	var err error

	//申请创建一个Gossip服务节点
	list, err = memberlist.Create(conf)
	if err != nil {
		mqLog.Loger.Println("Failed to create memberlist: " + err.Error())
		panic(err)
		return
	}

	//将节点加入到已存在的集群（即注册中心所在集群）
	_, err = list.Join([]string{registryAddr + ":" + registryGossipPort})
	if err != nil {
		mqLog.Loger.Println("Failed to join cluster: " + err.Error())
		panic(err)
		return
	}
}

// GetNodes 获取除了注册中心之外的集群所有节点
func GetNodes(c *gin.Context) {

	if config.Get.Cluster.IsCluster == 0 {
		c.String(0, "")
		return
	}

	var nodes []model.Node

	// 获取当前集群的消息队列节点信息（除去注册中心）
	for _, member := range list.Members() {
		m := strings.Split(member.Name, ":")
		//如果该节点是注册中心，跳过
		if m[0] == "r" {
			continue
		}

		node := model.Node{
			Name: member.Name,
			Addr: m[1],
			Port: m[2],
		}
		nodes = append(nodes, node)
	}

	data, err := json.Marshal(nodes)
	if err != nil {
		c.String(0, "")
		return
	}
	c.String(0, string(data))
}
