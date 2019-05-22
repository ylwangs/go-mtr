package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/ylwang1122/go-mtr/mtr"
)

var targets = []string{"216.58.200.78", "52.74.223.119", "123.125.114.144"}

func main() {
	var wg sync.WaitGroup
	for {
		for _, val := range targets {
			wg.Add(1)
			go func(target string) {
				mm, err := mtr.Mtr(target, 30, 10, 5)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(mm)
				wg.Done()
			}(val)
		}

		wg.Wait()

		time.Sleep(60 * time.Second)
	}
}
