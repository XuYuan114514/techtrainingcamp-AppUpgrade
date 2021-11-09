package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 配置Handler
	customizeouter(r)
	r.Run()
}
