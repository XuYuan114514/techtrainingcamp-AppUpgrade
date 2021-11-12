package main

import (
	"GaryReleaseProject/src/cache"
	"GaryReleaseProject/src/deploy_service"
	"GaryReleaseProject/src/model"
	"GaryReleaseProject/src/update_service"
	"github.com/gin-gonic/gin"
)

func customizeouter(r *gin.Engine) {
	model.InitAll()
	cache.WarmUpCache()
	r.GET("/ping", update_service.Pong)
	r.GET("/update", update_service.DealCReport)

	// 配置接口，可能需要HTTPS连接，先保证可用
	// 有点类似一个用户管理接口, 理想状态是管理员认证后登录，可以进行规则的更改
	r.GET("/deploy", deploy_service.DeployRule)
}



