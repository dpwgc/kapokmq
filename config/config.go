package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Conf struct {
	Server  Server  `yaml:"server"`
	Mq      Mq      `yaml:"mq"`
	Cluster Cluster `yaml:"cluster"`
	Sync    Sync    `yaml:"sync"`
}
type Server struct {
	Addr string `yaml:"addr"`
	Port string `yaml:"port"`
}
type Mq struct {
	SecretKey         string `yaml:"secretKey"`
	PushType          int    `yaml:"pushType"`
	MessageChanBuffer int    `yaml:"messageChanBuffer"`
	PushMessagesSpeed int    `yaml:"pushMessagesSpeed"`
	PushCount         int    `yaml:"pushCount"`
	PushRetryTime     int64  `yaml:"pushRetryTime"`
	PersistentFile    string `yaml:"persistentFile"`
	IsPersistent      int    `yaml:"isPersistent"`
	RecoveryStrategy  int    `yaml:"recoveryStrategy"`
	PersistentTime    int    `yaml:"persistentTime"`
	IsCleanConsumed   int    `yaml:"isCleanConsumed"`
	IsRePush          int    `yaml:"isRePush"`
	IsClean           int    `yaml:"isClean"`
	CheckSpeed        int    `yaml:"checkSpeed"`
	CleanTime         int64  `yaml:"cleanTime"`
}
type Cluster struct {
	IsCluster          int    `yaml:"isCluster"`
	GossipPort         int    `yaml:"gossipPort"`
	RegistryAddr       string `yaml:"registryAddr"`
	RegistryGossipPort string `yaml:"registryGossipPort"`
}
type Sync struct {
	IsSync         int    `yaml:"isSync"`
	IsSlave        int    `yaml:"isSlave"`
	MasterAddr     string `yaml:"masterAddr"`
	MasterPort     string `yaml:"masterPort"`
	MasterProtocol string `yaml:"masterProtocol"`
}

var Get Conf

// InitConfig 初始化项目配置
func InitConfig() {
	yamlFile, err := ioutil.ReadFile("application.yaml")
	if err != nil {
		panic(err)
	} // 将读取的yaml文件解析为响应的 struct
	err = yaml.Unmarshal(yamlFile, &Get)
	if err != nil {
		panic(err)
	}
}
