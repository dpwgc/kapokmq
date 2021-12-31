package console

import (
	"DPMQ/model"
	"DPMQ/server"
	"DPMQ/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"sort"
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

	//获取消息列表
	//搜索区间
	startTime, _ := c.GetQuery("startTime")
	endTime, _ := c.GetQuery("endTime")

	startDate, _ := c.GetQuery("startDate")
	endDate, _ := c.GetQuery("endDate")

	//状态（1：已消费消息，0：未消费消息，3：全部消息）
	status, _ := c.GetQuery("status")
	intStatus, _ := strconv.Atoi(status)

	startTime = startDate + " " + startTime
	endTime = endDate + " " + endTime

	//转成时间戳
	start := utils.ToTimestamp(startTime)
	end := utils.ToTimestamp(endTime)

	var messageList []model.Message

	//返回全部状态的消息
	if intStatus == 3 {
		//遍历消息列表
		server.MessageList.Range(func(key, message interface{}) bool {

			ts := message.(model.Message).CreateTime

			//消息创建时间符合搜索时间
			if ts >= start && ts <= end {
				messageList = append(messageList, message.(model.Message))
			}
			return true
		})

	} else {

		//遍历消息列表
		server.MessageList.Range(func(key, message interface{}) bool {

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

	c.HTML(http.StatusOK, "Index.html", gin.H{
		"consumers":   consumers,
		"configMap":   configMap,
		"messageList": messageList,
		"tip":         "From " + startTime + " To " + endTime,
	})
}
