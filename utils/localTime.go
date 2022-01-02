package utils

import "time"

var cstSh, _ = time.LoadLocation("Asia/Shanghai")

func GetLocalDateTimestamp() int64 {
	return time.Now().In(cstSh).Local().Unix()
}
