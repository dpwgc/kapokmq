package utils

import (
	"fmt"
	"github.com/go-basic/uuid"
)

// CreateCode 生成消息标识码
func CreateCode(messageData string) string {
	uu := uuid.New()
	data := fmt.Sprintf("%s%s", messageData, uu)
	tokenPrefix := Md5Sign(data)
	return tokenPrefix
}
