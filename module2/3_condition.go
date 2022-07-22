package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	q := Queue{
		queue: []string{},
		cond:  sync.NewCond(&sync.Mutex{}),
	}
	go func() {
		for {
			q.Enqueue("aaa")
			time.Sleep(time.Second * 3)
		}
	}()
	for {
		res := q.Dequeue()
		fmt.Println("get from queue -> ", res)
		time.Sleep(time.Second * 2)
	}

}

type Queue struct {
	queue []string
	cond  *sync.Cond
}

func (q *Queue) Enqueue(item string) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	q.queue = append(q.queue, item)
	// 有数据后广播通知
	fmt.Printf("putting %s to queue, notify all \n", item)
	q.cond.Broadcast()

}

func (q *Queue) Dequeue() string {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	// 如果队列是空的，则等待
	if len(q.queue) == 0 {
		fmt.Println("no data available, please wait")
		q.cond.Wait()
	}
	// 返回第一个，同时截取第一个到最后一个
	result := q.queue[0]
	q.queue = q.queue[1:]
	return result
}
