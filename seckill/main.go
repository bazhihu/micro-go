package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"micro-go/seckill/local"
	"micro-go/seckill/remote"
	"micro-go/seckill/util"
	"net/http"
	"os"
	"strings"
)

/*
	秒杀服务
*/

var (
	localSpike  local.LocalSpike
	remoteSpike remote.RemoteSpikeKeys
	redisPool   *redis.Pool
	done        chan int
)

// 初始化要使用的结构体和redis 连接池
func init() {
	localSpike = local.LocalSpike{
		LocalInStock:     150,
		LocalSalesVolume: 0,
	}
	remoteSpike = remote.RemoteSpikeKeys{
		SpikeOrderHashKey:  "ticket_hash_key",
		TotalInventoryKey:  "ticket_total_nums",
		QuantityOfOrderKey: "ticket_sold_nums",
	}
	redisPool = remote.NewPool()
	done = make(chan int, 1)
	done <- 1
}

func main() {
	http.HandleFunc("/buy/ticket", handleReq)
	http.ListenAndServe(":3001", nil)
}

func handleReq(w http.ResponseWriter, r *http.Request) {
	var LogMsg string
	redisConn := redisPool.Get()

	<-done
	if localSpike.LocalDeductionStock() && remoteSpike.RemoteDeductionStock(redisConn) {
		util.RespJson(w, 1, "抢票成功", nil)
		LogMsg = fmt.Sprintf("result:1,localSales:%d/n", localSpike.LocalSalesVolume)
	} else {
		util.RespJson(w, -1, "已售馨", nil)
		LogMsg = fmt.Sprintf("result:0,localSales:%d/n", localSpike.LocalSalesVolume)
	}

	// 将抢票状态写入到log中
	done <- 1
	writeLog(LogMsg, "./stat.log")
}

func writeLog(msg string, logPath string) {
	fd, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer fd.Close()
	content := strings.Join([]string{msg, "\r\n"}, "3001")
	buf := []byte(content)
	fd.Write(buf)
}
