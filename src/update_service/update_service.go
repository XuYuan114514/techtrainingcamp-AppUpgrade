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
	// go有类似shared_ptr，结束了有人用也不回收，变量逃逸分析
	cacheMessage,_ := database.MatchRules(&cr)
	pos,_ := cache.MatchWhitelist(cacheMessage,cr.DeviceId)

	model.MySQLRWMutex.RUnlock()
	if pos == -1{
		c.JSON(200,"")
	}else{
		c.JSON(200,gin.H{
			"download_url": (*cacheMessage)[pos].DownloadUrl,
			"update_version_code":(*cacheMessage)[pos].UpdateVersionCode,
			"md5":(*cacheMessage)[pos].Md5,
			"title":(*cacheMessage)[pos].Title,
			"update_tips":(*cacheMessage)[pos].UpdateTips,
			})
	}
}

