package deploy_service

import (
	"github.com/gin-gonic/gin"
)


func DeployRule(c *gin.Context){
	// 前端接受的信息处理成Rule结构体，并写入MySQL

	/*
	rule status = 0
	UpdateVersionCodeInt64 = model.VersionToInt64( )
	...
	 */

	// 根据情况调用pkg database中的插入规则还是修改规则函数
	//database.InsertRule()
	//database.ModifyRule()
}

