package database

import (
	"GaryReleaseProject/src/model"
	"fmt"
	"sync"
	"testing"
)

func TestLinkMySQL(t *testing.T){
	model.DB = model.InitDatabase()
	defer model.DB.Close()
	model.MySQLRWMutex = new(sync.RWMutex)
	r1 := model.Rule{
		// rule condition
		Platform:             "Android",
		DownloadUrl:          "https://cdn.mysql.com//Downloads/MySQLInstaller/mysql-installer-community-5.7.36.0.msi",
		UpdateVersionCode:    "8.5.3.01",
		MD5:                  "94851e04b5eb17a4fdddc48b5d62de84",
		DeviceIdList:         "001,002,003,004,005,006,007,008,009,010",
		MaxUpdateVersionCode: "8.5",
		MinUpdateVersionCode: "8.1",
		MaxOsApi:             40,
		MinOsApi:             20,
		CpuArch:              64,
		Channel:              "HUAWEI",
		Title:                "update for the new version",
		UpdateTips:           "Please make sure that your phone is fully charged before updating!",

		UpdateVersionCodeInt64:    model.VersionToInt64("8.5.3.01"),
		MaxUpdateVersionCodeInt64: model.VersionToInt64("8.5"),
		MinUpdateVersionCodeInt64: model.VersionToInt64("8.1"),
		RuleStatus:                0,
	}

	err := InsertRule(r1)
	if err != nil {
		fmt.Println("here!")
	}

}
