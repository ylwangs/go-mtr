package mtr

import (
	"bytes"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ylwang1122/go-mtr/common"
	"github.com/ylwang1122/go-mtr/icmp"
)

// 执行traceroute操作 新增ipv6操作
func Mtr(ipAddr string, maxHops, sntSize, retries int) (result string, err error) {
	options := common.MtrOptions{}
	options.SetMaxHops(maxHops)
	options.SetSntSize(sntSize)
	options.SetRetries(retries)

	var out common.MtrResult
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Start: %v, DestAddr: %v\n", time.Now().Format("2006-01-02 15:04:05"), ipAddr))
	out, err = runMtr(ipAddr, &options)

	if err == nil {
		if len(out.Hops) == 0 {
			buffer.WriteString("TestMtr failed. Expected at least one hop\n")
			return buffer.String(), nil
		}
	} else {
		buffer.WriteString(fmt.Sprintf("TestMtr failed due to an error: %v\n", err))
		return buffer.String(), err
	}

	buffer.WriteString(fmt.Sprintf("%-3v %-48v  %10v%c  %10v  %10v  %10v  %10v  %10v\n", "", "HOST", "Loss", '%', "Snt", "Last", "Avg", "Best", "Wrst"))

	lastTTL := 1
	for _, hop := range out.Hops {
		for j := (lastTTL + 1); j < hop.TTL; j++ {
			buffer.WriteString(fmt.Sprintf("%-3d %-48v  %10.1f%c  %10v  %10.2f  %10.2f  %10.2f  %10.2f\n", j, "???", float32(100), '%', int(0), float32(0), float32(0), float32(0), float32(0)))
		}
		lastTTL = hop.TTL
		if hop.Success {
			buffer.WriteString(fmt.Sprintf("%-3d %-48v  %10.1f%c  %10v  %10.2f  %10.2f  %10.2f  %10.2f\n", hop.TTL, hop.Address, hop.Loss, '%', hop.Snt, time2Float(hop.LastTime), time2Float(hop.AvgTime), time2Float(hop.BestTime), time2Float(hop.WrstTime)))
		} else {
			buffer.WriteString(fmt.Sprintf("%-3d %-48v  %10.1f%c  %10v  %10.2f  %10.2f  %10.2f  %10.2f\n", hop.TTL, "???", float32(100), '%', int(0), float32(0), float32(0), float32(0), float32(0)))
		}
	}

	return buffer.String(), nil
}

// mtr的实现
func runMtr(destAddr string, options *common.MtrOptions) (result common.MtrResult, err error) {
	result.Hops = []common.IcmpHop{}
	result.DestAddress = destAddr
	pid := goid()
	timeout := time.Duration(options.TimeoutMs()) * time.Millisecond

	mtrResults := make([]*common.MtrReturn, options.MaxHops()+1)

	// 发起一次mtr操作
	retry := 0
	for snt := 0; snt < options.SntSize(); snt++ {
		for ttl := 1; ttl < options.MaxHops(); ttl++ {
			time.Sleep(time.Nanosecond)

			if mtrResults[ttl] == nil {
				mtrResults[ttl] = &common.MtrReturn{TTL: ttl, Host: "???", SuccSum: 0, Success: false, LastTime: time.Duration(0), AllTime: time.Duration(0), BestTime: time.Duration(0), WrstTime: time.Duration(0), AvgTime: time.Duration(0)}
			}

			hopReturn, err := icmp.Icmp(destAddr, "", ttl, pid, timeout, 1)
			if err != nil || !hopReturn.Success {
				retry++
				if retry >= options.Retries() {
					break
				}
				continue
			}

			retry = 0
			mtrResults[ttl].SuccSum = mtrResults[ttl].SuccSum + 1
			mtrResults[ttl].Host = hopReturn.Addr
			mtrResults[ttl].LastTime = hopReturn.Elapsed
			if mtrResults[ttl].WrstTime == time.Duration(0) || hopReturn.Elapsed > mtrResults[ttl].WrstTime {
				mtrResults[ttl].WrstTime = hopReturn.Elapsed
			}
			if mtrResults[ttl].BestTime == time.Duration(0) || hopReturn.Elapsed < mtrResults[ttl].BestTime {
				mtrResults[ttl].BestTime = hopReturn.Elapsed
			}
			mtrResults[ttl].AllTime += hopReturn.Elapsed
			mtrResults[ttl].AvgTime = time.Duration((int64)(mtrResults[ttl].AllTime/time.Microsecond)/(int64)(mtrResults[ttl].SuccSum)) * time.Microsecond
			mtrResults[ttl].Seq = 1
			mtrResults[ttl].Success = true

			if isEqualIp(hopReturn.Addr, destAddr) {
				break
			}
		}
	}

	retry = 0
	for _, mtrResult := range mtrResults {
		if mtrResult == nil {
			retry++
			if retry >= options.Retries() {
				break
			}
			continue
		}
		retry = 0
		hop := common.IcmpHop{TTL: mtrResult.TTL, Snt: options.SntSize()}
		hop.Address = mtrResult.Host
		hop.Host = mtrResult.Host
		hop.AvgTime = mtrResult.AvgTime
		hop.BestTime = mtrResult.BestTime
		hop.LastTime = mtrResult.LastTime
		failSum := options.SntSize() - mtrResult.SuccSum
		loss := (float32)(failSum) / (float32)(options.SntSize()) * 100
		hop.Loss = float32(loss)
		hop.WrstTime = mtrResult.WrstTime
		hop.Success = true

		result.Hops = append(result.Hops, hop)

		// 主要用于避免ipv6省略的情况
		if isEqualIp(hop.Host, destAddr) {
			break
		}
	}

	return result, nil
}

func time2Float(t time.Duration) float32 {
	return (float32)(t/time.Microsecond) / float32(1000)
}

func goid() int {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic recover:panic info:%v", err)
		}
	}()

	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

func isEqualIp(ips1, ips2 string) bool {
	ip1 := net.ParseIP(ips1)
	if ip1 == nil {
		return false
	}

	ip2 := net.ParseIP(ips2)
	if ip2 == nil {
		return false
	}

	if ip1.String() != ip2.String() {
		return false
	}

	return true
}
