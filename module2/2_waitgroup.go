package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("#######################")
	//waitByChannel()
	waitByWaitGroup()
	fmt.Println("#######################")
}

func waitByChannel() {
	ch := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("in waitByChannel ->", i)
			ch <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		// 通过 channel 的方式，这里如果没有数值会一直阻塞等待，100次
		<-ch
	}

}

// 使用 wait group 可以让线程等待其他线程结束，其用法类似与 java 的 countdownlatch
func waitByWaitGroup() {
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("in waitByWaitGroup ->", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

/* 注意一种写法的错误，这种写法的打印结果可能会是一堆10，原因是因为在执行 go 时，没有把
i 作为参数传进去，外层的主线程跑完了，里面的子线程才开始，那么这个时候子线程打印的就是外层主线程的 10
*/
func xxx() {
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println("in waitByWaitGroup ->", i)
			wg.Done()
		}()
	}
	wg.Wait()
}
