package main

import (
	"github.com/gin-gonic/gin"
	"iptables_manage/api_service"
)

type ResponseStruct struct {
	AdCode   string `json:"ad_code"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	JsonData string `json:"jsonData"`
}

func main() {
	r := gin.Default()
	//全局中间件
	//r.Use(middlewares.TraceReqByUuid()) // 请求标记
	//r.Use(middlewares.Cors())           // 跨域中间件
	//r.Use(middlewares.StatCost())
	//r.Use(middlewares.ControlFrequency()) // 频率控制中间件

	r.POST("/ip_create", api_service.IpRuleCreate)
	r.POST("/ip_renewal", api_service.IpRuleRenewal)
	r.POST("/ip_monitor", api_service.IpRuleMonitor)
	//r.GET("/ad_code_search", api_service.AdCodeSearch) // http://127.0.0.1:8090/ad_code_search?ad_sub_code=410100

	r.Run(":8090") // 监听并在 0.0.0.0:8080 上启动服务
}
