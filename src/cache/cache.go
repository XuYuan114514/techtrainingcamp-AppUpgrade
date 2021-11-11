package cache

import (
	"GaryReleaseProject/src/database"
	"GaryReleaseProject/src/model"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

// lua脚本教程：https://studygolang.com/articles/19741
// https://blog.csdn.net/qq_44910471/article/details/89295668

const(
	MATCH_DEVICEID_HASH = `
local ruleExist = redis.call('exists', KEYS[1])
if ruleExist == 0
then
    return 3
else
	if redis.call('sismember',KEYS[1],ARGV[1]) == 1
	then
		return 1
	else
		return 2
	end
end
`
)

const(
	SAFT_DEL_MUTEX = `
if redis.call("get",KEYS[1]) == ARGV[1] then
    return redis.call("del",KEYS[1])
else
    return 0
end
`
)

const(
	CACHE_AND_SEARCH = `

	local str = ARGV[1]
	local delimiter = KEYS[3]

    local dLen = string.len(delimiter)
    local newDeli = ''
    for i=1,dLen,1 do
        newDeli = newDeli .. "["..string.sub(delimiter,i,i).."]"
    end
    local locaStart,locaEnd = string.find(str,newDeli)
    --local arr = {}
    local n = 1
    while locaStart ~= nil
    do
        if locaStart>0 then
            --arr[n] = string.sub(str,1,locaStart-1)
			redis.call("sadd",KEYS[1],string.sub(str,1,locaStart-1))
            n = n + 1
        end

        str = string.sub(str,locaEnd+1,string.len(str))
        locaStart,locaEnd = string.find(str,newDeli)
    end
    if str ~= nil then
		redis.call("sadd",KEYS[1],str)
        --arr[n] = str
    end
	
	if redis.call("sismember",KEYS[1],KEYS[2]) == 1 then 
		return 1
	else 
		return 2
	end
`
)

var(
	waitForSetnx = 2 * time.Second // 想要拉取的rule正在被别人拉取时，一次睡眠等待的时间
	setnxHoldTime = "100" // 分布锁失效时间
)

func MatchWhitelist(candidates *[]model.CacheMessage,deviceId string)(pos int,err error){
	// 拿到rule切片之后
	if len(*candidates)==0{
		return -1,nil
	}
	// 用完后将连接放回连接池
	rc := model.RedisClient.Get()
	defer rc.Close()
	// 按照优先级，逐个查缓存
	for i:=0;i < len(*candidates);i++{
		rid := strconv.Itoa((*candidates)[i].RuleId)
		/*
		1：ruleId在缓存且deviceId命中,可以返回了
		2：ruleId在缓存里且deviceId未命中，可以continue
		3：ruleId不在缓存中
		 */
		lua := redis.NewScript(1,MATCH_DEVICEID_HASH)
		in,err :=redis.Int(lua.Do(rc, "ruleId:"+rid, deviceId))
		if err!=nil{
			return -1,fmt.Errorf("redis lua发生错误")
		}
		switch in {
		case 1:
			return i,nil
		case 2:
			continue
		case 3:
			t1:=time.Tick( waitForSetnx)
			dealDone,_ := dealMissWhitelist(rid,deviceId,&rc,lua)
			for dealDone == 0 {
				select {
				case <-t1:
					dealDone,_ = dealMissWhitelist(rid,deviceId,&rc,lua)
				}
			}
			if dealDone == 1{
				return i,nil
			} else if dealDone == 2{
				continue
			}
		default:
			return -1,fmt.Errorf("redis lua返回异常结果")
		}
	}
	return -1,err
}

// 当ruleid不在缓存中，尝试从mysql拉取数据到cache中
// 假如当前ruleid被阻塞了，等待
func dealMissWhitelist(ruleId string,deviceId string,rc *redis.Conn,lua *redis.Script)(status int,err error){
	/*
		lua脚本返回in
		1：ruleId在缓存且deviceId命中,可以返回了
		2：ruleId在缓存里且deviceId未命中，可以continue
		3：ruleId不在缓存中

		本函数返回
		0：拿锁失败，继续等待
		1：不用拿锁，ruleId在缓存且deviceId命中,可以返回了 || 拿到锁，写入缓存，命中
		2：不用拿锁，ruleId在缓存里且deviceId未命中，可以continue || 拿到锁，写入缓存，未命中
	*/
	in,err :=redis.Int(lua.Do(*rc, "ruleId:"+ruleId, deviceId))
	if err!=nil{
		return -1,fmt.Errorf("redis lua发生错误")
	}
	switch {
	// 等待的时候已经有人写到cache了
	case in==1 || in ==2:
		return in,nil
	// 尝试拿锁
	case in == 3:
		return tryPullCache(ruleId,deviceId,rc)
	default:
		return 0,nil
	}

}

// 数据库longtext转切片，目前使用split，是否有更高性能的方案呢？
func tryPullCache(ruleId string,deviceId string, rcc *redis.Conn)(int,error){
	rc := *rcc
	// setnx 分布锁
	in,_ := rc.Do("set","mutex:"+ruleId,deviceId,"ex",setnxHoldTime,"nx")
	if in == 1{
		// 拿到锁了，开始写缓存
		whiteList,_ := database.GetWhitelist(ruleId)
		// lua保证写入之后立刻判断
		luaPull := redis.NewScript(3,CACHE_AND_SEARCH)
		status,_ :=redis.Int(luaPull.Do(rc,"ruleId:"+ruleId,deviceId,",",whiteList))
		// lua防止进程超时之后误删锁
		luaDel := redis.NewScript(1,SAFT_DEL_MUTEX)
		luaDel.Do(rc, "mutex:"+ruleId,deviceId)
		return status,nil
	}else{
		// 拿锁失败, 继续等待
		return 0,nil
	}
}







