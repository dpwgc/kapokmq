package cluster

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"github.com/spf13/viper"
	"time"
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
	nodePort := viper.GetInt("cluster.nodePort")

	registryAddr := viper.GetString("cluster.registryAddr")
	registryPort := viper.GetString("cluster.registryPort")

	//配置节点信息
	conf := memberlist.DefaultLANConfig()
	conf.Name = addr + ":" + port
	conf.BindPort = nodePort
	conf.AdvertisePort = nodePort

	var err error

	//创建一个节点
	list, err = memberlist.Create(conf)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	//将list加入到已存在的集群（即注册中心所在集群）
	n, err := list.Join([]string{registryAddr + ":" + registryPort})
	fmt.Println(n)
	if err != nil {
		panic("Failed to join cluster: " + err.Error())
	}

	go func() {
		for {
			SyncPush()
			time.Sleep(time.Second * 3)
		}
	}()
}
