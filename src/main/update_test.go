package main

import (
	"GaryReleaseProject/src/cache"
	"GaryReleaseProject/src/database"
	"GaryReleaseProject/src/model"
	"GaryReleaseProject/src/update_service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"testing"
)


func TestUpdateService(t *testing.T) {
	r := gin.Default()
	model.InitAll()
	r.GET("/ping", update_service.Pong)
	r.GET("/update", cr2struct)
	r.Run()

}
func cr2struct(c *gin.Context) {
	cr:= model.CReport{
		DevicePlatform : c.Query("device_platform"),
		DeviceId : c.Query("device_id"),
		// ToInt会把空串转为0
		OsApi : cast.ToInt(c.Query("os_api")),
		Channel : c.Query("channel"),
		UpdateVersionCode : c.Query("update_version_code"),
		CpuArch : cast.ToInt(c.Query("cpu_arch")),
	}
	model.MySQLRWMutex.RLock()
	// go有类似shared_ptr，结束了有人用也不回收，变量逃逸分析
	cacheMessage,_ := database.MatchRules(&cr)
	pos,_ := cache.MatchWhitelist(cacheMessage,cr.DeviceId)
	model.MySQLRWMutex.RUnlock()
	c.JSON(200,gin.H{"pos":pos})
}

