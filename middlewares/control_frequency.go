package middlewares

import (
	"iptables_manage/conn_tool"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 控制频率
func ControlFrequency() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip1 := ctx.ClientIP()
		ip2 := ctx.Request.Header.Get("X-Forward-For")
		keyName := ip1 + ip2
		_, err := conn_tool.RDB.Get(keyName).Result()
		if err != nil {
			conn_tool.RDB.Set(keyName, ip1+ip2, time.Second*1)
			ctx.Next()
		} else {
			ctx.Abort()
			remainingTime, _ := conn_tool.RDB.TTL(keyName).Result()
			ctx.JSON(http.StatusOK, gin.H{
				"code":             http.StatusServiceUnavailable,
				"message":          "频率太高",
				"remaining second": remainingTime / 1000000000,
			})
			return
		}
	}

}
