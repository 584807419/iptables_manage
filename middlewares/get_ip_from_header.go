package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func GetIpFromHeader() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip1 := ctx.ClientIP()
		ip2 := ctx.Request.Header.Get("X-Forward-For")
		if ip2 != "" {
			// X-Forwarded-For可能包含多个IP地址，最左边的是离客户端最近的IP
			ip2 = strings.Split(ip2, ", ")[0]
		}
		ip3 := ctx.Request.RemoteAddr
		fmt.Printf("获取到ip1\n：%v", ip1)
		fmt.Printf("获取到ip2\n：%v", ip2)
		fmt.Printf("获取到ip3\n：%v", ip3)
		ctx.Next()
		//_, err := conn_tool.RDB.Get(keyName).Result()
		//if err != nil {
		//	conn_tool.RDB.Set(keyName, ip1+ip2, time.Second*1)
		//	ctx.Next()
		//} else {
		//	ctx.Abort()
		//	remainingTime, _ := conn_tool.RDB.TTL(keyName).Result()
		//	ctx.JSON(http.StatusOK, gin.H{
		//		"code":             http.StatusServiceUnavailable,
		//		"message":          "频率太高",
		//		"remaining second": remainingTime / 1000000000,
		//	})
		//	return
		//}
	}

}
