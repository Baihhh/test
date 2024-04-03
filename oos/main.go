package main

// 存储相关功能的引入包只有这两个，后面不再赘述
import (
	"fmt"

	"github.com/qiniu/go-sdk/v7/auth"
)

func main() {
	accessKey := "your access key"
	secretKey := "your secret key"
	mac := auth.New(accessKey, secretKey)
	fmt.Printf(mac.AccessKey)
}
