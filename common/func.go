package common

import (
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// 使用本地的localdns对域名进行解析获取其ip
func DestAddrs(addr string) ([]string, error) {
	addrs, err := net.LookupHost(addr)
	if err != nil {
		return nil, err
	}

	// 判定获取到的ip地址是否都是ip
	ipAddrs := make([]string, 0)
	for _, addr := range addrs {
		ipAddr, err := net.ResolveIPAddr("ip", addr)
		if err != nil {
			continue
		}

		ipAddrs = append(ipAddrs, ipAddr.IP.String())
	}

	return ipAddrs, nil
}

// 获取协程id
func Goid() int {
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

// 判定ip相等
func IsEqualIp(ips1, ips2 string) bool {
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

// 时间转换为float32
func Time2Float(t time.Duration) float32 {
	return (float32)(t/time.Microsecond) / float32(1000)
}
