package console

import (
	"DPMQ/model"
	"DPMQ/server"
	"DPMQ/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"strings"
)

/**
 * 控制台服务
 */

//获取首页数据
func Index(c *gin.Context) {

	consumers := []model.Consumer{}
	consumer := model.Consumer{}

	//获取消费者客户端列表
	for key := range server.Consumers {

		consumer.Topic = strings.Split(key, "|")[0]
		consumer.ConsumerId = strings.Split(key, "|")[1]

		consumers = append(consumers, consumer)
	}

	//获取消息队列详细配置
	configMap := make(map[string]interface{}, 3)
	configMap["messageChanBuffer"] = viper.GetInt("mq.messageChanBuffer")
	configMap["pushMessagesSpeed"] = viper.GetInt("mq.pushMessagesSpeed")
	configMap["sendCount"] = viper.GetInt("mq.sendCount")
	configMap["sendRetryCount"] = viper.GetInt("mq.sendRetryCount")
	configMap["persistentPath"] = viper.GetString("mq.persistentPath")
	configMap["isPersistent"] = viper.GetInt("mq.isPersistent")
	configMap["recoveryStrategy"] = viper.GetInt("mq.recoveryStrategy")
	configMap["persistentTime"] = viper.GetInt("mq.persistentTime")

	c.HTML(http.StatusOK, "Index.html", gin.H{
		"consumers": consumers,
		"configMap": configMap,
	})
}

//获取指定状态的消息记录列表（可用该接口进行消息记录持久化操作）
func GetMessageListPage(c *gin.Context) {

	//搜索区间
	startTime, _ := c.GetQuery("startTime")
	endTime, _ := c.GetQuery("endTime")

	//转成时间戳
	start := utils.ToTimestamp(startTime)
	end := utils.ToTimestamp(endTime)

	status, _ := c.GetPostForm("status")
	intStatus, _ := strconv.Atoi(status)

	var messageList []model.Message

	//遍历消息列表
	server.MessageList.Range(func(key, message interface{}) bool {

		ts := utils.ToTimestamp(message.(model.Message).CreateTime)

		//消息符合搜索条件
		if message.(model.Message).Status == intStatus && ts >= start && ts <= end {
			messageList = append(messageList, message.(model.Message))
		}
		return true
	})

	c.JSON(0, gin.H{
		"code": 0,
		"data": messageList,
	})
}

//获取所有状态的消息记录列表（可用该接口进行消息记录持久化操作）
func GetAllMessageListPage(c *gin.Context) {

	//搜索区间
	startTime, _ := c.GetQuery("startTime")
	endTime, _ := c.GetQuery("endTime")

	fmt.Println(startTime)

	startTime = strings.Join(strings.Split(startTime, "%"), " ")
	endTime = strings.Join(strings.Split(endTime, "%"), " ")

	//转成时间戳
	start := utils.ToTimestamp(startTime)
	end := utils.ToTimestamp(endTime)

	var messageList []model.Message

	//遍历消息列表
	server.MessageList.Range(func(key, message interface{}) bool {

		ts := utils.ToTimestamp(message.(model.Message).CreateTime)

		//消息创建时间符合搜索时间
		if ts >= start && ts <= end {
			messageList = append(messageList, message.(model.Message))
		}
		return true
	})

	c.HTML(http.StatusOK, "GetAllMessageList.html", gin.H{
		"messageList": messageList,
	})
}
