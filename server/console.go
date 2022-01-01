package server

import (
	"DPMQ/model"
	"DPMQ/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"sort"
	"strconv"
	"strings"
)

/**
 * 控制台服务
 */

//获取全部消费者客户端集合
func GetConsumers(c *gin.Context) {

	consumers := []model.Consumer{}
	consumer := model.Consumer{}

	//遍历消费者客户端列表
	for key := range Consumers {

		consumer.Topic = strings.Split(key, "|")[0]
		consumer.ConsumerId = strings.Split(key, "|")[1]

		consumers = append(consumers, consumer)
	}

	c.JSON(0, gin.H{
		"code": 0,
		"data": consumers,
	})
}

//获取消息队列详细配置
func GetConfig(c *gin.Context) {

	configMap := make(map[string]interface{}, 3)
	configMap["messageChanBuffer"] = viper.GetInt("mq.messageChanBuffer")
	configMap["pushMessagesSpeed"] = viper.GetInt("mq.pushMessagesSpeed")
	configMap["sendCount"] = viper.GetInt("mq.sendCount")
	configMap["sendRetryCount"] = viper.GetInt("mq.sendRetryCount")
	configMap["persistentPath"] = viper.GetString("mq.persistentPath")
	configMap["isPersistent"] = viper.GetInt("mq.isPersistent")
	configMap["recoveryStrategy"] = viper.GetInt("mq.recoveryStrategy")
	configMap["persistentTime"] = viper.GetInt("mq.persistentTime")
	configMap["rePushSpeed"] = viper.GetInt("mq.rePushSpeed")
	configMap["clearTime"] = viper.GetInt("mq.clearTime")

	c.JSON(0, gin.H{
		"code": 0,
		"data": configMap,
	})
}

//获取指定状态的消息记录列表（可用该接口进行消息记录持久化操作）
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

	var messageList []model.Message

	//返回全部状态的消息
	if intStatus == 3 {
		//遍历消息列表
		MessageList.Range(func(key, message interface{}) bool {

			ts := message.(model.Message).CreateTime

			//消息创建时间符合搜索时间
			if ts >= start && ts <= end {
				messageList = append(messageList, message.(model.Message))
			}
			return true
		})

	} else {

		//遍历消息列表
		MessageList.Range(func(key, message interface{}) bool {

			ts := message.(model.Message).CreateTime

			//消息符合搜索条件
			if message.(model.Message).Status == intStatus && ts >= start && ts <= end {
				messageList = append(messageList, message.(model.Message))
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

//统计各状态消息的数量
func CountMessage(c *gin.Context) {

	count := make(map[string]int64, 4)

	var all int64 = 0
	var consumed int64 = 0
	var notConsumed int64 = 0
	var stateless int64 = 0

	//遍历消息列表
	MessageList.Range(func(key, message interface{}) bool {

		if message.(model.Message).Status == -1 {
			stateless++
		}
		if message.(model.Message).Status == 0 {
			notConsumed++
		}
		if message.(model.Message).Status == 1 {
			consumed++
		}
		all++

		return true
	})

	count["all"] = all
	count["consumed"] = consumed
	count["notConsumed"] = notConsumed
	count["stateless"] = stateless

	c.JSON(0, gin.H{
		"code": 0,
		"data": count,
	})
}
