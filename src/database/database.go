package database

import (
	"GaryReleaseProject/src/model"
	"database/sql"
	"errors"
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
		"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
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
	result, err = model.DB.Exec("INSERT INTO white_lists("+
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

	//根据设备平台情况选择相应的匹配方法
	//Android平台带OsApi
	sqlAndroid := "SELECT rule_id, download_url, update_version_code, md5, title, update_tips " +
		"FROM config " +
		"WHERE max_update_version_code_int64 >=  ? AND min_update_version_code_int64 <= ? " +
		"AND platform = ? " +
		"AND channel = ? " +
		"AND max_os_api >= ? AND min_os_api <= ? " +
		"AND cpu_arch = ? " +
		"AND rule_status = 0 " +
		"ORDER BY update_version_code_int64 DESC "
	//ios平台不带OsApi
	sqlIos := "SELECT rule_id, download_url, update_version_code, md5, title, update_tips " +
		"FROM config " +
		"WHERE max_update_version_code_int64 >=  ? AND min_update_version_code_int64 <= ? " +
		"AND platform = ? " +
		"AND channel = ? " +
		"AND cpu_arch = ? " +
		"AND rule_status = 0 " +
		"ORDER BY update_version_code_int64 DESC "

	var rows *sql.Rows
	var err error
	if cr.DevicePlatform == "Android" {
		rows, err = model.DB.Query(sqlAndroid,
			model.VersionToInt64(cr.UpdateVersionCode), model.VersionToInt64(cr.UpdateVersionCode),
			cr.DevicePlatform, cr.Channel, cr.OsApi, cr.OsApi, cr.CpuArch)
	} else {
		rows, err = model.DB.Query(sqlIos,
			model.VersionToInt64(cr.UpdateVersionCode), model.VersionToInt64(cr.UpdateVersionCode),
			cr.DevicePlatform, cr.Channel, cr.CpuArch)
	}
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

// GetWhitelist 根据ruleid string读取mysql中的白名单
func GetWhitelist(rid string) (string, error) {
	id, err := strconv.Atoi(rid)
	if err != nil {
		fmt.Println("rule id error:", err)
		return "", err
	}
	//从主表config中检查规则状态
	var status int
	e := model.DB.QueryRow("SELECT rule_status FROM config WHERE rule_id = ?", id).Scan(&status)
	if e != nil {
		log.Fatalln(e)
	}
	if status != 0 {
		return "", errors.New("rule paused of failed\n")
	}
	//正常规则则返回白名单
	row := model.DB.QueryRow("SELECT device_id_list FROM white_lists WHERE rule_id = ?", id)
	var result string
	err = row.Scan(&result)
	if err != nil {
		fmt.Println("search for the rule id failed:", err)
		return "", err
	}
	return result, nil
}

// GetNIds 获取n个最大的白名单rule_id
func GetNIds(n int) []int {
	rows, err := model.DB.Query("SELECT rule_id FROM config WHERE rule_status = 0 ORDER BY rule_id DESC")
	if err != nil{
		log.Fatalln(err)
	}
	defer rows.Close()
	var result []int
	counter := 0
	for rows.Next(){
		var id int
		if err := rows.Scan(&id); err != nil{
			log.Fatalln(err)
		}
		if counter < n{
			result = append(result, id)
			counter++
		}else{
			break
		}
	}
	return result


	/* 未考虑status
	var count int
	err := model.DB.QueryRow("SELECT COUNT(*) FROM white_lists").Scan(&count)
	if err != nil {
		log.Fatalln(err)
	}
	if count < n {
		n = count
	}
	//建表时，rule_id是从0自增的，假定不会更改白名单表，最后一个id为count-1
	res := make([]int, n)
	for i := range res {
		res[i] = count - 1
		count--
	}
	return res*/
}
