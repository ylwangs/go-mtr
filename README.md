# Go-MTR
------
## 简介
### MTR
> 程序基于:
> https://github.com/tonobo/mtr
> https://github.com/liuxinglanyue/mtr
> 进行了一定的bug修复和优化开发。支持单机多协程并发mtr探测，同时支持ipv4和ipv6。

### PING
> 支持单机多协程发起ping探测，同时支持ipv4和ipv6

## 案例
```
var targets = []string{"216.58.200.78", "52.74.223.119", "123.125.114.144"}

func main() {
    // 发起mtr操作
	for _, val := range targets {
		go func(target string) {
			for {
				mm, err := mtr.Mtr(target, 30, 10, 800)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(mm)

				time.Sleep(60 * time.Second)
			}
		}(val)
	}

    // 发起ping操作
	for _, val := range targets {
		go func(target string) {
			for {
				mm, err := ping.Ping(target, 10, 800, 10)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(mm)

				time.Sleep(60 * time.Second)
			}
		}(val)
	}

	select {}
}

```

*注：针对mtr可设定的参数包括maxHops(最大跳数)、sntSize(发送数据包数量)、timeoutMs(icmp包超时时间)。针对ping探测可设定的参数包括count(数据包数量)、timeoutMs(超时时间)、intervalMs(发包间隔)*

## 原理
> mtr发送icmp数据报，先发送TTL为1的，到第一个路由器TTL减1，并返回一个超时的ICMP报文，这样就得到了第一个路由器的地址；  接下来发送TTL值为2的报文，得到第二个路由器的报文； 直到找到目的IP地址，mtr一次探测结束。然后根据发送数据包的数量，重新发起mtr探测，记录结果。

