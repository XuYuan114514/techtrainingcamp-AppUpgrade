package update_service

import (
	"GaryReleaseProject/src/cache"
	"GaryReleaseProject/src/database"
	"GaryReleaseProject/src/model"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)



func Pong(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}


func DealCReport(c *gin.Context) {

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

	cacheMessage,_ := database.MatchRules(&cr)
	returnMessage,_ := cache.MatchWhitelist(cacheMessage)

	defer c.JSON(200,returnMessage)
	defer model.MySQLRWMutex.RUnlock()
}

