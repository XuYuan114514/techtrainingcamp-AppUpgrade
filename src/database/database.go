package database

import (
	"GaryReleaseProject/src/model"
	"fmt"
	"log"
	"strconv"
)

// InsertRule 可能传指针好一点，因为规则可能较大
func InsertRule(rule model.Rule) error {
	model.MySQLRWMutex.Lock()
	// 通过线程池连接数据库

	//插入config表，不带rule_id,使用自增的id
	result, err := model.DB.Exec("INSERT INTO config("+
		"platform, download_url, update_version_code, md5,"+
		"max_update_version_code, min_update_version_code, max_os_api, min_os_api,"+
		"cpu_arch, channel, title, update_tips,"+
		"update_version_code_int64, max_update_version_code_int64, min_update_version_code_int64, rule_status)VALUES("+
		"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		rule.Platform, rule.DownloadUrl, rule.UpdateVersionCode, rule.MD5,
		rule.MaxUpdateVersionCode, rule.MinUpdateVersionCode, rule.MaxOsApi, rule.MinOsApi,
		rule.CpuArch, rule.Channel, rule.Title, rule.UpdateTips,
		rule.UpdateVersionCodeInt64, rule.MaxUpdateVersionCodeInt64, rule.MinUpdateVersionCodeInt64, rule.RuleStatus)
	if err != nil {
		fmt.Println("insert to config table failed!", err)
		return err
	}
	lastId, e := result.LastInsertId()
	if e != nil {
		log.Fatal(err)
	}
	fmt.Println("insert to config table succeed!", lastId)

	//白名单单独插入white_lists表,同样不带rule_id,使用自增id
	result, err = model.DB.Exec("INSERT INTO white_list("+
		"device_id_list)VALUES(?)", rule.DeviceIdList)
	if err != nil {
		fmt.Println("insert to config table failed!", err)
		return err
	}
	lastId, e = result.LastInsertId()
	if e != nil {
		log.Fatal(err)
	}
	fmt.Println("insert to white_lists table succeed!", lastId)

	// 结束时归还线程池
	defer model.MySQLRWMutex.Unlock()
	return nil
}

func ModifyRule(ruleId int, status int) error {
	model.MySQLRWMutex.Lock()
	// 通过线程池连接数据库

	res, err := model.DB.Exec("UPDATE config SET rule_status = ? WHERE rule_id = ?", status, ruleId)
	if err != nil {
		fmt.Println("modify the rule failed", err)
		return err
	}
	row, e := res.RowsAffected()
	if e != nil {
		fmt.Println("rows failed", err)
		return e
	}
	fmt.Println("update the rule succeed", row)

	// 结束时归还线程池
	defer model.MySQLRWMutex.Unlock()
	return nil
}

func MatchRules(cr *model.CReport) (*[]model.CacheMessage, error) {
	// 查询其他规则表中所有匹配的rule
	// 要求返回的切片中rule_id对应规则按照UpdateVersionCodeInt64降序排列

	rows, err := model.DB.Query("SELECT rule_id, download_url, update_version_code, md5, title, update_tips "+
		"FROM config "+
		"WHERE max_update_version_code_int64 >=  ? AND min_update_version_code_int64 <= ? "+
		"AND platform = ? "+
		"AND channel = ? "+
		"AND max_os_api >= ? AND min_os_api <= ? "+
		"AND cpu_arch = ? "+
		"ORDER BY update_version_code_int64 DESC ",
		model.VersionToInt64(cr.UpdateVersionCode), model.VersionToInt64(cr.UpdateVersionCode),
		cr.DevicePlatform, cr.Channel, cr.OsApi, cr.OsApi, cr.CpuArch)
	//关闭rows释放所持有的数据库链接
	defer rows.Close()
	if err != nil {
		fmt.Println("query failed, err:", err)
		return nil, err
	}

	var result []model.CacheMessage
	for rows.Next() {
		var temp model.CacheMessage
		err := rows.Scan(&temp.RuleId, &temp.DownloadUrl, &temp.UpdateVersionCode,
			&temp.Md5, &temp.Title, &temp.UpdateTips)
		if err != nil {
			fmt.Println("Scan failed", err)
			return nil, err
		}
		result = append(result, temp)
	}
	return &result, nil
	// 通过线程池连接数据库
	// 结束时归还线程池

}

func GetWhitelist(rid string) (string, error) {
	id, err := strconv.Atoi(rid)
	if err != nil {
		fmt.Println("rule id error:", err)
		return "", err
	}
	row := model.DB.QueryRow("SELECT device_id_list FROM white_list WHERE rule_id = ?", id)
	var result string
	err = row.Scan(&result)
	if err != nil {
		fmt.Println("search for the rule id failed:", err)
		return "", err
	}
	return result, nil
}
