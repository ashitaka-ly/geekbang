package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("#######################")
	go rLock()
	go wLock()
	go lock()
	time.Sleep(3 * time.Second)
	fmt.Println("#######################")
}

func rLock() {
	lock := sync.RWMutex{}
	for i := 0; i < 5; i++ {
		// 读锁不互斥，可以多次加锁
		lock.RLock()
		// 解锁只会在函数退出时执行一次
		defer lock.RUnlock()
		fmt.Println("rLock ->", i)
	}
}

func wLock() {
	lock := sync.RWMutex{}
	for i := 0; i < 5; i++ {
		// 错误写法，写锁是互斥锁，在第二次执行时会直接卡住
		/*
			lock.Lock()
			defer lock.Unlock()
		*/
		func() {
			lock.Lock()
			defer lock.Unlock()
			fmt.Println("wLock ->", i)
		}()
	}
}

func lock() {
	lock := sync.Mutex{}
	for i := 0; i < 5; i++ {
		// 错误写法，写锁是互斥锁，在第二次执行时会直接卡住
		/*
			lock.Lock()
			defer lock.Unlock()
		*/
		func() {
			lock.Lock()
			defer lock.Unlock()
			fmt.Println("lock ->", i)
		}()
	}
}
