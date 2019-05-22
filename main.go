package main

import (
	"fmt"
	"time"

	"github.com/ylwang1122/go-mtr/mtr"
)

var targets = []string{"216.58.200.78", "52.74.223.119", "123.125.114.144"}

func main() {
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

	select {}
}
