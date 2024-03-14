package redis_tool

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup
var RedisConn *redis.Client

func GetRedisConn(Addr, Password string, DB int) *redis.Client {
	//连接服务器
	return redis.NewClient(&redis.Options{
		Addr:     Addr,     // use default Addr
		Password: Password, // no password set
		DB:       DB,       // use default DB
	})
}

func testRedisBase() {
	defer wg.Done()

	//心跳
	pong, err := RedisConn.Ping().Result()
	log.Println(pong, err) // Output: PONG <nil>

	ExampleClient_String()
	ExampleClient_List()
	ExampleClient_Hash()
	ExampleClient_Set()
	ExampleClient_SortSet()
	ExampleClient_HyperLogLog()
	ExampleClient_CMD()
	ExampleClient_Scan()
	ExampleClient_Tx()
	ExampleClient_Script()
	ExampleClient_PubSub()
}

func ExampleClient_String() {
	log.Println("ExampleClient_String")
	defer log.Println("ExampleClient_String")

	//kv读写
	err := RedisConn.Set("key", "value", 1*time.Second).Err()
	log.Println(err)

	//获取过期时间
	tm, err := RedisConn.TTL("key").Result()
	log.Println(tm)

	val, err := RedisConn.Get("key").Result()
	log.Println(val, err)

	val2, err := RedisConn.Get("missing_key").Result()
	if err == redis.Nil {
		log.Println("missing_key does not exist")
	} else if err != nil {
		log.Println("missing_key", val2, err)
	}

	//不存在才设置 过期时间 nx ex
	value, err := RedisConn.SetNX("counter", 0, 1*time.Second).Result()
	log.Println("setnx", value, err)

	//Incr
	result, err := RedisConn.Incr("counter").Result()
	log.Println("Incr", result, err)
}

func ExampleClient_List() {
	log.Println("ExampleClient_List")
	defer log.Println("ExampleClient_List")

	//添加
	log.Println(RedisConn.RPush("list_test", "message1").Err())
	log.Println(RedisConn.RPush("list_test", "message2").Err())

	//设置
	log.Println(RedisConn.LSet("list_test", 2, "message set").Err())

	//remove
	ret, err := RedisConn.LRem("list_test", 3, "message1").Result()
	log.Println(ret, err)

	rLen, err := RedisConn.LLen("list_test").Result()
	log.Println(rLen, err)

	//遍历
	lists, err := RedisConn.LRange("list_test", 0, rLen-1).Result()
	log.Println("LRange", lists, err)

	//pop没有时阻塞
	result, err := RedisConn.BLPop(1*time.Second, "list_test").Result()
	log.Println("result:", result, err, len(result))
}

func ExampleClient_Hash() {
	log.Println("ExampleClient_Hash")
	defer log.Println("ExampleClient_Hash")

	datas := map[string]interface{}{
		"name": "LI LEI",
		"sex":  1,
		"age":  28,
		"tel":  123445578,
	}

	//添加
	if err := RedisConn.HMSet("hash_test", datas).Err(); err != nil {
		log.Fatal(err)
	}

	//获取
	rets, err := RedisConn.HMGet("hash_test", "name", "sex").Result()
	log.Println("rets:", rets, err)

	//成员
	retAll, err := RedisConn.HGetAll("hash_test").Result()
	log.Println("retAll", retAll, err)

	//存在
	bExist, err := RedisConn.HExists("hash_test", "tel").Result()
	log.Println(bExist, err)

	bRet, err := RedisConn.HSetNX("hash_test", "id", 100).Result()
	log.Println(bRet, err)

	//删除
	log.Println(RedisConn.HDel("hash_test", "age").Result())
}

func ExampleClient_Set() {
	log.Println("ExampleClient_Set")
	defer log.Println("ExampleClient_Set")

	//添加
	ret, err := RedisConn.SAdd("set_test", "11", "22", "33", "44").Result()
	log.Println(ret, err)

	//数量
	count, err := RedisConn.SCard("set_test").Result()
	log.Println(count, err)

	//删除
	ret, err = RedisConn.SRem("set_test", "11", "22").Result()
	log.Println(ret, err)

	//成员
	members, err := RedisConn.SMembers("set_test").Result()
	log.Println(members, err)

	bret, err := RedisConn.SIsMember("set_test", "33").Result()
	log.Println(bret, err)

	RedisConn.SAdd("set_a", "11", "22", "33", "44")
	RedisConn.SAdd("set_b", "11", "22", "33", "55", "66", "77")
	//差集
	diff, err := RedisConn.SDiff("set_a", "set_b").Result()
	log.Println(diff, err)

	//交集
	inter, err := RedisConn.SInter("set_a", "set_b").Result()
	log.Println(inter, err)

	//并集
	union, err := RedisConn.SUnion("set_a", "set_b").Result()
	log.Println(union, err)

	ret, err = RedisConn.SDiffStore("set_diff", "set_a", "set_b").Result()
	log.Println(ret, err)

	rets, err := RedisConn.SMembers("set_diff").Result()
	log.Println(rets, err)
}

