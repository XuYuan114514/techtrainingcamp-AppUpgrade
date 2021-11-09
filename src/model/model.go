package model

import (
	"strconv"
	"strings"
	"sync"
)

type Rule struct {
	// rule condition
	Platform string `json:"platform"`
	DownloadUrl string `json:"download_url"`
	UpdateVersionCode string `json:"update_version_code"`
	MD5 string `json:"md5"`
	DeviceIdList string `json:"device_id_list"`
	MaxUpdateVersionCode string `json:"max_update_version_code"`
	MinUpdateVersionCode string `json:"min_update_version_code"`
	MaxOsApi int `json:"max_os_api"`
	MinOsApi int `json:"min_os_api"`
	CpuArch int `json:"cpu_arch"`
	Channel string `json:"channel"`
	Title string `json:"title"`
	UpdateTips string `json:"update_tips"`
	//新添加的字段
	UpdateVersionCodeInt64 int64
	MaxUpdateVersionCodeInt64 int64
	MinUpdateVersionCodeInt64 int64
	RuleStatus int
}

type CReport struct {
	// App uploads when start
	DevicePlatform string `json:"device_platform"`
	DeviceId string `json:"device_id"`
	OsApi int `json:"os_api"`
	Channel string `json:"channel"`
	UpdateVersionCode string `json:"update_version_code"`
	CpuArch int `json:"cpu_arch"`
}

type CacheMessage struct {
	// message return ro App
	RuleId int
	DownloadUrl string `json:"download_url"`
	UpdateVersionCode string `json:"update_version_code"`
	Md5 string `json:"md5"`
	Title string `json:"title"`
	UpdateTips string `json:"update_tips"`
}

type RMessage struct {
	// message return ro App
	DownloadUrl string `json:"download_url"`
	UpdateVersionCode string `json:"update_version_code"`
	Md5 string `json:"md5"`
	Title string `json:"title"`
	UpdateTips string `json:"update_tips"`
}

var (
	MySQLRWMutex *sync.RWMutex
	//DB *sql.DB
	//redis
)

func InitAll(){
	MySQLRWMutex = new(sync.RWMutex)
	// DB = initDatabase()
	// redis = initCache()

}

func initDatabase() interface{}{
	// 使用mysql连接池实现
	return nil
}

func initCache(){
	//连接redis
	//缓存预热
}

func VersionToInt64(version string) int64{
	var res  int64 = 0
	versionSlices := strings.Split(version,".")
	bitMoves := [4]int{48,32,16,0}
	for i,v := range versionSlices{
		j,err := strconv.Atoi(v)
		if err != nil{
			panic(err)
		}
		res |= int64(j)<< bitMoves[i]
	}
	return res
}
