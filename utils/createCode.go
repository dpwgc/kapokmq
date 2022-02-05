package utils

import (
	"fmt"
	"github.com/satori/go.uuid"
)

// CreateCode 生成消息标识码
func CreateCode(messageData string) string {
	uu := uuid.NewV4()
	data := fmt.Sprintf("%s%s", messageData, uu)
	tokenPrefix := Md5Sign(data)
	return tokenPrefix
}
