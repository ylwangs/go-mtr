package ping

import (
	"bytes"
	"fmt"
	"time"

	"github.com/ylwang1122/go-mtr/common"
	"github.com/ylwang1122/go-mtr/icmp"
)

// 输入参数包括 目的地址 数据包数量 超时时间 发包间隔
func Ping(addr string, count, timeout, interval int) (result string, err error) {
	pingOptions := &PingOptions{}
	pingOptions.SetCount(count)
	pingOptions.SetTimeoutMs(timeout)
	pingOptions.SetIntervalMs(interval)

	// 针对域名进行解析
	ipAddrs, err := common.DestAddrs(addr)
	if err != nil || len(ipAddrs) == 0 {
		return
	}

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Start %v, PING %v (%v)\n", time.Now().Format("2006-01-02 15:04:05"), addr, ipAddrs[0]))
	begin := time.Now().UnixNano() / 1e6
	pingReturn := runPing(ipAddrs[0], pingOptions)
	end := time.Now().UnixNano() / 1e6

	buffer.WriteString(fmt.Sprintf("%v packets transmitted, %v packet loss, time %vms\n", count, pingReturn.dropRate, end-begin))
	buffer.WriteString(fmt.Sprintf("rtt min/avg/max = %v/%v/%v ms\n", common.Time2Float(pingReturn.wrstTime), common.Time2Float(pingReturn.avgTime), common.Time2Float(pingReturn.bestTime)))

	result = buffer.String()

	return
}

func runPing(ipAddr string, option *PingOptions) (pingReturn PingReturn) {
	pingReturn = PingReturn{}
	pingReturn.destAddr = ipAddr

	pid := common.Goid()
	timeout := time.Duration(option.TimeoutMs()) * time.Millisecond
	interval := option.IntervalMs()
	ttl := DEFAULT_TTL

	pingResult := PingResult{}

	seq := 0
	for cnt := 0; cnt < option.Count(); cnt++ {
		icmpReturn, err := icmp.Icmp(ipAddr, ttl, pid, timeout, seq)
		if err != nil || !icmpReturn.Success || !common.IsEqualIp(ipAddr, icmpReturn.Addr) {
			continue
		}

		pingResult.succSum++
		if pingResult.wrstTime == time.Duration(0) || icmpReturn.Elapsed > pingResult.wrstTime {
			pingResult.wrstTime = icmpReturn.Elapsed
		}
		if pingResult.bestTime == time.Duration(0) || icmpReturn.Elapsed < pingResult.bestTime {
			pingResult.bestTime = icmpReturn.Elapsed
		}
		pingResult.allTime += icmpReturn.Elapsed
		pingResult.avgTime = time.Duration((int64)(pingResult.allTime/time.Microsecond)/(int64)(pingResult.succSum)) * time.Microsecond
		pingResult.success = true

		seq++

		time.Sleep(time.Duration(interval) * time.Millisecond)
	}

	if !pingResult.success {
		pingReturn.success = false
		pingReturn.dropRate = 100.0

		return
	}

	pingReturn.success = pingResult.success
	pingReturn.dropRate = float64(option.Count()-pingResult.succSum) / float64(option.Count())
	pingReturn.avgTime = pingResult.avgTime
	pingReturn.bestTime = pingResult.bestTime
	pingReturn.wrstTime = pingResult.wrstTime

	return
}
