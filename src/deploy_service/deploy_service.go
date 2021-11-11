package deploy_service

import (
	"GaryReleaseProject/src/database"
	"GaryReleaseProject/src/model"
	"github.com/gin-gonic/gin"
)


func DeployRule(c *gin.Context) {
	// 前端接受的信息处理成Rule结构体，并写入MySQL

	rule := model.Rule{
		Platform:             c.Query("platform"),
		DownloadUrl:          c.Query("download_url"),
		UpdateVersionCode:    c.Query("update_version_code"),
		MD5:                  c.Query("md5"),
		DeviceIdList:         c.Query("device_id_list"),
		MaxUpdateVersionCode: c.Query("max_update_version_code"),
		MinUpdateVersionCode: c.Query("min_update_version_code"),
		MaxOsApi:             model.Str2Int(c.Query("max_os_api")),
		MinOsApi:             model.Str2Int(c.Query("min_os_api")),
		CpuArch:              model.Str2Int(c.Query("cpu_arch")),
		Channel:              c.Query("channel"),
		Title:                c.Query("title"),
		UpdateTips:           c.Query("update_tips"),
		//新添加的字段:int64版本信息直接生成，status需要传入
		UpdateVersionCodeInt64:    model.VersionToInt64(c.Query("update_version_code")),
		MaxUpdateVersionCodeInt64: model.VersionToInt64(c.Query("max_update_version_code")),
		MinUpdateVersionCodeInt64: model.VersionToInt64(c.Query("min_update_version_code")),
		RuleStatus:                model.Str2Int(c.Query("rule_status")),
	}
	RuleId := model.Str2Int(c.Query("rule_id"))
	//根据rule_status来决定是插入数据库还是修改数据库
	if rule.RuleStatus == 0 {
		_ = database.InsertRule(rule)
		c.JSON(200, gin.H{"message": "Insert Rule"})
	} else {
		_ = database.ModifyRule(RuleId, rule.RuleStatus)
		c.JSON(200, gin.H{"message": "Modify Rule"})
	}

}

