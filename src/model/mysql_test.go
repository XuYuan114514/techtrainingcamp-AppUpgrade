package model

import (
	"GaryReleaseProject/src/database"
	"testing"
)

func TestLinkMySQL(t *testing.T){
	sql := initDatabase()
	defer sql.Close()
	r1 := Rule{
		// rule condition
		Platform : "IOS",
		DownloadUrl :"/asa/s1s1s1/s12s21s2",
		UpdateVersionCode: "1.9.8.2",
		MD5                 :"sid12sd1w13dc",
		DeviceIdList        :"00000000001",
		MaxUpdateVersionCode :"1.9.8",
		MinUpdateVersionCode:"1.8.0",
		MaxOsApi             :15,
		MinOsApi            :10,
		CpuArch             :32,
		Channel             :"",
		Title               :"",
		UpdateTips          :"",
		//新添加的字段
		UpdateVersionCodeInt64    :123,
		MaxUpdateVersionCodeInt64 :120,
		MinUpdateVersionCodeInt64 :110,
		RuleStatus                : 0, //需要传入
	}
	database.InsertRule(r1)

}
