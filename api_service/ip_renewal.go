package api_service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"iptables_manage/conn_tool"
	"log"
	"net/http"
	"os/exec"
	"time"
)

type IpInfo struct {
	IpAddr  string `json:"ip_addr"`
	Message string `json:"message"`
}

var callCount int

func ipExec(ipAddressCmd string) {
	//需要执行的命令： free -mh
	cmd := exec.Command("/bin/bash", "-c", ipAddressCmd)

	// 获取管道输入
	output, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("无法获取命令的标准输出管道", err.Error())
		return
	}

	// 执行Linux命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Linux命令执行失败，请检查命令输入是否有误", err.Error())
		return
	}

	// 读取所有输出
	bytes, err := io.ReadAll(output)
	if err != nil {
		fmt.Println("打印异常，请检查")
		return
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("Wait", err.Error())
		return
	}

	fmt.Printf("打印内存信息：\n\n%s", bytes)
}

func ipMonitor() {
	for {
		log.Println("开始循环检查，间隔5秒")
		time.Sleep(5 * time.Second)
		valueSlice, err := conn_tool.RDB.SMembers("hash_ip_flag").Result()
		if err != nil {
			log.Println("SMembers获取hash_ip_flag这个set中现存的IP出错%v", err)
		} else {
			for index, value := range valueSlice {
				log.Println("查看TTL-第%v个IP：%v", index, value)
				ttlValue, err := conn_tool.RDB.TTL(value).Result()

				if err != nil {
					panic(err)
				}

				if ttlValue == time.Duration(-1)*time.Second {
					log.Println("键存在但是没有设置生存时间")
				} else if ttlValue == time.Duration(-2)*time.Second {
					log.Println("IP键不存在")
					ipExec(fmt.Sprintf("iptables -D INPUT -s %v -j ACCEPT", value))
					log.Printf("从iptables中删除 IP：%v 规则成功\n", value)
					result, err := conn_tool.RDB.SRem("hash_ip_flag", value).Result()
					if err != nil {
						panic(err)
					}
					log.Printf("成功从集合中移除了 %d 个成员\n", result)
					log.Printf("从redis set中srem IP：%v 成功\n", value)
				} else {
					log.Printf("键的剩余生存时间：%v\n", ttlValue)
				}
			}
		}
	}
}

func IpRuleMonitor(c *gin.Context) {
	// 监控redis，过期了就删除规则
	var newIpInfo IpInfo
	if err := c.BindJSON(&newIpInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		callCount++
		log.Printf("调用次数：callCount: %v \n", callCount)
		if callCount > 1 {
			c.JSON(http.StatusConflict, gin.H{"message": "方法已经调用过了开始运行了"})
			return
		} else {
			log.Println("开始执行ipMonitor")
			go ipMonitor()
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "开始监控清空规则"})

}

func IpRuleCreate(c *gin.Context) {
	var newIpInfo IpInfo
	if err := c.BindJSON(&newIpInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		// 收到请求后先去执行一下添加规则,添加一个set中做记录
		ipExec(fmt.Sprintf("iptables -A INPUT -s %v -j ACCEPT", newIpInfo.IpAddr))
		log.Println("添加iptables规则执行完成")
		//result, err := conn_tool.RDB.SetNX(newIpInfo.IpAddr, 1, conn_tool.RedisCacheExpire).Result()
		result, err := conn_tool.RDB.SetNX(newIpInfo.IpAddr, 1, conn_tool.RedisCacheExpire).Result()
		if err != nil {
			panic(err)
		}
		if result {
			fmt.Printf("redis设置IP %v 为key成功\n", newIpInfo.IpAddr)
		} else {
			fmt.Printf("redis设置IP %v 为key失败\n", newIpInfo.IpAddr)
		}
		resultCount, err := conn_tool.RDB.SAdd("hash_ip_flag", newIpInfo.IpAddr).Result()
		if err != nil {
			panic(err)
		}
		if resultCount == 1 {
			fmt.Printf("redis设置IP %v SAdd hash_ip_flag成功\n", newIpInfo.IpAddr)
		} else {
			fmt.Printf("redis设置IP %v SAdd hash_ip_flag失败\n", newIpInfo.IpAddr)
		}
		log.Println("redis添加IP数据执行完成")
		c.JSON(http.StatusOK, gin.H{"message": "添加iptables规则和redis添加IP数据执行完成"})
	}
}

