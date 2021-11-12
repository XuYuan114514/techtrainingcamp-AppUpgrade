package model

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Rule struct {
	// rule condition
	Platform             string `json:"platform"`
	DownloadUrl          string `json:"download_url"`
	UpdateVersionCode    string `json:"update_version_code"`
	MD5                  string `json:"md5"`
	DeviceIdList         string `json:"device_id_list"`
	MaxUpdateVersionCode string `json:"max_update_version_code"`
	MinUpdateVersionCode string `json:"min_update_version_code"`
	MaxOsApi             int    `json:"max_os_api"`
	MinOsApi             int    `json:"min_os_api"`
	CpuArch              int    `json:"cpu_arch"`
	Channel              string `json:"channel"`
	Title                string `json:"title"`
	UpdateTips           string `json:"update_tips"`
	//新添加的字段
	UpdateVersionCodeInt64    int64
	MaxUpdateVersionCodeInt64 int64
	MinUpdateVersionCodeInt64 int64
	//需要传入，0是正常，1是暂停，2是下线
	RuleStatus                int `json:"rule_status"`
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



var (
	// 全局变量
	MySQLRWMutex *sync.RWMutex
	DB *sql.DB
	RedisClient *redis.Pool

)
// 设定参数
var(
	// redis相关
	RedisHost string = "localhost:6379"
	RedisDb   int    = 0    // 使用的redis数据库0~15
	MaxIdle   int    =  20  // 最大空闲连接数，即会有这么多个连接提前等待着，但过了超时时间也会关闭。
	MaxActive int    = 1000 //最大连接数，即最多的tcp连接数
	IdleTimeout = 300 * time.Second  //闲置空闲链接持续时间
	// mysql相关
	ConnMaxLifetime = 300 * time.Second
	MaxIdleConns = 10 // 最大闲置
	MaxOpenConns = 1000 // 链接池最大并发数
	DataSourceName = "root:SZK_6848439@(127.0.0.1:3306)/garyrelease"
)

func InitAll(){
	MySQLRWMutex = new(sync.RWMutex)
	DB = InitDatabase()
	RedisClient,_ = InitCache()
}

func InitDatabase() *sql.DB {
	// 使用mysql连接池实现
	//"用户名：密码@[连接方式](主机名：端口号）/数据库名”
	db, err := sql.Open("mysql", DataSourceName)
	if err != nil {
		log.Fatalln("open db fail", err)
		return nil
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln("ping db failed", err)
		return nil
	}
	//设置连接细节,可适当调大
	db.SetConnMaxLifetime(ConnMaxLifetime)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)
	CreateTable(db)
	return db
}

func CloseDatabase(db *sql.DB) {
	_ = db.Close()
	fmt.Println("database closed")
}

func CreateTable(db *sql.DB) {

	//创建其他规则表config
	sql := "CREATE TABLE IF NOT EXISTS config(" +
		"platform					char(64)	NOT	NULL," +
		"download_url				char(255)	NOT NULL," +
		"update_version_code		char(64)	NOT NULL," +
		"md5 						char(255)	NOT NULL," +
		"rule_id		 			int		 	NOT NULL	AUTO_INCREMENT," + //字段更改，用rule_id查表white_lists得到白名单
		"max_update_version_code	char(64)	NOT NULL," +
		"min_update_version_code	char(64)	NOT NULL," +
		"max_os_api					int			DEFAULT NULL," +
		"min_os_api					int			DEFAULT NULL," +
		"cpu_arch					int			NOT NULL," +
		"channel					char(64)	NOT NULL," +
		"title						char(255)	NOT NULL," +
		"update_tips				char(255)	NOT NULL," +
		//自此为新加字段
		"update_version_code_int64		bigint			NOT NULL," +
		"max_update_version_code_int64	bigint			NOT NULL," +
		"min_update_version_code_int64	bigint			NOT NULL," +
		"rule_status					int			NOT NULL," +
		"PRIMARY KEY(rule_id)" +
		")ENGINE=InnoDB"

	result, _ := db.Exec(sql)
	if result != nil {
		fmt.Println("create table config succeed")
	} else {
		fmt.Println("create table config failed")
	}

	//创建白名单表white_lists
	sql2 := "CREATE TABLE If Not Exists white_lists(" +
		"rule_id			int				NOT NULL	AUTO_INCREMENT," +
		"device_id_list		longtext		NOT NULL," + //rule_id对应的白名单，注意名称
		"PRIMARY KEY(rule_id)" +
		")ENGINE=InnoDB"
	res, _ := db.Exec(sql2)
	if res != nil {
		fmt.Println("create table white_lists succeed")
	} else {
		fmt.Println("create table white_lists failed")
	}
}

func InitCache()(*redis.Pool, error){
	rc := &redis.Pool{
		MaxIdle:     MaxIdle,
		MaxActive:   MaxActive,
		IdleTimeout: IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", RedisHost)
			if err != nil {
				return nil, err
			}
			// 选择db
			c.Do("SELECT", RedisDb)
			return c, nil
		},
	}
	//缓存预热
	return rc,nil
}

func VersionToInt64(version string) int64{
	if version ==""{
		return -1
	}
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

func Str2Int(str string) int {
	if str ==""{
		return -1
	}
	res, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return res
}