package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"iptables_manage/helper"
	"log"
	"net/http"
	"time"
)

// 添加uuid借助Context传过去方便后续打印日志追踪请求整个生命周期
func TraceReqByUuid() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//ip1 := ctx.ClientIP()
		//ip2 := ctx.Request.Header.Get("X-Forward-For")
		//keyName := ip1 + ip2
		//fmt.Printf("请求来源%v ", keyName)

		req_uuid := helper.GetUUID()
		ctx.Set("req_uuid", req_uuid)
		fmt.Printf("uuid%v ", req_uuid)

		//method := ctx.Request.Method
		//fmt.Printf("请求方法%v ", method)
		ctx.Next()
	}

}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func StatCost() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		//可以设置一些公共参数
		c.Set("example", "12345")
		//等其他中间件先执行
		c.Next()
		//获取耗时
		latency := time.Since(t)
		log.Printf("请求耗时:%d us", latency/1000)
	}
}

// 拦截器
func MyMiddleware(c *gin.Context) {
	//请求前逻辑
	c.Next()
	//请求后逻辑
}
