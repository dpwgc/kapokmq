package utils

import (
	"fmt"
	"time"
)

//生成消息标识码

func CreateCode(messageData string) string {
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := Md5Sign(messageData + ts + "_DPMQ")
	return tokenPrefix + ts[:8]
}
