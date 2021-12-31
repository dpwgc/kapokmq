package utils

import "time"

var cstSh, _ = time.LoadLocation("Asia/Shanghai")

func GetLocalDateTime() string {
	return time.Now().In(cstSh).Local().Format("2006-01-02 15:04:05")
}

func GetLocalDateTimestamp() int64 {
	return time.Now().In(cstSh).Local().Unix()
}
