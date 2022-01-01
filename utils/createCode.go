package utils

import (
	"fmt"
	"github.com/go-basic/uuid"
	"time"
)

//生成消息标识码

func CreateCode(messageData string) string {
	ts := fmt.Sprintf("%x", time.Now().Unix())
	uuid := uuid.New()
	tokenPrefix := Md5Sign(messageData + ts + uuid)
	return tokenPrefix + ts[:8]
}
