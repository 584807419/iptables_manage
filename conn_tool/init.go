package conn_tool

import (
	"github.com/go-redis/redis"
	"iptables_manage/config"
	"iptables_manage/redis_tool"
	"runtime"
	"time"
)

// var DB = Init() // 常量
var RDB = InitRedisDB()
var SysType = runtime.GOOS

var RedisCacheExpire = 30 * time.Second

//func Init() *sql.DB {
//	// 连接数据库
//	//DB, err := database_tool.GetDb("127.0.0.1", 3306, "root", "root")
//	DB, err := database_tool.GetDb("192.168.199.17", 49153, "root", "root")
//	if err != nil {
//		log.Println("gorm Init error:", err)
//	}
//	//defer DB.Close() //延迟关闭数据库控制器,释放数据库连接
//	return DB
//}

func InitRedisDB() *redis.Client {
	//redisConn := redis_tool.GetRedisConn("127.0.0.1:6379", "jTq7AApSED8MkH5q", 15)
	redisConn := redis_tool.GetRedisConn(config.RedisAddr, "", 0)
	//redisConn := redis_tool.GetRedisConn("127.0.0.1:6379", "", 0)
	return redisConn
}