func ExampleClient_SortSet() {
	log.Println("ExampleClient_SortSet")
	defer log.Println("ExampleClient_SortSet")

	addArgs := make([]redis.Z, 100)
	for i := 1; i < 100; i++ {
		addArgs = append(addArgs, redis.Z{Score: float64(i), Member: fmt.Sprintf("a_%d", i)})
	}
	//log.Println(addArgs)

	Shuffle := func(slice []redis.Z) {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		for len(slice) > 0 {
			n := len(slice)
			randIndex := r.Intn(n)
			slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
			slice = slice[:n-1]
		}
	}

	//随机打乱
	Shuffle(addArgs)

	//添加
	ret, err := RedisConn.ZAddNX("sortset_test", addArgs...).Result()
	log.Println(ret, err)

	//获取指定成员score
	score, err := RedisConn.ZScore("sortset_test", "a_10").Result()
	log.Println(score, err)

	//获取制定成员的索引
	index, err := RedisConn.ZRank("sortset_test", "a_50").Result()
	log.Println(index, err)

	count, err := RedisConn.SCard("sortset_test").Result()
	log.Println(count, err)

	//返回有序集合指定区间内的成员
	rets, err := RedisConn.ZRange("sortset_test", 10, 20).Result()
	log.Println(rets, err)

	//返回有序集合指定区间内的成员分数从高到低
	rets, err = RedisConn.ZRevRange("sortset_test", 10, 20).Result()
	log.Println(rets, err)

	//指定分数区间的成员列表
	rets, err = RedisConn.ZRangeByScore("sortset_test", redis.ZRangeBy{Min: "(30", Max: "(50", Offset: 1, Count: 10}).Result()
	log.Println(rets, err)
}

//用来做基数统计的算法，HyperLogLog 的优点是，在输入元素的数量或者体积非常非常大时，计算基数所需的空间总是固定 的、并且是很小的。
//每个 HyperLogLog 键只需要花费 12 KB 内存，就可以计算接近 2^64 个不同元素的基 数
func ExampleClient_HyperLogLog() {
	log.Println("ExampleClient_HyperLogLog")
	defer log.Println("ExampleClient_HyperLogLog")

	for i := 0; i < 10000; i++ {
		RedisConn.PFAdd("pf_test_1", fmt.Sprintf("pfkey%d", i))
	}
	ret, err := RedisConn.PFCount("pf_test_1").Result()
	log.Println(ret, err)

	for i := 0; i < 10000; i++ {
		RedisConn.PFAdd("pf_test_2", fmt.Sprintf("pfkey%d", i))
	}
	ret, err = RedisConn.PFCount("pf_test_2").Result()
	log.Println(ret, err)

	RedisConn.PFMerge("pf_test", "pf_test_2", "pf_test_1")
	ret, err = RedisConn.PFCount("pf_test").Result()
	log.Println(ret, err)
}

func ExampleClient_PubSub() {
	log.Println("ExampleClient_PubSub")
	defer log.Println("ExampleClient_PubSub")
	//发布订阅
	pubsub := RedisConn.Subscribe("subkey")
	_, err := pubsub.Receive()
	if err != nil {
		log.Fatal("pubsub.Receive")
	}
	ch := pubsub.Channel()
	time.AfterFunc(1*time.Second, func() {
		log.Println("Publish")

		err = RedisConn.Publish("subkey", "test publish 1").Err()
		if err != nil {
			log.Fatal("RedisConn.Publish", err)
		}

		RedisConn.Publish("subkey", "test publish 2")
	})
	for msg := range ch {
		log.Println("recv channel:", msg.Channel, msg.Pattern, msg.Payload)
	}
}

func ExampleClient_CMD() {
	log.Println("ExampleClient_CMD")
	defer log.Println("ExampleClient_CMD")

	//执行自定义redis命令
	Get := func(rdb *redis.Client, key string) *redis.StringCmd {
		cmd := redis.NewStringCmd("get", key)
		RedisConn.Process(cmd)
		return cmd
	}

	v, err := Get(RedisConn, "NewStringCmd").Result()
	log.Println("NewStringCmd", v, err)

	v, err = RedisConn.Do("get", "RedisConn.do").String()
	log.Println("RedisConn.Do", v, err)
}

func ExampleClient_Scan() {
	log.Println("ExampleClient_Scan")
	defer log.Println("ExampleClient_Scan")

	//scan
	for i := 1; i < 1000; i++ {
		RedisConn.Set(fmt.Sprintf("skey_%d", i), i, 0)
	}

	cusor := uint64(0)
	for {
		keys, retCusor, err := RedisConn.Scan(cusor, "skey_*", int64(100)).Result()
		log.Println(keys, cusor, err)
		cusor = retCusor
		if cusor == 0 {
			break
		}
	}
}

func ExampleClient_Tx() {
	pipe := RedisConn.TxPipeline()
	incr := pipe.Incr("tx_pipeline_counter")
	pipe.Expire("tx_pipeline_counter", time.Hour)

	// Execute
	//
	//     MULTI
	//     INCR pipeline_counter
	//     EXPIRE pipeline_counts 3600
	//     EXEC
	//
	// using one rdb-server roundtrip.
	_, err := pipe.Exec()
	fmt.Println(incr.Val(), err)
}

func ExampleClient_Script() {
	IncrByXX := redis.NewScript(`
        if redis.call("GET", KEYS[1]) ~= false then
            return redis.call("INCRBY", KEYS[1], ARGV[1])
        end
        return false
    `)

	n, err := IncrByXX.Run(RedisConn, []string{"xx_counter"}, 2).Result()
	fmt.Println(n, err)

	err = RedisConn.Set("xx_counter", "40", 0).Err()
	if err != nil {
		panic(err)
	}

	n, err = IncrByXX.Run(RedisConn, []string{"xx_counter"}, 2).Result()
	fmt.Println(n, err)
}
