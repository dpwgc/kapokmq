package utils

import "time"

//时间戳转换

func ToTimestamp(strTime string) int64 {
	timeLayout := "2006-01-02 15:04:05"  //转化模板
	loc, _ := time.LoadLocation("Local") //获取时区
	theTime, _ := time.ParseInLocation(timeLayout, strTime, loc)
	ts := theTime.Unix() //时间戳
	return ts
}
