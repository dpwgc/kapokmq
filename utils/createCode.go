package utils

import (
	"fmt"
	"github.com/go-basic/uuid"
	"time"
)

//生成消息标识码

func CreateCode(messageData string) string {
	ts := fmt.Sprintf("%x", time.Now().Unix())
	uu := uuid.New()
	data := fmt.Sprintf("%s%s%s", messageData, ts, uu)
	tokenPrefix := Md5Sign(data)
	return tokenPrefix + ts[:8]
}
