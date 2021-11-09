package database

import (
	"GaryReleaseProject/src/model"
)


func InsertRule(rule model.Rule)error{
	model.MySQLRWMutex.Lock()
	// 通过线程池连接数据库
	// 结束时归还线程池
	defer model.MySQLRWMutex.Unlock()
	return nil
}

func ModifyRule(ruleId int,status int)error{
	model.MySQLRWMutex.Lock()
	// 通过线程池连接数据库
	// 结束时归还线程池
	defer model.MySQLRWMutex.Unlock()
	return nil
}

func MatchRules(cr *model.CReport) ([]*model.CacheMessage,error){
	// 查询其他规则表中所有匹配的rule
	// 要求返回的切片中ruleid对应规则按照UpdateVersionCodeInt64降序排列

	// 通过线程池连接数据库
	// 结束时归还线程池
	return nil,nil
}