func IpRuleRenewal(c *gin.Context) {
	var newIpInfo IpInfo
	if err := c.BindJSON(&newIpInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		value, err := conn_tool.RDB.Expire(newIpInfo.IpAddr, conn_tool.RedisCacheExpire).Result()
		log.Println("redis给IP续期结果", value, err)
		c.JSON(http.StatusOK, gin.H{"message": "redis给IP续期完成"})
		//var resIns = new(ResponseStruct)
		//adCode, _ := c.GetQuery("ad_code")
		////单行数据查询
		//if adCode != "" {
		//
		//	// 去redis获取缓存
		//	jsonStr, _ := conn_tool.RDB.Get(adCode).Result()
		//	if jsonStr != "" { // 如果有缓存才去json反序列化 把读取出来的二进制字节流转化为结构体
		//		err := json.Unmarshal([]byte(jsonStr), &resIns)
		//		if err != nil {
		//			c.JSON(502, gin.H{})
		//		} else {
		//			c.JSON(200, gin.H{"ad_code": resIns.AdCode, "name": resIns.Name, "path": resIns.Path, "json_data": resIns.JsonData})
		//		}
		//	} else {
		//		codelSql := fmt.Sprintf("select path,name, json_data from ad_data_map.new_map_data where code = '%v'", adCode)
		//		showMasterStatus := database_tool.QueryAndParse(conn_tool.DB, codelSql)
		//		name := showMasterStatus["name"]
		//		path := showMasterStatus["path"]
		//		jsonData := showMasterStatus["json_data"]
		//		//dadadasdsa := database_tool.Data2Json(showMasterStatus)
		//		//print(dadadasdsa)
		//		c.JSON(200, gin.H{"ad_code": adCode, "name": name, "path": path, "json_data": jsonData})
		//
		//		// 放入redis缓存
		//		resIns.AdCode = adCode
		//		resIns.Name = name
		//		resIns.Path = path
		//		resIns.JsonData = jsonData
		//		byteSlice, err := json.Marshal(resIns)
		//		value, err := conn_tool.RDB.SetNX(adCode, byteSlice, conn_tool.RedisCacheExpire).Result()
		//		log.Println("setnx", value, err)
		//	}
		//}
		//adName, _ := c.GetQuery("ad_name")
		//if adName != "" {
		//
		//	// 去redis获取缓存
		//	jsonStr, _ := conn_tool.RDB.Get(adName).Result()
		//	if jsonStr != "" { // 如果有缓存才去json反序列化 把读取出来的二进制字节流转化为结构体
		//		err := json.Unmarshal([]byte(jsonStr), &resIns)
		//		if err != nil {
		//			c.JSON(502, gin.H{})
		//		} else {
		//			c.JSON(200, gin.H{"ad_code": resIns.AdCode, "name": resIns.Name, "path": resIns.Path, "json_data": resIns.JsonData})
		//		}
		//	} else {
		//		codelSql := fmt.Sprintf("select code,path,name, json_data from ad_data_map.new_map_data where name = '%v'", adName)
		//		showMasterStatus := database_tool.QueryAndParse(conn_tool.DB, codelSql)
		//		code := showMasterStatus["code"]
		//		name := showMasterStatus["name"]
		//		path := showMasterStatus["path"]
		//		jsonData := showMasterStatus["json_data"]
		//		c.JSON(200, gin.H{"ad_code": code, "name": name, "path": path, "json_data": jsonData})
		//
		//		// 放入redis缓存
		//		resIns.AdCode = code
		//		resIns.Name = name
		//		resIns.Path = path
		//		resIns.JsonData = jsonData
		//		byteSlice, err := json.Marshal(resIns)
		//		value, err := conn_tool.RDB.SetNX(adName, byteSlice, conn_tool.RedisCacheExpire).Result()
		//		log.Println("setnx", value, err)
		//	}
		//
		//}
		//if adName == "" && adCode == "" {
		//	c.JSON(404, gin.H{
		//		"message": "参数缺失",
		//	})
		//}
	}

}
