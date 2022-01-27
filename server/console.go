package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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

// GetConsumers 获取全部消费者客户端集合
func GetConsumers(c *gin.Context) {

	var consumers []model.Consumer

	//遍历获取消费者客户端列表
	for key := range Consumers {

		consumers = append(consumers, key)
	}

	c.JSON(0, gin.H{
		"code": 0,
		"data": consumers,
	})
}

// GetProducers 获取全部生产者客户端集合
func GetProducers(c *gin.Context) {

	var producers []model.Producer

	//遍历获取生产者客户端列表
	for key := range Producers {

		producers = append(producers, key)
	}

	c.JSON(0, gin.H{
		"code": 0,
		"data": producers,
	})
}

// GetConfig 获取消息队列详细配置
func GetConfig(c *gin.Context) {

	configMap := make(map[string]interface{}, 18)

	configMap["pushType"] = viper.GetInt("mq.pushType")

	configMap["messageChanBuffer"] = viper.GetInt("mq.messageChanBuffer")

	configMap["pushMessagesSpeed"] = viper.GetInt("mq.pushMessagesSpeed")
	configMap["pushCount"] = viper.GetInt("mq.pushCount")
	configMap["pushRetryCount"] = viper.GetInt("mq.pushRetryCount")

	configMap["persistentFile"] = viper.GetString("mq.persistentFile")
	configMap["isPersistent"] = viper.GetInt("mq.isPersistent")
	configMap["recoveryStrategy"] = viper.GetInt("mq.recoveryStrategy")
	configMap["persistentTime"] = viper.GetInt("mq.persistentTime")

	configMap["isCleanConsumed"] = viper.GetInt("mq.isCleanConsumed")

	configMap["isRePush"] = viper.GetInt("mq.isRePush")
	configMap["isClean"] = viper.GetInt("mq.isClean")

	configMap["checkSpeed"] = viper.GetInt("mq.checkSpeed")

	configMap["cleanTime"] = viper.GetInt("mq.cleanTime")

	configMap["isCluster"] = viper.GetInt("cluster.isCluster")
	configMap["gossipPort"] = viper.GetInt("cluster.gossipPort")
	configMap["registryAddr"] = viper.GetString("cluster.registryAddr")
	configMap["registryGossipPort"] = viper.GetString("cluster.registryGossipPort")

	c.JSON(0, gin.H{
		"code": 0,
		"data": configMap,
	})
}

// GetMessageList 获取指定状态的消息记录列表
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

			//搜索全部状态的消息-搜索全部主题的消息-消息符合搜索条件
			if intStatus == 3 && ts >= start && ts <= end {
				messageList = append(messageList, message.(model.Message))
				return true
			}

			//搜索指定状态的消息-搜索全部主题的消息-消息符合搜索条件
			if message.(model.Message).Status == intStatus && ts >= start && ts <= end {
				messageList = append(messageList, message.(model.Message))
				return true
			}
			return true
		})

	} else {

		//遍历消息列表
		MessageList.Range(func(key, message interface{}) bool {

			ts := message.(model.Message).CreateTime

			//搜索全部状态的消息-搜索指定主题的消息-消息符合搜索条件
			if intStatus == 3 && ts >= start && ts <= end && message.(model.Message).Topic == topic {
				messageList = append(messageList, message.(model.Message))
				return true
			}

			//搜索指定状态的消息-搜索指定主题的消息-消息符合搜索条件
			if message.(model.Message).Status == intStatus && ts >= start && ts <= end && message.(model.Message).Topic == topic {
				messageList = append(messageList, message.(model.Message))
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

// GetMessageEasy 获取指定状态的简易消息记录列表（不包含消息主体，用于绘制折线图）
func GetMessageEasy(c *gin.Context) {

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

// CountMessage 统计各状态消息的数量
func CountMessage(c *gin.Context) {

	count := make(map[string]int64, 4)

	var all int64 = 0
	var consumed int64 = 0
	var failure int64 = 0
	var unconsumed int64 = 0

	//遍历消息列表
	MessageList.Range(func(key, message interface{}) bool {

		if message.(model.Message).Status == -1 {
			unconsumed++
		}
		if message.(model.Message).Status == 0 {
			failure++
		}
		if message.(model.Message).Status == 1 {
			consumed++
		}
		all++

		return true
	})

	count["all"] = all
	count["consumed"] = consumed
	count["failure"] = failure
	count["unconsumed"] = unconsumed

	c.JSON(0, gin.H{
		"code": 0,
		"data": count,
	})
}
