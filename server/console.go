package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kapokmq/config"
	"kapokmq/model"
	"kapokmq/utils"
	"sort"
	"strconv"
)

/**
 * 控制台服务接口
 */

// Ping 检查MQ是否在运行
func Ping(c *gin.Context) {

	c.JSON(0, gin.H{
		"code": 0,
	})
}

// GetClients 获取全部客户端集合
func GetClients(c *gin.Context) {

	var consumers []model.Consumer
	var producers []model.Producer

	//遍历获取消费者客户端列表
	for key := range Consumers {

		consumers = append(consumers, key)
	}

	//遍历获取生产者客户端列表
	for key := range Producers {

		producers = append(producers, key)
	}

	c.JSON(0, gin.H{
		"code":      0,
		"consumers": consumers,
		"producers": producers,
	})
}

// GetConfig 获取消息队列详细配置
func GetConfig(c *gin.Context) {

	configMap := make(map[string]interface{}, 18)

	configMap["pushType"] = config.Get.Mq.PushType

	configMap["messageChanBuffer"] = config.Get.Mq.MessageChanBuffer

	configMap["pushMessagesSpeed"] = config.Get.Mq.PushMessagesSpeed
	configMap["pushCount"] = config.Get.Mq.PushCount
	configMap["pushRetryTime"] = config.Get.Mq.PushRetryTime

	configMap["persistentFile"] = config.Get.Mq.PersistentFile
	configMap["isPersistent"] = config.Get.Mq.IsPersistent
	configMap["recoveryStrategy"] = config.Get.Mq.RecoveryStrategy
	configMap["persistentTime"] = config.Get.Mq.PersistentTime

	configMap["isCleanConsumed"] = config.Get.Mq.IsCleanConsumed

	configMap["isRePush"] = config.Get.Mq.IsRePush
	configMap["isClean"] = config.Get.Mq.IsClean

	configMap["checkSpeed"] = config.Get.Mq.CheckSpeed

	configMap["cleanTime"] = config.Get.Mq.CleanTime

	configMap["isCluster"] = config.Get.Cluster.IsCluster
	configMap["gossipPort"] = config.Get.Cluster.GossipPort
	configMap["registryAddr"] = config.Get.Cluster.RegistryAddr
	configMap["registryGossipPort"] = config.Get.Cluster.RegistryGossipPort

	c.JSON(0, gin.H{
		"code": 0,
		"data": configMap,
	})
}

// GetMessageList 获取指定时间指定状态的消息记录列表（不包含消息主体）
func GetMessageList(c *gin.Context) {

	//搜索区间
	startTime, _ := c.GetPostForm("startTime")
	endTime, _ := c.GetPostForm("endTime")

	//转成时间戳
	start := utils.ToTimestamp(startTime)
	end := utils.ToTimestamp(endTime)

	//状态（-1：无状态消息，1：已消费消息，0：未消费消息，3：全部消息）
	status, _ := c.GetPostForm("status")
	intStatus, _ := strconv.Atoi(status)

	//主题
	topic, _ := c.GetPostForm("topic")

	var messageList []model.Message

	//如果主题为空-返回全部主题的消息
	if len(topic) == 0 {
		//遍历消息列表
		MessageList.Range(func(key, message interface{}) bool {

			ts := message.(model.Message).CreateTime

			//删除消息主体内容
			msg := message.(model.Message)
			msg.MessageData = ""

			//搜索全部状态的消息-搜索全部主题的消息-消息符合搜索条件
			if intStatus == 3 && ts >= start && ts <= end {
				messageList = append(messageList, msg)
				return true
			}

			//搜索指定状态的消息-搜索全部主题的消息-消息符合搜索条件
			if message.(model.Message).Status == intStatus && ts >= start && ts <= end {
				messageList = append(messageList, msg)
				return true
			}
			return true
		})

	} else {

		//遍历消息列表
		MessageList.Range(func(key, message interface{}) bool {

			ts := message.(model.Message).CreateTime

			//删除消息主体内容
			msg := message.(model.Message)
			msg.MessageData = ""

			//搜索全部状态的消息-搜索指定主题的消息-消息符合搜索条件
			if intStatus == 3 && ts >= start && ts <= end && message.(model.Message).Topic == topic {
				messageList = append(messageList, msg)
				return true
			}

			//搜索指定状态的消息-搜索指定主题的消息-消息符合搜索条件
			if message.(model.Message).Status == intStatus && ts >= start && ts <= end && message.(model.Message).Topic == topic {
				messageList = append(messageList, msg)
				return true
			}
			return true
		})
	}

	//消息列表排序 按创建时间降序 由大到小
	sort.SliceStable(messageList, func(i int, j int) bool {
		return messageList[i].CreateTime > messageList[j].CreateTime
	})

	c.JSON(0, gin.H{
		"code": 0,
		"data": messageList,
	})
}

// GetMessage 根据消息标识码获取消息
func GetMessage(c *gin.Context) {
	messageCode, _ := c.GetPostForm("messageCode")

	message, _ := MessageList.Load(messageCode)

	fmt.Println(message)

	c.JSON(0, gin.H{
		"code": 0,
		"data": message.(model.Message),
	})
}

// CountMessage 统计各状态消息的数量
func CountMessage(c *gin.Context) {

	count := make(map[string]int64, 4)

	var all int64 = 0
	var consumed int64 = 0
	var delay int64 = 0
	var unconsumed int64 = 0

	//遍历消息列表
	MessageList.Range(func(key, message interface{}) bool {

		if message.(model.Message).Status == -1 {
			unconsumed++
		}
		if message.(model.Message).Status == 0 {
			delay++
		}
		if message.(model.Message).Status == 1 {
			consumed++
		}
		all++

		return true
	})

	count["all"] = all
	count["consumed"] = consumed
	count["delay"] = delay
	count["unconsumed"] = unconsumed

	c.JSON(0, gin.H{
		"code": 0,
		"data": count,
	})
}
